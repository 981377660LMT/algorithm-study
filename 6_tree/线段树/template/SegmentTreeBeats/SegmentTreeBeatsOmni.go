// SegmentTreeBeatsOmni.go
// 查询、修改操作时间复杂度均为 O(logn^2).
// api:
//  1. NewSegmentTreeBeatsOmni(n int32, f func(int32) int) *SegmentTreeBeatsOmni
//  2. (seg *SegmentTreeBeatsOmni) Build(f func(int32) int)
//  3. (seg *SegmentTreeBeatsOmni) UpdateMin(start, end int32, x int)
//  4. (seg *SegmentTreeBeatsOmni) UpdateMax(start, end int32, x int)
//	5. (seg *SegmentTreeBeatsOmni) Update(start, end int32, x int)
//  6. (seg *SegmentTreeBeatsOmni) QueryMin(start, end int32) int
//  7. (seg *SegmentTreeBeatsOmni) QueryMax(start, end int32) int
//  8. (seg *SegmentTreeBeatsOmni) QuerySum(start, end int32) int
//  9. (seg *SegmentTreeBeatsOmni) Enumerate(f func(*Node))

package main

import (
	"fmt"
	"runtime/debug"
	"time"
)

func init() {
	debug.SetGCPercent(-1)
}

func main() {
	seg := NewSegmentTreeBeatsOmni(10, func(i int32) int { return int(i) })
	fmt.Println(seg.QuerySum(0, 10), seg.QueryMax(0, 10))
	seg.UpdateMin(0, 10, 1)
	fmt.Println(seg.QuerySum(0, 10), seg.QueryMax(0, 10))
	seg.Build(func(i int32) int { return int(i * i) })
	fmt.Println(seg.QuerySum(0, 3), seg.QueryMax(0, 3), seg.QueryMin(0, 3))
	seg.Update(0, 3, 1)
	fmt.Println(seg.QuerySum(0, 3), seg.QueryMax(0, 3), seg.QueryMin(0, 3))
	seg.UpdateMax(0, 3, 10)
	fmt.Println(seg.QuerySum(0, 3), seg.QueryMax(0, 3), seg.QueryMin(0, 3))
	seg.Enumerate(func(node *Node) { fmt.Println(node.firstLarge, node.secondLarge, node.firstLargeCount, node.sum) })

	time1 := time.Now()
	seg = NewSegmentTreeBeatsOmni(int32(1e5), func(i int32) int { return int(i) })
	fmt.Println(time.Since(time1)) // 15ms

	time1 = time.Now()
	for i := 0; i < 1e5; i++ {
		seg.Update(0, int32(i), 1)
		seg.UpdateMax(0, int32(i), 1)
		seg.UpdateMin(0, int32(i), 1)
		seg.QuerySum(0, int32(i))
		seg.QueryMax(0, int32(i))
		seg.QueryMin(0, int32(i))
	}
	fmt.Println(time.Since(time1)) // 110ms
}

const INF int = 2e18

type Node struct {
	firstLargeCount, fisrtSmallCount int32
	size                             int32
	firstLarge, secondLarge          int
	firstSmall, secondSmall          int
	sum                              int
	dirty                            int
	left, right                      *Node
}

type SegmentTreeBeatsOmni struct {
	root *Node
	n    int32
}

func NewSegmentTreeBeatsOmni(n int32, f func(int32) int) *SegmentTreeBeatsOmni {
	res := &SegmentTreeBeatsOmni{n: n}
	res.Build(f)
	return res
}

func (seg *SegmentTreeBeatsOmni) Build(f func(int32) int) {
	seg.root = seg._build(seg.root, f, 0, seg.n-1)
}

func (seg *SegmentTreeBeatsOmni) UpdateMin(start, end int32, x int) {
	if start >= end {
		return
	}
	end--
	seg._updateMin(seg.root, start, end, 0, seg.n-1, x)
}

func (seg *SegmentTreeBeatsOmni) UpdateMax(start, end int32, x int) {
	if start >= end {
		return
	}
	end--
	seg._updateMax(seg.root, start, end, 0, seg.n-1, x)
}

func (seg *SegmentTreeBeatsOmni) Update(start, end int32, x int) {
	if start >= end {
		return
	}
	end--
	seg._update(seg.root, start, end, 0, seg.n-1, x)
}

func (seg *SegmentTreeBeatsOmni) QuerySum(start, end int32) int {
	if start >= end {
		return 0
	}
	end--
	return seg._querySum(seg.root, start, end, 0, seg.n-1)
}

func (seg *SegmentTreeBeatsOmni) QueryMin(start, end int32) int {
	if start >= end {
		return INF
	}
	end--
	return seg._queryMin(seg.root, start, end, 0, seg.n-1)
}

func (seg *SegmentTreeBeatsOmni) QueryMax(start, end int32) int {
	if start >= end {
		return -INF
	}
	end--
	return seg._queryMax(seg.root, start, end, 0, seg.n-1)
}

func (seg *SegmentTreeBeatsOmni) Enumerate(f func(*Node)) {
	seg._enumerate(seg.root, f)
}

func (seg *SegmentTreeBeatsOmni) _build(node *Node, f func(int32) int, l, r int32) *Node {
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
		node.firstLarge = node.sum
		node.firstSmall = node.sum
		node.firstLargeCount = 1
		node.fisrtSmallCount = 1
		node.secondLarge = -INF
		node.secondSmall = INF
		node.size = 1
	}
	return node
}

func (seg *SegmentTreeBeatsOmni) _setMin(node *Node, x int) {
	if node.firstLarge <= x {
		return
	}
	node.sum -= (node.firstLarge - x) * int(node.firstLargeCount)
	node.firstLarge = x
	if node.firstSmall >= x {
		node.firstSmall = x
	}
	node.secondSmall = min(node.secondSmall, x)
	if node.secondSmall == node.firstSmall {
		node.secondSmall = INF
	}
}

