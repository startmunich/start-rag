###
# Script which syncs the VectorDB database with the neo4j graph database.
# Script offers a REST Api endpoint with which a local queue can be filled with IDs.
# The queue is then processed by the script, the content of the respective IDs is collected from the neo4j database.
# Each individual content is bein preprocessed by langchain and inserted into a qdrant collection. 
# The embeddings are created with infinity.
###

from flask import Flask, request, jsonify
from neo4j import GraphDatabase
import langchain.text_splitter
from langchain_community.embeddings import InfinityEmbeddings
from langchain_community.document_loaders.csv_loader import CSVLoader
from qdrant_client import QdrantClient
from qdrant_client.http.models import PointStruct
import threading
import json
import queue
import time
import os


app = Flask(__name__)


# Initialize Neo4j driver
neo4j_uri = os.environ.get("NEO4J_URL")
neo4j_user = os.environ.get("NEO4J_USER")
neo4j_password = os.environ.get("NEO4J_PASS")
neo4j_driver = GraphDatabase.driver(neo4j_uri, auth=(neo4j_user, neo4j_password))

# Initialize Qdrant client
qdrant_uri = os.environ.get("QDRANT_URL")
qdrant_collection_name = os.environ.get("QDRANT_COLLECTION_NAME")
qdrant_client = QdrantClient(qdrant_uri, qdrant_collection_name)

# initialize infinity
infinity_api_url = os.environ.get("INFINITY_URL")
infinity_model = os.environ.get("INFINITY_MODEL")

# langchain text splitter

text_splitter = langchain.text_splitter.RecursiveCharacterTextSplitter(
                    chunk_size=1000,
                    chunk_overlap=100,
                    length_function=len)
                
headers_to_split_on = [
    ("#", "Header 1"),
    ("##", "Header 2")
]

markdown_splitter = langchain.text_splitter.MarkdownHeaderTextSplitter(headers_to_split_on=headers_to_split_on,
                                                                        strip_headers=False)

# Queue for storing IDs to be processed
id_queue = queue.Queue()


def process_queue():
    while True:
        if not id_queue.empty():
            id_to_process = id_queue.get()
            with neo4j_driver.session() as session:
                result = session.run(
                    "MATCH (n:CrawledPage {page_id: $id}) RETURN n.content AS content",
                    id=id_to_process,
                )
                app.logger.info('result of session run')
                record = result.single()
                app.logger.info(record)
                record_content = record.get("content")
                content = json.loads(record_content)

                complete_content = []

                for element in content:
                    # check if content_type is markdown or database
                    if element['content_type'] == 'markdown':
                        # apply markdown_splitter and add to complete_content
                        complete_content += markdown_splitter.split_text(element['content'])
                    elif element['content_type'] == 'database':
                        # create temp csv file and apply csv_loader and add to complete_content
                        with open('temp.csv', 'w') as file:
                            file.write(element['content'])
                        csv_loader = CSVLoader("temp.csv")
                        database_elements = []
                        for document in csv_loader.load():
                            database_elements.append(document.page_content)
                        database_content = "\n\n".join(database_elements)
                        complete_content += text_splitter.split_text(database_content)
                    else:
                        continue

                
                

                # lower the case of the chunks in chunks
                chunks = [chunk.lower() for chunk in complete_content]


                # Embed the chunks using Infinity

                embeddings = InfinityEmbeddings(model=infinity_model, 
                                                infinity_api_url=infinity_api_url
                )
                try:
                    chunks_embedded = embeddings.embed_documents(chunks)
                    print(f"embeddings of {id_queue} created successful")
                except Exception as ex:
                    print(
                        "Make sure the infinity instance is running. Verify by clicking on "
                        f"{infinity_api_url.replace('v1','docs')} Exception: {ex}. "
                    )
                
                # delete all points with the same id_to_process
                qdrant_client.delete_by_id(collection_name=qdrant_collection_name, ids=[f"{id_to_process}_*"])

                # Insert the preprocessed chunk into Qdrant
                qdrant_client.upsert(
                    collection_name=qdrant_collection_name,
                    wait=True,
                    points=[PointStruct(id=f"{id_to_process}_{count}", 
                                        vector=chunk_embedding, payload={"content": chunk}) for count, chunk_embedding, chunk in zip(enumerate(chunks_embedded), chunks)]
                )
                
        else:
            # wait 10 seconds before checking the queue again
            print("Queue is empty, waiting for 10 seconds")
            time.sleep(10)



@app.route('/enqueue', methods=['POST'])
def enqueue_ids():
    app.logger.info('json payload')
    app.logger.info(request.json)
    data = request.json
    if 'ids' in data:
        for id_to_enqueue in data['ids']:
            id_queue.put(id_to_enqueue)
            print(f"ID {id_to_enqueue} enqueued successfully")
        return jsonify({"message": "IDs enqueued successfully"}), 200
    else:
        print("No IDs provided")
        return jsonify({"error": "No IDs provided"}), 400


if __name__ == '__main__':
    # Start a separate thread to continuously process the queue
    queue_processor_thread = threading.Thread(target=process_queue)
    queue_processor_thread.daemon = True
    queue_processor_thread.start()

    # Start Flask API
    app.run(debug=True, host='0.0.0.0')
