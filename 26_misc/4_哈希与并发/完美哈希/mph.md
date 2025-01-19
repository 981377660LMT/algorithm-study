下面给出对 **mph（Minimal Perfect Hashing）** 库中 **CHD**（CHD = "Compressed Hashing, Dynamic" 或 “CHD minimal perfect hashing” 的一种变体）实现的**详细讲解**。这段代码演示了如何使用**随机化方法**构造一个在**key -> value** 映射中无冲突（perfect），且占用空间相对紧凑（minimal-ish）的哈希结构。你可以借助它实现对一批静态键的高效查表，以及可持久化存储(序列化到磁盘)并支持后续查询。

主要内容：

1. **核心数据结构**及构建流程( `CHDBuilder` + `CHD` )
2. **随机哈希函数**(`chdHasher`)
3. **用到的桶(buckets)与多次试验**
4. **序列化/反序列化**( `Write()`、`Read()`、`Mmap()` )
5. **Get 查询方法**
6. **Iterator 遍历**

---

## 一、概览

### 1.1 Minimal Perfect Hash (MPH)

- 给定一组**不可重复**的键 (key)，MPH 构造出一个哈希函数（或类似结构）使得**没有冲突**，同时哈希表大小正好等于键数（或近似）。
- 在本例中，代码称之为 `CHD`（Compressed Hashing with possible dynamic approach），可将 key -> value 储存在一份紧凑表中，并通过**一次性构建**获得**无冲突**的哈希索引。

### 1.2 数据结构

- **`CHDBuilder`**：用于**收集所有 (key, value)** 并构建最终的 `CHD`。
- **`CHD`**：构建完成的只读结构，可进行：

  - `Get(key)`：检索对应 value
  - `Iterate()`：遍历所有记录
  - `Write()`：序列化到 `io.Writer`
  - `Mmap()` / `Read()`：从内存片或 stream 里加载

- **核心思路**：
  1. 对 `n` 个 key 用一个**大哈希函数**( `HashIndexFromKey` ) 把它们分配到 `m` 个桶(buckets)；
  2. 依桶(buckets)进行处理：在一个**全局**哈希表(`size = n`)里寻找对本桶内 keys 都能无冲突写入的位置( `Table(r, key)` )；
  3. 如出现冲突则换另一随机哈希直到成功。
  4. 最终对每个桶确定一个随机值 `r`；用这个 `r` + 全局 `r[0]` 来把 key->索引映射到无冲突位置。

---

## 二、建表流程：`CHDBuilder`

### 2.1 `CHDBuilder` 结构

```go
type CHDBuilder struct {
	keys   [][]byte
	values [][]byte
	seed   int64
	seeded bool
}
```

- 存放待构建的所有 `(key, value)` 对，以及可选的随机数种子 `seed`。

### 2.2 接口

- `Builder() *CHDBuilder`：创建空 builder。
- `Seed(seed int64)`：指定随机种子(否则默认为当前时间 `UnixNano()`)。
- `Add(key, value []byte)`：把一对 KV 加入队列里。
- `Build() (*CHD, error)`：构建最终 `CHD` 对象，若发现重复 key 或无法找到无冲突哈希，报错。

### 2.3 `Build()` 主逻辑

```go
func (b *CHDBuilder) Build() (*CHD, error) {
    n := uint64(len(b.keys))
    // 1) m = n/2, 也可做别的策略
    m := n/2
    if m == 0 { m = 1 }

    // 2) keys, values: 全局索引 -> (k,v),大小 = n
    keys := make([][]byte, n)
    values := make([][]byte, n)

    // 3) newCHDHasher(n, m, seed, seeded) => 构造随机hash对象
    hasher := newCHDHasher(n, m, b.seed, b.seeded)

    // 4) bucketVector: 大小 m, each bucket has index + slice of (keys, values)
    buckets := make(bucketVector, m)
    indices := make([]uint16, m) // for each bucket => which r index we used?
    // seen map => 记录哪些全局位置已被占用
    seen := make(map[uint64]bool)

    // 5) 先按照 hasher.HashIndexFromKey(key) 把所有key分到 buckets
    //    "hashIndexFromKey" => ( FNV64(key) ^ hasher.r[0] ) % m
    ...
    // 6) sort buckets by size (descending). 大桶优先处理
    sort.Sort(buckets)

    // 7) For each bucket: try existing or new random function r => "tryHash"
    ...
    // 8) if success => fill keys[]/values[] with those positions
    //    if no success after N tries => error

    // 9) build CHD object => 先把 keys[][]byte, values[][]byte 内容写到 mmap[] bytes + 记录 slices
    ...
    return &CHD{ ... }, nil
}
```

- **要点**：
  - `m = n/2`：把 n 个键 roughly 分到 n/2 个桶(每桶预期 2 个键，但可能不均)。
  - 先**“outer hash”**( `HashIndexFromKey` ) 将 key 分桶；再**“inner hash”**( `Table(r, key)`) 为桶内 key 找到**全局**位置(0..n-1)存储。
  - `sort.Sort(buckets)` => 让大桶(包含更多 key)先处理(更难找到无冲突的 r，需要优先处理)
  - **tryHash** => 给定一个随机数 `r`, 对桶内 key 做 `Table(r,k)` = ((FNV64(k) ^ r[0] ^ r) mod n ) 看是否与 `seen[]` 冲突；若不冲突就写入 `keys[]`+`values[]` 位置。

### 2.4 `tryHash` 函数

