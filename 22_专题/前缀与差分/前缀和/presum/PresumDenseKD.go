package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// abc366d()
	demo()
}

func demo() {
	e := func() int { return 0 }
	op := func(a, b int) int { return a + b }
	inv := func(a int) int { return -a }
	S := NewPresumDenseKD(e, op, inv)
	grid := [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	S.Build([]int32{int32(len(grid)), int32(len(grid[0]))}, func(ds []int32) int { return grid[ds[0]][ds[1]] })
	fmt.Println(S.Query([]int32{0, 2}, []int32{3, 3}))
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
	S := NewPresumDenseKD(e, op, inv)

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
	S.Build([]int32{n, n, n}, func(ds []int32) int { return mat[ds[0]][ds[1]][ds[2]] })

	var q int
	fmt.Fscan(in, &q)
	for i := 0; i < q; i++ {
		var lx, rx, ly, ry, lz, rz int32
		fmt.Fscan(in, &lx, &rx, &ly, &ry, &lz, &rz)
		res := S.Query([]int32{lx - 1, ly - 1, lz - 1}, []int32{rx, ry, rz})
		fmt.Fprintln(out, res)
	}
}

// K维前缀和/高维前缀和.
type PresumDenseKD[E any] struct {
	dimensions []int32
	presum     []E
	e          func() E
	op         func(a, b E) E
	inv        func(a E) E
}

func NewPresumDenseKD[E any](e func() E, op func(a, b E) E, inv func(a E) E) *PresumDenseKD[E] {
	return &PresumDenseKD[E]{e: e, op: op, inv: inv}
}

func (p *PresumDenseKD[E]) Build(sizes []int32, f func(indices []int32) E) {
	e, op := p.e, p.op
	size := int32(1)
	for _, d := range sizes {
		size *= (d + 1)
	}
	presum := make([]E, size)
	for i := range presum {
		presum[i] = e()
	}

	enumeratePrefix(sizes, func(indices []int32, num int32) bool {
		presum[num] = f(indices)
		return false
	})
	fmt.Println(presum)

	var dfs func(int32, int32, int32)
	dfs = func(pos, num, base int32) {
		if pos == -1 {
			return
		}
		dfs(pos-1, num, base*(sizes[pos]+1))
		dfs(pos-1, num+base, base*(sizes[pos]+1))
		for i := sizes[pos] - 1; i >= 0; i-- {
			presum[num+(i+1)*base] = op(presum[num+(i+1)*base], presum[num+i*base])
		}
	}
	dfs(int32(len(sizes))-1, 0, 1)

	fmt.Println(presum)
	p.dimensions = sizes
	p.presum = presum
}

// [lowers, uppers)
func (p *PresumDenseKD[E]) Query(lowers, uppers []int32) E {
	for i, d := range p.dimensions {
		if lowers[i] < 0 {
			lowers[i] = 0
		}
		if uppers[i] > d {
			uppers[i] = d
		}
		if lowers[i] >= uppers[i] {
			return p.e()
		}
	}
	e, op, inv := p.e, p.op, p.inv
	dimensions := p.dimensions

	res := e()
	var dfs func(int32, int32, int32, bool)
	dfs = func(pos, num int32, base int32, flag bool) {
		if pos == -1 {
			if flag {
				res = op(res, p.presum[num])
			} else {
				res = op(res, inv(p.presum[num]))
			}
			return
		}
		dfs(pos-1, num+(lowers[pos])*base, base*(dimensions[pos]+1), !flag)
		dfs(pos-1, num+(uppers[pos])*base, base*(dimensions[pos]+1), flag)
	}
	dfs(int32(len(dimensions))-1, 0, 1, true)
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

func enumeratePrefix(sizes []int32, f func(digits []int32, num int32) bool) {
	n := int32(len(sizes))
	var dfs func(int32, []int32, int32, int32) bool
	dfs = func(pos int32, digits []int32, num int32, base int32) bool {
		if pos == -1 {
			return f(digits, num)
		}
		for digits[pos] = 0; digits[pos] < sizes[pos]; digits[pos]++ {
			if dfs(pos-1, digits, num+(digits[pos]+1)*base, base*(sizes[pos]+1)) {
				return true
			}
		}
		return false
	}
	dfs(n-1, make([]int32, len(sizes)), 0, 1)
}
