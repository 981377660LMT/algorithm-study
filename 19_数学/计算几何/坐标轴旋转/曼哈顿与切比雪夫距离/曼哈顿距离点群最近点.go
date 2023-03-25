// 曼哈顿距离点群搜索(NNS, nearest neighbor search)

package main

import (
	"fmt"
	"math/rand"
	"sort"
)

func main() {

	n := rand.Intn(200) + 1
	m := rand.Intn(200) + 1
	from := make([][2]int, n)
	to := make([][2]int, m)
	for i := 0; i < n; i++ {
		from[i] = [2]int{rand.Intn(1000) - 500, rand.Intn(1000) - 500}
	}
	for i := 0; i < m; i++ {
		to[i] = [2]int{rand.Intn(1000) - 500, rand.Intn(1000) - 500}
	}

	dist, indexes := ManhattanNNS(from, to)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			x1, y1 := from[i][0], from[i][1]
			x2, y2 := to[j][0], to[j][1]
			d := abs(x1-x2) + abs(y1-y2)
			if dist[i] > d {
				fmt.Println(from[i], to[j], dist[i], d)
				panic(fmt.Sprintf("dist[i] > d, i=%d, j=%d, dist[i]=%d, d=%d", i, j, dist[i], d))
			}
			if j == indexes[i] && dist[i] != d {
				panic("j == indexes[i] && dist[i] != d")
			}
		}
	}
	fmt.Println("OK")
}

const INF int = 1e18

// 给定点集 from 和 to，求 from 中的每个点到 to 中的点的曼哈顿距离的最小值, 以及最小值对应的索引.
func ManhattanNNS(from [][2]int, to [][2]int) (dist []int, indexes []int) {
	if len(to) == 0 {
		return
	}

	n, m := len(from), len(to)
	points := make([][2]int, n+m)
	for i, p := range from {
		points[i] = p
	}
	for i, p := range to {
		points[n+i] = p
	}
	ys := make([]int, m)
	for i, p := range to {
		ys[i] = p[1]
	}

	var getRank func(int) int
	ys, getRank = sortedSet(ys)

	indexes = make([]int, n)
	dist = make([]int, n)
	for i := 0; i < n; i++ {
		indexes[i] = -1
		dist[i] = INF
	}

	addResult := func(i, j int) {
		if j == -1 {
			return
		}
		dx := points[i][0] - points[j][0]
		dy := points[i][1] - points[j][1]
		cand := abs(dx) + abs(dy)
		if dist[i] > cand {
			dist[i] = cand
			indexes[i] = j - n
		}
	}

	order := make([]int, n+m)
	for i := range order {
		order[i] = i
	}

	sort.Slice(order, func(i, j int) bool {
		return points[order[i]][0] < points[order[j]][0]
	})

	leaves := make([]E, len(ys))
	for i := 0; i < len(ys); i++ {
		leaves[i] = E{min_: INF, index: -1}
	}
	cal := func() {
		seg1, seg2 := NewSegmentTree(leaves), NewSegmentTree(leaves)
		for _, i := range order {
			x, y := points[i][0], points[i][1]
			idx := getRank(y)
			if i < n {
				addResult(i, seg1.Query(idx, len(ys)).index)
				addResult(i, seg2.Query(0, idx).index)
			} else {
				seg1.Set(idx, E{y - x, i})
				seg2.Set(idx, E{-(x + y), i})
			}
		}
	}

	cal()
	for i, j := 0, len(order)-1; i < j; i, j = i+1, j-1 {
		order[i], order[j] = order[j], order[i]
	}
	for i := 0; i < n+m; i++ {
		points[i][0] = -points[i][0]
	}
	cal()

	return
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func sortedSet(xs []int) (sorted []int, getRank func(int) int) {
	set := make(map[int]struct{}, len(xs))
	for _, v := range xs {
		set[v] = struct{}{}
	}
	sorted = make([]int, 0, len(set))
	for k := range set {
		sorted = append(sorted, k)
	}
	sort.Ints(sorted)
	getRank = func(x int) int { return sort.SearchInts(sorted, x) }
	return
}

// RangeMinWithIndex

type E = struct{ min_, index int }

func (*SegmentTree) e() E { return E{min_: INF, index: -1} }
func (*SegmentTree) op(a, b E) E {
	if a.min_ < b.min_ {
		return a
	}
	if a.min_ > b.min_ {
		return b
	}
	return E{min_: a.min_, index: min(a.index, b.index)}
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
				if predicate(st.op(res, st.seg[left])) {
					res = st.op(res, st.seg[left])
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
				if predicate(st.op(st.seg[right], res)) {
					res = st.op(st.seg[right], res)
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
