// 版本树/操作树 OperationTree
//
// 理解：
//  版本树（VersionTree）或操作树（OperationTree）是一种数据结构，用于`离线`管理和跟踪数据的版本历史.
//  !1. 在这里，版本指的是数据结构在特定时间点的状态.
//  !2. 每次变更都会生成一个新的版本，版本是immutable的.
//      这种方式允许我们快速回溯到历史状态，实现撤销（Undo）和重做（Redo）操作.
//  !3. 每个版本内可以有多个查询，但是只有一个变更(原子操作).
//
// Api:
//  NewVersionTree(maxMutation int32) *VersionTree
//  AddStep(apply func() bool, invert func()) int32
//  AddSwitchVersionStep(version int32) int32
//  SwitchVersion(version int32)
//  Commit()

package main

import (
	"bufio"
	"fmt"
	"os"
)

// Persistent Bookcase
// https://www.luogu.com.cn/problem/CF707D
// 1 i j : 将(i,j)置为1
// 2 i j : 将(i,j)置为0
// 3 i j : 将i行01反转
// 4 k : 回到第k次操作后的状态(0 <= k < q)
// 每次操作后，输出全局 1 的个数
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

	history := make([]int32, q) // !每次操作后全局 1 的个数
	onesCountAll := int32(0)
	add := func(i, j int32) bool {
		if sets[i].Add(j) {
			onesCountAll++
			return true
		}
		return false
	}
	discard := func(i, j int32) bool {
		if sets[i].Discard(j) {
			onesCountAll--
			return true
		}
		return false
	}
	flipAll := func(i int32) bool {
		onesCountAll -= sets[i].OnesCount()
		sets[i].FlipAll()
		onesCountAll += sets[i].OnesCount()
		return true
	}

	tree := NewVersionTree(q)
	for qi := int32(0); qi < q; qi++ {
		var op uint8
		fmt.Fscan(in, &op)

		switch op {
		case 1:
			var i, j int32
			fmt.Fscan(in, &i, &j)
			i--
			j--
			tree.AddStep(
				func() bool {
					return add(i, j)
				},
				func() {
					discard(i, j)
				},
			)
		case 2:
			var i, j int32
			fmt.Fscan(in, &i, &j)
			i--
			j--
			tree.AddStep(
				func() bool {
					return discard(i, j)
				},
				func() {
					add(i, j)
				},
			)
		case 3:
			var i int32
			fmt.Fscan(in, &i)
			i--
			tree.AddStep(
				func() bool {
					return flipAll(i)
				},
				func() {
					flipAll(i)
				},
			)
		case 4:
			var k int32
			fmt.Fscan(in, &k)
			tree.AddSwitchVersionStep(k)
		}

		tree.AddQuery(
			func(kth, _ int32) {
				history[kth] = onesCountAll
			},
		)
	}

	tree.Commit()

	for _, v := range history {
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
	nodes      []*treeNode
	version    int32
	queryCount int32
}

// 创建一个新的版本树，maxStepCount 表示最大操作数.
// !初始时版本号为0(没有任何修改)，第一次操作后版本号为1，以此类推.
func NewVersionTree(maxStepCount int32) *VersionTree {
	nodes := make([]*treeNode, 0, maxStepCount+1)
	nodes = append(nodes, newTreeNode(emptyStep))
	return &VersionTree{nodes: nodes}
}

// 在当前版本上添加(apply)一个变更，返回新版本号.
// AddMutation.
//
//	apply: 变更操作，返回是否成功.
//	invert: 变更的逆操作.如果操作不可逆，则需要拷贝变更前的状态以实现撤销.
func (t *VersionTree) AddStep(apply func() bool, invert func()) (newVersion int32) {
	newNode := newTreeNode2(apply, invert)
	t.nodes = append(t.nodes, newNode)
	t.nodes[t.version].children = append(t.nodes[t.version].children, newNode)
	newVersion = int32(len(t.nodes) - 1)
	t.version = newVersion
	return
}

// !在当前版本上添加一个切换版本的操作，视为一次修改操作.
func (t *VersionTree) AddSwitchVersionStep(version int32) (newVersion int32) {
	newNode := newTreeNode(emptyStep)
	t.nodes = append(t.nodes, newNode)
	t.nodes[version].children = append(t.nodes[version].children, newNode)
	newVersion = int32(len(t.nodes) - 1)
	t.version = version
	return
}

// !切换到指定版本，不视为一次修改操作.
func (t *VersionTree) SwitchVersion(version int32) {
	t.version = version
}

// 在当前版本上添加一个查询.
//
//	kth: 第 k 次查询(0-based).
//	version: 查询时的版本号.
func (t *VersionTree) AddQuery(query func(kth, version int32)) {
	q, v := t.queryCount, t.version
	t.nodes[t.version].queries = append(t.nodes[t.version].queries, func() { query(q, v) })
	t.queryCount++
}

// 应用所有操作.
func (t *VersionTree) Commit() { t.dfs(t.nodes[0]) }

func (t *VersionTree) dfs(root *treeNode) {
	ok := root.step.apply()
	for _, query := range root.queries {
		query()
	}
	for _, child := range root.children {
		t.dfs(child)
	}
	if ok {
		root.step.invert()
	}
}

type treeNode struct {
	queries  []func()
	step     *step
	children []*treeNode
}

func newTreeNode(step *step) *treeNode { return &treeNode{step: step} }
func newTreeNode2(apply func() bool, invert func()) *treeNode {
	return newTreeNode(newStep(apply, invert))
}

type step struct {
	apply  func() bool // stepResult
	invert func()
}

var emptyStep = newStep(func() bool { return false }, func() {})

func newStep(apply func() bool, invert func()) *step { return &step{apply: apply, invert: invert} }
