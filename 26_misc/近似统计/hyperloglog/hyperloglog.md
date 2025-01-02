下面是一份**详细**的代码解读，涵盖了 HyperLogLog 算法背景、数据结构设计和实现原理等部分。由于代码非常长，内容繁多，解读也会相应地分模块进行。希望能帮助你更好地理解该实现的工作机制。

---

# 总览

这段代码主要实现了**HyperLogLog**（以及其稀疏版本）的数据结构和操作，用来**估计集合基数（distinct count）**。HyperLogLog 在大数据场景下以非常低的内存开销，完成近似的去重计数。

在 Go 语言中，HyperLogLog 常见的开源实现如 `clarkduvall/hyperloglog`，这里的代码看起来与其存在部分相似之处，并且额外实现了自定义的**稀疏表示**（sparse representation）和**序列化/反序列化**功能，方便在大数据环境中使用。

---

# 主体结构：HyperLogLog Sketch

## 1. 常量与全局

```go
const (
    pp      = uint8(25)        // 全局常量, 用作 HLL 内部的 bit 操作, 主要与稀疏模式相关
    mp      = uint32(1) << pp  // 1 << 25, = 33554432
    version = 2                // 数据结构序列化的版本号
)
```

- `pp`：代表了稀疏模式（Sparse Representation）的一些内部处理所需的位宽（25 位），用来截取哈希值的部分位段。
- `mp = 1 << pp`：把 `pp` 看成 25，则 `mp = 2^25`，在一些地方用来做稀疏结构的计算（如 `linearCount(mp, mp - sk.sparseList.count)`）。
- `version`：序列化版本号。

## 2. `Sketch` 结构体

```go
type Sketch struct {
    p          uint8
    m          uint32
    alpha      float64
    tmpSet     *set
    sparseList *compressedList
    regs       []uint8
}
```

HyperLogLog 的关键字段：

- `p`：**精度**（precision）。决定了 `m = 2^p` 个桶。p 越大，HLL 的精度越高，估计越准确，但内存也越大。
- `m`：桶数量，`m = 1 << p`。
- `alpha`：常数因子，根据 `m` 不同取不同值，一般公式是 `\(\alpha = 0.7213 / (1 + 1.079/m)\)`.
- `tmpSet` / `sparseList`：如果启用了 **稀疏模式**，前期会把数据以一种“hash -> set/list”的方式存储，等稀疏数据较多时再切换到普通模式。
- `regs`：当不使用稀疏模式时，直接使用 `regs` 这个长度为 `m` 的数组，每个元素记录某个桶的 \(\rho\) 值（见后）。

**HyperLogLog 原理简述**

- 将每个元素做哈希，哈希值高位的 `p` 位用于定位“第几个桶”，然后根据后续低位计算前导零的个数 `\(\rho\)`；
- 记录在 `regs[bucket]` 中的值是该桶所见到的最大前导零数；
- 最终通过合并所有桶的估计值来得出全局基数估计。

---

# 构造与初始化

```go
func New() *Sketch           { return New14() }
func New14() *Sketch         { return newSketchNoError(14, true) }
func New16() *Sketch         { return newSketchNoError(16, true) }
func NewNoSparse() *Sketch   { return newSketchNoError(14, false) }
func New16NoSparse() *Sketch { return newSketchNoError(16, false) }

func newSketchNoError(precision uint8, sparse bool) *Sketch {
    sk, _ := NewSketch(precision, sparse)
    return sk
}

func NewSketch(precision uint8, sparse bool) (*Sketch, error) {
    if precision < 4 || precision > 18 {
        return nil, fmt.Errorf("p has to be >= 4 and <= 18")
    }
    m := uint32(1) << precision
    s := &Sketch{
        m:     m,
        p:     precision,
        alpha: alpha(float64(m)),
    }
    if sparse {
        s.tmpSet = newSet(0)
        s.sparseList = newCompressedList(0)
    } else {
        s.regs = make([]uint8, m)
    }
    return s, nil
}
```

