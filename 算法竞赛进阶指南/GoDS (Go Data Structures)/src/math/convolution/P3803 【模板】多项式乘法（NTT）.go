// !NTT卷积
// https://github.dev/EndlessCheng/codeforces-go/blob/3dd70515200872705893d52dc5dad174f2c3b5f3/copypasta/math_ntt.go#L350

// 给定一个n次多项式F(x)，和一个m次多项式G(x)。
// 请求出 F(x)和G(x)的卷积。
// n,m<=1e6
package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

// https://atcoder.jp/contests/practice2/tasks/practice2_f
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)

	poly1, poly2 := make(poly, n), make(poly, m)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &poly1[i])
	}
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &poly2[i])
	}

	res := poly1.Convolution(poly2)
	for i := 0; i < n+m-1; i++ {
		fmt.Fprint(out, res[i], " ")
	}
}

type poly []int64

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

const MOD = 998244353

var omega, omegaInv [31]int64 // 多开一点空间

func init() {
	const g, invG = 3, 332748118
	for i := 1; i < len(omega); i++ {
		omega[i] = qpow(g, (MOD-1)/(1<<i), MOD)
		omegaInv[i] = qpow(invG, (MOD-1)/(1<<i), MOD)
	}
}

type ntt struct {
	n    int
	invN int64
}

func newNTT(n int) ntt {
	return ntt{n, qpow(int64(n), MOD-2, MOD)}
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
			for i, w := 0, int64(1); i < m; i++ {
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

func qpow(base int64, exp int, mod int64) (res int64) {
	res = 1
	for ; exp > 0; exp >>= 1 {
		if exp&1 == 1 {
			res = res * base % mod
		}
		base = base * base % mod
	}
	return
}

// !多项式api

func (a poly) Reverse() poly {
	for i, n := 0, len(a); i < n/2; i++ {
		a[i], a[n-1-i] = a[n-1-i], a[i]
	}
	return a
}

func (a poly) ReverseCopy() poly {
	n := len(a)
	b := make(poly, n)
	for i, v := range a {
		b[n-1-i] = v
	}
	return b
}

func (a poly) Neg() poly {
	b := make(poly, len(a))
	for i, v := range a {
		if v > 0 {
			b[i] = MOD - v
		}
	}
	return b
}

// 多项式相加
func (a poly) Add(b poly) poly {
	c := make(poly, len(a))
	for i, v := range a {
		c[i] = (v + b[i]) % MOD
	}
	return c
}

// 多项式相减
func (a poly) Sub(b poly) poly {
	c := make(poly, len(a))
	for i, v := range a {
		c[i] = (v - b[i] + MOD) % MOD
	}
	return c
}

// 多项式相乘
func (a poly) Mul(k int64) poly {
	k %= MOD
	b := make(poly, len(a))
	for i, v := range a {
		b[i] = v * k % MOD
	}
	return b
}

// 多项式除法
// https://blog.orzsiyuan.com/archives/Polynomial-Division-and-Modulo/
// https://oi-wiki.org/math/poly/Div-mod/
// 模板题 https://www.luogu.com.cn/problem/P4512
func (a poly) Div(b poly) poly {
	k := len(a) - len(b) + 1
	if k <= 0 {
		return make(poly, 1)
	}
	A := a.ReverseCopy().resize(k)
	B := b.ReverseCopy().resize(k)
	return A.Convolution(B.Inv())[:k].Reverse()
}

func (a poly) Lsh(k int) poly {
	b := make(poly, len(a))
	if k > len(a) {
		return b
	}
	copy(b[k:], a)
	return b
}

func (a poly) Rsh(k int) poly {
	b := make(poly, len(a))
	if k > len(a) {
		return b
	}
	copy(b, a[k:])
	return b
}

// 多项式求导
func (a poly) Derivative() poly {
	n := len(a)
	d := make(poly, n)
	for i := 1; i < n; i++ {
		d[i-1] = a[i] * int64(i) % MOD
	}
	return d
}

// 多项式积分
func (a poly) Integral() poly {
	n := len(a)
	s := make(poly, n)
	s[0] = 0 // C
	// 线性求逆元，详见 math.go 中的 initAllInv
	inv := make([]int64, n)
	inv[1] = 1
	for i := 2; i < n; i++ {
		inv[i] = int64(MOD-MOD/i) * inv[MOD%i] % MOD
	}
	for i := 1; i < n; i++ {
		s[i] = a[i-1] * inv[i] % MOD
	}
	return s
}

// 多项式乘法逆 (mod x^n, 下同)
// 参考 https://blog.orzsiyuan.com/archives/Polynomial-Inversion/
// https://oi-wiki.org/math/poly/Inv/
// 模板题 https://www.luogu.com.cn/problem/P4238
func (a poly) Inv() poly {
	n := len(a)
	m := 1 << bits.Len(uint(n))
	A := a.resize(m)
	invA := make(poly, m)
	invA[0] = qpow(A[0], MOD-2, MOD)
	for l := 2; l <= m; l <<= 1 {
		ll := l << 1
		b := A[:l].resize(ll)
		iv := invA[:l].resize(ll)
		t := newNTT(ll)
		t.dft(b)
		t.dft(iv)
		for i, v := range iv {
			b[i] = v * (2 - v*b[i]%MOD + MOD) % MOD
		}
		t.idft(b)
		copy(invA, b[:l])
	}
	return invA[:n]
}

// 多项式取模
func (a poly) Mod(b poly) poly {
	m := len(b)
	return a[:m-1].Sub(a.Div(b).Convolution(b)[:m-1])
}

func (a poly) DivMod(b poly) (div, mod poly) {
	m := len(b)
	div = a.Div(b)
	mod = a[:m-1].Sub(div.Convolution(b)[:m-1])
	return
}
