// 注意cpp里的迭代器:
// !Begin指向第一个元素,
// !End指向最后一个元素的下一个位置,
// 这里的迭代器设计为:
// !Begin指向第一个元素的前一个位置,First指向第一个元素
// !Last指向最后一个元素,End指向最后一个元素的下一个位置

// 这里删除元素会引起迭代器失效

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

func lineAddGetMin() {
	// https://judge.yosupo.jp/problem/line_add_get_min
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	L := NewLineContainer(true)
	for i := 0; i < n; i++ {
		var k, m int
		fmt.Fscan(in, &k, &m)
		L.Add(k, m)
	}
	for i := 0; i < q; i++ {
		var op int
		fmt.Fscan(in, &op)
		if op == 0 {
			var a, b int
			fmt.Fscan(in, &a, &b)
			L.Add(a, b)
		} else {
			var x int
			fmt.Fscan(in, &x)
			fmt.Fprintln(out, L.Query(x))
		}
	}
}

func main() {
	lineAddGetMin()
}

const INF int = 2e18

type Line struct {
	k, m, p int
}

type LineContainer struct {
	minimize bool
	sl       *_SL
}

func NewLineContainer(minimize bool) *LineContainer {
	return &LineContainer{
		minimize: minimize,
		sl:       _NSL(func(a, b Value) int { return a.k - b.k }, 16),
	}
}

// 向集合中添加一条线，表示为y = kx + m
func (lc *LineContainer) Add(k, m int) {
	if lc.minimize {
		k, m = -k, -m
	}
	newLine := &Line{k: k, m: m}
	lc.sl.Add(newLine)
	it1 := lc.sl.BisectRight(newLine) - 1
	it2 := it1
	it1++
	for lc.insect(it2, it1) {
		lc.sl.Pop(it1)
	}

	it3 := it2
	if it3 != 0 {
		it3--
		if lc.insect(it3, it2) {
			lc.sl.Pop(it2)
			lc.insect(it3, it2)
		}
	}

	for it2 := it3; it2 != 0 && lc.sl.At(it3-1).p >= lc.sl.At(it2).p; it2 = it3 {
		it3--
		lc.sl.Pop(it2)
		lc.insect(it3, it2)
	}
}

// 查询 kx + m 的最小值（或最大值).
func (lc *LineContainer) Query(x int) int {
	// !这里有一个关键点：尽管Line<T>结构体中的operator<按k值对线性函数进行排序，
	// !但LineContainer类在维护这些线性函数时，确保了它们的交点的x坐标（p值）是有序的。
	// 这使得query函数可以通过调用lower_bound(x)来找到给定x值对应的最大（或最小）y值。
	if lc.sl.Len() == 0 {
		panic("empty container")
	}
	pos := lc.sl.LowerBoundWithP(x)
	line := lc.sl.At(pos)
	v := line.k*x + line.m
	if lc.minimize {
		return -v
	}
	return v
}

// 这个函数在向集合添加新线或删除旧线时用于计算交点。
// 计算线性函数x和y的交点，并将结果存储在x->p中。
func (lc *LineContainer) insect(posX, posY int) bool {
	if posY == lc.sl.Len() {
		line := lc.sl.At(posX)
		line.p = INF
		return false
	}

	line1, line2 := lc.sl.At(posX), lc.sl.At(posY)
	if line1.k == line2.k {
		if line1.m > line2.m {
			line1.p = INF
		} else {
			line1.p = -INF
		}
	} else {
		// lc_div
		a, b := line2.m-line1.m, line1.k-line2.k
		tmp := 0
		if (a^b) < 0 && a%b != 0 {
			tmp = 1
		}
		line1.p = a/b - tmp
	}
	return line1.p >= line2.p
}

// DIY: 传入自定义比较函数的LowerBound
func (sl *_SL) LowerBoundWithP(p int) int {
	return sl.BisectLeftWith(p)
}

