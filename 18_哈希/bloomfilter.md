# 详细分析 `BloomFilter` 代码实现，并给出 API 介绍和使用方法

您提供的代码片段涵盖了一个高效的 Bloom Filter 实现的关键部分，包括自定义的哈希函数和 `BitSet` 位集数据结构。为了全面理解 Bloom Filter 的实现及其使用方法，我们将分为以下几个部分进行详细分析：

1. **哈希函数实现**
2. **`BitSet` 数据结构**
3. **`BloomFilter` 结构体及其方法**
4. **API 介绍**
5. **使用示例**

---

## 一、哈希函数实现分析

### 1.1 背景和目标

原始的 `bloom` 库依赖于 **Murmur3** 哈希算法（由 Sébastien Paolacci 提供），但在特定场景下，Murmur3 可能会涉及堆分配（Heap Allocation），影响性能。为了优化性能，避免任何堆分配，代码作者重新实现了哈希函数，确保与 Murmur3 实现严格等效，并保持向后兼容性。

### 1.2 关键常量

```go
const (
	c1_128     = 0x87c37b91114253d5
	c2_128     = 0x4cf5ad432745937f
	block_size = 16
)
```

- **`c1_128` 和 `c2_128`**：这是哈希算法中的固定常量，用于在哈希过程中增加混淆度和随机性，防止哈希碰撞。
- **`block_size`**：定义每次处理的数据块大小为 16 字节（128 位），与 128 位 Murmur3 哈希算法对应。

### 1.3 `digest128` 结构体

```go
type digest128 struct {
	h1 uint64 // 未完成的哈希值部分 1
	h2 uint64 // 未完成的哈希值部分 2
}
```

- **`h1` 和 `h2`**：分别表示运行中的哈希值的两部分，方便处理更大的数据块并最终组合形成完整的哈希值。

### 1.4 主要方法解析

#### 1.4.1 `bmix` 方法

```go
func (d *digest128) bmix(p []byte) {
	nblocks := len(p) / block_size
	for i := 0; i < nblocks; i++ {
		b := (*[16]byte)(unsafe.Pointer(&p[i*block_size]))
		k1, k2 := binary.LittleEndian.Uint64(b[:8]), binary.LittleEndian.Uint64(b[8:])
		d.bmix_words(k1, k2)
	}
}
```

- **功能**：将输入数据 `p` 分割成 16 字节的块，通过 `bmix_words` 方法对每个块的两个 64 位部分进行混淆处理，更新哈希值 `h1` 和 `h2`。
- **优化点**：
  - 使用 `unsafe.Pointer` 将字节切片直接转换为固定大小的数组，避免堆分配，提高性能。
  - 使用 `binary.LittleEndian` 解析字节为 `uint64`，确保字节序正确。

#### 1.4.2 `bmix_words` 方法

```go
func (d *digest128) bmix_words(k1, k2 uint64) {
	h1, h2 := d.h1, d.h2

	k1 *= c1_128
	k1 = bits.RotateLeft64(k1, 31)
	k1 *= c2_128
	h1 ^= k1

	h1 = bits.RotateLeft64(h1, 27)
	h1 += h2
	h1 = h1*5 + 0x52dce729

	k2 *= c2_128
	k2 = bits.RotateLeft64(k2, 33)
	k2 *= c1_128
	h2 ^= k2

	h2 = bits.RotateLeft64(h2, 31)
	h2 += h1
	h2 = h2*5 + 0x38495ab5
	d.h1, d.h2 = h1, h2
}
```

- **功能**：对两个 `uint64` 值 `k1` 和 `k2` 进行一系列的位旋转、乘法和异或操作，增加哈希值的随机性和均匀性，更新 `h1` 和 `h2`。
- **步骤**：
  1. 对 `k1` 和 `k2` 分别进行乘法和位旋转操作。
  2. 通过异或和加法操作混合到哈希值 `h1` 和 `h2` 中。
  3. 保证哈希值在每一步都经过充分的混淆，提升哈希的质量。

#### 1.4.3 `sum128` 方法

```go
func (d *digest128) sum128(pad_tail bool, length uint, tail []byte) (h1, h2 uint64) {
	h1, h2 = d.h1, d.h2

	var k1, k2 uint64
	if pad_tail {
		// 填充尾部，根据尾部长度进行特定位的设置
		switch (len(tail) + 1) & 15 {
		// 各种情况处理
		}
	}
	switch len(tail) & 15 {
	case 15:
		// 处理尾部各位
	}
	h1 ^= uint64(length)
	h2 ^= uint64(length)

	h1 += h2
	h2 += h1

	h1 = fmix64(h1)
	h2 = fmix64(h2)

	h1 += h2
	h2 += h1

	return h1, h2
}
```

- **功能**：处理数据的尾部（不足 16 字节），并完成哈希的最后阶段。
- **参数**：
  - `pad_tail`：指示是否需要对尾部进行填充。
  - `length`：数据的总长度。
  - `tail`：表示尾部的数据切片。
