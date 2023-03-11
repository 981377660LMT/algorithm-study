package main

func smallestMissingValueSubtree(parents []int, nums []int) []int {
	n := len(nums)
	adjList := make([][]int, n)
	for i := 1; i < n; i++ {
		p := parents[i]
		adjList[p] = append(adjList[p], i)
	}

	ins, outs, dfsId := make([]int, n+5), make([]int, n+5), 1
	depth := make([]int, n+5)
	var dfsOrder func(cur, pre int)
	dfsOrder = func(cur, pre int) {
		ins[cur] = dfsId
		for _, next := range adjList[cur] {
			if next != pre {
				depth[next] = depth[cur] + 1
				dfsOrder(next, cur)
			}
		}
		outs[cur] = dfsId
		dfsId += 1
	}
	dfsOrder(0, -1)
	newNums := make([]int, n)
	for i := 0; i < n; i++ {
		id := outs[i] - 1
		newNums[id] = nums[i]
	}

	rmq := NewRangeMexQuery(newNums)
	for i := 0; i < n; i++ {
		rmq.AddQuery(ins[i]-1, outs[i])
	}

	return rmq.Run(1)
}

const INF int = 1e18

type E = int

func e() E        { return INF }
func op(a, b E) E { return min(a, b) }
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

type RangeMexQuery struct {
	nums  []int
	query [][2]int
}

func NewRangeMexQuery(nums []int) *RangeMexQuery {
	return &RangeMexQuery{nums: nums}
}

// [start, end)
//  0 <= start <= end <= n
func (rmq *RangeMexQuery) AddQuery(start, end int) {
	rmq.query = append(rmq.query, [2]int{start, end})
}

// mexStart: mex的起始值(从0开始还是从1开始)
func (rmq *RangeMexQuery) Run(mexStart int) []int {
	n := len(rmq.nums)
	leaves := make([]E, n+2)
	for i := 0; i < n+2; i++ {
		leaves[i] = -1
	}
	seg := NewSegmentTree(leaves)

	q := len(rmq.query)
	res := make([]int, q)
	ids := make([][]int, n+1)
	for i := 0; i < q; i++ {
		end := rmq.query[i][1]
		ids[end] = append(ids[end], i)
	}

	for i := 0; i < n+1; i++ {
		for _, q := range ids[i] {
			start := rmq.query[q][0]
			mex := seg.MaxRight(mexStart, func(x int) bool { return x >= start })
			res[q] = mex
		}
		if i < n && rmq.nums[i] < n+2 {
			seg.Set(rmq.nums[i], i)
		}
	}

	return res
}

type SegmentTree struct {
	n, size int
	seg     []E
}

func NewSegmentTree(leaves []E) *SegmentTree {
	res := &SegmentTree{}
	n := len(leaves)
	size := 1
	for size < n {
		size <<= 1
	}
	seg := make([]E, size<<1)
	for i := 0; i < n; i++ {
		seg[i+size] = leaves[i]
	}
	for i := size - 1; i > 0; i-- {
		seg[i] = op(seg[i<<1], seg[i<<1|1])
	}
	res.n = n
	res.size = size
	res.seg = seg
	return res
}

func (st *SegmentTree) Get(index int) E {
	if index < 0 || index >= st.n {
		return e()
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
		st.seg[index] = op(st.seg[index<<1], st.seg[index<<1|1])
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
		return e()
	}
	leftRes, rightRes := e(), e()
	start += st.size
	end += st.size
	for start < end {
		if start&1 == 1 {
			leftRes = op(leftRes, st.seg[start])
			start++
		}
		if end&1 == 1 {
			end--
			rightRes = op(st.seg[end], rightRes)
		}
		start >>= 1
		end >>= 1
	}
	return op(leftRes, rightRes)
}

func (st *SegmentTree) QueryAll() E { return st.seg[1] }

// 二分查询最大的 right 使得切片 [left:right] 内的值满足 predicate
func (st *SegmentTree) MaxRight(left int, predicate func(E) bool) int {
	if left == st.n {
		return st.n
	}
	left += st.size
	res := e()
	for {
		for left&1 == 0 {
			left >>= 1
		}
		if !predicate(op(res, st.seg[left])) {
			for left < st.size {
				left <<= 1
				if predicate(op(res, st.seg[left])) {
					res = op(res, st.seg[left])
					left++
				}
			}
			return left - st.size
		}
		res = op(res, st.seg[left])
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
	res := e()
	for {
		right--
		for right > 1 && right&1 == 1 {
			right >>= 1
		}
		if !predicate(op(st.seg[right], res)) {
			for right < st.size {
				right = right<<1 | 1
				if predicate(op(st.seg[right], res)) {
					res = op(st.seg[right], res)
					right--
				}
			}
			return right + 1 - st.size
		}
		res = op(st.seg[right], res)
		if right&-right == right {
			break
		}
	}
	return 0
}
