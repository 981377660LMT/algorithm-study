/*
Package tree provides implementation to get and print formatted tree.

Example:

	import "github.com/shivamMg/ppds/tree"

	// a tree node.
	type Node struct {
		data int
		children []*Node
	}

	func (n *Node) Data() interface{} {
		return strconv.Itoa(n.data)
	}

	// cannot return n.children directly.
	// https://github.com/golang/go/wiki/InterfaceSlice
	func (n *Node) Children() (c []tree.Node) {
		for _, child := range n.children {
			c = append(c, tree.Node(child))
		}
		return
	}

	// n1, n2 := Node{data: "b"}, Node{data: "c"}
	// n3 := Node{"a", []*Node{&n1, &n2}}
	// tree.Print(&n3)

*/

package main

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"
)

type MyNode struct {
	data     int
	children []*MyNode
}

func (n *MyNode) Data() interface{} {
	return n.data
}

func (n *MyNode) Children() []INode {
	var c []INode
	for _, child := range n.children {
		c = append(c, INode(child))
	}
	return c
}

func main() {
	n1, n2 := MyNode{data: 11}, MyNode{data: 12}
	n3 := MyNode{13, []*MyNode{&n1, &n2}}
	Print(&n3)
}

const (
	BoxVer       = "│"
	BoxHor       = "─"
	BoxVerRight  = "├"
	BoxDownLeft  = "┐"
	BoxDownRight = "┌"
	BoxDownHor   = "┬"
	BoxUpRight   = "└"
	// Gutter is number of spaces between two adjacent child nodes.
	Gutter = 2
)

// ErrDuplicateNode indicates that a duplicate Node (node with same hash) was
// encountered while going through the tree. As of now Sprint/Print and
// SprintWithError/PrintWithError cannot operate on such trees.
//
// This error is returned by SprintWithError/PrintWithError. It's also used
// in Sprint/Print as error for panic for the same case.
//
// FIXME: create internal representation of trees that copies data
var ErrDuplicateNode = errors.New("duplicate node")

// INode represents a node in a tree. Type that satisfies INode must be a hashable type.
type INode interface {
	// Data must return a value representing the node. It is stringified using "%v".
	// If empty, a space is used.
	Data() interface{}
	// Children must return a list of all child nodes of the node.
	Children() []INode
}

type queue struct {
	arr []INode
}

func (q queue) empty() bool {
	return len(q.arr) == 0
}

func (q queue) len() int {
	return len(q.arr)
}

func (q *queue) push(n INode) {
	q.arr = append(q.arr, n)
}

func (q *queue) pop() INode {
	if q.empty() {
		return nil
	}
	ele := q.arr[0]
	q.arr = q.arr[1:]
	return ele
}

func (q *queue) peek() INode {
	if q.empty() {
		return nil
	}
	return q.arr[0]
}

// Print prints the formatted tree to standard output. To handle ErrDuplicateNode use PrintWithError.
func Print(root INode) {
	fmt.Print(Sprint(root))
}

// Sprint returns the formatted tree. To handle ErrDuplicateNode use SprintWithError.
func Sprint(root INode) string {
	parents := map[INode]INode{}
	if err := setParents(parents, root); err != nil {
		panic(err)
	}
	return sprint(parents, root)
}

// PrintWithError prints the formatted tree to standard output.
func PrintWithError(root INode) error {
	s, err := SprintWithError(root)
	if err != nil {
		return err
	}
	fmt.Print(s)
	return nil
}

// SprintWithError returns the formatted tree.
func SprintWithError(root INode) (string, error) {
	parents := map[INode]INode{}
	if err := setParents(parents, root); err != nil {
		return "", err
	}
	return sprint(parents, root), nil
}

func sprint(parents map[INode]INode, root INode) string {
	isLeftMostChild := func(n INode) bool {
		p, ok := parents[n]
		if !ok {
			// root
			return true
		}
		return p.Children()[0] == n
	}

	paddings := map[INode]int{}
	setPaddings(paddings, map[INode]int{}, 0, root)

	q := queue{}
	q.push(root)
	lines := []string{}
	for !q.empty() {
		// line storing branches, and line storing nodes
		branches, nodes := "", ""
		// runes covered
		covered := 0
		qLen := q.len()
		for i := 0; i < qLen; i++ {
			n := q.pop()
			for _, c := range n.Children() {
				q.push(c)
			}

			spaces := paddings[n] - covered
			data := safeData(n)
			nodes += strings.Repeat(" ", spaces) + data

			w := utf8.RuneCountInString(data)
			covered += spaces + w
			current, next := isLeftMostChild(n), isLeftMostChild(q.peek())
			if current {
				branches += strings.Repeat(" ", spaces)
			} else {
				branches += strings.Repeat(BoxHor, spaces)
			}

			if current && next {
				branches += BoxVer
			} else if current {
				branches += BoxVerRight
			} else if next {
				branches += BoxDownLeft
			} else {
				branches += BoxDownHor
			}

			if next {
				branches += strings.Repeat(" ", w-1)
			} else {
				branches += strings.Repeat(BoxHor, w-1)
			}
		}
		lines = append(lines, branches, nodes)
	}

	s := ""
	// ignore first line since it's the branch above root
	for _, line := range lines[1:] {
		s += strings.TrimRight(line, " ") + "\n"

	}
	return s
}

