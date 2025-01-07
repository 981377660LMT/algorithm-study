这段代码实现了一个支持模糊查询（基于 Levenshtein 距离）的 Trie 树（字典树），旨在用于自动补全功能。相比于普通的前缀搜索，模糊查询允许返回与查询关键词在编辑距离（插入、删除、替换操作次数）范围内的结果。以下是对该代码的详细分析和讲解。

## 目录

- [目录](#目录)
- [高层次概述](#高层次概述)
- [主要数据结构](#主要数据结构)
  - [Trie](#trie)
  - [Node](#node)
  - [DoublyLinkedList 和 dllNode](#doublylinkedlist-和-dllnode)
- [Trie 的基本操作](#trie-的基本操作)
  - [插入 (Put)](#插入-put)
  - [删除 (Delete)](#删除-delete)
- [搜索功能](#搜索功能)
  - [搜索选项 (SearchOptions)](#搜索选项-searchoptions)
  - [搜索结果](#搜索结果)
  - [Levenshtein 距离的实现](#levenshtein-距离的实现)
    - [主要步骤：](#主要步骤)
    - [关键函数：](#关键函数)
  - [前缀搜索](#前缀搜索)
- [辅助功能](#辅助功能)
  - [堆 (Heap)](#堆-heap)
  - [打印 Trie 树](#打印-trie-树)
- [主函数示例](#主函数示例)
- [总结](#总结)
- [Levenshtein 距离搜索实现](#levenshtein-距离搜索实现)
  - [详细分析 Levenshtein 距离搜索实现](#详细分析levenshtein距离搜索实现)
- [1. Levenshtein 距离与 Trie 树](#1-levenshtein距离与trie树)
  - [1.1 什么是 Levenshtein 距离](#11-什么是levenshtein距离)
  - [1.2 Trie 树中的模糊搜索](#12-trie树中的模糊搜索)
- [2. 代码中的 Levenshtein 距离搜索实现](#2-代码中的levenshtein距离搜索实现)
  - [2.1 搜索接口与选项配置](#21-搜索接口与选项配置)
  - [2.2 Levenshtein 距离的计算方法](#22-levenshtein距离的计算方法)
    - [2.2.1 动态规划矩阵](#221-动态规划矩阵)
  - [2.3 `searchWithEditDistance` 与 `buildWithEditDistance` 函数解析](#23-searchwitheditdistance-与-buildwitheditdistance-函数解析)
    - [2.3.1 `searchWithEditDistance`](#231-searchwitheditdistance)
    - [2.3.2 `buildWithEditDistance`](#232-buildwitheditdistance)
    - [2.3.3 `getEditOps` 函数](#233-geteditops-函数)
  - [2.4 最大堆（Max Heap）用于 Top-K 结果维护](#24-最大堆max-heap用于top-k结果维护)
    - [2.4.1 定义](#241-定义)
    - [2.4.2 使用方式](#242-使用方式)
  - [2.5 优化策略](#25-优化策略)
    - [2.5.1 剪枝](#251-剪枝)
    - [2.5.2 优先遍历匹配节点](#252-优先遍历匹配节点)
    - [2.5.3 滚动编辑距离矩阵](#253-滚动编辑距离矩阵)
    - [2.5.4 使用双向链表维护子节点的插入顺序](#254-使用双向链表维护子节点的插入顺序)
  - [2.6 复杂度分析](#26-复杂度分析)
  - [2.7 代码示例分析](#27-代码示例分析)
- [3. 总结](#3-总结)
  - [可能的改进方向](#可能的改进方向)
  - [结语](#结语)

---

## 高层次概述

Trie（字典树）是一种多叉树结构，常用于存储和检索键（通常是字符串）集合。本实现的 Trie 支持以下功能：

- **插入和删除**：支持键（由字符串切片组成）和值的插入和删除。
- **前缀搜索**：基于键的前缀进行搜索。
- **模糊搜索**：基于 Levenshtein 距离（编辑距离）的模糊搜索，允许返回与查询键在一定编辑距离范围内的结果。
- **结果排序**：搜索结果按插入顺序确定，支持返回编辑距离最小的前 K 个结果。
- **编辑操作记录**：在模糊搜索中，记录将 Trie 中键转换为查询键所需的具体编辑操作。

此外，还包括辅助函数用于打印 Trie 树的结构，以便调试和展示。

---

## 主要数据结构

### Trie

```go
type Trie struct {
    root *Node
}
```

Trie 结构体包含一个根节点，代表 Trie 树的起点。根节点的`keyPart`通常为"^"表示根。

### Node

```go
type Node struct {
    keyPart     string
    isTerminal  bool
    value       interface{}
    dllNode     *dllNode
    children    map[string]*Node
    childrenDLL *doublyLinkedList
}
```

Node 结构体代表 Trie 树中的每个节点，包含以下字段：

- `keyPart`：当前节点代表的键部分（字符串）。
- `isTerminal`：布尔值，指示是否是一个完整键的终结节点。
- `value`：存储与键关联的值。
- `dllNode`：指向在兄弟节点双向链表中的节点，用于保持子节点的插入顺序。
- `children`：映射，键为子节点的`keyPart`，值为对应的子节点。
- `childrenDLL`：双向链表，按插入顺序存储子节点，用于有序遍历。

### DoublyLinkedList 和 dllNode

```go
type doublyLinkedList struct {
    head, tail *dllNode
}

type dllNode struct {
    trieNode   *Node
    next, prev *dllNode
}
```

双向链表用于维护子节点的插入顺序，确保搜索结果的确定性。`doublyLinkedList`有`head`和`tail`指针，`dllNode`代表链表中的每个节点，包含指向对应`Trie`节点的指针以及前后节点的指针。

---

## Trie 的基本操作

### 插入 (Put)

```go
func (t *Trie) Put(key []string, value interface{}) (existed bool) {
    node := t.root
    for i, part := range key {
        child, ok := node.children[part]
        if (!ok) {
            child = newNode(part)
            child.dllNode = newDLLNode(child)
            node.children[part] = child
            node.childrenDLL.append(child.dllNode)
        }
        if i == len(key)-1 {
            existed = child.isTerminal
            child.isTerminal = true
            child.value = value
        }
        node = child
    }
    return existed
}
```

- **功能**：向 Trie 中插入一个键（由字符串切片组成）和对应的值。如果键已经存在，则更新其值，并返回`true`；否则返回`false`。
- **实现细节**：
  - 从根节点开始，逐步遍历键的每一部分。
  - 如果当前键部分的子节点不存在，则创建一个新节点，并将其添加到`children`映射和`childrenDLL`双向链表中。
  - 如果到达键的末尾，将节点标记为终结节点（`isTerminal`为`true`），并设置其值。

### 删除 (Delete)

```go
func (t *Trie) Delete(key []string) (value interface{}, existed bool) {
    node := t.root
    parent := make(map[*Node]*Node)
    for _, keyPart := range key {
        child, ok := node.children[keyPart]
        if !ok {
            return nil, false
        }
        parent[child] = node
        node = child
    }
    if (!node.isTerminal) {
        return nil, false
    }
    node.isTerminal = false
    value = node.value
    node.value = nil
    for node != nil && !node.isTerminal && len(node.children) == 0 {
        delete(parent[node].children, node.keyPart)
        parent[node].childrenDLL.pop(node.dllNode)
        node = parent[node]
    }
    return value, true
}
```

- **功能**：从 Trie 中删除一个键，并返回其对应的值。如果键不存在，返回`false`。
- **实现细节**：
  - 首先，遍历键的每一部分，记录每个节点的父节点。
  - 如果找到对应的终结节点，取消其终结节点标记，并移除其值。
  - 从后向前，删除那些不再是终结节点且没有子节点的节点，以保持 Trie 的紧凑性。

---

## 搜索功能

Trie 的搜索功能支持多种选项，包括前缀搜索和基于编辑距离的模糊搜索。

### 搜索选项 (SearchOptions)

```go
type SearchOptions struct {
    exactKey        bool
    maxResults      bool
    maxResultsCount int
    editDistance    bool
    maxEditDistance int
    editOps         bool
    topKLeastEdited bool
}
```

- **字段解释**：
  - `exactKey`：是否仅返回完全匹配的结果。
  - `maxResults`：是否限制返回结果的最大数量。
  - `maxEditDistance`：最大允许的编辑距离，用于模糊搜索。
  - `editOps`：是否返回编辑操作（插入、删除、替换信息）。
  - `topKLeastEdited`：是否仅返回编辑距离最小的前 K 个结果。

### 搜索结果

```go
type SearchResults struct {
    Results         []*SearchResult
    heap            *searchResultMaxHeap
    tiebreakerCount int
}

type SearchResult struct {
    Key          []string
    Value        interface{}
    EditDistance int
    EditOps      []*EditOp
    tiebreaker   int
}
```

- **SearchResult**：
  - `Key`：匹配的键。
  - `Value`：对应的值。
  - `EditDistance`：键与查询键的编辑距离。
  - `EditOps`：将键转换为查询键所需的具体编辑操作。
  - `tiebreaker`：用于在堆中解决编辑距离相同的情况，保证插入顺序。

### Levenshtein 距离的实现

Levenshtein 距离计算文本之间的最小编辑操作次数。该实现通过动态规划在 Trie 中高效计算编辑距离，结合 Trie 的结构剪枝不必要的路径。

#### 主要步骤：

1. **初始化**：

   - 构建初始的编辑距离矩阵，其中第一行`newRow`代表从空字符串到查询键的距离。

2. **递归遍历 Trie**：

   - 优先遍历与查询键匹配的子节点，以提高搜索效率和结果的相关性。
   - 对每个节点，计算当前路径的编辑距离，并根据编辑距离是否在允许范围内决定是否将结果加入到搜索结果中。

3. **剪枝**：

   - 如果当前路径的最小编辑距离已经超过了最大允许距离，则停止进一步遍历该路径。

4. **记录编辑操作**：
   - 如果启用了`WithEditOps`选项，则在找到匹配结果时，回溯编辑距离矩阵，生成具体的编辑操作列表。

#### 关键函数：

- `searchWithEditDistance`：主函数，负责初始化并遍历 Trie，收集符合条件的搜索结果。
- `buildWithEditDistance`：递归函数，遍历 Trie 并计算编辑距离，收集结果。
- `getEditOps`：根据编辑距离矩阵回溯生成编辑操作列表。

```go
func (t *Trie) Search(key []string, options ...func(*SearchOptions)) *SearchResults {
    // 省略选项解析和校验
    if opts.editDistance {
        return t.searchWithEditDistance(key, opts)
    }
    return t.search(key, opts)
}
```

### 前缀搜索

在没有启用编辑距离选项时，搜索功能执行标准的前缀搜索，即返回所有以查询键为前缀的键。

```go
func (t *Trie) search(prefixKey []string, opts *SearchOptions) *SearchResults {
    results := &SearchResults{}
    node := t.root
    for _, keyPart := range prefixKey {
        child, ok := node.children[keyPart]
        if !ok {
            return results
        }
        node = child
    }
    if opts.exactKey {
        if node.isTerminal {
            result := &SearchResult{Key: prefixKey, Value: node.value}
            results.Results = append(results.Results, result)
        }
        return results
    }
    t.build(results, node, &prefixKey, opts)
    return results
}
```

- **功能**：定位到前缀对应的节点，然后递归遍历其所有子节点，收集符合条件的结果。
- **选项支持**：支持`exactKey`（仅返回完全匹配的结果）和`maxResults`等选项。

---

## 辅助功能

### 堆 (Heap)

用于在模糊搜索中维护一个最大堆，以便高效地获取编辑距离最小的前 K 个结果。

```go
type searchResultMaxHeap []*SearchResult

func (s searchResultMaxHeap) Len() int { return len(s) }

func (s searchResultMaxHeap) Less(i, j int) bool {
    if s[i].EditDistance == s[j].EditDistance {
        return s[i].tiebreaker > s[j].tiebreaker
    }
    return s[i].EditDistance > s[j].EditDistance
}

func (s searchResultMaxHeap) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s *searchResultMaxHeap) Push(x interface{}) {
    *s = append(*s, x.(*SearchResult))
}

func (s *searchResultMaxHeap) Pop() interface{} {
    old := *s
    n := len(old)
    x := old[n-1]
    *s = old[0 : n-1]
    return x
}
```

- **特性**：
  - `Less`方法定义了堆的顺序，首先按编辑距离降序排列；如果编辑距离相同，则按插入顺序（`tiebreaker`）降序排列，以便在达到最大堆容量时，能够替换掉编辑距离更大的结果。
  - 使用`container/heap`包提供的堆操作接口，实现了标准的堆逻辑。

### 打印 Trie 树

提供了多种方式来打印 Trie 树的结构，便于调试和可视化。

- **主要接口**：

  - `Print`：直接打印到标准输出。
  - `Sprint`：返回格式化后的字符串。
  - `PrintWithError`和`SprintWithError`：类似但能够处理错误（如存在重复节点）。
  - **其他格式**：
    - `PrintHr` 和 `SprintHr`：水平格式化打印。
    - `PrintHrn` 和 `SprintHrn`：带有水平换行的格式化打印。

- **实现细节**：
  - 使用深度优先或广度优先遍历，结合字符图形（如`├─`，`└─`，`│`）来表示树的层次和分支。
  - 处理节点的宽度以对齐子节点，确保输出的美观性。

```go
func printEditOps(ops []*EditOp) {
    for _, op := range ops {
        switch op.Type {
        case EditOpTypeNoEdit:
            fmt.Printf("- don't edit %q\n", op.KeyPart)
        case EditOpTypeInsert:
            fmt.Printf("- insert %q\n", op.KeyPart)
        case EditOpTypeDelete:
            fmt.Printf("- delete %q\n", op.KeyPart)
        case EditOpTypeReplace:
            fmt.Printf("- replace %q with %q\n", op.KeyPart, op.ReplaceWith)
        }
    }
}
```

- **示例**：将编辑操作列表打印为易读的文本描述，如“替换某个部分”或“删除某个部分”。

---

## 主函数示例

```go
func main() {
    tri := NewTrieFuzzy()
    // 插入键和值
    tri.Put([]string{"the"}, 1)
    tri.Put([]string{"the", "quick", "brown", "fox"}, 2)
    tri.Put([]string{"the", "quick", "sports", "car"}, 3)
    tri.Put([]string{"the", "green", "tree"}, 4)
    tri.Put([]string{"an", "apple", "tree"}, 5)
    tri.Put([]string{"an", "umbrella"}, 6)

    tri.Root().Print()
    // 打印Trie树的结构

    // 前缀搜索
    results := tri.Search([]string{"the", "quick"})
    for _, res := range results.Results {
        fmt.Println(res.Key, res.Value)
    }
    // 输出匹配前缀为["the", "quick"]的所有结果

    // 模糊搜索
    key := []string{"the", "tree"}
    results = tri.Search(key, WithMaxEditDistance(2), WithEditOps())
    for _, res := range results.Results {
        fmt.Println(res.Key, res.EditDistance)
    }
    // 输出与["the", "tree"]编辑距离不超过2的结果

    // 打印编辑操作
    result := results.Results[2]
    fmt.Printf("To convert %v to %v:\n", result.Key, key)
    printEditOps(result.EditOps)

    // 组合选项：返回编辑距离最小的前2个结果
    results = tri.Search(key, WithMaxEditDistance(2), WithTopKLeastEdited(), WithMaxResults(2))
    for _, res := range results.Results {
        fmt.Println(res.Key, res.Value, res.EditDistance)
    }
}
```

- **功能**：

  - 插入一些示例键和值。
  - 打印整个 Trie 树的结构。
  - 执行前缀搜索，输出匹配的结果。
  - 执行模糊搜索，允许编辑距离为 2，并输出相应的编辑距离和编辑操作。
  - 执行组合选项的搜索，获取编辑距离最小的前 2 个结果。

- **输出示例**：
  - **Trie 树结构**：
    ```
    ^
    ├─ the ($)
    │  ├─ quick
    │  │  ├─ brown
    │  │  │  └─ fox ($)
    │  │  └─ sports
    │  │     └─ car ($)
    │  └─ green
    │     └─ tree ($)
    └─ an
       ├─ apple
       │  └─ tree ($)
       └─ umbrella ($)
    ```
  - **前缀搜索结果**：
    ```
    [the quick brown fox] 2
    [the quick sports car] 3
    ```
  - **模糊搜索结果**：
    ```
    [the] 1
    [the green tree] 1
    [an apple tree] 2
    [an umbrella] 2
    ```
  - **编辑操作示例**：
    ```
    To convert [an apple tree] to [the tree]:
    - delete "an"
    - replace "apple" with "the"
    - don't edit "tree"
    ```
  - **组合选项的搜索结果**：
    ```
    [the] 1 1
    [the green tree] 4 1
    ```

---

## 总结

这段代码实现了一个功能丰富的 Trie 树，支持标准的前缀搜索和基于编辑距离的模糊搜索，适用于自动补全等应用场景。主要特点包括：

- **插入和删除**：高效地管理键值对，确保 Trie 结构的紧凑性。
- **模糊搜索**：基于 Levenshtein 距离的模糊搜索，允许在一定编辑距离范围内返回匹配结果，增强了搜索的灵活性。
- **排序和结果控制**：支持按插入顺序排序，限制返回结果的数量，以及优先返回编辑距离最小的结果。
- **编辑操作记录**：在模糊搜索中记录具体的编辑操作，提供更丰富的输出信息。
- **可视化工具**：提供多种打印方式，帮助开发者直观地了解 Trie 的结构，便于调试和优化。

整体而言，这段代码提供了一个强大且灵活的 Trie 实现，适用于需要高效搜索和自动补全的应用场景。

---

## Levenshtein 距离搜索实现

### 详细分析 Levenshtein 距离搜索实现

在这段代码中，实现了一个支持模糊查询（基于 Levenshtein 距离）的 Trie 树。Levenshtein 距离，也称为编辑距离，是衡量两个字符串之间差异的指标，表示将一个字符串转换成另一个字符串所需的最少编辑操作次数（插入、删除、替换）。在 Trie 树中结合 Levenshtein 距离，可以高效地实现模糊搜索，允许用户输入近似关键词时仍能获得相关结果。

以下是对 Levenshtein 距离搜索实现的详尽分析：

---

## 1. Levenshtein 距离与 Trie 树

### 1.1 什么是 Levenshtein 距离

Levenshtein 距离是计算将一个字符串转换为另一个字符串所需的最少的单字符编辑操作数，这些操作包括：

- **插入**一个字符
- **删除**一个字符
- **替换**一个字符

例如，将单词 "kitten" 转换为 "sitting" 的 Levenshtein 距离是 3：

1. 替换 'k' 为 's'。
2. 替换 'e' 为 'i'。
3. 插入 'g'。

### 1.2 Trie 树中的模糊搜索

Trie 树是一种高效的多叉树数据结构，主要用于存储和检索字符串。传统的 Trie 树支持前缀搜索，即快速查找以某个前缀开头的所有字符串。然而，前缀搜索要求精确匹配前缀，这在实际应用中可能不够灵活。

为了增强 Trie 树的搜索能力，引入了基于 Levenshtein 距离的模糊搜索。这允许用户输入近似词汇时，仍能获得相关的搜索结果。例如，输入 "helo" 也能匹配 "hello"。

## 2. 代码中的 Levenshtein 距离搜索实现

本文将重点分析以下几个关键部分：

1. **搜索接口与选项配置**
2. **Levenshtein 距离计算**
3. **Trie 遍历与距离计算**
4. **编辑操作的记录**
5. **结果排序与限制**

### 2.1 搜索接口与选项配置

```go
type SearchOptions struct {
    exactKey        bool
    maxResults      bool
    maxResultsCount int
    editDistance    bool
    maxEditDistance int
    editOps         bool
    topKLeastEdited bool
}
```

`SearchOptions` 结构体定义了搜索的各种可选参数，通过函数式选项模式进行配置。这些选项允许调用者自定义搜索行为，例如是否启用模糊搜索、设置最大编辑距离、是否返回编辑操作列表、以及结果的排序和限制等。

相关的选项配置函数包括：

- `WithExactKey()`
- `WithMaxResults(maxResults int)`
- `WithMaxEditDistance(maxDistance int)`
- `WithEditOps()`
- `WithTopKLeastEdited()`

### 2.2 Levenshtein 距离的计算方法

Levenshtein 距离通常通过动态规划算法实现，计算两个字符串之间的编辑距离。对于 Trie 树中的模糊搜索，算法需要在 Trie 的遍历过程中动态地更新距离矩阵，以确定当前路径是否可接受。

在这段代码中，结合了动态规划和 Trie 的 DFS 遍历来高效计算 Levenshtein 距离。

#### 2.2.1 动态规划矩阵

为了计算 Levenshtein 距离，通常使用一个二维矩阵 `dp[i][j]` 表示将字符串 A 的前 `i` 个字符转换为字符串 B 的前 `j` 个字符所需的最少编辑操作数。

然而，在 Trie 的模糊搜索中，字符串 A 是树中从根到某个节点的路径，而字符串 B 是用户的查询关键词。为了节省空间和提高效率，代码中采用了一种滚动数组优化，只维护当前行和前一行的编辑距离。

### 2.3 `searchWithEditDistance` 与 `buildWithEditDistance` 函数解析

#### 2.3.1 `searchWithEditDistance`

```go
func (t *Trie) searchWithEditDistance(key []string, opts *SearchOptions) *SearchResults {
    // 初始化编辑距离矩阵的第一行
    columns := len(key) + 1
    newRow := make([]int, columns)
    for i := 0; i < columns; i++ {
        newRow[i] = i
    }
    m := len(key)
    if m == 0 {
        m = 1
    }
    rows := make([][]int, 1, m)
    rows[0] = newRow
    results := &SearchResults{}
    if opts.topKLeastEdited {
        results.heap = &searchResultMaxHeap{}
    }

    // 优先遍历与查询键首字母匹配的节点
    keyColumn := make([]string, 1, m)
    stop := false
    var prioritizedNode *Node
    if len(key) > 0 {
        if prioritizedNode = t.root.children[key[0]]; prioritizedNode != nil {
            keyColumn[0] = prioritizedNode.keyPart
            t.buildWithEditDistance(&stop, results, prioritizedNode, &keyColumn, &rows, key, opts)
        }
    }
    // 遍历其他子节点
    for dllNode := t.root.childrenDLL.head; dllNode != nil; dllNode = dllNode.next {
        node := dllNode.trieNode
        if node == prioritizedNode {
            continue
        }
        keyColumn[0] = node.keyPart
        t.buildWithEditDistance(&stop, results, node, &keyColumn, &rows, key, opts)
    }
    if opts.topKLeastEdited {
        n := results.heap.Len()
        results.Results = make([]*SearchResult, n)
        for n != 0 {
            result := heap.Pop(results.heap).(*SearchResult)
            result.tiebreaker = 0
            results.Results[n-1] = result
            n--
        }
        results.heap = nil
        results.tiebreakerCount = 0
    }
    return results
}
```

**功能**：

- 初始化编辑距离矩阵的第一行，表示将空字符串转换为查询键所需的编辑操作数（即插入操作数）。
- 使用一个变量 `stop` 控制遍历的提前终止，当达到最大编辑距离限制时，可以停止进一步的遍历以节省计算资源。
- 优先遍历与查询键首字母匹配的节点，以提高搜索效率和结果的相关性。
- 如果启用了 `topKLeastEdited`，使用最大堆（`heap`）来维护前 K 个编辑距离最小的搜索结果。

**关键点**：

- **滚动矩阵**：仅保存了当前行和前一行，节省了空间。
- **优先遍历**：首先遍历与查询键首字母匹配的子节点，可以更快地找到更相关的结果。
- **最大堆**：用于维护顶 K 个最小编辑距离的结果。

#### 2.3.2 `buildWithEditDistance`

```go
func (t *Trie) buildWithEditDistance(stop *bool, results *SearchResults, node *Node, keyColumn *[]string, rows *[][]int, key []string, opts *SearchOptions) {
    if *stop {
        return
    }
    prevRow := (*rows)[len(*rows)-1]
    columns := len(key) + 1
    newRow := make([]int, columns)
    newRow[0] = prevRow[0] + 1
    for i := 1; i < columns; i++ {
        replaceCost := 1
        if key[i-1] == (*keyColumn)[len(*keyColumn)-1] {
            replaceCost = 0
        }
        newRow[i] = mins(
            newRow[i-1]+1,            // 插入
            prevRow[i]+1,             // 删除
            prevRow[i-1]+replaceCost, // 替换
        )
    }
    *rows = append(*rows, newRow)

    if newRow[columns-1] <= opts.maxEditDistance && node.isTerminal {
        editDistance := newRow[columns-1]
        lazyCreate := func() *SearchResult {
            resultKey := make([]string, len(*keyColumn))
            copy(resultKey, *keyColumn)
            result := &SearchResult{Key: resultKey, Value: node.value, EditDistance: editDistance}
            if opts.editOps {
                result.EditOps = t.getEditOps(rows, keyColumn, key)
            }
            return result
        }
        if opts.topKLeastEdited {
            results.tiebreakerCount++
            if results.heap.Len() < opts.maxResultsCount {
                result := lazyCreate()
                result.tiebreaker = results.tiebreakerCount
                heap.Push(results.heap, result)
            } else if (*results.heap)[0].EditDistance > editDistance {
                result := lazyCreate()
                result.tiebreaker = results.tiebreakerCount
                heap.Pop(results.heap)
                heap.Push(results.heap, result)
            }
        } else {
            result := lazyCreate()
            results.Results = append(results.Results, result)
            if opts.maxResults && len(results.Results) == opts.maxResultsCount {
                *stop = true
                return
            }
        }
    }

    if mins(newRow...) <= opts.maxEditDistance {
        var prioritizedNode *Node
        m := len(*keyColumn)
        if m < len(key) {
            if prioritizedNode = node.children[key[m]]; prioritizedNode != nil {
                *keyColumn = append(*keyColumn, prioritizedNode.keyPart)
                t.buildWithEditDistance(stop, results, prioritizedNode, keyColumn, rows, key, opts)
                *keyColumn = (*keyColumn)[:len(*keyColumn)-1]
            }
        }
        for dllNode := node.childrenDLL.head; dllNode != nil; dllNode = dllNode.next {
            child := dllNode.trieNode
            if child == prioritizedNode {
                continue
            }
            *keyColumn = append(*keyColumn, child.keyPart)
            t.buildWithEditDistance(stop, results, child, keyColumn, rows, key, opts)
            *keyColumn = (*keyColumn)[:len(*keyColumn)-1]
        }
    }

    *rows = (*rows)[:len(*rows)-1]
}
```

**功能**：

- 递归地遍历 Trie，计算当前节点对应路径的编辑距离。
- 根据编辑距离是否在允许范围内，将结果添加到搜索结果中。
- 使用剪枝策略，减少不必要的遍历，提高搜索效率。
- 记录编辑操作列表（启用 `WithEditOps` 时）。

**关键步骤**：

1. **计算当前行的编辑距离**：

   - `newRow[0] = prevRow[0] + 1`：表示从父路径到当前添加一个键部分需要的编辑操作数，即删除操作。
   - 对于每个查询键位置 `i`，计算插入、删除和替换的成本，选择最小的成本作为 `newRow[i]`。
   - 如果当前键部分与查询键的第 `i-1` 个部分匹配，则替换成本为 0。

2. **检查并记录结果**：

   - 如果当前路径的最后一个字符的编辑距离 `newRow[columns-1]` 小于或等于最大允许编辑距离，且当前节点是终结节点，则将该结果记录下来。
   - 如果启用了 `WithEditOps`，则调用 `getEditOps` 函数记录具体的编辑操作。

3. **使用最大堆维护 Top-K 结果**（如果启用 `WithTopKLeastEdited`）：

   - 如果堆中元素数量小于 K，直接加入堆。
   - 如果堆已满且当前编辑距离小于堆顶的最大编辑距离，则替换堆顶元素。

4. **剪枝**：

   - 如果当前路径最小的编辑距离仍然小于或等于最大允许编辑距离，则继续遍历子节点。
   - 否则，停止遍历该子树，以节省计算资源。

5. **优先遍历与查询键匹配的子节点**：

   - 这有助于更快地找到相关性更高的结果，提高搜索效率和结果质量。

6. **回溯编辑距离矩阵**：
   - 在递归调用结束后，弹出当前行，回溯到父路径。

#### 2.3.3 `getEditOps` 函数

```go
func (t *Trie) getEditOps(rows *[][]int, keyColumn *[]string, key []string) []*EditOp {
    // 回溯编辑距离矩阵，生成编辑操作列表
    ops := make([]*EditOp, 0, len(key))
    r, c := len(*rows)-1, len((*rows)[0])-1
    for r > 0 || c > 0 {
        insertionCost, deletionCost, substitutionCost := math.MaxInt, math.MaxInt, math.MaxInt
        if c > 0 {
            insertionCost = (*rows)[r][c-1]
        }
        if r > 0 {
            deletionCost = (*rows)[r-1][c]
        }
        if r > 0 && c > 0 {
            substitutionCost = (*rows)[r-1][c-1]
        }
        minCost := mins(insertionCost, deletionCost, substitutionCost)
        if minCost == substitutionCost {
            if (*rows)[r][c] > (*rows)[r-1][c-1] {
                ops = append(ops, &EditOp{Type: EditOpTypeReplace, KeyPart: (*keyColumn)[r-1], ReplaceWith: key[c-1]})
            } else {
                ops = append(ops, &EditOp{Type: EditOpTypeNoEdit, KeyPart: (*keyColumn)[r-1]})
            }
            r -= 1
            c -= 1
        } else if minCost == deletionCost {
            ops = append(ops, &EditOp{Type: EditOpTypeDelete, KeyPart: (*keyColumn)[r-1]})
            r -= 1
        } else if minCost == insertionCost {
            ops = append(ops, &EditOp{Type: EditOpTypeInsert, KeyPart: key[c-1]})
            c -= 1
        }
    }
    for i, j := 0, len(ops)-1; i < j; i, j = i+1, j-1 {
        ops[i], ops[j] = ops[j], ops[i]
    }
    return ops
}
```

**功能**：

- 根据编辑距离矩阵，回溯生成将 Trie 中的键转换为查询键所需的编辑操作列表。
- 编辑操作包括：替换、删除、插入和不编辑。

**关键步骤**：

1. **初始化**：

   - `r` 和 `c` 分别指向编辑距离矩阵的最后一个位置（行和列）。
   - `ops` 用于存储生成的编辑操作。

2. **回溯过程**：

   - **替换/无编辑**：
     - 如果当前位置的编辑距离大于左上角对角线的位置，表示发生了替换操作；否则，为无编辑操作。
     - 操作被添加到 `ops` 列表中，并同时 `r` 和 `c` 递减。
   - **删除**：
     - 如果当前位置的编辑距离来自上方，则是删除操作，操作被添加到 `ops` 中，并 `r` 递减。
   - **插入**：
     - 如果当前位置的编辑距离来自左侧，则是插入操作，操作被添加到 `ops` 中，并 `c` 递减。

3. **反转操作列表**：
   - 因为回溯是从后往前进行的，最终需要将操作列表反转，以确保编辑操作按从头到尾的顺序排列。

**示例输出**：

假设将关键词 `[an apple tree]` 转换为查询键 `[the tree]`，可能的编辑操作列表为：

1. 删除 `"an"`
2. 替换 `"apple"` 为 `"the"`
3. 不编辑 `"tree"`

---

### 2.4 最大堆（Max Heap）用于 Top-K 结果维护

在进行模糊搜索时，尤其是需要返回编辑距离最小的前 K 个结果时，使用堆（Heap）可以高效地维护这些结果。

#### 2.4.1 定义

```go
type searchResultMaxHeap []*SearchResult

func (s searchResultMaxHeap) Len() int { ... }
func (s searchResultMaxHeap) Less(i, j int) bool { ... }
func (s searchResultMaxHeap) Swap(i, j int) { ... }
func (s *searchResultMaxHeap) Push(x interface{}) { ... }
func (s *searchResultMaxHeap) Pop() interface{} { ... }
```

**功能**：

- 定义了一个最大堆 `searchResultMaxHeap`，用于存储搜索结果。
- `Less` 方法定义了堆的排序规则，优先级按 `EditDistance` 降序排列；如果编辑距离相同，则按 `tiebreaker` 降序排列。这确保了堆顶部始终是当前最大的编辑距离。
- 当堆的大小超过 K 时，可以将堆顶（即最大的编辑距离）弹出，以维持只保存前 K 个最小编辑距离的结果。

**关键点**：

- **最大堆结构**：通过最大堆，可以在保证每次插入操作仅 O(log K)时间复杂度的情况下，维护前 K 个最小编辑距离的结果。
- **tiebreaker**：用于在编辑距离相同的情况下，有明确的优先级顺序（通常按插入顺序）。

#### 2.4.2 使用方式

在`buildWithEditDistance`函数中，当启用了`WithTopKLeastEdited`选项时，维护一个最大堆来保存前 K 个最小编辑距离的结果：

1. 如果堆中元素少于 K，则直接将当前结果加入堆。
2. 如果堆已满且当前编辑距离小于堆顶的编辑距离，则弹出堆顶元素，并将当前结果加入堆。

完成遍历后，从堆中弹出所有元素，得到排序好的前 K 个最小编辑距离的结果。

---

### 2.5 优化策略

为了提高模糊搜索的效率并减少不必要的计算，代码中采用了以下优化策略：

#### 2.5.1 剪枝

在遍历 Trie 时，如果当前路径的最小编辑距离已经超过了最大允许编辑距离，则可以停止对该子树的进一步遍历。这显著降低了搜索空间，减少了不必要的计算。

```go
if mins(newRow...) <= opts.maxEditDistance {
    // 继续遍历子节点
} else {
    // 剪枝，不再遍历该子树
}
```

#### 2.5.2 优先遍历匹配节点

通过优先遍历与查询键当前部分匹配的子节点，能够更快地找到相关性更高的搜索结果。这不仅提高了搜索效率，还提高了搜索结果的质量。

```go
// 优先遍历与查询键首字母匹配的节点
if len(key) > 0 {
    if prioritizedNode = t.root.children[key[0]]; prioritizedNode != nil {
        keyColumn[0] = prioritizedNode.keyPart
        t.buildWithEditDistance(&stop, results, prioritizedNode, &keyColumn, &rows, key, opts)
    }
}
// 遍历其他子节点
```

#### 2.5.3 滚动编辑距离矩阵

通过仅保存当前行和前一行的编辑距离矩阵，减少了内存的使用量，提高了计算效率。

```go
rows := make([][]int, 1, m)
rows[0] = newRow
// 在递归调用结束后，回溯编辑距离矩阵
*rows = (*rows)[:len(*rows)-1]
```

#### 2.5.4 使用双向链表维护子节点的插入顺序

通过使用双向链表（`doublyLinkedList`）维护子节点的插入顺序，确保搜索结果的确定性和一致性。

---

### 2.6 复杂度分析

**时间复杂度**：

- **最坏情况**：当不进行剪枝时，搜索需要遍历整个 Trie，每个节点都需要计算编辑距离。对于一个最坏情况的 Trie，时间复杂度接近于 O(N \* M)，其中 N 是 Trie 中节点的总数，M 是查询键的长度。
- **优化后的情况**：通过剪枝和优先遍历，实际的时间复杂度通常远低于最坏情况，尤其是当 Trie 较大但查询键的相关路径有限时。

**空间复杂度**：

- **编辑距离矩阵**：由于采用了滚动矩阵技术，空间复杂度为 O(M)，其中 M 是查询键的长度。
- **递归堆栈**：最坏情况下，递归深度等于 Trie 的高度。对于平衡的 Trie，递归深度约为 log(N)，其中 N 是 Trie 中节点的总数。

- **结果存储**：如果启用了`WithTopKLeastEdited`，最大空间需求为 O(K)，其中 K 是前 K 个结果的数量。

---

### 2.7 代码示例分析

让我们通过代码中的主函数示例更深入地理解 Levenshtein 距离搜索的实现。

```go
func main() {
    tri := NewTrieFuzzy()
    // 插入键和值
    tri.Put([]string{"the"}, 1)
    tri.Put([]string{"the", "quick", "brown", "fox"}, 2)
    tri.Put([]string{"the", "quick", "sports", "car"}, 3)
    tri.Put([]string{"the", "green", "tree"}, 4)
    tri.Put([]string{"an", "apple", "tree"}, 5)
    tri.Put([]string{"an", "umbrella"}, 6)

    tri.Root().Print()
    // 打印Trie树的结构

    // 前缀搜索
    results := tri.Search([]string{"the", "quick"})
    for _, res := range results.Results {
        fmt.Println(res.Key, res.Value)
    }
    // 输出匹配前缀为["the", "quick"]的所有结果

    // 模糊搜索
    key := []string{"the", "tree"}
    results = tri.Search(key, WithMaxEditDistance(2), WithEditOps())
    for _, res := range results.Results {
        fmt.Println(res.Key, res.EditDistance)
    }
    // 输出与["the", "tree"]编辑距离不超过2的结果

    // 打印编辑操作
    result := results.Results[2]
    fmt.Printf("To convert %v to %v:\n", result.Key, key)
    printEditOps(result.EditOps)

    // 组合选项：返回编辑距离最小的前2个结果
    results = tri.Search(key, WithMaxEditDistance(2), WithTopKLeastEdited(), WithMaxResults(2))
    for _, res := range results.Results {
        fmt.Println(res.Key, res.Value, res.EditDistance)
    }
}
```

**执行流程**：

1. **Trie 树构建**：

   - 插入了六个键值对，如 `"the quick brown fox"` 对应值 `2`。
   - Trie 树结构如下：
     ```
     ^
     ├─ the ($)
     │  ├─ quick
     │  │  ├─ brown
     │  │  │  └─ fox ($)
     │  │  └─ sports
     │  │     └─ car ($)
     │  └─ green
     │     └─ tree ($)
     └─ an
        ├─ apple
        │  └─ tree ($)
        └─ umbrella ($)
     ```

2. **前缀搜索**：

   - 查询键前缀为 `["the", "quick"]`，应返回：
     - `[the quick brown fox] 2`
     - `[the quick sports car] 3`

3. **模糊搜索**：

   - 查询键为 `["the", "tree"]`，允许最大编辑距离为 `2`，并记录编辑操作。
   - 预期的结果包括：
     - `[the] 1`
     - `[the green tree] 1`
     - `[an apple tree] 2`
     - `[an umbrella] 2`
   - 具体到 `[an apple tree]` 转换为 `[the tree]` 的编辑操作：
     ```
     To convert [an apple tree] to [the tree]:
     - delete "an"
     - replace "apple" with "the"
     - don't edit "tree"
     ```

4. **组合选项的模糊搜索**：
   - 查询键为 `["the", "tree"]`，允许最大编辑距离为 `2`，要求返回编辑距离最小的前 `2` 个结果。
   - 预期输出：
     - `[the] 1 1`
     - `[the green tree] 4 1`

**关键分析**：

- **匹配与编辑距离**：

  - `[the]` 只需删除 `"the tree"` 中的 `"tree"`，编辑距离为 `1`。
  - `[the green tree]` 需要删除 `"green"`，编辑距离也为 `1`。
  - `[an apple tree]` 需要删除 `"an"` 和替换 `"apple"` 为 `"the"`，编辑距离为 `2`。
  - `[an umbrella]` 需要删除 `"an"` 和 `"umbrella"` 替换为 `"the"` 和 `"tree"`，编辑距离为 `2`。

- **编辑操作记录**：

  - 通过`getEditOps`函数，能够详细记录每个编辑操作，提供更丰富的信息。

- **结果排序与限制**：
  - 当启用了 `WithTopKLeastEdited` 和 `WithMaxResults(2)` 后，只返回编辑距离最小的前 `2` 个结果，即 `[the]` 和 `[the green tree]`。

---

## 3. 总结

这段代码通过结合 Trie 树和动态规划的 Levenshtein 距离计算，成功实现了一个高效的模糊搜索功能。以下是主要亮点：

1. **高效的编辑距离计算**：

   - 使用滚动矩阵优化空间复杂度，仅维护当前和前一行的编辑距离。
   - 通过剪枝和优先遍历优化时间复杂度，避免不必要的计算。

2. **灵活的搜索选项**：

   - 允许调用者自定义搜索行为，如设定最大编辑距离、是否记录编辑操作、限制返回结果数量等。

3. **编辑操作记录**：

   - 提供详细的编辑操作列表，帮助理解键之间的转换过程。

4. **结果排序与有限制**：

   - 支持返回编辑距离最小的前 K 个结果，通过最大堆实现高效的结果维护。

5. **确定性的结果顺序**：

   - 使用双向链表维护子节点的插入顺序，确保搜索结果的顺序一致性。

6. **可扩展性和维护性**：
   - 通过模块化设计和清晰的函数划分，代码易于维护和扩展。

### 可能的改进方向

尽管目前的实现已经相当完善，但仍有一些潜在的优化和改进方向：

1. **更高效的编辑距离算法**：

   - 使用更先进的编辑距离算法，如加权 Levenshtein 距离或 Damerau-Levenshtein 距离，更好地适应具体应用场景。

2. **多线程/并行化支持**：

   - 对 Trie 遍历和编辑距离计算进行并行化处理，进一步提高大规模数据集上的搜索性能。

3. **缓存优化**：

   - 对经常访问的子树或查询结果进行缓存，减少重复计算。

4. **动态调整最大编辑距离**：
   - 根据用户输入的实时需求动态调整允许的最大编辑距离，以更好地平衡搜索质量和性能。

### 结语

通过本次详细分析，可以看出该 Trie 树模糊搜索的 Levenshtein 距离实现结合了多种优化策略，既保证了搜索结果的准确性和丰富性，又在性能上进行了有效的优化。这使得该数据结构非常适合用于自动补全、拼写纠正等需要高效模糊匹配的应用场景。
