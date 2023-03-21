// MK平均值
// 滑动窗口topK之和 (这里topK是最小值)

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

// 每个滑动窗口topK之和(这里topK是最小值)
func windowTopkSum(nums []int, windowSize, k int) []int {
	n := len(nums)
	sl := NewSortedList(func(a, b int) int {
		return a - b
	}, n)

	res, topKSum := make([]int, 0, n-windowSize+1), 0

	add := func(x int) {
		pos := sl.BisectLeft(x)
		if pos < k {
			topKSum += x
			if k-1 < sl.Len() {
				topKSum -= sl.At(k - 1)
			}
		}
		sl.Add(x)
	}

	remove := func(x int) {
		pos := sl.BisectLeft(x)
		if pos < k {
			topKSum -= x
			if k < sl.Len() {
				topKSum += sl.At(k)
			}
		}
		sl.Discard(x)
	}

	query := func() int {
		return topKSum
	}

	for i := 0; i < n; i++ {
		add(nums[i])
		if i >= windowSize {
			remove(nums[i-windowSize])
		}
		if i >= windowSize-1 {
			res = append(res, query())
		}
	}

	return res
}

func main() {

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var n, m, k int
	fmt.Fscan(in, &n, &m, &k)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	res := windowTopkSum(nums, m, k)
	fmt.Fprint(out, strings.Trim(fmt.Sprint(res), "[]"))

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
		nodes:      make([]node, 0, max(initCapacity, 16)),
	}

	dummy := &node{size: 0, priority: sl.nextRand()} // dummy node 0
	sl.nodes = append(sl.nodes, *dummy)
	return sl
}

func (sl *SortedList) Build(nums []Value) int {
	n := len(nums)
	keys := make([]int, 0, n)
	for i := 0; i < n; i++ {
		keys = append(keys, sl.newNode(nums[i]))
	}

	// 将元素按照键值 key 排序，然后一个一个插入到当前的笛卡尔树中
	sort.Slice(keys, func(i, j int) bool {
		return sl.comparator(sl.nodes[keys[i]].value, sl.nodes[keys[j]].value) < 0
	})

	stack := []int{}
	pre := make([]int, n)
	for i := 0; i < n; i++ {
		pre[i] = -1
	}

	for i := 0; i < n; i++ {
		last := -1
		for len(stack) > 0 && sl.nodes[stack[len(stack)-1]].priority > sl.nodes[keys[i]].priority {
			last = stack[len(stack)-1]
			stack = stack[:len(stack)-1]
		}

		if len(stack) > 0 {
			pre[i] = stack[len(stack)-1]
		}
		if last != -1 {
			pre[last] = i
		}

		stack = append(stack, i)
	}

	root := -1
	for i := 0; i < n; i++ {
		if pre[i] != -1 {
			if i < pre[i] {
				sl.nodes[keys[pre[i]]].left = keys[i]
			} else {
				sl.nodes[keys[pre[i]]].right = keys[i]
			}
		} else {
			root = i
		}
	}

	sl.root = keys[root]
	sl.build(sl.root)
	return sl.root
}

func (sl *SortedList) build(root int) {
	nodeRef := &sl.nodes[root]
	if nodeRef.left != 0 {
		sl.build(nodeRef.left)
	}
	if nodeRef.right != 0 {
		sl.build(nodeRef.right)
	}
	sl.pushUp(root)
}

func (sl *SortedList) pushUp(root int) {
	sl.nodes[root].size = sl.nodes[sl.nodes[root].left].size + sl.nodes[sl.nodes[root].right].size + 1
}

func (sl *SortedList) Add(value Value) {
	var x, y, z int
	sl.splitByValue(sl.root, value, &x, &y, false)
	z = sl.newNode(value)
	sl.root = sl.merge(sl.merge(x, z), y)
}

func (sl *SortedList) At(index int) Value {
	n := sl.Len()
	if index < 0 {
		index += n
	}
	if index < 0 || index >= n {
		panic(fmt.Sprintf("%d index out of range: [%d,%d]", index, 0, n-1))
	}
	return sl.nodes[sl.kthNode(sl.root, index+1)].value
}

func (sl *SortedList) Pop(index int) Value {
	n := sl.Len()
	if index < 0 {
		index += n
	}

	index += 1 // dummy offset
	var x, y, z int
	sl.splitByRank(sl.root, index, &y, &z)
	sl.splitByRank(y, index-1, &x, &y)
	res := sl.nodes[y].value
	sl.root = sl.merge(x, z)
	return res
}

