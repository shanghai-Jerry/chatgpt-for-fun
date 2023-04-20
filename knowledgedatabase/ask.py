from gpt_index import SimpleDirectoryReader, GPTListIndex, readers, GPTSimpleVectorIndex, LLMPredictor, PromptHelper
from langchain import OpenAI
import sys
import os
from IPython.display import Markdown, display


def ask_question(text):
    index = GPTSimpleVectorIndex.load_from_disk('index.json')
    query = text
    response = index.query(query, response_mode="compact")
    print("Bot says:%s" % response.response)


if __name__ == '__main__':
    # 检查参数个数是否符合要求
    if len(sys.argv) < 1:
        print("Usage: python script.py arg1")
        sys.exit(1)

    # API Key for OpenAI
    os.environ["OPENAI_API_KEY"] = input(
        "Paste your OpenAI API key here and hit enter:")
    # 通过索引访问参数
    arg1 = sys.argv[1]
    ask_question(arg1)
