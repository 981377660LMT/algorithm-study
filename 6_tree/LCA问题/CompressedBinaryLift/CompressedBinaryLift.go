package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	tree := [][]int32{
		{1, 2},
		{3, 4},
		{5, 6},
		{},
		{},
		{},
		{},
	}
	bl := NewCompressedBinaryLiftFromTree(tree, 0)
	fmt.Println(bl.UpToDepth(6, 0)) // 2
	// yosupo()
}

// https://judge.yosupo.jp/problem/lca
func yosupo() {
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
	bl := NewCompressedBinaryLiftFromTree(tree, 0)
	for i := 0; i < q; i++ {
		var u, v int32
		fmt.Fscan(in, &u, &v)
		fmt.Fprintln(out, bl.Lca(u, v))
	}
}

// 空间复杂度`O(n)`的树上倍增.
//   - https://taodaling.github.io/blog/2020/03/18/binary-lifting/
//   - https://codeforces.com/blog/entry/74847
//   - https://codeforces.com/blog/entry/100826
type CompressedBinaryLift struct {
	Depth  []int32
	Parent []int32
	jump   []int32 // 指向当前节点的某个祖先节点.
}

func NewCompressedBinaryLift(n int32, depthOnTree, parentOnTree []int32) *CompressedBinaryLift {
	res := &CompressedBinaryLift{
		Depth:  depthOnTree,
		Parent: parentOnTree,
		jump:   make([]int32, n),
	}
	for i := int32(0); i < n; i++ {
		res.jump[i] = -1
	}
	for i := int32(0); i < n; i++ {
		res._consider(i)
	}
	return res
}

func NewCompressedBinaryLiftFromTree(tree [][]int32, root int32) *CompressedBinaryLift {
	n := int32(len(tree))
	res := &CompressedBinaryLift{
		Depth:  make([]int32, n),
		Parent: make([]int32, n),
		jump:   make([]int32, n),
	}
	res.Parent[root] = -1
	res.jump[root] = root
	res._setUp(tree, root)
	return res
}

func (bl *CompressedBinaryLift) FirstTrue(start int32, predicate func(end int32) bool) int32 {
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

func (bl *CompressedBinaryLift) LastTrue(start int32, predicate func(end int32) bool) int32 {
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

func (bl *CompressedBinaryLift) UpToDepth(root int32, toDepth int32) int32 {
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

func (bl *CompressedBinaryLift) KthAncestor(node, k int32) int32 {
	targetDepth := bl.Depth[node] - k
	return bl.UpToDepth(node, targetDepth)
}

func (bl *CompressedBinaryLift) Lca(a, b int32) int32 {
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

func (bl *CompressedBinaryLift) Dist(a, b int32) int32 {
	return bl.Depth[a] + bl.Depth[b] - 2*bl.Depth[bl.Lca(a, b)]
}

func (bl *CompressedBinaryLift) _consider(root int32) {
	if root == -1 || bl.jump[root] != -1 {
		return
	}
	p := bl.Parent[root]
	bl._consider(p)
	bl._addLeaf(root, p)
}

func (bl *CompressedBinaryLift) _addLeaf(leaf, parent int32) {
	if parent == -1 {
		bl.jump[leaf] = leaf
	} else if tmp := bl.jump[parent]; bl.Depth[parent]-bl.Depth[tmp] == bl.Depth[tmp]-bl.Depth[bl.jump[tmp]] {
		bl.jump[leaf] = bl.jump[tmp]
	} else {
		bl.jump[leaf] = parent
	}
}

func (bl *CompressedBinaryLift) _setUp(tree [][]int32, root int32) {
	queue := []int32{root}
	head := 0
	for head < len(queue) {
		cur := queue[head]
		head++
		nexts := tree[cur]
		for _, next := range nexts {
			if next == bl.Parent[cur] {
				continue
			}
			bl.Depth[next] = bl.Depth[cur] + 1
			bl.Parent[next] = cur
			queue = append(queue, next)
			bl._addLeaf(next, cur)
		}
	}
}
