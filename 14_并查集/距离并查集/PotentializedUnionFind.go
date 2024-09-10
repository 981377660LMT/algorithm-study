package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// demo()
	yosupo()
}

func demo() {
	e, op, inv := func() int32 { return 0 }, func(a, b int32) int32 { return a + b }, func(a int32) int32 { return -a }
	uf := NewPotentializedUnionFind(10, e, op, inv)
	uf.Union(1, 2, 10)
	fmt.Println(uf.Find(0))
	fmt.Println(uf.Find(1))
	fmt.Println(uf.Find(2))
	fmt.Println(uf.Diff(1, 2))
	fmt.Println(uf.Find(3))
}

// UnionfindwithPotential
// https://judge.yosupo.jp/problem/unionfind_with_potential
// 0 u v x: 判断A[u]=A[v]+x(mod Mod)是否成立. 如果与现有信息矛盾,则不进行任何操作,否则将该条件加入.
// 1 u v: 输出A[u]-A[v].如果不能确定,输出-1.
func yosupo() {
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
	v1, x1 := uf.Find(a)
	v2, x2 := uf.Find(b)
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

func (uf *PotentializedUnionFind[E]) Find(v int32) (root int32, dist E) {
	dist = uf.e()
	vs, ps := uf.values, uf.parents
	for v != ps[v] {
		dist = uf.op(vs[v], dist)
		dist = uf.op(vs[ps[v]], dist)
		vs[v] = uf.op(vs[ps[v]], vs[v])
		ps[v] = ps[ps[v]]
		v = ps[v]
	}
	root = v
	return
}

// Diff(a, b) = P[a] - P[b]
func (uf *PotentializedUnionFind[E]) Diff(a, b int32) (E, bool) {
	ru, xu := uf.Find(a)
	rv, xv := uf.Find(b)
	if ru != rv {
		return uf.e(), false
	}
	return uf.op(uf.inv(xu), xv), true
}
