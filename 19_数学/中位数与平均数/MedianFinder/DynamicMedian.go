// 动态中位数，用两个平衡树(对顶)实现
// api:
// 1. Insert(x T)
// 2. Discard(x T) bool
// 3. Median() (low, high T)
// 4. DistToMedian() T
// 5. Size() int32

package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"unsafe"
)

func main() {
	yuki738()
	// test()
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

// https://leetcode.cn/problems/minimum-operations-to-make-subarray-elements-equal/description/
func minOperations(nums []int, k int) int64 {
	n := len(nums)
	M := NewDynamicMedian()
	res := INF
	for i := 0; i < n; i++ {
		M.Insert(nums[i])
		if i >= k {
			M.Discard(nums[i-k])
		}
		if i >= k-1 {
			res = min(res, M.DistToMedian())
		}
	}
	return int64(res)
}

// 1e5 -> 100, 2e5 -> 200
const _LOAD int = 75

type S = int

var EMPTY S

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
		upper: NewSortedList(func(a, b S) bool { return a > b }),
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
	if d.lower.Discard(value) {
		d.lowerSum -= value
		d.size--
		d.balance()
		return true
	} else if d.upper.Discard(value) {
		d.upperSum -= value
		d.size--
		d.balance()
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
		high = d.upper.Max()
	} else {
		low = d.upper.Max()
		high = low
	}
	return
}

func (d *DynamicMedian) DistToMedian() S {
	if d.size == 0 {
		return 0
	}
	low, _ := d.Median()
	sum1 := low*S(d.lower.Len()) - d.lowerSum
	sum2 := d.upperSum - low*S(d.upper.Len())
	return sum1 + sum2
}

func (d *DynamicMedian) Size() int32 { return d.size }

