// 三维差分

package main

import "fmt"

func main() {
	e := func() int { return 0 }
	op := func(a, b int) int { return a + b }
	inv := func(a int) int { return -a }
	D := NewDiff3D(e, op, inv)
	x, y, z := int32(2), int32(3), int32(4)
	D.Init(x, y, z, func(i, j, k int32) int { return int(i*y*z + j*z + k) })

	printAll := func() {
		grid := make([][][]int, x)
		for i := int32(0); i < x; i++ {
			grid[i] = make([][]int, y)
			for j := int32(0); j < y; j++ {
				grid[i][j] = make([]int, z)
				for k := int32(0); k < z; k++ {
					grid[i][j][k] = D.Get(i, j, k)
				}
			}
		}
		fmt.Println(grid)
	}

	printAll()
	D.Add(1, 2, 1, 3, 1, 4, 1)
	printAll()
}

type Diff3D[E any] struct {
	Data    [][][]E
	dirty   bool
	x, y, z int32
	diff    [][][]E
	e       func() E
	op      func(a, b E) E
	inv     func(a E) E
}

func NewDiff3D[E any](e func() E, op func(a, b E) E, inv func(a E) E) *Diff3D[E] {
	return &Diff3D[E]{e: e, op: op, inv: inv}
}

func (d *Diff3D[E]) Init(x, y, z int32, f func(i, j, k int32) E) {
	data := make([][][]E, x)
	for i := int32(0); i < x; i++ {
		data[i] = make([][]E, y)
		for j := int32(0); j < y; j++ {
			data[i][j] = make([]E, z)
			for k := int32(0); k < z; k++ {
				data[i][j][k] = f(i, j, k)
			}
		}
	}
	diff := make([][][]E, x+1)
	for i := int32(0); i <= x; i++ {
		diff[i] = make([][]E, y+1)
		for j := int32(0); j <= y; j++ {
			diff[i][j] = make([]E, z+1)
			for k := int32(0); k <= z; k++ {
				diff[i][j][k] = d.e()
			}
		}
	}

	d.dirty = false
	d.x, d.y, d.z = x, y, z
	d.diff = diff
	d.Data = data
}

// [x1, x2) x [y1, y2)
func (d *Diff3D[E]) Add(x1, x2, y1, y2, z1, z2 int32, v E) {
	x1, y1, z1 = max32(x1, 0), max32(y1, 0), max32(z1, 0)
	x2, y2, z2 = min32(x2, d.x), min32(y2, d.y), min32(z2, d.z)
	if x1 >= x2 || y1 >= y2 || z1 >= z2 {
		return
	}
	d.dirty = true
	d.diff[x1][y1][z1] = d.op(d.diff[x1][y1][z1], v)
	d.diff[x1][y2][z1] = d.op(d.diff[x1][y2][z1], d.inv(v))
	d.diff[x2][y1][z1] = d.op(d.diff[x2][y1][z1], d.inv(v))
	d.diff[x2][y2][z1] = d.op(d.diff[x2][y2][z1], v)
	d.diff[x1][y1][z2] = d.op(d.diff[x1][y1][z2], d.inv(v))
	d.diff[x1][y2][z2] = d.op(d.diff[x1][y2][z2], v)
	d.diff[x2][y1][z2] = d.op(d.diff[x2][y1][z2], v)
	d.diff[x2][y2][z2] = d.op(d.diff[x2][y2][z2], d.inv(v))
}

func (d *Diff3D[E]) Get(x, y, z int32) E {
	if d.dirty {
		d.Build()
	}
	return d.Data[x][y][z]
}

func (d *Diff3D[E]) Build() {
	if !d.dirty {
		return
	}
	data, diff, e, op := d.Data, d.diff, d.e, d.op
	x, y, z := d.x, d.y, d.z

	for i := int32(1); i < x; i++ {
		for j := int32(0); j < y; j++ {
			for k := int32(0); k < z; k++ {
				diff[i][j][k] = op(diff[i][j][k], diff[i-1][j][k])
			}
		}
	}
	for i := int32(0); i < x; i++ {
		for j := int32(1); j < y; j++ {
			for k := int32(0); k < z; k++ {
				diff[i][j][k] = op(diff[i][j][k], diff[i][j-1][k])
			}
		}
	}
	for i := int32(0); i < x; i++ {
		for j := int32(0); j < y; j++ {
			for k := int32(1); k < z; k++ {
				diff[i][j][k] = op(diff[i][j][k], diff[i][j][k-1])
			}
		}
	}

	for i := int32(0); i < x; i++ {
		for j := int32(0); j < y; j++ {
			for k := int32(0); k < z; k++ {
				data[i][j][k] = op(data[i][j][k], diff[i][j][k])
				diff[i][j][k] = e()
			}
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
