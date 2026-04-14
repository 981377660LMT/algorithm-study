// G - Copy Query
// https://atcoder.jp/contests/abc453/tasks/abc453_g
//
// 有 $N$ 个长度为 $M$ 的数列 $A_1, A_2, \dots, A_N$。初始所有元素均为 $0$。
// 你需要按顺序处理 $Q$ 个操作：
//
// 1.  **复制 (Type 1)**：`1 X Y`
//     将数列 $A_Y$ 的内容完全复制给 $A_X$（即 $A_X = A_Y$）。
// 2.  **单点更新 (Type 2)**：`2 X Y Z`
//     将第 $X$ 个数列的第 $Y$ 个元素修改为 $Z$（即 $A_{X,Y} = Z$）。
// 3.  **区间查询 (Type 3)**：`3 X L R`
//     计算并输出第 $X$ 个数列在区间 $[L, R]$ 内的元素之和。
//
// 数据范围
// - $1 \le N, M, Q \le 2 \times 10^5$
// - $0 \le Z \le 10^9$
// - $1 \le L \le R \le M$
// - 时间限制：2s，内存限制：1024MiB

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m, q int
	fmt.Fscan(in, &n, &m, &q)

	vt := NewVersionTree(int32(q))
	bit := NewBitArray(m)
	actualValues := make([]int, m)

	// heads[i] 表示第 i 个分支当前指向的 commit.
	heads := make([]int32, n)
	for i := range heads {
		heads[i] = 0
	}

	res := make([]int, q)
	for i := range res {
		res[i] = -1
	}

	for i := 0; i < q; i++ {
		var t int
		fmt.Fscan(in, &t)
		switch t {
		case 1:
			var x, y int
			fmt.Fscan(in, &x, &y)
			x, y = x-1, y-1
			heads[x] = heads[y]
		case 2:
			var x, y, z int
			fmt.Fscan(in, &x, &y, &z)
			x, y = x-1, y-1

			targetPos, targetVal := y, z
			var oldVal int
			heads[x] = vt.AddStep(heads[x], func() bool {
				oldVal = actualValues[targetPos]
				bit.Add(targetPos, targetVal-oldVal)
				actualValues[targetPos] = targetVal
				return true
			}, func() {
				bit.Add(targetPos, oldVal-targetVal)
				actualValues[targetPos] = oldVal
			})

		case 3:
			var x, l, r int
			fmt.Fscan(in, &x, &l, &r)
			x--
			l--
			queryIdx := i
			vt.AddQuery(heads[x], func() {
				res[queryIdx] = bit.QueryRange(l, r)
			})
		}
	}

	vt.Commit()

	for _, v := range res {
		if v != -1 {
			fmt.Fprintln(out, v)
		}
	}
}

type VersionTree struct {
	nodes []*treeNode
}

// 创建一个新的版本树，maxStepCount 表示最大操作数.
// !初始时版本号为0(没有任何修改)，第一次操作后版本号为1，以此类推.
func NewVersionTree(maxStepCount int32) *VersionTree {
	nodes := make([]*treeNode, 0, maxStepCount+1)
	nodes = append(nodes, newTreeNode(emptyStep))
	return &VersionTree{nodes: nodes}
}

// 从 parent 这个版本派生一个新版本.
// AddMutation.
//
//	apply: 变更操作，返回是否成功.
//	invert: 变更的逆操作.如果操作不可逆，则需要拷贝变更前的状态以实现撤销.
func (t *VersionTree) AddStep(parent int32, apply func() bool, invert func()) (newVersion int32) {
	newNode := newTreeNode2(apply, invert)
	t.nodes = append(t.nodes, newNode)
	t.nodes[parent].children = append(t.nodes[parent].children, newNode)
	newVersion = int32(len(t.nodes) - 1)
	return
}

func (t *VersionTree) AddQuery(version int32, query func()) {
	t.nodes[version].queries = append(t.nodes[version].queries, query)
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

// !Point Add Range Sum, 0-based.
type BITArray struct {
	n     int
	total int
	data  []int
}

func NewBitArray(n int) *BITArray {
	res := &BITArray{n: n, data: make([]int, n)}
	return res
}

func (b *BITArray) Add(index int, v int) {
	b.total += v
	for index++; index <= b.n; index += index & -index {
		b.data[index-1] += v
	}
}

// [0, end).
func (b *BITArray) QueryPrefix(end int) int {
	if end > b.n {
		end = b.n
	}
	res := 0
	for ; end > 0; end -= end & -end {
		res += b.data[end-1]
	}
	return res
}

// [start, end).
func (b *BITArray) QueryRange(start, end int) int {
	if start < 0 {
		start = 0
	}
	if end > b.n {
		end = b.n
	}
	if start >= end {
		return 0
	}
	if start == 0 {
		return b.QueryPrefix(end)
	}
	pos, neg := 0, 0
	for end > start {
		pos += b.data[end-1]
		end &= end - 1
	}
	for start > end {
		neg += b.data[start-1]
		start &= start - 1
	}
	return pos - neg
}
