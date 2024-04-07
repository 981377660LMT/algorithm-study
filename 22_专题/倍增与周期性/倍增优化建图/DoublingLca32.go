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
	// lca()
	// jump()
	test()
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

// https://judge.yosupo.jp/problem/lca
func lca() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)

	tree := make([][]int32, n)
	for i := 1; i < n; i++ {
		var parent int32
		fmt.Fscan(in, &parent)
		tree[parent] = append(tree[parent], int32(i))
	}
	bl := NewDoublingLca32(tree, []int32{0})
	for i := 0; i < q; i++ {
		var u, v int32
		fmt.Fscan(in, &u, &v)
		fmt.Fprintln(out, bl.Lca(u, v))
	}
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
		tree[v] = append(tree[v], u)
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
		lca.dfsAndInitDp(root, -1)
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

func (lca *DoublingLca32) FirstTrue(start int32, predicate func(end int32) bool) (step, to int32) {
	// 判定条件取反，然后向上跳一层.
	if predicate(start) {
		step, to = 0, start
		return
	}

	for k := lca.log; k >= 0; k-- {
		tmp := lca.jump[k][start]
		if tmp == -1 {
			continue
		}
		if !predicate(tmp) {
			step |= 1 << k
			start = tmp
		}
	}

	if p := lca.jump[0][start]; p == -1 {
		step, to = -1, -1
	} else {
		step, to = step+1, p
	}
	return
}

func (lca *DoublingLca32) LastTrue(start int32, predicate func(end int32) bool) (step, to int32) {
	if !predicate(start) {
		step, to = -1, -1
		return
	}
	for k := lca.log; k >= 0; k-- {
		tmp := lca.jump[k][start]
		if tmp == -1 {
			continue
		}
		if predicate(tmp) {
			step |= 1 << k
			start = tmp
		}
	}
	to = start
	return
}

func (lca *DoublingLca32) makeDp() {
	n, log := lca.n, lca.log
	jump := make([][]int32, log+1)
	for k := int32(0); k < log+1; k++ {
		nums := make([]int32, n)
		// for i := range nums {
		// 	nums[i] = -1		// e()
		// }
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

func test() {

	//          0
	//        /   \
	//       1     2
	//      / \     \
	//     3   4     5
	//         /
	//        6

	n := 7
	edges := [][]int32{{0, 1}, {0, 2}, {1, 3}, {1, 4}, {2, 5}, {4, 6}}
	tree := make([][]int32, n)
	for _, e := range edges {
		u, v := e[0], e[1]
		tree[u] = append(tree[u], v)
		tree[v] = append(tree[v], u)
	}

	lca := NewDoublingLca32(tree, []int32{0})
	fmt.Println(lca.Lca(3, 6)) // 1

	step1, to1 := lca.FirstTrue(6, func(i int32) bool { return lca.Depth[i] <= -1 })
	expect[int32](step1, -1)
	expect[int32](to1, -1)
	step1, to1 = lca.FirstTrue(6, func(i int32) bool { return lca.Depth[i] <= 1 })
	expect[int32](step1, 2)
	expect[int32](to1, 1)
	step1, to1 = lca.FirstTrue(6, func(i int32) bool { return lca.Depth[i] <= 2 })
	expect[int32](step1, 1)
	expect[int32](to1, 4)
	step1, to1 = lca.FirstTrue(6, func(i int32) bool { return lca.Depth[i] <= 3 })
	expect[int32](step1, 0)
	expect[int32](to1, 6)
	step1, to1 = lca.FirstTrue(6, func(i int32) bool { return lca.Depth[i] <= 4 })
	expect[int32](step1, 0)
	expect[int32](to1, 6)
	step1, to1 = lca.FirstTrue(6, func(i int32) bool { return lca.Depth[i] >= 2 })
	expect[int32](step1, 0)
	expect[int32](to1, 6)
	step1, to1 = lca.FirstTrue(6, func(i int32) bool { return lca.Depth[i] >= 4 })
	expect[int32](step1, -1)
	expect[int32](to1, -1)

	step2, to2 := lca.LastTrue(6, func(i int32) bool { return lca.Depth[i] >= 1 })
	expect[int32](step2, 2)
	expect[int32](to2, 1)
	step2, to2 = lca.LastTrue(6, func(i int32) bool { return lca.Depth[i] >= 2 })
	expect[int32](step2, 1)
	expect[int32](to2, 4)
	step2, to2 = lca.LastTrue(6, func(i int32) bool { return lca.Depth[i] >= 3 })
	expect[int32](step2, 0)
	expect[int32](to2, 6)
	step2, to2 = lca.LastTrue(6, func(i int32) bool { return lca.Depth[i] >= 4 })
	expect[int32](step2, -1)
	expect[int32](to2, -1)
	step2, to2 = lca.LastTrue(6, func(i int32) bool { return lca.Depth[i] <= 2 })
	expect[int32](step2, -1)
	expect[int32](to2, -1)
	step2, to2 = lca.LastTrue(6, func(i int32) bool { return lca.Depth[i] <= 4 })
	expect[int32](step2, 3)
	expect[int32](to2, 0)
	step2, to2 = lca.LastTrue(6, func(i int32) bool { return lca.Depth[i] <= -1 })
	expect[int32](step2, -1)
	expect[int32](to2, -1)

	fmt.Println("test passed")
}

func expect[S comparable](actual, expected S) {
	if actual != expected {
		panic(fmt.Sprintf("actual: %v, expected: %v", actual, expected))
	}
}
