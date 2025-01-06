https://github.com/siongui/go-succinct-data-structure-trie?tab=readme-ov-file
https://zhuanlan.zhihu.com/p/38385054

---

下面是一份较为详细的代码解读，涵盖了从构建原始 Trie、将其压缩成 Succinct Trie、再到使用 RankDirectory 和 FrozenTrie 进行快速查询的完整流程。代码结构可大致分为以下几部分：

1. **main 函数（示例流程）**
2. **Trie 和 TrieNode（构建原始 Trie）**
3. **字母表及字符到整数的映射（alphabet）**
4. **Base64 编码相关（base64）**
5. **BitString（位串的读操作）**
6. **BitWriter（位串的写操作）**
7. **FrozenTrie 及其节点（解码并快速查询）**
8. **RankDirectory（Rank 和 Select 的快速查询结构）**
9. **搜索函数（GetSuggestedWords）**

下面将从主流程到各个数据结构和函数实现，进行拆解说明。

---

## 1. main 函数（示例流程）

```go
func main() {
    insertNotInAlphabeticalOrder := func(te *Trie) {
        te.Insert("apple")
        te.Insert("orange")
        te.Insert("alphapha")
        te.Insert("lamp")
        te.Insert("hello")
        te.Insert("jello")
        te.Insert("quiz")
    }

    // 1) 构建原始的 Trie
    te := Trie{}
    te.Init()

    // 2) 插入单词（无须按字母顺序，虽然按顺序会更快）
    insertNotInAlphabeticalOrder(&te)

    // 3) 将 Trie 编码为 Succinct Trie
    teData := te.Encode()
    println(teData)                    // 输出紧凑编码后的字符串
    println(te.GetNodeCount())         // 输出节点数

    // 4) 为快速 Rank/Select 查询，构建 RankDirectory
    rd := CreateRankDirectory(teData, te.GetNodeCount()*2+1, L1, L2)
    println(rd.GetData())             // 输出构建的 rank directory

    // 5) 构建 FrozenTrie，用于后续的解码和快速查询
    ft := FrozenTrie{}
    ft.Init(teData, rd.GetData(), te.GetNodeCount())

    // 6) 测试查找单词
    println(ft.Lookup("apple"))  // true
    println(ft.Lookup("appl"))   // false（因为 “appl” 并未标记为终止单词）
    println(ft.Lookup("applee")) // false （不存在）

    // 7) 前缀搜索：找出以 “a” 开头的单词，最多返回10个
    for _, word := range ft.GetSuggestedWords("a", 10) {
        println(word)
    }
}
```

主流程中做了以下几件事情：

1. **构建原始 Trie** 并插入若干单词。
2. 通过 `Encode()` 方法，把原始 Trie 转为 “Succinct” 的 bitstring 字符串 `teData`。
3. 基于这段字符串，用 `CreateRankDirectory` 构建一个 `RankDirectory`，供后续快速 Rank/Select 查询。
4. 用 `FrozenTrie.Init` 将字符串和 RankDirectory 数据结合，形成可用于查找的不可变、紧凑的 Trie：**FrozenTrie**。
5. 测试单词查询（`Lookup`）和前缀建议（`GetSuggestedWords`）。

---

## 2. Trie 和 TrieNode（构建原始 Trie）

### 2.1 TrieNode

```go
type TrieNode struct {
    letter   string      // 当前节点存储的字符
    final    bool        // 标识该节点是否对应一个完整单词
    children []*TrieNode // 子节点列表
}
```

- `letter`: 表示节点中存储的字符，比如 `'a'` / `'b'` / `' '` 等。
- `final`: 若为 `true`，代表从根到此节点的路径构成的单词是一个完整的单词。
- `children`: 子节点数组。每个子节点包含了某个后继字符。

### 2.2 Trie 结构

```go
type Trie struct {
    previousWord string
    root         *TrieNode
    cache        []*TrieNode
    nodeCount    uint
}
```

- `previousWord`: 存储上一次插入的单词（用来做公共前缀计算）。
- `root`: 根节点，通常是一个哨兵节点（其 `letter` 为 `" "`，不参与实际字符匹配）。
- `cache`: 在批量插入中，用于加速公共前缀定位的缓存数组；存储了从根节点开始，到当前节点这条路径上的节点引用。
- `nodeCount`: 记录 Trie 中节点总数。

