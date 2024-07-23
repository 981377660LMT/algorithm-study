// P3792 由乃与大母神原型和偶像崇拜-线段树维护区间二次幂(pow2)的和
// 线段树维护区间 pow2Sum
// 区间值域连续

package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

const INF int = 1e18
const INF32 int32 = 1e9 + 10
const MOD int = 1e9 + 7

var fp *FastPow

func init() {
	fp = NewFastPow(2, 1e9+10)
}

// 1 start end : 判断区间[start, end)是否可以重排为值域上连续的一段(区间值域连续)
// 2 pos val : 将 nums[pos] 修改为 val
//
// !通过求区间的 max,min 求出这个区间的左右端点，然后用公式验证这个区间的和是否符合连续区间的性质
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

	query := func(start, end int32) bool {
		cur := seg.Query(start, end)
		min, max := cur.min, cur.max
		pow2 := cur.pow2
		if ok1 := max-min+1 == end-start; !ok1 {
			return false
		}
		if ok2 := fp.RangePow2Sum(min, max+1) == pow2; !ok2 {
			return false
		}
		return true
	}

	for i := int32(0); i < q; i++ {
		var op int8
		fmt.Fscan(in, &op)
		if op == 1 {
			var pos, val int32
			fmt.Fscan(in, &pos, &val)
			pos--
			update(pos, val)
		} else {
			var start, end int32
			fmt.Fscan(in, &start, &end)
			start--
			if query(start, end) {
				fmt.Fprintln(out, "damushen")
			} else {
				fmt.Fprintln(out, "yuanxing")
			}
		}
	}
}

// 光速幂.
type FastPow struct {
	max     int32
	divData []int
	modData []int
}

// O(sqrt(maxN))预处理,O(1)查询.
//
//	base: 幂运算的基.
//	maxN: 最大的幂.
func NewFastPow(base int, maxN int32) *FastPow {
	max := int32(math.Ceil(math.Sqrt(float64(maxN))))
	res := &FastPow{max: max, divData: make([]int, max+1), modData: make([]int, max+1)}
	cur := 1
	for i := int32(0); i <= max; i++ {
		res.modData[i] = cur
		cur = cur * base % MOD
	}
	cur = 1
	last := res.modData[max]
	for i := int32(0); i <= max; i++ {
		res.divData[i] = cur
		cur = cur * last % MOD
	}
	return res
}

// n<=maxN.
func (fp *FastPow) Pow(n int32) int {
	return (fp.divData[n/fp.max] * fp.modData[n%fp.max] % MOD)
}

// 区间以2为底的幂和 (2^start + 2^(start+1) + ... + 2^(end-1)) % MOD.
func (fp *FastPow) RangePow2Sum(start, end int32) int {
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
	min, max int32
	pow2     int
}

func FromElement(v int32) E {
	return E{
		min: v, max: v,
		pow2: fp.Pow(v),
	}
}

func (*SegmentTree) e() E {
	return E{
		min: INF32, max: -INF32,
		pow2: 0,
	}
}

func (*SegmentTree) op(a, b E) E {
	newMin := min32(a.min, b.min)
	newMax := max32(a.max, b.max)
	newPow2 := a.pow2 + b.pow2
	if newPow2 >= MOD {
		newPow2 -= MOD
	}
	return E{min: newMin, max: newMax, pow2: newPow2}
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
func max32(a, b int32) int32 {
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
