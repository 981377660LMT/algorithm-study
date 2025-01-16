下面这段代码实现了一个 Elias-Fano 编码/解码结构 (**EliasFano**) 来压缩并迭代有序整数序列。Elias-Fano 是一种常见的针对**单调非降序列**（如倒排索引的文档列表）进行紧凑存储并支持**快速随机访问**的编码方法。

为方便理解，下文会先简要回顾 Elias-Fano 编码原理，再逐行解读代码中最关键的部分，包括各字段的含义、`Compress` 如何编码、`readCurrentValue` 如何解码，以及 `Move/Next` 如何实现随机与顺序访问。

---

## 1. Elias-Fano 原理回顾

给定一个单调非降序列 \(\langle s_1, s_2, \dots, s_n \rangle\)，并且已知所有元素都不超过一个上界 \( U \)（代码中用 `universe` 表示）。Elias-Fano 的做法大致是：

1. 取一个参数 \( L \approx \lfloor \log_2(\frac{U}{n}) \rfloor \)（这里由 `msb(universe / n)` 近似得到）。
2. 对每个元素 \( s_i \)，把它的二进制拆成：
   - **低位 (low bits)**：取 \( L \) 个比特
   - **高位 (high bits)**：剩余的更高部分
3. 把所有元素的低位部分**集中存储**在一片连续的位空间中 (在代码里叫做 `lowerBits` 区)。
4. 把所有元素的高位部分用一种**一元编码**(0 串 + 1 标记)的思路存储到另外一片位空间里 (在代码里叫做“高位区”或通过 `ef.b.Set(high)` 标识出来)，这样可以在解码时快速用“找第 i 个 set bit”的方式恢复第 i 个元素的高位。

并且，为了在代码实现里**简化**对高位的存储，Elias-Fano 会在高位区记录一些「偏移量」，如 \((s_i \gg L) + i + 1\) 之类，从而在解码时把它恢复成实际的高位值。这样做之后，**解码**时的核心步骤就是：

- 找到第 \(i\) 个 `1` 在高位区的下标（做 rank/select）；
- 从对应的“低位区”取出 \( L \) 个比特；
- 合并高位和低位，就可以得到原来的 \( s_i \)。

---

## 2. 主要字段含义

```go
type EliasFano struct {
	universe         uint64 // 序列可能出现的最大值（上界）
	n                uint64 // 序列长度
	lowerBits        uint64 // L, 拆分到“低位”的比特数
	higherBitsLength uint64 // 高位区所需的bit长度(大约是 n + (universe >> lowerBits) + 2)
	mask             uint64 // (1 << lowerBits) - 1, 用于提取/屏蔽低位
	lowerBitsOffset  uint64 // 低位区在整体bitset中的起始位置
	bvLen            uint64 // 整个bitset所占的总bit数
	b                *BitSet// 用于存储高位+低位的位图
	curValue         uint64 // 当前迭代到的值
	position         uint64 // 当前迭代器指向的序列下标(0-based)
	highBitsPos      uint64 // 当前迭代到的高位部分在位图中的位置
}
```

- **universe / n / lowerBits**：根据上文 Elias-Fano 原理来决定如何拆分高/低位。
- **higherBitsLength**：高位区大约需要 \(n + (universe >> L)\) 个比特，再加一些边界。
- **mask**：用来从整数中提取 \(L\) 个低位 (形如 `elem & mask`)。
- **lowerBitsOffset**：在整张位图中，前面那一段用于存储“高位区”（即 set bit 的信息），而从 `lowerBitsOffset` 开始才是真正的“低位区”。
- **bvLen**：bitset 的总长度，= 高位区长度 + \(n \times L\)。
- **curValue / position / highBitsPos**：用于迭代访问时，记录当前值、当前下标，以及在高位区对应的 bit 位置。

---

## 3. 初始化 (`NewEliasFano`)

```go
func NewEliasFano(universe uint64, n uint64) *EliasFano {
	var lowerBits uint64
	if lowerBits = 0; universe > n {
		lowerBits = msb(universe / n)
	}
	higherBitsLength := n + (universe >> lowerBits) + 2
	mask := (uint64(1) << lowerBits) - 1
	lowerBitsOffset := higherBitsLength
	bvLen := lowerBitsOffset + n*uint64(lowerBits)
	b := NewBitSet(uint(bvLen))
	return &EliasFano{
		universe, n, lowerBits, higherBitsLength,
		mask, lowerBitsOffset, bvLen, b,
		0, 0, 0,
	}
}
```

