package main

func main() {
	uk := NewUndoKnapsack(10)
	op := uk.Add(1, 2)
	op.Apply()
	println(uk.Query(1) == 2)
	op.Undo()
	println(uk.Query(1) == -INF)
}

type Operation struct{ apply, undo func() }

func NewOperation(apply, undo func()) *Operation { return &Operation{apply: apply, undo: undo} }
func (op *Operation) Apply()                     { op.apply() }
func (op *Operation) Undo()                      { op.undo() }

const INF int = 2e18

// 可撤销01背包.
type UndoKnapsack struct {
	dp []int
}

func NewUndoKnapsack(size int32) *UndoKnapsack {
	dp := make([]int, size+1)
	for i := int32(0); i <= size; i++ {
		dp[i] = -INF
	}
	dp[0] = 0
	return &UndoKnapsack{dp: dp}
}

// 添加一个物品，重量为 weight，价值为 value.
func (uk *UndoKnapsack) Add(weight int, value int) *Operation {
	next := make([]int, len(uk.dp))
	return NewOperation(
		func() {
			for i := 0; i < len(uk.dp); i++ {
				next[i] = uk.dp[i]
				if i >= weight {
					next[i] = max(next[i], uk.dp[i-weight]+value)
				}
			}
			uk.dp, next = next, uk.dp
		},
		func() {
			uk.dp, next = next, uk.dp
		},
	)
}

// 背包容量为 capacity 时的最大价值.
func (uk *UndoKnapsack) Query(capacity int) int {
	return uk.dp[capacity]
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
