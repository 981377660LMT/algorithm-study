// IntervalNearestPair 区间绝对值差的最小值.
// 区间绝对值差的最小值.
// 区间距离最小的二元组.
// 离线解法.
//
// 时间复杂度: O(nlog_2nlog_2M+qlog_2n)
// 空间复杂度: O(n\log_2M)
// M 是 a 中的最大绝对值.

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	Cf765F()
}

// https://codeforces.com/problemset/submission/765/240821486
// 给定一个数组和q组查询，每组查询包含两个整数start,end，求出[start,end)区间内的 abs(a[i]-a[j])的最小值(i!=j)。
func Cf765F() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	var q int
	fmt.Fscan(in, &q)
	queries := make([]*Query, q)
	for i := 0; i < q; i++ {
		var l, r int32
		fmt.Fscan(in, &l, &r)
		l--
		r--
		queries[i] = NewQuery(l, r, int32(i))
	}

	RangeNearestPair(nums, queries)
	res := make([]int, q)
	for _, q := range queries {
		res[q.id] = q.Res
	}
	for _, v := range res {
		fmt.Fprintln(out, v)
	}
}

const INF int = 2e18

type Query struct {
	Res         int
	left, right int32
	id          int32
}

func NewQuery(left, right int32, id int32) *Query {
	return &Query{Res: INF, left: left, right: right, id: id}
}

func (q *Query) String() string {
	return fmt.Sprintf("[%d, %d]: %d", q.left, q.right, q.Res)
}

func RangeNearestPair(nums []int, queries []*Query) {
	nums = append(nums[:0:0], nums...)
	sort.Slice(queries, func(i, j int) bool { return queries[i].right < queries[j].right })
	_rangeNearestPair(nums, queries)
	const max int = 1e18
	for i := 0; i < len(nums); i++ {
		nums[i] = max - nums[i]
	}
	_rangeNearestPair(nums, queries)
}

func _rangeNearestPair(nums []int, queries []*Query) {
	n := int32(len(nums))
	ceil := make([][]int32, n)
	buf := make([]int32, 0, 60)
	indexes := make([]int32, n)
	for i := int32(0); i < n; i++ {
		indexes[i] = i
	}
	sort.Slice(indexes, func(i, j int) bool {
		if nums[indexes[i]] == nums[indexes[j]] {
			return indexes[i] < indexes[j]
		}
		return nums[indexes[i]] > nums[indexes[j]]
	})
	segBs := newSegmentBs(0, n-1)
	segBs.reset(0, n-1, INF)
	for _, i := range indexes {
		threshold := INF - 1
		j := segBs.query(0, i, 0, n-1, threshold)
		buf = buf[:0]
		for j != -1 {
			buf = append(buf, j)
			threshold = (nums[j]+nums[i]+1)/2 - 1
			j = segBs.query(0, j-1, 0, n-1, threshold)
		}
		ceil[i] = append(ceil[i], buf...)
		segBs.update(i, i, 0, n-1, nums[i])
	}
	segtree := newSegment(0, n-1)
	l := int32(0)
	for i := int32(0); i < n; i++ {
		for _, j := range ceil[i] {
			segtree.update(0, j, 0, n-1, nums[j]-nums[i])
		}
		for l < int32(len(queries)) && queries[l].right == i {
			queries[l].Res = min(queries[l].Res, segtree.query(queries[l].left, queries[l].left, 0, n-1))
			l++
		}
	}
}

type segment struct {
	left, right *segment
	min         int
}

func (s *segment) modify(x int) {
	s.min = min(s.min, x)
}

func (s *segment) pushUp() {}

func (s *segment) pushDown() {
	s.left.modify(s.min)
	s.right.modify(s.min)
	s.min = INF
}

func newSegment(l, r int32) *segment {
	if l < r {
		m := (l + r) >> 1
		res := &segment{left: newSegment(l, m), right: newSegment(m+1, r), min: INF}
		res.pushUp()
		return res
	}
	return &segment{min: INF}
}

func (s *segment) covered(ll, rr, l, r int32) bool {
	return ll <= l && rr >= r
}

func (s *segment) noIntersection(ll, rr, l, r int32) bool {
	return ll > r || rr < l
}

func (s *segment) update(ll, rr, l, r int32, x int) {
	if s.noIntersection(ll, rr, l, r) {
		return
	}
	if s.covered(ll, rr, l, r) {
		s.modify(x)
		return
	}
	s.pushDown()
	m := (l + r) >> 1
	s.left.update(ll, rr, l, m, x)
	s.right.update(ll, rr, m+1, r, x)
	s.pushUp()
}

func (s *segment) query(ll, rr, l, r int32) int {
	if s.noIntersection(ll, rr, l, r) {
		return INF
	}
	if s.covered(ll, rr, l, r) {
		return s.min
	}
	s.pushDown()
	m := (l + r) >> 1
	return min(s.left.query(ll, rr, l, m), s.right.query(ll, rr, m+1, r))
}

type segmentBs struct {
	left, right *segmentBs
	min         int
}

func (s *segmentBs) modify(x int) {
	s.min = x
}

func (s *segmentBs) pushUp() {
	s.min = min(s.left.min, s.right.min)
}

func (s *segmentBs) pushDown() {}

func newSegmentBs(l, r int32) *segmentBs {
	if l < r {
		m := (l + r) >> 1
		res := &segmentBs{left: newSegmentBs(l, m), right: newSegmentBs(m+1, r)}
		res.pushUp()
		return res
	}
	return &segmentBs{}
}

func (s *segmentBs) reset(l, r int32, x int) {
	if l < r {
		m := (l + r) >> 1
		s.left.reset(l, m, x)
		s.right.reset(m+1, r, x)
		s.pushUp()
	} else {
		s.min = x
	}
}

func (s *segmentBs) covered(ll, rr, l, r int32) bool {
	return ll <= l && rr >= r
}

func (s *segmentBs) noIntersection(ll, rr, l, r int32) bool {
	return ll > r || rr < l
}

func (s *segmentBs) update(ll, rr, l, r int32, x int) {
	if s.noIntersection(ll, rr, l, r) {
		return
	}
	if s.covered(ll, rr, l, r) {
		s.modify(x)
		return
	}
	s.pushDown()
	m := (l + r) >> 1
	s.left.update(ll, rr, l, m, x)
	s.right.update(ll, rr, m+1, r, x)
	s.pushUp()
}

func (s *segmentBs) query(ll, rr, l, r int32, threshold int) int32 {
	if s.noIntersection(ll, rr, l, r) || s.min > threshold {
		return -1
	}
	if s.covered(ll, rr, l, r) && l == r {
		return l
	}
	s.pushDown()
	m := (l + r) >> 1
	res := s.right.query(ll, rr, m+1, r, threshold)
	if res == -1 {
		res = s.left.query(ll, rr, l, m, threshold)
	}
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
