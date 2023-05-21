// UnionFindWithUndoAndWeight
// https://nyaannyaan.github.io/library/data-structure/rollback-union-find.hpp
// 可撤销并查集(时间旅行)

// RollbackUnionFind(int sz)：
//   サイズszのUnionFindを生成する。
// find(int k)：
//   kの根を返す。
// same(int x, int y)：
//   xとyが同じ連結成分に所属しているかを返す。
// snapshot()：
//   現在のUnionFindの状態を保存する。(保存した状態はrollback()で再現できる。)計算量
// get_state()：
//   現在のuniteが呼ばれた回数を返す。
// rollback(int state = -1)：UnionFindをロールバックする。計算量は状況による。(ボトルネックにはならない)
//   state = -1のとき：snapshot()で保存した状態にロールバックする。
//   そうでないとき：uniteがstate回呼び出された時の状態にロールバックする。

package main

import (
	"fmt"
	"sort"
)

func main() {
	uf := NewRollBackUnionFind(10)
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
	fmt.Println(uf)
}

func NewRollBackUnionFind(n int) *RollbackUnionFind {
	data := make([]int, n)
	for i := range data {
		data[i] = -1
	}
	return &RollbackUnionFind{data: data}
}

type RollbackUnionFind struct {
	innerSnap int
	data      []int
	history   []struct{ a, b int }
}

// !撤销上一次合并操作，没合并成功也要撤销.
func (uf *RollbackUnionFind) Undo() bool {
	if len(uf.history) == 0 {
		return false
	}
	uf.data[uf.history[len(uf.history)-1].a] = uf.history[len(uf.history)-1].b
	uf.history = uf.history[:len(uf.history)-1]
	uf.data[uf.history[len(uf.history)-1].a] = uf.history[len(uf.history)-1].b
	uf.history = uf.history[:len(uf.history)-1]
	return true
}

// 回滚到指定的状态.
//  state 为 -1 表示回滚到上一次 `SnapShot` 时保存的状态.
//  其他值表示回滚到合并(Union) `state` 次后的状态.
func (uf *RollbackUnionFind) Rollback(state int) bool {
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

// 获取当前合并(Union)被调用的次数.
func (uf *RollbackUnionFind) GetState() int {
	return len(uf.history) >> 1
}

// 保存并查集当前的状态.
func (uf *RollbackUnionFind) Snapshot() {
	uf.innerSnap = len(uf.history) >> 1
}

func (uf *RollbackUnionFind) Union(x, y int) bool {
	x, y = uf.Find(x), uf.Find(y)
	uf.history = append(uf.history, struct{ a, b int }{x, uf.data[x]})
	uf.history = append(uf.history, struct{ a, b int }{y, uf.data[y]})
	if x == y {
		return false
	}
	if uf.data[x] > uf.data[y] {
		x, y = y, x
	}
	uf.data[x] += uf.data[y]
	uf.data[y] = x
	return true
}

func (uf *RollbackUnionFind) Find(x int) int {
	cur := x
	for uf.data[cur] >= 0 {
		cur = uf.data[cur]
	}
	return cur
}

func (uf *RollbackUnionFind) IsConnected(x, y int) bool { return uf.Find(x) == uf.Find(y) }

func (uf *RollbackUnionFind) GetSize(x int) int { return -uf.data[uf.Find(x)] }

func (uf *RollbackUnionFind) GetGroups() [][]int {
	mp := make(map[int][]int)
	for i := range uf.data {
		mp[uf.Find(i)] = append(mp[uf.Find(i)], i)
	}
	var res [][]int
	for _, g := range mp {
		res = append(res, g)
	}
	return res
}

func (uf *RollbackUnionFind) String() string {
	groups := uf.GetGroups()
	sort.Slice(groups, func(i, j int) bool { return groups[i][0] < groups[j][0] })
	res := []string{}
	for _, g := range groups {
		res = append(res, fmt.Sprintf("%v", g))
	}
	return fmt.Sprintf("state = %d, groups = %v", uf.GetState(), res)
}
