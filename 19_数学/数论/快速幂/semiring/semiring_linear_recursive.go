// 半环上的线性递推
// O(n^2logk)
// ref: https://nyaannyaan.github.io/library/math/semiring-linear-recursive.hpp

package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
)

// Semiring
type AddFunc[T any] func(a, b T) T
type MulFunc[T any] func(a, b T) T
type IdFunc[T any] func() T

type LinearRecursive[T any] struct {
	C   []T        // 递推系数 c_1, c_2, ...
	N   int        // 递推阶数
	Add AddFunc[T] // 半环加法
	Mul MulFunc[T] // 半环乘法
	Id0 IdFunc[T]  // 加法单位元
	Id1 IdFunc[T]  // 乘法单位元
}

func NewLinearRecursive[T any](
	c []T,
	add AddFunc[T], mul MulFunc[T], id0, id1 IdFunc[T],
) *LinearRecursive[T] {
	cc := slices.Clone(c)
	return &LinearRecursive[T]{
		C:   cc,
		N:   len(cc),
		Add: add,
		Mul: mul,
		Id0: id0,
		Id1: id1,
	}
}

// KthTerm 计算递推关系的第 k 项
// a 是初始项 a_0, a_1, ..., a_{N-1}
func (lr *LinearRecursive[T]) KthTerm(a []T, k int) T {
	if len(a) != lr.N {
		panic("len(a) != len(C)")
	}
	if k < len(a) {
		return a[k]
	}
	coeff := lr.getCoeff(k)
	res := lr.Id0()
	for i := range lr.N {
		res = lr.Add(res, lr.Mul(a[i], coeff[lr.N-1-i]))
	}
	return res
}

func (lr *LinearRecursive[T]) getCoeff(k int) []T {
	if k == 0 {
		b := make([]T, lr.N)
		for i := range b {
			b[i] = lr.Id0()
		}
		b[lr.N-1] = lr.Id1()
		return b
	}
	half := lr.getCoeff(k / 2)
	half = lr.doubling(half)
	if (k & 1) == 1 {
		half = lr.increment(half)
	}
	return half
}

func (lr *LinearRecursive[T]) increment(b []T) []T {
	v := make([]T, lr.N)
	for i := 0; i+1 < lr.N; i++ {
		v[i] = b[i+1]
	}
	v[lr.N-1] = lr.Id0()
	t := b[0]
	for i := range lr.N {
		v[i] = lr.Add(v[i], lr.Mul(t, lr.C[i]))
	}
	return v
}

func (lr *LinearRecursive[T]) doubling(b []T) []T {
	v := make([]T, lr.N)
	for i := range lr.N {
		v[i] = lr.Id0()
	}
	bb := slices.Clone(b)
	for i := range lr.N {
		mul := b[lr.N-1-i]
		for j := range lr.N {
			v[j] = lr.Add(v[j], lr.Mul(bb[j], mul))
		}
		bb = lr.increment(bb)
	}
	return v
}

func fib() {
	// 示例：斐波那契数列
	// a_n = 1*a_{n-1} + 1*a_{n-2}
	// C = {1, 1} (c_1, c_2)
	// a = {0, 1} (a_0, a_1)

	add := func(a, b int) int { return a + b }
	mul := func(a, b int) int { return a * b }
	id0 := func() int { return 0 }
	id1 := func() int { return 1 }

	lr := NewLinearRecursive([]int{1, 1}, add, mul, id0, id1)
	a := []int{0, 1}
	fmt.Println("Calculating Fibonacci numbers:")
	for k := range 10 {
		fmt.Printf("F(%d) = %d\n", k, lr.KthTerm(a, k))
	}
	fmt.Printf("F(30) = %d\n", lr.KthTerm(a, 30)) // 应该输出 832040
}

func yuki1460() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	const INF int = 2e18

	var k, n int
	fmt.Fscan(in, &k, &n)
	a := make([]int, k)
	b := make([]int, k)
	for i := range k {
		fmt.Fscan(in, &a[i])
	}
	for i := range k {
		fmt.Fscan(in, &b[i])
	}
	for i, j := 0, k-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}

	add := func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}
	mul := func(a, b int) int {
		if a < b {
			return a
		}
		return b
	}
	id0 := func() int { return -INF }
	id1 := func() int { return INF }
	lr := NewLinearRecursive(b, add, mul, id0, id1)
	fmt.Fprintln(out, lr.KthTerm(a, n))
}

func main() {
	yuki1460()
}
