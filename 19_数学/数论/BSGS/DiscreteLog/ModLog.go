package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func main() {
	// https://judge.yosupo.jp/problem/discrete_logarithm_mod
	// https://www.luogu.com.cn/problem/P4195
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for t := 0; t < T; t++ {
		var x, y, m int
		fmt.Fscan(in, &x, &y, &m)
		fmt.Fprintln(out, ModLog(x, y, m))
	}

	// for {
	// 	var a, p, b int
	// 	fmt.Fscan(in, &a, &p, &b)
	// 	if a == 0 && p == 0 && b == 0 {
	// 		break
	// 	}
	// 	res := ModLog(a, b, p)
	// 	if res == -1 {
	// 		fmt.Fprintln(out, "No Solution")
	// 	} else {
	// 		fmt.Fprintln(out, res)
	// 	}
	// }
}

type G = int
type S = int

// 求离散对数 a^n = b (mod m) 的最小非负整数解 n.
//  如果不存在则返回-1.
func ModLog(a, b, mod int) int {
	a %= mod
	b %= mod
	return discreteLogActed(
		func() G { return 1 % mod },
		func(g1, g2 G) G { return (g1 * g2) % mod },
		func(s, g G) S { return s * g % mod },
		a,
		1,
		b,
		0,
		mod,
	)
}

func discreteLogActed(
	e func() G,
	op func(g1, g2 G) G,

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
