// 01 on Tree（树上一类全序问题/树上拓扑序问题）
// https://blog.csdn.net/ez_lcw/article/details/120202160

package main

import (
	"bufio"
	"fmt"
	"os"
)

// 01 on Tree.
// 给定一棵树，每个节点有一个 Monoid. Monoid 满足全序关系.
// 求出一种结点的拓扑序，最大化 Monoid 的聚合值.
func OptimalProductOnTree[V any](
	tree [][]int32, root int32,
	values []V, op func(a, b V) V, less func(a, b V) bool,
) (order []int32, best V) {
	values = append(values[:0:0], values...)
	n := int32(len(tree))

	parent := make([]int32, n)
	parent[root] = -1
	{
		var stack []int32
		stack = append(stack, root)
		for len(stack) > 0 {
			v := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			for _, to := range tree[v] {
				if to == parent[v] {
					continue
				}
				parent[to] = v
				stack = append(stack, to)
			}
		}
	}

	head, tail, next := make([]int32, n), make([]int32, n), make([]int32, n)
	for v := int32(0); v < n; v++ {
		head[v], tail[v], next[v] = v, v, -1
	}

	uf := NewUnionFindArraySimple32(n)

	size := int32(1)
	for size < n {
		size <<= 1
	}
	seg := make([]int32, 2*size)
	for i := range seg {
		seg[i] = -1
	}
	update := func(i int32) {
		a, b := seg[i<<1], seg[i<<1|1]
		if a == -1 {
			seg[i] = b
		} else if b == -1 {
			seg[i] = a
		} else {
			if less(values[a], values[b]) {
				seg[i] = a
			} else {
				seg[i] = b
			}
		}
	}
	set := func(i, x int32) {
		i += size
		seg[i] = x
		for i > 1 {
			i >>= 1
			update(i)
		}
	}
	for i := int32(0); i < n; i++ {
		if i != root {
			seg[size+i] = i
		}
	}
	for i := size - 1; i > 0; i-- {
		update(i)
	}

	for i := int32(0); i < n-1; i++ {
		v := seg[1]
		a, b := head[uf.Find(v)], tail[uf.Find(v)]
		p := uf.Find(parent[a])
		c, d := head[p], tail[p]
		pv := op(values[p], values[v])
		uf.Union(p, v, nil)
		w := uf.Find(v)
		values[w] = pv
		head[v], tail[v], next[d] = c, b, a
		av, ap := int32(-1), int32(-1)
		if v == w && uf.Find(v) != uf.Find(root) {
			av = v
		}
		if p == w && uf.Find(p) != uf.Find(root) {
			ap = p
		}
		set(v, av)
		set(p, ap)
	}

	order = []int32{root}
	for nxt := next[root]; nxt != -1; nxt = next[nxt] {
		order = append(order, nxt)
	}
	best = values[uf.Find(root)]
	return
}

type UnionFindArraySimple32 struct {
	Part int32
	n    int32
	data []int32
}

func NewUnionFindArraySimple32(n int32) *UnionFindArraySimple32 {
	data := make([]int32, n)
	for i := int32(0); i < n; i++ {
		data[i] = -1
	}
	return &UnionFindArraySimple32{Part: n, n: n, data: data}
}

func (u *UnionFindArraySimple32) Union(key1, key2 int32, beforeMerge func(big, small int32)) bool {
	root1, root2 := u.Find(key1), u.Find(key2)
	if root1 == root2 {
		return false
	}
	if u.data[root1] > u.data[root2] {
		root1, root2 = root2, root1
	}
	if beforeMerge != nil {
		beforeMerge(root1, root2)
	}
	u.data[root1] += u.data[root2]
	u.data[root2] = root1
	u.Part--
	return true
}

func (u *UnionFindArraySimple32) Find(key int32) int32 {
	root := key
	for u.data[root] >= 0 {
		root = u.data[root]
	}
	for key != root {
		key, u.data[key] = u.data[key], root
	}
	return root
}

func (u *UnionFindArraySimple32) Size(key int32) int32 {
	return -u.data[u.Find(key)]
}

func exgcd(a, b int) (gcd, x, y int) {
	if b == 0 {
		return a, 1, 0
	}
	gcd, y, x = exgcd(b, a%b)
	y -= a / b * x
	return
}

