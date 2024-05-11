
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


# Initialize Qdrant client
qdrant_uri = os.environ.get("QDRANT_URL")
qdrant_collection_name = os.environ.get("QDRANT_COLLECTION_NAME")
qdrant_client = QdrantClient(url=qdrant_uri,port=6333)

# initialize infinity
infinity_api_url = os.environ.get("INFINITY_URL")
infinity_model = os.environ.get("INFINITY_MODEL")

# Initialize Neo4j driver
neo4j_uri = os.environ.get("NEO4J_URL")
neo4j_user = os.environ.get("NEO4J_USER")
neo4j_password = os.environ.get("NEO4J_PASS")
neo4j_driver = GraphDatabase.driver(neo4j_uri, auth=(neo4j_user, neo4j_password))


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

def notion_to_qdrant(id_to_process) -> None:
    with neo4j_driver.session() as session:
        result = session.run(
            "MATCH (n:CrawledPage {page_id: $id}) RETURN n.content AS content",
            id=id_to_process,
        )
        
        record = result.single()
        
        
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

        
        chunks_embedded = [requests.post(url=f"{infinity_api_url}/embeddings", json={"model": f"{infinity_model}", "input":[chunk]}).json()["data"][0]["embedding"] for chunk in chunks]
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
                                        payload={"content": chunk, "metadata":{"page_id": id_to_process, "type": "notion"}}) for chunk_embedding, chunk in zip(chunks_embedded, chunks)]
        
        

        qdrant_client.upsert(
            collection_name=qdrant_collection_name,
            points= points_to_update
            )
        
def web_to_qdrant(id_to_process):
    ## TODO
    # 1. get content from neo4j
    # 2. preprocess content as far as necessary
    # 3. embedd content with infinity
    # 4. insert into qdrant

    id_to_process = f'{id_to_process}'
    with neo4j_driver.session() as session:
        # returns the entry from the neo4j db for this id
        result = session.run(
            "MATCH (n:CrawledPage {page_id: $id}) RETURN n.content AS content",
            id=id_to_process,
        )
        # check if result is empty, then throw error
        if result.peek() is None:
            raise ValueError("No result found for id_to_process for web_to_qdrant")

        record = result.single()

        record_content = record.get("content")
        content = str(record_content)
        complete_content = []

        complete_content += text_splitter.split_text(text=content)

        # check if complete_content is empty, then throw error
        if len(complete_content) == 0:
            # stop function if no complete_content is found
            raise ValueError("No complete_content on web_to_qdrant found")

        # lower the case of the chunks in chunks
        chunks = [chunk.lower() for chunk in complete_content]

        # Embed the chunks using Infinity
        chunks_embedded = [requests.post(url=f"{infinity_api_url}/embeddings",
                                         json={"model": f"{infinity_model}", "input": [chunk]}).json()["data"][0]["embedding"] for chunk in chunks]
        print("embeddings for web_to_qdrant created successful")

        # delete all points created from web pages
        qdrant_client.delete(collection_name=qdrant_collection_name,
                             points_selector=FilterSelector(
                                 filter=Filter(
                                     must=[
                                         FieldCondition(
                                             key="type",
                                             match=MatchValue(value="web"),
                                         ),
                                     ],
                                 )
                             ),
                             )

        # Insert the preprocessed chunk into Qdrant
        points_to_update = [PointStruct(id=str(uuid.uuid4()),
                                        vector=chunk_embedding,
                                        payload={"content": chunk, "metadata":{"page_id": id_to_process, "type": "web"}}) for
                            chunk_embedding, chunk in zip(chunks_embedded, chunks)]
        
        qdrant_client.upsert(
            collection_name=qdrant_collection_name,
            points=points_to_update
        )

def slack_to_qdrant(id_to_process):
    ## TODO
    # 1. get content from neo4j
    # 2. preprocess content as far as necessary
    # 3. see if concatenation of some messages is necessary depending on their length
    # 4. embedd content with infinity
    # 5. insert into qdrant
    id_to_process = f'{id_to_process}'
    with neo4j_driver.session() as session:
        result = session.run(
            "MATCH (n:CrawledPage {page_id: $id}) RETURN n.content AS content",
            id=id_to_process,
        )

        # check if result is empty, then throw error
        if result.peek() is None:
            raise ValueError("No result found for id_to_process")
        
        record = result.single()
        
        record_content = record.get("content") # should be string?
        content = str(record_content)

        


        complete_content = []


        complete_content += text_splitter.split_text(text=content)

        # check if complete_content is empty, then throw error
        if len(complete_content) == 0:
            # stop function if no complete_content is found
            raise ValueError("No complete_content found")

            
        # lower the case of the chunks in chunks
        chunks = [chunk.lower() for chunk in complete_content]

        


        # Embed the chunks using Infinity

        
        chunks_embedded = [requests.post(url=f"{infinity_api_url}/embeddings", json={"model": f"{infinity_model}", "input":[chunk]}).json()["data"][0]["embedding"] for chunk in chunks]
        print("embeddings created successful")

        
        
        

        # delete all points created from slack messages
        qdrant_client.delete(collection_name=qdrant_collection_name,     
                            points_selector=FilterSelector(
                                filter=Filter(
                                    must=[
                                        FieldCondition(
                                            key="type",
                                            match=MatchValue(value="slack"),
                                        ),
                                    ],
                                )
                            ),
                            )

        # Insert the preprocessed chunk into Qdrant
        
        points_to_update = [PointStruct(id=str(uuid.uuid4()), 
                                        vector=chunk_embedding,
                                        payload={"content": chunk, "metadata":{"page_id": id_to_process, "type": "slack"}}) for chunk_embedding, chunk in zip(chunks_embedded, chunks)]
        
        
        

        qdrant_client.upsert(
            collection_name=qdrant_collection_name,
            points= points_to_update
            )
