// !Deprecated 使用 "差分约束"代替.
// https://taodaling.github.io/blog/2019/06/14/%E5%B7%AE%E5%88%86%E7%BA%A6%E6%9D%9F%E7%B3%BB%E7%BB%9F/
// https://www.luogu.com.cn/training/22983
// api:
// 1. NewDifferSystem(n int32) *DifferSystem
// 2. (ds *DifferSystem) LessThanOrEqualTo(i, j int32, d int)
// 3. (ds *DifferSystem) GreaterThanOrEqualTo(i, j int32, d int)
// 4. (ds *DifferSystem) EqualTo(i, j int32, d int)
// 5. (ds *DifferSystem) LessThan(i, j int32, d int)
// 6. (ds *DifferSystem) GreaterThan(i, j int32, d int)
// 7. (ds *DifferSystem) HasSolution() bool
// 8. (ds *DifferSystem) FindMaxDifferenceBetween(i, j int32) int
// 9. (ds *DifferSystem) FindMinDifferenceBetween(i, j int32) int
// 10. (ds *DifferSystem) RunSince(j int32) bool
// 11. (ds *DifferSystem) PossibleSolutionOf(i int32) int
// 12. (ds *DifferSystem) Clear(n int32)
//
// -note:
// 1. 找到一组可行解：虚拟源点n向其它所有顶点连长度为0的有向边，之后跑单源最短路径即可。
// 2. 找到更多的可行解：假设我们找到了一组可行解，那么我们将所有未知变量加上某个相同常量t，就可以得到无数的可行解了
// 3. 判断差分约束系统有解：无负环。
// 4. 计算 xa−xb 的最小（大）值: 直接以xb为起点跑最短路算法，这时候顶点b到顶点a的最短距离就是max(xa-xb)。
//
// !默认最大化解.如果要最小化解, 需要令 y = -x, 然后求最大化解(即将原来的边反向).

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func demo() {
	ds := NewDifferSystem(3)
	ds.GreaterThanOrEqualTo(0, 1, 10)
	ds.LessThanOrEqualTo(1, 2, 2)
	fmt.Println(ds.HasSolution())
	fmt.Println(ds)
	fmt.Println(ds.PossibleSolutionOf(2))
	fmt.Println(ds.FindMinDifferenceBetween(0, 1))
	// fmt.Println(ds.FindMaxDifferenceBetween(0, 1))
}

const INF int = 2e18

// !默认最大化解.如果要最小化解, 需要令 y = -x, 然后求最大化解(即将原来的边反向).
type DifferSystem struct {
	nodes       []*DifferNode
	queue       *Deque
	n           int32
	allPos      bool
	hasSolution bool
}

func NewDifferSystem(n int32) *DifferSystem {
	res := &DifferSystem{
		nodes:  make([]*DifferNode, n),
		queue:  NewDeque(n),
		n:      n,
		allPos: true,
	}
	for i := int32(0); i < n; i++ {
		node := NewNode()
		node.id = i
		res.nodes[i] = node
	}
	return res
}

// d[i] - d[j] <= d
func (ds *DifferSystem) LessThanOrEqualTo(i, j int32, d int) {
	ds.nodes[j].adj = append(ds.nodes[j].adj, NewNeighbor(ds.nodes[i], d))
	ds.allPos = ds.allPos && d >= 0
}

// d[i] - d[j] >= d
func (ds *DifferSystem) GreaterThanOrEqualTo(i, j int32, d int) {
	ds.LessThanOrEqualTo(j, i, -d)
}

// d[i] - d[j] = d
func (ds *DifferSystem) EqualTo(i, j int32, d int) {
	ds.GreaterThanOrEqualTo(i, j, d)
	ds.LessThanOrEqualTo(i, j, d)
}

// d[i] - d[j] < d
func (ds *DifferSystem) LessThan(i, j int32, d int) {
	ds.LessThanOrEqualTo(i, j, d-1)
}

// d[i] - d[j] > d
func (ds *DifferSystem) GreaterThan(i, j int32, d int) {
	ds.GreaterThanOrEqualTo(i, j, d+1)
}

func (ds *DifferSystem) HasSolution() bool {
	ds._prepare(0)
	for i := int32(0); i < ds.n; i++ {
		ds.nodes[i].inQueue = true
		ds.queue.Append(ds.nodes[i])
	}
	ds.hasSolution = ds._spfa()
	return ds.hasSolution
}

// Find max(ai - aj), if INF is returned, it means no constraint between ai and aj.
func (ds *DifferSystem) FindMaxDifferenceBetween(i, j int32) int {
	ds.RunSince(j)
	return ds.nodes[i].dist
}

