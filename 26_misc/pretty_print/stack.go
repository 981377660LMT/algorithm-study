/*
Package stack provides implementation to get and print formatted stack.

Example:

	import "github.com/shivamMg/ppds/stack"

	// a stack implementation.
	type myStack struct {
		elems []int
	}

	func (s *myStack) push(ele int) {
		s.elems = append(s.elems, ele)
	}

	func (s *myStack) pop() (int, bool) {
		l := len(s.elems)
		if l == 0 {
			return 0, false
		}
		ele := s.elems[l-1]
		s.elems = s.elems[:l-1]
		return ele, true
	}

	// myStack implements stack.Stack. Notice that the receiver is of *myStack
	// type - since Push and Pop are required to modify s.
	func (s *myStack) Pop() (interface{}, bool) {
		return s.pop()
	}

	func (s *myStack) Push(ele interface{}) {
		s.push(ele.(int))
	}

	// s := myStack{}
	// s.push(11)
	// s.push(12)
	// stack.Print(&s)
*/
package main

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type myStack struct {
	elems []any
}

func (s *myStack) Push(ele any) {
	s.elems = append(s.elems, ele)
}

func (s *myStack) Pop() (any, bool) {
	l := len(s.elems)
	if l == 0 {
		return 0, false
	}
	ele := s.elems[l-1]
	s.elems = s.elems[:l-1]
	return ele, true
}

func main() {
	s := myStack{}
	s.Push(11)
	s.Push(12)
	Print(&s)
}

const (
	BoxVer     = "│"
	BoxHor     = "─"
	BoxUpHor   = "┴"
	BoxUpLeft  = "┘"
	BoxUpRight = "└"
)

// IStack represents a stack of elements.
type IStack interface {
	// Pop must pop the top element out of the stack and return it. In case of
	// an empty stack, ok should be false, else true.
	Pop() (ele interface{}, ok bool)
	// Push must insert an element in the stack. Since ele is of interface type,
	// type assertion must be done before inserting in the stack.
	Push(ele interface{})
}

// Print prints the formatted stack to standard output.
func Print(s IStack) {
	fmt.Print(Sprint(s))
}

// Sprint returns the formatted stack.
func Sprint(s IStack) (str string) {
	elems := []interface{}{}
	e, ok := s.Pop()
	for ok {
		elems = append(elems, e)
		e, ok = s.Pop()
	}

	for i := len(elems) - 1; i >= 0; i-- {
		s.Push(elems[i])
	}

	data, maxWidth := []string{}, 0
	for _, e := range elems {
		d := fmt.Sprintf("%v", e)
		data = append(data, d)
		w := utf8.RuneCountInString(d)
		if w > maxWidth {
			maxWidth = w
		}
	}
	// position column maxWidth
	posW := digitWidth(len(data))

	for i, d := range data {
		str += fmt.Sprintf("%s %*s %s%*d%s\n", BoxVer, maxWidth, d, BoxVer, posW, len(data)-i, BoxVer)
	}
	str += fmt.Sprint(BoxUpRight, strings.Repeat("─", maxWidth+2), BoxUpHor, strings.Repeat("─", posW), BoxUpLeft, "\n")

	return
}

func digitWidth(d int) (w int) {
	for d != 0 {
		d = d / 10
		w++
	}
	return
}
