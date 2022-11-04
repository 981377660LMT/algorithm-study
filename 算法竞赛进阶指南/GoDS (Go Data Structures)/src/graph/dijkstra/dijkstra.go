// dijkstra
package dijkstra

import "cmnx/src/heap/binaryheap"

const inf int = 1e18

type Edge struct{ next, weight int }

func Dijkstra(n int, adjList [][]Edge, start int) (dist []int) {
	dist = make([]int, n)
	for i := range dist {
		dist[i] = inf
	}
	dist[start] = 0

	type pair struct{ node, dist int }
	pq := binaryheap.NewBinaryHeap(func(a, b interface{}) int {
		return a.(pair).dist - b.(pair).dist
	}, []interface{}{pair{start, 0}})

	for pq.Len() > 0 {
		curNode := pq.Pop().(pair)
		cur, curDist := curNode.node, curNode.dist
		if curDist > dist[cur] {
			continue
		}

		for _, edge := range adjList[cur] {
			next, weight := edge.next, edge.weight
			if cand := curDist + weight; cand < dist[next] {
				dist[next] = cand
				pq.Push(pair{next, cand})
			}
		}
	}

	return
}
