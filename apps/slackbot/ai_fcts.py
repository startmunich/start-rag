import os
from dotenv import load_dotenv
from langchain.prompts import PromptTemplate
from openai import OpenAI
from langchain.callbacks.streaming_stdout import StreamingStdOutCallbackHandler
from langchain_community.vectorstores.qdrant import Qdrant
from qdrant_client import QdrantClient
from langchain_community.embeddings import InfinityEmbeddings





# Load and export the environment variables
# load_dotenv(dotenv_path="apps/slackbot/.env.local")

# trigger new build
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
    metadata_payload_key="metadata",
    distance_strategy="Cosine",
    embeddings=InfinityEmbeddings(model=infinity_model, infinity_api_url=infinity_api_url)
)

# Create the retriever
retriever = qdrant_db.as_retriever(search_kwargs={"k": 3})

# Create the language model
client = OpenAI(
    api_key=os.environ.get("OPENAI_KEY"),
    )


def format_docs(docs):
    context = ""
    for index, doc in enumerate(docs):
        if doc.metadata['type'] == 'notion:':
            context += f"Document Rank {index + 1}. Source: https://www.notion.so/{doc.metadata['page_id']}. Content: {doc.page_content}\n\n"
        else:
            context += f"Document Rank {index + 1}. Source: {doc.metadata['type']}. Content: {doc.page_content}\n\n"
    return context

# create function to invoke the retrievalQA
def get_answer(query: str) -> str:

    enhanced_query = client.chat.completions.create(
        model="gpt-4o",
        messages=[
            {"role": "system", "content": """You are an assistant to StartGPT, an assistant for question-answering tasks of the student initative STARTMunich and your job is to add more detailed inquiry to the questions provided by a user. Here are some information about the initiative: STARTMunich is a student led initiative empowering students and young founders on their entrepreneurial journey. They organize various events and hackathons to help students extend their network and follow up their dream of becoming a founder themselves. The main database used is Notion which contains most of the information regarding the members, events, an FAQ for often asked question and much more.
            In the following you will get a question from one of the members of the initiative and should extend the question by adding more details to it like asking about a specific time frame, people involved, the location and current status so it can be querried better."""}, # <-- This is the system message that provides context to the model
            {"role": "user", "content": f"""<Beginning of question> {query} <End of question>"""}  # <-- This is the user message for which the model will generate a response
        ]
    )
    question = enhanced_query.choices[0].message.content
    docs = retriever.invoke(question)

    # create string of all the documents
    
    context = format_docs(docs)

    # enter context into the prompt
    completion = client.chat.completions.create(
        model="gpt-4o",
        messages=[
            {"role": "system", "content": """You are StartGPT, an assistant for question-answering tasks.
            The context you get will be from our Notion, Website and Slack. Use this context to answer the question.
            If you utilize context from Notion to answer the question, please provide the source link in your answer.
            If you utilize context from slack or the website, do not provide a link."""}, # <-- This is the system message that provides context to the model
            {"role": "user", "content": f"""<Beginning of context> {context} <End of context> \n
            <Beginning of question> {question} <End of question>"""}  # <-- This is the user message for which the model will generate a response
        ]
    )
    
    response = completion.choices[0].message.content
    # print(f"Retrieved docs: {docs}\n Retrieved context: {context}\n Generated response: {response}")

    return f"{response}"



