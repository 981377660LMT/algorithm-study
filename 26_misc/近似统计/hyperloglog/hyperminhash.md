下面的代码实现了一个基于 _HyperMinHash_（更确切地说是 LogLog-Beta 算法风格）的**基数（Cardinality）估计**以及**集合运算**（相似度、交集）功能。整体思路类似于 HyperLogLog 的算法：

1. 用一个大数组（寄存器 `register`）来记录哈希值中的信息；
2. 通过统计寄存器数组中记录到的信息（如最长前导零、签名等）来估计集合的基数；
3. 通过特殊的碰撞概率推导，来估计不同两个集合的交集大小或相似度。

下面我们从整体结构和关键函数的角度，**逐行剖析**、**解释关键逻辑**。

---

## 1. 常量定义

```go
const (
    p     = 14
    m     = uint32(1 << p) // 16384
    max   = 64 - p
    maxX  = math.MaxUint64 >> max
    alpha = 0.7213 / (1 + 1.079/float64(m))
    q     = 6  // the number of bits for the LogLog hash
    r     = 10 // number of bits for the bbit hash
    _2q   = 1 << q
    _2r   = 1 << r
    c     = 0.169919487159739093975315012348
)
```

- **p=14**
  - `m = 1 << p = 16384`。这是寄存器数组的大小，每个元素是 `register` 类型，用来存放某条哈希信息。
- **max = 64 - p**
  - `maxX` = 将 `math.MaxUint64` 右移 `max` 位，后续用来辅助提取哈希值某些部分。
- **alpha**：LogLog 类算法用到的一个校正因子，常见的 HLL 算法会有不同的 α (或称 \(\alpha_m\))。
- **q=6、r=10**
  - 这两个值主要用来在一个 16bit 的 `register` 中，前 `q` 位用来存“leading zeros”（记作 lz），后 `r` 位用来存签名（signature）。
  - \(\_2^q = 64\)，\(\_2^r = 1024\)。
- **c**：在估计碰撞时的一些修正系数（用于算法中的碰撞概率公式，具体细节来自论文/研究结果）。

---

## 2. `Sketch` 数据结构

```go
type register uint16

func (reg register) lz() uint8 {
    return uint8(uint16(reg) >> (16 - q))
}

func newReg(lz uint8, sig uint16) register {
    return register((uint16(lz) << r) | sig)
}

// Sketch is a sketch for cardinality estimation based on LogLog counting
type Sketch struct {
    reg [m]register
}
```

### 2.1 `register`

- `register` 用一个 `uint16` 来存储两部分信息：

  1. **lz**：前导零（Leading Zeros）的数量，用 `q` 位表示。
  2. **sig**：签名（Signature），用 `r` 位表示。

- `lz()`：取出寄存器里高 `q` 位。
- `newReg(lz, sig)`：将前导零数 `lz` 放到高 `q` 位，把签名 `sig` 放到低 `r` 位，并组合成一个 `uint16`。

### 2.2 `Sketch`

- `Sketch` 就是一组固定大小为 `m=16384` 的 `register` 数组，用来记录所有插入值的哈希信息。

---

## 3. 各种核心函数

### 3.1 `NewHyperMinHash`

```go
func NewHyperMinHash() *Sketch {
    return new(Sketch)
}
```

- 返回一个空的 `Sketch` 对象，内部的寄存器数组默认都是 0。

### 3.2 `Add` / `AddHash`

```go
func (sk *Sketch) Add(value []byte) {
    h1, h2 := Hash128(value, 1337)
    sk.AddHash(h1, h2)
}
```

1. `Add(value)`:
   - 先调用 `Hash128`，对传进来的 `value` 做 128 位哈希，返回两个 64 位整数 `h1, h2`。
   - 再调用 `AddHash(h1, h2)`。

