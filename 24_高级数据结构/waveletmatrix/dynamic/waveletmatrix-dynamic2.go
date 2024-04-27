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

	CF455D()
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

// 1e5 -> 150 , 2e5 -> 250
const _LOAD int32 = 200 // 块尺寸越大，修改越快，查询越慢

// 维护[0,maxValue].
type WaveletMatrixDynamic struct {
	size     int32
	maxValue int
	bitLen   int32
	mid      []int32
	bv       []*DynamicBitvector
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
	res.bv = make([]*DynamicBitvector, res.bitLen)
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
			s := wm.bv[bit]._insertAndCount1(index, 1)
			index = s + wm.mid[bit]
		} else {
			s := wm.bv[bit]._insertAndCount1(index, 0)
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
		wm.bv[bit] = NewDynamicBitvector(n, func(i int32) int8 { return bits[i] })
		wm.mid[bit] = p
		zero, data = data, zero
		copy(data[p:], one[:q])
	}
}

// 使用分块+树状数组维护的动态Bitvector.
type DynamicBitvector struct {
	size              int32
	totalOnes         int32
	blocks            [][]int8
	blockOnes         []int32
	preLen            []int32 // 每个块块长的前缀和
	preOnes           []int32 // 每个块01的前缀和
	shouldRebuildTree bool
}

func NewDynamicBitvector(n int32, f func(i int32) int8) *DynamicBitvector {
	blocks := make([][]int8, 0, n/_LOAD+1)
	blockOnes := make([]int32, 0, n/_LOAD+1)
	totalOnes := int32(0)
	for start := int32(0); start < n; start += _LOAD {
		end := start + _LOAD
		if end > n {
			end = n
		}
		block := make([]int8, end-start)
		ones := int32(0)
		for i := start; i < end; i++ {
			block[i-start] = f(i)
			ones += int32(block[i-start])
		}
		blocks = append(blocks, block)
		blockOnes = append(blockOnes, ones)
		totalOnes += ones
	}
	res := &DynamicBitvector{
		size:              n,
		totalOnes:         totalOnes,
		blocks:            blocks,
		blockOnes:         blockOnes,
		shouldRebuildTree: true,
	}
	return res
}

func (sl *DynamicBitvector) _accessPopAndCount1(index int32) (res int32) {
	pos, startIndex, preOnes := sl._findKthOneAndPreOnes(index)
	block := sl.blocks[pos]
	if startIndex < int32(len(block))>>1 {
		// 统计前缀
		onesInBlockPrefix := int16(0)
		for i := int32(0); i < startIndex; i++ {
			onesInBlockPrefix += int16(block[i])
		}
		res = preOnes + int32(onesInBlockPrefix)
	} else {
		onesInBlockSuffix := int16(0)
		for i := int32(len(block) - 1); i >= startIndex; i-- {
			onesInBlockSuffix += int16(block[i])
		}
		res = preOnes + sl.blockOnes[pos] - int32(onesInBlockSuffix)
	}
	res = res<<1 | int32(block[startIndex])

	// pop
	// !delete element
	sl.size--
	if block[startIndex] == 1 {
		sl._updatePreLenAndPreOnes(pos, false)
		sl.blockOnes[pos]--
		sl.totalOnes--
	} else {
		sl._updatePreLen(pos, false)
	}
	copy(sl.blocks[pos][startIndex:], sl.blocks[pos][startIndex+1:])
	sl.blocks[pos] = sl.blocks[pos][:len(sl.blocks[pos])-1]
	if len(sl.blocks[pos]) == 0 {
		// !delete block
		copy(sl.blocks[pos:], sl.blocks[pos+1:])
		sl.blocks = sl.blocks[:len(sl.blocks)-1]
		copy(sl.blockOnes[pos:], sl.blockOnes[pos+1:])
		sl.blockOnes = sl.blockOnes[:len(sl.blockOnes)-1]
		sl.shouldRebuildTree = true
	}
	return
}

