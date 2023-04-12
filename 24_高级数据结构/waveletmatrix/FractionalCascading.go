// 分数级联
// 查询区间 [start, end) 中，[floor, ceiling) 范围内的数的个数.

package main

import (
	"fmt"
	"sort"
)

const INF int = 1e18

func countQuadruplets(nums []int) int64 {
	W := NewSegmentTreeFractionalCascading(nums)
	res := 0
	for j := 1; j < len(nums)-2; j++ {
		for k := j + 1; k < len(nums)-1; k++ {
			if !(nums[k] < nums[j]) {
				continue
			}
			left := W.Query(0, j, 0, nums[k])
			right := W.Query(k+1, len(nums), nums[j]+1, INF)
			res += (left * right)
		}
	}
	return int64(res)
}

func main() {
	tree := NewSegmentTreeFractionalCascading([]int{3, 2, 5, 4, 6, 7, 8})
	fmt.Println(tree.Query(0, 1, 2, 3))
}

type SegmentTreeFractionalCascading struct {
	seg, ll, rr [][]int
	sz          int
}

func NewSegmentTreeFractionalCascading(array []int) *SegmentTreeFractionalCascading {
	res := &SegmentTreeFractionalCascading{}
	n := len(array)
	sz := 1
	for sz < n {
		sz <<= 1
	}
	tmp := 2*sz - 1
	seg := make([][]int, tmp)
	ll := make([][]int, tmp)
	rr := make([][]int, tmp)
	for k := 0; k < n; k++ {
		seg[k+sz-1] = append(seg[k+sz-1], array[k])
	}
	for k := sz - 2; k >= 0; k-- {
		a, b := 2*k+1, 2*k+2
		len_ := len(seg[a]) + len(seg[b])
		ll[k] = make([]int, len_+1)
		rr[k] = make([]int, len_+1)
		seg[k] = make([]int, 0, len_)
		i, j := 0, 0
		for i < len(seg[a]) && j < len(seg[b]) {
			if seg[a][i] < seg[b][j] {
				seg[k] = append(seg[k], seg[a][i])
				i++
			} else {
				seg[k] = append(seg[k], seg[b][j])
				j++
			}
		}
		seg[k] = append(seg[k], seg[a][i:]...)
		seg[k] = append(seg[k], seg[b][j:]...)
		tail1, tail2 := 0, 0
		for i := 0; i < len(seg[k]); i++ {
			for tail1 < len(seg[a]) && seg[a][tail1] < seg[k][i] {
				tail1++
			}
			for tail2 < len(seg[b]) && seg[b][tail2] < seg[k][i] {
				tail2++
			}
			ll[k][i] = tail1
			rr[k][i] = tail2
		}
		ll[k][len(seg[k])] = len(seg[a])
		rr[k][len(seg[k])] = len(seg[b])
	}
	res.seg = seg
	res.ll = ll
	res.rr = rr
	res.sz = sz
	return res
}

// 查询区间 [start, end) 中，[floor, ceiling) 范围内的数的个数.
func (st *SegmentTreeFractionalCascading) Query(start, end, floor, ceiling int) int {
	floor = sort.SearchInts(st.seg[0], floor)
	ceiling = sort.SearchInts(st.seg[0], ceiling)
	return st._query(start, end, floor, ceiling, 0, 0, st.sz)
}

func (st *SegmentTreeFractionalCascading) _query(a, b, lower, upper, k, l, r int) int {
	if a >= r || b <= l {
		return 0
	}
	if a <= l && r <= b {
		return upper - lower
	}
	return st._query(a, b, st.ll[k][lower], st.ll[k][upper], 2*k+1, l, (l+r)>>1) + st._query(a, b, st.rr[k][lower], st.rr[k][upper], 2*k+2, (l+r)>>1, r)
}