```go
func (sk *Sketch) AddHash(x, y uint64) {
    k := x >> uint(max)                 // 取哈希值x的最高 p 位，作为寄存器的下标
    lz := uint8(bits.LeadingZeros64((x<<p)^maxX)) + 1
    sig := uint16(y << (64 - r) >> (64 - r))
    reg := newReg(lz, sig)
    if sk.reg[k] < reg {
        sk.reg[k] = reg
    }
}
```

2. `AddHash(x, y)`：
   - **`k`**：`x >> (64 - p)`，这会把 `x` 的高 `p` 位当作数组索引；这样能把所有输入值分散到 `m = 2^p` 个桶（寄存器）里。
   - **`lz`**：先做 `(x << p) ^ maxX`，然后数它的前导零，最后加 1。`lz` 表示我们关心的哈希值“看起来有多少个前导 0”。LogLog 算法核心就是统计这个“最大前导零数”。
   - **`sig`**：从 `y` 中再提取后 `r` 位，当作签名。“签名”在 HyperMinHash 或其它类似算法里，用来帮助做相交和相似度估计。
   - **更新寄存器**：`if sk.reg[k] < reg { sk.reg[k] = reg }`
     - 只有当新的 `reg`（也就是 `(lz, sig)`) 大于已有的寄存器时才更新。这也是 HyperLogLog 里获取“最大前导零”的思路，只不过这里还带着签名部分。

---

### 3.3 `Cardinality`

```go
func (sk *Sketch) Cardinality() uint64 {
    sum, ez := regSumAndZeros(sk.reg[:])
    m := float64(m)
    return uint64(alpha * m * (m - ez) / (beta(ez) + sum))
}
```

- 先通过 `regSumAndZeros` 得到：
  1. `sum`：把每个寄存器的 lz 取出来，然后做 \(\sum \frac{1}{2^{lz}}\)。
  2. `ez`：有多少寄存器的 lz=0（表明这个桶里还没记录到有效信息，或者只记录到非常小的前导零数）。
- 计算方式：  
  \[
  \text{Cardinality} = \alpha \times m \times \frac{m - ez}{\beta(ez) + \sum}
  \]
  - `alpha`、`beta` 都是经验或理论校正函数，用来让估计值更准确。
  - `beta(ez)`：当空桶（`ez`）较多时，估计中要做进一步修正。

#### `regSumAndZeros`

```go
func regSumAndZeros(registers []register) (float64, float64) {
    var sum, ez float64
    for _, val := range registers {
        lz := val.lz()
        if lz == 0 {
            ez++
        }
        sum += 1 / math.Pow(2, float64(lz))
    }
    return sum, ez
}
```

- 挨个寄存器检查 `lz` 值，如果是 0，`ez++`；否则按 `1 / 2^lz` 累加到 `sum`。

#### `beta(ez)`

```go
func beta(ez float64) float64 {
    zl := math.Log(ez + 1)
    return -0.370393911*ez +
        0.070471823*zl +
        0.17393686*math.Pow(zl, 2) +
        0.16339839*math.Pow(zl, 3) +
        -0.09237745*math.Pow(zl, 4) +
        0.03738027*math.Pow(zl, 5) +
        -0.005384159*math.Pow(zl, 6) +
        0.00042419*math.Pow(zl, 7)
}
```

- 对 `ez` 做多项式拟合的一种经验修正。

---

### 3.4 `Merge`（合并）

```go
func (sk *Sketch) Merge(other *Sketch) *Sketch {
    m := *sk
    for i := range m.reg {
        if m.reg[i] < other.reg[i] {
            m.reg[i] = other.reg[i]
        }
    }
    return &m
}
```

- 创建一个新 `Sketch`（拷贝 self），然后把对应位置的寄存器取“更大”的那个（这也是 LogLog 合并集合的方式：在同一个桶里，你只要保留最大的前导零值，因为越大的前导零值代表更大的哈希分布跨度）。
- 最终返回合并好的新 `Sketch`。
- 这样合并后的 `Sketch` 就相当于**两个集合的并集**对应的 HLL(或者 HMH)状态。

