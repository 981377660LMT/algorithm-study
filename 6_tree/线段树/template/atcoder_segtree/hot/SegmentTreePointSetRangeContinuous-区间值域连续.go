// SegmentTreePointSetRangeContinuous-区间值域连续

package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

const INF int = 1e18
const MOD int = 1e9 + 7

var fp *FastPow

func init() {
	fp = NewFastPow(2, 1e9+10)
}

type SegmentTreePointSetRangeContinuous struct {
	tree *_SegTree
}

// 0<=nums[i]<=1e9.
func NewSegmentTreePointSetRangeContinuous(nums []int) *SegmentTreePointSetRangeContinuous {
	leaves := make([]E, len(nums))
	for i := range nums {
		leaves[i] = FromElement(nums[i])
	}
	tree := _NewSegTree(leaves)
	return &SegmentTreePointSetRangeContinuous{tree: tree}
}

// 0<=value<=1e9.
func (st *SegmentTreePointSetRangeContinuous) Set(index int, value int) {
	st.tree.Set(index, FromElement(value))
}

// 判断区间[start, end)是否可以重排为值域上连续的一段(区间值域连续).
func (st *SegmentTreePointSetRangeContinuous) Query(start, end int) bool {
	cur := st.tree.Query(start, end)
	min, max := cur.min, cur.max
	pow2 := cur.pow2
	if !(max-min+1 == end-start) {
		return false
	}
	return fp.RangePow2Sum(min, max+1) == pow2
}

// 光速幂.
type FastPow struct {
	max     int
	divData []int
	modData []int
}

// O(sqrt(maxN))预处理,O(1)查询.
//
//	base: 幂运算的基.
//	maxN: 最大的幂.
func NewFastPow(base int, maxN int) *FastPow {
	max := int(math.Ceil(math.Sqrt(float64(maxN))))
	res := &FastPow{max: max, divData: make([]int, max+1), modData: make([]int, max+1)}
	cur := 1
	for i := 0; i <= max; i++ {
		res.modData[i] = cur
		cur = cur * base % MOD
	}
	cur = 1
	last := res.modData[max]
	for i := 0; i <= max; i++ {
		res.divData[i] = cur
		cur = cur * last % MOD
	}
	return res
}

// n<=maxN.
func (fp *FastPow) Pow(n int) int {
	return (fp.divData[n/fp.max] * fp.modData[n%fp.max] % MOD)
}

// 区间以2为底的幂和 (2^start + 2^(start+1) + ... + 2^(end-1)) % MOD.
func (fp *FastPow) RangePow2Sum(start, end int) int {
	if start >= end {
		return 0
	}
	res := (fp.Pow(end) - fp.Pow(start)) % MOD
	if res < 0 {
		res += MOD
	}
	return res
}

type E = struct {
	min, max int
	pow2     int
}

func FromElement(v int) E {
	return E{
		min: v, max: v,
		pow2: fp.Pow(v),
	}
}

func (*_SegTree) e() E {
	return E{
		min: INF, max: -INF,
	}
}

func (*_SegTree) op(a, b E) E {
	return E{
		min: min(a.min, b.min), max: max(a.max, b.max),
		pow2: (a.pow2 + b.pow2) % MOD,
	}
}

type _SegTree struct {
	n, size int
	data    []E
}

func _NewSegTree(leaves []E) *_SegTree {
	res := &_SegTree{}
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
	res.data = seg
	return res
}

func (st *_SegTree) Set(index int, value E) {
	if index < 0 || index >= st.n {
		return
	}
	index += st.size
	st.data[index] = value
	for index >>= 1; index > 0; index >>= 1 {
		st.data[index] = st.op(st.data[index<<1], st.data[index<<1|1])
	}
}

// [start, end)
func (st *_SegTree) Query(start, end int) E {
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
			leftRes = st.op(leftRes, st.data[start])
			start++
		}
		if end&1 == 1 {
			end--
			rightRes = st.op(st.data[end], rightRes)
		}
		start >>= 1
		end >>= 1
	}
	return st.op(leftRes, rightRes)
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

// P3792 由乃与大母神原型和偶像崇拜
// https://www.luogu.com.cn/problem/P3792
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}
	st := NewSegmentTreePointSetRangeContinuous(nums)
	for i := 0; i < q; i++ {
		var op int
		fmt.Fscan(in, &op)
		if op == 1 {
			var pos, val int
			fmt.Fscan(in, &pos, &val)
			pos--
			st.Set(pos, val)
		} else {
			var start, end int
			fmt.Fscan(in, &start, &end)
			start--
			if st.Query(start, end) {
				fmt.Fprintln(out, "damushen")
			} else {
				fmt.Fprintln(out, "yuanxing")
			}
		}
	}
}
