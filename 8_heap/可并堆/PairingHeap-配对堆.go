// https://noshi91.github.io/Library/data_structure/pairing_heap.cpp
// https://scrapbox.io/data-structures/Pairing_Heap
// https://oi-wiki.org/ds/pairing-heap/

// 配对堆是一棵满足堆性质的带权多叉树（如下图），即每个节点的权值都小于或等于他的所有儿子
// 一个节点的所有儿子节点形成一个单向链表。每个节点储存第一个儿子的指针，即链表的头节点；和他的右兄弟的指针。

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	P3066()
}

// https://www.luogu.com.cn/problem/P3066
// 给出以0号点为根的一棵有根树,问每个点的子树中与它距离小于等于k的点有多少个
func P3066() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k int
	fmt.Fscan(in, &n, &k)

	tree := make([][][2]int, n)
	for i := 1; i < n; i++ {
		var parent, w int
		fmt.Fscan(in, &parent, &w)
		parent--
		tree[parent] = append(tree[parent], [2]int{i, w})
	}

	subHeaps := make([]*PairingHeap, n)
	for i := 0; i < n; i++ {
		subHeaps[i] = NewPairingHeap(func(a, b int) bool { return a > b })
	}
	subHeapSize := make([]int, n)
	res := make([]int, n)
	var dfs func(cur, pre int, dist int)
	dfs = func(cur, pre int, dist int) {
		subHeaps[cur].Push(dist)
		subHeapSize[cur]++
		for _, e := range tree[cur] {
			next, weight := e[0], e[1]
			if next == pre {
				continue
			}
			dfs(next, cur, dist+weight)
			subHeapSize[cur] += subHeapSize[next]
			subHeaps[cur] = Meld(subHeaps[cur], subHeaps[next])
		}
		for subHeapSize[cur] > 0 && subHeaps[cur].Top()-dist > k {
			subHeaps[cur].Pop()
			subHeapSize[cur]--
		}
		res[cur] = subHeapSize[cur]
	}
	dfs(0, -1, 0)

	for i := 0; i < n; i++ {
		fmt.Fprintln(out, res[i])
	}
}

func demo() {
	pq1 := NewPairingHeap(func(a, b int) bool { return a < b })
	pq1.Push(1)
	pq1.Push(2)
	pq2 := NewPairingHeap(func(a, b int) bool { return a < b })
	pq2.Push(3)
	pq2.Push(4)
	pq3 := Meld(pq1, pq2)
	fmt.Println(pq3.Pop())
	fmt.Println(pq3.Pop())
	fmt.Println(pq3.Pop())
	fmt.Println(pq3.Pop())
	fmt.Println(pq3.Empty())
}

type P = int
type PairingHeap struct {
	root *pNode
	less func(P, P) bool
}

type pNode struct {
	value P
	head  *pNode
	next  *pNode // sibling
}

func NewPairingHeap(less func(P, P) bool) *PairingHeap {
	return &PairingHeap{less: less}
}

// 融合两个堆, 返回新的堆, 原来的堆被破坏.
//
//	!注意两个堆的比较函数必须相同.
func Meld(h1, h2 *PairingHeap) *PairingHeap {
	return &PairingHeap{root: _merge(h1.root, h2.root, h1.less), less: h1.less}
}

func (h *PairingHeap) Empty() bool {
	return h.root == nil
}

func (h *PairingHeap) Top() P {
	if h.Empty() {
		panic("empty heap")
	}
	return h.root.value
}

func (h *PairingHeap) Push(x P) {
	h.root = _merge(h.root, &pNode{value: x}, h.less)
}

func (h *PairingHeap) Pop() P {
	if h.Empty() {
		panic("empty heap")
	}
	res := h.root.value
	h.root = _mergeList(h.root.head, h.less)
	return res
}

func _merge(x, y *pNode, less func(P, P) bool) *pNode {
	if x == nil {
		return y
	}
	if y == nil {
		return x
	}

	if !less(x.value, y.value) {
		x, y = y, x
	}
	y.next = x.head
	x.head = y
	return x
}

func _mergeList(list *pNode, less func(P, P) bool) *pNode {
	if list == nil || list.next == nil {
		return list
	}
	next := list.next
	rem := next.next
	return _merge(_merge(list, next, less), _mergeList(rem, less), less)
}
