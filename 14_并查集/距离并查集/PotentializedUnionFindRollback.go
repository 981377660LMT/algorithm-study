// 可撤销带权并查集.
// Api:
//   - NewPotentializedUnionFind(n int32, e func() E, op func(E, E) E, inv func(E) E) *PotentializedUnionFind[E]
//   - (uf *PotentializedUnionFind[E]) Union(a, b int32, x E) bool
//     !合并a, b所在的集合, 并且满足 P[a] - P[b] = x.
//   - (uf *PotentializedUnionFind[E]) Find(v int32) (root int32, diff E)
//     !返回v所在集合的根节点, 以及 P[v] - P[root].
//   - (uf *PotentializedUnionFind[E]) Diff(a, b int32) (E, bool)
//     !返回 P[a] - P[b] 以及是否在同一个集合中.
//
//	 - GetTime() int32
//	 - Rollback(time int32)

package main

import (
	"bufio"
	"fmt"

	"os"
)

func main() {
	yuki2293()
	// yosupo()
	// yosupoUnionfindwithPotentialNonCommutativeGroup()
}

// No.2293 無向辺 2SAT
// https://yukicoder.me/problems/no/2293
// 给定n个逻辑变量X[i].每个变量可以为0或1.
// 给定一个栈和q个操作.
// 每个操作形如：
// 1 u v: 向栈中添加条件 (X[u]=1 or X[v]=0) and (X[u]=0 or X[v]=1).
// 2 u v: 向栈中添加条件 (X[u]=0 or X[v]=0) and (X[u]=1 or X[v]=1).
// 3: 将栈清空.
// 每次结束后，求满足栈中所有条件的逻辑变量的方案数.
func yuki2293() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	const MOD int = 998244353
	pow := func(base, exp, mod int) int {
		base %= mod
		res := 1 % mod
		for ; exp > 0; exp >>= 1 {
			if exp&1 == 1 {
				res = res * base % mod
			}
			base = base * base % mod
		}
		return res
	}

	e := func() uint8 { return 0 }
	op := func(a, b uint8) uint8 { return a ^ b }
	inv := func(a uint8) uint8 { return a }

	var n, q int32
	fmt.Fscan(in, &n, &q)
	uf := NewPotentializedUnionFindRollback(n, e, op, inv)

	initTime := uf.GetTime()
	part := n
	ok := true
	for i := int32(0); i < q; i++ {
		var op int32
		fmt.Fscan(in, &op)
		if op == 1 {
			var u, v int32
			fmt.Fscan(in, &u, &v)
			u, v = u-1, v-1

			diff, same := uf.Diff(u, v)
			if !same {
				uf.Union(u, v, 0)
				part--
			} else {
				if diff == 1 {
					ok = false
				}
			}
		} else if op == 2 {
			var u, v int32
			fmt.Fscan(in, &u, &v)
			u, v = u-1, v-1

			diff, same := uf.Diff(u, v)
			if !same {
				uf.Union(u, v, 1)
				part--
			} else {
				if diff == 0 {
					ok = false
				}
			}
		} else {
			uf.Rollback(initTime)
			part = n
			ok = true
		}

		if !ok {
			fmt.Fprintln(out, 0)
		} else {
			fmt.Fprintln(out, pow(2, int(part), MOD))
		}
	}
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

	uf := NewPotentializedUnionFindRollback(n, e, op, inv)

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

// https://judge.yosupo.jp/problem/unionfind_with_potential_non_commutative_group
func yosupoUnionfindwithPotentialNonCommutativeGroup() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	const MOD int = 998244353

	var n, q int32
	fmt.Fscan(in, &n, &q)

	type E = [2][2]int

	// monoid为矩阵乘法.
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

	uf := NewPotentializedUnionFindRollback(n, e, op, inv)

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

type item[E any] struct {
	root int32
	diff E
}

// 可撤销势能并查集/距离并查集.
type PotentializedUnionFindRollback[E comparable] struct {
	data *RollbackArray[item[E]]
	e    func() E
	op   func(E, E) E
	inv  func(E) E
}

