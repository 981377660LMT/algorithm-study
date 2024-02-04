package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"sort"
)

const INF int = 1e18

func main() {
	// ColouredMountainHut()
	// CF613D()
	// Yuki3407()
	P2495()
}

func demo() {
	n := 5
	rawTree := NewTree(n)
	rawTree.AddEdge(0, 1, 1)
	rawTree.AddEdge(0, 2, 2)
	rawTree.AddEdge(1, 3, 3)
	rawTree.AddEdge(1, 4, 4)
	rawTree.Build(0)

	isCritical := make([]bool, n)
	criticals := []int{0, 1, 4}
	for _, v := range criticals {
		isCritical[v] = true
	}
	rawId, newTree := CompressTree(rawTree, criticals, false)
	inCriticals := make([]bool, len(rawId)) // 虚树上的某个节点是否在criticals中
	for i := 0; i < len(rawId); i++ {
		inCriticals[i] = isCritical[rawId[i]]
	}
	fmt.Println(rawId, newTree.Dist(0, 1, false))
	for _, v := range criticals {
		isCritical[v] = false
	}
}

// https://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=0439
// 给定一棵树,每个点有一个颜色。
// 对每一种颜色相同的点，求出每个点到其他点距离的最小值。保证每种颜色的点至少有两个。
// !虚树上求点对距离.
func ColouredMountainHut() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	colors := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &colors[i])
	}
	tree := NewTree(n)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		tree.AddEdge(u-1, v-1, 1)
	}
	tree.Build(0)

	groupByColor := make(map[int][]int)
	for i, c := range colors {
		groupByColor[c] = append(groupByColor[c], i)
	}

	res := make([]int, n)
	for i := 0; i < n; i++ {
		res[i] = INF
	}

	isCritical := make([]bool, n)
	for _, criticals := range groupByColor {
		for _, v := range criticals {
			isCritical[v] = true
		}

		rawId, newTree := CompressTree(tree, criticals, false)
		adjList := newTree.Tree
		starts := make([]int, 0, len(criticals)) // !获取critials 在新树上的编号
		for i := 0; i < len(rawId); i++ {
			if isCritical[rawId[i]] {
				starts = append(starts, i)
			}
		}
		minDistToOther, _ := MinDistToOther(adjList, starts)
		for i := 0; i < len(starts); i++ {
			node := rawId[starts[i]]
			res[node] = min(res[node], minDistToOther[i])
		}

		for _, v := range criticals {
			isCritical[v] = false
		}
	}

	for _, v := range res {
		fmt.Fprintln(out, v)
	}
}

// 给定一棵树，每次询问给定 k个特殊点，找出尽量少的非特殊点使得删去这些点后特殊点两两不连通。∑k≤n.
// 如果无法使得特殊点两两不连通，输出-1.
// https://codeforces.com/problemset/problem/613/D
func CF613D() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	tree := NewTree(n)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		tree.AddEdge(u-1, v-1, 1)
	}
	tree.Build(0)

	// !dp[i] 表示子树中保留i个关键点时的最小删除点数
	// ①:如果一个点被标记了,那么就要把他所有子树里有标记点的儿子都去掉
	// ②:如果一个点没有被标记,但是这个点有两颗以上的子树里有标记点，那么这个点就要去掉,然后这棵子树就没有可标记点了
	// ③:如果一个点子树里只有一个/没有标记点，那么就标记点的贡献挪到这个点上面来
	solve := func(adjList [][][2]int, inCriticals []bool) int {
		var dfs func(cur, pre int) [2]int // (zero, one)
		dfs = func(cur, pre int) [2]int {
			removeCost := 1
			dp := [2]int{INF, INF}
			if inCriticals[cur] {
				removeCost = INF // 无法删除
				dp[1] = 0
			} else {
				dp[0] = 0
			}

			for _, e := range adjList[cur] {
				next := e[0]
				if next == pre {
					continue
				}
				subDp := dfs(next, cur)
				ndp := [2]int{INF, INF}
				for a := 0; a < 2; a++ {
					for b := 0; b < 2; b++ {
						if a == 1 && b == 1 { // !不能>=2个关键点
							continue
						}
						ndp[a+b] = min(ndp[a+b], dp[a]+subDp[b])
					}
				}
				dp = ndp
				removeCost += min(subDp[0], subDp[1])
			}

			dp[0] = min(dp[0], removeCost)
			return dp
		}

		dp := dfs(0, -1)
		res := min(dp[0], dp[1])
		if res >= INF {
			res = -1
		}
		return res
	}

	isCritical := make([]bool, n)
	var q int
	fmt.Fscan(in, &q)
	for i := 0; i < q; i++ {
		var k int
		fmt.Fscan(in, &k)
		criticals := make([]int, k)
		for j := 0; j < k; j++ {
			var v int
			fmt.Fscan(in, &v)
			v--
			criticals[j] = v
			isCritical[v] = true
		}

		nodes := append(criticals[:0:0], criticals...)
		for _, v := range criticals {
			if v != 0 {
				nodes = append(nodes, tree.Parent[v]) // !父节点加进来
			}
		}
		nodes = unique(nodes)

		rawId, newTree := CompressTree(tree, nodes, true)
		m := len(rawId)
		inCriticals := make([]bool, m) // !压缩后的树中的节点是否在points中
		for i := 0; i < m; i++ {
			inCriticals[i] = isCritical[rawId[i]]
		}
		fmt.Println(solve(newTree.Tree, inCriticals))
		for _, v := range criticals {
			isCritical[v] = false
		}
	}
}

