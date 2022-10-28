package main

import (
	"fmt"

	"github.com/emirpasic/gods/stacks/linkedliststack"
)

func main() {
	stack := linkedliststack.New()
	stack.Push(1)
	stack.Push(2)
	stack.Push(2)
	stack.Push(2)
	fmt.Println(stack.Peek()) // 2 true
	fmt.Println(stack.Pop())  // 2 true
	fmt.Println(stack.Values()...)
}
