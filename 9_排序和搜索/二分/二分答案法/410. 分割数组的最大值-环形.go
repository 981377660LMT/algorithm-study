// No.1211 円環はお断り(圆环)
// https://yukicoder.me/problems/no/1211
// F - Cake Division
// https://atcoder.jp/contests/abc370/tasks/abc370_f
// !给定一个环形数组,分成k个非空连续子数组,使得这k个子数组的和的最小值最大,求出最大值.并求出哪些点不能作为第一段的起点.
// 1<=k<=n<=10^5 1<=nums[i]<=10^9
//
// 0.断环成链
// 1.二分mid
// 2.预处理出从每个位置出发,最远可以走多远,使得和不超过mid
// 3.倍增.

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

	var n, k int32
	fmt.Fscan(in, &n, &k)
	nums := make([]int, n)
	for i := int32(0); i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	presum := make([]int, 2*n+1)
	for i := int32(0); i < n; i++ {
		presum[i+1] = presum[i] + nums[i]
	}
	for i := int32(0); i < n; i++ {
		presum[n+i+1] = presum[n+i] + nums[i]
	}

	// 每段的和是否 >= mid
	check := func(mid int) bool {
		D := NewDoublingSimple(n+n+1, int(k))
		right := int32(0)
		for left := int32(0); left < n+n+1; left++ {
			for right < n+n && presum[right]-presum[left] < mid {
				right++
			}
			D.Add(left, right)
		}
		D.Build()
		for i := int32(0); i < n; i++ {
			if D.Jump(i, int(k)) <= i+n {
				return true
			}
		}
		return false
	}

	left, right := 1, presum[len(presum)-1]/int(k)+1
	for left <= right {
		mid := (left + right) / 2
		if check(mid) {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}

	// !此时最小值为right.
	mid := right
	validStart := int32(0)
	{
		D := NewDoublingSimple(n+n+1, int(k))
		right := int32(0)
		for left := int32(0); left < n+n+1; left++ {
			for right < n+n && presum[right]-presum[left] < mid {
				right++
			}
			D.Add(left, right)
		}
		D.Build()
		for i := int32(0); i < n; i++ {
			if D.Jump(i, int(k)) <= i+n {
				validStart++
			}
		}
	}

	fmt.Fprintln(out, right, n-validStart)
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
