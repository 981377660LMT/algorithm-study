/*
桃太郎电铁游戏里有+-两种格子
开始时左上角有一个棋子，每次可以向右或向下移动一格
移动到+格子得1分,移动到-格子得-1分
每个人都最佳移动策略,问最后谁会赢

ROW,COL<=2000 博弈dp计算分数差值
*/
package main

import (
	"bufio"
	"fmt"
	"os"
)

const INF int = 1e9

func gameInMomotetsuWorld(grid []string) int {
	ROW := len(grid)
	COL := len(grid[0])

	memo := [2010][2010]int{}
	for i := 0; i < 2010; i++ {
		for j := 0; j < 2010; j++ {
			memo[i][j] = -INF
		}
	}

	var dfs func(row, col int) int
	dfs = func(row, col int) int {
		if row == ROW-1 && col == COL-1 {
			return 0
		}

		if memo[row][col] != -INF {
			return memo[row][col]
		}

		res := -INF
		if row+1 < ROW {
			if grid[row+1][col] == '+' {
				res = max(res, 1-dfs(row+1, col))
			} else {
				res = max(res, -1-dfs(row+1, col))
			}
		}
		if col+1 < COL {
			if grid[row][col+1] == '+' {
				res = max(res, 1-dfs(row, col+1))
			} else {
				res = max(res, -1-dfs(row, col+1))
			}
		}

		memo[row][col] = res
		return res
	}

	return dfs(0, 0)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	var ROW, COL int
	fmt.Fscan(reader, &ROW, &COL)
	grid := make([]string, ROW)
	for i := 0; i < ROW; i++ {
		fmt.Fscan(reader, &grid[i])
	}
	res := gameInMomotetsuWorld(grid)
	if res > 0 {
		fmt.Fprint(writer, "Takahashi")
	} else if res < 0 {
		fmt.Fprint(writer, "Aoki")
	} else {
		fmt.Fprint(writer, "Draw")
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
