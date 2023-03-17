// bfs01/01bfs

package main

func main() {

}

const INF int = 1e18

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

func bfs01(adjList [][]int, start int) (dist, pre []int) {
	n := len(adjList)
	dist = make([]int, n)
	pre = make([]int, n)
	for i := 0; i < n; i++ {
		dist[i] = INF
		pre[i] = -1
	}
	queue := NewDeque(n)
	queue.Append(start)
	dist[start] = 0
	for queue.Size() > 0 {
		cur := queue.PopLeft()
		for _, next := range adjList[cur] {
			cand := dist[cur] + next
			if cand < dist[next] {
				dist[next] = cand
				pre[next] = cur
				if next == 0 {
					queue.AppendLeft(next)
				} else {
					queue.Append(next)
				}
			}
		}
	}
	return
}

func bfs01MultiStart(adjList [][]int, starts []int) (dist, pre, root []int) {
	n := len(adjList)
	dist = make([]int, n)
	pre = make([]int, n)
	root = make([]int, n)
	for i := 0; i < n; i++ {
		dist[i] = INF
		pre[i] = -1
		root[i] = -1
	}
	queue := NewDeque(n)
	for _, start := range starts {
		queue.Append(start)
		dist[start] = 0
		root[start] = start
	}
	for queue.Size() > 0 {
		cur := queue.PopLeft()
		for _, next := range adjList[cur] {
			cand := dist[cur] + next
			if cand < dist[next] {
				dist[next] = cand
				pre[next] = cur
				root[next] = root[cur]
				if next == 0 {
					queue.AppendLeft(next)
				} else {
					queue.Append(next)
				}
			}
		}
	}
	return
}

type D = int
type Deque struct{ l, r []D }

func NewDeque(cap int) *Deque { return &Deque{make([]D, 0, 1+cap/2), make([]D, 0, 1+cap/2)} }

func (q Deque) Empty() bool {
	return len(q.l) == 0 && len(q.r) == 0
}

func (q Deque) Size() int {
	return len(q.l) + len(q.r)
}

func (q *Deque) AppendLeft(v D) {
	q.l = append(q.l, v)
}

func (q *Deque) Append(v D) {
	q.r = append(q.r, v)
}

func (q *Deque) PopLeft() (v D) {
	if len(q.l) > 0 {
		q.l, v = q.l[:len(q.l)-1], q.l[len(q.l)-1]
	} else {
		v, q.r = q.r[0], q.r[1:]
	}
	return
}

func (q *Deque) Pop() (v D) {
	if len(q.r) > 0 {
		q.r, v = q.r[:len(q.r)-1], q.r[len(q.r)-1]
	} else {
		v, q.l = q.l[0], q.l[1:]
	}
	return
}

func (q Deque) Front() D {
	if len(q.l) > 0 {
		return q.l[len(q.l)-1]
	}
	return q.r[0]
}

func (q Deque) Back() D {
	if len(q.r) > 0 {
		return q.r[len(q.r)-1]
	}
	return q.l[0]
}

// 0 <= i < q.Size()
func (q Deque) At(i int) D {
	if i < len(q.l) {
		return q.l[len(q.l)-1-i]
	}
	return q.r[i-len(q.l)]
}