#### `Init()` 初始化

```go
func (t *Trie) Init() {
    t.previousWord = ""
    t.root = &TrieNode{ letter: " ", final: false }
    t.cache = append(t.cache, t.root)
    t.nodeCount = 1
}
```

初始化时，`root` 是一个空字符节点；将其放入 `cache` 里，并将节点计数设置为 1。

#### `GetNodeCount()`

```go
func (t *Trie) GetNodeCount() uint {
    return t.nodeCount
}
```

直接返回节点总数。

#### `Insert(word string)`

```go
func (t *Trie) Insert(word string) {
    // 1) 找到 word 与 t.previousWord 的公共前缀长度
    // 2) 截断 cache，仅保留公共前缀部分（因为公共前缀后可能需要创建新的节点）
    // 3) 逐字符向下创建 / 查找子节点
    // 4) 最后一个节点标记 final = true
    // 5) 更新 t.previousWord
}
```

在注释里已写得很清楚，关键逻辑是：

1. 先根据上一次插入的单词 `previousWord`，找出跟当前要插入的 `word` 的公共前缀部分（以节省重复节点）。
2. `t.cache = t.cache[:commonRuneCount+1]` 会把缓存数组截断到公共前缀节点。
3. 然后从公共前缀节点向下，根据当前单词剩余的字符，依次检查是否已有相同子节点：
   - 如果已经存在，直接复用，继续向下。
   - 如果不存在，则新建一个子节点，追加进 `children` 数组里，并让 `t.nodeCount++`。
4. 最后一个节点标记 `final = true`，表示这个路径是一个完整单词。
5. 更新 `t.previousWord = word`。

#### `Apply(fn func(*TrieNode))`

```go
func (t *Trie) Apply(fn func(*TrieNode)) {
    // 以队列的形式进行层序遍历，每访问一个节点就执行 fn(node)
}
```

通过层序遍历 (Level Order Traversal) 对每个节点执行传入的回调函数 `fn`。在后面的 `Encode()` 会用到它。

#### `Encode()`

```go
func (t *Trie) Encode() string {
    // 1) 首先写入整棵 Trie 的“结构”信息：对于每个节点，依次写一串 1（代表有子节点）+ 0 （收尾）
    // 2) 写入每个节点的 dataBits 位信息：第1位是 final，后面 (dataBits-1) 位是字符编码
}
```

1. **结构的 Unary 编码**：
   - 在层序遍历中，依次对每个节点写入它的孩子个数个 `1`，然后再写一个 `0`。
   - 这样就形成了一个“01 序列”，让我们可以在解码时知道各个节点有多少个孩子。
2. **具体节点信息**：
   - `dataBits` 由 `getDataBits(allowedCharacters)` 得到。它 = 1（表示是否是终结） + 若干位（可唯一表示字母表中的字符）。
   - 假设字母表大小为 27，那么需要 5bits 来表示其中任意字符（\(2^5 = 32\) > 27），再加上 1bit 的终结标志，总共 6bits = `dataBits`。
   - 最终会把每个节点 (final_bit + 字符编码) 写入到位串中。

`Encode()` 返回的字符串即是**紧凑表示**后的整棵 Trie——它由两大部分组成：

1. **结构位串**：存放在最前面的“1...10...0”序列。长度为节点数 \(\*\) 多少? 每个节点一串 `children_count` 个 `1` + 1 个 `0`。
2. **字符位串**：对于每个节点，写入 `dataBits` 位字符信息。

---

## 3. 字母表及字符到整数的映射（alphabet）

```go
var allowedCharacters = "abcdefghijklmnopqrstuvwxyz "
var mapCharToUint = getCharToUintMap(allowedCharacters)
var mapUintToChar = getUintToCharMap(mapCharToUint)
var dataBits = getDataBits(allowedCharacters)
```

- `allowedCharacters`: 允许出现的字符，比如示例中设为 `[a-z] + 空格` 共 27 个字符。
- `mapCharToUint` / `mapUintToChar`: 字符与整型 ID 的双向映射；用于在写入位串时快速映射为整数，在解码时再映射回字符。
- `dataBits = getDataBits(alphabet)`: 计算写一个节点所需的位数。其计算逻辑是：
  1. 先求能容纳 `alphabet` 大小的最小 \(\lceil \log_2(\text{alphabet_size}) \rceil\)。
  2. 再加 1bit 用于标记 `final`。

