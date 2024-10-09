// TODO

package main

import (
	"bufio"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"sort"
)

func main() {
	abc212_g()
}

func demo() {
	arr := NewArrayOnDivisors(9)
	arr.SetEulerPhi()
	arr.Enumerate(func(d, fd int) {
		fmt.Println(d, fd)
	})
}

// https://atcoder.jp/contests/abc212/tasks/abc212_g
func abc212_g() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	const MOD int = 998244353

	var p int
	fmt.Fscan(in, &p)

	A := NewArrayOnDivisors(p - 1)
	A.SetEulerPhi()
	res := 1
	A.Enumerate(func(d, fd int) { res = (res + (d%MOD)*(fd%MOD)) % MOD })
	fmt.Fprintln(out, res%MOD)
}

// https://atcoder.jp/contests/abc349/tasks/abc349_f
func abc349_f() {}

type ArrayOnDivisors struct {
	divs []int
	data []int
	pfs  [][2]int
	mp   map[int]int
}

func NewArrayOnDivisors(n int) *ArrayOnDivisors {
	return NewArrayOnDivisorsFromPfs(PfsSorted(n))
}

func NewArrayOnDivisorsFromPfs(pfs [][2]int) *ArrayOnDivisors {
	res := &ArrayOnDivisors{}
	res.build(pfs)
	return res
}

func (aod *ArrayOnDivisors) SetMultiplicative(f func(p, k int) int) {
	aod.data = aod.data[:0]
	aod.data = append(aod.data, 1)
	n := len(aod.divs)
	for _, pe := range aod.pfs {
		p, e := pe[0], pe[1]
		for k := 1; k <= e; k++ {
			for i := 0; i < n; i++ {
				aod.data = append(aod.data, aod.data[i]*f(p, k))
			}
		}
	}
}

func (aod *ArrayOnDivisors) SetEulerPhi() {
	aod.data = append(aod.divs[:0:0], aod.divs...)
	aod.DivisorMobius()
}

func (aod *ArrayOnDivisors) SetMobius() {
	aod.SetMultiplicative(func(_, k int) int {
		if k == 1 {
			return -1
		}
		return 0
	})
}

// 倍数前缀和.
func (aod *ArrayOnDivisors) MultiplierZeta() {
	k := 1
	n := len(aod.divs)
	for _, pe := range aod.pfs {
		mod := k * (pe[1] + 1)
		for i := 0; i < n/mod; i++ {
			for j := mod - k - 1; j >= 0; j-- {
				aod.data[mod*i+j] += aod.data[mod*i+j+k]
			}
		}
		k *= pe[1] + 1
	}
}

// 倍数前缀和.
func (aod *ArrayOnDivisors) MultiplierZetaFunc(add func(a, b int) int) {
	k := 1
	n := len(aod.divs)
	for _, pe := range aod.pfs {
		mod := k * (pe[1] + 1)
		for i := 0; i < n/mod; i++ {
			for j := mod - k - 1; j >= 0; j-- {
				aod.data[mod*i+j] = add(aod.data[mod*i+j], aod.data[mod*i+j+k])
			}
		}
		k *= pe[1] + 1
	}
}

// 倍数前缀和逆操作。
func (aod *ArrayOnDivisors) MultiplierMobius() {
	k := 1
	n := len(aod.divs)
	for _, pe := range aod.pfs {
		mod := k * (pe[1] + 1)
		for i := 0; i < n/mod; i++ {
			for j := 0; j < mod-k; j++ {
				aod.data[mod*i+j] -= aod.data[mod*i+j+k]
			}
		}
		k *= pe[1] + 1
	}
}

// 倍数前缀和逆操作。
func (aod *ArrayOnDivisors) MultiplierMobiusFunc(sub func(a, b int) int) {
	k := 1
	n := len(aod.divs)
	for _, pe := range aod.pfs {
		mod := k * (pe[1] + 1)
		for i := 0; i < n/mod; i++ {
			for j := 0; j < mod-k; j++ {
				aod.data[mod*i+j] = sub(aod.data[mod*i+j], aod.data[mod*i+j+k])
			}
		}
		k *= pe[1] + 1
	}
}

// 约数前缀和。
func (aod *ArrayOnDivisors) DivisorZeta() {
	k := 1
	n := len(aod.divs)
	for _, pe := range aod.pfs {
		mod := k * (pe[1] + 1)
		for i := 0; i < n/mod; i++ {
			for j := 0; j < mod-k; j++ {
				aod.data[mod*i+j+k] += aod.data[mod*i+j]
			}
		}
		k *= pe[1] + 1
	}
}

