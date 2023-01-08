package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	const MOD int = 998244353

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var q int
	fmt.Fscan(in, &q)

	// f2(f1(x)) = mul2*(mul1*x+add1)+add2 = mul2*mul1*x+mul2*add1+add2
	window := NewSlidingWindowAggregationDeque(
		func() E { return E{mul: 1, add: 0} },
		func(f1, f2 E) E {
			return E{mul: f1.mul * f2.mul % MOD, add: (f1.add*f2.mul + f2.add) % MOD}
		})

	for i := 0; i < q; i++ {
		var op, mul, add, x int
		fmt.Fscan(in, &op)
		if op == 0 {
			fmt.Fscan(in, &mul, &add)
			window.AppendLeft(E{mul: mul, add: add})
		} else if op == 1 {
			fmt.Fscan(in, &mul, &add)
			window.Append(E{mul: mul, add: add})
		} else if op == 2 {
			window.PopLeft()
		} else if op == 3 {
			window.Pop()
		} else if op == 4 {
			fmt.Fscan(in, &x)
			aggregation := window.Query()
			fmt.Fprintln(out, (aggregation.mul*x+aggregation.add)%MOD)
		}
	}

}

// type E = int
type E = struct{ mul, add int }

type SlidingWindowAggregationDeque struct {
	e                                    func() E
	op                                   func(left, right E) E
	leftVal, rightVal, leftSum, rightSum []E
}

func NewSlidingWindowAggregationDeque(e func() E, op func(left, right E) E) *SlidingWindowAggregationDeque {
	return &SlidingWindowAggregationDeque{e: e, op: op}
}

// Aggregate all operations from left to right and return the result.
//  i.e. op(op(op(...op((x)))))...
func (s *SlidingWindowAggregationDeque) Query() E {
	if len(s.leftSum) == 0 && len(s.rightSum) == 0 {
		return s.e()
	}
	if len(s.rightSum) == 0 {
		return s.leftSum[len(s.leftSum)-1]
	}
	if len(s.leftSum) == 0 {
		return s.rightSum[len(s.rightSum)-1]
	}
	return s.op(s.leftSum[len(s.leftSum)-1], s.rightSum[len(s.rightSum)-1])
}

func (s *SlidingWindowAggregationDeque) Append(val E) {
	if len(s.rightSum) == 0 {
		s.rightVal = append(s.rightVal, val)
		s.rightSum = append(s.rightSum, val)
	} else {
		s.rightVal = append(s.rightVal, val)
		s.rightSum = append(s.rightSum, s.op(s.rightSum[len(s.rightSum)-1], val))
	}
}

func (s *SlidingWindowAggregationDeque) AppendLeft(val E) {
	if len(s.leftSum) == 0 {
		s.leftVal = append(s.leftVal, val)
		s.leftSum = append(s.leftSum, val)
	} else {
		s.leftVal = append(s.leftVal, val)
		s.leftSum = append(s.leftSum, s.op(val, s.leftSum[len(s.leftSum)-1]))
	}
}

// Reconstruct by brute force
func (s *SlidingWindowAggregationDeque) Pop() {
	if len(s.rightSum) == 0 {
		ln := len(s.leftSum) / 2
		rn := len(s.leftSum) - ln
		lv := make([]E, 0, ln)
		s.leftSum = s.leftSum[:0]
		for i := 0; i < ln; i++ {
			lv = append(lv, s.leftVal[len(s.leftVal)-1])
			s.leftVal = s.leftVal[:len(s.leftVal)-1]
		}

		for i := 0; i < rn; i++ {
			x := s.leftVal[len(s.leftVal)-1]
			s.leftVal = s.leftVal[:len(s.leftVal)-1]
			s.rightVal = append(s.rightVal, x)
			if len(s.rightSum) == 0 {
				s.rightSum = append(s.rightSum, x)
			} else {
				s.rightSum = append(s.rightSum, s.op(s.rightSum[len(s.rightSum)-1], x))
			}
		}

		for i := 0; i < ln; i++ {
			x := lv[len(lv)-1]
			lv = lv[:len(lv)-1]
			s.leftVal = append(s.leftVal, x)
			if len(s.leftSum) == 0 {
				s.leftSum = append(s.leftSum, x)
			} else {
				s.leftSum = append(s.leftSum, s.op(x, s.leftSum[len(s.leftSum)-1]))
			}
		}
	}

	s.rightVal = s.rightVal[:len(s.rightVal)-1]
	s.rightSum = s.rightSum[:len(s.rightSum)-1]
}

// Reconstruct by brute force
func (s *SlidingWindowAggregationDeque) PopLeft() {
	if len(s.leftSum) == 0 {
		rn := len(s.rightSum) / 2
		ln := len(s.rightSum) - rn
		rv := make([]E, 0, rn)
		s.rightSum = s.rightSum[:0]
		for i := 0; i < rn; i++ {
			rv = append(rv, s.rightVal[len(s.rightVal)-1])
			s.rightVal = s.rightVal[:len(s.rightVal)-1]
		}

		for i := 0; i < ln; i++ {
			x := s.rightVal[len(s.rightVal)-1]
			s.rightVal = s.rightVal[:len(s.rightVal)-1]
			s.leftVal = append(s.leftVal, x)
			if len(s.leftSum) == 0 {
				s.leftSum = append(s.leftSum, x)
			} else {
				s.leftSum = append(s.leftSum, s.op(x, s.leftSum[len(s.leftSum)-1]))
			}
		}

		for i := 0; i < rn; i++ {
			x := rv[len(rv)-1]
			rv = rv[:len(rv)-1]
			s.rightVal = append(s.rightVal, x)
			if len(s.rightSum) == 0 {
				s.rightSum = append(s.rightSum, x)
			} else {
				s.rightSum = append(s.rightSum, s.op(s.rightSum[len(s.rightSum)-1], x))
			}
		}
	}

	s.leftVal = s.leftVal[:len(s.leftVal)-1]
	s.leftSum = s.leftSum[:len(s.leftSum)-1]
}

func (s *SlidingWindowAggregationDeque) Len() int {
	return len(s.leftSum) + len(s.rightSum)
}
