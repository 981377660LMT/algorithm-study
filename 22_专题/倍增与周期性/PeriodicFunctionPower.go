// PeriodicFunctionPower
// 周期函数的幂(k次转移后的状态)

package main

// https://leetcode.cn/problems/prison-cells-after-n-days/submissions/
func prisonAfterNDays(cells []int, n int) []int {
	s0 := [8]int{cells[0], cells[1], cells[2], cells[3], cells[4], cells[5], cells[6], cells[7]}
	F := NewPeriodicFunctionPower(s0, func(cur State) State {
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

type State = [8]int
type PeriodicFunctionPower struct {
	Offset int // 周期开始的位置
	Cycle  int // 周期长度
	data   []State
	used   map[State]int
}

// NewPeriodicFunctionPower 创建周期函数的幂(k次转移后的状态).
//  s0 初始状态.
//  next 状态转移函数.
func NewPeriodicFunctionPower(s0 State, next func(State) State) *PeriodicFunctionPower {
	res := &PeriodicFunctionPower{
		data: []State{},
		used: map[State]int{},
	}
	for _, ok := res.used[s0]; !ok; _, ok = res.used[s0] {
		res.used[s0] = len(res.data)
		res.data = append(res.data, s0)
		s0 = next(s0)
	}
	res.Offset = res.used[s0]
	res.Cycle = len(res.data) - res.Offset
	return res
}

// Query 查询k次转移后的状态(第k项).
//  k>=0.
func (p *PeriodicFunctionPower) Query(k int) State {
	index := k
	if k >= p.Offset {
		index = (k-p.Offset)%p.Cycle + p.Offset
	}
	return p.data[index]
}