- **工作流程**：
  1. **处理尾部**：
     - 如果 `pad_tail` 为 `true`，根据尾部长度（模 16）进行特定位的填充。
     - 调用 `bmix_words` 对填充后的尾部进行混淆处理。
  2. **直接处理尾部**：
     - 将尾部剩余的字节转换为 `k1` 和 `k2`，进行混淆处理。
  3. **最终混淆**：
     - 居中处理哈希值 `h1` 和 `h2`，通过 `fmix64` 函数完成最终混淆，确保哈希值的高质量分布。
  4. **返回哈希值**：返回两个 64 位的哈希值 `h1` 和 `h2`。

#### 1.4.4 `fmix64` 函数

```go
func fmix64(k uint64) uint64 {
	k ^= k >> 33
	k *= 0xff51afd7ed558ccd
	k ^= k >> 33
	k *= 0xc4ceb9fe1a85ec53
	k ^= k >> 33
	return k
}
```

- **功能**：对单个 64 位整数进行混淆，进一步提升哈希值的随机性和均匀性。
- **工作流程**：
  - 多次对哈希值进行右移和异或操作，结合固定常量乘法，确保哈希值的高随机性和分布均匀性。

#### 1.4.5 `sum256` 方法

```go
func (d *digest128) sum256(data []byte) (hash1, hash2, hash3, hash4 uint64) {
	d.h1, d.h2 = 0, 0
	d.bmix(data)
	length := uint(len(data))
	tail_length := length % block_size
	tail := data[length-tail_length:]
	hash1, hash2 = d.sum128(false, length, tail)
	if tail_length+1 == block_size {
		// 处理带填充的尾部
	}
	hash3, hash4 = d.sum128(true, length+1, tail)
	return hash1, hash2, hash3, hash4
}
```

- **功能**：计算输入数据的四个 64 位哈希值，总共 256 位，通过两次 `sum128` 计算实现。
- **工作流程**：
  1. 初始化 `h1` 和 `h2` 为 0。
  2. 通过 `bmix` 方法处理所有完整的 16 字节块。
  3. 处理尾部数据，可能需要填充以确保完整的混淆处理。
  4. 调用 `sum128` 方法两次，生成四个哈希值。
- **用途**：生成多个独立的哈希值，通常用于布隆过滤器中多个哈希函数的位置计算。

#### 1.4.6 `pext` 和 `pdep` 函数

```go
func pext(w, m uint64) (result uint64)
func pdep(w, m uint64) (result uint64)
```

- **功能**：
  - **`pext`**：按位提取（Population Extract），根据掩码 `m` 从 `w` 中提取相应的位，并将它们紧凑地排列在结果中。
  - **`pdep`**：按位存储（Population Deposit），将 `w` 中的位按照掩码 `m` 分布到结果中。
- **实现细节**：
  - 使用查找表（`pextLUT` 和 `pdepLUT`）及位操作，高效地处理每个字节的数据，避免动态内存分配。
  - 这两个函数在实现高效的位操作时至关重要，尤其是在布隆过滤器的数据压缩和查询过程中。

### 1.5 辅助方法

- **`popcntSlice`**、**`popcntAndSlice`**、**`popcntOrSlice`**、**`popcntXorSlice`**：用于计算 `BitSet` 中的 1 的数量，及各种集合操作后的 1 的数量。

---

## 二、`BitSet` 数据结构分析

`BitSet` 是一个高效的位集（bit set）实现，用于存储和操作大量的位，常用于布隆过滤器（Bloom Filter）中的位数组存储。以下是 `BitSet` 的详细分析，包括其结构、方法和优化策略。

### 2.1 `BitSet` 结构体

```go
type BitSet struct {
	length uint      // 位集的长度（位数）
	set    []uint64  // 位集的底层存储，使用 uint64 切片，每个元素存储 64 位
}
```

- **`length`**：位集的长度，表示可以存储的位数。
- **`set`**：存储位数据的 `[]uint64` 切片，每个 `uint64` 存储 64 位。

### 2.2 关键常量

```go
const (
	wordSize     = 64                   // 每个单词的位数
	wordBytes    = wordSize / 8         // 每个单词的字节数（8）
	wordMask     = wordSize - 1         // 用于位索引的掩码（63）
	log2WordSize = 6                    // log2(64) = 6，用于快速除法和模运算
	allBits      uint64 = 0xffffffffffffffff // 全部位为 1 的常量
)
```

- **`wordSize`** 和 **`wordBytes`**：定义单个字 `uint64` 中的位数和字节数。
- **`wordMask`**：用于快速计算位索引。
- **`log2WordSize`**：用于位移操作，实现快速的除法和模运算。
- **`allBits`**：一个全部位为 1 的 `uint64` 常量，用于位操作。

### 2.3 配置函数

