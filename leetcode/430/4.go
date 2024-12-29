package main

var E *Enumeration

func init() {
	const SIZE int = 1e5 + 10
	const MOD int = 1e9 + 7
	E = NewEnumeration(SIZE, MOD)
}

type Enumeration struct {
	fac, ifac, inv []int
	mod            int
}

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

func (e *Enumeration) Fac(k int) int {
	e.expand(k)
	return e.fac[k]
}

func (e *Enumeration) Ifac(k int) int {
	e.expand(k)
	return e.ifac[k]
}

func (e *Enumeration) Inv(k int) int {
	e.expand(k)
	return e.inv[k]
}

func (e *Enumeration) C(n, k int) int {
	if n < 0 || k < 0 || n < k {
		return 0
	}
	return e.Fac(n) * e.Ifac(k) % e.mod * e.Ifac(n-k) % e.mod
}

func (e *Enumeration) P(n, k int) int {
	if n < 0 || k < 0 || n < k {
		return 0
	}
	return e.Fac(n) * e.Ifac(n-k) % e.mod
}

func (e *Enumeration) H(n, k int) int {
	if n == 0 {
		if k == 0 {
			return 1
		}
		return 0
	}
	return e.C(n+k-1, k)
}

func (e *Enumeration) Put(n, k int) int {
	return e.C(n+k-1, n)
}

func (e *Enumeration) Catalan(n int) int {
	return e.C(2*n, n) * e.Inv(n+1) % e.mod
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
		e.ifac[size] = Pow(e.fac[size], mod-2, mod)
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
	for exp > 0 {
		if exp&1 == 1 {
			res = res * base % mod
		}
		base = base * base % mod
		exp >>= 1
	}
	return res
}

func countGoodArrays(n, m, k int) int {
	const MOD = 1e9 + 7
	if k > n-1 || m == 0 {
		return 0
	}
	select_ := E.C(n-1, k)
	first := m % MOD
	diff := Pow(m-1, (n-1)-k, MOD)
	res := select_ % MOD
	res = res * first % MOD
	res = res * diff % MOD
	return res
}
