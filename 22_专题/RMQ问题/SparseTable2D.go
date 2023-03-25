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

type S = int

//  query: 查询 [row1,col1,row2,col2] 闭区间的贡献值
//     0 <= row1 <= row2 < len(matrix)
//     0 <= col1 <= col2 < len(matrix[0])
func NewSparseTable2D(matrix [][]S, op func(S, S) S) (query func(row1, col1, row2, col2 int) S) {
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

	query = func(row1, col1, row2, col2 int) S {
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
