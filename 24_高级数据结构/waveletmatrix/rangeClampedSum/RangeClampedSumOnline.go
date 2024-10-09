// RangeClampSum/RangeClampedSum
// SumWithMin/SumWithMax/SumWithMinAndMax

package main

import (
	"math/bits"
	"sort"
)

// 100404. 统计满足 K 约束的子字符串数量 II
// https://leetcode.cn/problems/count-substrings-that-satisfy-k-constraint-ii/description/
//
// 思路：对每个结尾，找到最小的起点，使得区间内的 0 和 1 的数量都不超过 k
// !区间查询[l,r]时，相当于对结尾为l+1,l+2,...,r进行求和（RangeClampedSum/RangeClampSum）
func countKConstraintSubstrings(s string, k int, queries [][]int) []int64 {
	n := len(s)
	leftBound := func() []int {
		left, counter := 0, [2]int{}
		res := make([]int, n)
		for right := 0; right < n; right++ {
			counter[s[right]-'0']++
			for counter[0] > k && counter[1] > k {
				counter[s[left]-'0']--
				left++
			}
			res[right] = left
		}
		return res
	}()

	res := make([]int64, len(queries))
	S := NewRangeClampedSumOnline(leftBound)
	for i, q := range queries {
		l, r := q[0], q[1]
		rightSum := (l + r + 2) * (r - l + 1) / 2
		leftSum := S.SumWithMin(int32(l), int32(r)+1, l)
		res[i] = int64(rightSum - leftSum)
	}
	return res
}

type RangeClampedSumOnline struct {
	n  int32
	wm *waveletMatrixWithSum
}

func NewRangeClampedSumOnline(nums []int) *RangeClampedSumOnline {
	return &RangeClampedSumOnline{n: int32(len(nums)), wm: NewWaveletMatrixWithSum(nums, nums, -1, true)}
}

// [min, ?)
func (rcs *RangeClampedSumOnline) SumWithMin(start, end int32, min int) int {
	if start < 0 {
		start = 0
	}
	if end > rcs.n {
		end = rcs.n
	}
	if start >= end {
		return 0
	}
	count, sum := rcs.wm.RangeCountAndSum(start, end, min, INF, 0)
	return sum + min*int((end-start)-count)
}

// (?, max]
func (rcs *RangeClampedSumOnline) SumWithMax(start, end int32, max int) int {
	if start < 0 {
		start = 0
	}
	if end > rcs.n {
		end = rcs.n
	}
	if start >= end {
		return 0
	}
	count, sum := rcs.wm.RangeCountAndSum(start, end, -INF, max+1, 0)
	return sum + max*int((end-start)-count)
}

// [min, max]
func (rcs *RangeClampedSumOnline) SumWithMinAndMax(start, end int32, min, max int) int {
	if start < 0 {
		start = 0
	}
	if end > rcs.n {
		end = rcs.n
	}
	if start >= end {
		return 0
	}
	if min >= max {
		return 0
	}
	count1 := rcs.wm.PrefixCount(start, end, min)
	count2, sum2 := rcs.wm.RangeCountAndSum(start, end, min, max+1, 0)
	return sum2 + min*int(count1) + max*int((end-start)-count2-count1)
}

const INF WmValue = 1e18

type WmValue = int
type WmSum = int

func (*waveletMatrixWithSum) e() WmSum            { return 0 }
func (*waveletMatrixWithSum) op(a, b WmSum) WmSum { return a + b }
func (*waveletMatrixWithSum) inv(a WmSum) WmSum   { return -a }

type waveletMatrixWithSum struct {
	setLog   bool
	compress bool
	useSum   bool
	n, log   int32
	mid      []int32
	bv       []*BitVector
	key      []WmValue
	presum   [][]WmSum
}

// nums: 数组元素.
// sumData: 和数据，nil表示不需要和数据.
// log: 如果需要支持异或查询则需要传入log，-1表示默认.
// compress: 是否对nums进行离散化(值域较大(1e9)时可以离散化加速).
func NewWaveletMatrixWithSum(nums []WmValue, sumData []WmSum, log int32, compress bool) *waveletMatrixWithSum {
	wm := &waveletMatrixWithSum{}
	wm.build(nums, sumData, log, compress)
	return wm
}

