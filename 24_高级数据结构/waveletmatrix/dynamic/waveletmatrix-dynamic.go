// WaveletMatrixDynamic/DynamicWaveletMatrix
// 继承自waveletMatrixStatic，增加了 insert/pop/set 三个修改操作，单次操作时间复杂度O(lognlogU).

// api:
//  1. PrefixCount(end int32, x int) int32
//  2. RangeCount(start, end int32, x int) int32
//  3. RangeFreq(start, end int32, x, y int) int32
//
//  4. Kth(k int32, x int) int32
//
//  5. KthSmallest(start, end int32, k int32) int
//     KthSmallestIndex(start, end int32, k int32) int32
//  6. KthLargest(start, end int32, k int32) int
//     KthLargestIndex(start, end int32, k int32) int32
//
//  7. TopK(start, end int32, k int32) []topKPair
//
//  8. Floor(start, end int32, x int) (int, bool)
//  9. Lower(start, end int32, x int) (int, bool)
//  10. Ceil(start, end int32, x int) (int, bool)
//  11. Higher(start, end int32, x int) (int, bool)
//
//  12. CountAll(start, end int32, x int) (same, less, more int32)
//  13. CountSame(start, end int32, x int) int32
//  14. CountLess(start, end int32, x int) int32
//  15. CountMore(start, end int32, x int) int32
//
//  16. Insert(index int32, x int)
//  17. Pop(index int32) int
//  18. Set(index int32, x int)

// !常数较大.最好预先离散化，减少值域大小.
// - 计数问题，需要将x离散化成BisectLeft(origin, x).
// - 区间第k小问题，需要答案转换成origin[答案].

// TODO: golang 太慢，需要使用 Rust 版本

package main

import (
	"bufio"
	"fmt"
	"math"
	"math/bits"
	"math/rand"
	"os"
	"sort"
	"time"
)

func main() {
	// test()
	// testTime()

	// CF455D()
	// libraryQuery()
	// yosupo()
}

// Serega and Fun
// https://www.luogu.com.cn/problem/CF455D
// 1 start end: 区间[start, end)向右移动.
// 2 start end v: 区间[start, end)中k出现次数.
// 强制在线.
// n,q<=1e5.
func CF455D() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int32
	fmt.Fscan(in, &n)
	nums := make([]int, n)
	for i := int32(0); i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}
	set := make(map[int]struct{})
	for _, num := range nums {
		set[num] = struct{}{}
	}

	newNums, origin := DiscretizeFast(nums)
	wm := NewWaveletMatrixDynamic(int32(len(nums)), func(i int32) int { return int(newNums[i]) }, len(origin))

	rotateRight := func(start, end int32) {
		value := wm.Get(end - 1)
		wm.Pop(end - 1)
		wm.Insert(start, value)
	}

	query := func(start, end int32, v int) int32 {
		if _, has := set[v]; !has {
			return 0
		}
		v = int(BisectLeft(origin, v))
		return wm.RangeCount(start, end, v)
	}

	preRes := int32(0)
	var q int32
	fmt.Fscan(in, &q)
	for i := int32(0); i < q; i++ {
		var op int32
		fmt.Fscan(in, &op)
		if op == 1 {
			var l, r int32
			fmt.Fscan(in, &l, &r)
			l, r = (l+preRes-1)%n+1, (r+preRes-1)%n+1
			if l > r {
				l, r = r, l
			}
			l--
			rotateRight(l, r)
		} else {
			var l, r, v int32
			fmt.Fscan(in, &l, &r, &v)
			l, r, v = (l+preRes-1)%n+1, (r+preRes-1)%n+1, (v+preRes-1)%n+1
			if l > r {
				l, r = r, l
			}
			l--
			preRes = query(l, r, int(v))
			fmt.Fprintln(out, preRes)
		}
	}
}

// https://www.hackerrank.com/challenges/library-query/problem
// 0 start end k: 输出第k小的元素.
// 1 x v: 将第x个元素修改为v.
func libraryQuery() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)

	solve := func() {
		var n int32
		fmt.Fscan(in, &n)
		nums := make([]int, n)
		for i := int32(0); i < n; i++ {
			fmt.Fscan(in, &nums[i])
		}
		var q int32
		fmt.Fscan(in, &q)
		operations := make([][4]int32, q)
		for i := int32(0); i < q; i++ {
			var op int
			fmt.Fscan(in, &op)
			if op == 0 {
				var start, end, k int32
				fmt.Fscan(in, &start, &end, &k)
				start--
				k--
				operations[i] = [4]int32{0, start, end, k}
			} else {
				var x, v int32
				fmt.Fscan(in, &x, &v)
				x--
				operations[i] = [4]int32{1, x, v, 0}
			}
		}

		allNums := make([]int, 0, n+q)
		allNums = append(allNums, nums...)
		for _, op := range operations {
			if op[0] == 1 {
				allNums = append(allNums, int(op[2]))
			}
		}
		newNums, origin := DiscretizeFast(allNums)
		wm := NewWaveletMatrixDynamic(int32(len(nums)), func(i int32) int { return int(newNums[i]) }, len(origin))
		for _, op := range operations {
			if op[0] == 0 {
				start, end, k := op[1], op[2], op[3]
				res := wm.KthSmallest(start, end, k)
				fmt.Fprintln(out, origin[res])
			} else {
				index, value := op[1], op[2]
				value = BisectLeft(origin, int(value))
				wm.Set(index, int(value))
			}
		}
	}

	for t := 0; t < T; t++ {
		solve()
	}
}

// https://judge.yosupo.jp/problem/range_kth_smallest
func yosupo() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	newNums, origin := DiscretizeFast(nums)
	wm := NewWaveletMatrixDynamic(int32(len(nums)), func(i int32) int { return int(newNums[i]) }, len(origin))

	for i := 0; i < q; i++ {
		var start, end, x int32
		fmt.Fscan(in, &start, &end, &x)
		res := wm.KthSmallest(start, end, x)
		fmt.Fprintln(out, origin[res])
	}
}

// 将nums中的元素进行离散化，返回新的数组和对应的原始值.
// origin[newNums[i]] == nums[i]
func DiscretizeFast(nums []int) (newNums []int32, origin []int) {
	newNums = make([]int32, len(nums))
	origin = make([]int, 0, len(newNums))
	order := argSort(int32(len(nums)), func(i, j int32) bool { return nums[i] < nums[j] })
	for _, i := range order {
		if len(origin) == 0 || origin[len(origin)-1] != nums[i] {
			origin = append(origin, nums[i])
		}
		newNums[i] = int32(len(origin) - 1)
	}
	origin = origin[:len(origin):len(origin)]
	return
}

