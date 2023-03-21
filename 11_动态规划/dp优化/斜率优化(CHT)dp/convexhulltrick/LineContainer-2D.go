// 注意cpp里的迭代器:
// !Begin指向第一个元素,
// !End指向最后一个元素的下一个位置,
// 这里的迭代器设计为:
// !Begin指向第一个元素的前一个位置,First指向第一个元素
// !Last指向最后一个元素,End指向最后一个元素的下一个位置

// https://maspypy.github.io/library/convex/cht.hpp

// 在 C++ 中，long double 类型不等同于 float64。
// !long double 是一种浮点数类型，具有比 double 类型更高的精度和范围 (18位)。
// int 通常对应于 C++ 中的 double 类型，而非 long double 类型。
// long double 类型的精度和范围因编译器和平台而异。
// 在某些实现中，long double 可能与 double 类型具有相同的精度，
// 而在其他实现中，它可能具有更高的精度。
// 例如，在 x86 和 x86_64 架构上，long double 通常具有 80 位的扩展精度。

package main

import (
	"bufio"
	"fmt"
	stdio "io"
	"os"
	"strconv"
	"strings"
	"time"
)

// from https://atcoder.jp/users/ccppjsrb
var io *Iost

type Iost struct {
	Scanner *bufio.Scanner
	Writer  *bufio.Writer
}

func NewIost(fp stdio.Reader, wfp stdio.Writer) *Iost {
	const BufSize = 2000005
	scanner := bufio.NewScanner(fp)
	scanner.Split(bufio.ScanWords)
	scanner.Buffer(make([]byte, BufSize), BufSize)
	return &Iost{Scanner: scanner, Writer: bufio.NewWriter(wfp)}
}
func (io *Iost) Text() string {
	if !io.Scanner.Scan() {
		panic("scan failed")
	}
	return io.Scanner.Text()
}
func (io *Iost) Atoi(s string) int                 { x, _ := strconv.Atoi(s); return x }
func (io *Iost) Atoi64(s string) int64             { x, _ := strconv.ParseInt(s, 10, 64); return x }
func (io *Iost) Atof64(s string) float64           { x, _ := strconv.ParseFloat(s, 64); return x }
func (io *Iost) NextInt() int                      { return io.Atoi(io.Text()) }
func (io *Iost) NextInt64() int64                  { return io.Atoi64(io.Text()) }
func (io *Iost) NextFloat64() float64              { return io.Atof64(io.Text()) }
func (io *Iost) Print(x ...interface{})            { fmt.Fprint(io.Writer, x...) }
func (io *Iost) Printf(s string, x ...interface{}) { fmt.Fprintf(io.Writer, s, x...) }
func (io *Iost) Println(x ...interface{})          { fmt.Fprintln(io.Writer, x...) }

func main() {
	最大三角形面积()
}

func abc244_h() {
	// https://atcoder.jp/contests/abc244/tasks/abc244_h
	// - 向点集中追加一个点(a,b), 表示为 a*x + b*y
	// - 查询 x=xi,y=yi 时的最大值
	in := os.Stdin
	out := os.Stdout
	io = NewIost(in, out)
	defer func() {
		io.Writer.Flush()
	}()

	q := io.NextInt()
	cht := NewLineContainer2D(q)
	for i := 0; i < q; i++ {
		a, b, x, y := io.NextInt(), io.NextInt(), io.NextInt(), io.NextInt()
		cht.Add(a, b)
		io.Println(cht.QueryMax(x, y))
	}
}

func 最大三角形面积() {
	// https://yukicoder.me/problems/no/2012
	// 平面上有n个点, 问最其中两点和原点组成的三角形的最大面积的2倍
	// !将(a,-b)加入点集，对每个点(x,y)查询最大的a*y-b*x
	in := os.Stdin
	out := os.Stdout
	io = NewIost(in, out)
	defer func() {
		io.Writer.Flush()
	}()

	n := io.NextInt()
	points := make([][2]int, n)
	for i := 0; i < n; i++ {
		points[i] = [2]int{io.NextInt(), io.NextInt()}
	}

	cht := NewLineContainer2D(n)
	for _, p := range points {
		a, b := p[0], p[1]
		cht.Add(a, -b)
	}

	res := 0
	for _, p := range points {
		x, y := p[0], p[1]
		res = max(res, cht.QueryMax(y, x))
	}
	io.Println(res)
}

const INF int = 1e18

