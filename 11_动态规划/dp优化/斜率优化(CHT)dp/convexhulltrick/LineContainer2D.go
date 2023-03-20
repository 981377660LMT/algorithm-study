// 注意cpp里的迭代器:
// !Begin指向第一个元素,
// !End指向最后一个元素的下一个位置,
// 这里的迭代器设计为:
// !Begin指向第一个元素的前一个位置,First指向第一个元素
// !Last指向最后一个元素,End指向最后一个元素的下一个位置

// TODO
// https://maspypy.github.io/library/convex/cht.hpp
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

// -309827997742946804
// 219304001222283930
// 238136177859802308
// 525969950014529719
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var q int
	fmt.Fscan(in, &q)
	cht := NewLineContainer2D()
	for i := 0; i < q; i++ {
		var a, b, x, y int
		fmt.Fscan(in, &a, &b, &x, &y)
		cht.Add(a, b)
		fmt.Fprintln(out, cht.GetMax(x, y))
	}

}

const INF_INT int = 4e18

type Line struct {
	k, m, p int
}

type LineContainer2D struct {
	minCHT, maxCHT *LineContainer
	kMax, kMin     int
	bMax, bMin     int
	empty          bool
}

func NewLineContainer2D() *LineContainer2D {
	return &LineContainer2D{
		minCHT: NewLineContainer(true),
		maxCHT: NewLineContainer(false),
		kMax:   -INF_INT,
		kMin:   INF_INT,
		bMax:   -INF_INT,
		bMin:   INF_INT,
		empty:  true,
	}
}

// 向集合中添加一条线，表示为y = kx + b
func (lc *LineContainer2D) Add(k, b int) {
	lc.empty = false
	lc.minCHT.Add(b, k)
	lc.maxCHT.Add(b, k)
	lc.kMax = max(lc.kMax, k)
	lc.kMin = min(lc.kMin, k)
	lc.bMax = max(lc.bMax, b)
	lc.bMin = min(lc.bMin, b)
}

func lcDivInt(a, b int) int {
	tmp := 0
	if (a^b) < 0 && a%b != 0 {
		tmp = 1
	}
	return a/b - tmp
}

func (lc *LineContainer2D) GetMax(x, y int) int {
	if lc.minCHT.Size() == 0 {
		return -INF_INT
	}

	if x == 0 {
		return max(lc.bMax*y, lc.bMin*y)
	}

	z := lcDivInt(y, x)
	if x > 0 {
		l := lc.maxCHT.sl.LowerBoundWithP(z, func(p1, p2 int) bool { return p1 < p2 })
		line := lc.maxCHT.sl.At(l)
		a := line.m
		b := line.k
		return (a*x + b*y)
	}
	l := lc.minCHT.sl.LowerBoundWithP(z, func(p1, p2 int) bool { return p1 < p2 })
	line := lc.minCHT.sl.At(l)
	a := -line.m
	b := -line.k
	return (a*x + b*y)
}

func (lc *LineContainer2D) GetMin(x, y int) int { return -lc.GetMax(-x, -y) }

type LineContainer struct {
	minimize bool
	sl       *SortedList
}

func NewLineContainer(minimize bool) *LineContainer {
	return &LineContainer{
		minimize: minimize,
		sl:       NewSortedList(func(a, b Value) bool { return a.k < b.k }, 16),
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
	it3 := it2
	for lc.insect(it2, it1) {
		lc.sl.Pop(it1)
	}

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
	pos := lc.sl.LowerBoundWithP(x, func(p1, p2 int) bool { return p1 < p2 })

	line := lc.sl.At(pos)
	v := line.k*x + line.m
	if lc.minimize {
		return -v
	}
	return v
}

func (lc *LineContainer) Size() int { return lc.sl.Len() }

// 这个函数在向集合添加新线或删除旧线时用于计算交点。
// 计算线性函数x和y的交点，并将结果存储在x->p中。
func (lc *LineContainer) insect(posX, posY int) bool {
	if posY == lc.sl.Len() {
		line := lc.sl.At(posX)
		line.p = INF_INT
		return false
	}

	line1, line2 := lc.sl.At(posX), lc.sl.At(posY)
	if line1.k == line2.k {
		if line1.m > line2.m {
			line1.p = INF_INT
		} else {
			line1.p = -INF_INT
		}
	} else {
		// lc_div
		a, b := line2.m-line1.m, line1.k-line2.k
		line1.p = lcDivInt(a, b)
	}
	return line1.p >= line2.p
}

// DIY: 传入自定义比较函数的LowerBound
func (sl *SortedList) LowerBoundWithP(p int, less func(p1, p2 int) bool) int {
	return sl.BisectLeftPWith(p, less)
}

type Value = *Line

type node struct {
	left, right int
	size        int
	priority    uint64
	value       Value
}

type SortedList struct {
	seed  uint64
	root  int
	less  func(a, b Value) bool
	nodes []node
}

func NewSortedList(less func(a, b Value) bool, initCapacity int) *SortedList {
	sl := &SortedList{
		seed:  uint64(time.Now().UnixNano()/2 + 1),
		less:  less,
		nodes: make([]node, 0, max(initCapacity, 16)),
	}

	dummy := &node{size: 0, priority: sl.nextRand()} // dummy node 0
	sl.nodes = append(sl.nodes, *dummy)
	return sl
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

func (sl *SortedList) BisectLeftPWith(p int, less func(a, b int) bool) int {
	var x, y int
	sl.splitByValueWith(sl.root, p, &x, &y, less)
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
		if sl.less(sl.nodes[root].value, value) {
			*x = root
			sl.splitByValue(sl.nodes[root].right, value, &sl.nodes[root].right, y, strictLess)
		} else {
			*y = root
			sl.splitByValue(sl.nodes[root].left, value, x, &sl.nodes[root].left, strictLess)
		}
	} else {
		if !sl.less(value, sl.nodes[root].value) {
			*x = root
			sl.splitByValue(sl.nodes[root].right, value, &sl.nodes[root].right, y, strictLess)
		} else {
			*y = root
			sl.splitByValue(sl.nodes[root].left, value, x, &sl.nodes[root].left, strictLess)
		}
	}

	sl.pushUp(root)
}

func (sl *SortedList) splitByValueWith(root int, p int, x, y *int, less func(a, b int) bool) {
	if root == 0 {
		*x, *y = 0, 0
		return
	}
	if less(sl.nodes[root].value.p, p) {
		*x = root
		sl.splitByValueWith(sl.nodes[root].right, p, &sl.nodes[root].right, y, less)
	} else {
		*y = root
		sl.splitByValueWith(sl.nodes[root].left, p, x, &sl.nodes[root].left, less)
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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
