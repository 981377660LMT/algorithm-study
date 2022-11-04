package kruskal

import "sort"

type Edge struct {
	u, v   int
	weight int
	id     int // 某些题目需要
}

// https://github.dev/EndlessCheng/codeforces-go/blob/016834c19c4289ae5999988585474174224f47e2/copypasta/graph.go#L1474

func Kruskal(n int, sortedEdges []Edge) int {
}

// 曼哈顿距离最小生成树 O(nlogn)
// LC1584 https://leetcode-cn.com/problems/min-cost-to-connect-all-points/
// 做法见官方题解 https://leetcode-cn.com/problems/min-cost-to-connect-all-points/solution/lian-jie-suo-you-dian-de-zui-xiao-fei-yo-kcx7/
func (*graph) manhattanMST(points []struct{ x, y, i int }, abs func(int) int) (mst int) {
	n := len(points)
	// 读入时把 points 加上下标

	type edge struct{ v, w, dis int }
	edges := []edge{}

	build := func() {
		sort.Slice(points, func(i, j int) bool { a, b := points[i], points[j]; return a.x < b.x || a.x == b.x && a.y < b.y })

		// 离散化 y-x
		type pair struct{ v, i int }
		ps := make([]pair, n)
		for i, p := range points {
			ps[i] = pair{p.y - p.x, i}
		}
		sort.Slice(ps, func(i, j int) bool { return ps[i].v < ps[j].v })
		kth := make([]int, n)
		k := 1
		kth[ps[0].i] = k
		for i := 1; i < n; i++ {
			if ps[i].v != ps[i-1].v {
				k++
			}
			kth[ps[i].i] = k
		}

		const inf int = 2e9
		tree := make([]int, k+1)
		idRec := make([]int, k+1)
		for i := range tree {
			tree[i] = inf
			idRec[i] = -1
		}
		update := func(pos, val, id int) {
			for ; pos > 0; pos &= pos - 1 {
				if val < tree[pos] {
					tree[pos] = id
					idRec[pos] = id
				}
			}
		}
		query := func(pos int) int {
			minVal, minID := inf, -1
			for ; pos < len(tree); pos += pos & -pos {
				if tree[pos] < minVal {
					minVal = tree[pos]
					minID = idRec[pos]
				}
			}
			return minID
		}

		for i := n - 1; i >= 0; i-- {
			p := points[i]
			pos := kth[i]
			if j := query(pos); j != -1 {
				q := points[j]
				dis := abs(p.x-q.x) + abs(p.y-q.y)
				edges = append(edges, edge{p.i, q.i, dis})
			}
			update(pos, p.x+p.y, i)
		}
	}
	build()
	for i := range points {
		points[i].x, points[i].y = points[i].y, points[i].x
	}
	build()
	for i := range points {
		points[i].x = -points[i].x
	}
	build()
	for i := range points {
		points[i].x, points[i].y = points[i].y, points[i].x
	}
	build()

	sort.Slice(edges, func(i, j int) bool { return edges[i].dis < edges[j].dis })

	uf := NewUnionFind(n)
	left := n - 1
	for _, e := range edges {
		if uf.Merge(e.v, e.w) >= 0 {
			mst += e.dis // int64
			left--
			if left == 0 {
				break
			}
		}
	}
	return
}