因此，若字母表含 27 个字符，则 `\(\lceil \log_2(27)\rceil = 5\)`，再加 1bit = 6bits。

---

## 4. Base64 编码相关（base64）

```go
var BASE64 = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"
var W uint = 6
```

- 这里用的是自定义的 Base64 字符表（最后两个符号是 `-` 和 `_`），并定义 `W = 6`，说明每个字符可表示 6 个比特。
- `CHR(id uint) string`: 根据一个 6bit 数值，返回对应的 Base64 字符；
- `ORD(ch string) uint`: 反向操作；给定一个字符，返回对应的 6bit 数值。

**注意**：无论后续怎么在位级别上读写，最终都会被压缩成一串 Base64（或说 6-bit chunk）的形式存储。

---

## 5. BitString（位串的读操作）

`BitString` 主要管理“从 Base64 字符串中按位读取”的逻辑。它的关键字段和方法：

```go
type BitString struct {
    base64DataString string
    length           uint
}
```

- `base64DataString`: 底层存储的编码字符串（Base64）。
- `length`: 该字符串总共有多少位，通常 = `len(base64DataString) * W`，其中 `W=6`。

### `Init(data string)`

```go
func (bs *BitString) Init(data string) {
    bs.base64DataString = data
    bs.length = uint(len(bs.base64DataString)) * W
}
```

### `Get(p, n uint) uint`

```go
func (bs *BitString) Get(p, n uint) uint {
    // 从第 p 位开始，连续取 n 位。可能会跨越多个 6-bit 块
    // 需要按照 p%W、(p/W) 等计算索引和偏移
}
```

- 若读取范围 `(p%W)+n <= W`，说明在同一个 Base64 字符之内，就可一次性取到结果。
- 否则，需要分多次取，不断拼接结果。

### `Count(p, n uint) uint`

```go
func (bs *BitString) Count(p, n uint) uint {
    // 统计从 p 开始的 n 位里，“1”的个数
    // 每次最多比较 8 位，然后利用 BitsInByte 预先计算的表来加速计数
}
```

### `Rank(x uint) uint`

```go
func (bs *BitString) Rank(x uint) uint {
    // 从开头到 x 位置（含 x）之间，1 的个数。
    // 这是一个相对慢的实现(循环逐位Get)，用来测试验证。
}
```

**小结**：`BitString` 负责底层的位级读取操作，包括提取指定位置的位数据和统计 1 的数量。

---

## 6. BitWriter（位串的写操作）

`BitWriter` 用于在**构建**紧凑编码时，把一个个整数拆分成位写入。它没有太多优化，仅用于**编码阶段**。

```go
type BitWriter struct {
    bits []uint  // 用一个uint数组，临时存储每一位
}

func (bw *BitWriter) Write(data, numBits uint) {
    // 逐位写入 bits[] 中
}

func (bw *BitWriter) GetData() string {
    // 最后把 bits[] 打包成 Base64 字符串
    // 每6位组成一个 0~63 的数，再映射到 BASE64字符
}

func (bw *BitWriter) GetDebugString(group uint) string {
    // 调试用，输出如 "01001 1011 ..." 这种可读的二进制形式
}
```

- `Write(data, numBits)`: 把 `data` 的低 `numBits` 位，依次压入 `bits[]`。
- `GetData()`: 以 6 位为一组，通过 `CHR()` 转成 Base64 字符串。

---

## 7. FrozenTrie 及其节点（解码并快速查询）

`FrozenTrie` 是**紧凑 Trie** 的只读结构；配合 `RankDirectory` 可以快速地从编码里**反向**定位到各个节点并进行查询。它包含三个主要部分：

```go
type FrozenTrie struct {
    data        BitString     // 整棵 Trie 的编码(包括结构与字符)
    directory   RankDirectory // Rank/Select 索引
    letterStart uint          // 指向 Trie 各节点字符数据区起始位置
}
```

