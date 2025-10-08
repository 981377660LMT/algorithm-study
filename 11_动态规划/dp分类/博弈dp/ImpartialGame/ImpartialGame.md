https://nyaannyaan.github.io/library/game/impartial-game.hpp

---

### C++ 代码详细讲解

这段 C++ 代码实现了一个通用的、用于解决**不偏游戏 (Impartial Game)** 的求解器。不偏游戏是指在任何局面下，可以执行的合法操作仅由局面本身决定，而与轮到哪位玩家无关。例如，尼姆游戏 (Nim)、井字棋等都是不偏游戏。

这个求解器基于 **Sprague-Grundy 定理**。该定理指出，任何不偏游戏都可以等效于一个特定大小的尼姆堆。这个等效的尼姆堆的大小被称为该局面的 **Grundy 数** (或 nimber)。

- Grundy 数为 0 的局面是 **P-position** (Previous player winning)，即先手必败局面。
- Grundy 数非 0 的局面是 **N-position** (Next player winning)，即先手必胜局面。

代码的核心思想是：通过递归和记忆化搜索，计算出游戏中每个局面的 Grundy 数。

#### 1. 模板参数

```cpp
template <typename Board, typename Move = void, bool splittable = false>
```

- `typename Board`: 代表游戏**盘面**的类型。它可以是任何可比较的类型（因为要用作 `map` 的键），例如 `int`, `string`, `vector<int>` 或自定义结构体。
- `typename Move = void`: 代表一次**着手**（移动）的类型。
  - 如果只是想判断胜负，不需要知道具体怎么走，可以使用默认的 `void`。
  - 如果想找出最佳着手，需要定义一个能表示着手的类型，例如 `pair<int, int>`。
- `bool splittable = false`: 游戏是否是**可分割的**。
  - `false` (默认): 游戏是一个整体，比如井字棋。
  - `true`: 游戏可以分解为多个独立的子游戏，比如尼姆游戏。这种情况下，整个游戏的 Grundy 数是所有子游戏 Grundy 数的**异或和 (XOR sum)**。

#### 2. 核心类型别名

```cpp
using Game = conditional_t<splittable, vector<Board>, Board>;
using State = conditional_t<is_void_v<Move>, Game, pair<Game, Move>>;
using States = vector<State>;
using F = function<States(Board)>;
```

- `Game`: 根据 `splittable` 的值，定义一个“游戏”的类型。如果是可分割的，一个游戏就是一堆盘面 (`vector<Board>`)；否则，就是一个盘面 (`Board`)。
- `State`: 定义一个“状态”的类型。如果 `Move` 是 `void`，一个状态就是一个 `Game`；否则，它是一个 `pair<Game, Move>`，同时包含了后继局面和到达该局面的着手。
- `States`: 所有可能的后继状态的集合 (`vector<State>`)。
- `F`: 一个函数类型，代表游戏的**转移函数**。它接收一个盘面 `Board`，返回所有可能的后继状态 `States`。这是用户需要为特定游戏提供的核心逻辑。

#### 3. 成员变量

```cpp
map<Board, Nimber> mp;
F f;
```

- `mp`: 一个 `map`，用作**记忆化缓存**。它存储已经计算过的 `Board` 的 Grundy 数 (`Nimber`，即 `long long`)，避免重复计算。
- `f`: 用户提供的游戏转移函数。

#### 4. `get(const T& t)` 方法

这是计算 Grundy 数的核心递归函数。它使用 `if constexpr` 根据输入类型 `T` 执行不同的逻辑。

1.  **输入是 `Board`**:
    - 这是递归的基础情况。
    - 首先检查缓存 `mp` 中是否存在该 `Board` 的 Grundy 数。如果存在，直接返回。
    - 如果不存在，调用私有方法 `_get(b)` 进行计算，将结果存入缓存，然后返回。
