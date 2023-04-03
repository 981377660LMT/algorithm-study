// 寻找前驱后继/区间删除
package main

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

const INF int = 1e18

// 2612. 最少翻转操作数
// https://leetcode.cn/problems/minimum-reverse-operations/
func minReverseOperations(n int, p int, banned []int, k int) []int {
	finder := [2]*Finder{
		NewFinder(func(a, b int) int { return a - b }, n/2),
		NewFinder(func(a, b int) int { return a - b }, n/2),
	}

	for i := 0; i < n; i++ {
		finder[i&1].Insert(i)
	}
	for _, i := range banned {
		finder[i&1].Erase(i)
	}

	getRange := func(i int) (int, int) {
		return max(i-k+1, k-i-1), min(i+k-1, 2*n-k-i-1)
	}
	setUsed := func(u int) {
		finder[u&1].Erase(u)
	}

	findUnused := func(u int) int {
		left, right := getRange(u)
		pre, ok := finder[(u+k+1)&1].Prev(right)
		if ok && left <= pre && pre <= right {
			return pre
		}
		next, ok := finder[(u+k+1)&1].Next(left)
		if ok && left <= next && next <= right {
			return next
		}
		return -1
	}

	dist := OnlineBfs(n, p, setUsed, findUnused)
	res := make([]int, n)
	for i, d := range dist {
		if d == INF {
			res[i] = -1
		} else {
			res[i] = d
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

type Finder struct {
	seed       uint64
	root       int
	comparator func(a, b Value) int
	nodes      []node
}

func NewFinder(comparator func(a, b Value) int, initCapacity int) *Finder {
	sl := &Finder{
		seed:       uint64(time.Now().UnixNano()/2 + 1),
		comparator: comparator,
		nodes:      make([]node, 0, max(initCapacity, 16)),
	}
	sl.nodes = append(sl.nodes, node{size: 0, priority: sl.nextRand()}) // dummy node 0
	return sl
}

func (sl *Finder) Build(nums []Value) int {
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

func (sl *Finder) build(root int) {
	nodeRef := &sl.nodes[root]
	if nodeRef.left != 0 {
		sl.build(nodeRef.left)
	}
	if nodeRef.right != 0 {
		sl.build(nodeRef.right)
	}
	sl.pushUp(root)
}

func (sl *Finder) pushUp(root int) {
	sl.nodes[root].size = sl.nodes[sl.nodes[root].left].size + sl.nodes[sl.nodes[root].right].size + 1
}

func (sl *Finder) Insert(value Value) {
	var x, y, z int
	sl.splitByValue(sl.root, value, &x, &y, false)
	z = sl.newNode(value)
	sl.root = sl.merge(sl.merge(x, z), y)
}

func (sl *Finder) Erase(value Value) {
	var x, y, z int
	sl.splitByValue(sl.root, value, &x, &z, false)
	sl.splitByValue(x, value, &x, &y, true)
	y = sl.merge(sl.nodes[y].left, sl.nodes[y].right)
	sl.root = sl.merge(sl.merge(x, y), z)
}

// 求小于等于 value 的最大值.不存在则第二个返回值为 false.
func (sl *Finder) Prev(value Value) (res Value, ok bool) {
	var x, y int
	sl.splitByValue(sl.root, value, &x, &y, false)
	if x == 0 {
		ok = false
		return
	}
	res = sl.nodes[sl.kthNode(x, sl.nodes[x].size)].value
	sl.root = sl.merge(x, y)
	ok = true
	return
}

// 求大于等于 value 的最小值.不存在则第二个返回值为 false.
func (sl *Finder) Next(value Value) (res Value, ok bool) {
	var x, y int
	sl.splitByValue(sl.root, value, &x, &y, true)
	if y == 0 {
		ok = false
		return
	}
	res = sl.nodes[sl.kthNode(y, 1)].value
	sl.root = sl.merge(x, y)
	ok = true
	return
}

func (sl *Finder) String() string {
	sb := []string{"SortedList{"}
	values := []string{}
	for i := 0; i < sl.Len(); i++ {
		values = append(values, fmt.Sprintf("%v", sl.at(i)))
	}
	sb = append(sb, strings.Join(values, ","), "}")
	return strings.Join(sb, "")
}

func (sl *Finder) Len() int {
	return sl.nodes[sl.root].size
}

func (sl *Finder) at(index int) Value {
	n := sl.Len()
	if index < 0 {
		index += n
	}
	if index < 0 || index >= n {
		panic(fmt.Sprintf("%d index out of range: [%d,%d]", index, 0, n-1))
	}
	return sl.nodes[sl.kthNode(sl.root, index+1)].value
}

func (sl *Finder) kthNode(root int, k int) int {
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

func (sl *Finder) splitByValue(root int, value Value, x, y *int, strictLess bool) {
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
func (sl *Finder) splitByRank(root, k int, x, y *int) {
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

func (sl *Finder) merge(x, y int) int {
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
func (sl *Finder) InOrder() []Value {
	res := make([]Value, 0, sl.Len())
	sl.inOrder(sl.root, &res)
	return res
}

func (sl *Finder) inOrder(root int, res *[]Value) {
	if root == 0 {
		return
	}
	sl.inOrder(sl.nodes[root].left, res)
	*res = append(*res, sl.nodes[root].value)
	sl.inOrder(sl.nodes[root].right, res)
}

func (sl *Finder) newNode(value Value) int {
	sl.nodes = append(sl.nodes, node{
		value:    value,
		size:     1,
		priority: sl.nextRand(),
	})
	return len(sl.nodes) - 1
}

// https://nyaannyaan.github.io/library/misc/rng.hpp
func (sl *Finder) nextRand() uint64 {
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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// 在线bfs.
//   不预先给出图，而是通过两个函数 setUsed 和 findUnused 来在线寻找边.
//   setUsed(u)：将 u 标记为已访问。
//   findUnused(u)：找到和 u 邻接的一个未访问过的点。如果不存在, 返回 `-1`。

func OnlineBfs(
	n int, start int,
	setUsed func(u int), findUnused func(cur int) (next int),
) (dist []int) {
	dist = make([]int, n)
	for i := range dist {
		dist[i] = INF
	}
	dist[start] = 0
	queue := []int{start}
	setUsed(start)

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		for {
			next := findUnused(cur)
			if next == -1 {
				break
			}
			dist[next] = dist[cur] + 1 // weight
			queue = append(queue, next)
			setUsed(next)
		}
	}

	return
}
