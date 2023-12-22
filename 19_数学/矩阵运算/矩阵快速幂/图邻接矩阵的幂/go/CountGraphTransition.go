package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int = 1e9 + 7

// https://atcoder.jp/contests/dp/tasks/dp_r
func main() {
	in, out := bufio.NewReader(os.Stdin), bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k int
	fmt.Fscan(in, &n, &k)
	adjMatrix := make([][]int, n)
	for i := range adjMatrix {
		adjMatrix[i] = make([]int, n)
		for j := range adjMatrix[i] {
			var v int
			fmt.Fscan(in, &v)
			if v == 1 {
				adjMatrix[i][j]++
			}
		}
	}

	transition := CountGraphTransition(adjMatrix, k, MOD)
	res := 0
	for _, row := range transition {
		for _, v := range row {
			res = (res + v) % MOD
		}
	}

	fmt.Println(res)

}

// !给定一个图的临接矩阵，求转移k次(每次转移可以从一个点到到任意一个点)后到达目的地的方案数(矩阵).
// !adjMatrix[i][j]表示从i到j的边数.
// 转移1次就是floyd.
func CountGraphTransition(adjMatrixCount [][]int, k int, mod int) [][]int {
	n := len(adjMatrixCount)
	countCopy := make([][]int, n)
	for i := range countCopy {
		countCopy[i] = make([]int, n)
		copy(countCopy[i], adjMatrixCount[i])
	}

	res := newMatrix(n, n, 0)
	for i := 0; i < n; i++ {
		res[i][i] = 1
	}

	for k > 0 {
		if k&1 == 1 {
			res = MatMul(res, countCopy, mod)
		}
		k >>= 1
		countCopy = MatMul(countCopy, countCopy, mod)
	}

	return res
}

// 转移的自定义函数.
// ed:内部的结合律为取max(Floyd).
func MatMul(m1, m2 [][]int, mod int) [][]int {
	res := newMatrix(len(m1), len(m2[0]), 0)
	for i := 0; i < len(m1); i++ {
		for k := 0; k < len(m2); k++ {
			for j := 0; j < len(m2[0]); j++ {
				res[i][j] = (res[i][j] + m1[i][k]*m2[k][j]) % mod
				if res[i][j] < 0 {
					res[i][j] += mod
				}
			}
		}
	}
	return res
}

func newMatrix(row, col int, fill int) [][]int {
	res := make([][]int, row)
	for i := range res {
		row := make([]int, col)
		for j := range row {
			row[j] = fill
		}
		res[i] = row
	}
	return res
}

func max(a, b int) int {
	if a >= b {
		return a
	}
	return b
}
