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
	yuki738()
}

// 1703. 得到连续 K 个 1 的最少相邻交换次数
// https://leetcode.cn/problems/minimum-adjacent-swaps-for-k-consecutive-ones/description/
// 给你一个整数数组 nums 和一个整数 k 。 nums 仅包含 0 和 1 。每一次移动，你可以选择 相邻 两个数字并将它们交换。
// 请你返回使 nums 中包含 k 个 连续 1 的 最少 交换次数。
//
// !将每个1转换为i-len(onesIndex).
// 然后变成大小为k的滑动窗口所有数到中位数的距离和.
func minMoves(nums []int, k int) int {
	onesIndex := make([]int, 0, len(nums))
	for i, v := range nums {
		if v == 1 {
			onesIndex = append(onesIndex, i-len(onesIndex))
		}
	}
	onesIndex = onesIndex[:len(onesIndex):len(onesIndex)]
	dist := SlidingWindowDistSumToMedian(onesIndex, k)
	res := INF
	for _, d := range dist {
		res = min(res, d)
	}
	return res
}

// No.738 平らな農地
// https://yukicoder.me/problems/no/738
// !大小为k的滑动窗口所有数到中位数的距离和
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
	dist := SlidingWindowDistSumToMedian(nums, int(k))
	res := INF
	for _, d := range dist {
		res = min(res, d)
	}
	fmt.Fprintln(out, res)
}

// 对每个大小为k的滑动窗口，求出所有数到中位数的距离和.
func SlidingWindowDistSumToMedian(nums []int, k int) []int {
	if len(nums) < k {
		return nil
	}
	n := len(nums)
	wm := NewWaveletMatrixWithSum(nums, nums, -1, true)
	mf := NewMedianFinderWaveletMatrix(wm)
	res := make([]int, 0, n-k+1)
	for i := 0; i < n-k+1; i++ {
		start, end := int32(i), int32(i+k)
		res = append(res, mf.DistSumToMedianRange(start, end))
	}
	return res
}

// waveletMatrix维护区间中位数信息.
type MedianFinderWaveletMatrix struct {
	Wm *WaveletMatrixWithSum
}

// WaveletMatrix 维护区间中位数信息.
// `Proxy 的内部持有一个对 WaveletMatrix 的引用`.
func NewMedianFinderWaveletMatrix(wm *WaveletMatrixWithSum) *MedianFinderWaveletMatrix {
	return &MedianFinderWaveletMatrix{Wm: wm}
}

// upper: 如果有两个中位数，返回较大的那个.
func (mf *MedianFinderWaveletMatrix) Median(upper bool) int {
	return mf.MedianRange(0, mf.Wm.n, upper)
}

// 返回区间 [start, end) 中的中位数.
// upper: 如果有两个中位数，返回较大的那个.
func (mf *MedianFinderWaveletMatrix) MedianRange(start, end int32, upper bool) int {
	if start < 0 {
		start = 0
	}
	if end > mf.Wm.n {
		end = mf.Wm.n
	}
	if start >= end {
		return 0
	}
	return mf.Wm.Median(start, end, upper, 0)
}

func (mf *MedianFinderWaveletMatrix) DistSum(to int) int {
	return mf.DistSumRange(to, 0, mf.Wm.n)
}

func (mf *MedianFinderWaveletMatrix) DistSumRange(to int, start, end int32) int {
	if start < 0 {
		start = 0
	}
	if end > mf.Wm.n {
		end = mf.Wm.n
	}
	if start >= end {
		return 0
	}
	m := end - start
	lowerCount, lowerSum := mf.Wm.RangeCountAndSum(start, end, -INF, WmValue(to), 0) // bisect_left
	allSum := mf.Wm.SumAll(start, end)
	if lowerCount == 0 {
		return allSum - int(m)*to
	}
	if lowerCount == m {
		return int(m)*to - allSum
	}
	upperSum := allSum - lowerSum
	leftSum := to*int(lowerCount) - lowerSum
	rightSum := upperSum - to*int(m-lowerCount)
	return leftSum + rightSum
}

func (mf *MedianFinderWaveletMatrix) DistSumToMedian() int {
	return mf.DistSumToMedianRange(0, mf.Wm.n)
}

