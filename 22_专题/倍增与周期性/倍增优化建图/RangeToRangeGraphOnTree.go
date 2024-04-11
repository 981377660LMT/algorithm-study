package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {
	// P5344()
	// CF1904F()
	P9520()
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

	R := NewRangeToRangeGraphOnTree(tree, -1)
	size := R.Size()
	newGraph := make([][]Neighbour, size)

	R.Init(func(from, to int32) { newGraph[from] = append(newGraph[from], Neighbour{to, 0}) })

	for i := int32(0); i < q; i++ {
		op := &operations[i]
		if !valid[i] {
			continue
		}
		if op[0] == 1 {
			u1, v1, u2, v2, w := op[1], op[2], op[3], op[4], op[5]
			R.AddRangeToRange(u1, v1, u2, v2, func(from, to int32) {
				newGraph[from] = append(newGraph[from], Neighbour{to, w})
			})
		} else {
			u, v, w := op[1], op[2], op[3]
			R.Add(u, v, func(from, to int32) { newGraph[from] = append(newGraph[from], Neighbour{to, w}) })
			R.Add(v, u, func(from, to int32) { newGraph[from] = append(newGraph[from], Neighbour{to, w}) })
		}
	}

	dist := DijkstraSiftHeap1(int32(len(newGraph)), newGraph, start)
	for i := int32(0); i < n; i++ {
		d := dist[i] // !出点
		if d == INF {
			d = -1
		}
		fmt.Fprint(out, d, " ")
	}
}

// Beautiful Tree
// https://www.luogu.com.cn/problem/CF1904F
// 给出一棵树，与 m 条限制，每条限制为一条路径上点权最大/小的点的编号固定。
// 请你为图分配 1∼n 的点权使得满足所有限制。
// 限制可以看成规定点点权大/于路径上的"其它点"，我们把 a 的点权小于 b 的点权的限制视作一个有向边a→b。
// 则有解当且仅当整张图没有环，拓扑排序分配即可。
// !倍增优化建图，优化成 O(nlogn) 条边。
func CF1904F() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int32
	fmt.Fscan(in, &n, &m)
	tree := make([][]int32, n)
	for i := int32(0); i < n-1; i++ {
		var u, v int32
		fmt.Fscan(in, &u, &v)
		u, v = u-1, v-1
		tree[u] = append(tree[u], v)
		tree[v] = append(tree[v], u)
	}

	D := NewRangeToRangeGraphOnTree(tree, 0)
	size := D.Size()
	newGraph := make([][]int32, size)
	indeg := make([]int32, size)
	addEdge := func(from, to int32) {
		newGraph[from] = append(newGraph[from], to)
		indeg[to]++
	}
	D.Init(addEdge)

	addPointToRangeWithoutPoint := func(point, from, to int32) {
		from1, to1, from2, to2 := SplitPath(from, to, point, D.depth, D.kthAncestor, D.lca)
		if from1 != -1 && to1 != -1 {
			D.AddToRange(point, from1, to1, addEdge)
		}
		if from2 != -1 && to2 != -1 {
			D.AddToRange(point, from2, to2, addEdge)
		}
	}

	addRangeToPointWithoutPoint := func(from, to, point int32) {
		from1, to1, from2, to2 := SplitPath(from, to, point, D.depth, D.kthAncestor, D.lca)
		if from1 != -1 && to1 != -1 {
			D.AddFromRange(from1, to1, point, addEdge)
		}
		if from2 != -1 && to2 != -1 {
			D.AddFromRange(from2, to2, point, addEdge)
		}
	}

	for i := int32(0); i < m; i++ {
		var op, a, b, c int32
		fmt.Fscan(in, &op, &a, &b, &c)
		a, b, c = a-1, b-1, c-1
		// 点c的点权是路径a到b上的最小值
		if op == 1 {
			addPointToRangeWithoutPoint(c, a, b)
		} else {
			// 点c的点权是路径a到b上的最大值
			addRangeToPointWithoutPoint(a, b, c)
		}
	}

	queue := make([]int32, 0, size)
	for i := int32(0); i < size; i++ {
		if indeg[i] == 0 {
			queue = append(queue, i)
		}
	}
	topoOrder := make([]int32, 0, n)
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		if cur < n {
			topoOrder = append(topoOrder, cur)
		}
		for _, next := range newGraph[cur] {
			indeg[next]--
			if indeg[next] == 0 {
				queue = append(queue, next)
			}
		}
	}

	for _, d := range indeg {
		if d > 0 {
			fmt.Fprintln(out, -1)
			return
		}
	}

	res := make([]int32, n)
	alloc := int32(1)
	for _, cur := range topoOrder {
		res[cur] = alloc
		alloc++
	}

	for _, r := range res {
		fmt.Fprint(out, r, " ")
	}
}

