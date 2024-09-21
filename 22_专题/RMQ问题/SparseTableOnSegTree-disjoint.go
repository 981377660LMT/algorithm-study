// DisjointSparseTableOnSegTree
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

	st := NewDisjointSparseTableOnSegTreeFrom(dp, func() int32 { return 0 }, max32)

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

// 3148. 矩阵中的最大得分
// https://leetcode.cn/problems/maximum-difference-score-in-a-grid/description/
// 给你一个由 正整数 组成、大小为 m x n 的矩阵 grid。你可以从矩阵中的任一单元格移动到另一个位于正下方或正右侧的任意单元格（不必相邻）。从值为 c1 的单元格移动到值为 c2 的单元格的得分为 c2 - c1 。
// 你可以从 任一 单元格开始，并且必须至少移动一次。
// 返回你能得到的 最大 总得分。

const INF32 int32 = 1e9 + 10

func maxScore(grid [][]int) int {
	grid32 := make([][]int32, len(grid))
	for i := range grid {
		grid32[i] = make([]int32, len(grid[i]))
		for j := range grid[i] {
			grid32[i][j] = int32(grid[i][j])
		}
	}

	ROW, COL := int32(len(grid32)), int32(len(grid32[0]))
	res := -INF32
	st := NewDisjointSparseTableOnSegTreeFrom(grid32, func() int32 { return INF32 }, min32)
	for i := int32(0); i < ROW; i++ {
		for j := int32(0); j < COL; j++ {
			// up
			if i > 0 {
				upMin := st.Query(0, i, j, j+1)
				res = max32(res, grid32[i][j]-upMin)
			}
			// left
			if j > 0 {
				leftMin := st.Query(i, i+1, 0, j)
				res = max32(res, grid32[i][j]-leftMin)
			}

			// left and up
			if i > 0 && j > 0 {
				upLeftMin := st.Query(0, i, 0, j)
				res = max32(res, grid32[i][j]-upLeftMin)
			}
		}
	}

	return int(res)
}

// 更快的 DisjointSparseTableFast.
type DisjointSparseTableOnSegTree[E any] struct {
	row, col int32
	e        func() E
	op       func(E, E) E
	data     []*DisjointSparseTable[E]
}

func NewDisjointSparseTableOnSegTreeFrom[E any](grid [][]E, e func() E, op func(E, E) E) *DisjointSparseTableOnSegTree[E] {
	row := int32(len(grid))
	col := int32(0)
	if row > 0 {
		col = int32(len(grid[0]))
	}
	data := make([]*DisjointSparseTable[E], 2*row)
	for i := int32(0); i < row; i++ {
		data[row+i] = NewDisjointSparseTableFrom(grid[i], e, op)
	}
	for i := row - 1; i > 0; i-- {
		data[i] = NewDisjointSparseTable(
			col,
			func(j int32) E {
				x := data[2*i].Query(j, j+1)
				y := data[2*i+1].Query(j, j+1)
				return op(x, y)
			},
			e, op,
		)
	}
	return &DisjointSparseTableOnSegTree[E]{row: row, col: col, e: e, op: op, data: data}
}

func (st *DisjointSparseTableOnSegTree[E]) Query(rowStart, rowEnd, colStart, colEnd int32) E {
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

type DisjointSparseTable[E any] struct {
	n    int32
	e    func() E
	op   func(E, E) E
	data [][]E
}

// DisjointSparseTable 支持幺半群的区间静态查询.
//
//	eg: 区间乘积取模/区间仿射变换...
func NewDisjointSparseTable[E any](n int32, f func(int32) E, e func() E, op func(E, E) E) *DisjointSparseTable[E] {
	res := &DisjointSparseTable[E]{}
	log := int32(1)
	for (1 << log) < n {
		log++
	}
	data := make([][]E, log)
	data[0] = make([]E, 0, n)
	for i := int32(0); i < n; i++ {
		data[0] = append(data[0], f(i))
	}
	for i := int32(1); i < log; i++ {
		data[i] = append(data[i], data[0]...)
		tmp := data[i]
		b := int32(1 << i)
		for m := b; m <= n; m += 2 * b {
			l, r := m-b, min32(m+b, n)
			for j := m - 1; j >= l+1; j-- {
				tmp[j-1] = op(tmp[j-1], tmp[j])
			}
			for j := m; j < r-1; j++ {
				tmp[j+1] = op(tmp[j], tmp[j+1])
			}
		}
	}
	res.n = n
	res.e = e
	res.op = op
	res.data = data
	return res
}

func NewDisjointSparseTableFrom[E any](leaves []E, e func() E, op func(E, E) E) *DisjointSparseTable[E] {
	return NewDisjointSparseTable(int32(len(leaves)), func(i int32) E { return leaves[i] }, e, op)
}

func (ds *DisjointSparseTable[E]) Query(start, end int32) E {
	if start >= end {
		return ds.e()
	}
	end--
	if start == end {
		return ds.data[0][start]
	}
	lca := bits.Len32(uint32(start^end)) - 1
	return ds.op(ds.data[lca][start], ds.data[lca][end])
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
