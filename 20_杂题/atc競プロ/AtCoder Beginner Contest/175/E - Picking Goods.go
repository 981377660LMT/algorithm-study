// https://atcoder.jp/contests/abc175/tasks/abc175_e

// !输入 n m (1≤n,m≤3000) k(≤min(2e5,r*c))，表示一个 n*m 的网格，和网格中的 k 个物品。
// 接下来 k 行，每行三个数 x y v(≤1e9) 表示物品的行号、列号和价值（行列号从 1 开始）。
// 每个网格至多有一个物品。

// 你从 (1,1) 出发走到 (n,m)，每步只能向下或向右。
// !经过物品时，你可以选或不选，且每行至多可以选三个物品。
// !输出你选到的物品的价值和的最大值。

package main

import (
	"bufio"
	"fmt"
	"os"
)

func solve(grid [][]int) int {
	ROW, COL := len(grid), len(grid[0])
	memo := [4][3010][3010]int{}
	// memo fill -1
	for i := 0; i < 4; i++ {
		for j := 0; j < 3010; j++ {
			for k := 0; k < 3010; k++ {
				memo[i][j][k] = -1
			}
		}
	}

	var dfs func(row, col, count int) int
	dfs = func(row, col, count int) int {
		if row == ROW-1 && col == COL-1 {
			return 0
		}
		if memo[count][row][col] != -1 {
			return memo[count][row][col]
		}
		res := 0
		if row+1 < ROW {
			res = max(res, dfs(row+1, col, 0))
			if grid[row+1][col] > 0 {
				res = max(res, grid[row+1][col]+dfs(row+1, col, 1))
			}
		}
		if col+1 < COL {
			res = max(res, dfs(row, col+1, count))
			if grid[row][col+1] > 0 && count <= 2 {
				res = max(res, grid[row][col+1]+dfs(row, col+1, count+1))
			}
		}
		memo[count][row][col] = res
		return res
	}

	res := dfs(0, 0, 0)
	if grid[0][0] > 0 { // 有物品
		res = max(res, grid[0][0]+dfs(0, 0, 1))
	}
	return res
}

func main() {

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m, k int
	fmt.Fscan(in, &n, &m, &k)
	grid := make([][]int, n)
	for i := range grid {
		grid[i] = make([]int, m)
	}
	for i := 0; i < k; i++ {
		var x, y, v int
		fmt.Fscan(in, &x, &y, &v)
		grid[x-1][y-1] = v
	}
	fmt.Fprintln(out, solve(grid))
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