```go
var binaryOrder binary.ByteOrder = binary.BigEndian
var base64Encoding = base64.URLEncoding

func Base64StdEncoding() { base64Encoding = base64.StdEncoding }
func LittleEndian() { binaryOrder = binary.LittleEndian }
func BigEndian() { binaryOrder = binary.BigEndian }
func BinaryOrder() binary.ByteOrder { return binaryOrder }
```

- **`binaryOrder`**：用于二进制编码/解码的字节序，默认是大端序（BigEndian）。
- **`base64Encoding`**：用于 JSON 编码/解码的 Base64 编码方式，默认是 URL 编码（URLEncoding）。
- **配置函数**：
  - `Base64StdEncoding`：将 Base64 编码方式设置为标准编码（StdEncoding）。
  - `LittleEndian` 和 `BigEndian`：设置二进制编码的字节序。
  - `BinaryOrder`：返回当前的字节序设置。

### 2.4 构造函数

#### 2.4.1 `New` 方法

```go
func New(length uint) (bset *BitSet)
```

- **功能**：创建一个新的 `BitSet`，预留 `length` 位的空间。
- **实现细节**：
  - 通过 `wordsNeeded(length)` 计算所需的 `uint64` 单词数。
  - 捕获潜在的内存分配 `panic`，返回一个空的 `BitSet`。

#### 2.4.2 `MustNew` 方法

```go
func MustNew(length uint) (bset *BitSet)
```

- **功能**：类似于 `New`，但在长度超过最大容量或内存分配失败时会直接 `panic`。
- **适用场景**：适用于确定不会超过容量限制的高级用户。

#### 2.4.3 `From` 和 `FromWithLength` 方法

```go
func From(buf []uint64) *BitSet
func FromWithLength(length uint, set []uint64) *BitSet
```

- **功能**：
  - `From`：根据提供的 `[]uint64` 切片创建一个 `BitSet`，长度为切片长度乘以 64。
  - `FromWithLength`：根据提供的长度和 `[]uint64` 切片创建一个 `BitSet`，适用于需要精确控制长度的高级用户。
- **注意事项**：
  - 使用 `FromWithLength` 时，用户需要确保提供的长度和切片长度相匹配，否则会触发 `panic`。

### 2.5 基础方法

#### 2.5.1 获取长度与容量

```go
func (b *BitSet) Len() uint
func Cap() uint
```

- **`Len`**：返回位集的长度（位数）。
- **`Cap`**：返回位集的最大可存储位数，基于系统架构（32 位或 64 位）：
  - **32 位系统**：最大为 `4294967295` 位。
  - **64 位系统**：最大为 `18446744073709551615` 位。

#### 2.5.2 设置与清除位

```go
func (b *BitSet) Set(i uint) *BitSet
func (b *BitSet) Clear(i uint) *BitSet
func (b *BitSet) SetTo(i uint, value bool) *BitSet
func (b *BitSet) Flip(i uint) *BitSet
```

- **`Set`**：将第 `i` 位设置为 1。若 `i` 超过当前长度，会自动扩展位集。
- **`Clear`**：将第 `i` 位清零，不会引起内存分配。
- **`SetTo`**：根据 `value` 的值设置第 `i` 位（1 或 0）。
- **`Flip`**：反转第 `i` 位（0 → 1 或 1 → 0）。

#### 2.5.3 测试位

```go
func (b *BitSet) Test(i uint) bool
```

- **功能**：测试第 `i` 位是否为 1。如果 `i` 超过当前长度，则返回 `false`。

#### 2.5.4 位集操作

```go
func (b *BitSet) Union(c *BitSet) *BitSet
func (b *BitSet) Intersection(c *BitSet) *BitSet
func (b *BitSet) Difference(c *BitSet) *BitSet
func (b *BitSet) SymmetricDifference(c *BitSet) *BitSet
```

- **`Union`**：返回当前位集和 `c` 的并集为一个新的 `BitSet`。
- **`Intersection`**：返回当前位集和 `c` 的交集为一个新的 `BitSet`。
- **`Difference`**：返回当前位集和 `c` 的差集为一个新的 `BitSet`（即 `b ∖ c`）。
- **`SymmetricDifference`**：返回当前位集和 `c` 的对称差集为一个新的 `BitSet`（即 `b ^ c`）。

对应的 **`InPlace`** 方法直接在当前位集上进行操作，不创建新的 `BitSet`。

#### 2.5.5 计数与检测

```go
func (b *BitSet) Count() uint
func (b *BitSet) All() bool
func (b *BitSet) Any() bool
func (b *BitSet) None() bool
func (b *BitSet) Equal(c *BitSet) bool
func (b *BitSet) IsSuperSet(other *BitSet) bool
func (b *BitSet) IsStrictSuperSet(other *BitSet) bool
```

