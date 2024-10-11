package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
)

func main() {
	aizu_alds1_14_c()
}

// https://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=ALDS1_14_C
// !检测 二维矩阵中是否存在子矩阵与给定的特征矩阵相同,输出左上角坐标
// ROW,COL<=1e3
func aizu_alds1_14_c() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var row1, col1 int32
	fmt.Fscan(in, &row1, &col1)
	grid1 := make([]string, row1)
	for i := int32(0); i < row1; i++ {
		var s string
		fmt.Fscan(in, &s)
		grid1[i] = s
	}

	var winRow, winCol int32
	fmt.Fscan(in, &winRow, &winCol)
	grid2 := make([]string, winRow)
	for i := int32(0); i < winRow; i++ {
		var s string
		fmt.Fscan(in, &s)
		grid2[i] = s
	}

	H := NewRollingHash2D(0, 0)
	table1 := H.Build(row1, col1, func(i, j int32) uint64 { return uint64(grid1[i][j]) })
	table2 := H.Build(winRow, winCol, func(i, j int32) uint64 { return uint64(grid2[i][j]) })

	target := H.Query(table2, 0, winRow, 0, winCol)

	for i := int32(0); i < row1-winRow+1; i++ {
		for j := int32(0); j < col1-winCol+1; j++ {
			if H.Query(table1, i, i+winRow, j, j+winCol) == target {
				fmt.Fprintln(out, i, j)
			}
		}
	}
}

func demo() {
	arr := [][]uint64{{0, 1, 0}, {0, 1, 0}, {1, 0, 1}}

	R := NewRollingHash2D(0, 0)
	table := R.Build(int32(len(arr)), int32(len(arr[0])), func(i, j int32) uint64 { return arr[i][j] })
	fmt.Println(R.Query(table, 0, 1, 0, 2))
	fmt.Println(R.Query(table, 2, 3, 1, 3))
}

const (
	mod61  uint64 = (1 << 61) - 1
	mask30 uint64 = (1 << 30) - 1
	mask31 uint64 = (1 << 31) - 1
	mask61 uint64 = mod61
)

type RollingHash2D struct {
	b1, b2     uint64
	pow1, pow2 []uint64
}

// base: 0 表示随机生成
func NewRollingHash2D(base1, base2 uint64) *RollingHash2D {
	for base1 == 0 {
		base1 = rand.Uint64() % mod61 // rng61
	}
	for base2 == 0 {
		base2 = rand.Uint64() % mod61
	}
	return &RollingHash2D{b1: base1, b2: base2, pow1: []uint64{1}, pow2: []uint64{1}}
}

func (rh *RollingHash2D) Build(r, c int32, f func(i, j int32) uint64) (table [][]uint64) {
	table = make([][]uint64, r+1)
	for i := int32(0); i < r+1; i++ {
		table[i] = make([]uint64, c+1)
	}
	for i := int32(0); i < r; i++ {
		for j := int32(0); j < c; j++ {
			table[i+1][j+1] = rh.add(rh.mul(table[i+1][j], rh.b2), rh.mod(1+f(i, j)))
		}
		for j := int32(0); j < c+1; j++ {
			table[i+1][j] = rh.add(table[i+1][j], rh.mul(table[i][j], rh.b1))
		}
	}

	// expand
	for int32(len(rh.pow1)) <= r {
		rh.pow1 = append(rh.pow1, rh.mul(rh.pow1[len(rh.pow1)-1], rh.b1))
	}
	for int32(len(rh.pow2)) <= c {
		rh.pow2 = append(rh.pow2, rh.mul(rh.pow2[len(rh.pow2)-1], rh.b2))
	}
	return
}

// [xStart,xEnd) x [yStart,yEnd)
func (rh *RollingHash2D) Query(table [][]uint64, xStart, xEnd, yStart, yEnd int32) uint64 {
	a := rh.sub(table[xEnd][yEnd], rh.mul(table[xEnd][yStart], rh.pow2[yEnd-yStart]))
	b := rh.sub(table[xStart][yEnd], rh.mul(table[xStart][yStart], rh.pow2[yEnd-yStart]))
	return rh.sub(a, rh.mul(b, rh.pow1[xEnd-xStart]))
}

// x % (2^61-1)
func (rh *RollingHash2D) mod(x uint64) uint64 {
	xu := x >> 61
	xd := x & mask61
	res := xu + xd
	if res >= mod61 {
		res -= mod61
	}
	return res
}

// a*b % (2^61-1)
func (rh *RollingHash2D) mul(a, b uint64) uint64 {
	au := a >> 31
	ad := a & mask31
	bu := b >> 31
	bd := b & mask31
	mid := ad*bu + au*bd
	midu := mid >> 30
	midd := mid & mask30
	return rh.mod(au*bu<<1 + midu + (midd << 31) + ad*bd)
}

// a,b: modint61
func (rh *RollingHash2D) add(a, b uint64) uint64 {
	res := a + b
	if res >= mod61 {
		res -= mod61
	}
	return res
}

// a,b: modint61
func (rh *RollingHash2D) sub(a, b uint64) uint64 {
	res := a - b
	if res >= mod61 {
		res += mod61
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

func min32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func max32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}
