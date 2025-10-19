# 谈谈我对 ROSA 的理解：后缀自动机如何优化长文 LLM（含后缀自动机科普）

https://zhuanlan.zhihu.com/p/1962665268015788269

这篇文章的核心论点是：**Attention 机制中的 Softmax 检索过程，可以被看作一种“语义上的字符串匹配”，因此可以使用为字符串匹配设计的高效数据结构——后缀自动机（Suffix Automaton, SAM）——来进行优化，从而将长文本 LLM 推理中每个 token 的生成复杂度从 O(N) 降低到摊销 O(1)。**

下面我们将分步解析：

### 1. 后缀自动机 (SAM) 科普与 SuffixAutomaton.go 代码对应

文章首先科普了 SAM。我们可以将这些概念直接映射到您的 Go 代码中，以便更具体地理解。

- **核心概念**: SAM 是一个有向无环图（DAG），每个节点代表一个或多个子串的集合（end-pos 等价类），每条边代表一个字符的转移。

  - **节点/状态 (Node/State)**: 在代码中对应 `type Node struct`。

    ```go
    type Node struct {
        Next   [SIGMA]int32 // SAM 转移边 (文章中的状态转移)
        Link   int32        // 后缀链接 (文章中的 fail 指针)
        MaxLen int32        // 当前节点对应的最长子串的长度
        End    int32        // ...
    }
    ```

  - **Fail 指针 (后缀链接)**: 在代码中是 `Link` 字段。它的作用是，当在一个状态无法通过某个字符`c`转移时，可以跳到 `Link` 指向的节点。这个新节点代表了当前字符串集合的**最长后缀**所在的等价类。这正是文章中提到的“跳转到‘T+c 最长的后缀’”的关键。
  - **Fail 指针树 (Parent Tree)**: 代码中的 `BuildTree()` 函数构建了这棵树。它将所有节点的 `Link` 关系反向（从父亲指向儿子），形成一棵以初始状态为根的树。

    ```go
    func (sam *SuffixAutomaton) BuildTree() [][]int32 {
        // ...
        for v := int32(1); v < n; v++ {
            p := sam.Nodes[v].Link
            graph[p] = append(graph[p], v)
        }
        // ...
    }
    ```

  - **子树大小 = 出现次数**: 文章提到“u 的子树大小即 su 在 S 中出现的次数”。这在代码中由 `GetEndPosSize()` 实现。它通过在 Fail 树上自底向上（利用`dfsOrder`）累加，计算出每个状态代表的子串在原字符串中出现的总次数。

    ```go
    func (sam *SuffixAutomaton) GetEndPosSize(dfsOrder []int32) []int32 {
        // ...
        for i := size - 1; i >= 1; i-- {
            // ...
            pre := sam.Nodes[cur].Link
            endPosSize[pre] += endPosSize[cur] // 累加子树大小
        }
        // ...
    }
    ```

  - **摊销 O(1)的滑动窗口式查询**: 文章描述的“删掉前方字符，增加后方字符”的操作，在您的代码中由 `Move()` 和 `MoveLeft()` 函数实现，它们模拟了在 SAM 上的高效移动。

### 2. Attention 与字符串匹配的关系 (ROSA 的核心思想)

这是文章最具创新性的部分。作者提出了一个思想实验来类比：

1.  **量化**: 将 LLM 中高维的 Token Embedding 量化成离散的“字符”。例如，将词汇表中的每个 token 视为一个唯一的字符。这样，一个 token 序列就变成了一个字符串。
2.  **匹配即 Attention**: `softmax(q*k)` 的计算是在寻找与当前查询 token `q` **语义相关**的上下文 token `k`。作者将其类比为：在量化后的“字符串”中，寻找与 `q` 代表的“字符”或“短语”相匹配的所有历史位置。一个高的 Attention Score 意味着一次成功的“语义匹配”。
3.  **用 SAM 替代 Softmax**:
    - **传统 Attention**: 为了生成第 `n+1` 个 token，需要用它的 `q` 和前 `n` 个 token 的 `k` 全部计算一遍相似度，复杂度为 O(n)。
    - **ROSA (SAM) 方案**:
      - 维护一个关于前 `n` 个 token（即 KV Cache）构建的 SAM。
      - 当前 LLM 的“状态”对应 SAM 上的一个节点 `u`，这个节点 `u` 及其在 Fail 树上的子树，隐式地包含了所有历史匹配信息。
      - 当需要预测下一个 token 时，不再进行全局 Softmax。而是将不同的候选 token（作为“字符”）在 SAM 上从节点 `u` 进行转移。
      - 如果转移成功，说明这个“语义模式”在历史上出现过，可以利用新状态的信息来生成 token。这个转移操作是 O(1) 的。
      - 生成 token 后，将这个新 token 加入到 SAM 中，更新 SAM。这个 `Add()` 操作也是摊销 O(1) 的。

通过这种方式，ROSA 将昂贵的 O(n)全局检索，替换为了高效的 O(1)状态转移和更新，从而解决了长文本的性能瓶颈。

