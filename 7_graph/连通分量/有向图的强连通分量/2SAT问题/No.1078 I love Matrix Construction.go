// https://yukicoder.me/problems/no/1078
// 给定长为n的数组S,T,U
// 问能否构建出满足以下条件的n*n矩阵A
// 1. A[i][j] 为0/1
// 2. 对所有的 (i,j), A[S[i]][j]+A[j][T[i]]*2 != U[i]
// n<=500

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

	var n int
	fmt.Fscan(in, &n)

	S := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &S[i])
		S[i]--
	}
	T := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &T[i])
		T[i]--
	}
	U := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &U[i])
	}

	// 条件i为A[i][j]取0
	ts := NewTwoSat(n * n)
	for i := 0; i < n; i++ {
		si := S[i]
		ti := T[i]
		for j := 0; j < n; j++ {
			pos1 := si*n + j
			pos2 := j*n + ti

			if U[i] == 0 {
				ts.AddNand(pos1, pos2) // 0,0
			} else if U[i] == 1 {
				ts.AddNand(ts.Rev(pos1), pos2) //1,0
			} else if U[i] == 2 {
				ts.AddNand(pos1, ts.Rev(pos2)) //0,1
			} else if U[i] == 3 {
				ts.AddNand(ts.Rev(pos1), ts.Rev(pos2)) //1,1
			}

		}
	}

	res, ok := ts.Solve()
	if !ok {
		fmt.Fprintln(out, -1)
		return
	}

	matrix := make([][]int, n)
	for i := 0; i < n; i++ {
		matrix[i] = make([]int, n)
	}

	for i, v := range res {
		if !v {
			matrix[i/n][i%n] = 1
		}
	}

	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			fmt.Fprint(out, matrix[i][j], " ")
		}
		fmt.Fprintln(out)
	}
}

type TwoSat struct {
	sz  int
	scc *scc
}

func NewTwoSat(n int) *TwoSat {
	return &TwoSat{sz: n, scc: newScc(n + n)}
}

// u -> v <=> !v -> !u
func (ts *TwoSat) AddIf(u, v int) {
	ts.scc.AddEdge(u, v, 1)
	ts.scc.AddEdge(ts.Rev(v), ts.Rev(u), 1)
}

// u or v <=> !u -> v
func (ts *TwoSat) AddOr(u, v int) {
	ts.AddIf(ts.Rev(u), v)
}

// u nand v <=> u -> !v
func (ts *TwoSat) AddNand(u, v int) {
	ts.AddIf(u, ts.Rev(v))
}

// u <=> !u -> u
func (ts *TwoSat) SetTrue(u int) {
	ts.scc.AddEdge(ts.Rev(u), u, 1)
}

// !u <=> u -> !u
func (ts *TwoSat) SetFalse(u int) {
	ts.scc.AddEdge(u, ts.Rev(u), 1)
}

func (ts *TwoSat) Rev(u int) int {
	if u >= ts.sz {
		return u - ts.sz
	}
	return u + ts.sz
}

func (ts *TwoSat) Solve() (res []bool, ok bool) {
	ts.scc.Build()
	res = make([]bool, ts.sz)
	for i := 0; i < ts.sz; i++ {
		if ts.scc.Comp[i] == ts.scc.Comp[ts.Rev(i)] {
			return
		}
		res[i] = ts.scc.Comp[i] > ts.scc.Comp[ts.Rev(i)]
	}
	ok = true
	return
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

type scc struct {
	G     [][]int // 原图
	Comp  []int   //每个顶点所属的强连通分量的编号
	rg    [][]int
	order []int
	used  []bool
}

func newScc(n int) *scc {
	return &scc{G: make([][]int, n)}
}

func (scc *scc) AddEdge(from, to, cost int) {
	scc.G[from] = append(scc.G[from], to)
}

func (scc *scc) Build() {
	scc.rg = make([][]int, len(scc.G))
	for i := range scc.G {
		for _, e := range scc.G[i] {
			scc.rg[e] = append(scc.rg[e], i)
		}
	}

	scc.Comp = make([]int, len(scc.G))
	for i := range scc.Comp {
		scc.Comp[i] = -1
	}
	scc.used = make([]bool, len(scc.G))
	for i := range scc.G {
		scc.dfs(i)
	}
	for i, j := 0, len(scc.order)-1; i < j; i, j = i+1, j-1 {
		scc.order[i], scc.order[j] = scc.order[j], scc.order[i]
	}

	ptr := 0
	for _, v := range scc.order {
		if scc.Comp[v] == -1 {
			scc.rdfs(v, ptr)
			ptr++
		}
	}

}

// 获取顶点k所属的强连通分量的编号
func (scc *scc) Get(k int) int {
	return scc.Comp[k]
}

func (scc *scc) dfs(idx int) {
	tmp := scc.used[idx]
	scc.used[idx] = true
	if tmp {
		return
	}
	for _, e := range scc.G[idx] {
		scc.dfs(e)
	}
	scc.order = append(scc.order, idx)
}

func (scc *scc) rdfs(idx int, cnt int) {
	if scc.Comp[idx] != -1 {
		return
	}
	scc.Comp[idx] = cnt
	for _, e := range scc.rg[idx] {
		scc.rdfs(e, cnt)
	}
}
