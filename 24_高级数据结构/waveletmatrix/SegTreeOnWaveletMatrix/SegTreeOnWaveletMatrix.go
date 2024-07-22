package main

import (
	"fmt"
	"math/bits"
	"math/rand"
	"sort"
	"time"
)

func main() {
	test()
	wm := NewSegTreeOnWaveletMatrix(true)
	wm.Build(10, func(i int32) (int, int) { return int(i), 1 })
	fmt.Println(wm.CountPrefix(4, 10, 5))
	fmt.Println(wm.SumIndexRange(4, 10, 5, 7))

	testTime()
}

const INF int = 2e18

type E = int

func e() E        { return 0 }
func op(a, b E) E { return a + b }
func inv(a E) E   { return -a } // 可选, 查询范围和时用于差分

// (类)线段树的实现.
type SegTreeLike struct {
	seg *SegmentTree // StaticRangeProduct
}

func (st *SegTreeLike) Build(sum []E) {
	st.seg = NewSegmentTreeFrom(sum)
}

func (st *SegTreeLike) Query(start, end int32) E {
	return st.seg.Query(start, end)
}

func (st *SegTreeLike) Set(i int32, e E) {
	st.seg.Set(i, e)
}

func (st *SegTreeLike) Update(i int32, e E) {
	st.seg.Update(i, e)
}

// ------------------- SegTreeOnWaveletMatrix -------------------

var _ ISegTreeLike = (*SegTreeLike)(nil)

type ISegTreeLike interface {
	Build(sum []E)
	Query(start, end int32) E
	Set(i int32, e E)
	Update(i int32, e E)
}

// WaveletMatrix套线段树.
// 内部不使用接口约束，因为接口会导致性能下降.
type SegTreeOnWaveletMatrix struct {
	n, log, upper int32

	raw []int
	mid []int32
	bv  []*bitVector
	seg []*SegTreeLike

	smallRange bool
	index      func(int) int32
}

func NewSegTreeOnWaveletMatrix(smallRange bool) *SegTreeOnWaveletMatrix {
	return &SegTreeOnWaveletMatrix{smallRange: smallRange}
}

func (st *SegTreeOnWaveletMatrix) Build(m int32, f func(i int32) (v int, e E)) {
	arr, sum := make([]int, m), make([]E, m)
	for i := int32(0); i < m; i++ {
		arr[i], sum[i] = f(i)
	}
	st.build(arr, sum)
}

