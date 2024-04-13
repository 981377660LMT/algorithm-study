// 给定一颗树，树根为1，树上有n个顶点，第i个顶点有一笔价值wi的财富。
// 我们现在可以派出k个人，从k个顶点出发，向根出发，搜集路上所有的财富（同一个顶点的财富只会被搜集一次）。
// 问我们最多可以得到多少总财富？其中−10^9≤wi≤10^9,1≤k≤n≤10^6.

// https://taodaling.github.io/blog/2019/04/30/%E5%8C%BA%E9%97%B4%E6%93%8D%E4%BD%9C%E9%97%AE%E9%A2%98/
// 每次选择出发收益最大的顶点出发，k次后得到的就是最大收益。
// 假设v是初始时出发收益最大的顶点，那么如果最后的方案中不选择v，记方案中被搜集的顶点形成的树为T，
// 我们可以找到T中v的深度最小的祖先，回退之前T下的任意一个特殊顶点带来的影响，并加入v，可以证明这时候我们的收益是不会减少的。
// 首先我们可以利用所有顶点的dfs序构建线段树，之后利用线段树实现区间修改和弹出全局收益最大顶点的能力（即大根堆）。

package main

import "fmt"

func main() {
	//    0
	//   / \
	//  1   2
	//     / \
	//    3   4
	//   / \
	//  5   6

	n := int32(7)
	T := NewTreePickSpecialPointProblem(n)
	edges := [][2]int32{{1, 0}, {2, 0}, {3, 2}, {4, 2}, {5, 3}, {6, 3}}
	tree := make([][]int32, n)
	for _, e := range edges {
		u, v := e[0], e[1]
		tree[u] = append(tree[u], v)
		tree[v] = append(tree[v], u)
	}

	weight := []int{1, 2, 3, 4, 5, 6, 7}
	canSelect := []bool{true, true, true, true, true, true, true}
	select_ := make([]bool, n)
	k := int32(1)
	res := T.Apply(tree, weight, 0, k, canSelect, select_)
	fmt.Println(res, select_)
	res = T.Apply2(tree, weight, 0, canSelect, select_)
	fmt.Println(res, select_)
	res = T.ApplyNotMoreThan(tree, weight, 0, canSelect, select_, 3)
	fmt.Println(res, select_)

	// seg := newLongPriorityQueueBasedOnSegment(0, 7)
	// nums := []int{3, 1, 2, 4, 5, 6, 7, 8}
	// seg.reset(0, 7, func(i int32) int { return nums[i] })
	// fmt.Println(seg.pop(0, 7))
	// fmt.Println(seg.pop(0, 7))
	// fmt.Println(seg.pop(0, 7))
	// fmt.Println(seg.pop(0, 7))
	// fmt.Println(seg.pop(0, 7))
	// fmt.Println(seg.pop(0, 7))
	// fmt.Println(seg.pop(0, 7))
	// fmt.Println(seg.pop(0, 7))

}

const INF int = 2e18

type TreePickSpecialPointProblem struct {
	tree        [][]int32
	weight      []int
	sumOfWeight []int
	parents     []int32
	n           int32
	segment     *longPriorityQueueBasedOnSegment
	lid         []int32
	rid         []int32 // inclusive
	idToNode    []int32
	visited     []bool
	dfn         int32
}

func NewTreePickSpecialPointProblem(n int32) *TreePickSpecialPointProblem {
	return &TreePickSpecialPointProblem{
		n:           n,
		sumOfWeight: make([]int, n),
		parents:     make([]int32, n),
		segment:     newLongPriorityQueueBasedOnSegment(0, n),
		visited:     make([]bool, n),
		lid:         make([]int32, n),
		rid:         make([]int32, n),
		idToNode:    make([]int32, n),
	}
}

// 从中恰好选择k个特殊点.
// choice[i] 表示第i个点是否是特殊点.
// select_[i] 表示第i个点是否被选择.
func (t *TreePickSpecialPointProblem) Apply(tree [][]int32, weight []int, root int32, k int32, canSelect []bool, select_ []bool) int {
	t.prepare(tree, weight, root, canSelect)
	for i := 0; i < len(tree); i++ {
		select_[i] = false
	}
	res := 0
	for i := int32(0); i < k; i++ {
		res -= t.segment.minimum
		node := t.idToNode[t.segment.pop(0, t.n)]
		t.visit(node)
		select_[node] = true
	}
	return res
}

// 选择任意个特殊点.
func (t *TreePickSpecialPointProblem) Apply2(tree [][]int32, weight []int, root int32, canSelect []bool, select_ []bool) int {
	t.prepare(tree, weight, root, canSelect)
	for i := 0; i < len(tree); i++ {
		select_[i] = false
	}
	res := 0
	for t.segment.minimum < 0 {
		res -= t.segment.minimum
		node := t.idToNode[t.segment.pop(0, t.n)]
		t.visit(node)
		select_[node] = true
	}
	return res
}

