// 动态中位数，用两个可删除堆(对顶堆)实现
// api:
// 1. Insert(x T)
// 2. Erase(x T)
// 3. Median() (low, high T)
// 4. DistToMedian() T
// 5. Size() int32

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"sort"
)

func main() {
	yuki738()
}

const INF int = 1e18

// No.738 平らな農地
// https://yukicoder.me/problems/no/738
// !滑动窗口所有数到中位数的距离和
func yuki738() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k int32
	fmt.Fscan(in, &n, &k)
	nums := make([]int, n)
	for i := int32(0); i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	M := NewDynamicMedian()
	res := INF

	for i := int32(0); i < n; i++ {
		M.Insert(nums[i])
		if i >= k {
			M.Discard(nums[i-k])
		}
		if i >= k-1 {
			res = min(res, M.DistToMedian())
		}
	}
	fmt.Fprintln(out, res)
}

// 1e5 -> 200, 2e5 -> 400
const _LOAD int = 200

type S = int

type DynamicMedian struct {
	size     int32
	lower    *_sl
	upper    *_sl
	lowerSum S
	upperSum S
}

func NewDynamicMedian() *DynamicMedian {
	return &DynamicMedian{
		lower: NewSortedList(func(a, b S) bool { return a < b }),
		upper: NewSortedList(func(a, b S) bool { return a < b }),
	}
}

func (d *DynamicMedian) Insert(value S) {
	if d.size&1 == 0 {
		d.upper.Add(value)
		d.upperSum += value
	} else {
		d.lower.Add(value)
		d.lowerSum += value
	}
	d.size++
	d.balance()
}

func (d *DynamicMedian) Discard(value S) bool {
	if d.lower.Has(value) {
		d.lower.Discard(value)
		d.lowerSum -= value
		d.size--
		return true
	} else if d.upper.Has(value) {
		d.upper.Discard(value)
		d.upperSum -= value
		d.size--
		return true
	} else {
		return false
	}
}

// 返回中位数.如果元素个数为偶数,返回两个中位数.
func (d *DynamicMedian) Median() (low, high S) {
	if d.size == 0 {
		return
	}
	if d.size&1 == 0 {
		low = d.lower.Max()
		high = d.upper.Min()
	} else {
		low = d.upper.Min()
		high = low
	}
	return
}

func (d *DynamicMedian) DistToMedian() S {
	if d.size == 0 {
		return 0
	}
	fmt.Println(d.lower.Len(), d.upper.Len())
	low, _ := d.Median()
	sum1 := low*d.lower.Len() - d.lowerSum
	if sum1 < 0 {
		sum1 = -sum1
	}
	sum2 := low*d.upper.Len() - d.upperSum
	if sum2 < 0 {
		sum2 = -sum2
	}
	return sum1 + sum2
}

func (d *DynamicMedian) Size() int32 { return d.size }

func (d *DynamicMedian) balance() {
	// 偶数个数时，|lower heap| == |upper heap|
	// 奇数个数时，|lower heap| + 1 == |upper heap|
	for d.lower.Len()+1 < d.upper.Len() {
		d.lower.Add(d.upper.Pop(0))
	}
	for d.lower.Len() > d.upper.Len() {
		d.upper.Add(d.lower.Pop(0))
	}

	if d.size&1 == 0 {
		if d.lower.size != d.upper.size {
			panic("size error")
		}
	} else {
		if d.lower.size+1 != d.upper.size {
			panic("size error")
		}
	}

	if d.lower.Len() == 0 || d.upper.Len() == 0 {
		return
	}

	if d.lower.Max() > d.upper.Min() {
		upperMin := d.upper.Pop(0)
		d.lower.Add(upperMin)
		d.lowerSum += upperMin
		d.upperSum -= upperMin

		lowerMax := d.lower.Pop(d.lower.Len() - 1)
		d.upper.Add(lowerMax)
		d.upperSum += lowerMax
		d.lowerSum -= lowerMax
	}
}

// 使用分块+树状数组维护的有序序列.
type _sl struct {
	less              func(a, b S) bool
	size              int
	blocks            [][]S
	mins              []S
	tree              []int
	shouldRebuildTree bool
}

func NewSortedList(less func(a, b S) bool, elements ...S) *_sl {
	elements = append(elements[:0:0], elements...)
	res := &_sl{less: less}
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

func (sl *_sl) Add(value S) *_sl {
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
	if n := len(sl.blocks[pos]); _LOAD+_LOAD < n {
		sl.blocks = append(sl.blocks[:pos+1], append([][]S{sl.blocks[pos][_LOAD:]}, sl.blocks[pos+1:]...)...)
		sl.mins = append(sl.mins[:pos+1], append([]S{sl.blocks[pos][_LOAD]}, sl.mins[pos+1:]...)...)
		sl.blocks[pos] = sl.blocks[pos][:_LOAD:_LOAD] // !注意max的设置(为了让左右互不影响)
		sl.shouldRebuildTree = true
	}
	return sl
}

func (sl *_sl) Has(value S) bool {
	if len(sl.blocks) == 0 {
		return false
	}
	pos, index := sl._locLeft(value)
	return index < len(sl.blocks[pos]) && sl.blocks[pos][index] == value
}

func (sl *_sl) Pop(index int) S {
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

func (sl *_sl) Discard(value S) bool {
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

func (sl *_sl) Min() S {
	if sl.size == 0 {
		panic("Min() called on empty SortedList")
	}
	return sl.mins[0]
}

func (sl *_sl) Max() S {
	if sl.size == 0 {
		panic("Max() called on empty SortedList")
	}
	lastBlock := sl.blocks[len(sl.blocks)-1]
	return lastBlock[len(lastBlock)-1]
}

func (sl *_sl) Len() int {
	return sl.size
}

func (sl *_sl) _delete(pos, index int) {
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

func (sl *_sl) _locLeft(value S) (pos, index int) {
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

func (sl *_sl) _locRight(value S) (pos, index int) {
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

func (sl *_sl) _locBlock(value S) int {
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

func (sl *_sl) _buildTree() {
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

func (sl *_sl) _updateTree(index, delta int) {
	if sl.shouldRebuildTree {
		return
	}
	tree := sl.tree
	for i := index; i < len(tree); i |= i + 1 {
		tree[i] += delta
	}
}

func (sl *_sl) _queryTree(end int) int {
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

func (sl *_sl) _findKth(k int) (pos, index int) {
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