// P9520 [JOISC2022] 监狱
// https://www.luogu.com.cn/problem/P9520
// https://www.cnblogs.com/5k-sync-closer/p/18035300
// 对于n个点的树，有m条"起点与终点各不相同"的行进路线形如 si→ti，允许从某个点移动至相邻点
// !问能否在不存在某个点所在人数 >1的情况下完成所有行进路线。
// 1<=m<=n<=1.2e5
//
// 若 A 路径的起点在 B 路径上，则 A 必须比 B 先走，
// 若 A 路径的终点在 B 路径上，则 B 必须比 A 先走。
//
// !1.如果 A 的起点在 B 的路径上，那么 A 必须先于 B 走 =>
// 把每条路径向其起点连边，然后把每条路径除起点外的点向这条路径连边，
// 此时 A 连向 A 的起点，而 A 路径的起点在 B 路径上，所以 A 的起点连向 B。
// !2.如果 A 的终点在 B 的路径上，那么 B 必须先于 A 走 =>
// 把每个终点向其路径连边，然后把每条路径向这条路径除终点外的点连边，
// 此时 A 路径的终点在 B 路径上，所以 B 连向 A 的终点，而 A 的终点连向 A。
// TODO
func P9520() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	solve := func(tree [][]int32, routes [][2]int32) bool {
		D := NewRangeToRangeGraphOnTree(tree, 0)
		size := D.Size()
		newGraph := make([][]int32, size)
		indeg := make([]int32, size)
		addEdge := func(from, to int32) {
			newGraph[from] = append(newGraph[from], to)
			indeg[to]++
		}
		D.Init(addEdge)

		for _, route := range routes {
			start, end := route[0], route[1]
			// 每条路径向其起点连边
			D.AddFromRange(start, end, start, addEdge)
			// 每条路径除起点外的点向这条路径连边
			from1, to1, from2, to2 := SplitPath(start, end, start, D.depth, D.kthAncestor, D.lca)
			if from1 != -1 && to1 != -1 {
				D.AddRangeToRange(from1, to1, start, end, addEdge)
			}
			if from2 != -1 && to2 != -1 {
				D.AddRangeToRange(from2, to2, start, end, addEdge)
			}

			// 每个终点向其路径连边
			D.AddToRange(start, end, end, addEdge)
			// 每条路径向这条路径除终点外的点连边
			from1, to1, from2, to2 = SplitPath(start, end, end, D.depth, D.kthAncestor, D.lca)
			if from1 != -1 && to1 != -1 {
				D.AddRangeToRange(start, end, from1, to1, addEdge)
			}
			if from2 != -1 && to2 != -1 {
				D.AddRangeToRange(start, end, from2, to2, addEdge)
			}
		}

		queue := make([]int32, 0, size)
		for i := int32(0); i < size; i++ {
			if indeg[i] == 0 {
				queue = append(queue, i)
			}
		}
		for len(queue) > 0 {
			cur := queue[0]
			queue = queue[1:]
			for _, next := range newGraph[cur] {
				indeg[next]--
				if indeg[next] == 0 {
					queue = append(queue, next)
				}
			}
		}

		for _, d := range indeg {
			if d > 0 {
				return false
			}
		}
		return true
	}

	var T int32
	fmt.Fscan(in, &T)
	for i := int32(0); i < T; i++ {
		var n int32
		fmt.Fscan(in, &n)
		tree := make([][]int32, n)
		for i := int32(0); i < n-1; i++ {
			var u, v int32
			fmt.Fscan(in, &u, &v)
			u--
			v--
			tree[u] = append(tree[u], v)
			tree[v] = append(tree[v], u)
		}

		var m int32
		fmt.Fscan(in, &m)
		queries := make([][2]int32, m)
		for i := int32(0); i < m; i++ {
			var s, t int32
			fmt.Fscan(in, &s, &t)
			s--
			t--
			queries[i] = [2]int32{s, t}
		}

		ok := solve(tree, queries)
		if ok {
			fmt.Fprintln(out, "Yes")
		} else {
			fmt.Fprintln(out, "No")
		}
	}
}

