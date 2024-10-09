// SumWithMin/SumWithMax/SumWithMinAndMax

package main

import (
	"fmt"
	"math/bits"
	"sort"
)

func main() {
	wm := NewRangeClampedSumOnline([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, true)
	fmt.Println(wm.SumWithMin(0, 10, 5))          // 65
	fmt.Println(wm.SumWithMin(0, 10, 6))          // 70
	fmt.Println(wm.SumWithMax(0, 10, 5))          // 40
	fmt.Println(wm.SumWithMinAndMax(0, 10, 5, 6)) // 55
}

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
	S := NewRangeClampedSumOnline(leftBound, true)
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
	wm *WaveletMatrixWithSumFast
}

func NewRangeClampedSumOnline(nums []int, smallRange bool) *RangeClampedSumOnline {
	wm := NewWaveletMatrixWithSumFast(smallRange)
	wm.Build(int32(len(nums)), func(i int32) (v int, e E) { return nums[i], nums[i] })
	return &RangeClampedSumOnline{n: int32(len(nums)), wm: wm}
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
	count, sum := rcs.wm.CountAndSum(start, end, min, INF)
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
	count, sum := rcs.wm.CountAndSum(start, end, -INF, max+1)
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
	count1 := rcs.wm.CountPrefix(start, end, min)
	count2, sum2 := rcs.wm.CountAndSum(start, end, min, max+1)
	return sum2 + min*int(count1) + max*int((end-start)-count2-count1)
}

const INF int = 2e18

type E = int

func e() E        { return 0 }
func op(a, b E) E { return a + b }
func inv(a E) E   { return -a } // 可选, 查询范围和时用于差分

// (类)线段树的实现.
type SegTreeLike struct {
	seg *StaticRangeProductGroup
	// seg *SegmentTree // StaticRangeProduct
}

func (st *SegTreeLike) Build(sum []E) {
	st.seg = NewStaticRangeProductGroup(sum)
	// st.seg = NewSegmentTreeFrom(sum)
}

func (st *SegTreeLike) Query(start, end int32) E {
	return st.seg.Query(start, end)
}

func (st *SegTreeLike) Set(i int32, e E) {
	// st.seg.Set(i, e)
}

func (st *SegTreeLike) Update(i int32, e E) {
	// st.seg.Update(i, e)
}

// ------------------- SegTreeOnWaveletMatrix -------------------

var _ ISegTreeLike = (*SegTreeLike)(nil)

type ISegTreeLike interface {
	Build(sum []E)
	Query(start, end int32) E
	Set(i int32, e E)
	Update(i int32, e E)
}

// 内部不使用接口约束，因为接口会导致性能下降.
type WaveletMatrixWithSumFast struct {
	n, log, upper int32

	raw []int
	mid []int32
	bv  []*bitVector
	seg []*SegTreeLike

	smallRange bool
	index      func(int) int32
}

func NewWaveletMatrixWithSumFast(smallRange bool) *WaveletMatrixWithSumFast {
	return &WaveletMatrixWithSumFast{smallRange: smallRange}
}

func (st *WaveletMatrixWithSumFast) Build(m int32, f func(i int32) (v int, e E)) {
	arr, sum := make([]int, m), make([]E, m)
	for i := int32(0); i < m; i++ {
		arr[i], sum[i] = f(i)
	}
	st.build(arr, sum)
}

// [start, end) x [0, y)
func (st *WaveletMatrixWithSumFast) CountPrefix(start, end int32, y int) int32 {
	p := st.index(y)
	if p == 0 || start >= end {
		return 0
	}
	if p == st.upper {
		return end - start
	}
	res := int32(0)
	for d := st.log - 1; d >= 0; d-- {
		l0, r0 := st.bv[d].Rank0(start), st.bv[d].Rank0(end)
		l1, r1 := start+st.mid[d]-l0, end+st.mid[d]-r0
		if p>>d&1 == 1 {
			res += r0 - l0
			start, end = l1, r1
		} else {
			start, end = l0, r0
		}
	}
	return res
}

// [start, end) x [y1, y2)
func (st *WaveletMatrixWithSumFast) Count(start, end int32, y1, y2 int) int32 {
	if y1 >= y2 {
		return 0
	}
	return st.CountPrefix(start, end, y2) - st.CountPrefix(start, end, y1)
}