### 3. 对文章中几个问题的思考

- **如何逼近原始 Softmax？**: 文章提到 SAM 只能做 0/1 的硬匹配，而 Attention 是加权的。作者推测，SAM 节点 `u` 对应的 embedding 可以被设计为聚合其 Fail 树子树信息的加权结果，而不仅仅是简单的存在性判断。更新 SAM 时融合新 token 的 embedding，也类似于 RNN 的更新方式。
- **并行化挑战**: SAM 的增量构造 `Add()` 是串行的，这给训练时的并行化带来了巨大挑战。文章也指出了这是当前实现可能低效的原因。
- **qs=ks 的简化**: 在自回归模型中，Query 序列和 Key 序列来自同一个文本，所以 `qs=ks` 是一个自然而然的简化。这意味着模型在处理自己的历史输出时，是在进行“自我匹配”。

### 总结

这篇文章提供了一个非常新颖和深刻的视角：将看似无关的**后缀自动机**和**Transformer 的 Attention 机制**联系起来。它精准地指出了 Attention 的核心是一种“检索”操作，并创造性地提出用 SAM 这种为检索而生的数据结构去替代它。

您的 SuffixAutomaton.go 文件是这个思想的底层基石。ROSA 的构想正是建立在 `Add`, `Move`, `BuildTree`, `GetEndPosSize` 等函数所提供的强大而高效的字符串处理能力之上的。虽然将这个算法思想完美地融入 LLM 并解决所有工程挑战（如并行化、梯度流）还有很长的路要走，但它无疑为解决长文本问题开辟了一条极具潜力的道路。

---

以下是一个代码示例，演示如何使用您提供的 SuffixAutomaton.go 中的结构来解决一个经典问题：**计算一个模式串 (pattern) 在一个文本串 (text) 中作为子串出现的次数**。

这个例子直观地展示了后缀自动机如何高效地处理字符串匹配问题，这也是理解 ROSA 思想的基础。

### 示例：查找子串出现次数

这个程序将：

1.  根据文本串 `text` 构建后缀自动机。
2.  计算每个状态（节点）的 `endPos` 集合大小，这代表了该状态对应的子串在 `text` 中的出现次数。
3.  在构建好的自动机上匹配模式串 `pattern`。
4.  如果匹配成功，最终到达的状态的 `endPos` 大小就是 `pattern` 的出现次数。

```go
package main

import (
	"fmt"
)

// 假设 SuffixAutomaton 及其相关结构和方法已在当前包中定义。
// (代码来源于您提供的 SuffixAutomaton.go 文件)

func main() {
	// 1. 定义文本串和模式串
	text := "aababa"
	pattern := "aba"

	// 2. 为文本串构建后缀自动机
	sam := NewSuffixAutomaton()
	for _, char := range text {
		sam.Add(char)
	}

	// 3. 计算每个状态的 endPos 集合大小
	//    这需要先按 MaxLen 对节点排序（拓扑排序的替代）
	dfsOrder := sam.GetDfsOrder()
	//    然后在 parent tree 上自底向上累加，得到每个状态的出现次数
	endPosSize := sam.GetEndPosSize(dfsOrder)

	// 4. 在自动机上匹配模式串
	pos := int32(0) // 从初始状态开始
	matchSuccess := true
	for _, char := range pattern {
		c := char - OFFSET
		nextNode := sam.Nodes[pos].Next[c]
		if nextNode == -1 {
			// 无法转移，说明模式串不是文本串的子串
			matchSuccess = false
			break
		}
		pos = nextNode
	}

	// 5. 输出结果
	if matchSuccess {
		// 匹配成功，最终状态 pos 的 endPosSize 就是出现次数
		count := endPosSize[pos]
		fmt.Printf("模式串 '%s' 在文本串 '%s' 中出现了 %d 次。\n", pattern, text, count)
	} else {
		fmt.Printf("模式串 '%s' 未在文本串 '%s' 中出现。\n", pattern, text)
	}

	// 另一个例子
	pattern2 := "ac"
	pos = 0
	matchSuccess = true
	for _, char := range pattern2 {
		c := char - OFFSET
		nextNode := sam.Nodes[pos].Next[c]
		if nextNode == -1 {
			matchSuccess = false
			break
		}
		pos = nextNode
	}

	if matchSuccess {
		count := endPosSize[pos]
		fmt.Printf("模式串 '%s' 在文本串 '%s' 中出现了 %d 次。\n", pattern2, text, count)
	} else {
		fmt.Printf("模式串 '%s' 未在文本串 '%s' 中出现。\n", pattern2, text)
	}
}
```

#### 如何运行

您可以将此 `main` 函数替换掉 SuffixAutomaton.go 文件中现有的 `main` 函数，然后直接运行该文件，即可看到输出：

```text
模式串 'aba' 在文本串 'aababa' 中出现了 2 次。
模式串 'ac' 未在文本串 'aababa' 中出现。
```
