// EnumeratingC3

// !枚举无向图中的所有三元环(a,b,c) a<b<c
//  O(E^1.5)
// https://kopricky.github.io/code/Graph/EnumeratingC3.html

package main

import (
	"fmt"
	"math"
)

func minTrioDegree(n int, edges [][]int) int {
	deg := make([]int, n)
	for i := range edges {
		edges[i][0]--
		edges[i][1]--
		u, v := edges[i][0], edges[i][1]
		deg[u]++
		deg[v]++
	}

	res := math.MaxInt32
	enumerateC3(n, edges, func(u, v, w int) {
		res = min(res, deg[u]+deg[v]+deg[w]-6)
		fmt.Println(u, v, w)
	})
	if res == math.MaxInt32 {
		return -1
	}
	return res
}

// 枚举所有的三元环,传递回调函数减少空间开销
func enumerateC3(n int, edges [][]int, cb func(u, v, w int)) {
	threshold := 0
	adjList := make([][]int, n)
	for _, e := range edges {
		u, v := e[0], e[1]
		// 无向边定向
		if u < v {
			adjList[u] = append(adjList[u], v)
		} else {
			adjList[v] = append(adjList[v], u)
		}
		threshold++
	}

	threshold = int(math.Sqrt(float64(threshold / 4)))
	memo := make([][][2]int, n)
	visited := make([]bool, n)

	processHighDegree := func() {
		for i := 0; i < n; i++ {
			if len(adjList[i]) <= threshold {
				continue
			}
			for _, u := range adjList[i] {
				visited[u] = true
			}
			for _, u := range adjList[i] {
				for _, v := range adjList[u] {
					if visited[v] {
						cb(i, u, v)
					}
				}
			}
			for _, u := range adjList[i] {
				visited[u] = false
			}
		}
	}

	processLowDegree := func() {
		for i := 0; i < n; i++ {
			if len(adjList[i]) > threshold {
				continue
			}
			for _, u := range adjList[i] {
				for _, v := range adjList[i] {
					if v > u {
						memo[u] = append(memo[u], [2]int{i, v})
					}
				}
			}
		}
		for i := 0; i < n; i++ {
			for _, u := range adjList[i] {
				visited[u] = true
			}
			for j := 0; j < len(memo[i]); j++ {
				a, b := memo[i][j][0], memo[i][j][1]
				if visited[b] {
					cb(a, i, b)
				}
			}
			for _, u := range adjList[i] {
				visited[u] = false
			}
		}
	}

	processHighDegree()
	processLowDegree()
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
