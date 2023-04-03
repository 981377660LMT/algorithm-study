// !扩容影响了答案
// !以下为答案错误的代码

package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"time"
)

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
		return a.(int) - b.(int)
	}, q*60) // !开50-60倍空间

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
			fmt.Fprintln(out, sl.BisectLeft(roots[i], x)+1)
		case 4:
			roots[i] = curRoot
			x--
			fmt.Fprintln(out, sl.At(roots[i], x))
		case 5:
			// 前驱
			roots[i] = curRoot
			fmt.Fprintln(out, sl.Lower(roots[i], x))
		case 6:
			roots[i] = curRoot
			fmt.Fprintln(out, sl.Upper(roots[i], x))
		}
	}
}

type Value = interface{}

type node struct {
	left, right int
	size        int
	priority    uint64
	value       Value
}

type SortedList struct {
	seed       uint64
	root       int
	comparator func(a, b Value) int
	nodes      []node
}

func NewSortedList(comparator func(a, b Value) int, capacity int) *SortedList {
	sl := &SortedList{
		seed:       uint64(time.Now().UnixNano()/2 + 1),
		comparator: comparator,
		nodes:      make([]node, 0, max(capacity, 16)), // !扩容的时候是不是拿不到正确的长度
	}
	dummy := &node{size: 0, priority: sl.nextRand()}
	sl.nodes = append(sl.nodes, *dummy)
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
			sl.splitByValue(sl.nodes[*x].right, value, &sl.nodes[*x].right, y, strictLess)
			sl.pushUp(*x)
		} else {
			*y = sl.copyNode(root)
			sl.splitByValue(sl.nodes[*y].left, value, x, &sl.nodes[*y].left, strictLess)
			sl.pushUp(*y)
		}
	} else {
		if sl.comparator(sl.nodes[root].value, value) <= 0 {
			*x = sl.copyNode(root)
			sl.splitByValue(sl.nodes[*x].right, value, &sl.nodes[*x].right, y, strictLess)
			sl.pushUp(*x)
		} else {
			*y = sl.copyNode(root)
			sl.splitByValue(sl.nodes[*y].left, value, x, &sl.nodes[*y].left, strictLess)
			sl.pushUp(*y)
		}
	}
}

func (sl *SortedList) merge(x, y int) int {
	if x == 0 || y == 0 {
		return x + y
	}

	if sl.nodes[x].priority > sl.nodes[y].priority {
		p := sl.copyNode(x)
		sl.nodes[p].right = sl.merge(sl.nodes[p].right, y)
		sl.pushUp(p)
		return p
	} else {
		p := sl.copyNode(y)
		sl.nodes[p].left = sl.merge(x, sl.nodes[p].left)
		sl.pushUp(p)
		return p
	}
}

func (sl *SortedList) newNode(value Value) int {
	sl.nodes = append(sl.nodes, node{
		value:    value,
		size:     1,
		priority: sl.nextRand(),
	})
	return len(sl.nodes) - 1
}

func (sl *SortedList) copyNode(root int) int {
	nodeCopy := sl.nodes[root]
	sl.nodes = append(sl.nodes, nodeCopy)
	return len(sl.nodes) - 1
}

func (sl *SortedList) pushUp(x int) {
	sl.nodes[x].size = sl.nodes[sl.nodes[x].left].size + sl.nodes[sl.nodes[x].right].size + 1
}

// 求小于等于 value 的最大值.不存在则返回 nil
func (sl *SortedList) Prev(value Value) Value {
	var x, y int
	sl.splitByValue(sl.root, value, &x, &y, false)
	if x == 0 {
		return nil
	}
	res := sl.nodes[sl.kthNode(x, sl.nodes[x].size)].value
	sl.root = sl.merge(x, y)
	return res
}

// 求大于等于 value 的最小值.不存在则返回 nil
func (sl *SortedList) Next(value Value) Value {
	var x, y int
	sl.splitByValue(sl.root, value, &x, &y, true)
	if y == 0 {
		return nil
	}
	res := sl.nodes[sl.kthNode(y, 1)].value
	sl.root = sl.merge(x, y)
	return res
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
	res := sl.nodes[x].size
	sl.root = sl.merge(x, y)
	return res
}

func (sl *SortedList) BisectRight(rootVersion int, value Value) int {
	var x, y int
	sl.splitByValue(rootVersion, value, &x, &y, false)
	res := sl.nodes[x].size
	sl.root = sl.merge(x, y)
	return res
}

func (sl *SortedList) Len(rootVersion int) int {
	return sl.nodes[rootVersion].size
}

func (sl *SortedList) At(rootVersion int, index int) Value {
	n := sl.Len(rootVersion)
	if index < 0 {
		index += n
	}
	if index < 0 || index >= n {
		panic(fmt.Sprintf("%d index out of range: [%d,%d]", index, 0, n-1))
	}
	index++
	return sl.nodes[sl.kthNode(rootVersion, index)].value
}

func (sl *SortedList) Lower(rootVersion, val int) Value {
	var x, y int
	var res Value
	sl.splitByValue(rootVersion, val, &x, &y, true)
	if x == 0 {
		res = nil
	} else {
		res = sl.nodes[sl.kthNode(x, sl.nodes[x].size)].value
	}
	// sl.root = sl.merge(x, y) // 赋值rt其实灭有必要
	return res
}

func (sl *SortedList) Upper(rootVersion, val int) Value {
	var x, y int
	var res Value
	sl.splitByValue(rootVersion, val, &x, &y, false)
	if y == 0 {
		res = math.MaxInt32
	} else {
		res = sl.nodes[sl.kthNode(y, 1)].value
	}
	// sl.root = sl.merge(x, y) // 赋值rt其实灭有必要
	return res
}

func (sl *SortedList) kthNode(root int, k int) (nodeId int) {
	left := sl.nodes[root].left
	leftSize := sl.nodes[left].size
	if leftSize+1 > k {
		nodeId = sl.kthNode(left, k)
		return
	}
	if leftSize+1 == k {
		nodeId = root
		return
	}
	nodeId = sl.kthNode(sl.nodes[root].right, k-leftSize-1)
	return
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
