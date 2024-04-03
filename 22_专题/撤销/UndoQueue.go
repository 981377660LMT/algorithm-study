// https://codeforces.com/contest/1423/problem/H
// 给定n条顶点，m条无向边，之后提供q个请求，第i个请求为li,ri，要求回答仅考虑编号在[li,ri]之间的边，判断所有顶点是否连通。其中1≤n,m≤105,1≤q≤106

// https://codeforces.com/contest/1386/problem/C
// 提供n个商品，第i件商品的价格为wi，价值为vi。
// 我们总共有m单位金钱，希望能买到总价值最大的货物。
// 换言之，我们希望选择一些商品，这些商品的总价格不超过m的前提下保证总价值最大。
// 同时，如果我们同时购买k件物品，那么这k件物品中的一件我们可以免费以五折优化买走（如果价格不能整除2,则上取整），且上述优惠最多只能发生一次。
// 问最大能取走的货物的总价值。其中1≤n,m≤3000，1≤wi,vi≤109，1≤k≤n。

package main

import "fmt"

func main() {
	demo()
}

func demo() {
	sum := 0
	Q := NewUndoQueue(10)
	Q.Append(NewFlaggedCommutativeOperation(
		func() {
			sum++
		},
		func() {
			sum--
		},
	))
	fmt.Println(sum)
	Q.PopLeft()
	fmt.Println(sum)

}

// 满足交换律的操作，即(op1, op2) 和 (op2, op1) 作用效果相同.
type FlaggedCommutativeOperation struct {
	apply, undo func()
	flag        bool
}

func NewFlaggedCommutativeOperation(apply, undo func()) *FlaggedCommutativeOperation {
	return &FlaggedCommutativeOperation{apply: apply, undo: undo}
}
func (op *FlaggedCommutativeOperation) Apply() { op.apply() }
func (op *FlaggedCommutativeOperation) Undo()  { op.undo() }

// 支持撤销操作的队列，每个操作会被应用和撤销 O(logn) 次.
type UndoQueue struct {
	dq, bufA, bufB []*FlaggedCommutativeOperation
}

func NewUndoQueue(capacity int) *UndoQueue {
	return &UndoQueue{
		dq:   make([]*FlaggedCommutativeOperation, 0, capacity),
		bufA: make([]*FlaggedCommutativeOperation, 0, capacity),
		bufB: make([]*FlaggedCommutativeOperation, 0, capacity),
	}
}

func (uq *UndoQueue) Append(op *FlaggedCommutativeOperation) {
	op.flag = false
	uq._pushAndDo(op)
}

func (uq *UndoQueue) PopLeft() *FlaggedCommutativeOperation {
	if !uq.dq[len(uq.dq)-1].flag {
		uq._popAndUndo()
		for len(uq.dq) > 0 && len(uq.bufB) != len(uq.bufA) {
			uq._popAndUndo()
		}
		if len(uq.dq) == 0 {
			for len(uq.bufB) > 0 {
				res := uq.bufB[0]
				uq.bufB = uq.bufB[1:]
				res.flag = true
				uq._pushAndDo(res)
			}
		} else {
			for len(uq.bufB) > 0 {
				res := uq.bufB[len(uq.bufB)-1]
				uq.bufB = uq.bufB[:len(uq.bufB)-1]
				uq._pushAndDo(res)
			}
		}
		for len(uq.bufA) > 0 {
			uq._pushAndDo(uq.bufA[len(uq.bufA)-1])
			uq.bufA = uq.bufA[:len(uq.bufA)-1]
		}
	}

	res := uq.dq[len(uq.dq)-1]
	uq.dq = uq.dq[:len(uq.dq)-1]
	res.Undo()
	return res
}

func (uq *UndoQueue) Len() int {
	return len(uq.dq)
}

func (uq *UndoQueue) Empty() bool {
	return len(uq.dq) == 0
}

func (uq *UndoQueue) Clear() {
	n := len(uq.dq)
	for i := 0; i < n; i++ {
		uq._popAndUndo()
	}
	uq.bufA = uq.bufA[:0]
	uq.bufB = uq.bufB[:0]
}

func (uq *UndoQueue) _pushAndDo(op *FlaggedCommutativeOperation) {
	uq.dq = append(uq.dq, op)
	op.Apply()
}

func (uq *UndoQueue) _popAndUndo() {
	res := uq.dq[len(uq.dq)-1]
	uq.dq = uq.dq[:len(uq.dq)-1]
	res.Undo()
	if res.flag {
		uq.bufA = append(uq.bufA, res)
	} else {
		uq.bufB = append(uq.bufB, res)
	}
}