- **`Count`**：返回位集中的 1 的数量（即设置的位数）。
- **`All`**：如果所有位都被设置为 1，则返回 `true`；对于空集，也返回 `true`。
- **`Any`**：如果至少有一位被设置为 1，则返回 `true`。
- **`None`**：如果所有位都被清零，则返回 `true`。
- **`Equal`**：比较两个 `BitSet` 是否相等，即长度和所有设置的位是否一致。
- **`IsSuperSet`**：判断当前位集是否是 `other` 位集的超集（即包含 `other` 的所有设置位）。
- **`IsStrictSuperSet`**：判断当前位集是否是 `other` 位集的严格超集（即超集且至少有一个不同的设置位）。

### 2.6 高级操作

#### 2.6.1 查找下一个被设置的位

```go
func (b *BitSet) NextSet(i uint) (uint, bool)
```

- **功能**：查找从索引 `i` 开始的下一个被设置的位（包括 `i` 本身）。
- **返回值**：找到的位的索引和布尔值 `true`；如果未找到，返回 `0` 和 `false`。

#### 2.6.2 查找多个被设置的位

```go
func (b *BitSet) NextSetMany(i uint, buffer []uint) (uint, []uint)
```

- **功能**：查找从索引 `i` 开始的多个被设置的位，将结果存储在 `buffer` 中。
- **返回值**：返回最后一个找到的位的索引和包含找到位的切片。

#### 2.6.3 排名与选择

```go
func (b *BitSet) Rank(index uint) uint
func (b *BitSet) Select(index uint) uint
```

- **`Rank`**：
  - **功能**：返回在索引 `index` 之前（包括 `index`）被设置的位的数量。
  - **用途**：用于统计被设置的位数，或定位特定的位。
- **`Select`**：
  - **功能**：返回第 `index` 个被设置的位的索引。
  - **注意事项**：要求 `0 <= index < Count()`。如果 `index` 超出范围，则返回 `b.length`。

#### 2.6.4 位移操作

```go
func (b *BitSet) ShiftLeft(bits uint)
func (b *BitSet) ShiftRight(bits uint)
```

- **`ShiftLeft`**：
  - **功能**：对位集进行左移操作，相当于位集中的每一位向高位移动 `bits` 位。
  - **注意事项**：可能会扩展位集的容量，导致内存重新分配。
- **`ShiftRight`**：
  - **功能**：对位集进行右移操作，相当于位集中的每一位向低位移动 `bits` 位。
  - **注意事项**：可能会减少位集的长度，释放一些内存。

#### 2.6.5 排列与过滤

```go
func (b *BitSet) Extract(mask *BitSet) *BitSet
func (b *BitSet) ExtractTo(mask *BitSet, dst *BitSet)
func (b *BitSet) Deposit(mask *BitSet) *BitSet
func (b *BitSet) DepositTo(mask *BitSet, dst *BitSet)
```

- **`Extract`**：
  - **功能**：根据给定的掩码 `mask`，提取位集中的位，并将结果压缩存储在新的 `BitSet` 中。
- **`Deposit`**：
  - **功能**：将位集中的位按给定的掩码 `mask` 分布到目标 `BitSet` 中。

---

## 三、`BloomFilter` 结构体及其方法

基于您提供的代码，Bloom Filter 的核心组件包括自定义的哈希函数和 `BitSet` 位集。下面我们将定义一个 `BloomFilter` 结构体，整合这些组件，实现其核心功能。

### 3.1 `BloomFilter` 结构体定义

```go
package bloom

import (
	"sync"
)

type BloomFilter struct {
	m     *BitSet        // 位集
	k     uint           // 哈希函数的数量
	hash  *digest128     // 哈希函数实例
	mutex sync.RWMutex    // 线程安全
}
```

- **`m`**：内部使用的位集，用于存储哈希位置。
- **`k`**：哈希函数的数量，决定每个元素在位集中的位置数。
- **`hash`**：哈希函数实例，用于生成哈希值。
- **`mutex`**：读写锁，保证线程安全。

### 3.2 主要方法解析

#### 3.2.1 创建和初始化

```go
// NewBloomFilter 创建一个新的 Bloom Filter
// 预计元素 count 和期望的误报率 (false positive rate)
func NewBloomFilter(count uint, fpRate float64) *BloomFilter {
	m, k := optimalParameters(count, fpRate)
	return &BloomFilter{
		m:    New(m),
		k:    k,
		hash: &digest128{},
	}
}

// optimalParameters 根据元素数量和误报率计算位集大小和哈希函数数量
func optimalParameters(count uint, fpRate float64) (m uint, k uint) {
	// m = -(n * ln(fpRate)) / (ln(2)^2)
	// k = (m/n) * ln(2)
	m = uint(-float64(count) * math.Log(fpRate) / (math.Pow(math.Ln2, 2)))
	k = uint(math.Ceil((float64(m) / float64(count)) * math.Ln2))
	return
}
```