func (d *DynamicMedian) balance() {
	// 偶数个数时，|lower heap| == |upper heap|
	// 奇数个数时，|lower heap| + 1 == |upper heap|
	for d.lower.Len()+1 < d.upper.Len() {
		upperMin := d.upper._popLast()
		d.lower._appendLast(upperMin)
		d.lowerSum += upperMin
		d.upperSum -= upperMin
	}
	for d.lower.Len() > d.upper.Len() {
		lowerMin := d.lower._popLast()
		d.upper._appendLast(lowerMin)
		d.upperSum += lowerMin
		d.lowerSum -= lowerMin
	}

	// if d.size&1 == 0 {
	// 	if d.lower.size != d.upper.size {
	// 		panic("size error")
	// 	}
	// } else {
	// 	if d.lower.size+1 != d.upper.size {
	// 		panic("size error")
	// 	}
	// }

	if d.lower.Len() == 0 || d.upper.Len() == 0 {
		return
	}

	if d.lower.Max() > d.upper.Max() {
		upperMin := d.upper._popLast()
		d.lower.Add(upperMin)
		d.lowerSum += upperMin
		d.upperSum -= upperMin

		lowerMax := d.lower._popLast()
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
	sl.blocks[pos] = Insert(sl.blocks[pos], index, value)
	sl.mins[pos] = sl.blocks[pos][0]
	// n -> load + (n - load)
	if n := len(sl.blocks[pos]); _LOAD+_LOAD < n {
		left := append([]S(nil), sl.blocks[pos][:_LOAD]...)
		right := append([]S(nil), sl.blocks[pos][_LOAD:]...)
		sl.blocks = Replace(sl.blocks, pos, pos+1, left, right)
		sl.mins = Insert(sl.mins, pos+1, right[0])
		sl.shouldRebuildTree = true
	}
	return sl
}

func (sl *_sl) _appendLast(value S) *_sl {
	sl.size++
	if len(sl.blocks) == 0 {
		sl.blocks = append(sl.blocks, []S{value})
		sl.mins = append(sl.mins, value)
		sl.shouldRebuildTree = true
		return sl
	}
	pos := len(sl.blocks) - 1
	sl._updateTree(pos, 1)
	sl.blocks[pos] = append(sl.blocks[pos], value)
	// n -> load + (n - load)
	if n := len(sl.blocks[pos]); _LOAD+_LOAD < n {
		left := append([]S(nil), sl.blocks[pos][:_LOAD]...)
		right := append([]S(nil), sl.blocks[pos][_LOAD:]...)
		sl.blocks = Replace(sl.blocks, pos, pos+1, left, right)
		sl.mins = Insert(sl.mins, pos+1, right[0])
		sl.shouldRebuildTree = true
	}
	return sl
}

func (sl *_sl) _popLast() S {
	sl.size--
	pos := len(sl.blocks) - 1
	res := sl.blocks[pos][len(sl.blocks[pos])-1]
	sl._updateTree(pos, -1)
	sl.blocks[pos] = sl.blocks[pos][:len(sl.blocks[pos])-1]
	if len(sl.blocks[pos]) == 0 {
		// !delete block
		sl.blocks = sl.blocks[:pos]
		sl.mins = sl.mins[:pos]
		sl.shouldRebuildTree = true
	}
	return res
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
	sl.blocks[pos] = Replace(sl.blocks[pos], int(index), int(index+1))
	if len(sl.blocks[pos]) > 0 {
		sl.mins[pos] = sl.blocks[pos][0]
		return
	}

	// !delete block
	sl.blocks = Replace(sl.blocks, int(pos), int(pos)+1)
	sl.mins = Replace(sl.mins, int(pos), int(pos)+1)
	sl.shouldRebuildTree = true
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

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

// Replace replaces the elements s[i:j] by the given v, and returns the modified slice.
// !Like JavaScirpt's Array.prototype.splice.
func Replace[S ~[]E, E any](s S, i, j int, v ...E) S {
	if i < 0 {
		i = 0
	}
	if j > len(s) {
		j = len(s)
	}
	if i == j {
		return Insert(s, i, v...)
	}
	if j == len(s) {
		return append(s[:i], v...)
	}
	tot := len(s[:i]) + len(v) + len(s[j:])
	if tot > cap(s) {
		s2 := append(s[:i], make(S, tot-i)...)
		copy(s2[i:], v)
		copy(s2[i+len(v):], s[j:])
		return s2
	}
	r := s[:tot]
	if i+len(v) <= j {
		copy(r[i:], v)
		copy(r[i+len(v):], s[j:])
		// clear(s[tot:])
		return r
	}
	if !overlaps(r[i+len(v):], v) {
		copy(r[i+len(v):], s[j:])
		copy(r[i:], v)
		return r
	}
	y := len(v) - (j - i)
	if !overlaps(r[i:j], v) {
		copy(r[i:j], v[y:])
		copy(r[len(s):], v[:y])
		rotateRight(r[i:], y)
		return r
	}
	if !overlaps(r[len(s):], v) {
		copy(r[len(s):], v[:y])
		copy(r[i:j], v[y:])
		rotateRight(r[i:], y)
		return r
	}
	k := startIdx(v, s[j:])
	copy(r[i:], v)
	copy(r[i+len(v):], r[i+k:])
	return r
}

func rotateLeft[E any](s []E, r int) {
	for r != 0 && r != len(s) {
		if r*2 <= len(s) {
			swap(s[:r], s[len(s)-r:])
			s = s[:len(s)-r]
		} else {
			swap(s[:len(s)-r], s[r:])
			s, r = s[len(s)-r:], r*2-len(s)
		}
	}
}

func rotateRight[E any](s []E, r int) {
	rotateLeft(s, len(s)-r)
}

func swap[E any](x, y []E) {
	for i := 0; i < len(x); i++ {
		x[i], y[i] = y[i], x[i]
	}
}

func overlaps[E any](a, b []E) bool {
	if len(a) == 0 || len(b) == 0 {
		return false
	}
	elemSize := unsafe.Sizeof(a[0])
	if elemSize == 0 {
		return false
	}
	return uintptr(unsafe.Pointer(&a[0])) <= uintptr(unsafe.Pointer(&b[len(b)-1]))+(elemSize-1) &&
		uintptr(unsafe.Pointer(&b[0])) <= uintptr(unsafe.Pointer(&a[len(a)-1]))+(elemSize-1)
}

func startIdx[E any](haystack, needle []E) int {
	p := &needle[0]
	for i := range haystack {
		if p == &haystack[i] {
			return i
		}
	}
	panic("needle not found")
}

func Insert[S ~[]E, E any](s S, i int, v ...E) S {
	if i < 0 {
		i = 0
	}
	if i > len(s) {
		i = len(s)
	}

	m := len(v)
	if m == 0 {
		return s
	}
	n := len(s)
	if i == n {
		return append(s, v...)
	}
	if n+m > cap(s) {
		s2 := append(s[:i], make(S, n+m-i)...)
		copy(s2[i:], v)
		copy(s2[i+m:], s[i:])
		return s2
	}
	s = s[:n+m]
	if !overlaps(v, s[i+m:]) {
		copy(s[i+m:], s[i:])
		copy(s[i:], v)
		return s
	}
	copy(s[n:], v)
	rotateRight(s[i:], m)
	return s
}

func test() {
	for i := 0; i < 1000; i++ {
		M := NewDynamicMedian()
		sortedNums := make([]int, 0)

		add := func(x int) {

			sortedNums = append(sortedNums, x)
			sort.Ints(sortedNums)
		}

		discard := func(x int) {

			for i, v := range sortedNums {
				if v == x {
					sortedNums = append(sortedNums[:i], sortedNums[i+1:]...)
					break
				}
			}
		}

		median := func() (low, high int) {
			if len(sortedNums) == 0 {
				return
			}
			n := len(sortedNums)
			if n&1 == 0 {
				low = sortedNums[n/2-1]
				high = sortedNums[n/2]
			} else {
				low = sortedNums[n/2]
				high = low
			}
			return
		}

		distToMedian := func() int {
			if len(sortedNums) == 0 {
				return 0
			}
			low, _ := median()
			res := 0
			for _, v := range sortedNums {
				res += abs(v - low)
			}
			return res
		}

		size := func() int {
			return len(sortedNums)
		}

		for j := 0; j < 1000; j++ {
			x := rand.Intn(10000)

			// add
			M.Insert(x)
			add(x)

			// discard
			y := rand.Intn(10000)
			M.Discard(y)
			discard(y)

			// median
			low, high := M.Median()
			low2, high2 := median()
			if low != low2 || high != high2 {
				fmt.Println("error")
				fmt.Println(low, high, low2, high2)
				fmt.Println(sortedNums)
				panic("error")
			}

			// distToMedian
			res := M.DistToMedian()
			res2 := distToMedian()
			if res != res2 {
				fmt.Println("error")
				fmt.Println(res, res2)
				fmt.Println(sortedNums)
				panic("error")
			}

			// size
			sz := M.Size()
			sz2 := size()
			if sz != int32(sz2) {
				fmt.Println("error")
				fmt.Println(sz, sz2)
				fmt.Println(sortedNums)
				panic("error")
			}

		}

	}
	fmt.Println("pass")
}
