package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// demo()
	// abc351g()
	// p4719()

	yosupo()
}

func demo() {
	//    0
	//   / \
	//  1   2
	//      |
	//      3
	tree := NewTree32(4)
	tree.AddEdge(0, 1, 1)
	tree.AddEdge(0, 2, 1)
	tree.AddEdge(1, 3, 1)
	tree.Build(0)
	stt := NewStaticTopTree(tree)
	type Cluster = any
	single := func(v int32) Cluster {
		fmt.Println("single", v)
		return 1
	}
	rake := func(x, y Cluster, u, v int32) Cluster {
		fmt.Println("rake", x, y, u, v)
		return x.(int) + y.(int)
	}
	compress := func(x, y Cluster, a, b, c int32) Cluster {
		fmt.Println("compress", x, y, a, b, c)
		return x.(int) + y.(int)
	}
	dp := stt.TreeDp(single, rake, compress)
	fmt.Println(dp)
}

// [ABC351G] Hash on Tree (DynamicTreeHash，动态树哈希)
// https://www.luogu.com.cn/problem/AT_abc351_g
// !dp[parent] = value[parent] + dp[child1] * dp[child2] * ... * dp[childn]
//
// Data 维护(乘积,和)两个值.
func abc351g() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	const MOD int = 998244353

	var n, q int32
	fmt.Fscan(in, &n, &q)
	tree := NewTree32(n)
	for i := int32(1); i < n; i++ {
		var p int32
		fmt.Fscan(in, &p)
		p--
		tree.AddDirectedEdge(p, i, 0)
	}
	tree.Build(0)
	nums := make([]int, n)
	for i := int32(0); i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	type Data = struct{ mul, add int }
	single := func(v int32) Data {
		return Data{1, nums[v]}
	}
	rake := func(x, y Data, u, v int32) Data {
		x.mul = (x.mul * y.add) % MOD
		x.add = (x.add * y.add) % MOD
		return x
	}
	compress := func(x, y Data, a, b, c int32) Data {
		mul1, add1 := x.mul, x.add
		mul2, add2 := y.mul, y.add
		return Data{(mul1 * mul2) % MOD, (mul1*add2 + add1) % MOD}
	}

	dp := NewStaticTopTreeDP[Data](NewStaticTopTree(tree))
	dp.InitDP(single, rake, compress)

	for i := int32(0); i < q; i++ {
		var v, x int
		fmt.Fscan(in, &v, &x)
		v--
		nums[v] = x
		newRes := dp.Update(int32(v), single, rake, compress)
		fmt.Fprintln(out, newRes.add)
	}
}

// TODO
// https://www.luogu.com.cn/problem/P4719
// P4719 【模板】"动态 DP" & 动态树分治
//
// 给定一棵 n 个点的树，点带点权。
// 有 m 次操作，每次操作给定 x,y，表示修改点 x 的权值为y.
// 你需要在每次操作之后求出这棵树的最大权独立集的权值大小。
//
// dp[i][0/1] 表示选/不选这个点.
// dp[i][0] = sum(max(dp[j][0], dp[j][1])), j 是 i 的儿子.
// dp[i][1] = sum(dp[j][0]) + w[i], j 是 i 的儿子.
//
// !维护簇的两个端点分别选和不选时的答案
// 对于 base cluster 端点都选答案为 -INF，否则为 0。
func p4719() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	const INF int = 1e18

	var n, q int32
	fmt.Fscan(in, &n, &q)
	weights := make([]int, n)
	for i := range weights {
		fmt.Fscan(in, &weights[i])
	}
	tree := NewTree32(n)
	for i := int32(0); i < n-1; i++ {
		var u, v int32
		fmt.Fscan(in, &u, &v)
		u, v = u-1, v-1
		tree.AddEdge(u, v, 0)
	}
	tree.Build(0)

	// TODO
	type Data = struct{ exclude, include int } // 簇的两个端点分别都选和都不选时的答案
	single := func(v int32) Data {
		return Data{exclude: 0, include: -INF}
	}
	rake := func(x, y Data, u, v int32) Data {
		x.exclude += max(y.exclude, y.include)
		x.include += max(0, y.exclude)
		return x
	}
	compress := func(x, y Data, a, b, c int32) Data {
		x.exclude += max(y.exclude, y.include)
		x.include += y.exclude
		return x
	}

	dp := NewStaticTopTreeDP[Data](NewStaticTopTree(tree))
	dp.InitDP(single, rake, compress)
	for i := int32(0); i < q; i++ {
		var x, y int
		fmt.Fscan(in, &x, &y)
		x--
		weights[x] = y
		newRes := dp.Update(int32(x), single, rake, compress)
		fmt.Fprintln(out, max(newRes.exclude, newRes.include))
	}
}

