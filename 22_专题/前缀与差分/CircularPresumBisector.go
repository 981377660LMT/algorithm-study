package main

import (
	"bufio"
	"fmt"
	"os"
)

const INF int = 1e18

func main() {
	// abc368_c_1()
	abc368_c_2()
}

// 2875. 无限数组的最短子数组
// https://leetcode.cn/problems/minimum-size-subarray-in-infinite-array/
// 求循环数组中和为 target 的最短子数组的长度.不存在则返回 -1.
// 1 <= nums.length <= 1e5
// 1 <= nums[i] <= 1e5
// 1 <= target <= 1e9
func minSizeSubarray(nums []int, target int) int {
	Q := NewCircularPresumBisector(len(nums), func(i int) int { return nums[i] })
	res := INF
	for start := range nums {
		right := Q.MaxRight(start, func(sum int, right int) bool { return sum <= target }, 1e9+10)
		if sum := Q.Query(start, right); sum == target {
			res = min(res, right-start)
		}
	}
	if res == INF {
		return -1
	}
	return res
}

// C - Triple Attack
// https://atcoder.jp/contests/abc368/tasks/abc368_c_1
// !攻击怪兽。普通攻击-1，三连击-3，每三个回合可以使用一次三连击。怪兽血量为h[i]，问最少需要多少回合才能击败所有怪兽。
// 先处理循环节，然后处理剩余部分。
func abc368_c_1() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	nums := make([]int, n)
	for i := range nums {
		fmt.Fscan(in, &nums[i])
	}

	turn := 0
	for _, v := range nums {
		cycle := v / 5
		turn += cycle * 3
		v -= cycle * 5
		for v > 0 {
			turn++
			if turn%3 == 0 {
				v -= 3
			} else {
				v--
			}
		}
	}

	fmt.Fprintln(out, turn)
}

// C - Triple Attack
// https://atcoder.jp/contests/abc368/tasks/abc368_c_1
// !攻击怪兽。普通攻击-1，三连击-3，每三个回合可以使用一次三连击。怪兽血量为h[i]，问最少需要多少回合才能击败所有怪兽。
func abc368_c_2() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	nums := make([]int, n)
	for i := range nums {
		fmt.Fscan(in, &nums[i])
	}
	attack := [3]int{1, 1, 3}
	S := NewCircularPresumBisector(3, func(i int) int { return attack[i] })

	turn := 0
	for _, hp := range nums {
		right := S.MaxRight(turn, func(sum int, right int) bool { return sum <= hp }, 1e15)
		if s := S.Query(turn, right); s < hp {
			turn = right + 1
		} else {
			turn = right
		}
	}
	fmt.Fprintln(out, turn)
}

// 带有二分的环形前缀和，要求元素为非负数.
type CircularPresumBisector struct {
	n      int
	presum []int
}

func NewCircularPresumBisector(n int, f func(i int) int) *CircularPresumBisector {
	preSum := make([]int, n+1)
	for i := 1; i <= n; i++ {
		preSum[i] = preSum[i-1] + f(i-1)
	}
	return &CircularPresumBisector{n: n, presum: preSum}
}

func (c *CircularPresumBisector) Query(start, end int) int {
	if start < 0 {
		start = 0
	}
	if start >= end {
		return 0
	}
	return c.calc(end) - c.calc(start)
}

// 返回最大的 right 使得 [left,right) 内的值满足 check.
func (p *CircularPresumBisector) MaxRight(left int, check func(sum int, right int) bool, upper int) int {
	ok, ng := left, upper+1
	for ok+1 < ng {
		mid := (ok + ng) >> 1
		if check(p.Query(left, mid), mid) {
			ok = mid
		} else {
			ng = mid
		}
	}
	return ok
}

// 返回最小的 left 使得 [left,right) 内的值满足 check.
func (p *CircularPresumBisector) MinLeft(right int, check func(sum int, left int) bool, lower int) int {
	if right <= 0 {
		return 0
	}
	ok, ng := right, lower-1
	for ng+1 < ok {
		mid := (ok + ng) >> 1
		if check(p.Query(mid, right), mid) {
			ok = mid
		} else {
			ng = mid
		}
	}
	return ok
}

func (c *CircularPresumBisector) calc(r int) int {
	return c.presum[c.n]*(r/c.n) + c.presum[r%c.n]
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

func min32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func max32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}
