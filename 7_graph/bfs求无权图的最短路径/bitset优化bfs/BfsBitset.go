// https://maspypy.github.io/library/graph/shortest_path/bfs_bitset.hpp
// bitset优化bfs O(V^2>>6)
// 应用:稠密图无权最短路 例如2000*2000的邻接矩阵
// 密グラフの重みなし最短路問題
// 01 行列を vc<bitset> の形で渡す
// O(N^2/w)
// 参考：(4000,4000) を 4000 回で 2 秒以内？

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {
	P6328我是仙人掌()
}

func yukicoder1400() {
	// https://yukicoder.me/problems/no/1400
	// V<=2000,D<=1e18
	// 如果当前在第i行，Matrix[i][j]为1,就可以移动到第j列的任意一行
	// 问是否能做到:从任意一个点出发，经过D个回合后，可以到达任意一个点
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var V, D int
	fmt.Fscan(in, &V, &D)
	graph := make([]_BS, V+V)
	for i := range graph {
		graph[i] = _NewBS(V + V)
	}
	for i := 0; i < V; i++ {
		var S string
		fmt.Fscan(in, &S)
		for j := 0; j < V; j++ {
			if S[j] == '0' {
				continue
			}
			graph[i].Set(j + V)
			graph[i+V].Set(j)
		}
	}
	for s := 0; s < V; s++ {
		dist := BfsBitset(graph, s)
		if D%2 == 0 {
			for t := 0; t < V; t++ {
				if dist[t] == -1 || dist[t] > D {
					fmt.Fprintln(out, "No")
					return
				}
			}
		}
		if D%2 == 1 {
			for t := V; t < V+V; t++ {
				if dist[t] == -1 || dist[t] > D {
					fmt.Fprintln(out, "No")
					return
				}
			}
		}
	}
	fmt.Fprintln(out, "Yes")
}

// func P6328我是仙人掌() {
// 	// https://www.luogu.com.cn/problem/P6328
// 	// n<=1000,m,q<=1e5
// 	// 每个查询给定(x,dist),查询图中与x距离不超过dist的点的个数
// 	// 无向图，求最短路
// 	in := bufio.NewReader(os.Stdin)
// 	out := bufio.NewWriter(os.Stdout)
// 	defer out.Flush()

// 	var n, m, q int
// 	fmt.Fscan(in, &n, &m, &q)
// 	adjMatrix := make([]_BS, n)
// 	for i := range adjMatrix {
// 		adjMatrix[i] = _NewBS(n)
// 	}
// 	for i := 0; i < m; i++ {
// 		var u, v int
// 		fmt.Fscan(in, &u, &v)
// 		u--
// 		v--
// 		adjMatrix[u].Set(v)
// 		adjMatrix[v].Set(u)
// 	}

// 	res := make([]int, q)
// 	queries := make([][2]int, q)
// 	queryGroup := make([][]int, n)
// 	for qi := range queries {
// 		var x, dist int
// 		fmt.Fscan(in, &x, &dist)
// 		x--
// 		queries[qi] = [2]int{x, dist}
// 		queryGroup[x] = append(queryGroup[x], qi)
// 	}

// 	for s := 0; s < n; s++ {
// 		if len(queryGroup[s]) == 0 {
// 			continue
// 		}
// 		curQueries := queryGroup[s]
// 		sort.Slice(curQueries, func(i, j int) bool {
// 			return queries[curQueries[i]][1] < queries[curQueries[j]][1]
// 		})
// 		dist := BfsBitset(adjMatrix, s)
// 		sort.Ints(dist)
// 		ptr := 0
// 		for _, qi := range curQueries {
// 			for ptr < len(dist) && dist[ptr] <= queries[qi][1] {
// 				ptr++
// 			}
// 			res[qi] = ptr
// 		}
// 	}

// 	for _, v := range res {
// 		fmt.Fprintln(out, v)
// 	}
// }

// O(n^2/w)
// graph: 01邻接矩阵.
func BfsBitset(graph []_BS, start int) []int {
	n := len(graph)
	dist := make([]int, n)
	for i := range dist {
		dist[i] = -1
	}
	unused := _NewBS(n)
	for i := 0; i < n; i++ {
		unused.Set(i)
	}
	queue := _NewBS(n)
	queue.Set(start)
	d := 0
	for {
		p := queue.Index1()
		if p >= n {
			break
		}
		next := _NewBS(n)
		for p < n {
			dist[p] = d
			unused.Reset(p)
			next.IOr(graph[p])
			p = queue.Next1(p + 1)
		}
		queue = next.And(unused)
		d++
	}
	return dist
}

type _BS []uint64

func _NewBS(n int) _BS    { return make(_BS, n>>6+1) }  // (n+64-1)>>6
func (b _BS) Set(p int)   { b[p>>6] |= 1 << (p & 63) }  // 置 1
func (b _BS) Reset(p int) { b[p>>6] &^= 1 << (p & 63) } // 置 0

// 返回第一个 1 的下标，若不存在则返回一个不小于 n 的位置.
func (b _BS) Index1() int {
	for i, v := range b {
		if v != 0 {
			return i<<6 | bits.TrailingZeros64(v)
		}
	}
	return len(b) << 6
}

// 返回下标 >= p 的第一个 1 的下标，若不存在则返回一个不小于 n 的位置
func (b _BS) Next1(p int) int {
	if i := p >> 6; i < len(b) {
		v := b[i] & (^uint64(0) << (p & 63))
		if v != 0 {
			return i<<6 | bits.TrailingZeros64(v)
		}
		for i++; i < len(b); i++ {
			if b[i] != 0 {
				return i<<6 | bits.TrailingZeros64(b[i])
			}
		}
	}
	return len(b) << 6
}

func (b _BS) IOr(c _BS) {
	for i, v := range c {
		b[i] |= v
	}
}

func (b _BS) And(c _BS) _BS {
	res := make(_BS, len(b))
	for i, v := range b {
		res[i] = v & c[i]
	}
	return res
}
