## 目录

- [目录](#目录)
- [相似性（Similarity）概述](#相似性similarity概述)
  - [相似性的定义](#相似性的定义)
  - [常见的相似性度量方法](#常见的相似性度量方法)
  - [Jaccard 相似性](#jaccard-相似性)
- [MinHash 概述](#minhash-概述)
  - [MinHash 的原理](#minhash-的原理)
  - [MinHash 的应用场景](#minhash-的应用场景)
- [Go 实现 MinHash](#go-实现-minhash)
  - [实现步骤](#实现步骤)
  - [完整代码示例](#完整代码示例)
  - [代码说明](#代码说明)
- [示例与测试](#示例与测试)
  - [示例文档](#示例文档)
  - [测试结果](#测试结果)
- [总结](#总结)

---

## 相似性（Similarity）概述

### 相似性的定义

在计算机科学和数据分析中，相似性（Similarity）指的是衡量两个对象在某种度量下的相似程度。对象可以是文本、图像、音频或其他任何数据类型。相似性度量广泛应用于信息检索、推荐系统、聚类分析等领域。

### 常见的相似性度量方法

1. **欧几里得距离（Euclidean Distance）**：常用于数值型数据，通过计算两点之间的直线距离来衡量相似性。
2. **余弦相似性（Cosine Similarity）**：主要用于高维稀疏数据，如文本数据，衡量两个向量在空间中的夹角。
3. **Jaccard 相似性（Jaccard Similarity）**：用于集合数据，衡量两个集合的交集与并集的比率。
4. **曼哈顿距离（Manhattan Distance）**：计算两点在各坐标轴上的距离总和。

### Jaccard 相似性

**定义**：对于两个集合 A 和 B，Jaccard 相似性定义为它们交集的大小与并集的大小之比。

$$
J(A, B) = \frac{|A \cap B|}{|A \cup B|}
$$

**特点**：

- 范围在 [0,1] 之间，值越大表示相似度越高。
- 对于稀疏数据（如文本的词集）效果较好。
- 计算成本较高，尤其是对于大规模数据集。

## MinHash 概述

### MinHash 的原理

MinHash（Min-wise Independent Permutations Hashing）是一种用于估计 Jaccard 相似性的近似算法。它通过生成一组哈希函数，将高维稀疏数据映射到低维空间，从而有效地估计两个集合的相似性。

**核心思想**：

1. **哈希函数集合**：选择一组独立的哈希函数，每个函数能够对集合元素进行随机排列。
2. **签名矩阵**：对于每个集合，应用所有哈希函数，记录每个哈希函数在集合中的最小哈希值，形成一个签名向量。
3. **相似性估计**：通过比较两个签名向量中相同哈希值的比例，近似估计它们的 Jaccard 相似性。

**优势**：

- 大幅降低计算和存储成本。
- 支持高效的近似相似性计算。
- 易于并行化，适合大规模数据处理。

### MinHash 的应用场景

- **文档去重**：快速检测相似文档，避免重复内容。
- **推荐系统**：基于用户行为的相似性推荐。
- **聚类分析**：高效地将相似对象聚类。
- **大规模数据处理**：如搜索引擎中的网页相似性检测。

## Go 实现 MinHash

### 实现步骤

1. **选择哈希函数**：选择一组独立的哈希函数，通常可以通过参数化的单一哈希函数实现多组哈希值。
2. **构建集合表示**：将文档或对象表示为集合（如词集、shingles）。
3. **计算签名矩阵**：对于每个集合，使用所有哈希函数计算最小哈希值，形成签名向量。
4. **估计相似性**：通过比较签名向量的相似性，估计集合间的 Jaccard 相似性。

### 完整代码示例

以下是一个使用 Go 实现 MinHash 的完整示例，包括生成哈希函数、计算签名向量以及估计相似性。

```go
package main

import (
	"fmt"
	"hash/fnv"
	"math"
	"math/rand"
	"strings"
	"time"
)

// Constants
const (
	numHashFunctions = 100 // Number of hash functions
	maxHashValue     = math.MaxUint32
)

// MinHash struct
type MinHash struct {
	numHash uint32
	hashA   []uint32
	hashB   []uint32
}

// NewMinHash initializes a MinHash with a given number of hash functions
func NewMinHash(num uint32) *MinHash {
	mh := &MinHash{
		numHash: num,
		hashA:   make([]uint32, num),
		hashB:   make([]uint32, num),
	}

	// Initialize hash function coefficients with random values
	rand.Seed(time.Now().UnixNano())
	for i := uint32(0); i < num; i++ {
		mh.hashA[i] = rand.Uint32() | 1 // Ensure it's odd
		mh.hashB[i] = rand.Uint32()
	}

	return mh
}

// hashToken hashes a token to a uint32 using FNV hash
func hashToken(token string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(token))
	return h.Sum32()
}

// ComputeSignature computes the MinHash signature for a given set of tokens
func (mh *MinHash) ComputeSignature(tokens []string) []uint32 {
	signature := make([]uint32, mh.numHash)
	for i := range signature {
		signature[i] = maxHashValue
	}

	for _, token := range tokens {
		tokenHash := hashToken(token)
		for i := uint32(0); i < mh.numHash; i++ {
			combinedHash := (mh.hashA[i]*tokenHash + mh.hashB[i]) % maxHashValue
			if combinedHash < signature[i] {
				signature[i] = combinedHash
			}
		}
	}

	return signature
}

// EstimateJaccard estimates the Jaccard similarity between two signatures
func EstimateJaccard(sig1, sig2 []uint32) float64 {
	if len(sig1) != len(sig2) {
		return 0.0
	}
	match := 0
	for i := 0; i < len(sig1); i++ {
		if sig1[i] == sig2[i] {
			match++
		}
	}
	return float64(match) / float64(len(sig1))
}

// Helper function to convert text to tokens (e.g., word-level shingles)
func tokenize(text string) []string {
	// Simple whitespace tokenizer; can be replaced with more sophisticated tokenization
	return strings.Fields(text)
}

func main() {
	// Initialize MinHash
	mh := NewMinHash(numHashFunctions)

	// Example documents
	doc1 := "The quick brown fox jumps over the lazy dog"
	doc2 := "The quick brown fox leaps over the lazy dog"
	doc3 := "Lorem ipsum dolor sit amet, consectetur adipiscing elit"

	// Tokenize documents
	tokens1 := tokenize(doc1)
	tokens2 := tokenize(doc2)
	tokens3 := tokenize(doc3)

	// Compute signatures
	sig1 := mh.ComputeSignature(tokens1)
	sig2 := mh.ComputeSignature(tokens2)
	sig3 := mh.ComputeSignature(tokens3)

	// Estimate similarities
	sim12 := EstimateJaccard(sig1, sig2)
	sim13 := EstimateJaccard(sig1, sig3)
	sim23 := EstimateJaccard(sig2, sig3)

	// Output results
	fmt.Printf("Jaccard similarity between Doc1 and Doc2: %.4f\n", sim12)
	fmt.Printf("Jaccard similarity between Doc1 and Doc3: %.4f\n", sim13)
	fmt.Printf("Jaccard similarity between Doc2 and Doc3: %.4f\n", sim23)
}
```

### 代码说明

1. **MinHash 结构体**：

   - `numHash`: 哈希函数的数量。
   - `hashA` 和 `hashB`: 哈希函数的系数，用于生成独立的哈希函数。

2. **NewMinHash**：

   - 初始化 MinHash 实例，生成 `numHash` 个独立的哈希函数系数。
   - 使用随机数生成器生成 `hashA` 和 `hashB`，确保每个哈希函数的独立性。

3. **hashToken**：

   - 使用 FNV-1a 哈希算法将字符串 token 转换为 `uint32` 哈希值。
   - FNV-1a 是一种简单且快速的非加密哈希函数，适合用于此类应用。

4. **ComputeSignature**：

   - 对给定的 token 集合计算 MinHash 签名。
   - 初始化签名向量为最大哈希值。
   - 对每个 token 计算所有哈希函数的哈希值，并更新签名向量中的最小哈希值。

5. **EstimateJaccard**：

   - 通过比较两个签名向量中相同位置的哈希值，估计 Jaccard 相似性。
   - 相同的最小哈希值数量除以哈希函数的总数，得到相似性估计。

6. **tokenize**：

   - 简单的分词函数，将文本转换为 token 列表。
   - 可根据需要更改为更复杂的分词策略（如 n-gram、去除停用词等）。

7. **main 函数**：
   - 示例文档定义与分词。
   - 计算每个文档的 MinHash 签名。
   - 估计并输出文档间的 Jaccard 相似性。

## 示例与测试

### 示例文档

```plaintext
Doc1: "The quick brown fox jumps over the lazy dog"
Doc2: "The quick brown fox leaps over the lazy dog"
Doc3: "Lorem ipsum dolor sit amet, consectetur adipiscing elit"
```

### 测试结果

运行上述代码，可能得到如下输出（具体值可能因哈希函数的随机性略有不同）：

```
Jaccard similarity between Doc1 and Doc2: 0.8500
Jaccard similarity between Doc1 and Doc3: 0.0200
Jaccard similarity between Doc2 and Doc3: 0.0200
```

**解释**：

- **Doc1 和 Doc2** 有较高的相似性（约 85%），因为它们只有一个词的不同（"jumps" vs "leaps"）。
- **Doc1 和 Doc3** 以及 **Doc2 和 Doc3** 的相似性很低（约 2%），因为它们的内容完全不同。

## 总结

MinHash 是一种高效的近似算法，用于估计大规模集合间的 Jaccard 相似性。通过使用多个独立的哈希函数，MinHash 能够在保持较低计算和存储成本的同时，提供相对准确的相似性估计。在实际应用中，MinHash 广泛应用于文档去重、推荐系统、聚类分析等领域。

上述 Go 代码示例展示了如何实现 MinHash，包括哈希函数的生成、签名矩阵的计算以及相似性的估计。该实现可以根据具体需求进行扩展和优化，例如使用更复杂的哈希函数、优化存储结构或并行计算等。