// 约数前缀和。
func (aod *ArrayOnDivisors) DivisorZetaFunc(add func(a, b int) int) {
	k := 1
	n := len(aod.divs)
	for _, pe := range aod.pfs {
		mod := k * (pe[1] + 1)
		for i := 0; i < n/mod; i++ {
			for j := 0; j < mod-k; j++ {
				aod.data[mod*i+j+k] = add(aod.data[mod*i+j+k], aod.data[mod*i+j])
			}
		}
		k *= pe[1] + 1
	}
}

// 约数前缀和逆变换。
func (aod *ArrayOnDivisors) DivisorMobius() {
	k := 1
	n := len(aod.divs)
	for _, pe := range aod.pfs {
		mod := k * (pe[1] + 1)
		for i := 0; i < n/mod; i++ {
			for j := mod - k - 1; j >= 0; j-- {
				aod.data[mod*i+j+k] -= aod.data[mod*i+j]
			}
		}
		k *= pe[1] + 1
	}
}

// 约数前缀和逆变换。
func (aod *ArrayOnDivisors) DivisorMobiusFunc(sub func(a, b int) int) {
	k := 1
	n := len(aod.divs)
	for _, pe := range aod.pfs {
		mod := k * (pe[1] + 1)
		for i := 0; i < n/mod; i++ {
			for j := mod - k - 1; j >= 0; j-- {
				aod.data[mod*i+j+k] = sub(aod.data[mod*i+j+k], aod.data[mod*i+j])
			}
		}
		k *= pe[1] + 1
	}
}

func (aod *ArrayOnDivisors) Get(d int) int {
	return aod.data[aod.mp[d]]
}

func (aod *ArrayOnDivisors) Set(d int, f func(int) int) {
	pre := &aod.data[aod.mp[d]]
	*pre = f(*pre)
}

// (d, fd)
func (aod *ArrayOnDivisors) Enumerate(f func(d, fd int)) {
	for i, d := range aod.divs {
		f(d, aod.data[i])
	}
}

func (aod *ArrayOnDivisors) build(pfs [][2]int) {
	n := 1
	for _, pe := range pfs {
		n *= pe[1] + 1
	}
	divs := make([]int, n)
	for i := range divs {
		divs[i] = 1
	}
	data := make([]int, n)
	nxt := 1
	for _, pe := range pfs {
		l := nxt
		q := pe[0]
		for i := 0; i < pe[1]; i++ {
			for k := 0; k < l; k++ {
				divs[nxt] = divs[k] * q
				nxt++
			}
			q *= pe[0]
		}
	}
	mp := make(map[int]int, n)
	for i := 0; i < n; i++ {
		mp[divs[i]] = i
	}
	aod.divs = divs
	aod.data = data
	aod.pfs = pfs
	aod.mp = mp
}

func PfsSorted(n int) [][2]int {
	res := [][2]int{}
	for n > 1 {
		p := pollardRhoPrimeFactor(n)
		e := 0
		for n%p == 0 {
			n /= p
			e++
		}
		res = append(res, [2]int{p, e})
	}
	sort.Slice(res, func(i, j int) bool { return res[i][0] < res[j][0] })
	return res
}

func pollardRhoPrimeFactor(n int) int {
	if n == 4 {
		return 2
	}
	if isPrimeMillerRabin(n) {
		return n
	}

	cache := map[int]int{}
	var getPrimeFactor func(int) int
	getPrimeFactor = func(x int) (res int) {
		if cache[x] > 0 {
			return cache[x]
		}
		p := pollardRhoFactor(x)
		if p == x {
			cache[x] = x
			return p
		}
		res = max(getPrimeFactor(p), getPrimeFactor(x/p))
		cache[x] = res
		return
	}

	return getPrimeFactor(n)
}

func isPrimeMillerRabin(n int) bool {
	return big.NewInt(int64(n)).ProbablyPrime(0)
}

// Pollard-Rho 算法求出一个因子 O(n^1/4)
func pollardRhoFactor(n int) int {
	if n == 4 {
		return 2
	}
	if isPrimeMillerRabin(n) {
		return n
	}

	for {
		c := 1 + rand.Intn(n-1)
		f := func(x int) int { return (mul(x, x, n) + c) % n }
		for t, r := f(0), f(f(0)); t != r; t, r = f(t), f(f(r)) {
			if d := gcd(abs(t-r), n); d > 1 {
				return d
			}
		}
	}
}

func gcd(a, b int) int {
	for a != 0 {
		a, b = b%a, a
	}
	return b
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func mul(a, b int, n int) (res int) {
	for ; b > 0; b >>= 1 {
		if b&1 == 1 {
			res = (res + a) % n
		}
		a = (a + a) % n
	}
	return
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func max32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}
