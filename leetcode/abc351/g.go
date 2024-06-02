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

// # 座標平面上に
// # N 個の点
// # P
// # 1
// # ​
// #  ,P
// # 2
// # ​
// #  ,…,P
// # N
// # ​
// #   があり、点
// # P
// # i
// # ​
// #   の座標は
// # (X
// # i
// # ​
// #  ,Y
// # i
// # ​
// #  ) です。
// # 2 つの点
// # A,B の距離
// # dist(A,B) を次のように定義します。

// # 最初、ウサギが点
// # A にいる。
// # (x,y) にいるウサギは
// # (x+1,y+1),
// # (x+1,y−1),
// # (x−1,y+1),
// # (x−1,y−1) のいずれかに
// # 1 回のジャンプで移動することができる。
// # 点
// # A から点
// # B まで移動するために必要なジャンプの回数の最小値を
// # dist(A,B) として定義する。
// # ただし、何度ジャンプを繰り返しても点
// # A から点
// # B まで移動できないとき、
// # dist(A,B)=0 とする。

// # i=1
// # ∑
// # N−1
// # ​

// # j=i+1
// # ∑
// # N
// # ​
// #  dist(P
// # i
// # ​
// #  ,P
// # j
// # ​
// #  ) を求めてください。

// # 每次移动奇偶不变，所以考虑前缀中奇偶相同的点即可
// !切比雪夫距离转曼哈顿距离
func main() {
	in := os.Stdin
	out := os.Stdout
	io = NewIost(in, out)
	defer func() {
		io.Writer.Flush()
	}()

	n := io.NextInt()
	x := make([]int, n)
	y := make([]int, n)
	newX := make([]int, n)
	newY := make([]int, n)
	for i := 0; i < n; i++ {
		x[i] = io.NextInt()
		y[i] = io.NextInt()
		newX[i] = y[i] - x[i]
		newY[i] = x[i] + y[i]
	}

	sl00 := NewSortedListWithSum(func(a, b int) bool { return a < b })
	sl01 := NewSortedListWithSum(func(a, b int) bool { return a < b })
	sl10 := NewSortedListWithSum(func(a, b int) bool { return a < b })
	sl11 := NewSortedListWithSum(func(a, b int) bool { return a < b })
	M00 := NewMedianFinderSortedList(sl00)
	M01 := NewMedianFinderSortedList(sl01)
	M10 := NewMedianFinderSortedList(sl10)
	M11 := NewMedianFinderSortedList(sl11)
	res := 0
	for i := 0; i < n; i++ {
		sum_ := x[i] + y[i]
		x, y := newX[i], newY[i]
		if sum_&1 == 0 {
			res += M00.DistSum(x)
			sl00.Add(x)
			res += M01.DistSum(y)
			sl01.Add(y)
		} else {
			res += M10.DistSum(x)
			sl10.Add(x)
			res += M11.DistSum(y)
			sl11.Add(y)
		}
	}

	io.Println(res / 2)
}

// SortedList 动态维护中位数信息.
// `Proxy 的内部持有一个对 SortedList 的引用`.
type MedianFinderSortedList struct {
	List *SortedListWithSum
}

func NewMedianFinderSortedList(sortedList *SortedListWithSum) *MedianFinderSortedList {
	return &MedianFinderSortedList{List: sortedList}
}

// 如果有两个中位数，返回较小的那个.
func (mfs *MedianFinderSortedList) Median() (res int) {
	return mfs.MedianRange(0, mfs.List.Len())
}

func (mfs *MedianFinderSortedList) MedianRange(start, end int) (res int) {
	if start < 0 {
		start = 0
	}
	if n := mfs.List.Len(); end > n {
		end = n
	}
	if start >= end {
		return 0
	}
	mid := start + (end-start-1)>>1
	return mfs.List.At(mid)
}

// 返回所有数到`to`的距离和.
func (mfs *MedianFinderSortedList) DistSum(to int) int {
	sl := mfs.List
	pos := sl.BisectRight(to)
	allSum := sl.SumAll()
	var sum1, sum2 int
	if pos < sl.Len()>>1 {
		sum1 = sl.SumSlice(0, pos)
		sum2 = allSum - sum1
	} else {
		sum2 = sl.SumSlice(pos, sl.Len())
		sum1 = allSum - sum2
	}
	leftSum := to*pos - sum1
	rightSum := sum2 - to*(sl.Len()-pos)
	return leftSum + rightSum
}

