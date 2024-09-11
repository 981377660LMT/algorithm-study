// Api:
//  - NewPotentializedUnionFind(n int32, e func() E, op func(E, E) E, inv func(E) E) *PotentializedUnionFind[E]
//  - (uf *PotentializedUnionFind[E]) Diff(a, b int32) (diff E, same bool)
//    !返回 P[a] - P[b] 以及是否在同一个集合中.
//  - (uf *PotentializedUnionFind[E]) Union(a, b int32, x E) bool
//    !合并a, b所在的集合, 并且满足 P[a] - P[b] = x.
//  - (uf *PotentializedUnionFind[E]) Find(v int32) (root int32, diff E)
//    !返回v所在集合的根节点, 以及 P[v] - P[root].

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// dsl1B()
	// yuki1502()
	yuki2294()
	// demo()
	// abc280F()
	// yosupoUnionfindwithPotential()
	// yosupoUnionfindwithPotentialNonCommutativeGroup()
}

func demo() {
	// e, op, inv := func() int32 { return 0 }, func(a, b int32) int32 { return a + b }, func(a int32) int32 { return -a }
	e := func() int { return 0 }
	op := func(a, b int) int { return a ^ b }
	inv := func(a int) int { return a }
	uf := NewPotentializedUnionFind(10, e, op, inv)
	uf.Union(2, 1, 10)
	fmt.Println(uf.Find(0))
	fmt.Println(uf.Find(1))
	fmt.Println(uf.Find(2))
	fmt.Println(uf.Diff(2, 1))
	fmt.Println(uf.Diff(1, 2))
	fmt.Println(uf.Find(3))
}

// https://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=DSL_1_B
// relate(x,y,z): A[y] = A[x] + z
// diff(x,y): A[y] - A[x]
func dsl1B() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int32
	fmt.Fscan(in, &n, &q)

	e := func() int { return 0 }
	op := func(a, b int) int { return a + b }
	inv := func(a int) int { return -a }
	uf := NewPotentializedUnionFind(n, e, op, inv)

	relate := func(x, y int32, z int) {
		uf.Union(y, x, z)
	}

	diff := func(x, y int32) (int, bool) {
		return uf.Diff(y, x)
	}

	for i := int32(0); i < q; i++ {
		var op int32
		fmt.Fscan(in, &op)
		if op == 0 {
			var x, y int32
			var z int
			fmt.Fscan(in, &x, &y, &z)
			relate(x, y, z)
		} else {
			var x, y int32
			fmt.Fscan(in, &x, &y)
			res, same := diff(x, y)
			if same {
				fmt.Fprintln(out, res)
			} else {
				fmt.Fprintln(out, "?")
			}
		}
	}
}

// No.1502 Many Simple Additions
// https://yukicoder.me/problems/no/1502
func yuki1502() {

}

// No.2294 Union Path Query (Easy，异或和距离，所有点对的异或和)
// https://yukicoder.me/problems/no/2294
// 给定一张n个点的无向带权图.两点间的距离为异或和.
// 给定一个初始点X0.
// 给定q个查询，每个查询形如：
// 1 v w: 将顶点v与顶点X0用一条边权为w的边连接.保证连接后的图中没有环.
// 2 u v: 输出顶点u到顶点v的距离.如果无法到达，输出-1.
// 3 v: 求v所在联通分量的所有点对距离异或和模998244353.
// 4 add: 将X0增加add，然后对N取模.
// N<=2e5.w<=1e9.q<=2e5.
func yuki2294() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	const MOD int = 998244353

	var N, X, Q int32
	fmt.Fscan(in, &N, &X, &Q)

	const LOG int = 30
	type V = [LOG][2]int // 每一位染黑/白
	data := make([]V, N)
	for i := int32(0); i < N; i++ {
		for j := 0; j < LOG; j++ {
			data[i][j][0] = 1
			data[i][j][1] = 0
		}
	}

	// MonoidXor
	e := func() int { return 0 }
	op := func(a, b int) int { return a ^ b }
	inv := func(a int) int { return a }
	uf := NewPotentializedUnionFind(N, e, op, inv)

	link := func(a, b int32, w int) {
		ra, da := uf.Find(a)
		rb, db := uf.Find(b)
		if ra == rb {
			return
		}
		nd := da ^ db ^ w
		uf.Union2(a, b, w, func(big, small int32) {
			for i := 0; i < LOG; i++ {
				if (nd>>i)&1 == 1 {
					data[big][i][0] += data[small][i][1]
					data[big][i][1] += data[small][i][0]
				} else {
					data[big][i][0] += data[small][i][0]
					data[big][i][1] += data[small][i][1]
				}
			}
		})
	}

	dist := func(a, b int32) (int, bool) {
		return uf.Diff(a, b)
	}

	pairDist := func(a int32) int {
		root, _ := uf.Find(a)
		res := 0
		for i := 0; i < LOG; i++ {
			a, b := data[root][i][0], data[root][i][1]
			pair := a * b % MOD
			res += (1 << i) % MOD * pair % MOD
			res %= MOD
		}
		return res
	}

	for i := int32(0); i < Q; i++ {
		var op int32
		fmt.Fscan(in, &op)
		if op == 1 {
			var u int32
			var w int
			fmt.Fscan(in, &u, &w)
			link(u, X, w)
		} else if op == 2 {
			var u, v int32
			fmt.Fscan(in, &u, &v)
			d, ok := dist(u, v)
			if !ok {
				fmt.Fprintln(out, -1)
			} else {
				fmt.Fprintln(out, d)
				X = int32((int(X) + d) % int(N))
			}
		} else if op == 3 {
			var v int32
			fmt.Fscan(in, &v)
			fmt.Fprintln(out, pairDist(v))
		} else {
			var add int
			fmt.Fscan(in, &add)
			X = int32((int(X) + add) % int(N))
		}
	}
}

