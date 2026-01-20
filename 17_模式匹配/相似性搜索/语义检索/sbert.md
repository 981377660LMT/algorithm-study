SBERT (Sentence-BERT) 是目前最主流的语义相似度算法，它通过孪生网络（Siamese Network）将句子映射为固定维度的稠密向量，使得语义相近的句子在向量空间中的距离更近。

### 1. 基础实现：计算两个句子的相似度

使用 `sentence-transformers` 库，这是 SBERT 的官方实现。

```python
# 安装: pip install sentence-transformers
from sentence_transformers import SentenceTransformer, util

# 1. 加载预训练模型 (all-MiniLM-L6-v2 是一个兼顾速度与效果的轻量级模型)
model = SentenceTransformer('all-MiniLM-L6-v2')

sentences = [
    "The cat sits outside",
    "A man is playing guitar",
    "The feline is resting outdoors"
]

# 2. 将句子编码为向量 (Embeddings)
embeddings = model.encode(sentences)

# 3. 计算余弦相似度
# 计算第 0 个句子 ("The cat sits outside") 和第 2 个句子 ("The feline is resting outdoors")
# 尽管字面上没有重复词，但语义高度相似
cosine_score = util.cos_sim(embeddings[0], embeddings[2])

print(f"语义相似度得分: {cosine_score.item():.4f}")
```

### 2. 进阶实现：大规模语义搜索 (Semantic Search)

当你有海量文档（Corpus）需要检索时，可以使用 `util.semantic_search`。

```python
from sentence_transformers import SentenceTransformer, util
import torch

model = SentenceTransformer('all-MiniLM-L6-v2')

corpus = [
    "A man is eating food.",
    "A man is eating a piece of bread.",
    "The girl is carrying a baby.",
    "A man is riding a horse.",
    "A woman is playing violin.",
    "Two men pushed carts through the woods.",
    "A man is riding a white horse on an enclosed ground.",
    "A monkey is playing drums.",
    "Someone in a gorilla costume is playing a set of drums."
]

corpus_embeddings = model.encode(corpus, convert_to_tensor=True)

query = "A man eats pasta"
query_embedding = model.encode(query, convert_to_tensor=True)

# 使用 semantic_search 函数进行快速检索
# 它会自动计算余弦相似度并排序
hits = util.semantic_search(query_embedding, corpus_embeddings, top_k=3)


print(f"查询语句: {query}")
for hit in hits[0]:
    print(f"匹配文档: {corpus[hit['corpus_id']]} (得分: {hit['score']:.4f})")
```

### 核心优势

1.  **语义理解**：能识别“猫”和“猫科动物”、“吃饭”和“吃面”的关联，这是 BM25 等统计算法做不到的。
2.  **高性能**：预先计算好 Corpus 的向量并存储（如存入向量数据库），检索时只需计算一次 Query 向量，适合实时系统。
3.  **多语言支持**：使用 `paraphrase-multilingual-MiniLM-L12-v2` 等模型可以实现跨语言搜索（如用英文搜中文内容）。
