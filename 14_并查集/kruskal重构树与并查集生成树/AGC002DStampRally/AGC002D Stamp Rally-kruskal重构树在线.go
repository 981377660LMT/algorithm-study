// AGC002D Stamp Rally-kruskal重构树在线

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)
	graph := make([][]int, n)
	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		u--
		v--
		graph[u] = append(graph[u], v)
		graph[v] = append(graph[v], u)
	}

	var q int
	fmt.Fscan(in, &q)
	queries := make([][3]int, q)
	for i := 0; i < q; i++ {
		var x, y, z int
		fmt.Fscan(in, &x, &y, &z)
		x--
		y--
		queries[i] = [3]int{x, y, z}
	}

	res := StampRally(graph, queries)
	for _, v := range res {
		fmt.Fprintln(out, v)
	}
}

// https://atcoder.jp/contests/agc002/tasks/agc002_d
// 一张连通图，q 次询问从两个点 x 和 y 出发，
// 希望经过的点数量等于 z（每个点可以重复经过，但是重复经过只计算一次）
// 求经过的边最大编号最小是多少。
//
// !根据性质：kruskal重构树中两个点u,v路径上的最大边权等于lca(u,v)的点权; kruskal重构树是一个大根堆.
// 因此可以二分编号，对于询问(x,y,z),我们倍增向上跳到点权大于当前二分值的位置，
// 然后再判断此时跳到节点的子树中的叶子节点数量是否大于等于z.
func StampRally(graph [][]int, queries [][3]int) []int {}