func (seg *SegmentTreeBeatsOmni) _setMax(node *Node, x int) {
	if node.firstSmall >= x {
		return
	}
	node.sum += (x - node.firstSmall) * int(node.fisrtSmallCount)
	node.firstSmall = x
	if node.firstLarge <= x {
		node.firstLarge = x
	}
	node.secondLarge = max(node.secondLarge, x)
	if node.secondLarge == node.firstLarge {
		node.secondLarge = -INF
	}
}

func (seg *SegmentTreeBeatsOmni) _propagate(node *Node, x int) {
	node.dirty += x
	node.sum += x * int(node.size)
	node.firstSmall += x
	node.firstLarge += x
	node.secondSmall += x
	node.secondLarge += x
}

func (seg *SegmentTreeBeatsOmni) _pushDown(node *Node) {
	if node.dirty != 0 {
		seg._propagate(node.left, node.dirty)
		seg._propagate(node.right, node.dirty)
		node.dirty = 0
	}
	seg._setMin(node.left, node.firstLarge)
	seg._setMin(node.right, node.firstLarge)
	seg._setMax(node.left, node.firstSmall)
	seg._setMax(node.right, node.firstSmall)
}

func (seg *SegmentTreeBeatsOmni) _pushUp(node *Node) {
	node.firstLarge = max(node.left.firstLarge, node.right.firstLarge)
	tmp1 := node.left.firstLarge
	if tmp1 == node.firstLarge {
		tmp1 = node.left.secondLarge
	}
	tmp2 := node.right.firstLarge
	if tmp2 == node.firstLarge {
		tmp2 = node.right.secondLarge
	}
	node.secondLarge = max(tmp1, tmp2)
	tmp3 := int32(0)
	if node.left.firstLarge == node.firstLarge {
		tmp3 = node.left.firstLargeCount
	}
	tmp4 := int32(0)
	if node.right.firstLarge == node.firstLarge {
		tmp4 = node.right.firstLargeCount
	}
	node.firstLargeCount = tmp3 + tmp4

	node.firstSmall = min(node.left.firstSmall, node.right.firstSmall)
	tmp5 := node.left.firstSmall
	if tmp5 == node.firstSmall {
		tmp5 = node.left.secondSmall
	}
	tmp6 := node.right.firstSmall
	if tmp6 == node.firstSmall {
		tmp6 = node.right.secondSmall
	}
	node.secondSmall = min(tmp5, tmp6)
	tmp7 := int32(0)
	if node.left.firstSmall == node.firstSmall {
		tmp7 = node.left.fisrtSmallCount
	}
	tmp8 := int32(0)
	if node.right.firstSmall == node.firstSmall {
		tmp8 = node.right.fisrtSmallCount
	}
	node.fisrtSmallCount = tmp7 + tmp8

	node.sum = node.left.sum + node.right.sum
	node.size = node.left.size + node.right.size
}

func (seg *SegmentTreeBeatsOmni) _updateMin(node *Node, ll, rr, l, r int32, x int) {
	if _noIntersection(ll, rr, l, r) {
		return
	}
	if _covered(ll, rr, l, r) {
		if node.firstLarge <= x {
			return
		}
		if node.secondLarge < x {
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

func (seg *SegmentTreeBeatsOmni) _updateMax(node *Node, ll, rr, l, r int32, x int) {
	if _noIntersection(ll, rr, l, r) {
		return
	}
	if _covered(ll, rr, l, r) {
		if node.firstSmall >= x {
			return
		}
		if node.secondSmall > x {
			seg._setMax(node, x)
			return
		}
	}
	seg._pushDown(node)
	mid := (l + r) >> 1
	seg._updateMax(node.left, ll, rr, l, mid, x)
	seg._updateMax(node.right, ll, rr, mid+1, r, x)
	seg._pushUp(node)
}

func (seg *SegmentTreeBeatsOmni) _update(node *Node, ll, rr, l, r int32, x int) {
	if _noIntersection(ll, rr, l, r) {
		return
	}
	if _covered(ll, rr, l, r) {
		seg._propagate(node, x)
		return
	}
	seg._pushDown(node)
	mid := (l + r) >> 1
	seg._update(node.left, ll, rr, l, mid, x)
	seg._update(node.right, ll, rr, mid+1, r, x)
	seg._pushUp(node)
}

func (seg *SegmentTreeBeatsOmni) _queryMax(node *Node, ll, rr, l, r int32) int {
	if _noIntersection(ll, rr, l, r) {
		return -INF
	}
	if _covered(ll, rr, l, r) {
		return node.firstLarge
	}
	seg._pushDown(node)
	mid := (l + r) >> 1
	return max(seg._queryMax(node.left, ll, rr, l, mid), seg._queryMax(node.right, ll, rr, mid+1, r))
}

func (seg *SegmentTreeBeatsOmni) _queryMin(node *Node, ll, rr, l, r int32) int {
	if _noIntersection(ll, rr, l, r) {
		return INF
	}
	if _covered(ll, rr, l, r) {
		return node.firstSmall
	}
	seg._pushDown(node)
	mid := (l + r) >> 1
	return min(seg._queryMin(node.left, ll, rr, l, mid), seg._queryMin(node.right, ll, rr, mid+1, r))
}

func (seg *SegmentTreeBeatsOmni) _querySum(node *Node, ll, rr, l, r int32) int {
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

func (seg *SegmentTreeBeatsOmni) _enumerate(node *Node, f func(*Node)) {
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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