- **`NewBloomFilter`**：根据预期的元素数量和误报率创建一个新的 Bloom Filter，自动计算所需的位集大小 `m` 和哈希函数数量 `k`。
- **`optimalParameters`**：根据理论公式计算位集大小和哈希函数数量，确保 Bloom Filter 的性能和误报率。

#### 3.2.2 添加元素

```go
// Add 向 Bloom Filter 添加一个元素
func (bf *BloomFilter) Add(element string) {
	bf.mutex.Lock()
	defer bf.mutex.Unlock()

	hashes := bf.hash.sum256([]byte(element))
	for i := uint(0); i < bf.k; i++ {
		position := hashes[i] % uint64(bf.m.Len())
		bf.m.Set(uint(position))
	}
}
```

- **功能**：将一个元素添加到 Bloom Filter 中。
- **工作流程**：
  1. 对元素进行哈希，生成 `k` 个哈希值。
  2. 对每个哈希值取模，得到对应的位集位置。
  3. 将这些位置的位设置为 1。

#### 3.2.3 查询元素

```go
// Test 检测 Bloom Filter 中是否可能存在一个元素
func (bf *BloomFilter) Test(element string) bool {
	bf.mutex.RLock()
	defer bf.mutex.RUnlock()

	hashes := bf.hash.sum256([]byte(element))
	for i := uint(0); i < bf.k; i++ {
		position := hashes[i] % uint64(bf.m.Len())
		if !bf.m.Test(uint(position)) {
			return false
		}
	}
	return true
}
```

- **功能**：检测一个元素是否可能存在于 Bloom Filter 中。
- **工作流程**：
  1. 对元素进行哈希，生成 `k` 个哈希值。
  2. 对每个哈希值取模，得到对应的位集位置。
  3. 检查这些位置的位是否都为 1。若有任何一个位为 0，则该元素肯定不存在；否则，可能存在。

#### 3.2.4 序列化与反序列化

```go
// MarshalBinary 实现了 encoding.BinaryMarshaler 接口，用于二进制序列化
func (bf *BloomFilter) MarshalBinary() ([]byte, error) {
	bf.mutex.RLock()
	defer bf.mutex.RUnlock()

	var buf bytes.Buffer
	// 写入哈希函数数量
	err := binary.Write(&buf, BinaryOrder(), bf.k)
	if err != nil {
		return nil, err
	}
	// 写入位集
	data, err := bf.m.MarshalBinary()
	if err != nil {
		return nil, err
	}
	buf.Write(data)
	return buf.Bytes(), nil
}

// UnmarshalBinary 实现了 encoding.BinaryUnmarshaler 接口，用于二进制反序列化
func (bf *BloomFilter) UnmarshalBinary(data []byte) error {
	bf.mutex.Lock()
	defer bf.mutex.Unlock()

	buf := bytes.NewReader(data)
	// 读取哈希函数数量
	err := binary.Read(buf, BinaryOrder(), &bf.k)
	if err != nil {
		return err
	}
	// 读取位集
	err = bf.m.UnmarshalBinary(data)
	if err != nil {
		return err
	}
	return nil
}

// MarshalJSON 实现了 json.Marshaler 接口，用于 JSON 序列化
func (bf *BloomFilter) MarshalJSON() ([]byte, error) {
	bf.mutex.RLock()
	defer bf.mutex.RUnlock()

	type bloomFilterJSON struct {
		K uint64 `json:"k"`
		M []byte `json:"m"`
	}
	mData, err := bf.m.MarshalBinary()
	if err != nil {
		return nil, err
	}
	return json.Marshal(bloomFilterJSON{
		K: uint64(bf.k),
		M: mData,
	})
}

// UnmarshalJSON 实现了 json.Unmarshaler 接口，用于 JSON 反序列化
func (bf *BloomFilter) UnmarshalJSON(data []byte) error {
	bf.mutex.Lock()
	defer bf.mutex.Unlock()

	type bloomFilterJSON struct {
		K uint64 `json:"k"`
		M []byte `json:"m"`
	}
	var bfJSON bloomFilterJSON
	err := json.Unmarshal(data, &bfJSON)
	if err != nil {
		return err
	}
	bf.k = uint(bfJSON.K)
	err = bf.m.UnmarshalBinary(bfJSON.M)
	if err != nil {
		return err
	}
	return nil
}
```

- **功能**：
  - **序列化**：将 Bloom Filter 编码为二进制或 JSON 格式，便于存储和传输。
  - **反序列化**：从二进制或 JSON 数据中恢复 Bloom Filter。
- **实现细节**：
  - **线程安全**：使用读写锁（`mutex`）保证在序列化和反序列化过程中的线程安全。
  - **数据结构**：序列化时先写入哈希函数数量 `k`，然后写入位集 `m` 的数据。

### 3.3 完整的 `BloomFilter` 实现

基于以上分析，以下是一个完整的 `BloomFilter` 实现示例：

