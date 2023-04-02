// 离散对数-群G作用：
// DiscreteLogGroup/ModLog
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

func yuki() {
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
