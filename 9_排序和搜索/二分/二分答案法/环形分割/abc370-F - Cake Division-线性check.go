// F - Cake Division(环上最大化最小值)
// https://atcoder.jp/contests/abc370/tasks/abc370_f
// !给定一个环形数组,分成k个非空连续子数组,使得这k个子数组的和的最小值最大,求出最大值.
// !并求出哪些点不能作为第一段的起点.
// 1<=k<=n<=10^5 1<=nums[i]<=10^9

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

	check := func(mid int) bool {
		{
			// 先求解链上的问题(剪枝)
			count := 0
			left := 0
			for right := 0; right < n; right++ {
				if presum[right+1]-presum[left] >= mid {
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
			for right < n && presum[right]-presum[left] < mid {
				right++
			}
			if presum[right]-presum[left] >= mid {
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
			if presum[n]-(presum[end]-presum[i]) >= mid {
				count++
			}
			if count >= k {
				return true
			}
		}

		return false
	}

	left, right := 1, presum[n]/k+1
	for left <= right {
		mid := (left + right) / 2
		if check(mid) {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}

	fmt.Fprintln(out, right)

	// !求使得最小值>=right的起点的个数.
	// 可以做到O(n)，但是用倍增比较好理解.
	solve2 := func(target int) int {
		n32 := int32(n)
		presum2 := make([]int, 2*n32+1)
		for i := int32(0); i < n32; i++ {
			presum2[i+1] = presum2[i] + nums[i]
		}
		for i := int32(0); i < n32; i++ {
			presum2[n32+i+1] = presum2[n32+i] + nums[i]
		}

		res := 0
		D := NewDoublingSimple(n32+n32+1, int(k))
		right := int32(0)
		for left := int32(0); left < n32+n32+1; left++ {
			for right < n32+n32 && presum2[right]-presum2[left] < target {
				right++
			}
			D.Add(left, right)
		}
		D.Build()
		for i := int32(0); i < n32; i++ {
			if D.Jump(i, int(k)) <= i+n32 {
				res++
			}
		}
		return res
	}

	bad := n - solve2(right)
	fmt.Fprintln(out, bad)
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

// 求从 `from` 状态开始转移，满足 `check` 为 `true` 的最小的 `step` 以及最终状态的编号。
// 如果不存在，则返回 (-1, -1).
func (d *DoublingSimple) FirstTrue(from int32, check func(next int32) bool) (step int, to int32) {
	if check(from) {
		return 0, from
	}
	for k := d.log; k >= 0; k-- {
		if tmp := d.jump[k*d.n+from]; tmp != -1 && !check(tmp) {
			step |= 1 << k
			from = tmp
		}
	}
	step++
	to = d.jump[from]
	if to == -1 {
		step = -1
	}
	return
}

// !MaxStep.
// 求从 `from` 状态开始转移，满足 `check` 为 `true` 的最大的 `step` 以及最终状态的编号。
// 如果不存在，则返回 (-1, -1).
func (d *DoublingSimple) LastTrue(from int32, check func(next int32) bool) (step int, to int32) {
	if !check(from) {
		return -1, -1
	}
	for k := d.log; k >= 0; k-- {
		if tmp := d.jump[k*d.n+from]; tmp != -1 && check(tmp) {
			step |= 1 << k
			from = tmp
		}
	}
	to = from
	return
}
