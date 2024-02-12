package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// 不带修区间 Closest Equals 线段树+离线查询
// https://www.luogu.com.cn/problem/CF522D
// 给出一个长度为n的数组,有q个查询,每个查询包含两个整数l和r,求[l,r]中最近的两个相等的数的距离.
// 不存在的话输出-1.
// !1<=n,q<=5e5
//
// https://www.cnblogs.com/hypersq/articles/15706076.html
// last数组记录上一次这个数出现的位置,pre数组记录上一个位置的这个数的位置.
// !i与pre[i]的贡献记录在pre[i]的位置,那么这个区间查询的答案就是[l,r]中记录的贡献的最小值
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
	queries := make([][2]int, q)
	for i := 0; i < q; i++ {
		fmt.Fscan(in, &queries[i][0], &queries[i][1])
		queries[i][0]--
	}

	D := NewDictionary()
	for i := 0; i < n; i++ {
		nums[i] = D.Id(nums[i])
	}

	pre := make([]int, n) // 上一个nums[i]的位置
	for i := range pre {
		pre[i] = -1
	}
	last := make([]int, D.Size()) // 上一次这个数出现的位置
	for i := range last {
		last[i] = -1
	}
	for i, v := range nums {
		if tmp := last[v]; tmp != -1 {
			pre[i] = tmp
		}
		last[v] = i
	}

	order := make([]int, q)
	for i := range order {
		order[i] = i
	}
	sort.Slice(order, func(i, j int) bool { // 按右端点排序
		return queries[order[i]][1] < queries[order[j]][1]
	})

	seg := NewSegmentTree(n, func(i int) int { return INF })
	right := 0
	res := make([]int, q)
	for _, qid := range order {
		start, end := queries[qid][0], queries[qid][1]
		for right < end {
			if pre[right] != -1 {
				seg.Update(pre[right], right-pre[right])
			}
			right++
		}
		res[qid] = seg.Query(start, end)
	}

	for _, v := range res {
		if v == INF {
			v = -1
		}
		fmt.Fprintln(out, v)
	}
}

const INF int = 1e18

// PointSetRangeMin

type E = int

func (*SegmentTree) e() E        { return INF }
func (*SegmentTree) op(a, b E) E { return min(a, b) }
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

type V = int
type Dictionary struct {
	_idToValue []V
	_valueToId map[V]int
}

// A dictionary that maps values to unique ids.
func NewDictionary() *Dictionary {
	return &Dictionary{
		_valueToId: map[V]int{},
	}
}
func (d *Dictionary) Id(value V) int {
	res, ok := d._valueToId[value]
	if ok {
		return res
	}
	id := len(d._idToValue)
	d._idToValue = append(d._idToValue, value)
	d._valueToId[value] = id
	return id
}
func (d *Dictionary) Value(id int) V {
	return d._idToValue[id]
}
func (d *Dictionary) Has(value V) bool {
	_, ok := d._valueToId[value]
	return ok
}
func (d *Dictionary) Size() int {
	return len(d._idToValue)
}
