https://nyaannyaan.github.io/library/game/partisan-game.hpp

这段代码实现了一个用于解决**党派游戏**的通用求解器。与之前的不偏游戏（Impartial Game）不同，党派游戏的特点是：**在同一个局面下，不同玩家（通常称为 Left 和 Right）可以执行的操作集合是不同的。**

为了给这类游戏赋值，数学家约翰·康威（John Conway）引入了**超现实数**的概念。每个游戏局面都可以被赋予一个超现实数作为其“值”或“温度”。

- 如果值 > 0，则对 Left 玩家有利。
- 如果值 < 0，则对 Right 玩家有利。
- 如果值 = 0，则对后手玩家有利（先手必败）。

代码由两个核心部分组成：`SurrealNumber` 结构体和 `PartisanGameSolver` 模板类。

#### 1. `SurrealNumber` 结构体

这个结构体实现了超现实数的一个子集——**二进分数 (Dyadic Rationals)**。这是因为在有限步内能构造出的所有超现实数都可以表示为 `p / 2^q` 的形式，其中 p 是整数，q 是非负整数。

- **数据成员**:

  ```cpp
  i64 p, q; // 代表 p / 2^q
  ```

- **核心运算**:

  - `normalize`: 将分数化为最简形式。例如，`6 / 2^3` (即 6/8) 会被化简为 `3 / 2^2` (即 3/4)。
  - `operator+`, `operator-`: 实现加减法。通过通分（将 `q` 统一为较大的那个）来完成。
  - `operator<`, `operator==`, etc.: 实现比较。`a < b` 等价于 `b - a > 0`，而一个数 `s = p / 2^q` 是否大于 0，仅取决于其分子 `p` 是否大于 0。

- **游戏相关的核心函数**:
  - `child()`: 返回一个数的“左右子节点”。根据康威的构造法，一个数 `x` 可以看作是 `{ L | R }` 的形式，其中 `L` 是所有比 `x` “简单”且小于 `x` 的数，`R` 是所有比 `x` “简单”且大于 `x` 的数。`child()` 返回的就是最接近 `x` 的左右两个“子”选项。
  - `larger()`: 对于一个数 `x`，找到比它大的最“简单”的数。它通过从 0 开始，不断取右子节点直到超过 `x` 来实现。
  - `smaller()`: 与 `larger()` 相对，找到比 `x` 小的最“简单”的数。
  - `reduce(l, r)`: **这是最关键的函数**。给定一个左边界 `l` 和一个右边界 `r` (满足 `l < r`)，它寻找位于 `(l, r)` 区间内的**最简单**的超现实数。这里的“简单”指的是构造它所需的天数最少（或者说，分母中的 `q` 最小）。这完美对应了游戏 `{ L | R }` 的值的定义：`G = { G_L | G_R }` 的值是满足 `max(G_L) < x < min(G_R)` 的最简单的数 `x`。

#### 2. `PartisanGameSolver` 模板类

这是一个通用的党派游戏求解器。

- **模板参数**:

  - `typename Game`: 代表游戏盘面的类型，需要是可比较的（用作 `map` 的键）。
  - `typename F`: 一个函数类型，代表游戏的**转移函数**。它接收一个 `Game` 盘面，返回一个 `pair<vector<Game>, vector<Game>>`。`pair` 的第一个元素是 **Left** 玩家的所有后继局面，第二个元素是 **Right** 玩家的所有后继局面。

- **成员变量**:

  - `map<Game, S> mp`: 记忆化缓存，存储已计算过的盘面及其对应的超现实数值。
  - `F f`: 用户提供的转移函数。

- **核心方法**:
  - `get(Game g)`: 公共接口，带记忆化地计算盘面 `g` 的值。
  - `_get(Game g)`: 私有方法，实现实际的计算逻辑。
    1.  调用 `f(g)` 获取 Left 的后继局面 `gl` 和 Right 的后继局面 `gr`。
    2.  如果 `gl` 和 `gr` 都为空，说明是终止局面，值为 0。
    3.  递归调用 `get()` 计算 `gl` 和 `gr` 中所有局面的超现实数值，分别存入 `l` 和 `r`。
    4.  找到 Left 后继局面的最大值 `sl = max(l)` 和 Right 后继局面的最小值 `sr = min(r)`。
    5.  **处理边界情况**:
        - 如果 Right 没有走法 (`r` 为空)，游戏的值就是比 `sl` 大的最简单的数，即 `sl.larger()`。
        - 如果 Left 没有走法 (`l` 为空)，游戏的值就是比 `sr` 小的最简单的数，即 `sr.smaller()`。
    6.  **一般情况**: 游戏的值就是 `(sl, sr)` 区间内最简单的数，通过 `reduce(sl, sr)` 计算得出。

