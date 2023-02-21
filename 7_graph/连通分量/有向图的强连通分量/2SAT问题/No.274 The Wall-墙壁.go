// https://yukicoder.me/problems/no/274
// 用n个砖块构建墙壁,每个砖长度为m
// 每个砖块上有一段颜色,砖块可以180度旋转
// 将这n个砖块拼接成一面墙壁,使得每一列存在颜色的部分最多只有一个
// 问是否能够拼接成功
// n<=2000,m<=4000
// n个命题分别为[第i个砖块不旋转]
// 每两个之间验证四种情况

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
	color := make([][]int, n) // [left,right]
	for i := 0; i < n; i++ {
		var l, r int
		fmt.Fscan(in, &l, &r)
		color[i] = []int{l, r}
	}

	isOverlapped := func(left1, right1, left2, right2 int) bool {
		start := max(left1, left2)
		end := min(right1, right2)
		return start <= end
	}

	// n个命题分别为[第i个砖块不旋转]
	// 每两个之间验证四种情况
	ts := NewTwoSat(n)
	for i := 0; i < n; i++ {
		left1, right1 := color[i][0], color[i][1]
		revLeft1, revRight1 := m-1-right1, m-1-left1
		for j := i + 1; j < n; j++ {
			left2, right2 := color[j][0], color[j][1]
			revLeft2, revRight2 := m-1-right2, m-1-left2
			// 1. 两个砖块都不旋转
			if isOverlapped(left1, right1, left2, right2) {
				ts.AddNand(i, j)
			}

			// 2. 第一个砖块旋转,第二个砖块不旋转
			if isOverlapped(revLeft1, revRight1, left2, right2) {
				ts.AddNand(ts.Rev(i), j)
			}

			// 3. 第一个砖块不旋转,第二个砖块旋转
			if isOverlapped(left1, right1, revLeft2, revRight2) {
				ts.AddNand(i, ts.Rev(j))
			}

			// 4. 两个砖块都旋转
			if isOverlapped(revLeft1, revRight1, revLeft2, revRight2) {
				ts.AddNand(ts.Rev(i), ts.Rev(j))
			}
		}
	}

	_, ok := ts.Solve()
	if !ok {
		fmt.Fprintln(out, "NO")
		return
	}

	fmt.Fprintln(out, "YES")

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
