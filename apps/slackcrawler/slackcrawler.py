import os
import time
import json
from slack import WebClient
from neo4j import GraphDatabase
import logging
import uuid
import requests

logger = logging.getLogger('start_gpt')

neo4j_uri = os.getenv("NEO4J_URL")
neo4j_user = os.getenv("NEO4J_USER")
neo4j_pass = os.getenv("NEO4J_PASS")

# Define Slack credentials
## TODO: Add rleevant slack channels
SLACK_BOT_TOKEN = os.environ.get("SLACK_BOT_TOKEN")
# load channels from json file and store them in SLACK_CHANNELS
SLACK_CHANNELS = []
with open('relevant_channels.json') as f:
    SLACK_CHANNELS = list(dict(json.loads(f)).keys())

slack_client = WebClient(token=SLACK_BOT_TOKEN)


# Initialize Neo4j connection
neo4j_driver = GraphDatabase.driver(neo4j_uri, auth=(neo4j_user, neo4j_pass))


## TODO: Keep track of already crawled messages to avoid duplicates
def crawl_and_store_data():
    
    # Crawl messages from specific Slack channels
    for channel_id in SLACK_CHANNELS:
        response = slack_client.conversations_history(channel=channel_id, limit= 100)
        message_history = response['messages'] # list of messages

        message_ids = []

        for message in message_history:
            timestamp = message['ts']
            content = message['text']
            message_id = uuid.UUID(channel_id+timestamp)
            message_ids.append(message_id)

            with neo4j_driver.session() as session:
                logger.info(f"adding message: " + message_id + " to neo4j")
                session.run(
                    # merges the message with current message_id
                    # cf. notioncrawler/crawler/cache.go
                    "MERGE (n:CrawledPage { page_id: $message_id })\n" +
                    # if such an id doesn't yet exist create new node taking the pages attributes
                    "ON CREATE SET n.crawlerId=$crawler_id, n.content=$content, n.message_id =$message_id\n" +
                    # if such a node exists update the attributes that exists in page
                    "ON MATCH SET n.crawlerId=$crawler_id, n.content=$content, n.message_id =$message_id\n",
                    message_id=message_id,
                    crawler_id="Slackcrawler",
                    content=content
                )
                logger.info(f"adding message: " + message_id + " added to neo4j")

            # send post request to '/enqueue_slack' endpoint with message_id
            # to add the message to the queue
            # cf. vectordb_sync/vectordb_sync.py
            
        requests.post("http://vectordb_sync:5000/enqueue_slack", json={"ids": message_ids})
        logger.info(f"adding message_ids of channel: " + channel_id + " to queue")


            


            # Store the message in Neo4j if it is not already stored




if __name__ == "__main__":
    # Perform the crawling and storing process

    while True:
        try:
            neo4j_driver.verify_connectivity()
            break
        except:
            time.sleep(1)

    crawl_and_store_data()

    # Sleep for 10 minutes before running again
    time.sleep(10 * 60) 
