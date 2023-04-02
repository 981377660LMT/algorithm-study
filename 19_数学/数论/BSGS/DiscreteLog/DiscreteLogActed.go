// 离散对数-幺半群X与集合S作用
// DiscreteLogActed
// !时间复杂度 O(sqrt(mod))

package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func main() {
	abc186_e()
}

func abc270_g() {
	// https://atcoder.jp/contests/abc270/tasks/abc270_g
	// 给定一个由以下递推式定义的数列
	// 判断是否存在一个i使得Xi = TARGET，并求出满足条件的最小i。
	// !X[i] = (MUL * X[i-1] + ADD) % P if i > 0 else START
	// 100组样例,P<=1e9且为素数

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for t := 0; t < T; t++ {
		var P, MUL, ADD, START, TARGET int
		fmt.Fscan(in, &P, &MUL, &ADD, &START, &TARGET)
		fmt.Fprintln(out, SolveAffine(MUL, ADD, START, TARGET, P, 0)) // LOWER从0开始
	}
}

func abc222_g() {
	// https://atcoder.jp/contests/abc222/tasks/abc222_g
	// 有一个数字序列：2，22，222，2222，...，其中第i项是一个i位数，由数字2组成。
	// 对于每个给定的测试用例T，找到序列中第一个是K的倍数的项的索引。如果没有这样的项，则输出-1。
	// T<=200 k<=1e8
	// 类似于abc270_g
	// 等价于:
	// !X[i] = (10 * X[i-1] + 2) % P if i > 0 else 0
	// MUL=10,ADD=1,START=0,TARGET=0,P=K,LOWER=1

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for t := 0; t < T; t++ {
		var mod int // kth
		fmt.Fscan(in, &mod)
		fmt.Fprintln(out, SolveAffine(10, 2, 0, 0, mod, 1))
	}
}

func abc186_e() {
	// https://atcoder.jp/contests/abc186/tasks/abc186_e
	// 一个圆周上有n个椅子,其中一个是好的椅子
	// 开始时,从S号椅子出发,每次都会移动到下k个椅子
	// 问何时可以坐到好的椅子上
	// 如果不能,则输出-1
	// T<=100
	// n<=1e9
	// 等价于:
	// !X[i] = (X[i-1] + k) % P if i > 0 else 0
	// MUL=1,ADD=k,START=S,TARGET=0,P=n,LOWER=1

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for t := 0; t < T; t++ {
		var MOD, START, ADD int
		fmt.Fscan(in, &MOD, &START, &ADD)
		fmt.Fprintln(out, SolveAffine(1, ADD, START, 0, MOD, 1))
	}
}

// 是否存在一个i使得X[i] = TARGET，并求出满足条件的最小整数i, i >= lower.
// X[i] = (MUL * X[i-1] + ADD) % P if i > 0 else START.
//  如果不存在,返回-1.
func SolveAffine(mul, add int, start, target, mod int, lower int) int {
	if start == target && lower <= 0 {
		return 0
	}
	if mul == 0 {
		if add == target && lower <= 1 {
			return 1
		}
		return -1
	}
	if mod == 1 {
		if target == 0 && lower <= 1 {
			return 1
		}
		return -1
	}

	return DiscreteLogActed(
		func() G { return G{1, 0} },
		func(g1, g2 G) G { return G{g1.mul * g2.mul % mod, (g1.add*g2.mul + g2.add) % mod} },
		func(s S, g G) S { return (s*g.mul + g.add) % mod },
		G{mul, add},
		start,
		target,
		lower,
		mod+10,
	)
}

type G = struct{ mul, add int }
type S = int

// !求解`集合S`在`幺半群X`上的离散对数
// 给定一个幺半群X，一个作用在X上的集合S.
// 对于给定的幺半群元素 x in X，集合元素 s, t in S，求解 (x^n)*s = t 的整数解 n.
//  返回在[lower, higher)中的第一个解，如果不存在则返回-1.
//  - e: 幺半群X的单位元
//  - op: 幺半群X的结合律
//  - act: 幺半群X在集合S上的作用, 例如仿射变换X对整数集S的作用为
//         func act(s S, g G) S { return s*g.mul + g.add }
//  !即至少经过多少次群G的act作用，s可以到达t.
func DiscreteLogActed(
	/** 幺半群X */
	e func() G,
	op func(g1, g2 G) G,

	/** 集合S与幺半群X的作用 */
	act func(s S, g G) S,

	x G,
	s S,
	t S,
	lower int,
	higher int,
) int {
	if lower >= higher {
		return -1
	}
	UNIT := e()
	set := make(map[S]struct{})
	xpow := func(n int) G {
		p := x
		res := UNIT
		for n > 0 {
			if n&1 == 1 {
				res = op(res, p)
			}
			p = op(p, p)
			n /= 2
		}
		return res
	}

	ht := t
	s = act(s, xpow(lower))
	LIM := higher - lower
	K := int(math.Sqrt(float64(LIM))) + 1
	for k := 0; k < K; k++ {
		t = act(t, x)
		set[t] = struct{}{}
	}

	y := xpow(K)
	failed := false
	for k := 0; k <= K; k++ {
		s1 := act(s, y)
		if _, ok := set[s1]; ok {
			for i := 0; i < K; i++ {
				if s == ht {
					cand := k*K + i + lower
					if cand >= higher {
						return -1
					}
					return cand
				}
				s = act(s, x)
			}
			if failed {
				return -1
			}
			failed = true
		}
		s = s1
	}
	return -1
}
