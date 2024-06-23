// 区间频率查询
// 单点修改，查询区间某元素出现次数(PointSetRangeFreq)
//
// !如果要支持insert/pop操作，使用"WaveletMatrixDynamic"

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"math/rand"
	"os"
	"sort"
)

func main() {
	// test()
	pointSetRangeFreq()
}

// https://judge.yosupo.jp/problem/point_set_range_frequency
func pointSetRangeFreq() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int32
	fmt.Fscan(in, &n, &q)
	nums := make([]int32, n)
	for i := int32(0); i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	rfq := NewRangeFreqQueryPointSet(n, func(i int32) S { return S(nums[i]) }, func(a, b S) bool { return a < b })
	for i := int32(0); i < q; i++ {
		var op int32
		fmt.Fscan(in, &op)
		if op == 0 {
			var index, value int32
			fmt.Fscan(in, &index, &value)
			rfq.Set(index, value)
		} else {
			var start, end, value int32
			fmt.Fscan(in, &start, &end, &value)
			fmt.Fprintln(out, rfq.RangeFreq(start, end, value))
		}
	}
}

type RangeFreqQueryPointSet struct {
	arr  []S
	mp   map[S]*sortedList
	less func(a, b S) bool
}

func NewRangeFreqQueryPointSet(n int32, f func(i int32) S, less func(a, b S) bool) *RangeFreqQueryPointSet {
	arr := make([]S, n)
	tmp := make(map[S][]int32)
	for i := int32(0); i < n; i++ {
		arr[i] = f(i)
		tmp[arr[i]] = append(tmp[arr[i]], i)
	}
	res := &RangeFreqQueryPointSet{arr: arr, mp: make(map[S]*sortedList, len(tmp)), less: less}
	for k, v := range tmp {
		res.mp[k] = NewSortedList(less, v...)
	}
	return res
}

func (rfq *RangeFreqQueryPointSet) RangeFreq(start, end int32, value S) int32 {
	if start < 0 {
		start = 0
	}
	if end > int32(len(rfq.arr)) {
		end = int32(len(rfq.arr))
	}
	if start >= end {
		return 0
	}
	sl := rfq.getSl(value)
	return int32(sl.BisectLeft(end) - sl.BisectLeft(start))
}

func (rfq *RangeFreqQueryPointSet) Set(index int32, value S) {
	pre := rfq.arr[index]
	if pre == value {
		return
	}
	rfq.getSl(pre).Discard(index)
	rfq.getSl(value).Add(index)
	rfq.arr[index] = value
}

func (rfq *RangeFreqQueryPointSet) Get(index int32) S { return rfq.arr[index] }

func (rfq *RangeFreqQueryPointSet) Append(value S) {
	rfq.getSl(value)._appendLast(int32(len(rfq.arr)))
	rfq.arr = append(rfq.arr, value)
}

func (rfq *RangeFreqQueryPointSet) Pop() S {
	index := int32(len(rfq.arr) - 1)
	value := rfq.arr[index]
	rfq.getSl(value)._popLast()
	rfq.arr = rfq.arr[:len(rfq.arr)-1]
	return value
}

func (rfq *RangeFreqQueryPointSet) getSl(value S) *sortedList {
	if sl, has := rfq.mp[value]; has {
		return sl
	}
	sl := NewSortedList(rfq.less)
	rfq.mp[value] = sl
	return sl
}

// 1e5 -> 200, 2e5 -> 400
const _LOAD int = 100

type S = int32

var EMPTY S

// 使用分块+树状数组维护的有序序列.
type sortedList struct {
	less              func(a, b S) bool
	size              int
	blocks            [][]S
	mins              []S
	tree              []int
	shouldRebuildTree bool
}

