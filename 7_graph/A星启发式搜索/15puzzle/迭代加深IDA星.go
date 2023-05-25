// https://gist.github.com/E869120/c921a0605edd9af8f33cbe9d20c2c96a
// !迭代加深IDA星求解 15puzzle 问题.

package main

import (
	"fmt"
)

func main() {
	grid := [][]int{{1, 2, 3, 4}, {5, 6, 7, 8}, {9, 10, 11, 12}, {15, 14, 13, 0}}
	fmt.Println(Solve15Puzzle(grid))
}

var DIR4 = [4][2]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}
var TARGET [16]int
var DIST [16][16]int // DIST[i*j][k]:从(i,j)到最终空白格k的距离

func init() {
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			TARGET[i*4+j] = 4 * (i*4 + j)
		}
	}

	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			for k := 1; k < 16; k++ {
				sx, sy := ((k+15)&15)/4, ((k+15)&15)&3
				DIST[i*4+j][k] = abs(i-sx) + abs(j-sy)
			}
		}
	}
}

// 15puzzle问题的最少移动次数.
//  0 表示空格,1-15 表示对应的15个数字方块.
//  如果无解则返回 -1.
func Solve15Puzzle(grid [][]int) int {
	flag := false

	// exp表示当前状态到最终状态的最短距离估计值.
	var dfs func(dep int, s, px, py, preDir, exp int) bool
	dfs = func(dep int, s, px, py, preDir, exp int) bool {
		if flag {
			return true
		}
		if dep == 0 {
			flag = true
			return true
		}

		pos1 := px*4 + py
		for i := 0; i < 4; i++ {
			dir := DIR4[i]
			nx, ny := px+dir[0], py+dir[1]
			if nx < 0 || ny < 0 || nx > 3 || ny > 3 || (i+2)&3 == preDir {
				continue
			}

			pos2 := nx*4 + ny
			hash1, hash2 := (s>>TARGET[pos1])&15, (s>>TARGET[pos2])&15
			nextExp := exp - DIST[pos1][hash1] - DIST[pos2][hash2] + DIST[pos1][hash2] + DIST[pos2][hash1]
			if nextExp > dep-1 {
				continue
			}

			nextS := s + ((hash2 - hash1) << TARGET[pos1]) + ((hash1 - hash2) << TARGET[pos2])
			if dfs(dep-1, nextS, nx, ny, i, nextExp) {
				return true
			}
		}

		return false
	}

	initS, initPx, initPy, initExp := 0, 0, 0, 0
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			v := grid[i][j]
			initS += v << TARGET[i*4+j] // 状态哈希值
			if v != 0 {
				initExp += DIST[i*4+j][v]
			} else {
				initPx, initPy = i, j
			}
		}
	}

	if initExp == 0 {
		return 0
	}

	// 最大移动次数为 80.
	for i := 1; i <= 80; i++ {
		if dfs(i, initS, initPx, initPy, -1, initExp) {
			return i
		}
	}
	return -1
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
