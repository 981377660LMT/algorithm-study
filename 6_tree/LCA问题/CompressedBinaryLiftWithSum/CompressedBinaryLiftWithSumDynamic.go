package main

import (
	"fmt"
)

func main() {
	assert()
}

type S = int

func (*CompressedBinaryLiftWithSumDynamic) e() S          { return 0 }
func (*CompressedBinaryLiftWithSumDynamic) op(e1, e2 S) S { return e1 + e2 }

// 空间复杂度`O(n)`的树上倍增，用于倍增结构优化建图、查询路径聚合值.
//   - https://taodaling.github.io/blog/2020/03/18/binary-lifting/
//   - https://codeforces.com/blog/entry/74847
//   - https://codeforces.com/blog/entry/100826
type CompressedBinaryLiftWithSumDynamic struct {
	Depth       []int32
	Parent      []int32
	jump        []int32 // 指向当前节点的某个祖先节点.
	attachments []S     // 从当前结点到`jump`结点的路径上的聚合值(不包含`jump`结点).
	singles     []S     // 当前结点的聚合值.
}

// 不预先给出整棵树,而是动态添加叶子节点,维护树节点的LCA和k级祖先.
func NewCompressedBinaryLiftWithSumDynamic(n int32) *CompressedBinaryLiftWithSumDynamic {
	res := &CompressedBinaryLiftWithSumDynamic{
		Depth:       make([]int32, n),
		Parent:      make([]int32, n),
		jump:        make([]int32, n),
		attachments: make([]S, n),
		singles:     make([]S, n),
	}
	for i := int32(0); i < n; i++ {
		res.attachments[i] = res.e()
	}
	return res
}

func NewCompressedBinaryLiftWithSumDynamicWithRoot(n int32, root int32, weight S) *CompressedBinaryLiftWithSumDynamic {
	res := NewCompressedBinaryLiftWithSumDynamic(n)
	res.AddLeaf(root, -1, weight)
	return res
}

// 在树中添加一条从parent到leaf的边，要求parent已经存在于树中或为-1.
// parent为-1时，leaf为根节点.
func (bl *CompressedBinaryLiftWithSumDynamic) AddLeaf(leaf, parent int32, weight S) {
	if parent == -1 {
		bl.Parent[leaf] = -1
		bl.jump[leaf] = leaf
		bl.singles[leaf] = weight
		return
	}
	bl.Parent[leaf] = parent
	bl.Depth[leaf] = bl.Depth[parent] + 1
	bl.singles[leaf] = weight
	if tmp := bl.jump[parent]; bl.Depth[parent]-bl.Depth[tmp] == bl.Depth[tmp]-bl.Depth[bl.jump[tmp]] {
		bl.jump[leaf] = bl.jump[tmp]
		bl.attachments[leaf] = bl.op(weight, bl.attachments[parent])
		bl.attachments[leaf] = bl.op(bl.attachments[leaf], bl.attachments[tmp])
	} else {
		bl.jump[leaf] = parent
		bl.attachments[leaf] = weight // copy
	}
}

func (bl *CompressedBinaryLiftWithSumDynamic) FirstTrue(start int32, predicate func(end int32) bool) int32 {
	for !predicate(start) {
		if predicate(bl.jump[start]) {
			start = bl.Parent[start]
		} else {
			if start == bl.jump[start] {
				return -1
			}
			start = bl.jump[start]
		}
	}
	return start
}

