下面这段代码是 **Steve Hanov** 在其博客（[http://stevehanov.ca/blog/?id=115](http://stevehanov.ca/blog/?id=115)）中讨论的基于 **DAWG** (Directed Acyclic Word Graph) 的一个实现，侧重于**内存和磁盘效率**。它允许在构建完毕后，将数据结构序列化到文件或内存中，然后通过只读的方式“原地”使用该结构来进行查询。

代码相对庞大，下面分块做系统性分析。

---

## 目录

1. [总体流程和主要数据结构](#总体流程和主要数据结构)
2. [构建过程：Builder](#构建过程builder)
3. [完成构建：Finish()](#完成构建finish)
4. [DAWG 查找功能：Finder](#dawg-查找功能finder)
5. [序列化与反序列化](#序列化与反序列化)
   - [写出：Write() 和 Save()](#写出write-和-save)
   - [读入：Read()](#读入read)
6. [位操作与文件格式](#位操作与文件格式)
7. [重要函数和流程小结](#重要函数和流程小结)

---

## 1. 总体流程和主要数据结构

该实现将 DAWG 的流程分为两大阶段：

1. **构建** (Builder 阶段)：可以多次 `Add` 单词，但要求这些单词必须**严格按字典序**（strictly increasing）添加，不能重复。构建完毕后调用 `Finish()`。
2. **查询** (Finder 阶段)：`Finish()` 后，便得到一个 **Finder** 接口，可用于各种查询、序列化到文件等操作。此时不再允许添加单词。

### 1.1. 核心结构概览

- **`type dawg struct { ... }`**  
  既是构建期的 “Builder” 又是完成后的 “Finder”，通过内部标志 `finished` 区分所处阶段。

- **在“构建期”使用的字段**

  ```go
  lastWord       []rune
  nextID         int
  uncheckedNodes []uncheckedNode
  minimizedNodes map[string]int
  nodes          map[int]*node
  ```

  主要存储 Trie/DAWG 的临时结构，以及**最小化**所需的辅助信息。

- **在“完成后”使用的字段**

  ```go
  finished        bool
  numAdded        int
  numNodes        int
  numEdges        int

  // 读写相关
  r               io.ReaderAt
  size            int64
  cbits, abits    int64
  wbits           int64
  firstNodeOffset int64
  hasEmptyWord    bool
  ```

  这里包含：

  - `r`/`size`：指向序列化后文件或内存的只读接口；
  - `cbits`/`abits`/`wbits`：在写盘时确定的各类“位宽”，用于后续解析；
  - `hasEmptyWord`：记录根节点是否就是一个结束节点（空串）；
  - `numAdded`、`numNodes`、`numEdges` 用于统计信息等。

### 1.2. 辅助结构

- **`type node struct`**

  ```go
  type node struct {
      final bool
      count int
      edges []edgeStart
  }
  ```

  - `final`: 是否是一个单词的结尾。
  - `count`: 用来存储“从该节点可达的终止数目” （构建后计算，用于加速 IndexOf、AtIndex 等操作）。
  - `edges`: 子边列表，每个元素是 `edgeStart{node int, ch rune}`。

- **`type uncheckedNode struct`**  
  用于在最小化过程中暂时存储“尚未固定/合并”的路径信息：

  ```go
  type uncheckedNode struct {
      parent int
      ch     rune
      child  int
  }
  ```

- **`type edgeStart struct` 与 `type edgeEnd struct`**
  - `edgeStart`：查询时输入“起点+字符”，用于定位下一状态；
  - `edgeEnd`：则存储“下一节点 ID + 跳过多少个单词计数”等信息。

---

## 2. 构建过程：Builder

当我们调用 `dawg.New()` 时，得到一个处于“构建期”的 `dawg`。构建单词的核心方法是 `Add(wordIn string)`。

### 2.1. `Add(wordIn string)`

```go
func (d *dawg) Add(wordIn string) {
    // 1) 确保 wordIn > 上一个 word (严格字典序)，否则 panic
    // 2) 找到新单词和上个单词的公共前缀长度 commonPrefix
    // 3) 对 [commonPrefix, end) 的 uncheckedNodes 进行 minimize
    //    “把之前多余的后缀进行最小化合并”
    // 4) 把剩余的后缀插入 Trie / DAG
    // 5) 标记结尾节点 final = true
    // 6) 更新 lastWord、numAdded
}
```

其主要逻辑：

1. **字典序保证**：若当前插入单词小于等于上次插入的单词，直接报错。
2. **公共前缀**：计算当前单词和 `lastWord` 的公共前缀位置 `commonPrefix`。
3. **合并之前多余的后缀**：`d.minimize(commonPrefix)` 会把 `uncheckedNodes` 队列中，从末尾到 `commonPrefix` 处的节点逐个检查，如果有可以合并的子树就合并。
4. **对新产生的后缀**：从公共前缀之后，把字符一个个加进来，新建节点并插入 `edges`。同时将新建节点推到 `uncheckedNodes` 里等待后续合并。
5. **设置最终节点为 final**。
6. 更新 `lastWord`，`numAdded++`。

### 2.2. `minimize(downTo int)`

- 这个函数是**最小化**的关键：从 `uncheckedNodes` 的末尾往前处理，依次调用 `nameOf(child)` 构建子树的字符串标识（包含所有边及其子节点 ID），若已经在 `minimizedNodes` 映射里，说明有重复结构，就把当前的 child 换成已经存在的那个节点；否则把它插入 `minimizedNodes`。
- 处理完毕后，截断 `uncheckedNodes` 列表到指定长度 `downTo`。

### 2.3. `addChild/replaceChild`

- `addChild(parent, ch, child)`: 在 `parent` 节点的 `edges` 数组尾部追加一个 `{child, ch}`；并 `numEdges++`。
- `replaceChild(parent, ch, child)`: 在 `parent` 节点里找到对应 `ch` 的边，然后把它的 `node` 替换为 `child`。同时删除原先 `child` 对应的节点信息（因为已被合并）。

**注意**：构建时，要求 `edges` 数组中的字符必须是严格递增的，以保证文件序列化时更有效率。

---

## 3. 完成构建：`Finish()`

当所有单词都 `Add` 完，调用 `Finish()` 来：

1. 再次 `minimize(0)`：合并剩余的 `uncheckedNodes`。
2. `numNodes = len(d.minimizedNodes) + 1`：这里 `+1` 通常是根节点。
3. `d.calculateSkipped(rootNode)`: 给每个节点计算 `count`，代表“从此节点可达多少终止单词”。这在实现 `IndexOf`、`AtIndex` 时会用到。
4. `renumber()`: 重新给节点编号。构建阶段我们可能创建了很多“废弃”节点，被合并后就不再使用，`renumber()` 会把现有节点重新组织成 0,1,2... 的连贯编号，以便在后续序列化中更加紧凑、简洁、可预测。
5. 把所有与构建相关的临时结构 (`uncheckedNodes`, `minimizedNodes`, `lastWord`) 清空，只保留最终统计信息和节点映射。
6. **序列化到内存**：它调用一次 `d.Write(&buffer)`，将结果写到一个内存 buffer 里，并以只读方式回放赋给 `d.r`。这样 `dawg` 就已经是**只读**结构，可用来查询了。

最后返回一个 `Finder` 接口给用户。

---

## 4. DAWG 查找功能：Finder

一旦构建完成（`finished == true`），该结构提供多种查询方式。

### 4.1. `FindAllPrefixesOf(input string)`

```go
func (d *dawg) FindAllPrefixesOf(input string) []FindResult {
    // 1) 从 rootNode 开始
    // 2) 逐字符匹配 getEdge(...)
    // 3) 每次若当前节点是 final，就记录一个 (prefix, skippedCount)
    // 4) 若中途没有对应边，则停止
    // 5) 结尾若最后节点也 final，则也加入结果
}
```

- `skipped`：计数“在当前分支之前，被跳过的单词数目”。
- `edgeEnd.count`：每次沿某边下来，会加上相应的 skip 值。
- 最终得到 “该输入字符串所有前缀在 DAWG 中对应的索引”列表。

### 4.2. `IndexOf(input string)`

```go
func (d *dawg) IndexOf(input string) int {
    // 类似逐字符找对应边
    // 累加 edgeEnd.count 到 skipped
    // 如果到达终点是 final，则返回 skipped；否则 -1
}
```

### 4.3. `AtIndex(index int) (string, error)`

```go
func (d *dawg) AtIndex(index int) (string, error) {
    // 在已知 total count 下, DFS 或者按照 edges 的 skip, 判断哪个分支包含目标 index
    // 递归调用 d.atIndex(...)，每次看 skip 值决定向哪个子边走
    // 当遇到 final 节点且 skip == targetIndex, 即找到对应的单词
}
```

大致流程：

1. 如果当前节点是 final 并且 `atIndex == targetIndex`，说明命中；
2. 否则，根据 `node.edges[i].count` 的累加范围找到合适分支继续深入。
3. 递归直到找出目标或失败。

### 4.4. `Enumerate(fn EnumFn)`

```go
func (d *dawg) Enumerate(fn EnumFn) {
    // 遍历整棵 DAWG 的所有可能前缀
    // 对每个前缀，若是 final 则回调 fn(index, prefix, true)
    // 若 fn 返回 Continue，则继续往子边深入
    // 若 fn 返回 Skip，则跳过子边，不再往下
    // 若 fn 返回 Stop，则整个枚举结束
}
```

- 构造 `nodeResult` 来读取边信息，深度遍历 DAG。
- 适合做**批量遍历**、调试或其它需要遍历所有前缀的场景。

---

## 5. 序列化与反序列化

本代码非常注重**存储格式**的紧凑性，核心思路是**按位写入** (bit-level encoding)，最大程度节省空间。同时，它支持：

- **`Save(filename)`** or `Write(w io.Writer)`：将构建好的 DAWG 写出到文件 / 流。
- **`Load(filename)`** 或 `Read(f io.ReaderAt, offset int64)`：从文件 / 流中读取 DAWG，直接以只读方式访问。

### 5.1. 写出：`Write()` 和 `Save()`

1. **准备阶段**
   - 计算 `cbits`：存储字符所需的位数 (`bits.Len(uint(maxChar))`)。
   - 计算 `wbits`：存储单词数所需的位数 (`bits.Len(uint(numAdded))`)。
   - 计算 `abits`：存储节点地址所需的位数，需要保证能表示“所有 bit offset”。
   - 计算**文件总大小** `size`：遍历所有节点，估算写出时的 bit 数量，直到可以容纳所有地址。
2. **开始写文件头**：
   - 32 bits: 文件总大小 (字节数)。
   - 8 bits: `cbits`
   - 8 bits: `abits`
   - 变长编码 `numAdded`，`numNodes`，`numEdges`。
3. **写节点信息**：
   - 每个节点先写 “final” (1 bit)，再写 “fallthrough” (1 bit)。
   - 如果 `fallthrough`==1，则只需写这个字符 + 下一个节点位置；
   - 否则需要写“单边 singleEdge?” (1 bit)。
     - 若多条边则写“多少 edges + skipFieldLen + 每条边的 (字符 + skipValue + nodeAddress)”。
     - `skipValue` 就是 `count` 的差值，用于快速跳过若干单词。
4. **Flush** 并返回写出的总大小。

### 5.2. 读入：`Read(f io.ReaderAt, offset int64) (Finder, error)`

1. 读前 4 bytes 获得总大小；
2. 读后续的 `cbits`, `abits`，以及变长编码的 `numAdded`, `numNodes`, `numEdges`；
3. 读第一节点是否是空串等标记；
4. 保存到一个新的 `dawg` 结构中，设置 `finished=true`、`r=f`、`size=...` 等；
5. 返回该 `dawg`。之后所有查询都通过**位读取**来解析节点。

**关键之处**： 读出的 Finder 不会把整个结构加载到内存，而是**懒解析**：每次查询时，通过 `bitSeeker` 在文件中定位到某个节点的 bit offset，然后读其结构、决定下一步走向，因而在巨大文件场景下仍然只读“用到的部分”。

---

## 6. 位操作与文件格式

本代码有大量对“bit-level”数据的读写：

- **`type bitWriter`**

  - `WriteBits(data uint64, n int)`: 一次写 n bits 到底层流里；
  - `Flush()`: 把未写满的 byte 补足后写出。

- **`type bitSeeker`**
  - 维护一个 64-bit 缓存 `cache`，只要访问到某个 64-bit 边界，就 `ReadAt()` 读取 8 字节并 BigEndian 存入 `cache`；
  - `ReadBits(n int64)`: 在当前 `p`(bit 偏移) 下取 n bits。若超出当前 64-bit，就跨越下一个64-bit缓存；
  - `Seek(offset int64, whence int)`, `Tell()`：对 bit 偏移进行移动/查询。

**文件格式简述**（代码注释处有完整描述）：

```text
[0..31 bits]:  文件总大小 (单位：字节)
[32..39 bits]: cbits (存储字符用多少位)
[40..47 bits]: abits (存储节点偏移地址用多少位)
后面使用变长 7-bit 编码记录：numWords, numNodes, numEdges
然后开始写每个 node 的信息:
  - final (1 bit)
  - fallthrough (1 bit)
  - if fallthrough == 1:
      cbits: character
      ... -> node offset
    else:
      singleEdge (1 bit)
      if singleEdge == 0:
         write numberOfEdges (7code)
         write skipfieldLen
      for each edge:
         cbits: character
         if not the first edge:
             skipValue: #skipfieldLen bits
         abits: nextNodeAddress
  ...
```

---

## 7. 重要函数和流程小结

1. **`New()`**: 初始化一个待构建的空 DAWG（Builder）。
2. **`Add(word)`**: 按字典序插入单词，内部进行**最小化**合并。
3. **`Finish()`**: 标记构建完毕，计算并压缩，最终序列化到内存并作为 Finder 返回。
4. **`FindAllPrefixesOf(input)`**: 查找全部前缀对应的单词索引。
5. **`IndexOf(input)`**: 返回单词对应的插入顺序（若不存在返回 -1）。
6. **`AtIndex(index)`**: 根据插入顺序反向取出单词。
7. **`Save(filename)` / `Write(io.Writer)`**: 将 DAWG 的位编码完整写出到文件/流，以便下次重用。
8. **`Read(io.ReaderAt, offset)`**: 从外部文件/流中按同样格式读取，生成只读 DAWG。

**整体而言**，这是一个**面向大规模字典**、且关心**磁盘/内存占用**的 DAWG 实现。它通过**位级**编码和“跳过计数 (skip counts)”技巧，大幅减少存储体积，并支持按需在文件中定位节点，进行查询。与此同时，要求在构建期严格按字典序插入，以便实现在线最小化和单调 edges。

- 构建期使用 `nodes map[int]*node` 来维护基本 Trie + 后缀合并；
- 完成后进行序列化，序列化后再以 `bitSeeker` / `bitWriter` 的方式操作，真正做到了**按需访问**与**紧凑存储**相结合。

在实际应用中，如果想一次性加载到内存进行超快速查询，也可以将 `Write` 的结果读到内存（譬如 `bytes.Reader`），然后 `Read` 回来，就能在内存中以同样的方式“只读”查询，无需额外的数据结构。

---

**总结**：这段代码展示了一个完整的 **DAWG** (最小化 Trie) 的构建与序列化过程，兼顾**空间压缩**和**可外部化存储**。实现细节非常值得学习，尤其是位操作、最小化合并、和严格的字典序插入等部分。

---

下面这段来自 [Steve Hanov](http://stevehanov.ca/blog/?id=115) 的 **DAWG** 实现代码较为庞大，且包含了**从构建到序列化，再到只读查询**的完整功能。它的核心目标是：

1. **构建**：以最小化 Trie（即 DAWG）方式存储大量有序单词（需要按严格字典序插入），以减少冗余节点；
2. **序列化**：将构建好的结构用**按位（Bit-Level）**的方式紧凑存储到文件或内存；
3. **只读查询**：在完成构建后，可以以只读方式（无须重新加载全部节点到内存）进行各种查询：\(`IndexOf`, `FindAllPrefixesOf`, `AtIndex` 等\)；
4. **跳过计数 (skip count)**：每条边带一个“跳过了多少个单词”的计数，这使得 `IndexOf`、`AtIndex` 等功能可以在一条边上一次性地跳过若干单词索引，实现更快的定位。

下面从**设计思路**到**关键函数**，再到**文件格式和位操作**做一个全面、细致的解读。

---

## 目录

1. [整体结构与关键数据类型](#整体结构与关键数据类型)
2. [构建期（Builder 阶段）](#构建期builder-阶段)
   - 2.1. `Add` 方法
   - 2.2. 最小化合并 `minimize`
   - 2.3. 节点数据结构 `node`
3. [完成构建与转换（`Finish`）](#完成构建与转换finish)
   - 3.1. `calculateSkipped`
   - 3.2. `renumber`
   - 3.3. 预写入到内存并切换到 Finder
4. [查询期（Finder 阶段）](#查询期finder-阶段)
   - 4.1. `FindAllPrefixesOf`
   - 4.2. `IndexOf`
   - 4.3. `AtIndex`
   - 4.4. `Enumerate`
   - 4.5. 其它统计方法
5. [序列化与反序列化](#序列化与反序列化)
   - 5.1. 写出流程：`Write` / `Save`
   - 5.2. 读取流程：`Read` / `Load`
   - 5.3. 文件格式概览
6. [位操作辅助（bitWriter / bitSeeker）](#位操作辅助bitwriter--bitseeker)
   - 6.1. `bitWriter`
   - 6.2. `bitSeeker`
   - 6.3. 7-bit 编码（`writeUnsigned` / `readUnsigned`）
7. [小结与扩展](#小结与扩展)

---

## 1. 整体结构与关键数据类型

在 Go 代码最顶层，我们可以看到以下结构与接口：

- **`Finder`** 和 **`Builder`** 接口：

  - `Builder` 用于构建 DAWG（`Add` 单词、`Finish` 完成）。
  - `Finder` 用于在只读状态下进行查询、遍历、序列化等操作。

- **`type dawg struct`** 同时实现了 `Builder` 和 `Finder`：

  ```go
  type dawg struct {
      // 构建期使用的字段
      lastWord       []rune           // 上一次插入的单词，用于确保有序和找公共前缀
      uncheckedNodes []uncheckedNode  // 存放当前还未最小化的 “路径节点”
      minimizedNodes map[string]int   // 已最小化的状态(子树) 缓存: signature -> nodeID
      nodes          map[int]*node    // 所有节点的临时存储

      // 统计信息、标志
      finished  bool
      numAdded  int  // 已插入的单词数
      numNodes  int  // 节点数（完成后才确定）
      numEdges  int  // 边数

      // 位存储 / 文件结构相关
      r               io.ReaderAt
      size            int64
      cbits, abits    int64  // character bits / address bits
      wbits           int64  // word count bits
      firstNodeOffset int64
      hasEmptyWord    bool   // 根节点是否也是终止状态 (存空串)
  }
  ```

- **`type node struct`** 表示构建期的一个节点：

  ```go
  type node struct {
      final bool
      count int
      edges []edgeStart
  }
  ```

  - `final`: 是否是单词结束节点；
  - `count`: 在构建完成后，记录该节点下可达的单词数量，用于快速 skip；
  - `edges`: 子边数组，元素类型 `edgeStart{node int, ch rune}`，表示“(当前节点) -ch-> (下个节点ID)”。

- **`type edgeStart struct` 与 `type edgeEnd struct`**：
  - `edgeStart`：记录边的起点信息 `(node, ch)` 用于在查询时找下一节点；
  - `edgeEnd`：记录边的目的信息 `(node, count)`；`count` 表示“沿此边可以跳过多少单词”。

> 注意：当**构建完成**后，会把 `nodes`、`uncheckedNodes`、`minimizedNodes` 等通通清空，而是借助 **`io.ReaderAt`** + **位操作** 在文件或内存中**按需解析**节点。

---

## 2. 构建期（Builder 阶段）

### 2.1. `Add` 方法

当我们 `New()` 得到一个空的 `dawg` 后，就可以多次调用 `Add(word string)` 向其插入单词。关键约束是**单词必须严格按字典序**递增，否则会触发异常。

```go
func (d *dawg) Add(wordIn string) {
    // 1) 校验 wordIn > lastWord
    // 2) 求与上一个单词的公共前缀 commonPrefix
    // 3) d.minimize(commonPrefix)
    // 4) 从commonPrefix处开始为新单词的后缀创建节点
    //    并加入 uncheckedNodes
    // 5) 标记最后节点 final
    // 6) 记录 lastWord = word
    //    numAdded++
}
```

大体流程：

1. **公共前缀**：找到当前单词 `wordIn` 与 `lastWord` 相同的前缀长度 `commonPrefix`。这意味着 `[commonPrefix..end) `这一段相对于上一次插入的新内容需要继续插入。
2. **最小化合并**：先对“多余的” `uncheckedNodes` 中的节点（大于等于 `commonPrefix` 的部分）进行 `minimize`，把它们合并到 `minimizedNodes` 里，避免重复子树。
3. **新建后缀**：对新单词剩余的字符，逐个新建节点、添加子边 `addChild(...)`，并追加进 `uncheckedNodes`。
4. 最后一个节点设置 `final = true`，表示这是一个完整单词的结束位置。
5. 更新 `lastWord`、`numAdded` 等。

### 2.2. 最小化合并 `minimize`

```go
func (d *dawg) minimize(downTo int) {
    for i := len(d.uncheckedNodes) - 1; i >= downTo; i-- {
        u := d.uncheckedNodes[i]
        signature := d.nameOf(u.child)
        if node, ok := d.minimizedNodes[signature]; ok {
            // 如果已经有了相同结构
            d.replaceChild(u.parent, u.ch, node)
        } else {
            // 否则记录到 minimizedNodes
            d.minimizedNodes[signature] = u.child
        }
    }
    d.uncheckedNodes = d.uncheckedNodes[:downTo]
}
```

- 从 `uncheckedNodes` 的末尾往前处理，针对每个 `(parent, ch, child)`：
  1. `nameOf(u.child)`：生成 child 子树的“字符串标识”（包含了子边信息和 `final` 状态等）。
  2. 若在 `minimizedNodes` 中已存在同样的“字符串标识”，说明可以**复用**那个子树，则 `replaceChild(u.parent, u.ch, existingNodeID)` 并且删除重复节点；
  3. 否则 `minimizedNodes[signature] = u.child`。

这样的过程类似于许多最小化 Trie（DAWG）的实现方式：自底向上，对新增的路径做合并，从而避免重复子结构。

### 2.3. 节点数据结构 `node`

构建期的每个节点都保存在 `d.nodes[nodeID] = &node{ ... }` 当中，它包含了：

- `edges`: 一个 **有序** 的子边列表（按 `ch` 递增）。这是为了在最终序列化时能有效利用**顺序**，或者快速进行二分查找。
- `final`: 是否终止；
- `count`: 仅在最后收尾阶段才计算，用于记录从该节点可达多少个词（包括自身）。

---

## 3. 完成构建与转换（`Finish`）

当所有单词插入完毕后，调用 `Finish()` 才能获得一个可查询的结构。

### 3.1. `calculateSkipped`

`Finish()` 里会做：

1. `d.minimize(0)`：彻底合并剩余的 `uncheckedNodes`；
2. `d.calculateSkipped(rootNode)`：给每个节点算 `count` 值——从该节点可达多少终止单词。

```go
func (d *dawg) calculateSkipped(nodeid int) int {
    node := d.nodes[nodeid]
    if node.count >= 0 {
        return node.count
    }
    numReachable := 0
    if node.final {
        numReachable++
    }
    for _, edge := range node.edges {
        numReachable += d.calculateSkipped(edge.node)
    }
    node.count = numReachable
    return numReachable
}
```

- `node.count` 在查询中起重要作用：例如 `IndexOf` 时，需要知道每条边“跳过”多少个单词，以快速计算某单词的“索引”。

### 3.2. `renumber`

在最小化过程中，可能会产生很多合并后的“废弃” ID。为了让序列化更紧凑、地址映射更顺畅，需要对节点进行**重新编号**。

```go
func (d *dawg) renumber() {
    // DFS，从root开始，对访问到的节点进行递增编号
    // remap[id] = newID
    // 改完后，再把d.nodes根据newID重排
}
```

### 3.3. 预写入到内存并切换到 Finder

最后，`Finish()` 中执行：

```go
var buffer bytes.Buffer
d.size, _ = d.Write(&buffer)
d.r = bytes.NewReader(buffer.Bytes())
d.nodes = nil // 释放
```

- `d.Write` 会将整个 DAWG（节点信息）**按位**序列化到 `buffer`；
- 再让 `d.r` 指向这个 `buffer`，从而后续查询时只需要**bit-level 解析**；
- 同时把 `d.nodes` 置空，释放内存，表示后续全部“只读”操作都基于流 `r`。

此时，`dawg.finished = true`，它就变成了一个**只读 Finder**。

---

## 4. 查询期（Finder 阶段）

`Finish()` 返回的 `Finder` 接口就可执行各种查询，代码中提供了几个主要方法：

### 4.1. `FindAllPrefixesOf(input string)`

用来找出 “`input` 的所有前缀，在DAWG中属于已插入词的prefix” 及其对应的插入索引。例如，如果 DAWG 中有单词 `["a", "ab", "abc"]`，那么对输入 `"abcxyz"`，它会匹配出前缀 `"a"`, `"ab"`, `"abc"` 并给出这些前缀在 DAWG 中对应的 index。

```go
func (d *dawg) FindAllPrefixesOf(input string) []FindResult {
    // 1) 从根节点(node=0, skipped=0, final=hasEmptyWord)开始
    // 2) 对每个字符 letter:
    //    若节点是final，则加入 results
    //    查找是否存在 (node, letter) 边 => edgeEnd
    //    若不存在则停止
    //    否则 node = edgeEnd.node
    //    skipped += edgeEnd.count
    // 3) 若最后 node 也是 final，则也加进 results
}
```

- 在 `DAWG` 中查找边时，用到 `getEdge()`：它根据当前 `node` 的位偏移，解析“final bit, fallthrough bit, ...” 等信息找到对应字符 `letter` 的那条边以及 `count` 值。
- `skipped` 的含义是 “若我们走到这条边，说明我们跳过了多少个词的索引”。

### 4.2. `IndexOf(input string)`

返回某单词在 DAWG 中的插入顺序（0-based）。若不存在，则返回 -1。

```go
func (d *dawg) IndexOf(input string) int {
    // 与FindAllPrefixesOf类似，每次顺着边走
    // 并累加 edgeEnd.count 到 skipped
    // 最终若到达节点是 final，则返回 skipped，否则 -1
}
```

### 4.3. `AtIndex(index int) (string, error)`

反向：给定插入顺序 `index`，找出对应的单词。

- 原理是从根节点开始，根据 `node.count` 等信息定位应该走哪条边；
- 比如若当前节点 final 并且 `atIndex == targetIndex`，说明找到该单词；
- 否则根据子边的 `count` 判断要走哪条边。
- 递归一直往下，直到找到对应节点并把路径上的字符拼装起来。

### 4.4. `Enumerate(fn EnumFn)`

列举 DAWG 中的所有前缀（包括中间节点也会回调）：

```go
func (d *dawg) Enumerate(fn EnumFn) {
    // 以 DFS 方式遍历全部节点
    // 调用 fn(index, runes, final)
    // 其中 index 会累加 skipCount
    // 若 fn 返回 Continue 则继续往子节点走
    // 若 Skip 则跳过其子节点
    // 若 Stop 则立即退出整个遍历
}
```

此功能在调试或需要遍历整个词集时非常实用。

### 4.5. 其它统计方法

- `NumAdded()`：返回插入的单词总数；
- `NumNodes()`、`NumEdges()`：返回节点和边的数量；
- `Print()`：调用 `DumpFile(d.r)` 将序列化格式以**人类可读**的方式打印出来（主要用于调试）。

---

## 5. 序列化与反序列化

### 5.1. 写出流程：`Write` / `Save`

`Write(w io.Writer)` 函数会：

1. 若已存在 `d.r`（说明结构已经写过了），则直接 `io.Copy()` 把之前保存好的数据写到新的流里；否则就**现根据内存中的 `d.nodes` 做 bit-level 写出**。
2. 写文件头：
   - **32 bits** 存储文件总大小（字节数）；
   - **8 bits** 记录 `cbits`，再 8 bits 记录 `abits`；
   - 再用 **7-bit 变长编码**（`writeUnsigned`）写 `numAdded`, `numNodes`, `numEdges`；
3. 依次写每个节点信息：
   - `final` (1 bit)
   - `fallthrough` (1 bit)
   - 如果 `fallthrough==1`，只需写一个字符 + 下一个节点地址
   - 否则写是否 `singleEdge` (1 bit)；如果不止一条边，还要写“边数 + skipfieldLen”等。
   - 最后写每条边的 `ch`, `count`, `nodeAddress`。
4. 写完后 `Flush()`，返回总字节数。

`Save(filename)` 只是 `Create -> Write -> Close` 的简便包装。

### 5.2. 读取流程：`Read(f io.ReaderAt, offset int64) (Finder, error)`

- 首先读前 4 字节确定 `size`；
- 依次解析 `cbits, abits, numAdded, numNodes, numEdges` 等；
- 记录 `firstNodeOffset`：即写完这些头信息后，第一个节点在文件中的 bit 偏移；
- 读 `rootNode` 是否 `final` （判断空串），等等。
- 构造一个新的 `dawg`，设置 `r = f`, `size = size`, 并 `finished=true`。
- 这样一来，**所有查询**（`IndexOf`, `FindAllPrefixesOf`, ...）都会在 `r` 上做按位解析，而无需真正将节点一次性加载到内存。可在大型文件场景下省内存。

### 5.3. 文件格式概览

官方注释给出了完整格式说明，简要总结如下：

```
[0..31 bits] : 文件总大小 (单位字节)
[32..39 bits]: cbits (每个字符所需的位数)
[40..47 bits]: abits (节点地址所需的位数)
后面使用 7-bit 变长编码写: numWords, numNodes, numEdges

然后是按节点序列化的数据:
- 每个节点:
  1 bit: final
  1 bit: fallthrough?
  如果 fallthrough=1:
    写 cbits 个 bit 的字符
    ...
  否则:
    1 bit: singleEdge?
    if !singleEdge {
      7code: numberOfEdges
      读 nskiplen bits: skipfieldLen
    }
    对每个边:
      cbits: character
      if 不是第一个边:
        skipfieldLen: skipCount
      abits: 跳转目标节点地址
```

> 可以看到，这里有相当多的**条件分支**和**位数自适应**的编码方式，以追求最大化的紧凑度。

---

## 6. 位操作辅助（bitWriter / bitSeeker）

在这个实现中，有一大块逻辑专门处理“**按位读写**”。这是本代码最具特色、也最复杂的部分之一。

### 6.1. `bitWriter`

- 内部有 `cache uint8`, `used int` 表示缓存的“未写满bit数”；
- `WriteBits(data uint64, n int)`: 将 `data` 的后 `n` 位，按需要拼接到 `cache` 中，一旦 `cache` 满 8 bits，就写到底层 `io.Writer`。
- `Flush()`：把残余的没写满 8 bits 的数据补位并写出。

### 6.2. `bitSeeker`

- `p`：当前在**整个文件的 bit 偏移**（而非 byte 偏移）；
- `have` / `cache`: 维护一个 64-bit 的缓存；
- `ReadBits(n int64)`: 从 `p` 开始读 `n` 个bit，可能要跨过多次 64-bit 缓冲；
- `Seek(offset int64, whence int)`: 移动 bit 偏移；
- `Tell()`: 返回当前的 bit 偏移；
- 通过 `readUint32()` 等函数，我们可获取文件大小之类信息。

### 6.3. 7-bit 编码（`writeUnsigned` / `readUnsigned`）

- 为了存储如 `numAdded`, `numNodes`, `numEdges` 这些可能有很大范围的整数，但又希望**可变长**：
  - 采用一种“7-bit”格式，每字节只用 7 bits 来表示数值，第 8 bit 用做“是否继续”标记。
  - 这与常见的“varint”或 Protobuf 中的“变长整型”思路类似。

---

## 7. 小结与扩展

### 7.1. 架构要点

1. **构建期**（Builder 模式）

   - 要求**严格按字典序**添加单词；
   - 利用 `uncheckedNodes` 和 `minimizedNodes` 实现在线最小化；
   - 完成后调用 `Finish()`。

2. **完成后**（Finder 模式）

   - 不再保留全部 `nodes`；只读存储在 `io.ReaderAt`；
   - 每次查询都通过 `bitSeeker` 按需解析节点；
   - 大大降低了常驻内存消耗。

3. **跳过计数 (skip count)**

   - 每条边记录可达单词数目累计，用以快速定位单词索引；
   - 支持 `IndexOf`、`AtIndex` 这样需要根据顺序查找的操作。

4. **按位序列化**
   - 通过精心设计的文件格式，`cbits`、`abits` 等自适应位宽；
   - 变长 `7-bit` 编码；
   - 细致的条件写法（`fallthrough`、`singleEdge` 等），最大化节省空间。

### 7.2. 优势与应用

- **优势**

  1. 对于超大规模字典（百万甚至千万级单词），依赖 DAWG 最小化可显著减少冗余；
  2. 只读查询时无需加载整个结构，能在内存受限或文件非常大的场合使用；
  3. 可快速进行各种基于“索引”的操作（如 `IndexOf`, `AtIndex`）以及前缀查找。

- **不足 / 适用性**
  1. 不支持“动态插入或删除”——一旦 `Finish()`，就只能只读；要更新就得重建。
  2. 构建期要求**字典序输入**，若无法满足这一点，则需要先排序；
  3. 位操作与自定义格式使得实现相对复杂，调试或移植成本较高。

---

### 总体结论

这套代码把 **DAWG**（最小化 Trie） 的思想发挥到极致，不仅在线最小化，还支持**跳过计数**与**位级序列化**，做到了在**大规模单词集合**中的高效存储与只读查询。

- 在“构建”阶段，典型的 Trie + 后缀最小化流程被封装为 `Add / minimize / Finish`；
- 在“只读”阶段，通过 `bitSeeker` 在文件或内存中按需解析节点，极其紧凑且可高效完成 `IndexOf`、`FindAllPrefixesOf`、`AtIndex`、`Enumerate` 等操作。

如果你想要一个能**批量处理海量词典**、同时要**最小化存储**并且**只读查询**的结构，这段 **DAWG** 代码示例可提供很好的思路与实现参考。
