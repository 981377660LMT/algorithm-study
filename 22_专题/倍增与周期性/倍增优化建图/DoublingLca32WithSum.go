// DividePathOnTreeBinaryLift/DoublingLca
// 倍增拆树上路径`path(from,to)`：倍增拆点将树上的一段路径拆成logn个点
//
// 一共拆分成[0,log]层，每层有n个元素.
// !jumpId = level*n + index 表示第level层的第index个元素(0<=level<log+1,0<=index<n).

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {
	// jump()
	assert()
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
	D := NewDoublingLca32WithSum(tree, []int32{0}, make([]int, n))

	for i := int32(0); i < q; i++ {
		var from, to, k int32
		fmt.Fscan(in, &from, &to, &k)
		fmt.Fprintln(out, D.Jump(from, to, k))
	}
}

const INF int = 1e18

type S = int

func (*DoublingLca32WithSum) e() S          { return 0 }
func (*DoublingLca32WithSum) op(e1, e2 S) S { return e1 + e2 }

type DoublingLca32WithSum struct {
	Tree  [][]int32
	Depth []int32

	n, log, size int32
	jump         [][]int32 // 节点j向上跳2^i步的父节点
	data         [][]S
	values       []S
}

// values: 每个点的`点权`.
// 如果需要查询边权，则每个点的`点权`设为`该点与其父亲结点的边权`, 根节点的`点权`设为`幺元`.
func NewDoublingLca32WithSum(tree [][]int32, roots []int32, values []S) *DoublingLca32WithSum {
	n := int32(len(tree))
	depth := make([]int32, n)
	lca := &DoublingLca32WithSum{
		Tree:   tree,
		Depth:  depth,
		n:      n,
		log:    int32(bits.Len32(uint32(n))) - 1,
		values: values,
	}
	lca.size = n * (lca.log + 1)

	lca.makeDp()
	for _, root := range roots {
		lca.dfsAndInitDp(root, -1)
	}
	lca.updateDp()
	return lca
}

func (lca *DoublingLca32WithSum) Size() int32 { return lca.size }
func (lca *DoublingLca32WithSum) Log() int32  { return lca.log }

func (lca *DoublingLca32WithSum) FirstTrue(start int32, predicate func(end int32) bool) int32 {
	if predicate(start) {
		return start
	}
	for k := lca.log; k >= 0; k-- {
		tmp := lca.jump[k][start]
		if tmp != -1 && !predicate(tmp) {
			start = tmp
		}
	}
	return lca.jump[0][start] // 不存在则返回-1
}

func (lca *DoublingLca32WithSum) FirstTrueWithSum(start int32, predicate func(end int32, sum S) bool, isEdge bool) (int32, S) {
	if isEdge {
		sum := lca.e()
		if predicate(start, sum) {
			return start, sum
		}
		for k := lca.log; k >= 0; k-- {
			next := lca.jump[k][start]
			if next == -1 {
				continue
			}
			nextSum := lca.op(sum, lca.data[k][start])
			if !predicate(next, nextSum) {
				start = next
				sum = nextSum
			}
		}
		p := lca.jump[0][start]
		if p == -1 {
			return -1, sum
		}
		return p, lca.op(sum, lca.data[0][start])
	} else {
		if predicate(start, lca.values[start]) {
			return start, lca.values[start]
		}
		sum := lca.e()
		for k := lca.log; k >= 0; k-- {
			next := lca.jump[k][start]
			if next == -1 {
				continue
			}
			nextSum1 := lca.op(sum, lca.data[k][start])
			nextSum2 := lca.op(nextSum1, lca.values[next])
			if !predicate(next, nextSum2) {
				start = next
				sum = nextSum1
			}
		}
		p := lca.jump[0][start]
		if p == -1 {
			sum = lca.op(sum, lca.values[start])
			return -1, sum
		}
		sum = lca.op(sum, lca.data[0][start])
		sum = lca.op(sum, lca.values[p])
		return p, sum
	}
}

func (lca *DoublingLca32WithSum) LastTrue(start int32, predicate func(end int32) bool) int32 {
	if !predicate(start) {
		return -1
	}
	for k := lca.log; k >= 0; k-- {
		tmp := lca.jump[k][start]
		if tmp != -1 && predicate(tmp) {
			start = tmp
		}
	}
	return start
}

