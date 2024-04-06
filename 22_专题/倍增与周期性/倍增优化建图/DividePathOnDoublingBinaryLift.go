// DividePathOnDoublingBinaryLift
// 倍增拆分倍增结构上的路径 `link(from,len)`
// 需要处理若干个请求，每个请求要求修改路径link(from,len)上的所有结点。
// 在所有请求完成后，要求输出所有结点的权值。
//
// 一共拆分成[0,log]层，每层有n个元素.
// !jumpId = level*n + index 表示第level层的第index个元素(0<=level<log+1,0<=index<n).

package main

import (
	"fmt"
	"math/bits"
)

func main() {
	拆点()
}

func 拆点() {
	//    0
	//   / \
	//  1   2
	// / \   \
	// 3 4    5
	//         \
	//          6
	//           \
	//            7
	//             \
	//              8
	//               \
	//                9

	n := int32(10)
	edges := [][]int32{{0, 1}, {0, 2}, {1, 3}, {1, 4}, {2, 5}, {5, 6}, {6, 7}, {7, 8}, {8, 9}}
	db := NewDividePathOnDoublingBinaryLift(n, int(n))
	for _, e := range edges {
		db.Add(e[1], e[0])
	}
	db.Build()

	values := make([]int32, db.Size())
	fmt.Println(db.Jump(7, 4))
	db.EnumerateJump(7, 4, func(jumpId int32) {
		values[jumpId] = max32(values[jumpId], 999)
	})

	db.PushDown(func(parent, child1, child2 int32) {
		values[child1] = max32(values[child1], values[parent])
		values[child2] = max32(values[child2], values[parent])
	})
	fmt.Println(values[:n])
}

type DividePathOnDoublingBinaryLift struct {
	n    int32
	log  int32
	size int32
	jump []int32
}

func NewDividePathOnDoublingBinaryLift(n int32, maxStep int) *DividePathOnDoublingBinaryLift {
	res := &DividePathOnDoublingBinaryLift{n: n, log: int32(bits.Len(uint(maxStep))) - 1}
	res.size = n * (res.log + 1)
	res.jump = make([]int32, res.size)
	for i := range res.jump {
		res.jump[i] = -1
	}
	return res
}

func (d *DividePathOnDoublingBinaryLift) Add(from, to int32) {
	d.jump[from] = to
}

func (d *DividePathOnDoublingBinaryLift) Build() {
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

// 给定从 `from` 状态开始，转移 `len` 次的一段区间，遍历这段区间上的jumpId。
// O(log(n)).
func (d *DividePathOnDoublingBinaryLift) EnumerateJump(from int32, len int32, f func(jumpId int32)) {
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

// 下推路径信息，更新答案.
// O(n*log(n)).
func (d *DividePathOnDoublingBinaryLift) PushDown(f func(parent int32, child1, child2 int32)) {
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

func (d *DividePathOnDoublingBinaryLift) Size() int32 {
	return d.size
}

// 求从 `from` 状态开始转移 `step` 次的最终状态的编号。
// 不存在时返回 -1。
func (d *DividePathOnDoublingBinaryLift) Jump(from int32, step int) (to int32) {
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
func (d *DividePathOnDoublingBinaryLift) MaxStep(from int32, check func(next int32) bool) (step int, to int32) {
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
