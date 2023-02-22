// https://atcoder.jp/contests/practice2/tasks/practice2_h
// 1-N号旗设置位置(放置国旗)
// 第i号旗可以设置在xi位置或者yi位置
// !任意两面旗距离需要大于D
// 是否可以设置旗子
// 1≤N≤1000
// D,Xi,Yi<=1e9

// 命题i:第i个棋子放在xi位置,检查是否满足条件

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

	var n, d int
	fmt.Fscan(in, &n, &d)
	x, y := make([]int, n), make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &x[i], &y[i])
	}

	ts := NewTwoSat(n)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			// 00
			if abs(x[i]-x[j]) < d {
				ts.AddNand(i, j)
			}

			// 01
			if abs(x[i]-y[j]) < d {
				ts.AddNand(i, ts.Rev(j))
			}

			// 10
			if abs(y[i]-x[j]) < d {
				ts.AddNand(ts.Rev(i), j)
			}

			// 11
			if abs(y[i]-y[j]) < d {
				ts.AddNand(ts.Rev(i), ts.Rev(j))
			}
		}
	}

	res, ok := ts.Solve()
	if !ok {
		fmt.Fprintln(out, "No")
		return
	}
	fmt.Fprintln(out, "Yes")
	for i := 0; i < n; i++ {
		if res[i] {
			fmt.Fprintln(out, x[i])
		} else {
			fmt.Fprintln(out, y[i])
		}
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
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
