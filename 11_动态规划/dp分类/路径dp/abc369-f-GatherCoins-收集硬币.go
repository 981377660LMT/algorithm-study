package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

// F - Gather Coins (收集硬币，二维LIS)
// https://atcoder.jp/contests/abc369/tasks/abc369_f
// 给定一个H*W的网格, 有N格子有硬币.
// 有一个人在(0,0), 每次可以向右或向下走一格.
// 到达(H-1,W-1)，问最多能收集多少硬币.
// H,W<=2e5,N<=2e5.
//
// !按照(X,Y)坐标排序, dp[i]表示以第i个硬币结尾的最长路径.
// !就是按照X排序后，选区的Y单调不减的最长路径.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var H, W, N int
	fmt.Fscan(in, &H, &W, &N)
	X, Y := make([]int, N), make([]int, N)
	for i := 0; i < N; i++ {
		fmt.Fscan(in, &X[i], &Y[i])
		X[i]--
		Y[i]--
	}

	order := make([]int, N)
	for i := 0; i < N; i++ {
		order[i] = i
	}
	sort.Slice(order, func(i, j int) bool {
		if X[order[i]] != X[order[j]] {
			return X[order[i]] < X[order[j]]
		}
		return Y[order[i]] < Y[order[j]]
	})

	newY, originY := Discretize(Y)
	seg := NewSegmentTree(int32(len(originY)), func(i int32) E { return E{value: -INF, index: -1} })
	dp := make([]int, N)
	pre := make([]int, N)
	for i := 0; i < N; i++ {
		dp[i] = 1
		pre[i] = -1
	}

	for _, id := range order {
		preMax := seg.Query(0, newY[id]+1)
		if preMax.value+1 > dp[id] {
			dp[id] = preMax.value + 1
			pre[id] = preMax.index
		}
		seg.Set(newY[id], E{value: dp[id], index: id})
	}

	bestIndex := maxIndex(dp)
	best := dp[bestIndex]
	path := []int{bestIndex}
	for pre[bestIndex] != -1 {
		bestIndex = pre[bestIndex]
		path = append(path, bestIndex)
	}
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}

	xy := make([][2]int, 0, len(path)+2)
	xy = append(xy, [2]int{0, 0})
	for _, id := range path {
		xy = append(xy, [2]int{X[id], Y[id]})
	}
	xy = append(xy, [2]int{H - 1, W - 1})

	res := strings.Builder{}

	for i := 0; i < len(xy)-1; i++ {
		x1, y1 := xy[i][0], xy[i][1]
		x2, y2 := xy[i+1][0], xy[i+1][1]
		res.WriteString(strings.Repeat("D", x2-x1))
		res.WriteString(strings.Repeat("R", y2-y1))
	}

	fmt.Fprintln(out, best)
	fmt.Fprintln(out, res.String())
}

const INF int = 1e18

// PointSetRangeMaxIndex

type E = struct{ value, index int }

func (*SegmentTree) e() E { return E{value: -INF, index: -1} }
func (*SegmentTree) op(a, b E) E {
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

// 将nums中的元素进行离散化，返回新的数组和对应的原始值.
// origin[newNums[i]] == nums[i]
func Discretize(nums []int) (newNums []int32, origin []int) {
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

func maxIndex(nums []int) int {
	maxIndex := 0
	for i := 1; i < len(nums); i++ {
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

func mins(nums ...int) int {
	res := nums[0]
	for _, num := range nums {
		if num < res {
			res = num
		}
	}
	return res
}

func maxs(nums ...int) int {
	res := nums[0]
	for _, num := range nums {
		if num > res {
			res = num
		}
	}
	return res
}
