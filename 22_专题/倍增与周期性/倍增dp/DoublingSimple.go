// 只维护转移点，不维护转移的边权.

package main

import (
	"fmt"
	"math/bits"
)

type TreeAncestor struct {
	db *DoublingSimple
}

func Constructor(n int, parent []int) TreeAncestor {
	db := NewDoublingSimple(int32(n), n)
	for i, p := range parent {
		if p != -1 {
			db.Add(int32(i), int32(p))
		}
	}
	db.Build()
	return TreeAncestor{db: db}
}

func (this *TreeAncestor) GetKthAncestor(node int, k int) int {
	return int(this.db.Jump(int32(node), k))
}

func demo() {
	n := int32(100)
	D := NewDoublingSimple(n, int(100))
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

// 求从 `from` 状态开始转移 `step` 次，满足 `check` 为 `true` 的最大的 `step` 以及最终状态的编号。
func (d *DoublingSimple) MaxStep(from int32, check func(next int32) bool) (step int, to int32) {
	for k := d.log; k >= 0; k-- {
		tmp := d.jump[k*d.n+from]
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

func (d *DoublingSimple) Size() int32 {
	return d.size
}

// 倍增拆点.
// 需要处理若干个请求，每个请求要求修改路径link(from,len)上的所有结点。
// 在所有请求完成后，要求输出所有结点的权值。
func (d *DoublingSimple) EnumerateJump(from int32, len int32, f func(id int32)) {
	cur := from
	n, log := d.n, d.log
	for k := log; k >= 0; k-- {
		if cur == -1 {
			return
		}
		if len&(1<<k) != 0 {
			f(k*n + cur)
			cur = d.jump[k*n+cur]
		}
	}
	f(cur)
}

// O(n*log(n)).
func (d *DoublingSimple) PushDown(f func(parent int32, child1, child2 int32)) {
	n, log := d.n, d.log
	for k := log - 1; k >= 0; k-- {
		for i := int32(0); i < n; i++ {
			// push down jump(i,k+1) to jump(i,k) and jump(jump(i,k),k)
			parent := (k+1)*n + i
			if to := d.jump[parent]; to != -1 {
				left := k*n + i
				right := k*n + d.jump[left]
				f(parent, left, right)
			}
		}
	}
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
