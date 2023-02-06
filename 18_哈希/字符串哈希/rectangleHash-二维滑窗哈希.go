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

type RollingHash = func(left int, right int) int

// 二维滑窗哈希
func rectangleHash(ords [][]int, mod, base1, base2, winRow, winCol int) [][]int {
	ROW, COL := len(ords), len(ords[0])
	hs := make([]RollingHash, 0, ROW)
	for i := 0; i < ROW; i++ {
		hs = append(hs, StringHasher(ords[i], mod, base1))
	}

	res := make([][]int, ROW-winRow+1)
	for i := 0; i < ROW-winRow+1; i++ {
		res[i] = make([]int, COL-winCol+1)
	}

	for j := 0; j+winCol <= COL; j++ {
		ts := make([]int, ROW)
		for i := 0; i < ROW; i++ {
			ts[i] = hs[i](j, j+winCol)
		}
		rh := StringHasher(ts, mod, base2)
		for i := 0; i+winRow <= ROW; i++ {
			res[i][j] = rh(i, i+winRow)
		}
	}

	return res
}

func StringHasher(ords []int, mod int, base int) func(left int, right int) int {
	prePow := make([]int, len(ords)+1)
	prePow[0] = 1
	preHash := make([]int, len(ords)+1)
	for i, v := range ords {
		prePow[i+1] = (prePow[i] * base) % mod
		preHash[i+1] = (preHash[i]*base + v) % mod
	}

	sliceHash := func(left, right int) int {
		if left >= right {
			return 0
		}
		return (preHash[right] - preHash[left]*prePow[right-left]%mod + mod) % mod
	}

	return sliceHash
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

	res := rectangleHash(ords1, MOD, BASE1, BASE2, winRow, winCol)
	target := rectangleHash(ords2, MOD, BASE1, BASE2, winRow, winCol)

	for i := 0; i < row1-winRow+1; i++ {
		for j := 0; j < col1-winCol+1; j++ {
			if res[i][j] == target[0][0] {
				fmt.Fprintln(out, i, j)
			}
		}
	}
}
