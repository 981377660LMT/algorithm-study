// 有一个序列A。X,Y是给定的A的两个子串。
// 每次操作可以在X的开头或末尾增添或删除一个数字，且需满足任意时刻X非空且为A的子串，
// 求把X变成Y的最少操作次数。
// 题目保证解存在
// n<=2e5，1<=nums[i]<=n

// 1. 如果最长公共子串>=1，那么可以将X一直删到LCS，然后将X一直加到Y，操作len(X)+len(Y)-2*LCS次
// 2. 如果最长公共子串=0，那么就要求从X中某个顶点出发跑到Y中某个顶点的最短路

package main

import (
	"bufio"
	"fmt"
	"os"
)

func keepBeingSubstring(nums1, nums2, nums3 []int) int {
	s1, s2, _, _ := longestCommonSubstring2(nums2, nums3)
	lcs := s2 - s1
	if lcs >= 1 {
		return len(nums2) + len(nums3) - 2*lcs
	}

	// nums2和nums3不存在相同字符串,就要求从nums2中某个顶点出发跑到nums3中某个顶点的最短路
	// 多源bfs用虚拟源点汇点实现
	n := len(nums1)
	START, END := n+5, n+6
	graph := make([][]Edge, n+10)
	for _, v := range nums2 {
		graph[START] = append(graph[START], Edge{v, len(nums2) - 1}) // 删到只剩下一个字符
	}
	for _, v := range nums3 {
		graph[v] = append(graph[v], Edge{END, len(nums3) - 1}) // 删到只剩下一个字符
	}
	for i := 0; i < n-1; i++ {
		graph[nums1[i]] = append(graph[nums1[i]], Edge{nums1[i+1], 2}) // 权重为2(删1+加1)
		graph[nums1[i+1]] = append(graph[nums1[i+1]], Edge{nums1[i], 2})
	}

	dist := Dijkstra(n+10, graph, START)
	return dist[END]
}

const INF int = 1e18

type Edge struct{ to, weight int }

func Dijkstra(n int, adjList [][]Edge, start int) (dist []int) {
	dist = make([]int, n)
	for i := range dist {
		dist[i] = INF
	}
	dist[start] = 0

	pq := NewHeap(func(a, b H) int {
		return a.dist - b.dist
	}, []H{{start, 0}})

	for pq.Len() > 0 {
		curNode := pq.Pop()
		cur, curDist := curNode.node, curNode.dist
		if curDist > dist[cur] {
			continue
		}

		for _, edge := range adjList[cur] {
			next, weight := edge.to, edge.weight
			if cand := curDist + weight; cand < dist[next] {
				dist[next] = cand
				pq.Push(H{next, cand})
			}
		}
	}

	return
}

type H = struct{ node, dist int }

// Should return a number:
//    negative , if a < b
//    zero     , if a == b
//    positive , if a > b
type Comparator func(a, b H) int

func NewHeap(comparator Comparator, nums []H) *Heap {
	nums = append(nums[:0:0], nums...)
	heap := &Heap{comparator: comparator, data: nums}
	heap.heapify()
	return heap
}

type Heap struct {
	data       []H
	comparator Comparator
}

func (h *Heap) Push(value H) {
	h.data = append(h.data, value)
	h.pushUp(h.Len() - 1)
}

func (h *Heap) Pop() (value H) {
	if h.Len() == 0 {
		return
	}

	value = h.data[0]
	h.data[0] = h.data[h.Len()-1]
	h.data = h.data[:h.Len()-1]
	h.pushDown(0)
	return
}

