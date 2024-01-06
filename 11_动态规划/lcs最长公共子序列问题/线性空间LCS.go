// Hirschberg算法是一种用于解决最长公共子序列（LCS）问题的算法，
// 它是由Dan Hirschberg在1975年提出的。
// 这个算法的主要优点是它只需要线性的空间复杂度，这使得它在处理大数据集时非常有效。
//
// Hirschberg算法的基本思想是使用分治策略。
// 它首先将输入序列分为两半，然后找到跨越这两半的最长公共子序列。然后，算法递归地在每一半上应用相同的过程。
//
// 以下是Hirschberg算法的基本步骤：
//
// 1. 如果序列的长度为1，返回该序列。
// 2. 将序列分为两半。
// 3. 计算跨越两半的最长公共子序列。
// 4. 递归地在每一半上应用Hirschberg算法。
// 5. 返回结果。
//
// Hirschberg算法的主要优点是它的空间复杂度为O(n)，这使得它在处理大数据集时非常有效。
// 然而，它的时间复杂度为O(n^2)，这可能在处理非常大的数据集时成为一个问题。
// https://github.com/ShahjalalShohag/code-library/blob/main/Dynamic%20Programming%20Optimizations/Hirschbergs%20Algorithm.cpp

package main

import (
	"fmt"
	"sort"
)

// https://csacademy.com/contest/archive/task/classic-task/
func main() {
	arr1 := []int{57, 421, 1248, 68, 21}
	arr2 := []int{12478, 125, 41287, 57, 248, 4897}
	getValue := func(i, j int) int {
		return (arr1[i] + j + 1) ^ (arr2[j] + i + 1)
	}
	fmt.Println(LCSHirschberg(3, 1, getValue))
}

// 64. 最小路径和
// https://leetcode.cn/problems/0i0mDW/description/
func minPathSum(grid [][]int) int {
	path := LCSHirschberg(len(grid), len(grid[0]), func(i, j int) int {
		return grid[i][j]
	})
	r, c := 0, 0
	res := grid[0][0]
	for _, v := range path {
		if v == "D" {
			r++
		} else {
			c++
		}
		res += grid[r][c]
	}
	return res
}

const INF int = 1e18

// 给定一个长为row宽为col的隐式矩阵A,A[i][j] = getValue(i,j).
// 求从A[0][0]到A[row-1][col-1]的最短路径.
// 返回值是"R"和"D"组成的操作序列.
func LCSHirschberg(row, col int, getValue func(i, j int) int) []string {
	pos := [][2]int{}
	var dfs func(int, int, int, int)
	dfs = func(li, lj, ri, rj int) {
		mid := (lj + rj) / 2
		height := ri - li + 1
		if rj-lj < 1 {
			return
		}
		if height == 1 {
			pos = append(pos, [2]int{mid, li})
			dfs(li, lj, li, mid)
			dfs(li, mid+1, li, rj)
			return
		}

		wLeft := mid - lj + 1
		dp := make([][]int, 2)
		for i := 0; i < 2; i++ {
			dp[i] = make([]int, height)
		}
		dp[0][0] = getValue(li, lj)
		for i := 1; i < height; i++ {
			dp[0][i] = dp[0][i-1] + getValue(li+i, lj)
		}
		f := uint8(1)
		for j := 1; j < wLeft; j++ {
			for i := 0; i < height; i++ {
				dp[f][i] = INF
			}
			for i := 0; i < height; i++ {
				val := getValue(li+i, lj+j)
				dp[f][i] = min(dp[f][i], dp[f^1][i]+val)
				if i-1 >= 0 {
					dp[f][i] = min(dp[f][i], dp[f][i-1]+val)
				}
			}
			f ^= 1
		}
		f ^= 1
		m1 := make([]int, height)
		for i := 0; i < height; i++ {
			m1[i] = dp[f][i]
		}
		wRight := rj - mid
		dp[0][0] = getValue(ri, rj)
		for i := 1; i < height; i++ {
			dp[0][i] = dp[0][i-1] + getValue(ri-i, rj)
		}
		f = 1
		for j := 1; j < wRight; j++ {
			for i := 0; i < height; i++ {
				dp[f][i] = INF
			}
			for i := 0; i < height; i++ {
				val := getValue(ri-i, rj-j)
				dp[f][i] = min(dp[f][i], dp[f^1][i]+val)
				if i-1 >= 0 {
					dp[f][i] = min(dp[f][i], dp[f][i-1]+val)
				}
			}
			f ^= 1
		}
		f ^= 1
		m2 := make([]int, height)
		for i := 0; i < height; i++ {
			m2[height-i-1] = dp[f][i]
		}
		mi := INF
		res := -1
		for i := 0; i < height; i++ {
			sum := m1[i] + m2[i]
			if sum < mi {
				mi = sum
				res = i
			}
		}
		res += li
		pos = append(pos, [2]int{mid, res})
		dfs(li, lj, res, mid)
		dfs(res, mid+1, ri, rj)
	}
	dfs(0, 0, row-1, col-1)
	sort.Slice(pos, func(i, j int) bool { return pos[i][0] < pos[j][0] })

	x, y := 0, 0
	res := make([]string, 0)
	for {
		if x == col-1 {
			for y != row-1 {
				res = append(res, "D")
				y++
			}
			break
		}
		if pos[x][1] == y {
			x++
			res = append(res, "R")
		} else {
			y++
			res = append(res, "D")
		}
	}

	return res
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
