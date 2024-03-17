from langchain_community.document_loaders.async_html import AsyncHtmlLoader
from langchain_community.document_transformers import Html2TextTransformer
from langchain_text_splitters import RecursiveCharacterTextSplitter
from dotenv import load_dotenv
from neo4j import GraphDatabase
import logging
import os
import time

logger = logging.getLogger('start_gpt')

neo4j_uri = os.getenv("NEO4J_URL")
neo4j_user = os.getenv("NEO4J_USER")
neo4j_pass = os.getenv("NEO4J_PASS")
print(neo4j_uri)
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

    # Initialize urls to parse
    urls = [
        "https://www.startmunich.de",
        "https://www.startmunich.de/about-us",
        "https://www.startmunich.de/for-students",
        "https://www.startmunich.de/for-partners",
        "https://www.startmunich.de/membership",
        "https://www.startmunich.de/road-to-start-summit-2024",
        "https://www.startmunich.de/road-to-start-hack",
        "https://www.futurepynk.com"
    ]

    # Initialize loader with urls, load docs
    loader = AsyncHtmlLoader(urls)
    docs = loader.load()

    # Process raw HTML docs to text docs
    html2text = Html2TextTransformer()
    docs = html2text.transform_documents(docs)
    for page in docs:
        print(page.page_content)
    logger.info("load_web finished")
    return docs
def write_db():
    logger.info(f"start writing web_data")
    with neo4j_driver.session() as session:
        docs = load_web()
        # write each page into the neo4J database using the url as the id
        # TODO ask Khadim about the hash and url attributes (to they need to be set or can they be blank?)
        for page in docs:
            logger.info(f"adding page: " + page.metadata.get("source") + " to neo4j")
            session.run(
                # merges the page with the current page_id
                # cf. notioncrawler/crawler/cache.go
                "MERGE (n:CrawledPage { page_id: $page_id })\n" +
                # if such an id doesn't yet exist create new node taking the pages attributes
                "ON CREATE SET n.crawlerId=$crawler_id, n.url=$url, n.content=$content, n.child_pages=$child_pages, n.hash=$hash\n" +
                # if such a node exists update the attributes that exists in page
                "ON MATCH SET n.crawlerId=$crawler_id, n.url=$url, n.content=$content, n.child_pages=$child_pages, n.hash=$hash\n",
                page_id=page.metadata.get("source"),
                crawler_id="Webcrawler",
                content=page.page_content
            )
            logger.info(f"page: " + page.metadata.get("source") + " added to neo4j")



if __name__ == '__main__':
    load_web()
