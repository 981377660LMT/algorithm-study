// !https://github.dev/EndlessCheng/codeforces-go/blob/3dd70515200872705893d52dc5dad174f2c3b5f3/copypasta/graph.go#L2799
package main

import (
	"bufio"
	"fmt"
	"os"
)

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
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	boolToInt := func(b bool) int {
		if b {
			return 1
		}
		return 0
	}

	var n, m int
	fmt.Fscan(in, &n, &m)
	twoSat := NewTwoSat(n)
	// !假设初始时每个命题为xi为真
	for i := 0; i < m; i++ {
		var i, iState, j, jState int
		fmt.Fscan(in, &i, &iState, &j, &jState)
		i, j = i-1, j-1
		twoSat.AddLimit(i, iState, j, jState) // !加边方式1：至少满足一个
	}

	twoSat.Build()

	res := twoSat.Work()
	if res == nil {
		fmt.Fprintln(out, "IMPOSSIBLE")
	} else {
		fmt.Fprintln(out, "POSSIBLE")
		for _, b := range res {
			fmt.Fprint(out, boolToInt(b), " ")
		}
	}

}

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

type TwoSat struct {
	n      int
	graph  [][]int
	rGraph [][]int
	sccId  []int // 每个点所在的scc分组编号
}

func NewTwoSat(n int) *TwoSat {
	return &TwoSat{
		n:      n,
		graph:  make([][]int, 2*n),
		rGraph: make([][]int, 2*n),
	}
}

// 加边方式：根据命题间的限制条件 '至少满足一个' 添加边.
//  i=a 和 j=b 两个条件至少满足一个(a b 为 0/1 表示 假/真)
//  0 <= i < n, 0 <= j < n
func (ts *TwoSat) AddLimit(i, a, j, b int) {
	n, graph, rGraph := ts.n, ts.graph, ts.rGraph
	notU, v := i+a*n, j+(b^1)*n // ¬A⇒B
	graph[notU] = append(graph[notU], v)
	rGraph[v] = append(rGraph[v], notU)
	notV, u := j+b*n, i+(a^1)*n // ¬B⇒A
	graph[notV] = append(graph[notV], u)
	rGraph[u] = append(rGraph[u], notV)
}

// 加边方式2：根据命题的推导关系添加边.
// 如果 u => v (u成立可以推导出v成立)，那么就添加边 u => v 以及 ¬v => ¬u.。
// 注意当 u/v 为真命题时, 0 <= u/v < n
// 当 u/v 为假命题时, n <= u/v < 2*n
func (ts *TwoSat) AddEdge(u, v int) {
	var notU, notV int
	if u < ts.n {
		notU = u + ts.n
	} else {
		notU = u - ts.n
	}
	if v < ts.n {
		notV = v + ts.n
	} else {
		notV = v - ts.n
	}
	ts.graph[u] = append(ts.graph[u], v)
	ts.rGraph[v] = append(ts.rGraph[v], u)
	ts.graph[notV] = append(ts.graph[notV], notU)
	ts.rGraph[notU] = append(ts.rGraph[notU], notV)
}

func (ts *TwoSat) Build() {
	if ts.sccId == nil {
		ts.sccId = ts.getSccIds()
	}
}

// !返回每个命题 xi 的真假性，若无解(存在矛盾)则返回 nil
func (ts *TwoSat) Work() []bool {
	n, sccId := ts.n, ts.sccId
	res := make([]bool, n)
	for i, id := range sccId[:n] {
		// x 和 ¬x 处于同一个 SCC 时无解（因为 x ⇔ ¬x）
		if id == sccId[i+n] {
			return nil
		}

		// sid[x] > sid[¬x] ⇔ (¬x ⇒ x) ⇔ x 为真 (对应拓扑图上可以到达真命题的点)
		// sid[x] < sid[¬x] ⇔ (x ⇒ ¬x) ⇔ x 为假
		res[i] = id > sccId[i+n]
	}
	return res
}

func (ts *TwoSat) getSccIds() []int {
	m, graph, rGraph := ts.n*2, ts.graph, ts.rGraph
	rDfsOrder := make([]int, 0, m)
	visited := make([]bool, m)
	var dfs func(int)
	dfs = func(cur int) {
		visited[cur] = true
		for _, w := range graph[cur] {
			if !visited[w] {
				dfs(w)
			}
		}
		rDfsOrder = append(rDfsOrder, cur)
	}
	for i, ok := range visited {
		if !ok {
			dfs(i)
		}
	}

	visited = make([]bool, m)
	var group []int
	dfs = func(v int) {
		visited[v] = true
		group = append(group, v)
		for _, w := range rGraph[v] {
			if !visited[w] {
				dfs(w)
			}
		}
	}

	sccIDs := make([]int, m)
	cid := 0
	for i := len(rDfsOrder) - 1; i >= 0; i-- {
		if v := rDfsOrder[i]; !visited[v] {
			group = []int{}
			dfs(v)
			for _, v := range group {
				sccIDs[v] = cid
			}
			cid++
		}
	}

	return sccIDs
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