1. `data`: 就是从 `Encode()` 返回的 Base64 字符串，转为 `BitString` 后保存的。
2. `directory`: RankDirectory 对象，用于支持快速 “Select(0, rank)” 和 “Rank(1, x)” 等操作。
3. `letterStart`: 因为前面**结构位串**（那堆 1…10…0…）占用了一定 bit 数，所以要记录从哪一位开始才是各节点的字符位信息。

### `Init(data, directoryData string, nodeCount uint)`

```go
func (f *FrozenTrie) Init(data, directoryData string, nodeCount uint) {
    f.data.Init(data)
    f.directory.Init(directoryData, data, nodeCount*2+1, L1, L2)

    // letterStart = nodeCount*2 + 1
    // 这里 nodeCount*2+1 是结构位串的长度(1...10...0)，后面才是每个节点的 dataBits
}
```

1. `f.data.Init(data)` 初始化紧凑 Trie 的位串。
2. `f.directory.Init(directoryData, data, nodeCount*2+1, L1, L2)` 初始化 RankDirectory，以便后面快速查找结构位串中的“0/1”分布。
3. `f.letterStart` 用来告诉我们：在位串中，字符数据从哪一位开始。根据前面的写法，每个节点都会写入**若干个 1** + **1个 0**，于是总长度 = `nodeCount*2 + 1` 位（因为每个节点至少有一个“0”，外加它孩子的数目个“1”，总和正好等于 `2*nodeCount - 1` 再加一些常数，详见博客 / 论文中对 Trie 结构位串的说明）。

### `GetNodeByIndex(index uint) FrozenTrieNode`

```go
func (f *FrozenTrie) GetNodeByIndex(index uint) FrozenTrieNode {
    // 1) 取出该节点的 final 标记
    // 2) 取出该节点所代表的字符
    // 3) 通过 directory.Select(0, index+1) 定位到第 (index+1) 个 '0'
    //    并由此推断该节点第一个子节点的索引
}
```

逻辑如下：

1. `final = (f.data.Get(f.letterStart+index*dataBits, 1) == 1)`
   - 在“字符位串”区域里，先读 1bit，判断是否终结。
2. `letter = mapUintToChar[f.data.Get(f.letterStart+index*dataBits+1, dataBits-1)]`
   - 再读 (dataBits-1) 位，映射到具体字符。
3. `firstChild = f.directory.Select(0, index+1) - index`
   - `directory.Select(0, x)` 会找出“第 x 个 0 出现的位置”在结构位串中的下标。
   - “结构位串”里每个节点对应的“0”标志着这个节点孩子的结尾。
   - 减去 `index`，是因为在 level-order 下，每个节点在结构位串中出现的顺序和下标是相对应的。
4. `childOfNextNode = f.directory.Select(0, index+2) - (index+1)`
   - 这是下一个节点对应的 0 的位置，再减掉 `(index+1)`。可以得到当前节点孩子区域的结束边界。
   - 因此 `childCount = childOfNextNode - firstChild`。

返回 `FrozenTrieNode`，包含了该节点的必要信息。

### `GetRoot() FrozenTrieNode`

根节点是 `index = 0`。

### `Lookup(word string) bool`

```go
func (f *FrozenTrie) Lookup(word string) bool {
    node := f.GetRoot()
    for i, w := 0, 0; i < len(word); i += w {
        // 1) 取出下一个字符 runeValue
        // 2) 在 node 的所有子节点里找 child.letter == runeValue
        // 3) 若没找到，返回 false
        // 4) 否则更新 node = child
    }
    // 全部字符遍历完后，看 node.final 是否为 true
    return node.final
}
```

和普通 Trie 的遍历类似，只是内部要通过 `GetChildCount()` 和 `GetChild(j)` 获取子节点。因为是冻结结构，不能直接拿一个指针去 `children`，而要经过 `directory` 的 Rank/Select 计算。

---

## 8. RankDirectory（Rank 和 Select 的快速查询结构）

`RankDirectory` 用来在一个长位串（这里即结构位串）里，快速计算 “有多少个1” 和 “第 y 个 0 在哪” 之类的问题。

```go
type RankDirectory struct {
    directory   BitString // 存放L1/L2表
    data        BitString // 原始的位串(需要支持Rank/Select)
    l1Size      uint
    l2Size      uint
    l1Bits      uint
    l2Bits      uint
    sectionBits uint
    numBits     uint
}
```

