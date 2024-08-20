package main

import "fmt"

func main() {
	e := func() int { return 0 }
	op := func(a, b int) int { return a + b }
	inv := func(a int) int { return -a }
	S := NewPresumDense1D(e, op, inv)
	S.Build(10, func(i int32) int { return int(i) })
	fmt.Println(S.Query(0, 10))
	fmt.Println(S.Query(3, 5))
}

type PresumDense1D[E any] struct {
	n      int32
	presum []E
	e      func() E
	op     func(a, b E) E
	inv    func(a E) E
}

func NewPresumDense1D[E any](e func() E, op func(a, b E) E, inv func(a E) E) *PresumDense1D[E] {
	return &PresumDense1D[E]{e: e, op: op, inv: inv}
}

func (p *PresumDense1D[E]) Build(n int32, f func(i int32) E) {
	presum := make([]E, n+1)
	presum[0] = p.e()
	for i := int32(0); i < n; i++ {
		presum[i+1] = p.op(presum[i], f(i))
	}
	p.n = n
	p.presum = presum
}

func (p *PresumDense1D[E]) Query(start, end int32) E {
	if start < 0 {
		start = 0
	}
	if end > p.n {
		end = p.n
	}
	if start >= end {
		return p.e()
	}
	return p.op(p.inv(p.presum[start]), p.presum[end])
}