// F - Pay or Receive (爬山，势能模型)
// https://atcoder.jp/contests/abc280/tasks/abc280_f
// 给定n个点和m条无向边.
// 对每条边，从a->b, 可以获得c的收益，从b->a, 会获得-c的收益.
// 给定q个查询，每个查询形如a->b.
// 问从 a->b 可以获得的最大收益.
// 如果无法到达，输出nan.
// 如果可以获得无限收益，输出inf.
//
// 初始对每条边：
// 如果不连通 -> 连通并附加约束
// 如果连通 -> 检查约束是否满足，不满足则在正环上
func abc280F() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m, q int32
	fmt.Fscan(in, &n, &m, &q)

	e := func() int { return 0 }
	op := func(a, b int) int { return a + b }
	inv := func(a int) int { return -a }
	uf := NewPotentializedUnionFind(n, e, op, inv)

	inPosCycle := make([]bool, n)
	for i := int32(0); i < m; i++ {
		var a, b int32
		var c int
		fmt.Fscan(in, &a, &b, &c)
		a, b = a-1, b-1
		diff, same := uf.Diff(b, a)
		if !same {
			uf.Union(b, a, c) // P[b] - P[a] = c
			continue
		}
		if diff != c {
			inPosCycle[a] = true
			inPosCycle[b] = true
		}
	}

	// !transfer
	for i := int32(0); i < n; i++ {
		if inPosCycle[i] {
			root, _ := uf.Find(i)
			inPosCycle[root] = true
		}
	}

	for i := int32(0); i < q; i++ {
		var a, b int32
		fmt.Fscan(in, &a, &b)
		a, b = a-1, b-1
		diff, same := uf.Diff(b, a)
		if !same {
			fmt.Fprintln(out, "nan")
		} else {
			root, _ := uf.Find(a)
			if inPosCycle[root] {
				fmt.Fprintln(out, "inf")
			} else {
				fmt.Fprintln(out, diff)
			}
		}
	}
}

// UnionfindwithPotential
// https://judge.yosupoUnionfindwithPotential.jp/problem/unionfind_with_potential
// 0 u v x: 判断A[u]=A[v]+x(mod Mod)是否成立. 如果与现有信息矛盾,则不进行任何操作,否则将该条件加入.
// 1 u v: 输出A[u]-A[v].如果不能确定,输出-1.
func yosupoUnionfindwithPotential() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	const MOD int = 998244353

	var n, q int32
	fmt.Fscan(in, &n, &q)

	e := func() int { return 0 }
	op := func(a, b int) int {
		res := (a + b) % MOD
		if res < 0 {
			res += MOD
		}
		return res
	}
	inv := func(a int) int {
		res := -a % MOD
		if res < 0 {
			res += MOD
		}
		return res
	}

	uf := NewPotentializedUnionFind(n, e, op, inv)

	for i := int32(0); i < q; i++ {
		var op int32
		fmt.Fscan(in, &op)
		if op == 0 {
			var u, v int32
			var x int
			fmt.Fscan(in, &u, &v, &x)
			diff, same := uf.Diff(u, v)
			valid := !same || diff == x
			if valid {
				fmt.Fprintln(out, 1)
			} else {
				fmt.Fprintln(out, 0)
			}
			if !same {
				uf.Union(u, v, x)
			}
		} else {
			var u, v int32
			fmt.Fscan(in, &u, &v)
			if diff, ok := uf.Diff(u, v); ok {
				fmt.Fprintln(out, diff)
			} else {
				fmt.Fprintln(out, -1)
			}
		}
	}
}