// P2495 [SDOI2011] 消耗战
// 给定一棵树，每次询问给定 k个特殊点，需要断掉一些边使得从根节点无法到达任何特殊点，求最小需要断掉的边权之和。∑k≤2n.
// https://www.luogu.com.cn/problem/P2495
func P2495() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	tree := NewTree(n)
	for i := 0; i < n-1; i++ {
		var u, v, w int
		fmt.Fscan(in, &u, &v, &w)
		u, v = u-1, v-1
		tree.AddEdge(u, v, w)
	}
	tree.Build(0)

	lca := NewLCADoubling(tree.Tree, []int{0})

	// dp[i]表示i和以i为根的子树中的关键点都不相连的最小代价
	// 如果子节点是关键点，dp[i] += minWeight[i][child]
	// 如果子节点不是关键点，dp[i] += min(dp[child], minWeight[i][child])
	solve := func(adjList [][][2]int, inCriticals []bool, rawId []int) int {
		var dfs func(cur, pre int) int
		dfs = func(cur, pre int) int {
			res := 0
			for _, e := range adjList[cur] {
				next := e[0]
				if next == pre {
					continue
				}
				nextRes := dfs(next, cur)
				minWeight := lca.QueryMinWeight(rawId[cur], rawId[next], true)
				if inCriticals[next] {
					res += minWeight
				} else {
					res += min(nextRes, minWeight)
				}
			}
			return res
		}
		return dfs(0, -1)
	}

	var q int
	fmt.Fscan(in, &q)
	isRawIdCritical := make([]bool, n)
	for i := 0; i < q; i++ {
		var k int
		fmt.Fscan(in, &k)
		criticals := make([]int, k)
		for j := 0; j < k; j++ {
			var p int
			fmt.Fscan(in, &p)
			p--
			criticals[j] = p
			isRawIdCritical[p] = true
		}
		criticals = append(criticals, 0) // !构建虚树时加上根节点
		rawId, newTree := CompressTree(tree, criticals, true)
		inCriticals := make([]bool, len(rawId))
		for i := 0; i < len(rawId); i++ {
			inCriticals[i] = isRawIdCritical[rawId[i]]
		}
		fmt.Println(solve(newTree.Tree, inCriticals, rawId))
		for _, v := range criticals {
			isRawIdCritical[v] = false
		}
	}
}

// No.901 K-ary εxtrεεmε
// https://yukicoder.me/problems/3407
// !给定q个查询,求虚树(最小的包含指定点集的连通子图)组成的的边权之和
// !求虚树边权之和.
//
// 第二种解法是按照dfs序排序，求树链并, https://yukicoder.me/submissions/756376
func Yuki3407() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	tree := NewTree(n)
	for i := 0; i < n-1; i++ {
		var u, v, w int
		fmt.Fscan(in, &u, &v, &w)
		tree.AddEdge(u, v, w)
	}
	tree.Build(0)

	var q int
	fmt.Fscan(in, &q)
	for i := 0; i < q; i++ {
		var k int
		fmt.Fscan(in, &k)
		criticals := make([]int, k)
		for j := 0; j < k; j++ {
			fmt.Fscan(in, &criticals[j])
		}

		_, newTree := CompressTree(tree, criticals, true)
		adjList := newTree.Tree
		res := 0
		for _, nexts := range adjList {
			for _, e := range nexts {
				res += e[1]
			}
		}
		fmt.Fprintln(out, res)
	}
}

