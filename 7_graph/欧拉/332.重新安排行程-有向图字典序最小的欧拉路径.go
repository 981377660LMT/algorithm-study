package main

import (
	"sort"
)

func findItinerary(tickets [][]string) []string {
	// !1. 将字符串转换为数字并按照字典序分配id
	set := make(map[string]struct{})
	for _, e := range tickets {
		set[e[0]] = struct{}{}
		set[e[1]] = struct{}{}
	}
	allVertex := make([]string, 0, len(set))
	for k := range set {
		allVertex = append(allVertex, k)
	}
	sort.Strings(allVertex)
	mp := make(map[string]int)  // raw -> cur
	rmp := make(map[int]string) // cur -> raw
	for i, v := range allVertex {
		mp[v] = i
		rmp[i] = v
	}

	// 2. 求出字典序最小的欧拉路径
	et := NewEulerianTrail(len(mp), true) // 构建有向图
	for _, e := range tickets {
		et.AddEdge(mp[e[0]], mp[e[1]])
	}
	eids := et.GetSemiEulerianTrailStartsWith(mp["JFK"], true)
	if len(eids) == 0 {
		return nil
	}

	path := et.GetPathFromEdgeIds(eids)
	res := make([]string, 0, len(path))
	for _, v := range path {
		res = append(res, rmp[v])
	}
	return res
}

type EulerianTrail struct {
	g          [][][2]int
	es         [][2]int
	m          int
	usedVertex []bool
	usedEdge   []bool
	deg        []int
	directed   bool
}

func NewEulerianTrail(n int, directed bool) *EulerianTrail {
	res := &EulerianTrail{
		g:          make([][][2]int, n),
		usedVertex: make([]bool, n),
		deg:        make([]int, n),
		directed:   directed,
	}
	return res
}

func (e *EulerianTrail) AddEdge(a, b int) {
	e.es = append(e.es, [2]int{a, b})
	e.g[a] = append(e.g[a], [2]int{b, e.m})
	if e.directed {
		e.deg[a]++
		e.deg[b]--
	} else {
		e.g[b] = append(e.g[b], [2]int{a, e.m})
		e.deg[a]++
		e.deg[b]++
	}
	e.m++
}

func (e *EulerianTrail) GetPathFromEdgeIds(edgeIds []int) []int {
	path := make([]int, 0, len(edgeIds)+1)
	for i, id := range edgeIds {
		a, b := e.es[id][0], e.es[id][1]
		path = append(path, a)
		if i == len(edgeIds)-1 {
			path = append(path, b)
		}
	}
	return path
}

// 枚举所有连通块的`欧拉回路`,返回边的编号.
//  如果连通块内不存在欧拉回路,返回空.
//  minLex: 字典序最小.
func (e *EulerianTrail) EnumerateEulerianTrail(minLex bool) [][]int {
	if e.directed {
		for _, d := range e.deg {
			if d != 0 {
				return [][]int{}
			}
		}
	} else {
		for _, d := range e.deg {
			if d&1 == 1 {
				return [][]int{}
			}
		}
	}

	e.sortNeighbors(minLex)
	e.usedEdge = make([]bool, e.m)
	res := [][]int{}
	for i := 0; i < len(e.g); i++ {
		if !e.usedVertex[i] && len(e.g[i]) > 0 {
			res = append(res, e.work(i))
		}
	}
	return res
}

// 获取整张图的从任意点出发的`欧拉回路`,返回边的编号.
//  如果不存在欧拉回路,返回空.
//  minLex: 字典序最小.
func (e *EulerianTrail) GetEulerianTrail(minLex bool) (eids []int) {
	groups := e.EnumerateEulerianTrail(minLex)
	if len(groups) != 1 || len(groups[0]) != len(e.es) {
		return
	}
	eids = groups[0]
	return
}

// 获取整张图的从`start`出发的`欧拉回路`,返回边的编号.
//  如果从`start`出发不存在欧拉回路,返回空.
//  minLex: 字典序最小.
func (e *EulerianTrail) GetEulerianTrailStartsWith(start int, minLex bool) (eids []int) {
	if e.directed {
		for _, d := range e.deg {
			if d != 0 {
				return
			}
		}
	} else {
		for _, d := range e.deg {
			if d&1 == 1 {
				return
			}
		}
	}

	e.sortNeighbors(minLex)
	e.usedEdge = make([]bool, e.m)
	res := e.work(start)
	if len(res) != len(e.es) {
		return
	}
	eids = res
	return
}

