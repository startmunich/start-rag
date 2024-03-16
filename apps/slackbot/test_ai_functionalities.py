from ai_functionalitites import get_answer
import os
import unittest


# test the get_answer function

class TestAIFunctionalities(unittest.TestCase):
    def test_get_answer(self):

        start_munich_question = "What is the purpose of start munich?"

        # check if an answer is returned
        self.assertIsNotNone(get_answer(start_munich_question))

        # check if answer is a string
        self.assertIsInstance(get_answer(start_munich_question), str)

        # check if answer is not empty
        self.assertNotEqual(get_answer(start_munich_question), "")

        # print the answer
        print(get_answer(start_munich_question))

if __name__ == "__main__":

    unittest.main()

