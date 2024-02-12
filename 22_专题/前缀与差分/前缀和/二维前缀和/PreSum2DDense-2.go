package main

import "fmt"

func main() {
	S := [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	preSum := NewPreSum2DDenseFrom(S)
	fmt.Println(preSum.SumRegion(1, 1, 2, 2)) // 28
}

// 二维前缀和.
type PreSum2DDense struct {
	preSum [][]int
}

func NewPreSum2DDense(row, col int, f func(int, int) int) *PreSum2DDense {
	preSum := make([][]int, row+1)
	for i := range preSum {
		preSum[i] = make([]int, col+1)
	}
	for r := 0; r < row; r++ {
		for c := 0; c < col; c++ {
			preSum[r+1][c+1] = f(r, c) + preSum[r][c+1] + preSum[r+1][c] - preSum[r][c]
		}
	}
	return &PreSum2DDense{preSum}
}

func NewPreSum2DDenseFrom(matrix [][]int) *PreSum2DDense {
	return NewPreSum2DDense(len(matrix), len(matrix[0]), func(r, c int) int { return matrix[r][c] })
}

// 查询sum(A[r1:r2+1, c1:c2+1])的值.
// 0 <= r1 <= r2 < row, 0 <= c1 <= c2 < col.
func (ps *PreSum2DDense) SumRegion(row1, col1, row2, col2 int) int {
	return ps.preSum[row2+1][col2+1] - ps.preSum[row2+1][col1] - ps.preSum[row1][col2+1] + ps.preSum[row1][col1]
}
