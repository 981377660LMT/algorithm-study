// 二维RMQ/二维ST表
// https://codeforces.com/blog/entry/45485
// Preprocess : O(n*m*logn*logm)
//  we create a table[ 1+logn ][ n ][ 1+logm ][ m ]
//  Each box of the table[ 1+logn ][ n ] is a sparse table of size [ 1+logm ][ m ]
//  Let us see what table[ jr ][ ir ][ jc ][ ic ] actually contains:
//  It contains the minimum element from column ic to ic-1+2^jc of all rows from ir to ir-1+2^jr
//  In other words, it contain the minimum element in the submatrix [ (ir,ic), (ir-1+2^jr , ic-1+2^jc) ]
//  where submatrix [ (x1,y1),(x2,y2) ] denotes the submatrix with x1,y1 as its top left-most and x2,y2 as its bottom right-most point.
//  Now you can easily conclude that, table[ 0 ][ ir ][ jc ][ ic ] is nothing but the 1D RMQ table if we take our array as row ir

package main

import (
	"fmt"
	"math/bits"
)

func main() {
	matrix := [][]int{
		{1, 2, 3, 4, 5},
		{6, 7, 8, 9, 10},
		{11, 12, 13, 14, 15},
	}

	query := NewSparseTable2D(matrix, max)
	fmt.Println(query(1, 1, 2, 2))
}

// 3148. 矩阵中的最大得分
// https://leetcode.cn/problems/maximum-difference-score-in-a-grid/description/
// 给你一个由 正整数 组成、大小为 m x n 的矩阵 grid。你可以从矩阵中的任一单元格移动到另一个位于正下方或正右侧的任意单元格（不必相邻）。从值为 c1 的单元格移动到值为 c2 的单元格的得分为 c2 - c1 。
// 你可以从 任一 单元格开始，并且必须至少移动一次。
// 返回你能得到的 最大 总得分。

const INF int32 = 1e9 + 10

func maxScore(grid [][]int) int {
	grid32 := make([][]int32, len(grid))
	for i := range grid {
		grid32[i] = make([]int32, len(grid[i]))
		for j := range grid[i] {
			grid32[i][j] = int32(grid[i][j])
		}
	}

	ROW, COL := int32(len(grid32)), int32(len(grid32[0]))
	res := -INF
	st := NewSparseTable2D(grid32, min32)
	for i := int32(0); i < ROW; i++ {
		for j := int32(0); j < COL; j++ {
			// up
			if i > 0 {
				upMin := st(0, j, i-1, j)
				res = max32(res, grid32[i][j]-upMin)
			}
			// left
			if j > 0 {
				leftMin := st(i, 0, i, j-1)
				res = max32(res, grid32[i][j]-leftMin)
			}

			// left and up
			if i > 0 && j > 0 {
				upLeftMin := st(0, 0, i-1, j-1)
				res = max32(res, grid32[i][j]-upLeftMin)
			}
		}
	}

	return int(res)
}

// query: 查询 [row1,col1,row2,col2] 闭区间的贡献值
//
//	0 <= row1 <= row2 < len(matrix)
//	0 <= col1 <= col2 < len(matrix[0])
func NewSparseTable2D[S any](matrix [][]S, op func(S, S) S) (query func(row1, col1, row2, col2 int32) S) {
	n, m := int32(len(matrix)), int32(len(matrix[0]))
	rowSize := int32(bits.Len32(uint32(n))) // 1+logn
	colSize := int32(bits.Len32(uint32(m))) // 1+logm
	dp := make([][][][]S, rowSize)          // (rowSize * row * colSize * col)
	for i := int32(0); i < rowSize; i++ {
		dp[i] = make([][][]S, n)
		for j := int32(0); j < n; j++ {
			dp[i][j] = make([][]S, colSize)
			for k := int32(0); k < colSize; k++ {
				dp[i][j][k] = make([]S, m)
			}
		}
	}

	for ir := int32(0); ir < n; ir++ {
		for ic := int32(0); ic < m; ic++ {
			dp[0][ir][0][ic] = matrix[ir][ic]
		}
		for jc := int32(1); jc < colSize; jc++ {
			for ic := int32(0); ic+(1<<jc) <= m; ic++ {
				dp[0][ir][jc][ic] = op(dp[0][ir][jc-1][ic], dp[0][ir][jc-1][ic+(1<<(jc-1))])
			}
		}
	}

	for jr := int32(1); jr < rowSize; jr++ {
		for ir := int32(0); ir+(1<<jr) <= n; ir++ {
			for jc := int32(0); jc < colSize; jc++ {
				for ic := int32(0); ic+(1<<jc) <= m; ic++ {
					dp[jr][ir][jc][ic] = op(dp[jr-1][ir][jc][ic], dp[jr-1][ir+(1<<(jr-1))][jc][ic])
				}
			}
		}
	}

	query = func(row1, col1, row2, col2 int32) S {
		rowk := bits.Len32(uint32(row2-row1+1)) - 1
		colK := bits.Len32(uint32(col2-col1+1)) - 1
		res1 := op(dp[rowk][row1][colK][col1], dp[rowk][row1][colK][col2-(1<<colK)+1])
		res2 := op(dp[rowk][row2-(1<<rowk)+1][colK][col1], dp[rowk][row2-(1<<rowk)+1][colK][col2-(1<<colK)+1])
		return op(res1, res2)
	}

	return
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
