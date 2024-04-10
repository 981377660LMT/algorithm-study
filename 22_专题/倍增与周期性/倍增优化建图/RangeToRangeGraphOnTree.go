package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {
	P5344()
	// P9520()
	// CF1904F()

	// jump()
	// test()
}

// P5344 【XR-1】逛森林 (倍增优化建图)
// https://www.luogu.com.cn/problem/P5344
// 1 u1 v1 u2 v2 w : 路径u1v1上所有结点可以花费w的代价到达路径u2v2上的所有结点，如果路径不连通则无效。
// 2 u v w：结点u和v之间连接一条费用为w的无向边.如果u和v之间已经有边，则无效.
// 最后求从结点s出发，到每个节点的最小花费.
func P5344() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q, start int32
	fmt.Fscan(in, &n, &q, &start)
	start--
	operations := make([][6]int32, q)
	for i := int32(0); i < q; i++ {
		var op int32
		fmt.Fscan(in, &op)
		if op == 1 {
			var u1, v1, u2, v2, w int32
			fmt.Fscan(in, &u1, &v1, &u2, &v2, &w)
			u1, v1, u2, v2 = u1-1, v1-1, u2-1, v2-1
			operations[i] = [6]int32{op, u1, v1, u2, v2, w}
		} else {
			var u, v, w int32
			fmt.Fscan(in, &u, &v, &w)
			u, v = u-1, v-1
			operations[i] = [6]int32{op, u, v, w}
		}
	}

	uf := NewUnionFindArraySimple32(n)
	valid := make([]bool, q) // 每个1操作是否有效
	tree := make([][]int32, n)
	for i := int32(0); i < q; i++ {
		op := &operations[i]
		if op[0] == 1 {
			u1, v1, u2, v2 := op[1], op[2], op[3], op[4]
			valid[i] = uf.Find(u1) == uf.Find(v1) && uf.Find(u2) == uf.Find(v2)
		} else {
			u, v := op[1], op[2]
			if uf.Union(u, v) {
				tree[u] = append(tree[u], v)
				tree[v] = append(tree[v], u)
				valid[i] = true
			}
		}
	}

	R := NewRangeToRangeGraphOnTree(tree, -1)
	size := R.Size()
	newGraph := make([][]Neighbor, size)

	R.Init(func(from, to int32) { newGraph[from] = append(newGraph[from], Neighbor{to, 0}) })

	for i := int32(0); i < q; i++ {
		op := &operations[i]
		if !valid[i] {
			continue
		}
		if op[0] == 1 {
			u1, v1, u2, v2, w := op[1], op[2], op[3], op[4], op[5]
			R.AddRangeToRange(u1, v1, u2, v2, func(from, to int32) {
				newGraph[from] = append(newGraph[from], Neighbor{to, w})
			})
		} else {
			u, v, w := op[1], op[2], op[3]
			R.Add(u, v, func(from, to int32) { newGraph[from] = append(newGraph[from], Neighbor{to, w}) })
			R.Add(v, u, func(from, to int32) { newGraph[from] = append(newGraph[from], Neighbor{to, w}) })
		}
	}

	dist := Dijkstra(int32(len(newGraph)), newGraph, start)
	for i := int32(0); i < n; i++ {
		d := dist[i] // !出点
		if d == INF {
			d = -1
		}
		fmt.Fprint(out, d, " ")
	}
}

type RangeToRangeGraphOnTree struct {
	tree           [][]int32
	depth          []int32
	n, log, offset int32 // 底层真实点：[0,n)，倍增入点：[n,n+offset)，倍增出点：[n+offset,n+2*offset).
	root           int32
	jump           [][]int32 // 节点j向上跳2^i步的父节点
}

// root为-1.
func NewRangeToRangeGraphOnTree(tree [][]int32, root int32) *RangeToRangeGraphOnTree {
	n := int32(len(tree))
	depth := make([]int32, n)
	g := &RangeToRangeGraphOnTree{
		tree:  tree,
		depth: depth,
		n:     n,
		log:   int32(bits.Len32(uint32(n))) - 1,
		root:  root,
	}
	g.offset = n * (g.log + 1)
	return g
}

