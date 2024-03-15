import os
from langchain.prompts import PromptTemplate
from langchain_community.llms import replicate
from langchain.callbacks.streaming_stdout import StreamingStdOutCallbackHandler
from langchain.chains import RetrievalQA
from langchain_community.vectorstores.qdrant import Qdrant
from qdrant_client import QdrantClient
from langchain_community.embeddings import InfinityEmbeddings


# Load the environment variables
replicate_api_key = os.environ["REPLICATE_API_KEY"]
qdrant_uri = os.environ["QDRANT_URL"]
qdrant_collection_name = os.environ.get("QDRANT_COLLECTION_NAME")
infinity_api_url = os.environ.get("INFINITY_URL")
infinity_model = os.environ.get("INFINITY_MODEL")

# Create the Qdrant vector store
qdrant_db = Qdrant(
    client=QdrantClient(url=qdrant_uri, port=6333),
    collection_name=qdrant_collection_name,
    content_payload_key="content",
    metadata_payload_key="page_id",
    distance_strategy="Cosine",
    embedding_function=InfinityEmbeddings(model=infinity_model, infinity_api_url=infinity_api_url)
)

# Create the retriever
retriever = qdrant_db.as_retriever("mmr", search_kwargs={"k": 5, "fetch_k": 20})

# Create the language model
llm = replicate(
    streaming=True,
    callbacks=[StreamingStdOutCallbackHandler()],
    model="some_model",
    model_kwargs={"temperature": 0.75, "max_length": 500, "top_p": 1},
)

# Create the prompt template
prompt_template = """ 
It is December 2023. You are StartGPT, an assistant for question-answering tasks. Users reach out to you only via Slack. You serve a student led organization START Munich. 
The context you get will be from our Notionpage. Use the following pieces of retrieved context to answer the question. 
You decide what's more useful. If you don't know the answer, just say that you don't know.
Here's the question and the context:

<Beginning of question>
{question}
<End of question>

<Beginning of context>
{notion} 
<End of context>
"""

# Initialize prompt
prompt_template = PromptTemplate(input_variables=["question", "notion"], template=prompt_template)

qa_chain = RetrievalQA.from_chain_type(
    llm=llm,
    retriever=retriever,
    chain_type_kwargs={"prompt": prompt_template},
)

# create function to invoke the retrievalQA
def get_answer(question: str) -> str:
    
    response = qa_chain({"query": question})
    return response["result"]