// 选择任意个特殊点.
func (t *TreePickSpecialPointProblem) ApplyNotMoreThan(tree [][]int32, weight []int, root int32, canSelect []bool, select_ []bool, limit int32) int {
	t.prepare(tree, weight, root, canSelect)
	for i := 0; i < len(tree); i++ {
		select_[i] = false
	}
	res := 0
	for t.segment.minimum < 0 && limit > 0 {
		limit--
		res -= t.segment.minimum
		node := t.idToNode[t.segment.pop(0, t.n)]
		t.visit(node)
		select_[node] = true
	}
	return res
}

func (t *TreePickSpecialPointProblem) prepare(tree [][]int32, weight []int, root int32, canSelect []bool) {
	t.tree = tree
	t.weight = weight
	t.dfn = 0
	t.dfs(root, -1, 0)
	t.segment.reset(0, t.n, func(i int32) int {
		if i < int32(len(tree)) && canSelect[t.idToNode[i]] {
			return -t.sumOfWeight[t.idToNode[i]]
		}
		return INF
	})
	for i := 0; i < len(tree); i++ {
		t.visited[i] = false
	}
}

func (t *TreePickSpecialPointProblem) visit(root int32) {
	if root == -1 || t.visited[root] {
		return
	}
	t.visited[root] = true
	t.segment.update(t.lid[root], t.rid[root], 0, t.n, t.weight[root])
	t.visit(t.parents[root])
}

func (t *TreePickSpecialPointProblem) dfs(root, p int32, sum int) {
	t.lid[root] = t.dfn
	t.dfn++
	t.idToNode[t.lid[root]] = root
	t.sumOfWeight[root] = t.weight[root] + sum
	t.parents[root] = p
	for _, e := range t.tree[root] {
		if e == p {
			continue
		}
		t.dfs(e, root, t.sumOfWeight[root])
	}
	t.rid[root] = t.dfn - 1
}

type longPriorityQueueBasedOnSegment struct {
	left, right    *longPriorityQueueBasedOnSegment
	minimum, dirty int
}

func createNIL() *longPriorityQueueBasedOnSegment {
	res := &longPriorityQueueBasedOnSegment{}
	res.left = res
	res.right = res
	return res
}

func newLongPriorityQueueBasedOnSegment(l, r int32) *longPriorityQueueBasedOnSegment {
	if l < r {
		m := (l + r) >> 1
		res := &longPriorityQueueBasedOnSegment{}
		res.left = newLongPriorityQueueBasedOnSegment(l, m)
		res.right = newLongPriorityQueueBasedOnSegment(m+1, r)
		res.pushUp()
		return res
	}
	return &longPriorityQueueBasedOnSegment{minimum: INF}
}

func (pq *longPriorityQueueBasedOnSegment) reset(l, r int32, f func(int32) int) {
	pq.dirty = 0
	if l < r {
		m := (l + r) >> 1
		pq.left.reset(l, m, f)
		pq.right.reset(m+1, r, f)
		pq.pushUp()
	} else {
		pq.minimum = f(l)
	}
}

func (pq *longPriorityQueueBasedOnSegment) covered(ll, rr, l, r int32) bool {
	return ll <= l && rr >= r
}

func (pq *longPriorityQueueBasedOnSegment) noIntersection(ll, rr, l, r int32) bool {
	return ll > r || rr < l
}

func (pq *longPriorityQueueBasedOnSegment) update(ll, rr, l, r int32, val int) {
	if pq.noIntersection(ll, rr, l, r) {
		return
	}
	if pq.covered(ll, rr, l, r) {
		pq.modify(val)
		return
	}
	m := (l + r) >> 1
	pq.pushDown()
	pq.left.update(ll, rr, l, m, val)
	pq.right.update(ll, rr, m+1, r, val)
	pq.pushUp()
}

func (pq *longPriorityQueueBasedOnSegment) pop(l, r int32) int32 {
	if l == r {
		pq.minimum = INF
		return l
	}
	var res int32
	pq.pushDown()
	m := (l + r) >> 1
	if pq.left != nil && pq.left.minimum == pq.minimum {
		res = pq.left.pop(l, m)
	} else {
		res = pq.right.pop(m+1, r)
	}
	pq.pushUp()
	return res
}

func (pq *longPriorityQueueBasedOnSegment) modify(x int) {
	pq.minimum += x
	pq.dirty += x
}

func (pq *longPriorityQueueBasedOnSegment) pushUp() {
	pq.minimum = min(pq.left.minimum, pq.right.minimum)
}

func (pq *longPriorityQueueBasedOnSegment) pushDown() {
	if pq.dirty != 0 {
		pq.left.modify(pq.dirty)
		pq.right.modify(pq.dirty)
		pq.dirty = 0
	}
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
