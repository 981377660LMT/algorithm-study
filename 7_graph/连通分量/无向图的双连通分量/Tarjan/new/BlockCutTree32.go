// BlockCutTree-圆方树点双
// 处理点双、割点相关逻辑直接使用圆方树即可
// https://oi-wiki.org/graph/block-forest/
// https://oi-wiki.org/graph/images/block-forest2.svg
// 如果原图连通，则「圆方树」才是一棵树，如果原图有 k 个连通分量，则它的圆方树也会形成 k 棵树形成的森林。
//
// 例如:
//
//	     原图        圆方树
//			0  —  1     0     1
//			 \	 /       \	 /
//			  \ /   =>     5 (block)
//			   2           |
//			   |           2 (原图割点,与block相连)
//			   3           |
//		                 4 (block)
//		                 |
//		                 3
//
// 原图的割点`至少`在两个不同的 v-BCC 中
// 原图不是割点的点都`只存在`于一个 v-BCC 中
// v-BCC 形成的子图内没有割点

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	yosupo()
}

// https://judge.yosupo.jp/problem/biconnected_components
func yosupo() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int32
	fmt.Fscan(in, &n, &m)
	graph := make([][]int32, n)
	for i := int32(0); i < m; i++ {
		var u, v int32
		fmt.Fscan(in, &u, &v)
		graph[u] = append(graph[u], v)
		graph[v] = append(graph[v], u)
	}

	tree := BlockCutTree32(graph)
	fmt.Fprintln(out, int32(len(tree))-n) // block数量=点双连通分量数量
	for i := n; i < int32(len(tree)); i++ {
		fmt.Fprint(out, len(tree[i]))
		for _, v := range tree[i] {
			fmt.Fprint(out, " ", v)
		}
		fmt.Fprintln(out)
	}
}

// 求出无向图的blockCutTree, 用于解决点双连通分量相关问题.
// 在blockCutTree中, 满足性质:
// !1.[0, n)为原图中的点, [n, n+n_block)为block.每一个点双连通分量连接、对应一个block.
// !2.割点 <=> [0, n)中满足degree>=2的点.
func BlockCutTree32(graph [][]int32) (tree [][]int32) {
	n := int32(len(graph))
	low := make([]int32, n)
	order := make([]int32, n)
	stack := make([]int32, 0, n)
	used := make([]bool, n)
	id := n
	now := int32(0)
	edges := [][2]int32{}

	var dfs func(int32, int32)
	dfs = func(cur, pre int32) {
		stack = append(stack, cur)
		used[cur] = true
		low[cur] = now
		order[cur] = now
		now++
		child := 0
		for _, to := range graph[cur] {
			if to == pre {
				continue
			}
			if !used[to] {
				child++
				s := len(stack)
				dfs(int32(to), cur)
				low[cur] = min32(low[cur], low[to])
				if (pre == -1 && child > 1) || (pre != -1 && low[to] >= order[cur]) {
					edges = append(edges, [2]int32{id, cur})
					for len(stack) > s {
						edges = append(edges, [2]int32{id, stack[len(stack)-1]})
						stack = stack[:len(stack)-1]
					}
					id++
				}
			} else {
				low[cur] = min32(low[cur], order[to])
			}
		}
	}

	for i := int32(0); i < n; i++ {
		if !used[i] {
			dfs(i, -1)
			for _, v := range stack {
				edges = append(edges, [2]int32{id, v})
			}
			id++
			stack = stack[:0]
		}
	}

	tree = make([][]int32, id)
	for _, e := range edges {
		u, v := e[0], e[1]
		tree[u] = append(tree[u], v)
		tree[v] = append(tree[v], u)
	}
	return
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func min32(a, b int32) int32 {
	if a <= b {
		return a
	}
	return b
}
