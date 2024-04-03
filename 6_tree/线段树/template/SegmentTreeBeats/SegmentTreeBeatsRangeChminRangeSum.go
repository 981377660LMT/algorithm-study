// SegmentTreeBeatsRangeChminRangeSum.go
// 查询、修改操作时间复杂度均为 O(logn).
// api:
//  1. NewSegmentTreeBeatsRangeChminRangeSum(n int32, f func(int32) int) *SegmentTreeBeatsRangeChminRangeSum
//  2. (seg *SegmentTreeBeatsRangeChminRangeSum) Build(f func(int32) int)
//  3. (seg *SegmentTreeBeatsRangeChminRangeSum) UpdateMin(start, end int32, x int)
//  4. (seg *SegmentTreeBeatsRangeChminRangeSum) QuerySum(start, end int32) int
//  5. (seg *SegmentTreeBeatsRangeChminRangeSum) QueryMax(start, end int32) int
//  6. (seg *SegmentTreeBeatsRangeChminRangeSum) Enumerate(f func(*Node))

package main

import (
	"fmt"
	"runtime/debug"
)

func init() {
	debug.SetGCPercent(-1)
}

func main() {
	seg := NewSegmentTreeBeatsRangeChminRangeSum(10, func(i int32) int { return int(i) })
	fmt.Println(seg.QuerySum(0, 10), seg.QueryMax(0, 10))
	seg.UpdateMin(0, 10, 1)
	fmt.Println(seg.QuerySum(0, 10), seg.QueryMax(0, 10))
	seg.Enumerate(func(node *Node) { fmt.Println(node.first, node.second, node.firstCount, node.sum) })
	seg.Build(func(i int32) int { return int(i * i) })
	fmt.Println(seg.QuerySum(0, 3), seg.QueryMax(0, 3))
}

const INF int = 2e18

type Node struct {
	left, right   *Node
	first, second int
	firstCount    int32
	sum           int
}

type SegmentTreeBeatsRangeChminRangeSum struct {
	root *Node
	n    int32
}

func NewSegmentTreeBeatsRangeChminRangeSum(n int32, f func(int32) int) *SegmentTreeBeatsRangeChminRangeSum {
	res := &SegmentTreeBeatsRangeChminRangeSum{n: n}
	res.Build(f)
	return res
}

func (seg *SegmentTreeBeatsRangeChminRangeSum) Build(f func(int32) int) {
	seg.root = seg._build(seg.root, f, 0, seg.n-1)
}

func (seg *SegmentTreeBeatsRangeChminRangeSum) UpdateMin(start, end int32, x int) {
	if start >= end {
		return
	}
	end--
	seg._updateMin(seg.root, start, end, 0, seg.n-1, x)
}

func (seg *SegmentTreeBeatsRangeChminRangeSum) QuerySum(start, end int32) int {
	if start >= end {
		return 0
	}
	end--
	return seg._querySum(seg.root, start, end, 0, seg.n-1)
}

func (seg *SegmentTreeBeatsRangeChminRangeSum) QueryMax(start, end int32) int {
	if start >= end {
		return -INF
	}
	end--
	return seg._queryMax(seg.root, start, end, 0, seg.n-1)
}

func (seg *SegmentTreeBeatsRangeChminRangeSum) Enumerate(f func(*Node)) {
	seg._enumerate(seg.root, f)
}

func (seg *SegmentTreeBeatsRangeChminRangeSum) _build(node *Node, f func(int32) int, l, r int32) *Node {
	if node == nil {
		node = &Node{}
	}
	if l < r {
		mid := (l + r) >> 1
		node.left = seg._build(node.left, f, l, mid)
		node.right = seg._build(node.right, f, mid+1, r)
		seg._pushUp(node)
	} else {
		node.sum = f(l)
		node.first = node.sum
		node.second = -INF
		node.firstCount = 1
	}
	return node
}

func (seg *SegmentTreeBeatsRangeChminRangeSum) _setMin(node *Node, x int) {
	if node.first <= x {
		return
	}
	node.sum -= (node.first - x) * int(node.firstCount)
	node.first = x
}

func (seg *SegmentTreeBeatsRangeChminRangeSum) _pushDown(node *Node) {
	seg._setMin(node.left, node.first)
	seg._setMin(node.right, node.first)
}

func (seg *SegmentTreeBeatsRangeChminRangeSum) _pushUp(node *Node) {
	node.first = max(node.left.first, node.right.first)
	tmp1 := node.left.first
	if tmp1 == node.first {
		tmp1 = node.left.second
	}
	tmp2 := node.right.first
	if tmp2 == node.first {
		tmp2 = node.right.second
	}
	node.second = max(tmp1, tmp2)
	tmp3 := int32(0)
	if node.left.first == node.first {
		tmp3 = node.left.firstCount
	}
	tmp4 := int32(0)
	if node.right.first == node.first {
		tmp4 = node.right.firstCount
	}
	node.firstCount = tmp3 + tmp4
	node.sum = node.left.sum + node.right.sum
}

func (seg *SegmentTreeBeatsRangeChminRangeSum) _updateMin(node *Node, ll, rr, l, r int32, x int) {
	if _noIntersection(ll, rr, l, r) {
		return
	}
	if _covered(ll, rr, l, r) {
		if node.first <= x {
			return
		}
		if node.second < x {
			seg._setMin(node, x)
			return
		}
	}
	seg._pushDown(node)
	mid := (l + r) >> 1
	seg._updateMin(node.left, ll, rr, l, mid, x)
	seg._updateMin(node.right, ll, rr, mid+1, r, x)
	seg._pushUp(node)
}

func (seg *SegmentTreeBeatsRangeChminRangeSum) _queryMax(node *Node, ll, rr, l, r int32) int {
	if _noIntersection(ll, rr, l, r) {
		return -INF
	}
	if _covered(ll, rr, l, r) {
		return node.first
	}
	seg._pushDown(node)
	mid := (l + r) >> 1
	return max(seg._queryMax(node.left, ll, rr, l, mid), seg._queryMax(node.right, ll, rr, mid+1, r))
}

func (seg *SegmentTreeBeatsRangeChminRangeSum) _querySum(node *Node, ll, rr, l, r int32) int {
	if _noIntersection(ll, rr, l, r) {
		return 0
	}
	if _covered(ll, rr, l, r) {
		return node.sum
	}
	seg._pushDown(node)
	mid := (l + r) >> 1
	return seg._querySum(node.left, ll, rr, l, mid) + seg._querySum(node.right, ll, rr, mid+1, r)
}

func (seg *SegmentTreeBeatsRangeChminRangeSum) _enumerate(node *Node, f func(*Node)) {
	if node.left == nil && node.right == nil {
		f(node)
		return
	}
	seg._pushDown(node)
	seg._enumerate(node.left, f)
	seg._enumerate(node.right, f)
}

func _covered(ll, rr, l, r int32) bool {
	return ll <= l && rr >= r
}

func _noIntersection(ll, rr, l, r int32) bool {
	return ll > r || rr < l
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
