package main

import (
	"fmt"
	"math/bits"
)

func main() {
	fi := NewFastInv(998244353)
	fmt.Println(fi.Inv(2), (998244353+1)/2)

}

const B int = 2048

// 光速逆元.
type FastInv struct {
	mod        int
	isModPrime bool
	ipow2      [64]int
	pre        [B / 2]int
}

// mod必须是小于2**30的奇数.
func NewFastInv(mod int) *FastInv {
	res := &FastInv{}
	res.mod = mod
	res.isModPrime = res._testPrime(mod)

	res.ipow2[0] = 1 % mod
	for i := 1; i < 64; i++ {
		res.ipow2[i] = res.ipow2[i-1] * ((mod + 1) / 2) % mod
	}
	res.pre[0] = 1 % mod
	cal := func(i int) int {
		x, y, u, v, t, tmp := i, mod, 1, 0, 0, 0
		for y > 0 {
			t = x / y
			x -= t * y
			u -= t * v
			tmp = x
			x = y
			y = tmp
			tmp = u
			u = v
			v = tmp
		}
		if u < 0 {
			u += mod
		}
		return u
	}
	for i := 3; i < B; i += 2 {
		res.pre[i>>1] = cal(i)
	}
	return res
}

// 要求: 0 <= x < mod 且 x 存在模逆元.
func (fi *FastInv) Inv(x int) int {
	if fi.mod == 1 {
		return 0
	}
	b, s, t := fi.mod, 1, 0
	n := bits.TrailingZeros(uint(x))
	x >>= n
	for x-b != 0 {
		if fi.isModPrime {
			if x < B {
				break
			}
		}

		m := bits.TrailingZeros(uint(x - b))
		f := x > b
		n += m
		var u int
		if f {
			x = (x - b) >> m
			u = s - t
			t = t << m
		} else {
			b = -(x - b) >> m
			u = s << m
			t = -(s - t)
		}
		s = u
	}

	return s * fi.pre[x>>1] % fi.mod * fi.ipow2[n] % fi.mod
}

func (fi *FastInv) _testPrime(v int) bool {
	if v == 1 {
		return false
	}
	for p := 3; p*p <= v; p += 2 {
		if v%p == 0 {
			return false
		}
	}
	return true
}
