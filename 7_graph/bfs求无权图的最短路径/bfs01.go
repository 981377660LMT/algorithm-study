// bfs01/01bfs

package main

const INF int = 1e18
const INF32 int32 = 1e9 + 10

func bfs(adjList [][]int, start int) (dist, pre []int) {
	n := len(adjList)
	dist = make([]int, n)
	pre = make([]int, n)
	for i := 0; i < n; i++ {
		dist[i] = INF
		pre[i] = -1
	}
	queue := []int{start}
	dist[start] = 0
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		for _, next := range adjList[cur] {
			cand := dist[cur] + next
			if cand < dist[next] {
				dist[next] = cand
				pre[next] = cur
				queue = append(queue, next)
			}
		}
	}

	return
}

func bfs01(adjList [][][2]int, start int) (dist, pre []int) {
	n := len(adjList)
	dist = make([]int, n)
	pre = make([]int, n)
	for i := 0; i < n; i++ {
		dist[i] = INF
		pre[i] = -1
	}
	queue := NewDeque[int](int32(n))
	queue.Append(start)
	dist[start] = 0
	for queue.Size() > 0 {
		cur := queue.PopLeft()
		for _, e := range adjList[cur] {
			next, weight := e[0], e[1]
			cand := dist[cur] + weight
			if cand < dist[next] {
				dist[next] = cand
				pre[next] = cur
				if weight == 0 {
					queue.AppendLeft(next)
				} else {
					queue.Append(next)
				}
			}
		}
	}
	return
}

func bfs0132(adjList [][][2]int32, start int32) (dist, pre []int32) {
	n := int32(len(adjList))
	dist = make([]int32, n)
	pre = make([]int32, n)
	for i := int32(0); i < n; i++ {
		dist[i] = INF32
		pre[i] = -1
	}
	queue := NewDeque[int32](n)
	queue.Append(start)
	dist[start] = 0
	for queue.Size() > 0 {
		cur := queue.PopLeft()
		for _, e := range adjList[cur] {
			next, weight := e[0], e[1]
			cand := dist[cur] + weight
			if cand < dist[next] {
				dist[next] = cand
				pre[next] = cur
				if weight == 0 {
					queue.AppendLeft(next)
				} else {
					queue.Append(next)
				}
			}
		}
	}
	return
}

func bfs01MultiStart(adjList [][][2]int, starts []int) (dist, pre, root []int) {
	n := int32(len(adjList))
	dist = make([]int, n)
	pre = make([]int, n)
	root = make([]int, n)
	for i := int32(0); i < n; i++ {
		dist[i] = INF
		pre[i] = -1
		root[i] = -1
	}
	queue := NewDeque[int](n)
	for _, start := range starts {
		queue.Append(start)
		dist[start] = 0
		root[start] = start
	}
	for queue.Size() > 0 {
		cur := queue.PopLeft()
		for _, e := range adjList[cur] {
			next, weight := e[0], e[1]
			cand := dist[cur] + weight
			if cand < dist[next] {
				dist[next] = cand
				pre[next] = cur
				root[next] = root[cur]
				if weight == 0 {
					queue.AppendLeft(next)
				} else {
					queue.Append(next)
				}
			}
		}
	}
	return
}

type Deque[D any] struct{ l, r []D }

func NewDeque[D any](cap int32) *Deque[D] {
	return &Deque[D]{make([]D, 0, 1+cap/2), make([]D, 0, 1+cap/2)}
}

func (q Deque[D]) Empty() bool {
	return len(q.l) == 0 && len(q.r) == 0
}

func (q Deque[D]) Size() int {
	return len(q.l) + len(q.r)
}

func (q *Deque[D]) AppendLeft(v D) {
	q.l = append(q.l, v)
}

func (q *Deque[D]) Append(v D) {
	q.r = append(q.r, v)
}

func (q *Deque[D]) PopLeft() (v D) {
	if len(q.l) > 0 {
		q.l, v = q.l[:len(q.l)-1], q.l[len(q.l)-1]
	} else {
		v, q.r = q.r[0], q.r[1:]
	}
	return
}

func (q *Deque[D]) Pop() (v D) {
	if len(q.r) > 0 {
		q.r, v = q.r[:len(q.r)-1], q.r[len(q.r)-1]
	} else {
		v, q.l = q.l[0], q.l[1:]
	}
	return
}

func (q Deque[D]) Front() D {
	if len(q.l) > 0 {
		return q.l[len(q.l)-1]
	}
	return q.r[0]
}

func (q Deque[D]) Back() D {
	if len(q.r) > 0 {
		return q.r[len(q.r)-1]
	}
	return q.l[0]
}

// 0 <= i < q.Size()
func (q Deque[D]) At(i int) D {
	if i < len(q.l) {
		return q.l[len(q.l)-1-i]
	}
	return q.r[i-len(q.l)]
}
