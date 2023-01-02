// VertexAddPathSum
// https://judge.yosupo.jp/problem/vertex_add_path_sum
// 单点加/路径和查询
// 0 vertex add => 顶点加
// 1 root1 root2 => 路径和查询
package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)

	values := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &values[i])
	}

	tree := make([][]int, n)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		tree[u] = append(tree[u], v)
		tree[v] = append(tree[v], u)
	}

	hld := HeavyLightDecomposition(n, 0, tree, values)
	for i := 0; i < q; i++ {
		var op, vertex, add, root1, root2 int
		fmt.Fscan(in, &op)
		if op == 0 {
			fmt.Fscan(in, &vertex, &add)
			hld.UpdatePath(vertex, vertex, add)
		} else {
			fmt.Fscan(in, &root1, &root2)
			fmt.Fprintln(out, hld.QueryPath(root1, root2))
		}
	}
}

type HLD struct {
	UpdatePath    func(root1 int, root2 int, add int)
	QueryPath     func(root1 int, root2 int) int
	UpdateSubtree func(root int, add int)
	QuerySubtree  func(root int) int
	GetHeavyPath  func(start int) []int
}

// HeavyLightDecomposition returns a HLD struct.
//  n: the number of nodes in the tree.
//  root: the root of the tree.
//  adjList: the adjacency list of the tree.
//  vals: the values of the nodes.
// https://github.dev/EndlessCheng/codeforces-go/blob/016834c19c4289ae5999988585474174224f47e2/copypasta/graph_tree.go#L1008
func HeavyLightDecomposition(n, root int, tree [][]int, vals []int) *HLD {
	type node struct {
		// !重儿子，不存在时为 -1
		heavySon int
		// !所处轻/重链的顶点（深度最小），轻链的顶点为自身
		top int
		// 深度
		depth int
		// 子树大小
		subSize int
		// 父节点
		parent int
		// DFS 序（作为线段树中的编号，从 1 开始）
		//  !子树区间为 [dfn, dfn+size-1]
		dfn int
	}
	nodes := make([]node, n)
	dfnToNode := make([]int, n+1) // dfs序在树中所对应的节点  dfnToNode[nodes[v].dfn] == v

	// 寻找重儿子，拆分出重链与轻链
	var build func(cur, parent, depth int) int
	build = func(cur, parent, depth int) int {
		subSize, heavySize, heavySon := 1, 0, -1
		for _, next := range tree[cur] {
			if next != parent {
				nextSize := build(next, cur, depth+1)
				subSize += nextSize
				if nextSize > heavySize {
					heavySize, heavySon = nextSize, next
				}
			}
		}
		nodes[cur] = node{depth: depth, subSize: subSize, heavySon: heavySon, parent: parent}
		return subSize
	}
	build(root, -1, 0)

	// 对这些链进行维护，就要确保每个链上的节点都是连续的
	// 注意在进行重新编号的时候先访问重链，这样可以保证重链内的节点编号连续
	dfn := 0
	var markTop func(cur, top int)
	markTop = func(cur, top int) {
		node := &nodes[cur]
		node.top = top
		dfn++
		node.dfn = dfn
		dfnToNode[dfn] = cur
		if node.heavySon != -1 {
			// 优先遍历重儿子，保证在同一条重链上的点的 DFS 序是连续的
			markTop(node.heavySon, top)
			for _, next := range tree[cur] {
				if next != node.parent && next != node.heavySon {
					markTop(next, next)
				}
			}
		}
	}
	markTop(root, root)

	// 按照 DFS 序对应的点权初始化线段树 维护重链信息
	bit := NewSliceBIT(n + 10)
	for i, v := range vals {
		dfsId := nodes[i].dfn
		bit.Add(dfsId, dfsId, v)
	}

	// 利用重链上跳到top来加速路径处理
	walk := func(root1, root2 int, cb func(left, right int)) {
		node1, node2 := nodes[root1], nodes[root2]
		for ; node1.top != node2.top; node1, node2 = nodes[root1], nodes[root2] {
			top1, top2 := nodes[node1.top], nodes[node2.top]
			// root1 所处的重链顶点必须比 root2 的深
			if top1.depth < top2.depth {
				root1, root2 = root2, root1
				node1, node2 = node2, node1
				top1, top2 = top2, top1
			}
			cb(top1.dfn, node1.dfn)
			// TODO: 边权下，处理轻边的情况
			root1 = top1.parent
		}
		if node1.depth > node2.depth {
			root1, root2 = root2, root1
			node1, node2 = node2, node1
		}
		cb(node1.dfn, node2.dfn)
		// TODO: 边权下，处理轻边的情况
	}

	// 寻找以 start 为 top 的重链 ,heavyPath[-1] 即为重链末端节点
	getHeavyPath := func(start int) []int {
		heavyPath := []int{start}
		for node := nodes[start]; node.heavySon != -1; node = nodes[node.heavySon] {
			heavyPath = append(heavyPath, node.heavySon)
		}
		return heavyPath
	}

	// 将节点 root1 和节点 root2 之间路径上的所有节点（包括这两个节点）的权值增加 add。
	updatePath := func(root1, root2 int, add int) {
		walk(root1, root2, func(left, right int) { bit.Add(left, right, add) })
	}

	// 询问节点 root1 和节点 root2 之间路径上的所有节点（包括这两个节点）的权值和。
	// % mod
	queryPath := func(root1, root2 int) (sum int) {
		walk(root1, root2, func(left, right int) { sum += bit.Query(left, right) })
		return
	}

	// 将以节点 root 为根的子树上的所有节点的权值增加 add。
	updateSubtree := func(root int, add int) {
		node := nodes[root]
		bit.Add(node.dfn, node.dfn+node.subSize-1, add)
	}

	// 询问以节点 root 为根的子树上的所有节点的权值和。
	querySubtree := func(root int) (sum int) {
		node := nodes[root]
		return bit.Query(node.dfn, node.dfn+node.subSize-1)
	}

	return &HLD{
		QueryPath:     queryPath,
		UpdatePath:    updatePath,
		QuerySubtree:  querySubtree,
		UpdateSubtree: updateSubtree,
		GetHeavyPath:  getHeavyPath,
	}
}

type SliceBIT struct {
	n     int
	tree1 []int
	tree2 []int
}

func NewSliceBIT(n int) *SliceBIT {
	return &SliceBIT{
		n:     n,
		tree1: make([]int, n+1),
		tree2: make([]int, n+1),
	}
}

func (bit *SliceBIT) Add(left, right, delta int) {
	bit.add(left, delta)
	bit.add(right+1, -delta)
}

func (bit *SliceBIT) Query(left, right int) int {
	return bit.query(right) - bit.query(left-1)
}

func (bit *SliceBIT) add(index, delta int) {
	if index <= 0 {
		errorInfo := fmt.Sprintf("index must be greater than 0, but got %d", index)
		panic(errorInfo)
	}

	rawIndex := index
	for index <= bit.n {
		bit.tree1[index] += delta
		bit.tree2[index] += (rawIndex - 1) * delta
		index += index & -index
	}
}

func (bit *SliceBIT) query(index int) int {
	if index > bit.n {
		index = bit.n
	}

	rawIndex := index
	res := 0
	for index > 0 {
		res += rawIndex*bit.tree1[index] - bit.tree2[index]
		index -= index & -index
	}
	return res
}

func (bit *SliceBIT) String() string {
	nums := make([]int, bit.n+1)
	for i := 0; i < bit.n; i++ {
		nums[i+1] = bit.Query(i+1, i+1)
	}
	return fmt.Sprint(nums)
}
