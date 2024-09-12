package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	yuki1640()
}

// No.1640 簡単な色塗り
// https://yukicoder.me/problems/no/1640
// 给定n个点和n对条件，每个条件形如(u,v)，表示选择u或v中的一个点。
// 问能否选出所有点.
// 如果能，输出选择方案.
func yuki1640() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int32
	fmt.Fscan(in, &n)
	pairs := make([][2]int32, n)
	for i := range pairs {
		fmt.Fscan(in, &pairs[i][0], &pairs[i][1])
		pairs[i][0]--
		pairs[i][1]--
	}

	selects, ok := SelectOneFromEachPairRestore(n, pairs)
	if !ok {
		fmt.Fprintln(out, "No")
		return
	}
	fmt.Fprintln(out, "Yes")
	for _, v := range selects {
		fmt.Fprintln(out, v+1)
	}
}

// SelectOneFromEachEdgeRestore.
// !bfs，优先选择度数小的没选过的点.
func SelectOneFromEachPairRestore(n int32, pairs [][2]int32) (selects []int32, ok bool) {
	graph, deg := make([][][3]int32, n), make([]int32, n)
	for i := int32(0); i < n; i++ {
		u, v := pairs[i][0], pairs[i][1]
		graph[u] = append(graph[u], [3]int32{u, v, i})
		graph[v] = append(graph[v], [3]int32{v, u, i})
		deg[u]++
		deg[v]++
	}

	type E = [2]int32 // (deg,v)
	pq := NewHeap(func(a, b E) bool { return a[0] < b[0] }, nil)
	for i := int32(0); i < n; i++ {
		pq.Push(E{deg[i], i})
	}

	m := int32(len(pairs))
	selects = make([]int32, m)
	for i := int32(0); i < m; i++ {
		selects[i] = -1
	}
	visited := make([]bool, n)
	for pq.Len() > 0 {
		item := pq.Pop()
		curD, cur := item[0], item[1]
		if deg[cur] != curD {
			continue
		}
		if visited[cur] {
			continue
		}
		for _, e := range graph[cur] {
			from, to, id := e[0], e[1], e[2]
			if selects[id] != -1 {
				continue
			}
			selects[id] = cur
			visited[cur] = true
			deg[from]--
			deg[to]--
			pq.Push(E{deg[from], from})
			pq.Push(E{deg[to], to})
			break
		}
	}

	sum := int32(0)
	for _, v := range visited {
		if v {
			sum++
		}
	}
	ok = sum == m
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
