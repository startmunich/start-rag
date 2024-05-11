###
# Script which syncs the VectorDB database with the neo4j graph database.
# Script offers a REST Api endpoint with which a local queue can be filled with IDs.
# The queue is then processed by the script, the content of the respective IDs is collected from the neo4j database.
# Each individual content is bein preprocessed by langchain and inserted into a qdrant collection. 
# The embeddings are created with infinity.
###
# trigger build
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
from vectordb_sync_fcts import notion_to_qdrant, web_to_qdrant, slack_to_qdrant


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

# Redis instance for storing IDs to be processed
redis = Redis(host='redis', port=6379, db=0)

# model = SentenceTransformer(infinity_model)


def process_queue():
    app.logger.info(f"start process_queue")
    while True:
        if redis.llen('content_queue') != 0:

            

            task = json.loads(redis.rpop('content_queue'))

            print(f"One ID from redis queue popped")

            id_to_process = task["id"]

            if task["type"] == "notion":
            
                notion_to_qdrant(id_to_process)

            if task["type"] == "web":
                
                web_to_qdrant(id_to_process)

            if task["type"] == "slack":
                
                slack_to_qdrant(id_to_process)
            
                
                
                
        else:
            # wait 10 seconds before checking the queue again
            app.logger.info("Queue is empty, waiting for 10 seconds")
            time.sleep(10)

@app.route('/empty_redis', methods=['POST'])
def empty_redis():
    redis.delete('content_queue')
    return jsonify({"message": "Redis queue cleared"}), 200

@app.route('/redis_length', methods=['POST'])
def redis_length():
    len = redis.llen('content_queue')
    return jsonify({"message": len}), 200

@app.route('/empty_qdrant', methods=['POST'])
def empty_qdrant():
    qdrant_client.delete_collection(collection_name=qdrant_collection_name)
    return jsonify({"message": "Qdrant collection cleared"}), 200


@app.route('/enqueue', methods=['POST']) # rename to enqueue_notion together with notion_Crwaler
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
    
@app.route('/enqueue_web', methods=['POST'])
def enqueue_web(): # make data to dict with id and type before pushing to redis list
    app.logger.info('json payload')
    app.logger.info(request.json)
    data = request.json
    if 'ids' in data:
        for id_to_enqueue in data['ids']:
            sendeable_data = json.dumps({"id": id_to_enqueue, "type": "web"})
            redis.lpush("content_queue", sendeable_data)
            print(f"ID {id_to_enqueue} enqueued successfully")
        return jsonify({"message": "IDs for web enqueued successfully"}), 200
    else:
        print("No IDs provided")
        return jsonify({"error": "No IDs for web provided"}), 400

@app.route('/enqueue_slack', methods=['POST'])
def enqueue_slack(): # make data to dict with id and type before pushing to redis list
    app.logger.info('json payload')
    app.logger.info(request.json)
    data = request.json
    if 'ids' in data:
        for id_to_enqueue in data['ids']:
            sendeable_data = json.dumps({"id" : id_to_enqueue, "type" : "slack"})
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
        # delete collection if exists
        # qdrant_client.delete_collection(collection_name=qdrant_collection_name)

        # create collection if not exists
        qdrant_client.create_collection(collection_name=qdrant_collection_name, 
                                        vectors_config=VectorParams(size=1024, distance="Cosine"),
                                        on_disk_payload=True,)

    while True:
        try:
            requests.get(url=f"{infinity_api_url}/health")
            neo4j_driver.verify_connectivity()
            break
        except:
            time.sleep(1)
            

    # Start Flask API
    app.run(debug=True, host='0.0.0.0')
