// 离散对数-群G作用 discrete_log
// !时间复杂度 O(sqrt(mod))

package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func main() {
	// https://yukicoder.me/problems/no/1339
	// 求有理数1/n的循环节长度

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for t := 0; t < T; t++ {
		var n int
		fmt.Fscan(in, &n)

		// 让模与10互质
		for n%2 == 0 {
			n /= 2
		}
		for n%5 == 0 {
			n /= 5
		}

		// !模为1时不存在模逆元
		if n == 1 {
			fmt.Fprintln(out, 1)
			continue
		}

		// !求解 10^k mod n = 1 的最小整数解 k, 其中 1<=k<=n+1
		k := DiscreteLogGroup(
			func() G { return 1 % n },
			func(g1, g2 G) G { return g1 * g2 % n },
			func(g G) G { return modInv(g, n) },
			10,
			1,
			1,
			n+10,
		)
		fmt.Fprintln(out, k)
	}
}

type G = int

// 给定一个群G，群元素a, b in G，求解 a^n = b 的最小非负整数解 n.
//	返回在[lower, higher)中的第一个解，如果不存在则返回-1.
//  !可以理解为 a 经过多少次群运算后可以到达 b.
func DiscreteLogGroup(
	/** 群G */
	e func() G,
	op func(g1, g2 G) G,
	inv func(g G) G,
	a G,
	b G,
	lower, higher int,
) int {
	if lower >= higher {
		return -1
	}
	UNIT := e()
	s := UNIT
	mp := make(map[G]int)
	aPow := func(n int) G {
		p := a
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

	s = op(s, aPow(lower))
	LIM := higher - lower
	K := int(math.Sqrt(float64(LIM))) + 1
	for i := 0; i <= K; i++ {
		key := s
		if _, ok := mp[key]; !ok {
			mp[key] = i
		}
		if i != K {
			s = op(s, a)
		}
	}

	a = inv(aPow(K))
	for i := 0; i <= K; i++ {
		key := b
		if v, ok := mp[key]; ok {
			res := i*K + v + lower
			if res >= higher {
				return -1
			}
			return res
		}
		b = op(b, a)
	}

	return -1
}

func exgcd(a, b int) (gcd, x, y int) {
	if b == 0 {
		return a, 1, 0
	}
	gcd, y, x = exgcd(b, a%b)
	y -= a / b * x
	return
}

// 模逆元，注意模为1时不存在逆元
func modInv(a, mod int) int {
	gcd, x, _ := exgcd(a, mod)
	if gcd != 1 {
		panic(fmt.Sprintf("no inverse element for %d", a))
	}
	return (x%mod + mod) % mod
}

// 求离散对数 a^n = b (mod m) 的最小非负整数解 n.
// 其中模数为质数.
//  如果不存在则返回-1.
func ModLog(a, b, mod int) int {
	a %= mod
	b %= mod
	p := 1 % mod
	for k := 0; k < 32; k++ {
		if p == b {
			return k
		}
		p = p * a % mod
	}
	if a == 0 || b == 0 {
		return -1
	}

	gcd := func(a, b int) int {
		for b != 0 {
			a, b = b, a%b
		}
		return a
	}

	g := gcd(mod, p)
	if b%g != 0 {
		return -1
	}
	mod /= g
	a %= mod
	b %= mod
	if gcd(b, mod) > 1 {
		return -1
	}

	return DiscreteLogGroup(
		func() G { return 1 % mod },
		func(g1, g2 G) G { return (g1 * g2) % mod },
		func(g G) G { return modInv(g, mod) },
		a,
		b,
		32,
		mod,
	)
}

func testModLog() {
	// https://www.luogu.com.cn/problem/P4195
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	for {
		var base, p, target int
		_, _ = fmt.Fscan(in, &base, &p, &target)
		if base == target && target == p && p == 0 {
			break
		}
		res := ModLog(base, target, p)
		if res == -1 {
			fmt.Fprintln(out, "No Solution")
		} else {
			fmt.Fprintln(out, res)
		}
	}
}
