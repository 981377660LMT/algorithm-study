package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"sort"
)

const INF int = 1e18

// G - Sum of Tree Distance (虚树+换根dp)
// https://atcoder.jp/contests/abc359/tasks/abc359_g
// 给定一棵树,每个点有一个颜色。
// 对每一种颜色相同的点，求出每个点到其他所有颗色相同的点的距离和。
// !虚树上求点对距离.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	edges := make([][2]int, n-1)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		edges[i] = [2]int{u - 1, v - 1}
	}
	colors := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &colors[i])
	}

	groupByColor := make(map[int][]int)
	for i, c := range colors {
		groupByColor[c] = append(groupByColor[c], i)
	}
	tree := NewTree(n)
	for _, e := range edges {
		tree.AddEdge(e[0], e[1], 1)
	}
	tree.Build(0)

	res := 0
	isCritical := make([]bool, n)
	for _, criticals := range groupByColor {
		for _, v := range criticals {
			isCritical[v] = true
		}

		rawId, newTree := CompressTree(tree, criticals, false)
		adjList := newTree.Tree
		starts := make([]int, 0, len(criticals))
		for i := 0; i < len(rawId); i++ {
			if isCritical[rawId[i]] {
				starts = append(starts, i)
			}
		}

		res += GetDistSumToSpecials(newTree, adjList, starts)

		for _, v := range criticals {
			isCritical[v] = false
		}
	}

	fmt.Fprintln(out, res/2)
}

// 换根dp，求specials中所有点到其他所有有specials的点的距离和
func GetDistSumToSpecials(tree *Tree, adjList [][][2]int, specials []int) int {
	isSpecial := make([]bool, len(adjList))
	for _, v := range specials {
		isSpecial[v] = true
	}
	type E = [2]int
	R := NewRerooting[E](len(adjList))

	e := func(root int) E { return [2]int{0, 0} }
	op := func(childRes1, childRes2 E) E {
		dist1, cnt1 := childRes1[0], childRes1[1]
		dist2, cnt2 := childRes2[0], childRes2[1]
		return [2]int{dist1 + dist2, cnt1 + cnt2}
	}
	composition := func(fromRes E, parent, cur int, direction uint8) E {
		preDist, preCnt := fromRes[0], fromRes[1]
		from_ := cur
		if direction == 1 {
			from_ = parent
		}
		curCnt := 0
		if isSpecial[from_] {
			curCnt = 1
		}
		return [2]int{preDist + tree.Dist(parent, cur, true)*(preCnt+curCnt), preCnt + curCnt}
	}

	for u, e := range adjList {
		for _, v := range e {
			if u < v[0] {
				R.AddEdge(u, v[0])
			}
		}
	}

	dp := R.ReRooting(e, op, composition)

	res := 0
	for i := 0; i < len(dp); i++ {
		if isSpecial[i] {
			res += dp[i][0]
		}
	}
	return res

}

type Rerooting[E any] struct {
	Tree [][]int
	n    int
}

func NewRerooting[E any](n int) *Rerooting[E] {
	return &Rerooting[E]{Tree: make([][]int, n), n: n}
}

func NewRerootingFromTree[E any](tree [][]int) *Rerooting[E] {
	return &Rerooting[E]{Tree: tree, n: len(tree)}
}

func (r *Rerooting[E]) AddEdge(u, v int) {
	r.Tree[u] = append(r.Tree[u], v)
	r.Tree[v] = append(r.Tree[v], u)
}

