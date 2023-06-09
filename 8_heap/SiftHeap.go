// 原地升降的二叉堆.用于高速化dijkstra算法.

package main

import (
	"fmt"
	"time"
)

func main() {
	values := []string{"banana", "apple", "pear", "orange", "grape"}
	heap := NewSiftHeap(len(values), func(i, j int32) bool {
		return values[i] < values[j]
	})
	for i := 0; i < len(values); i++ {
		heap.Push(i)
	}

	for heap.Size() > 0 {
		fmt.Println(heap.Pop()) // 1 0 4 3 2 按照字典序输出
	}

	n := int(1e7)
	time1 := time.Now()
	pq := NewSiftHeap(n, func(i, j int32) bool { return i < j })
	for i := 0; i < n; i++ {
		pq.Push(i)
		pq.Push(i)
	}
	for pq.Size() > 0 {
		pq.Pop()
		pq.Peek()
	}
	fmt.Println(time.Since(time1)) // 1.7246415s

}

type SiftHeap struct {
	heap []int32
	pos  []int32
	less func(i, j int32) bool
	ptr  int32
}

func NewSiftHeap(n int, less func(i, j int32) bool) *SiftHeap {
	pos := make([]int32, n)
	for i := 0; i < n; i++ {
		pos[i] = -1
	}
	return &SiftHeap{
		heap: make([]int32, n),
		pos:  pos,
		less: less,
	}
}

func (h *SiftHeap) Push(i int) {
	if h.pos[i] == -1 {
		h.pos[i] = h.ptr
		h.heap[h.ptr] = int32(i)
		h.ptr++
	}
	h._siftUp(int32(i))
}

// 如果不存在,则返回-1.
func (h *SiftHeap) Pop() int {
	if h.ptr == 0 {
		return -1
	}
	res := h.heap[0]
	h.pos[res] = -1
	h.ptr--
	ptr := h.ptr
	if ptr > 0 {
		tmp := h.heap[ptr]
		h.pos[tmp] = 0
		h.heap[0] = tmp
		h._siftDown(tmp)
	}
	return int(res)
}

// 如果不存在,则返回-1.
func (h *SiftHeap) Peek() int {
	if h.ptr == 0 {
		return -1
	}
	return int(h.heap[0])
}

func (h *SiftHeap) Size() int {
	return int(h.ptr)
}

func (h *SiftHeap) _siftUp(i int32) {
	curPos := h.pos[i]
	p := int32(0)
	for curPos != 0 {
		p = h.heap[(curPos-1)>>1]
		if !h.less(i, p) {
			break
		}
		h.pos[p] = curPos
		h.heap[curPos] = p
		curPos = (curPos - 1) >> 1
	}
	h.pos[i] = curPos
	h.heap[curPos] = i
}

func (h *SiftHeap) _siftDown(i int32) {
	curPos := h.pos[i]
	c := int32(0)
	for {
		c = (curPos << 1) | 1
		if c >= h.ptr {
			break
		}
		if c+1 < h.ptr && h.less(h.heap[c+1], h.heap[c]) {
			c++
		}
		if !h.less(h.heap[c], i) {
			break
		}
		tmp := h.heap[c]
		h.heap[curPos] = tmp
		h.pos[tmp] = curPos
		curPos = c
	}
	h.pos[i] = curPos
	h.heap[curPos] = i
}