// 不可交换群.
// https://judge.yosupo.jp/problem/unionfind_with_potential_non_commutative_group
func yosupoUnionfindwithPotentialNonCommutativeGroup() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	const MOD int = 998244353

	var n, q int32
	fmt.Fscan(in, &n, &q)

	type E = [2][2]int

	// 群为矩阵乘法，不可交换.
	e := func() E { return E{{1, 0}, {0, 1}} }
	op := func(a, b E) E {
		v00 := a[0][0]*b[0][0] + a[0][1]*b[1][0]
		v01 := a[0][0]*b[0][1] + a[0][1]*b[1][1]
		v10 := a[1][0]*b[0][0] + a[1][1]*b[1][0]
		v11 := a[1][0]*b[0][1] + a[1][1]*b[1][1]
		v00, v01, v10, v11 = v00%MOD, v01%MOD, v10%MOD, v11%MOD
		return E{{v00, v01}, {v10, v11}}
	}
	inv := func(a E) E {
		v00, v01, v10, v11 := a[0][0], a[0][1], a[1][0], a[1][1]
		v00, v01, v10, v11 = v11, -v01, -v10, v00
		if v01 < 0 {
			v01 += MOD
		}
		if v10 < 0 {
			v10 += MOD
		}
		return E{{v00, v01}, {v10, v11}}
	}

	uf := NewPotentializedUnionFind(n, e, op, inv)

	for i := int32(0); i < q; i++ {
		var op int32
		fmt.Fscan(in, &op)
		if op == 0 {
			var u, v int32
			var v00, v01, v10, v11 int
			fmt.Fscan(in, &u, &v, &v00, &v01, &v10, &v11)
			x := E{{v00, v01}, {v10, v11}}
			diff, same := uf.Diff(u, v) // !注意非交换性 P[u] = P[v] * x
			valid := !same || diff == x
			if valid {
				fmt.Fprintln(out, 1)
			} else {
				fmt.Fprintln(out, 0)
			}
			if !same {
				uf.Union(u, v, x)
			}
		} else {
			var u, v int32
			fmt.Fscan(in, &u, &v)
			if diff, ok := uf.Diff(u, v); ok {
				v00, v01, v10, v11 := diff[0][0], diff[0][1], diff[1][0], diff[1][1]
				fmt.Fprintln(out, v00, v01, v10, v11)
			} else {
				fmt.Fprintln(out, -1)
			}
		}
	}
}

// 势能并查集/距离并查集.
type PotentializedUnionFind[E any] struct {
	n, Part int32
	parents []int32
	sizes   []int32
	values  []E
	e       func() E
	op      func(E, E) E
	inv     func(E) E
}

func NewPotentializedUnionFind[E any](n int32, e func() E, op func(E, E) E, inv func(E) E) *PotentializedUnionFind[E] {
	values, parents, sizes := make([]E, n), make([]int32, n), make([]int32, n)
	for i := int32(0); i < n; i++ {
		parents[i] = i
		sizes[i] = 1
		values[i] = e()
	}
	return &PotentializedUnionFind[E]{n: n, Part: n, parents: parents, sizes: sizes, values: values, e: e, op: op, inv: inv}
}

// P[a] - P[b] = x
func (uf *PotentializedUnionFind[E]) Union(a, b int32, x E) bool {
	v1, x1 := uf.Find(b)
	v2, x2 := uf.Find(a)
	if v1 == v2 {
		return false
	}
	if uf.sizes[v1] < uf.sizes[v2] {
		v1, v2 = v2, v1
		x1, x2 = x2, x1
		x = uf.inv(x)
	}
	x = uf.op(x1, x)
	x = uf.op(x, uf.inv(x2))
	uf.values[v2] = x
	uf.parents[v2] = v1
	uf.sizes[v1] += uf.sizes[v2]
	uf.Part--
	return true
}

// 返回根节点root, 以及 diff = P[v] - P[root]
func (uf *PotentializedUnionFind[E]) Find(v int32) (root int32, diff E) {
	diff = uf.e()
	vs, ps := uf.values, uf.parents
	for v != ps[v] {
		diff = uf.op(vs[v], diff)
		diff = uf.op(vs[ps[v]], diff)
		vs[v] = uf.op(vs[ps[v]], vs[v])
		ps[v] = ps[ps[v]]
		v = ps[v]
	}
	root = v
	return
}

// Diff(a, b) = P[a] - P[b]
func (uf *PotentializedUnionFind[E]) Diff(a, b int32) (E, bool) {
	ru, xu := uf.Find(b)
	rv, xv := uf.Find(a)
	if ru != rv {
		return uf.e(), false
	}
	return uf.op(uf.inv(xu), xv), true
}

// P[a] - P[b] = x
func (uf *PotentializedUnionFind[E]) Union2(a, b int32, x E, beforeUnion func(big, small int32)) bool {
	v1, x1 := uf.Find(b)
	v2, x2 := uf.Find(a)
	if v1 == v2 {
		return false
	}
	if uf.sizes[v1] < uf.sizes[v2] {
		v1, v2 = v2, v1
		x1, x2 = x2, x1
		x = uf.inv(x)
	}
	if beforeUnion != nil {
		beforeUnion(v1, v2)
	}
	x = uf.op(x1, x)
	x = uf.op(x, uf.inv(x2))
	uf.values[v2] = x
	uf.parents[v2] = v1
	uf.sizes[v1] += uf.sizes[v2]
	uf.Part--
	return true
}