func (sl *SortedList) Discard(value Value) {
	var x, y, z int
	sl.splitByValue(sl.root, value, &x, &z, false)
	sl.splitByValue(x, value, &x, &y, true)
	y = sl.merge(sl.nodes[y].left, sl.nodes[y].right)
	sl.root = sl.merge(sl.merge(x, y), z)
}

// Remove [start, stop) from list.
func (sl *SortedList) Erase(start, stop int) {
	var x, y, z int
	start++ // dummy offset
	sl.splitByRank(sl.root, stop, &y, &z)
	sl.splitByRank(y, start-1, &x, &y)
	sl.root = sl.merge(x, z)
}

func (sl *SortedList) BisectLeft(value Value) int {
	var x, y int
	sl.splitByValue(sl.root, value, &x, &y, true)
	res := sl.nodes[x].size
	sl.root = sl.merge(x, y)
	return res
}

func (sl *SortedList) BisectRight(value Value) int {
	var x, y int
	sl.splitByValue(sl.root, value, &x, &y, false)
	res := sl.nodes[x].size
	sl.root = sl.merge(x, y)
	return res
}

func (sl *SortedList) String() string {
	sb := []string{"SortedList{"}
	values := []string{}
	for i := 0; i < sl.Len(); i++ {
		values = append(values, fmt.Sprintf("%v", sl.At(i)))
	}
	sb = append(sb, strings.Join(values, ","), "}")
	return strings.Join(sb, "")
}

func (sl *SortedList) Len() int {
	return sl.nodes[sl.root].size
}

func (sl *SortedList) kthNode(root int, k int) int {
	cur := root
	for cur != 0 {
		if sl.nodes[sl.nodes[cur].left].size+1 == k {
			break
		} else if sl.nodes[sl.nodes[cur].left].size >= k {
			cur = sl.nodes[cur].left
		} else {
			k -= sl.nodes[sl.nodes[cur].left].size + 1
			cur = sl.nodes[cur].right
		}
	}
	return cur
}

func (sl *SortedList) splitByValue(root int, value Value, x, y *int, strictLess bool) {
	if root == 0 {
		*x, *y = 0, 0
		return
	}

	if strictLess {
		if sl.comparator(sl.nodes[root].value, value) < 0 {
			*x = root
			sl.splitByValue(sl.nodes[root].right, value, &sl.nodes[root].right, y, strictLess)
		} else {
			*y = root
			sl.splitByValue(sl.nodes[root].left, value, x, &sl.nodes[root].left, strictLess)
		}
	} else {
		if sl.comparator(sl.nodes[root].value, value) <= 0 {
			*x = root
			sl.splitByValue(sl.nodes[root].right, value, &sl.nodes[root].right, y, strictLess)
		} else {
			*y = root
			sl.splitByValue(sl.nodes[root].left, value, x, &sl.nodes[root].left, strictLess)
		}
	}

	sl.pushUp(root)
}

// Split by rank.
// Split the tree rooted at root into two trees, x and y, such that the size of x is k.
// x is the left subtree, y is the right subtree.
func (sl *SortedList) splitByRank(root, k int, x, y *int) {
	if root == 0 {
		*x, *y = 0, 0
		return
	}

	if k <= sl.nodes[sl.nodes[root].left].size {
		*y = root
		sl.splitByRank(sl.nodes[root].left, k, x, &sl.nodes[root].left)
		sl.pushUp(*y)
	} else {
		*x = root
		sl.splitByRank(sl.nodes[root].right, k-sl.nodes[sl.nodes[root].left].size-1, &sl.nodes[root].right, y)
		sl.pushUp(*x)
	}
}

func (sl *SortedList) merge(x, y int) int {
	if x == 0 || y == 0 {
		return x + y
	}

	if sl.nodes[x].priority < sl.nodes[y].priority {
		sl.nodes[x].right = sl.merge(sl.nodes[x].right, y)
		sl.pushUp(x)
		return x
	}
	sl.nodes[y].left = sl.merge(x, sl.nodes[y].left)
	sl.pushUp(y)
	return y
}

// Return all elements in index order.
func (sl *SortedList) InOrder() []Value {
	res := make([]Value, 0, sl.Len())
	sl.inOrder(sl.root, &res)
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

func (sl *SortedList) newNode(value Value) int {
	sl.nodes = append(sl.nodes, node{
		value:    value,
		size:     1,
		priority: sl.nextRand(),
	})
	return len(sl.nodes) - 1
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