func (mf *MedianFinderWaveletMatrix) DistSumToMedianRange(start, end int32) int {
	if start < 0 {
		start = 0
	}
	if end > mf.Wm.n {
		end = mf.Wm.n
	}
	if start >= end {
		return 0
	}
	m := end - start
	count1 := m / 2
	count2 := m - count1
	mid, sum1 := mf.Wm.KthValueAndSum(start, end, count1, 0)
	allSum := mf.Wm.SumAll(start, end)
	sum2 := allSum - sum1
	res := 0
	res += mid*int(count1) - sum1
	res += sum2 - mid*int(count2)
	return res
}

const INF WmValue = 1e18

type WmValue = int
type WmSum = int

func (*WaveletMatrixWithSum) e() WmSum            { return 0 }
func (*WaveletMatrixWithSum) op(a, b WmSum) WmSum { return a + b }
func (*WaveletMatrixWithSum) inv(a WmSum) WmSum   { return -a }

type WaveletMatrixWithSum struct {
	n, log   int32
	setLog   bool
	compress bool
	useSum   bool
	mid      []int32
	bv       []*BitVector
	key      []WmValue
	presum   [][]WmSum
}

// nums: 数组元素.
// sumData: 和数据，nil表示不需要和数据.
// log: 如果需要支持异或查询则需要传入log，-1表示默认.
// compress: 是否对nums进行离散化(值域较大(1e9)时可以离散化加速).
func NewWaveletMatrixWithSum(nums []WmValue, sumData []WmSum, log int32, compress bool) *WaveletMatrixWithSum {
	wm := &WaveletMatrixWithSum{}
	wm.build(nums, sumData, log, compress)
	return wm
}