func (h *Heap) Peek() (value H) {
	if h.Len() == 0 {
		return
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
	for parent := (root - 1) >> 1; parent >= 0 && h.comparator(h.data[root], h.data[parent]) < 0; parent = (root - 1) >> 1 {
		h.data[root], h.data[parent] = h.data[parent], h.data[root]
		root = parent
	}
}

func (h *Heap) pushDown(root int) {
	n := h.Len()
	for left := (root<<1 + 1); left < n; left = (root<<1 + 1) {
		right := left + 1
		minIndex := root

		if h.comparator(h.data[left], h.data[minIndex]) < 0 {
			minIndex = left
		}

		if right < n && h.comparator(h.data[right], h.data[minIndex]) < 0 {
			minIndex = right
		}

		if minIndex == root {
			return
		}

		h.data[root], h.data[minIndex] = h.data[minIndex], h.data[root]
		root = minIndex
	}
}

// 最长公共子串
func longestCommonSubstring1(s1, s2 string) (start1, end1, start2, end2 int) {
	ords1, ords2 := make([]int, len(s1)), make([]int, len(s2))
	for i, x := range s1 {
		ords1[i] = int(x)
	}
	for i, x := range s2 {
		ords2[i] = int(x)
	}
	return longestCommonSubstring2(ords1, ords2)
}

func longestCommonSubstring2(ords1, ords2 []int) (start1, end1, start2, end2 int) {
	if len(ords1) == 0 || len(ords2) == 0 {
		return
	}

	dummy := max(maxs(ords1...), maxs(ords2...)) + 1
	sb := make([]int, 0, len(ords1)+len(ords2)+1)
	sb = append(sb, ords1...)
	sb = append(sb, dummy)
	sb = append(sb, ords2...)
	sa, _, lcp := UseSA(sb)

	len_ := 0
	for i := 1; i < len(sb); i++ {
		if (sa[i-1] < len(ords1)) == (sa[i] < len(ords1)) {
			continue
		}
		if lcp[i] <= len_ {
			continue
		}
		len_ = lcp[i]

		// 来自s和t的不同子串
		// 找到了(严格)更长的公共子串,更新答案
		i1, i2 := sa[i-1], sa[i]
		if i1 > i2 {
			i1, i2 = i2, i1
		}

		start1 = i1
		end1 = start1 + len_
		start2 = i2 - len(ords1) - 1
		end2 = start2 + len_
	}

	return
}

func GetSA(ords []int) (sa []int) {
	if len(ords) == 0 {
		return []int{}
	}

	mn := mins(ords...)
	for i, x := range ords {
		ords[i] = x - mn + 1
	}
	ords = append(ords, 0)
	n := len(ords)
	m := maxs(ords...) + 1
	isS := make([]bool, n)
	isLms := make([]bool, n)
	lms := make([]int, 0, n)
	for i := 0; i < n; i++ {
		isS[i] = true
	}
	for i := n - 2; i > -1; i-- {
		if ords[i] == ords[i+1] {
			isS[i] = isS[i+1]
		} else {
			isS[i] = ords[i] < ords[i+1]
		}
	}
	for i := 1; i < n; i++ {
		isLms[i] = !isS[i-1] && isS[i]
	}
	for i := 0; i < n; i++ {
		if isLms[i] {
			lms = append(lms, i)
		}
	}
	bin := make([]int, m)
	for _, x := range ords {
		bin[x]++
	}

	induce := func() []int {
		sa := make([]int, n)
		for i := 0; i < n; i++ {
			sa[i] = -1
		}

		saIdx := make([]int, m)
		copy(saIdx, bin)
		for i := 0; i < m-1; i++ {
			saIdx[i+1] += saIdx[i]
		}
		for j := len(lms) - 1; j > -1; j-- {
			i := lms[j]
			x := ords[i]
			saIdx[x]--
			sa[saIdx[x]] = i
		}

		copy(saIdx, bin)
		s := 0
		for i := 0; i < m; i++ {
			s, saIdx[i] = s+saIdx[i], s
		}
		for j := 0; j < n; j++ {
			i := sa[j] - 1
			if i < 0 || isS[i] {
				continue
			}
			x := ords[i]
			sa[saIdx[x]] = i
			saIdx[x]++
		}

		copy(saIdx, bin)
		for i := 0; i < m-1; i++ {
			saIdx[i+1] += saIdx[i]
		}
		for j := n - 1; j > -1; j-- {
			i := sa[j] - 1
			if i < 0 || !isS[i] {
				continue
			}
			x := ords[i]
			saIdx[x]--
			sa[saIdx[x]] = i
		}

		return sa
	}

	sa = induce()

	lmsIdx := make([]int, 0, len(sa))
	for _, i := range sa {
		if isLms[i] {
			lmsIdx = append(lmsIdx, i)
		}
	}
	l := len(lmsIdx)
	order := make([]int, n)
	for i := 0; i < n; i++ {
		order[i] = -1
	}
	ord := 0
	order[n-1] = ord
	for i := 0; i < l-1; i++ {
		j, k := lmsIdx[i], lmsIdx[i+1]
		for d := 0; d < n; d++ {
			jIsLms, kIsLms := isLms[j+d], isLms[k+d]
			if ords[j+d] != ords[k+d] || jIsLms != kIsLms {
				ord++
				break
			}
			if d > 0 && (jIsLms || kIsLms) {
				break
			}
		}
		order[k] = ord
	}
	b := make([]int, 0, l)
	for _, i := range order {
		if i >= 0 {
			b = append(b, i)
		}
	}
	var lmsOrder []int
	if ord == l-1 {
		lmsOrder = make([]int, l)
		for i, ord := range b {
			lmsOrder[ord] = i
		}
	} else {
		lmsOrder = GetSA(b)
	}
	buf := make([]int, len(lms))
	for i, j := range lmsOrder {
		buf[i] = lms[j]
	}
	lms = buf
	return induce()[1:]
}

//  sa : 排第几的后缀是谁.
//  rank : 每个后缀排第几.
//  lcp : 排名相邻的两个后缀的最长公共前缀.
// 	lcp[0] = 0
// 	lcp[i] = LCP(s[sa[i]:], s[sa[i-1]:])
func UseSA(ords []int) (sa, rank, lcp []int) {
	n := len(ords)
	sa = GetSA(ords)

	rank = make([]int, n)
	for i := range rank {
		rank[sa[i]] = i
	}

	// !高度数组 lcp 也就是排名相邻的两个后缀的最长公共前缀。
	// lcp[0] = 0
	// lcp[i] = LCP(s[sa[i]:], s[sa[i-1]:])
	lcp = make([]int, n)
	h := 0
	for i, rk := range rank {
		if h > 0 {
			h--
		}
		if rk > 0 {
			for j := int(sa[rk-1]); i+h < n && j+h < n && ords[i+h] == ords[j+h]; h++ {
			}
		}
		lcp[rk] = h
	}

	return
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func mins(a ...int) int {
	mn := a[0]
	for _, x := range a {
		if x < mn {
			mn = x
		}
	}
	return mn
}

func maxs(a ...int) int {
	mx := a[0]
	for _, x := range a {
		if x > mx {
			mx = x
		}
	}
	return mx
}

func sum(a ...int) int {
	s := 0
	for _, x := range a {
		s += x
	}
	return s
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	nums1 := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums1[i])
	}

	var p int
	fmt.Fscan(in, &p)
	nums2 := make([]int, p)
	for i := 0; i < p; i++ {
		fmt.Fscan(in, &nums2[i])
	}

	var q int
	fmt.Fscan(in, &q)
	nums3 := make([]int, q)
	for i := 0; i < q; i++ {
		fmt.Fscan(in, &nums3[i])
	}

	fmt.Fprintln(out, keepBeingSubstring(nums1, nums2, nums3))
}