func (wm *waveletMatrixWithSum) build(nums []WmValue, sumData []WmSum, log int32, compress bool) {
	numsCopy := append(nums[:0:0], nums...)
	sumDataCopy := append(sumData[:0:0], sumData...)

	wm.n = int32(len(numsCopy))
	wm.log = log
	wm.setLog = log != -1
	wm.compress = compress
	wm.useSum = len(sumData) > 0
	if wm.n == 0 {
		wm.log = 0
		wm.presum = [][]WmSum{{wm.e()}}
		return
	}

	if compress {
		if wm.setLog {
			panic("compress and log should not be set at the same time")
		}
		wm.key = make([]WmValue, 0, wm.n)
		order := wm._argSort(numsCopy)
		for _, i := range order {
			if len(wm.key) == 0 || wm.key[len(wm.key)-1] != numsCopy[i] {
				wm.key = append(wm.key, numsCopy[i])
			}
			numsCopy[i] = WmValue(len(wm.key) - 1)
		}
		wm.key = wm.key[:len(wm.key):len(wm.key)]
	}
	if wm.log == -1 {
		tmp := wm._maxs(numsCopy)
		if tmp < 1 {
			tmp = 1
		}
		wm.log = int32(bits.Len(uint(tmp)))
	}
	wm.mid = make([]int32, wm.log)
	wm.bv = make([]*BitVector, wm.log)
	for i := range wm.bv {
		wm.bv[i] = NewBitVector(wm.n)
	}
	if wm.useSum {
		wm.presum = make([][]WmSum, 1+wm.log)
		for i := range wm.presum {
			sums := make([]WmSum, wm.n+1)
			for j := range sums {
				sums[j] = wm.e()
			}
			wm.presum[i] = sums
		}
	}
	if len(sumDataCopy) == 0 {
		sumDataCopy = make([]WmSum, len(numsCopy))
	}

	A, S := numsCopy, sumDataCopy
	A0, A1 := make([]WmValue, wm.n), make([]WmValue, wm.n)
	S0, S1 := make([]WmSum, wm.n), make([]WmSum, wm.n)
	for d := wm.log - 1; d >= -1; d-- {
		p0, p1 := int32(0), int32(0)
		if wm.useSum {
			tmp := wm.presum[d+1]
			for i := int32(0); i < wm.n; i++ {
				tmp[i+1] = wm.op(tmp[i], S[i])
			}
		}
		if d == -1 {
			break
		}
		for i := int32(0); i < wm.n; i++ {
			f := (A[i] >> d & 1) == 1
			if !f {
				if wm.useSum {
					S0[p0] = S[i]
				}
				A0[p0] = A[i]
				p0++
			} else {
				if wm.useSum {
					S1[p1] = S[i]
				}
				wm.bv[d].Set(i)
				A1[p1] = A[i]
				p1++
			}
		}
		wm.mid[d] = p0
		wm.bv[d].Build()
		A, A0 = A0, A
		S, S0 = S0, S
		for i := int32(0); i < p1; i++ {
			A[p0+i] = A1[i]
			S[p0+i] = S1[i]
		}
	}
}

func (wm *waveletMatrixWithSum) PrefixCount(start, end int32, x WmValue) int32 {
	if start < 0 {
		start = 0
	}
	if end > wm.n {
		end = wm.n
	}
	if start >= end {
		return 0
	}
	if wm.compress {
		x = wm._lowerBound(wm.key, x)
	}
	if x == 0 {
		return 0
	}
	count := int32(0)
	for d := wm.log - 1; d >= 0; d-- {
		l0, r0 := wm.bv[d].Rank(start, false), wm.bv[d].Rank(end, false)
		l1, r1 := start+wm.mid[d]-l0, end+wm.mid[d]-r0
		if x>>d&1 == 1 {
			count += r0 - l0
			start, end = l1, r1
		} else {
			start, end = l0, r0
		}
	}
	return count
}

