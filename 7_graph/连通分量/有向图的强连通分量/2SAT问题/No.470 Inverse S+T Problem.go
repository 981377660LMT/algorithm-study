// https://yukicoder.me/problems/no/470/editorial
// 给定n个长度为3的字符串,由小写字母和大写字母组成
// 求2n个非空串Si和Ti，使得 S[i] + T[i] = U[i],且这2n个串都不相同
// n<=1e5

// 长度为3的字符串划分成两个非空部分,要么是1/2,要么是2/1 => 2SAT
// n很大的时候一定有重复串,可以排除
// !n不大的时候用2SAT解决 :命题i代表S[i]用1个字符

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

	words := make([]string, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &words[i])
	}

	// 只使用a-z和A-Z的字符,一个字符一定有重复的
	if n > 26*2 {
		fmt.Fprintln(out, "Impossible")
		return
	}

	ts := NewTwoSat(n)
	// 枚举所有的对,看哪些s[i]不能同时用一个字符划分
	// かぶさる可能性のあるものを反転させたものをグラフに追加する  ??
	for i := 0; i < n; i++ {
		w1 := words[i]
		for j := i + 1; j < n; j++ {
			w2 := words[j]

			// 1 1
			s1, t1, s2, t2 := w1[0:1], w1[1:], w2[0:1], w2[1:]
			if s1 == s2 || t1 == t2 {
				ts.AddNand(i, j)
			}

			// 1 2
			s1, t1, s2, t2 = w1[0:1], w1[1:], w2[1:], w2[0:1]
			if s1 == t2 || t1 == s2 {
				ts.AddNand(i, ts.Rev(j))
			}

			// 2 1
			s1, t1, s2, t2 = w1[1:], w1[0:1], w2[0:1], w2[1:]
			if s1 == t2 || t1 == s2 {
				ts.AddNand(ts.Rev(i), j)
			}

			// 2 2
			s1, t1, s2, t2 = w1[1:], w1[0:1], w2[1:], w2[0:1]
			if s1 == s2 || t1 == t2 {
				ts.AddNand(ts.Rev(i), ts.Rev(j))
			}
		}
	}

	res, ok := ts.Solve()
	if !ok {
		fmt.Fprintln(out, "Impossible")
		return
	}

	for i := 0; i < n; i++ {
		if res[i] {
			s, t := words[i][0:1], words[i][1:]
			fmt.Fprint(out, s, " ", t)
		} else {
			s, t := words[i][1:], words[i][0:1]
			fmt.Fprint(out, s, " ", t)
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
