// UnionFindArrayWithUndoAndWeight/UnionFindWithUndoAndWeight
// https://hitonanode.github.io/cplib-cpp/unionfind/undo_monoid_unionfind.hpp
// 可撤销并查集 / 维护 `满足可换律的monoid` 权值 的并查集
// SetGroupWeight: 将下标为index元素`所在集合`的权值置为value.
// GetGroupWeight: 获取下标为index元素`所在集合`的权值.
// !Undo: 撤销上一次合并(Union)或者修改权值(Set)操作，没合并成功也要撤销
// Reset: 撤销所有操作

package main

import "fmt"

func main() {
	uf := NewUnionFindArrayWithUndoAndWeight(make([]int, 5))
	uf.Union(0, 1)
	uf.Union(2, 3)
	fmt.Println(uf.GetGroupWeight(0), uf.GetGroupWeight(1), uf.GetGroupWeight(2), uf.GetGroupWeight(3), uf.GetGroupWeight(4))
	uf.SetGroupWeight(0, 1)
	fmt.Println(uf.GetGroupWeight(0), uf.GetGroupWeight(1), uf.GetGroupWeight(2), uf.GetGroupWeight(3), uf.GetGroupWeight(4))
	uf.SetGroupWeight(0, 2)
	fmt.Println(uf.GetGroupWeight(0), uf.GetGroupWeight(1), uf.GetGroupWeight(2), uf.GetGroupWeight(3), uf.GetGroupWeight(4))
	uf.Undo()
	fmt.Println(uf.GetGroupWeight(0), uf.GetGroupWeight(1), uf.GetGroupWeight(2), uf.GetGroupWeight(3), uf.GetGroupWeight(4))
}

type S = int

func op(s1, s2 S) S { return s1 + s2 }

func NewUnionFindArrayWithUndoAndWeight(weight []S) *UnionFindArrayWithUndoAndWeight {
	n := len(weight)
	parent, rank, ws := make([]int, n), make([]int, n), make([]S, n)
	for i := 0; i < n; i++ {
		parent[i], rank[i], ws[i] = i, 1, weight[i]
	}
	history := []historyItem{}
	return &UnionFindArrayWithUndoAndWeight{Part: n, rank: rank, parent: parent, weight: ws, history: history}
}

type historyItem struct {
	root, rank int
	weight     S
}

type UnionFindArrayWithUndoAndWeight struct {
	Part    int
	rank    []int
	parent  []int
	weight  []S
	history []historyItem
}

// 将下标为index元素`所在集合`的权值置为value.
func (uf *UnionFindArrayWithUndoAndWeight) SetGroupWeight(index int, value S) {
	index = uf.Find(index)
	uf.history = append(uf.history, historyItem{index, uf.rank[index], uf.weight[index]})
	uf.weight[index] = value
}

// 获取下标为index元素`所在集合`的权值.
func (uf *UnionFindArrayWithUndoAndWeight) GetGroupWeight(index int) S {
	return uf.weight[uf.Find(index)]
}

// 撤销上一次合并(Union)或者修改权值(Set)操作
func (uf *UnionFindArrayWithUndoAndWeight) Undo() {
	if len(uf.history) == 0 {
		return
	}
	last := len(uf.history) - 1
	small := uf.history[last].root
	ps := uf.parent[small]
	uf.weight[ps] = uf.history[last].weight
	uf.rank[ps] = uf.history[last].rank
	if ps != small {
		uf.parent[small] = small
		uf.Part++
	}
	uf.history = uf.history[:last]
}

// 撤销所有操作
func (uf *UnionFindArrayWithUndoAndWeight) Reset() {
	for len(uf.history) > 0 {
		uf.Undo()
	}
}

func (uf *UnionFindArrayWithUndoAndWeight) Find(x int) int {
	if uf.parent[x] == x {
		return x
	}
	return uf.Find(uf.parent[x])
}

func (uf *UnionFindArrayWithUndoAndWeight) Union(x, y int) bool {
	x, y = uf.Find(x), uf.Find(y)
	if uf.rank[x] < uf.rank[y] {
		x ^= y
		y ^= x
		x ^= y
	}
	uf.history = append(uf.history, historyItem{y, uf.rank[x], uf.weight[x]})
	if x != y {
		uf.parent[y] = x
		uf.rank[x] += uf.rank[y]
		uf.weight[x] = op(uf.weight[x], uf.weight[y])
		uf.Part--
		return true
	}
	return false
}

func (ufa *UnionFindArrayWithUndoAndWeight) SetPart(part int) { ufa.Part = part }

func (uf *UnionFindArrayWithUndoAndWeight) IsConnected(x, y int) bool {
	return uf.Find(x) == uf.Find(y)
}

func (uf *UnionFindArrayWithUndoAndWeight) Size(x int) int { return uf.rank[uf.Find(x)] }
