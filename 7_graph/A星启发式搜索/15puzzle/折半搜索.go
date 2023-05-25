// https://gist.github.com/E869120/707c539013346aaa881092aa6ad06f44
// !折半搜索/meet in the middle 求解 15puzzle 问题.
// !时间复杂度为 O(3^(m/2)),m 为最大移动次数.

package main

import (
	"fmt"
	"sort"
)

func main() {
	grid := [][]int{{1, 2, 3, 4}, {5, 6, 7, 8}, {9, 10, 11, 12}, {13, 14, 0, 15}}
	fmt.Println(Solve15Puzzle(grid))
}

var DIR4 = [4][2]int{{-1, 0}, {0, 1}, {1, 0}, {0, -1}}

// 15puzzle问题的最少移动次数.
//  0 表示空格,1-15 表示对应的15个数字方块.
//  如果无解则返回 -1.
func Solve15Puzzle(grid [][]int) int {
	target := [16]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 0}
	initial := State{}
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			// target[i*4+j] = (i*4 + j + 1) & 15 // % 16
			initial[i*4+j] = grid[i][j]
		}
	}

	var left, right []State

	var dfs func(dep int, s State, preDir int, record *[]State)
	dfs = func(dep int, s State, preDir int, record *[]State) {
		if dep == 0 {
			*record = append(*record, s)
			return
		}

		sx, sy := 0, 0
		for i := 0; i < 4; i++ {
			for j := 0; j < 4; j++ {
				if s[i*4+j] == 0 {
					sx, sy = i, j
					break
				}
			}
		}

		for i := 0; i < 4; i++ {
			d4 := DIR4[i]
			nx, ny := sx+d4[0], sy+d4[1]
			if nx < 0 || ny < 0 || nx > 3 || ny > 3 || (i+2)&3 == preDir { // 与上一步相反的方向
				continue
			}
			nextS := s
			nextS[sx*4+sy], nextS[nx*4+ny] = nextS[nx*4+ny], nextS[sx*4+sy]
			dfs(dep-1, nextS, i, record)
		}
	}

	// 最多移动 move 次.
	// 求出从起点move/2步内能到达的状态集合 W1,从终点move-move/2步内能到达的状态集合 W2.
	// 二分判断 W1[i] = W2[j] 是否存在.
	solve := func(maxMove int) bool {
		left = left[:0]
		right = right[:0]
		dfs(maxMove/2, initial, -1, &left)
		dfs(maxMove-maxMove/2, target, -1, &right)
		sort.Slice(left, func(i, j int) bool { return StateLess(left[i], left[j]) })
		sort.Slice(right, func(i, j int) bool { return StateLess(right[i], right[j]) })

		// 判断 W1[i] = W2[j] 是否存在.
		for i := 0; i < len(left); i++ {
			pos := sort.Search(len(right), func(j int) bool { return StateLess(left[i], right[j]) }) - 1
			if pos >= 0 && left[i] == right[pos] {
				return true
			}
		}
		return false
	}

	// 最大移动次数为 80.
	for i := 0; i <= 80; i++ {
		if solve(i) {
			return i
		}
	}
	return -1
}

type State = [16]int

func StateLess(a1, a2 State) bool {
	for i := 0; i < 16; i++ {
		if a1[i] != a2[i] {
			return a1[i] < a2[i]
		}
	}
	return false
}
