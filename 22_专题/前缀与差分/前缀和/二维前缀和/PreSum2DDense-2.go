package main

import "fmt"

func main() {
	S := [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	preSum := NewPreSum2DDenseFrom(S)
	fmt.Println(preSum.SumRegion(1, 1, 2, 2)) // 6
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

// 查询sum(A[x1:x2+1, y1:y2+1])的值(包含边界).
// 0 <= x1 <= x2 < row, 0 <= y1 <= y2 < col.
func (ps *PreSum2DDense) SumRegion(x1, x2, y1, y2 int) int {
	if x1 > x2 || y1 > y2 {
		return 0
	}
	return ps.preSum[x2+1][y2+1] - ps.preSum[x2+1][y1] -
		ps.preSum[x1][y2+1] + ps.preSum[x1][y1]
}
