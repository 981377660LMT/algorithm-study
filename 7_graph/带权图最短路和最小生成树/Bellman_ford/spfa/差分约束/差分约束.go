// !如果所有边权非负,可以把spfa换成dijkstra

package main

import (
	"fmt"
)

const INF int = 1e18

type DualShortestPath struct {
	n      int
	g      [][][2]int
	min    bool
	hasNeg bool
}

func NewDualShortestPath(n int, min bool) *DualShortestPath {
	return &DualShortestPath{
		n:   n,
		g:   make([][][2]int, n),
		min: min,
	}
}

// f(j) <= f(i) + w
func (d *DualShortestPath) AddEdge(i, j, w int) {
	if d.min {
		d.g[i] = append(d.g[i], [2]int{j, w})
	} else {
		d.g[j] = append(d.g[j], [2]int{i, w})
	}
	if w < 0 {
		d.hasNeg = true
	}
}

// 求 `f(i) - f(0)` 的最小值/最大值, 并检测是否有负环/正环
func (d *DualShortestPath) Run() (dist []int, ok bool) {
	if d.min {
		return d.spfaMin()
	}
	if !d.hasNeg {
		return d.dijkMax()
	}
	return d.spfaMax()
}

func (d *DualShortestPath) spfaMin() (dist []int, ok bool) {
	dist = make([]int, d.n)
	queue := NewDeque(d.n)
	count := make([]int, d.n)
	inQueue := make([]bool, d.n)
	for i := 0; i < d.n; i++ {
		queue.Append(i)
		inQueue[i] = true
		count[i] = 1
	}
	for queue.Size() > 0 {
		cur := queue.PopLeft()
		inQueue[cur] = false
		for _, e := range d.g[cur] {
			next, weight := e[0], e[1]
			cand := dist[cur] + weight
			if cand < dist[next] {
				dist[next] = cand
				if !inQueue[next] {
					count[next]++
					if count[next] >= d.n+1 {
						return nil, false
					}
					inQueue[next] = true
					queue.AppendLeft(next)
				}
			}
		}
	}

	for i := 0; i < d.n; i++ {
		dist[i] = -dist[i]
	}
	ok = true
	return
}

func (d *DualShortestPath) spfaMax() (dist []int, ok bool) {
	dist = make([]int, d.n)
	inQueue := make([]bool, d.n)
	count := make([]int, d.n)
	for i := 0; i < d.n; i++ {
		dist[i] = INF
	}

	queue := NewDeque(d.n)
	queue.Append(0)
	dist[0] = 0
	inQueue[0] = true
	count[0] = 1
	for queue.Size() > 0 {
		cur := queue.PopLeft()
		inQueue[cur] = false
		for _, e := range d.g[cur] {
			next, weight := e[0], e[1]
			cand := dist[cur] + weight
			if cand < dist[next] {
				dist[next] = cand
				if !inQueue[next] {
					count[next]++
					if count[next] >= d.n+1 {
						return nil, false
					}
					inQueue[next] = true
					queue.AppendLeft(next)
				}
			}
		}
	}

	ok = true
	return
}

func (dsp *DualShortestPath) dijkMax() (dist []int, ok bool) {
	dist = make([]int, dsp.n)
	for i := 0; i < dsp.n; i++ {
		dist[i] = INF
	}
	pq := NewHeap(func(a, b H) bool { return a[0] < b[0] }, nil)
	pq.Push([2]int{0, 0})
	dist[0] = 0
	for pq.Len() > 0 {
		cur := pq.Pop()
		v, d := cur[1], cur[0]
		if dist[v] < d {
			continue
		}
		for _, e := range dsp.g[v] {
			next, weight := e[0], e[1]
			cand := dist[v] + weight
			if cand < dist[next] {
				dist[next] = cand
				pq.Push([2]int{cand, next})
			}
		}
	}
	ok = true
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

type H = [2]int // (dist,node)

func NewHeap(less func(a, b H) bool, nums []H) *Heap {
	nums = append(nums[:0:0], nums...)
	heap := &Heap{less: less, data: nums}
	heap.heapify()
	return heap
}

type Heap struct {
	data []H
	less func(a, b H) bool
}

func (h *Heap) Push(value H) {
	h.data = append(h.data, value)
	h.pushUp(h.Len() - 1)
}

func (h *Heap) Pop() (value H) {
	if h.Len() == 0 {
		panic("heap is empty")
	}

	value = h.data[0]
	h.data[0] = h.data[h.Len()-1]
	h.data = h.data[:h.Len()-1]
	h.pushDown(0)
	return
}

func (h *Heap) Peek() (value H) {
	if h.Len() == 0 {
		panic("heap is empty")
	}
	value = h.data[0]
	return
}

func (h *Heap) Len() int { return len(h.data) }

func (h *Heap) heapify() {
	n := h.Len()
	for i := (n >> 1) - 1; i > -1; i-- {
		h.pushDown(i)
	}
}

func (h *Heap) pushUp(root int) {
	for parent := (root - 1) >> 1; parent >= 0 && h.less(h.data[root], h.data[parent]); parent = (root - 1) >> 1 {
		h.data[root], h.data[parent] = h.data[parent], h.data[root]
		root = parent
	}
}

func (h *Heap) pushDown(root int) {
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

//
//
//
func main() {
	// 	// https://www.luogu.com.cn/problem/P5960
	// 	// 求任意一组满足这个不等式组的解
	// 	// 如果有多组解，请输出任意一组，无解请输出 NO。
	// 	in := bufio.NewReader(os.Stdin)
	// 	out := bufio.NewWriter(os.Stdout)
	// 	defer out.Flush()

	// 	var n, m int
	// 	fmt.Fscan(in, &n, &m)
	// 	DSP := NewDualShortestPath(n+10, true)
	// 	for i := 0; i < m; i++ {
	// 		// f(u)-f(v)<x
	// 		var u, v, x int // 1<=u,v<=n
	// 		fmt.Fscan(in, &u, &v, &x)
	// 		DSP.AddEdge(u, v, x)
	// 	}

	// 	res, ok := DSP.Run()
	// 	if !ok {
	// 		fmt.Fprintln(out, "NO")
	// 		return
	// 	}
	// 	for i := 1; i <= n; i++ {
	// 		fmt.Fprint(out, res[i], " ")
	// 	}

	n := 10
	limits := [][]int{{1, 4, 2}, {3, 6, 2}, {10, 10, 1}}
	D := NewDualShortestPath(n+10, false)
	for _, v := range limits {
		D.AddEdge(v[1], v[0]-1, v[2])
	}
	for i := 1; i <= n+1; i++ {
		D.AddEdge(i-1, i, 0)
		D.AddEdge(i, i-1, 1)
	}
	dist, ok := D.Run()
	if !ok {
		fmt.Println("No solution")
	}
	fmt.Println(dist[n])
}
