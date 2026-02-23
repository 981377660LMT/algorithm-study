import math
from collections import Counter


class BM25:
    def __init__(self, corpus, k1=1.5, b=0.75):
        self.k1 = k1
        self.b = b
        self.corpus_size = len(corpus)
        self.avgdl = sum(len(doc) for doc in corpus) / self.corpus_size
        self.doc_freqs = []
        self.idf = {}
        self.doc_len = []
        self._initialize(corpus)

    def _initialize(self, corpus):
        nd = {}  # 存储包含词 t 的文档数量
        for doc in corpus:
            self.doc_len.append(len(doc))
            frequencies = Counter(doc)
            self.doc_freqs.append(frequencies)
            for word in frequencies:
                nd[word] = nd.get(word, 0) + 1

        # 计算 IDF
        for word, n_q in nd.items():
            self.idf[word] = math.log((self.corpus_size - n_q + 0.5) / (n_q + 0.5) + 1)

    def get_score(self, query, index):
        score = 0.0
        doc_freq = self.doc_freqs[index]
        d_len = self.doc_len[index]

        for word in query:
            if word not in doc_freq:
                continue

            idf = self.idf.get(word, 0)
            # BM25 核心公式
            tf = doc_freq[word]
            numerator = idf * tf * (self.k1 + 1)
            denominator = tf + self.k1 * (1 - self.b + self.b * d_len / self.avgdl)
            score += numerator / denominator

        return score


corpus = [["hello", "world"], ["hello", "bm25", "is", "cool"], ["world", "is", "windy"]]
query = ["hello", "cool"]

bm25_manual = BM25(corpus)
scores = [bm25_manual.get_score(query, i) for i in range(len(corpus))]
print(f"手动计算文档得分: {scores}")
