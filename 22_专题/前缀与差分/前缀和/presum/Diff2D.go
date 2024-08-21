// 二维差分

package main

import "fmt"

func main() {
	e := func() int { return 0 }
	op := func(a, b int) int { return a + b }
	inv := func(a int) int { return -a }
	D := NewDiff2D(e, op, inv)
	row, col := int32(2), int32(3)
	D.Init(row, col, func(r, c int32) int { return int(r*col + c) })

	printAll := func() {
		grid := make([][]int, row)
		for i := int32(0); i < row; i++ {
			grid[i] = make([]int, col)
			for j := int32(0); j < col; j++ {
				grid[i][j] = D.Get(i, j)
			}
		}
		fmt.Println(grid)
	}

	printAll()
	D.Add(0, 2, 0, 3, 1)
	fmt.Println(D.diff)
	printAll()
}

type Diff2D[E any] struct {
	dirty bool
	n, m  int32
	diff  [][]E
	data  [][]E
	e     func() E
	op    func(a, b E) E
	inv   func(a E) E
}

func NewDiff2D[E any](e func() E, op func(a, b E) E, inv func(a E) E) *Diff2D[E] {
	return &Diff2D[E]{e: e, op: op, inv: inv}
}

func (d *Diff2D[E]) Init(n, m int32, f func(r, c int32) E) {
	diff := make([][]E, n)
	data := make([][]E, n)
	for i := int32(0); i < n; i++ {
		diffRow, row := make([]E, m), make([]E, m)
		for j := int32(0); j < m; j++ {
			diffRow[j] = d.e()
			row[j] = f(i, j)
		}
		diff[i], data[i] = diffRow, row
	}

	d.dirty = false
	d.n, d.m = n, m
	d.diff = diff
	d.data = data
}

// [x1, x2) x [y1, y2)
func (d *Diff2D[E]) Add(x1, x2, y1, y2 int32, v E) {
	x1, y1 = max32(x1, 0), max32(y1, 0)
	x2, y2 = min32(x2, d.n), min32(y2, d.m)
	if x1 >= x2 || y1 >= y2 {
		return
	}
	d.dirty = true
	d.diff[x1][y1] = d.op(d.diff[x1][y1], v)
	if x2 < d.n {
		d.diff[x2][y1] = d.op(d.diff[x2][y1], d.inv(v))
	}
	if y2 < d.m {
		d.diff[x1][y2] = d.op(d.diff[x1][y2], d.inv(v))
	}
	if x2 < d.n && y2 < d.m {
		d.diff[x2][y2] = d.op(d.diff[x2][y2], v)
	}
}

func (d *Diff2D[E]) Get(r, c int32) E {
	if d.dirty {
		d.rebuild()
	}
	return d.data[r][c]
}

func (d *Diff2D[E]) rebuild() {
	if !d.dirty {
		return
	}
	data, diff, e, op := d.data, d.diff, d.e, d.op
	for i := int32(0); i < d.n; i++ {
		for j := int32(1); j < d.m; j++ {
			diff[i][j] = op(diff[i][j], diff[i][j-1])
		}
	}
	for i := int32(1); i < d.n; i++ {
		for j := int32(0); j < d.m; j++ {
			diff[i][j] = op(diff[i][j], diff[i-1][j])
		}
	}
	for i := int32(0); i < d.n; i++ {
		for j := int32(0); j < d.m; j++ {
			data[i][j] = op(data[i][j], diff[i][j])
			diff[i][j] = e()
		}
	}
	d.dirty = false
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
