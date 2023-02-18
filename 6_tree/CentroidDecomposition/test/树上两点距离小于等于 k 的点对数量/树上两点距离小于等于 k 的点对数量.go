package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// https://www.luogu.com.cn/problem/P4178
// 给定一棵 n 个节点的树，每条边有边权，求出树上两点距离小于等于 k 的点对数量。
// !n<=4e4 O(n*logn*logn)

const INF int = 1e18

type Edge struct{ next, weight int }

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)

	tree := make([][]Edge, n)
	for i := 0; i < n-1; i++ {
		var u, v, w int
		fmt.Fscan(in, &u, &v, &w)
		u, v = u-1, v-1
		tree[u] = append(tree[u], Edge{v, w})
		tree[v] = append(tree[v], Edge{u, w})
	}

	var k int
	fmt.Fscan(in, &k)

	res := 树上两点距离小于等于k的点对数量(n, tree, 0, k)
	fmt.Fprintln(out, res)
}

// https://github.dev/EndlessCheng/codeforces-go/blob/cca30623b9ac0f3333348ca61b4894cd00b753cc/copypasta/graph_tree.go#L533
func 树上两点距离小于等于k的点对数量(n int, tree [][]Edge, root int, K int) int {
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
		// !容斥原理统计`经过重心`的长度等于k的路径数(点对数)
		// 所有点对数 - 子树内的点对数(不经过重心) = 经过重心的点对数
		allDist := []int{0} // 0方便统计包含重心的情况
		for _, e := range tree[centroid] {
			if next := e.next; !visited[next] {
				distToCentroid = []int{}
				calDistToCentroid(next, centroid, e.weight)
				res -= countPair(distToCentroid, K) // !减去不经过重心的点对数(不合法的路径)
				allDist = append(allDist, distToCentroid...)
			}
		}

		res += countPair(allDist, K) // !合在一起统计所有的对数
		return
	}

	return dfs(root, -1)
}

// 排序+头尾双指针求(dist[i]+dist[j]<=k,i<j)的对数
func countPair(dist []int, k int) int {
	sort.Ints(dist)
	res, left, right := 0, 0, len(dist)-1
	for left < right {
		if dist[left]+dist[right] <= k {
			res += right - left
			left++
		} else {
			right--
		}
	}
	return res
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
