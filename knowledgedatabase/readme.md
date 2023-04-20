# 知识库

实践：如何构建基于指定数据集下的Q&A问答机器人, 让机器人来帮你消化这些数据内容，有问题来问机器人。

## [GPT-Index](https://gpt-index.readthedocs.io/en/latest/)

用来将数据内容进行embedding， 方便检索和匹配。 

通过下面的方式配置依赖OpenAI的API key

```bash
os.environ["OPENAI_API_KEY"] = input("Paste your OpenAI API key here and hit enter:")
```

### 安装依赖
如果你本地pyton环境是python3.7的话，最好升级到3.8，避免一些不可预知的问题发生

```bash
pip3 install gpt_index
pip3 install langchain
```