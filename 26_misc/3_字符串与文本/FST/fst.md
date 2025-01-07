下列代码来自 [blevesearch/vellum](https://github.com/blevesearch/vellum) 项目，它实现了一个基于 **FST（Finite State Transducer）** 的键值存储结构。这个结构可以在插入时（Build 阶段）将一组 key（字节序列）与其对应的 value（`uint64`）写入到一个输出流中，最后形成一个紧凑的、有序的、可持久化/可内存映射的 FST。后续在查询阶段（Load/Open 阶段），可以在不加载全部内容到内存的情况下，通过 `Contains()` / `Get()` / `Iterator()` 等操作来检索数据。

以下将逐步为你**详细分析**整个代码的结构与核心原理。

---

## 1. 概览

1. **代码文件顶部**：有一个简单的示例 `main()` 函数，演示了如何构建 FST 并进行查询操作。
2. **核心接口**：
   - **Builder**：用于构建 FST，将 keys/values 按**字典序**插入，并最终 `Close()` 输出一个**有序的** FST。
   - **FST**：加载后的有限状态机/转移器，可进行 `Get()`、`Contains()`、`Iterator()` 等查询操作。
3. **代码分块**（用注释 `// #region ... // #endregion` 隔开）：
   - `builder`：实现了构造 FST 的主要逻辑；
   - `fst`：FST 的核心结构、查询操作；
   - `encoder` / `decoder`：实现 FST 在字节流中的序列化/反序列化；
   - `iterator`：实现区间/顺序遍历 (key-range iteration)；
   - `merge_iterator`：提供将多个迭代器的结果合并生成一个新的 FST 的能力；
   - 其它辅助数据结构：`registry`, `builderNode`, `automaton` 等。

本质上，这是一套**增量式**（但要求已排序）地生成 FST，并支持**静态查询**的库。它支持将生成结果写到文件里，然后再 `mmap`（内存映射）或用 `Load()` 方法将其加载到内存中。

---

## 2. main() 中的示例

```go
func main() {
    var buf bytes.Buffer
    builder, err := New(&buf, nil)
    // ...
    err = builder.Insert([]byte("cat"), 1)
    err = builder.Insert([]byte("dog"), 2)
    err = builder.Insert([]byte("fish"), 3)
    err = builder.Close()

    fst, err := Load(buf.Bytes())

    val, exists, err := fst.Get([]byte("dog"))
    if exists {
        fmt.Printf("contains dog with val: %d\n", val)
    }
}
```

- **Builder**：`New(&buf, nil)` 创建一个构造器，把构造输出写到 `&buf`（内存中的 `bytes.Buffer`）。
- **Insert**：按**字典序**插入 `"cat" < "dog" < "fish"`；若顺序不对，Builder 会返回错误（`ErrOutOfOrder`）。
- **Close()**：结束构造过程，写入 FST 的“footer”/收尾信息。
- **Load(buf.Bytes())**：从刚刚写入的二进制数据加载 FST，到内存里；
- **Get("dog")**：检查键 `"dog"` 是否存在，若存在则取出对应的 `uint64` 值。

---

## 3. FST 数据结构

```go
type FST struct {
    f       io.Closer
    ver     int
    len     int
    typ     int
    data    []byte
    decoder decoder
}
```

- `data []byte`：FST 的全部序列化内容（若是 `Open()` 大文件场景，可能使用 mmap，但此处简化为把数据都放进 `[]byte`）。
- `decoder`：版本相关的解码器，用于对 `data` 解析出各个状态/转移信息。
- `ver`, `typ`：编码版本、类型；在 `decodeHeader()` 中读出。
- `len`：FST 中存储的 key-value 条目的数量。

### 3.1 加载与查询

- `Load(data []byte) (*FST, error)`：把 `[]byte` 解析为 `FST`；`f.ver, f.typ = decodeHeader(data[:16])` 后调用 `loadDecoder(...)` 等初始化操作。
- `Get(input []byte) (uint64, bool, error)`：遍历状态机以找对应路径：
  1. 从 `decoder.getRoot()` 取根状态；
  2. 对每个字符 `c` 调用 `state.TransitionFor(c)`，若没找到则返回 `false`; 若找到则累计 “转移输出” 到 `total` 中；
  3. 最终若到达的状态是 `Final()`，再加上 `FinalOutput()`；
  4. 返回 `(total, true, nil)`，否则 `false`。

`Contains()` 函数只是 `Get()` 的简化调用（不需要 value，仅判断存在性）。

---

## 4. Builder 构造逻辑

**Builder** 用来“流式”地写 FST。主要字段：

```go
type Builder struct {
    unfinished *unfinishedNodes
    registry   *registry
    last       []byte
    len        int
    lastAddr   int
    encoder    encoder
    // ...
}
```

- `unfinished`：暂存尚未编译成“冻结”节点的那部分前缀路径；
- `registry`：用来去重/合并已经“冻结”的节点——因为多个分支可能合并成相同后缀；
- `encoder`：将节点写入输出流的策略（这里分为多版本，主要是 `encoderV1`）；
- `last`：记录上一次插入的 key，以检查新插入 key 是否按字典序递增。

### 4.1 插入

```go
func (b *Builder) Insert(key []byte, val uint64) error {
    // 1) 检查 key 是否比 b.last 更大，否则 ErrOutOfOrder
    // 2) 计算与旧 key 的公共前缀 prefixLen
    //    并做一些 Output 的差值处理
    // 3) compileFrom(prefixLen) 关闭 prefixLen 之后的分支
    // 4) 把剩余 suffix 的字符一一挂上 (unfinished.addSuffix...)
    // 5) b.len++
}
```

- **公共前缀**：若新插入 key 与上个 key 共享前缀，Builder 会把共享前缀之前的那些分支保留在栈中，只对分叉的节点做编译/冻结。
- **registry** 用来哈希合并：相同的“后缀”子树只需存一次。

### 4.2 完成与 Close

```go
func (b *Builder) Close() error {
    // compileFrom(0) —— 编译并冻结所有未完成节点
    // popRoot() —— 弹出根节点并 encode
    // encoder.finish(...) 写入 footer
    // ...
}
```

---

## 5. 编码与解码 (encoder / decoder)

**Vellum** 中定义了**多版本**的编码方案，这里实际用 `versionV1`。

### 5.1 头部与尾部

- **Header**（16 字节）：前 8 字节是版本号 `ver`，后 8 字节是 `typ`（用 0 表示 FST ？）。
- **Footer**（16 字节）：前 8 字节存 `count`（有多少 key-value），后 8 字节存 `rootAddr`（根状态在文件中的偏移地址）。

### 5.2 `encoderV1`：如何把构建好的状态写到流

- 一个“状态”里包含：
  1. 是否 `final`，以及可能的 `finalOutput`；
  2. 转移列表：若只有 1 个转移，就用**单转移**的特殊压缩（single transition）。若多个，就用“multi transition”格式。
  3. “deltaAddr” 存储下一状态的相对偏移，以减少数据空间。

核心函数如 `encodeStateOne()` / `encodeStateMany()` / `encodeStateOneFinish()`。

### 5.3 `decoderV1`：如何从二进制读取状态

```go
type decoderV1 struct {
    data []byte
}
```

- `getRoot()`：footer 中读 `rootAddr`;
- `stateAt(addr int, prealloc fstState)`: 解析在 `data` 里下标 `addr` 的状态；判断是单转移或多转移等。
- `TransitionFor(b byte)`：在多转移时，会在 `[transBottom, transTop]` 里搜字符 b；若找到则解析出 `dest` 和 `out`。
- `Final()`/`FinalOutput()`：若该状态是 final，就读取相应的 final output。

---

## 6. 迭代器 (Iterator)

`FSTIterator` 实现了对 key 的有序遍历：

1. `f.decoder.stateAt(f.decoder.getRoot(), ...)` 得到 root；
2. `startKeyInclusive` / `endKeyExclusive`：限制迭代范围；
3. `Next()`：沿着当前状态的转移关系前进，若没有更多转移就回溯；最终输出下一个可用的 final 状态对应的 key。
4. `Current()`：取得当前 key + 累计 output。

它结合**自动机**(`aut Automaton`) 可以进行更复杂的匹配（如模糊搜索），但默认是 `alwaysMatchAutomaton`（全部接受）。

---

## 7. Merge 功能

`Merge()` 函数允许将多个 `Iterator` 合并后再建新的 FST。也就是边遍历多个已排序的 `Iterator`，边把合并后的 key-value 写到新的 `Builder` 中：

```go
func Merge(w io.Writer, opts *BuilderOpts, itrs []Iterator, f MergeFunc) error {
    // 1) New builder
    // 2) NewMergeIterator(itrs, f)，交给 MergeIterator 合并多个 iterator
    // 3) 循环取 merged key/value -> builder.Insert(key,value)
    // 4) builder.Close()
}
```

- 当多个 iterator 有相同 key，会用 `MergeFunc`（如 `MergeMin`, `MergeMax`, `MergeSum`）来**合并 value**。

---

## 8. 主要数据结构与要点

1. **`unfinishedNodes`**：插入过程中，用一个 stack 结构维护尚未“冻结”的节点（BuilderNode）。
2. **`builderNode`**：表示一个 Trie 节点/状态，包含 `final`, `finalOutput` 以及 `trans`（转移）。
3. **`registry`**：一个哈希表，用来判断**相同**的“冻结节点”是否已存在（从而复用地址），实现**后缀合并**。
4. **`writer`**：内部封装 `bufio.Writer`，并维护写入计数 `counter`，用来知道当前写到文件/内存流的第几个字节。
5. **`packedSize/encodePackSize`**：为了尽量节省空间，对 `uint64` 做紧凑存储（小数字只用 1~2 字节，大数字才用更多字节）。

---

## 9. 整体流程串起来

**构建**阶段（Builder）：

1. 调用 `New(io.Writer, opts)`；
2. 多次调用 `Insert(key, value)`（键必须按字典序递增）；
   - 维护“未完成节点”（unfinished），对公共前缀进行复用；
   - 如果新 key 导致分叉过多，会将先前的节点“冻结”，写到输出流里；
   - 用 `registry` 判断是否已存在等价节点，若是则直接复用地址；
3. 调用 `Close()`：把根节点也写到流里，最后写入“Footer”，结束。

**查询**阶段（FST）：

1. 调用 `Load([]byte)` 或 `Open(filepath)` 读取 FST 数据；
2. `Get(key)`：从 root 开始按字符找转移，累加中间输出值与末端输出值；
3. `Iterator(...)`：从指定起始 key 开始向后遍历，直到没有下一个可用键或者超出 endKeyExclusive。

---

## 10. 关键特性与总结

1. **FST vs DFA**：一个纯粹的 DFA（确定性有穷自动机）只能识别“key 是否存在”，而 FST 在转移上附加了数值输出（`uint64`），从而在路径末尾可以得到最终的 value。
2. **写时压缩**：通过后缀合并（registry）和紧凑编码（single / multi transition、delta offset、packed uint），大幅减少存储空间，适合**大规模静态**键值对。
3. **只读**：建好后基本是只读结构；若要动态更新，需重建 FST。
4. **支持 mmap**：`Open()` 默认尝试 mmap 文件数据，这对于大规模只读词典很高效。
5. **Iterator**：可以在排序好的 key 集合上做范围遍历，常用于**前缀搜索**、区间搜索等场景（结合 `startKeyInclusive`, `endKeyExclusive`）。

---

### 小结

这份代码是一个**面向工业场景**的 FST（Finite State Transducer）实现，设计思路是：

- **静态**：键值对只需构建一次，之后高效查询；
- **紧凑**：带有多种压缩手段（delta offset、packed uint、single transition vs multiple transitions）；
- **可持久化/可 mmap**：支持大数据量不全部加载到内存；
- **可合并**：提供 Merge 功能，能将多个已排序的 iterators/fst 合并。

这让它尤其适用于搜索引擎、全文检索、分词、自动完成功能等场景，需要**大规模**字符串词典且只读的地方。插入操作需要严格的**字典序顺序**，这是为了保证构建过程可以顺序合并和压缩。

在阅读和使用该代码时，你可以重点关注：

- **Builder** 中的公共前缀处理和“冻结”节点逻辑；
- **registry** 如何对**等价节点**（final + trans 列表相同）进行哈希去重；
- **encoder/decoderV1** 如何以紧凑格式进行存储和加载；
- **Iterator** 在 FST 上做深度搜索/回溯，列举出所有可达并符合条件的 key。

这样就能理解 `vellum` 的核心原理与实现细节。希望以上详解对你理解这段 FST 代码有所帮助。

---