---

### 3.5 `Similarity`（相似度 / Jaccard Index）

```go
func (sk *Sketch) Similarity(other *Sketch) float64 {
    var C, N float64
    for i := range sk.reg {
        if sk.reg[i] != 0 && sk.reg[i] == other.reg[i] {
            C++
        }
        if sk.reg[i] != 0 || other.reg[i] != 0 {
            N++
        }
    }
    if C == 0 {
        return 0
    }

    n := float64(sk.Cardinality())
    m := float64(other.Cardinality())
    ec := sk.approximateExpectedCollisions(n, m)

    if C < ec {
        return 0
    }
    return (C - ec) / N
}
```

- 这里的思路：
  - **C**：有多少桶（寄存器）在两边都相同（相同且非 0）。
  - **N**：两边中至少有一边非 0 的桶数。
  - 直接用 `C/N` 还不够，因为 HyperMinHash 的原理里寄存器“相同”有一定概率是碰撞，需要扣除一下期望碰撞数 `ec`。
  - `ec` 使用 `approximateExpectedCollisions(n, m)` 做估算，然后再做 \((C - ec)/N\)。
  - 如果 `C < ec`，那就说明碰撞数量比“相同桶”还大，估计相似度为 0。

#### `approximateExpectedCollisions`

```go
func (sk *Sketch) approximateExpectedCollisions(n, m float64) float64 {
    if n < m {
        n, m = m, n
    }
    if n > math.Pow(2, math.Pow(2, q)+r) {
        return math.MaxUint64
    } else if n > math.Pow(2, p+5) {
        d := (4 * n / m) / math.Pow((1+n)/m, 2)
        return c*math.Pow(2, p-r)*d + 0.5
    } else {
        return sk.expectedCollision(n, m) / float64(p)
    }
}
```

- 对不同数量级的 \(n\)、\(m\)（即两个集合的基数），用不同的公式来近似碰撞数量。
- `expectedCollision` 是最精细的计算方式；`c` 是个实验/理论得出的修正系数。

#### `expectedCollision`

```go
func (sk *Sketch) expectedCollision(n, m float64) float64 {
    var x, b1, b2 float64
    for i := 1.0; i <= _2q; i++ {
        for j := 1.0; j <= _2r; j++ {
            if i != _2q {
                den := math.Pow(2, p+r+i)
                b1 = (_2r + j) / den
                b2 = (_2r + j + 1) / den
            } else {
                den := math.Pow(2, p+r+i-1)
                b1 = j / den
                b2 = (j + 1) / den
            }
            prx := math.Pow(1-b2, n) - math.Pow(1-b1, n)
            pry := math.Pow(1-b2, m) - math.Pow(1-b1, m)
            x += (prx * pry)
        }
    }
    return (x * float64(p)) + 0.5
}
```

- 这部分是把 \((i, j)\) 两个坐标对应的区域视为哈希区间，计算它们的碰撞概率，然后累加。
- 这里有点类似把**“哈希空间 + 前导零计数 + 签名”**当作一个二维概率分布积分，最后得到预期碰撞数。

---

### 3.6 `Intersection`（交集大小）

```go
func (sk *Sketch) Intersection(other *Sketch) uint64 {
    sim := sk.Similarity(other)
    return uint64((sim * float64(sk.Merge(other).Cardinality()) + 0.5))
}
```

- 先用 `Similarity` 算两者的估计 Jaccard 相似度 \(S\)。
- 再用 `sk.Merge(other).Cardinality()` 估计**并集**的大小 \(U\)。
- 交集大小约为 \(S \times U\)。因为  
  \[
  \text{Jaccard}(A, B) = \frac{|A \cap B|}{|A \cup B|}  
   \quad\Longrightarrow\quad  
   |A \cap B| = \text{Jaccard}(A, B) \times |A \cup B|.
  \]
- 再加个 0.5（类似四舍五入），最后强转成 `uint64`。

---