// 返回树压缩后保留的节点编号和新的树.
// !新的树保留了原树的边权.
func CompressTree(rawTree *Tree, nodes []int, directed bool) (rawId []int, newTree *Tree) {
	rawId = append(nodes[:0:0], nodes...)
	sort.Slice(rawId, func(i, j int) bool { return rawTree.LID[rawId[i]] < rawTree.LID[rawId[j]] })
	n := len(rawId)
	for i := 0; i < n; i++ {
		j := i + 1
		if j == n {
			j = 0
		}
		rawId = append(rawId, rawTree.LCA(rawId[i], rawId[j]))
	}
	// remainNodes = append(remainNodes, rawTree.IdToNode[0])
	sort.Slice(rawId, func(i, j int) bool { return rawTree.LID[rawId[i]] < rawTree.LID[rawId[j]] })

	unique := func(a []int) []int {
		visited := make(map[int]struct{})
		newNums := []int{}
		for _, v := range a {
			if _, ok := visited[v]; !ok {
				visited[v] = struct{}{}
				newNums = append(newNums, v)
			}
		}
		return newNums
	}

	rawId = unique(rawId)
	n = len(rawId)
	newTree = NewTree(n)

	stack := []int{0}
	for i := 1; i < n; i++ {
		for {
			p := rawId[stack[len(stack)-1]]
			v := rawId[i]
			if rawTree.IsInSubtree(v, p) {
				break
			}
			stack = stack[:len(stack)-1]
		}
		p := rawId[stack[len(stack)-1]]
		v := rawId[i]
		d := rawTree.DepthWeighted[v] - rawTree.DepthWeighted[p]
		newTree.AddDirectedEdge(stack[len(stack)-1], i, d)
		if !directed {
			newTree.AddDirectedEdge(i, stack[len(stack)-1], d)
		}
		stack = append(stack, i)
	}
	newTree.Build(0)
	return
}

type Tree struct {
	Tree                 [][][2]int // (next, weight)
	Depth, DepthWeighted []int
	Parent               []int
	LID, RID             []int // 欧拉序[in,out)
	IdToNode             []int
	top, heavySon        []int
	timer                int
}

func NewTree(n int) *Tree {
	tree := make([][][2]int, n)
	lid := make([]int, n)
	rid := make([]int, n)
	IdToNode := make([]int, n)
	top := make([]int, n)   // 所处轻/重链的顶点（深度最小），轻链的顶点为自身
	depth := make([]int, n) // 深度
	depthWeighted := make([]int, n)
	parent := make([]int, n)   // 父结点
	heavySon := make([]int, n) // 重儿子
	for i := range parent {
		parent[i] = -1
	}

	return &Tree{
		Tree:          tree,
		Depth:         depth,
		DepthWeighted: depthWeighted,
		Parent:        parent,
		LID:           lid,
		RID:           rid,
		IdToNode:      IdToNode,
		top:           top,
		heavySon:      heavySon,
	}
}

// 添加无向边 u-v, 边权为w.
func (tree *Tree) AddEdge(u, v, w int) {
	tree.Tree[u] = append(tree.Tree[u], [2]int{v, w})
	tree.Tree[v] = append(tree.Tree[v], [2]int{u, w})
}

// 添加有向边 u->v, 边权为w.
func (tree *Tree) AddDirectedEdge(u, v, w int) {
	tree.Tree[u] = append(tree.Tree[u], [2]int{v, w})
}

// root:0-based
//
//	当root设为-1时，会从0开始遍历未访问过的连通分量
func (tree *Tree) Build(root int) {
	if root != -1 {
		tree.build(root, -1, 0, 0)
		tree.markTop(root, root)
	} else {
		for i := 0; i < len(tree.Tree); i++ {
			if tree.Parent[i] == -1 {
				tree.build(i, -1, 0, 0)
				tree.markTop(i, i)
			}
		}
	}
}

