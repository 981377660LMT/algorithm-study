// WaveletMatrixStatic/StaticWaveletMatrix
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

// !值域较大时，可以预先离散化以减少值域大小.

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
	yosupo()
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
	wm := NewWaveletMatrixStatic(int32(len(nums)), func(i int32) int { return int(newNums[i]) }, len(origin))

	for i := 0; i < q; i++ {
		var start, end, x int32
		fmt.Fscan(in, &start, &end, &x)
		res := wm.KthSmallest(start, end, x)
		fmt.Fprintln(out, origin[res])
	}
}

func demo() {
	nums := []int{3, 1, 4, 1, 5, 9, 2, 6}
	wm := NewWaveletMatrixStatic(int32(len(nums)), func(i int32) int { return nums[i] }, maxs(nums)+1)
	fmt.Println(wm.PrefixCount(10, 2))
	fmt.Println(wm.Kth(0, 2))
	fmt.Println(wm.KthSmallest(1, 4, 2))
	fmt.Println(wm.TopK(0, 8, 3))
	fmt.Println(wm.RangeFreq(0, 8, 1, 4))
	fmt.Println(wm.RangeCount(0, 8, 1))
	fmt.Println(wm.RangeCount(0, 8, 2))
	fmt.Println(wm.Floor(0, 8, 3))
	fmt.Println(wm.Ceil(0, 8, 3))
	fmt.Println(wm.Lower(0, 8, 3))
	fmt.Println(wm.Higher(0, 8, 3))
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
type WaveletMatrixStatic struct {
	size     int32
	maxValue int
	bitLen   int32
	mid      []int32
	bv       []*bitVector
}

type topKPair = struct {
	value int
	count int32
}

func NewWaveletMatrixStatic(n int32, f func(int32) int, maxValue int) *WaveletMatrixStatic {
	if maxValue <= 0 {
		maxValue = 1
	}
	res := &WaveletMatrixStatic{
		size:     n,
		maxValue: maxValue,
		bitLen:   int32(bits.Len(uint(maxValue))),
	}
	res.mid = make([]int32, res.bitLen)
	res.bv = make([]*bitVector, res.bitLen)
	if n > 0 {
		res._build(n, f)
	}
	return res
}

func (wm *WaveletMatrixStatic) PrefixCount(end int32, x int) int32 {
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

func (wm *WaveletMatrixStatic) RangeCount(start, end int32, x int) int32 {
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

func (wm *WaveletMatrixStatic) RangeFreq(start, end int32, floor, higher int) int32 {
	if floor >= higher {
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
func (wm *WaveletMatrixStatic) Kth(k int32, x int) int32 {
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

func (wm *WaveletMatrixStatic) KthSmallest(start, end int32, k int32) int {
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

func (wm *WaveletMatrixStatic) KthSmallestIndex(start, end int32, k int32) int32 {
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

func (wm *WaveletMatrixStatic) KthLargest(start, end int32, k int32) int {
	return wm.KthSmallest(start, end, end-start-k-1)
}
func (wm *WaveletMatrixStatic) KthLargestIndex(start, end int32, k int32) int32 {
	return wm.KthSmallestIndex(start, end, end-start-k-1)
}

// 按照出现次数排序，返回出现次数最多的k个元素.
func (wm *WaveletMatrixStatic) TopK(start, end int32, k int32) (res []topKPair) {
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

func (wm *WaveletMatrixStatic) Floor(start, end int32, x int) (int, bool) {
	same, less, _ := wm.CountAll(start, end, x)
	if same > 0 {
		return x, true
	}
	if less == 0 {
		return -1, false
	}
	return wm.KthSmallest(start, end, less-1), true
}
func (wm *WaveletMatrixStatic) Lower(start, end int32, x int) (int, bool) {
	less := wm.CountLess(start, end, x)
	if less == 0 {
		return -1, false
	}
	return wm.KthSmallest(start, end, less-1), true
}

func (wm *WaveletMatrixStatic) Ceil(start, end int32, x int) (int, bool) {
	same, less, _ := wm.CountAll(start, end, x)
	if same > 0 {
		return x, true
	}
	if less == end-start {
		return -1, false
	}
	return wm.KthSmallest(start, end, less), true
}

func (wm *WaveletMatrixStatic) Higher(start, end int32, x int) (int, bool) {
	less := wm.CountLess(start, end, x+1)
	if less == end-start {
		return -1, false
	}
	return wm.KthSmallest(start, end, less), true
}

// 返回[start, end)中等于x的个数，小于x的个数，大于x的个数.
func (wm *WaveletMatrixStatic) CountAll(start, end int32, x int) (same, less, more int32) {
	if start < 0 {
		start = 0
	}
	if end > wm.size {
		end = wm.size
	}
	if start >= end {
		return 0, 0, 0
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

func (wm *WaveletMatrixStatic) Get(index int32) int {
	if index < 0 {
		index += wm.size
	}
	res := 0
	for bit := wm.bitLen - 1; bit >= 0; bit-- {
		if wm.bv[bit].Get(index) {
			res |= 1 << bit
			index = wm.bv[bit].Count1(index) + wm.mid[bit]
		} else {
			index = wm.bv[bit].Count0(index)
		}
	}
	return res
}

// 区间[start, end)中小于x的元素个数.
func (wm *WaveletMatrixStatic) CountLess(start, end int32, x int) int32 {
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
func (wm *WaveletMatrixStatic) CountMore(start, end int32, x int) int32 {
	return end - start - wm.CountLess(start, end, x+1)
}

func (wm *WaveletMatrixStatic) CountSame(start, end int32, x int) int32 {
	same, _, _ := wm.CountAll(start, end, x)
	return same
}

func (wm *WaveletMatrixStatic) Len() int32 {
	return wm.size
}

func (wm *WaveletMatrixStatic) _build(n int32, f func(int32) int) {
	data := make([]int, n)
	for i := int32(0); i < n; i++ {
		data[i] = f(i)
	}
	zero, one := make([]int, n), make([]int, n)
	for bit := wm.bitLen - 1; bit >= 0; bit-- {
		wm.bv[bit] = newBitVector(n)
		v := wm.bv[bit]
		p, q := int32(0), int32(0)
		for i, e := range data {
			if e>>bit&1 == 1 {
				v.Set(int32(i))
				one[q] = e
				q++
			} else {
				zero[p] = e
				p++
			}
		}
		v.Build()
		wm.mid[bit] = p
		zero, data = data, zero
		copy(data[p:], one[:q])
	}
}

type bitVector struct {
	n      int32
	size   int32
	bit    []uint64
	preSum []int32
}

func newBitVector(n int32) *bitVector {
	size := (n + 63) >> 6
	bit := make([]uint64, size+1)
	preSum := make([]int32, size+1)
	return &bitVector{n: n, size: size, bit: bit, preSum: preSum}
}

func (bv *bitVector) Set(i int32) {
	bv.bit[i>>6] |= 1 << (i & 63)
}

func (bv *bitVector) Build() {
	for i := int32(0); i < bv.size; i++ {
		bv.preSum[i+1] = bv.preSum[i] + int32(bits.OnesCount64(bv.bit[i]))
	}
}

func (bv *bitVector) Get(i int32) bool {
	return bv.bit[i>>6]>>(i&63)&1 == 1
}

func (bv *bitVector) Count0(end int32) int32 {
	return end - (bv.preSum[end>>6] + int32(bits.OnesCount64(bv.bit[end>>6]&(1<<(end&63)-1))))
}

func (bv *bitVector) Count1(end int32) int32 {
	return bv.preSum[end>>6] + int32(bits.OnesCount64(bv.bit[end>>6]&(1<<(end&63)-1)))
}

func (bv *bitVector) Count(end int32, value int8) int32 {
	if value == 1 {
		return bv.Count1(end)
	}
	return end - bv.Count1(end)
}

func (bv *bitVector) Kth0(k int32) int32 {
	if k < 0 || bv.Count0(bv.n) <= k {
		return -1
	}
	l, r := int32(0), bv.size+1
	for r-l > 1 {
		m := (l + r) >> 1
		if m<<6-bv.preSum[m] > k {
			r = m
		} else {
			l = m
		}
	}
	indx := l << 6
	k -= (l<<6 - bv.preSum[l]) - bv.Count0(indx)
	l, r = indx, indx+64
	for r-l > 1 {
		m := (l + r) >> 1
		if bv.Count0(m) > k {
			r = m
		} else {
			l = m
		}
	}
	return l
}

// k>=0.
func (bv *bitVector) Kth1(k int32) int32 {
	if k < 0 || bv.Count1(bv.n) <= k {
		return -1
	}
	l, r := int32(0), bv.size+1
	for r-l > 1 {
		m := (l + r) >> 1
		if bv.preSum[m] > k {
			r = m
		} else {
			l = m
		}
	}
	indx := l << 6
	k -= bv.preSum[l] - bv.Count1(indx)
	l, r = indx, indx+64
	for r-l > 1 {
		m := (l + r) >> 1
		if bv.Count1(m) > k {
			r = m
		} else {
			l = m
		}
	}
	return l
}

func (bv *bitVector) Kth(k int32, v int8) int32 {
	if v == 1 {
		return bv.Kth1(k)
	}
	return bv.Kth0(k)
}

func (bv *bitVector) GetAll() []bool {
	res := make([]bool, 0, bv.n)
	for i := int32(0); i < bv.n; i++ {
		res = append(res, bv.Get(i))
	}
	return res
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

func test() {
	for i := 0; i < 100; i++ {
		nums := make([]int, 1000)
		for j := 0; j < 1000; j++ {
			nums[j] = rand.Intn(1000)
		}
		wm := NewWaveletMatrixStatic(int32(len(nums)), func(i int32) int { return nums[i] }, maxs(nums))

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

		for j := 0; j < 100; j++ {
			start, end := rand.Intn(1000), rand.Intn(1000)
			if start > end {
				start, end = end, start
			}
			x := rand.Intn(1000)
			if prefixBf(int32(end), x) != wm.PrefixCount(int32(end), x) {
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
		}
	}

	fmt.Println("pass")
}

func testTime() {
	n := int32(2e5)
	nums := make([]int, n)
	for i := int32(0); i < n; i++ {
		nums[i] = rand.Intn(1e9 + 10)
	}

	time1 := time.Now()
	wm := NewWaveletMatrixStatic(int32(len(nums)), func(i int32) int { return nums[i] }, maxs(nums))
	fmt.Println("build time:", time.Since(time1))

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
	}

	fmt.Println(time.Since(time1)) // 1.0175364s

	time1 = time.Now()
	wm.TopK(0, n, n/2)
	fmt.Println(time.Since(time1)) // 47.2957ms
}
