package main

const INF int = 1e18

type Presum2DMonoid[T any] struct {
	preSum   [][]T
	row, col int32
	e        func() T
}

func NewPresum2DMonoid[T any](matrix [][]T, e func() T, op func(T, T) T) *Presum2DMonoid[T] {
	row, col := int32(len(matrix)), int32(len(matrix[0]))
	preSum := make([][]T, row+1)
	for i := range preSum {
		row := make([]T, col+1)
		for j := range row {
			row[j] = e()
		}
		preSum[i] = row
	}
	for r := int32(0); r < row; r++ {
		tmpSum0, tmpSum1 := preSum[r], preSum[r+1]
		tmpM := matrix[r]
		for c := int32(0); c < col; c++ {
			tmpSum1[c+1] = op(op(tmpSum0[c+1], tmpSum1[c]), tmpM[c])
		}
	}
	return &Presum2DMonoid[T]{preSum: preSum, row: row, col: col, e: e}
}

// 查询[0:x, 0:y]的聚合值.
func (p *Presum2DMonoid[T]) Query(x, y int32) T {
	if x <= 0 || y <= 0 {
		return p.e()
	}
	if x > p.row {
		x = p.row
	}
	if y > p.col {
		y = p.col
	}
	return p.preSum[x][y]
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

// 3148. 矩阵中的最大得分
// https://leetcode.cn/problems/maximum-difference-score-in-a-grid/
func maxScore(grid [][]int) int {
	S := NewPresum2DMonoid(grid, func() int { return INF }, min)
	res := -INF
	row, col := int32(len(grid)), int32(len(grid[0]))
	for i := int32(0); i < row; i++ {
		for j := int32(0); j < col; j++ {
			res = max(res, grid[i][j]-min(S.Query(i, j+1), S.Query(i+1, j)))
		}
	}
	return res
}
