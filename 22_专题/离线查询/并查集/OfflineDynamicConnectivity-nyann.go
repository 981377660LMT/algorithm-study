// https://nyaannyaan.github.io/library/graph/offline-dynamic-connectivity.hpp

package main

import (
	"fmt"
	"sort"
)

func main() {
	odc := NewOffLineDynamicConnectivity(5, 5)
	odc.AddEdge(0, 0, 1)
	odc.AddEdge(1, 1, 2)
	odc.AddEdge(2, 2, 3)
	odc.Build()
	add := func(u, v int) { fmt.Println(fmt.Sprintf("add %d<->%d", u, v)) }
	remove := func(u, v int) { fmt.Println(fmt.Sprintf("remove %d<->%d", u, v)) }
	query := func(t int) { fmt.Println(fmt.Sprintf("query %d", t)) }
	odc.Run(add, remove, query)
}

type pair = struct{ first, second int }
type OffLineDynamicConnectivity struct {
	n, q, segsz     int
	uf              *rollbackUnionFind
	seg, qadd, qdel [][]pair
	cnt             map[pair]*pair // map的value是不可寻址的,需要使用指针类型代替
}

// Init with n nodes and q queries.
func NewOffLineDynamicConnectivity(n, q int) *OffLineDynamicConnectivity {
	res := &OffLineDynamicConnectivity{
		n:     n,
		q:     q,
		segsz: 1,
		uf:    newRollBackUnionFind(n),
		qadd:  make([][]pair, q),
		qdel:  make([][]pair, q),
		cnt:   make(map[pair]*pair),
	}
	for res.segsz < q {
		res.segsz *= 2
	}
	res.seg = make([][]pair, res.segsz*2)
	return res
}

func (odc *OffLineDynamicConnectivity) Build() {
	for i := 0; i < odc.q; i++ {
		for _, e := range odc.qadd[i] {
			if _, ok := odc.cnt[e]; !ok {
				odc.cnt[e] = &pair{0, 0}
			}
			dat := odc.cnt[e]
			if dat.second++; dat.second == 1 {
				dat.first = i
			}
		}

		for _, e := range odc.qdel[i] {
			if _, ok := odc.cnt[e]; !ok {
				odc.cnt[e] = &pair{0, 0}
			}
			dat := odc.cnt[e]
			if dat.second--; dat.second == 0 {
				odc.segment(e, dat.first, i)
			}
		}
	}

	for e, dat := range odc.cnt {
		if dat.second != 0 {
			odc.segment(e, dat.first, odc.q)
		}
	}
}

// Add an edge u<->v at time t.
func (odc *OffLineDynamicConnectivity) AddEdge(t, u, v int) {
	if u > v {
		u, v = v, u
	}
	odc.qadd[t] = append(odc.qadd[t], pair{u, v})
}

// Remove an edge u<->v at time t.
func (odc *OffLineDynamicConnectivity) RemoveEdge(t, u, v int) {
	if u > v {
		u, v = v, u
	}
	odc.qdel[t] = append(odc.qdel[t], pair{u, v})
}

func (odc *OffLineDynamicConnectivity) Run(add, remove func(u, v int), query func(qi int)) {
	odc.dfs(add, remove, query, 1, 0, odc.segsz)
}

func (odc *OffLineDynamicConnectivity) dfs(add, remove func(u, v int), query func(qi int), id, l, r int) {
	if odc.q <= l {
		return
	}
	state := odc.uf.GetState()
	var es []pair
	for _, e := range odc.seg[id] {
		u, v := e.first, e.second
		if !odc.uf.IsConnected(u, v) {
			odc.uf.Union(u, v)
			add(u, v)
			es = append(es, e)
		}
	}

	if l+1 == r {
		query(l)
	} else {
		odc.dfs(add, remove, query, id<<1, l, (l+r)>>1)
		odc.dfs(add, remove, query, id<<1|1, (l+r)>>1, r)
	}

	for _, e := range es {
		remove(e.first, e.second)
	}
	odc.uf.Rollback(state)
}

func (odc *OffLineDynamicConnectivity) segment(e pair, l, r int) {
	left, right := l+odc.segsz, r+odc.segsz
	for left < right {
		if left&1 == 1 {
			odc.seg[left] = append(odc.seg[left], e)
			left++
		}
		if right&1 == 1 {
			right--
			odc.seg[right] = append(odc.seg[right], e)
		}
		left >>= 1
		right >>= 1
	}
}

func newRollBackUnionFind(n int) *rollbackUnionFind {
	data := make([]int, n)
	for i := range data {
		data[i] = -1
	}
	return &rollbackUnionFind{data: data}
}

type rollbackUnionFind struct {
	innerSnap int
	data      []int
	history   []struct{ a, b int }
}

// 撤销上一次合并操作.
func (uf *rollbackUnionFind) Undo() bool {
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
func (uf *rollbackUnionFind) Rollback(state int) bool {
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
func (uf *rollbackUnionFind) GetState() int {
	return len(uf.history) >> 1
}

// 保存并查集当前的状态.
func (uf *rollbackUnionFind) Snapshot() {
	uf.innerSnap = len(uf.history) >> 1
}

func (uf *rollbackUnionFind) Union(x, y int) bool {
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

func (uf *rollbackUnionFind) Find(x int) int {
	cur := x
	for uf.data[cur] >= 0 {
		cur = uf.data[cur]
	}
	return cur
}

func (uf *rollbackUnionFind) IsConnected(x, y int) bool { return uf.Find(x) == uf.Find(y) }

func (uf *rollbackUnionFind) GetSize(x int) int { return -uf.data[uf.Find(x)] }

func (uf *rollbackUnionFind) GetGroups() [][]int {
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

func (uf *rollbackUnionFind) String() string {
	groups := uf.GetGroups()
	sort.Slice(groups, func(i, j int) bool { return groups[i][0] < groups[j][0] })
	res := []string{}
	for _, g := range groups {
		res = append(res, fmt.Sprintf("%v", g))
	}
	return fmt.Sprintf("state = %d, groups = %v", uf.GetState(), res)
}