// 枚举所有连通块的`欧拉路径`(半欧拉回路),返回边的编号.
//  如果连通块内不存在欧拉路径,返回空.
//  minLex: 字典序最小.
func (e *EulerianTrail) EnumerateSemiEulerianTrail(minLex bool) [][]int {
	e.sortNeighbors(minLex)

	uf := newUnionFindArray(len(e.g))
	for _, es := range e.es {
		uf.Union(es[0], es[1])
	}
	group := make([][]int, len(e.g))
	for i := 0; i < len(e.g); i++ {
		group[uf.Find(i)] = append(group[uf.Find(i)], i)
	}

	res := [][]int{}
	e.usedEdge = make([]bool, e.m)
	for _, vs := range group {
		if len(vs) == 0 {
			continue
		}

		latte, malta := -1, -1
		if e.directed {
			for _, p := range vs {
				if abs(e.deg[p]) > 1 {
					return [][]int{}
				} else if e.deg[p] == 1 {
					if latte >= 0 {
						return [][]int{}
					}
					latte = p
				}
			}
		} else {
			for _, p := range vs {
				if e.deg[p]&1 == 1 {
					if latte == -1 {
						latte = p
					} else if malta == -1 {
						malta = p
					} else {
						return [][]int{}
					}
				}
			}
		}

		var cur []int
		if latte == -1 {
			cur = e.work(vs[0]) // 起点任意
		} else {
			cur = e.work(latte) // 起点选latte(有向图必须是latte,无向图可以是latte或malta)
		}

		if len(cur) > 0 {
			res = append(res, cur)
		}
	}

	return res
}

// 获取整张图的从任意点出发的`欧拉路径`,返回边的编号.
//  如果不存在欧拉路径,返回空.
//  minLex: 字典序最小.
func (e *EulerianTrail) GetSemiEulerianTrail(minLex bool) (eids []int) {
	groups := e.EnumerateSemiEulerianTrail(minLex)
	if len(groups) == 0 || len(groups[0]) != len(e.es) {
		return
	}
	eids = groups[0]
	return
}

// 获取整张图的从`start`出发的`欧拉路径`,返回边的编号.
//  如果从`start`出发不存在欧拉路径,返回空.
//  minLex: 字典序最小.
func (e *EulerianTrail) GetSemiEulerianTrailStartsWith(start int, minLex bool) (eids []int) {
	e.sortNeighbors(minLex)

	e.usedEdge = make([]bool, e.m)

	latte, malta := -1, -1
	if e.directed {
		for i := 0; i < len(e.g); i++ {
			if abs(e.deg[i]) > 1 {
				return
			} else if e.deg[i] == 1 {
				if latte >= 0 {
					return
				}
				latte = i
			}
		}
	} else {
		for i := 0; i < len(e.g); i++ {
			if e.deg[i]&1 == 1 {
				if latte == -1 {
					latte = i
				} else if malta == -1 {
					malta = i
				} else {
					return
				}
			}
		}
	}

	if e.directed {
		if latte != -1 && latte != start {
			return
		}
	} else {
		if latte != -1 && (latte != start && malta != start) {
			return
		}
	}

	res := e.work(start)
	if len(res) != len(e.es) {
		return
	}
	eids = res
	return
}

func (e *EulerianTrail) GetEdge(index int) (int, int) {
	return e.es[index][0], e.es[index][1]
}

func (e *EulerianTrail) work(s int) []int {
	st := [][2]int{}
	ord := []int{}
	st = append(st, [2]int{s, -1})
	for len(st) > 0 {
		index := st[len(st)-1][0]
		e.usedVertex[index] = true
		if len(e.g[index]) == 0 {
			ord = append(ord, st[len(st)-1][1])
			st = st[:len(st)-1]
		} else {
			e_ := e.g[index][len(e.g[index])-1]
			e.g[index] = e.g[index][:len(e.g[index])-1]
			if e.usedEdge[e_[1]] {
				continue
			}
			e.usedEdge[e_[1]] = true
			st = append(st, [2]int{e_[0], e_[1]})
		}
	}

	ord = ord[:len(ord)-1]
	for i, j := 0, len(ord)-1; i < j; i, j = i+1, j-1 {
		ord[i], ord[j] = ord[j], ord[i]
	}
	return ord
}

// 排在邻接表后面的点先出来.
func (e *EulerianTrail) sortNeighbors(minLex bool) {
	if minLex {
		for _, es := range e.g {
			sort.Slice(es, func(i, j int) bool {
				return es[i][0] > es[j][0]
			})
		}
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func newUnionFindArray(n int) *unionFindArray {
	parent, rank := make([]int, n), make([]int, n)
	for i := 0; i < n; i++ {
		parent[i] = i
		rank[i] = 1
	}

	return &unionFindArray{
		Part:   n,
		size:   n,
		Rank:   rank,
		parent: parent,
	}
}

type unionFindArray struct {
	size   int
	Part   int
	Rank   []int
	parent []int
}

func (ufa *unionFindArray) Union(key1, key2 int) bool {
	root1, root2 := ufa.Find(key1), ufa.Find(key2)
	if root1 == root2 {
		return false
	}
	if ufa.Rank[root1] > ufa.Rank[root2] {
		root1, root2 = root2, root1
	}
	ufa.parent[root1] = root2
	ufa.Rank[root2] += ufa.Rank[root1]
	ufa.Part--
	return true
}

func (ufa *unionFindArray) Find(key int) int {
	for ufa.parent[key] != key {
		ufa.parent[key] = ufa.parent[ufa.parent[key]]
		key = ufa.parent[key]
	}
	return key
}

func (ufa *unionFindArray) IsConnected(key1, key2 int) bool {
	return ufa.Find(key1) == ufa.Find(key2)
}

func (ufa *unionFindArray) Size(key int) int {
	return ufa.Rank[ufa.Find(key)]
}