// [start, end) x [0, y)
func (st *WaveletMatrixWithSumFast) CountAndSumPrefix(start, end int32, y int) (int32, E) {
	p := st.index(y)
	if p == 0 {
		return 0, e()
	}
	if p == st.upper {
		return end - start, st.seg[st.log].Query(start, end)
	}
	count := int32(0)
	sum := e()
	for d := st.log - 1; d >= 0; d-- {
		l0, r0 := st.bv[d].Rank0(start), st.bv[d].Rank0(end)
		l1, r1 := start+st.mid[d]-l0, end+st.mid[d]-r0
		if p>>d&1 == 1 {
			count += r0 - l0
			sum = op(sum, st.seg[d].Query(l0, r0))
			start, end = l1, r1
		} else {
			start, end = l0, r0
		}
	}
	return count, sum
}

// [start, end) x [y1, y2)
func (st *WaveletMatrixWithSumFast) CountAndSum(start, end int32, y1, y2 int) (int32, E) {
	if start < 0 {
		start = 0
	}
	if end > st.n {
		end = st.n
	}
	if start >= end || y1 >= y2 {
		return 0, e()
	}
	lo, hi := st.index(y1), st.index(y2)
	count := int32(0)
	sum := e()
	var dfs func(int32, int32, int32, int32, int32)
	dfs = func(d, L, R, a, b int32) {
		if hi <= a || b <= lo {
			return
		}
		if lo <= a && b <= hi {
			count += R - L
			sum = op(sum, st.seg[d].Query(L, R))
			return
		}
		d--
		c := (a + b) >> 1
		l0, r0 := st.bv[d].Rank0(L), st.bv[d].Rank0(R)
		l1, r1 := L+st.mid[d]-l0, R+st.mid[d]-r0
		dfs(d, l0, r0, a, c)
		dfs(d, l1, r1, c, b)
	}
	dfs(st.log, start, end, 0, 1<<st.log)
	return count, sum
}

// [start, end) x [0, y)
func (st *WaveletMatrixWithSumFast) SumPrefix(start, end int32, y int) E {
	_, sum := st.CountAndSumPrefix(start, end, y)
	return sum
}

// [start, end) x [y1, y2)
func (st *WaveletMatrixWithSumFast) Sum(start, end int32, y1, y2 int) E {
	_, sum := st.CountAndSum(start, end, y1, y2)
	return sum
}

func (st *WaveletMatrixWithSumFast) SumAll(start, end int32) E {
	return st.seg[st.log].Query(start, end)
}

// 排名在[k1, k2)间的元素的和.要求运算存在逆元.
func (st *WaveletMatrixWithSumFast) SumIndexRange(start, end int32, k1, k2 int32) E {
	if k1 < 0 {
		k1 = 0
	}
	if k2 > end-start {
		k2 = end - start
	}
	if start < 0 {
		start = 0
	}
	if end > st.n {
		end = st.n
	}
	if start >= end || k1 >= k2 {
		return e()
	}
	_, sum1 := st.KthValueAndSum(start, end, k1)
	_, sum2 := st.KthValueAndSum(start, end, k2)
	return op(inv(sum1), sum2)
}

// [start, end)区间内第k(k>=0)小的元素.
func (st *WaveletMatrixWithSumFast) Kth(start, end int32, k int32) int {
	if k < 0 {
		k = 0
	}
	if n := end - start - 1; k > n {
		k = n
	}
	p := int32(0)
	for d := st.log - 1; d >= 0; d-- {
		l0, r0 := st.bv[d].Rank0(start), st.bv[d].Rank0(end)
		l1, r1 := start+st.mid[d]-l0, end+st.mid[d]-r0
		if k < r0-l0 {
			start, end = l0, r0
		} else {
			k -= r0 - l0
			start, end = l1, r1
			p |= 1 << d
		}
	}
	return st.raw[p]
}

// 返回区间 [start, end) 中的 (第k小的元素, 前k个元素(不包括第k小的元素) 的 op 的结果).
// 如果k >= end-start, 返回 (INF, 区间 op 的结果).
func (st *WaveletMatrixWithSumFast) KthValueAndSum(start, end int32, k int32) (int, E) {
	if start < 0 {
		start = 0
	}
	if end > st.n {
		end = st.n
	}
	if start >= end {
		return INF, e()
	}
	if k >= end-start {
		return INF, st.SumAll(start, end)
	}
	p := int32(0)
	sum := e()
	for d := st.log - 1; d >= 0; d-- {
		l0, r0 := st.bv[d].Rank0(start), st.bv[d].Rank0(end)
		l1, r1 := start+st.mid[d]-l0, end+st.mid[d]-r0
		if k < r0-l0 {
			start, end = l0, r0
		} else {
			sum = op(sum, st.seg[d].Query(l0, r0))
			k -= r0 - l0
			start, end = l1, r1
			p |= 1 << d
		}
	}
	sum = op(sum, st.seg[0].Query(start, start+k))
	return st.raw[p], sum
}

