package main

// 在线二分图检测.
type BipartiteChecker struct {
	n            int32
	parent       []int32
	rank         []int32
	color        []int32
	version      int32
	firstViolate int32
	history      []*setValueStep // plugin
}

func NewBipartiteChecker(n int32) *BipartiteChecker {
	res := &BipartiteChecker{
		n:            n,
		parent:       make([]int32, n),
		rank:         make([]int32, n),
		color:        make([]int32, n),
		firstViolate: -1,
	}
	for i := int32(0); i < n; i++ {
		res.parent[i] = i
	}
	return res
}

func (b *BipartiteChecker) IsBipartite() bool {
	return b.firstViolate == -1
}

// (leader, color)
func (b *BipartiteChecker) Find(x int32) (int32, int32) {
	if x == b.parent[x] {
		return x, 0
	}
	leader, color := b.Find(b.parent[x])
	color ^= b.color[x]
	return leader, color
}

func (b *BipartiteChecker) Union(x, y int32) {
	b.version++
	color := int32(1)
	leaderX, distX := b.Find(x)
	x, color = leaderX, color^distX
	leaderY, distY := b.Find(y)
	y, color = leaderY, color^distY
	if x == y {
		if color == 1 && b.firstViolate == -1 {
			b.firstViolate = b.version
		}
		b.setValue(&b.parent[0], b.parent[0])
		return
	}
	if b.rank[x] < b.rank[y] {
		b.setValue(&b.parent[x], y)
		b.setValue(&b.color[x], color)
	} else {
		b.setValue(&b.parent[y], x)
		b.setValue(&b.color[y], color)
		if b.rank[x] == b.rank[y] {
			b.setValue(&b.rank[x], b.rank[x]+1)
		}
	}
}

func (b *BipartiteChecker) Undo() {
	if len(b.history) == 0 {
		return
	}
	v := b.history[len(b.history)-1].version
	if b.firstViolate == v {
		b.firstViolate = -1
	}
	for len(b.history) > 0 && b.history[len(b.history)-1].version == v {
		b.history[len(b.history)-1].Revert()
		b.history = b.history[:len(b.history)-1]
	}
}

func (b *BipartiteChecker) setValue(cell *int32, newValue int32) {
	step := newSetValueStep(cell, *cell, b.version)
	*cell = newValue // apply
	b.history = append(b.history, step)
}

type setValueStep struct {
	cell     *int32
	oldValue int32
	version  int32
}

func newSetValueStep(cell *int32, oldValue int32, version int32) *setValueStep {
	return &setValueStep{cell: cell, oldValue: oldValue, version: version}
}

// func (s *SetValueStep) Apply()  { *s.cell = s.newValue }
func (s *setValueStep) Revert() { *s.cell = s.oldValue }
