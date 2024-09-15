package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"sort"
)

// CF1181D Irrigation 灌溉 (二分矩形和+waveletmatrix)
// https://www.luogu.com.cn/problem/CF1181D
// 有n 个宽度相等高度不同的矩形，给出初始高度，每次选取最矮的矩形中编号最小的，并把它的高度+1，
// 然后不断重复这个操作，要求回答第K 次选取的矩形的编号
// !3266. K 次乘运算后的最终数组 II 的加法版本
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var M, N, Q int
	fmt.Fscan(in, &M, &N, &Q)
	A := make([]int, N) // 每个城市举办比赛的次数
	for i := 0; i < M; i++ {
		var city int
		fmt.Fscan(in, &city)
		A[city-1]++
	}

	order := make([]int, N)
	for i := 0; i < N; i++ {
		order[i] = i
	}
	sort.Slice(order, func(i, j int) bool {
		if A[order[i]] != A[order[j]] {
			return A[order[i]] < A[order[j]]
		}
		return i < j
	})
	sortedNums := make([]int, N)
	for i := 0; i < N; i++ {
		sortedNums[i] = A[order[i]]
	}
	presum := make([]int, N+1)
	for i := 0; i < N; i++ {
		presum[i+1] = presum[i] + sortedNums[i]
	}

	wm := NewWaveletMatrixFast(true)
	wm.Build(int32(N), func(i int32) int { return order[i] })

	// !二分可以与哪个数齐平
	// !再填充剩余的
	// k>=1.
	query := func(k int) int {
		right := MaxRight(0, func(right int) bool {
			return sortedNums[right-1]*right-presum[right] < k
		}, N)
		sum := sortedNums[right-1]*right - presum[right]
		remain := k - sum
		remain %= right
		if remain == 0 {
			remain = right
		}
		kth := wm.Kth(0, int32(right), int32(remain)-1)
		return kth
	}

	for i := 0; i < Q; i++ {
		var K int
		fmt.Fscan(in, &K)
		K -= M
		fmt.Fprintln(out, query(K)+1)
	}
}

// 返回最大的 right 使得 [left,right) 内的值满足 check.
// !注意check内的right不包含，使用时需要right-1.
// right<=upper.
func MaxRight(left int, check func(right int) bool, upper int) int {
	ok, ng := left, upper+1
	for ok+1 < ng {
		mid := (ok + ng) >> 1
		if check(mid) {
			ok = mid
		} else {
			ng = mid
		}
	}
	return ok
}

const INF int = 1e18

type WaveletMatrixFast struct {
	n, log, upper int32

	raw []int
	mid []int32
	bv  []*bitVector

	smallRange bool
	index      func(int) int32
}

func NewWaveletMatrixFast(smallRange bool) *WaveletMatrixFast {
	return &WaveletMatrixFast{smallRange: smallRange}
}

func (st *WaveletMatrixFast) Build(m int32, f func(i int32) (v int)) {
	arr := make([]int, m)
	for i := int32(0); i < m; i++ {
		arr[i] = f(i)
	}
	st.build(arr)
}

// [start, end) x [0, y)
func (st *WaveletMatrixFast) CountPrefix(start, end int32, y int) int32 {
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
func (st *WaveletMatrixFast) Count(start, end int32, y1, y2 int) int32 {
	if y1 >= y2 {
		return 0
	}
	return st.CountPrefix(start, end, y2) - st.CountPrefix(start, end, y1)
}

// [start, end)区间内第k(k>=0)小的元素.
func (st *WaveletMatrixFast) Kth(start, end int32, k int32) int {
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

// <= y 的最大值. 不存在则返回 -INF.
func (st *WaveletMatrixFast) Prev(start, end int32, y int) int {
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
func (st *WaveletMatrixFast) Next(start, end int32, y int) int {
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
func (st *WaveletMatrixFast) Median(start, end int32, upper bool) int {
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

func (st *WaveletMatrixFast) build(arr []int) {
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
	for d := log - 1; d >= 0; d-- {
		p0, p1 := int32(0), int32(0)
		for i := int32(0); i < n; i++ {
			f := (compressed[i] >> d & 1) == 1
			if !f {
				arr0[p0] = compressed[i]

				p0++
			} else {
				bv[d].Set(i)
				arr1[p1] = compressed[i]
				p1++
			}
		}
		compressed, arr0 = arr0, compressed
		copy(compressed[p0:], arr1[:p1])
		mid[d] = p0
		bv[d].Build()
	}

	st.n, st.log, st.upper = n, log, upper
	st.raw, st.mid, st.bv = raw, mid, bv
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
