// 3425. 最长特殊路径
// https://leetcode.cn/problems/longest-special-path/description/
// 给你一棵根节点为节点 0 的无向树，树中有 n 个节点，编号为 0 到 n - 1 ，这棵树通过一个长度为 n - 1 的二维数组 edges 表示，其中 edges[i] = [ui, vi, lengthi] 表示节点 ui 和 vi 之间有一条长度为 lengthi 的边。
// 同时给你一个整数数组 nums ，其中 nums[i] 表示节点 i 的值。
// 特殊路径 指的是树中一条从祖先节点 往下 到后代节点且经过节点的值 互不相同 的路径。
// 注意 ，一条路径可以开始和结束于同一节点。
// 请你返回一个长度为 2 的数组 result ，其中 result[0] 是 最长 特殊路径的 长度 ，result[1] 是所有 最长特殊路径中的 最少 节点数目。
//
// !无重复字符的最长子串树上版本.

package main

func longestSpecialPath(edges [][]int, nums []int) []int {
	n := len(nums)

	type edge struct{ to, weight int }
	tree := make([][]edge, n)
	for _, e := range edges {
		u, v, w := e[0], e[1], e[2]
		tree[u] = append(tree[u], edge{v, w})
		tree[v] = append(tree[v], edge{u, w})
	}

	distSum := make([]int, n)
	max_ := maxs(nums...)
	preDepth := make([]int, max_+1) // 记录每个值上一次出现的深度
	for i := range preDepth {
		preDepth[i] = -1
	}

	resLen := -1
	resCount := n + 1
	ptr := 0 // 滑窗左端点
	var dfs func(cur, pre, depth, dist int)
	dfs = func(cur, pre, depth, dist int) {
		distSum[depth] = dist

		v := nums[cur]
		preD := preDepth[v]
		prePtr := ptr
		if preD >= ptr {
			ptr = preD + 1
		}
		preDepth[v] = depth

		curLen := distSum[depth] - distSum[ptr]
		curCount := depth - ptr + 1
		if curLen > resLen {
			resLen = curLen
			resCount = curCount
		} else if curLen == resLen && curCount < resCount {
			resCount = curCount
		}

		for _, next := range tree[cur] {
			next, weight := next.to, next.weight
			if next == pre {
				continue
			}
			dfs(next, cur, depth+1, dist+weight)
		}

		preDepth[v] = preD
		ptr = prePtr
	}

	dfs(0, -1, 0, 0)
	return []int{resLen, resCount}
}

func maxs(nums ...int) int {
	res := nums[0]
	for _, num := range nums {
		if num > res {
			res = num
		}
	}
	return res
}