// <= y 的最大值. 不存在则返回 -INF.
func (st *WaveletMatrixWithSumFast) Prev(start, end int32, y int) int {
	if start < 0 {
		start = 0
	}
	if end > st.n {
		end = st.n
	}
	if start >= end {
		return -INF
	}
	k := st.index(y + 1)
	p := int32(-1)
	var dfs func(int32, int32, int32, int32, int32)
	dfs = func(d, L, R, a, b int32) {
		if b-1 <= p || L == R || k <= a {
			return
		}
		if d == 0 {
			if p < a {
				p = a
			}
			return
		}
		d--
		c := (a + b) >> 1
		l0, r0 := st.bv[d].Rank0(L), st.bv[d].Rank0(R)
		l1, r1 := L+st.mid[d]-l0, R+st.mid[d]-r0
		dfs(d, l1, r1, c, b)
		dfs(d, l0, r0, a, c)
	}
	dfs(st.log, start, end, 0, 1<<st.log)
	if p == -1 {
		return -INF
	}
	return st.raw[p]
}

// >= y 的最小值. 不存在则返回 INF.
func (st *WaveletMatrixWithSumFast) Next(start, end int32, y int) int {
	if start < 0 {
		start = 0
	}
	if end > st.n {
		end = st.n
	}
	if start >= end {
		return INF
	}
	k := st.index(y)
	p := st.upper
	var dfs func(int32, int32, int32, int32, int32)
	dfs = func(d, L, R, a, b int32) {
		if p <= a || L == R || b <= k {
			return
		}
		if d == 0 {
			if a < p {
				p = a
			}
			return
		}
		d--
		c := (a + b) >> 1
		l0, r0 := st.bv[d].Rank0(L), st.bv[d].Rank0(R)
		l1, r1 := L+st.mid[d]-l0, R+st.mid[d]-r0
		dfs(d, l0, r0, a, c)
		dfs(d, l1, r1, c, b)
	}
	dfs(st.log, start, end, 0, 1<<st.log)
	if p == st.upper {
		return INF
	}
	return st.raw[p]
}

// upper: 向上取中位数还是向下取中位数.
func (st *WaveletMatrixWithSumFast) Median(start, end int32, upper bool) int {
	if start < 0 || start >= end || end > st.n {
		panic("invalid range")
	}
	var k int32
	if upper {
		k = (end - start) >> 1
	} else {
		k = (end - start - 1) >> 1
	}
	return st.Kth(start, end, k)
}

// [start, end) x [0, y) 使得 check(count, sum) 为真的最大的 (count, sum)
func (st *WaveletMatrixWithSumFast) MaxRight(start, end int32, check func(int32, E) bool) (int32, E) {
	if start >= end {
		return 0, e()
	}
	if s := st.SumAll(start, end); check(end-start, s) {
		return end - start, s
	}
	count := int32(0)
	sum := e()
	for d := st.log - 1; d >= 0; d-- {
		l0, r0 := st.bv[d].Rank0(start), st.bv[d].Rank0(end)
		l1, r1 := start+st.mid[d]-l0, end+st.mid[d]-r0
		count1 := count + r0 - l0
		sum1 := op(sum, st.seg[d].Query(l0, r0))
		if check(count1, sum1) {
			count, sum = count1, sum1
			start, end = l1, r1
		} else {
			start, end = l0, r0
		}
	}
	return count, sum
}

// 设置第i个元素的和.
func (st *WaveletMatrixWithSumFast) Set(i int32, e E) {
	left, right := i, i+1
	st.seg[st.log].Set(left, e)
	for d := st.log - 1; d >= 0; d-- {
		l0, r0 := st.bv[d].Rank0(left), st.bv[d].Rank0(right)
		l1, r1 := left+st.mid[d]-l0, right+st.mid[d]-r0
		if l0 < r0 {
			left, right = l0, r0
		}
		if l0 == r0 {
			left, right = l1, r1
		}
		st.seg[d].Set(left, e)
	}
}

