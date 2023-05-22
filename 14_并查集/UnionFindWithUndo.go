// UnionFindWithUndoAndWeight
// https://nyaannyaan.github.io/library/data-structure/rollback-union-find.hpp
// 可撤销并查集(时间旅行)

// API:
// RollbackUnionFind(int sz)：

// Union(int x, int y)：
// Find(int k)：
// IsConnected(int x, int y)：

// Undo()：撤销上一次合并操作，没合并成功也要撤销.

// Snapshot():内部保存当前状态。
//  !Snapshot() 之后可以调用 Rollback(-1) 回滚到这个状态.
// Rollback(int state = -1)：回滚到指定状态。
//   state等于-1时，会回滚到snapshot()中保存的状态。
//   否则，会回滚到指定的state次union调用时的状态。
// GetState()：
//   返回当前状态，即union调用的次数。

package main

import (
	"fmt"
	"sort"
	"strings"
)

func main() {
	uf := NewUnionFindArrayWithUndo(10)
	uf.Union(0, 1)
	uf.Union(2, 3)
	fmt.Println(uf.Find(0), uf.Find(1), uf.Find(2), uf.Find(3))
	uf.Union(0, 2)
	fmt.Println(uf.Find(0), uf.Find(1), uf.Find(2), uf.Find(3))
	uf.Undo()
	fmt.Println(uf.Find(0), uf.Find(1), uf.Find(2), uf.Find(3))
	fmt.Println(uf)
	uf.Snapshot()
	fmt.Println(uf)
	uf.Union(0, 2)
	fmt.Println(uf)
	uf.Union(4, 5)
	fmt.Println(uf)
	uf.Rollback(-1)
	fmt.Println(uf)
	uf.Rollback(0)

	fmt.Println(uf, uf.Part)
}

func NewUnionFindArrayWithUndo(n int) *UnionFindArrayWithUndo {
	data := make([]int, n)
	for i := range data {
		data[i] = -1
	}
	return &UnionFindArrayWithUndo{Part: n, n: n, data: data}
}

type UnionFindArrayWithUndo struct {
	Part      int
	n         int
	innerSnap int
	data      []int
	history   []struct{ a, b int } // (root,data)
}

// !撤销上一次合并操作，没合并成功也要撤销.
func (uf *UnionFindArrayWithUndo) Undo() bool {
	if len(uf.history) == 0 {
		return false
	}
	big, bigData := uf.history[len(uf.history)-1].a, uf.history[len(uf.history)-1].b
	uf.data[big] = bigData
	uf.history = uf.history[:len(uf.history)-1]
	small, smallData := uf.history[len(uf.history)-1].a, uf.history[len(uf.history)-1].b
	uf.data[small] = smallData
	uf.history = uf.history[:len(uf.history)-1]
	if big != small {
		uf.Part++
	}
	return true
}

// 保存并查集当前的状态.
//  !Snapshot() 之后可以调用 Rollback(-1) 回滚到这个状态.
func (uf *UnionFindArrayWithUndo) Snapshot() {
	uf.innerSnap = len(uf.history) >> 1
}

// 回滚到指定的状态.
//  state 为 -1 表示回滚到上一次 `SnapShot` 时保存的状态.
//  其他值表示回滚到状态id为state时的状态.
func (uf *UnionFindArrayWithUndo) Rollback(state int) bool {
	if state == -1 {
		state = uf.innerSnap
	}
	state <<= 1
	if state < 0 || state > len(uf.history) {
		return false
	}
	for state < len(uf.history) {
		uf.Undo()
	}
	return true
}

// 获取当前并查集的状态id.
//  也就是当前合并(Union)被调用的次数.
func (uf *UnionFindArrayWithUndo) GetState() int {
	return len(uf.history) >> 1
}

func (uf *UnionFindArrayWithUndo) Reset() {
	for len(uf.history) > 0 {
		uf.Undo()
	}
}

func (uf *UnionFindArrayWithUndo) Union(x, y int) bool {
	x, y = uf.Find(x), uf.Find(y)
	uf.history = append(uf.history, struct{ a, b int }{x, uf.data[x]})
	uf.history = append(uf.history, struct{ a, b int }{y, uf.data[y]})
	if x == y {
		return false
	}
	if uf.data[x] > uf.data[y] {
		x ^= y
		y ^= x
		x ^= y
	}
	uf.data[x] += uf.data[y]
	uf.data[y] = x
	uf.Part--
	return true
}

func (uf *UnionFindArrayWithUndo) Find(x int) int {
	cur := x
	for uf.data[cur] >= 0 {
		cur = uf.data[cur]
	}
	return cur
}

func (uf *UnionFindArrayWithUndo) IsConnected(x, y int) bool { return uf.Find(x) == uf.Find(y) }

func (uf *UnionFindArrayWithUndo) GetSize(x int) int { return -uf.data[uf.Find(x)] }

func (ufa *UnionFindArrayWithUndo) GetGroups() map[int][]int {
	groups := make(map[int][]int)
	for i := 0; i < ufa.n; i++ {
		root := ufa.Find(i)
		groups[root] = append(groups[root], i)
	}
	return groups
}

func (ufa *UnionFindArrayWithUndo) String() string {
	sb := []string{"UnionFindArray:"}
	groups := ufa.GetGroups()
	keys := make([]int, 0, len(groups))
	for k := range groups {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, root := range keys {
		member := groups[root]
		cur := fmt.Sprintf("%d: %v", root, member)
		sb = append(sb, cur)
	}
	sb = append(sb, fmt.Sprintf("Part: %d", ufa.Part))
	return strings.Join(sb, "\n")
}
