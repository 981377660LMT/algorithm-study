// https://sotanishy.github.io/cp-library-cpp/dp/hu_tucker.cpp

package main

import (
	"bufio"
	"fmt"
	"os"
)

// https://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=2415
func main() {
	in, out := bufio.NewReader(os.Stdin), bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	weights := make([]int, n)
	for i := range weights {
		fmt.Fscan(in, &weights[i])
	}

	fmt.Fprintln(out, HuTucker(weights))
}

// 给定每个顶点的权值，求最优二叉树的最小代价.
func HuTucker(weights []int) int {
	n := len(weights)
	heaps := make([]*LHeap, n-1)
	for i := 0; i < n-1; i++ {
		heaps[i] = NewLHeap(nil)
	}
	left := make([]int, n)
	right := make([]int, n)
	cs := make([]int, n-1)
	pq := NewHeap(func(a, b H) bool {
		return a.dist < b.dist
	}, nil)
	for i := 0; i < n-1; i++ {
		left[i] = i - 1
		right[i] = i + 1
		cs[i] = weights[i] + weights[i+1]
		pq.Push(H{cs[i], i})
	}

	res := 0
	for k := 0; k < n-1; k++ {
		item := pq.Pop()
		c, i := item.dist, item.id
		for right[i] == -1 || cs[i] != c {
			item = pq.Pop()
			c, i = item.dist, item.id
		}
		mergeL, mergeR := false, false
		if weights[i]+weights[right[i]] == c {
			mergeL, mergeR = true, true
		} else {
			_, top := heaps[i].Top()
			heaps[i].Pop()
			if weights[i]+top == c {
				mergeL = true
			} else if top+weights[right[i]] == c {
				mergeR = true
			} else {
				heaps[i].Pop()
			}
		}
		res += c
		heaps[i].Push(-1, c)
		if mergeL {
			weights[i] = INF
		}
		if mergeR {
			weights[right[i]] = INF
		}
		if mergeL && i > 0 {
			j := left[i]
			heaps[j] = MeldHeap(heaps[j], heaps[i])
			right[j] = right[i]
			right[i] = -1
			left[right[j]] = j
			i = j
		}
		if mergeR && right[i] < n-1 {
			j := right[i]
			heaps[i] = MeldHeap(heaps[i], heaps[j])
			right[i] = right[j]
			right[j] = -1
			left[right[i]] = i
		}
		cs[i] = weights[i] + weights[right[i]]
		if !heaps[i].Empty() {
			_, top := heaps[i].Top()
			heaps[i].Pop()
			cs[i] = min(cs[i], min(weights[i], weights[right[i]])+top)
			if !heaps[i].Empty() {
				_, second := heaps[i].Top()
				cs[i] = min(cs[i], top+second)
			}
			heaps[i].Push(-1, top)
		}
		pq.Push(H{cs[i], i})
	}

	return res
}

const INF int = 1e18

type Node struct {
	left, right *Node
	s           int
	id          int
	val, lazy   int
}

func NewNode(id, x int) *Node {
	return &Node{id: id, val: x}
}

type LHeap struct {
	root *Node
}

func NewLHeap(root *Node) *LHeap {
	return &LHeap{root: root}
}

func (h *LHeap) Meld(a, b *LHeap) *LHeap {
	return &LHeap{root: MeldNode(a.root, b.root)}
}

func (h *LHeap) Top() (int, int) {
	pushDown(h.root)
	return h.root.id, h.root.val
}

func (h *LHeap) Pop() {
	pushDown(h.root)
	h.root = MeldNode(h.root.left, h.root.right)
}

func (h *LHeap) Push(id, x int) {
	h.root = MeldNode(h.root, &Node{id: id, val: x})
}

func (h *LHeap) Empty() bool {
	return h.root == nil
}

func (h *LHeap) Add(x int) {
	h.root.lazy += x
}

func MeldHeap(a, b *LHeap) *LHeap {
	return &LHeap{root: MeldNode(a.root, b.root)}
}

func MeldNode(a, b *Node) *Node {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	pushDown(a)
	pushDown(b)
	if a.val > b.val {
		a, b = b, a
	}
	a.right = MeldNode(a.right, b)
	if a.left == nil || a.left.s < a.right.s {
		a.left, a.right = a.right, a.left
	}
	if a.right == nil {
		a.s = 1
	} else {
		a.s = a.right.s + 1
	}
	return a
}

func pushDown(t *Node) {
	if t.left != nil {
		t.left.lazy += t.lazy
	}
	if t.right != nil {
		t.right.lazy += t.lazy
	}
	t.val += t.lazy
	t.lazy = 0
}

type H = struct{ dist, id int }

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

func (h *Heap) Top() (value H) {
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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