- 先根据 `universe` 和 `n` 估计 `lowerBits`。如果 `universe <= n`，就让它为 0；否则用 `msb(universe / n)`。
  - `msb(x)` 返回 \( \lfloor \log_2(x) \rfloor \)，即最高位所在位置。
- `higherBitsLength` 约等于 `n + (universe >> lowerBits) + 2`；这是 Elias-Fano 中存高位区的一种常规大小估算。
- `lowerBitsOffset` 就是“高位区”所占位数的末端，下标从这里开始就进入“低位区”了。
- `bvLen` = 高位区大小 + 低位区大小(\(n \times lowerBits\))。
- 新建一个 `BitSet` 来容纳这 `bvLen` 位。

---

## 4. `Compress`：往 Elias-Fano 结构里写入序列

```go
func (ef *EliasFano) Compress(elems []uint64) {
	last := uint64(0)
	for i, elem := range elems {
		if i > 0 && elem < last {
			log.Fatal("Sequence is not sorted")
		}
		if elem > ef.universe {
			log.Fatalf("Element %d is greater than universe", elem)
		}

		// 1) 计算高位 (elem >> L) + i + 1
		//    之所以加 i+1, 是为了编码时在bitset中对应位置, 方便后续解码
		high := (elem >> ef.lowerBits) + uint64(i) + 1

		// 2) 计算低位
		low := elem & ef.mask

		// 3) 在高位区, 对应 'high' 的那个位置打上 set bit (1)
		ef.b.Set(uint(high))

		// 4) 把 low bits 写到低位区 offset 处
		offset := ef.lowerBitsOffset + uint64(i)*ef.lowerBits
		setBits(ef.b, offset, low, ef.lowerBits)

		last = elem
		if i == 0 {
			// 初始化迭代状态
			ef.curValue = elem
			ef.highBitsPos = high
		}
	}
}
```

- **高位存法**：`ef.b.Set(uint(high))`，相当于在整个 bitset 的索引为 `high` 的位置置 1。这样对第 i 个元素，我们把 \((elem >> L) + i + 1\) 用来标识它的高位信息。
- **低位存法**：在位置 `offset = lowerBitsOffset + i * L` 里写入 `low`；调用 `setBits` 把这 `L` 位一个个设到 bitset 里。
- 注意这里 `i + 1` 这一步，使得**真正的高位 = ( highBitsPos - i - 1 )**。这样做是 Elias-Fano 编码的一种常见技巧，可以减少额外的 rank/select 开销或简化解码公式。

---

### 4.1 `setBits` 函数

```go
func setBits(b *BitSet, offset uint64, bits uint64, length uint64) {
	for i := uint64(0); i < length; i++ {
		val := bits & (1 << (length - i - 1))
		b.SetTo(uint(offset+i+1), val > 0)
	}
}
```

- 逐位写入：从最高位到最低位，将 `bits` 中的每一位拷贝到 bitset 的 `[offset+1 ... offset+length]` 区域。
- 注意这里的下标是 `offset + i + 1`，在这份实现里所有位似乎都往后再偏移了 1；这是具体实现的细节。

---

## 5. 解码：`readCurrentValue`、`Move`、`Next`

### 5.1 `readCurrentValue`

这是最核心的**解码**操作，专门用于更新 `ef.curValue`：

```go
func (ef *EliasFano) readCurrentValue() {
	pos := uint(ef.highBitsPos)
	if pos > 0 {
		pos++
	}
	pos, _ = ef.b.NextSet(pos)  // (1) 找到下一个 set bit
	ef.highBitsPos = uint64(pos)

	// (2) 从低位区读出 ef.lowerBits 位
	low := uint64(0)
	offset := ef.lowerBitsOffset + ef.position*ef.lowerBits
	for i := uint64(0); i < ef.lowerBits; i++ {
		if ef.b.Test(uint(offset + i + 1)) {
			low++
		}
		low = low << 1
	}
	low = low >> 1

	// (3) 把高位和低位合并
	//     高位 = (ef.highBitsPos - ef.position - 1)
	//     后面 << ef.lowerBits, 再或上 low
	ef.curValue = uint64(((ef.highBitsPos - ef.position - 1) << ef.lowerBits) | low)
}
```

1. **(1) 找下一个 set bit**：通过 `pos, _ = ef.b.NextSet(pos)`，定位到高位区的下一个“1”的位置，赋给 `ef.highBitsPos`。
2. **(2) 读低位**：从 `offset = lowerBitsOffset + position*lowerBits` 开始，读出 `L` 位合成一个 `low`。
3. **(3) 合并**：`ef.curValue = ((highBitsPos - position - 1) << L) | low`
   - 这里 `highBitsPos - position - 1` 就还原了**真正的高位 = (elem >> L)**，因为在写入时做了 `+ i + 1`。

