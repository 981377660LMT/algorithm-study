package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
)

func main() {
	// acwing2534()
	cf852I()

}

// 树上莫队，点权.
func acwing2534() {
	// 2534. 树上计数2
	// https://www.acwing.com/problem/content/description/2536/
	// 对每个查询，求u到v的路径上顶点种类数
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int32
	fmt.Fscan(in, &n, &q)

	nums := make([]int32, n) // 顶点权值
	for i := int32(0); i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	newNums, origin := DiscretizeId32(nums)

	tree := NewTree32(n)
	for i := int32(0); i < n-1; i++ {
		var u, v int32
		fmt.Fscan(in, &u, &v)
		u--
		v--
		tree.AddEdge(u, v, 1)
	}
	tree.Build(0)

	mo := NewMoOnTreeV2(tree)
	for i := int32(0); i < q; i++ {
		var u, v int32
		fmt.Fscan(in, &u, &v)
		u--
		v--
		mo.AddQuery(u, v)
	}

	res := make([]int32, q)
	counter, kind := make([]int32, len(origin)), int32(0)
	add := func(i int32) {
		counter[newNums[i]]++
		if counter[newNums[i]] == 1 {
			kind++
		}
	}
	remove := func(i int32) {
		counter[newNums[i]]--
		if counter[newNums[i]] == 0 {
			kind--
		}
	}
	query := func(qid int32) {
		res[qid] = kind
	}
	init := func() {
		add(0)
	}

	mo.CalcVertex(init, add, add, remove, remove, query)
	for _, v := range res {
		fmt.Fprintln(out, v)
	}
}

// 树上莫队，点权.
// https://www.luogu.com.cn/problem/CF852I
// 给定n个节点的树，每个节点有一个男生或女生。每个人都有一个喜欢的数字。
// 然后m次询问，每次询问树上从a到b的路径中，有多少对男女喜欢的数字相同
// n,q<=1e5.
func cf852I() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int32
	fmt.Fscan(in, &n)
	A, B := make([]int32, n), make([]int32, n) // !A: 1/0表示男/女, B: 喜欢的数字
	for i := int32(0); i < n; i++ {
		fmt.Fscan(in, &A[i])
	}
	for i := int32(0); i < n; i++ {
		fmt.Fscan(in, &B[i])
	}
	D := NewDictionary[int32]()
	for i := int32(0); i < n; i++ {
		B[i] = int32(D.Id(B[i]))
	}

	tree := NewTree32(n)
	for i := int32(0); i < n-1; i++ {
		var u, v int32
		fmt.Fscan(in, &u, &v)
		tree.AddEdge(u-1, v-1, 1)
	}
	tree.Build(0)

	mo := NewMoOnTreeV2(tree)
	var q int32
	fmt.Fscan(in, &q)
	for i := int32(0); i < q; i++ {
		var u, v int32
		fmt.Fscan(in, &u, &v)
		mo.AddQuery(u-1, v-1)
	}

	res := make([]int, q)
	counter := make([][]int, 2)
	for i := range counter {
		counter[i] = make([]int, D.Size())
	}

	cur := 0

	add := func(v int32) {
		counter[A[v]][B[v]]++
		cur += counter[A[v]^1][B[v]]
	}

	remove := func(v int32) {
		counter[A[v]][B[v]]--
		cur -= counter[A[v]^1][B[v]]
	}

	query := func(qid int32) {
		res[qid] = cur
	}

	init := func() { add(0) }

	mo.CalcVertex(init, add, add, remove, remove, query)
	for _, v := range res {
		fmt.Fprintln(out, v)
	}
}