// 返回 root 的欧拉序区间, 左闭右开, 0-indexed.
func (tree *Tree) Id(root int) (int, int) {
	return tree.LID[root], tree.RID[root]
}

// 返回返回边 u-v 对应的 欧拉序起点编号, 1 <= eid <= n-1., 0-indexed.
func (tree *Tree) Eid(u, v int) int {
	if tree.LID[u] > tree.LID[v] {
		return tree.LID[u]
	}
	return tree.LID[v]
}

func (tree *Tree) LCA(u, v int) int {
	for {
		if tree.LID[u] > tree.LID[v] {
			u, v = v, u
		}
		if tree.top[u] == tree.top[v] {
			return u
		}
		v = tree.Parent[tree.top[v]]
	}
}

func (tree *Tree) RootedLCA(u, v int, root int) int {
	return tree.LCA(u, v) ^ tree.LCA(u, root) ^ tree.LCA(v, root)
}

func (tree *Tree) RootedParent(u int, root int) int {
	return tree.Jump(u, root, 1)
}

func (tree *Tree) Dist(u, v int, weighted bool) int {
	if weighted {
		return tree.DepthWeighted[u] + tree.DepthWeighted[v] - 2*tree.DepthWeighted[tree.LCA(u, v)]
	}
	return tree.Depth[u] + tree.Depth[v] - 2*tree.Depth[tree.LCA(u, v)]
}

// k: 0-based
//
//	如果不存在第k个祖先，返回-1
//	kthAncestor(root,0) == root
func (tree *Tree) KthAncestor(root, k int) int {
	if k > tree.Depth[root] {
		return -1
	}
	for {
		u := tree.top[root]
		if tree.LID[root]-k >= tree.LID[u] {
			return tree.IdToNode[tree.LID[root]-k]
		}
		k -= tree.LID[root] - tree.LID[u] + 1
		root = tree.Parent[u]
	}
}

// 从 from 节点跳向 to 节点,跳过 step 个节点(0-indexed)
//
//	返回跳到的节点,如果不存在这样的节点,返回-1
func (tree *Tree) Jump(from, to, step int) int {
	if step == 1 {
		if from == to {
			return -1
		}
		if tree.IsInSubtree(to, from) {
			return tree.KthAncestor(to, tree.Depth[to]-tree.Depth[from]-1)
		}
		return tree.Parent[from]
	}
	c := tree.LCA(from, to)
	dac := tree.Depth[from] - tree.Depth[c]
	dbc := tree.Depth[to] - tree.Depth[c]
	if step > dac+dbc {
		return -1
	}
	if step <= dac {
		return tree.KthAncestor(from, step)
	}
	return tree.KthAncestor(to, dac+dbc-step)
}

func (tree *Tree) CollectChild(root int) []int {
	res := []int{}
	for _, e := range tree.Tree[root] {
		next := e[0]
		if next != tree.Parent[root] {
			res = append(res, next)
		}
	}
	return res
}

// 返回沿着`路径顺序`的 [起点,终点] 的 欧拉序 `左闭右闭` 数组.
//
//	!eg:[[2 0] [4 4]] 沿着路径顺序但不一定沿着欧拉序.
func (tree *Tree) GetPathDecomposition(u, v int, vertex bool) [][2]int {
	up, down := [][2]int{}, [][2]int{}
	for {
		if tree.top[u] == tree.top[v] {
			break
		}
		if tree.LID[u] < tree.LID[v] {
			down = append(down, [2]int{tree.LID[tree.top[v]], tree.LID[v]})
			v = tree.Parent[tree.top[v]]
		} else {
			up = append(up, [2]int{tree.LID[u], tree.LID[tree.top[u]]})
			u = tree.Parent[tree.top[u]]
		}
	}
	edgeInt := 1
	if vertex {
		edgeInt = 0
	}
	if tree.LID[u] < tree.LID[v] {
		down = append(down, [2]int{tree.LID[u] + edgeInt, tree.LID[v]})
	} else if tree.LID[v]+edgeInt <= tree.LID[u] {
		up = append(up, [2]int{tree.LID[u], tree.LID[v] + edgeInt})
	}
	for i := 0; i < len(down)/2; i++ {
		down[i], down[len(down)-1-i] = down[len(down)-1-i], down[i]
	}
	return append(up, down...)
}

