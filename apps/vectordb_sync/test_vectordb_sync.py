import os
import queue
import unittest
import uuid
import time

import requests
from qdrant_client.http import models
from scipy import spatial
from qdrant_client import QdrantClient
from qdrant_client.http.models import PointStruct, FilterSelector, Filter, FieldCondition, MatchValue, VectorParams
from langchain_community.embeddings import InfinityEmbeddings

# initialize infinity
# infinity_api_url = os.environ.get("INFINITY_URL")
infinity_api_url = "http://infinity:7997"
# infinity_model = os.environ.get("INFINITY_MODEL")
infinity_model = "BAAI/bge-small-en"

# Initialize Qdrant client
# qdrant_uri = os.environ.get("QDRANT_URL")
qdrant_uri = "http://qdrant:6333"
# qdrant_collection_name = os.environ.get("QDRANT_COLLECTION_NAME")
qdrant_collection_name = "startgpt"
qdrant_client = QdrantClient(url=qdrant_uri, port=6333)

sentence1 = 'Hallo, wie gehts?'
sentence2 = 'Ich mag Hunde'
sentence3 = 'Ich mag Katzen'

sentences = [sentence1, sentence2, sentence3]
while True:
    try:
        requests.get(url=f"{infinity_api_url}/ready")
        break
    except:
        time.sleep(1)
sentences_embedded = [requests.post(url=f"{infinity_api_url}/embeddings",
                                   json={"model": "bge-small-en-v1.5", "input": [sentence]}).json()["data"][0][
                          "embedding"] for sentence in sentences]


page_ids = [0, 1, 2]
points_to_update = [PointStruct(id=str(uuid.uuid4()),
                                vector=sentence_embedding,
                                payload={"content": sentence, "page_id": page_id}) for
                    sentence_embedding, sentence, page_id
                    in zip(sentences_embedded, sentences, page_ids)]


class EmbeddingTest(unittest.TestCase):

    def test_infinity_embedding(self):

        print(infinity_api_url)

        result_12 = 1 - spatial.distance.cosine(sentences_embedded[0], sentences_embedded[1])
        result_13 = 1 - spatial.distance.cosine(sentences_embedded[0], sentences_embedded[2])
        result_23 = 1 - spatial.distance.cosine(sentences_embedded[1], sentences_embedded[2])

        print(result_12)
        print(result_13)
        print(result_23)
        self.assertTrue(result_23 > result_12)


class QdrantTest(unittest.TestCase):

    def test_post_data(self):
        # inserts points into database
        qdrant_client.upsert(
            collection_name=qdrant_collection_name,
            points=points_to_update
        )

        # should return list of the points that have been inserted
        output = QdrantTest.find_elements_by_ids(0, 3)

        # checks if the objects in both lists are the same
        for x in range(0, 3):
            self.assertTrue(points_to_update[x] in output)

    def test_delete_data(self):

        # deletes the points with the ids 0,1,2
        for x in range(0, 3):
            qdrant_client.delete(collection_name=qdrant_collection_name,
                                 points_selector=FilterSelector(
                                     filter=Filter(
                                         must=[
                                             FieldCondition(
                                                 key="page_id",
                                                 match=MatchValue(value=x),
                                             ),
                                         ],
                                     )
                                 ),
                                 )

        # should return empty list since points should be deleted
        del_elem = QdrantTest.find_elements_by_ids(0, 3)
        # checks if list is empty
        self.assertTrue(not del_elem)

    def test_retrievel_of_points(self):

        # inserts points into database
        qdrant_client.upsert(collection_name=qdrant_collection_name,
                             points=points_to_update)

        # for each point with ids 0,1,2
        for x in range(0, 3):

            # retrieve the point
            point = qdrant_client.scroll(collection_name=qdrant_collection_name,
                                         scroll_filter=models.Filter(
                                             must=[
                                                 models.FieldCondition(
                                                     key="page_id",
                                                     match=models.MatchValue(value=x)
                                                 )
                                             ]
                                         ))[0]
            # check if equal to point that should have been inserted
            self.assertTrue(point == points_to_update[x])


    # function retrieving points in id range
    def find_elements_by_ids(lower_id, upper_id):
        output = []
        for x in range(lower_id, upper_id):
            result = qdrant_client.scroll(
                collection_name=qdrant_collection_name,
                scroll_filter=models.Filter(
                    must=[
                        models.FieldCondition(
                            key="page_id",
                            match=models.MatchValue(value=x),
                        )
                    ]
                ),
            )
            output = output + result
        return output


if __name__ == '__main__':
    unittest.main()
