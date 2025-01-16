下面给出这段 Go 代码的**详细解读**。它实现了两种**空间填充曲线 (Space-Filling Curve)** 的 2D 变体：**Hilbert 曲线**和**Peano 曲线**，并提供了将一维整数映射到曲线上的二维坐标以及反向映射的功能。代码中还定义了一个通用的 `SpaceFilling` 接口，用于对不同曲线进行统一管理。

---

## 一、背景概念

1. **空间填充曲线 (Space-Filling Curve)**  
   是指一条曲线在其极限状态下能“覆盖”一个多维空间区域。像 Hilbert 曲线、Peano 曲线、Z-order (Morton) 等，都能将一维值与二维/多维坐标互相映射，并在某种程度上保持空间邻近性。这在数据库、图形、GIS 等场景很有用。

2. **Hilbert 曲线 (2D)**

   - 输入一个整数 \(t\in[0, n^2-1]\)，输出二维坐标 \((x, y)\in[0,n-1]\times[0,n-1]\)。
   - 反之给出 \((x, y)\)，也可得到其在 Hilbert 曲线上的序数 \(t\)。
   - 要求 \(n\) 为 2 的幂(例如 2,4,8,16...)，因为 Hilbert 二维实现通常在递归层级中，层级深度= \(\log_2(n)\)。

3. **Peano 曲线 (2D)**
   - 类似，但必须要求 \(n\) 是 3 的幂(例如 3,9,27...)。
   - 该代码中只实现了从一维序数 \(t\) 到 \((x,y)\) 的映射 `Map()`，未完成反向映射 `MapInverse()`。

---

## 二、Hilbert 代码解析

### 2.1 `Hilbert` 结构

```go
type Hilbert struct {
	N int
}
```

- `N`: 二维空间的大小(宽度=高度= N)，必须是**2 的幂**。

### 2.2 构造与元数据

```go
func NewHilbert(n int) (*Hilbert, error) {
	if n <= 0 {
		return nil, ErrNotPositive
	}
	if (n & (n - 1)) != 0 {
		return nil, ErrNotPowerOfTwo
	}
	return &Hilbert{N: n}, nil
}

func (s *Hilbert) GetDimensions() (int, int) {
	return s.N, s.N
}
```

- `NewHilbert(n)`：检查 `n` 是否 > 0 且为 2 的幂。
- `GetDimensions()`：返回 2D 空间大小 `(N,N)`。

### 2.3 一维->二维：`Map(t)`

```go
func (s *Hilbert) Map(t int) (x, y int, err error) {
	if t < 0 || t >= s.N*s.N {
		return -1, -1, ErrOutOfRange
	}
	// 类似于经典Hilbert算法: 逐层处理(t的高层bits和低层bits)
	for i := 1; i < s.N; i *= 2 {
		// t&2==2 => 取出 t 的特定bit判断rx
		// t&1==1 => 取出 t 的特定bit判断ry
		rx := (t & 2) == 2
		ry := (t & 1) == 1
		if rx {
			ry = !ry
		}
		x, y = s.rotate(i, x, y, rx, ry)
		// 根据 rx, ry 决定向 quadrant 偏移
		if rx {
			x += i
		}
		if ry {
			y += i
		}
		t /= 4 // 每层处理2 bits(=4)
	}
	return
}
```

- 关键原理：Hilbert 曲线的离散实现常把 `(rx, ry)` 看做“当前层的象限索引”，然后**旋转 + 偏移**以得到下层的坐标。
- 循环 `i := 1; i < s.N; i *= 2` 表示从最小格到最大格（或从最低层到最高层）依次处理 2 bits：
  - `(t & 2) == 2` => 第 bit 1 决定 `rx`; `(t & 1) == 1` => 第 bit 0 决定 `ry`。
  - 若 `rx==true`，则翻转 `ry`；
  - 调用 `rotate(i, x, y, rx, ry)` 做相应的子方块旋转；
  - 最后若 `rx`, `ry` 为真，x / y 需要加上 `i` 偏移到正确象限。
  - `t /= 4` 去掉已经处理的 2 bits，进入下一层。

### 2.4 二维->一维：`MapInverse(x, y)`

```go
func (s *Hilbert) MapInverse(x, y int) (t int, err error) {
	if x < 0 || x >= s.N || y < 0 || y >= s.N {
		return -1, ErrOutOfRange
	}
	for i := s.N / 2; i > 0; i /= 2 {
		rx := (x & i) > 0
		ry := (y & i) > 0

		a := 0
		if rx {
			a = 3
		}
		// 关键：a ^ b2i(ry) => 组合bits
		t += i * i * (a ^ b2i(ry))

		x, y = s.rotate(i, x, y, rx, ry)
	}
	return
}
```

- 反向过程与 `Map` 相似，每层用 `rx = (x & i)>0`、`ry=(y & i)>0` 判断该层象限，累加到 `t`。
- `a` 初步是 `3`(二进制 11) 若 `rx=true`，否则 0(二进制 00)；然后 `a ^ b2i(ry)` 处理 `ry`。
- `t += i * i * (...)` => 这相当于把 2 bits 写到 `t` 最后(因为每一层处理2 bits)，只是以 `i*i` 做位置控制。
- `x, y = s.rotate(...)` 继续将坐标旋转到下一层原点。

### 2.5 辅助：`rotate(n, x, y, rx, ry)`