// 返回切片`[start,end)`中所有数到`to`的距离和.
func (mfs *MedianFinderSortedList) DistSumRange(to int, start, end int) int {
	sl := mfs.List
	if start < 0 {
		start = 0
	}
	if end > sl.Len() {
		end = sl.Len()
	}
	if start >= end {
		return 0
	}
	pos := sl.BisectLeft(to)
	if pos <= start {
		return sl.SumSlice(start, end) - to*(end-start)
	}
	if pos >= end {
		return to*(end-start) - sl.SumSlice(start, end)
	}
	leftSum := to*(pos-start) - sl.SumSlice(start, pos)
	rightSum := sl.SumSlice(pos, end) - to*(end-pos)
	return leftSum + rightSum
}

// 返回所有数到中位数的距离和.
func (mfs *MedianFinderSortedList) DistSumToMedian() int {
	if mfs.List.Len() == 0 {
		return 0
	}
	return mfs.DistSum(mfs.Median())
}

// 返回切片`[start,end)`中所有数到中位数的距离和.
func (mfs *MedianFinderSortedList) DistSumToMedianRange(start, end int) int {
	if mfs.List.Len() == 0 {
		return 0
	}
	return mfs.DistSumRange(mfs.Median(), start, end)
}

// 1e5 -> 200, 2e5 -> 400
const _LOAD int = 200

type E = int

func e() E        { return 0 }
func op(a, b E) E { return a + b }
func inv(a E) E   { return -a }

// 使用分块+树状数组维护的有序序列.
type SortedListWithSum struct {
	less              func(a, b E) bool
	size              int
	blocks            [][]E
	mins              []E
	tree              []int
	shouldRebuildTree bool

	sums   []E
	allSum E
}

func NewSortedListWithSum(less func(a, b E) bool, elements ...E) *SortedListWithSum {
	elements = append(elements[:0:0], elements...)
	res := &SortedListWithSum{less: less}
	sort.Slice(elements, func(i, j int) bool { return less(elements[i], elements[j]) })
	n := len(elements)
	blocks := [][]E{}
	sums := []E{}
	allSum := e()
	for start := 0; start < n; start += _LOAD {
		end := min(start+_LOAD, n)
		newBlock := elements[start:end:end] // !各个块互不影响, max参数也需要指定为end
		blocks = append(blocks, newBlock)
		cur := e()
		for _, v := range newBlock {
			cur = op(cur, v)
		}
		sums = append(sums, cur)
		allSum = op(allSum, cur)
	}
	mins := make([]E, len(blocks))
	for i, cur := range blocks {
		mins[i] = cur[0]
	}
	res.size = n
	res.blocks = blocks
	res.mins = mins
	res.shouldRebuildTree = true
	res.sums = sums
	res.allSum = allSum
	return res
}

// 返回区间`[start, end)`的和.
func (sl *SortedListWithSum) SumSlice(start, end int) E {
	if start < 0 {
		start += sl.size
	}
	if start < 0 {
		start = 0
	}
	if end < 0 {
		end += sl.size
	}
	if end > sl.size {
		end = sl.size
	}
	if start >= end {
		return e()
	}

	bid1, startIndex1 := sl._findKth(start)
	bid2, startIndex2 := sl._findKth(end)
	start, end = startIndex1, startIndex2
	res := e()
	if bid1 == bid2 {
		block := sl.blocks[bid1]
		for i := start; i < end; i++ {
			res = op(res, block[i])
		}
	} else {
		if start < len(sl.blocks[bid1]) {
			block1 := sl.blocks[bid1]
			for i := start; i < len(block1); i++ {
				res = op(res, block1[i])
			}
		}
		for i := bid1 + 1; i < bid2; i++ {
			res = op(res, sl.sums[i])
		}
		if m := len(sl.blocks); bid2 < m && end > 0 {
			block2 := sl.blocks[bid2]
			tmp := e()
			for i := 0; i < end; i++ {
				tmp = op(tmp, block2[i])
			}
			res = op(res, tmp)
		}
	}
	return res
}