// Point Set Tree Path Composite Sum (Fixed Root)
// https://judge.yosupo.jp/problem/point_set_tree_path_composite_sum_fixed_root
// 0 p x: 将点 p 的值设为 x
// 1 ei mul add: 将边 ei 的值设为 (x -> mul * x + add)
func yosupo() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	const MOD int = 998244353

	var n, q int32
	fmt.Fscan(in, &n, &q)
	weights := make([]int, n)
	for i := range weights {
		fmt.Fscan(in, &weights[i])
	}
	mul, add := make([]int, n-1), make([]int, n-1)
	tree := NewTree32(n)
	for i := int32(0); i < n-1; i++ {
		var u, v int32
		fmt.Fscan(in, &u, &v)
		tree.AddEdge(u, v, 0)
		var m, a int
		fmt.Fscan(in, &m, &a)
		mul[i], add[i] = m, a
	}
	tree.Build(0)

	type Data = struct{ mul, add, count, res int }
	single := func(v int32) Data {
		if v == 0 {
			return Data{mul: 1, add: 0, count: 1, res: weights[v]}
		}
		eid := tree.VToE(v) // 父边的id
		m, a := mul[eid], add[eid]
		return Data{mul: m, add: a, count: 1, res: (m*weights[v] + a) % MOD}
	}
	rake := func(x, y Data, u, v int32) Data {
		return Data{mul: x.mul, add: x.add, count: x.count + y.count, res: (x.res + y.res) % MOD}
	}
	compress := func(x, y Data, a, b, c int32) Data {
		mul1, add1 := x.mul, x.add
		mul2, add2 := y.mul, y.add
		// x -> (cx+d) -> a(cx+d)+b
		aa, bb := mul1*mul2%MOD, (mul1*add2+add1)%MOD
		count := x.count + y.count
		res := (x.res + mul1*y.res + add1*y.count) % MOD
		return Data{mul: aa, add: bb, count: count, res: res}
	}

	dp := NewStaticTopTreeDP[Data](NewStaticTopTree(tree))
	dp.InitDP(single, rake, compress)
	for i := int32(0); i < q; i++ {
		var op int32
		fmt.Fscan(in, &op)
		if op == 0 {
			var p, x int
			fmt.Fscan(in, &p, &x)
			weights[p] = x
			dp.Update(int32(p), single, rake, compress)
		} else {
			var eid, m, a int
			fmt.Fscan(in, &eid, &m, &a)
			mul[eid], add[eid] = m, a
			v := tree.EToV(int32(eid))
			dp.Update(v, single, rake, compress)
		}
		x := dp.Get()
		fmt.Fprintln(out, x.res)
	}
}

// !Cluster中维护的信息不包含子树根节点的信息.
//
//	!Single(v) : v 与其父边组成的 Cluster.
//        x
//       /
//      ◯
//	!Rake(x, y, u, v) : 形成以 uv 为上下边界的 Cluster, v可能为-1.
//        x           x
//       / \   ->    /
//      ◯   ◯       ◯
//	!Compress(x, y, a, b, c) : 从上到下，依次为 (a,b] + (b,c]. a或c可能为-1.
//        x             x
//       /   		       /
//      ◯  +    ->    /
//     /             /
//    ◯             ◯