// 更新第i个元素的和.
func (st *WaveletMatrixWithSumFast) Update(i int32, e E) {
	left, right := i, i+1
	st.seg[st.log].Update(left, e)
	for d := st.log - 1; d >= 0; d-- {
		l0, r0 := st.bv[d].Rank0(left), st.bv[d].Rank0(right)
		l1, r1 := left+st.mid[d]-l0, right+st.mid[d]-r0
		if l0 < r0 {
			left, right = l0, r0
		}
		if l0 == r0 {
			left, right = l1, r1
		}
		st.seg[d].Update(left, e)
	}
}

func (st *WaveletMatrixWithSumFast) build(arr []int, sum []E) {
	n := int32(len(arr))
	compressed, index := createIndexCompressionSame(arr, st.smallRange)
	upper := int32(0)
	for _, v := range compressed {
		if v > upper {
			upper = v
		}
	}
	upper++
	raw := make([]int, upper)
	for i, v := range arr {
		raw[compressed[i]] = v
	}
	log := int32(0)
	for 1<<log < upper {
		log++
	}
	mid := make([]int32, log)
	bv := make([]*bitVector, log)
	for i := range bv {
		bv[i] = newBitVector(n)
	}
	arr0, arr1 := make([]int32, n), make([]int32, n)
	sum0, sum1 := make([]E, n), make([]E, n)
	seg := make([]*SegTreeLike, log+1)
	for i := range seg {
		seg[i] = &SegTreeLike{}
	}
	seg[log].Build(sum)
	for d := log - 1; d >= 0; d-- {
		p0, p1 := int32(0), int32(0)
		for i := int32(0); i < n; i++ {
			f := (compressed[i] >> d & 1) == 1
			if !f {
				arr0[p0] = compressed[i]
				sum0[p0] = sum[i]
				p0++
			} else {
				bv[d].Set(i)
				arr1[p1] = compressed[i]
				sum1[p1] = sum[i]
				p1++
			}
		}
		compressed, arr0 = arr0, compressed
		sum, sum0 = sum0, sum
		copy(compressed[p0:], arr1[:p1])
		copy(sum[p0:], sum1[:p1])
		mid[d] = p0
		bv[d].Build()
		seg[d].Build(sum)
	}

	st.n, st.log, st.upper = n, log, upper
	st.raw, st.mid, st.bv, st.seg = raw, mid, bv, seg
	st.index = index
}

type bitVector struct {
	bits   []uint64
	preSum []int32
}

func newBitVector(n int32) *bitVector {
	return &bitVector{bits: make([]uint64, n>>6+1), preSum: make([]int32, n>>6+1)}
}

func (bv *bitVector) Set(i int32) {
	bv.bits[i>>6] |= 1 << (i & 63)
}

func (bv *bitVector) Build() {
	for i := 0; i < len(bv.bits)-1; i++ {
		bv.preSum[i+1] = bv.preSum[i] + int32(bits.OnesCount64(bv.bits[i]))
	}
}

func (bv *bitVector) Rank0(k int32) int32 {
	return k - bv.preSum[k>>6] - int32(bits.OnesCount64(bv.bits[k>>6]&((1<<(k&63))-1)))
}

func createIndexCompressionSame(arr []int, smallRange bool) (compressedArr []int32, index func(int) int32) {
	if smallRange {
		return indexCompressionSameSmall(arr)
	} else {
		return indexCompressionSameLarge(arr)
	}
}

func indexCompressionSameSmall(arr []int) (compressedArr []int32, index func(int) int32) {
	var min_, max_ int
	var data []int32
	compressedArr = make([]int32, len(arr))
	for i, v := range arr {
		compressedArr[i] = int32(v)
	}
	min32, max32 := int32(0), int32(-1)
	if len(compressedArr) > 0 {
		for _, x := range compressedArr {
			if x < min32 {
				min32 = x
			}
			if x > max32 {
				max32 = x
			}
		}
	}
	data = make([]int32, max32-min32+2)
	for _, x := range compressedArr {
		data[x-min32+1] = 1
	}
	for i := 0; i < len(data)-1; i++ {
		data[i+1] += data[i]
	}
	for i, v := range compressedArr {
		compressedArr[i] = data[v-min32]
	}
	min_, max_ = int(min32), int(max32)
	index = func(x int) int32 { return data[clamp(x-min_, 0, max_-min_+1)] }
	return
}

