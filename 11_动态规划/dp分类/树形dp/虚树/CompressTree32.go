package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type int32 = int

const INF32 int32 = 1e18

func main() {
	CF613D()
}

// Kingdom and its Cities
// 给定一棵树，每次询问给定 k个特殊点，找出尽量少的非特殊点使得删去这些点后特殊点两两不连通。∑k≤n.
// 如果无法使得特殊点两两不连通，输出-1.
// https://www.luogu.com.cn/problem/CF613D
func CF613D() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	unique := func(nums []int32) []int32 {
		visited := make(map[int32]struct{})
		newNums := []int32{}
		for _, v := range nums {
			if _, ok := visited[v]; !ok {
				visited[v] = struct{}{}
				newNums = append(newNums, v)
			}
		}
		return newNums
	}

	const INF32 int32 = 1e9

	var n int32
	fmt.Fscan(in, &n)
	tree := NewTree32(n)
	for i := int32(0); i < n-1; i++ {
		var u, v int32
		fmt.Fscan(in, &u, &v)
		tree.AddEdge(u-1, v-1)
	}
	tree.Build(0)

	// !dp[i] 表示子树中保留i个关键点时的最小删除点数
	// ①:如果一个点被标记了,那么就要把他所有子树里有标记点的儿子都去掉
	// ②:如果一个点没有被标记,但是这个点有两颗以上的子树里有标记点，那么这个点就要去掉,然后这棵子树就没有可标记点了
	// ③:如果一个点子树里只有一个/没有标记点，那么就标记点的贡献挪到这个点上面来
	solve := func(adjList [][]int32, inCriticals []bool) int32 {
		var dfs func(cur, pre int32) [2]int32 // (zero, one)
		dfs = func(cur, pre int32) [2]int32 {
			removeCost := int32(1)
			dp := [2]int32{INF32, INF32}
			if inCriticals[cur] {
				removeCost = INF32 // 无法删除
				dp[1] = 0
			} else {
				dp[0] = 0
			}

			for _, next := range adjList[cur] {
				if next == pre {
					continue
				}
				subDp := dfs(next, cur)
				ndp := [2]int32{INF32, INF32}
				for a := 0; a < 2; a++ {
					for b := 0; b < 2; b++ {
						if a == 1 && b == 1 { // !不能>=2个关键点
							continue
						}
						ndp[a+b] = min32(ndp[a+b], dp[a]+subDp[b])
					}
				}
				dp = ndp
				removeCost += min32(subDp[0], subDp[1])
			}

			dp[0] = min32(dp[0], removeCost)
			return dp
		}

		dp := dfs(0, -1)
		res := min32(dp[0], dp[1])
		if res >= INF32 {
			res = -1
		}
		return res
	}

	visitedTime := make([]int32, n)
	for i := range visitedTime {
		visitedTime[i] = -1
	}
	var q int32
	fmt.Fscan(in, &q)
	for qi := int32(0); qi < q; qi++ {
		var k int32
		fmt.Fscan(in, &k)
		criticals := make([]int32, k)
		for j := int32(0); j < k; j++ {
			var v int32
			fmt.Fscan(in, &v)
			v--
			criticals[j] = v
			visitedTime[v] = qi
		}

		nodes := append(criticals[:0:0], criticals...)
		for _, v := range criticals {
			if v != 0 {
				nodes = append(nodes, tree.Parent[v]) // !父节点加进来
			}
		}
		nodes = unique(nodes)

		rawId, newTree := CompressTree32(tree, nodes, true)
		m := len(rawId)
		inCriticals := make([]bool, m) // !压缩后的树中的节点是否在points中
		for i := 0; i < m; i++ {
			inCriticals[i] = visitedTime[rawId[i]] == qi
		}
		fmt.Println(solve(newTree.Tree, inCriticals))
	}
}

