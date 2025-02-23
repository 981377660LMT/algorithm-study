// Lucas 定理用于求解大组合数取模的问题，但模数必须是质数。
// !扩展 Lucas 定理用于求解大组合数取模的问题，模数不一定要求是质数。
//
// P4720 【模板】扩展卢卡斯定理(EXLucas)
// https://www.luogu.com.cn/problem/P4720
// 1≤k≤n≤1e18 ，2≤mod≤1e6 ，不保证 mod 是质数。
//
// https://judge.yosupo.jp/problem/binomial_coefficient

package main

import (
	"bufio"
	"fmt"
	"os"
)

// https://leetcode.cn/problems/check-if-digits-are-equal-in-string-after-operations-ii/description/
func hasSameDigits(s string) bool {
	MOD := 10
	C := NewBinomialCoefficient(MOD)
	n := len(s)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		nums[i] = int(s[i] - '0')
	}

	sum1 := 0
	for i := 0; i < n-1; i++ {
		sum1 = (sum1 + nums[i]*C.C(n-2, i)) % MOD
	}

	sum2 := 0
	for i := 0; i < n-1; i++ {
		sum2 = (sum2 + nums[i+1]*C.C(n-2, i)) % MOD
	}

	return sum1 == sum2
}

// https://judge.yosupo.jp/problem/binomial_coefficient
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T, MOD int
	fmt.Fscan(in, &T, &MOD)
	bc := NewBinomialCoefficient(MOD)
	for i := 0; i < T; i++ {
		var N, K int
		fmt.Fscan(in, &N, &K)
		fmt.Fprintln(out, bc.C(N, K))
	}
}

type factor struct {
	p  int
	pe int // p 的幂，比如 p^e
}

type BinomialCoefficient struct {
	mod           int
	factorization []factor
	facs          [][]int
	invs          [][]int
	coeffs        []int
	pows          [][]int
}

func NewBinomialCoefficient(mod int) *BinomialCoefficient {
	res := &BinomialCoefficient{
		mod: mod,
	}
	res.factorization = res.factorize(mod)
	nFac := len(res.factorization)
	res.facs = make([][]int, nFac)
	res.invs = make([][]int, nFac)
	res.coeffs = make([]int, nFac)
	res.pows = make([][]int, nFac)

	for i, fac := range res.factorization {
		p := fac.p
		pe := fac.pe

		facArr := make([]int, pe)
		for j := 0; j < pe; j++ {
			facArr[j] = 1
		}
		for j := 1; j < pe; j++ {
			mult := 1
			if j%p != 0 {
				mult = j
			}
			facArr[j] = (facArr[j-1] * mult) % pe
		}

		invArr := make([]int, pe)
		invArr[pe-1] = facArr[pe-1]
		for j := pe - 1; j >= 1; j-- {
			mult := 1
			if j%p != 0 {
				mult = j
			}
			invArr[j-1] = (invArr[j] * mult) % pe
		}

		res.facs[i] = facArr
		res.invs[i] = invArr

		c := modinv(mod/pe, pe)
		res.coeffs[i] = (mod / pe * c) % mod
		powp := []int{1}
		for powp[len(powp)-1]*p != pe {
			powp = append(powp, powp[len(powp)-1]*p)
		}
		res.pows[i] = powp
	}

	return res
}

func (bc *BinomialCoefficient) C(n, k int) int {
	if k < 0 || k > n {
		return 0
	}
	if k == 0 || k == n {
		return 1 % bc.mod
	}
	res := 0
	for i, fac := range bc.factorization {
		p := fac.p
		pe := fac.pe
		cp := bc.choosePe(n, k, p, pe, bc.facs[i], bc.invs[i], bc.pows[i])
		res = (res + cp*bc.coeffs[i]) % bc.mod
	}
	return res
}

func (bc *BinomialCoefficient) e(n, k, r, p int) int {
	res := 0
	for n > 0 {
		n /= p
		k /= p
		r /= p
		res += n - k - r
	}
	return res
}

func (bc *BinomialCoefficient) choosePe(n, k, p, pe int, fac, inv, powp []int) int {
	r := n - k
	e0 := bc.e(n, k, r, p)
	if e0 >= len(powp) {
		return 0
	}
	res := powp[e0]
	nDiv := n / (pe / p)
	kDiv := k / (pe / p)
	rDiv := r / (pe / p)
	if (p != 2 || pe == 4) && (bc.e(nDiv, kDiv, rDiv, p)%2 != 0) {
		res = pe - res
	}
	for n > 0 {
		res = (res * fac[n%pe]) % pe
		res = (res * inv[k%pe]) % pe
		res = (res * inv[r%pe]) % pe
		n /= p
		k /= p
		r /= p
	}
	return res
}

func (bc *BinomialCoefficient) factorize(N int) []factor {
	factors := []factor{}
	for i := 2; i*i <= N; i++ {
		if N%i != 0 {
			continue
		}
		cnt := 0
		for N%i == 0 {
			N /= i
			cnt++
		}
		power := 1
		for j := 0; j < cnt; j++ {
			power *= i
		}
		factors = append(factors, factor{p: i, pe: power})
	}
	if N != 1 {
		factors = append(factors, factor{p: N, pe: N})
	}
	return factors
}

func modinv(a, MOD int) int {
	r0, r1 := a, MOD
	s0, s1 := 1, 0
	for r1 != 0 {
		q := r0 / r1
		r0, r1 = r1, r0%r1
		s0, s1 = s1, s0-q*s1
	}
	return (s0%MOD + MOD) % MOD
}
