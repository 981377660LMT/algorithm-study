package main

type Operation struct{ apply, undo func() }

func NewOperation(apply, undo func()) *Operation { return &Operation{apply: apply, undo: undo} }
func (op *Operation) Apply()                     { op.apply() }
func (op *Operation) Undo()                      { op.undo() }

type UndoStack struct {
	stack []*Operation
}

func NewUndoStack(capacity int32) *UndoStack {
	return &UndoStack{stack: make([]*Operation, 0, capacity)}
}

func (us *UndoStack) Push(op *Operation) {
	us.stack = append(us.stack, op)
	op.Apply()
}

func (us *UndoStack) Pop() *Operation {
	n := len(us.stack)
	op := us.stack[n-1]
	us.stack = us.stack[:n-1]
	op.Undo()
	return op
}

func (us *UndoStack) Len() int {
	return len(us.stack)
}

func (us *UndoStack) Empty() bool {
	return us.Len() == 0
}

func (us *UndoStack) Clear() {
	for !us.Empty() {
		us.Pop()
	}
}
