// !动态维护k级祖先，可以动态添加叶子节点(在树上做递归/回溯操作时比较有用)

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {
	// tree := NewDoublingLCAOnline(10, 0)
	// tree.AddDirectedEdge(0, 1, 0)
	// tree.AddDirectedEdge(1, 2, 0)
	R()
}

func R() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var q int
	fmt.Fscan(in, &q)
	operations := make([][2]int, q)
	for i := 0; i < q; i++ {
		var op string
		fmt.Fscan(in, &op)
		switch op {
		case "+":
			var x int
			fmt.Fscan(in, &x)
			operations[i] = [2]int{1, x}
		case "-":
			var k int
			fmt.Fscan(in, &k)
			operations[i] = [2]int{2, k}
		case "!":
			operations[i] = [2]int{3, 0}
		case "?":
			operations[i] = [2]int{4, 0}
		}
	}

	res := Rollbacks(operations)
	for _, v := range res {
		fmt.Fprintln(out, v)
	}
}

// https://www.luogu.com.cn/problem/CF1858E1
// 给定一个初始时为空的数组nums, 需要实现下面四种类型的操作：
// [1, x]: 将x添加到nums尾部
// [2, k]: 将尾部的k个数删除.保证存在k个数.
// [3, 0]: 撤销上一次操作1或2操作
// [4, 0]: 查询当前nums中有多少个不同的数
//
// 1<=q<=1e6,询问次数不超过1e5
//
// !1. 因为要支持撤销，所以需要`保存版本`或者action, 例如immer.js两种方式都会提供.
// 这里保存版本(PersistentStack)比较合适。不同版本构成了一棵树.
// !2. 离线方法是建立一棵版本树，+操作加边，-操作通过倍增上跳到对应节点，
// 操作回退到上个节点，?操作记录当前节点需要记录答案，最后dfs整棵树求解.
// !3. 这个倍增也非常巧妙，不预处理而是是动态的
// !4. 树节点个数最多为n个，因为只有1操作会增加节点.
func Rollbacks(operations [][2]int) []int {
	q := len(operations)
	tree := make([][]int, q+1) // 0是初始空状态的虚拟根节点
	nodeValue := make([]int, q+1)
	nodeQuery := make([][]int, q+1)
	queryIndex := 0
	history := []int{0} // 保存节点位置的stack
	nodeMex := 0
	curNode := 0
	lca := NewDoublingLCAOnline(q + 1)
	for _, op := range operations {
		kind, x := op[0], op[1]
		if kind == 1 {
			nodeMex++
			nodeValue[nodeMex] = x
			tree[curNode] = append(tree[curNode], nodeMex)
			lca.AddDirectedEdge(curNode, nodeMex, 0)
			curNode = nodeMex
			history = append(history, curNode)
		} else if kind == 2 {
			curNode = lca.KthAncestor(curNode, x)
			history = append(history, curNode)
		} else if kind == 3 {
			if len(history) <= 1 {
				continue
			}
			history = history[:len(history)-1]
			curNode = history[len(history)-1]
		} else {
			nodeQuery[curNode] = append(nodeQuery[curNode], queryIndex)
			queryIndex++
		}
	}

	res := make([]int, queryIndex)
	counter := make(map[int]int)
	unique := 0
	var dfs func(int)
	dfs = func(node int) {
		curValue := nodeValue[node]
		if node != 0 {
			counter[curValue]++
			if counter[curValue] == 1 {
				unique++
			}
		}
		for _, q := range nodeQuery[node] {
			res[q] = unique
		}
		for _, next := range tree[node] {
			dfs(next)
		}
		if node != 0 {
			counter[curValue]--
			if counter[curValue] == 0 {
				unique--
			}
		}
	}
	dfs(0)
	return res
}

type DoublingLCAOnline struct {
	Depth         []int32
	DepthWeighted []int
	n             int
	bitLen        int
	dp            [][]int32
}

// 不预先给出整棵树,而是动态添加叶子节点,维护树节点的LCA和k级祖先.
func NewDoublingLCAOnline(n int) *DoublingLCAOnline {
	n += 1 // 防止越界
	bit := bits.Len(uint(n))
	dp := make([][]int32, bit)
	for i := range dp {
		cur := make([]int32, n)
		for j := range cur {
			cur[j] = -1
		}
		dp[i] = cur
	}
	depth := make([]int32, n)
	depthWeighted := make([]int, n)
	return &DoublingLCAOnline{n: n, bitLen: bit, dp: dp, Depth: depth, DepthWeighted: depthWeighted}
}

// 在树中添加一条从parent到child的边.要求parent已经存在于树中(或者为根节点),且child不存在于树中.
func (lca *DoublingLCAOnline) AddDirectedEdge(parent, child int, weight int) {
	lca.Depth[child] = lca.Depth[parent] + 1
	lca.DepthWeighted[child] = lca.DepthWeighted[parent] + weight
	lca.dp[0][child] = int32(parent)
	for i := 0; i < lca.bitLen-1; i++ {
		pre := lca.dp[i][child]
		if pre == -1 {
			break
		}
		lca.dp[i+1][child] = lca.dp[i][pre]
	}
}

// 查询节点node的第k个祖先(向上跳k步).如果不存在,返回-1.
func (lca *DoublingLCAOnline) KthAncestor(node, k int) int {
	node32 := int32(node)
	if k > int(lca.Depth[node32]) {
		return -1
	}
	bit := 0
	for k > 0 {
		if k&1 == 1 {
			node32 = lca.dp[bit][node32]
			if node32 == -1 {
				return -1
			}
		}
		bit++
		k >>= 1
	}
	return int(node32)
}

// 从 root 开始向上跳到指定深度 toDepth,toDepth<=depth[v],返回跳到的节点.
func (lca *DoublingLCAOnline) UpToDepth(root, toDepth int) int {
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

func (lca *DoublingLCAOnline) LCA(root1, root2 int) int {
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

// 从start节点跳向target节点,跳过step个节点(0-indexed)
// 返回跳到的节点,如果不存在这样的节点,返回-1
func (lca *DoublingLCAOnline) Jump(start, target, step int) int {
	lca_ := lca.LCA(start, target)
	dep1, dep2, deplca := lca.Depth[start], lca.Depth[target], lca.Depth[lca_]
	dist := int(dep1 + dep2 - 2*deplca)
	if step > dist {
		return -1
	}
	if step <= int(dep1-deplca) {
		return lca.KthAncestor(start, step)
	}
	return lca.KthAncestor(target, dist-step)
}

func (lca *DoublingLCAOnline) Dist(root1, root2 int, weighted bool) int {
	if weighted {
		return lca.DepthWeighted[root1] + lca.DepthWeighted[root2] - 2*lca.DepthWeighted[lca.LCA(root1, root2)]
	}
	return int(lca.Depth[root1] + lca.Depth[root2] - 2*lca.Depth[lca.LCA(root1, root2)])
}