// 树上莫队，边权.
// 2846. 边权重均等查询
// https://leetcode.cn/problems/minimum-edge-weight-equilibrium-queries-in-a-tree/description/
// 在一次操作中，你可以选择树上的任意一条边，并将其权重更改为任意值。
// 请你找出使从 ai 到 bi 路径上每条边的权重相等所需的 最小操作次数
func minOperationsQueries(n int, edges [][]int, queries [][]int) []int {
	D := NewDictionary[int]()
	for i := range edges {
		edges[i][2] = D.Id(edges[i][2])
	}
	tree := NewTree32(int32(n))
	for _, edge := range edges {
		u, v, w := edge[0], edge[1], edge[2]
		tree.AddEdge(int32(u), int32(v), w)
	}
	tree.Build(0)

	mo := NewMoOnTreeV2(tree)
	for _, q := range queries {
		mo.AddQuery(int32(q[0]), int32(q[1]))
	}

	res := make([]int, len(queries))
	weightCounter := make([]int32, D.Size())
	sl := NewSortedListRangeBlock32Simple(int32(n - 1))
	edgeCount := int32(0)

	init := func() {}

	add := func(from int32, to int32) {
		eid := tree.Eid(from, to)
		weight := edges[eid][2]
		preCount := weightCounter[weight]
		weightCounter[weight] = preCount + 1
		if preCount > 0 {
			sl.Remove(preCount)
		}
		sl.Add(preCount + 1)
		edgeCount++
	}
	remove := func(from int32, to int32) {
		eid := tree.Eid(from, to)
		weight := edges[eid][2]
		preCount := weightCounter[weight]
		weightCounter[weight] = preCount - 1
		sl.Remove(preCount)
		if preCount > 1 {
			sl.Add(preCount - 1)
		}
		edgeCount--
	}
	query := func(qid int32) {
		if sl.Len() > 0 {
			max_ := sl.Max()
			res[qid] = int(edgeCount - max_)
		}
	}

	mo.CalcEdge(init, add, add, remove, remove, query)
	return res
}

// 树上莫队.
type MoOnTreeV2 struct {
	tree          *Tree32
	lefts, rights []int32
}

func NewMoOnTreeV2(tree *Tree32) *MoOnTreeV2 {
	return &MoOnTreeV2{tree: tree}
}

func (m *MoOnTreeV2) AddQuery(u, v int32) {
	if m.tree.Lid[u] > m.tree.Lid[v] {
		u, v = v, u
	}
	m.lefts = append(m.lefts, m.tree.ELid(u)+1)
	m.rights = append(m.rights, m.tree.ELid(v)+1)
}

func (m *MoOnTreeV2) CalcVertex(
	init func(), // !仅有根节点的情况(点权时需要考虑)
	addL func(v int32), addR func(v int32), // 路径开头/结尾添加节点v
	removeL func(v int32), removeR func(v int32), // 路径开头/结尾删除节点v
	query func(qid int32),
) {
	tree := m.tree
	n := tree.n
	order := getMoOrder(m.lefts, m.rights)
	from, to, idx := make([]int32, 2*n), make([]int32, 2*n), make([]int32, 2*n)
	visited := make([]bool, n)
	path := newDeque32(n)
	path.Append(0)

	for v := int32(0); v < n; v++ {
		a, b := tree.ELid(v), tree.ERid(v)
		from[a], to[a] = tree.Parent[v], v
		from[b], to[b] = v, tree.Parent[v]
		idx[a], idx[b] = v, v
	}

	flipLeft := func(i int32) {
		a, b, c := from[i], to[i], idx[i]
		if !visited[c] {
			v := path.Front() ^ a ^ b
			path.AppendLeft(v)
			addL(v)
		} else {
			v := path.Front()
			path.PopLeft()
			removeL(v)
		}
		visited[c] = !visited[c]
	}

	flipRight := func(i int32) {
		a, b, c := from[i], to[i], idx[i]
		if !visited[c] {
			v := path.Back() ^ a ^ b
			path.Append(v)
			addR(v)
		} else {
			v := path.Back()
			path.Pop()
			removeR(v)
		}
		visited[c] = !visited[c]
	}

	init()

	l, r := int32(1), int32(1)
	for _, idx := range order {
		left, right := m.lefts[idx], m.rights[idx]
		for l > left {
			l--
			flipLeft(l)
		}
		for r < right {
			flipRight(r)
			r++
		}
		for l < left {
			flipLeft(l)
			l++
		}
		for r > right {
			r--
			flipRight(r)
		}
		query(idx)
	}
}

