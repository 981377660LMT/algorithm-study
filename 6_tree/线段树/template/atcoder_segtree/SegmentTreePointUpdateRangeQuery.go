package main

const INF int = 1e18

// PointUpdateRangeMax

type E = int

func (*SegmentTreePointUpdateRangeQuery) e() E        { return 0 }
func (*SegmentTreePointUpdateRangeQuery) op(a, b E) E { return max(a, b) }
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

type SegmentTreePointUpdateRangeQuery struct {
	n, size int
	data    []E
}

func NewSegmentTreePointUpdateRangeQuery(leaves []E) *SegmentTreePointUpdateRangeQuery {
	res := &SegmentTreePointUpdateRangeQuery{}
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
func (st *SegmentTreePointUpdateRangeQuery) Get(index int) E {
	if index < 0 || index >= st.n {
		return st.e()
	}
	return st.data[index+st.size]
}
func (st *SegmentTreePointUpdateRangeQuery) Set(index int, value E) {
	if index < 0 || index >= st.n {
		return
	}
	index += st.size
	st.data[index] = value
	for index >>= 1; index > 0; index >>= 1 {
		st.data[index] = st.op(st.data[index<<1], st.data[index<<1|1])
	}
}
func (st *SegmentTreePointUpdateRangeQuery) Update(index int, value E) {
	if index < 0 || index >= st.n {
		return
	}
	index += st.size
	st.data[index] = st.op(st.data[index], value)
	for index >>= 1; index > 0; index >>= 1 {
		st.data[index] = st.op(st.data[index<<1], st.data[index<<1|1])
	}
}

// [start, end)
func (st *SegmentTreePointUpdateRangeQuery) Query(start, end int) E {
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
func (st *SegmentTreePointUpdateRangeQuery) QueryAll() E { return st.data[1] }

// 二分查询最大的 right 使得切片 [left:right] 内的值满足 predicate
func (st *SegmentTreePointUpdateRangeQuery) MaxRight(left int, predicate func(E) bool) int {
	if left == st.n {
		return st.n
	}
	left += st.size
	res := st.e()
	for {
		for left&1 == 0 {
			left >>= 1
		}
		if !predicate(st.op(res, st.data[left])) {
			for left < st.size {
				left <<= 1
				if predicate(st.op(res, st.data[left])) {
					res = st.op(res, st.data[left])
					left++
				}
			}
			return left - st.size
		}
		res = st.op(res, st.data[left])
		left++
		if (left & -left) == left {
			break
		}
	}
	return st.n
}

// 二分查询最小的 left 使得切片 [left:right] 内的值满足 predicate
func (st *SegmentTreePointUpdateRangeQuery) MinLeft(right int, predicate func(E) bool) int {
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
		if !predicate(st.op(st.data[right], res)) {
			for right < st.size {
				right = right<<1 | 1
				if predicate(st.op(st.data[right], res)) {
					res = st.op(st.data[right], res)
					right--
				}
			}
			return right + 1 - st.size
		}
		res = st.op(st.data[right], res)
		if right&-right == right {
			break
		}
	}
	return 0
}