type Line struct {
	k, b   int
	p1, p2 int // p=p1/p2
}

type LineContainer2D struct {
	minCHT, maxCHT *_LineContainer
	kMax, kMin     int
	bMax, bMin     int
}

func NewLineContainer2D(capacity int) *LineContainer2D {
	return &LineContainer2D{
		minCHT: _NewLineContainer(true, capacity),
		maxCHT: _NewLineContainer(false, capacity),
		kMax:   -INF,
		kMin:   INF,
		bMax:   -INF,
		bMin:   INF,
	}
}

// 追加 a*x + b*y.
func (lc *LineContainer2D) Add(a, b int) {
	lc.minCHT.Add(b, a)
	lc.maxCHT.Add(b, a)
	lc.kMax = max(lc.kMax, a)
	lc.kMin = min(lc.kMin, a)
	lc.bMax = max(lc.bMax, b)
	lc.bMin = min(lc.bMin, b)
}

// 查询 x=xi,y=yi 时的最大值 max_{a,b} (ax + by).
func (lc *LineContainer2D) QueryMax(x, y int) int {
	if lc.minCHT.Size() == 0 {
		return -INF
	}

	if x == 0 {
		if y > 0 {
			return lc.bMax * y
		}
		return lc.bMin * y
	}
	if y == 0 {
		if x > 0 {
			return lc.kMax * x
		}
		return lc.kMin * x
	}

	// y/x
	if x > 0 {
		l := lc.maxCHT.sl.LowerBoundWith(y, x)
		line := lc.maxCHT.sl.At(l)
		a := line.b
		b := line.k
		return a*x + b*y
	}
	l := lc.minCHT.sl.LowerBoundWith(y, x)
	line := lc.minCHT.sl.At(l)
	a := -line.b
	b := -line.k
	return a*x + b*y
}

// 查询 x=xi,y=yi 时的最小值 min_{a,b} (ax + by).
func (lc *LineContainer2D) QueryMin(x, y int) int { return -lc.QueryMax(-x, -y) }

type _LineContainer struct {
	minimize bool
	sl       *_SL
}

func _NewLineContainer(minimize bool, capacity int) *_LineContainer {
	return &_LineContainer{
		minimize: minimize,
		sl:       _NSL(capacity),
	}
}

func (lc *_LineContainer) Add(k, m int) {
	if lc.minimize {
		k, m = -k, -m
	}

	newLine := Line{k: k, b: m}
	lc.sl.Add(newLine)
	it1 := lc.sl.BisectRight(newLine.k) - 1
	it2 := it1
	line2 := lc.sl.At(it2)
	it1++
	it3 := it2
	for lc.insect(line2, lc.sl.At(it1)) {
		lc.sl.Pop(it1)
	}

	if it3 != 0 {
		it3--
		line3 := lc.sl.At(it3)
		if lc.insect(line3, line2) {
			lc.sl.Pop(it2)
			lc.insect(line3, lc.sl.At(it2))
		}
	}

	if it3 == 0 {
		return
	}

	dp1, dp2 := lc.sl.At(it3-1), lc.sl.At(it3)
	for it3 != 0 {
		it2 := it3
		if less(dp1.p1, dp1.p2, dp2.p1, dp2.p2) {
			break
		}
		it3--
		lc.sl.Pop(it2)
		lc.insect(dp1, lc.sl.At(it2))
		dp1, dp2 = lc.sl.At(it3-1), dp1
	}
}

// 查询 kx + m 的最小值（或最大值).
func (lc *_LineContainer) Query(x int) int {
	if lc.sl.Len() == 0 {
		panic("empty container")
	}
	pos := lc.sl.LowerBoundWith(x, 1)
	line := lc.sl.At(pos)
	v := line.k*x + line.b
	if lc.minimize {
		return -v
	}
	return v
}

func (lc *_LineContainer) Size() int { return lc.sl.Len() }

// 这个函数在向集合添加新线或删除旧线时用于计算交点。
// 计算线性函数x和y的交点，并将结果存储在x->p中。
func (lc *_LineContainer) insect(line1, line2 *Line) bool {
	if line2 == nil {
		line1.p1 = INF
		line1.p2 = 1
		return false
	}
	if line1.k == line2.k {
		if line1.b > line2.b {
			line1.p1 = INF
			line1.p2 = 1
		} else {
			line1.p1 = INF
			line1.p2 = -1
		}
	} else {
		// lc_div
		line1.p1 = line2.b - line1.b
		line1.p2 = line1.k - line2.k
	}
	return !less(line1.p1, line1.p2, line2.p1, line2.p2)
}

