package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	abc366d()
}

func demo() {
	e := func() int { return 0 }
	op := func(a, b int) int { return a + b }
	inv := func(a int) int { return -a }
	S := NewPresumDense3D(e, op, inv)
	mat := make([][][]int, 2)
	mat[0] = [][]int{
		{1, 2, 3},
		{4, 5, 6},
	}
	mat[1] = [][]int{
		{7, 8, 9},
		{10, 11, 12},
	}
	S.Build(int32(len(mat)), int32(len(mat[0])), int32(len(mat[0][0])), func(x, y, z int32) int { return mat[x][y][z] })
	fmt.Println(S.Query(0, 1, 0, 2, 0, 3))
}

// https://atcoder.jp/contests/abc366/tasks/abc366_d
// https://atcoder.jp/contests/abc366/editorial/10643
// 三维前缀和.
func abc366d() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int32
	fmt.Fscan(in, &n)
	e := func() int { return 0 }
	op := func(a, b int) int { return a + b }
	inv := func(a int) int { return -a }
	S := NewPresumDense3D(e, op, inv)

	mat := make([][][]int, n)
	for i := int32(0); i < n; i++ {
		mat[i] = make([][]int, n)
		for j := int32(0); j < n; j++ {
			mat[i][j] = make([]int, n)
			for k := int32(0); k < n; k++ {
				var a int
				fmt.Fscan(in, &a)
				mat[i][j][k] = a
			}
		}
	}
	S.Build(n, n, n, func(x, y, z int32) int { return mat[x][y][z] })

	var q int
	fmt.Fscan(in, &q)
	for i := 0; i < q; i++ {
		var lx, rx, ly, ry, lz, rz int32
		fmt.Fscan(in, &lx, &rx, &ly, &ry, &lz, &rz)
		res := S.Query(lx-1, rx, ly-1, ry, lz-1, rz)
		fmt.Fprintln(out, res)
	}
}

type PresumDense3D[E any] struct {
	x1, x2, x3 int32
	presum     [][][]E
	e          func() E
	op         func(a, b E) E
	inv        func(a E) E
}

func NewPresumDense3D[E any](e func() E, op func(a, b E) E, inv func(a E) E) *PresumDense3D[E] {
	return &PresumDense3D[E]{e: e, op: op, inv: inv}
}

func (p *PresumDense3D[E]) Build(x1, x2, x3 int32, f func(x1, x2, x3 int32) E) {
	e, op := p.e, p.op
	presum := make([][][]E, x1+1)
	for x := range presum {
		presum[x] = make([][]E, x2+1)
		for y := range presum[x] {
			row := make([]E, x3+1)
			for i := range row {
				row[i] = e()
			}
			presum[x][y] = row
		}
	}

	for x := int32(1); x <= x1; x++ {
		for y := int32(1); y <= x2; y++ {
			for z := int32(1); z <= x3; z++ {
				presum[x][y][z] = op(f(x-1, y-1, z-1), presum[x][y-1][z])
			}
		}
		for y := int32(1); y <= x2; y++ {
			for z := int32(1); z <= x3; z++ {
				presum[x][y][z] = op(presum[x][y][z], presum[x][y][z-1])
			}
		}
		for y := int32(1); y <= x2; y++ {
			for z := int32(1); z <= x3; z++ {
				presum[x][y][z] = op(presum[x][y][z], presum[x-1][y][z])
			}
		}
	}

	p.x1, p.x2, p.x3 = x1, x2, x3
	p.presum = presum
}

// [x1, x2) x [y1, y2) x [z1, z2)
func (p *PresumDense3D[E]) Query(x1, x2 int32, y1, y2 int32, z1, z2 int32) E {
	x1, y1, z1 = max32(x1, 0), max32(y1, 0), max32(z1, 0)
	x2, y2, z2 = min32(x2, p.x1), min32(y2, p.x2), min32(z2, p.x3)
	if x1 >= x2 || y1 >= y2 || z1 >= z2 {
		return p.e()
	}
	op, inv, presum := p.op, p.inv, p.presum
	res := presum[x2][y2][z2]
	res = op(res, inv(presum[x1][y2][z2]))
	res = op(res, inv(presum[x2][y1][z2]))
	res = op(res, inv(presum[x2][y2][z1]))
	res = op(res, presum[x1][y1][z2])
	res = op(res, presum[x1][y2][z1])
	res = op(res, presum[x2][y1][z1])
	res = op(res, inv(presum[x1][y1][z1]))
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
