//  求组合数

package main

var E *Enumeration

func init() {
	const SIZE int = 1e6 + 10
	const MOD int = 1e9 + 7
	E = NewEnumeration(SIZE, MOD)
}

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

func (e *Enumeration) expand(size int) {
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
	res := 1
	for ; exp > 0; exp >>= 1 {
		if exp&1 == 1 {
			res = res * base % mod
		}
		base = base * base % mod
	}
	return res
}

// if __name__ == "__main__":
//     # https://yukicoder.me/problems/no/117
//     import sys

//     sys.setrecursionlimit(int(1e9))
//     input = lambda: sys.stdin.readline().rstrip("\r\n")
//     T = int(input())
//     C = Combination(2 * 10**6 + 10)
//     for _ in range(T):
//         s = input()
//         op = s[0]
//         inner = s[2:-1]
//         n, k = map(int, inner.split(","))
//         if op == "C":
//             print(C(n, k))
//         elif op == "P":
//             print(C.P(n, k))
//         elif op == "H":
//             print(C.H(n, k))