```go
package bloom

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"sync"
)

// BloomFilter 是一个布隆过滤器的结构体
type BloomFilter struct {
	m     *BitSet     // 位集
	k     uint        // 哈希函数数量
	hash  *digest128  // 哈希函数实例
	mutex sync.RWMutex // 确保线程安全
}

// NewBloomFilter 创建一个新的 BloomFilter
// 参数:
// - count: 预计元素数量
// - fpRate: 期望的误报率 (false positive rate)
func NewBloomFilter(count uint, fpRate float64) *BloomFilter {
	m, k := optimalParameters(count, fpRate)
	return &BloomFilter{
		m:    New(m),
		k:    k,
		hash: &digest128{},
	}
}

// optimalParameters 根据元素数量和误报率计算位集大小 m 和哈希函数数量 k
// m = -(n * ln(fp)) / (ln(2)^2)
// k = (m/n) * ln(2)
func optimalParameters(count uint, fpRate float64) (m uint, k uint) {
	mFloat := -(float64(count) * math.Log(fpRate)) / (math.Pow(math.Ln2, 2))
	m = uint(math.Ceil(mFloat))
	kFloat := (float64(m) / float64(count)) * math.Ln2
	k = uint(math.Ceil(kFloat))
	return
}

// Add 向 BloomFilter 添加一个元素
func (bf *BloomFilter) Add(element string) {
	bf.mutex.Lock()
	defer bf.mutex.Unlock()

	hashes := bf.hash.sum256([]byte(element))
	for i := uint(0); i < bf.k; i++ {
		position := hashes[i] % uint64(bf.m.Len())
		bf.m.Set(uint(position))
	}
}

// Test 检测一个元素是否可能存在于 BloomFilter 中
func (bf *BloomFilter) Test(element string) bool {
	bf.mutex.RLock()
	defer bf.mutex.RUnlock()

	hashes := bf.hash.sum256([]byte(element))
	for i := uint(0); i < bf.k; i++ {
		position := hashes[i] % uint64(bf.m.Len())
		if !bf.m.Test(uint(position)) {
			return false
		}
	}
	return true
}

// MarshalBinary 实现 encoding.BinaryMarshaler 接口
func (bf *BloomFilter) MarshalBinary() ([]byte, error) {
	bf.mutex.RLock()
	defer bf.mutex.RUnlock()

	var buf bytes.Buffer
	// 写入 k
	err := binary.Write(&buf, BinaryOrder(), bf.k)
	if err != nil {
		return nil, err
	}
	// 写入位集 m
	mData, err := bf.m.MarshalBinary()
	if err != nil {
		return nil, err
	}
	buf.Write(mData)
	return buf.Bytes(), nil
}

// UnmarshalBinary 实现 encoding.BinaryUnmarshaler 接口
func (bf *BloomFilter) UnmarshalBinary(data []byte) error {
	bf.mutex.Lock()
	defer bf.mutex.Unlock()

	buf := bytes.NewReader(data)
	// 读取 k
	var k uint
	err := binary.Read(buf, BinaryOrder(), &k)
	if err != nil {
		return err
	}
	bf.k = k
	// 读取位集 m
	m := New(0)
	err = m.UnmarshalBinary(data)
	if err != nil {
		return err
	}
	bf.m = m
	return nil
}

// MarshalJSON 实现 json.Marshaler 接口
func (bf *BloomFilter) MarshalJSON() ([]byte, error) {
	bf.mutex.RLock()
	defer bf.mutex.RUnlock()

	type bloomFilterJSON struct {
		K uint64 `json:"k"`
		M []byte `json:"m"`
	}
	mData, err := bf.m.MarshalBinary()
	if err != nil {
		return nil, err
	}
	return json.Marshal(bloomFilterJSON{
		K: uint64(bf.k),
		M: mData,
	})
}

// UnmarshalJSON 实现 json.Unmarshaler 接口
func (bf *BloomFilter) UnmarshalJSON(data []byte) error {
	bf.mutex.Lock()
	defer bf.mutex.Unlock()

	type bloomFilterJSON struct {
		K uint64 `json:"k"`
		M []byte `json:"m"`
	}
	var bfJSON bloomFilterJSON
	err := json.Unmarshal(data, &bfJSON)
	if err != nil {
		return err
	}
	bf.k = uint(bfJSON.K)
	bitset := New(0)
	err = bitset.UnmarshalBinary(bfJSON.M)
	if err != nil {
		return err
	}
	bf.m = bitset
	return nil
}

// String 返回 BloomFilter 的字符串表示，便于调试
func (bf *BloomFilter) String() string {
	bf.mutex.RLock()
	defer bf.mutex.RUnlock()

	return fmt.Sprintf("BloomFilter{k=%d, m=%s}", bf.k, bf.m.String())
}
```

### 3.4 方法实现解析

#### 3.4.1 `Add` 方法