// 注意模为1时不存在逆元
func modInv(a, mod int) int {
	gcd, x, _ := exgcd(a, mod)
	if gcd != 1 {
		panic(fmt.Sprintf("no inverse element for %d", a))
	}
	return (x%mod + mod) % mod
}

//
//
//
//
//
//
//
//
//
//
//
//

func main() {
	// agc023_f()
	abc376_g()
}

// F - 01 on Tree/Tree01/01Tree
// https://atcoder.jp/contests/agc023/tasks/agc023_f
// 给定一颗权值为 0/1 的树，求出一种结点的拓扑序，最小化逆序对数.
// n<=2e5.
//
// !合并优先级: 0的个数/1的个数
func agc023_f() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var N int32
	fmt.Fscan(in, &N)
	P := make([]int32, N-1)
	for i := int32(0); i < N-1; i++ {
		fmt.Fscan(in, &P[i])
		P[i]--
	}
	V := make([]int32, N)
	for i := int32(0); i < N; i++ {
		fmt.Fscan(in, &V[i])
	}

	tree := make([][]int32, N)
	for i := int32(0); i < N-1; i++ {
		tree[P[i]] = append(tree[P[i]], i+1)
	}
	root := int32(0)

	type E struct {
		c0, c1, inv int
	}
	values := make([]E, N)
	for i := int32(0); i < N; i++ {
		if V[i] == 0 {
			values[i] = E{c0: 1}
		} else {
			values[i] = E{c1: 1}
		}
	}
	op := func(a, b E) E {
		res := E{}
		res.c0 = a.c0 + b.c0
		res.c1 = a.c1 + b.c1
		res.inv = a.inv + b.inv + a.c1*b.c0
		return res
	}
	less := func(a, b E) bool {
		return a.c0*b.c1 > a.c1*b.c0
	}

	_, best := OptimalProductOnTree(tree, root, values, op, less)
	fmt.Fprintln(out, best.inv)
}

// G - Treasure Hunting
// https://atcoder.jp/contests/abc376/editorial/11196
// 给定一颗有n+1个顶点的树. 顶点编号为0,1,...,n.
// 顶点0是根. 顶点i的父亲是Pi.
// 顶点i有一个价值Vi.(1<=i<=n)
// !求一个拓扑序，使得 ∑i*Vi 最小.
// !转换为01Tree问题.
// !等价于 每个顶点的价值为(0,0,...,0,1)，有V[i]个0.最小化逆序对数.
func abc376_g() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	const MOD int = 998244353

	solve := func(n int32, parents []int32, values []int32) int {
		tree := make([][]int32, n+1)
		for i := int32(0); i < n; i++ {
			tree[parents[i]] = append(tree[parents[i]], i+1)
		}
		root := int32(0)

		type E = struct {
			c0, c1, inv int
		}
		newValues := make([]E, n+1)
		for i := int32(1); i < n+1; i++ {
			newValues[i] = E{c0: int(values[i-1]) % MOD, c1: 1}
		}
		op := func(a, b E) E {
			res := E{}
			res.c0 = (a.c0 + b.c0) % MOD
			res.c1 = (a.c1 + b.c1) % MOD
			res.inv = (a.inv + b.inv + a.c1*b.c0) % MOD
			return res
		}
		less := func(a, b E) bool {
			return a.c0*b.c1 > a.c1*b.c0
		}

		_, best := OptimalProductOnTree(tree, root, newValues, op, less)

		res := best.inv
		sum := 0
		for _, v := range values {
			sum += int(v)
		}
		res *= modInv(sum, MOD)
		res += 1
		res %= MOD
		return res
	}

	var T int32
	fmt.Fscan(in, &T)
	for i := int32(0); i < T; i++ {
		var N int32
		fmt.Fscan(in, &N)
		P := make([]int32, N)
		for i := int32(0); i < N; i++ {
			fmt.Fscan(in, &P[i])
		}
		V := make([]int32, N)
		for i := int32(0); i < N; i++ {
			fmt.Fscan(in, &V[i])
		}
		fmt.Fprintln(out, solve(N, P, V))
	}
}