func (wm *WaveletMatrixWithSum) build(nums []WmValue, sumData []WmSum, log int32, compress bool) {
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

// 返回区间 [start, end) 中 值在 [a, b) 中的元素个数以及这些元素的和.
func (wm *WaveletMatrixWithSum) RangeCountAndSum(start, end int32, a, b WmValue, xorValue WmValue) (int32, WmSum) {
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

// 返回区间 [start, end) 中的 (第k小的元素, 前k个元素(不包括第k小的元素) 的 op 的结果).
// 如果k >= end-start, 返回 (-1, 区间 op 的结果).
func (wm *WaveletMatrixWithSum) KthValueAndSum(start, end, k int32, xorVal WmValue) (WmValue, WmSum) {
	if start < 0 {
		start = 0
	}
	if end > wm.n {
		end = wm.n
	}
	if start >= end {
		return -1, wm.e()
	}
	if k >= end-start {
		return -1, wm.SumAll(start, end)
	}
	if xorVal != 0 {
		if !wm.setLog {
			panic("log should be set when xor is used")
		}
	}

	sum, val := wm.e(), WmValue(0)
	for d := wm.log - 1; d >= 0; d-- {
		l0, r0 := wm.bv[d].Rank(start, false), wm.bv[d].Rank(end, false)
		l1, r1 := start+wm.mid[d]-l0, end+wm.mid[d]-r0
		if (xorVal>>d)&1 == 1 {
			l0, l1 = l1, l0
			r0, r1 = r1, r0
		}
		if k < r0-l0 {
			start, end = l0, r0
		} else {
			k -= r0 - l0
			val |= 1 << d
			start, end = l1, r1
			if wm.useSum {
				sum = wm.op(sum, wm._get(d, l0, r0))
			}
		}
	}
	if wm.useSum {
		sum = wm.op(sum, wm._get(0, start, start+k))
	}
	if wm.compress {
		val = wm.key[val]
	}
	return val, sum
}

// [start, end)区间内第k(k>=0)小的元素.
func (wm *WaveletMatrixWithSum) Kth(start, end, k int32, xorVal WmValue) WmValue {
	if k < 0 {
		k = 0
	}
	if n := end - start - 1; k > n {
		k = n
	}
	v, _ := wm.KthValueAndSum(start, end, k, xorVal)
	return v
}

// upper: 向上取中位数还是向下取中位数.
func (wm *WaveletMatrixWithSum) Median(start, end int32, upper bool, xorVal WmValue) WmValue {
	n := end - start
	var k int32
	if upper {
		k = n >> 1
	} else {
		k = (n - 1) >> 1
	}
	return wm.Kth(start, end, k, xorVal)
}

// [start, end) 中小于等于 x 的数中最大的数.
//
//	如果不存在则返回-INF.
func (wm *WaveletMatrixWithSum) Floor(start, end int32, x WmValue, xor WmValue) WmValue {
	if xor != 0 {
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
	if start >= end {
		return -INF
	}
	res := -INF
	x++
	if wm.compress {
		x = wm._lowerBound(wm.key, x)
	}
	var dfs func(d, l, r int32, lx, rx WmValue)
	dfs = func(d, l, r int32, lx, rx WmValue) {
		if rx-1 <= res || l == r || x <= lx {
			return
		}
		if d == 0 {
			res = max(res, lx)
			return
		}
		d--
		mx := (lx + rx) >> 1
		l0, r0 := wm.bv[d].Rank(l, false), wm.bv[d].Rank(r, false)
		l1, r1 := l+wm.mid[d]-l0, r+wm.mid[d]-r0
		if xor>>d&1 == 1 {
			l0, l1 = l1, l0
			r0, r1 = r1, r0
		}
		dfs(d, l1, r1, mx, rx)
		dfs(d, l0, r0, lx, mx)
	}
	dfs(wm.log, start, end, 0, 1<<wm.log)
	if wm.compress && res != -INF {
		res = wm.key[res]
	}
	return res
}

// [start, end) 中大于等于 x 的数中最小的数
//
//	如果不存在则返回INF
func (wm *WaveletMatrixWithSum) Ceil(start, end int32, x WmValue, xor WmValue) int {
	if xor != 0 {
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
	if start >= end {
		return INF
	}
	if wm.compress {
		x = wm._lowerBound(wm.key, x)
	}
	res := INF
	var dfs func(d, l, r int32, lx, rx WmValue)
	dfs = func(d, l, r int32, lx, rx WmValue) {
		if res <= lx || l == r || rx <= x {
			return
		}
		if d == 0 {
			res = min(res, lx)
			return
		}
		d--
		mx := (lx + rx) >> 1
		l0, r0 := wm.bv[d].Rank(l, false), wm.bv[d].Rank(r, false)
		l1, r1 := l+wm.mid[d]-l0, r+wm.mid[d]-r0
		if xor>>d&1 == 1 {
			l0, l1 = l1, l0
			r0, r1 = r1, r0
		}
		dfs(d, l0, r0, lx, mx)
		dfs(d, l1, r1, mx, rx)
	}
	dfs(wm.log, start, end, 0, 1<<wm.log)
	if wm.compress && res < INF {
		res = wm.key[res]
	}
	return res
}

// 返回区间 [start, end) 中 范围在 [a, b) 中的元素的和.
func (wm *WaveletMatrixWithSum) SumRange(start, end int32, a, b WmValue, xorVal WmValue) WmSum {
	if !wm.useSum {
		panic("sum data must be provided")
	}
	if start < 0 {
		start = 0
	}
	if end > wm.n {
		end = wm.n
	}
	if start >= end || a >= b {
		return wm.e()
	}
	_, sum := wm.RangeCountAndSum(start, end, a, b, xorVal)
	return sum
}

// 返回区间 [start, end) 中 排名在 [k1, k2) 中的元素的和.
func (wm *WaveletMatrixWithSum) SumSlice(start, end, k1, k2 int32, xorVal WmValue) WmSum {
	if !wm.useSum {
		panic("sum data must be provided")
	}
	if k1 < 0 {
		k1 = 0
	}
	if k2 > end-start {
		k2 = end - start
	}
	if start < 0 {
		start = 0
	}
	if end > wm.n {
		end = wm.n
	}
	if start >= end || k1 >= k2 {
		return wm.e()
	}
	_, sum1 := wm.KthValueAndSum(start, end, k1, xorVal)
	_, sum2 := wm.KthValueAndSum(start, end, k2, xorVal)
	return wm.op(sum2, wm.inv(sum1))
}

func (wm *WaveletMatrixWithSum) SumAll(start, end int32) WmSum {
	if start < 0 {
		start = 0
	}
	if end > wm.n {
		end = wm.n
	}
	if start >= end {
		return wm.e()
	}
	return wm._get(wm.log, start, end)
}

// 使得predicate(count, sum)为true的最大的(count, sum).
func (wm *WaveletMatrixWithSum) MaxRight(predicate func(int32, WmSum) bool, start, end int32, xorVal WmValue) (int32, WmSum) {
	if xorVal != 0 {
		if !wm.setLog {
			panic("log should be set when xor is used")
		}
	}
	if start >= end {
		return 0, wm.e()
	}
	if s := wm._get(wm.log, start, end); predicate(end-start, s) {
		return end - start, s
	}
	count, sum := int32(0), wm.e()
	for d := wm.log - 1; d >= 0; d-- {
		l0, r0 := wm.bv[d].Rank(start, false), wm.bv[d].Rank(end, false)
		l1, r1 := start+wm.mid[d]-l0, end+wm.mid[d]-r0
		if xorVal>>d&1 == 1 {
			l0, l1 = l1, l0
			r0, r1 = r1, r0
		}
		if s := wm.op(sum, wm._get(d, l0, r0)); predicate(count+r0-l0, s) {
			count += r0 - l0
			sum = s
			start, end = l1, r1
		} else {
			start, end = l0, r0
		}
	}
	k := wm._binarySearch(func(k int32) bool {
		return predicate(count+k, wm.op(sum, wm._get(0, start, start+k)))
	}, 0, end-start)
	count += k
	sum = wm.op(sum, wm._get(0, start, start+k))
	return count, sum
}

func (wm *WaveletMatrixWithSum) _get(d, l, r int32) WmSum {
	if wm.useSum {
		return wm.op(wm.presum[d][r], wm.inv(wm.presum[d][l]))
	}
	return wm.e()
}

func (wm *WaveletMatrixWithSum) _argSort(nums []WmValue) []int32 {
	order := make([]int32, len(nums))
	for i := range order {
		order[i] = int32(i)
	}
	sort.Slice(order, func(i, j int) bool { return nums[order[i]] < nums[order[j]] })
	return order
}

func (wm *WaveletMatrixWithSum) _maxs(nums []WmValue) WmValue {
	res := nums[0]
	for _, v := range nums {
		if v > res {
			res = v
		}
	}
	return res
}

func (wm *WaveletMatrixWithSum) _lowerBound(nums []WmValue, target WmValue) WmValue {
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

func (wm *WaveletMatrixWithSum) _binarySearch(f func(int32) bool, ok, ng int32) int32 {
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
func test() {

	medianBruteForce := func(nums []int) int {
		if len(nums) == 0 {
			return 0
		}
		sortedNums := append([]int(nil), nums...)
		sort.Ints(sortedNums)
		mid := (len(sortedNums) - 1) / 2
		return sortedNums[mid]
	}

	distSumBruteForce := func(nums []int, to int) int {
		res := 0
		for _, num := range nums {
			res += abs(num - to)
		}
		return res
	}

	distSumRangeBruteForce := func(nums []int, to, start, end int) int {
		res := 0
		for i := start; i < end; i++ {
			res += abs(nums[i] - to)
		}
		return res
	}

	distSumToMedianBruteForce := func(nums []int) int {
		median := medianBruteForce(nums)
		return distSumBruteForce(nums, median)
	}

	distSumToMedianRangeBruteForce := func(nums []int, start, end int) int {
		median := medianBruteForce(nums[start:end])
		return distSumRangeBruteForce(nums, median, start, end)
	}

	for tc := 0; tc < 1000; tc++ {
		n := rand.Intn(1000) + 1
		nums := make([]int, n)
		for i := 0; i < n; i++ {
			nums[i] = rand.Intn(1e6) - 5e5
		}
		mf := NewMedianFinderWaveletMatrix(NewWaveletMatrixWithSum(nums, nums, -1, true))

		for i := 0; i < 100; i++ {
			start, end := rand.Intn(n), rand.Intn(n)+1
			if start > end {
				start, end = end, start
			}
			to := rand.Intn(1e5) - 5e4

			if mf.MedianRange(int32(start), int32(end), false) != medianBruteForce(nums[start:end]) {
				panic("err0")
			}
			if mf.DistSum(int(to)) != distSumBruteForce(nums, to) {
				panic("err1")
			}
			if mf.DistSumRange(int(to), int32(start), int32(end)) != distSumRangeBruteForce(nums, to, start, end) {
				panic("err2")
			}
			if mf.DistSumToMedian() != distSumToMedianBruteForce(nums) {
				panic("err3")
			}
			if mf.DistSumToMedianRange(int32(start), int32(end)) != distSumToMedianRangeBruteForce(nums, start, end) {
				fmt.Println(mf.DistSumToMedianRange(int32(start), int32(end)), distSumToMedianRangeBruteForce(nums, start, end))
				panic("err4")
			}
		}
	}

	fmt.Println("pass")
}