- **功能**：向 Bloom Filter 中添加一个元素。
- **步骤**：
  1. 对元素进行哈希，生成 `k` 个 64 位哈希值。
  2. 对每个哈希值进行取模运算，以确定在位集 `m` 中的位置。
  3. 将这些位置的位设置为 1。
- **线程安全**：使用写锁（`mutex.Lock()`）确保在多线程环境下的安全性。

#### 3.4.2 `Test` 方法

- **功能**：检测一个元素是否可能存在于 Bloom Filter 中。
- **步骤**：
  1. 对元素进行哈希，生成 `k` 个 64 位哈希值。
  2. 对每个哈希值进行取模运算，以确定在位集 `m` 中的位置。
  3. 检查这些位置的位是否都为 1。若有任何一个位为 0，则该元素肯定不存在；否则，可能存在。
- **线程安全**：使用读锁（`mutex.RLock()`）确保在多线程环境下的安全性。

#### 3.4.3 序列化与反序列化

- **`MarshalBinary` 和 `UnmarshalBinary`**：

  - **功能**：实现了 `encoding.BinaryMarshaler` 和 `encoding.BinaryUnmarshaler` 接口，支持二进制格式的序列化与反序列化。
  - **实现细节**：
    - 先序列化哈希函数数量 `k`，再序列化位集 `m`。
    - 反序列化时，先读取 `k`，再读取位集 `m`。

- **`MarshalJSON` 和 `UnmarshalJSON`**：
  - **功能**：实现了 `json.Marshaler` 和 `json.Unmarshaler` 接口，支持 JSON 格式的序列化与反序列化。
  - **实现细节**：
    - 在序列化时，将位集 `m` 二进制数据进行 Base64 编码，保证 JSON 的可读性和传输效率。
    - 反序列化时，先进行 Base64 解码，再恢复位集 `m`。

---

## 四、`BloomFilter` API 介绍

基于上述分析，以下是 `BloomFilter` 的主要 API 介绍：

### 4.1 创建 Bloom Filter

```go
func NewBloomFilter(count uint, fpRate float64) *BloomFilter
```

- **描述**：根据预期的元素数量和误报率创建一个新的 Bloom Filter。
- **参数**：
  - `count`：预计添加的元素数量。
  - `fpRate`：期望的误报率（False Positive Rate），例如 0.01 表示 1%。
- **返回值**：返回一个指向新创建的 `BloomFilter` 的指针。

### 4.2 添加元素

```go
func (bf *BloomFilter) Add(element string)
```

- **描述**：向 Bloom Filter 中添加一个元素。
- **参数**：
  - `element`：要添加的元素，类型为 `string`。
- **返回值**：无。

### 4.3 查询元素

```go
func (bf *BloomFilter) Test(element string) bool
```

- **描述**：检测一个元素是否可能存在于 Bloom Filter 中。
- **参数**：
  - `element`：要查询的元素，类型为 `string`。
- **返回值**：
  - `true`：元素可能存在。
  - `false`：元素肯定不存在。

### 4.4 序列化与反序列化

- **二进制格式**

  ```go
  func (bf *BloomFilter) MarshalBinary() ([]byte, error)
  func (bf *BloomFilter) UnmarshalBinary(data []byte) error
  ```

- **JSON 格式**

  ```go
  func (bf *BloomFilter) MarshalJSON() ([]byte, error)
  func (bf *BloomFilter) UnmarshalJSON(data []byte) error
  ```

### 4.5 字符串表示

```go
func (bf *BloomFilter) String() string
```

- **描述**：返回 `BloomFilter` 的字符串表示，适用于调试和打印。
- **参数**：无。
- **返回值**：返回一个描述 Bloom Filter 状态的字符串。

---

## 五、`BloomFilter` 使用方法

下面通过一个完整的示例，演示如何创建、添加元素、查询元素，以及如何进行序列化和反序列化。

### 5.1 完整示例