func NewPotentializedUnionFindRollback[E comparable](n int32, e func() E, op func(E, E) E, inv func(E) E) *PotentializedUnionFindRollback[E] {
	initData := make([]item[E], n)
	for i := int32(0); i < n; i++ {
		initData[i] = item[E]{root: -1, diff: e()}
	}
	return &PotentializedUnionFindRollback[E]{
		data: NewRollbackArrayFrom(initData),
		e:    e, op: op, inv: inv,
	}
}

func (uf *PotentializedUnionFindRollback[E]) GetTime() int32 {
	return uf.data.GetTime()
}

func (uf *PotentializedUnionFindRollback[E]) Rollback(time int32) {
	uf.data.Rollback(time)
}

// P[a] - P[b] = x
func (uf *PotentializedUnionFindRollback[E]) Union(a, b int32, x E) bool {
	v1, x1 := uf.Find(b)
	v2, x2 := uf.Find(a)
	if v1 == v2 {
		return false
	}
	item1, item2 := uf.data.Get(v1), uf.data.Get(v2)
	s1, s2 := -item1.root, -item2.root
	if s1 < s2 {
		s1, s2 = s2, s1
		v1, v2 = v2, v1
		x1, x2 = x2, x1
		x = uf.inv(x)
	}
	x = uf.op(x1, x)
	x = uf.op(x, uf.inv(x2))
	uf.data.Set(v2, item[E]{root: v1, diff: x})
	uf.data.Set(v1, item[E]{root: -(s1 + s2), diff: uf.e()})
	return true
}

// root，P[v] - P[root]
func (uf *PotentializedUnionFindRollback[E]) Find(v int32) (root int32, diff E) {
	diff = uf.e()
	for {
		item := uf.data.Get(v)
		if item.root < 0 {
			break
		}
		diff = uf.op(item.diff, diff)
		v = item.root
	}
	root = v
	return
}

// Diff(a, b) = P[a] - P[b]
func (uf *PotentializedUnionFindRollback[E]) Diff(a, b int32) (E, bool) {
	ru, xu := uf.Find(b)
	rv, xv := uf.Find(a)
	if ru != rv {
		return uf.e(), false
	}
	return uf.op(uf.inv(xu), xv), true
}

type HistoryItem[V comparable] struct {
	index int32
	value V
}

type RollbackArray[V comparable] struct {
	n       int32
	data    []V
	history []HistoryItem[V]
}

func NewRollbackArray[V comparable](n int32, f func(i int32) V) *RollbackArray[V] {
	data := make([]V, n)
	for i := int32(0); i < n; i++ {
		data[i] = f(i)
	}
	return &RollbackArray[V]{
		n:    n,
		data: data,
	}
}

func NewRollbackArrayFrom[V comparable](data []V) *RollbackArray[V] {
	return &RollbackArray[V]{n: int32(len(data)), data: data}
}

func (r *RollbackArray[V]) GetTime() int32 {
	return int32(len(r.history))
}

func (r *RollbackArray[V]) Rollback(time int32) {
	for i := int32(len(r.history)) - 1; i >= time; i-- {
		pair := r.history[i]
		r.data[pair.index] = pair.value
	}
	r.history = r.history[:time]
}

func (r *RollbackArray[V]) Undo() bool {
	if len(r.history) == 0 {
		return false
	}
	pair := r.history[len(r.history)-1]
	r.history = r.history[:len(r.history)-1]
	r.data[pair.index] = pair.value
	return true
}

func (r *RollbackArray[V]) Get(index int32) V {
	return r.data[index]
}

func (r *RollbackArray[V]) Set(index int32, value V) bool {
	if r.data[index] == value {
		return false
	}
	r.history = append(r.history, HistoryItem[V]{index: index, value: r.data[index]})
	r.data[index] = value
	return true
}

func (r *RollbackArray[V]) GetAll() []V {
	return append(r.data[:0:0], r.data...)
}

func (r *RollbackArray[V]) Len() int32 {
	return r.n
}