func (m *MoOnTreeV2) CalcEdge(
	init func(),
	addL func(from, to int32), addR func(from, to int32), // 路径开头/结尾添加边(from,to)
	removeL func(from, to int32), removeR func(from, to int32), // 路径开头/结尾删除边(from,to)
	query func(qid int32),
) {
	tree := m.tree
	n := tree.n
	order := getMoOrder(m.lefts, m.rights)
	from, to, idx := make([]int32, 2*n), make([]int32, 2*n), make([]int32, 2*n)
	visited := make([]bool, n)
	path := newDeque32(n)
	path.Append(0)

	for v := int32(0); v < n; v++ {
		a, b := tree.ELid(v), tree.ERid(v)
		from[a], to[a] = tree.Parent[v], v
		from[b], to[b] = v, tree.Parent[v]
		idx[a], idx[b] = v, v
	}

	flipLeft := func(i int32) {
		a, b, c := from[i], to[i], idx[i]
		if !visited[c] {
			v := path.Front() ^ a ^ b
			path.AppendLeft(v)
			addL(v, v^a^b)
		} else {
			v := path.Front()
			path.PopLeft()
			removeL(v, v^a^b)
		}
		visited[c] = !visited[c]
	}

	flipRight := func(i int32) {
		a, b, c := from[i], to[i], idx[i]
		if !visited[c] {
			v := path.Back() ^ a ^ b
			path.Append(v)
			addR(v^a^b, v)
		} else {
			v := path.Back()
			path.Pop()
			removeR(v^a^b, v)
		}
		visited[c] = !visited[c]
	}

	init()

	l, r := int32(1), int32(1)
	for _, idx := range order {
		left, right := m.lefts[idx], m.rights[idx]
		for l > left {
			l--
			flipLeft(l)
		}
		for r < right {
			flipRight(r)
			r++
		}
		for l < left {
			flipLeft(l)
			l++
		}
		for r > right {
			r--
			flipRight(r)
		}
		query(idx)
	}
}

func getMoOrder(lefts, rights []int32) []int32 {
	n := int32(1)
	for i := 0; i < len(lefts); i++ {
		n = max32(n, lefts[i])
		n = max32(n, rights[i])
	}
	q := len(lefts)
	if q == 0 {
		return []int32{}
	}
	bs := int32(math.Sqrt(3) * float64(n) / math.Sqrt(2*float64(q)))
	bs = max32(bs, 1)
	order := make([]int32, q)
	for i := 0; i < q; i++ {
		order[i] = int32(i)
	}
	belong := make([]int32, q)
	for i := 0; i < q; i++ {
		belong[i] = lefts[i] / bs
	}
	sort.Slice(order, func(a, b int) bool {
		oa, ob := order[a], order[b]
		bida, bidb := belong[oa], belong[ob]
		if bida != bidb {
			return bida < bidb
		}
		if bida&1 == 1 {
			return rights[oa] > rights[ob]
		}
		return rights[oa] < rights[ob]
	})

	{
		cost := func(a, b int32) int32 {
			oa, ob := order[a], order[b]
			return abs32(lefts[oa]-lefts[ob]) + abs32(rights[oa]-rights[ob])
		}
		for k := int32(0); k < int32(q-5); k++ {
			if cost(k, k+2)+cost(k+1, k+3) < cost(k, k+1)+cost(k+2, k+3) {
				order[k+1], order[k+2] = order[k+2], order[k+1]
			}
			if cost(k, k+3)+cost(k+1, k+4) < cost(k, k+1)+cost(k+3, k+4) {
				order[k+1], order[k+3] = order[k+3], order[k+1]
			}
		}
	}

	return order
}

type neighbor = struct {
	to   int32
	eid  int32
	cost int
}

