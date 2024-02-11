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
	// 放置国旗()
	// yosupo()
	// yuki274()
	P5782()
	// yuki1078()
}

// P5782 [POI2001] 和平委员会
// https://www.luogu.com.cn/problem/P5782
// 每个党在议会中有 2个代表。代表从1编号到2n。编号为 2i-1 和 2i 的代表属于第i个党派。
// 委员会必须满足下列条件：
// 1. 每个党派都在委员会中恰有 1 个代表。
// 2. 如果 2 个代表彼此厌恶，则他们不能都属于委员会。
// 如果不能创立委员会，则输出信息 NIE。
// 若能够成立，则输出包括 n 个从区间 1 到 2n 选出的整数，按升序写出，每行一个，这些数字为委员会中代表的编号。
//
// !命题i代表选择第i个代表
func P5782() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)
	badPairs := make([][2]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &badPairs[i][0], &badPairs[i][1])
		badPairs[i][0]--
		badPairs[i][1]--
	}

	ts := NewTwoSat(2 * n)
	for i := 0; i < n; i++ {
		a, b := 2*i, 2*i+1
		// ts.AddOr(a, b)
		// ts.AddNand(a, b)
		ts.AddXor(a, b)
	}
	for i := 0; i < m; i++ {
		u, v := badPairs[i][0], badPairs[i][1]
		ts.AddNand(u, v)
	}

	res, ok := ts.Solve()
	if !ok {
		fmt.Fprintln(out, "NIE")
		return
	}
	for i := 0; i < 2*n; i++ {
		if res[i] {
			fmt.Fprintln(out, i+1)
		}
	}
}

// https://atcoder.jp/contests/practice2/tasks/practice2_h
// 1-N号旗设置位置(放置国旗)
// 第i号旗可以设置在xi位置或者yi位置
// !任意两面旗距离需要大于D
// 是否可以设置旗子
// 1≤N≤1000
// D,Xi,Yi<=1e9
//
// 命题i:第i个棋子放在xi位置,检查是否满足条件
func 放置国旗() {
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

// No.274 The Wall-墙壁
// https://yukicoder.me/problems/no/274
// 用n个砖块构建墙壁,每个砖长度为m
// 每个砖块上有一段颜色,砖块可以180度旋转
// 将这n个砖块拼接成一面墙壁,使得每一列存在颜色的部分最多只有一个
// 问是否能够拼接成功
// n<=2000,m<=4000
// n个命题分别为[第i个砖块不旋转]
// 每两个之间验证四种情况
func yuki274() {
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

// No.470 Inverse S+T Problem
// https://yukicoder.me/problems/no/470/editorial
// 给定n个长度为3的字符串,由小写字母和大写字母组成
// 求2n个非空串Si和Ti，使得 S[i] + T[i] = U[i],且这2n个串都不相同
// n<=1e5

// 长度为3的字符串划分成两个非空部分,要么是1/2,要么是2/1 => 2SAT
// n很大的时候一定有重复串,可以排除
// !n不大的时候用2SAT解决 :命题i代表S[i]用1个字符
func yuki470() {
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
			s1, t1, s2, t2 = w1[0:1], w1[1:], w2[0:2], w2[2:]
			if s1 == t2 || t1 == s2 {
				ts.AddNand(i, ts.Rev(j))
			}

			// 2 1
			s1, t1, s2, t2 = w1[0:2], w1[2:], w2[0:1], w2[1:]
			if s1 == t2 || t1 == s2 {
				ts.AddNand(ts.Rev(i), j)
			}

			// 2 2
			s1, t1, s2, t2 = w1[0:2], w1[2:], w2[0:2], w2[2:]
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
			s, t := words[i][0:2], words[i][2:]
			fmt.Fprint(out, s, " ", t)
		}
		fmt.Fprintln(out)
	}
}

