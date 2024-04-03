package main

import (
	"fmt"
	"sort"
)

func main() {
	count := 0
	ops := make([]*PriorityCommutativeOperation, 10)
	for i := 0; i < 10; i++ {
		ops[i] = NewPriorityCommutativeOperation(
			func() {
				count += i
			},
			func() {
				count -= i
			},
			10-i,
		)
	}

	uq := NewUndoPriorityQueue(10)
	for _, op := range ops {
		uq.Push(op)
	}
	uq.Pop()
	uq.Pop()
	uq.Pop()
	fmt.Println(count)
}

// 满足交换律的操作，即(op1, op2) 和 (op2, op1) 作用效果相同.
type PriorityCommutativeOperation struct {
	apply, undo    func()
	priority       int
	offsetToBottom int32
}

func NewPriorityCommutativeOperation(apply, undo func(), priority int) *PriorityCommutativeOperation {
	return &PriorityCommutativeOperation{apply: apply, undo: undo, priority: priority}
}
func (op *PriorityCommutativeOperation) Apply() { op.apply() }
func (op *PriorityCommutativeOperation) Undo()  { op.undo() }

// 支持撤销操作的优先队列，每个操作会被应用和撤销 O(logn) 次.
// 注意操作的优先级必须唯一.
type UndoPriorityQueue struct {
	sl                    *specializedSortedList
	stack                 *undoStack
	bufferForHighPriority []*PriorityCommutativeOperation
	bufferForLowPriority  []*PriorityCommutativeOperation
}

func NewUndoPriorityQueue(capacity int) *UndoPriorityQueue {
	return &UndoPriorityQueue{
		sl:                    newSpecializedSortedList(func(a, b S) bool { return a.priority > b.priority }),
		stack:                 newUndoStack(capacity),
		bufferForLowPriority:  make([]*PriorityCommutativeOperation, 0, capacity),
		bufferForHighPriority: make([]*PriorityCommutativeOperation, 0, capacity),
	}
}

func (uq *UndoPriorityQueue) Push(op *PriorityCommutativeOperation) {
	uq.sl.Add(op)
	uq._pushStack(op)
}

func (uq *UndoPriorityQueue) Pop() *PriorityCommutativeOperation {
	k := int32(0)
	size := int32(uq.Len())
	uq.bufferForLowPriority = uq.bufferForLowPriority[:0]
	uq.sl.ForEachRev(func(value S) bool {
		uq.bufferForLowPriority = append(uq.bufferForLowPriority, value)
		k = max32(k, size-value.offsetToBottom)
		value.offsetToBottom = -1
		return int32(len(uq.bufferForLowPriority)*2) >= k
	})
	if k > 1 {
		uq.bufferForHighPriority = uq.bufferForHighPriority[:0]
		for i := int32(0); i < k; i++ {
			op := uq.stack.Pop()
			if op.offsetToBottom != -1 {
				uq.bufferForHighPriority = append(uq.bufferForHighPriority, op)
			}
		}
		for _, op := range uq.bufferForHighPriority {
			uq._pushStack(op)
		}
		for i, j := 0, len(uq.bufferForLowPriority)-1; i < j; i, j = i+1, j-1 {
			uq.bufferForLowPriority[i], uq.bufferForLowPriority[j] = uq.bufferForLowPriority[j], uq.bufferForLowPriority[i]
		}
		for _, op := range uq.bufferForLowPriority {
			uq._pushStack(op)
		}
	}
	res := uq.sl.PopLast()
	uq.stack.Pop()
	return res
}

func (uq *UndoPriorityQueue) Len() int { return uq.stack.Len() }

func (uq *UndoPriorityQueue) _pushStack(op *PriorityCommutativeOperation) {
	op.offsetToBottom = int32(uq.stack.Len())
	uq.stack.Push(op)
}

// 1e5 -> 200, 2e5 -> 400
const _LOAD int = 200

type S = *PriorityCommutativeOperation

type specializedSortedList struct {
	less   func(a, b S) bool
	size   int
	blocks [][]S
	mins   []S
}

