import os
from dotenv import load_dotenv
from langchain.prompts import PromptTemplate
from langchain_community.llms import Replicate
from langchain.callbacks.streaming_stdout import StreamingStdOutCallbackHandler
from langchain.chains import RetrievalQA
from langchain_community.vectorstores.qdrant import Qdrant
from qdrant_client import QdrantClient
from langchain_community.embeddings import InfinityEmbeddings
from langchain import hub




# Load and export the environment variables
load_dotenv(dotenv_path="apps/slackbot/.env.local")



# REPLICATE_API_TOKEN  = os.environ["REPLICATE_API_KEY"]
qdrant_uri = os.environ["QDRANT_URL"]
qdrant_collection_name = os.environ.get("QDRANT_COLLECTION_NAME")
infinity_api_url = os.environ.get("INFINITY_URL")
infinity_model = os.environ.get("INFINITY_MODEL")

# Create the Qdrant vector store
qdrant_db = Qdrant(
    client=QdrantClient(url=qdrant_uri, port=6333),
    collection_name=qdrant_collection_name,
    content_payload_key="content",
    metadata_payload_key=None,
    distance_strategy="Cosine",
    embeddings=InfinityEmbeddings(model=infinity_model, infinity_api_url=infinity_api_url)
)

# Create the retriever
retriever = qdrant_db.as_retriever(search_type = "mmr", search_kwargs={"k": 5, "fetch_k": 20})

# Create the language model
llm = Replicate(
    streaming=True,
    callbacks=[StreamingStdOutCallbackHandler()],
    model="a16z-infra/llama13b-v2-chat:df7690f1994d94e96ad9d568eac121aecf50684a0b0963b25a41cc40061269e5",
    model_kwargs={"temperature": 0.75, "max_length": 500, "top_p": 1}
)

# Create the prompt template
prompt_template = """ [INST]
It is December 2023. You are StartGPT, an assistant for question-answering tasks. Users reach out to you only via Slack. You serve a student led organization START Munich. 
The context you get will be from our Notionpage. Use the following pieces of retrieved context to answer the question. 
You decide what's more useful. If you don't know the answer, just say that you don't know.
Here's the question and the context:

<Beginning of question>
{query}
<End of question>

<Beginning of context>
{context} 
<End of context>
[INST]
"""

# Initialize prompt
prompt_template = PromptTemplate(input_variables=["query", "context"], template=prompt_template)

prompt = hub.pull("rlm/rag-prompt")

qa_chain = RetrievalQA.from_chain_type(
    llm=llm,
    retriever=retriever,
    chain_type_kwargs={"prompt": prompt},
)

# create function to invoke the retrievalQA
def get_answer(query: str) -> str:
    
    response = qa_chain.invoke({'query': query})
    return response["result"]
