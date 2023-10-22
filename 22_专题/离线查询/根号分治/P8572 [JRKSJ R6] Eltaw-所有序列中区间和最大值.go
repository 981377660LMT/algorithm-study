// https://www.luogu.com.cn/problem/P8572
// https://www.luogu.com.cn/blog/cyffff/solution-JRKSJ-Eltaw

// P8572 [JRKSJ R6] Eltaw
// 又ROW个长为COL的序列nums_1,nums_2,...,nums_k
// 有q次询问,每次询问给出一个区间[left,right]
// !求出所有序列在区间[left,right]的和的最大值
// ROW,COL,q<=5e5 ROW*COL<=5e5

// !(注意到ROW*COL<=5e5这个奇怪的条件)

// !COL<=ROW时(列数很少，预处理所有查询)
// 询问的区间只有O(COL^2)种 `预处理所有查询`一共(ROW*COL*COL) 即O(5e5*sqrt(5e5))

// !COL>=ROW时(行数很少，直接每行都遍历一遍)
// 每次询问都要(ROW)回答 一共O(q*ROW) 即O(q*sqrt(5e5))

package main

import (
	"bufio"
	"fmt"
	"os"
)

const INF int = 1e18

func P8572(grid [][]int, queries [][2]int) []int {
	ROW, COL := len(grid), len(grid[0])
	if COL <= ROW {
		return Solve1(grid, queries)
	} else {
		return Solve2(grid, queries)
	}
}

// 列数很少，预处理所有查询
func Solve1(grid [][]int, queries [][2]int) []int {
	ROW, COL := len(grid), len(grid[0])
	preSums := make([][]int, ROW)
	for i, row := range grid {
		curPreSum := make([]int, COL+1)
		for j, v := range row {
			curPreSum[j+1] = curPreSum[j] + v
		}
		preSums[i] = curPreSum
	}

	memo := make([][]int, COL+1) // 左端点为i，右端点为j的区间和的最大值
	for i := 0; i <= COL; i++ {
		cur := make([]int, COL+1)
		for j := i; j <= COL; j++ {
			max_ := -INF
			for _, preSum := range preSums {
				max_ = max(max_, preSum[j]-preSum[i])
			}
			cur[j] = max_
		}
		memo[i] = cur
	}

	res := make([]int, len(queries))
	for i, query := range queries {
		start, end := query[0], query[1]
		res[i] = memo[start][end]
	}
	return res
}

// 行数很少，直接每行都遍历一遍
func Solve2(grid [][]int, queries [][2]int) []int {
	ROW, COL := len(grid), len(grid[0])
	preSums := make([][]int, ROW)
	for i, row := range grid {
		curPreSum := make([]int, COL+1)
		for j, v := range row {
			curPreSum[j+1] = curPreSum[j] + v
		}
		preSums[i] = curPreSum
	}

	res := make([]int, len(queries))
	for i, query := range queries {
		start, end := query[0], query[1]
		max_ := -INF
		for _, preSum := range preSums {
			max_ = max(max_, preSum[end]-preSum[start])
		}
		res[i] = max_
	}
	return res
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

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k, q int
	fmt.Fscan(in, &n, &k, &q)
	grid := make([][]int, k)
	for i := 0; i < k; i++ {
		curRow := make([]int, n)
		for j := 0; j < n; j++ {
			fmt.Fscan(in, &curRow[j])
		}
		grid[i] = curRow
	}

	queries := make([][2]int, q)
	for i := 0; i < q; i++ {
		var start, end int
		fmt.Fscan(in, &start, &end)
		start--
		queries[i] = [2]int{start, end}
	}

	res := P8572(grid, queries)
	for _, v := range res {
		fmt.Fprintln(out, v)
	}
}