```go
func (s *Hilbert) rotate(n, x, y int, rx, ry bool) (int, int) {
	if !ry {
		// ry == false => 需要旋转/翻转
		if rx {
			x = n - 1 - x
			y = n - 1 - y
		}
		// x,y 互换
		x, y = y, x
	}
	return x, y
}
```

- 根据 `rx, ry` 决定是否要把坐标 `(x,y)` 做 90 度旋转或对称变换。
- 这里 `n` 是当前子格大小(例如 1,2,4,...)。

---

## 三、Peano 代码解析

### 3.1 `Peano` 结构

```go
type Peano struct {
	N int // 必须是3的幂
}
```

### 3.2 构造与元数据

```go
func NewPeano(n int) (*Peano, error) {
	if n <= 0 {
		return nil, ErrNotPositive
	}
	if !isPow3(float64(n)) {
		return nil, ErrNotPowerOfThree
	}
	return &Peano{N: n}, nil
}
func (p *Peano) GetDimensions() (int, int) { return p.N, p.N }
```

- 检查 `n>0` 且 `n` 是否 3 的幂(`isPow3`)。

### 3.3 一维->二维：`Map(t)`

```go
func (p *Peano) Map(t int) (x, y int, err error) {
	if t < 0 || t >= p.N*p.N {
		return -1, -1, ErrOutOfRange
	}
	for i := 1; i < p.N; i *= 3 {
		s := t % 9 // 取出基于 9 进制的2 bits(或3 bits)
		rx := int(s / 3) // 3x3小方格里的列
		ry := int(s % 3) // 行
		// 如果 rx==1 => flip vertical
		if rx == 1 {
			ry = 2 - ry
		}

		// 根据 s 决定对 (x,y) 的旋转
		if i > 1 {
			x, y = p.rotate(i, x, y, s)
		}

		x += rx * i
		y += ry * i

		t /= 9 // 已处理2位(这里是 3^2=9个子方格)
	}
	return x, y, nil
}
```

- Peano 曲线以 3 为基，通常把空间分成 `3x3 = 9` 的小方格，再进一步细分，每层处理“mod 9”。
- `(rx, ry)` 表示当前层中 3x3 网格的列(row)与行；然后看 `rx==1` 时要做翻转。
- 可能在大多数 Peano 实现里，会有更复杂的旋转方式，但此处是一个示例性的实现。
- `t /= 9` 进入下一层。
- 注意**这段仅实现了正向映射**，反向(`MapInverse`)还未完成。

### 3.4 辅助：`rotate`

```go
func (p *Peano) rotate(n, x, y, s int) (int, int) {
	if n == 1 { return x, y }

	n = n - 1
	switch s {
	case 0: return x, y
	case 1: return n - x, y        // flip horizontal
	case 2: return x, y
	case 3: return x, n - y        // flip vertical
	case 4: return n - x, n - y    // flip both
	case 5: return x, n - y
	case 6: return x, y
	case 7: return n - x, y
	case 8: return x, y
	}
	panic("should not reach")
}
```

- 根据 s(0~8) 选择不同的翻转组合：
  - 0,2,6,8 => no transform
  - 1,7 => flip horizontal
  - 3,5 => flip vertical
  - 4 => flip both horizontal+vertical
- 这些是作者对 3x3 子方格位置的归纳，从而“拼接”出 Peano 曲线的形状。

---

## 四、接口与辅助

1. **`SpaceFilling` 接口**

   ```go
   type SpaceFilling interface {
       Map(t int) (x, y int, err error)
       MapInverse(x, y int) (t int, err error)
       GetDimensions() (int, int)
   }
   ```

   - Hilbert、Peano 都实现了该接口，故可在外部统一使用。

2. **错误定义**

   ```go
   var (
       ErrNotPositive     = errors.New("N must be greater than zero")
       ErrNotPowerOfTwo   = errors.New("N must be a power of two")
       ErrNotPowerOfThree = errors.New("N must be a power of three")
       ErrOutOfRange      = errors.New("value is out of range")
   )
   ```

3. **`b2i(b bool) int`**  
   小辅助函数，把 `bool` 转成 0 或 1。Hilbert 的反向映射中用于 XOR。

4. **`isPow3(n float64) bool`**  
   循环除以 3 检查是否最终为 1。用浮点实现有些不太精确，但在该场景够用。

---

## 五、main 函数

```go
func main() {
    // Nothing done here in the snippet
}
```

- 当前例子中的 `main()` 并没有实际演示任何用法(只留空)。
- 实际使用中，你可 `NewHilbert(8)` 或 `NewPeano(9)`, 然后调用 `Map(...)` / `MapInverse(...)` 来做坐标转换。

---

## 六、总结

- **Hilbert**：

  - 需要 `N` = 2 的幂；
  - 提供 `Map(t)->(x,y)` 和 `MapInverse(x,y)->t`；
  - 通过逐层提取 bits `(rx, ry)` 并做**子方块旋转**实现。

- **Peano**：

  - 需要 `N` = 3 的幂；
  - 提供 `Map(t)->(x,y)`，但尚未完成 `MapInverse(x,y)->t`；
  - 每层将 `t mod 9` 分解到 3x3 小网格 `(rx,ry)`；可能对 `(x,y)` 做翻转以保持曲线连续性。

- **通用接口** `SpaceFilling`：
  - `Map` / `MapInverse` / `GetDimensions()`。

这段代码示例演示了**如何用离散化的 Hilbert / Peano 曲线**在 2D 空间中进行一维-二维映射。它主要是示例性质，还可以扩展到更高维、或更复杂的实现，但对于基本概念演示和简单应用场景足够了。
