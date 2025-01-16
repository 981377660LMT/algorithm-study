下面这份代码演示了**Morton 编码（Z-Order 曲线）的多维实现**，即把多维坐标（例如 2D、3D、4D）**交错插入**到一个 64 位整数里，用于在一维空间中**近似保留**多维坐标的局部性。这在计算机图形、空间索引（R 树、Quad 树等）或大数据处理（类似 Hilbert 曲线、Z-order curve 等）中很常见。

我会按照以下顺序进行讲解：

1. **概念简介**：什么是 Morton 编码(Z-Order)；为什么需要“交错插入”。
2. **核心类型**：`Morton`、`Table`、`Bit` 及其作用。
3. **多维表构建流程**：`CreateTables`、`CreateTable`、`InterleaveBits`。
4. **编码/解码实现**：`Encode`、`Decode` 函数原理。
5. **其他辅助**：`MakeMagic` 用于构造解码时的“magic bits”。
6. **main 函数示例**。

---

## 1. 概念简介

**Morton 编码(Z-Order)** 的主要思想是：

- 给定一个 \(d\)-维坐标 \((x_1, x_2, \dots, x_d)\)。
- 把每个坐标用二进制表示，并将它们的位交错插入到一个单一的整数(通常 64 位)中。例如，2D 情况下把 `x` 的第 0 位放在结果第 0 位、`y` 的第 0 位放在结果第 1 位、`x` 的第 1 位放在结果第 2 位、`y` 的第 1 位放在结果第 3 位……形成一个“Z 字形”。
- 这样得到的 64 位整数称为**Morton code**或**Z-order value**。在空间索引中，可以利用 Morton code 来保留一定的空间局部性——相邻坐标在一维排序后尽量保持在一起。

在这份代码里，为了支持**多个维度**，也就是由 `dimensions` 决定插入的方式，还需要一些**查找表**(`Tables`)和**magic bits** 来帮助加速或简化插入/拆分位的操作。

---

## 2. 核心类型

### 2.1 `Morton` 结构

```go
type Morton struct {
	Dimensions uint8
	Tables     []Table
	Magic      []uint64
}
```

- **`Dimensions`**：维度数，如 2、3、4。
- **`Tables`**：针对每个维度，存储了一个 `Table` 用于快速编码。
- **`Magic`**：用来协助**解码**（从 Morton code 反推多维坐标），是一系列的掩码或拆位信息。

### 2.2 `Table` 结构

```go
type Table struct {
	Index  uint8
	Length uint32
	Encode []Bit
}
```

- 一个 `Table` 对应一维，比如 `Index=0` 可能对应 x 坐标，`Index=1` 可能对应 y 坐标，等等。
- `Length`：这张表支持的坐标最大值 + 1。比如 `Length=512` 表示可编码到 511。
- `Encode`：是个 `[]Bit`，其中下标对应“具体的坐标值”，比如 `Encode[10]` 存放了**如何**将坐标值 10 的二进制插入到 Morton 码中。

### 2.3 `Bit` 结构

```go
type Bit struct {
	Index uint32
	Value uint64
}
```

- `Index`：这里表面上记录了在 `Table.Encode` 里的下标，也就是坐标的实际值；
- `Value`：这个 64 位整数就是“该坐标值对应的交错结果 bits”，其实就是把坐标值的二进制按特定方式插入了位置(shift)后所得。

举例：若 `Index=10` (二进制 `1010`) 在 3D 里，这个 `Value` 可能是一串交错过后的位，如 `(bits for x=? y=? z=? )`。实际上它只对单一维度来说是一段**插入**(interleave)好了的 bits(留了空给其他维度)。

---

## 3. 多维表构建流程

### 3.1 `Create` 方法

```go
func (m *Morton) Create(dimensions uint8, size uint32) {
	done := make(chan struct{})
	mch := make(chan []uint64)

	go func() {
		m.CreateTables(dimensions, size)
		done <- struct{}{}
	}()
	go func() {
		mch <- MakeMagic(dimensions)
		m.Magic = MakeMagic(dimensions)
	}()
	m.Magic = <-mch
	close(mch)
	<-done
	close(done)
}
```

- 并发地做了两件事：
  1. `m.CreateTables(dimensions, size)`：为每个维度生成一个 `Table`；
  2. `m.Magic = MakeMagic(dimensions)`：生成解码用的“magic bits”。
- `dimensions` 表示多少维度，`size` 指定每个维度的坐标最大值(不严格是最大，但决定了表的长度)。
- 这里用两个 goroutine 并发完成后再汇总。