2.  **输入是 `Boards` (即 `vector<Board>`)**:
    - 这对应 `splittable = true` 的情况。
    - 根据 Sprague-Grundy 定理，总 Grundy 数是所有子盘面 Grundy 数的异或和。
    - 遍历所有子盘面，递归调用 `get()` 计算它们的 Grundy 数，然后将结果全部异或起来。
3.  **输入是 `pair<Game, Move>`**:
    - 这对应 `Move != void` 的情况。
    - 我们只关心后继局面的 Grundy 数，所以直接对 `t.first` (即 `Game`) 递归调用 `get()`。

#### 5. `_get(const Board& b)` 方法

这是实际计算单个盘面 Grundy 数的函数。

1.  调用转移函数 `f(b)` 获取所有可能的后继状态 `gs`。
2.  对 `gs` 中的每一个状态，递归调用 `get()` 计算出它们的 Grundy 数，存入 `ns` 数组。
3.  `ns` 数组现在包含了所有后继局面的 Grundy 数集合。
4.  计算这个集合的 **mex (Minimum Excluded value)**，即不包含在该集合中的最小非负整数。
    - 为了高效计算 mex，代码先对 `ns` 进行排序和去重。
    - 然后遍历 `ns`，找到第一个 `ns[i] != i` 的 `i`，这个 `i` 就是 mex。
    - 如果 `ns` 是 `0, 1, 2, ..., k-1`，那么 mex 就是 `k` (即 `ns.size()`)。
5.  返回的 mex 值就是当前盘面 `b` 的 Grundy 数。

#### 6. `get_best_move(...)` 和 `change_x(...)` 方法

这两个方法仅在 `Move` 不是 `void` 时有效，用于找出**最佳着手**。

- **前提**: 当前局面必须是先手必胜的 (Grundy 数 `n != 0`)。
- **逻辑**: 一个最佳着手，是能将当前局面转移到一个 Grundy 数为 0 的后继局面的着手。根据异或的性质 `a ^ b = c <=> a ^ c = b`，如果当前 Grundy 数为 `n`，我们需要找到一个后继局面，其 Grundy 数为 `g`，使得 `n ^ g = 0`，即 `g = n`。但这是错的，我们是想让对手面对的局面的异或和为 0。
  - 正确的逻辑是：如果当前总 Grundy 数为 `n`，我们需要找到一个后继局面，其 Grundy 数为 `g`，使得 `n ^ g` 最小。对于必胜态，总能找到一个 `g` 使得 `n ^ g < n`。最佳策略是移动到 Grundy 数为 `g'` 的局面，使得 `n ^ g'` 为 0。不，这是对于组合游戏。对于单个游戏，是移动到 Grundy 数为 0 的局面。
  - 让我们再读一下代码：`if (get(st) == (x ^ n))`，其中 `x` 是 `get_best_move` 传入的 `n`，`n` 是当前局面的 Grundy 数。所以 `x^n` 就是 `n^n = 0`。**因此，代码的逻辑是寻找一个能转移到 Grundy 数为 0 的后继状态的着手。**
- `change_x(b, x)`: 尝试在盘面 `b` 上找到一个着手，使得新局面的 Grundy 数 `g` 满足 `g = get(b) ^ x`。在 `get_best_move` 中，`x` 被设为当前局面的 Grundy 数 `n`，所以它寻找的是能转移到 Grundy 数为 `n ^ n = 0` 的着手。
- `get_best_move(t)`:
  - 计算当前局面 `t` 的 Grundy 数 `n`。
  - 断言 `n != 0`，因为败局没有最佳着手。
  - 调用 `change_x(t, n)` 寻找那个能将 Grundy 数变为 0 的着手并返回。
  - 对于可分割游戏，它会遍历所有子游戏，尝试在其中一个子游戏上执行一个操作，使得整个游戏的异或和变为 0。

---

### Go 语言翻译与实现

Go 语言没有 C++ 那样强大的模板元编程和泛型能力（尤其是在 Go 1.18 之前）。Go 1.18 引入了泛型，但仍然无法完全复制 C++ `conditional_t` 这样的编译期类型选择。