// 遍历路径上的 `[起点,终点)` 欧拉序 `左闭右开` 区间.
func (tree *Tree) EnumeratePathDecomposition(u, v int, vertex bool, f func(start, end int)) {
	for {
		if tree.top[u] == tree.top[v] {
			break
		}
		if tree.LID[u] < tree.LID[v] {
			a, b := tree.LID[tree.top[v]], tree.LID[v]
			if a > b {
				a, b = b, a
			}
			f(a, b+1)
			v = tree.Parent[tree.top[v]]
		} else {
			a, b := tree.LID[u], tree.LID[tree.top[u]]
			if a > b {
				a, b = b, a
			}
			f(a, b+1)
			u = tree.Parent[tree.top[u]]
		}
	}

	edgeInt := 1
	if vertex {
		edgeInt = 0
	}

	if tree.LID[u] < tree.LID[v] {
		a, b := tree.LID[u]+edgeInt, tree.LID[v]
		if a > b {
			a, b = b, a
		}
		f(a, b+1)
	} else if tree.LID[v]+edgeInt <= tree.LID[u] {
		a, b := tree.LID[u], tree.LID[v]+edgeInt
		if a > b {
			a, b = b, a
		}
		f(a, b+1)
	}
}

func (tree *Tree) GetPath(u, v int) []int {
	res := []int{}
	composition := tree.GetPathDecomposition(u, v, true)
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

// 以root为根时,结点v的子树大小.
func (tree *Tree) SubSize(v, root int) int {
	if root == -1 {
		return tree.RID[v] - tree.LID[v]
	}
	if v == root {
		return len(tree.Tree)
	}
	x := tree.Jump(v, root, 1)
	if tree.IsInSubtree(v, x) {
		return tree.RID[v] - tree.LID[v]
	}
	return len(tree.Tree) - tree.RID[x] + tree.LID[x]
}

// child 是否在 root 的子树中 (child和root不能相等)
func (tree *Tree) IsInSubtree(child, root int) bool {
	return tree.LID[root] <= tree.LID[child] && tree.LID[child] < tree.RID[root]
}

// 寻找以 start 为 top 的重链 ,heavyPath[-1] 即为重链底端节点.
func (tree *Tree) GetHeavyPath(start int) []int {
	heavyPath := []int{start}
	cur := start
	for tree.heavySon[cur] != -1 {
		cur = tree.heavySon[cur]
		heavyPath = append(heavyPath, cur)
	}
	return heavyPath
}

// 结点v的重儿子.如果没有重儿子,返回-1.
func (tree *Tree) GetHeavyChild(v int) int {
	k := tree.LID[v] + 1
	if k == len(tree.Tree) {
		return -1
	}
	w := tree.IdToNode[k]
	if tree.Parent[w] == v {
		return w
	}
	return -1
}

func (tree *Tree) ELID(u int) int {
	return 2*tree.LID[u] - tree.Depth[u]
}

func (tree *Tree) ERID(u int) int {
	return 2*tree.RID[u] - tree.Depth[u] - 1
}

func (tree *Tree) build(cur, pre, dep, dist int) int {
	subSize, heavySize, heavySon := 1, 0, -1
	for _, e := range tree.Tree[cur] {
		next, weight := e[0], e[1]
		if next != pre {
			nextSize := tree.build(next, cur, dep+1, dist+weight)
			subSize += nextSize
			if nextSize > heavySize {
				heavySize, heavySon = nextSize, next
			}
		}
	}
	tree.Depth[cur] = dep
	tree.DepthWeighted[cur] = dist
	tree.heavySon[cur] = heavySon
	tree.Parent[cur] = pre
	return subSize
}

func (tree *Tree) markTop(cur, top int) {
	tree.top[cur] = top
	tree.LID[cur] = tree.timer
	tree.IdToNode[tree.timer] = cur
	tree.timer++
	heavySon := tree.heavySon[cur]
	if heavySon != -1 {
		tree.markTop(heavySon, top)
		for _, e := range tree.Tree[cur] {
			next := e[0]
			if next != heavySon && next != tree.Parent[cur] {
				tree.markTop(next, next)
			}
		}
	}
	tree.RID[cur] = tree.timer
}

type LCADoubling struct {
	Tree          [][][2]int
	Depth         []int32
	DepthWeighted []int
	n             int
	bitLen        int
	dp            [][]int32 // 节点j向上跳2^i步的父节点
	dpWeight2     [][]int   // 节点j向上跳2^i步经过的最小边权
}

func NewLCADoubling(tree [][][2]int, roots []int) *LCADoubling {
	n := len(tree)
	depth := make([]int32, n)
	for i := range depth {
		depth[i] = -1
	}
	lca := &LCADoubling{
		Tree:          tree,
		Depth:         depth,
		DepthWeighted: make([]int, n),
		n:             n,
		bitLen:        bits.Len(uint(n)),
	}
	lca.dp, lca.dpWeight2 = makeDp(lca)
	for _, root := range roots {
		lca.dfsAndInitDp(int32(root), -1, 0, 0)
	}
	lca.fillDp()
	return lca
}

// 查询树节点两点的最近公共祖先
func (lca *LCADoubling) QueryLCA(root1, root2 int) int {
	if lca.Depth[root1] < lca.Depth[root2] {
		root1, root2 = root2, root1
	}
	root1 = lca.UpToDepth(root1, int(lca.Depth[root2]))
	if root1 == root2 {
		return root1
	}
	root132, root232 := int32(root1), int32(root2)
	for i := lca.bitLen - 1; i >= 0; i-- {
		if lca.dp[i][root132] != lca.dp[i][root232] {
			root132 = lca.dp[i][root132]
			root232 = lca.dp[i][root232]
		}
	}
	return int(lca.dp[0][root132])
}

// 查询树节点两点间距离
//
//	weighted: 是否将边权计入距离
func (lca *LCADoubling) QueryDist(root1, root2 int, weighted bool) int {
	if weighted {
		return lca.DepthWeighted[root1] + lca.DepthWeighted[root2] - 2*lca.DepthWeighted[lca.QueryLCA(root1, root2)]
	}
	return int(lca.Depth[root1] + lca.Depth[root2] - 2*lca.Depth[lca.QueryLCA(root1, root2)])
}

// 查询树节点两点路径上最小边权(倍增的时候维护其他属性)
//
//	isEdge 为true表示查询路径上边权,为false表示查询路径上点权
func (lca *LCADoubling) QueryMinWeight(root1, root2 int, isEdge bool) int {
	res := INF
	if lca.Depth[root1] < lca.Depth[root2] {
		root1, root2 = root2, root1
	}
	toDepth := lca.Depth[root2]
	root132, root232 := int32(root1), int32(root2)
	for i := lca.bitLen - 1; i >= 0; i-- { // upToDepth
		if (lca.Depth[root132]-toDepth)&(1<<i) > 0 {
			res = min(res, lca.dpWeight2[i][root132])
			root132 = lca.dp[i][root132]
		}
	}
	if root132 == root232 {
		return res
	}
	for i := lca.bitLen - 1; i >= 0; i-- {
		if lca.dp[i][root132] != lca.dp[i][root232] {
			res = min(res, min(lca.dpWeight2[i][root132], lca.dpWeight2[i][root232]))
			root132 = lca.dp[i][root132]
			root232 = lca.dp[i][root232]
		}
	}
	res = min(res, min(lca.dpWeight2[0][root132], lca.dpWeight2[0][root232]))
	if !isEdge {
		lca_ := lca.dp[0][root132]
		res = min(res, lca.dpWeight2[0][lca_])
	}
	return res
}

// 查询树节点root的第k个祖先(向上跳k步),如果不存在这样的祖先节点,返回 -1
func (lca *LCADoubling) QueryKthAncestor(root, k int) int {
	root32 := int32(root)
	if k > int(lca.Depth[root32]) {
		return -1
	}
	bit := 0
	for k > 0 {
		if k&1 == 1 {
			root32 = lca.dp[bit][root32]
			if root32 == -1 {
				return -1
			}
		}
		bit++
		k >>= 1
	}
	return int(root32)
}

// 从 root 开始向上跳到指定深度 toDepth,toDepth<=dep[v],返回跳到的节点
func (lca *LCADoubling) UpToDepth(root, toDepth int) int {
	toDepth32 := int32(toDepth)
	if toDepth32 >= lca.Depth[root] {
		return root
	}
	root32 := int32(root)
	for i := lca.bitLen - 1; i >= 0; i-- {
		if (lca.Depth[root32]-toDepth32)&(1<<i) > 0 {
			root32 = lca.dp[i][root32]
		}
	}
	return int(root32)
}

// 从start节点跳向target节点,跳过step个节点(0-indexed)
// 返回跳到的节点,如果不存在这样的节点,返回-1
func (lca *LCADoubling) Jump(start, target, step int) int {
	lca_ := lca.QueryLCA(start, target)
	dep1, dep2, deplca := lca.Depth[start], lca.Depth[target], lca.Depth[lca_]
	dist := int(dep1 + dep2 - 2*deplca)
	if step > dist {
		return -1
	}
	if step <= int(dep1-deplca) {
		return lca.QueryKthAncestor(start, step)
	}
	return lca.QueryKthAncestor(target, dist-step)
}

func (lca *LCADoubling) dfsAndInitDp(cur, pre, dep int32, dist int) {
	lca.Depth[cur] = dep
	lca.dp[0][cur] = pre
	lca.DepthWeighted[cur] = dist
	for _, e := range lca.Tree[cur] {
		next, weight := int32(e[0]), e[1]
		if next != pre {
			lca.dpWeight2[0][next] = weight
			lca.dfsAndInitDp(next, cur, dep+1, dist+weight)
		}
	}
}

func makeDp(lca *LCADoubling) (dp [][]int32, dpWeight2 [][]int) {
	dp, dpWeight2 = make([][]int32, lca.bitLen), make([][]int, lca.bitLen)
	for i := 0; i < lca.bitLen; i++ {
		dp[i], dpWeight2[i] = make([]int32, lca.n), make([]int, lca.n)
		for j := 0; j < lca.n; j++ {
			dp[i][j] = -1
			dpWeight2[i][j] = INF
		}
	}
	return
}

func (lca *LCADoubling) fillDp() {
	for i := 0; i < lca.bitLen-1; i++ {
		for j := 0; j < lca.n; j++ {
			pre := lca.dp[i][j]
			if pre == -1 {
				lca.dp[i+1][j] = -1
			} else {
				lca.dp[i+1][j] = lca.dp[i][pre]
				lca.dpWeight2[i+1][j] = min(lca.dpWeight2[i][j], lca.dpWeight2[i][lca.dp[i][j]])
			}
		}
	}

	return
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

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func unique(nums []int) []int {
	visited := make(map[int]struct{})
	newNums := []int{}
	for _, v := range nums {
		if _, ok := visited[v]; !ok {
			visited[v] = struct{}{}
			newNums = append(newNums, v)
		}
	}
	return newNums
}

func MinDistToOther(adjList [][][2]int, points []int) (dist []int, nearest []int) {
	n := len(adjList)
	dist = make([]int, n)
	source1, source2 := make([]int, n), make([]int, n)
	for i := 0; i < n; i++ {
		dist[i] = INF
		source1[i], source2[i] = -1, -1
	}

	pq := NewHeap(func(a, b H) bool { return a.dist < b.dist }, nil)
	for _, v := range points {
		pq.Push(H{dist: 0, node: v, source: v})
	}

	for pq.Len() > 0 {
		item := pq.Pop()
		curDist, cur, curSource := item.dist, item.node, item.source
		if curSource == source1[cur] || curSource == source2[cur] {
			continue
		}
		if source1[cur] == -1 {
			source1[cur] = curSource
		} else if source2[cur] == -1 {
			source2[cur] = curSource
		} else {
			continue
		}

		if curSource != cur { // 出发点不为自己时，更新距离
			dist[cur] = min(dist[cur], curDist)
		}
		for _, e := range adjList[cur] {
			next, weight := e[0], e[1]
			nextDist := curDist + weight
			pq.Push(H{nextDist, next, curSource})
		}
	}

	nearest = source2
	for i, v := range points {
		dist[i] = dist[v]
		nearest[i] = nearest[v]
	}
	dist = dist[:len(points)]
	nearest = nearest[:len(points)]
	return
}

type H = struct{ dist, node, source int }

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
