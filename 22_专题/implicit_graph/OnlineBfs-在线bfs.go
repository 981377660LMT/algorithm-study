// OnlineBfs-在线bfs
// https://leetcode.cn/problems/minimum-reverse-operations/solution/python-zai-xian-bfs-jie-jue-bian-shu-hen-y58m/

package main

const INF int = 1e18

// 在线bfs.
//   不预先给出图，而是通过两个函数 setUsed 和 findUnused 来在线寻找边.
//   setUsed(u)：将 u 标记为已访问。
//   findUnused(u)：找到和 u 邻接的一个未访问过的点。如果不存在, 返回 `-1`。
// https://leetcode.cn/problems/minimum-reverse-operations/solution/python-zai-xian-bfs-jie-jue-bian-shu-hen-y58m/
func OnlineBfs(
	n int, start int,
	setUsed func(u int), findUnused func(cur int) (next int),
) (dist []int) {
	dist = make([]int, n)
	for i := range dist {
		dist[i] = INF
	}
	dist[start] = 0
	queue := []int{start}
	setUsed(start)

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		for {
			next := findUnused(cur)
			if next == -1 {
				break
			}
			dist[next] = dist[cur] + 1 // weight
			queue = append(queue, next)
			setUsed(next)
		}
	}

	return
}
