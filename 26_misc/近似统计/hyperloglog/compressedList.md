下面详细讲解 **`compressedList`** 数据结构在这段 HyperLogLog 实现中的用途、内部原理和关键代码逻辑。它是用于**稀疏模式**下保存哈希值（编码后）的重要部分，通过**差分 + 可变长**的方式压缩存储大量 32-bit 整数。

---

# 1. 结构定义与字段含义

在代码中可以看到：

```go
type compressedList struct {
    count uint32
    last  uint32
    b     variableLengthList
}
```

- **`count`**：当前存储的元素数量（逻辑上插入了多少个数）。
- **`last`**：记录最后插入的元素值（未差分前的原始值），帮助在差分编码/解码时进行累加。
- **`b`**：是一个自定义的 `variableLengthList`，底层是一个 `[]uint8` 切片，用来以**可变长编码（varint-like）**的方式存储每个元素**与上一个元素**的“差值”（即 `delta = x - last`）。

所以，`compressedList` 的本质是：

1. **维护一个递增序列**（通常是已排序的 32-bit 整数，且不会乱序插入）；
2. 在插入每个新数 `x` 时，存储 `delta = x - last` 的可变长编码到 `b` 里，并更新 `last = x`；
3. `count` 跟踪插入了多少个数。

---

# 2. 差分可变长编码 (variable-length encoding)

为了节省空间，`compressedList` 不直接存储 32-bit 整数，而是**只存储相邻元素差值**，并用**可变长**方式对差值进行编码。

