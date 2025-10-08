[群论](https://oi-wiki.org/math/algebra/basic/)

https://nyaannyaan.github.io/library/math/semiring-linear-recursive.hpp

### 一、详细讲解

这份代码库定义了一个在算法竞赛中非常强大的抽象工具：**半环（Semiring）**。

#### 1. 什么是半环 (Semiring)？

半环是一种代数结构，可以看作是“环 (Ring)”的弱化版。一个环我们知道有加减乘，而半环则更加宽松。

一个集合 `R` 和两种运算 `+`（加法）和 `·`（乘法）构成一个半环，需要满足以下条件：

1.  **`(R, +)` 是一个可换独异点 (Commutative Monoid)**

    - **封闭性**: `a, b` 在 `R` 中，则 `a + b` 也在 `R` 中。
    - **结合律**: `(a + b) + c = a + (b + c)`。
    - **单位元**: 存在一个“加法单位元” `0`，使得 `a + 0 = a`。
    - **交换律**: `a + b = b + a`。
    - **注意**: 它不要求有逆元（减法）。

2.  **`(R, ·)` 是一个独异点 (Monoid)**

    - **封闭性**: `a, b` 在 `R` 中，则 `a · b` 也在 `R` 中。
    - **结合律**: `(a · b) · c = a · (b · c)`。
    - **单位元**: 存在一个“乘法单位元” `1`，使得 `a · 1 = a`。
    - **注意**: 不要求有逆元（除法），也不要求交换律。

3.  **乘法对加法满足分配律 (Distributive Law)**

    - `a · (b + c) = (a · b) + (a · c)`
    - `(a + b) · c = (a · c) + (b · c)`

4.  **加法单位元 `0` 是乘法的零元 (Annihilator)**
    - `a · 0 = 0 · a = 0`。

#### 2. 半环有什么用？—— 矩阵乘法的泛化

标准矩阵乘法的定义是 `C[i][j] = Σ (A[i][k] * B[k][j])`。
这里的 `Σ` (求和) 和 `*` (相乘) 正是普通算术中的“加法”和“乘法”。

半环的强大之处在于，我们可以**替换**这里的“加法”和“乘法”为半环中定义的任意 `+` 和 `·` 运算，只要它们满足半环的公理，那么所有基于矩阵乘法的算法（如矩阵快速幂）就依然成立！

**最重要的例子：最短路径问题**

- 考虑 **Min-Plus 半环** (也叫热带半环 Tropical Semiring):
  - 集合 `R`: 实数加上 `+∞`
  - “加法” `+` : 定义为 `min` 函数。`a + b = min(a, b)`。
  - “乘法” `·` : 定义为 `+` 函数。`a · b = a + b`。
- 我们来验证一下：
  - 加法单位元 `0`: `min(a, +∞) = a`，所以这里的 `0` 是 `+∞`。
  - 乘法单位元 `1`: `a + 0 = a`，所以这里的 `1` 是 `0`。
  - 分配律: `a + min(b, c) = min(a+b, a+c)`，成立。
- 现在，把这套运算代入矩阵乘法公式：
  `C[i][j] = min_k (A[i][k] + B[k][j])`
- 这个公式是不是非常眼熟？这正是 **Floyd-Warshall 算法**的核心步骤！如果 `A` 是图的邻接矩阵（边权），那么 `A²` 的 `(i, j)` 元素就表示从 `i` 到 `j` 经过最多一条中转边的最短路径。`A^k` 就表示路径长度最多为 `k` 的最短路。
- **结论**: 图的所有顶点对最短路径问题，可以被抽象为在 Min-Plus 半环上的矩阵乘法问题。我们可以用**矩阵快速幂**在 `O(N³ log k)` 的时间内求出路径长度恰好为 `k` 的最短路。

#### 3. C++ 代码解析

- **`template <typename T, T (*add)(T, T), T (*mul)(T, T), T (*I0)(), T (*I1)()> struct semiring`**

  - 这是一个非常巧妙的 C++ 模板设计。它通过模板参数接受了：
    - `T`: 半环中元素的底层数据类型 (如 `long long`, `pair`)。
    - `add`, `mul`: 实现半环加法和乘法的**函数指针**。
    - `I0`, `I1`: 返回加法和乘法单位元的**函数指针**。
  - 这种设计使得 `semiring` 结构体完全泛化，可以适用于任何满足半环定义的运算。

- **`struct Mat`**
  - 这是一个泛型矩阵类，它的模板参数 `rig` 就是上面定义的 `semiring` 类型。
  - 它重载了 `+` 和 `*` 运算符，其实现完全依赖于 `rig` 提供的 `+=` 和 `*=`。
  - `operator*=` 的实现 `C[i][j] += A[i][k] * B[k][j]` 看似普通，但实际上这里的 `+=` 和 `*` 都是 `semiring` 中定义的、被泛化了的运算。
  - `operator^=` 实现了矩阵快速幂，同样适用于任何半环。

### 二、Go 语言翻译与实现

Go 语言没有 C++ 那样的模板和运算符重载，因此我们无法做到完全相同的语法糖。但是，我们可以通过**接口 (Interface)** 来实现同样的核心思想：定义一套行为规范，让不同的结构体去实现它。

#### 1. 定义接口

首先，我们定义一个 `Semiring` 接口，它规定了一个半环类型必须具备的行为。

```go
package main

// Semiring 接口定义了半环需要满足的行为。
// T 是泛型参数，代表半环中元素的基础类型。
type Semiring[T any] interface {
	// Add 执行半环的“加法”操作
	Add(other T) T
	// Mul 执行半环的“乘法”操作
	Mul(other T) T
	// Value 返回元素的底层值
	Value() T
}

// SemiringFactory 接口定义了如何创建半环的单位元。
// 这是因为单位元是类型的属性，而不是实例的属性。
type SemiringFactory[T any] interface {
	// Zero 返回加法单位元 ("0")
	Zero() T
	// One 返回乘法单位元 ("1")
	One() T
	// New 从一个底层值创建一个新的半环元素
	New(val T) T
}
```

#### 2. 实现具体的半环

现在，我们来实现文档中提到的 **Min-Plus 半环**。

```go
// ...existing code...
import "math"

// --- Min-Plus Semiring 实现 ---

const infLL = int64(math.MaxInt64 / 2) // 防止溢出

// MinPlus 实现了 Semiring[int64] 接口
type MinPlus struct {
	val int64
}

func (mp MinPlus) Add(other MinPlus) MinPlus {
	return MinPlus{min(mp.val, other.val)}
}

func (mp MinPlus) Mul(other MinPlus) MinPlus {
	return MinPlus{mp.val + other.val}
}

func (mp MinPlus) Value() int64 {
	return mp.val
}

// MinPlusFactory 实现了 SemiringFactory[MinPlus] 接口
type MinPlusFactory struct{}

func (f MinPlusFactory) Zero() MinPlus {
	// Min-Plus 的加法是 min，单位元是无穷大
	return MinPlus{infLL}
}

func (f MinPlusFactory) One() MinPlus {
	// Min-Plus 的乘法是 +，单位元是 0
	return MinPlus{0}
}

func (f MinPlusFactory) New(val int64) MinPlus {
	return MinPlus{val}
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}
```

#### 3. 实现泛型矩阵

接下来，我们创建一个泛型矩阵结构体，它可以在任何实现了我们定义的 `Semiring` 接口的类型上工作。

```go
// ...existing code...

// Mat 是一个泛型矩阵，可以在任何半环类型上工作。
type Mat[T Semiring[T]] struct {
	A       [][]T
	Factory SemiringFactory[T] // 需要一个工厂来创建单位元
	N       int
}

// NewMat 创建一个新的 N x N 矩阵，并用加法单位元填充。
func NewMat[T Semiring[T]](n int, factory SemiringFactory[T]) *Mat[T] {
	m := &Mat[T]{
		A:       make([][]T, n),
		Factory: factory,
		N:       n,
	}
	zero := factory.Zero()
	for i := 0; i < n; i++ {
		m.A[i] = make([]T, n)
		for j := 0; j < n; j++ {
			m.A[i][j] = zero
		}
	}
	return m
}

// NewIdentityMat 创建一个 N x N 的单位矩阵。
func NewIdentityMat[T Semiring[T]](n int, factory SemiringFactory[T]) *Mat[T] {
	m := NewMat(n, factory)
	one := factory.One()
	for i := 0; i < n; i++ {
		m.A[i][i] = one
	}
	return m
}

// Mul 执行矩阵乘法。
func (m *Mat[T]) Mul(other *Mat[T]) *Mat[T] {
	c := NewMat(m.N, m.Factory)
	for i := 0; i < m.N; i++ {
		for j := 0; j < m.N; j++ {
			// c.A[i][j] 初始为 zero
			for k := 0; k < m.N; k++ {
				// C[i][j] = C[i][j] + (A[i][k] * B[k][j])
				// 这里的 Add 和 Mul 都是半环定义的运算
				prod := m.A[i][k].Mul(other.A[k][j])
				c.A[i][j] = c.A[i][j].Add(prod)
			}
		}
	}
	return c
}

// Pow 执行矩阵快速幂。
func (m *Mat[T]) Pow(k int64) *Mat[T] {
	res := NewIdentityMat(m.N, m.Factory)
	base := m
	for k > 0 {
		if k&1 == 1 {
			res = res.Mul(base)
		}
		base = base.Mul(base)
		k >>= 1
	}
	return res
}
```

#### 4. 使用示例

下面是如何使用我们创建的 Min-Plus 半环和矩阵来解决一个简单的最短路问题。

```go
package main

import "fmt"

func main() {
	// 假设有一个3个节点的图，我们想求任意两点间路径长度恰好为2的最短路
	n := 3
	factory := MinPlusFactory{}

	// 创建邻接矩阵
	adj := NewMat(n, factory)
	adj.A = [][]MinPlus{
		{factory.New(0), factory.New(1), factory.New(10)},      // 0->0=0, 0->1=1, 0->2=10
		{factory.New(infLL), factory.New(0), factory.New(2)},   // 1->1=0, 1->2=2
		{factory.New(4), factory.New(infLL), factory.New(0)},   // 2->0=4, 2->2=0
	}

	// 计算 adj^2
	// 在 Min-Plus 半环上，这就是路径长度为2的最短路矩阵
	adj2 := adj.Pow(2)

	fmt.Println("Shortest paths of length 2:")
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			fmt.Printf("%d ", adj2.A[i][j].Value())
		}
		fmt.Println()
	}

	// 例如，从 0 到 2 长度为 2 的最短路是 0->1->2，长度为 1+2=3
	// 矩阵计算结果 adj2.A[0][2] 的值应该是 3
	// min(adj[0][0]+adj[0][2], adj[0][1]+adj[1][2], adj[0][2]+adj[2][2])
	// min(0+10, 1+2, 10+0) = min(10, 3, 10) = 3
}
```

这个 Go 实现通过接口和泛型，成功地复刻了 C++ 库的核心思想：将代数结构（半环）与算法（矩阵乘法/快速幂）解耦，从而可以用一套代码解决多种看似不相关的问题。
