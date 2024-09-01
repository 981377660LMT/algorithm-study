// 二维上升子序列(二维LIS)
// 按照横坐标从小到大排序后，选取子序列最长且满足纵坐标单调不减，
// 那么题目就转化成了二维偏序问题，可以用最长不下降子序列解决。

package main

import (
	"fmt"
	"sort"
)

func main() {
	xs, ys := []int{3, 2, 2, 1}, []int{3, 1, 3, 4}
	fmt.Println(Lis2D(xs, ys, true, true))   // [1 0] [2 1 1 1]
	fmt.Println(Lis2D(xs, ys, true, false))  // [2 0] [2 1 1 1]
	fmt.Println(Lis2D(xs, ys, false, true))  // [1 0] [2 1 2 1]
	fmt.Println(Lis2D(xs, ys, false, false)) // [1 2 0] [3 1 2 1]
}

const INF int = 1e18
const INF32 int32 = 1e9 + 10

// 返回非严格最长上升子序列的点的下标、每个点为结尾的LIS长度.
func Lis2D(xs, ys []int, xStrict bool, yStrict bool) (lis []int32, dp []int32) {
	n := int32(len(xs))
	order := make([]int32, n)
	for i := int32(0); i < n; i++ {
		order[i] = i
	}
	sort.Slice(order, func(i, j int) bool {
		if xs[order[i]] != xs[order[j]] {
			return xs[order[i]] < xs[order[j]]
		}
		return ys[order[i]] < ys[order[j]]
	})

	newY, originY := discretize(ys)
	seg := newSegmentTree(int32(len(originY)), func(i int32) E { return E{value: -INF32, index: -1} })
	dp = make([]int32, n)
	pre := make([]int32, n)
	for i := int32(0); i < n; i++ {
		dp[i] = 1
		pre[i] = -1
	}

	if xStrict {
		// x相同的元素一起处理
		enumerateGroupByKey(order, func(index int32) int { return xs[order[index]] }, func(start, end int32) {
			preMaxs := make([]E, 0, end-start)
			for _, id := range order[start:end] {
				upper := newY[id] + 1
				if yStrict {
					upper--
				}
				preMaxs = append(preMaxs, seg.Query(0, upper))
			}
			for i, id := range order[start:end] {
				preMax := preMaxs[i]
				if preMax.value+1 > dp[id] {
					dp[id] = preMax.value + 1
					pre[id] = preMax.index
				}
				seg.Set(newY[id], E{value: dp[id], index: id})
			}
		})
	} else {
		for _, id := range order {
			upper := newY[id] + 1
			if yStrict {
				upper--
			}
			preMax := seg.Query(0, upper)
			if preMax.value+1 > dp[id] {
				dp[id] = preMax.value + 1
				pre[id] = preMax.index
			}
			seg.Set(newY[id], E{value: dp[id], index: id})
		}

	}

	bestIndex := maxIndex(dp)
	lis = []int32{bestIndex}
	for pre[bestIndex] != -1 {
		bestIndex = pre[bestIndex]
		lis = append(lis, bestIndex)
	}
	for i, j := 0, len(lis)-1; i < j; i, j = i+1, j-1 {
		lis[i], lis[j] = lis[j], lis[i]
	}
	return
}

// PointSetRangeMaxIndex

type E = struct{ value, index int32 }

func (*segmentTree) e() E { return E{value: -INF32, index: -1} }
func (*segmentTree) op(a, b E) E {
	if a.value > b.value {
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
func max32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

type segmentTree struct {
	n, size int32
	seg     []E
}

func newSegmentTree(n int32, f func(int32) E) *segmentTree {
	res := &segmentTree{}
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
func newSegmentTreeFrom(leaves []E) *segmentTree {
	res := &segmentTree{}
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
func (st *segmentTree) Get(index int32) E {
	if index < 0 || index >= st.n {
		return st.e()
	}
	return st.seg[index+st.size]
}
func (st *segmentTree) Set(index int32, value E) {
	if index < 0 || index >= st.n {
		return
	}
	index += st.size
	st.seg[index] = value
	for index >>= 1; index > 0; index >>= 1 {
		st.seg[index] = st.op(st.seg[index<<1], st.seg[index<<1|1])
	}
}
func (st *segmentTree) Update(index int32, value E) {
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
func (st *segmentTree) Query(start, end int32) E {
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
func (st *segmentTree) QueryAll() E { return st.seg[1] }
func (st *segmentTree) GetAll() []E {
	res := make([]E, st.n)
	copy(res, st.seg[st.size:st.size+st.n])
	return res
}

// 将nums中的元素进行离散化，返回新的数组和对应的原始值.
// origin[newNums[i]] == nums[i]
func discretize(nums []int) (newNums []int32, origin []int) {
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

func argSort(n int32, less func(i, j int32) bool) []int32 {
	order := make([]int32, n)
	for i := range order {
		order[i] = int32(i)
	}
	sort.Slice(order, func(i, j int) bool { return less(order[i], order[j]) })
	return order
}

func enumerateGroupByKey[E any, K comparable](arr []E, key func(index int32) K, f func(start, end int32)) {
	n := int32(len(arr))
	end := int32(0)
	for end < n {
		start := end
		leader := key(end)
		end++
		for end < n && key(end) == leader {
			end++
		}
		f(start, end) // [start, end)
	}
}

func maxIndex(nums []int32) int32 {
	maxIndex := int32(0)
	for i := int32(0); i < int32(len(nums)); i++ {
		if nums[i] > nums[maxIndex] {
			maxIndex = i
		}
	}
	return maxIndex
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