func NewSortedList(less func(a, b S) bool, elements ...S) *sortedList {
	elements = append(elements[:0:0], elements...)
	res := &sortedList{less: less}
	sort.Slice(elements, func(i, j int) bool { return less(elements[i], elements[j]) })
	n := len(elements)
	blocks := [][]S{}
	for start := 0; start < n; start += _LOAD {
		end := min(start+_LOAD, n)
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

func (sl *sortedList) _appendLast(value S) {
	sl.size++
	if len(sl.blocks) == 0 {
		sl.blocks = append(sl.blocks, []S{value})
		sl.mins = append(sl.mins, value)
		sl.shouldRebuildTree = true
		return
	}
	pos := len(sl.blocks) - 1
	sl._updateTree(pos, 1)
	sl.blocks[pos] = append(sl.blocks[pos], value)
	sl._adjust(pos)
	return
}
func (sl *sortedList) _popFirst() S {
	pos, startIndex := 0, 0
	value := sl.blocks[pos][startIndex]
	sl._delete(pos, startIndex)
	return value
}
func (sl *sortedList) _popLast() S {
	pos := len(sl.blocks) - 1
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
	// sl.shouldRebuildTree = true // TODO: 能否不重建树
	return value
}
func (sl *sortedList) _adjust(pos int) {
	// n -> load + (n - load)
	if n := len(sl.blocks[pos]); _LOAD+_LOAD < n {
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

func (sl *sortedList) Add(value S) {
	sl.size++
	if len(sl.blocks) == 0 {
		sl.blocks = append(sl.blocks, []S{value})
		sl.mins = append(sl.mins, value)
		sl.shouldRebuildTree = true
		return
	}
	pos, index := sl._locRight(value)
	sl._updateTree(pos, 1)
	sl.blocks[pos] = append(sl.blocks[pos][:index], append([]S{value}, sl.blocks[pos][index:]...)...)
	sl.mins[pos] = sl.blocks[pos][0]
	// n -> load + (n - load)
	if n := len(sl.blocks[pos]); _LOAD+_LOAD < n {
		sl.blocks = append(sl.blocks, nil)
		copy(sl.blocks[pos+2:], sl.blocks[pos+1:])
		sl.blocks[pos+1] = sl.blocks[pos][_LOAD:]
		sl.blocks[pos] = sl.blocks[pos][:_LOAD:_LOAD]
		sl.mins = append(sl.mins, EMPTY)
		copy(sl.mins[pos+2:], sl.mins[pos+1:])
		sl.mins[pos+1] = sl.blocks[pos+1][0]
		sl.shouldRebuildTree = true
	}
	return
}

func (sl *sortedList) Discard(value S) bool {
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

// 返回第一个大于等于 `value` 的元素的索引/严格小于 `value` 的元素的个数.
func (sl *sortedList) BisectLeft(value S) int {
	pos, index := sl._locLeft(value)
	return sl._queryTree(pos) + index
}

func (sl *sortedList) Len() int {
	return sl.size
}

func (sl *sortedList) _delete(pos, index int) {
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

func (sl *sortedList) _locLeft(value S) (pos, index int) {
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

func (sl *sortedList) _locRight(value S) (pos, index int) {
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

func (sl *sortedList) _buildTree() {
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

func (sl *sortedList) _updateTree(index, delta int) {
	if sl.shouldRebuildTree {
		return
	}
	tree := sl.tree
	for i := index; i < len(tree); i |= i + 1 {
		tree[i] += delta
	}
}

func (sl *sortedList) _queryTree(end int) int {
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

func (sl *sortedList) _findKth(k int) (pos, index int) {
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

func test() {
	for i := 0; i < 100; i++ {
		n := rand.Intn(1000) + 1001
		arr := make([]int32, n)
		for i := 0; i < n; i++ {
			arr[i] = int32(rand.Intn(1000))
		}

		rfq := NewRangeFreqQueryPointSet(int32(n), func(i int32) S { return S(arr[i]) }, func(a, b S) bool { return a < b })
		for i := 0; i < 10000; i++ {
			{
				// RangeFreq
				start := int32(rand.Intn(len(arr)))
				end := int32(rand.Intn(len(arr)))
				if start > end {
					start, end = end, start
				}
				value := S(rand.Intn(1000))
				res1 := rfq.RangeFreq(start, end, value)
				res2 := int32(0)
				for i := start; i < end; i++ {
					if arr[i] == value {
						res2++
					}
				}
				if res1 != res2 {
					panic("WA RangeFreq")
				}
			}

			// Set
			{
				index := int32(rand.Intn(n))
				value := S(rand.Intn(1000))
				rfq.Set(index, value)
				arr[index] = value
			}

			// Append
			{
				value := S(rand.Intn(1000))
				rfq.Append(value)
				arr = append(arr, value)
			}

			// Pop
			{
				value1 := rfq.Pop()
				value2 := arr[len(arr)-1]
				arr = arr[:len(arr)-1]
				if value1 != value2 {
					panic("WA Pop")
				}
			}

			// Get
			{
				for i := 0; i < len(arr); i++ {
					if rfq.Get(int32(i)) != arr[i] {
						panic("WA Get")
					}
				}
			}
		}
	}

	fmt.Println("pass")
}
