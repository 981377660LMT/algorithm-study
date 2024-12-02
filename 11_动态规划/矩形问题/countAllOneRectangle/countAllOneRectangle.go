package main

// 3359. 查找最大元素不超过 K 的有序子矩阵
// https://leetcode.cn/problems/find-sorted-submatrices-with-maximum-element-at-most-k/
func countSubmatrices(grid [][]int, k int) int64 {
	row, col := len(grid), len(grid[0])
	check := func(r, c int) bool { return grid[r][c] <= k }
	checkLeft := func(r, c int) bool { return grid[r][c] <= grid[r][c-1] }
	return int64(CountAllOneRectangle(row, col, check, checkLeft))
}

// https://leetcode.cn/problems/count-submatrices-with-all-ones/description/
func numSubmat(mat [][]int) int {
	row, col := len(mat), len(mat[0])
	check := func(r, c int) bool { return mat[r][c] == 1 }
	checkLeft := func(r, c int) bool { return true }
	return CountAllOneRectangle(row, col, check, checkLeft)
}

// O(row*col) 统计全 1 子矩形.
func CountAllOneRectangle(
	row, col int,
	check func(r, c int) bool, // 当前单元格是否合法
	checkLeft func(r, c int) bool, // 当前单元格左侧是否合法, c > 0
) int {
	rowDp := make([][]int, row)
	for i := 0; i < row; i++ {
		rowDp[i] = make([]int, col)
	}
	for r := 0; r < row; r++ {
		for c := 0; c < col; c++ {
			if check(r, c) {
				if c > 0 && checkLeft(r, c) {
					rowDp[r][c] = rowDp[r][c-1] + 1
				} else {
					rowDp[r][c] = 1
				}
			}
		}
	}

	res := 0
	for j := 0; j < col; j++ {
		var stack [][2]int
		total := 0
		for i := 0; i < row; i++ {
			height := 1
			for len(stack) > 0 && stack[len(stack)-1][0] > rowDp[i][j] {
				total -= stack[len(stack)-1][1] * (stack[len(stack)-1][0] - rowDp[i][j])
				height += stack[len(stack)-1][1]
				stack = stack[:len(stack)-1]
			}
			total += rowDp[i][j]
			res += total
			stack = append(stack, [2]int{rowDp[i][j], height})
		}
	}
	return res
}
