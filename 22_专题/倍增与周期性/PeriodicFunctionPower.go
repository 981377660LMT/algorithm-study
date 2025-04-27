// PeriodicFunctionPower/CycleDecomposition
// 周期函数的幂(k次转移后的状态)

package main

import "fmt"

func main() {
	S := NewPeriodicFunctionPower[int](0, func(cur int) int { return 2 + cur%3 })
	fmt.Println(S.CycleStart, S.CycleLength, S.PreCycle, S.Cycle) // 1 3 [0] [2 4 3]
}

// https://leetcode.cn/problems/prison-cells-after-n-days/submissions/
func prisonAfterNDays(cells []int, n int) []int {
	s0 := [8]int{cells[0], cells[1], cells[2], cells[3], cells[4], cells[5], cells[6], cells[7]}
	type State = [8]int
	F := NewPeriodicFunctionPower[State](s0, func(cur State) State {
		var next State
		for i := 1; i < 7; i++ {
			if cur[i-1] == cur[i+1] {
				next[i] = 1
			}
		}
		return next
	})
	res := F.Query(n)
	return []int{res[0], res[1], res[2], res[3], res[4], res[5], res[6], res[7]}
}

type PeriodicFunctionPower[S comparable] struct {
	CycleStart  int // 周期开始的位置
	CycleLength int // 周期长度
	PreCycle    []S // 周期开始前的状态
	Cycle       []S // 周期内的状态
}

// NewPeriodicFunctionPower 创建周期函数的幂(k次转移后的状态).
//  s0 初始状态.
//  next 状态转移函数.
func NewPeriodicFunctionPower[S comparable](s0 S, next func(S) S) *PeriodicFunctionPower[S] {
	res := &PeriodicFunctionPower[S]{}
	history, visited := make([]S, 0), make(map[S]int)
	for _, ok := visited[s0]; !ok; _, ok = visited[s0] {
		visited[s0] = len(history)
		history = append(history, s0)
		s0 = next(s0)
	}
	res.CycleStart = visited[s0]
	res.CycleLength = len(history) - res.CycleStart
	res.PreCycle = history[:res.CycleStart]
	res.Cycle = history[res.CycleStart:]
	return res
}

// Query 查询k次转移后的状态(第k项).
//  k>=0.
func (p *PeriodicFunctionPower[S]) Query(k int) S {
	if k < p.CycleStart {
		return p.PreCycle[k]
	}
	k -= p.CycleStart
	k %= p.CycleLength
	return p.Cycle[k]
}
