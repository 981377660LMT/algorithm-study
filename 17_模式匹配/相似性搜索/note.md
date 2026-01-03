# 相似性搜索（Similarity Search）

https://www.bilibili.com/video/BV1a84beFEfS
https://www.bilibili.com/video/BV1N7YMegEaf

### 第一层级：基于字符的算法 (Morphological)

**核心逻辑**：把字符串看作字符序列，关注拼写差异。
**适用场景**：拼写纠错、模糊匹配、短文本（如人名、地址）。

1.  **Levenshtein Distance (编辑距离)**
    - **原理**：计算将 A 变成 B 需要的最少单字符编辑操作次数（插入、删除、替换）。
    - **地位**：衡量字符串差异的“黄金标准”。
2.  **Hamming Distance (海明距离)**
    - **原理**：两个**等长**字符串对应位置不同字符的个数。
    - **场景**：SimHash 指纹比对、数据传输纠错。
3.  **Jaro-Winkler Distance**
    - **原理**：在 Jaro 距离（考虑字符匹配数和换位）基础上，提高了**前缀匹配**的权重。
    - **场景**：人名匹配（因为人名通常姓氏在前或名字在前，前缀很重要）。
4.  **LCS (Longest Common Subsequence)**
    - **原理**：寻找两个序列中最长的公共子序列（不需要连续）。
    - **场景**：代码 Diff 工具（如 Git）、生物序列比对。

---

### 第二层级：基于集合/词块的算法 (Token-based)

**核心逻辑**：把字符串切分为词（Token）或 N-gram 的集合，关注重合度，忽略语序。
**适用场景**：短语匹配、简单文本去重。

1.  **Jaccard Similarity**
    - **原理**：交集大小 / 并集大小。
    - **公式**：$J(A,B) = \frac{|A \cap B|}{|A \cup B|}$
2.  **Sørensen–Dice Coefficient**
    - **原理**：类似 Jaccard，但给交集部分双倍权重。
    - **公式**：$DSC = \frac{2|A \cap B|}{|A| + |B|}$
3.  **Overlap Coefficient**
    - **原理**：判断一个字符串是否是另一个的子集。

---

### 第三层级：基于统计权重的算法 (Statistical / Sparse Vector)

**核心逻辑**：不仅看词是否重合，还看词的**重要性**。这里就是 **TF-IDF 和 BM25** 的领地。它们将文本转换为**稀疏向量**（维度极高，大部分为 0）。
**适用场景**：搜索引擎、关键词检索、文档相关性排序。

1.  **TF-IDF (Term Frequency - Inverse Document Frequency)**
    - **原理**：
      - **TF**：词频。词在当前文档出现越多，越重要。
      - **IDF**：逆文档频率。词在所有文档中出现越少（如生僻词），权重越高；常见词（如 "the", "是"）权重极低。
    - **缺点**：没有考虑词频饱和（词出现 100 次和出现 10 次，权重线性增长，容易被长文主导）。
2.  **BM25 (Best Matching 25)**
    - **地位**：**当前工业界文本检索（如 Elasticsearch, Lucene）的默认标准算法**。
    - **原理**：它是 TF-IDF 的改良版。
      - **词频饱和**：当一个词出现次数达到一定程度后，得分增长会趋于平缓（不再无限增加）。
      - **长度归一化**：惩罚过长的文档，避免长文档仅仅因为词多而得分高。

---

### 第四层级：基于语义的算法 (Semantic / Dense Vector)

**核心逻辑**：将文本映射为**稠密向量**，关注“意思”是否相近，而非“字”是否一样。
**适用场景**：智能问答、推荐系统、跨语言检索。

1.  **Static Embeddings (Word2Vec, GloVe)**
    - **原理**：每个词对应一个固定向量。
    - **缺点**：无法解决多义词（如“苹果”是水果还是手机）。
2.  **Contextual Embeddings (BERT, RoBERTa)**
    - **原理**：基于 Transformer，根据上下文动态生成向量。
    - **SBERT (Sentence-BERT)**：专门优化用于生成句子级别的向量，计算余弦相似度极快。
3.  **Cross-Encoder**
    - **原理**：不生成向量，直接把两个句子一起扔进模型打分。
    - **特点**：精度最高，但速度最慢（适合重排序 Re-ranking）。

---

