// https://ei1333.github.io/library/other/mo-tree.hpp
// https://oi-wiki.org/misc/mo-algo-on-tree/
// https://github.com/EndlessCheng/codeforces-go/blob/53262fb81ffea176cd5f039cec71e3bd266dce83/copypasta/mo.go#L301
// https://github.com/EndlessCheng/codeforces-go/blob/53262fb81ffea176cd5f039cec71e3bd266dce83/copypasta/mo.go#L301
//
// 处理树上的路径相关的离线查询.
// 一般的莫队只能处理线性问题，我们要把树强行压成序列。
// 通过欧拉序(括号序)转化成序列上的查询，然后用莫队解决。

package main

import (
	"math"
	"sort"
)

// 维护点权的树上莫队.
type MoOnTree struct {
	tree    [][]int
	root    int
	queries [][2]int
}

func NewMoOnTree(tree [][]int, root int) *MoOnTree {
	return &MoOnTree{tree: tree, root: root}
}

// 添加从顶点u到顶点v的查询.
func (mo *MoOnTree) AddQuery(u, v int) { mo.queries = append(mo.queries, [2]int{u, v}) }

// 处理每个查询.
//
//	add: 将数据添加到窗口.
//	remove: 将数据从窗口移除.
//	query: 查询窗口内的数据.
func (mo *MoOnTree) Run(add func(rootId int), remove func(rootId int), query func(qid int)) {
	if len(mo.queries) == 0 {
		return
	}

	n := len(mo.tree)
	dfnToNode := make([]int, 0, 2*n)
	ins := make([]int, n)
	outs := make([]int, n)

	var dfs func(cur, pre int)
	dfs = func(cur, pre int) {
		ins[cur] = len(dfnToNode)
		dfnToNode = append(dfnToNode, cur)
		for _, to := range mo.tree[cur] {
			if to != pre {
				dfs(to, cur)
			}
		}
		outs[cur] = len(dfnToNode)
		dfnToNode = append(dfnToNode, cur)
	}
	dfs(mo.root, -1)

	lca := _offlineLCA(mo.tree, mo.queries, mo.root)
	blockSize := int(math.Ceil(float64(2*n) / math.Sqrt(float64(len(mo.queries)))))
	type Q struct{ bid, l, r, lca, qid int }
	qs := make([]Q, len(mo.queries))
	for i := range qs {
		v, w := mo.queries[i][0], mo.queries[i][1]
		if ins[v] > ins[w] {
			v, w = w, v
		}
		if lca_ := lca[i]; lca_ != v {
			qs[i] = Q{outs[v] / blockSize, outs[v], ins[w] + 1, lca_, i}
		} else {
			qs[i] = Q{ins[v] / blockSize, ins[v], ins[w] + 1, -1, i}
		}
	}

	sort.Slice(qs, func(i, j int) bool {
		a, b := qs[i], qs[j]
		if a.bid != b.bid {
			return a.bid < b.bid
		}
		if a.bid&1 == 0 {
			return a.r < b.r
		}
		return a.r > b.r
	})

	flip := make([]bool, n)
	f := func(u int) {
		flip[u] = !flip[u]
		if flip[u] {
			add(u)
		} else {
			remove(u)
		}
	}

	l, r := 0, 0
	for _, q := range qs {
		for ; r < q.r; r++ {
			f(dfnToNode[r])
		}
		for ; l < q.l; l++ {
			f(dfnToNode[l])
		}
		for l > q.l {
			l--
			f(dfnToNode[l])
		}
		for r > q.r {
			r--
			f(dfnToNode[r])
		}
		if q.lca >= 0 {
			f(q.lca)
		}
		query(q.qid)
		if q.lca >= 0 {
			f(q.lca)
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// LCA离线.
func _offlineLCA(tree [][]int, queries [][2]int, root int) []int {
	n := len(tree)
	ufa := NewUnionFindArray(n)
	st, mark, ptr, res := make([]int, n), make([]int, n), make([]int, n), make([]int, len(queries))
	for i := 0; i < len(queries); i++ {
		res[i] = -1
	}
	top := 0
	st[top] = root
	for _, q := range queries {
		mark[q[0]]++
		mark[q[1]]++
	}
	q := make([][][2]int, n)
	for i := 0; i < n; i++ {
		q[i] = make([][2]int, 0, mark[i])
		mark[i] = -1
		ptr[i] = len(tree[i])
	}
	for i := range queries {
		u, v := queries[i][0], queries[i][1]
		q[u] = append(q[u], [2]int{v, i})
		q[v] = append(q[v], [2]int{u, i})
	}
	run := func(u int) bool {
		for ptr[u] != 0 {
			v := tree[u][ptr[u]-1]
			ptr[u]--
			if mark[v] == -1 {
				top++
				st[top] = v
				return true
			}
		}
		return false
	}

	for top != -1 {
		u := st[top]
		if mark[u] == -1 {
			mark[u] = u
		} else {
			ufa.Union(u, tree[u][ptr[u]])
			mark[ufa.Find(u)] = u
		}

		if !run(u) {
			for _, v := range q[u] {
				if mark[v[0]] != -1 && res[v[1]] == -1 {
					res[v[1]] = mark[ufa.Find(v[0])]
				}
			}
			top--
		}
	}

	return res
}

type _unionFindArray struct {
	data []int
}

func NewUnionFindArray(n int) *_unionFindArray {
	data := make([]int, n)
	for i := 0; i < n; i++ {
		data[i] = -1
	}
	return &_unionFindArray{data: data}
}

func (ufa *_unionFindArray) Union(key1, key2 int) bool {
	root1, root2 := ufa.Find(key1), ufa.Find(key2)
	if root1 == root2 {
		return false
	}
	if ufa.data[root1] > ufa.data[root2] {
		root1 ^= root2
		root2 ^= root1
		root1 ^= root2
	}
	ufa.data[root1] += ufa.data[root2]
	ufa.data[root2] = root1
	return true
}

func (ufa *_unionFindArray) Find(key int) int {
	if ufa.data[key] < 0 {
		return key
	}
	ufa.data[key] = ufa.Find(ufa.data[key])
	return ufa.data[key]
}
