package main

// 树的直径个数.
func CountDiameter(tree [][]int, start int) (diameter, diameterCount int) {
	var dfs func(cur, pre int) (int, int)
	dfs = func(cur, pre int) (int, int) {
		maxDepth, count := 0, 1
		for _, next := range tree[cur] {
			if next != pre {
				nextDepth, nextCount := dfs(next, cur)
				if tmp := maxDepth + nextDepth; tmp > diameter {
					diameter, diameterCount = tmp, count*nextCount
				} else if tmp == diameter {
					diameterCount += count * nextCount
				}
				if nextDepth > maxDepth {
					maxDepth, count = nextDepth, nextCount
				} else if nextDepth == maxDepth {
					count += nextCount
				}
			}
		}
		return maxDepth + 1, count
	}
	dfs(start, -1)
	return
}
