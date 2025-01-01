下面这段代码来自 GitHub 上的一个 Cuckoo Filter 实现（[seiflotfy/cuckoofilter](https://github.com/seiflotfy/cuckoofilter)），它演示了如何在 Go 语言中构建并使用布谷鸟过滤器（Cuckoo Filter）以及可扩展的 Scalable Cuckoo Filter。下面我会按照核心模块来进行详细的讲解。

---

## 目录

1. **main 函数（使用示例）**
2. **Filter 结构（基本的 Cuckoo Filter）**
3. **ScalableCuckooFilter（可扩展的 Cuckoo Filter）**
4. **bucket 结构（存储指纹）**
5. **metroHash（哈希函数）**
6. **辅助函数**

---

## 1. main 函数（使用示例）

```go
func main() {
    cf := NewFilter(1000)             // 创建一个容量大约为 1000 的 Cuckoo Filter
    cf.InsertUnique([]byte("geeky ogre"))

    // 查询是否存在（查 "hello" 不存在）
    cf.Lookup([]byte("hello"))

    count := cf.Count()
    fmt.Println(count) // count == 1

    // 尝试删除 "hello"（其实并不存在）
    cf.Delete([]byte("hello"))

    count = cf.Count()
    fmt.Println(count) // count == 1 (无影响)

    // 删除真正存在的 "geeky ogre"
    cf.Delete([]byte("geeky ogre"))

    count = cf.Count()
    fmt.Println(count) // count == 0

    // 重置整个过滤器
    cf.Reset()
}
```

- `NewFilter(1000)`：初始化一个 Cuckoo Filter，目标容纳 1000 个元素。内部会取大于等于 `1000` 的**下一个 2 的幂**（并除以桶大小 4），从而分配实际的桶数组大小。
- `InsertUnique()`：将字符串插入过滤器，如果已经存在则返回 false；否则插入并返回 true。
- `Lookup()`：判断某个字符串是否可能存在。若返回 false，说明一定不存在；若返回 true，则“极大概率存在”（极小概率假阳性）。
- `Delete()`：删除某个元素。若元素的指纹在桶中存在，则真正删除并更新计数；否则什么也不做。
- `Count()`：返回过滤器中元素（指纹）计数（大概的）。
- `Reset()`：重置过滤器，清空所有桶。

---

## 2. Filter 结构（基本的 Cuckoo Filter）

### 2.1 Filter 数据结构

```go
type Filter struct {
    buckets   []bucket // 存储指纹的桶数组，每个桶固定大小 4
    count     uint     // 当前过滤器中（可能）存储的元素总数
    bucketPow uint     // 记录桶数量是 2^bucketPow，方便位运算
}
```

- `buckets []bucket`：实际存储指纹的容器，每个 `bucket` 大小为 4（见后面 `bucketSize = 4`）。
- `count`：Cuckoo Filter 中存储的元素计数，用于统计和简单判断负载。
- `bucketPow`：因为桶的总数为 2 的幂，所以用 `bucketPow` 表示 `\(\log_2(\text{桶数量})\)`；在查找插入时可以用位运算快速取模。

### 2.2 NewFilter(capacity uint)

```go
func NewFilter(capacity uint) *Filter {
    // 先获得大于等于 `capacity` 的 2^x，并除以 bucketSize 4
    capacity = getNextPow2(uint64(capacity)) / bucketSize
    if capacity == 0 {
        capacity = 1
    }
    buckets := make([]bucket, capacity)
    return &Filter{
        buckets:   buckets,
        count:     0,
        bucketPow: uint(bits.TrailingZeros(capacity)),
    }
}
```

- `getNextPow2()`：找下一个 `2^x` 方便做“取模”操作（用掩码来替代 `%`）。
- 除以 `bucketSize` 是因为每个桶可以容纳 4 个指纹，相当于把容量分摊到各桶中。
- `bits.TrailingZeros(capacity)`：找 `capacity` 的最低位 1 的索引位置，这个值就是 \(\log_2(\text{capacity})\)。

### 2.3 Lookup(data []byte) bool

```go
func (cf *Filter) Lookup(data []byte) bool {
    i1, fp := getIndexAndFingerprint(data, cf.bucketPow)
    if cf.buckets[i1].getFingerprintIndex(fp) > -1 {
        return true
    }
    i2 := getAltIndex(fp, i1, cf.bucketPow)
    return cf.buckets[i2].getFingerprintIndex(fp) > -1
}
```

- 首先对 `data` 做哈希，得到 `(i1, fp)`：
  - `i1`：桶索引 1
  - `fp`：指纹（1 个字节）
- 如果 `i1` 号桶中含有该指纹，就返回 `true`；否则计算出它的替代桶 `i2`（替代桶索引常见公式：`i2 = i1 XOR hash(fp)`），检查是否有同样的指纹。
- 若两个桶都没有，则返回 `false`。

### 2.4 Insert(data []byte) bool

```go
func (cf *Filter) Insert(data []byte) bool {
    i1, fp := getIndexAndFingerprint(data, cf.bucketPow)
    if cf.insert(fp, i1) {
        return true
    }
    i2 := getAltIndex(fp, i1, cf.bucketPow)
    if cf.insert(fp, i2) {
        return true
    }
    return cf.reinsert(fp, randi(i1, i2))
}
```

1. 计算指纹和第一个桶索引 `i1`；尝试往 `i1` 插入指纹。
2. 若失败则去替代桶 `i2`（由 `getAltIndex()` 算出），再尝试插入。
3. 如果两个桶都满了，就使用 `reinsert()` 进行布谷鸟式的踢出过程。

#### insert(fp fingerprint, i uint) bool

```go
func (cf *Filter) insert(fp fingerprint, i uint) bool {
    if cf.buckets[i].insert(fp) {
        cf.count++
        return true
    }
    return false
}
```

- 尝试将指纹插入到 bucket（若有空位则插入成功），成功则 `count++`。

#### reinsert(fp fingerprint, i uint) bool

```go
func (cf *Filter) reinsert(fp fingerprint, i uint) bool {
    for k := 0; k < maxCuckooCount; k++ {
        j := rand.Intn(bucketSize)
        oldfp := fp
        fp = cf.buckets[i][j]
        cf.buckets[i][j] = oldfp

        // compute the alternate location for the 'evicted' fingerprint
        i = getAltIndex(fp, i, cf.bucketPow)
        if cf.insert(fp, i) {
            return true
        }
    }
    return false
}
```

- 随机踢出桶中某个指纹（用 `rand.Intn(bucketSize)` 选定位置），然后将新指纹插入其位置；
- 对被踢出的指纹 `fp`，计算它的替代桶索引 `i = getAltIndex(...)`，再尝试插入；
- 如果循环一定次数（`maxCuckooCount = 500`）还没成功，插入失败。可以视为“过滤器已近满载”。

### 2.5 Delete(data []byte) bool

```go
func (cf *Filter) Delete(data []byte) bool {
    i1, fp := getIndexAndFingerprint(data, cf.bucketPow)
    if cf.delete(fp, i1) {
        return true
    }
    i2 := getAltIndex(fp, i1, cf.bucketPow)
    return cf.delete(fp, i2)
}
```

- 和查找类似，算出 `i1, fp`，若能在 `i1` 上找到并删除则返回 true；否则尝试 `i2`。
- 删除成功就 `count--`。

### 2.6 其他方法

- `Reset()`：清空所有桶中的指纹，并将计数 `count` 置 0。
- `Count()`：返回计数 `count`。
- `Encode()/Decode()`：序列化和反序列化，用于将 Filter 转成字节流或从字节流恢复。

---

## 3. ScalableCuckooFilter（可扩展的 Cuckoo Filter）

当我们希望过滤器容量不足时能自动扩容，而不是返回插入失败，就可以用 **ScalableCuckooFilter**。它内部维护一个 `[]*Filter`，每个 `Filter` 可以看作一层，当某一层（最后一个 Filter）负载过高，则创建一个新的、容量更大的 Filter，继续存储后续元素。

```go
type ScalableCuckooFilter struct {
    filters    []*Filter
    loadFactor float32
    scaleFactor func(capacity uint) uint
}
```

- `filters`：多个 `Filter` 组成的层级；
- `loadFactor`：可接受的负载因子，比如 0.9；当最后一层超过这个负载时，就扩容；
- `scaleFactor`：扩容策略函数，如“当前大小 × 2”之类。

### 3.1 关键方法

- `NewScalableCuckooFilter(opts ...option)`：创建可扩展过滤器，默认用初始容量为 `DefaultCapacity = 10000` 的 Filter。
- `Insert(data []byte)`：
  1. 先尝试插入到最后一个 Filter；
  2. 如果负载超限或插入失败，就扩容（新建一个更大的 Filter 并添加到 `filters` 里），再插入。
- `Lookup(data []byte) bool`：按顺序在 `filters` 列表里查找，**只要有一个 Filter 返回 true，就返回 true**。
- `Delete(data []byte) bool`：依次在所有 Filter 中尝试删除，若找到就停。
- `Count()`：把所有 Filter 的 `count` 加和。

这样就可以在运行时无限扩容，不用重建整个结构。

---

## 4. bucket 结构（存储指纹）

```go
type fingerprint byte
type bucket [bucketSize]fingerprint

const (
    nullFp     = 0
    bucketSize = 4
)
```

- 每个 `bucket` 能存 4 个指纹（`fingerprint` 1 字节）。
- `insert(fp) bool`：找空位（`nullFp`）插入，若成功返回 true，否则 false。
- `delete(fp) bool`：匹配到则置空，返回 true。
- `getFingerprintIndex(fp)`：返回指纹所在的下标（否则 -1）。
- `reset()`：清空整个 bucket。

---

## 5. metroHash（哈希函数）

```go
func Hash64(buffer []byte, seed uint64) uint64 {
    // ...
}
```

- 这是一个 **非加密、高速哈希** 的实现，类似 CityHash、xxHash、MurmurHash 等。
- 在这个项目里，`Hash64` 用作 Cuckoo Filter 的主哈希函数来计算 `i1` 和指纹。
- 注意：这里并非真正的 MetroHash 官方实现，而是 repo 主根据 MetroHash 思路写的简化版。

```go
func getAltIndex(fp fingerprint, i uint, bucketPow uint) uint {
    mask := masks[bucketPow]
    hash := altHash[fp] & mask
    return (i & mask) ^ hash
}
```

- `getAltIndex`：根据当前索引 i1、指纹 fp，算出替代索引 i2。
- `altHash[fp]` 预先存好 `fp` (0–255) 对应的哈希值。
- 这里用异或来获得与 i1 不同的桶索引 i2。

```go
func getIndexAndFingerprint(data []byte, bucketPow uint) (uint, fingerprint) {
    hash := defaultHasher.Hash64(data)
    fp := getFingerprint(hash)
    i1 := uint(hash>>32) & masks[bucketPow]
    return i1, fingerprint(fp)
}
```

- 将 64 位哈希拆分成两部分：
  - 高 32 位：用于桶索引 `i1`
  - 低 8 位：用于指纹 `fp`（并且还会 mod 255 以防 0）
- `masks[bucketPow]` 用于“取模”，因为桶数为 2^bucketPow。

---

## 6. 辅助函数

### 6.1 getNextPow2(n uint64) uint

```go
func getNextPow2(n uint64) uint {
    n--
    n |= n >> 1
    n |= n >> 2
    n |= n >> 4
    n |= n >> 8
    n |= n >> 16
    n |= n >> 32
    n++
    return uint(n)
}
```

- 经典做法，把数字变成下一次方：
  1. `n--`
  2. 然后不断向右移位并或等操作，把低位的 1 “扩展”到更高位。
  3. 最后 `n++` 就得到大于或等于原来的 2^x。

### 6.2 Encode() / Decode()

- 用于序列化和反序列化 `Filter` 或 `ScalableCuckooFilter`：
  1. `Encode()`：遍历每个 bucket，每个位置（4 个指纹），存到一个 byte slice 里。
  2. `Decode()`：从 byte slice 恢复出一个 `Filter`。
  3. `ScalableCuckooFilter.Encode()`：会把内部的多个 `Filter` 转成二维 byte slice，然后用 `gob` 序列化。
  4. `DecodeScalableFilter()`：反向恢复。

---

## 总结

1. **核心原理**

   - **Cuckoo Filter** 基于布谷鸟哈希思想：每个元素映射到两个可能桶中，若桶满则随机踢出一个已存在指纹到它的替代桶，反复直到插入成功或踢出次数过多。
   - 通过仅存储“指纹”而非完整元素，可以节省大量空间。

2. **主要流程**

   - **Insert**：先尝试两个桶；若都满则“踢出”已有指纹。
   - **Lookup**：只要任意一个桶包含该指纹，就返回 true。
   - **Delete**：若找到指纹就将其移除。
   - **Count**：维护简单的计数来统计大致元素数量。

3. **ScalableCuckooFilter**

   - 当一个 Filter 的负载超限或插入失败次数过多，就新建一个更大的过滤器接力，以实现**动态扩容**。

4. **哈希函数**
   - 使用了一个类似 MetroHash 的高速函数，适合非安全场景下的快速哈希。

在实际项目中，这段代码给出了一个相对精简的 Cuckoo Filter 实现示例。如果要在生产环境中使用，建议根据业务需求来**选择指纹长度**、**桶大小**、**扩容策略**、**并发安全**等进行更深度的优化或改造。
