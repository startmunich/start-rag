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
from qdrant_client.http.models import PointStruct, FilterSelector, Filter, FieldCondition, MatchValue, VectorParams
import threading
import requests
import json
import time
import os
import uuid
import numpy as np
from redis import Redis


app = Flask(__name__)


# Initialize Neo4j driver
neo4j_uri = os.environ.get("NEO4J_URL")
neo4j_user = os.environ.get("NEO4J_USER")
neo4j_password = os.environ.get("NEO4J_PASS")
neo4j_driver = GraphDatabase.driver(neo4j_uri, auth=(neo4j_user, neo4j_password))

# Initialize Qdrant client
qdrant_uri = os.environ.get("QDRANT_URL")
qdrant_collection_name = os.environ.get("QDRANT_COLLECTION_NAME")
qdrant_client = QdrantClient(url=qdrant_uri,port=6333)

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

# Redis instance for storing IDs to be processed
redis = Redis(host='redis', port=6379, db=0)

# model = SentenceTransformer(infinity_model)


def process_queue():
    app.logger.info(f"start process_queue")
    while True:
        if redis.llen('content_queue') != 0:

            task = json.loads(redis.rpop('content_queue'))

            id_to_process = task["id"]

            if task["type"] == "notion":
            
            
                
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
                            markdown_documents = markdown_splitter.split_text(element['content'])
                            markdown_strings = [document.page_content for document in markdown_documents]
                            complete_content += markdown_strings
                        elif element['content_type'] == 'database':
                            # create temp csv file and apply csv_loader and add to complete_content
                            with open('temp.csv', 'w') as file:
                                file.write(element['content'])
                            csv_loader = CSVLoader("temp.csv")
                            database_elements = [document.page_content for document in csv_loader.load()]
                            database_content = "\n\n".join(database_elements)
                            complete_content += text_splitter.split_text(database_content)
                        else:
                            continue

                    
                    

                    # lower the case of the chunks in chunks
                    chunks = [chunk.lower() for chunk in complete_content]


                    # Embed the chunks using Infinity

                    
                    chunks_embedded = [requests.post(url=f"{infinity_api_url}/embeddings", json={"model": "bge-small-en-v1.5", "input":[chunk]}).json()["data"][0]["embedding"] for chunk in chunks]
                    print("embeddings created successful")
                    
                    

                    # delete all points with the same id_to_process
                    qdrant_client.delete(collection_name=qdrant_collection_name,     
                                        points_selector=FilterSelector(
                                            filter=Filter(
                                                must=[
                                                    FieldCondition(
                                                        key="page_id",
                                                        match=MatchValue(value=id_to_process),
                                                    ),
                                                ],
                                            )
                                        ),
                                        )

                    # Insert the preprocessed chunk into Qdrant
                    
                    points_to_update = [PointStruct(id=str(uuid.uuid4()), 
                                                    vector=chunk_embedding,
                                                    payload={"content": chunk, "page_id": id_to_process}) for chunk_embedding, chunk in zip(chunks_embedded, chunks)]
                    
                    app.logger.info(f"points_to_update: {points_to_update}")

                    qdrant_client.upsert(
                        collection_name=qdrant_collection_name,
                        points= points_to_update
                        )
                
        else:
            # wait 10 seconds before checking the queue again
            app.logger.info("Queue is empty, waiting for 10 seconds")
            time.sleep(10)



@app.route('/enqueue', methods=['POST'])
def enqueue_ids():
    app.logger.info('json payload')
    app.logger.info(request.json)
    data = request.json
    if 'ids' in data:
        for id_to_enqueue in data['ids']:
            sendeable_data = json.dumps({"id" : id_to_enqueue, "type" : "notion"})
            redis.lpush("content_queue", sendeable_data)
            print(f"ID {id_to_enqueue} enqueued successfully")
        return jsonify({"message": "IDs enqueued successfully"}), 200
    else:
        print("No IDs provided")
        return jsonify({"error": "No IDs provided"}), 400
    
@app.route('/ready', methods=['GET'])
def ready():
    return jsonify({"status": "ready"}), 200


if __name__ == '__main__':
    # Start a separate thread to continuously process the queue
    queue_processor_thread = threading.Thread(target=process_queue)
    queue_processor_thread.daemon = True
    queue_processor_thread.start()

    # check if qdrant collection exists
    if requests.get(url=f"{qdrant_uri}/collections/{qdrant_collection_name}/exists").json()["result"]["exists"] == False:
        # create collection if not exists
        qdrant_client.create_collection(collection_name=qdrant_collection_name, 
                                        vectors_config=VectorParams(size=384, distance="Cosine"),
                                        on_disk_payload=True,)

    while True:
        try:
            requests.get(url=f"{infinity_api_url}/ready")
            neo4j_driver.verify_connectivity()
            break
        except:
            time.sleep(1)
            

    # Start Flask API
    app.run(debug=True, host='0.0.0.0')