func BisectLeft(nums []int, target int) int32 {
	left, right := int32(0), int32(len(nums)-1)
	for left <= right {
		mid := (left + right) >> 1
		if nums[mid] < target {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return left
}

func argSort(n int32, less func(i, j int32) bool) []int32 {
	order := make([]int32, n)
	for i := range order {
		order[i] = int32(i)
	}
	sort.Slice(order, func(i, j int) bool { return less(order[i], order[j]) })
	return order
}

func maxs(nums []int) int {
	max := nums[0]
	for _, num := range nums {
		if num > max {
			max = num
		}
	}
	return max
}

func maxs32(nums []int32) int32 {
	max := nums[0]
	for _, num := range nums {
		if num > max {
			max = num
		}
	}
	return max
}

// 维护[0,maxValue].
type WaveletMatrixDynamic struct {
	size     int32
	maxValue int
	bitLen   int32
	mid      []int32
	bv       []*AVLTreeBitVector
}

type topKPair = struct {
	value int
	count int32
}

func NewWaveletMatrixDynamic(n int32, f func(int32) int, maxValue int) *WaveletMatrixDynamic {
	if maxValue <= 0 {
		maxValue = 1
	}
	res := &WaveletMatrixDynamic{
		size:     n,
		maxValue: maxValue,
		bitLen:   int32(bits.Len(uint(maxValue))),
	}
	res.mid = make([]int32, res.bitLen)
	res.bv = make([]*AVLTreeBitVector, res.bitLen)
	if n > 0 {
		res._build(n, f)
	}
	return res
}

func (wm *WaveletMatrixDynamic) Insert(index int32, x int) {
	if index < 0 {
		index += wm.size
	}
	if index < 0 || index > wm.size {
		panic(fmt.Sprintf("index out of range: %d", index))
	}
	for bit := wm.bitLen - 1; bit >= 0; bit-- {
		if x>>bit&1 == 1 {
			s := wm.bv[bit]._insertAndCount(index, 1)
			index = s + wm.mid[bit]
		} else {
			s := wm.bv[bit]._insertAndCount(index, 0)
			index -= s
			wm.mid[bit]++
		}
	}
	wm.size++
}

func (wm *WaveletMatrixDynamic) Pop(index int32) int {
	if index < 0 {
		index += wm.size
	}
	if index < 0 || index >= wm.size {
		panic(fmt.Sprintf("index out of range: %d", index))
	}
	res := 0
	for bit := wm.bitLen - 1; bit >= 0; bit-- {
		sb := wm.bv[bit]._accessPopAndCount1(index)
		s := sb >> 1
		if sb&1 == 1 {
			res |= 1 << bit
			index = s + wm.mid[bit]
		} else {
			wm.mid[bit]--
			index -= s
		}
	}
	wm.size--
	return res
}

func (wm *WaveletMatrixDynamic) Set(index int32, x int) {
	wm.Pop(index)
	wm.Insert(index, x)
}

func (wm *WaveletMatrixDynamic) PrefixCount(end int32, x int) int32 {
	if x > wm.maxValue {
		return 0
	}
	if end > wm.size {
		end = wm.size
	}
	if end <= 0 {
		return 0
	}
	start := int32(0)
	mid := wm.mid
	for bit := wm.bitLen - 1; bit >= 0; bit-- {
		if x>>bit&1 == 1 {
			start = wm.bv[bit].Count1(start) + mid[bit]
			end = wm.bv[bit].Count1(end) + mid[bit]
		} else {
			start = wm.bv[bit].Count0(start)
			end = wm.bv[bit].Count0(end)
		}
	}
	return end - start
}

func (wm *WaveletMatrixDynamic) RangeCount(start, end int32, x int) int32 {
	if x > wm.maxValue {
		return 0
	}
	if start < 0 {
		start = 0
	}
	if end > wm.size {
		end = wm.size
	}
	if start >= end {
		return 0
	}
	same, _, _ := wm.CountAll(start, end, x)
	return same
}

func (wm *WaveletMatrixDynamic) RangeFreq(start, end int32, floor, higher int) int32 {
	if floor >= higher {
		return 0
	}
	if floor > wm.maxValue {
		return 0
	}
	if start < 0 {
		start = 0
	}
	if end > wm.size {
		end = wm.size
	}
	if start >= end {
		return 0
	}
	return wm.CountLess(start, end, higher) - wm.CountLess(start, end, floor)
}

// 返回第k个x所在的位置.
func (wm *WaveletMatrixDynamic) Kth(k int32, x int) int32 {
	s := int32(0)
	for bit := wm.bitLen - 1; bit >= 0; bit-- {
		if x>>bit&1 == 1 {
			s = wm.bv[bit].Count0(wm.size) + wm.bv[bit].Count1(s)
		} else {
			s = wm.bv[bit].Count0(s)
		}
	}
	s += k
	for bit := int32(0); bit < wm.bitLen; bit++ {
		if x>>bit&1 == 1 {
			s = wm.bv[bit].Kth1(s - wm.bv[bit].Count0(wm.size))
		} else {
			s = wm.bv[bit].Kth0(s)
		}
	}
	return s
}

func (wm *WaveletMatrixDynamic) KthSmallest(start, end int32, k int32) int {
	if k < 0 || k >= end-start {
		return -1
	}
	res := 0
	for bit := wm.bitLen - 1; bit >= 0; bit-- {
		l0, r0 := wm.bv[bit].Count0(start), wm.bv[bit].Count0(end)
		if c := r0 - l0; c <= k {
			res |= 1 << bit
			k -= c
			start += wm.mid[bit] - l0
			end += wm.mid[bit] - r0
		} else {
			start, end = l0, r0
		}
	}
	return res
}

func (wm *WaveletMatrixDynamic) KthSmallestIndex(start, end int32, k int32) int32 {
	if k < 0 || k >= end-start {
		return -1
	}
	val := 0
	for i := wm.bitLen - 1; i >= 0; i-- {
		numOfZeroBegin := wm.bv[i].Count0(start)
		numOfZeroEnd := wm.bv[i].Count0(end)
		numOfZero := numOfZeroEnd - numOfZeroBegin
		bit := 0
		if k >= numOfZero {
			bit = 1
		}
		if bit == 1 {
			k -= numOfZero
			start = wm.mid[i] + start - numOfZeroBegin
			end = wm.mid[i] + end - numOfZeroEnd
		} else {
			start = numOfZeroBegin
			end = numOfZeroBegin + numOfZero
		}
		val = (val << 1) | bit
	}

	left := int32(0)
	for i := wm.bitLen - 1; i >= 0; i-- {
		bit := int8((val >> i) & 1)
		left = wm.bv[i].Count(left, bit)
		if bit == 1 {
			left += wm.mid[i]
		}
	}
	rank := start + k - left
	return wm.Kth(rank, val)
}

func (wm *WaveletMatrixDynamic) KthLargest(start, end int32, k int32) int {
	return wm.KthSmallest(start, end, end-start-k-1)
}
func (wm *WaveletMatrixDynamic) KthLargestIndex(start, end int32, k int32) int32 {
	return wm.KthSmallestIndex(start, end, end-start-k-1)
}

// 按照出现次数排序，返回出现次数最多的k个元素.
func (wm *WaveletMatrixDynamic) TopK(start, end int32, k int32) (res []topKPair) {
	if k == 0 {
		return nil
	}
	type item struct {
		len int32
		x   int
		l   int32
		bit int8
	}

	pq := NewHeap[item](func(a, b item) bool {
		// 频率相同，返回靠前的元素.
		return a.len > b.len
		// 频率相同，返回较小的元素.
		// if a.len != b.len {
		// 	return a.len > b.len
		// }
		// return a.x < b.x
	}, nil)
	pq.Push(item{len: end - start, x: 0, l: start, bit: int8(wm.bitLen - 1)})
	for pq.Len() > 0 {
		v := pq.Pop()
		length, x, l, bit := v.len, v.x, v.l, v.bit
		if bit == -1 {
			res = append(res, topKPair{x, length})
			k--
			if k == 0 {
				break
			}
		} else {
			r := l + length
			l0 := wm.bv[bit].Count0(l)
			r0 := wm.bv[bit].Count0(r)
			if l0 < r0 {
				pq.Push(item{len: r0 - l0, x: x, l: l0, bit: bit - 1})
			}
			l1 := wm.bv[bit].Count1(l) + wm.mid[bit]
			r1 := wm.bv[bit].Count1(r) + wm.mid[bit]
			if l1 < r1 {
				pq.Push(item{len: r1 - l1, x: x | 1<<bit, l: l1, bit: bit - 1})
			}
		}
	}
	return
}

func (wm *WaveletMatrixDynamic) Floor(start, end int32, x int) (int, bool) {
	same, less, _ := wm.CountAll(start, end, x)
	if same > 0 {
		return x, true
	}
	if less == 0 {
		return -1, false
	}
	return wm.KthSmallest(start, end, less-1), true
}
func (wm *WaveletMatrixDynamic) Lower(start, end int32, x int) (int, bool) {
	less := wm.CountLess(start, end, x)
	if less == 0 {
		return -1, false
	}
	return wm.KthSmallest(start, end, less-1), true
}

func (wm *WaveletMatrixDynamic) Ceil(start, end int32, x int) (int, bool) {
	same, less, _ := wm.CountAll(start, end, x)
	if same > 0 {
		return x, true
	}
	if less == end-start {
		return -1, false
	}
	return wm.KthSmallest(start, end, less), true
}

func (wm *WaveletMatrixDynamic) Higher(start, end int32, x int) (int, bool) {
	less := wm.CountLess(start, end, x+1)
	if less == end-start {
		return -1, false
	}
	return wm.KthSmallest(start, end, less), true
}

// 返回[start, end)中等于x的个数，小于x的个数，大于x的个数.
func (wm *WaveletMatrixDynamic) CountAll(start, end int32, x int) (same, less, more int32) {
	if start < 0 {
		start = 0
	}
	if end > wm.size {
		end = wm.size
	}
	if start >= end {
		return 0, 0, 0
	}
	if x > wm.maxValue {
		return 0, end - start, 0
	}
	num := end - start
	for i := wm.bitLen - 1; i >= 0 && start < end; i-- {
		bit := x >> i & 1
		rank0Begin := wm.bv[i].Count0(start)
		rank0End := wm.bv[i].Count0(end)
		rank1Begin := start - rank0Begin
		rank1End := end - rank0End
		if bit == 1 {
			less += rank0End - rank0Begin
			start = wm.mid[i] + rank1Begin
			end = wm.mid[i] + rank1End
		} else {
			more += rank1End - rank1Begin
			start = rank0Begin
			end = rank0End
		}
	}
	same = num - less - more
	return
}

func (wm *WaveletMatrixDynamic) Get(index int32) int {
	if index < 0 {
		index += wm.size
	}
	res := 0
	for bit := wm.bitLen - 1; bit >= 0; bit-- {
		if wm.bv[bit].Get(index) == 1 {
			res |= 1 << bit
			index = wm.bv[bit].Count1(index) + wm.mid[bit]
		} else {
			index = wm.bv[bit].Count0(index)
		}
	}
	return res
}

// 区间[start, end)中小于x的元素个数.
func (wm *WaveletMatrixDynamic) CountLess(start, end int32, x int) int32 {
	if start < 0 {
		start = 0
	}
	if end > wm.size {
		end = wm.size
	}
	if start >= end {
		return 0
	}
	if x > wm.maxValue {
		return end - start
	}
	res := int32(0)
	for bit := wm.bitLen - 1; bit >= 0; bit-- {
		l0, r0 := wm.bv[bit].Count0(start), wm.bv[bit].Count0(end)
		if x>>bit&1 == 1 {
			res += r0 - l0
			start += wm.mid[bit] - l0
			end += wm.mid[bit] - r0
		} else {
			start, end = l0, r0
		}
	}
	return res
}

// 区间[start, end)中大于x的元素个数.
func (wm *WaveletMatrixDynamic) CountMore(start, end int32, x int) int32 {
	return end - start - wm.CountLess(start, end, x+1)
}

func (wm *WaveletMatrixDynamic) CountSame(start, end int32, x int) int32 {
	return wm.RangeCount(start, end, x)
}

func (wm *WaveletMatrixDynamic) Len() int32 {
	return wm.size
}

func (wm *WaveletMatrixDynamic) _build(n int32, f func(int32) int) {
	data := make([]int, n)
	for i := int32(0); i < n; i++ {
		data[i] = f(i)
	}
	bits := make([]int8, n)
	zero, one := make([]int, n), make([]int, n)
	for bit := wm.bitLen - 1; bit >= 0; bit-- {
		p, q := int32(0), int32(0)
		for i, e := range data {
			if e>>bit&1 == 1 {
				bits[i] = 1
				one[q] = e
				q++
			} else {
				bits[i] = 0
				zero[p] = e
				p++
			}
		}
		wm.bv[bit] = newAVLTreeBitVector(bits)
		wm.mid[bit] = p
		zero, data = data, zero
		copy(data[p:], one[:q])
	}

}

type AVLTreeBitVector struct {
	root      int32
	end       int32 // 使用的结点数
	bitLen    []int8
	key       []uint64 // 结点mask
	total     []int32  // 子树onesCount之和
	size      []int32
	left      []int32
	right     []int32
	balance   []int8 // 左子树高度-右子树高度
	pathStack []int32
}

const W int32 = 63
const W8 int8 = 63

func newAVLTreeBitVector(data []int8) *AVLTreeBitVector {
	res := &AVLTreeBitVector{
		root:      0,
		end:       1,
		bitLen:    []int8{0},
		key:       []uint64{0},
		total:     []int32{0},
		size:      []int32{0},
		left:      []int32{0},
		right:     []int32{0},
		balance:   []int8{0},
		pathStack: make([]int32, 0, 128),
	}
	if len(data) > 0 {
		res._build(data)
	}
	return res
}

func (t *AVLTreeBitVector) Reserve(n int32) {
	n = n/W + 1
	t.bitLen = append(t.bitLen, make([]int8, n)...)
	t.key = append(t.key, make([]uint64, n)...)
	t.size = append(t.size, make([]int32, n)...)
	t.total = append(t.total, make([]int32, n)...)
	t.left = append(t.left, make([]int32, n)...)
	t.right = append(t.right, make([]int32, n)...)
	t.balance = append(t.balance, make([]int8, n)...)
}

func (t *AVLTreeBitVector) Insert(index int32, v int8) {
	if t.root == 0 {
		t.root = t._makeNodeWithBitLen1(uint64(v))
		return
	}

	n := t.Len()
	if index < 0 {
		index += n
	}
	if index < 0 {
		index = 0
	}
	if index > n {
		index = n
	}

	v32 := int32(v)
	node := t.root

	t.pathStack = t.pathStack[:0]
	path := t.pathStack

	d := int32(0)
	for node != 0 {
		b32 := int32(t.bitLen[node])
		tmp := t.size[t.left[node]] + b32
		if tmp-b32 <= index && index <= tmp {
			break
		}
		d <<= 1
		t.size[node]++
		t.total[node] += v32
		path = append(path, node)
		if tmp > index {
			node = t.left[node]
			d |= 1
		} else {
			node = t.right[node]
			index -= tmp
		}
	}
	index -= t.size[t.left[node]]
	b32 := int32(t.bitLen[node])
	if b32 < W {
		mask := t.key[node]
		bl := b32 - index
		t.key[node] = (((mask>>bl)<<1 | uint64(v)) << bl) | (mask & ((1 << bl) - 1))
		t.bitLen[node]++
		t.size[node]++
		t.total[node] += v32
		return
	}
	path = append(path, node)
	t.size[node]++
	t.total[node] += v32
	mask := t.key[node]
	bl := W - index
	mask = (((mask>>bl)<<1 | uint64(v)) << bl) | (mask & ((1 << bl) - 1))
	leftKey := mask >> W
	leftKeyPopcount := int32(leftKey & 1)
	t.key[node] = mask & ((1 << W) - 1)
	node = t.left[node]
	d <<= 1
	d |= 1
	if node == 0 {
		last := path[len(path)-1]
		if t.bitLen[last] < W8 {
			t.bitLen[last]++
			t.key[last] = (t.key[last] << 1) | leftKey
			return
		} else {
			t.left[last] = t._makeNodeWithBitLen1(leftKey)
		}
	} else {
		path = append(path, node)
		t.size[node]++
		t.total[node] += leftKeyPopcount
		d <<= 1
		for t.right[node] != 0 {
			node = t.right[node]
			path = append(path, node)
			t.size[node]++
			t.total[node] += leftKeyPopcount
			d <<= 1
		}
		if t.bitLen[node] < W8 {
			t.bitLen[node]++
			t.key[node] = (t.key[node] << 1) | leftKey
			return
		} else {
			t.right[node] = t._makeNodeWithBitLen1(leftKey)
		}
	}
	newNode := int32(0)
	for len(path) > 0 {
		node = path[len(path)-1]
		path = path[:len(path)-1]
		if d&1 == 1 {
			t.balance[node]++
		} else {
			t.balance[node]--
		}
		d >>= 1
		if t.balance[node] == 0 {
			break
		}
		if t.balance[node] == 2 {
			if t.balance[t.left[node]] == -1 {
				newNode = t._rotateLR(node)
			} else {
				newNode = t._rotateL(node)
			}
			break
		} else if t.balance[node] == -2 {
			if t.balance[t.right[node]] == 1 {
				newNode = t._rotateRL(node)
			} else {
				newNode = t._rotateR(node)
			}
			break
		}
	}
	if newNode != 0 {
		if len(path) > 0 {
			if d&1 == 1 {
				t.left[path[len(path)-1]] = newNode
			} else {
				t.right[path[len(path)-1]] = newNode
			}
		} else {
			t.root = newNode
		}
	}
}

func (t *AVLTreeBitVector) Pop(index int32) int8 {
	n := t.Len()
	if index < 0 {
		index += n
	}
	if index < 0 || index >= n {
		panic(fmt.Sprintf("index out of range: %d", index))
	}
	left, right, size := t.left, t.right, t.size
	bitLen, keys, total := t.bitLen, t.key, t.total
	node := t.root
	d := int32(0)
	t.pathStack = t.pathStack[:0]
	path := t.pathStack
	for node != 0 {
		b32 := int32(bitLen[node])
		t := size[left[node]] + b32
		if t-b32 <= index && index < t {
			break
		}
		path = append(path, node)
		d <<= 1
		if t > index {
			node = left[node]
			d |= 1
		} else {
			node = right[node]
			index -= t
		}
	}
	index -= size[left[node]]
	v := keys[node]
	b32 := int32(bitLen[node])
	res := int32(v >> (b32 - index - 1) & 1)
	if b32 == 1 {
		t._popUnder(path, d, node, res)
		return int8(res)
	}
	keys[node] = ((v >> (b32 - index)) << (b32 - index - 1)) | (v & ((1 << (b32 - index - 1)) - 1))
	bitLen[node]--
	size[node]--
	total[node] -= res
	for _, p := range path {
		size[p]--
		total[p] -= res
	}
	return int8(res)
}

func (t *AVLTreeBitVector) Set(index int32, v int8) {
	n := t.Len()
	if index < 0 {
		index += n
	}
	if index < 0 || index >= n {
		panic(fmt.Sprintf("index out of range: %d", index))
	}

	left, right, bitLen, size, key, total := t.left, t.right, t.bitLen, t.size, t.key, t.total
	node := t.root

	t.pathStack = t.pathStack[:0]
	path := t.pathStack
	for true {
		b32 := int32(bitLen[node])
		tmp := size[left[node]] + b32
		path = append(path, node)
		if tmp-b32 <= index && index < tmp {
			index -= size[left[node]]
			index = b32 - index - 1
			if v == 1 {
				key[node] |= 1 << index
			} else {
				key[node] &= ^(1 << index)
			}
			break
		} else if tmp > index {
			node = left[node]
		} else {
			node = right[node]
			index -= tmp
		}
	}
	for len(path) > 0 {
		node = path[len(path)-1]
		path = path[:len(path)-1]
		total[node] = t._popcount(key[node]) + total[left[node]] + total[right[node]]
	}
}

func (t *AVLTreeBitVector) Get(index int32) int8 {
	if index < 0 {
		index += t.Len()
	}
	left, right, bitLen, size, key := t.left, t.right, t.bitLen, t.size, t.key
	node := t.root
	for true {
		b32 := int32(bitLen[node])
		tmp := size[left[node]] + b32
		if tmp-b32 <= index && index < tmp {
			index -= size[left[node]]
			return int8(key[node] >> (b32 - index - 1) & 1)
		}
		if tmp > index {
			node = left[node]
		} else {
			node = right[node]
			index -= tmp
		}
	}
	panic("unreachable")
}

func (t *AVLTreeBitVector) Count0(end int32) int32 {
	if end < 0 {
		return 0
	}
	if n := t.Len(); end > n {
		end = n
	}
	return end - t._pref(end)
}

func (t *AVLTreeBitVector) Count1(end int32) int32 {
	if end < 0 {
		return 0
	}
	if n := t.Len(); end > n {
		end = n
	}
	return t._pref(end)
}
func (t *AVLTreeBitVector) Count(end int32, v int8) int32 {
	if v == 1 {
		return t.Count1(end)
	}
	return t.Count0(end)
}
func (t *AVLTreeBitVector) Kth0(k int32) int32 {
	n := t.Len()
	if k < 0 || t.Count0(n) <= k {
		return -1
	}
	l, r := int32(0), n
	for r-l > 1 {
		m := (l + r) >> 1
		if m-t._pref(m) > k {
			r = m
		} else {
			l = m
		}
	}
	return l
}
func (t *AVLTreeBitVector) Kth1(k int32) int32 {
	n := t.Len()
	if k < 0 || t.Count1(n) <= k {
		return -1
	}
	l, r := int32(0), n
	for r-l > 1 {
		m := (l + r) >> 1
		if t._pref(m) > k {
			r = m
		} else {
			l = m
		}
	}
	return l
}
func (t *AVLTreeBitVector) Kth(k int32, v int8) int32 {
	if v == 1 {
		return t.Kth1(k)
	}
	return t.Kth0(k)
}
func (t *AVLTreeBitVector) Len() int32 { return t.size[t.root] }

func (t *AVLTreeBitVector) ToList() []int8 {
	if t.root == 0 {
		return nil
	}
	left, right, key, bitLen := t.left, t.right, t.key, t.bitLen
	res := make([]int8, 0, t.Len())
	var rec func(node int32)
	rec = func(node int32) {
		if left[node] != 0 {
			rec(left[node])
		}
		for i := bitLen[node] - 1; i >= 0; i-- {
			res = append(res, int8(key[node]>>i&1))
		}
		if right[node] != 0 {
			rec(right[node])
		}
	}
	rec(t.root)
	return res
}

func (t *AVLTreeBitVector) Debug() {
	left, right, key := t.left, t.right, t.key
	var rec func(node int32) int32
	rec = func(node int32) int32 {
		acc := t._popcount(key[node])
		if left[node] != 0 {
			acc += rec(left[node])
		}
		if right[node] != 0 {
			acc += rec(right[node])
		}
		if acc != t.total[node] {
			// fmt.Println(node, acc, t.total[node])
			panic("error")
		}
		return acc
	}
	rec(t.root)
}

func (t *AVLTreeBitVector) _build(data []int8) {
	n := int32(len(data))
	bit := uint64(bits.Len32(uint32(n)) + 2)
	mask := uint64(1<<bit - 1)
	end := t.end
	t.Reserve(n)
	index := end
	for i := int32(0); i < n; i += W {
		j, v := int32(0), uint64(0)
		for j < W && i+j < n {
			v <<= 1
			v |= uint64(data[i+j])
			j++
		}
		t.key[index] = v
		t.bitLen[index] = int8(j)
		t.size[index] = j
		t.total[index] = t._popcount(v)
		index++
	}
	t.end = index

	var rec func(lr uint64) uint64
	rec = func(lr uint64) uint64 {
		l, r := lr>>bit, lr&mask
		mid := (l + r) >> 1
		hl, hr := uint64(0), uint64(0)
		if l != mid {
			le := rec(l<<bit | mid)
			t.left[mid], hl = int32(le>>bit), le&mask
			t.size[mid] += t.size[t.left[mid]]
			t.total[mid] += t.total[t.left[mid]]
		}
		if mid+1 != r {
			ri := rec((mid+1)<<bit | r)
			t.right[mid], hr = int32(ri>>bit), ri&mask
			t.size[mid] += t.size[t.right[mid]]
			t.total[mid] += t.total[t.right[mid]]
		}
		t.balance[mid] = int8(hl - hr)
		return mid<<bit | (max64(hl, hr) + 1)
	}
	t.root = int32(rec(uint64(end)<<bit|uint64(t.end)) >> bit)
}

func (t *AVLTreeBitVector) _rotateL(node int32) int32 {
	left, right, size, balance, total := t.left, t.right, t.size, t.balance, t.total
	u := left[node]
	size[u] = size[node]
	total[u] = total[node]
	size[node] -= size[left[u]] + int32(t.bitLen[u])
	total[node] -= total[left[u]] + t._popcount(t.key[u])
	left[node] = right[u]
	right[u] = node
	if balance[u] == 1 {
		balance[u] = 0
		balance[node] = 0
	} else {
		balance[u] = -1
		balance[node] = 1
	}
	return u
}

func (t *AVLTreeBitVector) _rotateR(node int32) int32 {
	left, right, size, balance, total := t.left, t.right, t.size, t.balance, t.total
	u := right[node]
	size[u] = size[node]
	total[u] = total[node]
	size[node] -= size[right[u]] + int32(t.bitLen[u])
	total[node] -= total[right[u]] + t._popcount(t.key[u])
	right[node] = left[u]
	left[u] = node
	if balance[u] == -1 {
		balance[u] = 0
		balance[node] = 0
	} else {
		balance[u] = 1
		balance[node] = -1
	}
	return u
}

func (t *AVLTreeBitVector) _rotateLR(node int32) int32 {
	left, right, size, total := t.left, t.right, t.size, t.total
	B := left[node]
	E := right[B]
	size[E] = size[node]
	size[node] -= size[B] - size[right[E]]
	size[B] -= size[right[E]] + int32(t.bitLen[E])
	total[E] = total[node]
	total[node] -= total[B] - total[right[E]]
	total[B] -= total[right[E]] + t._popcount(t.key[E])
	right[B] = left[E]
	left[E] = B
	left[node] = right[E]
	right[E] = node
	t._updateBalance(E)
	return E
}

func (t *AVLTreeBitVector) _rotateRL(node int32) int32 {
	left, right, size, total := t.left, t.right, t.size, t.total
	C := right[node]
	D := left[C]
	size[D] = size[node]
	size[node] -= size[C] - size[left[D]]
	size[C] -= size[left[D]] + int32(t.bitLen[D])
	total[D] = total[node]
	total[node] -= total[C] - total[left[D]]
	total[C] -= total[left[D]] + t._popcount(t.key[D])
	left[C] = right[D]
	right[D] = C
	right[node] = left[D]
	left[D] = node
	t._updateBalance(D)
	return D
}

func (t *AVLTreeBitVector) _updateBalance(node int32) {
	balance := t.balance
	if b := balance[node]; b == 1 {
		balance[t.right[node]] = -1
		balance[t.left[node]] = 0
	} else if b == -1 {
		balance[t.right[node]] = 0
		balance[t.left[node]] = 1
	} else {
		balance[t.right[node]] = 0
		balance[t.left[node]] = 0
	}
	balance[node] = 0
}

func (t *AVLTreeBitVector) _pref(r int32) int32 {
	left, right, bitLen, size, key, total := t.left, t.right, t.bitLen, t.size, t.key, t.total
	node := t.root
	s := int32(0)
	for r > 0 {
		b32 := int32(bitLen[node])
		tmp := size[left[node]] + b32
		if tmp-b32 < r && r <= tmp {
			r -= size[left[node]]
			s += total[left[node]] + t._popcount(key[node]>>(b32-r))
			break
		}
		if tmp > r {
			node = left[node]
		} else {
			s += total[left[node]] + t._popcount(key[node])
			node = right[node]
			r -= tmp
		}
	}
	return s
}

func (t *AVLTreeBitVector) _makeNodeWithBitLen1(v uint64) int32 {
	end := t.end
	if end >= int32(len(t.key)) {
		t.key = append(t.key, v)
		t.bitLen = append(t.bitLen, 1)
		t.size = append(t.size, 1)
		t.total = append(t.total, t._popcount(v))
		t.left = append(t.left, 0)
		t.right = append(t.right, 0)
		t.balance = append(t.balance, 0)
	} else {
		t.key[end] = v
		t.bitLen[end] = 1
		t.size[end] = 1
		t.total[end] = t._popcount(v)
	}
	t.end++
	return end
}

// 这里的path可以不用*[]int32
func (t *AVLTreeBitVector) _popUnder(path []int32, d int32, node int32, res int32) {
	left, right, size, bitLen, balance, keys, total := t.left, t.right, t.size, t.bitLen, t.balance, t.key, t.total
	fd, lmaxTotal, lmaxBitLen := int32(0), int32(0), int8(0)

	if left[node] != 0 && right[node] != 0 {
		path = append(path, node)
		d <<= 1
		d |= 1
		lmax := left[node]
		for right[lmax] != 0 {
			path = append(path, lmax)
			d <<= 1
			fd <<= 1
			fd |= 1
			lmax = right[lmax]
		}
		lmaxTotal = t._popcount(keys[lmax])
		lmaxBitLen = bitLen[lmax]
		keys[node] = keys[lmax]
		bitLen[node] = lmaxBitLen
		node = lmax
	}
	var cNode int32
	if left[node] == 0 {
		cNode = right[node]
	} else {
		cNode = left[node]
	}
	if len(path) > 0 {
		if d&1 == 1 {
			left[path[len(path)-1]] = cNode
		} else {
			right[path[len(path)-1]] = cNode
		}
	} else {
		t.root = cNode
		return
	}
	for len(path) > 0 {
		newNode := int32(0)
		node = path[len(path)-1]
		path = path[:len(path)-1]
		if d&1 == 1 {
			balance[node]--
		} else {
			balance[node]++
		}
		if fd&1 == 1 {
			size[node] -= int32(lmaxBitLen)
			total[node] -= lmaxTotal
		} else {
			size[node]--
			total[node] -= res
		}

		d >>= 1
		fd >>= 1
		if balance[node] == 2 {
			if balance[left[node]] < 0 {
				newNode = t._rotateLR(node)
			} else {
				newNode = t._rotateL(node)
			}
		} else if balance[node] == -2 {
			if balance[right[node]] > 0 {
				newNode = t._rotateRL(node)
			} else {
				newNode = t._rotateR(node)
			}
		} else if balance[node] != 0 {
			break
		}
		if newNode != 0 {
			if len(path) == 0 {
				t.root = newNode
				return
			}
			if d&1 == 1 {
				left[path[len(path)-1]] = newNode
			} else {
				right[path[len(path)-1]] = newNode
			}
			if balance[newNode] != 0 {
				break
			}
		}
	}

	ptr := len(path) - 1
	for ptr >= 0 {
		node := path[ptr]
		ptr--
		if fd&1 == 1 {
			size[node] -= int32(lmaxBitLen)
			total[node] -= lmaxTotal
		} else {
			size[node]--
			total[node] -= res
		}
		fd >>= 1
	}
}

func (t *AVLTreeBitVector) _popcount(v uint64) int32 {
	return int32(bits.OnesCount64(v))
}

func (t *AVLTreeBitVector) _insertAndCount(index int32, v int8) int32 {
	if t.root == 0 {
		t.root = t._makeNodeWithBitLen1(uint64(v))
		return 0
	}
	v32 := int32(v)
	node := t.root
	s := int32(0)

	t.pathStack = t.pathStack[:0]
	path := t.pathStack
	d := int32(0)
	for node != 0 {
		b32 := int32(t.bitLen[node])
		tmp := t.size[t.left[node]] + b32
		if tmp-b32 <= index && index <= tmp {
			break
		}
		if tmp <= index {
			s += t.total[t.left[node]] + t._popcount(t.key[node])
		}
		d <<= 1
		t.size[node]++
		t.total[node] += v32
		path = append(path, node)
		if tmp > index {
			node = t.left[node]
			d |= 1
		} else {
			node = t.right[node]
			index -= tmp
		}
	}

	index -= t.size[t.left[node]]
	b32 := int32(t.bitLen[node])
	s += t.total[t.left[node]] + t._popcount(t.key[node]>>(b32-index))
	if b32 < W {
		mask := t.key[node]
		bl := b32 - index
		t.key[node] = (((mask>>bl)<<1 | uint64(v)) << bl) | (mask & ((1 << bl) - 1))
		t.bitLen[node]++
		t.size[node]++
		t.total[node] += v32
		return s
	}
	path = append(path, node)
	t.size[node]++
	t.total[node] += v32
	mask := t.key[node]
	bl := W - index
	mask = (((mask>>bl)<<1 | uint64(v)) << bl) | (mask & ((1 << bl) - 1))
	leftKey := mask >> W
	leftKeyPopcount := int32(leftKey & 1)
	t.key[node] = mask & ((1 << W) - 1)
	node = t.left[node]
	d <<= 1
	d |= 1
	if node == 0 {
		last := path[len(path)-1]
		if t.bitLen[last] < W8 {
			t.bitLen[last]++
			t.key[last] = (t.key[last] << 1) | leftKey
			return s
		} else {
			t.left[last] = t._makeNodeWithBitLen1(leftKey)
		}
	} else {
		path = append(path, node)
		t.size[node]++
		t.total[node] += leftKeyPopcount
		d <<= 1
		for t.right[node] != 0 {
			node = t.right[node]
			path = append(path, node)
			t.size[node]++
			t.total[node] += leftKeyPopcount
			d <<= 1
		}
		if t.bitLen[node] < W8 {
			t.bitLen[node]++
			t.key[node] = (t.key[node] << 1) | leftKey
			return s
		} else {
			t.right[node] = t._makeNodeWithBitLen1(leftKey)
		}
	}
	newNode := int32(0)
	for len(path) > 0 {
		node = path[len(path)-1]
		path = path[:len(path)-1]
		if d&1 == 1 {
			t.balance[node]++
		} else {
			t.balance[node]--
		}
		d >>= 1
		if t.balance[node] == 0 {
			break
		}
		if t.balance[node] == 2 {
			if t.balance[t.left[node]] == -1 {
				newNode = t._rotateLR(node)
			} else {
				newNode = t._rotateL(node)
			}
			break
		} else if t.balance[node] == -2 {
			if t.balance[t.right[node]] == 1 {
				newNode = t._rotateRL(node)
			} else {
				newNode = t._rotateR(node)
			}
			break
		}
	}
	if newNode != 0 {
		if len(path) > 0 {
			if d&1 == 1 {
				t.left[path[len(path)-1]] = newNode
			} else {
				t.right[path[len(path)-1]] = newNode
			}
		} else {
			t.root = newNode
		}
	}
	return s
}

func (t *AVLTreeBitVector) _accessPopAndCount1(index int32) int32 {
	left, right, size := t.left, t.right, t.size
	bitLen, keys, total := t.bitLen, t.key, t.total
	s := int32(0)
	node := t.root
	d := int32(0)

	t.pathStack = t.pathStack[:0]
	path := t.pathStack
	for node != 0 {
		b32 := int32(bitLen[node])
		tmp := size[left[node]] + b32
		if tmp-b32 <= index && index < tmp {
			break
		}
		if tmp <= index {
			s += total[left[node]] + t._popcount(keys[node])
		}
		path = append(path, node)
		d <<= 1
		if tmp > index {
			node = left[node]
			d |= 1
		} else {
			node = right[node]
			index -= tmp
		}
	}
	index -= size[left[node]]
	b32 := int32(bitLen[node])
	s += total[left[node]] + t._popcount(keys[node]>>(b32-index))
	v := keys[node]
	res := int32(v >> (b32 - index - 1) & 1)
	if b32 == 1 {
		t._popUnder(path, d, node, res)
		return s<<1 | res
	}
	keys[node] = ((v >> (b32 - index)) << (b32 - index - 1)) | (v & ((1 << (b32 - index - 1)) - 1))
	bitLen[node]--
	size[node]--
	total[node] -= res
	for _, p := range path {
		size[p]--
		total[p] -= res
	}
	return s<<1 | res
}

func NewHeap[H any](less func(a, b H) bool, nums []H) *Heap[H] {
	nums = append(nums[:0:0], nums...)
	heap := &Heap[H]{less: less, data: nums}
	heap.heapify()
	return heap
}

type Heap[H any] struct {
	data []H
	less func(a, b H) bool
}

func (h *Heap[H]) Push(value H) {
	h.data = append(h.data, value)
	h.pushUp(h.Len() - 1)
}

func (h *Heap[H]) Pop() (value H) {
	if h.Len() == 0 {
		panic("heap is empty")
	}
	value = h.data[0]
	h.data[0] = h.data[h.Len()-1]
	h.data = h.data[:h.Len()-1]
	h.pushDown(0)
	return
}

func (h *Heap[H]) Top() (value H) {
	value = h.data[0]
	return
}

func (h *Heap[H]) Len() int { return len(h.data) }

func (h *Heap[H]) heapify() {
	n := h.Len()
	for i := (n >> 1) - 1; i > -1; i-- {
		h.pushDown(i)
	}
}

func (h *Heap[H]) pushUp(root int) {
	for parent := (root - 1) >> 1; parent >= 0 && h.less(h.data[root], h.data[parent]); parent = (root - 1) >> 1 {
		h.data[root], h.data[parent] = h.data[parent], h.data[root]
		root = parent
	}
}

func (h *Heap[H]) pushDown(root int) {
	n := h.Len()
	for left := (root<<1 + 1); left < n; left = (root<<1 + 1) {
		right := left + 1
		minIndex := root
		if h.less(h.data[left], h.data[minIndex]) {
			minIndex = left
		}
		if right < n && h.less(h.data[right], h.data[minIndex]) {
			minIndex = right
		}
		if minIndex == root {
			return
		}
		h.data[root], h.data[minIndex] = h.data[minIndex], h.data[root]
		root = minIndex
	}
}

func max32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

func max64(a, b uint64) uint64 {
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

func test() {
	for i := 0; i < 10; i++ {
		nums := make([]int, 1000)
		for j := 0; j < 1000; j++ {
			nums[j] = rand.Intn(1000)
		}
		wm := NewWaveletMatrixDynamic(int32(len(nums)), func(i int32) int { return nums[i] }, maxs(nums))

		prefixBf := func(end int32, x int) int32 {
			res := int32(0)
			for i := int32(0); i < end; i++ {
				if nums[i] == x {
					res++
				}
			}
			return res
		}

		rangeBf := func(start, end int32, x int) int32 {
			res := int32(0)
			for i := start; i < end; i++ {
				if nums[i] == x {
					res++
				}
			}
			return res
		}

		rangeFreqBf := func(start, end int32, x, y int) int32 {
			res := int32(0)
			for i := start; i < end; i++ {
				if nums[i] >= x && nums[i] < y {
					res++
				}
			}
			return res
		}

		kthBf := func(k int32, x int) int32 {
			cnt := int32(0)
			for i := int32(0); i < int32(len(nums)); i++ {
				if nums[i] == x {
					if cnt == k {
						return i
					}
					cnt++
				}
			}
			return -1
		}

		kthSmallestBf := func(start, end, k int32) int {
			arr := make([]int, 0, end-start)
			for i := start; i < end; i++ {
				arr = append(arr, nums[i])
			}
			sort.Ints(arr)
			if int(k) >= len(arr) {
				return -1
			}
			return arr[k]
		}

		kthSmallestIndexBf := func(start, end, k int32) int32 {
			arrWithIndex := make([][2]int, 0, end-start)
			for i := start; i < end; i++ {
				arrWithIndex = append(arrWithIndex, [2]int{nums[i], int(i)})
			}
			sort.Slice(arrWithIndex, func(i, j int) bool {
				if arrWithIndex[i][0] != arrWithIndex[j][0] {
					return arrWithIndex[i][0] < arrWithIndex[j][0]
				}
				return arrWithIndex[i][1] < arrWithIndex[j][1]
			})
			if int(k) >= len(arrWithIndex) {
				return -1
			}
			return int32(arrWithIndex[k][1])
		}

		kthLargestBf := func(start, end, k int32) int {
			arr := make([]int, 0, end-start)
			for i := start; i < end; i++ {
				arr = append(arr, nums[i])
			}
			sort.Ints(arr)
			if int(k) >= len(arr) {
				return -1
			}
			return arr[len(arr)-1-int(k)]
		}

		topKBf := func(start, end, k int32) []topKPair {
			m := map[int]int32{}
			for i := start; i < end; i++ {
				m[nums[i]]++
			}
			arr := make([]topKPair, 0, len(m))
			for k, v := range m {
				arr = append(arr, topKPair{k, v})
			}
			sort.Slice(arr, func(i, j int) bool {
				if arr[i].count != arr[j].count {
					return arr[i].count > arr[j].count
				}
				return arr[i].value < arr[j].value
			})
			if int(k) >= len(arr) {
				return arr
			}
			return arr[:k]
		}

		floorBf := func(start, end int32, x int) (int, bool) {
			res := math.MinInt
			for i := start; i < end; i++ {
				if nums[i] <= x && nums[i] > res {
					res = nums[i]
				}
			}
			if res == math.MinInt {
				return -1, false
			}
			return res, true
		}

		lowerBf := func(start, end int32, x int) (int, bool) {
			res := math.MinInt
			for i := start; i < end; i++ {
				if nums[i] < x && nums[i] > res {
					res = nums[i]
				}
			}
			if res == math.MinInt {
				return -1, false
			}
			return res, true
		}

		ceilBf := func(start, end int32, x int) (int, bool) {
			res := math.MaxInt
			for i := start; i < end; i++ {
				if nums[i] >= x && nums[i] < res {
					res = nums[i]
				}
			}
			if res == math.MaxInt {
				return -1, false
			}
			return res, true
		}

		higherBf := func(start, end int32, x int) (int, bool) {
			res := math.MaxInt
			for i := start; i < end; i++ {
				if nums[i] > x && nums[i] < res {
					res = nums[i]
				}
			}

			if res == math.MaxInt {
				return -1, false
			}
			return res, true
		}

		countAllBf := func(start, end int32, v int) (int32, int32, int32) {
			rank, rankLessThan, rankMoreThan := int32(0), int32(0), int32(0)
			for i := start; i < end; i++ {
				if nums[i] == v {
					rank++
				}
				if nums[i] < v {
					rankLessThan++
				}
				if nums[i] > v {
					rankMoreThan++
				}
			}
			return rank, rankLessThan, rankMoreThan
		}

		countLessBf := func(start, end int32, x int) int32 {
			res := int32(0)
			for i := start; i < end; i++ {
				if nums[i] < x {
					res++
				}
			}
			return res
		}

		countMoreBf := func(start, end int32, x int) int32 {
			res := int32(0)
			for i := start; i < end; i++ {
				if nums[i] > x {
					res++
				}
			}
			return res
		}

		insertBf := func(index int32, v int) {
			nums = append(nums, 0)
			copy(nums[index+1:], nums[index:])
			nums[index] = v
		}

		popBf := func(index int32) int {
			res := nums[index]
			nums = append(nums[:index], nums[index+1:]...)
			return res
		}

		setBf := func(index int32, v int) {
			nums[index] = v
		}

		for j := 0; j < 100; j++ {
			start, end := rand.Intn(1000), rand.Intn(1000)
			if start > end {
				start, end = end, start
			}
			x := rand.Intn(1000)
			if prefixBf(int32(end), x) != wm.PrefixCount(int32(end), x) {
				fmt.Println(prefixBf(int32(end), x), wm.PrefixCount(int32(end), x), end, x)
				panic("prefixBf")
			}
			if rangeBf(int32(start), int32(end), x) != wm.RangeCount(int32(start), int32(end), x) {
				panic("rangeBf")
			}

			y := rand.Intn(1000)
			if res1, res2 := rangeFreqBf(int32(start), int32(end), x, y), wm.RangeFreq(int32(start), int32(end), x, y); res1 != res2 {
				fmt.Println(res1, res2, start, end, x, y)
				panic("rangeFreqBf")
			}

			k := rand.Intn(1000)
			if res1, res2 := kthBf(int32(k), x), wm.Kth(int32(k), x); res1 != res2 {
				fmt.Println(res1, res2, k, x)
				panic("kthBf")
			}

			if res1, res2 := kthSmallestBf(int32(start), int32(end), int32(k)), wm.KthSmallest(int32(start), int32(end), int32(k)); res1 != res2 {
				fmt.Println(res1, res2, start, end, k)
				panic("kthSmallestBf")
			}

			if res1, res2 := kthSmallestIndexBf(int32(start), int32(end), int32(k)), wm.KthSmallestIndex(int32(start), int32(end), int32(k)); res1 != res2 {
				fmt.Println(res1, res2, start, end, k)
				panic("kthSmallestIndexBf")
			}

			if res1, res2 := kthLargestBf(int32(start), int32(end), int32(k)), wm.KthLargest(int32(start), int32(end), int32(k)); res1 != res2 {
				fmt.Println(res1, res2, start, end, k)
				panic("kthLargestBf")
			}

			topK1 := topKBf(int32(start), int32(end), int32(k))
			topK2 := wm.TopK(int32(start), int32(end), int32(k))
			if len(topK1) != len(topK2) {
				fmt.Println(len(topK1), len(topK2), start, end, k)
				panic("topKBf")
			}
			for i := 0; i < len(topK1); i++ {
				if topK1[i].count != topK2[i].count {
					fmt.Println(topK1[i], topK2[i], start, end, k)
					panic("topKBf")
				}
			}

			funcs1 := []func(int32, int32, int) (int, bool){floorBf, lowerBf, ceilBf, higherBf}
			funcs2 := []func(int32, int32, int) (int, bool){wm.Floor, wm.Lower, wm.Ceil, wm.Higher}
			for i := 0; i < len(funcs1); i++ {
				res1, ok1 := funcs1[i](int32(start), int32(end), x)
				res2, ok2 := funcs2[i](int32(start), int32(end), x)
				if res1 != res2 || ok1 != ok2 {
					fmt.Println(res1, res2, start, end, x)
					panic("funcs")
				}
			}

			c1, c2, c3 := countAllBf(int32(start), int32(end), x)
			c4, c5, c6 := wm.CountAll(int32(start), int32(end), x)
			if c1 != c4 || c2 != c5 || c3 != c6 {
				fmt.Println(c1, c4, c2, c5, c3, c6, start, end, x)
				panic("countAllBf")
			}

			if res1, res2 := countLessBf(int32(start), int32(end), x), wm.CountLess(int32(start), int32(end), x); res1 != res2 {
				fmt.Println(res1, res2, start, end, x)
				panic("countLessBf")
			}

			if res1, res2 := countMoreBf(int32(start), int32(end), x), wm.CountMore(int32(start), int32(end), x); res1 != res2 {
				fmt.Println(res1, res2, start, end, x)
				panic("countMoreBf")
			}

			insertIndex := int32(rand.Intn(1000))
			insertValue := rand.Intn(1000)
			insertBf(insertIndex, insertValue)
			wm.Insert(insertIndex, insertValue)

			// _ = popBf
			popIndex := int32(rand.Intn(len(nums)))
			if popBf(popIndex) != wm.Pop(popIndex) {
				panic("popBf")
			}

			setIndex := int32(rand.Intn(len(nums)))
			setValue := rand.Intn(1000)
			setBf(setIndex, setValue)
			wm.Set(setIndex, setValue)
		}
	}

	fmt.Println("pass")
}

func testTime() {
	n := int32(1e5)
	nums := make([]int, n)
	for i := int32(0); i < n; i++ {
		nums[i] = rand.Intn(1e9)
	}

	wm := NewWaveletMatrixDynamic(int32(len(nums)), func(i int32) int { return nums[i] }, maxs(nums))
	// fmt.Println("build time:", time.Since(time1))
	time1 := time.Now()

	for i := int32(0); i < n; i++ {
		wm.PrefixCount(i, nums[i])
		wm.RangeCount(0, i, nums[i])
		wm.RangeFreq(0, i, nums[i], nums[i]+1)
		wm.Kth(i, nums[i])
		wm.KthSmallest(0, i, i)
		wm.KthSmallestIndex(0, i, i)
		wm.KthLargest(0, i, i)
		wm.Floor(0, i, nums[i])
		wm.Lower(0, i, nums[i])
		wm.Ceil(0, i, nums[i])
		wm.Higher(0, i, nums[i])
		wm.CountAll(0, i, nums[i])
		wm.CountLess(0, i, nums[i])
		wm.CountMore(0, i, nums[i])
		wm.CountSame(0, i, nums[i])
		wm.Set(i, nums[i])
		wm.Insert(i, nums[i])
		wm.Pop(i)
	}

	fmt.Println(time.Since(time1)) // 9.2127399s

	time1 = time.Now()
	wm.TopK(0, n, n/2)
	fmt.Println(time.Since(time1)) // 102.8108ms
}
