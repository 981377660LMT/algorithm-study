package main

func main() {
	e := func() int { return 0 }
	op := func(a, b int) int { return a + b }
	pow := func(a int, b int) int { return a * b }
	s0 := 0
	next := func(s int) (int, int) {
		return (s + 1) % 10, s + 1
	}
	pf := NewPeriodFunctionWeighted(e, op, pow, s0, next)

	// Example usage
	for k := 0; k < 10; k++ {
		state, totalWeight := pf.Query(k)
		println("State:", state)
		println("Total Weight:", totalWeight)
	}
}

// https://leetcode.cn/problems/prison-cells-after-n-days/submissions/
func prisonAfterNDays(cells []int, n int) []int {
	type State = [8]int

	e := func() int { return 0 }
	op := func(a, b int) int { return a + b }
	pow := func(a int, b int) int { return a * b }
	s0 := [8]int{cells[0], cells[1], cells[2], cells[3], cells[4], cells[5], cells[6], cells[7]}
	nextFn := func(cur State) (State, int) {
		var next State
		for i := 1; i < 7; i++ {
			if cur[i-1] == cur[i+1] {
				next[i] = 1
			}
		}
		return next, 1
	}
	F := NewPeriodFunctionWeighted(e, op, pow, s0, nextFn)
	res, _ := F.Query(n)
	return []int{res[0], res[1], res[2], res[3], res[4], res[5], res[6], res[7]}
}

type PeriodicFunctionPowerWeighted[S comparable, E any] struct {
	CycleStart   int
	CycleLen     int
	PreStates    []S
	CycleStates  []S
	PreWeights   []E
	CycleWeights []E
	CycleSum     E

	e   func() E
	op  func(E, E) E
	pow func(E, int) E
}

// NewPeriodFunctionWeighted constructs from initial state s0, transition next, and weight operations.
// next(s) -> (nextState, weight)
func NewPeriodFunctionWeighted[S comparable, E any](
	e func() E,
	op func(E, E) E,
	pow func(E, int) E,
	s0 S,
	next func(S) (S, E),
) *PeriodicFunctionPowerWeighted[S, E] {
	pf := &PeriodicFunctionPowerWeighted[S, E]{e: e, op: op, pow: pow}
	visited := make(map[S]int)
	states := []S{}
	weights := make([]E, 0)
	cur := s0

	for _, ok := visited[cur]; !ok; _, ok = visited[cur] {
		visited[cur] = len(states)
		states = append(states, cur)
		nxt, w := next(cur)
		weights = append(weights, w)
		cur = nxt
	}

	cycleStart := visited[cur]
	cycleLen := len(states) - cycleStart
	preStates := states[:cycleStart]
	cycleStates := states[cycleStart:]

	preWeights := make([]E, cycleStart+1)
	acc := e()
	preWeights[0] = acc
	for i := 0; i < cycleStart; i++ {
		acc = op(acc, weights[i])
		preWeights[i+1] = acc
	}

	cycleWeights := make([]E, cycleLen+1)
	acc = e()
	cycleWeights[0] = acc
	for i := 0; i < cycleLen; i++ {
		acc = op(acc, weights[cycleStart+i])
		cycleWeights[i+1] = acc
	}

	pf.CycleStart = cycleStart
	pf.CycleLen = cycleLen
	pf.PreStates = preStates
	pf.CycleStates = cycleStates
	pf.PreWeights = preWeights
	pf.CycleWeights = cycleWeights
	pf.CycleSum = acc

	return pf
}

// Query returns (state, totalWeight) after k transitions from the initial state.
func (pf *PeriodicFunctionPowerWeighted[S, E]) Query(k int) (S, E) {
	if k < pf.CycleStart {
		return pf.PreStates[k], pf.PreWeights[k]
	}

	steps := k - pf.CycleStart
	d := steps / pf.CycleLen
	r := steps % pf.CycleLen
	w := pf.PreWeights[pf.CycleStart]
	if d > 0 {
		w = pf.op(w, pf.pow(pf.CycleSum, d))
	}
	w = pf.op(w, pf.CycleWeights[r])
	return pf.CycleStates[r], w
}