// 建立内部连接.
func (g *RangeToRangeGraphOnTree) Init(f func(from, to int32)) {
	g.makeDp()
	if g.root == -1 {
		for i := range g.depth {
			g.depth[i] = -1
		}
		for i := int32(0); i < g.n; i++ {
			if g.depth[i] == -1 {
				g.depth[i] = 0
				g.dfsAndInitDp(i, -1, f)
			}
		}
	} else {
		g.dfsAndInitDp(g.root, -1, f)
	}
	g.updateDp()

	// pushDown jump
	n, log, offset := g.n, g.log, g.offset
	for k := log - 1; k >= 0; k-- {
		for i := int32(0); i < n; i++ {
			if to := g.jump[k][i]; to != -1 {
				c1 := k*n + i + n
				c2 := k*n + to + n
				p := c1 + n
				f(c1, p)
				f(c2, p)
				f(p+offset, c1+offset)
				f(p+offset, c2+offset)
			}
		}
	}
}

// 添加一条从from到to的边.
func (g *RangeToRangeGraphOnTree) Add(from, to int32, f func(from, to int32)) {
	f(from, to)
}

// 从区间 [fromStart, fromEnd) 中的每个点到 to 都添加一条边.
func (g *RangeToRangeGraphOnTree) AddFromRange(fromStart, fromEnd, to int32, f func(from, to int32)) {
	to += g.offset
	g.enumerateJumpDangerously(fromStart, to, func(id int32) { f(id, to) })
}

// 从 from 到区间 [toStart, toEnd) 中的每个点都添加一条边.
func (g *RangeToRangeGraphOnTree) AddToRange(from, toStart, toEnd int32, f func(from, to int32)) {
	g.enumerateJumpDangerously(toStart, toEnd, func(id int32) { f(from, id+g.offset) })
}

// 从区间 [fromStart, fromEnd) 中的每个点到 [toStart, toEnd) 中的每个点都添加一条边.
func (g *RangeToRangeGraphOnTree) AddRangeToRange(fromStart, fromEnd, toStart, toEnd int32, f func(from, to int32)) {
	from, to := make([]int32, 0, 2), make([]int32, 0, 2)
	g.enumerateJumpDangerously(fromStart, fromEnd, func(id int32) { from = append(from, id) })
	g.enumerateJumpDangerously(toStart, toEnd, func(id int32) { to = append(to, id+g.offset) })
	for _, a := range from {
		for _, b := range to {
			f(a, b)
		}
	}
}

// 总结点数.
func (g *RangeToRangeGraphOnTree) Size() int32 { return g.n + g.offset*2 }

func (g *RangeToRangeGraphOnTree) makeDp() {
	n, log := g.n, g.log
	jump := make([][]int32, log+1)
	for k := int32(0); k < log+1; k++ {
		nums := make([]int32, n)
		jump[k] = nums
	}
	g.jump = jump
}

func (g *RangeToRangeGraphOnTree) dfsAndInitDp(cur, pre int32, f func(from, to int32)) {
	g.jump[0][cur] = pre
	for _, next := range g.tree[cur] {
		if next != pre {
			g.depth[next] = g.depth[cur] + 1
			// push down jump(0,next) to origin node.
			in := g.n + cur
			out := in + g.offset
			f(next, in)
			f(cur, in)
			f(out, next)
			f(out, cur)
			g.dfsAndInitDp(next, cur, f)
		}
	}
}

func (g *RangeToRangeGraphOnTree) updateDp() {
	n, log := g.n, g.log
	jump := g.jump
	for k := int32(0); k < log; k++ {
		for v := int32(0); v < n; v++ {
			j := jump[k][v]
			if j == -1 {
				jump[k+1][v] = -1
			} else {
				jump[k+1][v] = jump[k][j]
			}
		}
	}
}