func (bl *CompressedBinaryLiftWithSumDynamic) FirstTrueWithSum(start int32, predicate func(end int32, sum S) bool, isEdge bool) (int32, S) {
	if isEdge {
		sum := bl.e() // 不包含_singles[start]
		for {
			if predicate(start, sum) {
				return start, sum
			}
			jumpStart, jumpSum := bl.jump[start], bl.op(sum, bl.attachments[start])
			if predicate(jumpStart, jumpSum) {
				sum = bl.op(sum, bl.singles[start])
				start = bl.Parent[start]
			} else {
				if start == jumpStart {
					return -1, jumpSum
				}
				sum = jumpSum
				start = jumpStart
			}
		}
	} else {
		sum := bl.e() // 不包含_singles[start]
		for {
			sumWithSingle := bl.op(sum, bl.singles[start])
			if predicate(start, sumWithSingle) {
				return start, sumWithSingle
			}
			jumpStart, jumpSum1 := bl.jump[start], bl.op(sum, bl.attachments[start])
			jumpSum2 := bl.op(jumpSum1, bl.singles[jumpStart])
			if predicate(jumpStart, jumpSum2) {
				sum = sumWithSingle
				start = bl.Parent[start]
			} else {
				if start == jumpStart {
					return -1, jumpSum2
				}
				sum = jumpSum1
				start = jumpStart
			}
		}
	}
}

func (bl *CompressedBinaryLiftWithSumDynamic) LastTrue(start int32, predicate func(end int32) bool) int32 {
	if !predicate(start) {
		return -1
	}
	for {
		if predicate(bl.jump[start]) {
			if start == bl.jump[start] {
				return start
			}
			start = bl.jump[start]
		} else if predicate(bl.Parent[start]) {
			start = bl.Parent[start]
		} else {
			return start
		}
	}
}

func (bl *CompressedBinaryLiftWithSumDynamic) LastTrueWithSum(start int32, predicate func(end int32, sum S) bool, isEdge bool) (int32, S) {
	if isEdge {
		sum := bl.e() // 不包含_singles[start]
		if !predicate(start, sum) {
			return -1, sum
		}
		for {
			jumpStart, jumpSum := bl.jump[start], bl.op(sum, bl.attachments[start])
			if predicate(jumpStart, jumpSum) {
				if start == jumpStart {
					return start, sum
				}
				sum = jumpSum
				start = jumpStart
			} else {
				parentStart, parentSum := bl.Parent[start], bl.op(sum, bl.singles[start])
				if predicate(parentStart, parentSum) {
					sum = parentSum
					start = parentStart
				} else {
					return start, sum
				}
			}
		}
	} else {
		if !predicate(start, bl.singles[start]) {
			return -1, bl.singles[start]
		}
		sum := bl.e() // 不包含_singles[start]
		for {
			jumpStart, jumpSum1 := bl.jump[start], bl.op(sum, bl.attachments[start])
			jumpSum2 := bl.op(jumpSum1, bl.singles[jumpStart])
			if predicate(jumpStart, jumpSum2) {
				if start == jumpStart {
					return start, jumpSum2
				}
				sum = jumpSum1
				start = jumpStart
			} else {
				parentStart, parentSum1 := bl.Parent[start], bl.op(sum, bl.singles[start])
				parentSum2 := bl.op(parentSum1, bl.singles[parentStart])
				if predicate(parentStart, parentSum2) {
					sum = parentSum1
					start = parentStart
				} else {
					return start, parentSum1
				}
			}
		}
	}
}

func (bl *CompressedBinaryLiftWithSumDynamic) UpToDepth(root int32, toDepth int32) int32 {
	if !(0 <= toDepth && toDepth <= bl.Depth[root]) {
		return -1
	}
	for bl.Depth[root] > toDepth {
		if bl.Depth[bl.jump[root]] < toDepth {
			root = bl.Parent[root]
		} else {
			root = bl.jump[root]
		}
	}
	return root
}

func (bl *CompressedBinaryLiftWithSumDynamic) UpToDepthWithSum(root int32, toDepth int32, isEdge bool) (int32, S) {
	sum := bl.e() // 不包含_singles[root]
	if !(0 <= toDepth && toDepth <= bl.Depth[root]) {
		return -1, sum
	}
	for bl.Depth[root] > toDepth {
		if bl.Depth[bl.jump[root]] < toDepth {
			sum = bl.op(sum, bl.singles[root])
			root = bl.Parent[root]
		} else {
			sum = bl.op(sum, bl.attachments[root])
			root = bl.jump[root]
		}
	}
	if !isEdge {
		sum = bl.op(sum, bl.singles[root])
	}
	return root, sum
}

