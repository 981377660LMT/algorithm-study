// 离散对数-群G与集合S作用 discrete_log
// !时间复杂度 O(sqrt(mod))

// 注意几个边界:
// 1.invGcd需要模和a互质
// 2.仿射变换群的mul为0时不存在逆元(注意,数列的线性递推式可以写成仿射变换群)
//   !X[i] = (MUL * X[i-1] + ADD) % P if i > 0 else START
// 3.模为1时不存在模逆元

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
	// !X[i] = 2 * (10 * X[i-1] + 1) % P if i > 0 else 0
	// MUL=10,ADD=1,START=0,TARGET=0,P=K,LOWER=1

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for t := 0; t < T; t++ {
		var mod int // kth
		fmt.Fscan(in, &mod)

		// 为了让10可逆,mod不能包含2/5
		if mod%4 == 0 || mod%5 == 0 {
			fmt.Fprintln(out, -1)
			continue
		}
		if mod%2 == 0 { // !所有数去除2,变为2*(1,11,111,1111,...)
			mod /= 2
		}

		fmt.Fprintln(out, SolveAffine(10, 1, 0, 0, mod, 1))

	}
}

func abc186_e() {
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
// !mul 必须和 mod 互质.
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
		// !仿射变换的逆元
		func(g G) G {
			mul, add := g.mul, g.add
			mul = modInv(mul, mod)
			return G{mul, mul * (mod - add) % mod}
		},
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

// !求解`集合S`在`群G`上的离散对数
// 给定一个群G，一个作用在G上的集合S.
// 对于给定的群元素 x in G，集合元素 s, t in S，求解 (x^n)*s = t 的整数解 n.
//  返回在[lower, higher)中的第一个解，如果不存在则返回-1.
//  - e: 群G的单位元
//  - op: 群G的结合律
//  - inv: 群G的逆元
//  - act: 群G在集合S上的作用, 例如仿射变换G对整数集S的作用为
//         func act(s S, g G) S { return s*g.mul + g.add }
//  !即至少经过多少次群G的act作用，s可以到达t.
func DiscreteLogActed(
	/** 群G */
	e func() G,
	op func(g1, g2 G) G,
	inv func(g G) G,

	/** 集合S与群G的作用 */
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
	mp := make(map[S]int)
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

	s = act(s, xpow(lower))
	LIM := higher - lower
	K := int(math.Sqrt(float64(LIM))) + 1
	for i := 0; i <= K; i++ {
		key := s
		if _, ok := mp[key]; !ok {
			mp[key] = i
		}
		if i != K {
			s = act(s, x)
		}
	}

	x = inv(xpow(K))
	for i := 0; i <= K; i++ {
		key := t
		if v, ok := mp[key]; ok {
			res := i*K + v + lower
			if res >= higher {
				return -1
			}
			return res
		}
		t = act(t, x)
	}

	return -1
}

func Pow(base, exp, mod int) int {
	if exp == -1 {
		return modInv(base, mod)
	}

	base %= mod
	res := 1
	for ; exp > 0; exp >>= 1 {
		if exp&1 == 1 {
			res = res * base % mod
		}
		base = base * base % mod
	}
	return res
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
