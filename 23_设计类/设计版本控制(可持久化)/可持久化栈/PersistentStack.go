package main

import "fmt"

func main() {
	stack := NewPersistentStack()
	stack = stack.Push(1)
	stack = stack.Push(2)
	stack = stack.Push(3)
	fmt.Println(stack)
}

type S int

type PersistentStack struct {
	root *StackNode
}

type StackNode struct {
	value S
	pre   *StackNode
}

// 创建一个新的可持久化栈.
func NewPersistentStack() *PersistentStack {
	return &PersistentStack{}
}

func (s *PersistentStack) Push(value S) *PersistentStack {
	return &PersistentStack{root: &StackNode{value: value, pre: s.root}}
}

func (s *PersistentStack) Pop() *PersistentStack {
	if s.root == nil {
		panic("stack is empty")
	}
	return &PersistentStack{root: s.root.pre}
}

func (s *PersistentStack) Top() S {
	if s.root == nil {
		panic("stack is empty")
	}
	return s.root.value
}

func (s *PersistentStack) Empty() bool {
	return s.root == nil
}

func (s *PersistentStack) Reverse() *PersistentStack {
	res := NewPersistentStack()
	x := s
	for !x.Empty() {
		res = res.Push(x.Top())
		x = x.Pop()
	}
	return res
}

func (s *PersistentStack) String() string {
	sb := []string{}
	x := s
	for !x.Empty() {
		sb = append(sb, fmt.Sprintf("%v", x.Top()))
		x = x.Pop()
	}
	for i, j := 0, len(sb)-1; i < j; i, j = i+1, j-1 {
		sb[i], sb[j] = sb[j], sb[i]
	}
	return fmt.Sprintf("Stack%v", sb)
}
