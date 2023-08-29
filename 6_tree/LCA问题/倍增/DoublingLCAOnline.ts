// !动态维护k级祖先，可以动态添加叶子节点(在树上做递归/回溯操作时比较有用)

// type DoublingLCAOnline struct {
// 	Depth         []int
// 	DepthWeighted []int
// 	n             int
// 	bitLen        int
// 	dp            [][]int
// 	root          int
// }

// // 不预先给出整棵树,而是动态添加叶子节点,维护树节点的LCA和k级祖先.
// // 初始时只有一个根节点root.
// func NewDoublingLCAOnline(n int, root int) *DoublingLCAOnline {
// 	n += 1 // 防止越界
// 	bit := bits.Len(uint(n))
// 	dp := make([][]int, bit)
// 	for i := range dp {
// 		cur := make([]int, n)
// 		for j := range cur {
// 			cur[j] = -1
// 		}
// 		dp[i] = cur
// 	}
// 	depth := make([]int, n)
// 	depthWeighted := make([]int, n)
// 	return &DoublingLCAOnline{n: n, bitLen: bit, dp: dp, Depth: depth, DepthWeighted: depthWeighted, root: root}
// }

// // 在树中添加一条从parent到child的边.要求parent已经存在于树中.
// func (lca *DoublingLCAOnline) AddDirectedEdge(parent, child int, weight int) {
// 	if parent != lca.root && lca.Depth[parent] == 0 {
// 		panic(fmt.Sprintf("parent %d not exists", parent))
// 	}
// 	lca.Depth[child] = lca.Depth[parent] + 1
// 	lca.DepthWeighted[child] = lca.DepthWeighted[parent] + weight
// 	lca.dp[0][child] = parent
// 	for i := 0; i < lca.bitLen-1; i++ {
// 		pre := lca.dp[i][child]
// 		if pre == -1 {
// 			break
// 		}
// 		lca.dp[i+1][child] = lca.dp[i][pre]
// 	}
// }

// // 查询节点node的第k个祖先(向上跳k步).如果不存在,返回-1.
// func (lca *DoublingLCAOnline) KthAncestor(node, k int) int {
// 	if k > lca.Depth[node] {
// 		return -1
// 	}
// 	bit := 0
// 	for k > 0 {
// 		if k&1 == 1 {
// 			node = lca.dp[bit][node]
// 			if node == -1 {
// 				return -1
// 			}
// 		}
// 		bit++
// 		k >>= 1
// 	}
// 	return node
// }

// // 从 root 开始向上跳到指定深度 toDepth,toDepth<=depth[v],返回跳到的节点.
// func (lca *DoublingLCAOnline) UpToDepth(root, toDepth int) int {
// 	if toDepth >= lca.Depth[root] {
// 		return root
// 	}
// 	for i := lca.bitLen - 1; i >= 0; i-- {
// 		if (lca.Depth[root]-toDepth)&(1<<i) > 0 {
// 			root = lca.dp[i][root]
// 		}
// 	}
// 	return root
// }

// func (lca *DoublingLCAOnline) LCA(root1, root2 int) int {
// 	if lca.Depth[root1] < lca.Depth[root2] {
// 		root1, root2 = root2, root1
// 	}
// 	root1 = lca.UpToDepth(root1, lca.Depth[root2])
// 	if root1 == root2 {
// 		return root1
// 	}
// 	for i := lca.bitLen - 1; i >= 0; i-- {
// 		if lca.dp[i][root1] != lca.dp[i][root2] {
// 			root1 = lca.dp[i][root1]
// 			root2 = lca.dp[i][root2]
// 		}
// 	}
// 	return lca.dp[0][root1]
// }

// // 从start节点跳向target节点,跳过step个节点(0-indexed)
// // 返回跳到的节点,如果不存在这样的节点,返回-1
// func (lca *DoublingLCAOnline) Jump(start, target, step int) int {
// 	lca_ := lca.LCA(start, target)
// 	dep1, dep2, deplca := lca.Depth[start], lca.Depth[target], lca.Depth[lca_]
// 	dist := dep1 + dep2 - 2*deplca
// 	if step > dist {
// 		return -1
// 	}
// 	if step <= dep1-deplca {
// 		return lca.KthAncestor(start, step)
// 	}
// 	return lca.KthAncestor(target, dist-step)
// }

// func (lca *DoublingLCAOnline) Dist(root1, root2 int, weighted bool) int {
// 	if weighted {
// 		return lca.DepthWeighted[root1] + lca.DepthWeighted[root2] - 2*lca.DepthWeighted[lca.LCA(root1, root2)]
// 	}
// 	return lca.Depth[root1] + lca.Depth[root2] - 2*lca.Depth[lca.LCA(root1, root2)]
// }

class DoublingLCAOnline {}

export { DoublingLCAOnline }