- `New14() / New16()`：创建 p=14 (即 m=2^14) 或 p=16 (m=2^16) 的 HyperLogLog，并默认**启用稀疏模式**。
- `NewNoSparse()`：创建 p=14 但**不启用稀疏模式**。
- `NewSketch(precision, sparse)`：真正的构造函数。
  - 若 `precision < 4` 或 `>18` 会返回错误，防止极端不合理的值。
  - 根据 `sparse`，要么初始化 `tmpSet + sparseList`，要么直接分配 `regs`。

**稀疏模式**（sparse）：

- 当数据量不大时，完全没有必要开辟一个大数组 `regs`；
- 此时先用 `tmpSet` + `sparseList` 的组合以稀疏的形式（仅存储插入的 hash）记录；
- 当数据量逐渐增大，超过一定阈值，就会**转换**为传统的“密集（dense）”数组模式（`regs`）。

---

# Clone / Merge / Convert

## 1. Clone

```go
func (sk *Sketch) Clone() *Sketch {
    clone := *sk
    clone.regs = append([]uint8(nil), sk.regs...)
    clone.tmpSet = sk.tmpSet.Clone()
    clone.sparseList = sk.sparseList.Clone()
    return &clone
}
```

- 复制一份 Sketch：浅拷贝结构体本身，再分别**深拷贝** `regs`、`tmpSet`、`sparseList`。
- 这样不会产生共享底层切片导致的后续修改冲突。

## 2. Merge

```go
func (sk *Sketch) Merge(other *Sketch) error {
    if other == nil {
        return nil
    }
    if sk.p != other.p {
        return errors.New("precisions must be equal")
    }
    if sk.sparse() && other.sparse() {
        sk.mergeSparseSketch(other)
    } else {
        sk.mergeDenseSketch(other)
    }
    return nil
}
```

- **Merge** 用于将另一个 HyperLogLog 的分布信息合并到当前对象，以便综合统计。
- 必须保证两者 `p` 相同，否则桶数不同，算法不兼容。
- 如果都是稀疏，就合并稀疏结构；否则转为合并稠密结构。

### 2.1 mergeSparseSketch

```go
func (sk *Sketch) mergeSparseSketch(other *Sketch) {
    sk.tmpSet.Merge(other.tmpSet)
    for iter := other.sparseList.Iter(); iter.HasNext(); {
        sk.tmpSet.add(iter.Next())
    }
    sk.maybeToNormal()
}
```

- 将 `other` 的稀疏内容全部塞进 `sk.tmpSet`，然后调用 `maybeToNormal()` 看是否需要转为稠密模式。

### 2.2 mergeDenseSketch

```go
func (sk *Sketch) mergeDenseSketch(other *Sketch) {
    if sk.sparse() {
        sk.toNormal()
    }

    if other.sparse() {
        // 把 other 的稀疏数据一个个 decode 再插入 sk
        ...
    } else {
        // 双方都是稠密模式，直接合并桶 regs 中的最大前导零
        for i, v := range other.regs {
            if v > sk.regs[i] {
                sk.regs[i] = v
            }
        }
    }
}
```

- 如果当前 `sk` 还在稀疏模式，需要先转为稠密模式，然后再合并。
- 若 `other` 仍是稀疏，则遍历其稀疏内容——对每个 hash decode 出 `(bucket, \rho)`，再更新到 `sk.regs[bucket]`。
- 如果双方都是稠密，则只需 `max()` 合并对应桶的 \(\rho\) 值即可。

---

## 3. 转换至稠密模式

```go
func (sk *Sketch) maybeToNormal() {
    if uint32(sk.tmpSet.Len())*100 > sk.m {
        sk.mergeSparse()
        if uint32(sk.sparseList.Len()) > sk.m {
            sk.toNormal()
        }
    }
}
```

- 当稀疏数据的数量达到一定阈值（`tmpSet.Len() * 100 > sk.m`），就先做一次合并（把 `tmpSet` 合并进 `sparseList`），然后如果 `sparseList.Len()` 依旧大于 `sk.m`，说明数据已经很大，直接切换到稠密模式 `toNormal()`。

