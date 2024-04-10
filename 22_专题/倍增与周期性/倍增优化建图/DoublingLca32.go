// DividePathOnTreeBinaryLift/DoublingLca
// 倍增拆树上路径`path(from,to)`：倍增拆点将树上的一段路径拆成logn个点
//
// 一共拆分成[0,log]层，每层有n个元素.
// !jumpId = level*n + index 表示第level层的第index个元素(0<=level<log+1,0<=index<n).

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {
	P5344()

	// jump()
	// test()
}

// P5344 【XR-1】逛森林
// https://www.luogu.com.cn/problem/P5344
// 1 u1 v1 u2 v2 w : 路径u1v1上所有结点可以花费w的代价到达路径u2v2上的所有结点，如果路径不连通则无效。
// 2 u v w：结点u和v之间连接一条费用为w的无向边.如果u和v之间已经有边，则无效.
// 最后求从结点s出发，到每个节点的最小花费.
//
// 1. 倍增上的一段路径拆分成两个jump.
// !2.每个jump拆分成入点和出点，最后pushDown时，入点从儿子向父亲连边，出点从父亲向儿子连边.
func P5344() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q, start int32
	fmt.Fscan(in, &n, &q, &start)
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
	valid := make([]bool, q) // 每个操作是否有效
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

	// !每个jump拆分成入点和出点，最后pushDown时，入点从儿子向父亲连边，出点从父亲向儿子连边.
	D := NewDoublingLca32(tree, -1)
	size := D.Size()
	newGraph := make([][]Neighbour, size*2) // 入点:[0,size), 出点:[size,2*size)
	addRangeToRange := func(u1, v1, u2, v2, w int32) {
		from, to := make([]int32, 0, 2), make([]int32, 0, 2)
		D.EnumerateJumpDangerously(u1, v1, func(level, index int32) {
			id := level*D.n + index
			from = append(from, id)
		})
		D.EnumerateJumpDangerously(u2, v2, func(level, index int32) {
			id := (level*D.n + index) + size
			to = append(to, id)
		})
		for _, f := range from {
			for _, t := range to {
				newGraph[f] = append(newGraph[f], Neighbour{t, w})
			}
		}
	}

	for i := int32(0); i < q; i++ {
		op := &operations[i]
		if !valid[i] {
			continue
		}
		if op[0] == 1 {
			u1, v1, u2, v2, w := op[1], op[2], op[3], op[4], op[5]
			addRangeToRange(u1, v1, u2, v2, w)
		} else {
			u, v, w := op[1], op[2], op[3]
			newGraph[u] = append(newGraph[u], Neighbour{v + size, w})
			newGraph[v] = append(newGraph[v], Neighbour{u + size, w})
		}
	}

	D.PushDown(func(pLevel, pIndex, cLevel, cIndex1, cIndex2 int32) {
		p, c1, c2 := pLevel*D.n+pIndex, cLevel*D.n+cIndex1, cLevel*D.n+cIndex2
		// 入点从儿子向父亲连边，出点从父亲向儿子连边.
		newGraph[c1] = append(newGraph[c1], Neighbour{p + size, 0})
		newGraph[c2] = append(newGraph[c2], Neighbour{p + size, 0})
		newGraph[p] = append(newGraph[p], Neighbour{c1 + size, 0})
		newGraph[p] = append(newGraph[p], Neighbour{c2 + size, 0})
	})

	dist := DijkstraSiftHeap1(int32(len(newGraph)), newGraph, start)
	for i := size; i < size+n; i++ {
		d := dist[i]
		if d == INF {
			fmt.Fprintln(out, -1)
		} else {
			fmt.Fprintln(out, d)
		}
	}
}

func jump() {
	// https://judge.yosupo.jp/problem/jump_on_tree
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int32
	fmt.Fscan(in, &n, &q)

	tree := make([][]int32, n)
	for i := int32(0); i < n-1; i++ {
		var u, v int32
		fmt.Fscan(in, &u, &v)
		tree[u] = append(tree[u], v)
		tree[v] = append(tree[v], u)
	}
	D := NewDoublingLca32(tree, 0)

	for i := int32(0); i < q; i++ {
		var from, to, k int32
		fmt.Fscan(in, &from, &to, &k)
		fmt.Fprintln(out, D.Jump(from, to, k))
	}
}

const INF int = 1e18

type DoublingLca32 struct {
	Tree  [][]int32
	Depth []int32

	n, log, size int32
	jump         [][]int32 // 节点j向上跳2^i步的父节点
}