// index前插入value，并统计插入后[0, index) 中 value 的个数.
func (sl *DynamicBitvector) _insertAndCount1(index int32, value int8) (res int32) {
	if len(sl.blocks) == 0 {
		sl.blocks = append(sl.blocks, []int8{value})
		sl.blockOnes = append(sl.blockOnes, int32(value))
		sl.shouldRebuildTree = true
		sl.size++
		if value == 1 {
			sl.totalOnes++
		}
		return 0
	}

	pos, startIndex, preOnes := sl._findKthOneAndPreOnes(index)
	block := sl.blocks[pos]
	if startIndex < int32(len(block))>>1 {
		// 统计前缀
		onesInBlockPrefix := int16(0)
		for i := int32(0); i < startIndex; i++ {
			onesInBlockPrefix += int16(block[i])
		}
		res = preOnes + int32(onesInBlockPrefix)
	} else {
		onesInBlockSuffix := int16(0)
		for i := int32(len(block) - 1); i >= startIndex; i-- {
			onesInBlockSuffix += int16(block[i])
		}
		res = preOnes + sl.blockOnes[pos] - int32(onesInBlockSuffix)
	}

	if value == 1 {
		sl._updatePreLenAndPreOnes(pos, true)
		sl.blockOnes[pos]++
		sl.totalOnes++
	} else {
		sl._updatePreLen(pos, true)
	}
	sl.blocks[pos] = append(sl.blocks[pos], 0)
	copy(sl.blocks[pos][startIndex+1:], sl.blocks[pos][startIndex:])
	sl.blocks[pos][startIndex] = value

	// n -> load + (n - load)
	if n := int32(len(sl.blocks[pos])); _LOAD+_LOAD < n {
		totalOnes := sl.blockOnes[pos]
		sl.blocks = append(sl.blocks, nil)
		copy(sl.blocks[pos+2:], sl.blocks[pos+1:])
		sl.blocks[pos+1] = sl.blocks[pos][_LOAD:] // !注意max的设置(为了让左右互不影响)
		sl.blocks[pos] = sl.blocks[pos][:_LOAD:_LOAD]
		ones1 := int32(0)
		for _, v := range sl.blocks[pos] {
			ones1 += int32(v)
		}
		ones2 := totalOnes - ones1
		sl.blockOnes = append(sl.blockOnes, 0)
		copy(sl.blockOnes[pos+2:], sl.blockOnes[pos+1:])
		sl.blockOnes[pos] = ones1
		sl.blockOnes[pos+1] = ones2
		sl.shouldRebuildTree = true
	}

	sl.size++
	return
}

func (sl *DynamicBitvector) Insert(index int32, value int8) {
	if len(sl.blocks) == 0 {
		sl.blocks = append(sl.blocks, []int8{value})
		sl.blockOnes = append(sl.blockOnes, int32(value))
		sl.shouldRebuildTree = true
		sl.size++
		if value == 1 {
			sl.totalOnes++
		}
		return
	}

	pos, startIndex := sl._findKth(index)
	if value == 1 {
		sl._updatePreLenAndPreOnes(pos, true)
		sl.blockOnes[pos]++
		sl.totalOnes++
	} else {
		sl._updatePreLen(pos, true)
	}
	sl.blocks[pos] = append(sl.blocks[pos], 0)
	copy(sl.blocks[pos][startIndex+1:], sl.blocks[pos][startIndex:])
	sl.blocks[pos][startIndex] = value

	// n -> load + (n - load)
	if n := int32(len(sl.blocks[pos])); _LOAD+_LOAD < n {
		totalOnes := sl.blockOnes[pos]
		sl.blocks = append(sl.blocks, nil)
		copy(sl.blocks[pos+2:], sl.blocks[pos+1:])
		sl.blocks[pos+1] = sl.blocks[pos][_LOAD:] // !注意max的设置(为了让左右互不影响)
		sl.blocks[pos] = sl.blocks[pos][:_LOAD:_LOAD]
		ones1 := int32(0)
		for _, v := range sl.blocks[pos] {
			ones1 += int32(v)
		}
		ones2 := totalOnes - ones1
		sl.blockOnes = append(sl.blockOnes, 0)
		copy(sl.blockOnes[pos+2:], sl.blockOnes[pos+1:])
		sl.blockOnes[pos] = ones1
		sl.blockOnes[pos+1] = ones2
		sl.shouldRebuildTree = true
	}

	sl.size++
	return
}

