package main

import (
	"bufio"
	"fmt"
	"os"
)

// https://www.luogu.com.cn/problem/P3806
// 给定一棵有 n 个点的树，q次询问树上距离为 k 的点对是否存在。
// !n<=1e4 q<=100

const INF int = 1e18

type Edge struct{ next, weight int }

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)

	tree := make([][]Edge, n)
	for i := 0; i < n-1; i++ {
		var u, v, w int
		fmt.Fscan(in, &u, &v, &w)
		u, v = u-1, v-1
		tree[u] = append(tree[u], Edge{v, w})
		tree[v] = append(tree[v], Edge{u, w})
	}

	for i := 0; i < q; i++ {
		var k int
		fmt.Fscan(in, &k)
		res := 树上两点距离等于k的点对数量(n, tree, 0, k)
		if res > 0 {
			fmt.Fprintln(out, "AYE")
		} else {
			fmt.Fprintln(out, "NAY")
		}
	}
}

// https://github.dev/EndlessCheng/codeforces-go/blob/cca30623b9ac0f3333348ca61b4894cd00b753cc/copypasta/graph_tree.go#L533
func 树上两点距离等于k的点对数量(n int, tree [][]Edge, root int, upperDist int) int {
	visited := make([]bool, n)
	subSizes := make([]int, n)
	var distToCentroid []int

	// !计算子树大小 calSubSize
	var calSubSize func(cur, pre int) int
	calSubSize = func(cur, pre int) int {
		subSize := 1
		for _, e := range tree[cur] {
			if next := e.next; next != pre && !visited[next] {
				subSize += calSubSize(next, cur)
			}
		}
		subSizes[cur] = subSize
		return subSize
	}

	// !找重心 findCentroid
	var findCentroid func(int, int, int) int
	findCentroid = func(cur, pre, mid int) int {
		for _, e := range tree[cur] {
			if next := e.next; next != pre && !visited[next] && subSizes[next] > mid {
				return findCentroid(next, cur, mid)
			}
		}
		return cur
	}

	// !计算子树到重心的距离 calDistToCentroid
	var calDistToCentroid func(int, int, int)
	calDistToCentroid = func(cur, pre, dist int) {
		distToCentroid = append(distToCentroid, dist)
		for _, e := range tree[cur] {
			if next := e.next; next != pre && !visited[next] {
				calDistToCentroid(next, cur, dist+e.weight)
			}
		}
	}

	// !分治+业务逻辑 dfs
	var dfs func(cur, pre int) int
	dfs = func(cur, pre int) (res int) {
		calSubSize(cur, pre)
		centroid := findCentroid(cur, pre, subSizes[cur]>>1)
		visited[centroid] = true
		defer func() { visited[centroid] = false }()

		// !统计按照重心分割后子树中的答案
		for _, e := range tree[centroid] {
			if next := e.next; !visited[next] {
				res += dfs(next, cur)
			}
		}

		// 业务逻辑
		// !哈希表统计`经过重心`的长度等于k的路径数(点对数)
		counter := map[int]int{0: 1} // 统计经过重心的路径长度
		for _, e := range tree[centroid] {
			if next := e.next; !visited[next] {
				distToCentroid = []int{}
				calDistToCentroid(next, centroid, e.weight)
				for _, dist := range distToCentroid {
					need := upperDist - dist
					res += counter[need]
				}
				for _, dist := range distToCentroid {
					counter[dist]++
				}
			}
		}

		return
	}

	return dfs(root, -1)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
