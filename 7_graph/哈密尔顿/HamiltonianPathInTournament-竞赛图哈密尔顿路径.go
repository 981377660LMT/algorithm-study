// !每对顶点之间都有一条边相连的有向图称为竞赛图(Tournament)
// 通过在`无向完全图`中为每个边缘分配方向而获得的有向图.
// Tournament图可以用来表示比赛结果.

// !任何有限数量n个顶点的竞赛图都包含一个哈密尔顿路径.
// 注意Tournament是完全图。
// https://yukicoder.me/problems/no/2085

package main

import (
	"fmt"
	"math/rand"
)

func main() {
	// https://maspypy.github.io/library/test/mytest/tournament.test.cpp
	gen := func(n int) [][]bool {
		mat := make([][]bool, n)
		for i := range mat {
			mat[i] = make([]bool, n)
		}
		for i := 0; i < n; i++ {
			for j := 0; j < i; j++ {
				b := rand.Intn(2) == 0
				if b {
					mat[i][j] = true
				}
				if !b {
					mat[j][i] = true
				}
			}
		}
		return mat
	}

	sum := func(use []bool) int {
		s := 0
		for _, v := range use {
			if v {
				s++
			}
		}
		return s
	}

	for round := 0; round < 1; round++ {
		for n := 5000; n <= 5000; n++ {
			adjMatrix := gen(n)
			isConnected := func(i, j int) bool { return adjMatrix[i][j] }
			path := HamiltonianPathInTournament(n, isConnected)
			used := make([]bool, n)
			for _, x := range path {
				used[x] = true
			}
			if len(path) != n {
				panic("len(P) != n")
			}
			if sum(used) != n {
				panic("sum(use) != n")
			}
			for i := 0; i < n-1; i++ {
				a, b := path[i], path[i+1]
				if !adjMatrix[a][b] {
					panic("!G[a][b]")
				}
			}
		}
	}

	fmt.Println("OK")
}

// 返回竞赛图中的哈密尔顿路径.任何有限数量n个顶点的竞赛图都包含一个哈密尔顿路径.
//  check: i 和 j 之间是否有边.
func HamiltonianPathInTournament(n int, isConnected func(i, j int) bool) []int {
	var dfs func(int, int) []int
	dfs = func(left, right int) []int {
		if right == left+1 {
			return []int{left}
		}
		mid := (left + right) / 2
		leftPath := dfs(left, mid)
		rightPath := dfs(mid, right)
		path := make([]int, 0, right-left)
		i, j := 0, 0
		for len(path) < right-left {
			if i == len(leftPath) {
				path = append(path, rightPath[j])
				j++
			} else if j == len(rightPath) {
				path = append(path, leftPath[i])
				i++
			} else {
				if isConnected(leftPath[i], rightPath[j]) {
					path = append(path, leftPath[i])
					i++
				} else {
					path = append(path, rightPath[j])
					j++
				}
			}
		}
		return path
	}
	return dfs(0, n)
}
