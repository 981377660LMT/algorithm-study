// One Occurrence (离线查询+扫描线,区间只出现一次的数,区间频率为1)
// https://www.luogu.com.cn/problem/CF1000F
// 给定一个长度为 n 的序列，q 次询问，每次询问给定一个区间 [l,r]，
// 如果这个区间里存在只出现一次的数，输出这个数（如果有多个就输出任意一个），没有就输出 0。
// n,q<=5e5

// !定义一个新序列 last, last[v]表示之前与v相同的数中最大的下标（如果不存在则设为 −1）。
// !出现一次等价于 last[v] < start。
// 只保留每个数最后一次出现的 last 值，而将前面的值作废,经过这样的操作，样例处理结果即为：
// nums: 1   1   2   3   2   4
// last: INF 0  INF  -1  2   -1
// 由于在维护靠后的 last 时会对前面已经维护的 last 造成影响，故我们需要将询问离线，并按区间右端点排序.

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

	var n int32
	fmt.Fscan(in, &n)
	nums := make([]int32, n)
	for i := int32(0); i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	var q int32
	fmt.Fscan(in, &q)
	groupByEnd := make([][]int32, n+1)
	queries := make([][2]int32, q)
	for i := int32(0); i < q; i++ {
		var start, end int32
		fmt.Fscan(in, &start, &end)
		start--
		queries[i] = [2]int32{start, end}
		groupByEnd[end] = append(groupByEnd[end], i)
	}

	res := make([]int32, q)
	for i := range res {
		res[i] = -1
	}
	D := NewDictionary()
	for i, v := range nums {
		nums[i] = D.Id(v)
	}
	last := make([]int32, D.Size()) // last[v] 表示之前与v相同的数中最大的下标（如果不存在则设为 −1）。
	for i := range last {
		last[i] = -1
	}
	seg := NewSegmentTree(n, func(i int32) int32 { return INF32 })

	for end := int32(0); end < n+1; end++ {
		for _, qid := range groupByEnd[end] {
			start := queries[qid][0]
			minLast := seg.Query(start, end)
			if minLast >= start {
				res[qid] = -1
			} else {
				index := seg.MaxRight(start, func(x int32) bool { return x > minLast })
				res[qid] = nums[index]
			}
		}

		if end == n {
			break
		}
		num := nums[end]
		pre := last[num]
		last[num] = end
		if pre >= 0 {
			seg.Set(pre, INF32)
		}
		seg.Set(end, pre)
	}

	for _, v := range res {
		if v == -1 {
			fmt.Fprintln(out, 0)
		} else {
			fmt.Fprintln(out, D.Value(v))
		}
	}
}

type V = int32
type Dictionary struct {
	_idToValue []V
	_valueToId map[V]int32
}

// A dictionary that maps values to unique ids.
func NewDictionary() *Dictionary {
	return &Dictionary{
		_valueToId: map[V]int32{},
	}
}
func (d *Dictionary) Id(value V) int32 {
	res, ok := d._valueToId[value]
	if ok {
		return res
	}
	id := int32(len(d._idToValue))
	d._idToValue = append(d._idToValue, value)
	d._valueToId[value] = id
	return id
}
func (d *Dictionary) Value(id int32) V {
	return d._idToValue[id]
}
func (d *Dictionary) Has(value V) bool {
	_, ok := d._valueToId[value]
	return ok
}
func (d *Dictionary) Size() int32 {
	return int32(len(d._idToValue))
}

const INF32 int32 = 1 << 30

// PointSetRangeMin

type E = int32

func (*SegmentTree) e() E        { return INF32 }
func (*SegmentTree) op(a, b E) E { return min32(a, b) }
func min32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}
func max(a, b int32) int32 {
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

// 二分查询最大的 right 使得切片 [left:right] 内的值满足 predicate
func (st *SegmentTree) MaxRight(left int32, predicate func(E) bool) int32 {
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
func (st *SegmentTree) MinLeft(right int32, predicate func(E) bool) int32 {
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