```go
func tryHash(hasher *chdHasher, seen map[uint64]bool,
             keys [][]byte, values [][]byte, indices []uint16,
             bucket *bucket, ri uint16, r uint64) bool {
    // 1) 计算hashes[] = for each k in bucket => hasher.Table(r, k)
    // 2) 若已在 seen[] => 冲突
    // 3) no conflict => update seen, indices[bucket.index]=ri
    //                  update keys[hashes[i]] = bucket.keys[i]
    //                           values[hashes[i]] = bucket.values[i]
    return true/false
}
```

- `ri` 是**hash函数编号**(也即在 `hasher.r` 里的索引)；
- `r` 是**随机值**；
- `bucket.index` 是外层 `m` 中桶的索引，用来记录在 `indices[]` 里——表示**这个桶最终使用了 hasher.r[ri]** 这个随机值。

---

## 三、`chdHasher`

```go
type chdHasher struct {
	r       []uint64 // array of random seeds
	size    uint64   // n = # keys
	buckets uint64   // m = #buckets
	rand    *rand.Rand
}
```

- `r[0]` 用于**outer hash**: `(FNV64(k) ^ r[0]) % m` => which bucket
- 其余 `r[i]` 用于**inner hash**: `(FNV64(k) ^ r[0] ^ r[i]) % size`

### 3.1 `newCHDHasher` + `Add(r)`

- `newCHDHasher()` 初始化 `c.r = []uint64{ c.rand.Uint64() }` => `r[0]`
- `Add(r)` 把新的随机值 append 到 `r[]` 用于**新的** inner hash function。

### 3.2 `Generate()`

- 返回 `(uint16(len(r)), c.rand.Uint64())` => (ri, rVal)
- 也就是**新的** `rVal` + 索引 = `ri`

---

## 四、CHD 完成后：`CHD` 类型

- `r []uint64`: hash seeds; `indices []uint16`: for each bucket => which r index used
- `mmap []byte`: 存放**合并**后的 keys & values 数据
- `keys` `values` : slice of (start,end) 指向 `mmap[]` 中具体存储

### 4.1 `Get(key)`

```go
func (c *CHD) Get(key []byte) []byte {
  r0 := c.r[0]
  h := hasher(key) ^ r0
  i := h % uint64(len(c.indices)) // which bucket
  ri := c.indices[i]
  if ri >= uint16(len(c.r)) { return nil }
  r := c.r[ri]
  ti := (h ^ r) % uint64(len(c.keys)) // the final position in [0..n-1]
  kSlice := c.keys[ti]
  if bytes.Compare(c.slice(kSlice), key) != 0 {
    return nil
  }
  return c.slice(c.values[ti])
}
```

1. Outer hash: `i = (FNV64(key)^r[0]) % #buckets`
2. Find `ri = indices[i]` => which random seed used for that bucket
3. Inner hash: `ti = ( (FNV64(key)^r[0]) ^ r ) % #keys` => position in `keys[]/values[]`
4. Compare stored key with actual, if match => get value

---

## 五、序列化/反序列化

### 5.1 `Write(w io.Writer) error`

- 先写 `len(r)`, 再写 r[]；
- 写 `len(indices)`, 再写 indices[]；
- 写 `len(keys)`, 对于每个 i：
  - 写 key 长度, value 长度 => 然后写 key bytes, value bytes

### 5.2 `Read(r io.Reader) (*CHD, error)`

```go
func Read(r io.Reader) (*CHD, error) {
  b, err := ioutil.ReadAll(r)
  if err != nil {return nil, err}
  return Mmap(b)
}
```

- 读取全部到内存，然后调用 `Mmap(b)`

### 5.3 `Mmap(b []byte) (*CHD, error)`

- 不真正复制，只在 `CHD` 里保存 `mmap = b`；
- 依序读 `rl = ReadInt()`, 读 `rl` 个 uint64 => `c.r`
- 读 `il = ReadInt()`, 读 `il` 个 uint16 => `c.indices`
- 读 `el = ReadInt()`, 遍历 e 次：
  - keyLen, valLen => sliceReader移动指针 => c.keys[i].start/end = ??? => c.values[i].start/end => ???

这样就实现了**零拷贝**(zero-copy) 解析，只记录每段在 `mmap` 中的位置。

---

## 六、迭代器

`CHD` 提供 `Iterate() *Iterator` 来遍历所有 (key, value)。在 `Iterator.Next()` 中切换下一个下标 `i` 并返回 `(c.slice(c.keys[i]), c.slice(c.values[i]))`。

---

## 七、总结

1. **`CHDBuilder`** 通过**两级哈希**+**随机碰撞重试**来构建**无冲突**的哈希表：
   - Outer hash 将 n 个 key roughly 分到 m 个桶；
   - Inner hash(多次随机尝试) 保证桶内 key 不冲突地映射到全局表中。
2. 构建成功后，`CHD` 里有：
   - `r[0]` 用于 outer hash；
   - `r[i]` (i>0) 用于各桶的 inner hash(通过 `indices[bucketIndex]` 找到具体 i)。
   - `mmap[]` 存储所有 (key, value) 的拼接数据；
   - `keys[i]`, `values[i]` 标识各记录在 `mmap` 里的切片位置。
3. `Get(key)` 时只需要**一次** outer hash + inner hash，即可 O(1) 找到全局下标，然后比对 key 是否真匹配。
4. `Write()` / `Read()` 让该结构持久化到文件或内存字节片；`Mmap()` 可以在**外部内存映射**下零拷贝地解析。

这是一个**静态**(只读)哈希结构；一旦构建好就不能轻易插入/删除(除非重新构建)。它在需要**大规模 key->value** 只读查询的场景（如编译器表、字典文件等）非常有用。
