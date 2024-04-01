package main

import "fmt"

func main() {
	n := int32(10)
	dsu := NewUndoDSU(n)
	operation := dsu.Union(1, 2)
	operation.Apply()
	fmt.Println(dsu.Size(1) == 2)
	operation.Undo()
	fmt.Println(dsu.Size(1) == 1)
	fmt.Println(dsu.Size(2) == 1)
}

type IOperation interface {
	Apply()
	Undo()
}

type Operation struct{ apply, undo func() }

func NewOperation(apply, undo func()) *Operation { return &Operation{apply: apply, undo: undo} }
func (op *Operation) Apply()                     { op.apply() }
func (op *Operation) Undo()                      { op.undo() }

type UndoDSU struct {
	rank, parent []int32
}

func NewUndoDSU(n int32) *UndoDSU {
	rank := make([]int32, n)
	parent := make([]int32, n)
	for i := int32(0); i < n; i++ {
		rank[i] = 1
		parent[i] = -1
	}
	return &UndoDSU{rank: rank, parent: parent}
}

func (dsu *UndoDSU) Find(x int32) int32 {
	for dsu.parent[x] != -1 {
		x = dsu.parent[x]
	}
	return x
}

func (dsu *UndoDSU) Size(x int32) int32 {
	return dsu.rank[dsu.Find(x)]
}

func (dsu *UndoDSU) Union(a, b int32) IOperation {
	var x, y int32
	// to NamedTuple
	return NewOperation(
		func() {
			x = dsu.Find(a)
			y = dsu.Find(b)
			if x == y {
				return
			}
			if dsu.rank[x] < dsu.rank[y] {
				x, y = y, x
			}
			dsu.parent[y] = x
			dsu.rank[x] += dsu.rank[y]
		},
		func() {
			cur := y
			for dsu.parent[cur] != -1 {
				cur = dsu.parent[cur]
				dsu.rank[cur] -= dsu.rank[y]
			}
			dsu.parent[y] = -1
		},
	)
}
