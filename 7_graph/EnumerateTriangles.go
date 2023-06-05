// EnumeratingC3

// 求出无向图中的所有三元环(a,b,c) a<b<c
//  O(E^1.5)
// https://kopricky.github.io/code/Graph/EnumeratingC3.html

package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

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
	memo := make([][]int, n)
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
						memo[u] = append(memo[u], i*(1<<32)+v)
					}
				}
			}
		}
		for i := 0; i < n; i++ {
			for _, u := range adjList[i] {
				visited[u] = true
			}
			for j := 0; j < len(memo[i]); j++ {
				hash := memo[i][j]
				a, b := hash>>32, hash&0xffffffff
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

func main() {

	// https://judge.yosupo.jp/problem/enumerate_triangles

	const MOD int = 998244353

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)
	values := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &values[i])
	}

	edges := make([][]int, m)
	for i := 0; i < m; i++ {
		edges[i] = make([]int, 2)
		fmt.Fscan(in, &edges[i][0], &edges[i][1])
	}

	res := 0
	enumerateC3(n, edges, func(u, v, w int) {
		res += values[u] * values[v] % MOD * values[w] % MOD
		res %= MOD
	})
	fmt.Println(res)
}