### 第五层级：基于哈希的速算算法 (Hashing)

**核心逻辑**：将长文本压缩为短指纹，用于海量数据快速去重。

1.  **SimHash**
    - **原理**：基于 LSH（局部敏感哈希），使得相似文本的哈希值海明距离很小。
    - **场景**：Google 网页去重。
2.  **MinHash**
    - **原理**：快速估算 Jaccard 相似度。

---

### 总结与对比图谱

| 算法层级   | 代表算法          | 核心特征       | 优点                   | 缺点           | 典型应用          |
| :--------- | :---------------- | :------------- | :--------------------- | :------------- | :---------------- |
| **字符级** | Levenshtein       | 物理编辑距离   | 精确，容忍拼写错误     | 慢，无语义     | 拼写检查          |
| **集合级** | Jaccard           | 词汇重合度     | 简单快速               | 忽略词频权重   | 简单去重          |
| **统计级** | **TF-IDF / BM25** | **关键词权重** | **可解释性强，抓重点** | **无语义理解** | **搜索引擎 (ES)** |
| **语义级** | BERT / SBERT      | 深度语义理解   | 懂同义词、上下文       | 计算资源大     | 智能客服、RAG     |
| **哈希级** | SimHash           | 指纹摘要       | 极快，内存小           | 精度略低       | 海量网页去重      |

### 代码演示：TF-IDF 与 BM25 (Python)

```python
from sklearn.feature_extraction.text import TfidfVectorizer
from rank_bm25 import BM25Okapi # 需要安装 rank_bm25

corpus = [
    "The quick brown fox jumps over the lazy dog",
    "The quick brown fox",
    "The lazy dog",
    "A fast brown wolf"
]
tokenized_corpus = [doc.split(" ") for doc in corpus]
query = "brown fox"

# --- 1. TF-IDF ---
vectorizer = TfidfVectorizer()
tfidf_matrix = vectorizer.fit_transform(corpus)
# (此处省略余弦相似度计算步骤，通常结合 sklearn 使用)

# --- 2. BM25 ---
bm25 = BM25Okapi(tokenized_corpus)
tokenized_query = query.split(" ")
doc_scores = bm25.get_scores(tokenized_query)

print(f"Query: {query}")
print("BM25 Scores:", doc_scores)
# 结果会显示包含 "brown" 和 "fox" 的文档得分最高
```

---

**Jaro-Winkler Distance** 是一种专门用于衡量两个短字符串（特别是人名）相似度的算法。它是 **Jaro Distance** 的改进版。

它的核心思想是：**不仅考虑字符匹配的数量和顺序，还特别给予“共同前缀”更高的权重**。因为在实际应用中（如匹配人名），如果两个字符串的开头相同，它们是同一个东西的概率就非常大。

---

### 1. 基础：Jaro Distance (Jaro 距离)

在理解 Jaro-Winkler 之前，必须先理解 Jaro 距离。Jaro 距离主要基于两个因素：

1.  **匹配字符数 ($m$)**：两个字符串中相同的字符数量。
2.  **换位数 ($t$)**：虽然字符相同，但位置不对应的数量。

#### 匹配窗口 (Matching Window)

两个字符被认为是“匹配”的，不仅要求字符相同，还要求它们的位置索引差值不超过一个特定距离。
这个距离阈值是：
$$ \text{Window Size} = \left\lfloor \frac{\max(|S_1|, |S_2|)}{2} \right\rfloor - 1 $$

- 如果两个相同字符的距离超过这个范围，就不算匹配。

#### Jaro 公式

$$ d_j = \frac{1}{3} \left( \frac{m}{|S_1|} + \frac{m}{|S_2|} + \frac{m - t}{m} \right) $$

- $|S_1|, |S_2|$：两个字符串的长度。
- $m$：匹配的字符数。
- $t$：换位数的一半（Transpositions / 2）。即需要交换多少次才能让顺序一致。

---

### 2. 进阶：Jaro-Winkler Distance

Jaro-Winkler 在 Jaro 的基础上引入了**前缀缩放（Prefix Scale）**。

#### 核心逻辑

如果两个字符串的开头（前缀）完全相同，那么它们的相似度得分应该在 Jaro 得分的基础上再加分。

#### Jaro-Winkler 公式

