package main

import "fmt"

func main() {
	grid := [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	diagonalPresum := NewDiagonalPresum(grid)
	fmt.Println(diagonalPresum.QueryDiagonal([2]int{0, 0}, [2]int{2, 2}))
	fmt.Println(diagonalPresum.QueryAntiDiagonal([2]int{2, 0}, [2]int{0, 2}))
}

var DIR4 = [][4]int{{-1, 0}, {0, -1}, {1, 0}, {0, 1}}

// 对角线前缀和.
type DiagonalPresum struct {
	preSum1 [][]int
	preSum2 [][]int
}

func NewDiagonalPresum(matrix [][]int) *DiagonalPresum {
	row, col := len(matrix), len(matrix[0])
	preSum1 := make([][]int, row+1) // 主对角线前缀和 ↘
	preSum2 := make([][]int, row+1) // 反对角线前缀和 ↙
	for i := range preSum1 {
		preSum1[i] = make([]int, col+1)
		preSum2[i] = make([]int, col+1)
	}
	for i, r := range matrix {
		tmp1 := preSum1[i]
		tmp2 := preSum1[i+1]
		tmp3 := preSum2[i]
		tmp4 := preSum2[i+1]
		for j, v := range r {
			tmp2[j+1] = tmp1[j] + v
			tmp4[j] = tmp3[j+1] + v
		}
	}
	return &DiagonalPresum{preSum1: preSum1, preSum2: preSum2}
}

// 正对角线左上角到右下角的前缀和 ↘.
func (d *DiagonalPresum) QueryDiagonal(leftUp, rightDown [2]int) int {
	r1, c1, r2, c2 := leftUp[0], leftUp[1], rightDown[0], rightDown[1]
	return d.preSum1[r2+1][c2+1] - d.preSum1[r1][c1]
}

// 副对角线左下角到右上角的前缀和 ↙.
func (d *DiagonalPresum) QueryAntiDiagonal(leftDown, rightUp [2]int) int {
	r1, c1, r2, c2 := leftDown[0], leftDown[1], rightUp[0], rightUp[1]
	return d.preSum2[r1+1][c1] - d.preSum2[r2][c2+1]
}
