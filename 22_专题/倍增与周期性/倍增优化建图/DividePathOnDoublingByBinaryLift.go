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
	db := NewDividePathOnDoublingByBinaryLift(n, int(n))
	for _, e := range edges {
		db.Add(e[1], e[0])
	}
	db.Build()

	values := make([]int32, db.Size())
	fmt.Println(db.Jump(9, 6))
	db.EnumerateJump(9, 6, func(level, index int32) {
		fmt.Println(level, index, "as")
		jumpId := level*n + index
		values[jumpId] = max32(values[jumpId], 999)
	})

	db.PushDown(func(pLevel, pIndex, cLevel, cIndex1, cIndex2 int32) {
		p, c1, c2 := pLevel*n+pIndex, cLevel*n+cIndex1, cLevel*n+cIndex2
		values[c1] = max32(values[c1], values[p])
		values[c2] = max32(values[c2], values[p])
	})
	fmt.Println(values[:n])

	// db.EnumerateJump2(7, 5, 3, func(level, index1, index2 int32) {
	// 	fmt.Println(level, index1, index2)
	// })
}

type DividePathOnDoublingByBinaryLift struct {
	n    int32
	log  int32
	size int32
	jump []int32
}

func NewDividePathOnDoublingByBinaryLift(n int32, maxStep int) *DividePathOnDoublingByBinaryLift {
	res := &DividePathOnDoublingByBinaryLift{n: n, log: int32(bits.Len(uint(maxStep))) - 1}
	res.size = n * (res.log + 1)
	res.jump = make([]int32, res.size)
	for i := range res.jump {
		res.jump[i] = -1
	}
	return res
}

func (d *DividePathOnDoublingByBinaryLift) Add(from, to int32) {
	d.jump[from] = to
}

func (d *DividePathOnDoublingByBinaryLift) Build() {
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

// 给定从 `from` 状态开始，转移 `step` 次的一段区间(一共step+1个点)，遍历这段区间上的jump。
// O(log(n)).
func (d *DividePathOnDoublingByBinaryLift) EnumerateJump(from int32, step int, f func(level, index int32)) {
	cur := from
	n, log := d.n, d.log
	for k := log; k >= 0; k-- {
		if step&(1<<k) != 0 {
			f(k, cur)
			cur = d.jump[k*n+cur]
			if cur == -1 {
				return
			}
		}
	}
	f(0, cur)

	// TODO: 这里能否拆成两段区间.
	// k := int32(bits.Len(uint(step+1)) - 1)
	// jumpLen := (step + 1) - (1 << k)
	// from2 := d.Jump(from, jumpLen)
	// f(k, from)
	// f(k, from2)
	// fmt.Println("jumpLen", jumpLen, "from2", from2)
	// fmt.Println("k", k)
}

// 从 `from1` 和 `from2` 状态开始转移 `step` 次(每段step+1个点)，遍历这两区间上的jump。
// O(log(n)).
func (d *DividePathOnDoublingByBinaryLift) EnumerateJump2(from1, from2 int32, step int, f func(level, index1, index2 int32)) {
	cur1, cur2 := from1, from2
	n, log := d.n, d.log
	for k := log; k >= 0; k-- {
		if step&(1<<k) != 0 {
			f(k, cur1, cur2)
			cur1, cur2 = d.jump[k*n+cur1], d.jump[k*n+cur2]
			if cur1 == -1 || cur2 == -1 {
				return
			}
		}
	}
	f(0, cur1, cur2)
}

// 下推路径信息，更新答案.s
// O(n*log(n)).
func (d *DividePathOnDoublingByBinaryLift) PushDown(f func(pLevel, pIndex int32, cLevel, cIndex1, cIndex2 int32)) {
	n, log := d.n, d.log
	for k := log - 1; k >= 0; k-- {
		for i := int32(0); i < n; i++ {
			// push down jump(i,k+1) to jump(i,k) and jump(jump(i,k),k)
			if to := d.jump[(k+1)*n+i]; to != -1 {
				fmt.Println("push down", k+1, i, k, i, to)
				f(k+1, i, k, i, d.jump[k*n+i])
			}
		}
	}
}

func (d *DividePathOnDoublingByBinaryLift) Size() int32 { return d.size }
func (d *DividePathOnDoublingByBinaryLift) Log() int32  { return d.log }

// 求从 `from` 状态开始转移 `step` 次的最终状态的编号。
// 不存在时返回 -1。
func (d *DividePathOnDoublingByBinaryLift) Jump(from int32, step int) (to int32) {
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
func (d *DividePathOnDoublingByBinaryLift) MaxStep(from int32, check func(next int32) bool) (step int, to int32) {
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
