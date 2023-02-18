package bfs

const INF int = 1e18

type Edge struct{ to, weight int }

func Bfs01(n int, adjList [][]Edge, start int) (dist []int) {
	dist = make([]int, n)
	for i := range dist {
		dist[i] = INF
	}

	queue := Deque{}
	queue.Append(start)
	dist[start] = 0
	for !queue.Empty() {
		cur := queue.PopLeft()
		for _, e := range adjList[cur] {
			next, cost := e.to, e.weight
			if dist[next] > dist[cur]+cost {
				dist[next] = dist[cur] + cost
				if cost == 0 {
					queue.AppendLeft(next)
				} else {
					queue.Append(next)
				}
			}
		}
	}

	return
}

func Bfs01Point(n int, adjList [][]Edge, start, target int) int {
	dist := make([]int, n)
	for i := range dist {
		dist[i] = INF
	}

	queue := Deque{}
	queue.Append(start)
	dist[start] = 0
	for !queue.Empty() {
		cur := queue.PopLeft()
		if cur == target {
			return dist[cur]
		}
		for _, e := range adjList[cur] {
			next, cost := e.to, e.weight
			if dist[next] > dist[cur]+cost {
				dist[next] = dist[cur] + cost
				if cost == 0 {
					queue.AppendLeft(next)
				} else {
					queue.Append(next)
				}
			}
		}
	}

	return INF
}

type D = int
type Deque struct{ l, r []D }

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