type Tree32 struct {
	n             int32
	Lid, Rid      []int32
	IdToNode      []int32
	Depth         []int32
	Parent        []int32
	Head          []int32 // 重链头
	vToE          []int32 // 节点v的父边的id
	DepthWeighted []int
	Tree          [][]neighbor
	Edges         [][2]int32
}

func NewTree32(n int32) *Tree32 {
	res := &Tree32{n: n, Tree: make([][]neighbor, n), Edges: make([][2]int32, 0, n-1)}
	return res
}

func (t *Tree32) AddEdge(u, v int32, w int) {
	eid := int32(len(t.Edges))
	t.Tree[u] = append(t.Tree[u], neighbor{to: v, eid: eid, cost: w})
	t.Tree[v] = append(t.Tree[v], neighbor{to: u, eid: eid, cost: w})
	t.Edges = append(t.Edges, [2]int32{u, v})
}

func (t *Tree32) AddDirectedEdge(from, to int32, cost int) {
	eid := int32(len(t.Edges))
	t.Tree[from] = append(t.Tree[from], neighbor{to: to, eid: eid, cost: cost})
	t.Edges = append(t.Edges, [2]int32{from, to})
}

func (t *Tree32) Build(root int32) {
	if root != -1 && int32(len(t.Edges)) != t.n-1 {
		panic("edges count != n-1")
	}
	n := t.n
	t.Lid = make([]int32, n)
	t.Rid = make([]int32, n)
	t.IdToNode = make([]int32, n)
	t.Depth = make([]int32, n)
	t.DepthWeighted = make([]int, n)
	t.Parent = make([]int32, n)
	t.Head = make([]int32, n)
	t.vToE = make([]int32, n)
	for i := int32(0); i < n; i++ {
		t.Depth[i] = -1
		t.Head[i] = root
		t.vToE[i] = -1
	}
	if root != -1 {
		t._dfsSize(root, -1)
		time := int32(0)
		t._dfsHld(root, &time)
	} else {
		time := int32(0)
		for i := int32(0); i < n; i++ {
			if t.Depth[i] == -1 {
				t._dfsSize(i, -1)
				t._dfsHld(i, &time)
			}
		}
	}
}

// 从v开始沿着重链向下收集节点.
func (t *Tree32) HeavyPathAt(v int32) []int32 {
	path := []int32{v}
	for {
		a := path[len(path)-1]
		for _, e := range t.Tree[a] {
			if e.to != t.Parent[a] && t.Head[e.to] == v {
				path = append(path, e.to)
				break
			}
		}
		if path[len(path)-1] == a {
			break
		}
	}
	return path
}

// 返回重儿子，如果没有返回 -1.
func (t *Tree32) HeavyChild(v int32) int32 {
	k := t.Lid[v] + 1
	if k == t.n {
		return -1
	}
	w := t.IdToNode[k]
	if t.Parent[w] == v {
		return w
	}
	return -1
}

// 从v开始向上走k步.
func (t *Tree32) KthAncestor(v, k int32) int32 {
	if k > t.Depth[v] {
		return -1
	}
	for {
		u := t.Head[v]
		if t.Lid[v]-k >= t.Lid[u] {
			return t.IdToNode[t.Lid[v]-k]
		}
		k -= t.Lid[v] - t.Lid[u] + 1
		v = t.Parent[u]
	}
}

func (t *Tree32) Lca(u, v int32) int32 {
	for {
		if t.Lid[u] > t.Lid[v] {
			u, v = v, u
		}
		if t.Head[u] == t.Head[v] {
			return u
		}
		v = t.Parent[t.Head[v]]
	}
}

func (t *Tree32) LcaRooted(u, v, root int32) int32 {
	return t.Lca(u, v) ^ t.Lca(u, root) ^ t.Lca(v, root)
}

func (t *Tree32) Dist(a, b int32) int32 {
	c := t.Lca(a, b)
	return t.Depth[a] + t.Depth[b] - 2*t.Depth[c]
}