// 返回范围`[min, max]`的和.
func (sl *SortedListWithSum) SumRange(min, max E) E {
	if sl.less(max, min) {
		return e()
	}
	res := e()
	pos, start := sl._locLeft(min)
	for i := pos; i < len(sl.blocks); i++ {
		block := sl.blocks[i]
		if sl.less(max, block[0]) {
			break
		}
		if start == 0 && !sl.less(block[len(block)-1], max) {
			res = op(res, sl.sums[i])
		} else {
			for j := start; j < len(block); j++ {
				cur := block[j]
				if sl.less(max, cur) {
					break
				}
				res = op(res, cur)
			}
		}
		start = 0
	}
	return res
}

func (sl *SortedListWithSum) SumAll() E {
	return sl.allSum
}

func (sl *SortedListWithSum) Add(value E) *SortedListWithSum {
	sl.size++
	sl.allSum = op(sl.allSum, value)
	if len(sl.blocks) == 0 {
		sl.blocks = append(sl.blocks, []E{value})
		sl.mins = append(sl.mins, value)
		sl.shouldRebuildTree = true
		sl.sums = append(sl.sums, value)
		return sl
	}

	pos, index := sl._locRight(value)
	sl._updateTree(pos, 1)
	sl.blocks[pos] = append(sl.blocks[pos][:index], append([]E{value}, sl.blocks[pos][index:]...)...)
	sl.mins[pos] = sl.blocks[pos][0]
	sl.sums[pos] = op(sl.sums[pos], value)

	// n -> load + (n - load)
	if n := len(sl.blocks[pos]); _LOAD+_LOAD < n {
		oldSum := sl.sums[pos]
		left, right := make([]E, _LOAD), make([]E, n-_LOAD)
		copy(left, sl.blocks[pos][:_LOAD])
		copy(right, sl.blocks[pos][_LOAD:])
		sl.blocks = append(sl.blocks[:pos], append([][]E{left, right}, sl.blocks[pos+1:]...)...)
		sl.mins = append(sl.mins[:pos], append([]E{left[0], right[0]}, sl.mins[pos+1:]...)...)
		sl.shouldRebuildTree = true

		sl._rebuildSum(pos)
		newSum := op(oldSum, inv(sl.sums[pos]))
		sl.sums = append(sl.sums[:pos+1], append([]E{newSum}, sl.sums[pos+1:]...)...)
	}

	return sl
}

func (sl *SortedListWithSum) Has(value E) bool {
	if len(sl.blocks) == 0 {
		return false
	}
	pos, index := sl._locLeft(value)
	return index < len(sl.blocks[pos]) && sl.blocks[pos][index] == value
}