func (sl *DynamicBitvector) Pop(index int32) int8 {
	pos, startIndex := sl._findKth(index)
	value := sl.blocks[pos][startIndex]
	// !delete element
	sl.size--
	if value == 1 {
		sl._updatePreLenAndPreOnes(pos, false)
		sl.blockOnes[pos]--
		sl.totalOnes--
	} else {
		sl._updatePreLen(pos, false)
	}

	copy(sl.blocks[pos][startIndex:], sl.blocks[pos][startIndex+1:])
	sl.blocks[pos] = sl.blocks[pos][:len(sl.blocks[pos])-1]

	if len(sl.blocks[pos]) == 0 {
		// !delete block
		copy(sl.blocks[pos:], sl.blocks[pos+1:])
		sl.blocks = sl.blocks[:len(sl.blocks)-1]
		copy(sl.blockOnes[pos:], sl.blockOnes[pos+1:])
		sl.blockOnes = sl.blockOnes[:len(sl.blockOnes)-1]
		sl.shouldRebuildTree = true
	}
	return value
}

func (sl *DynamicBitvector) Get(index int32) int8 {
	pos, startIndex := sl._findKth(index)
	return sl.blocks[pos][startIndex]
}

func (sl *DynamicBitvector) Set(index int32, value int8) {
	pos, startIndex := sl._findKth(index)
	curValue := sl.blocks[pos][startIndex]
	if curValue == value {
		return
	}
	sl.blocks[pos][startIndex] = value
	if value == 1 {
		sl._updatePreOnes(pos, true)
		sl.blockOnes[pos]++
		sl.totalOnes++
	} else {
		sl._updatePreOnes(pos, false)
		sl.blockOnes[pos]--
		sl.totalOnes--
	}
}

func (sl *DynamicBitvector) Count0(end int32) int32 {
	if end <= 0 {
		return 0
	}
	if end > sl.size {
		end = sl.size
	}
	return end - sl.Count1(end)
}

func (sl *DynamicBitvector) Count1(end int32) int32 {
	if end <= 0 {
		return 0
	}
	if end > sl.size {
		end = sl.size
	}
	if sl.shouldRebuildTree {
		sl._buildPreLenAndPreOnes()
	}
	pos, startIndex, preOnes := sl._findKthOneAndPreOnes(end - 1)
	block := sl.blocks[pos]
	if startIndex < int32(len(block))>>1 {
		// 统计前缀
		onesInBlockPrefix := int16(0)
		for i := int32(0); i <= startIndex; i++ {
			onesInBlockPrefix += int16(block[i])
		}
		return preOnes + int32(onesInBlockPrefix)
	} else {
		onesInBlockSuffix := int16(0)
		for i := int32(len(block) - 1); i > startIndex; i-- {
			onesInBlockSuffix += int16(block[i])
		}
		return preOnes + sl.blockOnes[pos] - int32(onesInBlockSuffix)
	}
}

func (sl *DynamicBitvector) Count(end int32, value int8) int32 {
	if value == 1 {
		return sl.Count1(end)
	}
	return end - sl.Count1(end)
}

func (sl *DynamicBitvector) Kth0(k int32) int32 {
	if k < 0 || sl.size-sl.totalOnes <= k {
		return -1
	}
	if sl.shouldRebuildTree {
		sl._buildPreLenAndPreOnes()
	}
	pos, preLen, newK := sl._findKthZero(k)
	block := sl.blocks[pos]
	for i := int32(0); i < int32(len(block)); i++ {
		if block[i] == 0 {
			if newK == 0 {
				return preLen + i
			}
			newK--
		}
	}
	return -1
}