这一技巧在很多地方都用过，比如 [Protocol Buffers 中的 Varint 编码](https://developers.google.com/protocol-buffers/docs/encoding)、[压缩索引/日志的差分存储](https://en.wikipedia.org/wiki/Variable-length_quantity) 等，都能显著减少存储开销。

- **差分**：假设现在已插入过一个元素 `last=100`，当你要插入新的 `x=108`，则只需要存储 `108 - 100 = 8`；如果下一次插入 `x=110`，只需存 `2`，依次类推。
- **可变长**：如果差值小（比如 1~127），只占 1 个字节；如果差值很大（超 127，需要更多字节），就会继续写第二字节、第三字节，直到写完差值的所有 7-bit 块。

代码里的可变长编码实现见 `variableLengthList.Append(x uint32)` 和 `decode(i int, last uint32) (uint32, int)` 等函数：

```go
func (v variableLengthList) Append(x uint32) variableLengthList {
    // 当 x 超过 7 bit，就在高位继续分割
    // 最后以一个 0x7F 范围内的字节结束
    for x & 0xffffff80 != 0 {
        v = append(v, uint8((x & 0x7f) | 0x80))
        x >>= 7
    }
    // 剩余的低 7 bit
    return append(v, uint8(x & 0x7f))
}
```

- 若 `x <= 127`, 只写 1 字节；
- 若大于 127, 则拆分为多个 7-bit 块，每个块最高位 `0x80` 表示“还有后续字节”。直到最后一个块把最高位清零 (`0x80` 位置为 0)。

然后在读取时（decode）要做相反的操作：

```go
func (v variableLengthList) decode(i int, last uint32) (uint32, int) {
    // 从 v[i] 开始，读连续字节，每个字节低 7 位拼起来，高位 bit=1 表示下一字节
    var x uint32
    j := i
    for ; v[j] & 0x80 != 0; j++ {
        x |= uint32(v[j] & 0x7f) << (uint(j - i) * 7)
    }
    // 处理最后一个字节
    x |= uint32(v[j]) << (uint(j - i) * 7)

    // x 是差值，恢复原值需加上 last
    return x + last, j + 1
}
```

- 上述函数返回 `(新元素值, 下一次读取起点)`。

---

# 3. 插入新元素 (`Append`)

在 `compressedList` 层面，每次插入一个 32-bit 整数 `x`，做的操作是：

```go
func (v *compressedList) Append(x uint32) {
    v.count++                     // 记录插入次数+1
    v.b = v.b.Append(x - v.last) // 差分 (x - last) 做可变长编码
    v.last = x                   // 更新 last
}
```

也就是把 `(x - last)` 的可变长编码写到 `b`。这样如果顺序插入的 `x` 们是递增且相近，就能非常高效地压缩。

---

# 4. 解码（遍历）

当需要遍历或读取 `compressedList` 中所有元素时，会用到迭代器 `Iter()`：

```go
func (v *compressedList) Iter() *iterator {
    return &iterator{0, 0, v}
}
```

其中 `iterator` 结构：

```go
type iterator struct {
    i    int        // 当前在 variableLengthList 的索引
    last uint32     // 用于累计恢复原值
    v    iterable   // 这里 v = compressedList
}

func (iter *iterator) Next() uint32 {
    n, i := iter.v.decode(iter.i, iter.last)
    iter.last = n
    iter.i = i
    return n
}

func (iter *iterator) Peek() uint32 {
    n, _ := iter.v.decode(iter.i, iter.last)
    return n
}

func (iter iterator) HasNext() bool {
    return iter.i < iter.v.Len()
}
```

- 每次 `Next()` 调用都会：
  1. 用 `decode(iter.i, iter.last)` 把当前差分片段解析出真正的值 `n`；
  2. 更新 `iter.last = n`；
  3. 更新 `iter.i` 到下一位置 `i`；
  4. 返回新的真实数 `n`。
- 这样我们就能顺序迭代刚才存进去的所有 32-bit 元素。

---

# 5. 与 HyperLogLog 稀疏模式的关系

在本实现中，`Sketch` 在稀疏模式下会将哈希值**编码**后（即 `encodeHash()` 得到 32-bit 整数），存到两个地方：

- **`tmpSet`**：一个普通的 Set，用来快速去重并存储尚未合并的哈希值；
- **`sparseList`**：也就是 `compressedList`，这里做了“有序合并”和“差分可变长编码”，适合批量存储大量哈希值。

当 `tmpSet` 达到一定规模时，就会 **mergeSparse()**：

- 将 `tmpSet` 里的元素取出来排序，与 `sparseList` 做 **归并**，结果新的 `sparseList` 就是有序合并后的数据；
- 然后清空 `tmpSet`；
- 一旦 `sparseList` 依然很大时（超过 `m`），就会转为稠密模式（`toNormal()`）。

因此，`compressedList` 主要是让大量的哈希值在**稀疏阶段**得到紧凑的存储，以节省内存和加快操作。

---

# 6. (反)序列化

`compressedList` 也提供了 `MarshalBinary` / `UnmarshalBinary`，以便在 HyperLogLog 整体序列化时能把稀疏列表一并保存/恢复。

```go
func (v *compressedList) MarshalBinary() (data []byte, err error) {
    // 先序列化 v.count 和 v.last (各4字节)
    // 再序列化 v.b (variableLengthList), 其自身也带长度信息
    ...
}

func (v *compressedList) UnmarshalBinary(data []byte) error {
    // 读出 count, last
    // 读出 v.b 的大小, 再读切片
    ...
    return nil
}
```

读取时就能恢复该列表中的所有差分编码。后续在 HyperLogLog 的 `Estimate()` 或 `mergeSparse()` 阶段，如果需要真正遍历所有元素，才会通过 `Iter()` 一个个 decode 出来。

---

# 7. 优点 & 总结

**`compressedList`** 利用“**差分** + **可变长**”存储来保存一系列有序 32-bit 整数，在 HLL 的稀疏模式中非常有用，主要优势：

1. **节省空间**：大部分哈希值的相邻差距可能较小，因此绝大多数差值可以在 1~2 个字节内就能表达，远比直接存 4 字节要节省。
2. **顺序读取**：设计了迭代器来顺序 decode，写操作也只在末尾追加，不需要复杂的中间插入操作；
3. **有序特性**：在 `mergeSparse()` 阶段，会做一次有序归并，为进一步切换到稠密模式（或下次合并）提供方便。

因此在本实现中，`compressedList` 既是 HyperLogLog 稀疏模式下的一种**持久存储**方式，也能**序列化**到磁盘或网络传输并在后续解码使用。这种技巧（差分 + varint）对于任何 **单调递增序列** 的存储都十分常见且高效。

---

## 综上

- **`compressedList`** 是 HyperLogLog 稀疏模式的核心数据结构，用**差分 + varint**压缩存储一批有序的 32-bit 整数。
- 每次 `Append(x)` 会把 `(x - last)` 编码追加进 `b`，并更新 `last`、`count`；
- 遍历时 (`Iter().Next()`) 通过解码把差值还原成真正的整数。
- 在稀疏模式下，`compressedList` 配合 `tmpSet` 能够在小数据量阶段减少内存占用，延迟真正的稠密模式分配。
- 当数据量过大时，代码会将 `compressedList` 里的内容解码并写到稠密数组 `regs`，从而完成模式切换。
