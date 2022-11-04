package heavyLightDecomposition

import "cmnx/src/segmentTree"



// 树链剖分/重链剖分 (HLD, Heavy Light Decomposition）
//  性质：
//    1.树上每个结点都属于且仅属于一条重链
//    2.从根结点到任意结点所经过的重链数为 O(logn)，轻边数为 O(logn)
//    3.如果边(u,v),为轻边,那么Size(v)≤Size(u)/2。
// 树链剖分详解 https://www.cnblogs.com/zwfymqz/p/8094500.html
// 树链剖分详解 https://www.luogu.com.cn/blog/communist/shu-lian-pou-fen-yang-xie
//
// https://github.dev/EndlessCheng/codeforces-go/blob/016834c19c4289ae5999988585474174224f47e2/copypasta/graph_tree.go#L1008
func HeavyLightDecomposition(n, root int, tree [][]int, vals []int) { // vals 为点权
	type node struct {
		// !重儿子，不存在时为 -1
		hSon int
		// !所处轻/重链的顶点（深度最小），轻链的顶点为自身
		top int
		// 深度
		depth int
		// 子树大小
		size int
		// 父节点
		fa int
		// DFS 序（作为线段树中的编号，从 1 开始）
		//  !子树区间为 [dfn, dfn+size-1]
		dfn int
	}
	nodes := make([]node, n)
	dfnToNode := make([]int, n+1) // dfs序在树中所对应的节点  dfnToNode[nodes[v].dfn] == v

	// 寻找重儿子，拆分出重链与轻链
	var build func(cur, pre, dep int) int
	build = func(cur, pre, dep int) int {
		subSize, hSize, hSon := 1, 0, -1
		for _, next := range tree[cur] {
			if next != pre {
				nextSize := build(next, cur, dep+1)
				subSize += nextSize
				if nextSize > hSize {
					hSize, hSon = nextSize, next
				}
			}
		}
		nodes[cur] = node{depth: dep, size: subSize, hSon: hSon, fa: pre}
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
		if node.hSon != -1 {
			// 优先遍历重儿子，保证在同一条重链上的点的 DFS 序是连续的
			markTop(node.hSon, top)
			for _, next := range tree[cur] {
				if next != node.fa && next != node.hSon {
					markTop(next, next)
				}
			}
		}
	}
	markTop(root, root)

	// 按照 DFS 序对应的点权初始化线段树 维护重链信息
	dfnVals := make([]int, n)
	for i, v := range vals {
		dfnVals[nodes[i].dfn-1] = v
	}
	segTree := segmentTree.New(dfnVals)

	// 利用重链上跳到top来加速路径处理
	doPath := func(root1, root2 int, cb func(left, right int)) {
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
			root1 = top1.fa
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
		for node := nodes[start]; node.hSon != -1; node = nodes[node.hSon] {
			heavyPath = append(heavyPath, node.hSon)
		}
		return heavyPath
	}

	// 将节点 root1 和节点 root2 之间路径上的所有节点（包括这两个节点）的权值增加 add。
	updatePath := func(root1, root2 int, add int) {
		doPath(root1, root2, func(left, right int) { segTree.update(left, right, add) })
	}

	// 询问节点 root1 和节点 root2 之间路径上的所有节点（包括这两个节点）的权值和。
	// % mod
	queryPath := func(root1, root2 int) (sum int) {
		doPath(root1, root2, func(left, right int) { sum += segTree.query(left, right) })
		return
	}

	// 将以节点 root 为根的子树上的所有节点的权值增加 add。
	updateSubtree := func(root int, add int) {
		node := nodes[root]
		segTree.update(node.dfn, node.dfn+node.size-1, add)
	}

	// 询问以节点 root 为根的子树上的所有节点的权值和。
	querySubtree := func(root int) (sum int) {
		node := nodes[root]
		return segTree.query(1, node.dfn, node.dfn+node.size-1)
	}

	_ = []interface{}{
		getHeavyPath,
		updatePath, queryPath, updateSubtree, querySubtree,
	}
}
