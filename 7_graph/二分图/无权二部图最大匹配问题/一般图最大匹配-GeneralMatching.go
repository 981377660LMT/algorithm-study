// TODO
// 定义和yosupo模板

//
// General Graph Matching (Gabow-Edmonds)
//
// Description:
//
//   For a graph G = (V, E), a matching M is a set of edges
//   such that any vertex is contained in M at most once.
//   The matching with maximum cardinality is computed by
//   the Edmonds blossom algorithm.
//
//   This implementation is the Gabow's simplified version
//   with the lazy update technique to improve the complexity
//   in sparse graphs.
//
//
// Complexity:
//
//   O(n m log n)
//
//
// Verified:
//
//   SPOJ ADABLOOM
//
//
// References:
//   H.Gabow (1976):
//   An efficient implementation of Edmonds' algorithm for maximum matching on graphs.
//   Journal of the ACM, vol.23, no.2, pp.221-234.
//
//   http://min-25.hatenablog.com/entry/2016/11/21/222625
//   https://ei1333.github.io/library/graph/flow/gabow-edmonds.hpp
//
//   !一般图的最大匹配

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)
	flow := NewGeneralMatching(n)
	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		flow.AddEdge(u, v)
	}
	ret := flow.MaxMathing()
	fmt.Fprintln(out, len(ret))
	for _, p := range ret {
		fmt.Fprintln(out, p[0], p[1])
	}
}

// GabowEdmonds 算法求一般图的最大匹配.
//  O(VElogV).
type GeneralMatching struct {
	g                  [][][2]int // (to,idx)
	edges              [][2]int
	mate, label, first []int
	que                []int
}

func NewGeneralMatching(n int) *GeneralMatching {
	gm := &GeneralMatching{}
	gm.g = make([][][2]int, n+1)
	gm.mate = make([]int, n+1)
	gm.label = make([]int, n+1)
	for i := range gm.label {
		gm.label[i] = -1
	}
	gm.first = make([]int, n+1)
	return gm
}

// 连接无向边 u-v.
// 0 <= u,v < n.
func (gm *GeneralMatching) AddEdge(u, v int) {
	u++
	v++
	gm.g[u] = append(gm.g[u], [2]int{v, len(gm.edges) + len(gm.g)})
	gm.g[v] = append(gm.g[v], [2]int{u, len(gm.edges) + len(gm.g)})
	gm.edges = append(gm.edges, [2]int{u, v})
}

// 返回最大匹配的边集合.
func (gm *GeneralMatching) MaxMathing() [][2]int {
	for i := 1; i < len(gm.g); i++ {
		if gm.mate[i] != 0 {
			continue
		}
		if gm._augmentCheck(i) {
			for j := range gm.label {
				gm.label[j] = -1
			}
		}
	}

	var res [][2]int
	for i := 1; i < len(gm.g); i++ {
		if i < gm.mate[i] {
			res = append(res, [2]int{i - 1, gm.mate[i] - 1})
		}
	}
	return res
}

func (gm *GeneralMatching) _find(x int) int {
	if gm.label[gm.first[x]] < 0 {
		return gm.first[x]
	}
	gm.first[x] = gm._find(gm.first[x])
	return gm.first[x]
}

func (gm *GeneralMatching) _reMatch(v, w int) {
	t := gm.mate[v]
	gm.mate[v] = w
	if gm.mate[t] != v {
		return
	}
	if gm.label[v] < len(gm.g) {
		gm.mate[t] = gm.label[v]
		gm._reMatch(gm.label[v], t)
	} else {
		x := gm.edges[gm.label[v]-len(gm.g)][0]
		y := gm.edges[gm.label[v]-len(gm.g)][1]
		gm._reMatch(x, y)
		gm._reMatch(y, x)
	}
}

func (gm *GeneralMatching) _assignLabel(x, y, num int) {
	r := gm._find(x)
	s := gm._find(y)
	join := 0
	if r == s {
		return
	}
	gm.label[r] = -num
	gm.label[s] = -num
	for {
		if s != 0 {
			r, s = s, r
		}
		r = gm._find(gm.label[gm.mate[r]])
		if gm.label[r] == -num {
			join = r
			break
		}
		gm.label[r] = -num
	}
	v := gm.first[x]
	for v != join {
		gm.que = append(gm.que, v)
		gm.label[v] = num
		gm.first[v] = join
		v = gm.first[gm.label[gm.mate[v]]]
	}
	v = gm.first[y]
	for v != join {
		gm.que = append(gm.que, v)
		gm.label[v] = num
		gm.first[v] = join
		v = gm.first[gm.label[gm.mate[v]]]
	}
}

func (gm *GeneralMatching) _augmentCheck(u int) bool {
	gm.que = gm.que[:0]
	gm.first[u] = 0
	gm.label[u] = 0
	gm.que = append(gm.que, u)
	for len(gm.que) > 0 {
		x := gm.que[0]
		gm.que = gm.que[1:]
		for _, e := range gm.g[x] {
			y := e[0]
			if gm.mate[y] == 0 && y != u {
				gm.mate[y] = x
				gm._reMatch(x, y)
				return true
			} else if gm.label[y] >= 0 {
				gm._assignLabel(x, y, e[1])
			} else if gm.label[gm.mate[y]] < 0 {
				gm.label[gm.mate[y]] = x
				gm.first[gm.mate[y]] = y
				gm.que = append(gm.que, gm.mate[y])
			}
		}
	}
	return false
}
