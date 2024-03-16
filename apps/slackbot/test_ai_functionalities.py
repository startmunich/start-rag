from ai_functionalitites import get_answer
import os
import unittest


# test the get_answer function

class TestAIFunctionalities(unittest.TestCase):

    start_munich_question = "What is the purpose of start munich?"

    def test_answer_none(self, start_munich_question=start_munich_question):
        # check if an answer is returned
        answer = get_answer(start_munich_question)
        print(f"LLM answer is: {answer}")
        self.assertIsNotNone(answer)

    def test_answer_type(self, start_munich_question=start_munich_question):
        # check if answer is a string
        answer = get_answer(start_munich_question)
        print(f"LLM answer is: {answer}")
        self.assertIsInstance(answer, str)


    def test_answer_not_empty(self, start_munich_question=start_munich_question):
        # check if answer is not empty
        answer = get_answer(start_munich_question)
        print(f"LLM answer is: {answer}")
        self.assertNotEqual(answer, "")

if __name__ == "__main__":

    unittest.main()