// No.1078 I love Matrix Construction
// https://yukicoder.me/problems/no/1078
// 给定长为n的数组S,T,U
// 问能否构建出满足以下条件的n*n矩阵A
// 1. A[i][j] 为0/1
// 2. 对所有的 (i,j), A[S[i]][j]+A[j][T[i]]*2 != U[i]
// n<=500
func yuki1078() {
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

// https://judge.yosupo.jp/problem/two_sat
// N 変数  M 節の 2 Sat が与えられる。充足可能か判定し、可能ならば割り当てを一つ求めてください。
func yosupo() {
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
			a++
			a = ts.Rev(-a)
		} else {
			a--
		}
		if b < 0 {
			b++
			b = ts.Rev(-b)
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
	size  int
	graph [][]int32
}

func NewTwoSat(n int) *TwoSat {
	return &TwoSat{size: n, graph: make([][]int32, n*2)}
}

// u -> v <=> !v -> !u
func (ts *TwoSat) AddIf(u, v int) {
	ts.AddDirectedEdge(u, v)
	ts.AddDirectedEdge(ts.Rev(v), ts.Rev(u))
}

// u,v 中至少有一个为真.
// u or v <=> !u -> v
func (ts *TwoSat) AddOr(u, v int) {
	ts.AddIf(ts.Rev(u), v)
}

// u,v 中恰好有一个为真, 一个为假.
// u xor v <=> u -> !v, v -> !u, !u -> v, !v -> u
func (ts *TwoSat) AddXor(u, v int) {
	ts.AddOr(u, v)
	ts.AddNand(u, v)
}

// u,v 不同时为真.
// u nand v <=> u -> !v
func (ts *TwoSat) AddNand(u, v int) {
	ts.AddIf(u, ts.Rev(v))
}

// 手动添加边(不推荐).常用于优化建图时.
func (ts *TwoSat) AddDirectedEdge(u, v int) {
	ts.graph[u] = append(ts.graph[u], int32(v))
}

// u <=> !u -> u
func (ts *TwoSat) SetTrue(u int) {
	ts.AddDirectedEdge(ts.Rev(u), u)
}

// !u <=> u -> !u
func (ts *TwoSat) SetFalse(u int) {
	ts.AddDirectedEdge(u, ts.Rev(u))
}

func (ts *TwoSat) Rev(u int) int {
	if u >= ts.size {
		return u - ts.size
	}
	return u + ts.size
}

func (ts *TwoSat) Solve() (res []bool, ok bool) {
	_, belong := StronglyConnectedComponentInt32(ts.graph)
	res = make([]bool, ts.size)
	for i := 0; i < int(ts.size); i++ {
		if belong[i] == belong[ts.Rev(i)] {
			return
		}
		res[i] = belong[i] > belong[ts.Rev(i)]
	}
	ok = true
	return
}

// 有向图强连通分量分解.
func StronglyConnectedComponentInt32(graph [][]int32) (count int32, belong []int32) {
	n := int32(len(graph))
	belong = make([]int32, n)
	low := make([]int32, n)
	order := make([]int32, n)
	for i := range order {
		order[i] = -1
	}
	now := int32(0)
	path := []int32{}

	var dfs func(int32)
	dfs = func(v int32) {
		low[v] = now
		order[v] = now
		now++
		path = append(path, v)
		for _, to := range graph[v] {
			if order[to] == -1 {
				dfs(to)
				low[v] = min32(low[v], low[to])
			} else {
				low[v] = min32(low[v], order[to])
			}
		}
		if low[v] == order[v] {
			for {
				u := path[len(path)-1]
				path = path[:len(path)-1]
				order[u] = n
				belong[u] = count
				if u == v {
					break
				}
			}
			count++
		}
	}

	for i := int32(0); i < n; i++ {
		if order[i] == -1 {
			dfs(i)
		}
	}
	for i := int32(0); i < n; i++ {
		belong[i] = count - 1 - belong[i]
	}
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
func min32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func max32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