func (r *Rerooting[E]) ReRooting(e func(root int) E, op func(child1, child2 E) E, composition func(fromRes E, parent, cur int, direction uint8) E) []E {
	parents := make([]int, r.n)
	for i := range parents {
		parents[i] = -1
	}
	order := []int{0}
	stack := []int{0}
	for len(stack) > 0 {
		cur := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		for _, next := range r.Tree[cur] {
			if next != parents[cur] {
				parents[next] = cur
				order = append(order, next)
				stack = append(stack, next)
			}
		}
	}

	dp1, dp2 := make([]E, r.n), make([]E, r.n)
	for i := range dp1 {
		dp1[i] = e(i)
		dp2[i] = e(i)
	}
	for i := r.n - 1; i >= 0; i-- {
		cur := order[i]
		res := e(cur)
		for _, next := range r.Tree[cur] {
			if next != parents[cur] {
				dp2[next] = res
				res = op(res, composition(dp1[next], cur, next, 0))
			}
		}

		res = e(cur)
		for j := len(r.Tree[cur]) - 1; j >= 0; j-- {
			next := r.Tree[cur][j]
			if next != parents[cur] {
				dp2[next] = op(res, dp2[next])
				res = op(res, composition(dp1[next], cur, next, 0))
			}
		}

		dp1[cur] = res
	}

	for i := 1; i < r.n; i++ {
		newRoot := order[i]
		parent := parents[newRoot]
		dp2[newRoot] = composition(op(dp2[newRoot], dp2[parent]), parent, newRoot, 1)
		dp1[newRoot] = op(dp1[newRoot], dp2[newRoot])
	}

	return dp1
}

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
	Tree                 [][][2]int
	Depth, DepthWeighted []int
	Parent               []int
	LID, RID             []int
	IdToNode             []int
	top, heavySon        []int
	timer                int
}

func NewTree(n int) *Tree {
	tree := make([][][2]int, n)
	lid := make([]int, n)
	rid := make([]int, n)
	IdToNode := make([]int, n)
	top := make([]int, n)
	depth := make([]int, n)
	depthWeighted := make([]int, n)
	parent := make([]int, n)
	heavySon := make([]int, n)
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

func (tree *Tree) AddEdge(u, v, w int) {
	tree.Tree[u] = append(tree.Tree[u], [2]int{v, w})
	tree.Tree[v] = append(tree.Tree[v], [2]int{u, w})
}

func (tree *Tree) AddDirectedEdge(u, v, w int) {
	tree.Tree[u] = append(tree.Tree[u], [2]int{v, w})
}

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

func (tree *Tree) Id(root int) (int, int) {
	return tree.LID[root], tree.RID[root]
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

func (tree *Tree) RootedLCA(u, v int, w int) int {
	return tree.LCA(u, v) ^ tree.LCA(u, w) ^ tree.LCA(v, w)
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

func (tree *Tree) IsInSubtree(child, root int) bool {
	return tree.LID[root] <= tree.LID[child] && tree.LID[child] < tree.RID[root]
}

func (tree *Tree) GetHeavyPath(start int) []int {
	heavyPath := []int{start}
	cur := start
	for tree.heavySon[cur] != -1 {
		cur = tree.heavySon[cur]
		heavyPath = append(heavyPath, cur)
	}
	return heavyPath
}

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
	dp            [][]int32
	dpWeight2     [][]int
}

func NewLCADoubling(tree [][][2]int, roots []int) *LCADoubling {
	n := len(tree)
	depth := make([]int32, n)
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

func (lca *LCADoubling) QueryDist(root1, root2 int, weighted bool) int {
	if weighted {
		return lca.DepthWeighted[root1] + lca.DepthWeighted[root2] - 2*lca.DepthWeighted[lca.QueryLCA(root1, root2)]
	}
	return int(lca.Depth[root1] + lca.Depth[root2] - 2*lca.Depth[lca.QueryLCA(root1, root2)])
}

func (lca *LCADoubling) QueryMinWeight(root1, root2 int, isEdge bool) int {
	res := INF
	if lca.Depth[root1] < lca.Depth[root2] {
		root1, root2 = root2, root1
	}
	toDepth := lca.Depth[root2]
	root132, root232 := int32(root1), int32(root2)
	for i := lca.bitLen - 1; i >= 0; i-- {
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
