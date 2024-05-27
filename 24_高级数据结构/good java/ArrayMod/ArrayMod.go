// 数组乘积取模/数组乘积逆元
// 区间乘积取模/区间乘积逆元
// 线性时间求逆元
// MultiMod, MultiModInv
// allInv, allModInv

package main

import "fmt"

func main() {
	arrayMod := NewArrayMod(10, 1000000007)
	arrayMod.Init(10, func(i int) int { return i + 1 })
	fmt.Println(arrayMod.PrefixMod(5))
	fmt.Println(arrayMod.PrefixModInv(5)*arrayMod.PrefixMod(5)%1000000007 == 1)
	fmt.Println(arrayMod.RangeMod(2, 5))
	fmt.Println(arrayMod.RangeModInv(2, 5)*arrayMod.RangeMod(2, 5)%1000000007 == 1)
}

type ArrayMod struct {
	mod     int
	fact    []int
	invFact []int
	one     int
	n       int
}

func NewArrayMod(n, mod int) *ArrayMod {
	return &ArrayMod{
		mod:     mod,
		fact:    make([]int, n),
		invFact: make([]int, n),
		one:     1 % mod,
	}
}

func (a *ArrayMod) Init(n int, f func(int) int) {
	a.n = n
	if n == 0 {
		return
	}
	mod, fact, invFact := a.mod, a.fact, a.invFact
	for i := 0; i < n; i++ {
		fact[i] = f(i)
		if fact[i] == 0 {
			fact[i] = 1
		}
		if i > 0 {
			fact[i] = (fact[i] * fact[i-1]) % mod
		}
	}
	invFact[n-1] = modInv(fact[n-1], mod)
	for i := n - 2; i >= 0; i-- {
		invFact[i] = (invFact[i+1] * f(i+1)) % mod
	}
}

func (a *ArrayMod) PrefixMod(right int) int {
	if right >= a.n {
		right = a.n - 1
	}
	if right < 0 {
		return a.one
	}
	return a.fact[right]
}

func (a *ArrayMod) PrefixModInv(right int) int {
	if right >= a.n {
		right = a.n - 1
	}
	if right < 0 {
		return a.one
	}
	return a.invFact[right]
}

func (a *ArrayMod) RangeMod(left, right int) int {
	if left < 0 {
		left = 0
	}
	if right >= a.n {
		right = a.n - 1
	}
	if left > right {
		return a.one
	}
	res := a.fact[right]
	if left > 0 {
		res = (res * a.invFact[left-1]) % a.mod
	}
	return res
}

func (a *ArrayMod) RangeModInv(left, right int) int {
	if left < 0 {
		left = 0
	}
	if right >= a.n {
		right = a.n - 1
	}
	if left > right {
		return a.one
	}
	res := a.invFact[right]
	if left > 0 {
		res = (res * a.fact[left-1]) % a.mod
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

func modInv(a, m int) int {
	g, x, _ := exgcd(a, m)
	if g != 1 && g != -1 {
		return -1
	}
	res := x % m
	if res < 0 {
		res += m
	}
	return res
}