---

### Go 语言翻译与实现

Go 语言缺乏 C++ 的模板和操作符重载，因此我们需要用更明确的结构体和方法来模拟。

#### 1. `SurrealNumber` 的 Go 实现

我们将创建一个 `SurrealNumber` 结构体，并为其定义各种方法。

```go
package main

import (
	"fmt"
	"math"
	"math/bits"
)

// SurrealNumber 表示一个二进分数 p / 2^q。
type SurrealNumber struct {
	p int64 // 分子
	q uint  // 分母是 2^q
}

// NewSurrealNumber 创建一个新的超现实数。
func NewSurrealNumber(p int64, q uint) SurrealNumber {
	return SurrealNumber{p, q}
}

// Normalize 将分数化为最简形式。
func (s SurrealNumber) Normalize() SurrealNumber {
	if s.p == 0 {
		return SurrealNumber{0, 0}
	}
	// trailing zeros in p
	tz := uint(bits.TrailingZeros64(uint64(s.p)))
	if tz > s.q {
		tz = s.q
	}
	return SurrealNumber{s.p >> tz, s.q - tz}
}

// String 用于打印。
func (s SurrealNumber) String() string {
	if s.q == 0 {
		return fmt.Sprintf("%d", s.p)
	}
	return fmt.Sprintf("%d / %d", s.p, int64(1)<<s.q)
}

// Add 实现加法。
func (s SurrealNumber) Add(other SurrealNumber) SurrealNumber {
	cq := s.q
	if other.q > cq {
		cq = other.q
	}
	cp := (s.p << (cq - s.q)) + (other.p << (cq - other.q))
	return SurrealNumber{cp, cq}.Normalize()
}

// Sub 实现减法。
func (s SurrealNumber) Sub(other SurrealNumber) SurrealNumber {
	cq := s.q
	if other.q > cq {
		cq = other.q
	}
	cp := (s.p << (cq - s.q)) - (other.p << (cq - other.q))
	return SurrealNumber{cp, cq}.Normalize()
}

// Neg 返回相反数。
func (s SurrealNumber) Neg() SurrealNumber {
	return SurrealNumber{-s.p, s.q}
}

// LessThan 比较大小。
func (s SurrealNumber) LessThan(other SurrealNumber) bool {
	return other.Sub(s).p > 0
}

// EqualTo 比较是否相等。
func (s SurrealNumber) EqualTo(other SurrealNumber) bool {
	return s.Sub(other).p == 0
}

// Child 返回左右子节点。
func (s SurrealNumber) Child() (SurrealNumber, SurrealNumber) {
	if s.p == 0 {
		return NewSurrealNumber(-1, 0), NewSurrealNumber(1, 0)
	}
	if s.q == 0 && s.p > 0 {
		return NewSurrealNumber(s.p, 0).Add(NewSurrealNumber(-1, 0)),
			NewSurrealNumber(s.p+1, 0)
	}
	if s.q == 0 && s.p < 0 {
		return NewSurrealNumber(s.p-1, 0),
			NewSurrealNumber(s.p, 0).Add(NewSurrealNumber(1, 1))
	}
	return s.Sub(NewSurrealNumber(1, s.q+1)), s.Add(NewSurrealNumber(1, s.q+1))
}

// Larger 找到比 s 大的最简单的数。
func (s SurrealNumber) Larger() SurrealNumber {
	root := NewSurrealNumber(0, 0)
	for !s.LessThan(root) { // while s >= root
		_, rr := root.Child()
		root = rr
	}
	return root
}

// Smaller 找到比 s 小的最简单的数。
func (s SurrealNumber) Smaller() SurrealNumber {
	root := NewSurrealNumber(0, 0)
	for !root.LessThan(s) { // while root >= s
		lr, _ := root.Child()
		root = lr
	}
	return root
}

// Reduce 找到 (l, r) 区间内最简单的数。
func Reduce(l, r SurrealNumber) SurrealNumber {
	if !l.LessThan(r) {
		panic("l must be less than r for Reduce")
	}
	root := NewSurrealNumber(0, 0)
	for !l.LessThan(root) || !root.LessThan(r) { // while l >= root or root >= r
		lr, rr := root.Child()
		if !root.LessThan(r) { // if root >= r
			root = lr
		} else {
			root = rr
		}
	}
	return root
}
```

#### 2. `PartisanGameSolver` 的 Go 实现

