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
func CentroidDecomposition(n int, tree [][]Edge, root int) int {
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

	removed := make([]bool, n)
	subSize := make([]int, n)

	// !计算子树大小 getSize
	var getSize func(cur, pre int) int
	getSize = func(cur, pre int) int {
		sub := 1
		for _, e := range tree[cur] {
			if next := e.to; next != pre && !removed[next] {
				sub += getSize(next, cur)
			}
		}
		subSize[cur] = sub
		return sub
	}

	// !找重心 getCentroid
	var getCentroid func(int, int, int) int
	getCentroid = func(cur, pre, mid int) int {
		for _, e := range tree[cur] {
			if next := e.to; next != pre && !removed[next] && subSize[next] > mid {
				return getCentroid(next, cur, mid)
			}
		}
		return cur
	}

	var distToCentroid []int
	// !计算子树到重心的距离 collectDist
	var collectDist func(int, int, int)
	collectDist = func(cur, pre, dist int) {
		distToCentroid = append(distToCentroid, dist)
		for _, e := range tree[cur] {
			if next := e.to; next != pre && !removed[next] {
				collectDist(next, cur, dist+e.weight)
			}
		}
	}

	// !分治+业务逻辑 f
	var f func(cur, pre int) int
	f = func(cur, pre int) (res int) {
		getSize(cur, -1)
		c := getCentroid(cur, pre, subSize[cur]>>1)
		removed[c] = true
		defer func() { removed[c] = false }()

		// !统计`经过重心`的答案(为了确保经过重心，可以容斥原理计算)
		// 需要处理一些信息，例如每个点到重心的距离，等等
		// 业务逻辑写在这里...
		// ...

		// !统计按照重心分割后子树中的答案
		for _, e := range tree[c] {
			if next := e.to; !removed[next] {
				res += f(next, cur)
			}
		}

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