type StaticTopTreeDP[Data any] struct {
	stt *StaticTopTree
	dp  []Data
}

func NewStaticTopTreeDP[Data any](stt *StaticTopTree) *StaticTopTreeDP[Data] {
	return &StaticTopTreeDP[Data]{stt: stt}
}

func (stdp *StaticTopTreeDP[Data]) InitDP(
	single func(v int32) Data,
	rake func(x, y Data, u, v int32) Data,
	compress func(x, y Data, a, b, c int32) Data,
) Data {
	n := int32(len(stdp.stt.parent))
	stdp.dp = make([]Data, n)
	for i := int32(0); i < n; i++ {
		stdp._update(i, single, rake, compress)
	}
	return stdp.dp[n-1]
}

func (stdp *StaticTopTreeDP[Data]) Update(
	v int32,
	single func(v int32) Data,
	rake func(x, y Data, u, v int32) Data,
	compress func(x, y Data, a, b, c int32) Data,
) Data {
	for k := v; k != -1; k = stdp.stt.parent[k] {
		stdp._update(k, single, rake, compress)
	}
	return stdp.dp[len(stdp.dp)-1]
}

func (stdp *StaticTopTreeDP[Data]) Get() Data { return stdp.dp[len(stdp.dp)-1] }

func (stdp *StaticTopTreeDP[Data]) _update(
	v int32,
	single func(v int32) Data,
	rake func(x, y Data, u, v int32) Data,
	compress func(x, y Data, a, b, c int32) Data,
) {
	if 0 <= v && v < stdp.stt.n {
		stdp.dp[v] = single(v)
		return
	}
	stt := stdp.stt
	left, right := stt.leftChild[v], stt.rightChild[v]
	top, bottom := stt.topBound, stt.bottomBound
	if stt.isCompress[v] {
		a, b := top[left], bottom[left]
		d := bottom[right]
		stdp.dp[v] = compress(stdp.dp[left], stdp.dp[right], a, b, d)
	} else {
		stdp.dp[v] = rake(stdp.dp[left], stdp.dp[right], top[v], bottom[v])
	}
}

type StaticTopTree struct {
	n                             int32
	tree                          *Tree32
	parent, leftChild, rightChild []int32
	topBound, bottomBound         []int32
	isCompress                    []bool
}

func NewStaticTopTree(tree *Tree32) *StaticTopTree {
	stt := &StaticTopTree{n: tree.n, tree: tree}
	stt._build()
	return stt
}

// 获取整个树的dp值.
func (stt *StaticTopTree) TreeDp(
	single func(v int32) any,
	rake func(x, y any, u, v int32) any,
	compress func(x, y any, a, b, c int32) any,
) any {
	left, right, topBound, bottomBound, isCompress := stt.leftChild, stt.rightChild, stt.topBound, stt.bottomBound, stt.isCompress
	var dfs func(cur int32) any
	dfs = func(cur int32) any {
		if 0 <= cur && cur < stt.n {
			return single(cur)
		}
		x := dfs(left[cur])
		y := dfs(right[cur])
		if isCompress[cur] {
			return compress(x, y, topBound[left[cur]], bottomBound[left[cur]], bottomBound[right[cur]])
		}
		return rake(x, y, topBound[cur], bottomBound[cur])
	}
	return dfs(2*stt.n - 2)
}

func (stt *StaticTopTree) _build() {
	n := stt.n
	stt.parent = make([]int32, n)
	stt.leftChild = make([]int32, n)
	stt.rightChild = make([]int32, n)
	stt.topBound = make([]int32, n)
	stt.bottomBound = make([]int32, n)
	stt.isCompress = make([]bool, n)
	tree := stt.tree
	for i := int32(0); i < n; i++ {
		stt.parent[i] = -1
		stt.leftChild[i] = -1
		stt.rightChild[i] = -1
		stt.topBound[i] = tree.Parent[i]
		stt.bottomBound[i] = i
	}
	stt._buildDfs(tree.IdToNode[0])
	if int32(len(stt.parent)) != 2*n-1 {
		panic("len(stt.parent) != 2*n-1")
	}
}