func (t *Tree32) DistWeighted(a, b int32) int {
	c := t.Lca(a, b)
	return t.DepthWeighted[a] + t.DepthWeighted[b] - 2*t.DepthWeighted[c]
}

// c 是否在 p 的子树中.c和p不能相等.
func (t *Tree32) InSubtree(c, p int32) bool {
	return t.Lid[p] <= t.Lid[c] && t.Lid[c] < t.Rid[p]
}

// 从 a 开始走 k 步到 b.
func (t *Tree32) Jump(a, b, k int32) int32 {
	if k == 1 {
		if a == b {
			return -1
		}
		if t.InSubtree(b, a) {
			return t.KthAncestor(b, t.Depth[b]-t.Depth[a]-1)
		}
		return t.Parent[a]
	}
	c := t.Lca(a, b)
	dac := t.Depth[a] - t.Depth[c]
	dbc := t.Depth[b] - t.Depth[c]
	if k > dac+dbc {
		return -1
	}
	if k <= dac {
		return t.KthAncestor(a, k)
	}
	return t.KthAncestor(b, dac+dbc-k)
}

func (t *Tree32) SubtreeSize(v int32) int32 {
	return t.Rid[v] - t.Lid[v]
}

func (t *Tree32) SubtreeSizeRooted(v, root int32) int32 {
	if v == root {
		return t.n
	}
	x := t.Jump(v, root, 1)
	if t.InSubtree(v, x) {
		return t.Rid[v] - t.Lid[v]
	}
	return t.n - t.Rid[x] + t.Lid[x]
}

func (t *Tree32) CollectChild(v int32) []int32 {
	var res []int32
	for _, e := range t.Tree[v] {
		if e.to != t.Parent[v] {
			res = append(res, e.to)
		}
	}
	return res
}

// 收集与 v 相邻的轻边.
func (t *Tree32) CollectLight(v int32) []int32 {
	var res []int32
	skip := true
	for _, e := range t.Tree[v] {
		if e.to != t.Parent[v] {
			if !skip {
				res = append(res, e.to)
			}
			skip = false
		}
	}
	return res
}

func (tree *Tree32) RestorePath(from, to int32) []int32 {
	res := []int32{}
	composition := tree.GetPathDecomposition(from, to, 0)
	for _, e := range composition {
		a, b := e[0], e[1]
		if a <= b {
			for i := a; i <= b; i++ {
				res = append(res, tree.IdToNode[i])
			}
		} else {
			for i := a; i >= b; i-- {
				res = append(res, tree.IdToNode[i])
			}
		}
	}
	return res
}

// 返回沿着`路径顺序`的 [起点,终点] 的 欧拉序 `左闭右闭` 数组.
//
//	!eg:[[2 0] [4 4]] 沿着路径顺序但不一定沿着欧拉序.
func (tree *Tree32) GetPathDecomposition(u, v int32, edge int32) [][2]int32 {
	up, down := [][2]int32{}, [][2]int32{}
	lid, head, parent := tree.Lid, tree.Head, tree.Parent
	for {
		if head[u] == head[v] {
			break
		}
		if lid[u] < lid[v] {
			down = append(down, [2]int32{lid[head[v]], lid[v]})
			v = parent[head[v]]
		} else {
			up = append(up, [2]int32{lid[u], lid[head[u]]})
			u = parent[head[u]]
		}
	}
	if lid[u] < lid[v] {
		down = append(down, [2]int32{lid[u] + edge, lid[v]})
	} else if lid[v]+edge <= lid[u] {
		up = append(up, [2]int32{lid[u], lid[v] + edge})
	}
	for i := 0; i < len(down)/2; i++ {
		down[i], down[len(down)-1-i] = down[len(down)-1-i], down[i]
	}
	return append(up, down...)
}

