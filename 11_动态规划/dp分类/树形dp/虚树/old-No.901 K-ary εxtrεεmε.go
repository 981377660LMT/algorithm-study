// https://yukicoder.me/problems/3407
// !给定q个查询,求虚树(最小的包含指定点集的连通子图)组成的的边权之和
// !因为要从根结点出发的链上求路径边权之和,可以用前缀和(差分)来求

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	edges := make([][]int, 0, n-1)
	weight := make([]map[int]int, n)
	for i := 0; i < n; i++ {
		weight[i] = make(map[int]int)
	}
	for i := 0; i < n-1; i++ {
		var u, v, w int
		fmt.Fscan(in, &u, &v, &w)
		edges = append(edges, []int{u, v})
		weight[u][v] = w
		weight[v][u] = w
	}

	A := NewAuxiliaryTree(n, edges)
	adjList := A.RawTree
	dist := make([]int, n) // !深度(带权)
	for i := 0; i < n; i++ {
		dist[i] = -1
	}
	queue := []int{0}
	dist[0] = 0
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		for _, next := range adjList[cur] {
			if dist[next] != -1 {
				continue
			}
			dist[next] = dist[cur] + weight[cur][next]
			queue = append(queue, next)
		}
	}

	var q int
	fmt.Fscan(in, &q)
	for i := 0; i < q; i++ {
		var k int
		fmt.Fscan(in, &k)
		points := make([]int, k)
		for j := 0; j < k; j++ {
			fmt.Fscan(in, &points[j])
		}

		tree, root := A.Query(points)
		if root == -1 {
			fmt.Fprintln(out, 0)
			continue
		}

		res := 0
		var dfs func(cur, pre int)
		dfs = func(cur, pre int) {
			for _, next := range tree[cur] {
				if next == pre {
					continue
				}
				res += dist[next] - dist[cur]
				dfs(next, cur)
			}
		}
		dfs(root, root)

		fmt.Fprintln(out, res)
	}
}

type AuxiliaryTree struct {
	RawTree       [][]int // 原图邻接表(无向边)
	g0            [][]int // 虚树邻接表(有向边)
	s             []int
	fs, ls, depth []int

	lg []int
	st [][]int
}

// 给定顶点个数n和无向边集(u,v)构建.
//  O(nlogn)
func NewAuxiliaryTree(n int, edges [][]int) *AuxiliaryTree {
	g := make([][]int, n)
	for _, e := range edges {
		g[e[0]] = append(g[e[0]], e[1])
		g[e[1]] = append(g[e[1]], e[0])
	}
	res := &AuxiliaryTree{
		RawTree: g,
		g0:      make([][]int, n),
		s:       []int{},
		fs:      make([]int, n),
		ls:      make([]int, n),
		depth:   make([]int, n),
	}

	res.dfs(0, -1, 0)
	res.buildSt()
	return res
}

// 指定点集,返回虚树的(有向图邻接表,虚树的根).
//  如果虚树不存在`(len(points<=1))`,返回空邻接表和-1.
//  O(klogk) 构建虚树.
func (t *AuxiliaryTree) Query(points []int) ([][]int, int) {
	k := len(points)
	points = append(points[:0:0], points...)
	sort.Slice(points, func(i, j int) bool {
		return t.fs[points[i]] < t.fs[points[j]]
	})

	for i := 0; i < k-1; i++ {
		x, y := t.fs[points[i]], t.fs[points[i+1]]
		l := t.lg[y-x+1]
		w := t.st[l][x]
		if t.depth[t.st[l][y-(1<<l)+1]] < t.depth[t.st[l][x]] {
			w = t.st[l][y-(1<<l)+1]
		}
		points = append(points, w)
	}

	sort.Slice(points, func(i, j int) bool {
		return t.fs[points[i]] < t.fs[points[j]]
	})

	stk := []int{}
	pre := -1
	root := -1
	for _, v := range points {
		if pre == v {
			continue
		}
		for len(stk) > 0 && t.ls[stk[len(stk)-1]] < t.fs[v] {
			stk = stk[:len(stk)-1]
		}
		if len(stk) > 0 {
			parent := stk[len(stk)-1]
			t.g0[parent] = append(t.g0[parent], v)
			if root == -1 {
				root = parent
			}
		}

		t.g0[v] = t.g0[v][:0]
		stk = append(stk, v)
		pre = v
	}

	return t.g0, root
}

func (t *AuxiliaryTree) dfs(v, p, d int) {
	t.depth[v] = d
	t.fs[v] = len(t.s)
	t.s = append(t.s, v)
	for _, w := range t.RawTree[v] {
		if w == p {
			continue
		}
		t.dfs(w, v, d+1)
		t.s = append(t.s, v)
	}
	t.ls[v] = len(t.s)
	t.s = append(t.s, v)
}

func (t *AuxiliaryTree) buildSt() {
	l := len(t.s)
	lg := make([]int, l+1)
	for i := 2; i <= l; i++ {
		lg[i] = lg[i>>1] + 1
	}
	st := make([][]int, lg[l]+1)
	for i := range st {
		st[i] = make([]int, l-(1<<i)+1)
		for j := range st[i] {
			st[i][j] = l
		}
	}

	copy(st[0], t.s)
	b := 1
	for i := 0; i < lg[l]; i++ {
		st0, st1 := st[i], st[i+1]
		for j := 0; j < l-(b<<1)+1; j++ {
			st1[j] = st0[j]
			if t.depth[st0[j+b]] < t.depth[st0[j]] {
				st1[j] = st0[j+b]
			}
		}
		b <<= 1
	}

	t.lg = lg
	t.st = st
}
