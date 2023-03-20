// 注意cpp里的迭代器:
// !Begin指向第一个元素,
// !End指向最后一个元素的下一个位置,
// 这里的迭代器设计为:
// !Begin指向第一个元素的前一个位置,First指向第一个元素
// !Last指向最后一个元素,End指向最后一个元素的下一个位置

// 注意插入和删除都可能导致迭代器失效.
package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"sort"
	"strings"
)

func linAddGetMin() {
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

func demo() {
	// lc := NewLineContainer(true)
	// lc.Add(1, 1)
	// lc.Add(2, 2)
	// lc.Add(3, 3)
	// fmt.Println(lc.Query(1))
	// fmt.Println(lc.Query(1))
	// fmt.Println(lc.Query(2))
	// fmt.Println(lc.Query(3))
	// fmt.Println(lc.Query(4))
	// fmt.Println(lc.Query(5))
	// lc.Add(1, 2)
	// fmt.Println(lc.Query(1))
	// fmt.Println(lc.Query(2))
	// fmt.Println(lc.Query(3))
	// lc.Add(0, 0)
	// fmt.Println(lc.Query(1))
	// check with bf
	L := NewLineContainer(true)
	lines := make([]*Line, 0)
	for i := 0; i < 10; i++ {
		k := rand.Intn(10)
		m := rand.Intn(10)
		L.Add(k, m)
		lines = append(lines, &Line{k: k, m: m})

	}
	bf := func(x int, lines []*Line) int {
		res := INF
		for _, line := range lines {
			res = min(res, line.k*x+line.m)
		}
		return res
	}

	for i := 0; i < 100; i++ {
		x := rand.Intn(100)
		v1 := L.Query(x)
		v2 := bf(x, lines)
		if v1 != v2 {
			fmt.Println(v1, v2, x)
			for _, line := range lines {
				fmt.Println(line.k, line.m)
			}
			panic("error")
		}
	}

}

func main() {
	demo()
}

const INF int = 2e18

type Line struct {
	k, m, p int
}

type LineContainer struct {
	minimize bool
	sl       *SortedList
}

func NewLineContainer(minimize bool) *LineContainer {
	return &LineContainer{
		minimize: minimize,
		sl:       NewSortedList(func(a, b S) bool { return a.k < b.k }, nil),
	}
}

// 向集合中添加一条线，表示为y = kx + m
func (lc *LineContainer) Add(k, m int) {
	if lc.minimize {
		k, m = -k, -m
	}

	newLine := &Line{k: k, m: m}
	it1 := lc.sl.Insert(newLine)
	it2 := it1 // copy
	it1.Next()
	it3 := it2
	for lc.insect(&it2, &it1) {
		fmt.Println(2) // TODO
		it1 = lc.sl.Erase(it1)
	}

	if !it3.IsFirst() {
		it3.Prev()
		if lc.insect(&it3, &it2) {
			it2 = lc.sl.Erase(it2)
			lc.insect(&it3, &it2)
		}
	}

	for {
		it2 = it3
		if it2.IsFirst() {
			break
		}

		it3.Prev()
		if it3.Value().p >= it2.Value().p {
			it2 := lc.sl.Erase(it2)
			lc.insect(&it3, &it2)
		}
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

	it := lc.sl.LowerBoundWithP(x, func(p1, p2 int) bool { return p1 < p2 })
	line := it.Value()
	v := line.k*x + line.m
	if lc.minimize {
		return -v
	}
	return v
}

// 这个函数在向集合添加新线或删除旧线时用于计算交点。
// 计算线性函数x和y的交点，并将结果存储在x->p中。
func (lc *LineContainer) insect(x, y *Iterator) bool {
	if y.IsEnd() {
		x.Value().p = INF
		return false
	}

	line1, line2 := x.Value(), y.Value()
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
func (sl *SortedList) LowerBoundWithP(p int, less func(p1, p2 int) bool) Iterator {
	for i, block := range sl.blocks {
		if !less(block[len(block)-1].p, p) {
			pos := sort.Search(len(block), func(i int) bool {
				return !less(block[i].p, p)
			})
			return sl.IteratorAt(i, pos)
		}
	}
	res := sl.Iterator()
	res.End()
	return res
}

type S = *Line

// SortedListSQRT
type SortedList struct {
	less   func(a, b S) bool
	size   int
	blocks [][]S
}

func NewSortedList(less func(a, b S) bool, items []S) *SortedList {
	res := &SortedList{less: less}
	if len(items) > 0 {
		items = append(items[:0:0], items...)
		sort.Slice(items, func(i, j int) bool {
			return less(items[i], items[j])
		})
		res.blocks = res._initBlocks(items)
	}
	res.size = len(items)
	return res
}

// 50/170
const (
	_BLOCK_RATIO   = 25
	_REBUILD_RATIO = 70
)

type _positon byte

const (
	_begin, _between, _end = 0, 1, 2
)

func (sl *SortedList) Add(value S) {
	if sl.size == 0 {
		sl.blocks = append(sl.blocks[:0], []S{value})
		sl.size = 1
		return
	}

	hitIndex := sl._findBlockIndex(value)
	if hitIndex == -1 {
		sl.blocks[len(sl.blocks)-1] = append(sl.blocks[len(sl.blocks)-1], value)
		sl.size++
		if len(sl.blocks[len(sl.blocks)-1]) > _REBUILD_RATIO*len(sl.blocks) {
			sl.rebuild()
		}
		return
	}

	hitted := sl.blocks[hitIndex]
	pos := sl._bisectRight(hitted, value) // insort
	sl.blocks[hitIndex] = append(hitted[:pos], append([]S{value}, hitted[pos:]...)...)
	sl.size++
	if len(hitted) > _REBUILD_RATIO*len(sl.blocks) {
		sl.rebuild()
	}
}

func (sl *SortedList) Has(value S) bool {
	if sl.size == 0 {
		return false
	}
	hitIndex := sl._findBlockIndex(value)
	if hitIndex == -1 {
		return false
	}
	hitted := sl.blocks[hitIndex]
	pos := sl._bisectLeft(hitted, value)
	return pos < len(hitted) && hitted[pos] == value
}

func (sl *SortedList) Discard(value S) bool {
	if sl.size == 0 {
		return false
	}
	hitIndex := sl._findBlockIndex(value)
	if hitIndex == -1 {
		return false
	}
	hitted := sl.blocks[hitIndex]
	pos := sl._bisectLeft(hitted, value)
	if pos == len(hitted) || hitted[pos] != value {
		return false
	}
	sl.blocks[hitIndex] = append(hitted[:pos], hitted[pos+1:]...)
	sl.size--
	if len(sl.blocks[hitIndex]) == 0 {
		// !Splice When Empty, Do Not Rebuild
		sl.blocks = append(sl.blocks[:hitIndex], sl.blocks[hitIndex+1:]...)
	}
	return true
}

func (sl *SortedList) Pop(index int) S {
	if index < 0 {
		index += sl.size
	}
	if index < 0 || index >= sl.size {
		panic("index out of range")
	}
	for i, block := range sl.blocks {
		if index < len(block) {
			res := block[index]
			sl.blocks[i] = append(block[:index], block[index+1:]...)
			sl.size--
			if len(sl.blocks[i]) == 0 {
				// !Splice When Empty, Do Not Rebuild
				sl.blocks = append(sl.blocks[:i], sl.blocks[i+1:]...)
			}
			return res
		}
		index -= len(block)
	}
	panic("impossible")
}

func (sl *SortedList) At(index int) S {
	if index < 0 {
		index += sl.size
	}
	if index < 0 || index >= sl.size {
		panic("index out of range")
	}
	for _, block := range sl.blocks {
		if index < len(block) {
			return block[index]
		}
		index -= len(block)
	}
	panic("impossible")
}

// Count the number of elements < value or
// returns the index of the first element >= value.
func (sl *SortedList) BisectLeft(value S) int {
	res := 0
	for _, block := range sl.blocks {
		if !sl.less(block[len(block)-1], value) {
			return res + sl._bisectLeft(block, value)
		}
		res += len(block)
	}
	return res
}

// Count the number of elements <= value or
// returns the index of the first element > value.
func (sl *SortedList) BisectRight(value S) int {
	res := 0
	for _, block := range sl.blocks {
		if sl.less(value, block[len(block)-1]) {
			return res + sl._bisectRight(block, value)
		}
		res += len(block)
	}
	return res
}

func (sl *SortedList) Clear() {
	sl.blocks = sl.blocks[:0]
	sl.size = 0
}

func (sl *SortedList) Lower(value S) (res S, ok bool) {
	for i := len(sl.blocks) - 1; i >= 0; i-- {
		block := sl.blocks[i]
		if sl.less(block[0], value) {
			pos := sl._bisectLeft(block, value)
			return block[pos-1], true
		}
	}
	return
}

func (sl *SortedList) Higher(value S) (res S, ok bool) {
	for _, block := range sl.blocks {
		if sl.less(value, block[len(block)-1]) {
			pos := sl._bisectRight(block, value)
			return block[pos], true
		}
	}
	return
}

func (sl *SortedList) Floor(value S) (res S, ok bool) {
	for i := len(sl.blocks) - 1; i >= 0; i-- {
		block := sl.blocks[i]
		if !sl.less(value, block[0]) {
			pos := sl._bisectRight(block, value)
			return block[pos-1], true
		}
	}
	return
}

func (sl *SortedList) Ceiling(value S) (res S, ok bool) {
	for _, block := range sl.blocks {
		if !sl.less(block[len(block)-1], value) {
			pos := sl._bisectLeft(block, value)
			return block[pos], true
		}
	}
	return
}

func (sl *SortedList) ForEach(f func(value S, index int)) {
	pos := 0
	for _, block := range sl.blocks {
		for _, value := range block {
			f(value, pos)
			pos++
		}
	}
}

func (sl *SortedList) Len() int {
	return sl.size
}

func (sl *SortedList) String() string {
	res := make([]string, 0)
	sl.ForEach(func(value S, _ int) {
		res = append(res, fmt.Sprintf("%v", value))
	})
	return fmt.Sprintf("SortedList{%v}", strings.Join(res, ", "))
}

func (sl *SortedList) rebuild() {
	if sl.size == 0 {
		return
	}
	bc := int(math.Ceil(math.Sqrt(float64(sl.size) / _BLOCK_RATIO)))
	bs := (sl.size + bc - 1) / bc
	newB := make([][]S, bc)
	ptr := 0
	for i := 0; i < len(sl.blocks); i++ {
		b := sl.blocks[i]
		for j := 0; j < len(b); j++ {
			tmp := ptr / bs
			newB[tmp] = append(newB[tmp], b[j])
			ptr++
		}
	}
	sl.blocks = newB
}

func (sl *SortedList) _initBlocks(sorted []S) [][]S {
	bc := int(math.Ceil(math.Sqrt(float64(len(sorted)) / _BLOCK_RATIO)))
	bs := (len(sorted) + bc - 1) / bc
	res := make([][]S, bc)
	for i := 0; i < bc; i++ {
		res[i] = append(res[i], sorted[i*bs:min((i+1)*bs, len(sorted))]...)
	}
	return res
}

func (sl *SortedList) _bisectLeft(nums []S, value S) int {
	return sort.Search(len(nums), func(i int) bool {
		return !sl.less(nums[i], value)
	})
}

func (sl *SortedList) _bisectRight(nums []S, value S) int {
	return sort.Search(len(nums), func(i int) bool {
		return sl.less(value, nums[i])
	})
}

// 如果没有找到,返回-1
func (sl *SortedList) _findBlockIndex(x S) int {
	for i, block := range sl.blocks {
		if !sl.less(block[len(block)-1], x) {
			return i
		}
	}
	return -1
}

type Iterator struct {
	blocks   [][]S
	bid      int
	pos      int
	position _positon
}

func (sl *SortedList) Iterator() Iterator {
	return Iterator{blocks: sl.blocks, bid: 0, pos: -1, position: _begin}
}

func (sl *SortedList) IteratorAt(bid, pos int) Iterator {
	return Iterator{blocks: sl.blocks, bid: bid, pos: pos, position: _between}
}

func (it *Iterator) Next() bool {
	if it.position == _end {
		goto end
	}

	if nextPos := it.pos + 1; nextPos >= len(it.blocks[it.bid]) {
		nextBid := it.bid + 1
		if nextBid >= len(it.blocks) {
			goto end
		}
		it.bid = nextBid
		it.pos = 0
		goto between
	} else {
		it.pos = nextPos
		goto between
	}

end:
	it.bid = len(it.blocks)
	it.pos = 0
	it.position = _end
	return false
between:
	it.position = _between
	return true

}

func (it *Iterator) Prev() bool {
	if it.position == _begin {
		goto begin
	}

	if prevPos := it.pos - 1; prevPos < 0 {
		prevBid := it.bid - 1
		if prevBid < 0 {
			goto begin
		}
		it.bid = prevBid
		it.pos = len(it.blocks[it.bid]) - 1
		goto between
	} else {
		it.pos = prevPos
		goto between
	}

begin:
	it.bid = 0
	it.pos = -1
	it.position = _begin
	return false
between:
	it.position = _between
	return true
}

func (it *Iterator) Value() S {
	return it.blocks[it.bid][it.pos]
}

func (it *Iterator) Begin() {
	it.bid = 0
	it.pos = -1
	it.position = _begin
}

func (it *Iterator) End() {
	it.bid = len(it.blocks)
	it.pos = 0
	it.position = _end
}

func (it *Iterator) First() bool {
	it.Begin()
	return it.Next()
}

func (it *Iterator) Last() bool {
	it.End()
	return it.Prev()
}

func (it *Iterator) IsBegin() bool {
	return it.position == _begin
}

func (it *Iterator) IsEnd() bool {
	return it.position == _end
}

func (it *Iterator) IsFirst() bool {
	return it.bid == 0 && it.pos == 0
}

func (it *Iterator) IsLast() bool {
	return it.bid == len(it.blocks)-1 && it.pos == len(it.blocks[it.bid])-1
}

// 返回删除元素的后继元素的迭代器，如果删除的是最后一个元素，则返回end()迭代器。
func (sl *SortedList) Erase(it Iterator) Iterator {
	if it.position != _between {
		return it
	}
	value := it.Value()
	sl.Discard(value)
	return it
}

func (sl *SortedList) Insert(value S) Iterator {
	if sl.size == 0 {
		sl.blocks = append(sl.blocks[:0], []S{value})
		sl.size = 1
		return sl.IteratorAt(0, 0)
	}

	hitIndex := sl._findBlockIndex(value)
	if hitIndex == -1 {
		sl.blocks[len(sl.blocks)-1] = append(sl.blocks[len(sl.blocks)-1], value)
		sl.size++
		if len(sl.blocks[len(sl.blocks)-1]) > _REBUILD_RATIO*len(sl.blocks) {
			sl.rebuild()
		}
		return sl.IteratorAt(len(sl.blocks)-1, len(sl.blocks[len(sl.blocks)-1])-1)
	}

	hitted := sl.blocks[hitIndex]
	pos := sl._bisectLeft(hitted, value)
	sl.blocks[hitIndex] = append(hitted[:pos], append([]S{value}, hitted[pos:]...)...)
	sl.size++
	if len(hitted) > _REBUILD_RATIO*len(sl.blocks) {
		sl.rebuild()
	}
	return sl.IteratorAt(hitIndex, pos)
}

// 返回一个迭代器，指向键值>= value的第一个元素。
func (sl *SortedList) LowerBound(value S) (Iterator, bool) {
	for i, block := range sl.blocks {
		if !sl.less(block[len(block)-1], value) {
			pos := sl._bisectLeft(block, value)
			return sl.IteratorAt(i, pos), true
		}
	}
	return sl.Iterator(), false
}

// 返回一个迭代器，指向键值> value的第一个元素。
func (sl *SortedList) UpperBound(value S) (Iterator, bool) {
	for i, block := range sl.blocks {
		if sl.less(value, block[len(block)-1]) {
			pos := sl._bisectRight(block, value)
			return sl.IteratorAt(i, pos), true
		}
	}
	return sl.Iterator(), false
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
