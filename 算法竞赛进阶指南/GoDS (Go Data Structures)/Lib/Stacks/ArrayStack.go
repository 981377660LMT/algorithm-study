package stacks

import (
	"fmt"

	"github.com/emirpasic/gods/stacks/arraystack"
)

func a() {
	stack := arraystack.New()
	stack.Push(1)
	stack.Push(2)
	stack.Push(2)
	stack.Push(2)
	fmt.Println(stack.Peek()) // 2 true
	fmt.Println(stack.Values()...)
}
