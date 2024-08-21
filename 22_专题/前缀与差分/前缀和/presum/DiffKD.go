// 高维差分

package main

import "fmt"

func main() {
	e := func() int { return 0 }
	op := func(a, b int) int { return a + b }
	inv := func(a int) int { return -a }
	D := NewDiffKD(e, op, inv)
	x, y, z := int32(2), int32(3), int32(4)
	D.Init([]int32{x, y, z}, func(indices []int32) int { return int(indices[0]*y*z + indices[1]*z + indices[2]) })

	printAll := func() {
		grid := make([][][]int, x)
		for i := int32(0); i < x; i++ {
			grid[i] = make([][]int, y)
			for j := int32(0); j < y; j++ {
				grid[i][j] = make([]int, z)
				for k := int32(0); k < z; k++ {
					grid[i][j][k] = D.Get([]int32{i, j, k})
				}
			}
		}
		fmt.Println(grid)
	}

	printAll()
	D.Add([]int32{1, 1, 1}, []int32{2, 3, 4}, 1)
	printAll()
}

// 2536. 子矩阵元素加 1
// https://leetcode.cn/problems/increment-submatrices-by-one/description/
func rangeAddQueries(n int, queries [][]int) [][]int {
	e := func() int { return 0 }
	op := func(a, b int) int { return a + b }
	inv := func(a int) int { return -a }
	D := NewDiffKD(e, op, inv)
	D.Init([]int32{int32(n), int32(n)}, func(indices []int32) int { return 0 })
	for _, q := range queries {
		D.Add([]int32{int32(q[0]), int32(q[1])}, []int32{int32(q[2] + 1), int32(q[3] + 1)}, 1)
	}

	res := make([][]int, n)
	for i := 0; i < n; i++ {
		res[i] = make([]int, n)
		for j := 0; j < n; j++ {
			res[i][j] = D.Get([]int32{int32(i), int32(j)})
		}
	}

	return res
}

type DiffKD[E any] struct {
	dirty bool
	sizes []int32
	bases []int32
	diff  []E
	data  []E
	e     func() E
	op    func(a, b E) E
	inv   func(a E) E
}

func NewDiffKD[E any](e func() E, op func(a, b E) E, inv func(a E) E) *DiffKD[E] {
	return &DiffKD[E]{e: e, op: op, inv: inv}
}

func (d *DiffKD[E]) Init(sizes []int32, f func(indices []int32) E) {
	sizes = append(sizes[:0:0], sizes...)

	bases := make([]int32, len(sizes))
	base := int32(1)
	for i := len(sizes) - 1; i >= 0; i-- {
		bases[i] = base
		base *= sizes[i] + 1
	}

	data := make([]E, base)
	enumerateProductWithNum(sizes, bases, func(indices []int32, num int32) {
		data[num] = f(indices)
	})

	diff := make([]E, base)
	for i := range diff {
		diff[i] = d.e()
	}

	d.dirty = false
	d.sizes = sizes
	d.bases = bases
	d.diff = diff
	d.data = data
}

// [lowers[i], uppers[i]) += v.
func (d *DiffKD[E]) Add(lowers, uppers []int32, v E) {
	for i := range lowers {
		if lowers[i] < 0 {
			lowers[i] = 0
		}
		if uppers[i] > d.sizes[i] {
			uppers[i] = d.sizes[i]
		}
		if lowers[i] >= uppers[i] {
			return
		}
	}
	d.dirty = true
	n := int32(len(d.sizes))
	bases, diff := d.bases, d.diff
	var dfs func(int32, int32, bool)
	dfs = func(pos, num int32, flag bool) {
		if pos == n {
			if flag {
				diff[num] = d.op(diff[num], v)
			} else {
				diff[num] = d.op(diff[num], d.inv(v))
			}
			return
		}
		dfs(pos+1, num+bases[pos]*lowers[pos], flag)
		dfs(pos+1, num+bases[pos]*uppers[pos], !flag)
	}
	dfs(0, 0, true)
}

func (d *DiffKD[E]) Get(indices []int32) E {
	if d.dirty {
		d.Build()
	}
	num := int32(0)
	for i, v := range indices {
		num += v * d.bases[i]
	}
	return d.data[num]
}

func (d *DiffKD[E]) Build() {
	if !d.dirty {
		return
	}

	data, diff, e, op := d.data, d.diff, d.e, d.op
	sizes, bases := d.sizes, d.bases
	lowers := make([]int32, len(sizes))
	uppers := sizes
	for d := int32(0); d < int32(len(sizes)); d++ {
		lowers[d] = 1
		enumerateDigits(lowers, uppers, bases, func(num int32) {
			diff[num] = op(diff[num], diff[num-bases[d]])
		})
		lowers[d] = 0
	}

	enumerateDigits(lowers, uppers, bases, func(num int32) {
		data[num] = op(data[num], diff[num])
		diff[num] = e()
	})

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

// EnumerateDigits.
func enumerateProductWithNum(sizes []int32, bases []int32, f func(digits []int32, num int32)) {
	n := int32(len(sizes))
	var dfs func(int32, []int32, int32)
	dfs = func(pos int32, digits []int32, num int32) {
		if pos == n {
			f(digits, num)
			return
		}
		for digits[pos] = 0; digits[pos] < sizes[pos]; digits[pos]++ {
			dfs(pos+1, digits, num+bases[pos]*digits[pos])
		}
	}
	dfs(0, make([]int32, len(sizes)), 0)
}

func enumerateDigits(lowers []int32, uppers []int32, bases []int32, f func(num int32)) {
	n := int32(len(bases))
	var dfs func(int32, int32)
	dfs = func(pos int32, num int32) {
		if pos == n {
			f(num)
			return
		}
		for i := lowers[pos]; i < uppers[pos]; i++ {
			dfs(pos+1, num+bases[pos]*i)
		}
	}
	dfs(0, 0)
}
