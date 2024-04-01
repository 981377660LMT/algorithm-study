package main

import "fmt"

func main() {
	n := int32(10)
	dsu := NewUndoDSUXor(n)
	fmt.Println(dsu.Xor(1, 2) == 0)
	operation1 := dsu.Union(1, 2, 3)
	operation1.Apply()
	fmt.Println(dsu.Xor(1, 2) == 3)
	operation2 := dsu.Union(2, 3, 4)
	operation2.Apply()
	fmt.Println(dsu.Xor(1, 3) == 7)
	operation2.Undo()
	fmt.Println(dsu.Xor(1, 2) == 3)
	operation1.Undo()
	fmt.Println(dsu.Xor(1, 2) == 0)
}

type Operation struct{ apply, undo func() }

func NewOperation(apply, undo func()) *Operation { return &Operation{apply: apply, undo: undo} }
func (op *Operation) Apply()                     { op.apply() }
func (op *Operation) Undo()                      { op.undo() }

// 可撤销异或并查集.
type UndoDSUXor struct {
	rank     []int32
	parent   []int32
	xor      []int
	conflict bool
}

func NewUndoDSUXor(n int32) *UndoDSUXor {
	rank := make([]int32, n)
	parent := make([]int32, n)
	xor := make([]int, n)
	for i := int32(0); i < n; i++ {
		rank[i] = 1
		parent[i] = -1
	}
	return &UndoDSUXor{rank: rank, parent: parent, xor: xor}
}

func (dsu *UndoDSUXor) XorToRoot(x int32) int {
	res := 0
	for dsu.parent[x] != -1 {
		res ^= dsu.xor[x]
		x = dsu.parent[x]
	}
	return res
}

func (dsu *UndoDSUXor) Xor(a, b int32) int {
	return dsu.XorToRoot(a) ^ dsu.XorToRoot(b)
}

func (dsu *UndoDSUXor) Conflict() bool {
	return dsu.conflict
}

func (dsu *UndoDSUXor) Find(x int32) int32 {
	for dsu.parent[x] != -1 {
		x = dsu.parent[x]
	}
	return x
}

func (dsu *UndoDSUXor) Union(a, b int32, dist int) *Operation {
	var x, y int32
	var conflictSnapshot bool
	return NewOperation(
		func() {
			x, y = dsu.Find(a), dsu.Find(b)
			delta := dsu.XorToRoot(a) ^ dsu.XorToRoot(b) ^ dist
			conflictSnapshot = dsu.conflict
			if x == y {
				dsu.conflict = dsu.conflict || delta != 0
				return
			}
			if dsu.rank[x] < dsu.rank[y] {
				x, y = y, x
			}
			dsu.parent[y] = x
			dsu.xor[y] = delta
			dsu.rank[x] += dsu.rank[y]
		},
		func() {
			for cur := y; dsu.parent[cur] != -1; {
				cur = dsu.parent[cur]
				dsu.rank[cur] -= dsu.rank[y]
			}
			dsu.parent[y] = -1
			dsu.xor[y] = 0
			dsu.conflict = conflictSnapshot
		},
	)
}

func (dsu *UndoDSUXor) Size(x int32) int32 {
	return dsu.rank[dsu.Find(x)]
}
