package main

import "math/bits"

// 我们希望找出 grid 中每个 3 x 3 矩阵中的最大值。
func largestLocal(grid [][]int) [][]int {
	n, m := len(grid), len(grid[0])
	query := NewSparseTable2D(grid, max)
	res := make([][]int, n-2)
	for i := range res {
		res[i] = make([]int, m-2)
		for j := range res[i] {
			res[i][j] = query(i, j, i+2, j+2)
		}
	}
	return res
}

type S = int

//  query: 查询 [row1,col1,row2,col2] 闭区间的贡献值
//     0 <= row1 <= row2 < len(matrix)
//     0 <= col1 <= col2 < len(matrix[0])
func NewSparseTable2D(matrix [][]S, op func(S, S) S) (query func(row1, col1, row2, col2 int) int) {
	n, m := len(matrix), len(matrix[0])
	rowSize := bits.Len(uint(n))   // 1+logn
	colSize := bits.Len(uint(m))   // 1+logm
	dp := make([][][][]S, rowSize) // (rowSize * row * colSize * col)
	for i := 0; i < rowSize; i++ {
		dp[i] = make([][][]S, n)
		for j := 0; j < n; j++ {
			dp[i][j] = make([][]S, colSize)
			for k := 0; k < colSize; k++ {
				dp[i][j][k] = make([]S, m)
			}
		}
	}

	for ir := 0; ir < n; ir++ {
		for ic := 0; ic < m; ic++ {
			dp[0][ir][0][ic] = matrix[ir][ic]
		}
		for jc := 1; jc < colSize; jc++ {
			for ic := 0; ic+(1<<jc) <= m; ic++ {
				dp[0][ir][jc][ic] = op(dp[0][ir][jc-1][ic], dp[0][ir][jc-1][ic+(1<<(jc-1))])
			}
		}
	}

	for jr := 1; jr < rowSize; jr++ {
		for ir := 0; ir+(1<<jr) <= n; ir++ {
			for jc := 0; jc < colSize; jc++ {
				for ic := 0; ic+(1<<jc) <= m; ic++ {
					dp[jr][ir][jc][ic] = op(dp[jr-1][ir][jc][ic], dp[jr-1][ir+(1<<(jr-1))][jc][ic])
				}
			}
		}
	}

	query = func(row1, col1, row2, col2 int) int {
		rowk := bits.Len(uint(row2-row1+1)) - 1
		colK := bits.Len(uint(col2-col1+1)) - 1
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
