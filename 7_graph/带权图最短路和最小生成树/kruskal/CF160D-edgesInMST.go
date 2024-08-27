// https://codeforces.com/problemset/problem/160/D
// 给一个带权的无向图,保证没有自环和重边.
// 对每条边判断：
// 0:一定在所有最小生成树中
// 1:可能在某个最小生成树中
// 2:一定不在任何最小生成树中

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"sort"
	"strings"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)
	edges := make([][3]int, m)
	for i := 0; i < m; i++ {
		var u, v, w int
		fmt.Fscan(in, &u, &v, &w)
		u, v = u-1, v-1
		edges[i] = [3]int{u, v, w}
	}

	res := EdgesInMST(n, edges)
	for _, v := range res {
		if v == 0 {
			fmt.Fprintln(out, "any")
		} else if v == 1 {

			fmt.Fprintln(out, "at least one")
		} else {
			fmt.Fprintln(out, "none")
		}
	}
}

func EdgesInMST(n int, edges [][3]int) []int {
	m := len(edges)
	edgeGroups := make(map[int][]int) // 按照权值分组
	for i := range edges {
		w := edges[i][2]
		edgeGroups[w] = append(edgeGroups[w], i)
	}
	keys := make([]int, 0, len(edgeGroups))
	for k := range edgeGroups {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	graph := make([][]int, n)
	uf := NewUnionFindArray(n)
	inMST := make([]bool, m)
	for _, k := range keys {
		curEdges := edgeGroups[k]
		for _, eid := range curEdges {
			u, v := edges[eid][0], edges[eid][1]
			if uf.Union(u, v) {
				graph[u] = append(graph[u], v)
				graph[v] = append(graph[v], u)
				inMST[eid] = true
			}
		}
	}

	tree := _NewTree(graph, 0)
	parent, lid, rid := tree.Parent, tree.LID, tree.RID
	bit1, bit2 := NewBitArray(n), NewBitArray(n)
	res := make([]int, m)

	// 不在MST中的边
	for _, k := range keys {
		curEdges := edgeGroups[k]
		for _, eid := range curEdges {
			if inMST[eid] {
				continue
			}
			u, v := edges[eid][0], edges[eid][1]
			lca := tree.LCA(u, v)
			dist := tree.Dist(u, v)
			smallEdge := bit1.Query(lid[u]+1) + bit1.Query(lid[v]+1) - 2*bit1.Query(lid[lca]+1)
			if dist == smallEdge {
				res[eid] = 2
			} else {
				res[eid] = 1
			}
			bit2.Add(lid[u], 1)
			bit2.Add(lid[v], 1)
			bit2.Add(lid[lca], -2)
		}

		// 在MST中的边
		for _, eid := range curEdges {
			if !inMST[eid] {
				continue
			}
			u, v := edges[eid][0], edges[eid][1]
			var p int
			if parent[u] == v {
				p = u
			} else {
				p = v
			}
			x := bit2.QueryRange(lid[p], rid[p])
			if x == 0 {
				res[eid] = 0
			} else {
				res[eid] = 1
			}
			bit1.Add(lid[p], 1)
			bit1.Add(rid[p], -1)
		}

	}

	return res
}

type _Tree struct {
	Tree          [][]int // (next, weight)
	Depth         []int
	Parent        []int
	LID, RID      []int // 欧拉序[in,out)
	IdToNode      []int
	top, heavySon []int
	timer         int
}

func _NewTree(tree [][]int, root int) *_Tree {
	n := len(tree)
	lid := make([]int, n)
	rid := make([]int, n)
	IdToNode := make([]int, n)
	top := make([]int, n)      // 所处轻/重链的顶点（深度最小），轻链的顶点为自身
	depth := make([]int, n)    // 深度
	parent := make([]int, n)   // 父结点
	heavySon := make([]int, n) // 重儿子
	for i := range parent {
		parent[i] = -1
	}

	res := &_Tree{
		Tree:     tree,
		Depth:    depth,
		Parent:   parent,
		LID:      lid,
		RID:      rid,
		IdToNode: IdToNode,
		top:      top,
		heavySon: heavySon,
	}

	res.build(root, -1, 0)
	res.markTop(root, root)
	return res
}

// 返回 root 的欧拉序区间, 左闭右开, 0-indexed.
func (tree *_Tree) Id(root int) (int, int) {
	return tree.LID[root], tree.RID[root]
}

func (tree *_Tree) LCA(u, v int) int {
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

func (tree *_Tree) RootedLCA(u, v int, root int) int {
	return tree.LCA(u, v) ^ tree.LCA(u, root) ^ tree.LCA(v, root)
}

func (tree *_Tree) RootedParent(u int, root int) int {
	return tree.Jump(u, root, 1)
}

func (tree *_Tree) Dist(u, v int) int {
	return tree.Depth[u] + tree.Depth[v] - 2*tree.Depth[tree.LCA(u, v)]
}

// k: 0-based
//
//	如果不存在第k个祖先，返回-1
//	kthAncestor(root,0) == root
func (tree *_Tree) KthAncestor(root, k int) int {
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
func (tree *_Tree) Jump(from, to, step int) int {
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

// 返回沿着`路径顺序`的 [起点,终点] 的 欧拉序 `左闭右闭` 数组.
//
//	!eg:[[2 0] [4 4]] 沿着路径顺序但不一定沿着欧拉序.
func (tree *_Tree) GetPathDecomposition(u, v int, vertex bool) [][2]int {
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
func (tree *_Tree) EnumeratePathDecomposition(u, v int, vertex bool, f func(start, end int)) {
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

func (tree *_Tree) GetPath(u, v int) []int {
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
func (tree *_Tree) SubSize(v, root int) int {
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
func (tree *_Tree) IsInSubtree(child, root int) bool {
	return tree.LID[root] <= tree.LID[child] && tree.LID[child] < tree.RID[root]
}

// 寻找以 start 为 top 的重链 ,heavyPath[-1] 即为重链底端节点.
func (tree *_Tree) GetHeavyPath(start int) []int {
	heavyPath := []int{start}
	cur := start
	for tree.heavySon[cur] != -1 {
		cur = tree.heavySon[cur]
		heavyPath = append(heavyPath, cur)
	}
	return heavyPath
}

// 结点v的重儿子.如果没有重儿子,返回-1.
func (tree *_Tree) GetHeavyChild(v int) int {
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

func (tree *_Tree) ELID(u int) int {
	return 2*tree.LID[u] - tree.Depth[u]
}

func (tree *_Tree) ERID(u int) int {
	return 2*tree.RID[u] - tree.Depth[u] - 1
}

func (tree *_Tree) build(cur, pre, dep int) int {
	subSize, heavySize, heavySon := 1, 0, -1
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

func (tree *_Tree) markTop(cur, top int) {
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

// !Point Add Range Sum, 0-based.
type BITArray struct {
	n    int
	log  int
	data []int
}

func NewBitArray(n int) *BITArray {
	return &BITArray{n: n, log: bits.Len(uint(n)), data: make([]int, n+1)}
}

func NewBitArrayFrom(arr []int) *BITArray {
	res := NewBitArray(len(arr))
	res.Build(arr)
	return res
}

func (b *BITArray) Build(arr []int) {
	if b.n != len(arr) {
		panic("len of arr is not equal to n")
	}
	for i := 1; i <= b.n; i++ {
		b.data[i] = arr[i-1]
	}
	for i := 1; i <= b.n; i++ {
		j := i + (i & -i)
		if j <= b.n {
			b.data[j] += b.data[i]
		}
	}
}

func (b *BITArray) Add(i int, v int) {
	for i++; i <= b.n; i += i & -i {
		b.data[i] += v
	}
}

// [0, r)
func (b *BITArray) Query(r int) int {
	res := 0
	for ; r > 0; r -= r & -r {
		res += b.data[r]
	}
	return res
}

// [l, r).
func (b *BITArray) QueryRange(l, r int) int {
	return b.Query(r) - b.Query(l)
}

// 返回闭区间[0,k]的总和>=x的最小k.要求序列单调增加.
func (b *BITArray) LowerBound(x int) int {
	i := 0
	for k := 1 << b.log; k > 0; k >>= 1 {
		if i+k <= b.n && b.data[i+k] < x {
			x -= b.data[i+k]
			i += k
		}
	}
	return i
}

// 返回闭区间[0,k]的总和>x的最小k.要求序列单调增加.
func (b *BITArray) UpperBound(x int) int {
	i := 0
	for k := 1 << b.log; k > 0; k >>= 1 {
		if i+k <= b.n && b.data[i+k] <= x {
			x -= b.data[i+k]
			i += k
		}
	}
	return i
}

func (b *BITArray) String() string {
	sb := []string{}
	for i := 0; i < b.n; i++ {
		sb = append(sb, fmt.Sprintf("%d", b.QueryRange(i, i+1)))
	}
	return fmt.Sprintf("BitArray: [%v]", strings.Join(sb, ", "))
}

type UnionFindArray struct {
	// 连通分量的个数
	Part int
	n    int
	data []int
}

func NewUnionFindArray(n int) *UnionFindArray {
	data := make([]int, n)
	for i := 0; i < n; i++ {
		data[i] = -1
	}
	return &UnionFindArray{
		Part: n,
		n:    n,
		data: data,
	}
}

// 按秩合并.
func (ufa *UnionFindArray) Union(key1, key2 int) bool {
	root1, root2 := ufa.Find(key1), ufa.Find(key2)
	if root1 == root2 {
		return false
	}
	if ufa.data[root1] > ufa.data[root2] {
		root1 ^= root2
		root2 ^= root1
		root1 ^= root2
	}
	ufa.data[root1] += ufa.data[root2]
	ufa.data[root2] = root1
	ufa.Part--
	return true
}

func (ufa *UnionFindArray) UnionWithCallback(key1, key2 int, cb func(big, small int)) bool {
	root1, root2 := ufa.Find(key1), ufa.Find(key2)
	if root1 == root2 {
		return false
	}
	if ufa.data[root1] > ufa.data[root2] {
		root1 ^= root2
		root2 ^= root1
		root1 ^= root2
	}
	ufa.data[root1] += ufa.data[root2]
	ufa.data[root2] = root1
	ufa.Part--
	if cb != nil {
		cb(root1, root2)
	}
	return true
}

func (ufa *UnionFindArray) Find(key int) int {
	if ufa.data[key] < 0 {
		return key
	}
	ufa.data[key] = ufa.Find(ufa.data[key])
	return ufa.data[key]
}

func (ufa *UnionFindArray) IsConnected(key1, key2 int) bool {
	return ufa.Find(key1) == ufa.Find(key2)
}

func (ufa *UnionFindArray) GetSize(key int) int {
	return -ufa.data[ufa.Find(key)]
}

func (ufa *UnionFindArray) GetGroups() map[int][]int {
	groups := make(map[int][]int)
	for i := 0; i < ufa.n; i++ {
		root := ufa.Find(i)
		groups[root] = append(groups[root], i)
	}
	return groups
}

func (ufa *UnionFindArray) String() string {
	sb := []string{"UnionFindArray:"}
	groups := ufa.GetGroups()
	keys := make([]int, 0, len(groups))
	for k := range groups {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, root := range keys {
		member := groups[root]
		cur := fmt.Sprintf("%d: %v", root, member)
		sb = append(sb, cur)
	}
	sb = append(sb, fmt.Sprintf("Part: %d", ufa.Part))
	return strings.Join(sb, "\n")
}

var _pool = make(map[interface{}]int)

func id(o interface{}) int {
	if v, ok := _pool[o]; ok {
		return v
	}
	v := len(_pool)
	_pool[o] = v
	return v
}

type UnionFindMap struct {
	Part int
	data map[int]int
}

func NewUnionFindMap() *UnionFindMap {
	return &UnionFindMap{
		data: make(map[int]int),
	}
}

func (ufm *UnionFindMap) Union(key1, key2 int) bool {
	root1, root2 := ufm.Find(key1), ufm.Find(key2)
	if root1 == root2 {
		return false
	}
	if ufm.data[root1] > ufm.data[root2] {
		root1 ^= root2
		root2 ^= root1
		root1 ^= root2
	}
	ufm.data[root1] += ufm.data[root2]
	ufm.data[root2] = root1
	ufm.Part--
	return true
}

func (ufm *UnionFindMap) UnionWithCallback(key1, key2 int, cb func(big, small int)) bool {
	root1, root2 := ufm.Find(key1), ufm.Find(key2)
	if root1 == root2 {
		return false
	}
	if ufm.data[root1] > ufm.data[root2] {
		root1 ^= root2
		root2 ^= root1
		root1 ^= root2
	}
	ufm.data[root1] += ufm.data[root2]
	ufm.data[root2] = root1
	ufm.Part--
	if cb != nil {
		cb(root1, root2)
	}
	return true
}

func (ufm *UnionFindMap) Find(key int) int {
	if _, ok := ufm.data[key]; !ok {
		ufm.Add(key)
		return key
	}
	if ufm.data[key] < 0 {
		return key
	}
	ufm.data[key] = ufm.Find(ufm.data[key])
	return ufm.data[key]
}

func (ufm *UnionFindMap) IsConnected(key1, key2 int) bool {
	return ufm.Find(key1) == ufm.Find(key2)
}

func (ufm *UnionFindMap) GetSize(key int) int {
	return -ufm.data[ufm.Find(key)]
}

func (ufm *UnionFindMap) GetGroups() map[int][]int {
	groups := make(map[int][]int)
	for k := range ufm.data {
		root := ufm.Find(k)
		groups[root] = append(groups[root], k)
	}
	return groups
}

func (ufm *UnionFindMap) Has(key int) bool {
	_, ok := ufm.data[key]
	return ok
}

func (ufm *UnionFindMap) Add(key int) bool {
	if _, ok := ufm.data[key]; ok {
		return false
	}
	ufm.data[key] = -1
	ufm.Part++
	return true
}

func (ufm *UnionFindMap) String() string {
	sb := []string{"UnionFindMap:"}
	groups := ufm.GetGroups()
	keys := make([]int, 0, len(groups))
	for k := range groups {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, root := range keys {
		member := groups[root]
		cur := fmt.Sprintf("%d: %v", root, member)
		sb = append(sb, cur)
	}
	sb = append(sb, fmt.Sprintf("Part: %d", ufm.Part))
	return strings.Join(sb, "\n")
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