func (sl *SortedListWithSum) Discard(value E) bool {
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

func (sl *SortedListWithSum) Pop(index int) E {
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

func (sl *SortedListWithSum) At(index int) E {
	if index < 0 {
		index += sl.size
	}
	if index < 0 || index >= sl.size {
		panic("index out of range")
	}
	pos, startIndex := sl._findKth(index)
	return sl.blocks[pos][startIndex]
}

func (sl *SortedListWithSum) Erase(start, end int) {
	sl.Enumerate(start, end, nil, true)
}

func (sl *SortedListWithSum) Lower(value E) (res E, ok bool) {
	pos := sl.BisectLeft(value)
	if pos == 0 {
		return
	}
	return sl.At(pos - 1), true
}

func (sl *SortedListWithSum) Higher(value E) (res E, ok bool) {
	pos := sl.BisectRight(value)
	if pos == sl.size {
		return
	}
	return sl.At(pos), true
}

func (sl *SortedListWithSum) Floor(value E) (res E, ok bool) {
	pos := sl.BisectRight(value)
	if pos == 0 {
		return
	}
	return sl.At(pos - 1), true
}

func (sl *SortedListWithSum) Ceiling(value E) (res E, ok bool) {
	pos := sl.BisectLeft(value)
	if pos == sl.size {
		return
	}
	return sl.At(pos), true
}

// 返回第一个大于等于 `value` 的元素的索引/严格小于 `value` 的元素的个数.
func (sl *SortedListWithSum) BisectLeft(value E) int {
	pos, index := sl._locLeft(value)
	return sl._queryTree(pos) + index
}

// 返回第一个严格大于 `value` 的元素的索引/小于等于 `value` 的元素的个数.
func (sl *SortedListWithSum) BisectRight(value E) int {
	pos, index := sl._locRight(value)
	return sl._queryTree(pos) + index
}

func (sl *SortedListWithSum) Count(value E) int {
	return sl.BisectRight(value) - sl.BisectLeft(value)
}

func (sl *SortedListWithSum) Clear() {
	sl.size = 0
	sl.blocks = sl.blocks[:0]
	sl.mins = sl.mins[:0]
	sl.tree = sl.tree[:0]
	sl.shouldRebuildTree = true
	sl.sums = sl.sums[:0]
	sl.allSum = e()
}

func (sl *SortedListWithSum) ForEach(f func(value E, index int), reverse bool) {
	if !reverse {
		count := 0
		for i := 0; i < len(sl.blocks); i++ {
			block := sl.blocks[i]
			for j := 0; j < len(block); j++ {
				f(block[j], count)
				count++
			}
		}
		return
	}

	count := 0
	for i := len(sl.blocks) - 1; i >= 0; i-- {
		block := sl.blocks[i]
		for j := len(block) - 1; j >= 0; j-- {
			f(block[j], count)
			count++
		}
	}
}

func (sl *SortedListWithSum) Enumerate(start, end int, f func(value E), erase bool) {
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
	for ; count > 0 && pos < len(sl.blocks); pos++ {
		block := sl.blocks[pos]
		endIndex := min(len(block), startIndex+count)
		if f != nil {
			for j := startIndex; j < endIndex; j++ {
				f(block[j])
			}
		}
		deleted := endIndex - startIndex

		if erase {
			if deleted == len(block) {
				// !delete block
				sl.blocks = append(sl.blocks[:pos], sl.blocks[pos+1:]...)
				sl.mins = append(sl.mins[:pos], sl.mins[pos+1:]...)
				sl.allSum = op(sl.allSum, inv(sl.sums[pos]))
				sl.sums = append(sl.sums[:pos], sl.sums[pos+1:]...)
				sl.shouldRebuildTree = true
				pos--
			} else {
				// !delete [index, end)
				for i := startIndex; i < endIndex; i++ {
					inv_ := inv(block[i])
					sl.allSum = op(sl.allSum, inv_)
					sl.sums[pos] = op(sl.sums[pos], inv_)
				}
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

func (sl *SortedListWithSum) Slice(start, end int) []E {
	if start < 0 {
		start += sl.size
	}
	if start < 0 {
		start = 0
	}
	if end < 0 {
		end += sl.size
	}
	if end > sl.size {
		end = sl.size
	}
	if start >= end {
		return nil
	}
	count := end - start
	res := make([]E, 0, count)
	pos, index := sl._findKth(start)
	for ; count > 0 && pos < len(sl.blocks); pos++ {
		block := sl.blocks[pos]
		endPos := min(len(block), index+count)
		curCount := endPos - index
		res = append(res, block[index:endPos]...)
		count -= curCount
		index = 0
	}
	return res
}

func (sl *SortedListWithSum) Range(min, max E) []E {
	if sl.less(max, min) {
		return nil
	}
	res := []E{}
	pos := sl._locBlock(min)
	for i := pos; i < len(sl.blocks); i++ {
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

func (sl *SortedListWithSum) IteratorAt(index int) *Iterator {
	if index < 0 {
		index += sl.size
	}
	if index < 0 || index >= sl.size {
		panic("Index out of range")
	}
	pos, startIndex := sl._findKth(index)
	return sl._iteratorAt(pos, startIndex)
}

func (sl *SortedListWithSum) LowerBound(value E) *Iterator {
	pos, index := sl._locLeft(value)
	return sl._iteratorAt(pos, index)
}

func (sl *SortedListWithSum) UpperBound(value E) *Iterator {
	pos, index := sl._locRight(value)
	return sl._iteratorAt(pos, index)
}

func (sl *SortedListWithSum) Min() E {
	if sl.size == 0 {
		panic("Min() called on empty SortedListWithSum")
	}
	return sl.mins[0]
}

func (sl *SortedListWithSum) Max() E {
	if sl.size == 0 {
		panic("Max() called on empty SortedListWithSum")
	}
	lastBlock := sl.blocks[len(sl.blocks)-1]
	return lastBlock[len(lastBlock)-1]
}

func (sl *SortedListWithSum) String() string {
	sb := strings.Builder{}
	sb.WriteString("SortedListWithSum{")
	sl.ForEach(func(value E, index int) {
		if index > 0 {
			sb.WriteByte(',')
		}
		sb.WriteString(fmt.Sprintf("%v", value))
	}, false)
	sb.WriteByte('}')
	return sb.String()
}

func (sl *SortedListWithSum) Len() int {
	return sl.size
}

func (sl *SortedListWithSum) _delete(pos, index int) {
	// !delete element
	sl.size--
	sl._updateTree(pos, -1)
	block := sl.blocks[pos]
	deleted := block[index]
	sl.blocks[pos] = append(block[:index], block[index+1:]...)
	sl.allSum = op(sl.allSum, inv(deleted))
	if len(sl.blocks[pos]) > 0 {
		sl.mins[pos] = sl.blocks[pos][0]
		sl.sums[pos] = op(sl.sums[pos], inv(deleted))
		return
	}

	// !delete block
	sl.blocks = append(sl.blocks[:pos], sl.blocks[pos+1:]...)
	sl.mins = append(sl.mins[:pos], sl.mins[pos+1:]...)
	sl.shouldRebuildTree = true
	sl.sums = append(sl.sums[:pos], sl.sums[pos+1:]...)
}

func (sl *SortedListWithSum) _locLeft(value E) (pos, index int) {
	if sl.size == 0 {
		return
	}

	// find pos
	left := -1
	right := len(sl.blocks) - 1
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
	right = len(cur)
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

func (sl *SortedListWithSum) _locRight(value E) (pos, index int) {
	if sl.size == 0 {
		return
	}

	// find pos
	left := 0
	right := len(sl.blocks)
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
	right = len(cur)
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

func (sl *SortedListWithSum) _locBlock(value E) int {
	left, right := -1, len(sl.blocks)-1
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

func (sl *SortedListWithSum) _buildTree() {
	sl.tree = make([]int, len(sl.blocks))
	for i := 0; i < len(sl.blocks); i++ {
		sl.tree[i] = len(sl.blocks[i])
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

func (sl *SortedListWithSum) _updateTree(index, delta int) {
	if sl.shouldRebuildTree {
		return
	}
	tree := sl.tree
	for i := index; i < len(tree); i |= i + 1 {
		tree[i] += delta
	}
}

func (sl *SortedListWithSum) _queryTree(end int) int {
	if sl.shouldRebuildTree {
		sl._buildTree()
	}
	tree := sl.tree
	sum := 0
	for end > 0 {
		sum += tree[end-1]
		end &= end - 1
	}
	return sum
}

func (sl *SortedListWithSum) _findKth(k int) (pos, index int) {
	if k < len(sl.blocks[0]) {
		return 0, k
	}
	last := len(sl.blocks) - 1
	lastLen := len(sl.blocks[last])
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
		if next < len(tree) && k >= tree[next] {
			pos = next
			k -= tree[pos]
		}
	}
	return pos + 1, k
}

func (sl *SortedListWithSum) _rebuildSum(pos int) {
	block := sl.blocks[pos]
	cur := e()
	for _, v := range block {
		cur = op(cur, v)
	}
	sl.sums[pos] = cur
}

func (sl *SortedListWithSum) _iteratorAt(pos, index int) *Iterator {
	return &Iterator{sl: sl, pos: pos, index: index}
}

type Iterator struct {
	sl    *SortedListWithSum
	pos   int
	index int
}

func (it *Iterator) HasNext() bool {
	return it.pos < len(it.sl.blocks)-1 || it.index < len(it.sl.blocks[it.pos])-1
}

func (it *Iterator) Next() (res E, ok bool) {
	if !it.HasNext() {
		return
	}
	it.index++
	if it.index == len(it.sl.blocks[it.pos]) {
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

func (it *Iterator) Prev() (res E, ok bool) {
	if !it.HasPrev() {
		return
	}
	it.index--
	if it.index == -1 {
		it.pos--
		it.index = len(it.sl.blocks[it.pos]) - 1
	}
	res = it.sl.blocks[it.pos][it.index]
	ok = true
	return
}

func (it *Iterator) Remove() {
	it.sl._delete(it.pos, it.index)
}

func (it *Iterator) Value() (res E, ok bool) {
	if it.pos < 0 || it.pos >= it.sl.Len() {
		return
	}
	block := it.sl.blocks[it.pos]
	if it.index < 0 || it.index >= len(block) {
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