func (stt *StaticTopTree) _newNode(l, r, a, b int32, c bool) int32 {
	v := int32(len(stt.parent))
	stt.parent = append(stt.parent, -1)
	stt.leftChild = append(stt.leftChild, l)
	stt.rightChild = append(stt.rightChild, r)
	stt.topBound = append(stt.topBound, a)
	stt.bottomBound = append(stt.bottomBound, b)
	stt.isCompress = append(stt.isCompress, c)
	stt.parent[l] = v
	stt.parent[r] = v
	return v
}

func (stt *StaticTopTree) _buildDfs(v int32) int32 {
	path := stt.tree.HeavyPathAt(v)
	var dfs func(l, r int32) int32
	dfs = func(l, r int32) int32 {
		if l+1 < r {
			mid := (l + r) >> 1
			x := dfs(l, mid)
			y := dfs(mid, r)
			return stt._newNode(x, y, stt.topBound[x], stt.bottomBound[y], true)
		}
		if l == 0 {
			return path[l]
		}
		pq := newHeap(func(a, b [2]int32) bool { return a[0] < b[0] }, nil)
		p := path[l-1]
		for _, to := range stt.tree.CollectLight(p) {
			x := stt._buildDfs(to)
			pq.Push([2]int32{stt.tree.SubtreeSize(to), x})
		}
		if pq.Len() == 0 {
			return path[l]
		}
		for pq.Len() >= 2 {
			item1 := pq.Pop()
			item2 := pq.Pop()
			z := stt._newNode(item1[1], item2[1], p, -1, false)
			pq.Push([2]int32{item1[0] + item2[0], z})
		}
		item := pq.Pop()
		return stt._newNode(path[l], item[1], p, path[l], false)
	}
	return dfs(0, int32(len(path)))
}

type neighbor = struct {
	to   int32
	cost int
	eid  int32
}

type Tree32 struct {
	Lid, Rid      []int32
	IdToNode      []int32
	Depth         []int32
	DepthWeighted []int
	Parent        []int32
	Head          []int32 // 重链头
	Tree          [][]neighbor
	Edges         [][2]int32
	vToE          []int32 // 节点v的父边的id
	n             int32
}

func NewTree32(n int32) *Tree32 {
	res := &Tree32{Tree: make([][]neighbor, n), Edges: make([][2]int32, 0, n-1), n: n}
	return res
}

func (t *Tree32) AddEdge(u, v int32, w int) {
	eid := int32(len(t.Edges))
	t.Tree[u] = append(t.Tree[u], neighbor{to: v, cost: w, eid: eid})
	t.Tree[v] = append(t.Tree[v], neighbor{to: u, cost: w, eid: eid})
	t.Edges = append(t.Edges, [2]int32{u, v})
}

