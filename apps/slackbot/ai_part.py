import os
from langchain.chains import LLMChain
from langchain.prompts import PromptTemplate
from langchain_community.llms import replicate
from langchain.callbacks.streaming_stdout import StreamingStdOutCallbackHandler
from langchain.chains import RetrievalQA
from langchain_community.vectorstores.qdrant import Qdrant
from qdrant_client import QdrantClient
from langchain_core.vectorstores import VectorStoreRetriever

replicate_api_key = os.environ["REPLICATE_API_KEY"]
qdrant_uri = os.environ["QDRANT_URL"]


qdrant_db = Qdrant(
    client=QdrantClient(url=qdrant_uri, port=6333),
    collection_name="startgpt",
    content_payload_key="content",
    metadata_payload_key="page_id",
    distance_strategy="Cosine",
)

retriever = VectorStoreRetriever(vector_store=qdrant_db)

llm = replicate(
    streaming=True,
    callbacks=[StreamingStdOutCallbackHandler()],
    model="some_model",
    model_kwargs={"temperature": 0.75, "max_length": 500, "top_p": 1},
)

prompt = """
User: Answer the following yes/no question by reasoning step by step. Can a dog drive a car?
Assistant:
"""
_ = llm(prompt)