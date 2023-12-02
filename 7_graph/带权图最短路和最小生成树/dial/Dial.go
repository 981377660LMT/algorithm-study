package main

// https://leetcode.cn/problems/network-delay-time/
func networkDelayTime(times [][]int, n int, k int) int {
	adjList := make([][][2]int, n)
	maxWeight := 0
	for _, e := range times {
		u, v, w := e[0]-1, e[1]-1, e[2]
		adjList[u] = append(adjList[u], [2]int{v, w})
		maxWeight = max(maxWeight, w)
	}
	dist := Dial(adjList, k-1, maxWeight*n)
	maxDist := 0
	for _, d := range dist {
		if d == INF {
			return -1
		}
		maxDist = max(maxDist, d)
	}
	return maxDist
}

const INF int = 1e18

// Dial 算法求解带权图的单源最短路.
// !边权必须是非负整数.
// 时间复杂度 O(E + V * maxWeight).
// 特别地,maxWeight = 1 时,变为0-1 BFS.
func Dial(adjList [][][2]int, start int, maxWeight int) []int {
	n := len(adjList)
	dist := make([]int, n)
	for i := range dist {
		dist[i] = INF
	}
	dist[start] = 0
	buckets := make([][]int, maxWeight*n)
	buckets[0] = append(buckets[0], start)
	for d := 0; d < len(buckets); d++ {
		for len(buckets[d]) > 0 {
			v := buckets[d][len(buckets[d])-1]
			buckets[d] = buckets[d][:len(buckets[d])-1]
			if dist[v] < d {
				continue
			}
			for _, e := range adjList[v] {
				to, weight := e[0], e[1]
				if tmp := dist[v] + weight; tmp < dist[to] {
					dist[to] = tmp
					buckets[tmp] = append(buckets[tmp], to)
				}
			}
		}
	}
	return dist
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