func (t *Tree32) AddDirectedEdge(from, to int32, cost int) {
	eid := int32(len(t.Edges))
	t.Tree[from] = append(t.Tree[from], neighbor{to: to, cost: cost, eid: eid})
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
	composition := tree.GetPathDecomposition(from, to, true)
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
func (tree *Tree32) GetPathDecomposition(u, v int32, vertex bool) [][2]int32 {
	up, down := [][2]int32{}, [][2]int32{}
	for {
		if tree.Head[u] == tree.Head[v] {
			break
		}
		if tree.Lid[u] < tree.Lid[v] {
			down = append(down, [2]int32{tree.Lid[tree.Head[v]], tree.Lid[v]})
			v = tree.Parent[tree.Head[v]]
		} else {
			up = append(up, [2]int32{tree.Lid[u], tree.Lid[tree.Head[u]]})
			u = tree.Parent[tree.Head[u]]
		}
	}
	edgeInt := int32(1)
	if vertex {
		edgeInt = 0
	}
	if tree.Lid[u] < tree.Lid[v] {
		down = append(down, [2]int32{tree.Lid[u] + edgeInt, tree.Lid[v]})
	} else if tree.Lid[v]+edgeInt <= tree.Lid[u] {
		up = append(up, [2]int32{tree.Lid[u], tree.Lid[v] + edgeInt})
	}
	for i := 0; i < len(down)/2; i++ {
		down[i], down[len(down)-1-i] = down[len(down)-1-i], down[i]
	}
	return append(up, down...)
}

// 遍历路径上的 `[起点,终点)` 欧拉序 `左闭右开` 区间.
func (tree *Tree32) EnumeratePathDecomposition(u, v int32, vertex bool, f func(start, end int32)) {
	for {
		if tree.Head[u] == tree.Head[v] {
			break
		}
		if tree.Lid[u] < tree.Lid[v] {
			a, b := tree.Lid[tree.Head[v]], tree.Lid[v]
			if a > b {
				a, b = b, a
			}
			f(a, b+1)
			v = tree.Parent[tree.Head[v]]
		} else {
			a, b := tree.Lid[u], tree.Lid[tree.Head[u]]
			if a > b {
				a, b = b, a
			}
			f(a, b+1)
			u = tree.Parent[tree.Head[u]]
		}
	}

	edgeInt := int32(1)
	if vertex {
		edgeInt = 0
	}

	if tree.Lid[u] < tree.Lid[v] {
		a, b := tree.Lid[u]+edgeInt, tree.Lid[v]
		if a > b {
			a, b = b, a
		}
		f(a, b+1)
	} else if tree.Lid[v]+edgeInt <= tree.Lid[u] {
		a, b := tree.Lid[u], tree.Lid[v]+edgeInt
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

// 返回返回边 u-v 对应的 欧拉序起点编号, 1 <= eid <= n-1., 0-indexed.
func (tree *Tree32) Eid(u, v int32) int32 {
	if tree.Lid[u] > tree.Lid[v] {
		return tree.Lid[u]
	}
	return tree.Lid[v]
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

func (t *Tree32) _dfsSize(cur, pre int32) {
	size := t.Rid
	t.Parent[cur] = pre
	if pre != -1 {
		t.Depth[cur] = t.Depth[pre] + 1
	} else {
		t.Depth[cur] = 0
	}
	size[cur] = 1
	nexts := t.Tree[cur]
	for i := int32(len(nexts)) - 2; i >= 0; i-- {
		e := nexts[i+1]
		if t.Depth[e.to] == -1 {
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

func newHeap[H any](less func(a, b H) bool, nums []H) *heap[H] {
	nums = append(nums[:0:0], nums...)
	heap := &heap[H]{less: less, data: nums}
	heap.heapify()
	return heap
}

type heap[H any] struct {
	data []H
	less func(a, b H) bool
}

func (h *heap[H]) Push(value H) {
	h.data = append(h.data, value)
	h.pushUp(h.Len() - 1)
}

func (h *heap[H]) Pop() (value H) {
	if h.Len() == 0 {
		panic("heap is empty")
	}
	value = h.data[0]
	h.data[0] = h.data[h.Len()-1]
	h.data = h.data[:h.Len()-1]
	h.pushDown(0)
	return
}

func (h *heap[H]) Top() (value H) {
	value = h.data[0]
	return
}

func (h *heap[H]) Len() int { return len(h.data) }

func (h *heap[H]) heapify() {
	n := h.Len()
	for i := (n >> 1) - 1; i > -1; i-- {
		h.pushDown(i)
	}
}

func (h *heap[H]) pushUp(root int) {
	for parent := (root - 1) >> 1; parent >= 0 && h.less(h.data[root], h.data[parent]); parent = (root - 1) >> 1 {
		h.data[root], h.data[parent] = h.data[parent], h.data[root]
		root = parent
	}
}

func (h *heap[H]) pushDown(root int) {
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

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
