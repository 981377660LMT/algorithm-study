// https://www.luogu.com.cn/problem/solution/P3835
// 可持久化平衡树
// TODO
// !FIXME

package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"time"
)

// 您需要写一种数据结构（可参考题目标题），来维护一个可重整数集合，
// 其中需要提供以下操作（ 对于各个以往的历史版本 ）：
// 1、 插入 x
// 2、 删除 x（若有多个相同的数，应只删除一个，如果没有请忽略该操作）
// 3、 查询 x 的排名（排名定义为比当前数小的数的个数 +1）
// 4、查询排名为 x 的数
// 5、 求 x 的前驱（前驱定义为小于 x，且最大的数，如不存在输出 −2^31+1 ）
// 6、求 x 的后继（后继定义为大于 x，且最小的数，如不存在输出 2^31 −1 ）
// 和原本平衡树不同的一点是，每一次的任何操作都是基于某一个历史版本，
// 同时生成一个新的版本。（操作3, 4, 5, 6即保持原版本无变化）

// 每个版本的编号即为操作的序号（版本0即为初始状态，空树）
// https://www.cnblogs.com/dx123/p/16584604.html

// !哪里有问题
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var q int
	fmt.Fscan(in, &q)

	roots := make([]int, q+10)
	// !保存每个版本的根节点 初始为0
	// !每一次的任何操作都是基于某一个历史版本，同时生成一个新的版本。（操作3, 4, 5, 6即保持原版本无变化）
	sl := NewSortedList(func(a, b Value) int {
		return a - b
	}, q*60) // !开60倍空间

	for i := 1; i <= q; i++ {
		var version, op, x int
		fmt.Fscan(in, &version, &op, &x)
		curRoot := roots[version]
		switch op {
		case 1:
			sl.Add(curRoot, x)
			roots[i] = sl.Add(curRoot, x)
		case 2:
			sl.Discard(curRoot, x)
			roots[i] = sl.Discard(curRoot, x)
		case 3:
			roots[i] = curRoot
			fmt.Fprintln(out, sl.BisectLeft(roots[i], x))
		case 4:
			roots[i] = curRoot
			fmt.Fprintln(out, sl.At(roots[i], x))
		case 5:
			// 前驱
			roots[i] = curRoot
			fmt.Fprintln(out, sl.getPre(roots[i], x))
		case 6:
			roots[i] = curRoot
			fmt.Fprintln(out, sl.getNxt(roots[i], x))
		}
	}

}

type Value = int

// type Value = interface{}

type node struct {
	left, right int
	size        int
	priority    uint64
	value       Value
}

type SortedList struct {
	seed       uint64
	root       int
	nodeId     int
	comparator func(a, b Value) int
	nodes      []node
}

func NewSortedList(comparator func(a, b Value) int, initCapacity int) *SortedList {
	sl := &SortedList{
		seed:       uint64(time.Now().UnixNano()/2 + 1),
		comparator: comparator,
		nodes:      make([]node, max(initCapacity, 128)), // !定长开
	}
	return sl
}

func (sl *SortedList) splitByValue(root int, value Value, x, y *int, strictLess bool) {
	if root == 0 {
		*x, *y = 0, 0
		return
	}

	if strictLess {
		if sl.comparator(sl.nodes[root].value, value) < 0 {
			*x = sl.copyNode(root)
			sl.nodes[*x] = sl.nodes[root]
			sl.splitByValue(sl.nodes[*x].right, value, &sl.nodes[*x].right, y, strictLess)
			sl.pushUp(*x)
		} else {
			*y = sl.copyNode(root)
			sl.nodes[*y] = sl.nodes[root]
			sl.splitByValue(sl.nodes[*y].left, value, x, &sl.nodes[*y].left, strictLess)
			sl.pushUp(*y)
		}
	} else {
		if sl.comparator(sl.nodes[root].value, value) <= 0 {
			*x = sl.copyNode(root)
			sl.nodes[*x] = sl.nodes[root]
			sl.splitByValue(sl.nodes[*x].right, value, &sl.nodes[*x].right, y, strictLess)
			sl.pushUp(*x)
		} else {
			*y = sl.copyNode(root)
			sl.nodes[*y] = sl.nodes[root]
			sl.splitByValue(sl.nodes[*y].left, value, x, &sl.nodes[*y].left, strictLess)
			sl.pushUp(*y)
		}
	}

}