$$ d\_{jw} = d_j + (l \cdot p \cdot (1 - d_j)) $$

- $d_j$：原始的 Jaro 距离（0 到 1 之间）。
- $l$：共同前缀的长度（Length of common prefix）。
  - 通常规定最大为 4。如果前缀超过 4 个字符相同，也只算 4。
- $p$：缩放因子（Scaling Factor）。
  - 这是一个常量，通常设为 **0.1**。
  - 这意味着每多一个相同的前缀字符，分数就往上提一点。
  - $l \cdot p$ 最大不能超过 0.25（保证最终得分不超过 1）。

---

### 3. 举例计算 (手动推导)

假设我们要比较：

- $S_1$: **MARTHA**
- $S_2$: **MARHTA**

**第一步：计算 Jaro 距离**

1.  **长度**：$|S_1|=6, |S_2|=6$。
2.  **匹配窗口**：$\lfloor \frac{6}{2} \rfloor - 1 = 2$。
3.  **匹配字符 ($m$)**：
    - M, A, R, T, H, A 全部都在对方的窗口范围内找到对应字符。
    - $m = 6$。
4.  **换位数 ($t$)**：
    - $S_1$: M-A-R-**T-H**-A
    - $S_2$: M-A-R-**H-T**-A
    - T 和 H 的顺序反了。需要交换 1 次。
    - Transpositions = 1，所以公式里的 $t = 1/2 = 0.5$ (注意：有些定义里 $t$ 直接指换位次数的一半，这里为了严谨，通常指非匹配顺序对数的一半)。
    - _修正标准定义_：在 Jaro 中，$t$ 是匹配字符中，顺序不同的字符数除以 2。这里 T 和 H 顺序不同，H 和 T 顺序不同，共 2 个字符顺序不对。所以 $t = 2 / 2 = 1$。

$$ d_j = \frac{1}{3} \left( \frac{6}{6} + \frac{6}{6} + \frac{6 - 1}{6} \right) = \frac{1}{3} (1 + 1 + 0.833) = 0.944 $$

**第二步：计算 Jaro-Winkler 距离**

1.  **共同前缀 ($l$)**：
    - **MAR**THA vs **MAR**HTA
    - 前 3 个字符 "MAR" 完全一样。所以 $l = 3$。
2.  **缩放因子 ($p$)**：取标准值 0.1。

$$ d*{jw} = 0.944 + (3 \cdot 0.1 \cdot (1 - 0.944)) $$
$$ d*{jw} = 0.944 + (0.3 \cdot 0.056) $$
$$ d\_{jw} = 0.944 + 0.0168 = 0.961 $$

**结论**：因为前缀 "MAR" 相同，相似度从 0.944 提升到了 0.961。

---

### 4. 优缺点与适用场景

| 特性         | 说明                                                                                                                                                             |
| :----------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **优点**     | 1. **适合短文本**：特别是人名、地名。<br>2. **前缀敏感**：符合人类直觉（名字写错通常在后面，开头很少错）。<br>3. **归一化**：结果在 0-1 之间，方便作为概率使用。 |
| **缺点**     | 1. **不适合长文本**：对于句子或段落，效果不如 TF-IDF 或 Embedding。<br>2. **计算复杂度**：比简单的海明距离复杂，但比编辑距离快（通常情况下）。                   |
| **最佳场景** | **数据清洗（Data Linkage / Deduplication）**。<br>例如：判断 "Bill Gates" 和 "Bill Gutes" 是否为同一个人。                                                       |

### 5. 代码示例 (Python)

通常使用 `Levenshtein` 库或 `jellyfish` 库来计算。

```python
# 需要安装: pip install jellyfish
import jellyfish

s1 = "MARTHA"
s2 = "MARHTA"

# 1. Jaro Distance
jaro = jellyfish.jaro_similarity(s1, s2)
print(f"Jaro Distance: {jaro:.4f}")

# 2. Jaro-Winkler Distance
jw = jellyfish.jaro_winkler_similarity(s1, s2)
print(f"Jaro-Winkler:  {jw:.4f}")

# 对比一个前缀不同的例子
s3 = "ARTHA" # 少了首字母 M
jw_diff_prefix = jellyfish.jaro_winkler_similarity(s1, s3)
print(f"Prefix Mismatch: {jw_diff_prefix:.4f}")
```

### 总结

