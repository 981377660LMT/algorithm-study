package main

type RangeEqualNode[V comparable] struct {
	isEqual bool
	empty   bool
	lVal    V
	rVal    V
}

type RangeEqualQueryPointSet[V comparable] struct {
	st *SegmentTreeG[RangeEqualNode[V]]
}

func NewRangeEqualQueryPointSet[V comparable](arr []V) *RangeEqualQueryPointSet[V] {
	e := func() RangeEqualNode[V] {
		return RangeEqualNode[V]{isEqual: true, empty: true}
	}

	op := func(a, b RangeEqualNode[V]) RangeEqualNode[V] {
		if a.empty {
			return b
		}
		if b.empty {
			return a
		}
		return RangeEqualNode[V]{
			isEqual: a.isEqual && b.isEqual && (a.rVal == b.lVal),
			empty:   false,
			lVal:    a.lVal,
			rVal:    b.rVal,
		}
	}

	leaves := make([]RangeEqualNode[V], len(arr))
	for i, v := range arr {
		leaves[i] = RangeEqualNode[V]{
			isEqual: true,
			empty:   false,
			lVal:    v,
			rVal:    v,
		}
	}

	return &RangeEqualQueryPointSet[V]{
		st: NewSegmentTreeGFrom(leaves, op, e),
	}
}

func (req *RangeEqualQueryPointSet[V]) Set(index int, val V) {
	req.st.Set(index, RangeEqualNode[V]{
		isEqual: true,
		empty:   false,
		lVal:    val,
		rVal:    val,
	})
}

func (req *RangeEqualQueryPointSet[V]) Query(start, end int) bool {
	if start >= end-1 {
		return true
	}
	return req.st.Query(start, end).isEqual
}

type SegmentTreeG[E any] struct {
	n, size int
	seg     []E
	op      func(E, E) E
	e       func() E
}

func NewSegmentTreeGFrom[E any](leaves []E, op func(E, E) E, e func() E) *SegmentTreeG[E] {
	res := &SegmentTreeG[E]{
		op: op,
		e:  e,
	}
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

func (st *SegmentTreeG[E]) Get(index int) E {
	if index < 0 || index >= st.n {
		return st.e()
	}
	return st.seg[index+st.size]
}

func (st *SegmentTreeG[E]) Set(index int, value E) {
	if index < 0 || index >= st.n {
		return
	}
	index += st.size
	st.seg[index] = value
	for index >>= 1; index > 0; index >>= 1 {
		st.seg[index] = st.op(st.seg[index<<1], st.seg[index<<1|1])
	}
}

func (st *SegmentTreeG[E]) Update(index int, value E) {
	if index < 0 || index >= st.n {
		return
	}
	index += st.size
	st.seg[index] = st.op(st.seg[index], value)
	for index >>= 1; index > 0; index >>= 1 {
		st.seg[index] = st.op(st.seg[index<<1], st.seg[index<<1|1])
	}
}

func (st *SegmentTreeG[E]) Query(start, end int) E {
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
