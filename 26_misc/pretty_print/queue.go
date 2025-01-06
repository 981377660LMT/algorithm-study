/*
Package queue provides implementation to get and print formatted queue.

Example:

	import "github.com/shivamMg/ppds/queue"

	// a queue implementation.
	type myQueue struct {
		elems []int
	}

	func (q *myQueue) push(ele int) {
		q.elems = append([]int{ele}, q.elems...)
	}

	func (q *myQueue) pop() (int, bool) {
		if len(q.elems) == 0 {
			return 0, false
		}
		e := q.elems[len(q.elems)-1]
		q.elems = q.elems[:len(q.elems)-1]
		return e, true
	}

	// myQueue implements queue.Queue. Notice that the receiver is of *myQueue
	// type - since Push and Pop are required to modify q.
	func (q *myQueue) Push(ele interface{}) {
		q.push(ele.(int))
	}

	func (q *myQueue) Pop() (interface{}, bool) {
		return q.pop()
	}

	// q := myQueue{}
	// q.push(11)
	// q.push(12)
	// queue.Print(&q)
*/
package main

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type myQueue struct {
	elems []any
}

func (q *myQueue) Push(ele any) {
	q.elems = append([]any{ele}, q.elems...)
}

func (q *myQueue) Pop() (any, bool) {
	if len(q.elems) == 0 {
		return 0, false
	}
	e := q.elems[len(q.elems)-1]
	q.elems = q.elems[:len(q.elems)-1]
	return e, true
}

func main() {
	q := myQueue{}
	q.Push(11)
	q.Push(12)
	Print(&q)
}

const (
	Arrow        = "→"
	BoxVer       = "│"
	BoxHor       = "─"
	BoxDownLeft  = "┐"
	BoxDownRight = "┌"
	BoxUpLeft    = "┘"
	BoxUpRight   = "└"
	BoxVerLeft   = "┤"
	BoxVerRight  = "├"
)

// IQueue represents a queue of elements.
type IQueue interface {
	// Pop must pop and return the first element out of the queue. If queue is
	// empty ok should be false, else true.
	Pop() (ele interface{}, ok bool)
	// Push must insert ele in the queue. Since ele is of interface type, type
	// assertion must be done before inserting in the queue.
	Push(ele interface{})
}

// Print prints the formatted queue to standard output.
func Print(q IQueue) {
	fmt.Print(Sprint(q))
}

// Sprint returns the formatted queue.
func Sprint(q IQueue) (s string) {
	elems := []interface{}{}
	e, ok := q.Pop()
	for ok {
		elems = append(elems, e)
		e, ok = q.Pop()
	}

	for _, e := range elems {
		q.Push(e)
	}

	maxDigitWidth := digitWidth(len(elems))
	data, widths := []string{}, []int{}
	for i := len(elems) - 1; i >= 0; i-- {
		d := fmt.Sprintf("%v", elems[i])
		data = append(data, d)
		w := utf8.RuneCountInString(d)
		if w > maxDigitWidth {
			widths = append(widths, w)
		} else {
			widths = append(widths, maxDigitWidth)
		}
	}

	for _, w := range widths {
		s += " " + BoxDownRight + strings.Repeat(BoxHor, w) + BoxDownLeft
	}
	s += "\n"

	s += Arrow + BoxVer
	for i, d := range data {
		s += fmt.Sprintf("%*s", widths[i], d)
		if i == len(data)-1 {
			s += BoxVer + Arrow + "\n"
		} else {
			s += BoxVer + Arrow + BoxVer
		}
	}

	for _, w := range widths {
		s += " " + BoxVerRight + strings.Repeat(BoxHor, w) + BoxVerLeft
	}
	s += "\n"

	s += " " + BoxVer
	for i, w := range widths {
		s += fmt.Sprintf("%*d", w, i+1)
		if i == len(widths)-1 {
			s += BoxVer + "\n"
		} else {
			s += BoxVer + " " + BoxVer
		}
	}

	for _, w := range widths {
		s += " " + BoxUpRight + strings.Repeat(BoxHor, w) + BoxUpLeft
	}
	s += "\n"

	return
}

func digitWidth(d int) (w int) {
	for d != 0 {
		d = d / 10
		w++
	}
	return
}
