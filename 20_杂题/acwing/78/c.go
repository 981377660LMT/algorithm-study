// # n 个小朋友排成一排，从左到右依次编号为 1∼n。

// # 第 i 个小朋友的身高为 hi。

// # 虽然队伍已经排好，但是小朋友们对此并不完全满意。

// # 对于一个小朋友来说，如果存在其他小朋友身高比他更矮，却站在他右侧的情况，该小朋友就会感到不满。

// # 每个小朋友的不满程度都可以量化计算，具体来说，对于第 i 个小朋友：

// # 如果存在比他更矮且在他右侧的小朋友，那么他的不满值等于其中最靠右的那个小朋友与他之间的小朋友数量。
// # 如果不存在比他更矮且在他右侧的小朋友，那么他的不满值为 −1。
// # 请你计算并输出每个小朋友的不满值。

// # 注意，第 1 个小朋友和第 2 个小朋友之间的小朋友数量为 0，第 1 个小朋友和第 4 个小朋友之间的小朋友数量为 2。
package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

func main() {
	const INF int = int(1e18)
	const MOD int = 998244353

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	sl := NewSortedList(func(a, b Value) int {
		return a - b
	}, n)

	type Pair struct{ value, index int }
	pairs := make([]Pair, n)
	for i := 0; i < n; i++ {
		pairs[i].value = nums[i]
		pairs[i].index = i
	}

	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].value < pairs[j].value
	})

	res := make([]int, n)
	for i := 0; i < n; i++ {
		res[i] = -1
	}

	groupByValue := make(map[int][]Pair)
	for i := 0; i < n; i++ {
		groupByValue[pairs[i].value] = append(groupByValue[pairs[i].value], pairs[i])
	}

	// sortedKeys
	sortedKeys := make([]int, 0, len(groupByValue))
	for k := range groupByValue {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Ints(sortedKeys)

	for _, k := range sortedKeys {
		groups := groupByValue[k]
		for _, p := range groups {
			index := p.index
			if sl.Len() > 0 {
				last := sl.At(sl.Len() - 1)
				if index > last {
					res[index] = -1
				} else {
					res[index] = last - index - 1
				}
			}
		}

		for _, p := range groups {
			index := p.index
			sl.Add(index)
		}
	}

	for i := 0; i < n; i++ {
		fmt.Fprint(out, res[i], " ")
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
	node := &node{
		value:    value,
		size:     1,
		priority: sl.nextRand(),
	}
	sl.nodes = append(sl.nodes, *node)
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
