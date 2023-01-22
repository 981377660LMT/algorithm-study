package main

import (
	"fmt"
	"math/big"
	"math/bits"
)

const MOD int = 998244353

// FFT求大数乘法(高精度乘法/字符串相乘)
// github.com/EndlessCheng/codeforces-go
func main() {
	// '2'*2e5 * '3'*2e5
	nums1 := make([]byte, 2e5)
	for i := range nums1 {
		nums1[i] = '2'
	}
	nums2 := make([]byte, 2e5)
	for i := range nums2 {
		nums2[i] = '3'
	}

	res := string(BigMul(nums1, nums2))
	bigInt := new(big.Int)
	bigInt.SetString(res, 10)
	fmt.Println(bigInt.Mod(bigInt, big.NewInt(int64(MOD))))
}

// 大数乘法
func BigMul(s1, s2 []byte) []rune {
	n, m := len(s1), len(s2)
	limit := 1 << uint(bits.Len(uint(n+m-1)))
	f := newNTT(limit)
	A := make([]int, limit)
	for i, v := range s1 {
		A[n-1-i] = int(v & 15)
	}
	B := make([]int, limit)
	for i, v := range s2 {
		B[m-1-i] = int(v & 15)
	}
	f.dft(A)
	f.dft(B)
	for i, v := range A {
		A[i] = v * B[i] % MOD
	}
	f.idft(A)

	res := make([]rune, n+m)
	for i := 0; i < n+m-1; i++ {
		res[i] += rune(A[i])
		res[i+1] = res[i] / 10
		res[i] %= 10
	}
	r := n + m
	for res[r-1] == 0 {
		r--
	}
	res = res[:r]
	for i := 0; i < r/2; i++ {
		res[i], res[r-1-i] = res[r-1-i]+'0', res[i]+'0'
	}
	if r&1 > 0 {
		res[r/2] += '0'
	}
	return res
}

type poly []int

// 计算 A(x) 和 B(x) 的卷积 (convolution)
//  c[i] = ∑a[k]*b[i-k], k=0..i
//  入参出参都是次项从低到高的系数
// https://github.dev/EndlessCheng/codeforces-go/blob/3dd70515200872705893d52dc5dad174f2c3b5f3/copypasta/math_ntt.go#L350
// 模板题 https://www.luogu.com.cn/problem/P3803 https://www.luogu.com.cn/problem/P1919 https://atcoder.jp/contests/practice2/tasks/practice2_f
func (a poly) Convolution(b poly) poly {
	n, m := len(a), len(b)
	limit := 1 << bits.Len(uint(n+m-1))
	A := a.resize(limit)
	B := b.resize(limit)
	t := newNTT(limit)
	t.dft(A)
	t.dft(B)
	for i, v := range A {
		A[i] = v * B[i] % MOD
	}
	t.idft(A)
	return A[:n+m-1]
}

// 计算多个多项式的卷积
// 入参出参都是次项从低到高的系数
func MultiConvolution(coefs []poly) poly {
	n := len(coefs)
	if n == 1 {
		return coefs[0]
	}
	return MultiConvolution(coefs[:n/2]).Convolution(MultiConvolution(coefs[n/2:]))
}

func (a poly) resize(n int) poly {
	b := make(poly, n)
	copy(b, a)
	return b
}

var omega, omegaInv [31]int // 多开一点空间

func init() {
	// 常用素数及原根 http://blog.miskcoo.com/2014/07/fft-prime-table
	const g, invG = 3, 332748118
	for i := 1; i < len(omega); i++ {
		omega[i] = qpow(g, (MOD-1)/(1<<i), MOD)
		omegaInv[i] = qpow(invG, (MOD-1)/(1<<i), MOD)
	}
}

type ntt struct {
	n    int
	invN int
}

func newNTT(n int) ntt {
	return ntt{n, qpow(int(n), MOD-2, MOD)}
}

func (t ntt) transform(a, omega poly) {
	for i, j := 0, 0; i < t.n; i++ {
		if i > j {
			a[i], a[j] = a[j], a[i]
		}
		for l := t.n >> 1; ; l >>= 1 {
			j ^= l
			if j >= l {
				break
			}
		}
	}

	for l, li := 2, 1; l <= t.n; l <<= 1 {
		m := l >> 1
		wn := omega[li]
		li++
		for st := 0; st < t.n; st += l {
			b := a[st:]
			for i, w := 0, int(1); i < m; i++ {
				d := b[m+i] * w % MOD
				b[m+i] = (b[i] - d + MOD) % MOD
				b[i] = (b[i] + d) % MOD
				w = w * wn % MOD
			}
		}
	}
}

func (t ntt) dft(p poly) {
	t.transform(p, omega[:])
}

func (t ntt) idft(p poly) {
	t.transform(p, omegaInv[:])
	for i, v := range p {
		p[i] = v * t.invN % MOD
	}
}

func qpow(base int, exp int, mod int) (res int) {
	res = 1
	for ; exp > 0; exp >>= 1 {
		if exp&1 == 1 {
			res = res * base % mod
		}
		base = base * base % mod
	}
	return
}
