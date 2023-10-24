// 按位或严格大于最大值的子数组个数
package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	nums := make([]int, n)
	for i := range nums {
		fmt.Fscan(in, &nums[i])
	}

	fmt.Println(HighCry(nums))
}

// https://www.luogu.com.cn/problem/CF875D
// 求有多少个子数组满足：子数组的或严格大于子数组的最大值.
// 1<=nums.length<=2e5, 1<=nums[i]<=1e9.
func HighCry(nums []int) int {
	res := 0

	st := NewSparseTable(nums, func() int { return 0 }, max)

	// 右端点为right，左端点left的范围满足 pos1<=left<=pos2.
	// 求有多少个左端点left满足：子数组nums[left,right]最大值严格小于k.
	query := func(pos1 int, pos2 int, right int, k int) int {
		tmp := st.Query(pos2, right+1)
		if tmp >= k {
			return 0
		}
		lo, hi := right-pos2, right-pos1 // 二分最大长度
		for lo <= hi {
			mid := (lo + hi) >> 1
			if st.Query(right-mid, right+1) < k {
				lo = mid + 1
			} else {
				hi = mid - 1
			}
		}
		pos := right - lo
		return pos2 - pos
	}

	LogTrick(nums, func(a, b int) int { return a | b }, func(left []interval, right int) {
		for _, v := range left {
			start, end, or := v.leftStart, v.leftEnd, v.value
			res += query(start, end-1, right, or)
		}
	})

	return res
}

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

		// 去重
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

		// 将区间[0,pos]分成了dp.length个左闭右开区间.
		// 每一段区间的左端点left范围 在 [dp[i].leftStart,dp[i].leftEnd) 中。
		// 对应子数组 arr[left:pos+1] 的 op 值为 dp[i].value.
		for _, v := range dp {
			res[v.value] += v.leftEnd - v.leftStart
		}
		if f != nil {
			f(dp, pos)
		}
	}

	return res
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

type S = int

// SparseTable 稀疏表: st[j][i] 表示区间 [i, i+2^j) 的贡献值.
type SparseTable struct {
	st     [][]S
	lookup []int
	e      func() S
	op     func(S, S) S
	n      int
}

func NewSparseTable(leaves []S, e func() S, op func(S, S) S) *SparseTable {
	res := &SparseTable{}
	n := len(leaves)
	b := bits.Len(uint(n))
	st := make([][]S, b)
	for i := range st {
		st[i] = make([]S, n)
	}
	for i := range leaves {
		st[0][i] = leaves[i]
	}
	for i := 1; i < b; i++ {
		for j := 0; j+(1<<i) <= n; j++ {
			st[i][j] = op(st[i-1][j], st[i-1][j+(1<<(i-1))])
		}
	}
	lookup := make([]int, n+1)
	for i := 2; i < len(lookup); i++ {
		lookup[i] = lookup[i>>1] + 1
	}
	res.st = st
	res.lookup = lookup
	res.e = e
	res.op = op
	res.n = n
	return res
}

// 查询区间 [start, end) 的贡献值.
func (st *SparseTable) Query(start, end int) S {
	if start >= end {
		return st.e()
	}
	b := st.lookup[end-start]
	return st.op(st.st[b][start], st.st[b][end-(1<<b)])
}

// 返回最大的 right 使得 [left,right) 内的值满足 check.
func (ds *SparseTable) MaxRight(left int, check func(e S) bool) int {
	if left == ds.n {
		return ds.n
	}
	ok, ng := left, ds.n+1
	for ok+1 < ng {
		mid := (ok + ng) >> 1
		if check(ds.Query(left, mid)) {
			ok = mid
		} else {
			ng = mid
		}
	}
	return ok
}

// 返回最小的 left 使得 [left,right) 内的值满足 check.
func (ds *SparseTable) MinLeft(right int, check func(e S) bool) int {
	if right == 0 {
		return 0
	}
	ok, ng := right, -1
	for ng+1 < ok {
		mid := (ok + ng) >> 1
		if check(ds.Query(mid, right)) {
			ok = mid
		} else {
			ng = mid
		}
	}
	return ok
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
