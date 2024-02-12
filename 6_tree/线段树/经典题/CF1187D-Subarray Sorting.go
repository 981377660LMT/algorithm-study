package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// Subarray Sorting
// https://www.luogu.com.cn/problem/CF1187D
// 给定两个长度相等的数组a,b，每次操作可以将a的子数组[l,r]从小到大排序。问能否通过若干次操作使序列a变为序列b。
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	solve := func(nums1, nums2 []int) bool {
		n := len(nums1)
		mp1, mp2 := make(map[int][]int), make(map[int][]int)
		for i := 0; i < n; i++ {
			mp1[nums1[i]] = append(mp1[nums1[i]], i)
			mp2[nums2[i]] = append(mp2[nums2[i]], i)
		}

		to := make(map[int]int, n) // !映射的相对顺序
		for key := range mp1 {
			pos1, pos2 := mp1[key], mp2[key]
			if len(pos1) != len(pos2) {
				return false
			}
			for i := 0; i < len(pos1); i++ {
				to[pos1[i]] = pos2[i]
			}
		}

		seg := NewSegmentTree(n, func(i int) int { return -INF })

		sortedKeys := make([]int, 0, len(mp1))
		for key := range mp1 {
			sortedKeys = append(sortedKeys, key)
		}
		sort.Ints(sortedKeys)
		for _, key := range sortedKeys {
			pos := mp1[key]
			for _, i := range pos {
				if seg.Query(0, i) >= to[i] { // 最大index大于等于to[i]，说明改变了相对顺序
					return false
				}
			}
			for _, i := range pos {
				seg.Set(i, to[i])
			}
		}

		return true
	}

	var T int
	fmt.Fscan(in, &T)
	for i := 0; i < T; i++ {
		var n int
		fmt.Fscan(in, &n)
		nums1 := make([]int, n)
		for j := 0; j < n; j++ {
			fmt.Fscan(in, &nums1[j])
		}
		nums2 := make([]int, n)
		for j := 0; j < n; j++ {
			fmt.Fscan(in, &nums2[j])
		}
		res := solve(nums1, nums2)
		if res {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}

const INF int = 1e18

// PointSetRangeMax

type E = int

func (*SegmentTree) e() E        { return -INF }
func (*SegmentTree) op(a, b E) E { return max(a, b) }
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

func NewSegmentTree(n int, f func(int) E) *SegmentTree {
	res := &SegmentTree{}
	size := 1
	for size < n {
		size <<= 1
	}
	seg := make([]E, size<<1)
	for i := range seg {
		seg[i] = res.e()
	}
	for i := 0; i < n; i++ {
		seg[i+size] = f(i)
	}
	for i := size - 1; i > 0; i-- {
		seg[i] = res.op(seg[i<<1], seg[i<<1|1])
	}
	res.n = n
	res.size = size
	res.seg = seg
	return res
}
func NewSegmentTreeFrom(leaves []E) *SegmentTree {
	res := &SegmentTree{}
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
func (st *SegmentTree) Update(index int, value E) {
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
func (st *SegmentTree) GetAll() []E {
	res := make([]E, st.n)
	copy(res, st.seg[st.size:st.size+st.n])
	return res
}

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