```go
func (sk *Sketch) toNormal() {
    if sk.tmpSet.Len() > 0 {
        sk.mergeSparse()
    }
    sk.regs = make([]uint8, sk.m)
    for iter := sk.sparseList.Iter(); iter.HasNext(); {
        i, r := decodeHash(iter.Next(), sk.p, pp)
        sk.insert(i, r)
    }
    sk.tmpSet = nil
    sk.sparseList = nil
}
```

- `toNormal()`：
  1. 先合并 `tmpSet`。
  2. 分配 `regs` 数组。
  3. 遍历 `sparseList`，对每个 hash decode 出 `(i, r)`, 调用 `insert()` 写入 `regs[i]`。
  4. 把 `tmpSet` 和 `sparseList` 置空。

---

# Insert / Estimate

## 1. Insert

```go
func (sk *Sketch) Insert(e []byte) {
    sk.InsertHash(hash(e))
}

func (sk *Sketch) InsertHash(x uint64) {
    if sk.sparse() {
        if sk.tmpSet.add(encodeHash(x, sk.p, pp)) {
            sk.maybeToNormal()
        }
        return
    }
    i, r := getPosVal(x, sk.p)
    sk.insert(uint32(i), r)
}
```

- 将任意字节数组 `e` 哈希为 `x`，然后插入。
- 如果当前是稀疏模式，就把 `encodeHash(x, sk.p, pp)` 放进 `tmpSet`，若成功插入，就可能触发 `maybeToNormal()`。
- 否则直接计算 `(桶索引 i, 前导零数 r)` 并存进 `regs[i]`。

```go
func getPosVal(x uint64, p uint8) (uint64, uint8) {
    // i = 高位 p 位
    i := bextr(x, 64-p, p)
    // w = 其余部分(低位) + 1<<(p-1)
    w := x<<p | 1<<(p-1)
    rho := uint8(bits.LeadingZeros64(w)) + 1
    return i, rho
}
```

- `bextr`：bit extraction，取出 `x` 在 [start, start+length) 范围内的 bits。
- `rho`: 统计 `w` 的前导零数 + 1，是 HyperLogLog 经典的 \(\rho\) 计算方式。

## 2. Estimate

```go
func (sk *Sketch) Estimate() uint64 {
    if sk.sparse() {
        sk.mergeSparse()
        return uint64(linearCount(mp, mp - sk.sparseList.count))
    }

    sum, ez := sumAndZeros(sk.regs)
    m := float64(sk.m)
    est := sk.alpha * m * (m - ez) / (sum + beta(sk.p, ez))
    return uint64(est + 0.5)
}
```

- 如果还在稀疏模式（并且其实已经插入了一些数据），先调用 `mergeSparse()` 把 `tmpSet` 合并进 `sparseList`，然后用**线性计数**（Linear Counting）方法估计：`linearCount(mp, mp - sk.sparseList.count)`。
  - 线性计数适用于“尚未装满时”，适合稀疏的情况。
- 如果在**稠密模式**：
  1. `sumAndZeros(sk.regs)`：累加 \(\frac{1}{2^{regs[i]}}\) 并统计空桶数 `ez`；
  2. `est = alpha * m^2 / (sum + beta(...))` 的形式，有一项 `(m - ez)` 修正；另外 `beta(p, ez)` 用于微调小基数的估计偏差。
  3. 最后返回 `uint64(est + 0.5)` 做四舍五入。

```go
func sumAndZeros(regs []uint8) (res, ez float64) {
    for _, v := range regs {
        if v == 0 {
            ez++
        }
        res += 1.0 / math.Pow(2.0, float64(v))
    }
    return res, ez
}
```

---

# 稀疏部分

## 1. encodeHash / decodeHash

```go
func encodeHash(x uint64, p, pp uint8) uint32 { ... }
func decodeHash(k uint32, p, pp uint8) (uint32, uint8) { ... }
```