// 遍历路径(start,target)上的所有jump.
// !要求运算幂等(idempotent).
func (g *RangeToRangeGraphOnTree) enumerateJumpDangerously(start, target int32, f func(id int32)) {
	if start == target {
		f(start)
		return
	}
	divide := func(node, ancestor int32, f func(id int32)) {
		len_ := g.depth[node] - g.depth[ancestor] + 1
		k := int32(bits.Len32(uint32(len_))) - 1
		jumpLen := len_ - (1 << k)
		from2 := g.kthAncestor(node, jumpLen)
		n := g.n
		f(k*n + n + node)
		f(k*n + n + from2)
	}
	if g.depth[start] < g.depth[target] {
		start, target = target, start
	}
	lca_ := g.lca(start, target)
	if lca_ == target {
		divide(start, lca_, f)
	} else {
		divide(start, lca_, f)
		divide(target, lca_, f)
	}
}

func (g *RangeToRangeGraphOnTree) kthAncestor(root, k int32) int32 {
	if k > g.depth[root] {
		return -1
	}
	bit := 0
	for k > 0 {
		if k&1 == 1 {
			root = g.jump[bit][root]
			if root == -1 {
				return -1
			}
		}
		bit++
		k >>= 1
	}
	return root
}

func (g *RangeToRangeGraphOnTree) lca(root1, root2 int32) int32 {
	if g.depth[root1] < g.depth[root2] {
		root1, root2 = root2, root1
	}
	root1 = g.upToDepth(root1, g.depth[root2])
	if root1 == root2 {
		return root1
	}
	for i := g.log; i >= 0; i-- {
		if a, b := g.jump[i][root1], g.jump[i][root2]; a != b {
			root1, root2 = a, b
		}
	}
	return g.jump[0][root1]
}

func (g *RangeToRangeGraphOnTree) upToDepth(root, toDepth int32) int32 {
	if toDepth >= g.depth[root] {
		return root
	}
	for i := g.log; i >= 0; i-- {
		if (g.depth[root]-toDepth)&(1<<i) > 0 {
			root = g.jump[i][root]
		}
	}
	return root
}

type UnionFindArraySimple32 struct {
	Part int32
	n    int32
	data []int32
}

func NewUnionFindArraySimple32(n int32) *UnionFindArraySimple32 {
	data := make([]int32, n)
	for i := int32(0); i < n; i++ {
		data[i] = -1
	}
	return &UnionFindArraySimple32{Part: n, n: n, data: data}
}

func (u *UnionFindArraySimple32) Union(key1, key2 int32) bool {
	root1, root2 := u.Find(key1), u.Find(key2)
	if root1 == root2 {
		return false
	}
	if u.data[root1] > u.data[root2] {
		root1, root2 = root2, root1
	}
	u.data[root1] += u.data[root2]
	u.data[root2] = int32(root1)
	u.Part--
	return true
}

func (u *UnionFindArraySimple32) Find(key int32) int32 {
	if u.data[key] < 0 {
		return key
	}
	u.data[key] = u.Find(u.data[key])
	return u.data[key]
}

func (u *UnionFindArraySimple32) GetSize(key int32) int32 {
	return -u.data[u.Find(key)]
}

const INF int = 1e18

type Neighbor struct {
	to, weight int32
}

func Dijkstra(n int32, adjList [][]Neighbor, start int32) (dist []int) {
	dist = make([]int, n)
	for i := range dist {
		dist[i] = INF
	}
	dist[start] = 0

	pq := nhp(func(a, b H) int {
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
			if cand := curDist + int(weight); cand < dist[next] {
				dist[next] = cand
				pq.Push(H{next, cand})
			}
		}
	}

	return
}

type H = struct {
	node int32
	dist int
}

// Should return a number:
//
//	negative , if a < b
//	zero     , if a == b
//	positive , if a > b
type Comparator func(a, b H) int

func nhp(comparator Comparator, nums []H) *Heap {
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
