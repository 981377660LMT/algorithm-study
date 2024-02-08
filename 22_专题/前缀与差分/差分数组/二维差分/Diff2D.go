package main

import "fmt"

func main() {
	S := [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	preSum := NewPreSum2DFrom(S)
	fmt.Println(preSum.QueryRange(1, 1, 2, 2)) // 28
	D := NewDiff2DFrom(S)
	D.Add(1, 1, 2, 2, 10)
	fmt.Println(D.Query(1, 1)) // 15
}

type Diff2D struct {
	Matrix   [][]int
	diff     [][]int
	row, col int
	dirty    bool
}

// 二维差分.
func NewDiff2D(row, col int, f func(int, int) int) *Diff2D {
	matrix := make([][]int, row)
	for i := range matrix {
		matrix[i] = make([]int, col)
		for j := range matrix[i] {
			matrix[i][j] = f(i, j)
		}
	}
	diff := make([][]int, row+2)
	for i := range diff {
		diff[i] = make([]int, col+2)
	}
	return &Diff2D{Matrix: matrix, diff: diff, row: row, col: col}
}

func NewDiff2DFrom(mat [][]int) *Diff2D {
	return NewDiff2D(len(mat), len(mat[0]), func(r, c int) int { return mat[r][c] })
}

// 区间更新左上角`(row1, col1)` 到 右下角`(row2, col2)`闭区间所描述的子矩阵的元素.
// 0<=r1<=r2<row, 0<=c1<=c2<col.
func (d *Diff2D) Add(r1, c1, r2, c2, delta int) {
	d.diff[r1+1][c1+1] += delta
	d.diff[r1+1][c2+2] -= delta
	d.diff[r2+2][c1+1] -= delta
	d.diff[r2+2][c2+2] += delta
	d.dirty = true
}

// 查询矩阵中指定位置的元素.
func (d *Diff2D) Query(r, c int) int {
	if d.dirty {
		d.Build()
	}
	return d.Matrix[r][c]
}

// 遍历矩阵，还原对应元素的增量.
func (d *Diff2D) Build() {
	if !d.dirty {
		return
	}
	d.dirty = false
	for i := 0; i < d.row; i++ {
		tmpDiff0, tmpDiff1 := d.diff[i], d.diff[i+1]
		tmpMatrix := d.Matrix[i]
		for j := 0; j < d.col; j++ {
			tmpDiff1[j+1] += tmpDiff1[j] + tmpDiff0[j+1] - tmpDiff0[j]
			tmpMatrix[j] += tmpDiff1[j+1]
		}
	}
}

type PreSum2D struct {
	preSum [][]int
}

func NewPreSum2D(row, col int, f func(int, int) int) *PreSum2D {
	preSum := make([][]int, row+1)
	for i := range preSum {
		preSum[i] = make([]int, col+1)
	}
	for r := 0; r < row; r++ {
		for c := 0; c < col; c++ {
			preSum[r+1][c+1] = f(r, c) + preSum[r][c+1] + preSum[r+1][c] - preSum[r][c]
		}
	}
	return &PreSum2D{preSum}
}

func NewPreSum2DFrom(mat [][]int) *PreSum2D {
	return NewPreSum2D(len(mat), len(mat[0]), func(r, c int) int { return mat[r][c] })
}

// 查询sum(A[r1:r2+1, c1:c2+1])的值.
// 0 <= r1 <= r2 < row, 0 <= c1 <= c2 < col.
func (ps *PreSum2D) QueryRange(row1, col1, row2, col2 int) int {
	return ps.preSum[row2+1][col2+1] - ps.preSum[row2+1][col1] - ps.preSum[row1][col2+1] + ps.preSum[row1][col1]
}
