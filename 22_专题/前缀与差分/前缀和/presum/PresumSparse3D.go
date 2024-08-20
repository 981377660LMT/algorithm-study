package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	abc366d()
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
	S.Build(n*n*n, func(i int32) (x, y, z int, e int) {
		x = int(i / n / n)
		y = int(i / n % n)
		z = int(i % n)
		e = mat[x][y][z]
		return
	})

	var q int
	fmt.Fscan(in, &q)
	for i := 0; i < q; i++ {
		var lx, rx, ly, ry, lz, rz int
		fmt.Fscan(in, &lx, &rx, &ly, &ry, &lz, &rz)
		res := S.Query(lx-1, rx, ly-1, ry, lz-1, rz)
		fmt.Fprintln(out, res)
	}
}

type PresumSparse3D[E any] struct {
	x, y, z                   int32
	presum                    [][][]E
	originX, originY, originZ []int
	e                         func() E
	op                        func(a, b E) E
	inv                       func(a E) E
}

func NewPresumDense3D[E any](e func() E, op func(a, b E) E, inv func(a E) E) *PresumSparse3D[E] {
	return &PresumSparse3D[E]{e: e, op: op, inv: inv}
}

func (p *PresumSparse3D[E]) Build(n int32, f func(i int32) (x, y, z int, e E)) {
	xs, ys, zs, es := make([]int, n), make([]int, n), make([]int, n), make([]E, n)
	for i := int32(0); i < n; i++ {
		xs[i], ys[i], zs[i], es[i] = f(i)
	}
	newXs, originX := discretize1D(xs)
	newYs, originY := discretize1D(ys)
	newZs, originZ := discretize1D(zs)
	e, op := p.e, p.op
	x1, x2, x3 := int32(len(originX)), int32(len(originY)), int32(len(originZ))
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
	for i := int32(0); i < n; i++ {
		x, y, z, e := newXs[i], newYs[i], newZs[i], es[i]
		presum[x+1][y+1][z+1] = op(presum[x+1][y+1][z+1], e)
	}

	for x := int32(1); x <= x1; x++ {
		for y := int32(1); y <= x2; y++ {
			for z := int32(1); z <= x3; z++ {
				presum[x][y][z] = op(presum[x][y][z], presum[x][y-1][z])
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

	p.x, p.y, p.z = x1, x2, x3
	p.presum = presum
	p.originX, p.originY, p.originZ = originX, originY, originZ
}

// [x1, x2) x [y1, y2) x [z1, z2)
func (p *PresumSparse3D[E]) Query(x1, x2 int, y1, y2 int, z1, z2 int) E {
	if x1 >= x2 || y1 >= y2 || z1 >= z2 {
		return p.e()
	}
	nx1, nx2 := p.compressX(x1), p.compressX(x2)
	ny1, ny2 := p.compressY(y1), p.compressY(y2)
	nz1, nz2 := p.compressZ(z1), p.compressZ(z2)
	op, inv, presum := p.op, p.inv, p.presum
	res := presum[nx2][ny2][nz2]
	res = op(res, inv(presum[nx1][ny2][nz2]))
	res = op(res, inv(presum[nx2][ny1][nz2]))
	res = op(res, inv(presum[nx2][ny2][nz1]))
	res = op(res, presum[nx1][ny1][nz2])
	res = op(res, presum[nx1][ny2][nz1])
	res = op(res, presum[nx2][ny1][nz1])
	res = op(res, inv(presum[nx1][ny1][nz1]))
	return res
}

func (p *PresumSparse3D[E]) compressX(x int) int32 {
	return int32(sort.SearchInts(p.originX, x))
}

func (p *PresumSparse3D[E]) compressY(y int) int32 {
	return int32(sort.SearchInts(p.originY, y))
}

func (p *PresumSparse3D[E]) compressZ(z int) int32 {
	return int32(sort.SearchInts(p.originZ, z))
}

func discretize1D(nums []int) (newNums []int32, origin []int) {
	newNums = make([]int32, len(nums))
	origin = make([]int, 0, len(newNums))
	order := argSort(int32(len(nums)), func(i, j int32) bool { return nums[i] < nums[j] })
	for _, i := range order {
		if len(origin) == 0 || origin[len(origin)-1] != nums[i] {
			origin = append(origin, nums[i])
		}
		newNums[i] = int32(len(origin) - 1)
	}
	origin = origin[:len(origin):len(origin)]
	return
}

func argSort(n int32, less func(i, j int32) bool) []int32 {
	order := make([]int32, n)
	for i := range order {
		order[i] = int32(i)
	}
	sort.Slice(order, func(i, j int) bool { return less(order[i], order[j]) })
	return order
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