// Find min(ai - aj), if INF is returned, it means no constraint between ai and aj.
func (ds *DifferSystem) FindMinDifferenceBetween(i, j int32) int {
	r := ds.FindMaxDifferenceBetween(j, i) // !最小化解
	if r == INF {
		return INF
	}
	return -r
}

// After invoking this method, the value of i is max(ai - aj).
func (ds *DifferSystem) RunSince(j int32) bool {
	ds._prepare(INF)
	ds.queue.Clear()
	ds.queue.Append(ds.nodes[j])
	ds.nodes[j].dist = 0
	ds.nodes[j].inQueue = true
	ds.hasSolution = ds._spfa()
	return ds.hasSolution
}

func (ds *DifferSystem) PossibleSolutionOf(i int32) int {
	return ds.nodes[i].dist
}

func (ds *DifferSystem) Clear(n int32) {
	ds.n = n
	ds.allPos = true
	for i := int32(0); i < n; i++ {
		ds.nodes[i].adj = ds.nodes[i].adj[:0]
	}
}

func (ds *DifferSystem) _spfa() bool {
	if ds.allPos {
		ds._dijkstra()
		return true
	}
	queue := ds.queue
	for queue.Size() > 0 {
		head := queue.PopLeft()
		head.inQueue = false
		if head.times >= ds.n {
			return false
		}
		for _, edge := range head.adj {
			node := edge.next
			if node.dist <= edge.weight+head.dist {
				continue
			}
			node.dist = edge.weight + head.dist
			if node.inQueue {
				continue
			}
			node.times++
			node.inQueue = true
			queue.AppendLeft(node)
		}
	}
	return true
}

func (ds *DifferSystem) _prepare(initDist int) {
	ds.queue.Clear()
	for i := int32(0); i < ds.n; i++ {
		node := ds.nodes[i]
		node.dist = initDist
		node.times = 0
		node.inQueue = false
	}
}

func (ds *DifferSystem) _dijkstra() {
	pq := NewErasableHeapGeneric(
		func(a, b *DifferNode) bool {
			if a.dist == b.dist {
				return a.id < b.id
			}
			return a.dist < b.dist
		},
		ds.queue.GetAll()...,
	)
	for pq.Len() > 0 {
		head := pq.Pop()
		for _, e := range head.adj {
			if e.next.dist <= head.dist+e.weight {
				continue
			}
			pq.Erase(e.next)
			e.next.dist = head.dist + e.weight
			pq.Push(e.next)
		}
	}
}

func (ds *DifferSystem) String() string {
	sb := strings.Builder{}
	for i := int32(0); i < ds.n; i++ {
		for _, edge := range ds.nodes[i].adj {
			sb.WriteString(edge.String())
			sb.WriteByte('\n')
		}
	}
	sb.WriteString("-------------\n")
	if !ds.hasSolution {
		sb.WriteString("impossible")
	} else {
		for i := int32(0); i < ds.n; i++ {
			sb.WriteString(fmt.Sprintf("a%d=%d\n", i, ds.nodes[i].dist))
		}
	}
	return sb.String()
}

type Neighbor struct {
	next   *DifferNode
	weight int
}

func NewNeighbor(next *DifferNode, weight int) *Neighbor {
	return &Neighbor{next: next, weight: weight}
}

func (n *Neighbor) String() string {
	return fmt.Sprintf("next=%d, weight=%d", n.next.id, n.weight)
}

type DifferNode struct {
	adj     []*Neighbor
	dist    int
	inQueue bool
	times   int32
	id      int32
}

func NewNode() *DifferNode {
	return &DifferNode{}
}

type ErasableHeapGeneric[H comparable] struct {
	data   *HeapGeneric[H]
	erased *HeapGeneric[H]
	size   int
}

func NewErasableHeapGeneric[H comparable](less func(a, b H) bool, nums ...H) *ErasableHeapGeneric[H] {
	return &ErasableHeapGeneric[H]{NewHeapGeneric(less, nums...), NewHeapGeneric(less), len(nums)}
}

// 从堆中删除一个元素,要保证堆中存在该元素.
func (h *ErasableHeapGeneric[H]) Erase(value H) {
	h.erased.Push(value)
	h.normalize()
	h.size--
}

func (h *ErasableHeapGeneric[H]) Push(value H) {
	h.data.Push(value)
	h.normalize()
	h.size++
}

func (h *ErasableHeapGeneric[H]) Pop() (value H) {
	value = h.data.Pop()
	h.normalize()
	h.size--
	return
}

func (h *ErasableHeapGeneric[H]) Peek() (value H) {
	value = h.data.Top()
	return
}

