// 3419. 图的最大边权的最小值
// https://leetcode.cn/problems/minimize-the-maximum-edge-weight-of-graph/description/
//
// 给你两个整数 n 和 threshold ，同时给你一个 n 个节点的 有向 带权图，节点编号为 0 到 n - 1 。这个图用 二维 整数数组 edges 表示，其中 edges[i] = [Ai, Bi, Wi] 表示节点 Ai 到节点 Bi 之间有一条边权为 Wi的有向边。
//
// 你需要从这个图中删除一些边（也可能 不 删除任何边），使得这个图满足以下条件：
//
// 所有其他节点都可以到达节点 0 。
// 图中剩余边的 最大 边权值尽可能小。
// 每个节点都 至多 有 threshold 条出去的边。
// 请你返回删除必要的边后，最大 边权的 最小值 为多少。如果无法满足所有的条件，请你返回 -1 。
//
// 1. dfs生成树/bfs生成树 + 二分答案.
// 2. threshold 是多余条件.
// !3. 不重建visited 数组 => 使用时间戳，不重建图 => 遍历时判断边是否可用.

package main

func minMaxWeight(n int, edges [][]int, threshold int) int {
	if len(edges) < n-1 {
		return -1
	}

	type edge struct{ to, weight int }
	revG := make([][]edge, n)
	maxWeight := 0
	for _, e := range edges {
		u, v, w := e[0], e[1], e[2]
		revG[v] = append(revG[v], edge{to: u, weight: w})
		maxWeight = max(maxWeight, w)
	}

	visited := make([]int, n)
	timeStamp := 0
	check := func(mid int) bool {
		timeStamp++
		queue := []int{0}
		visited[0] = timeStamp
		count := 1
		for len(queue) > 0 {
			cur := queue[0]
			queue = queue[1:]
			for _, e := range revG[cur] {
				if e.weight <= mid && visited[e.to] != timeStamp {
					visited[e.to] = timeStamp
					queue = append(queue, e.to)
					count++
					if count == n {
						return true
					}
				}
			}
		}
		return count == n
	}

	ok := false
	left, right := 0, maxWeight
	for left <= right {
		mid := (left + right) / 2
		if check(mid) {
			right = mid - 1
			ok = true
		} else {
			left = mid + 1
		}
	}
	if !ok {
		return -1
	}
	return left
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