在这个实现里，`NextSet` 等操作来自 `b *BitSet` 结构，它可以找到从 pos 往后下一个置 1 的 bit 的索引。这样就相当于**顺序扫描**高位区，得到第 `position` 个元素对应的 high。

---

### 5.2 `Move(position uint64)`：随机访问

```go
func (ef *EliasFano) Move(position uint64) (uint64, error) {
	if position >= ef.Size() {
		return 0, errors.New("Out of bound")
	}
	if ef.position == position {
		return ef.Value(), nil
	}
	if position < ef.position {
		// 如果往回跳, 直接 Reset 后再前进
		ef.Reset()
	}
	skip := position - ef.position
	pos := uint(ef.highBitsPos)
	for i := uint64(0); i < skip; i++ {
		pos, _ = ef.b.NextSet(pos + 1)
	}
	ef.highBitsPos = uint64(pos - 1)
	ef.position = position
	ef.readCurrentValue()
	return ef.Value(), nil
}
```

- 目标：移动到第 `position` 个元素。
- 如果要往前跳，就直接 `Reset()` 再从头到尾扫描；如果要往后跳，就在高位区循环 `NextSet` 若干次。
- 最后更新 `ef.highBitsPos`、`ef.position` 并调用 `readCurrentValue()` 解得当前值。
- **注意**：这样做随机访问的复杂度是 \(O(position)\)，并不算高效。真正 Elias-Fano 支持可以用**rank/select**数据结构\*\*在 \(O(1)\) 或 \(O(\log n)\) 找到“第 i 个 set bit”，从而实现更快的随机访问。这里的实现是用简单的循环跳过 set bits。

---

### 5.3 `Next()`

```go
func (ef *EliasFano) Next() (uint64, error) {
	ef.position++
	if ef.position >= ef.Size() {
		return 0, errors.New("End reached")
	}
	ef.readCurrentValue()
	return ef.Value(), nil
}
```

- 简单地 `position++`，然后再次 `readCurrentValue()`。
- 因为 `readCurrentValue()` 内部会通过 `NextSet` 找到**下一个** set bit 并取相应的低位，所以就可以顺序迭代到下一个值。

---

## 6. 其它辅助函数

- `Reset()`：把迭代器恢复到开头 (position=0, highBitsPos=0)，然后 `readCurrentValue()` 解出第一个值。
- `Value()`：直接返回当前迭代器位置解码出的值 `curValue`。
- `Size()`：返回 `ef.n`。
- `Bitsize()`：返回实际 bitset 的占用大小。
- `msb(x)`：返回 \(\lfloor \log_2(x) \rfloor\)。

---

## 7. 小结

1. **存储格式**

   - 该实现将 Elias-Fano 的“高位区”与“低位区”放在同一个 `BitSet` 里，通过 `lowerBitsOffset` 划分。
   - 高位区：在下标为 `(elem>>L) + i + 1` 的地方置 1，表示第 i 个元素的高位。
   - 低位区：在 `lowerBitsOffset + i*L` 存第 i 个元素的低位（逐位写）。

2. **编码过程 (`Compress`)**

   - 对于第 `i` 个元素 `elem`，高位 `= (elem >> L) + i + 1` ，低位 `= elem & mask`。
   - 设置相应 bit。

3. **解码过程 (`readCurrentValue`)**

   - 从高位区找当前 position 对应的下一个 set bit 下标 `highBitsPos`；
   - 从低位区读出低位 `low`；
   - 合并得到最终值 `((highBitsPos - position - 1) << L) | low`。

4. **迭代/随机访问**
   - `Next()`：顺序前进一位，通过 `NextSet` 找下一组高位；
   - `Move(pos)`：若往后跳就循环调用 `NextSet(pos+1)`, 若往前就干脆 `Reset()` 再扫；这种方法较简单但效率不高(线性)。

> 如果需要高效的 `O(1)` 或 `O(log n)` 随机访问，就要在高位区的 bitset 上再构建一个**rank/select**索引结构\*\*。本代码仅做了演示性质的实现，用 `NextSet` 顺序去找。

因此，这份代码展示了 Elias-Fano 的**核心编码思路**和**最基本的解码(遍历)方法**。尽管随机访问实现得相对朴素，但已经能完整地将 Elias-Fano 的“高位打标 + 低位集中”原则体现出来。