## 4. `Hash128` 函数

```go
func Hash128(buffer []byte, seed uint64) (uint64, uint64) {
    const (
        k0 = 0xC83A91E1
        k1 = 0x8648DBDB
        k2 = 0x7BDEC03B
        k3 = 0x2F5870A5
    )
    ...
    return v[0], v[1]
}
```

- 一个自定义的 128 位哈希函数：

  1. 分块地从 `buffer` 里读出 64/32/16/8 位数据；
  2. 不断地做混合、旋转位运算（`rotate_right`）、乘常数等；
  3. 最终输出两个 64 位值 `v[0], v[1]`。

- 这个哈希函数本身和常见的 xxHash / MurmurHash / SipHash 原理类似，通过混合和旋转实现较好的分布性。

```go
func rotate_right(v uint64, k uint) uint64 {
    return (v >> k) | (v << (64 - k))
}
```

- 标准的右旋转操作。

---

## 5. 主函数中的调用逻辑

```go
func main() {
    sk1 := NewHyperMinHash()
    sk2 := NewHyperMinHash()

    for i := 0; i < 10000; i++ {
        sk1.Add([]byte(strconv.Itoa(i)))
    }
    sk1.Cardinality() // 估计数量 ~10000

    for i := 3333; i < 23333; i++ {
        sk2.Add([]byte(strconv.Itoa(i)))
    }
    sk2.Cardinality()     // 估计数量 ~20000
    sk1.Similarity(sk2)   // 估计相似度
    sk1.Intersection(sk2) // 估计两者的交集（~6667）

    sk1.Merge(sk2)
    sk1.Cardinality() // 合并后估计的并集大小 ~30000（实际是 23333）
}
```

1. 创建两个 `Sketch`：`sk1` / `sk2`。
2. 分别往里插入不同范围的数字（转成字符串、再转字节），
   - `sk1` = {0, 1, 2, ..., 9999}
   - `sk2` = {3333, 3334, ..., 23332}
   - 其中重复的部分是 {3333, ..., 9999}，大小为 `10000 - 3333 = 6667`。
3. 调 `Cardinality()` 得到基数估计值；
4. 调 `Similarity(sk2)` 得到近似的 Jaccard 相似度 \(\approx \frac{6667}{10000 + 20000 - 6667} \approx 0.286\)；
5. 调 `Intersection(sk2)` 就是 `Jaccard * unionSize`。
6. `sk1.Merge(sk2)` 返回的是并集的新 sketch，之后再取 `Cardinality()`，就能估计两者合并后的基数 \(\approx 23333\)。

---

# 总结

1. **核心数据结构**：用一个大小为 \(2^p\) 的 `register` 数组，每个 `register` 16 位里分成两部分：
   - **lz**（q位）记录哈希值的“前导零数”。
   - **sig**（r位）记录一个签名，用来在相似度、交集估计时处理碰撞。
2. **插入/合并**：
   - 插入时，根据哈希值计算“应插入哪一个桶（下标）”，再计算“前导零数 lz”与“签名 sig”，并和当前值作比较，只保留大的。
   - 合并时，对应位置取“最大”的 `register`。
3. **Cardinality**：
   - 通过“空桶数 + \(\sum \frac{1}{2^{lz}}\)”等信息，套 LogLog-Beta 公式估计基数。
   - `beta(ez)` 是修正函数。
4. **Similarity / Intersection**：
   - 扫描两组 `register`，数“非零且相同”的桶数 `C`，并结合预计碰撞数 `ec` 算出 Jaccard 相似度 \(\frac{C - ec}{N}\)。
   - 交集则用 `sim * unionCardinality` 得到。
5. **Hash128**：
   - 自定义的 128 位哈希，用分块 + 旋转混合来打散输入数据，输出两个 64 位哈希值给 `AddHash`。

这样就实现了一个**可做基数估计、并集、相似度、交集估计**的简易 HyperMinHash/LogLog-Beta 算法。