// 遍历路径上的 `[起点,终点)` 欧拉序 `左闭右开` 区间.
func (tree *Tree32) EnumeratePathDecomposition(u, v int32, edge int32, f func(start, end int32)) {
	head, lid, parent := tree.Head, tree.Lid, tree.Parent
	for {
		if head[u] == head[v] {
			break
		}
		if lid[u] < lid[v] {
			a, b := lid[head[v]], lid[v]
			if a > b {
				a, b = b, a
			}
			f(a, b+1)
			v = parent[head[v]]
		} else {
			a, b := lid[u], lid[head[u]]
			if a > b {
				a, b = b, a
			}
			f(a, b+1)
			u = parent[head[u]]
		}
	}
	if lid[u] < lid[v] {
		a, b := lid[u]+edge, lid[v]
		if a > b {
			a, b = b, a
		}
		f(a, b+1)
	} else if lid[v]+edge <= lid[u] {
		a, b := lid[u], lid[v]+edge
		if a > b {
			a, b = b, a
		}
		f(a, b+1)
	}
}

// 返回 root 的欧拉序区间, 左闭右开, 0-indexed.
func (tree *Tree32) Id(root int32) (int32, int32) {
	return tree.Lid[root], tree.Rid[root]
}

// 返回返回边 u-v 对应的 边id.
func (tree *Tree32) Eid(u, v int32) int32 {
	if tree.Parent[u] != v {
		u, v = v, u
	}
	return tree.vToE[u]
}

// 点v对应的父边的边id.如果v是根节点则返回-1.
func (tre *Tree32) VToE(v int32) int32 {
	return tre.vToE[v]
}

// 第i条边对应的深度更深的那个节点.
func (tree *Tree32) EToV(i int32) int32 {
	u, v := tree.Edges[i][0], tree.Edges[i][1]
	if tree.Parent[u] == v {
		return u
	}
	return v
}

func (tree *Tree32) ELid(u int32) int32 {
	return 2*tree.Lid[u] - tree.Depth[u]
}

func (tree *Tree32) ERid(u int32) int32 {
	return 2*tree.Rid[u] - tree.Depth[u] - 1
}

func (t *Tree32) _dfsSize(cur, pre int32) {
	size := t.Rid
	t.Parent[cur] = pre
	if pre == -1 {
		t.Depth[cur] = 0
	} else {
		t.Depth[cur] = t.Depth[pre] + 1
	}
	size[cur] = 1
	nexts := t.Tree[cur]
	for i := int32(len(nexts)) - 2; i >= 0; i-- {
		if t.Depth[nexts[i+1].to] == -1 {
			nexts[i], nexts[i+1] = nexts[i+1], nexts[i]
		}
	}
	hldSize := int32(0)
	for i, e := range nexts {
		to := e.to
		if t.Depth[to] == -1 {
			t.DepthWeighted[to] = t.DepthWeighted[cur] + e.cost
			t.vToE[to] = e.eid
			t._dfsSize(to, cur)
			size[cur] += size[to]
			if size[to] > hldSize {
				hldSize = size[to]
				if i != 0 {
					nexts[0], nexts[i] = nexts[i], nexts[0]
				}
			}
		}
	}
}

func (t *Tree32) _dfsHld(cur int32, times *int32) {
	t.Lid[cur] = *times
	*times++
	t.Rid[cur] += t.Lid[cur]
	t.IdToNode[t.Lid[cur]] = cur
	heavy := true
	for _, e := range t.Tree[cur] {
		to := e.to
		if t.Depth[to] > t.Depth[cur] {
			if heavy {
				t.Head[to] = t.Head[cur]
			} else {
				t.Head[to] = to
			}
			heavy = false
			t._dfsHld(to, times)
		}
	}
}

// 路径 [a,b] 与 [c,d] 的交集.
// 如果为空则返回 {-1,-1}，如果只有一个交点则返回 {x,x}，如果有两个交点则返回 {x,y}.
func (t *Tree32) PathIntersection(a, b, c, d int32) (int32, int32) {
	ab := t.Lca(a, b)
	ac := t.Lca(a, c)
	ad := t.Lca(a, d)
	bc := t.Lca(b, c)
	bd := t.Lca(b, d)
	cd := t.Lca(c, d)
	x := ab ^ ac ^ bc // meet(a,b,c)
	y := ab ^ ad ^ bd // meet(a,b,d)
	if x != y {
		return x, y
	}
	z := ac ^ ad ^ cd
	if x != z {
		x = -1
	}
	return x, x
}