// root=-1表示多个根节点.
func NewDoublingLca32(tree [][]int32, root int32) *DoublingLca32 {
	n := int32(len(tree))
	depth := make([]int32, n)
	lca := &DoublingLca32{
		Tree:  tree,
		Depth: depth,
		n:     n,
		log:   int32(bits.Len32(uint32(n))) - 1,
	}
	lca.size = n * (lca.log + 1)

	lca.makeDp()
	if root == -1 {
		for i := range depth {
			depth[i] = -1
		}
		for i := int32(0); i < n; i++ {
			if lca.Depth[i] == -1 {
				lca.Depth[i] = 0
				lca.dfsAndInitDp(i, -1)
			}
		}
	} else {
		lca.dfsAndInitDp(root, -1)
	}

	lca.updateDp()
	return lca
}

// 遍历路径(start,target)上的所有jump.
// 倍增拆点，将树上的一段路径拆成logn个点.
func (lca *DoublingLca32) EnumerateJump(start, target int32, f func(level, index int32)) {
	if start == target {
		f(0, start)
		return
	}
	if lca.Depth[start] < lca.Depth[target] {
		start, target = target, start
	}
	toDepth := lca.Depth[target]
	if lca.Depth[start] > toDepth {
		for i := lca.log; i >= 0; i-- {
			if (lca.Depth[start]-toDepth)&(1<<i) > 0 {
				f(i, start)
				start = lca.jump[i][start]
			}
		}
	}
	if start == target {
		f(0, start)
		return
	}
	for i := lca.log; i >= 0; i-- {
		if a, b := lca.jump[i][start], lca.jump[i][target]; a != b {
			f(i, start)
			f(i, target)
			start, target = a, b
		}
	}
	f(0, start)
	f(0, target)
	f(0, lca.jump[0][start])
}

// 遍历路径(start,target)上的所有jump.
// !要求运算幂等(idempotent).
func (lca *DoublingLca32) EnumerateJumpDangerously(start, target int32, f func(level, index int32)) {
	if start == target {
		f(0, start)
		return
	}

	divide := func(node, ancestor int32, f func(level, index int32)) {
		len_ := lca.Depth[node] - lca.Depth[ancestor] + 1
		k := int32(bits.Len32(uint32(len_))) - 1
		jumpLen := len_ - (1 << k)
		from2 := lca.KthAncestor(node, jumpLen)
		f(k, node)
		f(k, from2)
	}

	if lca.Depth[start] < lca.Depth[target] {
		start, target = target, start
	}
	lca_ := lca.Lca(start, target)
	if lca_ == target {
		divide(start, lca_, f)
	} else {
		divide(start, lca_, f)
		divide(target, lca_, f)
	}
}

// 下推路径信息，更新答案.
// O(n*log(n)).
func (lca *DoublingLca32) PushDown(f func(pLevel, pIndex int32, cLevel, cIndex1, cIndex2 int32)) {
	n, log := lca.n, lca.log
	for k := log - 1; k >= 0; k-- {
		for i := int32(0); i < n; i++ {
			// push down jump(i,k+1) to jump(i,k) and jump(jump(i,k),k)
			if to := lca.jump[k][i]; to != -1 {
				f(k+1, i, k, i, to)
			}
		}
	}
}

func (lca *DoublingLca32) Size() int32 { return lca.size }
func (lca *DoublingLca32) Log() int32  { return lca.log }

func (lca *DoublingLca32) Lca(root1, root2 int32) int32 {
	if lca.Depth[root1] < lca.Depth[root2] {
		root1, root2 = root2, root1
	}
	root1 = lca.UpToDepth(root1, lca.Depth[root2])
	if root1 == root2 {
		return root1
	}
	for i := lca.log; i >= 0; i-- {
		if a, b := lca.jump[i][root1], lca.jump[i][root2]; a != b {
			root1, root2 = a, b
		}
	}
	return lca.jump[0][root1]
}

func (lca *DoublingLca32) Dist(root1, root2 int32, weighted bool) int32 {
	return lca.Depth[root1] + lca.Depth[root2] - 2*lca.Depth[lca.Lca(root1, root2)]
}

// 查询树节点root的第k个祖先(向上跳k步),如果不存在这样的祖先节点,返回 -1
func (lca *DoublingLca32) KthAncestor(root, k int32) int32 {
	if k > lca.Depth[root] {
		return -1
	}
	bit := 0
	for k > 0 {
		if k&1 == 1 {
			root = lca.jump[bit][root]
			if root == -1 {
				return -1
			}
		}
		bit++
		k >>= 1
	}
	return root
}

// 从 root 开始向上跳到指定深度 toDepth,toDepth<=dep[v],返回跳到的节点
func (lca *DoublingLca32) UpToDepth(root, toDepth int32) int32 {
	if toDepth >= lca.Depth[root] {
		return root
	}
	for i := lca.log; i >= 0; i-- {
		if (lca.Depth[root]-toDepth)&(1<<i) > 0 {
			root = lca.jump[i][root]
		}
	}
	return root
}

