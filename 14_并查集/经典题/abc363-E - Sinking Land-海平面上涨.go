// abc363-E - Sinking Land-海平面上涨淹没岛屿, 接雨水
// https://atcoder.jp/contests/abc363/tasks/abc363_e
//
// 方格岛，四面环海。
// 给出岛的高度，每年海平面上升1。
// 问y年的每一年，没被淹的岛的数量。
//
//
// !用优先队列维护这些危险岛的高度，当有新的危险岛被淹时，新增周围变成危险岛的高度，模拟即可。

package main

import (
	"bufio"
	"fmt"
	"os"
)

var dir4 = [4][2]int{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var H, W, Y int
	fmt.Fscan(in, &H, &W, &Y)
	A := make([][]int, H)
	for i := range A {
		A[i] = make([]int, W)
		for j := range A[i] {
			fmt.Fscan(in, &A[i][j])
		}
	}

	pq := NewHeap(func(a, b [3]int) bool { return a[0] < b[0] }, nil)
	visited := make([][]bool, H)
	for i := range visited {
		visited[i] = make([]bool, W)
	}

	push := func(x, y int) {
		if 0 <= x && x < H && 0 <= y && y < W && !visited[x][y] {
			visited[x][y] = true
			pq.Push([3]int{A[x][y], x, y})
		}
	}

	for i := 0; i < H; i++ {
		push(i, 0)
		push(i, W-1)
	}
	for i := 0; i < W; i++ {
		push(0, i)
		push(H-1, i)
	}

	res := H * W
	for i := 1; i <= Y; i++ {
		for pq.Len() > 0 {
			hxy := pq.Top()
			h, x, y := hxy[0], hxy[1], hxy[2]
			if h <= i {
				pq.Pop()
				res--
				for _, d := range dir4 {
					push(x+d[0], y+d[1])
				}
			} else {
				break
			}
		}
		fmt.Fprintln(out, res)
	}
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
