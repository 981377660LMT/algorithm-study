// Dominator Tree of a Directed Graph
// O(ElogV)

// python版本:
// https://tjkendev.github.io/procon-library/python/graph/dominator-lengauer-tarjan.html

package main

import (
	"bufio"
	"fmt"
	"os"
)

//
//
// 求出有向图的支配树, 返回每个节点的`直接支配顶点idom`.
//  n: 顶点数
//  edges: 有向边
//  root: 支配树的根节点
// 考虑支配树中从root到达结点i的路径, 该路径上的结点都是i的支配者，要到达i就必须经过这些结点(有向图的割点).
// 其中, 从root到达结点i的路径上的最后一个结点就是i的直接支配者idom[i]，也是距离i`最近的支配者`.
// !idom[i]=-1：从根节点无法到达结点i
// !idom[i]=root：从根节点到达结点i的路径上不存在i的支配者(割点)
// !idom[i]为其他值：i的最近的支配者(割点)为idom[i]
func BuildDominatorTree(n int, edges [][]int, root int) []int {
	graph := make([][]int, n)
	rGraph := make([][]int, n)
	for _, e := range edges { // addEdge
		u, v := e[0], e[1]
		graph[u] = append(graph[u], v)
		rGraph[v] = append(rGraph[v], u)
	}

	// build
	par := make([]int, n)
	idom, semi := make([]int, n), make([]int, n) // idom: immediate dominator
	for i := 0; i < n; i++ {
		idom[i], semi[i] = -1, -1
	}
	ord := make([]int, 0, n)

	var dfs func(v int)
	dfs = func(v int) {
		semi[v] = len(ord)
		ord = append(ord, v)
		for _, u := range graph[v] {
			if semi[u] != -1 {
				continue
			}
			par[u] = v
			dfs(u)
		}
	}

	dfs(root)
	bkt := make([][]int, n)
	uf := newUnionFind(semi)
	us := make([]int, n)

	for i := len(ord) - 1; i >= 0; i-- {
		v := ord[i]
		for _, u := range rGraph[v] {
			if semi[u] < 0 {
				continue
			}
			u := uf.eval(u)
			if semi[u] < semi[v] {
				semi[v] = semi[u]
			}
		}
		bkt[ord[semi[v]]] = append(bkt[ord[semi[v]]], v)
		for _, u := range bkt[par[v]] {
			us[u] = uf.eval(u)
		}
		bkt[par[v]] = bkt[par[v]][:0] // clear
		uf.link(par[v], v)
	}

	for i := 1; i < len(ord); i++ {
		v := ord[i]
		u := us[v]
		if semi[v] == semi[u] {
			idom[v] = semi[v]
		} else {
			idom[v] = idom[u]
		}
	}

	for i := 1; i < len(ord); i++ {
		v := ord[i]
		idom[v] = ord[idom[v]]
	}

	idom[root] = root

	return idom
}

type unionFind struct {
	semi   []int
	ps, ms []int
}

func newUnionFind(semi []int) *unionFind {
	n := len(semi)
	ps, ms := make([]int, n), make([]int, n)
	for i := 0; i < n; i++ {
		ps[i], ms[i] = i, i
	}
	return &unionFind{semi, ps, ms}
}

func (uf *unionFind) find(v int) int {
	if uf.ps[v] == v {
		return v
	}
	r := uf.find(uf.ps[v])
	if uf.semi[uf.ms[v]] > uf.semi[uf.ms[uf.ps[v]]] {
		uf.ms[v] = uf.ms[uf.ps[v]]
	}
	uf.ps[v] = r
	return r
}

func (uf *unionFind) eval(v int) int {
	uf.find(v)
	return uf.ms[v]
}

func (uf *unionFind) link(p, v int) {
	uf.ps[v] = p
}

//
//
// 捕获赤牛大作战
// https://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=0294
// 给定一个带有n个结点的`连通的有向图`,基地位于0号点
// 给定q个询问
// 每个询问给定一个赤牛的位置qi,你需要选择一个点作为埋伏点来捕获赤牛
// 选择的点的要求为:
// 1. 不能在基地埋伏(太危险)
// 2. 埋伏点必须是从基地到达赤牛位置的所有路径中一定会经过的点
// 3. 在满足1,2的条件下,埋伏点距离赤牛的距离越近越好;如果不存在这样的点,就直接埋伏在赤牛的位置
// 对每个询问,输出一个整数,表示你选择的埋伏点的编号
// n<=1e5, m<=3e5, q<=1e5
func catchAkabeko(n int, edges [][]int, queries []int) []int {
	idom := BuildDominatorTree(n, edges, 0) // 每个点最近的支配者
	res := make([]int, 0, len(queries))
	for _, q := range queries {
		if idom[q] == 0 || idom[q] == -1 { // 不存在割点或者无法到达,直接埋伏在赤牛的位置
			res = append(res, q)
		} else {
			res = append(res, idom[q]) // 埋伏在最近的支配者
		}
	}
	return res
}

// https://leetcode.cn/problems/disconnect-path-in-a-binary-matrix-by-at-most-one-flip/
func isPossibleToCutPath(grid [][]int) bool {
	ROW, COL := len(grid), len(grid[0])
	edges := [][]int{}
	for i := 0; i < ROW; i++ {
		for j := 0; j < COL; j++ {
			if grid[i][j] == 0 {
				continue
			}
			if i+1 < ROW && grid[i+1][j] == 1 {
				edges = append(edges, []int{i*COL + j, (i+1)*COL + j})
			}
			if j+1 < COL && grid[i][j+1] == 1 {
				edges = append(edges, []int{i*COL + j, i*COL + j + 1})
			}
		}
	}

	idom := BuildDominatorTree(ROW*COL, edges, 0)
	return idom[ROW*COL-1] != 0
}

// 求每个点能支配的个数(包括自己)
// https://www.luogu.com.cn/problem/P5180
func DominatorTreeSize(n int, edges [][]int, root int) []int {
	parents := BuildDominatorTree(n, edges, root)
	adjList := make([][]int, n)
	for i := 0; i < n; i++ {
		p := parents[i]
		if p != -1 && p != i {
			adjList[p] = append(adjList[p], i)
		}
	}

	subSize := make([]int, n)
	var dfs func(cur int)
	dfs = func(cur int) {
		subSize[cur] = 1
		for _, next := range adjList[cur] {
			dfs(next)
			subSize[cur] += subSize[next]
		}
	}
	dfs(root)

	return subSize
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)
	edges := make([][]int, 0, m)
	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		u, v = u-1, v-1
		edges = append(edges, []int{u, v})
	}

	res := DominatorTreeSize(n, edges, 0)
	for i := 0; i < n; i++ {
		fmt.Fprint(out, res[i], " ")
	}
}