func newSpecializedSortedList(less func(a, b S) bool, elements ...S) *specializedSortedList {
	elements = append(elements[:0:0], elements...)
	res := &specializedSortedList{less: less}
	sort.Slice(elements, func(i, j int) bool { return less(elements[i], elements[j]) })
	n := len(elements)
	blocks := [][]S{}
	for start := 0; start < n; start += _LOAD {
		end := min(start+_LOAD, n)
		blocks = append(blocks, elements[start:end:end]) // !各个块互不影响, max参数也需要指定为end
	}
	mins := make([]S, len(blocks))
	for i, cur := range blocks {
		mins[i] = cur[0]
	}
	res.size = n
	res.blocks = blocks
	res.mins = mins
	return res
}

func (sl *specializedSortedList) Add(value S) *specializedSortedList {
	sl.size++
	if len(sl.blocks) == 0 {
		sl.blocks = append(sl.blocks, []S{value})
		sl.mins = append(sl.mins, value)
		return sl
	}
	pos, index := sl._locRight(value)
	sl.blocks[pos] = append(sl.blocks[pos][:index], append([]S{value}, sl.blocks[pos][index:]...)...)
	sl.mins[pos] = sl.blocks[pos][0]
	// n -> load + (n - load)
	if n := len(sl.blocks[pos]); _LOAD+_LOAD < n {
		sl.blocks = append(sl.blocks[:pos+1], append([][]S{sl.blocks[pos][_LOAD:]}, sl.blocks[pos+1:]...)...)
		sl.mins = append(sl.mins[:pos+1], append([]S{sl.blocks[pos][_LOAD]}, sl.mins[pos+1:]...)...)
		sl.blocks[pos] = sl.blocks[pos][:_LOAD:_LOAD] // !注意max的设置(为了让左右互不影响)
	}
	return sl
}

func (sl *specializedSortedList) PopLast() S {
	pos := len(sl.blocks) - 1
	value := sl.blocks[pos][len(sl.blocks[pos])-1]

	// !delete element
	sl.size--
	sl.blocks[pos] = sl.blocks[pos][:len(sl.blocks[pos])-1]
	if len(sl.blocks[pos]) > 0 {
		sl.mins[pos] = sl.blocks[pos][0]
		return value
	}
	// !delete block
	sl.blocks = append(sl.blocks[:pos], sl.blocks[pos+1:]...)
	sl.mins = append(sl.mins[:pos], sl.mins[pos+1:]...)

	return value
}

func (sl *specializedSortedList) ForEachRev(f func(value S) bool) {
	count := 0
	for i := len(sl.blocks) - 1; i >= 0; i-- {
		block := sl.blocks[i]
		for j := len(block) - 1; j >= 0; j-- {
			if f(block[j]) {
				return
			}
			count++
		}
	}
}

func (sl *specializedSortedList) Len() int {
	return sl.size
}

func (sl *specializedSortedList) _locRight(value S) (pos, index int) {
	if sl.size == 0 {
		return
	}

	// find pos
	left := 0
	right := len(sl.blocks)
	for left+1 < right {
		mid := (left + right) >> 1
		if sl.less(value, sl.mins[mid]) {
			right = mid
		} else {
			left = mid
		}
	}
	pos = left

	// find index
	cur := sl.blocks[pos]
	left = -1
	right = len(cur)
	for left+1 < right {
		mid := (left + right) >> 1
		if sl.less(value, cur[mid]) {
			right = mid
		} else {
			left = mid
		}
	}

	index = right
	return
}

type undoStack struct {
	stack []*PriorityCommutativeOperation
}

func newUndoStack(capacity int) *undoStack {
	return &undoStack{stack: make([]*PriorityCommutativeOperation, 0, capacity)}
}

func (us *undoStack) Push(op *PriorityCommutativeOperation) {
	us.stack = append(us.stack, op)
	op.Apply()
}

func (us *undoStack) Pop() *PriorityCommutativeOperation {
	n := len(us.stack)
	op := us.stack[n-1]
	us.stack = us.stack[:n-1]
	op.Undo()
	return op
}

func (us *undoStack) Len() int {
	return len(us.stack)
}

func (us *undoStack) Empty() bool {
	return us.Len() == 0
}

func (us *undoStack) Clear() {
	for !us.Empty() {
		us.Pop()
	}
}

func max32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}
