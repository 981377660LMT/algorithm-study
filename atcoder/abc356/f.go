package main

import (
	"bufio"
	"fmt"
	stdio "io"
	"math/bits"
	"os"
	"sort"
	"strconv"
	"strings"
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

// 整数
// K が与えられます。はじめ空である集合
// S に対して、次の
// 2 種類のクエリを順に
// Q 個処理してください。

// 1 x：整数
// x が与えられる。
// S に
// x が含まれているならば、
// S から
// x を取り除く。そうでないならば、
// S に
// x を追加する。
// 2 x：
// S に含まれる整数
// x が与えられる。
// S に含まれる数を頂点とし、差の絶対値が
// K 以下であるような数の間に辺を張ったグラフにおいて、
// x が属する連結成分の頂点数を出力する。
func main() {
	in := os.Stdin
	out := os.Stdout
	io = NewIost(in, out)
	defer func() {
		io.Writer.Flush()
	}()

	q, k := int32(io.NextInt()), io.NextInt()
	sl := NewSortedList32(func(a, b S) bool { return a < b })
	lct := NewLinkCutTreeSubTree(false)
	set := make(map[int]struct{})
	queries := make([][2]int, q)
	for i := int32(0); i < q; i++ {
		op, x := io.NextInt(), io.NextInt()
		queries[i] = [2]int{op, x}
		set[x] = struct{}{}
	}
	nodes := lct.Build(int32(len(set))+10, func(i int32) E { return 1 })
	allNums := make([]int, 0, len(set))
	for x := range set {
		allNums = append(allNums, x)
	}
	sort.Ints(allNums)
	rank := make(map[int]int32, len(allNums))
	for i, x := range allNums {
		rank[x] = int32(i)
	}

	add := func(x int) {
		lower, ok1 := sl.Lower(x)
		ok1 = ok1 && x-lower <= k
		higher, ok2 := sl.Higher(x)
		ok2 = ok2 && higher-x <= k
		if ok1 && ok2 && higher-lower <= k {
			lct.CutEdge(nodes[rank[lower]], nodes[rank[higher]])
		}

		if ok1 {
			lct.LinkEdge(nodes[rank[x]], nodes[rank[lower]])
		}
		if ok2 {
			lct.LinkEdge(nodes[rank[x]], nodes[rank[higher]])
		}

		sl.Add(x)
	}

	remove := func(x int) {
		lower, ok1 := sl.Lower(x)
		ok1 = ok1 && x-lower <= k
		higher, ok2 := sl.Higher(x)
		ok2 = ok2 && higher-x <= k
		if ok1 {
			lct.CutEdge(nodes[rank[x]], nodes[rank[lower]])
		}

		if ok2 {
			lct.CutEdge(nodes[rank[x]], nodes[rank[higher]])
		}

		if ok1 && ok2 && higher-lower <= k {
			lct.LinkEdge(nodes[rank[lower]], nodes[rank[higher]])
		}

		sl.Discard(x)
	}

	querySize := func(x int) int32 {
		node := nodes[rank[x]]
		lct.Evert(node)
		return lct.QuerySubTree(node)
	}

	for _, query := range queries {
		op, x := query[0], query[1]
		if op == 1 {
			if sl.Has(x) {
				remove(x)
			} else {
				add(x)
			}
		} else {
			io.Println(querySize(x))
		}
	}
}

// 1e5 -> 200, 2e5 -> 400
const _LOAD int32 = 200

type S = int

var EMPTY S

// 使用分块+树状数组维护的有序序列.
type SortedList32 struct {
	less              func(a, b S) bool
	size              int32
	blocks            [][]S
	mins              []S
	tree              []int32
	shouldRebuildTree bool
}

func NewSortedList32(less func(a, b S) bool, elements ...S) *SortedList32 {
	elements = append(elements[:0:0], elements...)
	res := &SortedList32{less: less}
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

func (sl *SortedList32) Add(value S) *SortedList32 {
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
		sl.blocks = append(sl.blocks, nil)
		copy(sl.blocks[pos+2:], sl.blocks[pos+1:])
		sl.blocks[pos+1] = sl.blocks[pos][_LOAD:]
		sl.blocks[pos] = sl.blocks[pos][:_LOAD:_LOAD]
		sl.mins = append(sl.mins, EMPTY)
		copy(sl.mins[pos+2:], sl.mins[pos+1:])
		sl.mins[pos+1] = sl.blocks[pos+1][0]
		sl.shouldRebuildTree = true
	}

	return sl
}

func (sl *SortedList32) Has(value S) bool {
	if len(sl.blocks) == 0 {
		return false
	}
	pos, index := sl._locLeft(value)
	return index < int32(len(sl.blocks[pos])) && sl.blocks[pos][index] == value
}

func (sl *SortedList32) Discard(value S) bool {
	if len(sl.blocks) == 0 {
		return false
	}
	pos, index := sl._locRight(value)
	if index > 0 && sl.blocks[pos][index-1] == value {
		sl._delete(pos, index-1)
		return true
	}
	return false
}

func (sl *SortedList32) Pop(index int32) S {
	if index < 0 {
		index += sl.size
	}
	if index < 0 || index >= sl.size {
		panic("index out of range")
	}
	pos, startIndex := sl._findKth(index)
	value := sl.blocks[pos][startIndex]
	sl._delete(pos, startIndex)
	return value
}

func (sl *SortedList32) At(index int32) S {
	if index < 0 {
		index += sl.size
	}
	if index < 0 || index >= sl.size {
		panic("index out of range")
	}
	pos, startIndex := sl._findKth(index)
	return sl.blocks[pos][startIndex]
}

func (sl *SortedList32) Erase(start, end int32) {
	sl.Enumerate(start, end, nil, true)
}

func (sl *SortedList32) Lower(value S) (res S, ok bool) {
	pos := sl.BisectLeft(value)
	if pos == 0 {
		return
	}
	return sl.At(pos - 1), true
}

func (sl *SortedList32) Higher(value S) (res S, ok bool) {
	pos := sl.BisectRight(value)
	if pos == sl.size {
		return
	}
	return sl.At(pos), true
}

func (sl *SortedList32) Floor(value S) (res S, ok bool) {
	pos := sl.BisectRight(value)
	if pos == 0 {
		return
	}
	return sl.At(pos - 1), true
}

func (sl *SortedList32) Ceiling(value S) (res S, ok bool) {
	pos := sl.BisectLeft(value)
	if pos == sl.size {
		return
	}
	return sl.At(pos), true
}

// 返回第一个大于等于 `value` 的元素的索引/严格小于 `value` 的元素的个数.
func (sl *SortedList32) BisectLeft(value S) int32 {
	pos, index := sl._locLeft(value)
	return sl._queryTree(pos) + index
}

// 返回第一个严格大于 `value` 的元素的索引/小于等于 `value` 的元素的个数.
func (sl *SortedList32) BisectRight(value S) int32 {
	pos, index := sl._locRight(value)
	return sl._queryTree(pos) + index
}

func (sl *SortedList32) Count(value S) int32 {
	return sl.BisectRight(value) - sl.BisectLeft(value)
}

func (sl *SortedList32) Clear() {
	sl.size = 0
	sl.blocks = sl.blocks[:0]
	sl.mins = sl.mins[:0]
	sl.tree = sl.tree[:0]
	sl.shouldRebuildTree = true
}

func (sl *SortedList32) ForEach(f func(value S, index int32) bool, reverse bool) {
	if !reverse {
		count := int32(0)
		for i := 0; i < len(sl.blocks); i++ {
			block := sl.blocks[i]
			for j := 0; j < len(block); j++ {
				if f(block[j], count) {
					return
				}
				count++
			}
		}
		return
	}
	count := int32(0)
	for i := len(sl.blocks) - 1; i >= 0; i-- {
		block := sl.blocks[i]
		for j := len(block) - 1; j >= 0; j-- {
			if f(block[j], count) {
				return
			}
			count++
		}
	}
}

func (sl *SortedList32) Enumerate(start, end int32, f func(value S), erase bool) {
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

func (sl *SortedList32) Slice(start, end int32) []S {
	if start < 0 {
		start = 0
	}
	if end > sl.size {
		end = sl.size
	}
	if start >= end {
		return nil
	}
	count := end - start
	res := make([]S, 0, count)
	pos, index := sl._findKth(start)
	m := int32(len(sl.blocks))
	for ; count > 0 && pos < m; pos++ {
		block := sl.blocks[pos]
		endPos := min32(int32(len(block)), index+count)
		curCount := endPos - index
		res = append(res, block[index:endPos]...)
		count -= curCount
		index = 0
	}
	return res
}

func (sl *SortedList32) Range(min, max S) []S {
	if sl.less(max, min) {
		return nil
	}
	res := []S{}
	pos := sl._locBlock(min)
	m := int32(len(sl.blocks))
	for i := pos; i < m; i++ {
		block := sl.blocks[i]
		for j := 0; j < len(block); j++ {
			x := block[j]
			if sl.less(max, x) {
				return res
			}
			if !sl.less(x, min) {
				res = append(res, x)
			}
		}
	}
	return res
}

func (sl *SortedList32) IteratorAt(index int32) *Iterator {
	if index < 0 {
		index += sl.size
	}
	if index < 0 || index >= sl.size {
		panic("Index out of range")
	}
	pos, startIndex := sl._findKth(index)
	return sl._iteratorAt(pos, startIndex)
}

func (sl *SortedList32) LowerBound(value S) *Iterator {
	pos, index := sl._locLeft(value)
	return sl._iteratorAt(pos, index)
}

func (sl *SortedList32) UpperBound(value S) *Iterator {
	pos, index := sl._locRight(value)
	return sl._iteratorAt(pos, index)
}

func (sl *SortedList32) Min() S {
	if sl.size == 0 {
		panic("Min() called on empty SortedList")
	}
	return sl.mins[0]
}

func (sl *SortedList32) Max() S {
	if sl.size == 0 {
		panic("Max() called on empty SortedList")
	}
	lastBlock := sl.blocks[len(sl.blocks)-1]
	return lastBlock[len(lastBlock)-1]
}

func (sl *SortedList32) String() string {
	sb := strings.Builder{}
	sb.WriteString("SortedList{")
	sl.ForEach(func(value S, index int32) bool {
		if index > 0 {
			sb.WriteByte(',')
		}
		sb.WriteString(fmt.Sprintf("%v", value))
		return false
	}, false)
	sb.WriteByte('}')
	return sb.String()
}

func (sl *SortedList32) Len() int32 {
	return sl.size
}

func (sl *SortedList32) _delete(pos, index int32) {
	// !delete element
	sl.size--
	sl._updateTree(pos, -1)
	copy(sl.blocks[pos][index:], sl.blocks[pos][index+1:])
	sl.blocks[pos] = sl.blocks[pos][:len(sl.blocks[pos])-1]
	if len(sl.blocks[pos]) > 0 {
		sl.mins[pos] = sl.blocks[pos][0]
		return
	}

	// !delete block
	copy(sl.blocks[pos:], sl.blocks[pos+1:])
	sl.blocks = sl.blocks[:len(sl.blocks)-1]
	copy(sl.mins[pos:], sl.mins[pos+1:])
	sl.mins = sl.mins[:len(sl.mins)-1]
	sl.shouldRebuildTree = true
}

func (sl *SortedList32) _locLeft(value S) (pos, index int32) {
	if sl.size == 0 {
		return
	}

	// find pos
	left := int32(-1)
	right := int32(len(sl.blocks) - 1)
	for left+1 < right {
		mid := (left + right) >> 1
		if !sl.less(sl.mins[mid], value) {
			right = mid
		} else {
			left = mid
		}
	}
	if right > 0 {
		block := sl.blocks[right-1]
		if !sl.less(block[len(block)-1], value) {
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
		if !sl.less(cur[mid], value) {
			right = mid
		} else {
			left = mid
		}
	}

	index = right
	return
}

func (sl *SortedList32) _locRight(value S) (pos, index int32) {
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

func (sl *SortedList32) _locBlock(value S) int32 {
	left, right := int32(-1), int32(len(sl.blocks)-1)
	for left+1 < right {
		mid := (left + right) >> 1
		if !sl.less(sl.mins[mid], value) {
			right = mid
		} else {
			left = mid
		}
	}
	if right > 0 {
		block := sl.blocks[right-1]
		if !sl.less(block[len(block)-1], value) {
			right--
		}
	}
	return right
}

func (sl *SortedList32) _buildTree() {
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

func (sl *SortedList32) _updateTree(index, delta int32) {
	if sl.shouldRebuildTree {
		return
	}
	tree := sl.tree
	m := int32(len(tree))
	for i := index; i < m; i |= i + 1 {
		tree[i] += delta
	}
}

func (sl *SortedList32) _queryTree(end int32) int32 {
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

func (sl *SortedList32) _findKth(k int32) (pos, index int32) {
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
	bitLength := bits.Len32(uint32(len(tree)))
	for d := bitLength - 1; d >= 0; d-- {
		next := pos + (1 << d)
		if next < int32(len(tree)) && k >= tree[next] {
			pos = next
			k -= tree[pos]
		}
	}
	return pos + 1, k
}

func (sl *SortedList32) _iteratorAt(pos, index int32) *Iterator {
	return &Iterator{sl: sl, pos: pos, index: index}
}

type Iterator struct {
	sl    *SortedList32
	pos   int32
	index int32
}

func (it *Iterator) HasNext() bool {
	return it.pos < int32(len(it.sl.blocks))-1 || it.index < int32(len(it.sl.blocks[it.pos]))-1
}

func (it *Iterator) Next() (res S, ok bool) {
	if !it.HasNext() {
		return
	}
	it.index++
	if it.index == int32(len(it.sl.blocks[it.pos])) {
		it.pos++
		it.index = 0
	}
	res = it.sl.blocks[it.pos][it.index]
	ok = true
	return
}

func (it *Iterator) HasPrev() bool {
	return it.pos > 0 || it.index > 0
}

func (it *Iterator) Prev() (res S, ok bool) {
	if !it.HasPrev() {
		return
	}
	it.index--
	if it.index == -1 {
		it.pos--
		it.index = int32(len(it.sl.blocks[it.pos]) - 1)
	}
	res = it.sl.blocks[it.pos][it.index]
	ok = true
	return
}

func (it *Iterator) Remove() {
	it.sl._delete(it.pos, it.index)
}

func (it *Iterator) Value() (res S, ok bool) {
	if it.pos < 0 || it.pos >= it.sl.Len() {
		return
	}
	block := it.sl.blocks[it.pos]
	if it.index < 0 || it.index >= int32(len(block)) {
		return
	}
	res = block[it.index]
	ok = true
	return
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

func min32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func max32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

// hack.如果已知元素是最大/最小的, 可以使用下面的方法.
func (sl *SortedList32) _appendFirst(value S) {
	sl.size++
	if len(sl.blocks) == 0 {
		sl.blocks = append(sl.blocks, []S{value})
		sl.mins = append(sl.mins, value)
		sl.shouldRebuildTree = true
		return
	}
	pos := int32(0)
	sl._updateTree(pos, 1)
	sl.blocks[pos] = append(sl.blocks[pos], EMPTY)
	copy(sl.blocks[pos][1:], sl.blocks[pos])
	sl.blocks[pos][0] = value
	sl._adjust(pos)
	return
}
func (sl *SortedList32) _appendLast(value S) {
	sl.size++
	if len(sl.blocks) == 0 {
		sl.blocks = append(sl.blocks, []S{value})
		sl.mins = append(sl.mins, value)
		sl.shouldRebuildTree = true
		return
	}
	pos := int32(len(sl.blocks) - 1)
	sl._updateTree(pos, 1)
	sl.blocks[pos] = append(sl.blocks[pos], value)
	sl._adjust(pos)
	return
}
func (sl *SortedList32) _popFirst() S {
	pos, startIndex := int32(0), int32(0)
	value := sl.blocks[pos][startIndex]
	sl._delete(pos, startIndex)
	return value
}
func (sl *SortedList32) _popLast() S {
	pos := int32(len(sl.blocks) - 1)
	startIndex := len(sl.blocks[pos]) - 1
	value := sl.blocks[pos][startIndex]
	// !delete element
	sl.size--
	sl._updateTree(pos, -1)
	sl.blocks[pos] = sl.blocks[pos][:len(sl.blocks[pos])-1]
	if len(sl.blocks[pos]) > 0 {
		return value
	}

	// !delete block
	sl.blocks = sl.blocks[:pos]
	sl.mins = sl.mins[:pos]
	sl.shouldRebuildTree = true // TODO: 能否不重建树
	return value
}
func (sl *SortedList32) _adjust(pos int32) {
	// n -> load + (n - load)
	if n := int32(len(sl.blocks[pos])); _LOAD+_LOAD < n {
		sl.blocks = append(sl.blocks, nil)
		copy(sl.blocks[pos+2:], sl.blocks[pos+1:])
		sl.blocks[pos+1] = sl.blocks[pos][_LOAD:]
		sl.blocks[pos] = sl.blocks[pos][:_LOAD:_LOAD]
		sl.mins = append(sl.mins, EMPTY)
		copy(sl.mins[pos+2:], sl.mins[pos+1:])
		sl.mins[pos+1] = sl.blocks[pos+1][0]
		sl.shouldRebuildTree = true
	}
}

type E = int32 // 子树和

func (*treeNode) e() E                { return 0 }
func (*treeNode) op(this, other E) E  { return this + other }
func (*treeNode) inv(this, other E) E { return this - other }

type LinkCutTreeSubTree struct {
	nodeId int32
	edges  map[[2]int32]struct{}
	check  bool
}

// check: AddEdge/RemoveEdge で辺の存在チェックを行うかどうか.
func NewLinkCutTreeSubTree(check bool) *LinkCutTreeSubTree {
	return &LinkCutTreeSubTree{edges: make(map[[2]int32]struct{}), check: check}
}

func (lct *LinkCutTreeSubTree) Build(n int32, f func(i int32) E) []*treeNode {
	nodes := make([]*treeNode, n)
	for i := int32(0); i < n; i++ {
		nodes[i] = lct.Alloc(f(i))
	}
	return nodes
}

// 要素の値を v としたノードを生成する.
func (lct *LinkCutTreeSubTree) Alloc(key E) *treeNode {
	res := newTreeNode(key, lct.nodeId)
	lct.nodeId++
	lct.update(res)
	return res
}

// t を根に変更する.
func (lct *LinkCutTreeSubTree) Evert(t *treeNode) {
	lct.expose(t)
	lct.toggle(t)
	lct.push(t)
}

func (lct *LinkCutTreeSubTree) LinkEdge(child, parent *treeNode) (ok bool) {
	if lct.check {
		if lct.IsConnected(child, parent) {
			return
		}
		id1, id2 := child.id, parent.id
		if id1 > id2 {
			id1, id2 = id2, id1
		}
		tuple := [2]int32{id1, id2}
		lct.edges[tuple] = struct{}{}
	}

	lct.Evert(child)
	lct.expose(parent)
	child.p = parent
	parent.r = child
	lct.update(parent)
	return true
}

func (lct *LinkCutTreeSubTree) CutEdge(u, v *treeNode) (ok bool) {
	if lct.check {
		id1, id2 := u.id, v.id
		if id1 > id2 {
			id1, id2 = id2, id1
		}
		tuple := [2]int32{id1, id2}
		if _, has := lct.edges[tuple]; !has {
			return
		}
		delete(lct.edges, tuple)
	}

	lct.Evert(u)
	lct.expose(v)
	parent := v.l
	v.l = nil
	lct.update(v)
	parent.p = nil
	return true
}

// u と v の lca を返す.
//
//	u と v が異なる連結成分なら nullptr を返す.
//	!上記の操作は根を勝手に変えるため, 事前に Evert する必要があるかも.
func (lct *LinkCutTreeSubTree) QueryLCA(u, v *treeNode) *treeNode {
	if !lct.IsConnected(u, v) {
		return nil
	}
	lct.expose(u)
	return lct.expose(v)
}

func (lct *LinkCutTreeSubTree) KthAncestor(x *treeNode, k int32) *treeNode {
	lct.expose(x)
	for x != nil {
		lct.push(x)
		if x.r != nil && x.r.cnt > k {
			x = x.r
		} else {
			if x.r != nil {
				k -= x.r.cnt
			}
			if k == 0 {
				return x
			}
			k--
			x = x.l
		}
	}
	return nil
}

func (lct *LinkCutTreeSubTree) GetParent(t *treeNode) *treeNode {
	lct.expose(t)
	p := t.l
	if p == nil {
		return nil
	}
	for {
		lct.push(p)
		if p.r == nil {
			return p
		}
		p = p.r
	}
}

func (lct *LinkCutTreeSubTree) Jump(from, to *treeNode, k int32) *treeNode {
	lct.Evert(to)
	lct.expose(from)
	for {
		lct.push(from)
		rs := int32(0)
		if from.r != nil {
			rs = from.r.cnt
		}
		if k < rs {
			from = from.r
			continue
		}
		if k == rs {
			break
		}
		k -= rs + 1
		from = from.l
	}
	lct.splay(from)
	return from
}

func (lct *LinkCutTreeSubTree) Dist(u, v *treeNode) int32 {
	lct.Evert(u)
	lct.expose(v)
	return v.cnt - 1
}

// t を根とする部分木の要素の値の和を返す.
//
//	!Evert を忘れない！
func (lct *LinkCutTreeSubTree) QuerySubTree(t *treeNode) E {
	lct.expose(t)
	return t.op(t.key, t.sub)
}

// t の値を v に変更する.
func (lct *LinkCutTreeSubTree) Set(t *treeNode, key E) *treeNode {
	lct.expose(t)
	t.key = key
	lct.update(t)
	return t
}

// t の値を返す.
func (lct *LinkCutTreeSubTree) Get(t *treeNode) E {
	return t.key
}

// u と v が同じ連結成分に属する場合は true, そうでなければ false を返す.
func (lct *LinkCutTreeSubTree) IsConnected(u, v *treeNode) bool {
	return u == v || lct.GetRoot(u) == lct.GetRoot(v)
}

func (lct *LinkCutTreeSubTree) GetRoot(t *treeNode) *treeNode {
	lct.expose(t)
	for t.l != nil {
		lct.push(t)
		t = t.l
	}
	return t
}

func (lct *LinkCutTreeSubTree) expose(t *treeNode) *treeNode {
	rp := (*treeNode)(nil)
	for cur := t; cur != nil; cur = cur.p {
		lct.splay(cur)
		if cur.r != nil {
			cur.Add(cur.r)
		}
		cur.r = rp
		if cur.r != nil {
			cur.Erase(cur.r)
		}
		lct.update(cur)
		rp = cur
	}
	lct.splay(t)
	return rp
}

func (lct *LinkCutTreeSubTree) update(t *treeNode) {
	t.cnt = 1
	if t.l != nil {
		t.cnt += t.l.cnt
	}
	if t.r != nil {
		t.cnt += t.r.cnt
	}

	t.Merge(t.l, t.r)
}

func (lct *LinkCutTreeSubTree) rotr(t *treeNode) {
	x := t.p
	y := x.p
	x.l = t.r
	if t.r != nil {
		t.r.p = x
	}
	t.r = x
	x.p = t
	lct.update(x)
	lct.update(t)
	t.p = y
	if y != nil {
		if y.l == x {
			y.l = t
		}
		if y.r == x {
			y.r = t
		}
		lct.update(y)
	}
}

func (lct *LinkCutTreeSubTree) rotl(t *treeNode) {
	x := t.p
	y := x.p
	x.r = t.l
	if t.l != nil {
		t.l.p = x
	}
	t.l = x
	x.p = t
	lct.update(x)
	lct.update(t)
	t.p = y
	if y != nil {
		if y.l == x {
			y.l = t
		}
		if y.r == x {
			y.r = t
		}
		lct.update(y)
	}
}

func (lct *LinkCutTreeSubTree) toggle(t *treeNode) {
	t.l, t.r = t.r, t.l
	t.rev = !t.rev
}

func (lct *LinkCutTreeSubTree) push(t *treeNode) {
	if t.rev {
		if t.l != nil {
			lct.toggle(t.l)
		}
		if t.r != nil {
			lct.toggle(t.r)
		}
		t.rev = false
	}
}

func (lct *LinkCutTreeSubTree) splay(t *treeNode) {
	lct.push(t)
	for !t.IsRoot() {
		q := t.p
		if q.IsRoot() {
			lct.push(q)
			lct.push(t)
			if q.l == t {
				lct.rotr(t)
			} else {
				lct.rotl(t)
			}
		} else {
			r := q.p
			lct.push(r)
			lct.push(q)
			lct.push(t)
			if r.l == q {
				if q.l == t {
					lct.rotr(q)
					lct.rotr(t)
				} else {
					lct.rotl(t)
					lct.rotr(t)
				}
			} else {
				if q.r == t {
					lct.rotl(q)
					lct.rotl(t)
				} else {
					lct.rotr(t)
					lct.rotl(t)
				}
			}
		}
	}
}

type treeNode struct {
	key, sum, sub E
	rev           bool
	cnt           int32
	id            int32
	l, r, p       *treeNode
}

func newTreeNode(key E, id int32) *treeNode {
	res := &treeNode{key: key, sum: key, cnt: 1, id: id}
	res.sub = res.e()
	return res
}

func (n *treeNode) IsRoot() bool {
	return n.p == nil || (n.p.l != n && n.p.r != n)
}

func (n *treeNode) Add(other *treeNode)   { n.sub = n.op(n.sub, other.sum) }
func (n *treeNode) Erase(other *treeNode) { n.sub = n.inv(n.sub, other.sum) }
func (n *treeNode) Merge(n1, n2 *treeNode) {
	var tmp1, tmp2 E
	if n1 != nil {
		tmp1 = n1.sum
	} else {
		tmp1 = n.e()
	}

	if n2 != nil {
		tmp2 = n2.sum
	} else {
		tmp2 = n.e()
	}

	n.sum = n.op(n.op(tmp1, n.key), n.op(n.sub, tmp2))
}
