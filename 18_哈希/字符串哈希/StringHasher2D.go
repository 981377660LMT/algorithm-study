// rectangleHash
// 二维滑窗哈希

package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int = 1e9 + 7
const BASE1 int = 1777771
const BASE2 int = 1e8 + 7

// 二维哈希
//  query(x1, y1, x2, y2 int) int:返回(x1,y1)和(x2,y2)组成的闭区间的哈希值
func stringHash2D(ords [][]int, mod, baseX, baseY int) func(x1, y1, x2, y2 int) int {
	ROW, COL := len(ords), len(ords[0])

	// genPow
	powX, powY := make([]int, ROW+1), make([]int, COL+1)
	powX[0], powY[0] = 1, 1
	for i := 1; i <= ROW; i++ {
		powX[i] = powX[i-1] * baseX % mod
	}
	for i := 1; i <= COL; i++ {
		powY[i] = powY[i-1] * baseY % mod
	}

	// genHash
	hash := make([][]int, ROW)
	for i := 0; i < ROW; i++ {
		hash[i] = make([]int, COL)
	}

	get := func(x, y int) int {
		if x < 0 || y < 0 || x >= ROW || y >= COL {
			return 0
		}
		return hash[x][y]
	}

	for i := 0; i < ROW; i++ {
		for j := 0; j < COL; j++ {
			hash[i][j] = (get(i, j-1)*baseY%mod + ords[i][j]) % mod
		}
	}

	for i := 0; i < ROW; i++ {
		for j := 0; j < COL; j++ {
			hash[i][j] = (get(i-1, j)*baseX%mod + hash[i][j]) % mod
		}
	}

	query := func(x1, y1, x2, y2 int) int {
		x2, y2 = x2+1, y2+1
		res := get(x2-1, y2-1)
		res = (res - get(x1-1, y2-1)*powX[x2-x1]%mod + mod) % mod
		res = (res - get(x2-1, y1-1)*powY[y2-y1]%mod + mod) % mod
		res = (res + (get(x1-1, y1-1)*powX[x2-x1]%mod)*powY[y2-y1]%mod) % mod
		return res
	}

	return query
}

// https://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=ALDS1_14_C
// !检测 二维矩阵中是否存在子矩阵与给定的特征矩阵相同,输出左上角坐标
// ROW,COL<=1e3
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var row1, col1 int
	fmt.Fscan(in, &row1, &col1)

	ords1 := make([][]int, row1)
	for i := 0; i < row1; i++ {
		ords1[i] = make([]int, col1)
		var s string
		fmt.Fscan(in, &s)
		for j, v := range s {
			ords1[i][j] = int(v)
		}
	}

	var winRow, winCol int
	fmt.Fscan(in, &winRow, &winCol)
	ords2 := make([][]int, winRow)
	for i := 0; i < winRow; i++ {
		ords2[i] = make([]int, winCol)
		var s string
		fmt.Fscan(in, &s)
		for j, v := range s {
			ords2[i][j] = int(v)
		}
	}

	query1 := stringHash2D(ords1, MOD, BASE1, BASE2)
	target := stringHash2D(ords2, MOD, BASE1, BASE2)(0, 0, winRow-1, winCol-1)

	for i := 0; i < row1-winRow+1; i++ {
		for j := 0; j < col1-winCol+1; j++ {
			if query1(i, j, i+winRow-1, j+winCol-1) == target {
				fmt.Fprintln(out, i, j)
			}
		}
	}
}