因此，我们将采用更符合 Go 语言习惯的方式来实现，使用接口 (interface) 和结构体来模拟这种泛化能力。

我们将定义一个 `ImpartialGame` 接口，任何具体游戏都需要实现它。

```go
package main

import (
	"fmt"
	"sort"
)

// Board 代表任何可比较的盘面状态。
// Go 泛型约束它必须是可比较的，这样才能用作 map 的键。
type Board comparable

// Move 代表任何类型的着手。
type Move any

// Nimber 是 Grundy 数的类型。
type Nimber int64

// State 封装了一个后继局面和到达该局面的着手。
type State[B Board, M Move] struct {
	Game  any // 可以是 B (单个盘面) 或 []B (分割后的游戏)
	Move M
}

// ImpartialGame 接口定义了一个不偏游戏需要提供的核心逻辑。
// 用户需要为自己的游戏实现这个接口。
type ImpartialGame[B Board, M Move] interface {
	// NextStates 接收一个盘面，返回所有可能的后继状态。
	NextStates(board B) []State[B, M]
	// IsSplittable 指示游戏是否可分割。
	IsSplittable() bool
}

// ImpartialGameSolver 是通用的求解器。
type ImpartialGameSolver[B Board, M Move] struct {
	memo map[B]Nimber
	game ImpartialGame[B, M]
}

// NewImpartialGameSolver 创建一个新的求解器实例。
func NewImpartialGameSolver[B Board, M Move](game ImpartialGame[B, M]) *ImpartialGameSolver[B, M] {
	return &ImpartialGameSolver[B, M]{
		memo: make(map[B]Nimber),
		game: game,
	}
}

// getGameNimber 计算一个 "Game" (单个盘面或盘面切片) 的 Grundy 数。
func (s *ImpartialGameSolver[B, M]) getGameNimber(game any) Nimber {
	switch g := game.(type) {
	case B: // 单个盘面
		return s.GetNimber(g)
	case []B: // 盘面切片 (可分割游戏)
		var n Nimber = 0
		for _, board := range g {
			n ^= s.GetNimber(board)
		}
		return n
	default:
		// 不应该发生
		panic("invalid game type")
	}
}

// GetNimber 计算单个盘面的 Grundy 数 (带记忆化)。
func (s *ImpartialGameSolver[B, M]) GetNimber(board B) Nimber {
	// 检查缓存
	if nimber, ok := s.memo[board]; ok {
		return nimber
	}

	// 计算
	nimber := s.calculateNimber(board)
	s.memo[board] = nimber
	return nimber
}

// calculateNimber 是计算 Grundy 数的核心逻辑 (mex)。
func (s *ImpartialGameSolver[B, M]) calculateNimber(board B) Nimber {
	nextStates := s.game.NextStates(board)
	if len(nextStates) == 0 {
		return 0 // 终止局面，Grundy 数为 0
	}

	nextNimbers := make(map[Nimber]struct{})
	for _, state := range nextStates {
		nextNimbers[s.getGameNimber(state.Game)] = struct{}{}
	}

	// 计算 mex
	var mex Nimber = 0
	for {
		if _, exists := nextNimbers[mex]; !exists {
			return mex
		}
		mex++
	}
}

// GetBestMove 查找最佳着手。
// 对于不可分割游戏，返回一个 Move。
// 对于可分割游戏，返回一个 pair {subgame_index, Move}。
// 如果是败局，会 panic。
func (s *ImpartialGameSolver[B, M]) GetBestMove(game any) any {
	totalNimber := s.getGameNimber(game)
	if totalNimber == 0 {
		panic("No best move in a losing position.")
	}

	if !s.game.IsSplittable() {
		// 不可分割游戏
		board := game.(B)
		for _, state := range s.game.NextStates(board) {
			if s.getGameNimber(state.Game) == 0 {
				return state.Move
			}
		}
	} else {
		// 可分割游戏
		boards := game.([]B)
		for i, subBoard := range boards {
			subNimber := s.GetNimber(subBoard)
			// 我们需要找到一个移动，将 subBoard 的 nimber 变为 subNimber'
			// 使得 totalNimber ^ subNimber ^ subNimber' = 0
			// 即 subNimber' = totalNimber ^ subNimber
			targetNimber := totalNimber ^ subNimber
			if targetNimber < subNimber { // 必须移动到更小的 Grundy 数
				for _, state := range s.game.NextStates(subBoard) {
					if s.getGameNimber(state.Game) == targetNimber {
						return struct {
							Index int
							Move  M
						}{i, state.Move}
					}
				}
			}
		}
	}

	panic("Error in GetBestMove: logic failed.")
}

// --- 示例：尼姆游戏 ---

// NimMove 定义了尼姆游戏的着手：从第 i 堆拿走 k 个。
type NimMove struct {
	Index int
	Count int
}

// NimGame 实现了 ImpartialGame 接口。
// 盘面是 []int，代表每堆石子的数量。
type NimGame struct{}

func (ng *NimGame) IsSplittable() bool {
	return true // 尼姆游戏是可分割的
}

// NextStates 对于尼姆游戏来说，盘面是整个 []int。
// 但我们的求解器是基于单个 Board 的，所以这里传入的 board 是一个 int (一堆石子)。
func (ng *NimGame) NextStates(board int) []State[int, NimMove] {
	states := []State[int, NimMove]{}
	// 从一堆数量为 board 的石子中，可以拿走 1 到 board 个
	for i := 1; i <= board; i++ {
		// 移动后，这堆石子剩下 board - i 个
		// 注意：这里的 Move 实际上没有意义，因为我们不知道这是第几堆。
		// 最佳着手是在 GetBestMove 中计算的。
		states = append(states, State[int, NimMove]{
			Game: board - i,
			Move: NimMove{Count: i},
		})
	}
	return states
}

func main() {
	nimGame := &NimGame{}
	// 注意：Board 是 int, Move 是 NimMove
	solver := NewImpartialGameSolver[int, NimMove](nimGame)

	// 游戏局面：3堆，分别有 3, 4, 5 个石子
	initialBoard := []int{3, 4, 5}

	// 计算总 Grundy 数
	nimber := solver.getGameNimber(initialBoard)
	fmt.Printf("Initial board %v has Nimber: %d\n", initialBoard, nimber)

	if nimber == 0 {
		fmt.Println("It's a losing position (P-position).")
	} else {
		fmt.Println("It's a winning position (N-position).")
		bestMove := solver.GetBestMove(initialBoard)
		fmt.Printf("Best move is: %+v\n", bestMove)
	}
}
```

