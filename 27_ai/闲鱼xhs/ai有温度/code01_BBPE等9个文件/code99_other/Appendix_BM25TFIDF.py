from sklearn.feature_extraction.text import TfidfVectorizer
import numpy as np
from rank_bm25 import BM25Okapi
 
# 示例文档
corpus = [
    "the cat sat on the mat",
    "the cat sat on the mat mat mat mat",
    "the dog chased the cat",
    "cats and dogs are friends"
]
 
# 查询
query = "cat mat".split()
 
# 1. TF-IDF
vectorizer = TfidfVectorizer()
tfidf_matrix = vectorizer.fit_transform(corpus)
 
# 将查询向量化
query_vec = vectorizer.transform(["cat mat"])
 
# 计算余弦相似度
cosine_sim = (tfidf_matrix @ query_vec.T).toarray().flatten()
 
print("TF-IDF 排序结果:")
for idx in np.argsort(-cosine_sim):
    print(f"Doc {idx}: {corpus[idx]} -> Score {cosine_sim[idx]:.4f}")
 
# 2. BM25
tokenized_corpus = [doc.split(" ") for doc in corpus]
bm25 = BM25Okapi(tokenized_corpus, k1=1.2, b=0)
 
bm25_scores = bm25.get_scores(query)
 
print("\nBM25 排序结果:")
for idx in np.argsort(-bm25_scores):
    print(f"Doc {idx}: {corpus[idx]} -> Score {bm25_scores[idx]:.4f}")