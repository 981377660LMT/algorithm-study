package main

import (
	"fmt"
	"strings"
)

// queries = [[1,7],[2,7,6],[1,2],[2,7,5],[2,7,6]]
func main() {
	queries := [][]int{{1, 7}, {2, 7, 6}, {1, 2}, {2, 7, 5}, {2, 7, 6}}
	fmt.Println(getResults(queries))
	// queries = [[1,2],[2,3,3],[2,3,1],[2,2,2]]
	queries = [][]int{{1, 2}, {2, 3, 3}, {2, 3, 1}, {2, 2, 2}}
	fmt.Println(getResults(queries))
	// [[2,1,2]]
	queries = [][]int{{2, 1, 2}}
	fmt.Println(getResults(queries))
	// [[2,5,1],[1,3],[2,6,4]]
	queries = [][]int{{2, 5, 1}, {1, 3}, {2, 6, 4}}
	fmt.Println(getResults(queries)) // [ture, false]
}

// 有一条无限长的数轴，原点在 0 处，沿着 x 轴 正 方向无限延伸。

// 给你一个二维数组 queries ，它包含两种操作：

// 操作类型 1 ：queries[i] = [1, x] 。在距离原点 x 处建一个障碍物。数据保证当操作执行的时候，位置 x 处 没有 任何障碍物。
// 操作类型 2 ：queries[i] = [2, x, sz] 。判断在数轴范围 [0, x] 内是否可以放置一个长度为 sz 的物块，这个物块需要 完全 放置在范围 [0, x] 内。如果物块与任何障碍物有重合，那么这个物块 不能 被放置，但物块可以与障碍物刚好接触。注意，你只是进行查询，并 不是 真的放置这个物块。每个查询都是相互独立的。
// 请你返回一个 boolean 数组results ，如果第 i 个操作类型 2 的操作你可以放置物块，那么 results[i] 为 true ，否则为 false 。
func getResults(queries [][]int) []bool {
	q := len(queries)
	res := make([]bool, 0, q)

	leaves := make([]E, 3*q+10)
	for i := range leaves {
		leaves[i] = FromElement(1)
	}

	seg := NewSegmentTreeLongestOne(leaves)
	for _, query := range queries {
		if query[0] == 1 {
			seg.Set(query[1], FromElement(0))
		} else {
			if query[1] < query[2] {
				res = append(res, false)
				continue
			}
			res = append(res, seg.Query(1, query[1]).longestOne+1 >= query[2])
		}
	}
	return res
}

type V = int

type E = struct {
	size                       int
	preOne, sufOne, longestOne int // 前缀1的个数, 后缀1的个数, 最长1的个数
	leftV, rightV              V   // 左端点值, 右端点值

	// !TODO
	pairCount int // 区间内所有极长连续1段的贡献和 sum(len_i*(len_i+1)/2)
}

func FromElement(v V) E {
	if v == 1 {
		return E{
			size: 1, preOne: 1, sufOne: 1, longestOne: 1, leftV: 1, rightV: 1,
			// !TODO
			pairCount: 1,
		}
	}

	return E{
		size: 1, leftV: v, rightV: v,
		// !TODO
		pairCount: 0,
	}
}

func (*SegmentTreeLongestOne) e() E { return E{} }
func (*SegmentTreeLongestOne) op(a, b E) E {
	res := E{leftV: a.leftV, rightV: b.rightV, size: a.size + b.size}
	if a.rightV == b.leftV {
		res.preOne = a.preOne
		if a.preOne == a.size {
			res.preOne += b.preOne
		}
		res.sufOne = b.sufOne
		if b.sufOne == b.size {
			res.sufOne += a.sufOne
		}
		res.longestOne = max(max(a.longestOne, b.longestOne), a.sufOne+b.preOne)

		// !TODO
		n1, n2 := a.sufOne, b.preOne
		n3 := n1 + n2
		res.pairCount = a.pairCount + b.pairCount + n3*(n3+1)/2 - n1*(n1+1)/2 - n2*(n2+1)/2
	} else {
		res.preOne = a.preOne
		res.sufOne = b.sufOne
		res.longestOne = max(a.longestOne, b.longestOne)

		// !TODO
		res.pairCount = a.pairCount + b.pairCount
	}
	return res
}

// 维护区间最长1的个数, 区间前缀1的个数，区间后缀1的个数.
type SegmentTreeLongestOne struct {
	n, size int
	seg     []E
}

func NewSegmentTreeLongestOne(leaves []E) *SegmentTreeLongestOne {
	res := &SegmentTreeLongestOne{}
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
func (st *SegmentTreeLongestOne) Get(index int) E {
	if index < 0 || index >= st.n {
		return st.e()
	}
	return st.seg[index+st.size]
}
func (st *SegmentTreeLongestOne) Set(index int, value E) {
	if index < 0 || index >= st.n {
		return
	}
	index += st.size
	st.seg[index] = value
	for index >>= 1; index > 0; index >>= 1 {
		st.seg[index] = st.op(st.seg[index<<1], st.seg[index<<1|1])
	}
}
func (st *SegmentTreeLongestOne) Update(index int, value E) {
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
func (st *SegmentTreeLongestOne) Query(start, end int) E {
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
func (st *SegmentTreeLongestOne) QueryAll() E { return st.seg[1] }

// 二分查询最大的 right 使得切片 [left:right] 内的值满足 predicate
func (st *SegmentTreeLongestOne) MaxRight(left int, predicate func(E) bool) int {
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
				if tmp := st.op(res, st.seg[left]); predicate(tmp) {
					res = tmp
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
func (st *SegmentTreeLongestOne) MinLeft(right int, predicate func(E) bool) int {
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
				if tmp := st.op(st.seg[right], res); predicate(tmp) {
					res = tmp
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

func (tree *SegmentTreeLongestOne) String() string {
	var sb []string
	sb = append(sb, "[")
	for i := 0; i < tree.n; i++ {
		if i != 0 {
			sb = append(sb, ", ")
		}
		sb = append(sb, fmt.Sprintf("%v", tree.Get(i)))
	}
	sb = append(sb, "]")
	return strings.Join(sb, "")
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
