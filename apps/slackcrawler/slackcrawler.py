import os
import time
from datetime import datetime
from llama_index import download_loader
from neo4j import GraphDatabase
import logging

logger = logging.getLogger('start_gpt')

neo4j_uri = os.getenv("NEO4J_URL")
neo4j_user = os.getenv("NEO4J_USER")
neo4j_pass = os.getenv("NEO4J_PASS")

# Define Slack credentials
## TODO: Add rleevant slack channels
SLACK_BOT_TOKEN = os.environ.get("SLACK_BOT_TOKEN")
SLACK_CHANNELS = ["", ""]  

reader = download_loader("SlackReader")
loader = reader(SLACK_BOT_TOKEN)


# Initialize Neo4j connection
neo4j_driver = GraphDatabase.driver(neo4j_uri, auth=(neo4j_user, neo4j_pass))


## TODO: Keep track of already crawled messages to avoid duplicates
def crawl_and_store_data():
    # Initialize Slack crawler
    

    # Crawl messages from specific Slack channels
    for channel_ids in SLACK_CHANNELS:
        lama_docs = loader.load_data(channel_ids)
        lang_docs = [doc.to_langchain_format() for doc in lama_docs]

        # Store messages in Neo4j
        for doc in lang_docs:
            store_doc_in_neo4j(doc)

## TODO: Adapt Khadims code to pythons cript or ask dawid about the functionality of his code 
def store_doc_in_neo4j(message):
    pass


if __name__ == "__main__":
    # Perform the crawling and storing process
    crawl_and_store_data()

    while True:
        try:
            neo4j_driver.verify_connectivity()
            break
        except:
            time.sleep(1)

    # Sleep for 5 minutes before running again
    time.sleep(5 * 60)  # 5 minutes in seconds