func (h *ErasableHeapGeneric[H]) Len() int {
	return h.size
}

func (h *ErasableHeapGeneric[H]) Clear() {
	h.data.Clear()
	h.erased.Clear()
	h.size = 0
}

func (h *ErasableHeapGeneric[H]) normalize() {
	for h.data.Len() > 0 && h.erased.Len() > 0 && h.data.Top() == h.erased.Top() {
		h.data.Pop()
		h.erased.Pop()
	}
}

type HeapGeneric[H comparable] struct {
	data []H
	less func(a, b H) bool
}

func NewHeapGeneric[H comparable](less func(a, b H) bool, nums ...H) *HeapGeneric[H] {
	nums = append(nums[:0:0], nums...)
	heap := &HeapGeneric[H]{less: less, data: nums}
	if len(nums) > 1 {
		heap.heapify()
	}
	return heap
}

func (h *HeapGeneric[H]) Push(value H) {
	h.data = append(h.data, value)
	h.pushUp(h.Len() - 1)
}

func (h *HeapGeneric[H]) Pop() (value H) {
	if h.Len() == 0 {
		panic("heap is empty")
	}

	value = h.data[0]
	h.data[0] = h.data[h.Len()-1]
	h.data = h.data[:h.Len()-1]
	h.pushDown(0)
	return
}

func (h *HeapGeneric[H]) Top() (value H) {
	if h.Len() == 0 {
		panic("heap is empty")
	}
	value = h.data[0]
	return
}

func (h *HeapGeneric[H]) Len() int { return len(h.data) }

func (h *HeapGeneric[H]) Clear() {
	h.data = h.data[:0]
}

func (h *HeapGeneric[H]) heapify() {
	n := h.Len()
	for i := (n >> 1) - 1; i > -1; i-- {
		h.pushDown(i)
	}
}

func (h *HeapGeneric[H]) pushUp(root int) {
	for parent := (root - 1) >> 1; parent >= 0 && h.less(h.data[root], h.data[parent]); parent = (root - 1) >> 1 {
		h.data[root], h.data[parent] = h.data[parent], h.data[root]
		root = parent
	}
}

func (h *HeapGeneric[H]) pushDown(root int) {
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

type D = *DifferNode
type Deque struct{ l, r []D }

func NewDeque(cap int32) *Deque { return &Deque{make([]D, 0, 1+cap/2), make([]D, 0, 1+cap/2)} }

func (q *Deque) Empty() bool {
	return len(q.l) == 0 && len(q.r) == 0
}

func (q *Deque) Size() int {
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

func (q *Deque) Front() D {
	if len(q.l) > 0 {
		return q.l[len(q.l)-1]
	}
	return q.r[0]
}

func (q *Deque) Back() D {
	if len(q.r) > 0 {
		return q.r[len(q.r)-1]
	}
	return q.l[0]
}

// 0 <= i < q.Size()
func (q *Deque) At(i int) D {
	if i < len(q.l) {
		return q.l[len(q.l)-1-i]
	}
	return q.r[i-len(q.l)]
}

func (q *Deque) Clear() {
	q.l = q.l[:0]
	q.r = q.r[:0]
}

func (q *Deque) GetAll() []D {
	return append(q.l, q.r...)
}

func main() {
	P3275()
	// P1993()
	// P5960()
	// abc216g()
}

// P3275 [SCOI2011] 糖果 (最小化)
// https://www.luogu.com.cn/problem/P3275
// 如果X=1，表示第 A 个小朋友分到的糖果必须和第 B 个小朋友分到的糖果一样多；
// 如果X=2，表示第 A 个小朋友分到的糖果必须少于第 B 个小朋友分到的糖果；
// 如果X=3，表示第 A 个小朋友分到的糖果必须不少于第 B 个小朋友分到的糖果；
// 如果X=4，表示第 A 个小朋友分到的糖果必须多于第 B 个小朋友分到的糖果；
// 如果X=5，表示第 A 个小朋友分到的糖果必须不多于第 B 个小朋友分到的糖果；
// 输出一行，表示老师至少需要准备的糖果数，如果不能满足小朋友们的所有要求，就输出−1。
func P3275() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int32
	fmt.Fscan(in, &n, &m)
	D := NewDifferSystem(n + 1)
	for i := int32(0); i < m; i++ {
		var op, a, b int32
		fmt.Fscan(in, &op, &a, &b)
		a, b = a-1, b-1
		a, b = b, a // !反向边
		if op == 1 {
			D.EqualTo(a, b, 0)
		} else if op == 2 {
			D.LessThan(a, b, 0)
		} else if op == 3 {
			D.GreaterThanOrEqualTo(a, b, 0)
		} else if op == 4 {
			D.GreaterThan(a, b, 0)
		} else {
			D.LessThanOrEqualTo(a, b, 0)
		}
	}

	DUMMY := n // 虚拟源点
	// xi>=1 => xi- x0 >=1
	for i := int32(0); i < n; i++ {
		D.GreaterThanOrEqualTo(DUMMY, i, 1) // !反向边
	}

	ok := D.RunSince(DUMMY)
	if !ok {
		fmt.Fprintln(out, -1)
		return
	}

	v := int32(0)
	for i := int32(0); i < n; i++ {
		v += int32(D.PossibleSolutionOf(i))
	}
	fmt.Fprintln(out, -v)
}

// P1993 小 K 的农场
// https://luogu.com.cn/problem/P1993
// 小 K 在 MC 里面建立很多很多的农场，总共 n 个，
// 以至于他自己都忘记了每个农场中种植作物的具体数量了，
// 他只记得一些含糊的信息（共 m 个），
// 以下列三种形式描述：
// 农场 a 比农场 b 至少多种植了 c 个单位的作物；
// 农场 a 比农场 b 至多多种植了 c 个单位的作物；
// 农场 a 与农场 b 种植的作物数一样多。
func P1993() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int32
	fmt.Fscan(in, &n, &m)
	D := NewDifferSystem(n)
	for i := int32(0); i < m; i++ {
		var op int32
		fmt.Fscan(in, &op)
		if op == 1 {
			var a, b, c int32
			fmt.Fscan(in, &a, &b, &c)
			D.GreaterThanOrEqualTo(a-1, b-1, int(c))
		} else if op == 2 {
			var a, b, c int32
			fmt.Fscan(in, &a, &b, &c)
			D.LessThanOrEqualTo(a-1, b-1, int(c))
		} else {
			var a, b int32
			fmt.Fscan(in, &a, &b)
			D.EqualTo(a-1, b-1, 0)
		}
	}

	ok := D.HasSolution()
	if ok {
		fmt.Fprintln(out, "Yes")
	} else {
		fmt.Fprintln(out, "No")
	}
}

