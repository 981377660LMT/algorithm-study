// https://www.luogu.com.cn/problem/P8572
// # https://www.luogu.com.cn/blog/cyffff/solution-JRKSJ-Eltaw

// # P8572 [JRKSJ R6] Eltaw
// # 又ROW个长为COL的序列nums_1,nums_2,...,nums_k
// # 有q次询问,每次询问给出一个区间[left,right]
// # !求出所有序列在区间[left,right]的和的最大值
// # ROW,COL,q<=5e5 ROW*COL<=5e5

// # !(注意到ROW*COL<=5e5这个奇怪的条件)

// # COL<=ROW时
// # 询问的区间只有O(COL^2)种 `预处理所有查询`一共(ROW*COL*COL) 即O(5e5*sqrt(5e5))

// # COL>=ROW时
// # 每次询问都要(ROW)回答 一共O(q*ROW) 即O(q*sqrt(5e5))

package main

import (
	"bufio"
	"fmt"
	"os"
)

func P8572(grid [][]int, queries [][2]int) []int {
	ROW, COL := len(grid), len(grid[0])

	res := make([]int, len(queries))
	return res
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
