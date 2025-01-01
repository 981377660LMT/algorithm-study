下面这份代码来自 [go-adaptive-radix-tree](https://github.com/plar/go-adaptive-radix-tree)，它实现了自适应基数树（ART, Adaptive Radix Tree）这一高效的前缀树（Radix Tree）数据结构。阅读这段代码可能会觉得很长、复杂，但它实际上遵循了一个非常清晰的设计思路。下面会分模块逐步解释，包括：

1. **基本概念与常量定义**
2. **接口设计（Node、Tree、Iterator）**
3. **核心结构：nodeRef、leaf、nodeX（node4、node16、node48、node256）**
4. **树的增删查遍历逻辑**
5. **辅助函数与调试工具**

通过这些部分，可以帮助你理解整个 ART 的实现原理和代码结构。

---

## 1. 基本概念与常量定义

### Node types

```go
const (
    Leaf    Kind = 0
    Node4   Kind = 1
    Node16  Kind = 2
    Node48  Kind = 3
    Node256 Kind = 4
)
```

- **Leaf**：叶子节点，用来存储真正的 key-value。
- **Node4 / Node16 / Node48 / Node256**：内部节点，分别表示有最多 4、16、48、256 个分支的节点。

自适应基数树的核心思想是：

- 当节点分支数很少时（例如 <=4），就使用 Node4 这种“小容量”节点，减少内存浪费；
- 如果分支多了，就“扩容”成更大容量的节点 (Node16, Node48, Node256)；
- 如果分支又被大量删除，也可以“缩容”。

### Traverse Options

```go
const (
    TraverseLeaf    = 1
    TraverseNode    = 2
    TraverseAll     = TraverseLeaf | TraverseNode
    TraverseReverse = 4
)
```

- **TraverseLeaf**：遍历时只关注叶子节点。
- **TraverseNode**：遍历时只关注内部节点。
- **TraverseAll**：遍历时两者都要。
- **TraverseReverse**：标记需要逆序遍历。

这些选项会在树的遍历 `ForEach`、`Iterator` 中控制是正序 / 倒序、遍历哪些节点等。

### 错误定义

```go
var (
    ErrConcurrentModification = errors.New("concurrent modification has been detected")
    ErrNoMoreNodes            = errors.New("there are no more nodes in the tree")
)
```

- 在迭代器遍历的过程中，如果检测到底层树发生了结构性修改，就会抛出 `ErrConcurrentModification`。
- 如果迭代到尽头了，则会抛出 `ErrNoMoreNodes`。

---

## 2. 接口设计（Node、Tree、Iterator）

### Node 接口

```go
type Node interface {
    Kind() Kind
    Key() Key
    Value() Value
}
```

- **Kind()**：节点类型（Leaf、Node4、Node16、Node48、Node256）。
- **Key()**、**Value()**：只有在 `Kind == Leaf` 时才有实际意义；内部节点调用返回 nil。

### Iterator 接口

```go
type Iterator interface {
    HasNext() bool
    Next() (Node, error)
}
```

- **HasNext()**：是否还有下一个节点可遍历。
- **Next()**：返回下一个节点。如果没有更多，返回 `ErrNoMoreNodes`；如果树结构在遍历期间被改了，返回 `ErrConcurrentModification`。

### Tree 接口

```go
type Tree interface {
    Insert(key Key, value Value) (oldValue Value, updated bool)
    Delete(key Key) (value Value, deleted bool)
    Search(key Key) (value Value, found bool)
    ForEach(callback Callback, options ...int)
    ForEachPrefix(keyPrefix Key, callback Callback, options ...int)
    Iterator(options ...int) Iterator
    Minimum() (Value, bool)
    Maximum() (Value, bool)
    Size() int
}
```

- **Insert/Delete/Search**：对应增删查操作。
- **ForEach/ForEachPrefix**：遍历；前者遍历所有节点或指定类型节点（叶/内部节点），后者只遍历指定前缀的节点。
- **Iterator**：返回一个迭代器，用于在树上顺序或逆序地遍历。
- **Minimum/Maximum**：获取树中最小/最大 key 对应的叶子节点。
- **Size()**：当前存储的键值对数量。

### `NewTree`

```go
func NewTree() Tree {
    return newTree()
}
```

- `NewTree()` 是对外的构造函数，内部调用 `newTree()` 返回私有实现 `*tree`。

---

## 3. 核心结构

### 3.1 nodeRef

```go
type nodeRef struct {
    ref  unsafe.Pointer
    kind Kind
}
```

- `nodeRef` 是最核心的“指针包装结构”，它对所有类型的节点（leaf 或 node4/16/48/256）做了统一封装：
  - `ref` 存储了指向不同节点结构的 `unsafe.Pointer`；
  - `kind` 表示节点类型，用于在需要时进行类型转换。

#### 常见方法

- **`isLeaf()`**：是否是叶子。
- **`leaf()`**、`nodeX()`：将 `unsafe.Pointer` 转换成相应的具体类型 `(*leaf)`, `(*node4)`, `(*node16)` 等，方便调用特定方法。

> 设计上之所以要这样做，是为了在增删查找时，只需要对外暴露 `*nodeRef`，而不需要每次做显式的类型断言。通过 `nr.kind` 判断即可。

### 3.2 leaf

```go
type leaf struct {
    key   Key
    value interface{}
}
```

- 叶子节点结构，里面直接存储用户的 `key` 和 `value`。
- **`match(key)`**：判断是否和目标 key 完全一致。
- **`prefixMatch(key)`**：判断是否以某个前缀开头（前缀查询时用）。

### 3.3 node（内部节点的公共结构）

```go
type node struct {
    prefix      prefix // prefix是一个定长数组 [maxPrefixLen]byte
    prefixLen   uint16
    childrenLen uint16
}
```

- 自适应基数树内部节点常见的字段：
  - `prefix`：存储了该节点处的一段“公共前缀”（为了在搜索时能快速跳过前缀匹配）。
  - `prefixLen`：实际使用的前缀长度（有时不一定用满 `prefix` 数组）。
  - `childrenLen`：子节点的实际数量。

不同类型的节点（`node4`, `node16`, `node48`, `node256`）都“内嵌”了这个 `node` 结构，从而拥有共同的前缀、子节点计数等。

### 3.4 node4 / node16 / node48 / node256

它们都**内嵌**了 `node`，再根据需要有不同的容器来存储子节点引用。这里的 4、16、48、256 指的是**最多有多少个分支**。

- **node4**：

  ```go
  type node4 struct {
      node
      children [node4Max + 1]*nodeRef // 最多4个分支 + 第5个存储 “key为0的子节点”（invalid keyChar）
      keys     [node4Max]byte
      present  [node4Max]byte
  }
  ```

  - `childrenLen <= 4`。
  - `keys` 和 `present` 共同表示每个子节点的索引和存在状态。（`present[i]` 表示第 i 个位置是不是在用。）

- **node16**：

  ```go
  type node16 struct {
      node
      children [node16Max + 1]*nodeRef
      keys     [node16Max]byte
      present  present16
  }
  ```

  - `childrenLen <= 16`。
  - 用了一个 `present16 uint16` 的位图来标记子节点是否存在。

- **node48**：

  ```go
  type node48 struct {
      node
      children [node48Max + 1]*nodeRef
      keys     [node256Max]byte
      present  present48
  }
  ```

  - `childrenLen <= 48`。
  - 这里 `keys[ch]` 存储 “字符 ch 对应的 child 在 children 数组中是第几号下标”。
  - `present48` 则是一个 4×64=256 位的位图，用来标识 0~255 每个字符是否存在分支。

- **node256**：
  ```go
  type node256 struct {
      node
      children [node256Max + 1]*nodeRef
  }
  ```
  - `childrenLen <= 256`，直接用一个 256 长度的数组存储每个字符的子节点指针，没有额外的位图。
  - 当节点已经塞不下了就不会再扩容（256 是上限）。

#### Grow / Shrink

比如从 `node4` 长到 `node16`，再长到 `node48`，再长到 `node256`，或者反过来 shrink。这是自适应基数树最重要的地方：它保证在节点分支很少时不浪费空间，在分支较多时能保持查询效率。

- `grow()`:  
  如果当前节点满了，在插入一个孩子前，先调用 `grow()` 变成更大容量节点，再把这个新 child 加进去。
- `shrink()`:  
  如果一个节点的子节点数量小到一定阈值（如 `node16` 少于 5 个子节点，就可以 shrink 成 `node4`），则转换成更紧凑的数据结构。

---

## 4. 树的增删查遍历逻辑

### 4.1 整体树结构

```go
type tree struct {
    version int      // 每次结构修改(version++)，用于检测并发修改
    size    int      // 树中存储的 key-value 总数
    root    *nodeRef // 树的根节点(可能是leaf, 也可能是node4/16/48/256)
}
```

### 4.2 Insert

核心流程在 `Insert(key, value)` → `insertRecursively(...)`:

```go
func (tr *tree) Insert(key Key, value Value) (Value, bool) {
    oldVal, status := tr.insertRecursively(&tr.root, key, value, 0)
    if status == treeOpInserted {
        tr.version++  // 修改 version
        tr.size++
    }
    return oldVal, status == treeOpUpdated
}
```

- `insertRecursively(&tr.root, key, value, 0)`：从根节点开始，用递归插入。
  - 如果根节点是 `nil`，就直接新建一个叶子。
  - 如果根节点是叶子，则比较 key：
    - 如果相同，更新值即可；
    - 如果不同，说明需要**分裂**成一个内部节点（通常是 node4），同时把旧 leaf 和新 leaf 都挂在这个 node4 上。
  - 如果根节点是内部节点：
    - 先和节点上的 prefix 对比，看是否需要分裂（prefix 有不匹配时就要 `splitNode`）。
    - 如果 prefix 完全匹配，就继续往下找子节点（character 索引），递归插入。
    - 如果找不到子节点，就新建一个 leaf，挂在这个节点上。

#### 分裂 leaf

```go
func (tr *tree) splitLeaf(...)
```

- 当一个 leaf 遇到插入一个不同的 key，但之前 prefix 又相同，就会转化为 `node4`，把旧 leaf、新 leaf 分别挂到这个 node4 的两个分支上。

#### 分裂 nodeX

```go
func (tr *tree) splitNode(...)
```

- 当 node 前缀与 key 一部分匹配，但在某个字符处不匹配，就需要在此处分裂成更上一层的 node4。

### 4.3 Delete

同理，删除操作 `Delete(key)` → `deleteRecursively(&tr.root, key, 0)`：

1. 如果 root 是叶子并且 key 匹配，删除（置 `*nrp = nil`）。
2. 如果是内部节点，先校验 prefix，若不匹配则无事发生；若匹配，则继续往下找到相应分支的 child。
3. 删除 child 里的叶子后，如果节点分支数小于“最小阈值”，就 shrink。
4. shrink 可能把一个 node4 转成 leaf，或者 node16 转成 node4，等等。

### 4.4 Search

```go
func (tr *tree) Search(key Key) (Value, bool) {
    keyOffset := 0
    current := tr.root
    for current != nil {
        if current.isLeaf() {
            leaf := current.leaf()
            if leaf.match(key) {
                return leaf.value, true
            }
            return nil, false
        }

        // 比对 prefix
        curNode := current.node()
        prefixLen := current.match(key, keyOffset)
        ...
        keyOffset += ...

        // 找到子节点
        next := current.findChildByKey(key, keyOffset)
        if *next == nil {
            return nil, false
        }
        current = *next
        keyOffset++
    }
    return nil, false
}
```

- 一层一层地匹配 `prefix` 并找子节点，直到找到叶子 or 中途无法匹配。

### 4.5 遍历（ForEach / ForEachPrefix / Iterator）

- `ForEach(callback, options...)`：从根节点开始，递归调用 `forEachRecursively`。内部再根据 `TraverseLeaf` / `TraverseNode` / `TraverseAll` 和是否 `TraverseReverse` 来决定遍历策略。
- `ForEachPrefix(keyPrefix, callback, options...)`：带前缀过滤，只回调满足该前缀的叶子。
- `Iterator`：通过一个栈 `state.items` 记录当前节点和它的子节点遍历状态，每次 `Next()` 都会弹出下一个需要访问的节点。如果检测到 `tree.version != it.version`，则说明并发修改，抛出 `ErrConcurrentModification`。

---

## 5. 辅助函数与调试工具

### 5.1 DumpNode / TreeStringer

代码最后有一大段以 `DumpNode` / `TreeStringer` 结尾，它其实是一个**可视化打印** ART 的工具，生成类似于下图的树结构：

```
─── Node4 (0xc00011c2d0)
    prefix(0): [··········] [0 0 0 0 0 0 0 0 0 0]
    ...
    ├── Leaf (0xc00011c2a0)
    │   key(1): [A] [65]
    │   val: A
    │
    └── Node4 (0xc00011c300)
        prefix(0): ...
        ...
```

它通过维护一个 `nodeRegistry` 来给每个节点分配一个唯一 ID，另外还能打印内存地址、前缀信息等，以便在调试时更直观地了解树形结构。

### 5.2 工具函数

- **`minInt(a, b) int`**：简单的取最小值。
- **`copyNode(dst, src)`**：把节点的基本信息（prefixLen 等）从 src 复制到 dst。
- **`findLongestCommonPrefix`**：在 split 时需要用到，用于找出两个 key 的公共前缀长度。
- **`ternary(condition, ifTrue, ifFalse)`**：一个泛型三元运算符。

---

## 小结

整个自适应基数树（ART）的实现要点：

1. **nodeRef + kind**：用一个统一的 `nodeRef` 封装各种节点类型，配合 `kind` 来做类型判断、强转，可以显著减少到处做类型断言的复杂度。
2. **leaf**：真正存储 key-value 的地方；prefix 存不下时就会分裂为内部节点。
3. **内部节点 node4/16/48/256**：通过不同大小的分支存储方式，加速增删查效率并节省内存。
4. **Grow/Shrink**：在插入/删除时动态调整节点类型，使得节点分支数量和节点容量“相匹配”，提升性能。
5. **前缀 prefix**：存储在内部节点中，用来优化搜索时的比较性能（一次比较多个字符）。
6. **version**：用来检测迭代器遍历过程中的并发修改。

这套代码结构比较全面、系统。ART 作为一棵前缀树，能在字符串/字节序列等场景下提供很好的查询和插入性能，且在有大量公共前缀时节省空间。本实现也支持对不同层次的节点进行可视化打印，方便调试和理解。

如果你第一次接触这类自适应基数树实现，不必一下子啃完全部细节，先掌握**叶子节点 + 各种 nodeX 的增删查逻辑**即可，再回头看代码里是如何在不同类型节点间做动态扩容/缩容，就能逐步理解 ART 的设计思想和实现精髓。