// 返回树压缩后保留的节点编号和新的树.
// !新的树保留了原树的边权.
func CompressTree32(rawTree *Tree32, nodes []int32, directed bool) (rawId []int32, newTree *Tree32) {
	rawId = append(nodes[:0:0], nodes...)
	sort.Slice(rawId, func(i, j int) bool { return rawTree.LID[rawId[i]] < rawTree.LID[rawId[j]] })
	n := int32(len(rawId))
	for i := int32(0); i < n; i++ {
		j := i + 1
		if j == n {
			j = 0
		}
		rawId = append(rawId, rawTree.LCA(rawId[i], rawId[j]))
	}
	// remainNodes = append(remainNodes, rawTree.IdToNode[0])
	sort.Slice(rawId, func(i, j int) bool { return rawTree.LID[rawId[i]] < rawTree.LID[rawId[j]] })

	unique := func(a []int32) []int32 {
		visited := make(map[int32]struct{})
		newNums := []int32{}
		for _, v := range a {
			if _, ok := visited[v]; !ok {
				visited[v] = struct{}{}
				newNums = append(newNums, v)
			}
		}
		return newNums
	}

	rawId = unique(rawId)
	n = int32(len(rawId))
	newTree = NewTree32(n)

	stack := []int32{0}
	for i := int32(1); i < n; i++ {
		for {
			p := rawId[stack[len(stack)-1]]
			v := rawId[i]
			if rawTree.IsInSubtree(v, p) {
				break
			}
			stack = stack[:len(stack)-1]
		}
		newTree.AddDirectedEdge(stack[len(stack)-1], i)
		if !directed {
			newTree.AddDirectedEdge(i, stack[len(stack)-1])
		}
		stack = append(stack, i)
	}
	newTree.Build(0)
	return
}

type Tree32 struct {
	Tree          [][]int32 // (next)
	Depth         []int32
	Parent        []int32
	LID, RID      []int32 // 欧拉序[in,out)
	IdToNode      []int32
	top, heavySon []int32
	timer         int32
}

func NewTree32From(tree [][]int32) *Tree32 {
	n := int32(len(tree))
	lid := make([]int32, n)
	rid := make([]int32, n)
	IdToNode := make([]int32, n)
	top := make([]int32, n)      // 所处轻/重链的顶点（深度最小），轻链的顶点为自身
	depth := make([]int32, n)    // 深度
	parent := make([]int32, n)   // 父结点
	heavySon := make([]int32, n) // 重儿子
	for i := range parent {
		parent[i] = -1
	}

	return &Tree32{
		Tree:     tree,
		Depth:    depth,
		Parent:   parent,
		LID:      lid,
		RID:      rid,
		IdToNode: IdToNode,
		top:      top,
		heavySon: heavySon,
	}
}

func NewTree32(n int32) *Tree32 {
	return NewTree32From(make([][]int32, n))
}

// 添加无向边 u-v.
func (tree *Tree32) AddEdge(u, v int32) {
	tree.Tree[u] = append(tree.Tree[u], v)
	tree.Tree[v] = append(tree.Tree[v], u)
}

// 添加有向边 u->v.
func (tree *Tree32) AddDirectedEdge(u, v int32) {
	tree.Tree[u] = append(tree.Tree[u], v)
}

// root:0-based
//
//	当root设为-1时，会从0开始遍历未访问过的连通分量
func (tree *Tree32) Build(root int32) {
	if root != -1 {
		tree.build(root, -1, 0)
		tree.markTop(root, root)
	} else {
		for i := int32(0); i < int32(len(tree.Tree)); i++ {
			if tree.Parent[i] == -1 {
				tree.build(i, -1, 0)
				tree.markTop(i, i)
			}
		}
	}
}

// 返回 root 的欧拉序区间, 左闭右开, 0-indexed.
func (tree *Tree32) Id(root int32) (int32, int32) {
	return tree.LID[root], tree.RID[root]
}