- 这里把一个 64 位 hash 拆分并压缩进 32 位 `uint32`，作为稀疏结构的标记。
- `encodeHash`：取高位 `pp` 位当作 `idx`；如果再往下还有全 0 的话，把**前导零数**编码进后续 bits。
- `decodeHash`：还原 `(bucketIndex, \rho)`。

## 2. `set` + `compressedList`

HyperLogLog 稀疏模式下，会把新 hash 值先放进一个 `tmpSet`，然后在合并/转换时，把它们合并进 `sparseList`。

- `set` 是一个自定义结构，底层是 `Map[uint32, struct{}]`，用来去重。
- `compressedList` 则是一个**有序存储**的结构，里面用“前缀差分 + 可变长编码”来节省空间。

```go
type set struct {
    m *Set[uint32]
}

type compressedList struct {
    count uint32
    last  uint32
    b     variableLengthList
}
```

- `tmpSet` 存 hash。
- `sparseList`（compressedList）记录很多 hash，但采用**可变长差分**，把相邻元素的差值再用**7-bit**分块编码，减少空间。

### 2.1 mergeSparse

```go
func (sk *Sketch) mergeSparse() {
    if sk.tmpSet.Len() == 0 {
        return
    }

    keys := make(uint64Slice, 0, sk.tmpSet.Len())
    sk.tmpSet.ForEach(func(k uint32) {
        keys = append(keys, k)
    })
    sort.Sort(keys)

    newList := newCompressedList(4*sk.tmpSet.Len() + sk.sparseList.Len())
    ...
    sk.sparseList = newList
    sk.tmpSet = newSet(0)
}
```

- 把 `tmpSet` 里的所有元素取出来并排序，然后与 `sparseList` 做**归并**（类似归并排序的过程），最终生成一个新的 `newList`。
- 清空 `tmpSet`。

---

# 序列化与反序列化

## 1. MarshalBinary

```go
func (sk *Sketch) MarshalBinary() (data []byte, err error) {
    data = make([]byte, 0, 8+len(sk.regs))
    // 1) 版本号
    data = append(data, version)
    // 2) p
    data = append(data, sk.p)
    // 3) b (used in V1)
    data = append(data, 0)

    if sk.sparse() {
        // 标记稀疏模式 byte(1)
        data = append(data, byte(1))

        tsdata, err := sk.tmpSet.MarshalBinary()
        ...
        sdata, err := sk.sparseList.MarshalBinary()
        ...
        return append(data, sdata...), nil
    }

    // 否则标记稠密 byte(0)
    ...
    // 长度 + regs 数据
    ...
    return data, nil
}
```

- 写出：版本号、`p`、一个留给兼容老版本的占位、是否稀疏标记...
- 如果稀疏，序列化 `tmpSet` 和 `sparseList`；否则序列化 `regs`。

## 2. UnmarshalBinary

```go
func (sk *Sketch) UnmarshalBinary(data []byte) error {
    if len(data) < 8 {
        return ErrorTooShort
    }
    v := data[0]
    p := data[1]
    b := data[2]
    sparse := data[3] == byte(1)
    ...
    if sparse {
        // 读取 tmpSet.size
        // 依次读取 tmpSet 的元素
        // 再调用 sparseList.UnmarshalBinary
    } else {
        if v == 1 {
            return sk.unmarshalBinaryV1(data[8:], b)
        }
        return sk.unmarshalBinaryV2(data)
    }
}
```

- 先从 `data` 里读出 4 个字节（版本号 / p / b / 是否稀疏），然后根据稀疏 / 稠密做相应处理。
- 稀疏：先读 `tmpSet`，再读 `sparseList`。
- 稠密：V1 或 V2 不同的压缩格式。
  - `unmarshalBinaryV1()`：以前版本把每 2 个 4-bit 值打包进一个字节。
  - `unmarshalBinaryV2()`：直接拷贝后续数据到 `regs`。

---

# Beta 修正函数

