//  求组合数

package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int = 1e9 + 7

var E = NewEnumeration(1e5+10, MOD)

type Enumeration struct {
	fac, ifac, inv []int
	mod            int
}

// 模数为质数时的组合数计算.
func NewEnumeration(initSize, mod int) *Enumeration {
	res := &Enumeration{
		fac:  make([]int, 1, initSize+1),
		ifac: make([]int, 1, initSize+1),
		inv:  make([]int, 1, initSize+1),
		mod:  mod,
	}
	res.fac[0] = 1
	res.ifac[0] = 1
	res.inv[0] = 1
	res.expand(initSize)
	return res
}

// 阶乘.
func (e *Enumeration) Fac(k int) int {
	e.expand(k)
	return e.fac[k]
}

// 阶乘逆元.
func (e *Enumeration) Ifac(k int) int {
	e.expand(k)
	return e.ifac[k]
}

// 模逆元.
func (e *Enumeration) Inv(k int) int {
	e.expand(k)
	return e.inv[k]
}

// 组合数.
func (e *Enumeration) C(n, k int) int {
	if n < 0 || k < 0 || n < k {
		return 0
	}
	mod := e.mod
	return e.Fac(n) * e.Ifac(k) % mod * e.Ifac(n-k) % mod
}

// 排列数.
func (e *Enumeration) P(n, k int) int {
	if n < 0 || k < 0 || n < k {
		return 0
	}
	mod := e.mod
	return e.Fac(n) * e.Ifac(n-k) % mod
}

// 可重复选取元素的组合数.
func (e *Enumeration) H(n, k int) int {
	if n == 0 {
		if k == 0 {
			return 1
		}
		return 0
	}
	return e.C(n+k-1, k)
}

// n个相同的球放入k个不同的盒子(盒子可放任意个球)的方法数.
func (e *Enumeration) Put(n, k int) int {
	return e.C(n+k-1, n)
}

// 卡特兰数.
func (e *Enumeration) Catalan(n int) int {
	return e.C(2*n, n) * e.Inv(n+1) % e.mod
}

// lucas定理求解组合数.适合模数较小的情况.
func (e *Enumeration) Lucas(n, k int) int {
	if k == 0 {
		return 1
	}
	mod := e.mod
	return e.C(n%mod, k%mod) * e.Lucas(n/mod, k/mod) % mod
}

func (e *Enumeration) expand(size int) {
	if upper := e.mod - 1; size > upper {
		size = upper
	}
	if len(e.fac) < size+1 {
		mod := e.mod
		preSize := len(e.fac)
		diff := size + 1 - preSize
		e.fac = append(e.fac, make([]int, diff)...)
		e.ifac = append(e.ifac, make([]int, diff)...)
		e.inv = append(e.inv, make([]int, diff)...)
		for i := preSize; i < size+1; i++ {
			e.fac[i] = e.fac[i-1] * i % mod
		}
		e.ifac[size] = Pow(e.fac[size], mod-2, mod) // !modInv
		for i := size - 1; i >= preSize; i-- {
			e.ifac[i] = e.ifac[i+1] * (i + 1) % mod
		}
		for i := preSize; i < size+1; i++ {
			e.inv[i] = e.ifac[i] * e.fac[i-1] % mod
		}
	}
}

func Pow(base, exp, mod int) int {
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

func main() {
	// https://yukicoder.me/problems/no/117
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var s string
		fmt.Fscan(in, &s)
		op := s[0]
		inner := s[2 : len(s)-1]
		var n, k int
		fmt.Sscanf(inner, "%d,%d", &n, &k)
		switch op {
		case 'C':
			fmt.Fprintln(out, E.C(n, k))
		case 'P':
			fmt.Fprintln(out, E.P(n, k))
		case 'H':
			fmt.Fprintln(out, E.H(n, k))
		}
	}
}
