// F - Perfect Matching on a Tree (树上最大匹配/最小匹配问题)
// https://atcoder.jp/contests/abc362/tasks/abc362_f
// 给定n个顶点的树，我们要将所有顶点两两配对(一共floor(n/2)对)，其中顶点u和顶点v配对的权值为dist(u,v)
// 现在要求将所有顶点两两配对，且要求计算最大的可能权重和最小的可能权重
//
// !树上匹配问题
// https://taodaling.github.io/blog/2019/09/10/%E6%A0%91%E4%B8%8A%E7%AE%97%E6%B3%95/#heading-%E6%A0%91%E4%B8%8A%E5%8C%B9%E9%85%8D%E9%97%AE%E9%A2%98
//
// !设分组集合为P，则代价W(P) = dist(u,v) = depth(u) + depth(v) - 2 * depth(lca(u,v)) | (u,v) ∈ P
// 权重最大/小时，∑depth(lca(u,v))| (u,v) ∈ P 需要最小/大
//
// 1.在求最大的情况，我们可以对树进行DFS，并尽可能将LCA设置为当前遍历的顶点.
//   !实际上这里要让权重最大，我们可以直接找到树的重心即可，将重心作为根，这时候所有配对顶点的lca都是根
// 2.而在求最小的情况，我们同样对树进行DFS，但是尽可能在遍历子树的时候就将子树中的顶点配对.

package main

import (
	"bufio"
	"fmt"
	"os"
)

// 输出匹配方案
// https://atcoder.jp/contests/abc362/tasks/abc362_f
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int32
	fmt.Fscan(in, &n)
	tree := make([][]int32, n)
	for i := int32(0); i < n-1; i++ {
		var u, v int32
		fmt.Fscan(in, &u, &v)
		u--
		v--
		tree[u] = append(tree[u], v)
		tree[v] = append(tree[v], u)
	}

	matching := []int32{}
	var dfs func(cur, pre int32)
	dfs = func(cur, pre int32) {
		matching = append(matching, cur)
		for _, next := range tree[cur] {
			if next == pre {
				continue
			}
			dfs(next, cur)
		}
	}

	centroid := FindCentroids32(n, tree, 0)[0]
	for _, next := range tree[centroid] {
		dfs(next, centroid)
	}
	if n%2 == 0 {
		matching = append(matching, centroid)
	}

	m := len(matching) / 2
	for i := 0; i < m; i++ {
		fmt.Fprintln(out, matching[i]+1, matching[i+m]+1) // matching[i]和matching[i+m]位于重心的两个不同子树中
	}
}

func FindCentroids32(n int32, tree [][]int32, root int32) (centroids []int32) {
	weight := make([]int32, n)
	subSize := make([]int32, n)
	var dfs func(cur, pre int32)
	dfs = func(cur, pre int32) {
		subSize[cur] = 1
		for _, next := range tree[cur] {
			if next == pre {
				continue
			}
			dfs(next, cur)
			subSize[cur] += subSize[next]
			weight[cur] = max32(weight[cur], subSize[next])
		}
		weight[cur] = max32(weight[cur], n-subSize[cur])
		if weight[cur] <= n/2 {
			centroids = append(centroids, cur)
		}
	}

	dfs(root, -1)
	return
}

func max32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}