我们将使用泛型（Go 1.18+）和接口来构建求解器。

```go
package main

// GameState 代表任何可比较的游戏盘面状态。
type GameState comparable

// GameLogic 定义了党派游戏需要提供的核心逻辑。
type GameLogic[G GameState] interface {
	// NextStates 返回 Left 和 Right 玩家的后继局面。
	NextStates(g G) (left []G, right []G)
}

// PartisanGameSolver 是通用的党派游戏求解器。
type PartisanGameSolver[G GameState] struct {
	memo  map[G]SurrealNumber
	logic GameLogic[G]
}

// NewPartisanGameSolver 创建一个新的求解器实例。
func NewPartisanGameSolver[G GameState](logic GameLogic[G]) *PartisanGameSolver[G] {
	return &PartisanGameSolver[G]{
		memo:  make(map[G]SurrealNumber),
		logic: logic,
	}
}

// Get 计算一个盘面的超现实数值。
func (s *PartisanGameSolver[G]) Get(g G) SurrealNumber {
	if val, ok := s.memo[g]; ok {
		return val
	}
	val := s.calculate(g)
	s.memo[g] = val
	return val
}

func (s *PartisanGameSolver[G]) calculate(g G) SurrealNumber {
	gl, gr := s.logic.NextStates(g)

	if len(gl) == 0 && len(gr) == 0 {
		return NewSurrealNumber(0, 0)
	}

	var leftValues, rightValues []SurrealNumber
	for _, cg := range gl {
		leftValues = append(leftValues, s.Get(cg))
	}
	for _, cg := range gr {
		rightValues = append(rightValues, s.Get(cg))
	}

	var sl, sr SurrealNumber
	if len(leftValues) > 0 {
		sl = leftValues[0]
		for i := 1; i < len(leftValues); i++ {
			if sl.LessThan(leftValues[i]) {
				sl = leftValues[i]
			}
		}
	}
	if len(rightValues) > 0 {
		sr = rightValues[0]
		for i := 1; i < len(rightValues); i++ {
			if rightValues[i].LessThan(sr) {
				sr = rightValues[i]
			}
		}
	}

	if len(rightValues) == 0 {
		return sl.Larger()
	}
	if len(leftValues) == 0 {
		return sr.Smaller()
	}

	if !sl.LessThan(sr) {
		panic("game rule violation: max of left options is not less than min of right options")
	}

	return Reduce(sl, sr)
}

// --- 示例：一个简单的党派游戏 ---

// SimpleGameLogic 实现一个简单的游戏。
// 盘面是一个整数。
// Left 可以将 n 变为 n-2。
// Right 可以将 n 变为 n-1。
// 游戏在 n <= 0 时结束。
type SimpleGameLogic struct{}

func (sgl *SimpleGameLogic) NextStates(g int) (left []int, right []int) {
	// Left's move
	if g-2 > 0 {
		left = append(left, g-2)
	}
	// Right's move
	if g-1 > 0 {
		right = append(right, g-1)
	}
	return
}

func main() {
	logic := &SimpleGameLogic{}
	solver := NewPartisanGameSolver[int](logic)

	// 计算盘面为 5 时的值
	val := solver.Get(5)
	fmt.Printf("Value of game at state 5 is: %s\n", val) // 预期是 1/2

	val2 := solver.Get(6)
	fmt.Printf("Value of game at state 6 is: %s\n", val2) // 预期是 1
}
```

#### Go 版本的设计说明

1.  **接口与泛型**: 我们定义了 `GameLogic` 接口来解耦求解器和具体游戏逻辑，这类似于 C++ 中传递函数对象 `F`。求解器 `PartisanGameSolver` 使用泛型 `[G GameState]` 来接受任何可比较的盘面类型。
2.  **无操作符重载**: Go 不支持操作符重载，因此所有的数学和比较运算都实现为结构体的方法，例如 `s.Add(other)` 和 `s.LessThan(other)`。这使得代码更冗长，但意图清晰。
3.  **错误处理**: C++ 版本在 `reduce` 和 `_get` 中使用了 `assert`。在 Go 版本中，我们用 `panic` 来替代，因为这些断言失败通常意味着游戏规则本身存在矛盾或程序逻辑有误，属于不可恢复的错误。
4.  **最大/最小值**: Go 标准库没有直接获取 slice 最大/最小值的函数，所以我们用一个简单的循环来实现。
5.  **示例**: 提供了一个 `SimpleGameLogic` 的例子来演示如何使用这个求解器。这个例子清晰地展示了党派游戏的非对称性：Left 和 Right 的移动规则不同。
