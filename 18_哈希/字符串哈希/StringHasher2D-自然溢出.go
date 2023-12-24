// rectangleHash
// 二维滑窗哈希

package main

import (
	"bufio"
	"fmt"
	"os"
)

// https://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=ALDS1_14_C
// !检测 二维矩阵中是否存在子矩阵与给定的特征矩阵相同,输出左上角坐标
// ROW,COL<=1e3
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var row1, col1 int
	fmt.Fscan(in, &row1, &col1)
	grid1 := make([]string, row1)
	for i := 0; i < row1; i++ {
		var s string
		fmt.Fscan(in, &s)
		grid1[i] = s
	}

	var winRow, winCol int
	fmt.Fscan(in, &winRow, &winCol)
	grid2 := make([]string, winRow)
	for i := 0; i < winRow; i++ {
		var s string
		fmt.Fscan(in, &s)
		grid2[i] = s
	}

	H := NewRollingHash2D(BASE1, BASE2)
	table1 := H.Build(grid1)
	table2 := H.Build(grid2)

	target := H.Query(table2, 0, winRow, 0, winCol)

	for i := 0; i < row1-winRow+1; i++ {
		for j := 0; j < col1-winCol+1; j++ {
			if H.Query(table1, i, i+winRow, j, j+winCol) == target {
				fmt.Fprintln(out, i, j)
			}
		}
	}
}

const BASE1 uint = 1777771 // 131/13331/1713302033171(回文素数)
const BASE2 uint = 13331

type S = string

type RollingHash2D struct {
	base1, base2 uint
	power1       []uint
	power2       []uint
}

func NewRollingHash2D(base1, base2 uint) *RollingHash2D {
	return &RollingHash2D{
		base1:  base1,
		base2:  base2,
		power1: []uint{1},
		power2: []uint{1},
	}
}

func (r *RollingHash2D) Build(grid []S) (hashTable [][]uint) {
	row, col := len(grid), len(grid[0])
	hashTable = make([][]uint, row+1)
	for i := 0; i < row+1; i++ {
		hashTable[i] = make([]uint, col+1)
	}
	for i := 0; i < row; i++ {
		for j := 0; j < col; j++ {
			hashTable[i+1][j+1] = hashTable[i+1][j]*r.base2 + uint(grid[i][j])
		}
		for j := 0; j < col+1; j++ {
			hashTable[i+1][j] += hashTable[i][j] * r.base1
		}
	}
	// expand
	for len(r.power1) <= row {
		r.power1 = append(r.power1, r.power1[len(r.power1)-1]*r.base1)
	}
	for len(r.power2) <= col {
		r.power2 = append(r.power2, r.power2[len(r.power2)-1]*r.base2)
	}
	return hashTable
}

// [xStart,xEnd) x [yStart,yEnd)
func (r *RollingHash2D) Query(hashTable [][]uint, xStart, xEnd, yStart, yEnd int) uint {
	a := hashTable[xEnd][yEnd] - hashTable[xEnd][yStart]*r.power2[yEnd-yStart]
	b := hashTable[xStart][yEnd] - hashTable[xStart][yStart]*r.power2[yEnd-yStart]
	return a - b*r.power1[xEnd-xStart]
}