// 从start节点跳向target节点,跳过step个节点(0-indexed)
// 返回跳到的节点,如果不存在这样的节点,返回-1
func (lca *DoublingLca32) Jump(start, target, step int32) int32 {
	lca_ := lca.Lca(start, target)
	dep1, dep2, deplca := lca.Depth[start], lca.Depth[target], lca.Depth[lca_]
	dist := dep1 + dep2 - 2*deplca
	if step > dist {
		return -1
	}
	if step <= dep1-deplca {
		return lca.KthAncestor(start, step)
	}
	return lca.KthAncestor(target, dist-step)
}

func (lca *DoublingLca32) FirstTrue(start int32, predicate func(end int32) bool) int32 {
	// `LastTrue`判定条件取反，然后向上跳一层.
	if predicate(start) {
		return start
	}
	for k := lca.log; k >= 0; k-- {
		tmp := lca.jump[k][start]
		if tmp != -1 && !predicate(tmp) {
			start = tmp
		}
	}
	return lca.jump[0][start] // 不存在则返回-1
}

func (lca *DoublingLca32) LastTrue(start int32, predicate func(end int32) bool) int32 {
	if !predicate(start) {
		return -1
	}
	for k := lca.log; k >= 0; k-- {
		tmp := lca.jump[k][start]
		if tmp != -1 && predicate(tmp) {
			start = tmp
		}
	}
	return start
}

func (lca *DoublingLca32) makeDp() {
	n, log := lca.n, lca.log
	jump := make([][]int32, log+1)
	for k := int32(0); k < log+1; k++ {
		nums := make([]int32, n)
		// for i := range nums {
		// 	nums[i] = -1		// e()
		// }
		jump[k] = nums
	}
	lca.jump = jump
}

func (lca *DoublingLca32) dfsAndInitDp(cur, pre int32) {
	lca.jump[0][cur] = pre
	for _, next := range lca.Tree[cur] {
		if next != pre {
			lca.Depth[next] = lca.Depth[cur] + 1
			lca.dfsAndInitDp(next, cur)
		}
	}
}

