package main

import (
	"sort"
)

// 区间元素相等查询.
type RangeEqualQueryStatic struct {
	preDiff []int32 // 前 i 个相邻对中的不相等数量
}

func NewRangeEqualQueryStatic(n int, isAdjacentEqual func(i, j /** i+1 **/ int) bool) *RangeEqualQueryStatic {
	if n <= 1 {
		return &RangeEqualQueryStatic{preDiff: []int32{0}}
	}
	preDiff := make([]int32, n)
	curDiffCount := int32(0)
	for i := 0; i < n-1; i++ {
		if !isAdjacentEqual(i, i+1) {
			curDiffCount++
		}
		preDiff[i+1] = curDiffCount
	}
	return &RangeEqualQueryStatic{preDiff: preDiff}
}

func (ash *RangeEqualQueryStatic) Query(start, end int) bool {
	if start < 0 {
		start = 0
	}
	if end > len(ash.preDiff) {
		end = len(ash.preDiff)
	}
	if end-start <= 1 {
		return true
	}
	return ash.preDiff[start] == ash.preDiff[end-1]
}

// 适用于连续段 m 较少的情况，节省空间.
type RangeEqualQueryStaticRLE struct {
	n             int
	segmentStarts []int
}

func NewRangeEqualQueryStaticRLE(n int, isAdjacentEqual func(i, j /** i+1 **/ int) bool) *RangeEqualQueryStaticRLE {
	if n == 0 {
		return &RangeEqualQueryStaticRLE{}
	}
	segmentStarts := []int{0}
	for i := 0; i < n-1; i++ {
		if !isAdjacentEqual(i, i+1) {
			segmentStarts = append(segmentStarts, i+1)
		}
	}
	return &RangeEqualQueryStaticRLE{n: n, segmentStarts: segmentStarts}
}

func (q *RangeEqualQueryStaticRLE) Query(start, end int) bool {
	if start < 0 {
		start = 0
	}
	if end > q.n {
		end = q.n
	}
	if end-start <= 1 {
		return true
	}
	idx1 := sort.SearchInts(q.segmentStarts, start+1) - 1
	idx2 := sort.SearchInts(q.segmentStarts, end) - 1
	return idx1 == idx2
}
