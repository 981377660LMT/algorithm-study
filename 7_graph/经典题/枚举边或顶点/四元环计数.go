package main

// https://github.dev/EndlessCheng/codeforces-go/tree/master/copypasta
// https://blog.csdn.net/weixin_43466755/article/details/112985722

// !时间复杂度O(E*3/2)
// !四元环唯一表示两条无向边+两条有向边 u-v->w u-x->w 且 rank(u)最小rank(w)最大
// !这样才能保证一个四元环不被反复计算。
func countCycle4(n int, edges [][2]int) (res int64) {
	adjList1 := make([][]int, n)
	deg := make([]int, n)
	less := func(u, v int) bool { return deg[u] < deg[v] || (deg[u] == deg[v] && u < v) }

	for _, e := range edges {
		u, v := e[0], e[1]
		adjList1[u] = append(adjList1[u], v)
		adjList1[v] = append(adjList1[v], u)
		deg[u]++
		deg[v]++
	}

	adjList2 := make([][]int, n)
	for cur, nexts := range adjList1 {
		for _, next := range nexts {
			if less(cur, next) {
				adjList2[cur] = append(adjList2[cur], next)
			}
		}
	}

	count := make([]int, n)
	for cur, nexts := range adjList1 {
		for _, next1 := range nexts {
			for _, next2 := range adjList2[next1] {
				if less(cur, next2) {
					count[cur]++
					res += int64(count[cur])
				}

			}

			for _, next1 := range nexts {
				for _, next2 := range adjList2[next1] {
					if less(cur, next2) {
						count[cur] = 0
					}
				}
			}
		}
	}

	return
}

func main() {
	// https://atcoder.jp/contests/abc260/tasks/abc260_f
}