func (lca *DoublingLca32) updateDp() {
	n, log := lca.n, lca.log
	jump := lca.jump
	for k := int32(0); k < log; k++ {
		for v := int32(0); v < n; v++ {
			j := jump[k][v]
			if j == -1 {
				jump[k+1][v] = -1 // e()
			} else {
				jump[k+1][v] = jump[k][j] // op()
			}
		}
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func max32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func min32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
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

// 采用SiftHeap加速的dijkstra算法.求出起点到各点的最短距离.
type Neighbour struct {
	next   int32
	weight int32
}

func DijkstraSiftHeap1(n int32, graph [][]Neighbour, start int32) []int {
	dist := make([]int, n)
	for i := int32(0); i < n; i++ {
		dist[i] = INF
	}
	pq := NewSiftHeap32(n, func(i, j int32) bool { return dist[i] < dist[j] })
	dist[start] = 0
	pq.Push(start)
	for pq.Size() > 0 {
		cur := pq.Pop()
		for _, e := range graph[cur] {
			next, weight := e.next, e.weight
			cand := dist[cur] + int(weight)
			if cand < dist[next] {
				dist[next] = cand
				pq.Push(next)
			}
		}
	}
	return dist
}

type SiftHeap32 struct {
	heap []int32
	pos  []int32
	less func(i, j int32) bool
	ptr  int32
}

func NewSiftHeap32(n int32, less func(i, j int32) bool) *SiftHeap32 {
	pos := make([]int32, n)
	for i := int32(0); i < n; i++ {
		pos[i] = -1
	}
	return &SiftHeap32{
		heap: make([]int32, n),
		pos:  pos,
		less: less,
	}
}

func (h *SiftHeap32) Push(i int32) {
	if h.pos[i] == -1 {
		h.pos[i] = h.ptr
		h.heap[h.ptr] = i
		h.ptr++
	}
	h._siftUp(i)
}

// 如果不存在,则返回-1.
func (h *SiftHeap32) Pop() int32 {
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
	return res
}

// 如果不存在,则返回-1.
func (h *SiftHeap32) Peek() int32 {
	if h.ptr == 0 {
		return -1
	}
	return h.heap[0]
}

func (h *SiftHeap32) Size() int32 {
	return h.ptr
}

func (h *SiftHeap32) _siftUp(i int32) {
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

func (h *SiftHeap32) _siftDown(i int32) {
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

func test() {

	//          0
	//        /   \
	//       1     2
	//      / \     \
	//     3   4     5
	//         /
	//        6

	n := 7
	edges := [][]int32{{0, 1}, {0, 2}, {1, 3}, {1, 4}, {2, 5}, {4, 6}}
	tree := make([][]int32, n)
	for _, e := range edges {
		u, v := e[0], e[1]
		tree[u] = append(tree[u], v)
		tree[v] = append(tree[v], u)
	}

	lca := NewDoublingLca32(tree, 0)
	fmt.Println(lca.Lca(3, 6)) // 1

	to1 := lca.FirstTrue(6, func(i int32) bool { return lca.Depth[i] <= -1 })
	expect[int32](to1, -1)
	to1 = lca.FirstTrue(6, func(i int32) bool { return lca.Depth[i] <= 1 })
	expect[int32](to1, 1)
	to1 = lca.FirstTrue(6, func(i int32) bool { return lca.Depth[i] <= 2 })
	expect[int32](to1, 4)
	to1 = lca.FirstTrue(6, func(i int32) bool { return lca.Depth[i] <= 3 })
	expect[int32](to1, 6)
	to1 = lca.FirstTrue(6, func(i int32) bool { return lca.Depth[i] <= 4 })
	expect[int32](to1, 6)
	to1 = lca.FirstTrue(6, func(i int32) bool { return lca.Depth[i] >= 2 })
	expect[int32](to1, 6)
	to1 = lca.FirstTrue(6, func(i int32) bool { return lca.Depth[i] >= 4 })
	expect[int32](to1, -1)

	to2 := lca.LastTrue(6, func(i int32) bool { return lca.Depth[i] >= 1 })
	expect[int32](to2, 1)
	to2 = lca.LastTrue(6, func(i int32) bool { return lca.Depth[i] >= 2 })
	expect[int32](to2, 4)
	to2 = lca.LastTrue(6, func(i int32) bool { return lca.Depth[i] >= 3 })
	expect[int32](to2, 6)
	to2 = lca.LastTrue(6, func(i int32) bool { return lca.Depth[i] >= 4 })
	expect[int32](to2, -1)
	to2 = lca.LastTrue(6, func(i int32) bool { return lca.Depth[i] <= 2 })
	expect[int32](to2, -1)
	to2 = lca.LastTrue(6, func(i int32) bool { return lca.Depth[i] <= 4 })
	expect[int32](to2, 0)
	to2 = lca.LastTrue(6, func(i int32) bool { return lca.Depth[i] <= -1 })
	expect[int32](to2, -1)

	values := make([]int32, lca.Size())
	lca.EnumerateJump(6, 3, func(level, index int32) {
		values[level*lca.n+index] += 2
	})
	lca.EnumerateJump(3, 5, func(level, index int32) {
		values[level*lca.n+index] += 3
	})
	lca.EnumerateJump(4, 4, func(level, index int32) {
		values[level*lca.n+index] += 5
	})
	lca.PushDown(func(pLevel, pIndex, cLevel, cIndex1, cIndex2 int32) {
		p, c1, c2 := pLevel*lca.n+pIndex, cLevel*lca.n+cIndex1, cLevel*lca.n+cIndex2
		values[c1] += values[p]
		values[c2] += values[p]
	})

	expected := []int32{3, 5, 3, 5, 7, 3, 2}
	for i := 0; i < n; i++ {
		expect[int32](values[i], expected[i])
	}

	// test EnumerateJumpDangerously by idempoent function
	values = make([]int32, lca.Size())
	lca.EnumerateJumpDangerously(6, 3, func(level, index int32) {
		id := level*lca.n + index
		values[id] = max32(values[id], 2)
	})
	lca.EnumerateJumpDangerously(3, 5, func(level, index int32) {
		id := level*lca.n + index
		values[id] = max32(values[id], 3)
	})
	lca.EnumerateJumpDangerously(4, 4, func(level, index int32) {
		id := level*lca.n + index
		values[id] = max32(values[id], 5)
	})
	lca.EnumerateJumpDangerously(1, 6, func(level, index int32) {
		id := level*lca.n + index
		values[id] = max32(values[id], 7)
	})

	lca.PushDown(func(pLevel, pIndex, cLevel, cIndex1, cIndex2 int32) {
		p, c1, c2 := pLevel*lca.n+pIndex, cLevel*lca.n+cIndex1, cLevel*lca.n+cIndex2
		values[c1] = max32(values[c1], values[p])
		values[c2] = max32(values[c2], values[p])
	})

	expected = []int32{3, 7, 3, 3, 7, 3, 7}
	for i := 0; i < n; i++ {
		expect[int32](values[i], expected[i])
	}

	fmt.Println("test passed")
}

func expect[S comparable](actual, expected S) {
	if actual != expected {
		panic(fmt.Sprintf("actual: %v, expected: %v", actual, expected))
	}
}