Jaro-Winkler 是**记录链接（Record Linkage）**领域的王者算法。当你需要合并两个脏乱差的客户名单（Customer Database），且主要依靠**姓名**字段进行匹配时，Jaro-Winkler 通常比编辑距离（Levenshtein）效果更好。

---

基于哈希的相似性搜索算法旨在解决**海量数据（亿级以上）**下的去重和检索问题。它们的核心是将变长的文本映射为定长的“指纹”（Fingerprint），通过比较指纹的差异来估算原始文本的相似度。

---

### 1. SimHash (局部敏感哈希)

SimHash 是由 Google 设计用于网页去重的算法。它的精髓在于：**原始文本越相似，生成的哈希值（二进制位）相同的越多。**

#### 核心步骤：

1.  **分词与权重**：将文本分词，并给每个词赋予权重（如 TF-IDF 值）。
2.  **哈希化**：将每个词通过普通哈希函数映射为一个 $f$ 位的二进制序列（如 64 位）。
3.  **加权**：遍历哈希序列，如果是 1，则加上该词的权重；如果是 0，则减去该词的权重。
4.  **合并**：将所有词的加权结果按位累加，得到一个长度为 $f$ 的数值序列。
5.  **降维（二值化）**：对累加结果进行“二值化”：大于 0 的位设为 1，小于等于 0 的位设为 0。

#### 判定标准：

使用 **海明距离 (Hamming Distance)**。通常在 64 位指纹中，海明距离 $\le 3$ 即可认为两篇文档高度相似。

---

### 2. MinHash (最小哈希)

MinHash 主要用于快速估算两个集合的 **Jaccard 相似度**。

#### 核心原理：

Jaccard 相似度计算的是交集与并集的比值。对于海量集合，直接计算交集太慢。
MinHash 证明了一个数学定理：**两个集合经过随机置换后，其最小哈希值相等的概率，等于它们的 Jaccard 相似度。**
$$P(h_{min}(A) = h_{min}(B)) = J(A, B)$$

#### 核心步骤：

1.  **特征提取**：将文档转化为 N-gram（Shingling）集合。
2.  **哈希置换**：使用 $K$ 个不同的哈希函数，对集合中的每个元素计算哈希值。
3.  **取最小值**：对于每个哈希函数，记录集合中所有元素产生的最小哈希值。
4.  **生成签名**：$K$ 个最小值组成该文档的 **MinHash Signature**（签名向量）。
5.  **对比**：两个签名向量中相同元素的比例，即为 Jaccard 相似度的近似值。

---

### 3. SimHash vs MinHash 对比

| 特性           | SimHash                        | MinHash                      |
| :------------- | :----------------------------- | :--------------------------- |
| **相似度度量** | 余弦相似度 (Cosine Similarity) | Jaccard 相似度               |
| **关注点**     | 词频与权重 (TF-IDF)            | 集合元素的重合度             |
| **计算方式**   | 位运算 (Hamming Distance)      | 签名向量对比                 |
| **适用场景**   | 网页去重、长文本相似度         | 推荐系统、文档聚类、子集匹配 |

---

### 4. Python 代码示例

```python
# pip install simhash datasketch
from simhash import Simhash
from datasketch import MinHash

text1 = "GitHub Copilot is an AI coding assistant"
text2 = "GitHub Copilot is a great AI coding assistant"

# --- SimHash 示例 ---
hash1 = Simhash(text1)
hash2 = Simhash(text2)
# 计算海明距离 (越小越相似)
print(f"SimHash Hamming Distance: {hash1.distance(hash2)}")

# --- MinHash 示例 ---
m1, m2 = MinHash(num_perm=128), MinHash(num_perm=128)
for word in text1.split():
    m1.update(word.encode('utf8'))
for word in text2.split():
    m2.update(word.encode('utf8'))

# 估算 Jaccard 相似度 (越大越相似)
print(f"MinHash Jaccard Estimate: {m1.jaccard(m2):.4f}")
```

### 总结

- 如果你需要处理**带权重的长文本去重**，选 **SimHash**。
- 如果你需要处理**集合重合度或大规模聚类**，选 **MinHash**。
- 在工业界，MinHash 常配合 **LSH (Locality Sensitive Hashing)** 桶分技术，实现从数亿数据中秒级检索相似项。