func (sl *DynamicBitvector) Kth1(k int32) int32 {
	if k < 0 || sl.totalOnes <= k {
		return -1
	}
	if sl.shouldRebuildTree {
		sl._buildPreLenAndPreOnes()
	}
	pos, preLen, newK := sl._findKthOne(k)
	block := sl.blocks[pos]
	for i := int32(0); i < int32(len(block)); i++ {
		if block[i] == 1 {
			if newK == 0 {
				return preLen + i
			}
			newK--
		}
	}
	return -1
}

func (sl *DynamicBitvector) Kth(k int32, value int8) int32 {
	if value == 1 {
		return sl.Kth1(k)
	}
	return sl.Kth0(k)
}

func (sl *DynamicBitvector) Len() int32 {
	return sl.size
}

func (sl *DynamicBitvector) GetAll() []int8 {
	res := make([]int8, 0, sl.size)
	for _, block := range sl.blocks {
		res = append(res, block...)
	}
	return res
}

func (sl *DynamicBitvector) _buildPreLenAndPreOnes() {
	sl.preLen = make([]int32, len(sl.blocks))
	for i := 0; i < len(sl.blocks); i++ {
		sl.preLen[i] = int32(len(sl.blocks[i]))
	}
	sl.preOnes = append(sl.blockOnes[:0:0], sl.blockOnes...)
	tree1, tree2 := sl.preLen, sl.preOnes
	m := int32(len(tree1))
	for i := int32(0); i < m; i++ {
		j := i | (i + 1)
		if j < m {
			tree1[j] += tree1[i]
			tree2[j] += tree2[i]
		}
	}
	sl.shouldRebuildTree = false
}

func (sl *DynamicBitvector) _updatePreLen(index int32, addOne bool) {
	if sl.shouldRebuildTree {
		return
	}
	tree := sl.preLen
	m := int32(len(tree))
	if addOne {
		for i := index; i < m; i |= i + 1 {
			tree[i]++
		}
	} else {
		for i := index; i < m; i |= i + 1 {
			tree[i]--
		}
	}
}

func (sl *DynamicBitvector) _updatePreOnes(index int32, addOne bool) {
	if sl.shouldRebuildTree {
		return
	}
	tree := sl.preOnes
	m := int32(len(tree))
	if addOne {
		for i := index; i < m; i |= i + 1 {
			tree[i]++
		}
	} else {
		for i := index; i < m; i |= i + 1 {
			tree[i]--
		}
	}
}

func (sl *DynamicBitvector) _updatePreLenAndPreOnes(index int32, addOne bool) {
	if sl.shouldRebuildTree {
		return
	}
	tree1, tree2 := sl.preLen, sl.preOnes
	m := int32(len(tree1))
	if addOne {
		for i := index; i < m; i |= i + 1 {
			tree1[i]++
			tree2[i]++
		}
	} else {
		for i := index; i < m; i |= i + 1 {
			tree1[i]--
			tree2[i]--
		}
	}
}

func (sl *DynamicBitvector) _findKth(k int32) (pos, index int32) {
	if k < int32(len(sl.blocks[0])) {
		return 0, k
	}
	last := int32(len(sl.blocks) - 1)
	lastLen := int32(len(sl.blocks[last]))
	if k >= sl.size {
		return last, lastLen
	}
	if k >= sl.size-lastLen {
		return last, k + lastLen - sl.size
	}
	if sl.shouldRebuildTree {
		sl._buildPreLenAndPreOnes()
	}
	tree := sl.preLen
	pos = -1
	m := int32(len(tree))
	bitLen := int8(bits.Len32(uint32(m)))
	for d := bitLen - 1; d >= 0; d-- {
		next := pos + (1 << d)
		if next < m && k >= tree[next] {
			pos = next
			k -= tree[pos]
		}
	}
	return pos + 1, k
}

