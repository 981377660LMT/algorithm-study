package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {
	yuki1211()
}

// https://yukicoder.me/problems/no/1211
func yuki1211() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k int
	fmt.Fscan(in, &n, &k)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	presum := make([]int, n+1)
	for i := 0; i < n; i++ {
		presum[i+1] = presum[i] + nums[i]
	}

	cost := func(start, end int) int {
		return presum[end] - presum[start]
	}

	res := MaximizeMinCostOnCycleDoubling(n, cost, k, 1, presum[n]/k+1)
	fmt.Println(res)
}

// 给定一个n个点的环形数组, [start, end) 的代价为 cost(start, end), 且 cost 满足单调性.
// 将环形数组分成k个非空连续子数组, 最大化这k个子数组的代价的最小值.
// 返回这个最大的最小值.
// !O(n*log(upper-lower))时间复杂度.
func MaximizeMinCostOnCycleDp(
	n int, cost func(start, end int) int, k int,
	lower, upper int,
) int {
	if k > n {
		panic("k must be not greater than n")
	}

	check := func(mid int) bool {
		{
			// 先求解链上的问题(剪枝)
			count := 0
			left := 0
			for right := 0; right < n; right++ {
				if cost(left, right+1) >= mid {
					count++
					left = right + 1
				}
			}
			if count >= k {
				return true
			}
			if count <= k-2 {
				return false
			}
		}

		next := make([]int, n)
		right := 0
		for left := 0; left < n; left++ {
			for right < n && cost(left, right) < mid {
				right++
			}
			if cost(left, right) >= mid {
				next[left] = right
			} else {
				next[left] = -1
			}
		}

		type dpItem struct{ count, next int }
		dp := make([]dpItem, n+1)
		dp[n] = dpItem{next: n}
		for i := n - 1; i >= 0; i-- {
			if next[i] == -1 {
				dp[i] = dpItem{next: i}
			} else {
				dp[i] = dp[next[i]]
				dp[i].count++
			}
		}

		for i := 0; i < n; i++ {
			count := dp[i].count
			if count <= k-2 {
				break
			}
			end := dp[i].next
			// 检查最后一段是否满足条件.
			if cost(0, i)+cost(end, n) >= mid {
				count++
			}
			if count >= k {
				return true
			}
		}

		return false
	}

	left, right := lower, upper
	for left <= right {
		mid := (left + right) / 2
		if check(mid) {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return right
}

// 给定一个n个点的环形数组, [start, end) 的代价为 cost(start, end), 且 cost 满足单调性.
// 将环形数组分成k个非空连续子数组, 最大化这k个子数组的代价的最小值.
// 返回这个最大的最小值.
// !O(n*log(upper-lower)*logk)时间复杂度.
func MaximizeMinCostOnCycleDoubling(
	n int, cost func(start, end int) int, k int,
	lower, upper int,
) int {
	if k > n {
		panic("k must be not greater than n")
	}

	costWrapper := func(start, end int) int {
		if start >= end {
			return 0
		}
		if end <= n {
			return cost(start, end)
		}
		if start >= n {
			return cost(start-n, end-n)
		}
		return cost(start, n) + cost(0, end-n)
	}

	check := func(mid int) bool {
		{
			// 先求解链上的问题(剪枝)
			count := 0
			left := 0
			for right := 0; right < n; right++ {
				if cost(left, right+1) >= mid {
					count++
					left = right + 1
				}
			}
			if count >= k {
				return true
			}
			if count <= k-2 {
				return false
			}
		}

		n32 := int32(n)
		D := NewDoublingSimple(n32+n32+1, k)
		right := int32(0)
		for left := int32(0); left < n32+n32+1; left++ {
			for right < n32+n32 && costWrapper(int(left), int(right)) < mid {
				right++
			}
			D.Add(left, right)
		}
		D.Build()
		for i := int32(0); i < n32; i++ {
			if D.Jump(i, k) <= i+n32 {
				return true
			}
		}
		return false
	}

	left, right := lower, upper
	for left <= right {
		mid := (left + right) / 2
		if check(mid) {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return right
}

type DoublingSimple struct {
	n    int32
	log  int32
	size int32
	jump []int32
}

func NewDoublingSimple(n int32, maxStep int) *DoublingSimple {
	res := &DoublingSimple{n: n, log: int32(bits.Len(uint(maxStep))) - 1}
	res.size = n * (res.log + 1)
	res.jump = make([]int32, res.size)
	for i := range res.jump {
		res.jump[i] = -1
	}
	return res
}

func (d *DoublingSimple) Add(from, to int32) {
	d.jump[from] = to
}

func (d *DoublingSimple) Build() {
	n := d.n
	for k := int32(0); k < d.log; k++ {
		for v := int32(0); v < n; v++ {
			w := d.jump[k*n+v]
			next := (k+1)*n + v
			if w == -1 {
				d.jump[next] = -1
				continue
			}
			d.jump[next] = d.jump[k*n+w]
		}
	}
}

// 求从 `from` 状态开始转移 `step` 次的最终状态的编号。
// 不存在时返回 -1。
func (d *DoublingSimple) Jump(from int32, step int) (to int32) {
	to = from
	for k := int32(0); k < d.log+1; k++ {
		if to == -1 {
			return
		}
		if step&(1<<k) != 0 {
			to = d.jump[k*d.n+to]
		}
	}
	return
}
