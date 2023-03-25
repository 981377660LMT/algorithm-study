// 无区间更新的线段树

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// https://yukicoder.me/submissions/845112
	// op1:赋值
	// !op2:查询区间最小值处的索引(线段树上二分即可)

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	nums := make([]E, n)
	for i := range nums {
		fmt.Fscan(in, &nums[i])
	}

	seg := NewSegmentTree(nums)
	for i := 0; i < q; i++ {
		var op, left, right int
		fmt.Fscan(in, &op, &left, &right)
		if op == 1 {
			left--
			right--
			tmp := seg.Get(left)
			seg.Set(left, seg.Get(right))
			seg.Set(right, tmp)
		} else {
			left--
			min_ := seg.Query(left, right)
			res := seg.MaxRight(left, func(x int) bool { return x > min_ })
			fmt.Fprintln(out, res+1)
		}
	}
}

type E = int

const INF int = 1e18

func (*SegmentTree) e() E        { return INF }
func (*SegmentTree) op(a, b E) E { return min(a, b) }

type SegmentTree struct {
	n, log, size int
	seg          []E
}

func NewSegmentTree(leaves []E) *SegmentTree {
	res := &SegmentTree{}
	n := len(leaves)
	log := 1
	for 1<<log < n {
		log++
	}
	size := 1 << log
	seg := make([]E, 2*size)
	for i := 0; i < n; i++ {
		seg[i+size] = leaves[i]
	}
	for i := size - 1; i > 0; i-- {
		seg[i] = res.op(seg[2*i], seg[2*i+1])
	}
	res.n = n
	res.log = log
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
		st.seg[index] = st.op(st.seg[2*index], st.seg[2*index+1])
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

// maxRight returns the maximum r such that [start, r) satisfies the predicate.
func (st *SegmentTree) MaxRight(start int, predicate func(E) bool) int {
	if start == st.n {
		return st.n
	}

	start += st.size
	res := st.e()
	for {
		for start&1 == 0 {
			start >>= 1
		}
		if !predicate(st.op(res, st.seg[start])) {
			for start < st.size {
				start = 2 * start
				if predicate(st.op(res, st.seg[start])) {
					res = st.op(res, st.seg[start])
					start++
				}
			}

			return start - st.size
		}
		res = st.op(res, st.seg[start])
		start++
		if (start & -start) == start {
			break
		}
	}
	return st.n
}

// minLeft returns the minimum l such that [l, end) satisfies the predicate.
func (st *SegmentTree) MinLeft(end int, predicate func(E) bool) int {
	if end == 0 {
		return 0
	}
	end += st.size
	sm := st.e()
	for {
		end--
		for end > 1 && end&1 == 1 {
			end >>= 1
		}
		if !predicate(st.op(st.seg[end], sm)) {
			for end < st.size {
				end = 2*end + 1
				if predicate(st.op(st.seg[end], sm)) {
					sm = st.op(st.seg[end], sm)
					end--
				}
			}
			return end + 1 - st.size
		}
		sm = st.op(st.seg[end], sm)
		if end&-end == end {
			break
		}
	}
	return 0
}

// !如果 Monoid 满足交换律(commute), 可以求出 op(nums[i xor x]...) (l<=i<r) 的值
//  下标异或上xor的区间查询.
func (st *SegmentTree) XorQuery(start, end, indexXor int) E {
	x := st.e()
	for k := 0; k < st.log+1; k++ {
		if start >= end {
			break
		}
		if start&1 == 1 {
			x = st.op(x, st.seg[(st.size>>k)+((start)^indexXor)])
			start++
		}
		if end&1 == 1 {
			end--
			x = st.op(x, st.seg[(st.size>>k)+((end)^indexXor)])
		}
		start, end, indexXor = start/2, end/2, indexXor/2
	}
	return x
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