// P5960 【模板】差分约束 (差分约束系统)
// https://www.luogu.com.cn/problem/P5960
// 给定一些关系xi-xj<=d.
// 求任意一组满足这个不等式组的解
// 如果有多组解，请输出任意一组，无解请输出 NO。
func P5960() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int32
	fmt.Fscan(in, &n, &m)

	S := NewDifferSystem(n + 1)
	for i := int32(0); i < m; i++ {
		var u, v, d int32
		fmt.Fscan(in, &u, &v, &d)
		u, v = u-1, v-1
		S.LessThanOrEqualTo(u, v, int(d))
	}

	DUMMY := n // 赋值为0的虚拟源点
	for i := int32(0); i < n; i++ {
		S.GreaterThanOrEqualTo(DUMMY, i, 0)
	}

	ok := S.RunSince(DUMMY)
	if !ok {
		fmt.Fprintln(out, "NO")
		return
	}
	for i := int32(0); i < n; i++ {
		fmt.Fprint(out, S.PossibleSolutionOf(i), " ")
	}
}

// https://atcoder.jp/contests/abc216/tasks/abc216_g
// 一个长度为n的序列，只由0和1组成，给出m个约束条件l, r, c，表示l 到r中至少有c个1，
// 问满足条件的序列是什么，如果有多种，则输出1的数量最小的那种
// n<=2e5, m<=2e5
func abc216g() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int32
	fmt.Fscan(in, &n, &m)
	D := NewDifferSystem(n + 2)
	for i := int32(0); i < m; i++ {
		var l, r, c int32
		fmt.Fscan(in, &l, &r, &c)
		D.GreaterThanOrEqualTo(r, l-1, int(c)) // Sr-Sl-1>=c
	}
	for i := int32(1); i <= n; i++ {
		D.GreaterThanOrEqualTo(i, i-1, 0) // preSum[i] - preSum[i-1] >= 0
		D.LessThanOrEqualTo(i, i-1, 1)    // preSum[i] - preSum[i-1] <= 1
	}

	DUMMY := n + 1
	for i := int32(0); i <= n; i++ {
		D.GreaterThanOrEqualTo(DUMMY, i, 0) // Sd-Si>=0
	}
	D.RunSince(DUMMY)

	res := make([]int, 0, n)
	for i := int32(1); i <= n; i++ {
		a := D.PossibleSolutionOf(i)
		b := D.PossibleSolutionOf(i - 1)
		res = append(res, a-b) // !每个位置取0还是1
	}

	for _, v := range res {
		fmt.Fprint(out, v, " ")
	}
}