### 3.2 `CreateTables`

```go
func (m *Morton) CreateTables(dimensions uint8, length uint32) {
	ch := make(chan Table)

	m.Dimensions = dimensions
	for i := uint8(0); i < dimensions; i++ {
		go func(i uint8) {
			ch <- CreateTable(i, dimensions, length)
		}(i)
	}
	for i := uint8(0); i < dimensions; i++ {
		t := <-ch
		m.Tables = append(m.Tables, t)
	}
	close(ch)

	sort.Sort(ByTable(m.Tables))
}
```

- 对每个维度 `i`，异步地构建一个 `Table`；
- 然后把它们收集进 `m.Tables` 并按 `Index` 排序。

### 3.3 `CreateTable`

```go
func CreateTable(index, dimensions uint8, length uint32) Table {
	t := Table{Index: index, Length: length}
	bch := make(chan Bit)

	// Build interleave queue
	for i := uint32(0); i < length; i++ {
		go func(i uint32) {
			bch <- InterleaveBits(i, uint32(index), uint32(dimensions-1))
		}(i)
	}
	// Pull results
	for i := uint32(0); i < length; i++ {
		ib := <-bch
		t.Encode = append(t.Encode, ib)
	}
	close(bch)

	sort.Sort(ByBit(t.Encode))
	return t
}
```

- **核心**：对该维度 `index`，对 0..(length-1) 的每个坐标值 `i` 调用 `InterleaveBits(i, index, dimensions-1)`。
- 得到一个 `Bit{Index:i, Value: someMortonBits}`, 放入 `t.Encode[i]`。
- 最后对 `t.Encode` 排序(根据 `Bit.Index` 升序)，使得 `Encode[i]` 对应坐标值 `i`。

### 3.4 `InterleaveBits`

```go
func InterleaveBits(value, offset, spread uint32) Bit {
	ib := Bit{value, 0}

	// n是value，用来确定需要多少位来存储value
	n := value
	limit := uint64(0)
	for i := uint32(0); n != 0; i++ {
		n = n >> 1
		limit++
	}

	// offset = dimension index
	// spread = dimensions-1
	v, o, s := uint64(value), uint64(offset), uint64(spread)
	for i := uint64(0); i < limit; i++ {
		// 把v的第i位插入到 ib.Value 的合适位置
		// i * s 表示 每个bit要被空出(s-1)位, 并为这维度提供空间
		ib.Value |= (v & (1 << i)) << (i * s)
	}
	ib.Value = ib.Value << o

	return ib
}
```

- `value` 表示该维度坐标值(如 10)。
- `offset` = 该维度的索引(如 x=0, y=1, z=2...)。
- `spread` = `dimensions - 1`。
- 循环把 `value` 的第 `i` 位拷贝到 `ib.Value`。在多维情况下，每个 bit 之间要空 `spread` 位(或者说插入时，把它往左移 `(i * s)` 位)；最后再统一左移 `o` 来把它真正放到合适维度的“起始偏移”位置。

概念上，如果我们只有 2 维：

- `offset=0` 对应 x。每个 bit 往 `i*1` 移位，后再左移 `0` => 这样就把 x 的 bits 放在偶数位置(0,2,4,...)；
- `offset=1` 对应 y。则在 `InterleaveBits` 里会把 y 的 bits 放在奇数位置(1,3,5,...)。

对更多维度以此类推，只是空更多位。

这就**构造了**一个表：`Encode[value].Value` = 这维度的坐标值 bits 插入后(保留空间给其他维度)得到的 partial Morton bits。

---

## 4. 编码(`Encode`)和解码(`Decode`)

### 4.1 `Encode`

```go
func (m *Morton) Encode(vector []uint32) (result uint64, err error) {
	// vector长度不可超过m.Tables数量
	if len(vector) > len(m.Tables) {
		return 0, errors.New("Input vector slice length exceeds ...")
	}

	// 对于每个维度k, 坐标=vector[k]
	for k, v := range vector {
		if v > uint32(len(m.Tables[k].Encode)-1) {
			// 超过该table的预期范围
			return 0, errors.New("Input vector component ... exceeds the table size")
		}
		// 读取table中事先算好的bits
		result |= m.Tables[k].Encode[v].Value
	}

	return
}
```

- 给定一个多维坐标 `vector` (如 `[511, 472, 103, 7]`)，逐维去 `m.Tables[k].Encode[ coordinate ]` 拿到那部分**预先插入好的 bits**，然后用位或(`|=`)合并到最终的 64 位 `result`。
- 因为**不同行**维度的 bits 不会冲突(在 `InterleaveBits` 里已经预留位置)，合并就得到了完整 Morton code。