func (bl *CompressedBinaryLiftWithSumDynamic) KthAncestor(node, k int32) int32 {
	targetDepth := bl.Depth[node] - k
	return bl.UpToDepth(node, targetDepth)
}

func (bl *CompressedBinaryLiftWithSumDynamic) KthAncestorWithSum(node, k int32, isEdge bool) (int32, S) {
	targetDepth := bl.Depth[node] - k
	return bl.UpToDepthWithSum(node, targetDepth, isEdge)
}

func (bl *CompressedBinaryLiftWithSumDynamic) Lca(a, b int32) int32 {
	if bl.Depth[a] > bl.Depth[b] {
		a = bl.KthAncestor(a, bl.Depth[a]-bl.Depth[b])
	} else if bl.Depth[a] < bl.Depth[b] {
		b = bl.KthAncestor(b, bl.Depth[b]-bl.Depth[a])
	}
	for a != b {
		if bl.jump[a] == bl.jump[b] {
			a = bl.Parent[a]
			b = bl.Parent[b]
		} else {
			a = bl.jump[a]
			b = bl.jump[b]
		}
	}
	return a
}

// 查询路径`a`到`b`的聚合值.
// isEdge 是否是边权.
func (bl *CompressedBinaryLiftWithSumDynamic) LcaWithSum(a, b int32, isEdge bool) (int32, S) {
	var e S // 不包含_singles[a]和_singles[b]
	if bl.Depth[a] > bl.Depth[b] {
		end, sum := bl.UpToDepthWithSum(a, bl.Depth[b], true)
		a, e = end, sum
	} else if bl.Depth[a] < bl.Depth[b] {
		end, sum := bl.UpToDepthWithSum(b, bl.Depth[a], true)
		b, e = end, sum
	} else {
		e = bl.e()
	}
	for a != b {
		if bl.jump[a] == bl.jump[b] {
			e = bl.op(e, bl.singles[a])
			e = bl.op(e, bl.singles[b])
			a = bl.Parent[a]
			b = bl.Parent[b]
		} else {
			e = bl.op(e, bl.attachments[a])
			e = bl.op(e, bl.attachments[b])
			a = bl.jump[a]
			b = bl.jump[b]
		}
	}
	if !isEdge {
		e = bl.op(e, bl.singles[a])
	}
	return a, e
}

func (bl *CompressedBinaryLiftWithSumDynamic) Dist(a, b int32) int32 {
	return bl.Depth[a] + bl.Depth[b] - 2*bl.Depth[bl.Lca(a, b)]
}

func assert() {
	//          0
	//        /   \
	//       1     2
	//      / \     \
	//     3   4     5
	//         /
	//        6

	n := int32(7)
	edges := [][]int32{{0, 1}, {0, 2}, {1, 3}, {1, 4}, {2, 5}, {4, 6}}
	tree := make([][]int32, n)
	for _, e := range edges {
		u, v := e[0], e[1]
		tree[u] = append(tree[u], v)
		tree[v] = append(tree[v], u)
	}
	values := []int{1, 1, 2, 3, 4, 5, 6}
	bl := NewCompressedBinaryLiftWithSumDynamicWithRoot(n, 0, values[0])
	var dfs func(int32, int32)
	dfs = func(cur, pre int32) {
		for _, next := range tree[cur] {
			if next == pre {
				continue
			}
			bl.AddLeaf(next, cur, values[next])
			dfs(next, cur)
		}
	}
	dfs(0, -1)

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
			sum = bl.op(sum, bl.singles[u])
			u = bl.Parent[u]
		}
		for u != v {
			sum = bl.op(sum, bl.singles[u])
			sum = bl.op(sum, bl.singles[v])
			u = bl.Parent[u]
			v = bl.Parent[v]
		}
		if !isEdge {
			sum = bl.op(sum, bl.singles[u])
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
