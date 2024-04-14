package main

import "fmt"

// 给你两个数组 nums 和 andValues，长度分别为 n 和 m。
// 数组的 值 等于该数组的 最后一个 元素。
// 你需要将 nums 划分为 m 个 不相交的连续 子数组，对于第 ith 个子数组 [li, ri]，子数组元素的按位AND运算结果等于 andValues[i]，换句话说，对所有的 1 <= i <= m，nums[li] & nums[li + 1] & ... & nums[ri] == andValues[i] ，其中 & 表示按位AND运算符。
// 返回将 nums 划分为 m 个子数组所能得到的可能的 最小 子数组 值 之和。如果无法完成这样的划分，则返回 -1 。

type interval = struct{ leftStart, leftEnd, value int }

// 将 nums 的所有非空子数组的元素进行 op 操作，返回所有不同的结果和其出现次数.
//
// nums: 1 <= nums.length <= 1e5.
// op: 与/或/gcd/lcm 中的一种操作，具有单调性.
// f:
// 数组的右端点为right.
// interval 的 leftStart/leftEnd 表示子数组的左端点left的范围.
// interval 的 value 表示该子数组 arr[left,right] 的 op 结果.
func LogTrick(nums []int, op func(int, int) int, f func(left []interval, right int)) map[int]int {
	res := make(map[int]int)

	dp := []interval{}
	for pos, cur := range nums {
		for i, pre := range dp {
			dp[i].value = op(pre.value, cur)
		}
		dp = append(dp, interval{leftStart: pos, leftEnd: pos + 1, value: cur})

		ptr := 0
		for _, v := range dp[1:] {
			if v.value != dp[ptr].value {
				ptr++
				dp[ptr] = v
			} else {
				dp[ptr].leftEnd = v.leftEnd
			}
		}
		dp = dp[:ptr+1]

		for _, v := range dp {
			res[v.value] += v.leftEnd - v.leftStart
		}
		if f != nil {
			f(dp, pos)
		}
	}

	return res
}

func minimumValueSum(nums []int, andValues []int) int {
	n := len(nums)
	groupByRight := make([][]interval, n)
	LogTrick(nums, func(i1, i2 int) int { return i1 & i2 }, func(left []interval, right int) {
		groupByRight[right] = append(groupByRight[right], left...)
	})

	dp := NewSegmentTree(n+1, func(i int) int { return INF })
	dp.Set(0, 0)
	for _, andValue := range andValues {
		ndp := NewSegmentTree(n+1, func(i int) int { return INF })
		for i := 1; i <= n; i++ {
			groupInfo := groupByRight[i-1]
			for _, interval := range groupInfo {
				if interval.value == andValue {
					start, end := interval.leftStart, interval.leftEnd
					cand := dp.Query(start, end) + nums[i-1]
					ndp.Update(i, cand)
				}
			}
		}
		dp = ndp
	}

	res := dp.Get(n)
	if res == INF {
		return -1
	}
	return res
}

func main() {
	// nums = [1,4,3,3,2], andValues = [0,3,3,2]
	// [{0 1 1}] 0
	// [{0 1 0} {1 2 4}] 1
	// [{0 2 0} {2 3 3}] 2
	// [{0 2 0} {2 4 3}] 3
	// [{0 2 0} {2 5 2}] 4
	fmt.Println(minimumValueSum([]int{1, 4, 3, 3, 2}, []int{0, 3, 3, 2}))
}

const INF int = 1e18

// PointSetRangeMin

type E = int

func (*SegmentTree) e() E        { return INF }
func (*SegmentTree) op(a, b E) E { return min(a, b) }
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

func (st *SegmentTree) Copy() *SegmentTree {
	seg := make([]E, len(st.seg))
	copy(seg, st.seg)
	return &SegmentTree{st.n, st.size, seg}
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
