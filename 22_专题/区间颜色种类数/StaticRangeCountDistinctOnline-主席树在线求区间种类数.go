// StaticRangeCountDistinctOnline-主席树在线求区间种类数
// 较慢，推荐用二维数点的版本.

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {
	P1972()
}

// P1972 [SDOI2009] HH的项链
// https://www.luogu.com.cn/problem/P1972
// https://judge.yosupo.jp/problem/static_range_count_distinct
// 区间颜色种类数
func P1972() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	nums := make([]int32, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}
	S := NewStaticRangeCountDistinct(nums)

	for i := 0; i < q; i++ {
		var start, end int
		fmt.Fscan(in, &start, &end)
		// start--
		fmt.Fprintln(out, S.Query(start, end))
	}
}

// 持久化权值线段树求区间颜色种类数.
// !对于询问区间[i,j]，直接用[i,j]减去重复颜色的数量。
// 思路是记录每个数上一次出现的位置pre.
type StaticRangeCountDistinctOnline struct {
	n                   int32
	ptr                 int32
	root                int32
	left, right, preSum []int32
	roots               []int32
}

func NewStaticRangeCountDistinct(nums []int32) *StaticRangeCountDistinctOnline {
	n := int32(len(nums))
	maxLog := int32(bits.Len(uint(n))) + 1
	size := 2 * n * maxLog
	res := &StaticRangeCountDistinctOnline{
		n:      n,
		ptr:    1,
		root:   1,
		left:   make([]int32, size),
		right:  make([]int32, size),
		preSum: make([]int32, size),
		roots:  make([]int32, n+1),
	}
	last := make(map[int32]int32)
	for i := int32(0); i < n; i++ {
		res.makeRoot()
		x := nums[i]
		res.add(last[x], -1, res.roots[res.root-1], 0, n+1)
		last[x] = i + 1
		res.add(i+1, 1, res.roots[res.root-1], 0, n+1)
	}
	return res
}

func (sr *StaticRangeCountDistinctOnline) Query(start, end int) int {
	a, b := int32(start+1), int32(end+1)
	return int(sr.get(a, b, sr.roots[b-1], 0, sr.n+1))
}

func (sr *StaticRangeCountDistinctOnline) add(pos, delta, root, left, right int32) {
	sr.preSum[root] += delta
	if right-left > 1 {
		mid := (left + right) >> 1
		if pos < mid {
			sr.left[root] = sr.copy(sr.left[root])
			sr.add(pos, delta, sr.left[root], left, mid)
		} else {
			sr.right[root] = sr.copy(sr.right[root])
			sr.add(pos, delta, sr.right[root], mid, right)
		}
	}
}

func (sr *StaticRangeCountDistinctOnline) get(a, b, root, left, right int32) int32 {
	if right <= a || b <= left {
		return 0
	} else if a <= left && right <= b {
		return sr.preSum[root]
	} else {
		mid := (left + right) >> 1
		return sr.get(a, b, sr.left[root], left, mid) + sr.get(a, b, sr.right[root], mid, right)
	}
}

func (sr *StaticRangeCountDistinctOnline) copy(v int32) int32 {
	sr.left[sr.ptr] = sr.left[v]
	sr.right[sr.ptr] = sr.right[v]
	sr.preSum[sr.ptr] = sr.preSum[v]
	sr.ptr++
	return sr.ptr - 1
}

func (sr *StaticRangeCountDistinctOnline) makeRoot() {
	sr.roots[sr.root] = sr.copy(sr.roots[sr.root-1])
	sr.root++
}
