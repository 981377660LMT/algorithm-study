// !扩容影响了答案
// !以下为答案错误的代码

package main

import (
	"fmt"
	"math"
	"reflect"
	"time"
)

func main() {
	for i := 0; i < 1000000; i++ {
		res := run()
		if !reflect.DeepEqual(res, []int{9, 1, 2, 10, 3}) {
			fmt.Println("not ok", res)
		}
	}
	fmt.Println("ok")
}

func run() []int {
	q := 10
	roots := make([]int, q+10)
	Q := [10][3]int{
		{0, 1, 9},
		{1, 1, 3},
		{1, 1, 10},
		{2, 4, 2},
		{3, 3, 9},
		{3, 1, 2},
		{6, 4, 1},
		{6, 2, 9},
		{8, 6, 3},
		{4, 5, 8},
	}

	sl := NewSortedList(func(a, b Value) int {
		return a - b
	}, 10)

	res := []int{}
	for i := 1; i <= q; i++ {
		version, op, x := Q[i-1][0], Q[i-1][1], Q[i-1][2]
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
			res = append(res, sl.BisectLeft(roots[i], x))
		case 4:
			roots[i] = curRoot
			res = append(res, sl.kthNode(roots[i], x))
		case 5:
			// 前驱
			roots[i] = curRoot
			res = append(res, sl.getPre(roots[i], x))
		case 6:
			roots[i] = curRoot
			res = append(res, sl.getNxt(roots[i], x))
		}
	}

	return res
}

type Value = int

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

func NewSortedList(comparator func(a, b Value) int, initCapacity int) *SortedList {
	sl := &SortedList{
		seed:       uint64(time.Now().UnixNano()/2 + 1),
		comparator: comparator,
		nodes:      make([]node, 0, max(initCapacity, 16)), // !扩容的时候是不是拿不到正确的长度
	}
	dummy := &node{size: 0, priority: sl.nextRand()}
	sl.nodes = append(sl.nodes, *dummy)
	return sl
}

func (sl *SortedList) splitByValue(root int, value Value, x, y *int) {
	if root == 0 {
		*x, *y = 0, 0
		return
	}

	if sl.comparator(sl.nodes[root].value, value) <= 0 {
		*x = sl.copyNode(root)
		sl.nodes[*x] = sl.nodes[root]
		sl.splitByValue(sl.nodes[*x].right, value, &sl.nodes[*x].right, y)
		sl.pushUp(*x)
	} else {
		*y = sl.copyNode(root)
		sl.nodes[*y] = sl.nodes[root]
		sl.splitByValue(sl.nodes[*y].left, value, x, &sl.nodes[*y].left)
		sl.pushUp(*y)
	}

}

func (sl *SortedList) merge(x, y int) int {
	if x == 0 || y == 0 {
		return x + y
	}

	if sl.nodes[x].priority > sl.nodes[y].priority {
		p := sl.copyNode(x)
		sl.nodes[p] = sl.nodes[x]
		sl.nodes[p].right = sl.merge(sl.nodes[p].right, y)
		sl.pushUp(p)
		return p
	} else {
		p := sl.copyNode(y)
		sl.nodes[p] = sl.nodes[y]
		sl.nodes[p].left = sl.merge(x, sl.nodes[p].left)
		sl.pushUp(p)
		return p
	}
}

func (sl *SortedList) newNode(value Value) int {
	newNode := &node{
		size:     1,
		priority: sl.nextRand(),
		value:    value,
	}
	sl.nodes = append(sl.nodes, *newNode)
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

func (sl *SortedList) Add(rootVersion int, value Value) int {
	var x, y int
	sl.splitByValue(rootVersion, value, &x, &y)
	sl.root = sl.merge(sl.merge(x, sl.newNode(value)), y)
	return sl.root
}

func (sl *SortedList) Discard(rootVersion int, value Value) int {
	var x, y, z int
	sl.splitByValue(rootVersion, value-1, &x, &y)
	sl.splitByValue(y, value, &y, &z)
	y = sl.merge(sl.nodes[y].left, sl.nodes[y].right)
	sl.root = sl.merge(sl.merge(x, y), z)
	return sl.root
}

func (sl *SortedList) BisectLeft(rootVersion int, value Value) int {
	var x, y int
	sl.splitByValue(rootVersion, value-1, &x, &y)
	res := sl.nodes[x].value + 1
	sl.root = sl.merge(x, y)
	return res
}

func (sl *SortedList) BisectRight(rootVersion int, value Value) int {
	var x, y int
	sl.splitByValue(rootVersion, value, &x, &y)
	res := sl.nodes[x].value + 1 // !value 还是 size
	sl.root = sl.merge(x, y)
	return res
}

func (sl *SortedList) Len(rootVersion int) int {
	return sl.nodes[rootVersion].size
}

func (sl *SortedList) kthNode(root int, k int) int {
	left := sl.nodes[root].left
	leftSize := sl.nodes[left].size
	if leftSize+1 > k {
		return sl.kthNode(left, k)
	}
	if leftSize+1 == k {
		return sl.nodes[root].value
	}
	return sl.kthNode(sl.nodes[root].right, k-leftSize-1)
}

func (sl *SortedList) getPre(now, val int) int {
	var x, y, res int
	sl.splitByValue(now, val-1, &x, &y)
	if x == 0 {
		res = -math.MaxInt32
	} else {
		res = sl.kthNode(x, sl.nodes[x].size)
	}
	sl.root = sl.merge(x, y)
	return res
}

func (sl *SortedList) getNxt(now, val int) int {
	var x, y, res int
	sl.splitByValue(now, val, &x, &y)
	if y == 0 {
		res = math.MaxInt32
	} else {
		res = sl.kthNode(y, 1)
	}
	sl.root = sl.merge(x, y)
	return res
}

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
