package main

import "fmt"

const MOD int = 1e9 + 7
const N int = 2e5 + 10

var fac, ifac [N + 1]int

func init() {
	fac[0] = 1
	ifac[0] = 1
	for i := 1; i < N+1; i++ {
		fac[i] = fac[i-1] * i % MOD
	}
	ifac[N] = Pow(fac[N], MOD-2, MOD)
	for i := N; i > 0; i-- {
		ifac[i-1] = ifac[i] * i % MOD
	}
}

func C(n, k int) int {
	if n < 0 || k < 0 || n < k {
		return 0
	}
	return fac[n] * ifac[k] % MOD * ifac[n-k] % MOD
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

func main() {
	fmt.Println(C(500, 300))
}
