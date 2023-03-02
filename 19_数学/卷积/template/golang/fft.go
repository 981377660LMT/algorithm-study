// fft 快速傅里叶变换求卷积
// 浮点数运算产生的误差会比较大 一般用ntt

package main

import (
	"bufio"
	"fmt"
	"math"
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

	poly1 := make([]int, n) // 从低到高表示F(x)的系数
	poly2 := make([]int, m) // 从低到高表示G(x)的系数
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &poly1[i])
	}
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &poly2[i])
	}

	conv := Convolution(poly1, poly2)
	for i := 0; i < n+m-1; i++ {
		fmt.Fprint(out, conv[i], " ")
	}
}

// 计算 A(x) 和 B(x) 的卷积
//  c[i] = ∑a[k]*b[i-k], k=0..i
//  入参出参都是次项从低到高的系数
func Convolution(a, b []int) []int {
	n, m := len(a), len(b)
	limit := 1 << uint(bits.Len(uint(n+m-1)))
	f := newFFT(limit)
	cmplxA := make([]complex128, limit)
	for i, v := range a {
		cmplxA[i] = complex(float64(v), 0)
	}
	cmplxB := make([]complex128, limit)
	for i, v := range b {
		cmplxB[i] = complex(float64(v), 0)
	}
	f.dft(cmplxA)
	f.dft(cmplxB)
	for i := range cmplxA {
		cmplxA[i] *= cmplxB[i]
	}
	f.idft(cmplxA)
	conv := make([]int, n+m-1)
	for i := range conv {
		conv[i] = int(math.Round(real(cmplxA[i])))
	}
	return conv
}

// 计算多个多项式的卷积
// 入参出参都是次项从低到高的系数
func MultiConvolution(coefs [][]int) []int {
	n := len(coefs)
	if n == 1 {
		return coefs[0]
	}
	return Convolution(MultiConvolution(coefs[:n/2]), MultiConvolution(coefs[n/2:]))
}

// https://github.dev/EndlessCheng/codeforces-go/tree/master/copypasta
type fft struct {
	n               int
	omega, omegaInv []complex128
}

func newFFT(n int) *fft {
	omega := make([]complex128, n)
	omegaInv := make([]complex128, n)
	for i := range omega {
		sin, cos := math.Sincos(2 * math.Pi * float64(i) / float64(n))
		omega[i] = complex(cos, sin)
		omegaInv[i] = complex(cos, -sin)
	}
	return &fft{n, omega, omegaInv}
}

func (f *fft) transform(a, omega []complex128) {
	for i, j := 0, 0; i < f.n; i++ {
		if i > j {
			a[i], a[j] = a[j], a[i]
		}
		for l := f.n >> 1; ; l >>= 1 {
			j ^= l
			if j >= l {
				break
			}
		}
	}
	for l := 2; l <= f.n; l <<= 1 {
		m := l >> 1
		for st := 0; st < f.n; st += l {
			p := a[st:]
			for i := 0; i < m; i++ {
				t := omega[f.n/l*i] * p[m+i]
				p[m+i] = p[i] - t
				p[i] += t
			}
		}
	}
}

func (f *fft) dft(a []complex128) {
	f.transform(a, f.omega)
}

func (f *fft) idft(a []complex128) {
	f.transform(a, f.omegaInv)
	for i := range a {
		a[i] /= complex(float64(f.n), 0)
	}
}