type deque32 struct{ left, right []int32 }

func newDeque32(initCapacity int32) *deque32 {
	return &deque32{make([]int32, 0, 1+initCapacity/2), make([]int32, 0, 1+initCapacity/2)}
}

func (q *deque32) Empty() bool {
	return len(q.left) == 0 && len(q.right) == 0
}

func (q *deque32) Len() int {
	return len(q.left) + len(q.right)
}

func (q *deque32) AppendLeft(v int32) {
	q.left = append(q.left, v)
}

func (q *deque32) Append(v int32) {
	q.right = append(q.right, v)
}

func (q *deque32) PopLeft() (v int32) {
	if len(q.left) > 0 {
		q.left, v = q.left[:len(q.left)-1], q.left[len(q.left)-1]
	} else {
		v, q.right = q.right[0], q.right[1:]
	}
	return
}

func (q *deque32) Pop() (v int32) {
	if len(q.right) > 0 {
		q.right, v = q.right[:len(q.right)-1], q.right[len(q.right)-1]
	} else {
		v, q.left = q.left[0], q.left[1:]
	}
	return
}

func (q *deque32) Front() int32 {
	if len(q.left) > 0 {
		return q.left[len(q.left)-1]
	}
	return q.right[0]
}

func (q *deque32) Back() int32 {
	if len(q.right) > 0 {
		return q.right[len(q.right)-1]
	}
	return q.left[0]
}

// 0 <= i < q.Len()
func (q *deque32) At(i int) int32 {
	if i < len(q.left) {
		return q.left[len(q.left)-1-i]
	}
	return q.right[i-len(q.left)]
}

func (q *deque32) Clear() {
	q.left = q.left[:0]
	q.right = q.right[:0]
}

func (q *deque32) ForEach(f func(v int32)) {
	for i := len(q.left) - 1; i >= 0; i-- {
		f(q.left[i])
	}
	for i := 0; i < len(q.right); i++ {
		f(q.right[i])
	}
}