```go
var betaMap = map[uint8]func(float64) float64{
    4: beta4, 5: beta5, ..., 18: beta18,
}
func beta(p uint8, ez float64) float64 {
    f, ok := betaMap[p]
    return f(ez)
}
```

- 这些 `betaX()` 函数是一批多项式回归或插值函数，用来纠正 HyperLogLog 在小基数情况下的估计偏差，使得估计结果更接近真实值。
- 不同 `p`（桶数量）对应一组多项式系数，都是预先实验拟合好的。

---

# 其他辅助部分

## 1. `Map` / `Set` 实现

在文件下方能看到一个自定义的**哈希 Map**（支持任意整数类型 K）与**Set**。

- `type Set[K IntKey] Map[K, struct{}]`：只存 key 的集合，不关心 value。
- `Map[K, V]`：开放寻址哈希表，用 “phiMix64” 做哈希扰动。
- 当 `key = 0` 时，有特殊处理（用 `hasZeroKey` 标志），以免 0 与空位冲突。
- 这段实现大体上是一个简单的 Robin Hood / 线性探查变体。

## 2. `metro` Hash64

```go
func Hash64(buffer []byte, seed uint64) uint64 {
    // 这是一个自定义的哈希函数 (MetroHash 或类似)，对输入 buffer 分块处理，
    // 进行若干轮混合、旋转等操作，最终输出 64-bit hash。
}
```

- 不同于 Go 标准库的 `hash/fnv` 或 `crypto/sha256`，这里是一个自定义高效的 64 位哈希函数，对性能要求高的大数据场景更友好。

---

# HyperLogLog 核心原理回顾

1. **桶数量**：\( m = 2^p \)。
2. **哈希映射**：对元素做 64 位哈希。
   - 高位 `p` 位决定所属桶索引 `i`。
   - 剩余位用于计算前导零数 \(\rho\)。
3. **合并规则**：对于桶 `i`，保存观测到的最大 \(\rho\)。
4. **估计公式**：\[
   E \approx \alpha*m \cdot m^2 \left/ \sum*{j=1}^{m} 2^{-M[j]} \right.
   \]
   但在小基数时要做**线性计数**或**beta 校正**。

在本实现中，还通过**稀疏表示**在小数据量时节省内存，通过**转为稠密**在数据变大时保持 O(1) 查询/更新时间。

---

# 小结

1. **HyperLogLog 主结构**

   - `Sketch`：维护 `p, m, alpha`，以及**稀疏**或**稠密**存储。
   - 稀疏：`tmpSet + sparseList`；稠密：`regs` 数组。

2. **插入**

   - 稀疏模式：将编码后的 hash 放入 `tmpSet`，必要时合并进 `sparseList`，并可能转稠密。
   - 稠密模式：直接更新 `regs[bucket] = max(regs[bucket], ρ)`。

3. **合并 (Merge)**

   - 可以把另一个同精度的 HLL 合并进来；若双方都是稀疏则简单合并稀疏，否则转稠密再合并。

4. **估计 (Estimate)**

   - 稀疏：先做一次合并，然后用线性计数；
   - 稠密：根据经典 HLL 公式 + `beta()` 多项式修正。

5. **序列化**

   - 区分版本（V1 / V2）及稀疏 / 稠密模式，写入相应数据结构；
   - 反序列化时读出标记，还原 `tmpSet / sparseList / regs`。

6. **额外**
   - 自定义哈希函数 `Hash64`
   - 自定义 `Map` 与 `Set`（简单的开放寻址哈希表）
   - 大量位操作函数（`bextr, bextr32, bits.LeadingZeros64` 等）保证高效处理哈希值。

---

**总结**：  
这段代码实现了一个带**稀疏优化**的 HyperLogLog，并支持**序列化/反序列化**。它能在小数据量时节省空间（稀疏模式），在大数据量时自动转为稠密模式以保持性能，并通过多个多项式函数 `betaX` 来提升估计精度。核心思想仍是：**使用多个桶记录最大前导零数**，然后通过统计公式估计去重总数。该实现相当完整，也比较复杂，适合在生产环境进行大规模去重计数的需求。
