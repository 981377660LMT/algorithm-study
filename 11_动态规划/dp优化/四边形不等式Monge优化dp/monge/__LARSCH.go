// https://maspypy.github.io/library/convex/larsch.hpp
// https://noshi91.github.io/Library/algorithm/larsch.cpp.html

package main

func main() {
	larcsh := NewLARSCH(5, func(i, j int) int { return i + j })
	println(larcsh.GetArgmin())
}

type LARSCH struct {
	base *_reduceRow
}

func NewLARSCH(n int, f func(i, j int) int) *LARSCH {
	res := &LARSCH{base: newReduceRow(n)}
	res.base.setF(f)
	return res
}

func (l *LARSCH) GetArgmin() int {
	return l.base.getArgmin()
}

type _reduceRow struct {
	n      int
	f      func(i, j int) int
	curRow int
	state  int
	rec    *_reduceCol
}

func newReduceRow(n int) *_reduceRow {
	res := &_reduceRow{n: n}
	m := n / 2
	if m != 0 {

		res.rec = newReduceCol(m)
	}
	return res
}

func (r *_reduceRow) setF(f func(i, j int) int) {
	r.f = f
	if r.rec != nil {
		r.rec.setF(func(i, j int) int {
			return f(2*i+1, j)
		})
	}
}

func (r *_reduceRow) getArgmin() int {
	curRow := r.curRow
	r.curRow += 1
	if curRow%2 == 0 {
		prevArgmin := r.state
		var nextArgmin int
		if curRow+1 == r.n {
			nextArgmin = r.n - 1
		} else {
			nextArgmin = r.rec.getArgmin()
		}
		r.state = nextArgmin
		ret := prevArgmin
		for j := prevArgmin + 1; j <= nextArgmin; j += 1 {
			if r.f(curRow, ret) > r.f(curRow, j) {
				ret = j
			}
		}
		return ret
	}

	if r.f(curRow, r.state) <= r.f(curRow, curRow) {
		return r.state
	}
	return curRow
}

type _reduceCol struct {
	n      int
	f      func(i, j int) int
	curRow int
	cols   []int
	rec    *_reduceRow
}

func newReduceCol(n int) *_reduceCol {
	return &_reduceCol{n: n, rec: newReduceRow(n)}
}

func (c *_reduceCol) setF(f func(i, j int) int) {
	c.f = f
	c.rec.setF(func(i, j int) int {
		return f(i, c.cols[j])
	})
}

func (r *_reduceCol) getArgmin() int {
	curRow := r.curRow
	r.curRow += 1
	var cs []int
	if curRow == 0 {
		cs = []int{0}
	} else {
		cs = []int{2*curRow - 1, 2 * curRow}
	}

	for _, j := range cs {
		for {
			size := len(r.cols)
			flag := size != curRow && r.f(size-1, r.cols[size-1]) > r.f(size-1, j)
			if !flag {
				break
			}
			r.cols = r.cols[:size-1]
		}
		if len(r.cols) != r.n {
			r.cols = append(r.cols, j)
		}
	}
	return r.cols[r.rec.getArgmin()]
}
