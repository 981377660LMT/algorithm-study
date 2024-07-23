// 区间哈希(带修改，multiset哈希，区间多重集哈希)
// !两个数列区间排序后对应的相邻元素的差是否完全相等，并且支持单点修改
// https://www.luogu.com.cn/problem/P6688
// 0 index value -> 将第index个数变为value
// 1 l1 r1 l2 r2 -> 判断区间[l1, r1)和[l2, r2)是否本质相同
//
// 多项式哈希
// !hash(start,end) = ∑base^nums[i], i∈[start,end)
// 好处：区间所有数加上一个数后哈希值容易计算.

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

	var n, q int32
	fmt.Fscan(in, &n, &q)
	nums := make([]int32, n)
	for i := int32(0); i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	seg := NewSegmentTree(n, func(i int32) E { return FromElement(nums[i]) })

	update := func(index int32, value int32) {
		seg.Set(index, FromElement(value))
	}

	query := func(l1, r1, l2, r2 int32) bool {
		e1 := seg.Query(l1, r1)
		e2 := seg.Query(l2, r2)
		if e1.min_ > e2.min_ {
			e1, e2 = e2, e1
		}
		return qpow(BASE, int(e2.min_-e1.min_))*e1.hash%MOD == e2.hash
	}

	for i := int32(0); i < q; i++ {
		var kind int32
		fmt.Fscan(in, &kind)
		if kind == 0 {
			var index, value int32
			fmt.Fscan(in, &index, &value)
			index--
			update(index, value)
		} else {
			var l1, r1, l2, r2 int32
			fmt.Fscan(in, &l1, &r1, &l2, &r2)
			l1--
			l2--
			if query(l1, r1, l2, r2) {
				fmt.Fprintln(out, "YES")
			} else {
				fmt.Fprintln(out, "NO")
			}
		}
	}
}

const INF int = 1e18
const INF32 int32 = 1e9 + 10
const MOD int = 998244353
const BASE int = 131

func qpow(a, b int) int {
	res := 1
	for b > 0 {
		if b&1 == 1 {
			res = res * a % MOD
		}
		a = a * a % MOD
		b >>= 1
	}
	return res
}

// PointSetRangeMinPow

type E = struct {
	min_ int32
	hash int
}

func FromElement(v int32) E { return E{v, qpow(BASE, int(v))} }
func (*SegmentTree) e() E   { return E{INF32, 0} }
func (*SegmentTree) op(a, b E) E {
	newMin := min32(a.min_, b.min_)
	newHash := (a.hash + b.hash)
	if newHash >= MOD {
		newHash -= MOD
	}
	return E{newMin, newHash}
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

func min32(a, b int32) int32 {
	if a < b {
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
