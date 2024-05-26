// abc353-g-MerchangTakahashi
// 有 n个村庄，分别为 1,2,3…n。
// 从第 i 个村庄到第 j 个村庄所需代价为 c×|i−j|。
// 在这些村庄中先后要赶 m次集。每次赶集用一个二元组 (ti​,pi​) 表示，
// 其中ti​ 是赶集的地点，pi​ 是高桥去了能赚到的钱，可以不去.
// 高桥移动的时间忽略不计，初始钱数为 0 且钱数可以为负。求最终赚到钱的最大值。
//
// !dp[i][pos] = max(dp[i-1][j] + p[j] - c * abs(j - pos))
// 去绝对值后，线段树维护前缀/后缀最大值即可
// dp[i][pos] = max(dp[i-1][j] + p[j] - c * (pos - j)) (j < pos)
// dp[i][pos] = max(dp[i-1][j] + p[j] - c * (j - pos)) (j > pos)
// 线段树维护 dp[j] + c * j 和 dp[j] - c * j 的最大值即可

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var N, C int
	fmt.Fscan(in, &N, &C)
	var M int
	fmt.Fscan(in, &M)
	pairs := make([][2]int, M) // (pos, money)
	for i := 0; i < M; i++ {
		fmt.Fscan(in, &pairs[i][0], &pairs[i][1])
		pairs[i][0]--
	}

	seg0 := NewSegmentTree(N, func(i int) int { return -INF })
	seg1 := NewSegmentTree(N, func(i int) int { return -INF })
	// 初始金额 INF
	seg0.Set(0, INF)
	seg1.Set(0, INF)
	res := INF
	for _, pair := range pairs {
		t, p := pair[0], pair[1]
		best := max(seg0.Query(0, t+1)-C*t+p, seg1.Query(t, N)+C*t+p)
		res = max(res, best)
		seg0.Update(t, best+C*t)
		seg1.Update(t, best-C*t)
	}

	fmt.Fprintln(out, res-INF)
}

const INF int = 1e18

// PointSetRangeMax

type E = int

func (*SegmentTree) e() E        { return -INF }
func (*SegmentTree) op(a, b E) E { return max(a, b) }
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

type SegmentTree struct {
	n, size int
	seg     []E
}

func NewSegmentTree(n int, f func(int) E) *SegmentTree {
	res := &SegmentTree{}
	size := 1
	for size < n {
		size <<= 1
	}
	seg := make([]E, size<<1)
	for i := range seg {
		seg[i] = res.e()
	}
	for i := 0; i < n; i++ {
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
	n := len(leaves)
	size := 1
	for size < n {
		size <<= 1
	}
	seg := make([]E, size<<1)
	for i := range seg {
		seg[i] = res.e()
	}
	for i := 0; i < n; i++ {
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
func (st *SegmentTree) Get(index int) E {
	if index < 0 || index >= st.n {
		return st.e()
	}
	return st.seg[index+st.size]
}
func (st *SegmentTree) Set(index int, value E) {
	if index < 0 || index >= st.n {
		return
	}
	index += st.size
	st.seg[index] = value
	for index >>= 1; index > 0; index >>= 1 {
		st.seg[index] = st.op(st.seg[index<<1], st.seg[index<<1|1])
	}
}
func (st *SegmentTree) Update(index int, value E) {
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
func (st *SegmentTree) Query(start, end int) E {
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