type RangeToRangeGraphOnTree struct {
	tree           [][]int32
	depth          []int32
	n, log, offset int32 // !底层真实点：[0,n)，倍增入点：[n,n+offset)，倍增出点：[n+offset,n+2*offset).
	root           int32
	jump           [][]int32 // 节点j向上跳2^i步的父节点
}

// root为-1表示无根.
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

// 总结点数.
func (g *RangeToRangeGraphOnTree) Size() int32 { return g.n + g.offset*2 }

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

// 从路径 [fromStart, fromEnd] 中的每个点到 to 都添加一条边.
func (g *RangeToRangeGraphOnTree) AddFromRange(fromStart, fromEnd, to int32, f func(from, to int32)) {
	g.enumerateJumpDangerously(fromStart, fromEnd, func(id int32) { f(id, to) })
}

// 从 from 到路径 [toStart, toEnd] 中的每个点都添加一条边.
func (g *RangeToRangeGraphOnTree) AddToRange(from, toStart, toEnd int32, f func(from, to int32)) {
	g.enumerateJumpDangerously(toStart, toEnd, func(id int32) { f(from, id+g.offset) })
}

// 从路径 [fromStart, fromEnd] 中的每个点到 [toStart, toEnd] 中的每个点都添加一条边.
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
	// push down jump(0,cur).
	in := g.n + cur
	out := in + g.offset
	f(cur, in)
	f(out, cur)
	for _, next := range g.tree[cur] {
		if next != pre {
			g.depth[next] = g.depth[cur] + 1
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
		f(start + g.n)
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

func (g *RangeToRangeGraphOnTree) jumpFn(start, target, step int32) int32 {
	lca_ := g.lca(start, target)
	dep1, dep2, deplca := g.depth[start], g.depth[target], g.depth[lca_]
	dist := dep1 + dep2 - 2*deplca
	if step > dist {
		return -1
	}
	if step <= dep1-deplca {
		return g.kthAncestor(start, step)
	}
	return g.kthAncestor(target, dist-step)
}

func SplitPathByJumpFn(
	from, to, separator int32,
	jumpFn func(start, target, step int32) int32,
) (from1, to1, from2, to2 int32) {
	from1, to1, from2, to2 = -1, -1, -1, -1
	if from == to {
		return
	}
	if separator == from {
		from2 = jumpFn(from, to, 1)
		to2 = to
		return
	}
	if separator == to {
		from1 = from
		to1 = jumpFn(to, from, 1)
		return
	}
	from1 = from
	to1 = jumpFn(separator, from, 1)
	from2 = jumpFn(separator, to, 1)
	to2 = to
	return
}

func SplitPath(
	from, to int32, separator int32,
	depth []int32,
	kthAncestorFn func(node, k int32) int32,
	lcaFn func(node1, node2 int32) int32,
) (from1, to1, from2, to2 int32) {
	from1, to1, from2, to2 = -1, -1, -1, -1
	if from == to {
		return
	}

	down, top := from, to
	swapped := false
	if depth[down] < depth[top] {
		down, top = top, down
		swapped = true
	}

	lca := lcaFn(from, to)
	if lca == top {
		// down和top在一条链上.
		if separator == down {
			from2 = kthAncestorFn(separator, 1)
			to2 = top
		} else if separator == top {
			from1 = down
			to1 = kthAncestorFn(down, depth[down]-depth[separator]-1)
		} else {
			from1 = down
			to1 = kthAncestorFn(down, depth[down]-depth[separator]-1)
			from2 = kthAncestorFn(separator, 1)
			to2 = top
		}
	} else {
		// down和top在lca两个子树上.
		if separator == down {
			from2 = kthAncestorFn(separator, 1)
			to2 = top
		} else if separator == top {
			from1 = down
			to1 = kthAncestorFn(separator, 1)
		} else {
			var jump1, jump2 int32
			if separator == lca {
				jump1 = kthAncestorFn(down, depth[down]-depth[separator]-1)
				jump2 = kthAncestorFn(top, depth[top]-depth[separator]-1)
			} else if lcaFn(separator, down) == separator {
				jump1 = kthAncestorFn(down, depth[down]-depth[separator]-1)
				jump2 = kthAncestorFn(separator, 1)
			} else {
				jump1 = kthAncestorFn(separator, 1)
				jump2 = kthAncestorFn(top, depth[top]-depth[separator]-1)
			}
			from1 = down
			to1 = jump1
			from2 = jump2
			to2 = top
		}
	}

	if swapped {
		from1, to1, from2, to2 = to2, from2, to1, from1
	}
	return
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