这里的思路：

1. 把 `data` 分段（L1 / L2 分段），预先计算某些“分界点之前有多少个 1”。
2. 用两级表（L1：较大跨度；L2：较小跨度）来减少查询距离。
3. 这样就能在 `O(1)` 或 `O(log n)` 的时间内完成 `rank()`。
4. `select()` 则用二分法 + `rank()` 去定位第 `y` 个 0 或第 `y` 个 1 的具体位置。

#### `CreateRankDirectory(data string, numBits, l1Size, l2Size uint)`

```go
func CreateRankDirectory(data string, numBits, l1Size, l2Size uint) RankDirectory {
    // 1) 遍历 data，每 L2 位累计1的个数写入 directory
    // 2) 若恰好走了 L1 位，则把累计值写入 L1 表，并清零
    // 3) 生成的 directoryData 用 BitWriter 写好后
    // 4) RankDirectory.Init(...) 初始化
}
```

- `numBits = trie 的结构位串长度 = nodeCount*2 + 1`。
- `l1Size = L1`，`l2Size = L2`，本例中 `L1 = 32*32 = 1024`，`L2 = 32`。
- 最终在 `directory.GetData()` 拿到一串 Base64 字符串，里面存了所有分块的累积计数。

#### `Rank(which, x uint) uint`

```go
func (rd *RankDirectory) Rank(which, x uint) uint {
    // 计算前 x 位中 which=1 的个数(或 0 的个数)
    // which=0 时，可用 “x+1 - Rank(1, x)” 转化为 1 的计数
    // 再拆成 L1 段 + L2 段 + 剩余 bits 的计数之和
}
```

- 先根据 `x / l1Size` 去 directory 的 L1 表找“已知的总和”，然后余下部分根据 `x / l2Size` 找 L2 表，再加上剩余不满 L2 的位用 `data.Count()` 补足。

#### `Select(which, y uint) uint`

```go
func (rd *RankDirectory) Select(which, y uint) uint {
    // 找到第 y 个 which(0/1) 出现的位置
    // 用二分搜索 + Rank(1, mid) 来逐步收敛
}
```

---

## 9. 前缀搜索函数（GetSuggestedWords）

在 `FrozenTrie` 里，还有一个用来做**前缀搜索**的小函数：

```go
func (f *FrozenTrie) GetSuggestedWords(word string, limit int) []string {
    // 1) 找到与 word 匹配的节点 node
    // 2) 然后从 node 开始做一个 BFS / level-order，把能连到的单词收集起来
    // 3) 最多返回 limit 个
}
```

具体流程：

1. 先从 `GetRoot()` 出发，按每个字符依次在子节点里匹配。
2. 如果中途匹配失败，直接返回空结果。
3. 找到那个节点后，执行 `traverseSubTrie(node, prefix, limit)`。
   - 这里通过队列方式（与普通 BFS 类似）把后续所有可达节点都扫描一次。
   - 遇到 `final = true` 的节点，就把当前前缀加入结果数组。
   - 直到收集满 `limit` 个，或者把整棵子树遍历完。

---

# 总结

本代码展示了如何从传统 Trie 数据结构出发，**构建**出一棵**紧凑（Succinct）**的 Trie，并使用**RankDirectory** 来实现快速地对结构位串进行 Rank/Select 查询，进而在一个只读的 `FrozenTrie` 上能够快速执行：

- **Lookup(word)**：判断一个单词是否在 Trie 中。
- **GetSuggestedWords(prefix, limit)**：根据前缀做补全 / 提示。

其中，最核心的两个点在于：

1. **Trie 压缩为位串**：

   - 前半部分记录树的结构（各节点有多少孩子，何时终止等），采用“对每个节点写 `children_count` 个 1，再补上一个 0”的方式。
   - 后半部分记录节点的字符及是否为单词终结，一共 `dataBits` 位（1 位标记终结，其余位表示字符 ID）。

2. **RankDirectory** 的使用让“到第几个节点为止有多少 1/0”或者“第 n 个 0/1 在哪” 这种查询变得高效，从而快速定位层序遍历下任意节点的起止区间。

在真实的生产环境中，这种**Succinct Trie** 能在**大规模字典**（尤其是很多前缀重复）时，既节省存储空间，又能维持较快的查询性能。