#### Go 版本的设计说明

1.  **接口替代模板特化**: 我们使用 `ImpartialGame` 接口来强制用户提供游戏逻辑，这取代了 C++ 中的 `function<...>` 成员。
2.  **泛型参数**: `Board` 和 `Move` 使用了 Go 1.18+ 的泛型。`Board` 被约束为 `comparable`，以确保它可以作为 `map` 的键。
3.  **显式处理可分割性**: C++ 的 `conditional_t` 在编译期选择类型。在 Go 中，我们通过 `IsSplittable()` 方法和 `any` (或 `interface{}`) 在运行时进行类型断言来处理。`State.Game` 字段被定义为 `any`，在 `getGameNimber` 中通过 `switch game.(type)` 来区分是单个盘面还是盘面切片。
4.  **最佳着手**: `GetBestMove` 的返回值是 `any`，因为对于可分割和不可分割的游戏，返回值的结构不同。调用者需要根据游戏类型进行类型断言。
5.  **示例**: 提供了一个 `NimGame` 的完整示例，展示了如何实现 `ImpartialGame` 接口并使用求解器。在尼姆游戏的例子中，`Board` 是 `int`（一堆石子的数量），而整个游戏状态是 `[]int`。求解器通过对每个 `int` 调用 `GetNimber` 并将结果异或来处理整个游戏。
