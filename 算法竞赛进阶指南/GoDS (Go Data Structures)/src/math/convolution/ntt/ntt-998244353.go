// https://github.dev/EndlessCheng/codeforces-go/blob/3dd70515200872705893d52dc5dad174f2c3b5f3/copypasta/math_ntt.go#L350
// 模板题 https://www.luogu.com.cn/problem/P3803 https://www.luogu.com.cn/problem/P1919 https://atcoder.jp/contests/practice2/tasks/practice2_f
// !NTT卷积 受模数的限制 一般限定模数为998244353

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
	// n++
	// m++

	poly1, poly2 := make([]int, n), make([]int, m)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &poly1[i])
	}
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &poly2[i])
	}

	res := Convolution(poly1, poly2)
	for i := 0; i < n+m-1; i++ {
		fmt.Fprint(out, res[i], " ")
	}
}

const MOD = 998244353

type poly = []int

// 计算 A(x) 和 B(x) 的卷积 (convolution)
//  c[i] = ∑a[k]*b[i-k], k=0..i
//  入参出参都是次项从低到高的系数
func Convolution(a, b poly) poly {
	n, m := len(a), len(b)
	limit := 1 << bits.Len(uint(n+m-1))
	A := resize(a, limit)
	B := resize(b, limit)
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
	return Convolution(MultiConvolution(coefs[:n/2]), MultiConvolution(coefs[n/2:]))
}

func resize(a poly, n int) poly {
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
			for i, w := 0, 1; i < m; i++ {
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
	base %= mod
	res = 1
	for ; exp > 0; exp >>= 1 {
		if exp&1 == 1 {
			res = res * base % mod
		}
		base = base * base % mod
	}
	return
}
