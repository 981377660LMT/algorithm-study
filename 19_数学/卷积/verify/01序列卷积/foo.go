// 给定两个长度为 n 的数组 ，你可以执行任意次以下操作
// 将A元素所有元素向左轮转一位
// 求能得到的最大的 sum(A[i] | B[i]) 的值
// !2<=n<=5e5,0<=ai,bi<=31

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func XorSum(nums1 []int, nums2 []int) int {
	n := len(nums1)
	tmp := nums2
	nums2 = make([]int, n)
	for i := 0; i < n; i++ {
		nums2[n-1-i] = tmp[i]
	}
	res := make([]int, n) // 每个轮转对应的和
	for bit := 0; bit < 5; bit++ {
		A := make([]int, n)
		B := make([]int, n)
		sum := 0
		for i := 0; i < n; i++ {
			A[i] = (nums1[i] >> bit) & 1
			B[i] = (nums2[i] >> bit) & 1
			if A[i] == 1 {
				sum += 1
			}
			if B[i] == 1 {
				sum += 1
			}
		}

		conv := Convolution(A, B)
		for i := n; i < n+n-1; i++ {
			conv[i%n] += conv[i]
		}
		for i := 0; i < n; i++ {
			res[i] += (sum - conv[i]) * (1 << bit)
		}
	}

	return maxs(res...)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	nums1, nums2 := make([]int, n), make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums1[i])
	}
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums2[i])
	}

	fmt.Fprintln(out, XorSum(nums1, nums2))
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
	base %= mod // ! 防止overflow
	res = 1
	for ; exp > 0; exp >>= 1 {
		if exp&1 == 1 {
			res = res * base % mod
		}
		base = base * base % mod
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

func mins(nums ...int) int {
	res := nums[0]
	for _, num := range nums {
		if num < res {
			res = num
		}
	}
	return res
}

func maxs(nums ...int) int {
	res := nums[0]
	for _, num := range nums {
		if num > res {
			res = num
		}
	}
	return res
}
