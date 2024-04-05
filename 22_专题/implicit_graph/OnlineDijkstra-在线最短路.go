package main

const INF int = 1e18

// 在线dijsktra.
//
//	不预先给出图，而是通过两个函数 setUsed 和 findUnused 来在线寻找边.
//	setUsed(u)：将 u 标记为已访问。
//	findUnused(u)：找到和 u 邻接的一个未访问过的点。如果不存在, 返回 `-1`。
func OnlineDijkstra(
	n int, start int,
	setUsed func(u int),
	findUnused func(cur int) (next int, weight int),
) (dist []int) {
	dist = make([]int, n)
	for i := range dist {
		dist[i] = INF
	}
	dist[start] = 0
	queue := NewHeap[[2]int](func(a, b [2]int) bool { return a[1] < b[1] }, [][2]int{{start, 0}}) // [u, dist[u]]
	setUsed(start)

	for queue.Len() > 0 {
		top := queue.Pop()
		cur, curDist := top[0], top[1]
		if dist[cur] < curDist {
			continue
		}
		for {
			next, weight := findUnused(cur)
			if next == -1 {
				break
			}
			setUsed(next)
			if cand := curDist + weight; cand < dist[next] {
				dist[next] = cand
				queue.Push([2]int{next, cand})
			}
		}
	}

	return
}

func NewHeap[H any](less func(a, b H) bool, nums []H) *Heap[H] {
	nums = append(nums[:0:0], nums...)
	heap := &Heap[H]{less: less, data: nums}
	heap.heapify()
	return heap
}

type Heap[H any] struct {
	data []H
	less func(a, b H) bool
}

func (h *Heap[H]) Push(value H) {
	h.data = append(h.data, value)
	h.pushUp(h.Len() - 1)
}

func (h *Heap[H]) Pop() (value H) {
	if h.Len() == 0 {
		panic("heap is empty")
	}
	value = h.data[0]
	h.data[0] = h.data[h.Len()-1]
	h.data = h.data[:h.Len()-1]
	h.pushDown(0)
	return
}

func (h *Heap[H]) Top() (value H) {
	value = h.data[0]
	return
}

func (h *Heap[H]) Len() int { return len(h.data) }

func (h *Heap[H]) heapify() {
	n := h.Len()
	for i := (n >> 1) - 1; i > -1; i-- {
		h.pushDown(i)
	}
}

func (h *Heap[H]) pushUp(root int) {
	for parent := (root - 1) >> 1; parent >= 0 && h.less(h.data[root], h.data[parent]); parent = (root - 1) >> 1 {
		h.data[root], h.data[parent] = h.data[parent], h.data[root]
		root = parent
	}
}

func (h *Heap[H]) pushDown(root int) {
	n := h.Len()
	for left := (root<<1 + 1); left < n; left = (root<<1 + 1) {
		right := left + 1
		minIndex := root
		if h.less(h.data[left], h.data[minIndex]) {
			minIndex = left
		}
		if right < n && h.less(h.data[right], h.data[minIndex]) {
			minIndex = right
		}
		if minIndex == root {
			return
		}
		h.data[root], h.data[minIndex] = h.data[minIndex], h.data[root]
		root = minIndex
	}
}