// 因为split已经复制过了，所以这里不需要复制
func (sl *SortedList) merge(x, y int) int {
	if x == 0 || y == 0 {
		return x + y
	}

	if sl.nodes[x].priority > sl.nodes[y].priority {
		sl.nodes[x].right = sl.merge(sl.nodes[x].right, y)
		sl.pushUp(x)
		return x
	} else {
		sl.nodes[y].left = sl.merge(x, sl.nodes[y].left)
		sl.pushUp(y)
		return y
	}
}

func (sl *SortedList) newNode(value Value) int {
	sl.nodeId++
	sl.nodes[sl.nodeId].size = 1
	sl.nodes[sl.nodeId].value = value
	sl.nodes[sl.nodeId].priority = sl.nextRand()
	sl.nodes[sl.nodeId].left = 0
	sl.nodes[sl.nodeId].right = 0
	return sl.nodeId
}

func (sl *SortedList) copyNode(root int) int {
	sl.nodeId++
	sl.nodes[sl.nodeId] = sl.nodes[root]
	return sl.nodeId
}

func (sl *SortedList) pushUp(x int) {
	sl.nodes[x].size = sl.nodes[sl.nodes[x].left].size + sl.nodes[sl.nodes[x].right].size + 1
}

func (sl *SortedList) Add(rootVersion int, value Value) int {
	var x, y int
	sl.splitByValue(rootVersion, value, &x, &y, false)
	sl.root = sl.merge(sl.merge(x, sl.newNode(value)), y)
	return sl.root
}

func (sl *SortedList) Discard(rootVersion int, value Value) int {
	var x, y, z int
	sl.splitByValue(rootVersion, value, &x, &y, true)
	sl.splitByValue(y, value, &y, &z, false)
	y = sl.merge(sl.nodes[y].left, sl.nodes[y].right)
	sl.root = sl.merge(sl.merge(x, y), z)
	return sl.root
}

func (sl *SortedList) BisectLeft(rootVersion int, value Value) int {
	var x, y int
	sl.splitByValue(rootVersion, value, &x, &y, true)
	res := sl.nodes[x].size + 1
	sl.root = sl.merge(x, y)
	return res
}

func (sl *SortedList) BisectRight(rootVersion int, value Value) int {
	var x, y int
	sl.splitByValue(rootVersion, value, &x, &y, false)
	res := sl.nodes[x].size + 1
	sl.root = sl.merge(x, y)
	return res
}

func (sl *SortedList) Len(rootVersion int) int {
	return sl.nodes[rootVersion].size
}

// at k
func (sl *SortedList) At(root int, k int) int {
	left := sl.nodes[root].left
	leftSize := sl.nodes[left].size
	if leftSize+1 > k {
		return sl.At(left, k)
	}
	if leftSize+1 == k {
		return sl.nodes[root].value
	}
	return sl.At(sl.nodes[root].right, k-leftSize-1)
}

func (sl *SortedList) getPre(now, val int) int {
	var x, y, res int
	sl.splitByValue(now, val, &x, &y, true)
	if x == 0 {
		res = -math.MaxInt32
	} else {
		res = sl.At(x, sl.nodes[x].size)
	}
	// sl.root = sl.merge(x, y) // 赋值rt其实灭有必要
	return res
}

func (sl *SortedList) getNxt(now, val int) int {
	var x, y, res int
	sl.splitByValue(now, val, &x, &y, false)
	if y == 0 {
		res = math.MaxInt32
	} else {
		res = sl.At(y, 1)
	}
	// sl.root = sl.merge(x, y) // 赋值rt其实灭有必要
	return res
}

// Return all elements in index order.
func (sl *SortedList) InOrder(rootVersion int) []Value {
	res := make([]Value, 0, sl.Len(rootVersion))
	sl.inOrder(rootVersion, &res)
	return res
}

func (sl *SortedList) inOrder(root int, res *[]Value) {
	if root == 0 {
		return
	}
	sl.inOrder(sl.nodes[root].left, res)
	*res = append(*res, sl.nodes[root].value)
	sl.inOrder(sl.nodes[root].right, res)
}

// https://nyaannyaan.github.io/library/misc/rng.hpp
func (sl *SortedList) nextRand() uint64 {
	sl.seed ^= sl.seed << 7
	sl.seed ^= sl.seed >> 9
	return sl.seed
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
