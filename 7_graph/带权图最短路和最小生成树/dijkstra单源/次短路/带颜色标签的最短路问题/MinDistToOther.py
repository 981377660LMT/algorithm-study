# // !给定一个无负权的带权图和一个点集，对点集内的每个点V，求V到点集内其他点距离的最小值,以及到V最近的点是谁.
# // 换一种描述方法，给定一张图，有一些黑点和白点，对每个黑点，求出它到其他黑点的最近距离.
# // 按照points中点的顺序返回答案.
# //
# // !以黑点集合为源做一次多源次短路dij，然后每个黑点的次短路就是答案.
# // 注意次短路的出发点不为能自己.
# // 类似 abc245G - Foreign Friends-简洁写法.go
# func MinDistToOther(adjList [][][2]int, points []int) (dist []int, nearest []int) {
# 	n := len(adjList)
# 	dist = make([]int, n)
# 	source1, source2 := make([]int, n), make([]int, n)
# 	for i := 0; i < n; i++ {
# 		dist[i] = INF
# 		source1[i], source2[i] = -1, -1
# 	}

# 	pq := NewHeap(func(a, b H) bool { return a.dist < b.dist }, nil)
# 	for _, v := range points {
# 		pq.Push(H{dist: 0, node: v, source: v})
# 	}

# 	for pq.Len() > 0 {
# 		item := pq.Pop()
# 		curDist, cur, curSource := item.dist, item.node, item.source
# 		if curSource == source1[cur] || curSource == source2[cur] {
# 			continue
# 		}
# 		if source1[cur] == -1 {
# 			source1[cur] = curSource
# 		} else if source2[cur] == -1 {
# 			source2[cur] = curSource
# 		} else {
# 			continue
# 		}

# 		if curSource != cur { // 出发点不为自己时，更新距离
# 			dist[cur] = min(dist[cur], curDist)
# 		}
# 		for _, e := range adjList[cur] {
# 			next, weight := e[0], e[1]
# 			nextDist := curDist + weight
# 			pq.Push(H{nextDist, next, curSource})
# 		}
# 	}

# 	nearest = source2
# 	for i, v := range points {
# 		dist[i] = dist[v]
# 		nearest[i] = nearest[v]
# 	}
# 	dist = dist[:len(points)]
# 	nearest = nearest[:len(points)]
# 	return
# }
