// SparseTableOnSegTree
// 二维rmq/线段树套st表/线段树套rmq
// O(r*c*log(c))构建,O(log(r))查询.
// https://maspypy.github.io/library/ds/sparse_table/sparse_table_on_segtree.hpp

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {
	cf713D()
}

// Animals and Puzzle (区间最大正方形)
// https://www.luogu.com.cn/problem/CF713D
// 题意：给定一个 n×m 的地图 a，ai​ 为 0 或 1。
// 有 q 次询问，每次询问给定一个矩形，求出这个矩形中最大的由 1 构成的正方形的边长是多少。
//
// 1.dp 预处理出每个点为左上角的最大正方形边长.
// 2.线段树套st表,二分答案查询二维区间最大值.
func cf713D() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var row, col int32
	fmt.Fscan(in, &row, &col)
	grid := make([][]int32, row)
	for i := int32(0); i < row; i++ {
		grid[i] = make([]int32, col)
		for j := int32(0); j < col; j++ {
			fmt.Fscan(in, &grid[i][j])
		}
	}
	dp := make([][]int32, row)
	for i := int32(0); i < row; i++ {
		dp[i] = make([]int32, col)
	}
	for i := row - 1; i >= 0; i-- {
		for j := col - 1; j >= 0; j-- {
			cur := min32(row-i, col-j)
			if i+1 < row {
				cur = min32(cur, dp[i+1][j]+1)
			}
			if j+1 < col {
				cur = min32(cur, dp[i][j+1]+1)
			}
			if i+1 < row && j+1 < col {
				cur = min32(cur, dp[i+1][j+1]+1)
			}
			if grid[i][j] == 0 {
				cur = 0
			}
			dp[i][j] = cur
		}
	}

	st := NewSparseTableOnSegTreeFrom(dp, func() int32 { return 0 }, max32)

	query := func(rowStart, rowEnd, colStart, colEnd int32) int32 {
		check := func(mid int32) bool {
			return st.Query(rowStart, rowEnd-mid+1, colStart, colEnd-mid+1) >= mid
		}

		left, right := int32(1), min32(rowEnd-rowStart, colEnd-colStart)
		for left <= right {
			mid := (left + right) / 2
			if check(mid) {
				left = mid + 1
			} else {
				right = mid - 1
			}
		}
		return right
	}

	var q int32
	fmt.Fscan(in, &q)
	for i := int32(0); i < q; i++ {
		var x1, y1, x2, y2 int32
		fmt.Fscan(in, &x1, &y1, &x2, &y2)
		x1--
		y1--
		fmt.Fprintln(out, query(x1, x2, y1, y2))
	}
}

type SparseTableOnSegTree[E any] struct {
	row, col int32
	e        func() E
	op       func(E, E) E
	data     []*SparseTable[E]
}

func NewSparseTableOnSegTreeFrom[E any](grid [][]E, e func() E, op func(E, E) E) *SparseTableOnSegTree[E] {
	row := int32(len(grid))
	col := int32(0)
	if row > 0 {
		col = int32(len(grid[0]))
	}
	data := make([]*SparseTable[E], 2*row)
	for i := int32(0); i < row; i++ {
		data[row+i] = NewSparseTableFrom(grid[i], e, op)
	}
	for i := row - 1; i > 0; i-- {
		data[i] = NewSparseTable(
			col,
			func(j int32) E {
				x := data[2*i].Query(j, j+1)
				y := data[2*i+1].Query(j, j+1)
				return op(x, y)
			},
			e, op,
		)
	}
	return &SparseTableOnSegTree[E]{row: row, col: col, e: e, op: op, data: data}
}

func (st *SparseTableOnSegTree[E]) Query(rowStart, rowEnd, colStart, colEnd int32) E {
	if !(0 <= rowStart && rowStart <= rowEnd && rowEnd <= st.row) {
		return st.e()
	}
	if !(0 <= colStart && colStart <= colEnd && colEnd <= st.col) {
		return st.e()
	}
	res := st.e()
	rowStart += st.row
	rowEnd += st.row
	for rowStart < rowEnd {
		if rowStart&1 != 0 {
			res = st.op(res, st.data[rowStart].Query(colStart, colEnd))
			rowStart++
		}
		if rowEnd&1 != 0 {
			rowEnd--
			res = st.op(res, st.data[rowEnd].Query(colStart, colEnd))
		}
		rowStart >>= 1
		rowEnd >>= 1
	}
	return res
}

type SparseTable[E any] struct {
	st [][]E
	e  func() E
	op func(E, E) E
	n  int32
}

func NewSparseTable[E any](n int32, f func(int32) E, e func() E, op func(E, E) E) *SparseTable[E] {
	res := &SparseTable[E]{}

	b := bits.Len(uint(n))
	st := make([][]E, b)
	for i := range st {
		st[i] = make([]E, n)
	}
	for i := int32(0); i < n; i++ {
		st[0][i] = f(i)
	}
	for i := 1; i < b; i++ {
		for j := int32(0); j+(1<<i) <= n; j++ {
			st[i][j] = op(st[i-1][j], st[i-1][j+(1<<(i-1))])
		}
	}
	res.st = st
	res.e = e
	res.op = op
	res.n = n
	return res
}

func NewSparseTableFrom[E any](leaves []E, e func() E, op func(E, E) E) *SparseTable[E] {
	return NewSparseTable(int32(len(leaves)), func(i int32) E { return leaves[i] }, e, op)
}

// 查询区间 [start, end) 的贡献值.
func (st *SparseTable[E]) Query(start, end int32) E {
	if start >= end {
		return st.e()
	}
	b := bits.Len(uint(end-start)) - 1
	return st.op(st.st[b][start], st.st[b][end-(1<<b)])
}

func max32(a, b int32) int32 {
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