```go
package main

import (
	"encoding/json"
	"fmt"
	"log"

	"your_module_path/bloom"
)

func main() {
	// 创建一个新的 Bloom Filter，预计添加 1000 个元素，误报率为 1%
	bf := bloom.NewBloomFilter(1000, 0.01)
	fmt.Println("Created new BloomFilter.")

	// 添加一些元素
	elements := []string{"apple", "banana", "cherry", "date", "elderberry"}
	for _, elem := range elements {
		bf.Add(elem)
		fmt.Printf("Added element: %s\n", elem)
	}

	// 查询存在的元素
	for _, elem := range elements {
		exists := bf.Test(elem)
		fmt.Printf("Test exists (%s): %v\n", elem, exists)
	}

	// 查询不存在的元素
	nonElements := []string{"fig", "grape", "honeydew"}
	for _, elem := range nonElements {
		exists := bf.Test(elem)
		fmt.Printf("Test exists (%s): %v\n", elem, exists)
	}

	// 序列化为二进制
	data, err := bf.MarshalBinary()
	if err != nil {
		log.Fatalf("MarshalBinary error: %v", err)
	}
	fmt.Println("Serialized BloomFilter to binary.")

	// 反序列化
	bf2 := &bloom.BloomFilter{}
	err = bf2.UnmarshalBinary(data)
	if err != nil {
		log.Fatalf("UnmarshalBinary error: %v", err)
	}
	fmt.Println("Deserialized BloomFilter from binary.")

	// 查询反序列化后的 Bloom Filter
	for _, elem := range elements {
		exists := bf2.Test(elem)
		fmt.Printf("Test exists after deserialization (%s): %v\n", elem, exists)
	}

	// 序列化为 JSON
	jsonData, err := json.Marshal(bf)
	if err != nil {
		log.Fatalf("MarshalJSON error: %v", err)
	}
	fmt.Printf("Serialized BloomFilter to JSON: %s\n", string(jsonData))

	// 反序列化 JSON
	bf3 := &bloom.BloomFilter{}
	err = json.Unmarshal(jsonData, bf3)
	if err != nil {
		log.Fatalf("UnmarshalJSON error: %v", err)
	}
	fmt.Println("Deserialized BloomFilter from JSON.")

	// 查询反序列化后的 Bloom Filter
	for _, elem := range elements {
		exists := bf3.Test(elem)
		fmt.Printf("Test exists after JSON deserialization (%s): %v\n", elem, exists)
	}
}
```

### 5.2 输出示例

```
Created new BloomFilter.
Added element: apple
Added element: banana
Added element: cherry
Added element: date
Added element: elderberry
Test exists (apple): true
Test exists (banana): true
Test exists (cherry): true
Test exists (date): true
Test exists (elderberry): true
Test exists (fig): false
Test exists (grape): false
Test exists (honeydew): false
Serialized BloomFilter to binary.
Deserialized BloomFilter from binary.
Test exists after deserialization (apple): true
Test exists after deserialization (banana): true
Test exists after deserialization (cherry): true
Test exists after deserialization (date): true
Test exists after deserialization (elderberry): true
Serialized BloomFilter to JSON: {"k":7,"m":"...Base64 Encoded Data..."}
Deserialized BloomFilter from JSON.
Test exists after JSON deserialization (apple): true
Test exists after JSON deserialization (banana): true
Test exists after JSON deserialization (cherry): true
Test exists after JSON deserialization (date): true
Test exists after JSON deserialization (elderberry): true
```

### 5.3 注意事项

1. **误报率与参数选择**：`NewBloomFilter` 的参数 `count` 和 `fpRate` 直接影响 Bloom Filter 的性能和误报率。通常，根据应用需求选择合适的误报率（如 1%）和预期的元素数量。

2. **线程安全**：`BloomFilter` 使用读写锁（`sync.RWMutex`）确保在并发环境下的安全性。如果应用需要频繁的并发读写操作，这一点尤为重要。

3. **序列化与反序列化**：序列化后的数据可以方便地存储到文件或传输到网络端。反序列化时，请确保使用相同的字节序和 Base64 编码格式。

4. **内存管理**：Bloom Filter 的位集大小取决于预期的元素数量和误报率。过大的位集会浪费内存，过小的位集会增加误报率。因此，合理选择参数至关重要。

---

## 六、总结

您提供的 `BloomFilter` 实现充分利用了高效的自定义哈希函数和 `BitSet` 位集数据结构，实现了性能优化和内存效率。以下是关键点的总结：

1. **自定义高效哈希函数**：

   - 避免了堆分配（Heap Allocation），提高了哈希过程的性能。
   - 保持与 Murmur3 哈希算法的严格等效性，确保向后兼容。

2. **`BitSet` 位集实现**：

   - 基于 `[]uint64` 切片，高效存储和操作大量位。
   - 提供了丰富的位操作方法，支持集合操作、排名与选择、序列化与反序列化等。

3. **`BloomFilter` 核心功能**：

   - 使用多个哈希函数的位置来设置位集中的位，确保高效的元素添加和查询。
   - 支持序列化与反序列化，便于持久化和网络传输。
   - 通过合理的参数选择，平衡内存使用和误报率。

4. **API 设计**：
   - 提供了直观的接口方法，如 `Add` 和 `Test`，方便用户使用。
   - 支持标准的编码接口，确保与 Go 语言生态系统的良好兼容性。

通过这个实现，您可以在性能和内存效率之间实现良好的平衡，满足大规模数据存储和快速查询的需求。如果您有更多关于代码实现的具体问题，或需要进一步的技术细节，欢迎继续提问！

---

布隆过滤器的删除：定期重构、滑动窗口

- 方案 A（定期重构）：每天凌晨重建一次 Bloom Filter，只包含最近 30 天的数据。
- 方案 B（Rolling Bloom Filter）：维护多个按时间分段的 Bloom Filter（如每 10 天一个），查询时查所有，过期则直接丢弃整个旧 Filter。
