// https://nyaannyaan.github.io/library/set-function/subset-convolution.hpp

// 子集卷积:
// !H = F|G 且 |H| = |F| + |G|  (不相交的两个集合的并集)
// 为每个集合带上大小进行卷积运算 O(2^n*n^2)

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {
	// https://judge.yosupo.jp/problem/subset_convolution
	// 给定两个长度为2^n的序列A,B
	// !求 ck = ∑ai*bj (i&j==0,i|j==k)
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	A := make([]int, 1<<n)
	for i := 0; i < 1<<n; i++ {
		fmt.Fscan(in, &A[i])
	}
	B := make([]int, 1<<n)
	for i := 0; i < 1<<n; i++ {
		fmt.Fscan(in, &B[i])
	}

	sc := NewSubsetConvolution(n)
	res := sc.Multiply(A, B)
	for i := 0; i < 1<<n; i++ {
		fmt.Fprint(out, res[i], " ")
	}
}

const MOD int = 998244353

type SubsetConvolution struct {
	n  int
	pc []int // pc[i] = popcount(i)
}

func NewSubsetConvolution(n int) *SubsetConvolution {
	res := &SubsetConvolution{n: n, pc: make([]int, 1<<n)}
	for i := 1; i < (1 << n); i++ {
		res.pc[i] = res.pc[i-(i&-i)] + 1
	}
	return res
}

func (sc *SubsetConvolution) Multiply(a, b []int) []int {
	A := sc.lift(a)
	B := sc.lift(b)
	sc.zeta(A)
	sc.zeta(B)
	sc.prod(A, B)
	sc.mobius(A)
	return sc.unlift(A)
}

//  len(l) == len(r) == n + 1
func (sc *SubsetConvolution) add(l, r []int, d int) {
	for i := 0; i < d; i++ {
		l[i] = (l[i] + r[i]) % MOD
	}
}

//  len(l) == len(r) == n + 1
func (sc *SubsetConvolution) sub(l, r []int, d int) {
	for i := d; i <= sc.n; i++ {
		l[i] = (l[i] - r[i] + MOD) % MOD
	}
}

func (sc *SubsetConvolution) zeta(a [][]int) {
	n := len(a)
	for w := 1; w < n; w *= 2 {
		for k := 0; k < n; k += w * 2 {
			for i := 0; i < w; i++ {
				sc.add(a[k+w+i], a[k+i], sc.pc[k+w+i])
			}
		}
	}
}

func (sc *SubsetConvolution) mobius(a [][]int) {
	n := len(a)
	for w := n >> 1; w > 0; w >>= 1 {
		for k := 0; k < n; k += w * 2 {
			for i := 0; i < w; i++ {
				sc.sub(a[k+w+i], a[k+i], sc.pc[k+w+i])
			}
		}
	}
}

func (sc *SubsetConvolution) prod(a, b [][]int) {
	n := len(a)
	d := bits.TrailingZeros(uint(n))
	for i := 0; i < n; i++ {
		c := make([]int, sc.n+1)
		for j := 0; j <= d; j++ {
			for k := 0; k <= d-j; k++ {
				c[j+k] = (c[j+k] + a[i][j]*b[i][k]) % MOD
			}
		}
		a[i] = c
	}
}

// Wrap with popcount.
//  len(a) == 2^n
func (sc *SubsetConvolution) lift(a []int) [][]int {
	res := make([][]int, len(a))
	for i := 0; i < len(a); i++ {
		res[i] = make([]int, sc.n+1)
		res[i][sc.pc[i]] = a[i]
	}
	return res
}

// Unwrap popcount.
func (sc *SubsetConvolution) unlift(a [][]int) []int {
	res := make([]int, len(a))
	for i := 0; i < len(a); i++ {
		res[i] = a[i][sc.pc[i]]
	}
	return res
}
