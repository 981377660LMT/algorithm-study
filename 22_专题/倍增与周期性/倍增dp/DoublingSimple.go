// 只维护转移点，不维护转移的边权.

package main

import (
	"fmt"
	"math/bits"
)

func main() {
	n := int32(100)
	D := NewDoubling(n, int(100))
	values := make([]int32, n)
	for i := int32(0); i < n; i++ {
		values[i] = n - 10*i
	}
	for i := int32(0); i < n-1; i++ {
		D.Add(i, i+1)
	}
	D.Build()

	start := int32(0)
	step, to := D.MaxStep(start, func(next int32) bool { return values[next] >= 50 })
	fmt.Println(step, to)
	fmt.Println(D.Jump(start, step))
	fmt.Println(values[to])
}

type DoublingSimple struct {
	n   int32
	log int32
	to  []int32
}

func NewDoubling(n int32, maxStep int) *DoublingSimple {
	res := &DoublingSimple{n: n, log: int32(bits.Len(uint(maxStep)))}
	size := n * res.log
	res.to = make([]int32, size)
	for i := int32(0); i < size; i++ {
		res.to[i] = -1
	}
	return res
}

func (d *DoublingSimple) Add(from, to int32) {
	d.to[from] = to
}

func (d *DoublingSimple) Build() {
	n := d.n
	for k := int32(0); k < d.log-1; k++ {
		for v := int32(0); v < n; v++ {
			w := d.to[k*n+v]
			next := (k+1)*n + v
			if w == -1 {
				d.to[next] = -1
				continue
			}
			d.to[next] = d.to[k*n+w]
		}
	}
}

// 求从 `from` 状态开始转移 `step` 次的最终状态的编号。
// 不存在时返回 -1。
func (d *DoublingSimple) Jump(from int32, step int) (to int32) {
	to = from
	for k := int32(0); k < d.log; k++ {
		if to == -1 {
			return
		}
		if step&(1<<k) != 0 {
			to = d.to[k*d.n+to]
		}
	}
	return
}

// 求从 `from` 状态开始转移 `step` 次，满足 `check` 为 `true` 的最大的 `step` 以及最终状态的编号。
func (d *DoublingSimple) MaxStep(from int32, check func(next int32) bool) (step int, to int32) {
	for k := d.log - 1; k >= 0; k-- {
		tmp := d.to[k*d.n+from]
		if tmp == -1 {
			continue
		}
		if check(tmp) {
			step |= 1 << k
			from = tmp
		}
	}
	to = from
	return
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