func indexCompressionSameLarge(arr []int) (compressedArr []int32, index func(int) int32) {
	var data []int
	order := argSort(arr)
	compressedArr = make([]int32, len(arr))
	for _, v := range order {
		if len(data) == 0 || data[len(data)-1] != arr[v] {
			data = append(data, arr[v])
		}
		compressedArr[v] = int32(len(data) - 1)
	}
	data = data[:len(data):len(data)]
	index = func(x int) int32 { return int32(sort.SearchInts(data, x)) }
	return
}

func clamp(x, min, max int) int {
	if x < min {
		return min
	}
	if x > max {
		return max
	}
	return x
}

func argSort(nums []int) []int32 {
	order := make([]int32, len(nums))
	for i := int32(0); i < int32(len(order)); i++ {
		order[i] = i
	}
	sort.Slice(order, func(i, j int) bool { return nums[order[i]] < nums[order[j]] })
	return order
}

const INF32 int32 = 1 << 30

// PointAddRangeSum.

func (*SegmentTree) e() E        { return 0 }
func (*SegmentTree) op(a, b E) E { return a + b }
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

type SegmentTree struct {
	n, size int32
	seg     []E
}

func NewSegmentTree(n int32, f func(int32) E) *SegmentTree {
	res := &SegmentTree{}
	size := int32(1)
	for size < n {
		size <<= 1
	}
	seg := make([]E, size<<1)
	for i := range seg {
		seg[i] = res.e()
	}
	for i := int32(0); i < n; i++ {
		seg[i+size] = f(i)
	}
	for i := size - 1; i > 0; i-- {
		seg[i] = res.op(seg[i<<1], seg[i<<1|1])
	}
	res.n = n
	res.size = size
	res.seg = seg
	return res
}
func NewSegmentTreeFrom(leaves []E) *SegmentTree {
	res := &SegmentTree{}
	n := int32(len(leaves))
	size := int32(1)
	for size < n {
		size <<= 1
	}
	seg := make([]E, size<<1)
	for i := range seg {
		seg[i] = res.e()
	}
	for i := int32(0); i < n; i++ {
		seg[i+size] = leaves[i]
	}
	for i := size - 1; i > 0; i-- {
		seg[i] = res.op(seg[i<<1], seg[i<<1|1])
	}
	res.n = n
	res.size = size
	res.seg = seg
	return res
}
func (st *SegmentTree) Get(index int32) E {
	if index < 0 || index >= st.n {
		return st.e()
	}
	return st.seg[index+st.size]
}
func (st *SegmentTree) Set(index int32, value E) {
	if index < 0 || index >= st.n {
		return
	}
	index += st.size
	st.seg[index] = value
	for index >>= 1; index > 0; index >>= 1 {
		st.seg[index] = st.op(st.seg[index<<1], st.seg[index<<1|1])
	}
}
func (st *SegmentTree) Update(index int32, value E) {
	if index < 0 || index >= st.n {
		return
	}
	index += st.size
	st.seg[index] = st.op(st.seg[index], value)
	for index >>= 1; index > 0; index >>= 1 {
		st.seg[index] = st.op(st.seg[index<<1], st.seg[index<<1|1])
	}
}

// [start, end)
func (st *SegmentTree) Query(start, end int32) E {
	if start < 0 {
		start = 0
	}
	if end > st.n {
		end = st.n
	}
	if start >= end {
		return st.e()
	}
	leftRes, rightRes := st.e(), st.e()
	start += st.size
	end += st.size
	for start < end {
		if start&1 == 1 {
			leftRes = st.op(leftRes, st.seg[start])
			start++
		}
		if end&1 == 1 {
			end--
			rightRes = st.op(st.seg[end], rightRes)
		}
		start >>= 1
		end >>= 1
	}
	return st.op(leftRes, rightRes)
}
func (st *SegmentTree) QueryAll() E { return st.seg[1] }
func (st *SegmentTree) GetAll() []E {
	res := make([]E, st.n)
	copy(res, st.seg[st.size:st.size+st.n])
	return res
}

type StaticRangeProductGroup struct {
	data []E
}

func NewStaticRangeProductGroup(arr []E) *StaticRangeProductGroup {
	m := int32(len(arr))
	data := make([]E, m+1)
	data[0] = e()
	for i := int32(0); i < m; i++ {
		data[i+1] = op(data[i], arr[i])
	}
	return &StaticRangeProductGroup{data}
}

func (s *StaticRangeProductGroup) Query(start, end int32) E {
	if start >= end {
		return e()
	}
	return op(inv(s.data[start]), s.data[end])
}
