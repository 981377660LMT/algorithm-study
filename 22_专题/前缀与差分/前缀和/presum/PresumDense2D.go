package main

import "fmt"

func main() {
	e := func() int { return 0 }
	op := func(a, b int) int { return a + b }
	inv := func(a int) int { return -a }
	S := NewPresumDense2D(e, op, inv)
	grid := [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	S.Build(int32(len(grid)), int32(len(grid[0])), func(r, c int32) int { return grid[r][c] })
	fmt.Println(S.Query(0, 3, 0, 3))
	fmt.Println(S.Query(1, 2, 1, 2))
}

type PresumDense2D[E any] struct {
	n, m   int32
	presum [][]E
	e      func() E
	op     func(a, b E) E
	inv    func(a E) E
}

func NewPresumDense2D[E any](e func() E, op func(a, b E) E, inv func(a E) E) *PresumDense2D[E] {
	return &PresumDense2D[E]{e: e, op: op, inv: inv}
}

func (p *PresumDense2D[E]) Build(n int32, m int32, f func(r, c int32) E) {
	e, op := p.e, p.op
	presum := make([][]E, n+1)
	for i := int32(0); i < n+1; i++ {
		presum[i] = make([]E, m+1)
		for j := int32(0); j < m+1; j++ {
			presum[i][j] = e()
		}
	}
	for i := int32(1); i < n+1; i++ {
		for j := int32(1); j < m+1; j++ {
			presum[i][j] = op(f(i-1, j-1), presum[i][j-1])
		}
		for j := int32(1); j < m+1; j++ {
			presum[i][j] = op(presum[i][j], presum[i-1][j])
		}
	}
	p.n, p.m = n, m
	p.presum = presum
}

// [x1, x2) x [y1, y2)
func (p *PresumDense2D[E]) Query(x1, x2 int32, y1, y2 int32) E {
	x1, y1 = max32(x1, 0), max32(y1, 0)
	x2, y2 = min32(x2, p.n), min32(y2, p.m)
	if x1 >= x2 || y1 >= y2 {
		return p.e()
	}
	res := p.presum[x2][y2]
	res = p.op(res, p.inv(p.presum[x1][y2]))
	res = p.op(res, p.inv(p.presum[x2][y1]))
	res = p.op(res, p.presum[x1][y1])
	return res
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

func min32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func max32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}
