package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	abc366d()
	// demo()
}

func demo() {
	e := func() int { return 0 }
	op := func(a, b int) int { return a + b }
	inv := func(a int) int { return -a }
	S := NewPresumDenseKD(e, op, inv)
	grid := [][]int{{1, 2}, {3, 4}, {5, 6}}
	S.Build([]int32{int32(len(grid)), int32(len(grid[0]))}, func(ds []int32) int { return grid[ds[0]][ds[1]] })
	fmt.Println(S.Query([]int32{1, 1}, []int32{3, 2}))
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
	bases      []int32
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
	bases := make([]int32, len(sizes))
	base := int32(1)
	for i := len(sizes) - 1; i >= 0; i-- {
		bases[i] = base
		base *= sizes[i] + 1
	}
	presum := make([]E, size)
	for i := range presum {
		presum[i] = e()
	}

	{
		// 初始化.
		enumeratePrefix(sizes, func(indices []int32, num int32) {
			presum[num] = f(indices)
		})
		// 累加每一个维度.
		mins := make([]int32, len(sizes))
		for i := range mins {
			mins[i] = 1
		}
		for d := int32(0); d < int32(len(sizes)); d++ {
			enumerateDigits(mins, sizes, bases, func(num int32) {
				presum[num] = op(presum[num], presum[num-bases[d]])
			})
		}
	}

	p.dimensions = sizes
	p.presum = presum
	p.bases = bases
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
	bases, dimensions := p.bases, p.dimensions

	n := int32(len(dimensions))

	// 按照容斥原理遍历.
	res := e()
	var dfs func(int32, int32, bool)
	dfs = func(pos, num int32, flag bool) {
		if pos == n {
			if flag {
				res = op(res, p.presum[num])
			} else {
				res = op(res, inv(p.presum[num]))
			}
			return
		}
		dfs(pos+1, num+bases[pos]*lowers[pos], !flag)
		dfs(pos+1, num+bases[pos]*uppers[pos], flag)
	}
	dfs(0, 0, true)
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

func enumeratePrefix(sizes []int32, f func(indicies []int32, num int32)) {
	n := int32(len(sizes))
	var dfs func(int32, []int32, int32, int32)
	dfs = func(pos int32, indicies []int32, num int32, base int32) {
		if pos == -1 {
			f(indicies, num)
			return
		}
		for indicies[pos] = 0; indicies[pos] < sizes[pos]; indicies[pos]++ {
			dfs(pos-1, indicies, num+(indicies[pos]+1)*base, base*(sizes[pos]+1))
		}
	}
	dfs(n-1, make([]int32, len(sizes)), 0, 1)
}

func enumerateDigits(mins []int32, maxs []int32, bases []int32, f func(num int32)) {
	n := int32(len(bases))
	var dfs func(int32, int32)
	dfs = func(pos int32, num int32) {
		if pos == n {
			f(num)
			return
		}
		for i := mins[pos]; i <= maxs[pos]; i++ {
			dfs(pos+1, num+bases[pos]*i)
		}
	}
	dfs(0, 0)
}
