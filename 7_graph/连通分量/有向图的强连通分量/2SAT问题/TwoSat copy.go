// !在实际问题中，2-SAT问题在大多数时候表现成以下形式：
// 有N对物品，每对物品中必须选取一个，也只能选取一个，
// 并且它们之间存在某些`限制关系`，比如「A和B不能同时选取」、「C和D必须同时选取」。
// （如某两个物品不能都选，某两个物品不能都不选，某两个物品必须且只能选一个，某个物品必选）等，
// 这时，可以将每对物品当成一个布尔值（选取第一个物品相当于0，选取第二个相当于1），
// 如果所有的限制关系最多只对两个物品进行限制，
// 则它们都可以转化成9种基本限制关系，从而转化为2-SAT模型。
// !非常像种类并查集
// https://www.cnblogs.com/kuangbin/archive/2012/10/05/2712429.html

// 有n个布尔变量x1～xn，另有m个需要满足的条件，每个条件的形式都是`xi为true/false` 或 `xj为true/false`。
// 比如「x1为真或x3为假」、「x7为假或x2为假」。
// https://www.luogu.com.cn/problem/P4782
// https://ei1333.github.io/library/graph/others/two-satisfiability.hpp
// 2-Satisfiability (2-SAT)
// !https://zhuanlan.zhihu.com/p/50211772
// todo https://www.luogu.com.cn/blog/85514/post-2-sat-xue-xi-bi-ji
// 2-SAT 总结 by kuangbin https://www.cnblogs.com/kuangbin/archive/2012/10/05/2712429.html
// !NOTE: 一些建边的转换(命题为真对应0-n-1,命题为假对应n-2*n-1)：
//       A 为真          (A)     ¬A⇒A     注：A ⇔ A∨A ⇔ ¬A⇒A∧¬A⇒A ⇔ ¬A⇒A
//       A 为假          (¬A)    A⇒¬A
//       A 为真 B 就为真          A⇒B, ¬B⇒¬A
//       A 为假 B 就为假          ¬A⇒¬B, B⇒A
//       !A,B 至少存在一个 (A|B)    ¬A⇒B, ¬B⇒A 意思是一个为假的时候，另一个一定为真 https://www.luogu.com.cn/problem/P4782
//       A,B 不能同时存在 (¬A|¬B)  A⇒¬B, B⇒¬A 就是上面的式子替换了一下（一个为真，另一个一定为假）
//       A,B 必须且只一个 (A^B)    A⇒¬B, B⇒¬A, ¬A⇒B, ¬B⇒A
//       A,B 同时或都不在 (¬(A^B)) A⇒B, B⇒A, ¬A⇒¬B, ¬B⇒¬A
// !NOTE: 单独的条件 x为a 可以用 (x为a)∨(x为a) 来表示
// 模板题 https://www.luogu.com.cn/problem/P4782
// 建边练习【模板代码】 https://codeforces.com/contest/468/problem/B
// 定义 Ai 表示「选 Xi」，这样若两个旗子 i j 满足 |Xi-Xj|<D 时，就相当于 Ai,Aj 至少一个为假。其他情况类似 https://atcoder.jp/contests/practice2/tasks/practice2_h
// !github.com/EndlessCheng/codeforces-go

// TwoSatisfiability(N): N 個のリテラルで初期化する.
// AddIf(u, v): 条件 u ならば v を追加する.
// AddOr(u, v): 条件 u または v が true を追加する.
// AddNand(u, v): 条件 u または v が false を追加する.
// SetTrue(u): 条件 u が true を追加する.
// SetFalse(u): 条件 u が false を追加する.
// Rev(u): 変数 u の否定を返す.
// Solve(): 充足可能か判定し, 可能なら各リテラルの割り当ての例を格納した配列, 不能なら空配列を返す.
// O(V+E)

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// https://judge.yosupo.jp/problem/two_sat
	// N 変数  M 節の 2 Sat が与えられる。充足可能か判定し、可能ならば割り当てを一つ求めてください。
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var p, cnf string
	var n, m int
	fmt.Fscan(in, &p, &cnf, &n, &m)
	ts := NewTwoSat(n)
	for i := 0; i < m; i++ {
		var a, b, c int
		fmt.Fscan(in, &a, &b, &c)
		if a < 0 {
			a = ts.Rev(-a - 1)
		} else {
			a--
		}
		if b < 0 {
			b = ts.Rev(-b - 1)
		} else {
			b--
		}
		ts.AddOr(a, b)
	}

	res, ok := ts.Solve()
	if !ok {
		fmt.Fprintln(out, "s UNSATISFIABLE")
		return
	}
	fmt.Fprintln(out, "s SATISFIABLE")
	fmt.Fprint(out, "v ")
	for i, v := range res {
		if v {
			fmt.Fprint(out, i+1, " ")
		} else {
			fmt.Fprint(out, -(i + 1), " ")
		}
	}
	fmt.Fprintln(out, 0)
}

type TwoSat struct {
	sz  int
	scc *StronglyConnectedComponents
}

func NewTwoSat(n int) *TwoSat {
	return &TwoSat{sz: n, scc: NewStronglyConnectedComponents(n + n)}
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

type WeightedEdge struct{ from, to, cost int }
type StronglyConnectedComponents struct {
	G     [][]WeightedEdge // 原图
	Dag   [][]WeightedEdge // 强连通分量缩点后的顶点和边组成的DAG
	Comp  []int            //每个顶点所属的强连通分量的编号
	Group [][]int          // 每个强连通分量所包含的顶点
	rg    [][]WeightedEdge
	order []int
	used  []bool
}

func NewStronglyConnectedComponents(n int) *StronglyConnectedComponents {
	return &StronglyConnectedComponents{G: make([][]WeightedEdge, n)}
}

func (scc *StronglyConnectedComponents) AddEdge(from, to, cost int) {
	scc.G[from] = append(scc.G[from], WeightedEdge{from, to, cost})
}

func (scc *StronglyConnectedComponents) Build() {
	scc.rg = make([][]WeightedEdge, len(scc.G))
	for i := range scc.G {
		for _, e := range scc.G[i] {
			scc.rg[e.to] = append(scc.rg[e.to], WeightedEdge{e.to, e.from, e.cost})
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

	dag := make([][]WeightedEdge, ptr)
	for i := range scc.G {
		for _, e := range scc.G[i] {
			x, y := scc.Comp[e.from], scc.Comp[e.to]
			if x == y {
				continue
			}
			dag[x] = append(dag[x], WeightedEdge{x, y, e.cost})
		}
	}
	scc.Dag = dag

	scc.Group = make([][]int, ptr)
	for i := range scc.G {
		scc.Group[scc.Comp[i]] = append(scc.Group[scc.Comp[i]], i)
	}
}

// 获取顶点k所属的强连通分量的编号
func (scc *StronglyConnectedComponents) Get(k int) int {
	return scc.Comp[k]
}

func (scc *StronglyConnectedComponents) dfs(idx int) {
	tmp := scc.used[idx]
	scc.used[idx] = true
	if tmp {
		return
	}
	for _, e := range scc.G[idx] {
		scc.dfs(e.to)
	}
	scc.order = append(scc.order, idx)
}

func (scc *StronglyConnectedComponents) rdfs(idx int, cnt int) {
	if scc.Comp[idx] != -1 {
		return
	}
	scc.Comp[idx] = cnt
	for _, e := range scc.rg[idx] {
		scc.rdfs(e.to, cnt)
	}
}
