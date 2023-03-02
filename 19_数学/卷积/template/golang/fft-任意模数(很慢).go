package main

import (
	"bufio"
	"fmt"
	"math"
	"math/bits"
	"os"
)

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

	conv := Convolution(poly1, poly2, 1e9+7)
	for i := 0; i < n+m-1; i++ {
		fmt.Fprint(out, conv[i], " ")
	}
}

func Convolution(A, B []int, mod int) []int {
	s := 10
	s2 := s << 1
	mask := (1 << s) - 1
	m0, m1, m2 := make([]int, len(A)), make([]int, len(A)), make([]int, len(A))
	n0, n1, n2 := make([]int, len(B)), make([]int, len(B)), make([]int, len(B))
	for i, v := range A {
		m0[i] = v & mask
		m1[i] = (v >> s) & mask
		m2[i] = v >> s2
	}
	for i, v := range B {
		n0[i] = v & mask
		n1[i] = (v >> s) & mask
		n2[i] = v >> s2
	}

	p_0, p1, pm1, pm2 := make([]int, len(A)), make([]int, len(A)), make([]int, len(A)), make([]int, len(A))
	for i := range p_0 {
		p_0[i] = m0[i] + m2[i]
		p1[i] = p_0[i] + m1[i]
		pm1[i] = p_0[i] - m1[i]
		pm2[i] = ((pm1[i] + m2[i]) << 1) - m0[i]
	}
	p0 := m0
	pinf := m2

	q_0, q1, qm1, qm2 := make([]int, len(B)), make([]int, len(B)), make([]int, len(B)), make([]int, len(B))
	for i := range q_0 {
		q_0[i] = n0[i] + n2[i]
		q1[i] = q_0[i] + n1[i]
		qm1[i] = q_0[i] - n1[i]
		qm2[i] = ((qm1[i] + n2[i]) << 1) - n0[i]
	}
	q0 := n0
	qinf := n2

	r0 := _convolution(p0, q0)
	r1 := _convolution(p1, q1)
	rm1 := _convolution(pm1, qm1)
	rm2 := _convolution(pm2, qm2)
	rinf := _convolution(pinf, qinf)

	r_0 := r0
	r_4 := rinf
	r_3, r_1, r_2 := make([]int, len(rm2)), make([]int, len(rm1)), make([]int, len(rm2))
	for i := range r_3 {
		r_3[i] = (rm2[i] - r1[i]) / 3
		r_1[i] = (r1[i] - rm1[i]) >> 1
		r_2[i] = rm1[i] - r0[i]
		r_3[i] = ((r_2[i] - r_3[i]) >> 1) + (rinf[i] << 1)
		r_2[i] += r_1[i] - r_4[i]
		r_1[i] -= r_3[i]
	}

	res := make([]int, len(r_4))
	for i := range res {
		res[i] = (((r_4[i]<<s2)+(r_3[i]<<s))%mod + r_2[i]) % mod
		res[i] = ((res[i] << s2) + (r_1[i] << s) + r_0[i]) % mod
	}
	return res
}

// 计算多个多项式的卷积
// 入参出参都是次项从低到高的系数
func MultiConvolution(coefs [][]int, mod int) []int {
	n := len(coefs)
	if n == 1 {
		return coefs[0]
	}
	return Convolution(MultiConvolution(coefs[:n/2], mod), MultiConvolution(coefs[n/2:], mod), mod)
}

// 计算 A(x) 和 B(x) 的卷积
//  c[i] = ∑a[k]*b[i-k], k=0..i
//  入参出参都是次项从低到高的系数
func _convolution(a, b []int) []int {
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
