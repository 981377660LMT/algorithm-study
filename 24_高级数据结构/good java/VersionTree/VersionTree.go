// 版本树/操作树 OperationTree

package main

import (
	"bufio"
	"fmt"
	"os"
)

// https://www.luogu.com.cn/problem/CF707D
// Persistent Bookcase
// 1 i j : 将(i,j)置为1
// 2 i j : 将(i,j)置为0
// 3 i j : 将i行01反转
// 4 k : 回到第k次操作后的状态(0 <= k < q)
// 每次操作后，输出全局 1 的个数

// TODO: 添加查询的api
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m, q int32
	fmt.Fscan(in, &n, &m, &q)
	sets := make([]*bitsetFastFlipAll, n)
	for i := range sets {
		sets[i] = newBitsetFastFlipAll(m, false)
	}

	onesCountAll := int32(0)
	res := make([]int32, q)
	tree := NewVersionTree(q)
	for i := int32(0); i < q; i++ {
		var op uint8
		fmt.Fscan(in, &op)

		switch op {
		case 1:
			var i, j int32
			fmt.Fscan(in, &i, &j)
			i--
			j--
			tree.Apply(
				func() {
					if sets[i].Add(j) {
						onesCountAll++
					}
				},
				func() {
					if sets[i].Discard(j) {
						onesCountAll--
					}
				},
			)
		case 2:
			var i, j int32
			fmt.Fscan(in, &i, &j)
			i--
			j--
			tree.Apply(
				func() {
					if sets[i].Discard(j) {
						onesCountAll--
					}
				},
				func() {
					if sets[i].Add(j) {
						onesCountAll++
					}
				},
			)
		case 3:
			var i int32
			fmt.Fscan(in, &i)
			i--
			tree.Apply(
				func() {
					onesCountAll -= sets[i].OnesCount()
					sets[i].FlipAll()
					onesCountAll += sets[i].OnesCount()
				},
				func() {
					onesCountAll -= sets[i].OnesCount()
					sets[i].FlipAll()
					onesCountAll += sets[i].OnesCount()
				},
			)
		case 4:
			var k int32
			fmt.Fscan(in, &k)
			tree.SwitchVersion(k)
		}

		tree.Apply(
			func() {
				res[i] = onesCountAll
			},
			func() {},
		)
	}

	tree.Run()

	for _, v := range res {
		fmt.Fprintln(out, v)
	}
}

type bitsetFastFlipAll struct {
	flip      bool
	n         int32
	onesCount int32
	data      []uint64
}

func newBitsetFastFlipAll(n int32, filled bool) *bitsetFastFlipAll {
	data := make([]uint64, n>>6+1)
	onesCount := int32(0)
	if filled {
		for i := range data {
			data[i] = ^uint64(0)
		}
		if n != 0 {
			data[len(data)-1] >>= int32(len(data)<<6) - n
		}
		onesCount = n
	}
	return &bitsetFastFlipAll{n: n, data: data, onesCount: onesCount}
}
func (b *bitsetFastFlipAll) FlipAll() {
	b.flip = !b.flip
	b.onesCount = b.n - b.onesCount
}
func (b *bitsetFastFlipAll) Add(i int32) bool {
	if b.data[i>>6]>>(i&63)&1 == 1 != b.flip {
		return false
	}
	b.data[i>>6] ^= 1 << (i & 63)
	b.onesCount++
	return true
}
func (b *bitsetFastFlipAll) Discard(i int32) bool {
	if b.data[i>>6]>>(i&63)&1 == 1 == b.flip {
		return false
	}
	b.data[i>>6] ^= 1 << (i & 63)
	b.onesCount--
	return true
}
func (b *bitsetFastFlipAll) Flip(i int32) {
	if b.data[i>>6]>>(i&63)&1 == 1 == b.flip {
		b.data[i>>6] ^= 1 << (i & 63)
		b.onesCount++
	} else {
		b.data[i>>6] ^= 1 << (i & 63)
		b.onesCount--
	}
}
func (b *bitsetFastFlipAll) Has(i int32) bool { return b.data[i>>6]>>(i&63)&1 == 1 != b.flip }
func (b *bitsetFastFlipAll) OnesCount() int32 { return b.onesCount }

// ------------------- VersionTree -------------------

type VersionTree struct {
	nodes   []*treeNode
	version int32
}

// 创建一个新的版本树，maxOperation 表示最大操作数.
// 版本号从 0 开始.
func NewVersionTree(maxOperation int32) *VersionTree {
	nodes := make([]*treeNode, 0, maxOperation+1)
	nodes = append(nodes, newTreeNode(emptyOperation))
	return &VersionTree{nodes: nodes}
}

// 在当前版本上添加一个操作，返回新版本号.
func (t *VersionTree) Apply(apply, undo func()) int32 {
	newNode := newTreeNode2(apply, undo)
	t.nodes = append(t.nodes, newNode)
	t.nodes[t.version].children = append(t.nodes[t.version].children, newNode)
	t.version = int32(len(t.nodes) - 1)
	return t.version
}

// 切换到指定版本.
func (t *VersionTree) SwitchVersion(version int32) { t.version = version }

// 应用所有操作.
func (t *VersionTree) Run() { t.dfs(t.nodes[0]) }

// 获取当前版本号.
func (t *VersionTree) Version() int32 { return t.version }

func (t *VersionTree) dfs(root *treeNode) {
	root.operation.apply()
	for _, child := range root.children {
		t.dfs(child)
	}
	root.operation.undo()
}

type treeNode struct {
	children  []*treeNode
	operation *operation
}

func newTreeNode(op *operation) *treeNode       { return &treeNode{operation: op} }
func newTreeNode2(apply, undo func()) *treeNode { return newTreeNode(newOperation(apply, undo)) }

type operation struct{ apply, undo func() }

var emptyOperation = newOperation(func() {}, func() {})

func newOperation(apply, undo func()) *operation { return &operation{apply: apply, undo: undo} }