// [start, end) x [0, y)
func (st *SegTreeOnWaveletMatrix) CountPrefix(start, end int32, y int) int32 {
	p := st.index(y)
	if p == 0 {
		return 0
	}
	if p == st.upper {
		return end - start
	}
	res := int32(0)
	for d := st.log - 1; d >= 0; d-- {
		l0, r0 := st.bv[d].Rank(start, false), st.bv[d].Rank(end, false)
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
func (st *SegTreeOnWaveletMatrix) Count(start, end int32, y1, y2 int) int32 {
	if y1 >= y2 {
		return 0
	}
	return st.CountPrefix(start, end, y2) - st.CountPrefix(start, end, y1)
}

// [start, end) x [0, y)
func (st *SegTreeOnWaveletMatrix) CountAndSumPrefix(start, end int32, y int) (int32, E) {
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
		l0, r0 := st.bv[d].Rank(start, false), st.bv[d].Rank(end, false)
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
func (st *SegTreeOnWaveletMatrix) CountAndSum(start, end int32, y1, y2 int) (int32, E) {
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
		l0, r0 := st.bv[d].Rank(L, false), st.bv[d].Rank(R, false)
		l1, r1 := L+st.mid[d]-l0, R+st.mid[d]-r0
		dfs(d, l0, r0, a, c)
		dfs(d, l1, r1, c, b)
	}
	dfs(st.log, start, end, 0, 1<<st.log)
	return count, sum
}

// [start, end) x [0, y)
func (st *SegTreeOnWaveletMatrix) SumPrefix(start, end int32, y int) E {
	_, sum := st.CountAndSumPrefix(start, end, y)
	return sum
}

// [start, end) x [y1, y2)
func (st *SegTreeOnWaveletMatrix) Sum(start, end int32, y1, y2 int) E {
	_, sum := st.CountAndSum(start, end, y1, y2)
	return sum
}

func (st *SegTreeOnWaveletMatrix) SumAll(start, end int32) E {
	return st.seg[st.log].Query(start, end)
}

// 排名在[k1, k2)间的元素的和.要求运算存在逆元.
func (st *SegTreeOnWaveletMatrix) SumIndexRange(start, end int32, k1, k2 int32) E {
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
func (st *SegTreeOnWaveletMatrix) Kth(start, end int32, k int32) int {
	if k < 0 {
		k = 0
	}
	if n := end - start - 1; k > n {
		k = n
	}
	p := int32(0)
	for d := st.log - 1; d >= 0; d-- {
		l0, r0 := st.bv[d].Rank(start, false), st.bv[d].Rank(end, false)
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
func (st *SegTreeOnWaveletMatrix) KthValueAndSum(start, end int32, k int32) (int, E) {
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
		l0, r0 := st.bv[d].Rank(start, false), st.bv[d].Rank(end, false)
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
func (st *SegTreeOnWaveletMatrix) Prev(start, end int32, y int) int {
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
		l0, r0 := st.bv[d].Rank(L, false), st.bv[d].Rank(R, false)
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
func (st *SegTreeOnWaveletMatrix) Next(start, end int32, y int) int {
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
		l0, r0 := st.bv[d].Rank(L, false), st.bv[d].Rank(R, false)
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
func (st *SegTreeOnWaveletMatrix) Median(start, end int32, upper bool) int {
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
func (st *SegTreeOnWaveletMatrix) MaxRight(start, end int32, check func(int32, E) bool) (int32, E) {
	if start >= end {
		return 0, e()
	}
	if s := st.SumAll(start, end); check(end-start, s) {
		return end - start, s
	}
	count := int32(0)
	sum := e()
	for d := st.log - 1; d >= 0; d-- {
		l0, r0 := st.bv[d].Rank(start, false), st.bv[d].Rank(end, false)
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

func (st *SegTreeOnWaveletMatrix) Set(i int32, e E) {
	left, right := i, i+1
	st.seg[st.log].Set(left, e)
	for d := st.log - 1; d >= 0; d-- {
		l0, r0 := st.bv[d].Rank(left, false), st.bv[d].Rank(right, false)
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

func (st *SegTreeOnWaveletMatrix) Update(i int32, e E) {
	left, right := i, i+1
	st.seg[st.log].Update(left, e)
	for d := st.log - 1; d >= 0; d-- {
		l0, r0 := st.bv[d].Rank(left, false), st.bv[d].Rank(right, false)
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

func (st *SegTreeOnWaveletMatrix) build(arr []int, sum []E) {
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

func (bv *bitVector) Rank(k int32, f bool) int32 {
	m, s := bv.bits[k>>6], bv.preSum[k>>6]
	res := s + int32(bits.OnesCount64(m&((1<<(k&63))-1)))
	if f {
		return res
	}
	return k - res
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

func test() {

	for i := 0; i < 20; i++ {
		nums := make([]int, 3000)
		for j := 0; j < 3000; j++ {
			nums[j] = rand.Intn(1000)
		}
		wm := NewSegTreeOnWaveletMatrix(false)
		wm.Build(int32(len(nums)), func(i int32) (int, int) { return nums[i], nums[i] })

		rangeFreqBf := func(start, end int32, x, y int) int32 {
			res := int32(0)
			for i := start; i < end; i++ {
				if nums[i] >= x && nums[i] < y {
					res++
				}
			}
			return res
		}
		_ = rangeFreqBf

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

		setBf := func(index int32, v int) {
			nums[index] = v
		}

		for j := 0; j < 2000; j++ {
			start, end := rand.Intn(1000), rand.Intn(1000)
			if start > end {
				start, end = end, start
			}

			x := rand.Intn(1000)
			y := rand.Intn(1000)
			if res1, res2 := rangeFreqBf(int32(start), int32(end), x, y), wm.Count(int32(start), int32(end), x, y); res1 != res2 {
				fmt.Println(res1, res2, start, end, x, y)
				panic("rangeFreqBf")
			}

			k := rand.Intn(max(1, end-start))
			if k > 0 {
				if res1, res2 := kthSmallestBf(int32(start), int32(end), int32(k)), wm.Kth(int32(start), int32(end), int32(k)); res1 != res2 {
					fmt.Println(res1, res2, start, end, k)
					panic("kthSmallestBf")
				}
			}

			// setIndex := int32(rand.Intn(len(nums)))
			// setValue := rand.Intn(1000)
			// setBf(setIndex, setValue)
			// wm.Set(setIndex, setValue)
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

	wm := NewSegTreeOnWaveletMatrix( /*smallRange=*/ false)
	wm.Build(n, func(i int32) (int, int) { return nums[i], nums[i] })
	time1 := time.Now()

	for i := int32(0); i < n; i++ {
		wm.Count(0, n, 0, nums[i])
		wm.Sum(0, n, 0, nums[i])
		wm.CountPrefix(0, n, nums[i])
		wm.SumPrefix(0, n, nums[i])
		wm.Kth(0, n, i)
		wm.Prev(0, n, nums[i])
		wm.Next(0, n, nums[i])
		wm.Median(0, n, true)
		wm.Median(0, n, false)
		wm.KthValueAndSum(0, n, i)
		wm.SumIndexRange(0, n, i, i+1)
		wm.MaxRight(0, n, func(count int32, sum int) bool { return true })
		wm.Set(i, nums[i])
		wm.Update(i, nums[i])
	}

	fmt.Println(time.Since(time1)) // 726.0624ms

}
