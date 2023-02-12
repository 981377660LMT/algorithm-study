package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// https://yukicoder.me/problems/no/583
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)
	et := NewEulerianTrail(n, false)
	for i := 0; i < m; i++ {
		var a, b int
		fmt.Fscan(in, &a, &b)
		et.AddEdge(a, b)
	}
	res := et.EnumerateSemiEulerianTrail()
	if len(res) == 1 {
		fmt.Fprintln(out, "YES")
	} else {
		fmt.Fprintln(out, "NO")
	}
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

// 枚举所有连通块的`欧拉回路`,返回边的编号.
//  如果连通块内不存在欧拉回路,返回空.
func (e *EulerianTrail) EnumerateEulerianTrail() [][]int {
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

	e.usedEdge = make([]bool, e.m)
	res := [][]int{}
	for i := 0; i < len(e.g); i++ {
		if !e.usedVertex[i] && len(e.g[i]) > 0 {
			res = append(res, e.work(i))
		}
	}

	return res
}

// 枚举所有连通块的`欧拉路径`(半欧拉回路),返回边的编号.
//  如果连通块内不存在欧拉路径,返回空.
func (e *EulerianTrail) EnumerateSemiEulerianTrail() [][]int {
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
			cur = e.work(vs[0])
		} else {
			cur = e.work(latte)
		}

		if len(cur) > 0 {
			res = append(res, cur)
		}
	}

	return res
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