// safeData always returns non-empty representation of n's data. Empty data
// messes up tree structure, and ignoring such node will return incomplete
// tree output (tree without an entire subtree). So it returns a space.
func safeData(n INode) string {
	data := fmt.Sprintf("%v", n.Data())
	if data == "" {
		return " "
	}
	return data
}

// setPaddings sets left padding (distance of a node from the root)
// for each node in the tree.
func setPaddings(paddings map[INode]int, widths map[INode]int, pad int, root INode) {
	for _, c := range root.Children() {
		paddings[c] = pad
		setPaddings(paddings, widths, pad, c)
		pad += width(widths, c)
	}
}

// setParents sets child-parent relationships for the tree rooted
// at root.
func setParents(parents map[INode]INode, root INode) error {
	for _, c := range root.Children() {
		if _, ok := parents[c]; ok {
			return ErrDuplicateNode
		}
		parents[c] = root
		if err := setParents(parents, c); err != nil {
			return err
		}
	}
	return nil
}

// width returns either the sum of widths of it's children or its own
// data length depending on which one is bigger. widths is used in
// memoization.
func width(widths map[INode]int, n INode) int {
	if w, ok := widths[n]; ok {
		return w
	}

	w := utf8.RuneCountInString(safeData(n)) + Gutter
	widths[n] = w
	if len(n.Children()) == 0 {
		return w
	}

	sum := 0
	for _, c := range n.Children() {
		sum += width(widths, c)
	}
	if sum > w {
		widths[n] = sum
		return sum
	}
	return w
}

// PrintHr prints the horizontal formatted tree to standard output.
func PrintHr(root INode) {
	fmt.Print(SprintHr(root))
}

// SprintHr returns the horizontal formatted tree.
func SprintHr(root INode) (s string) {
	for _, line := range lines(root) {
		// ignore runes before root node
		line = string([]rune(line)[2:])
		s += strings.TrimRight(line, " ") + "\n"
	}
	return
}

func lines(root INode) (s []string) {
	data := fmt.Sprintf("%s %v ", BoxHor, root.Data())
	l := len(root.Children())
	if l == 0 {
		s = append(s, data)
		return
	}

	w := utf8.RuneCountInString(data)
	for i, c := range root.Children() {
		for j, line := range lines(c) {
			if i == 0 && j == 0 {
				if l == 1 {
					s = append(s, data+BoxHor+line)
				} else {
					s = append(s, data+BoxDownHor+line)
				}
				continue
			}

			var box string
			if i == l-1 && j == 0 {
				// first line of the last child
				box = BoxUpRight
			} else if i == l-1 {
				box = " "
			} else if j == 0 {
				box = BoxVerRight
			} else {
				box = BoxVer
			}
			s = append(s, strings.Repeat(" ", w)+box+line)
		}
	}
	return
}

// PrintHrn prints the horizontal-newline formatted tree to standard output.
func PrintHrn(root INode) {
	fmt.Print(SprintHrn(root))
}

// SprintHrn returns the horizontal-newline formatted tree.
func SprintHrn(root INode) (s string) {
	return strings.Join(lines2(root), "\n") + "\n"
}

func lines2(root INode) (s []string) {
	s = append(s, fmt.Sprintf("%v", root.Data()))
	l := len(root.Children())
	if l == 0 {
		return
	}

	for i, c := range root.Children() {
		for j, line := range lines2(c) {
			// first line of the last child
			if i == l-1 && j == 0 {
				s = append(s, BoxUpRight+BoxHor+" "+line)
			} else if j == 0 {
				s = append(s, BoxVerRight+BoxHor+" "+line)
			} else if i == l-1 {
				s = append(s, "   "+line)
			} else {
				s = append(s, BoxVer+"  "+line)
			}
		}
	}
	return
}