// 返回区间 [start, end) 中 值在 [a, b) 中的元素个数以及这些元素的和.
func (wm *waveletMatrixWithSum) RangeCountAndSum(start, end int32, a, b WmValue, xorValue WmValue) (int32, WmSum) {
	if xorValue != 0 {
		if !wm.setLog {
			panic("log should be set when xor is used")
		}
	}
	if start < 0 {
		start = 0
	}
	if end > wm.n {
		end = wm.n
	}
	if start >= end || a >= b {
		return 0, wm.e()
	}
	if wm.compress {
		a = wm._lowerBound(wm.key, a)
		b = wm._lowerBound(wm.key, b)
	}
	count, sum := int32(0), wm.e()
	var dfs func(d, l, r int32, lx, rx WmValue)
	dfs = func(d, l, r int32, lx, rx WmValue) {
		if rx <= a || b <= lx {
			return
		}
		if a <= lx && rx <= b {
			count += r - l
			if wm.useSum {
				sum = wm.op(sum, wm._get(d, l, r))
			}
			return
		}
		d--
		mx := (lx + rx) >> 1
		l0, r0 := wm.bv[d].Rank(l, false), wm.bv[d].Rank(r, false)
		l1, r1 := l+wm.mid[d]-l0, r+wm.mid[d]-r0
		if xorValue>>d&1 == 1 {
			l0, l1 = l1, l0
			r0, r1 = r1, r0
		}
		dfs(d, l0, r0, lx, mx)
		dfs(d, l1, r1, mx, rx)
	}
	dfs(wm.log, start, end, 0, 1<<wm.log)
	return count, sum
}

func (wm *waveletMatrixWithSum) _get(d, l, r int32) WmSum {
	if wm.useSum {
		return wm.op(wm.presum[d][r], wm.inv(wm.presum[d][l]))
	}
	return wm.e()
}

func (wm *waveletMatrixWithSum) _argSort(nums []WmValue) []int32 {
	order := make([]int32, len(nums))
	for i := range order {
		order[i] = int32(i)
	}
	sort.Slice(order, func(i, j int) bool { return nums[order[i]] < nums[order[j]] })
	return order
}

func (wm *waveletMatrixWithSum) _maxs(nums []WmValue) WmValue {
	res := nums[0]
	for _, v := range nums {
		if v > res {
			res = v
		}
	}
	return res
}

func (wm *waveletMatrixWithSum) _lowerBound(nums []WmValue, target WmValue) WmValue {
	left, right := int32(0), int32(len(nums)-1)
	for left <= right {
		mid := (left + right) >> 1
		if nums[mid] < target {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return WmValue(left)
}

func (wm *waveletMatrixWithSum) _binarySearch(f func(int32) bool, ok, ng int32) int32 {
	for abs32(ok-ng) > 1 {
		x := (ok + ng) >> 1
		if f(x) {
			ok = x
		} else {
			ng = x
		}
	}
	return ok
}

type BitVector struct {
	bits   []uint64
	preSum []int32
}

func NewBitVector(n int32) *BitVector {
	return &BitVector{bits: make([]uint64, n>>6+1), preSum: make([]int32, n>>6+1)}
}

func (bv *BitVector) Set(i int32) {
	bv.bits[i>>6] |= 1 << (i & 63)
}

func (bv *BitVector) Build() {
	for i := 0; i < len(bv.bits)-1; i++ {
		bv.preSum[i+1] = bv.preSum[i] + int32(bits.OnesCount64(bv.bits[i]))
	}
}

func (bv *BitVector) Rank(k int32, f bool) int32 {
	m, s := bv.bits[k>>6], bv.preSum[k>>6]
	res := s + int32(bits.OnesCount64(m&((1<<(k&63))-1)))
	if f {
		return res
	}
	return k - res
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func abs32(x int32) int32 {
	if x < 0 {
		return -x
	}
	return x
}

func min32(a, b int32) int32 {
	if a < b {
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