// DIY: 传入自定义比较函数的LowerBound
func (sl *_SL) LowerBoundWith(a, b int) int {
	return sl.BisectLeftWith(a, b)
}

type _Value = Line

type _node struct {
	left, right int
	size        int
	priority    uint64
	value       _Value
}

type _SL struct {
	seed  uint64
	root  int
	nodes []_node
}

func _NSL(initCapacity int) *_SL {
	sl := &_SL{
		seed:  uint64(time.Now().UnixNano()/2 + 1),
		nodes: make([]_node, 0, max(initCapacity, 16)),
	}
	dummy := &_node{size: 0, priority: sl.nextRand(), value: _Value{p2: 1}} // dummy node 0
	sl.nodes = append(sl.nodes, *dummy)
	return sl
}

func (sl *_SL) pushUp(root int) {
	sl.nodes[root].size = sl.nodes[sl.nodes[root].left].size + sl.nodes[sl.nodes[root].right].size + 1
}

func (sl *_SL) Add(value _Value) {
	var x, y, z int
	sl.splitByValue(sl.root, value.k, &x, &y, false)
	z = sl.newNode(value)
	sl.root = sl.merge(sl.merge(x, z), y)
}

func (sl *_SL) At(index int) *_Value {
	if index < 0 || index >= sl.Len() {
		return nil
	}
	return &sl.nodes[sl.kthNode(sl.root, index+1)].value
}

func (sl *_SL) Pop(index int) _Value {
	index += 1 // dummy offset
	var x, y, z int
	sl.splitByRank(sl.root, index, &y, &z)
	sl.splitByRank(y, index-1, &x, &y)
	res := sl.nodes[y].value
	sl.root = sl.merge(x, z)
	return res
}

func (sl *_SL) BisectLeft(k int) int {
	var x, y int
	sl.splitByValue(sl.root, k, &x, &y, true)
	res := sl.nodes[x].size
	sl.root = sl.merge(x, y)
	return res
}

func (sl *_SL) BisectLeftWith(a, b int) int {
	var x, y int
	sl.splitByValueWith(sl.root, a, b, &x, &y)
	res := sl.nodes[x].size
	sl.root = sl.merge(x, y)
	return res
}

func (sl *_SL) BisectRight(k int) int {
	var x, y int
	sl.splitByValue(sl.root, k, &x, &y, false)
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

func (sl *_SL) splitByValue(root int, k int, x, y *int, strictLess bool) {
	if root == 0 {
		*x, *y = 0, 0
		return
	}
	if strictLess {
		if sl.nodes[root].value.k < k {
			*x = root
			sl.splitByValue(sl.nodes[root].right, k, &sl.nodes[root].right, y, strictLess)
		} else {
			*y = root
			sl.splitByValue(sl.nodes[root].left, k, x, &sl.nodes[root].left, strictLess)
		}
	} else {
		if sl.nodes[root].value.k <= k {
			*x = root
			sl.splitByValue(sl.nodes[root].right, k, &sl.nodes[root].right, y, strictLess)
		} else {
			*y = root
			sl.splitByValue(sl.nodes[root].left, k, x, &sl.nodes[root].left, strictLess)
		}
	}
	sl.pushUp(root)
}

func (sl *_SL) splitByValueWith(root int, a1, b1 int, x, y *int) {
	if root == 0 {
		*x, *y = 0, 0
		return
	}
	a2, b2 := sl.nodes[root].value.p1, sl.nodes[root].value.p2
	if less(a2, b2, a1, b1) {
		*x = root
		sl.splitByValueWith(sl.nodes[root].right, a1, b1, &sl.nodes[root].right, y)
	} else {
		*y = root
		sl.splitByValueWith(sl.nodes[root].left, a1, b1, x, &sl.nodes[root].left)
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

func (sl *_SL) newNode(value _Value) int {
	sl.nodes = append(sl.nodes, _node{
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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// 分母不为0的分数比较大小
//  a1/b1 < a2/b2
func less(a1, b1, a2, b2 int) bool {
	if a1 == INF || a2 == INF { // 有一个是+-INF
		return a1/b1 < a2/b2
	}
	diff := a1*b2 - a2*b1
	mul := b1 * b2
	return diff^mul < 0
}