### 4.2 `Decode`

```go
func (m *Morton) Decode(code uint64) (result []uint32) {
	if m.Dimensions == 0 {return nil}

	d := uint64(m.Dimensions)
	r := make([]uint64, d)

	for i := uint64(0); i < d; i++ {
		// 1) 初步提取: 取 code 中从第 i bit开始, 间隔 d bits
		r[i] = (code >> i) & m.Magic[0]
		// 2) 依次使用 Magic[] 修正
		for j := uint64(0); int(j) < len(m.Magic)-1; j++ {
			r[i] = (r[i] ^ (r[i] >> ((d - 1) * (1 << j)))) & m.Magic[j+1]
		}
		result = append(result, uint32(r[i]))
	}
	return
}
```

- **`m.Magic`** 存了一系列掩码(`nth`数组)，用于从 `code` 中解开某一维度的 bits。
- 大致的流程是：
  1. `(code >> i) & m.Magic[0]`：先取出**那些**与该维度相关的位(把 x0, x1, x2... 或 y0, y1, y2... 等抽取出来)。
  2. 再经过若干次“向右移 + 异或 + &mask”操作，把这些分散在大的位间隔中的位**收拢**到低位连续位置。
  3. 最终得到 `uint32(r[i])` 即维度 i 的坐标值。

这段需要对 bit tricks 比较熟悉，其原理常被称为**解交错**(de-interleave)操作，用多轮移位/掩码把分散的 bits 拼回到一个连续的整数里。

---

## 5. `MakeMagic` 函数

```go
func MakeMagic(dimensions uint8) []uint64 {
	d := uint64(dimensions)
	limit := 64/d + 1
	nth := []uint64{0, 0, 0, 0, 0, 0}
	for i := uint64(0); i < limit; i++ {
		switch {
		case i <= 32:
			nth[0] |= 1 << (i * d)
			fallthrough
		case i <= 16:
			nth[1] |= 3 << (i * (d << 1))
			// ...
		}
		// ...
	}
	return nth
}
```

- 它生成若干**掩码**放进 `nth[]`，用于上面提到的 `Decode` 多轮移位和掩码运算。
- 不同 size(1,2,4,8,16,...) 对应的是**不同步长**(例如把偶数位取出来，再把隔 2 位的再组合起来...等)。
- 这个实现逻辑比较复杂，大致就是**对多维度情况**做了“1 bits, 2 bits, 4 bits, 8 bits, ...”的掩码，以在 `Decode` 中按层次拆解位。

这种方法在 2D/3D 的 Morton code decode 中也很常见，如 `_pdep_u64/_pext_u64` 指令在 x86-64 上也可实现硬件加速。

---

## 6. main 示例

```go
func main() {
	// Create a new Morton
	m := new(Morton)
	m.Create(4, 512)

	// coordinates
	c := []uint32{511, 472, 103, 7}

	// Encode
	e, err := m.Encode(c)
	// ...

	fmt.Printf("Coordinates: %v\n", c)
	fmt.Printf("Encoded Coordinates: %v\n", e)
	fmt.Printf("Decoded Coordinates: %v\n", m.Decode(e))
}
```

- 创建一个 `Morton` 对象，指定**维度=4**，且查找表大小=512。
- 准备了一个 4D 坐标 `[511, 472, 103, 7]`。
- `Encode` 得到一个 64 位整数 `e`；
- `Decode(e)` 再还原回 `[511, 472, 103, 7]`。

由此演示了**多维度**到**单个 64 位数**的互相转换过程。

---

## 结论

1. 该库使用了**预计算**的 `Tables` 以及**magic bits** 来快速完成多维 Morton 编码与解码。
2. `Encode` 流程：对每个维度 `k`，根据坐标值 `v` 去 `m.Tables[k].Encode[v].Value` 拿到已经**预插入**好的 partial bits，用 OR 合并到 `result`。
3. `Decode` 流程：用 `MakeMagic(dim)` 产生的一组掩码，通过多轮移位与掩码，将插入后的 bits 分离回单个维度的二进制。
4. 该实现支持**任意维度**，但**受限于 64 位**：坐标分辨率受 `64 / dimensions` 的影响。
5. main 函数展示了 4 维( 0~511 )的用例，并打印出坐标与 Morton code 的互相转换。

这种 Z-Order(Morton code) 对大规模的空间或多维数据，有着**较好的局部性**与**易于实现**的特性，在很多空间索引或计算几何中都会见到类似实现。