type Value = *Line

type node struct {
	left, right int
	size        int
	priority    uint64
	value       Value
}

type _SL struct {
	seed       uint64
	root       int
	comparator func(a, b Value) int
	nodes      []node
}

func _NSL(comparator func(a, b Value) int, initCapacity int) *_SL {
	sl := &_SL{
		seed:       uint64(time.Now().UnixNano()/2 + 1),
		comparator: comparator,
		nodes:      make([]node, 0, max(initCapacity, 16)),
	}
	dummy := &node{size: 0, priority: sl.nextRand()} // dummy node 0
	sl.nodes = append(sl.nodes, *dummy)
	return sl
}

func (sl *_SL) pushUp(root int) {
	sl.nodes[root].size = sl.nodes[sl.nodes[root].left].size + sl.nodes[sl.nodes[root].right].size + 1
}

func (sl *_SL) Add(value Value) {
	var x, y, z int
	sl.splitByValue(sl.root, value, &x, &y, false)
	z = sl.newNode(value)
	sl.root = sl.merge(sl.merge(x, z), y)
}

func (sl *_SL) At(index int) Value {
	return sl.nodes[sl.kthNode(sl.root, index+1)].value
}

func (sl *_SL) Pop(index int) Value {
	index += 1 // dummy offset
	var x, y, z int
	sl.splitByRank(sl.root, index, &y, &z)
	sl.splitByRank(y, index-1, &x, &y)
	res := sl.nodes[y].value
	sl.root = sl.merge(x, z)
	return res
}

func (sl *_SL) BisectLeft(value Value) int {
	var x, y int
	sl.splitByValue(sl.root, value, &x, &y, true)
	res := sl.nodes[x].size
	sl.root = sl.merge(x, y)
	return res
}

func (sl *_SL) BisectLeftWith(p int) int {
	var x, y int
	sl.splitByValueWith(sl.root, p, &x, &y)
	res := sl.nodes[x].size
	sl.root = sl.merge(x, y)
	return res
}

func (sl *_SL) BisectRight(value Value) int {
	var x, y int
	sl.splitByValue(sl.root, value, &x, &y, false)
	res := sl.nodes[x].size
	sl.root = sl.merge(x, y)
	return res
}

func (sl *_SL) String() string {
	sb := []string{"SortedList{"}
	values := []string{}
	for i := 0; i < sl.Len(); i++ {
		values = append(values, fmt.Sprintf("%v", sl.At(i)))
	}
	sb = append(sb, strings.Join(values, ","), "}")
	return strings.Join(sb, "")
}

func (sl *_SL) Len() int {
	return sl.nodes[sl.root].size
}

func (sl *_SL) kthNode(root int, k int) int {
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

func (sl *_SL) splitByValue(root int, value Value, x, y *int, strictLess bool) {
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

func (sl *_SL) splitByValueWith(root int, p int, x, y *int) {
	if root == 0 {
		*x, *y = 0, 0
		return
	}
	if sl.nodes[root].value.p < p {
		*x = root
		sl.splitByValueWith(sl.nodes[root].right, p, &sl.nodes[root].right, y)
	} else {
		*y = root
		sl.splitByValueWith(sl.nodes[root].left, p, x, &sl.nodes[root].left)
	}
	sl.pushUp(root)
}

// Split by rank.
// Split the tree rooted at root into two trees, x and y, such that the size of x is k.
// x is the left subtree, y is the right subtree.
func (sl *_SL) splitByRank(root, k int, x, y *int) {
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

func (sl *_SL) merge(x, y int) int {
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

func (sl *_SL) newNode(value Value) int {
	sl.nodes = append(sl.nodes, node{
		value:    value,
		size:     1,
		priority: sl.nextRand(),
	})
	return len(sl.nodes) - 1
}

// https://nyaannyaan.github.io/library/misc/rng.hpp
func (sl *_SL) nextRand() uint64 {
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
