package main

import (
	"bufio"
	"fmt"
	"os"
)

// https://judge.yosupo.jp/problem/static_range_sum
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int32
	fmt.Fscan(in, &n, &q)
	nums := make([]int, n)
	for i := range nums {
		fmt.Fscan(in, &nums[i])
	}

	e := func() int { return 0 }
	op := func(a, b int) int { return a + b }
	inv := func(a int) int { return -a }

	S := NewStaticRangeProductGroup(e, op, inv)
	S.BuildFrom(nums)

	for i := int32(0); i < q; i++ {
		var l, r int32
		fmt.Fscan(in, &l, &r)
		fmt.Fprintln(out, S.Query(l, r))
	}
}

type StaticRangeProductGroup[E any] struct {
	e    func() E
	op   func(a, b E) E
	inv  func(a E) E
	data []E
}

func NewStaticRangeProductGroup[E any](e func() E, op func(a, b E) E, inv func(E) E) *StaticRangeProductGroup[E] {
	return &StaticRangeProductGroup[E]{e: e, op: op, inv: inv}
}

func (s *StaticRangeProductGroup[E]) Build(m int32, f func(int32) E) {
	s.data = make([]E, m+1)
	s.data[0] = s.e()
	for i := int32(0); i < m; i++ {
		s.data[i+1] = s.op(s.data[i], f(i))
	}
}

func (s *StaticRangeProductGroup[E]) BuildFrom(a []E) {
	s.Build(int32(len(a)), func(i int32) E { return a[i] })
}

func (s *StaticRangeProductGroup[E]) Query(start, end int32) E {
	if start >= end {
		return s.e()
	}
	return s.op(s.inv(s.data[start]), s.data[end])
}
