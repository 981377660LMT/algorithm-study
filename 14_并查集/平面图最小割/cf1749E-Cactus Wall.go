// 2 种遇到有后效性动态规划的处理方法：高斯消元和最短路
// https://www.luogu.com.cn/article/urd6r0r7

package main

import (
	"bufio"
	"fmt"
	"os"
)

// https://www.luogu.com.cn/problem/CF1749E
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	solve := func(grid [][]bool) (res [][]bool, ok bool) {}

	var T int32
	fmt.Fscan(in, &T)
	for i := int32(0); i < T; i++ {
		var row, col int32
		fmt.Fscan(in, &row, &col)
		grid := make([][]bool, row)
		for i := int32(0); i < row; i++ {
			var s string
			fmt.Fscan(in, &s)
			grid[i] = make([]bool, col)
			for j, c := range s {
				grid[i][j] = c == '#'
			}
		}

		res, ok := solve(grid)
		if !ok {
			fmt.Fprintln(out, "NO")
		} else {
			fmt.Fprintln(out, "YES")
			for _, row := range res {
				for _, b := range row {
					if b {
						fmt.Fprint(out, "#")
					} else {
						fmt.Fprint(out, ".")
					}
				}
				fmt.Fprintln(out)
			}
		}
	}
}
