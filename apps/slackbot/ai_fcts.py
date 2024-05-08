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
from datetime import date





# Load and export the environment variables
# load_dotenv(dotenv_path="apps/slackbot/.env.local")

# trigger new build

# REPLICATE_API_TOKEN  = os.environ["REPLICATE_API_KEY"]
qdrant_uri = os.environ["QDRANT_URL"]
qdrant_collection_name = os.environ.get("QDRANT_COLLECTION_NAME")
infinity_api_url = os.environ.get("INFINITY_URL")
infinity_model = os.environ.get("INFINITY_MODEL")
llm_model = os.environ.get("LLM_MODEL")

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
retriever = qdrant_db.as_retriever(search_kwargs={"k": 5})

# Create the language model
llm = Replicate(
    streaming=True,
    callbacks=[StreamingStdOutCallbackHandler()],
    model=llm_model,
    model_kwargs={"temperature": 0.2, "max_length": 1500, "top_p": 0.9, "top_k": 50, "max_new_tokens": 400,
        "min_new_tokens": 20, "repetition_penalty": 0.1},
    verbose = False
)



# Create the prompt template
prompt_template = """ [INST]
You are StartGPT, an assistant for question-answering tasks.
The context you get will be from our Notion and Slack. Summarize the context and answer the question. Add whenever possible a link to the corrsponding Notion page.

<Beginning of context>
{context} 
<End of context>

<Beginning of question>
{question}
<End of question>
[INST]
"""

# Initialize prompt
prompt_template = PromptTemplate(input_variables=["question", "context"], template=prompt_template)

# prompt = hub.pull("rlm/rag-prompt")

# qa_chain = RetrievalQA.from_chain_type(
#     llm=llm,
#     retriever=retriever,
#     chain_type_kwargs={"prompt": prompt_template},
#     verbose=False
# )
def format_docs(docs):
    context = ""
    for index, doc in enumerate(docs):
        context += f"Document Rank {index + 1}: {doc.page_content}\n\n"
    return context

# create function to invoke the retrievalQA
def get_answer(query: str) -> str:
    
    docs = retriever.invoke(query)

    # create string of all the documents
    
    context = format_docs(docs)

    # enter context into the prompt

    prompt = prompt_template.format(context=context, question=query)

    response = llm.invoke(prompt)
    print(f"Retrieved docs: {docs}\n Retrieved context: {context}\n Generated response: {response}")

    return response