func (lca *DoublingLca32WithSum) LastTrueWithSum(start int32, predicate func(end int32, sum S) bool, isEdge bool) (int32, S) {
	if isEdge {
		sum := lca.e()
		if !predicate(start, sum) {
			return -1, sum
		}
		for k := lca.log; k >= 0; k-- {
			next := lca.jump[k][start]
			if next == -1 {
				continue
			}
			nextSum := lca.op(sum, lca.data[k][start])
			if predicate(next, nextSum) {
				start = next
				sum = nextSum
			}
		}
		return start, sum
	} else {
		if !predicate(start, lca.values[start]) {
			return -1, lca.values[start]
		}
		sum := lca.e()
		for k := lca.log; k >= 0; k-- {
			next := lca.jump[k][start]
			if next == -1 {
				continue
			}
			nextSum1 := lca.op(sum, lca.data[k][start])
			nextSum2 := lca.op(nextSum1, lca.values[next])
			if predicate(next, nextSum2) {
				start = next
				sum = nextSum1
			}
		}
		sum = lca.op(sum, lca.values[start])
		return start, sum
	}
}

func (lca *DoublingLca32WithSum) Lca(root1, root2 int32) int32 {
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

// 查询路径`a`到`b`的聚合值.
// isEdge 是否是边权.
func (lca *DoublingLca32WithSum) LcaWithSum(root1, root2 int32, isEdge bool) (int32, S) {
	var e S
	if lca.Depth[root1] > lca.Depth[root2] {
		end, sum := lca.UpToDepthWithSum(root1, lca.Depth[root2], true)
		root1, e = end, sum
	} else if lca.Depth[root1] < lca.Depth[root2] {
		end, sum := lca.UpToDepthWithSum(root2, lca.Depth[root1], true)
		root2, e = end, sum
	} else {
		e = lca.e()
	}
	if root1 == root2 {
		if !isEdge {
			e = lca.op(e, lca.values[root1])
		}
		return root1, e
	}
	for i := lca.log; i >= 0; i-- {
		if a, b := lca.jump[i][root1], lca.jump[i][root2]; a != b {
			e = lca.op(e, lca.data[i][root1])
			e = lca.op(e, lca.data[i][root2])
			root1, root2 = a, b
		}
	}
	e = lca.op(e, lca.values[root1])
	e = lca.op(e, lca.values[root2])
	p := lca.jump[0][root1]
	if !isEdge {
		e = lca.op(e, lca.values[p])
	}
	return p, e
}

// 从 root 开始向上跳到指定深度 toDepth,toDepth<=dep[v],返回跳到的节点
func (lca *DoublingLca32WithSum) UpToDepth(root, toDepth int32) int32 {
	if !(0 <= toDepth && toDepth <= lca.Depth[root]) {
		return -1
	}
	for i := lca.log; i >= 0; i-- {
		if (lca.Depth[root]-toDepth)&(1<<i) > 0 {
			root = lca.jump[i][root]
		}
	}
	return root
}

func (lca *DoublingLca32WithSum) UpToDepthWithSum(root, toDepth int32, isEdge bool) (int32, S) {
	sum := lca.e()
	if !(0 <= toDepth && toDepth <= lca.Depth[root]) {
		return -1, sum
	}
	for i := lca.log; i >= 0; i-- {
		if (lca.Depth[root]-toDepth)&(1<<i) > 0 {
			sum = lca.op(sum, lca.data[i][root])
			root = lca.jump[i][root]
		}
	}
	if !isEdge {
		sum = lca.op(sum, lca.values[root])
	}
	return root, sum
}

// 查询树节点root的第k个祖先(向上跳k步),如果不存在这样的祖先节点,返回 -1
func (lca *DoublingLca32WithSum) KthAncestor(root, k int32) int32 {
	targetDepth := lca.Depth[root] - k
	return lca.UpToDepth(root, targetDepth)
}

func (lca *DoublingLca32WithSum) KthAncestorWithSum(root, k int32, isEdge bool) (int32, S) {
	targetDepth := lca.Depth[root] - k
	return lca.UpToDepthWithSum(root, targetDepth, isEdge)
}

// 从start节点跳向target节点,跳过step个节点(0-indexed)
// 返回跳到的节点,如果不存在这样的节点,返回-1
func (lca *DoublingLca32WithSum) Jump(start, target, step int32) int32 {
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

func (lca *DoublingLca32WithSum) Dist(root1, root2 int32, weighted bool) int32 {
	return lca.Depth[root1] + lca.Depth[root2] - 2*lca.Depth[lca.Lca(root1, root2)]
}

func (lca *DoublingLca32WithSum) makeDp() {
	n, log := lca.n, lca.log
	jump := make([][]int32, log+1)
	data := make([][]S, log+1)
	for k := int32(0); k < log+1; k++ {
		nums1 := make([]int32, n)
		nums2 := make([]S, n)
		// for i := range nums1 {
		// 	nums1[i] = -1
		// 	nums2[i] = lca.e()
		// }
		jump[k] = nums1
		data[k] = nums2
	}
	lca.jump = jump
	lca.data = data
}

func (lca *DoublingLca32WithSum) dfsAndInitDp(cur, pre int32) {
	lca.jump[0][cur] = pre
	lca.data[0][cur] = lca.values[cur] // jump(0, cur)：cur到父节点的边权
	for _, next := range lca.Tree[cur] {
		if next != pre {
			lca.Depth[next] = lca.Depth[cur] + 1
			lca.dfsAndInitDp(next, cur)
		}
	}
}

func (lca *DoublingLca32WithSum) updateDp() {
	n, log := lca.n, lca.log
	jump, data, op := lca.jump, lca.data, lca.op
	for k := int32(0); k < log; k++ {
		for v := int32(0); v < n; v++ {
			j := jump[k][v]
			if j == -1 {
				jump[k+1][v] = -1 // e()
			} else {
				jump[k+1][v] = jump[k][j]
				data[k+1][v] = op(data[k][v], data[k][j]) // op()
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

func assert() {
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
	values := []int{1, 1, 2, 3, 4, 5, 6}
	bl := NewDoublingLca32WithSum(tree, []int32{0}, values)

	type pair struct {
		node int32
		sum  S
	}

	// firstTrueWithSum
	node, sum := bl.FirstTrueWithSum(6, func(i int32, sum S) bool { return sum >= 0 }, true)
	expect(pair{node, sum}, pair{6, 0})
	node, sum = bl.FirstTrueWithSum(6, func(i int32, sum S) bool { return sum >= 6 }, true)
	expect(pair{node, sum}, pair{4, 6})
	node, sum = bl.FirstTrueWithSum(6, func(i int32, sum S) bool { return sum >= 10 }, true)
	expect(pair{node, sum}, pair{1, 10})
	node, sum = bl.FirstTrueWithSum(6, func(i int32, sum S) bool { return sum >= 11 }, true)
	expect(pair{node, sum}, pair{0, 11})
	node, sum = bl.FirstTrueWithSum(6, func(i int32, sum S) bool { return sum >= 15 }, true)
	expect(pair{node, sum}, pair{-1, 11})

	node, sum = bl.FirstTrueWithSum(6, func(i int32, sum S) bool { return sum >= 0 }, false)
	expect(pair{node, sum}, pair{6, 6})
	node, sum = bl.FirstTrueWithSum(6, func(i int32, sum S) bool { return sum >= 6 }, false)
	expect(pair{node, sum}, pair{6, 6})
	node, sum = bl.FirstTrueWithSum(6, func(i int32, sum S) bool { return sum >= 10 }, false)
	expect(pair{node, sum}, pair{4, 10})
	node, sum = bl.FirstTrueWithSum(6, func(i int32, sum S) bool { return sum >= 11 }, false)
	expect(pair{node, sum}, pair{1, 11})
	node, sum = bl.FirstTrueWithSum(6, func(i int32, sum S) bool { return sum >= 12 }, false)
	expect(pair{node, sum}, pair{0, 12})
	node, sum = bl.FirstTrueWithSum(6, func(i int32, sum S) bool { return sum >= 15 }, false)
	expect(pair{node, sum}, pair{-1, 12})
	node, sum = bl.FirstTrueWithSum(6, func(i int32, sum S) bool { return bl.Depth[i] <= 1 }, true)
	expect(pair{node, sum}, pair{1, 10})
	node, sum = bl.FirstTrueWithSum(6, func(i int32, sum S) bool { return bl.Depth[i] <= 1 }, false)
	expect(pair{node, sum}, pair{1, 11})

	// lastTrueWithSum
	node, sum = bl.LastTrueWithSum(6, func(i int32, sum S) bool { return sum <= -1 }, true)
	expect(pair{node, sum}, pair{-1, 0})
	node, sum = bl.LastTrueWithSum(6, func(i int32, sum S) bool { return sum <= 0 }, true)
	expect(pair{node, sum}, pair{6, 0})
	node, sum = bl.LastTrueWithSum(6, func(i int32, sum S) bool { return sum <= 5 }, true)
	expect(pair{node, sum}, pair{6, 0})
	node, sum = bl.LastTrueWithSum(6, func(i int32, sum S) bool { return sum <= 6 }, true)
	expect(pair{node, sum}, pair{4, 6})
	node, sum = bl.LastTrueWithSum(6, func(i int32, sum S) bool { return sum <= 7 }, true)
	expect(pair{node, sum}, pair{4, 6})
	node, sum = bl.LastTrueWithSum(6, func(i int32, sum S) bool { return sum <= 10 }, true)
	expect(pair{node, sum}, pair{1, 10})
	node, sum = bl.LastTrueWithSum(6, func(i int32, sum S) bool { return sum <= 11 }, true)
	expect(pair{node, sum}, pair{0, 11})
	node, sum = bl.LastTrueWithSum(6, func(i int32, sum S) bool { return sum <= 12 }, true)
	expect(pair{node, sum}, pair{0, 11})
	node, sum = bl.LastTrueWithSum(6, func(i int32, sum S) bool { return sum <= 13 }, true)
	expect(pair{node, sum}, pair{0, 11})

	node, sum = bl.LastTrueWithSum(6, func(i int32, sum S) bool { return sum <= -1 }, false)
	expect(pair{node, sum}, pair{-1, 6})
	node, sum = bl.LastTrueWithSum(6, func(i int32, sum S) bool { return sum <= 0 }, false)
	expect(pair{node, sum}, pair{-1, 6})
	node, sum = bl.LastTrueWithSum(6, func(i int32, sum S) bool { return sum <= 5 }, false)
	expect(pair{node, sum}, pair{-1, 6})
	node, sum = bl.LastTrueWithSum(6, func(i int32, sum S) bool { return sum <= 6 }, false)
	expect(pair{node, sum}, pair{6, 6})
	node, sum = bl.LastTrueWithSum(6, func(i int32, sum S) bool { return sum <= 7 }, false)
	expect(pair{node, sum}, pair{6, 6})
	node, sum = bl.LastTrueWithSum(6, func(i int32, sum S) bool { return sum <= 10 }, false)
	expect(pair{node, sum}, pair{4, 10})
	node, sum = bl.LastTrueWithSum(6, func(i int32, sum S) bool { return sum <= 11 }, false)
	expect(pair{node, sum}, pair{1, 11})
	node, sum = bl.LastTrueWithSum(6, func(i int32, sum S) bool { return sum <= 12 }, false)
	expect(pair{node, sum}, pair{0, 12})
	node, sum = bl.LastTrueWithSum(6, func(i int32, sum S) bool { return sum <= 13 }, false)
	expect(pair{node, sum}, pair{0, 12})
	node, sum = bl.LastTrueWithSum(6, func(i int32, sum S) bool { return bl.Depth[i] >= 2 }, true)
	expect(pair{node, sum}, pair{4, 6})
	node, sum = bl.LastTrueWithSum(6, func(i int32, sum S) bool { return bl.Depth[i] >= 2 }, false)
	expect(pair{node, sum}, pair{4, 10})

	// upToDepthWithSum
	type uptoDepthWithSumArgs struct {
		root, toDepth int32
		isEdge        bool
	}
	args := []uptoDepthWithSumArgs{{6, 1, true}, {6, 1, false}, {6, 2, true}, {6, 2, false}, {6, 3, true}, {6, 3, false}, {6, 4, true}, {6, 4, false}}
	expected := []pair{{1, 10}, {1, 11}, {4, 6}, {4, 10}, {6, 0}, {6, 6}, {-1, 0}, {-1, 0}}
	for i, arg := range args {
		node, sum := bl.UpToDepthWithSum(arg.root, arg.toDepth, arg.isEdge)
		expect(pair{node, sum}, expected[i])
	}

	// kthAncestorWithSum
	args = []uptoDepthWithSumArgs{{6, 0, true}, {6, 0, false}, {6, 1, true}, {6, 1, false}, {6, 2, true}, {6, 2, false}, {6, 3, true}, {6, 3, false}, {6, 4, true}, {6, 4, false}}
	expected = []pair{{6, 0}, {6, 6}, {4, 6}, {4, 10}, {1, 10}, {1, 11}, {0, 11}, {0, 12}, {-1, 0}, {-1, 0}}
	for i, arg := range args {
		node, sum := bl.KthAncestorWithSum(arg.root, arg.toDepth, arg.isEdge)
		expect(pair{node, sum}, expected[i])
	}

	// lcaWithSum
	weigthSum := func(u, v int32, isEdge bool) S {
		if bl.Depth[u] < bl.Depth[v] {
			u, v = v, u
		}
		sum := bl.e()
		for bl.Depth[u] > bl.Depth[v] {
			sum = bl.op(sum, values[u])
			u = bl.jump[0][u]
		}
		for u != v {
			sum = bl.op(sum, values[u])
			sum = bl.op(sum, values[v])
			u = bl.jump[0][u]
			v = bl.jump[0][v]
		}
		if !isEdge {
			sum = bl.op(sum, values[u])
		}
		return sum
	}

	for i := int32(0); i < int32(n); i++ {
		for j := int32(0); j < int32(n); j++ {
			lca := bl.Lca(i, j)
			node, sum := bl.LcaWithSum(i, j, true)
			expect(pair{node, sum}, pair{lca, weigthSum(i, j, true)})
			node, sum = bl.LcaWithSum(i, j, false)
			expect(pair{node, sum}, pair{lca, weigthSum(i, j, false)})
		}
	}

	fmt.Println("pass")
}

func expect[S comparable](actual, expected S) {
	if actual != expected {
		panic(fmt.Sprintf("actual: %v, expected: %v", actual, expected))
	}
}
