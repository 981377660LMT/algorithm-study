package main

import (
	"bufio"
	"fmt"
	"os"
)

// CF1093G-Multidimensional Queries
// https://www.luogu.com.cn/problem/CF1093G
// 给定n个k维点(k<=5).
// 1 i x1 x2 ... xk: 第i个点的坐标更新为(x1,x2,...,xk);
// 2 start end ：查询区间[start,end)内最大的两点间曼哈顿距离.
//
// !习惯性的把曼哈顿距离的绝对值拆出来，用二进制表示
// 31 的二进制表示是 11111，表示 5 维的一个点的坐标加入的正负情况都为正
// 即 x[0] - y[0] + x[1] - y[1] + x[2] - y[2] + x[3] - y[3] + x[4] - y[4]
// 29 的二进制表示是 11101，表示 5 维的一个点的坐标加入的正负情况为正、负、正、正、正
// 即 x[0] - y[0] + x[1] - y[1] + x[2] - y[2] - x[3] + y[3] + x[4] - y[4]
// 那么要求的就是 max(f[0]+f[31],f[1]+f[30],f[2]+f[29],...,f[15]+f[16])
// 用线段树维护最大值即可.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k int
	fmt.Fscan(in, &n, &k)
	points := make([][]int, n)
	for i := 0; i < n; i++ {
		points[i] = make([]int, k)
		for j := 0; j < k; j++ {
			fmt.Fscan(in, &points[i][j])
		}
	}

	fromPoint := func(point []int) E {
		var res E
		res[0] = 0
		for i := 0; i < 1<<k; i++ {
			for j := 0; j < k; j++ {
				if i&(1<<j) > 0 {
					res[i] += point[j]
				} else {
					res[i] -= point[j]
				}
			}
		}
		return res
	}

	var q int
	fmt.Fscan(in, &q)
	seg := NewSegmentTree(n, func(i int) E { return fromPoint(points[i]) })
	for i := 0; i < q; i++ {
		var kind int
		fmt.Fscan(in, &kind)
		if kind == 1 {
			var pos int
			fmt.Fscan(in, &pos)
			pos--
			point := make([]int, k)
			for j := 0; j < k; j++ {
				fmt.Fscan(in, &point[j])
			}
			seg.Set(pos, fromPoint(point))
		} else {
			var start, end int
			fmt.Fscan(in, &start, &end)
			start--
			tmp := seg.Query(start, end)
			res := 0
			for i := 0; i < 1<<k; i++ {
				res = max(res, tmp[i]+tmp[(1<<k)-1-i])
			}
			fmt.Fprintln(out, res)
		}
	}
}

const INF int = 1e18

// PointSetRangeMax

type E = [32]int

func (*SegmentTree) e() E { return [32]int{INF} }
func (*SegmentTree) op(a, b E) E {
	if a[0] == INF {
		return b
	}
	if b[0] == INF {
		return a
	}
	for i := 0; i < 32; i++ {
		a[i] = max(a[i], b[i])
	}
	return a
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
