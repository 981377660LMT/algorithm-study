// 指配问题/线性分配问题.
package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	yosupo()
}

const INF int = 1e18

// 3276. 选择矩阵中单元格的最大得分
func maxScore(grid [][]int) int {
	max_ := 0
	for i := 0; i < len(grid); i++ {
		for j := 0; j < len(grid[0]); j++ {
			max_ = max(max_, grid[i][j])
		}
	}
	costMatrix := make([][]int, len(grid))
	for r, row := range grid {
		costMatrix[r] = make([]int, max_+1)
		for _, v := range row {
			costMatrix[r][v] = v
		}
	}
	res, _, _, _, _ := LinearSumAssignment(costMatrix, true)
	return res
}

// 2463. 最小移动总距离
// https://leetcode.cn/problems/minimum-total-distance-traveled/description/
func minimumTotalDistance(robot []int, factory [][]int) int64 {
	var targets []int
	for _, f := range factory {
		pos, limit := f[0], f[1]
		for i := 0; i < limit; i++ {
			targets = append(targets, pos)
		}
	}

	costMatrix := make([][]int, len(robot))
	for i, r := range robot {
		costMatrix[i] = make([]int, len(targets))
		for j, t := range targets {
			costMatrix[i][j] = abs(r - t)
		}
	}

	_, rowIndex, colIndex, _, _ := LinearSumAssignment(costMatrix, false)
	res := 0
	for i := 0; i < len(rowIndex); i++ {
		res += costMatrix[rowIndex[i]][colIndex[i]]
	}
	return int64(res)
}

// https://judge.yosupo.jp/problem/assignment
func yosupo() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)

	costMatrix := make([][]int, n)
	for i := 0; i < n; i++ {
		costMatrix[i] = make([]int, n)
		for j := 0; j < n; j++ {
			fmt.Fscan(in, &costMatrix[i][j])
		}
	}

	res, _, colIndex, _, _ := LinearSumAssignment(costMatrix, false)
	fmt.Fprintln(out, res)
	for _, v := range colIndex {
		fmt.Fprint(out, v, " ")
	}
}

// O(n^2m).
// 线性分配问题.
// 返回最大/最小权值和, 行索引, 列索引, 对偶问题的解X, Y.
//
// 这里的对偶问题是:
//
//	max(sum(x_i) + sum(y_j)), subj to x_i + y_j <= C_{ij}.
func LinearSumAssignment(costMatrix [][]int, maximize bool) (res int, rowIndex, colIndex []int32, X, Y []int) {
	n, m := int32(len(costMatrix)), int32(len(costMatrix[0]))
	if n > m {
		newCostMatrix := make([][]int, m)
		for i := int32(0); i < m; i++ {
			newCostMatrix[i] = make([]int, n)
			for j := int32(0); j < n; j++ {
				newCostMatrix[i][j] = costMatrix[j][i]
			}
		}
		res, colIndex, rowIndex, Y, X = LinearSumAssignment(newCostMatrix, maximize)
		return
	}

	A := make([][]int, n+1)
	for i := int32(0); i < n; i++ {
		A[i+1] = make([]int, m+1)
		for j := int32(0); j < m; j++ {
			if maximize {
				A[i+1][j+1] = -costMatrix[i][j]
			} else {
				A[i+1][j+1] = costMatrix[i][j]
			}
		}
	}
	n++
	m++

	P := make([]int32, m)
	way := make([]int32, m)
	X = make([]int, n)
	Y = make([]int, m)
	minV := make([]int, m)
	used := make([]bool, m)

	for i := int32(1); i < n; i++ {
		P[0] = i
		for j := int32(0); j < m; j++ {
			minV[j] = INF
			used[j] = false
		}
		j0 := int32(0)
		for P[j0] != 0 {
			i0, j1 := P[j0], int32(0)
			used[j0] = true
			delta := INF
			for j := int32(1); j < m; j++ {
				if used[j] {
					continue
				}
				curr := A[i0][j] - X[i0] - Y[j]
				if curr < minV[j] {
					minV[j] = curr
					way[j] = j0
				}
				if minV[j] < delta {
					delta = minV[j]
					j1 = j
				}
			}
			for j := int32(0); j < m; j++ {
				if used[j] {
					X[P[j]] += delta
					Y[j] -= delta
				} else {
					minV[j] -= delta
				}
			}
			j0 = j1
		}

		for {
			P[j0] = P[way[j0]]
			j0 = way[j0]
			if j0 == 0 {
				break
			}
		}
	}

	res = -Y[0]
	X = X[1:]
	Y = Y[1:]
	rowIndex = make([]int32, n-1)
	for i := int32(0); i < n-1; i++ {
		rowIndex[i] = i
	}
	colIndex = make([]int32, n)
	for i := int32(0); i < m; i++ {
		colIndex[P[i]] = i
	}
	colIndex = colIndex[1:]
	for i := range colIndex {
		colIndex[i]--
	}
	if maximize {
		res = -res
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

func mins(nums ...int) int {
	res := nums[0]
	for _, num := range nums {
		if num < res {
			res = num
		}
	}
	return res
}

func maxs(nums ...int) int {
	res := nums[0]
	for _, num := range nums {
		if num > res {
			res = num
		}
	}
	return res
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func abs32(x int32) int32 {
	if x < 0 {
		return -x
	}
	return x
}
