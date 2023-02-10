// UnionFindWithUndoAndWeight
// https://hitonanode.github.io/cplib-cpp/unionfind/undo_monoid_unionfind.hpp
// 可撤销 / 维护 `满足可换律的monoid` 权值 的并查集
// Set: 将下标为index元素`所在集合`的权值置为value.
// Get: 获取下标为index元素`所在集合`的权值.
// Undo: 撤销上一次合并(Union)或者修改权值(Set)操作
// Reset: 撤销所有操作

package main

import "fmt"

func main() {
	uf := NewUndoDSU(10)
	uf.Union(0, 1)
	uf.Union(2, 3)
	fmt.Println(uf.Find(0), uf.Find(1), uf.Find(2), uf.Find(3))
	uf.Union(0, 2)
	fmt.Println(uf.Find(0), uf.Find(1), uf.Find(2), uf.Find(3))
	uf.Undo()
	fmt.Println(uf.Find(0), uf.Find(1), uf.Find(2), uf.Find(3))
}

type S = int

func (*UndoDSU) op(s1, s2 S) S { return s1 + s2 }

func NewUndoDSU(n int) *UndoDSU { return NewUndoDSUWithWeights(make([]S, n)) }
func NewUndoDSUWithWeights(weights []S) *UndoDSU {
	n := len(weights)
	parent, rank, ws := make([]int, n), make([]int, n), make([]S, n)
	for i := 0; i < n; i++ {
		parent[i], rank[i], ws[i] = i, 1, weights[i]
	}
	history := []historyItem{}
	return &UndoDSU{Rank: rank, Part: n, parents: parent, weights: ws, history: history}
}

type historyItem struct {
	a, b int
	c    S
}

type UndoDSU struct {
	Rank    []int // 每个集合的大小
	Part    int   // 集合的个数
	parents []int
	weights []S
	history []historyItem
}

// 将下标为index元素`所在集合`的权值置为value.
func (uf *UndoDSU) Set(index int, value S) {
	index = uf.Find(index)
	uf.history = append(uf.history, historyItem{index, uf.Rank[index], uf.weights[index]})
	uf.weights[index] = value
}

// 获取下标为index元素`所在集合`的权值.
func (uf *UndoDSU) Get(index int) S { return uf.weights[uf.Find(index)] }

// 撤销上一次合并(Union)或者修改权值(Set)操作
func (uf *UndoDSU) Undo() {
	uf.weights[uf.parents[uf.history[len(uf.history)-1].a]] = uf.history[len(uf.history)-1].c
	uf.Rank[uf.parents[uf.history[len(uf.history)-1].a]] = uf.history[len(uf.history)-1].b
	uf.parents[uf.history[len(uf.history)-1].a] = uf.history[len(uf.history)-1].a
	uf.history = uf.history[:len(uf.history)-1]
}

// 撤销所有操作
func (uf *UndoDSU) Reset() {
	for len(uf.history) > 0 {
		uf.Undo()
	}
}

func (uf *UndoDSU) Find(x int) int {
	if uf.parents[x] == x {
		return x
	}
	return uf.Find(uf.parents[x])
}

func (uf *UndoDSU) Union(x, y int) bool {
	x, y = uf.Find(x), uf.Find(y)
	if uf.Rank[x] < uf.Rank[y] {
		x, y = y, x
	}
	uf.history = append(uf.history, historyItem{y, uf.Rank[x], uf.weights[x]})
	if x != y {
		uf.parents[y] = x
		uf.Rank[x] += uf.Rank[y]
		uf.weights[x] = uf.op(uf.weights[x], uf.weights[y])
		uf.Part--
		return true
	}
	return false
}

func (uf *UndoDSU) IsConnected(x, y int) bool { return uf.Find(x) == uf.Find(y) }
