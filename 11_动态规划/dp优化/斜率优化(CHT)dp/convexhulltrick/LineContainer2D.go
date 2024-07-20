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

//
// 为什么使用SortedList而不是Treap?
// -> .At()调用次数多，Treap的时间复杂度为O(logn)，而SortedList的时间复杂度接近O(1)

package main

import (
	"bufio"
	"fmt"
	stdio "io"
	"math/bits"
	"os"
	"sort"
	"strconv"
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
	// 最大三角形面积()
	abc244_h()
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
	cht := NewLineContainer2D()
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

	cht := NewLineContainer2D()
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

const INF int = 4e18

type Line struct {
	k, b   int
	p1, p2 int // p=p1/p2
}

type LineContainer2D struct {
	minCHT, maxCHT *_LineContainer
	kMax, kMin     int
	bMax, bMin     int
}

func NewLineContainer2D() *LineContainer2D {
	return &LineContainer2D{
		minCHT: _NewLineContainer(true),
		maxCHT: _NewLineContainer(false),
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
		line := lc.maxCHT.sl.BisectLeftByPairForValue(y, x)
		a := line.b
		b := line.k
		return a*x + b*y
	}
	line := lc.minCHT.sl.BisectLeftByPairForValue(y, x)
	a := -line.b
	b := -line.k
	return a*x + b*y
}

// 查询 x=xi,y=yi 时的最小值 min_{a,b} (ax + by).
func (lc *LineContainer2D) QueryMin(x, y int) int { return -lc.QueryMax(-x, -y) }

func (lc *LineContainer2D) Clear() {
	lc.minCHT.Clear()
	lc.maxCHT.Clear()
	lc.kMax = -INF
	lc.kMin = INF
	lc.bMax = -INF
	lc.bMin = INF
}

type _LineContainer struct {
	minimize bool
	sl       *SpecializedSortedList
}

func _NewLineContainer(minimize bool) *_LineContainer {
	return &_LineContainer{
		minimize: minimize,
		sl:       NewSpecializedSortedList(func(a, b *Line) bool { return a.k < b.k }),
	}
}

func (lc *_LineContainer) Add(k, m int) {
	if lc.minimize {
		k, m = -k, -m
	}

	newLine := &Line{k: k, b: m}
	lc.sl.Add(newLine)

	iter := lc.sl.BisectRightByKForIterator(k)
	iter.Prev()

	{
		probe := iter.Copy()
		probe.Next()
		start := probe.ToIndex()
		removeCount := int32(0)
		for lc.insect(iter, probe) {
			probe.Next()
			removeCount++
		}
		lc.sl.Erase(start, start+removeCount)
	}

	{
		probe := iter.Copy()
		if !iter.IsBegin() {
			iter.Prev()
			if lc.insect(iter, probe) {
				probIndex := probe.ToIndex()
				probe.Next()
				lc.insect(iter, probe)
				lc.sl.Pop(probIndex)
			}
		}
	}

	if iter.IsBegin() {
		return
	}

	{
		var pivot *Line
		if iter.HasNext() {
			pivot = iter.NextValue()
		}
		end := iter.ToIndex() + 1
		removeCount := int32(0)
		for !iter.IsBegin() {
			iter.Prev()
			if lessLine(iter.Value(), iter.NextValue()) {
				break
			}
			lc.insectLine(iter.Value(), pivot)
			removeCount++
		}
		lc.sl.Erase(end-removeCount, end)
	}
}

// 查询 kx + m 的最小值（或最大值).
func (lc *_LineContainer) Query(x int) int {
	if lc.sl.Len() == 0 {
		panic("empty container")
	}
	line := lc.sl.BisectLeftByPairForValue(x, 1)
	v := line.k*x + line.b
	if lc.minimize {
		return -v
	}
	return v
}

func (lc *_LineContainer) Clear() { lc.sl.Clear() }

func (lc *_LineContainer) Size() int32 { return lc.sl.Len() }

// 这个函数在向集合添加新线或删除旧线时用于计算交点。
// 计算线性函数x和y的交点，并将结果存储在x->p中。
func (lc *_LineContainer) insect(iterX, iterY *Iterator) bool {
	if iterY.IsEnd() {
		line1 := iterX.Value()
		line1.p1 = INF
		line1.p2 = 1
		return false
	}
	line1, line2 := iterX.Value(), iterY.Value()
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
	return !lessPair(line1.p1, line1.p2, line2.p1, line2.p2)
}

func (lc *_LineContainer) insectLine(line1, line2 *Line) bool {
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
	return !lessPair(line1.p1, line1.p2, line2.p1, line2.p2)
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
//
//	a1/b1 < a2/b2
func lessPair(a1, b1, a2, b2 int) bool {
	if a1 == INF || a2 == INF { // 有一个是+-INF
		return a1/b1 < a2/b2
	}
	diff := a1*b2 - a2*b1
	mul := b1 * b2
	return diff^mul < 0
}

func lessLine(a, b *Line) bool {
	return lessPair(a.p1, a.p2, b.p1, b.p2)
}

const _LOAD int32 = 75 // 75/100/150/200

type S = *Line

type SpecializedSortedList struct {
	less              func(a, b S) bool
	size              int32
	blocks            [][]S
	mins              []S
	tree              []int32
	shouldRebuildTree bool
}

func NewSpecializedSortedList(less func(a, b S) bool, elements ...S) *SpecializedSortedList {
	elements = append(elements[:0:0], elements...)
	res := &SpecializedSortedList{less: less}
	sort.Slice(elements, func(i, j int) bool { return less(elements[i], elements[j]) })
	n := int32(len(elements))
	blocks := [][]S{}
	for start := int32(0); start < n; start += _LOAD {
		end := min32(start+_LOAD, n)
		blocks = append(blocks, elements[start:end:end]) // !各个块互不影响, max参数也需要指定为end
	}
	mins := make([]S, len(blocks))
	for i, cur := range blocks {
		mins[i] = cur[0]
	}
	res.size = n
	res.blocks = blocks
	res.mins = mins
	res.shouldRebuildTree = true
	return res
}

func (sl *SpecializedSortedList) Erase(start, end int32) {
	sl.Enumerate(start, end, nil, true)
}

func (sl *SpecializedSortedList) Enumerate(start, end int32, f func(value S), erase bool) {
	if start < 0 {
		start = 0
	}
	if end > sl.size {
		end = sl.size
	}
	if start >= end {
		return
	}

	pos, startIndex := sl._findKth(start)
	count := end - start
	m := int32(len(sl.blocks))
	for ; count > 0 && pos < m; pos++ {
		block := sl.blocks[pos]
		endIndex := min32(int32(len(block)), startIndex+count)
		if f != nil {
			for j := startIndex; j < endIndex; j++ {
				f(block[j])
			}
		}
		deleted := endIndex - startIndex

		if erase {
			if deleted == int32(len(block)) {
				// !delete block
				sl.blocks = append(sl.blocks[:pos], sl.blocks[pos+1:]...)
				sl.mins = append(sl.mins[:pos], sl.mins[pos+1:]...)
				sl.shouldRebuildTree = true
				pos--
			} else {
				// !delete [index, end)
				sl._updateTree(pos, -deleted)
				sl.blocks[pos] = append(sl.blocks[pos][:startIndex], sl.blocks[pos][endIndex:]...)
				sl.mins[pos] = sl.blocks[pos][0]
			}
			sl.size -= deleted
		}

		count -= deleted
		startIndex = 0
	}
}

func (sl *SpecializedSortedList) Add(value S) *SpecializedSortedList {
	sl.size++
	if len(sl.blocks) == 0 {
		sl.blocks = append(sl.blocks, []S{value})
		sl.mins = append(sl.mins, value)
		sl.shouldRebuildTree = true
		return sl
	}

	pos, index := sl._locRight(value)

	sl._updateTree(pos, 1)
	sl.blocks[pos] = append(sl.blocks[pos][:index], append([]S{value}, sl.blocks[pos][index:]...)...)
	sl.mins[pos] = sl.blocks[pos][0]

	// n -> load + (n - load)
	if n := int32(len(sl.blocks[pos])); _LOAD+_LOAD < n {
		sl.blocks = append(sl.blocks[:pos+1], append([][]S{sl.blocks[pos][_LOAD:]}, sl.blocks[pos+1:]...)...)
		sl.mins = append(sl.mins[:pos+1], append([]S{sl.blocks[pos][_LOAD]}, sl.mins[pos+1:]...)...)
		sl.blocks[pos] = sl.blocks[pos][:_LOAD:_LOAD] // !注意max的设置(为了让左右互不影响)
		sl.shouldRebuildTree = true
	}

	return sl
}

func (sl *SpecializedSortedList) Pop(index int32) {
	pos, startIndex := sl._findKth(index)
	sl._delete(pos, startIndex)
}

func (sl *SpecializedSortedList) At(index int32) S {
	if index < 0 || index >= sl.size {
		return nil
	}
	pos, startIndex := sl._findKth(index)
	return sl.blocks[pos][startIndex]
}

func (sl *SpecializedSortedList) BisectRightByK(k int) int32 {
	pos, index := sl._locRightByK(k)
	return sl._queryTree(pos) + index
}

// 返回一个迭代器，指向键值> key的第一个元素.
// UpperBoundByK.
func (sl *SpecializedSortedList) BisectRightByKForIterator(k int) *Iterator {
	pos, index := sl._locRightByK(k)
	return &Iterator{sl: sl, pos: pos, index: index}
}

func (sl *SpecializedSortedList) BisectLeftByPair(a, b int) int32 {
	pos, index := sl._locLeftByPair(a, b)
	return sl._queryTree(pos) + index
}

func (sl *SpecializedSortedList) BisectLeftByPairForValue(a, b int) S {
	pos, index := sl._locLeftByPair(a, b)
	return sl.blocks[pos][index]
}

func (sl *SpecializedSortedList) Clear() {
	sl.size = 0
	sl.blocks = sl.blocks[:0]
	sl.mins = sl.mins[:0]
	sl.tree = sl.tree[:0]
	sl.shouldRebuildTree = true
}

func (sl *SpecializedSortedList) Len() int32 {
	return sl.size
}

func (sl *SpecializedSortedList) _delete(pos, index int32) {
	// !delete element
	sl.size--
	sl._updateTree(pos, -1)
	sl.blocks[pos] = append(sl.blocks[pos][:index], sl.blocks[pos][index+1:]...)
	if len(sl.blocks[pos]) > 0 {
		sl.mins[pos] = sl.blocks[pos][0]
		return
	}

	// !delete block
	sl.blocks = append(sl.blocks[:pos], sl.blocks[pos+1:]...)
	sl.mins = append(sl.mins[:pos], sl.mins[pos+1:]...)
	sl.shouldRebuildTree = true
}

func (sl *SpecializedSortedList) _locLeftByPair(a, b int) (pos, index int32) {
	if sl.size == 0 {
		return
	}

	// find pos
	left := int32(-1)
	right := int32(len(sl.blocks) - 1)
	for left+1 < right {
		mid := (left + right) >> 1
		if !lessPair(sl.mins[mid].p1, sl.mins[mid].p2, a, b) {
			right = mid
		} else {
			left = mid
		}
	}
	if right > 0 {
		block := sl.blocks[right-1]
		last := block[len(block)-1]
		if !lessPair(last.p1, last.p2, a, b) {
			right--
		}
	}
	pos = right

	// find index
	cur := sl.blocks[pos]
	left = -1
	right = int32(len(cur))
	for left+1 < right {
		mid := (left + right) >> 1
		if !lessPair(cur[mid].p1, cur[mid].p2, a, b) {
			right = mid
		} else {
			left = mid
		}
	}

	index = right
	return
}

func (sl *SpecializedSortedList) _locRight(value S) (pos, index int32) {
	if sl.size == 0 {
		return
	}

	// find pos
	left := int32(0)
	right := int32(len(sl.blocks))
	for left+1 < right {
		mid := (left + right) >> 1
		if sl.less(value, sl.mins[mid]) {
			right = mid
		} else {
			left = mid
		}
	}
	pos = left

	// find index
	cur := sl.blocks[pos]
	left = -1
	right = int32(len(cur))
	for left+1 < right {
		mid := (left + right) >> 1
		if sl.less(value, cur[mid]) {
			right = mid
		} else {
			left = mid
		}
	}

	index = right
	return
}

func (sl *SpecializedSortedList) _locRightByK(k int) (pos, index int32) {
	if sl.size == 0 {
		return
	}

	// find pos
	left := int32(0)
	right := int32(len(sl.blocks))
	for left+1 < right {
		mid := (left + right) >> 1
		if k < sl.mins[mid].k {
			right = mid
		} else {
			left = mid
		}
	}
	pos = left

	// find index
	cur := sl.blocks[pos]
	left = -1
	right = int32(len(cur))
	for left+1 < right {
		mid := (left + right) >> 1
		if k < cur[mid].k {
			right = mid
		} else {
			left = mid
		}
	}

	index = right
	return
}

func (sl *SpecializedSortedList) _buildTree() {
	sl.tree = make([]int32, len(sl.blocks))
	for i := 0; i < len(sl.blocks); i++ {
		sl.tree[i] = int32(len(sl.blocks[i]))
	}
	tree := sl.tree
	for i := 0; i < len(tree); i++ {
		j := i | (i + 1)
		if j < len(tree) {
			tree[j] += tree[i]
		}
	}
	sl.shouldRebuildTree = false
}

func (sl *SpecializedSortedList) _updateTree(index, delta int32) {
	if sl.shouldRebuildTree {
		return
	}
	tree := sl.tree
	for i := index; i < int32(len(tree)); i |= i + 1 {
		tree[i] += delta
	}
}

func (sl *SpecializedSortedList) _queryTree(end int32) int32 {
	if sl.shouldRebuildTree {
		sl._buildTree()
	}
	tree := sl.tree
	sum := int32(0)
	for end > 0 {
		sum += tree[end-1]
		end &= end - 1
	}
	return sum
}

func (sl *SpecializedSortedList) _findKth(k int32) (pos, index int32) {
	if k < int32(len(sl.blocks[0])) {
		return 0, k
	}
	last := int32(len(sl.blocks) - 1)
	lastLen := int32(len(sl.blocks[last]))
	if k >= sl.size-lastLen {
		return last, k + lastLen - sl.size
	}
	if sl.shouldRebuildTree {
		sl._buildTree()
	}
	tree := sl.tree
	pos = -1
	m := int32(len(tree))
	bitLength := bits.Len32(uint32(m))
	for d := bitLength - 1; d >= 0; d-- {
		next := pos + (1 << d)
		if next < m && k >= tree[next] {
			pos = next
			k -= tree[pos]
		}
	}
	return pos + 1, k
}

type Iterator struct {
	sl         *SpecializedSortedList
	pos, index int32
}

func (it *Iterator) HasNext() bool {
	b := it.sl.blocks
	m := int32(len(b))
	if it.pos < m-1 {
		return true
	}
	return it.pos == m-1 && it.index < int32(len(b[it.pos]))-1
}

func (it *Iterator) Next() {
	it.index++
	if it.index == int32(len(it.sl.blocks[it.pos])) {
		it.pos++
		it.index = 0
	}
}

func (it *Iterator) HasPrev() bool {
	if it.pos > 0 {
		return true
	}
	return it.pos == 0 && it.index > 0
}

func (it *Iterator) Prev() {
	it.index--
	if it.index == -1 {
		it.pos--
		it.index = int32(len(it.sl.blocks[it.pos]) - 1)
	}
}

// GetMut
func (it *Iterator) Value() S {
	return it.sl.blocks[it.pos][it.index]
}

func (it *Iterator) NextValue() S {
	newPos, newIndex := it.pos, it.index
	newIndex++
	if newIndex == int32(len(it.sl.blocks[it.pos])) {
		newPos++
		newIndex = 0
	}
	return it.sl.blocks[newPos][newIndex]
}

func (it *Iterator) PrevValue() S {
	newPos, newIndex := it.pos, it.index
	newIndex--
	if newIndex == -1 {
		newPos--
		newIndex = int32(len(it.sl.blocks[newPos]) - 1)
	}
	return it.sl.blocks[newPos][newIndex]
}

func (it *Iterator) ToIndex() int32 {
	res := it.sl._queryTree(it.pos)
	return res + it.index
}

func (it *Iterator) Copy() *Iterator {
	return &Iterator{sl: it.sl, pos: it.pos, index: it.index}
}

func (it *Iterator) IsBegin() bool {
	return it.pos == 0 && it.index == 0
}

func (it *Iterator) IsEnd() bool {
	m := int32(len(it.sl.blocks))
	return it.pos == m && it.index == 0
}
func min32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}
