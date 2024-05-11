from langchain_community.document_loaders.async_html import AsyncHtmlLoader
from langchain_community.document_transformers import Html2TextTransformer
from neo4j import GraphDatabase
import logging
import os
import time
import requests
from bs4 import BeautifulSoup

logger = logging.getLogger('start_gpt')

neo4j_uri = os.getenv("NEO4J_URL")
neo4j_user = os.getenv("NEO4J_USER")
neo4j_pass = os.getenv("NEO4J_PASS")
neo4j_driver = GraphDatabase.driver(neo4j_uri, auth=(neo4j_user, neo4j_pass))
# checks if connection to database can be established, else throws error
while True:
    try:
        neo4j_driver.verify_connectivity()
        break
    except:
        time.sleep(1)


def load_web():
    logger.info("load_web")
    print("Loading of webpages started")
    # Initialize urls to parse
    urls = [
        "https://www.startmunich.de",
        "https://www.startmunich.de/about-us",
        "https://www.startmunich.de/for-students",
        "https://www.startmunich.de/for-partners",
        "https://www.startmunich.de/apply",
        "https://www.futurepynk.com",
        "https://www.futurepynk.com/theroyaljungle",
        "https://www.futurepynk.com/hamburg",
        "https://www.futurepynk.com/munich",
        "https://www.futurepynk.com/berlin",
    ]

    # Initialize loader with urls, load docs
    loader = AsyncHtmlLoader(urls)
    docs = loader.load()

    # Process raw HTML docs to text docs
    html2text = Html2TextTransformer()
    docs = html2text.transform_documents(docs)
    print("Loading of webpages finished")
    logger.info("load_web finished")
    return docs


def write_db():
    logger.info(f"start writing web_data")
    print("Start writing web_data to neo4j")

    docs = load_web()
    # write each page into the neo4J database using the url as the id
    web_ids = []

    for page in docs:

        web_ids.append(page.metadata.get("source"))

        with neo4j_driver.session() as session:
            logger.info(f"adding page: " + page.metadata.get("source") + " to neo4j")
            session.run(
                # merges the page with the current page_id
                # cf. notioncrawler/crawler/cache.go
                "MERGE (n:CrawledPage { page_id: $page_id })\n" +
                # if such an id doesn't yet exist create new node taking the pages attributes
                "ON CREATE SET n.crawlerId=$crawler_id, n.content=$content, n.page_id=$page_id\n" +
                # if such a node exists update the attributes that exists in page
                "ON MATCH SET n.crawlerId=$crawler_id, n.content=$content, n.page_id=$page_id\n",
                page_id=page.metadata.get("source"),
                crawler_id="Webcrawler",
                content=page.page_content
            )
        logger.info(f"page: " + page.metadata.get("source") + " added to neo4j")
        print("page: " + page.metadata.get("source") + " added to neo4j")

        # send post request to '/enqueue_web' endpoint with web_ids
        # to add the webpage to the queue
        # cf. vectordb_sync/vectordb_sync.py

    requests.post("http://vectordb_sync:5000/enqueue_web", json={"ids": web_ids})
    logger.info(f"adding web_ids: to queue")
    print("adding web_ids: to queue")


if __name__ == '__main__':
    # Perform the crawling and storing process

    while True:
        try:
            neo4j_driver.verify_connectivity()
            requests.get(url="http://vectordb_sync:5000/ready")
            break
        except:
            time.sleep(1)

    write_db()
    # Sleep for 1 day before running again
    time.sleep(24 * 3600)