func (sl *DynamicBitvector) _findKthOneAndPreOnes(k int32) (pos, index, ones int32) {
	if k < int32(len(sl.blocks[0])) {
		return 0, k, 0
	}
	last := int32(len(sl.blocks) - 1)
	lastLen := int32(len(sl.blocks[last]))
	if k >= sl.size-lastLen {
		return last, k + lastLen - sl.size, sl.totalOnes - sl.blockOnes[last]
	}
	if sl.shouldRebuildTree {
		sl._buildPreLenAndPreOnes()
	}
	tree1, tree2 := sl.preLen, sl.preOnes
	pos = -1
	m := int32(len(tree1))
	bitLen := int8(bits.Len32(uint32(m)))
	for d := bitLen - 1; d >= 0; d-- {
		next := pos + (1 << d)
		if next < m && k >= tree1[next] {
			pos = next
			k -= tree1[pos]
			ones += tree2[pos]
		}
	}
	return pos + 1, k, ones
}

func (sl *DynamicBitvector) _findKthOne(k int32) (pos, preLen, newK int32) {
	if k < sl.blockOnes[0] {
		return 0, 0, k
	}
	last := int32(len(sl.blocks) - 1)
	lastOnes := sl.blockOnes[last]
	if tmp := sl.totalOnes - lastOnes; k >= tmp {
		return last, sl.size - int32(len(sl.blocks[last])), k - tmp
	}
	if sl.shouldRebuildTree {
		sl._buildPreLenAndPreOnes()
	}
	tree1, tree2 := sl.preLen, sl.preOnes
	pos = -1
	m := int32(len(tree2))
	bitLen := int8(bits.Len32(uint32(m)))
	newK = k
	for d := bitLen - 1; d >= 0; d-- {
		next := pos + (1 << d)
		if next < m && newK >= tree2[next] {
			pos = next
			preLen += tree1[pos]
			newK -= tree2[pos]
		}
	}
	pos++
	return
}

func (sl *DynamicBitvector) _findKthZero(k int32) (pos, preLen, newK int32) {
	if k < int32(len(sl.blocks[0]))-sl.blockOnes[0] {
		return 0, 0, k
	}
	last := int32(len(sl.blocks) - 1)
	lastLen := int32(len(sl.blocks[last]))
	lastZero := lastLen - sl.blockOnes[last]
	if tmp := sl.size - lastLen + lastZero; k >= tmp {
		return last, sl.size - lastLen, k - tmp
	}
	if sl.shouldRebuildTree {
		sl._buildPreLenAndPreOnes()
	}
	tree1, tree2 := sl.preLen, sl.preOnes
	pos = -1
	m := int32(len(tree2))
	bitLen := int8(bits.Len32(uint32(m)))
	newK = k
	for d := bitLen - 1; d >= 0; d-- {
		next := pos + (1 << d)
		if next < m && newK >= tree1[next]-tree2[next] {
			pos = next
			preLen += tree1[pos]
			newK -= tree1[pos] - tree2[pos]
		}
	}
	pos++
	return
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

func min32(a, b int32) int32 {
	if a < b {
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

	for i := 0; i < 20; i++ {
		nums := make([]int, 3000)
		for j := 0; j < 3000; j++ {
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

		for j := 0; j < 2000; j++ {
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
			_ = topKBf

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
			_ = insertBf

			popIndex := int32(rand.Intn(len(nums)))
			if res1, res2 := popBf(popIndex), wm.Pop(popIndex); res1 != res2 {
				fmt.Println(res1, res2, popIndex)
				panic("popBf")
			}
			_ = popBf

			setIndex := int32(rand.Intn(len(nums)))
			setValue := rand.Intn(1000)
			setBf(setIndex, setValue)
			wm.Set(setIndex, setValue)
			_ = setBf
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
