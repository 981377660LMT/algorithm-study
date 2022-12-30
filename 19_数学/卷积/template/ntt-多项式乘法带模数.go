// !NTT 解决的是多项式乘法带模数的情况，可以说有些受模数的限制，
// 数也比较大。目前最常见的模数是 998244353。

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)

	poly1, poly2 := make([]int, n+1), make([]int, m+1)
	for i := 0; i < n+1; i++ {
		fmt.Fscan(in, &poly1[i])
	}
	for i := 0; i < m+1; i++ {
		fmt.Fscan(in, &poly2[i])
	}

	res := Convolution(poly1, poly2)
	for i := 0; i < n+m+1; i++ {
		fmt.Fprint(out, res[i], " ")
	}
}

const MOD = 998244353

// https://github.dev/EndlessCheng/codeforces-go/tree/master/copypasta
// 计算 A(x) 和 B(x) 的卷积 (convolution)
// c[i] = ∑a[k]*b[i-k], k=0..i
// 入参出参都是次项从低到高的系数
// 模板题 https://www.luogu.com.cn/problem/P3803 https://www.luogu.com.cn/problem/P1919 https://atcoder.jp/contests/practice2/tasks/practice2_f
func Convolution(a, b []int) []int {
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

func resize(a []int, n int) []int {
	b := make([]int, n)
	copy(b, a)
	return b
}

var omega, omegaInv [31]int // 多开一点空间

func init() {
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

func (t ntt) transform(a, omega []int) {
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

func (t ntt) dft(a []int) {
	t.transform(a, omega[:])
}

func (t ntt) idft(a []int) {
	t.transform(a, omegaInv[:])
	for i, v := range a {
		a[i] = v * t.invN % MOD
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