func max(a, b int) int {
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

func max32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

// 将nums中的元素进行重映射，返回新的数组和对应的原始值.
// origin[newNums[i]] == nums[i]
func DiscretizeId32(nums []int32) (newNums []int32, origin []int32) {
	pool := make(map[int32]int32)
	newNums = make([]int32, len(nums))
	origin = make([]int32, 0, len(newNums))
	for i, v := range nums {
		if id, ok := pool[v]; ok {
			newNums[i] = id
		} else {
			id := int32(len(origin))
			pool[v] = id
			newNums[i] = id
			origin = append(origin, v)
		}
	}
	origin = origin[:len(origin):len(origin)]
	return
}

func abs32(x int32) int32 {
	if x < 0 {
		return -x
	}
	return x
}

type SortedListRangeBlock32 struct {
	_blockSize  int32   // 每个块的大小.
	_len        int32   // 所有数的个数.
	_counter    []int32 // 每个数出现的次数.
	_blockCount []int32 // 每个块的数的个数.
	_belong     []int32 // 每个数所在的块.
	_blockSum   []int   // 每个块的和.
}

// 值域分块模拟SortedList.
// `O(1)`add/remove，`O(sqrt(n))`查询.
// 一般配合莫队算法使用.
//
//	max:值域的最大值.0 <= max <= 1e6.
//	iterable:初始值.
func NewSortedListRangeBlock32Simple(max int32, nums ...int32) *SortedListRangeBlock32 {
	max += 5
	size := int32(math.Sqrt(float64(max)))
	count := 1 + (max / size)
	sl := &SortedListRangeBlock32{
		_blockSize:  size,
		_counter:    make([]int32, max),
		_blockCount: make([]int32, count),
		_belong:     make([]int32, max),
		_blockSum:   make([]int, count),
	}
	for i := int32(0); i < max; i++ {
		sl._belong[i] = i / size
	}
	if len(nums) > 0 {
		for _, num := range nums {
			sl.Add(num)
		}
	}
	return sl
}

// O(1).
func (sl *SortedListRangeBlock32) Add(value int32) {
	sl._counter[value]++
	pos := sl._belong[value]
	sl._blockCount[pos]++
	sl._blockSum[pos] += int(value)
	sl._len++
}

// O(1).
func (sl *SortedListRangeBlock32) Remove(value int32) {
	sl._counter[value]--
	pos := sl._belong[value]
	sl._blockCount[pos]--
	sl._blockSum[pos] -= int(value)
	sl._len--
}

// O(1).
func (sl *SortedListRangeBlock32) Discard(value int32) bool {
	if !sl.Has(value) {
		return false
	}
	sl.Remove(value)
	return true
}

// O(1).
func (sl *SortedListRangeBlock32) Has(value int32) bool {
	return sl._counter[value] > 0
}

// O(sqrt(n)).
func (sl *SortedListRangeBlock32) At(index int32) int32 {
	if index < 0 {
		index += sl._len
	}
	if index < 0 || index >= sl._len {
		panic(fmt.Sprintf("index out of range: %d", index))
	}
	for i := int32(0); i < int32(len(sl._blockCount)); i++ {
		count := sl._blockCount[i]
		if index < count {
			num := i * sl._blockSize
			for {
				numCount := sl._counter[num]
				if index < numCount {
					return num
				}
				index -= numCount
				num++
			}
		}
		index -= count
	}
	panic("unreachable")
}

// O(sqrt(n)).
func (sl *SortedListRangeBlock32) Pop(index int32) int32 {
	if index < 0 {
		index += sl._len
	}
	if index < 0 || index >= sl._len {
		panic(fmt.Sprintf("index out of range: %d", index))
	}
	value := sl.At(index)
	sl.Remove(value)
	return value
}

func (sl *SortedListRangeBlock32) Len() int32 {
	return sl._len
}

func (sl *SortedListRangeBlock32) Min() int32 {
	return sl.At(0)
}

func (sl *SortedListRangeBlock32) Max() int32 {
	if sl._len == 0 {
		panic("empty")
	}

	for i := int32(len(sl._blockCount) - 1); i >= 0; i-- {
		if sl._blockCount[i] == 0 {
			continue
		}
		num := (i+1)*sl._blockSize - 1
		for {
			if sl._counter[num] > 0 {
				return num
			}
			num--
		}
	}

	panic("unreachable")
}

// 返回索引在`kth`处的元素的`value`,以及该元素是`value`中的第几个(`index`).
func (sl *SortedListRangeBlock32) _findKth(kth int32) (value, index int32) {
	for i := int32(0); i < int32(len(sl._blockCount)); i++ {
		count := sl._blockCount[i]
		if kth < count {
			num := i * sl._blockSize
			for {
				numCount := sl._counter[num]
				if kth < numCount {
					return num, kth
				}
				kth -= numCount
				num++
			}
		}
		kth -= count
	}

	panic("unreachable")
}

type Dictionary[V comparable] struct {
	_idToValue []V
	_valueToId map[V]int32
}

// A dictionary that maps values to unique ids.
func NewDictionary[V comparable]() *Dictionary[V] {
	return &Dictionary[V]{
		_valueToId: map[V]int32{},
	}
}
func (d *Dictionary[V]) Id(value V) int {
	res, ok := d._valueToId[value]
	if ok {
		return int(res)
	}
	id := len(d._idToValue)
	d._idToValue = append(d._idToValue, value)
	d._valueToId[value] = int32(id)
	return id
}
func (d *Dictionary[V]) Value(id int) V {
	return d._idToValue[id]
}
func (d *Dictionary[V]) Has(value V) bool {
	_, ok := d._valueToId[value]
	return ok
}
func (d *Dictionary[V]) Size() int {
	return len(d._idToValue)
}
