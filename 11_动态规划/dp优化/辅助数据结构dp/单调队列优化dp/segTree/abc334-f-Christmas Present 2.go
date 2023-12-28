// F - Christmas Present 2
// https://atcoder.jp/contests/abc334/tasks/abc334_f

// 快递员送货,起始点为(sx, sy),需要到一些房子houses去送货.
// 送货需要按照顺序送,即先送第一个房子,再送第二个房子,以此类推.
// !快递员每次最多携带k个包裹,中途可以回到起点将包裹补充满.
// 问从起点出发,送完所有房子,回到起点的最短距离是多少.
// n,k<=2e5

// 思路:
// 类似"从仓库到码头运输箱子"
// dp[i] 表示前i个礼物运送完毕时的最短距离 (0<=i<=n)
// dp[i] = dp[j] + (preDist[i]-preDist[j+1]) + (distToStart[i]+distToStart[j+1]) | i - j <= k
// 合并同类项得
// d[i] = (dp[j] - preDist[j+1] + distToStart[j+1]) + preDist[i] + distToStart[i] | i - j <= k
// 线段树维护即可

package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k int
	fmt.Fscan(in, &n, &k)
	var sx, sy int
	fmt.Fscan(in, &sx, &sy)
	houses := make([][2]int, n)
	for i := range houses {
		fmt.Fscan(in, &houses[i][0], &houses[i][1])
	}
	fmt.Fprintln(out, ChristmasPresent2(sx, sy, houses, k))
}

func ChristmasPresent2(sx, sy int, houses [][2]int, k int) float64 {
	dist := func(x1, y1, x2, y2 int) float64 {
		return math.Sqrt(float64((x1-x2)*(x1-x2) + (y1-y2)*(y1-y2)))
	}

	n := len(houses)
	distToStart := make([]float64, n+1)
	for i, h := range houses {
		distToStart[i+1] = dist(sx, sy, h[0], h[1])
	}
	preDist := make([]float64, 0, n+1) // preDist[i] 表示运送前i个礼物的相邻移动距离
	preDist = append(preDist, 0, 0)
	for i := 0; i < n-1; i++ {
		x1, y1 := houses[i][0], houses[i][1]
		x2, y2 := houses[i+1][0], houses[i+1][1]
		preDist = append(preDist, preDist[len(preDist)-1]+dist(x1, y1, x2, y2))
	}

	dp := make([]float64, n+1)
	for i := range dp {
		dp[i] = -INF
	}
	dp[0] = 0
	seg := NewSegmentTree(n+1, func(i int) float64 { return -INF })
	seg.Set(0, dp[0]-preDist[1]+distToStart[1])
	for i := 1; i <= n; i++ {
		preMin := seg.Query(i-k, i)
		dp[i] = preMin + preDist[i] + distToStart[i]
		if i < n {
			seg.Set(i, dp[i]-preDist[i+1]+distToStart[i+1])
		}
	}

	return dp[n]
}

const INF float64 = 1e18

// PointSetRangeMin

type E = float64

func (*SegmentTree) e() E        { return INF }
func (*SegmentTree) op(a, b E) E { return min64(a, b) }
func min64(a, b float64) float64 {
	if a < b {
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

// 二分查询最大的 right 使得切片 [left:right] 内的值满足 predicate
func (st *SegmentTree) MaxRight(left int, predicate func(E) bool) int {
	if left == st.n {
		return st.n
	}
	left += st.size
	res := st.e()
	for {
		for left&1 == 0 {
			left >>= 1
		}
		if !predicate(st.op(res, st.seg[left])) {
			for left < st.size {
				left <<= 1
				if tmp := st.op(res, st.seg[left]); predicate(tmp) {
					res = tmp
					left++
				}
			}
			return left - st.size
		}
		res = st.op(res, st.seg[left])
		left++
		if (left & -left) == left {
			break
		}
	}
	return st.n
}

// 二分查询最小的 left 使得切片 [left:right] 内的值满足 predicate
func (st *SegmentTree) MinLeft(right int, predicate func(E) bool) int {
	if right == 0 {
		return 0
	}
	right += st.size
	res := st.e()
	for {
		right--
		for right > 1 && right&1 == 1 {
			right >>= 1
		}
		if !predicate(st.op(st.seg[right], res)) {
			for right < st.size {
				right = right<<1 | 1
				if tmp := st.op(st.seg[right], res); predicate(tmp) {
					res = tmp
					right--
				}
			}
			return right + 1 - st.size
		}
		res = st.op(st.seg[right], res)
		if right&-right == right {
			break
		}
	}
	return 0
}
