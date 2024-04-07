// DividePathOnTreeBinaryLift/DoublingLca
// 倍增拆树上路径`path(from,to)`：倍增拆点将树上的一段路径拆成logn个点
// TODO: 与`CompressedLCA`功能保持一致，并增加拆路径的功能.

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {
	jump()
}

type DivideIntervalBinaryLift struct {
	n, log int32
	size   int32
}

func NewDivideIntervalBinaryLift(n int32) *DivideIntervalBinaryLift {
	log := int32(bits.Len(uint(n))) - 1
	size := n * (log + 1)
	return &DivideIntervalBinaryLift{n: n, log: log, size: size}
}

func (d *DivideIntervalBinaryLift) EnumerateRange(start int32, end int, f func(jumpId int32)) {}

func (d *DivideIntervalBinaryLift) EnumerateRange2(start1, end1 int, start2, end2 int32, f func(jumpId int32)) {
}

func (d *DivideIntervalBinaryLift) Size() int32 {
	return d.size
}

// https://leetcode.cn/problems/closest-node-to-path-in-tree/
func closestNode(n int, edges [][]int, query [][]int) []int {
	tree := make([][]int32, n)
	for _, edge := range edges {
		u, v := int32(edge[0]), int32(edge[1])
		tree[u] = append(tree[u], v)
		tree[v] = append(tree[v], u)
	}

	lca := NewDoublingLca32(tree, []int32{0})
	res := make([]int, len(query))
	for i, q := range query {
		// lca里最深的那个
		tmp := maxWithKey32(
			func(x int32) int32 { return lca.Depth[x] },
			lca.Lca(int32(q[0]), int32(q[1])),
			lca.Lca(int32(q[0]), int32(q[2])),
			lca.Lca(int32(q[1]), int32(q[2])),
		)
		res[i] = int(tmp)
	}

	return res
}

func jump() {
	// https://judge.yosupo.jp/problem/jump_on_tree
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int32
	fmt.Fscan(in, &n, &q)

	tree := make([][]int32, n)
	for i := int32(0); i < n-1; i++ {
		var u, v int32
		fmt.Fscan(in, &u, &v)
		tree[u] = append(tree[u], v)
	}
	D := NewDoublingLca32(tree, []int32{0})

	for i := int32(0); i < q; i++ {
		var from, to, k int32
		fmt.Fscan(in, &from, &to, &k)
		fmt.Fprintln(out, D.Jump(from, to, k))
	}
}

const INF int = 1e18

type DoublingLca32 struct {
	Tree  [][]int32
	Depth []int32

	n, log, size int32
	jump         [][]int32 // 节点j向上跳2^i步的父节点
}

func NewDoublingLca32(tree [][]int32, roots []int32) *DoublingLca32 {
	n := int32(len(tree))
	depth := make([]int32, n)
	lca := &DoublingLca32{
		Tree:  tree,
		Depth: depth,
		n:     n,
		log:   int32(bits.Len32(uint32(n))) - 1,
	}

	lca.makeDp()
	for _, root := range roots {
		lca.dfsAndInitDp(int32(root), -1)
	}
	lca.updateDp()
	return lca
}

func (lca *DoublingLca32) Lca(root1, root2 int32) int32 {
	if lca.Depth[root1] < lca.Depth[root2] {
		root1, root2 = root2, root1
	}
	root1 = lca.UpToDepth(root1, lca.Depth[root2])
	if root1 == root2 {
		return root1
	}
	for i := lca.log; i >= 0; i-- {
		if a, b := lca.jump[i][root1], lca.jump[i][root2]; a != b {
			root1, root2 = a, b
		}
	}
	return lca.jump[0][root1]
}

func (lca *DoublingLca32) Dist(root1, root2 int32, weighted bool) int32 {
	return lca.Depth[root1] + lca.Depth[root2] - 2*lca.Depth[lca.Lca(root1, root2)]
}

// 查询树节点root的第k个祖先(向上跳k步),如果不存在这样的祖先节点,返回 -1
func (lca *DoublingLca32) KthAncestor(root, k int32) int32 {
	if k > lca.Depth[root] {
		return -1
	}
	bit := 0
	for k > 0 {
		if k&1 == 1 {
			root = lca.jump[bit][root]
			if root == -1 {
				return -1
			}
		}
		bit++
		k >>= 1
	}
	return root
}

// 从 root 开始向上跳到指定深度 toDepth,toDepth<=dep[v],返回跳到的节点
func (lca *DoublingLca32) UpToDepth(root, toDepth int32) int32 {
	if toDepth >= lca.Depth[root] {
		return root
	}
	for i := lca.log; i >= 0; i-- {
		if (lca.Depth[root]-toDepth)&(1<<i) > 0 {
			root = lca.jump[i][root]
		}
	}
	return root
}

// 从start节点跳向target节点,跳过step个节点(0-indexed)
// 返回跳到的节点,如果不存在这样的节点,返回-1
func (lca *DoublingLca32) Jump(start, target, step int32) int32 {
	lca_ := lca.Lca(start, target)
	dep1, dep2, deplca := lca.Depth[start], lca.Depth[target], lca.Depth[lca_]
	dist := dep1 + dep2 - 2*deplca
	if step > dist {
		return -1
	}
	if step <= dep1-deplca {
		return lca.KthAncestor(start, step)
	}
	return lca.KthAncestor(target, dist-step)
}

func (lca *DoublingLca32) makeDp() {
	n, log := lca.n, lca.log
	jump := make([][]int32, log+1)
	for k := int32(0); k < log+1; k++ {
		nums := make([]int32, n)
		for i := range nums {
			nums[i] = -1 // e()
		}
		jump[k] = nums
	}
	lca.jump = jump
}

func (lca *DoublingLca32) dfsAndInitDp(cur, pre int32) {
	lca.jump[0][cur] = pre
	for _, next := range lca.Tree[cur] {
		if next != pre {
			lca.Depth[next] = lca.Depth[cur] + 1
			lca.dfsAndInitDp(next, cur)
		}
	}
}

func (lca *DoublingLca32) updateDp() {
	n, log := lca.n, lca.log
	jump := lca.jump
	for k := int32(0); k < log; k++ {
		for v := int32(0); v < n; v++ {
			j := jump[k][v]
			if j == -1 {
				jump[k+1][v] = -1 // e()
			} else {
				jump[k+1][v] = jump[k][j] // op()
			}
		}
	}
}

func max(a, b int) int {
	if a > b {
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

func maxWithKey32(key func(x int32) int32, args ...int32) int32 {
	max := args[0]
	for _, v := range args[1:] {
		if key(max) < key(v) {
			max = v
		}
	}
	return max
}
