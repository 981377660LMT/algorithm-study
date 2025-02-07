// 有向图dfs三色标记找环
// https://leetcode.cn/problems/course-schedule/solutions/2992884/san-se-biao-ji-fa-pythonjavacgojsrust-by-pll7/
// https://leetcode.cn/problems/find-eventual-safe-states/solutions/916155/zhao-dao-zui-zhong-de-an-quan-zhuang-tai-yzfz/
//
// 白色（用 0 表示）：该节点尚未被访问；
// 灰色（用 1 表示）：该节点位于递归栈中，或者在某个环上；
// 黑色（用 2 表示）：该节点搜索完毕。

package main

func canFinish(numCourses int, prerequisites [][]int) bool {
	edgeAt := func(i int) (int, int) { return prerequisites[i][1], prerequisites[i][0] }
	return !HasCycleDirected(numCourses, len(prerequisites), edgeAt)
}

// https://leetcode.cn/problems/find-eventual-safe-states/description/
func eventualSafeNodes(graph [][]int) (res []int) {
	n := len(graph)
	color := make([]int, n)
	var safe func(int) bool
	safe = func(u int) bool {
		if color[u] > 0 {
			return color[u] == 2
		}
		color[u] = 1
		for _, v := range graph[u] {
			if !safe(v) {
				return false
			}
		}
		color[u] = 2
		return true
	}
	for i := 0; i < n; i++ {
		if safe(i) {
			res = append(res, i)
		}
	}
	return
}

// 有向图判断是否有环.
func HasCycleDirected(numVertex int, numEdge int, edgeAt func(int) (from, to int)) bool {
	graph := make([][]int, numVertex)
	for i := 0; i < numEdge; i++ {
		from, to := edgeAt(i)
		graph[from] = append(graph[from], to)
	}

	// 0：节点 x 尚未被访问到;
	// 1：节点 x 正在访问中，dfs(x) 尚未结束;
	// 2：节点 x 已经完全访问完毕，dfs(x) 已返回。
	color := make([]int, numVertex)

	var dfs func(int) bool
	dfs = func(u int) bool {
		color[u] = 1
		for _, v := range graph[u] {
			if color[v] == 1 || (color[v] == 0 && dfs(v)) {
				return true
			}
		}
		color[u] = 2
		return false
	}

	for i := 0; i < numVertex; i++ {
		if color[i] == 0 && dfs(i) {
			return true
		}
	}
	return false
}
