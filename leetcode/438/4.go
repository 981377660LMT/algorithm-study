package main

import (
	"math/bits"
	"slices"
)

func maxDistance(side int, points [][]int, k int) int {
	trans := func(x, y int) int {
		if x == 0 {
			return y
		}
		if y == side {
			return side + x
		}
		if x == side {
			return 3*side - y
		}
		return 4*side - x
	}

	n := int32(len(points))
	nums := make([]int, len(points))
	for i, p := range points {
		nums[i] = trans(p[0], p[1])
	}
	slices.Sort(nums)

	arr2 := make([]int, 2*len(nums)+1)
	for i := int32(0); i < n; i++ {
		arr2[i+1] = nums[i]
		arr2[n+i+1] = nums[i] + 4*side
	}

	check := func(mid int) bool {
		D := NewDoublingSimple(n+n+1, k)
		right := int32(0)
		for left := int32(0); left < n+n+1; left++ {
			for right < n+n && arr2[right]-arr2[left] < mid {
				right++
			}
			D.Add(left, right)
		}
		D.Build()
		for i := int32(0); i < n; i++ {
			if D.Jump(i, k) <= i+n {
				return true
			}
		}
		return false
	}

	left, right := 0, 2*side
	res := 0
	for left <= right {
		mid := (left + right) / 2
		if check(mid) {
			res = mid
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return res
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
