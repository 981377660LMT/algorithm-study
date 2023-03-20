package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// https://judge.yosupo.jp/problem/queue_operate_all_composite
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	swag := NewSlidingWindowAggregation()
	var q int
	fmt.Fscan(in, &q)
	for i := 0; i < q; i++ {
		var op int
		fmt.Fscan(in, &op)
		if op == 0 {
			var a, b int
			fmt.Fscan(in, &a, &b)
			swag.Append(E{a, b})
		} else if op == 1 {
			swag.PopLeft()
		} else {
			var x int
			fmt.Fscan(in, &x)
			res := swag.Query()
			fmt.Fprintln(out, (res.mul*x+res.add)%MOD)
		}
	}
}

const MOD int = 998244353

type E = struct{ mul, add int }

func (*SlidingWindowAggregation) e() E { return E{1, 0} }
func (*SlidingWindowAggregation) op(a, b E) E {
	return E{a.mul * b.mul % MOD, (a.add*b.mul + b.add) % MOD}
}

type SlidingWindowAggregation struct {
	cumL []E
	cumR E
	dat  []E
	sz   int
}

func NewSlidingWindowAggregation() *SlidingWindowAggregation {
	res := &SlidingWindowAggregation{}
	res.cumL = []E{res.e()}
	res.cumR = res.e()
	return res
}

func (s *SlidingWindowAggregation) Len() int {
	return s.sz
}

func (s *SlidingWindowAggregation) Append(x E) {
	s.sz++
	s.cumR = s.op(s.cumR, x)
	s.dat = append(s.dat, x)
}

func (s *SlidingWindowAggregation) PopLeft() {
	s.sz--
	s.cumL = s.cumL[:len(s.cumL)-1]
	if len(s.cumL) == 0 {
		s.cumL = []E{s.e()}
		s.cumR = s.e()
		for len(s.dat) > 1 {
			s.cumL = append(s.cumL, s.op(s.dat[len(s.dat)-1], s.cumL[len(s.cumL)-1]))
			s.dat = s.dat[:len(s.dat)-1]
		}
		s.dat = s.dat[:0]
	}
}

func (s *SlidingWindowAggregation) Query() E {
	return s.op(s.cumL[len(s.cumL)-1], s.cumR)
}
