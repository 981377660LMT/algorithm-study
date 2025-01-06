/*
Package list provides implementation to get and print formatted linked list.

Example:

	import "github.com/shivamMg/ppds/list"

	// a linked list node. implements list.Node.
	type Node struct {
		data int
		next *Node
	}

	func (n *Node) Data() interface{} {
		return strconv.Itoa(n.data)
	}

	func (n *Node) Next() list.Node {
		return n.next
	}

	// n1 := Node{data: 11}
	// n2 := Node{12, &n1}
	// n3 := Node{13, &n2}
	// list.Print(&n3)
*/
package main

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unicode/utf8"
)

type MyNode struct {
	data int
	next *MyNode
}

func (n *MyNode) Data() interface{} {
	return strconv.Itoa(n.data)
}

func (n *MyNode) Next() INode {
	return n.next
}

func main() {
	n1 := MyNode{data: 11}
	n2 := MyNode{12, &n1}
	n3 := MyNode{13, &n2}
	Print(&n3)
}

const (
	BoxVer       = "│"
	BoxHor       = "─"
	BoxDownLeft  = "┐"
	BoxDownRight = "┌"
	BoxUpLeft    = "┘"
	BoxUpRight   = "└"
	BoxVerLeft   = "┤"
	BoxVerRight  = "├"
)

// INode represents a linked list node.
type INode interface {
	// Data must return a value representing the node.
	Data() interface{}
	// Next must return a pointer to the next node. The last element
	// of the list must return a nil pointer.
	Next() INode
}

// Print prints the formatted linked list to standard output.
func Print(head INode) {
	fmt.Print(Sprint(head))
}

// Sprint returns the formatted linked list.
func Sprint(head INode) (s string) {
	data := []string{}
	for reflect.ValueOf(head) != reflect.Zero(reflect.TypeOf(head)) {
		data = append(data, fmt.Sprintf("%v", head.Data()))
		head = head.Next()
	}

	widths, maxDigitWidth := []int{}, digitWidth(len(data))
	for _, d := range data {
		w := utf8.RuneCountInString(d)
		if w > maxDigitWidth {
			widths = append(widths, w)
		} else {
			widths = append(widths, maxDigitWidth)
		}
	}

	for _, w := range widths {
		s += BoxDownRight + strings.Repeat(BoxHor, w) + BoxDownLeft
	}
	s += "\n"

	s += BoxVer
	for i, d := range data {
		s += fmt.Sprintf("%*s", widths[i], d)
		if i == len(data)-1 {
			s += BoxVer + "\n"
		} else {
			s += BoxVerRight + BoxVerLeft
		}
	}

	for _, w := range widths {
		s += BoxVerRight + strings.Repeat(BoxHor, w) + BoxVerLeft
	}
	s += "\n"

	s += BoxVer
	for i, w := range widths {
		s += fmt.Sprintf("%*d", w, i+1)
		if i == len(widths)-1 {
			s += BoxVer + "\n"
		} else {
			s += BoxVer + BoxVer
		}
	}

	for _, w := range widths {
		s += BoxUpRight + strings.Repeat(BoxHor, w) + BoxUpLeft
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