func (tree *Tree32) LCA(u, v int32) int32 {
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

func (tree *Tree32) RootedLCA(u, v int32, root int32) int32 {
	return tree.LCA(u, v) ^ tree.LCA(u, root) ^ tree.LCA(v, root)
}

func (tree *Tree32) RootedParent(u int32, root int32) int32 {
	return tree.Jump(u, root, 1)
}

func (tree *Tree32) Dist(u, v int32) int32 {
	return tree.Depth[u] + tree.Depth[v] - 2*tree.Depth[tree.LCA(u, v)]
}

// k: 0-based
//
//	如果不存在第k个祖先，返回-1
//	kthAncestor(root,0) == root
func (tree *Tree32) KthAncestor(root, k int32) int32 {
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
func (tree *Tree32) Jump(from, to, step int32) int32 {
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

func (tree *Tree32) CollectChild(root int32) []int32 {
	res := []int32{}
	for _, next := range tree.Tree[root] {
		if next != tree.Parent[root] {
			res = append(res, next)
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
		if tree.top[u] == tree.top[v] {
			break
		}
		if tree.LID[u] < tree.LID[v] {
			down = append(down, [2]int32{tree.LID[tree.top[v]], tree.LID[v]})
			v = tree.Parent[tree.top[v]]
		} else {
			up = append(up, [2]int32{tree.LID[u], tree.LID[tree.top[u]]})
			u = tree.Parent[tree.top[u]]
		}
	}
	edgeInt := int32(1)
	if vertex {
		edgeInt = 0
	}
	if tree.LID[u] < tree.LID[v] {
		down = append(down, [2]int32{tree.LID[u] + edgeInt, tree.LID[v]})
	} else if tree.LID[v]+edgeInt <= tree.LID[u] {
		up = append(up, [2]int32{tree.LID[u], tree.LID[v] + edgeInt})
	}
	for i := 0; i < len(down)/2; i++ {
		down[i], down[len(down)-1-i] = down[len(down)-1-i], down[i]
	}
	return append(up, down...)
}

// 遍历路径上的 `[起点,终点)` 欧拉序 `左闭右开` 区间.
func (tree *Tree32) EnumeratePathDecomposition(u, v int32, vertex bool, f func(start, end int32)) {
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

	edgeInt := int32(1)
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

func (tree *Tree32) GetPath(u, v int32) []int32 {
	res := []int32{}
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
func (tree *Tree32) SubSize(v, root int32) int32 {
	if root == -1 {
		return tree.RID[v] - tree.LID[v]
	}
	if v == root {
		return int32(len(tree.Tree))
	}
	x := tree.Jump(v, root, 1)
	if tree.IsInSubtree(v, x) {
		return tree.RID[v] - tree.LID[v]
	}
	return int32(len(tree.Tree)) - tree.RID[x] + tree.LID[x]
}

// child 是否在 root 的子树中 (child和root不能相等)
func (tree *Tree32) IsInSubtree(child, root int32) bool {
	return tree.LID[root] <= tree.LID[child] && tree.LID[child] < tree.RID[root]
}

// 寻找以 start 为 top 的重链 ,heavyPath[-1] 即为重链底端节点.
func (tree *Tree32) GetHeavyPath(start int32) []int32 {
	heavyPath := []int32{start}
	cur := start
	for tree.heavySon[cur] != -1 {
		cur = tree.heavySon[cur]
		heavyPath = append(heavyPath, cur)
	}
	return heavyPath
}

// 结点v的重儿子.如果没有重儿子,返回-1.
func (tree *Tree32) GetHeavyChild(v int32) int32 {
	k := tree.LID[v] + 1
	if k == int32(len(tree.Tree)) {
		return -1
	}
	w := tree.IdToNode[k]
	if tree.Parent[w] == v {
		return w
	}
	return -1
}

func (tree *Tree32) ELID(u int32) int32 {
	return 2*tree.LID[u] - tree.Depth[u]
}

func (tree *Tree32) ERID(u int32) int32 {
	return 2*tree.RID[u] - tree.Depth[u] - 1
}

func (tree *Tree32) build(cur, pre, dep int32) int32 {
	subSize, heavySize, heavySon := int32(1), int32(0), int32(-1)
	for _, next := range tree.Tree[cur] {
		if next != pre {
			nextSize := tree.build(next, cur, dep+1)
			subSize += nextSize
			if nextSize > heavySize {
				heavySize, heavySon = nextSize, next
			}
		}
	}
	tree.Depth[cur] = dep
	tree.heavySon[cur] = heavySon
	tree.Parent[cur] = pre
	return subSize
}

func (tree *Tree32) markTop(cur, top int32) {
	tree.top[cur] = top
	tree.LID[cur] = tree.timer
	tree.IdToNode[tree.timer] = cur
	tree.timer++
	heavySon := tree.heavySon[cur]
	if heavySon != -1 {
		tree.markTop(heavySon, top)
		for _, next := range tree.Tree[cur] {
			if next != heavySon && next != tree.Parent[cur] {
				tree.markTop(next, next)
			}
		}
	}
	tree.RID[cur] = tree.timer
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

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}
