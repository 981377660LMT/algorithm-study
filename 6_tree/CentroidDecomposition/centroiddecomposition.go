// !https://github.dev/EndlessCheng/codeforces-go/blob/cca30623b9ac0f3333348ca61b4894cd00b753cc/copypasta/graph_tree.go#L473
// !https://nyaannyaan.github.io/library/tree/centroid-decomposition.hpp

package centroiddecomposition

// 重心分解
// https://zhuanlan.zhihu.com/p/359209926
// https://oi-wiki.org/graph/tree-divide/
// https://www.luogu.com.cn/problem/P3806 树上距离为 k 的点对是否存在
// https://www.luogu.com.cn/problem/P4178 树上两点距离小于等于 k 的点对数量
// !应用：树上路径问题

const INF int = 1e18

type Edge struct{ to, weight int }

// 模板
// https://github.dev/EndlessCheng/codeforces-go/blob/cca30623b9ac0f3333348ca61b4894cd00b753cc/copypasta/graph_tree.go#L533
func usedCentroidDecomposition(n int, tree [][]Edge, root int) int {
	// !预处理parent和depth  build
	// parents := make([]int, n)
	// depths := make([]int, n)
	// var build func(cur, pre, dep int)
	// build = func(cur, pre, dep int) {
	// 	parents[cur] = pre
	// 	depths[cur] = dep
	// 	for _, e := range tree[cur] {
	// 		if e.next != pre {
	// 			build(e.next, cur, dep+1)
	// 		}
	// 	}
	// }
	// build(root, -1, 0)

	visited := make([]bool, n)
	subSize := make([]int, n)
	var distToCentroid []int

	// !计算子树大小 calSubSize
	var calSubSize func(cur, pre int) int
	calSubSize = func(cur, pre int) int {
		sub := 1
		for _, e := range tree[cur] {
			if next := e.to; next != pre && !visited[next] {
				sub += calSubSize(next, cur)
			}
		}
		subSize[cur] = sub
		return sub
	}

	// !找重心 findCentroid
	var findCentroid func(int, int, int) int
	findCentroid = func(cur, pre, mid int) int {
		for _, e := range tree[cur] {
			if next := e.to; next != pre && !visited[next] && subSize[next] > mid {
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
			if next := e.to; next != pre && !visited[next] {
				calDistToCentroid(next, cur, dist+e.weight)
			}
		}
	}

	// !分治+业务逻辑 f
	var f func(cur, pre int) int
	f = func(cur, pre int) (res int) {
		calSubSize(cur, -1)
		c := findCentroid(cur, pre, subSize[cur]>>1)
		visited[c] = true
		defer func() { visited[c] = false }()

		// !统计按照重心分割后子树中的答案
		for _, e := range tree[c] {
			if next := e.to; !visited[next] {
				res += f(next, cur)
			}
		}

		// !统计`经过重心`的答案(为了确保经过重心，可以容斥原理计算)
		// 需要处理一些信息，例如每个点到重心的距离，等等
		// 业务逻辑写在这里...
		// ...

		return
	}

	return f(root, -1)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
