// https://github.com/deadpixi/rope
//
// The rope package provides an immutable, value-oriented Rope.
// Ropes allow large sequences of text to be manipulated efficiently.

package main

import (
	"bytes"
	"io"
	"strings"
)

// https://leetcode.cn/problems/design-a-text-editor/
// 2296. 设计一个文本编辑器
type TextEditor struct {
	pos  int
	rope Rope
}

func Constructor() TextEditor {
	return TextEditor{rope: New()}
}

func (this *TextEditor) AddText(text string) {
	this.rope = this.rope.InsertString(this.pos, text)
	this.pos += len(text)
}

func (this *TextEditor) DeleteText(k int) int {
	res := min(k, this.pos)
	this.rope = this.rope.Delete(this.pos-res, res)
	this.pos -= res
	return res
}

func (this *TextEditor) CursorLeft(k int) string {
	this.pos = max(0, this.pos-k)
	start := max(0, this.pos-10)
	slice := this.rope.Slice(start, this.pos)
	return string(slice)
}

func (this *TextEditor) CursorRight(k int) string {
	this.pos = min(this.rope.Length(), this.pos+k)
	start := max(0, this.pos-10)
	slice := this.rope.Slice(start, this.pos)
	return string(slice)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

//

const (
	maxDepth    = 64
	maxLeafSize = 4096
)

var fibonacci []int

func init() {
	// The heurstic for whether a rope is balanced depends on the Fibonacci sequence;
	// we initialize the table of Fibonacci numbers here.
	first := 0
	second := 1

	for c := 0; c < maxDepth+3; c++ {
		next := 0
		if c <= 1 {
			next = c
		} else {
			next = first + second
			first = second
			second = next
		}
		fibonacci = append(fibonacci, next)
	}
}

// A Rope is a data structure for storing long runs of text.
// Ropes are persistent: there is no way to modify an existing rope.
// Instead, all operations return a new rope with the requested changes.
//
// This persistence makes it easy to store old versions of a Rope just by holding on to old roots.
type Rope struct {
	content       string
	length, depth int
	left, right   *Rope
}

// Return a new empty rope.
func New() Rope {
	return Rope{}
}

// Return a new rope with the contents of string s.
func NewString(s string) Rope {
	return Rope{content: s, length: len(s)}
}

// Notice that all of the methods take and return ropes by value.
// This is slightly less efficient than if we'd done pointers, but it
// seems cleaner from a "persistent data structure" point of view.
func (rope Rope) concat(other Rope) Rope {
	switch {
	case rope.length == 0:
		return other
	case other.length == 0:
		return rope
	case rope.length+other.length <= maxLeafSize:
		return NewString(rope.String() + other.String())
	default:
		depth := rope.depth
		if other.depth > depth {
			depth = other.depth
		}
		return Rope{
			length: rope.length + other.length,
			depth:  depth + 1,
			left:   &rope,
			right:  &other,
		}
	}
}

// Return a new rope that is the concatenation of this rope and the other rope.
func (rope Rope) Append(other Rope) Rope {
	return rope.concat(other).rebalanceIfNeeded()
}

// Return a new rope that is the concatenation of this rope and string s.
func (rope Rope) AppendString(other string) Rope {
	return rope.Append(NewString(other))
}

// Return a new rope with length bytes at offset deleted.
func (rope Rope) Delete(offset, length int) Rope {
	if length == 0 || offset == rope.length {
		return rope
	}

	left, right := rope.Split(offset)
	_, newRight := right.Split(length)
	return left.Append(newRight)
}

// Returns true if this rope is equal to other.
func (rope Rope) Equal(other Rope) bool {
	if rope == other {
		return true
	}

	if rope.length != other.length {
		return false
	}

	for i := 0; i < rope.length; i += maxLeafSize {
		if !bytes.Equal(rope.Slice(i, i+maxLeafSize), other.Slice(i, i+maxLeafSize)) {
			return false
		}
	}

	return true
}

// Return a new rope with the contents of other inserted at the given index.
func (rope Rope) Insert(at int, other Rope) Rope {
	switch at {
	case 0:
		return other.Append(rope)
	case rope.length:
		return rope.Append(other)
	default:
		left, right := rope.Split(at)
		return left.concat(other).Append(right)
	}
}

// Return a new rope with the contents of string other inserted at the given index.
func (rope Rope) InsertString(at int, other string) Rope {
	return rope.Insert(at, NewString(other))
}

// Return the length of the rope in bytes.
func (rope Rope) Length() int {
	return rope.length
}

// Return a new version of this rope that is balanced for better performance.
// Generally speaking, this will be invoked automatically during the course of other operations and
// thus only needs to be called if you know you'll be generating a lot of unbalanced ropes.
func (rope Rope) Rebalance() Rope {
	if rope.isBalanced() {
		return rope
	}

	var leaves []Rope
	rope.walk(func(node Rope) {
		leaves = append(leaves, node)
	})

	return merge(leaves, 0, len(leaves))
}

// Return the bytes in [a, b)
func (rope Rope) Slice(a, b int) []byte {
	p := make([]byte, b-a)
	n, _ := rope.ReadAt(p, int64(a))
	return p[:n]
}

// Returns two new ropes, one containing the content to the left of the given index and the other the content to the right.
func (rope Rope) Split(at int) (Rope, Rope) {
	switch {
	case rope.isLeaf():
		return NewString(rope.content[0:at]), NewString(rope.content[at:])

	case at == 0:
		return Rope{}, rope

	case at == rope.length:
		return rope, Rope{}

	case at < rope.left.length:
		left, right := rope.left.Split(at)
		return left, right.Append(*rope.right)

	case at > rope.left.length:
		left, right := rope.right.Split(at - rope.left.length)
		return rope.left.Append(left), right

	default:
		return *rope.left, *rope.right
	}
}

// Return the contents of the rope as a string.
func (rope Rope) String() string {
	if rope.isLeaf() {
		return rope.content
	}

	var builder strings.Builder
	rope.walk(func(node Rope) {
		builder.WriteString(node.content)
	})

	return builder.String()
}

func (rope Rope) isBalanced() bool {
	switch {
	case rope.isLeaf():
		return true
	case rope.depth >= len(fibonacci)-2:
		return false
	default:
		return fibonacci[rope.depth+2] <= rope.length
	}
}

func (rope Rope) isLeaf() bool {
	return rope.left == nil
}

func (rope Rope) leafForOffset(at int) (Rope, int) {
	switch {
	case rope.isLeaf():
		return rope, at
	case at < rope.left.length:
		return rope.left.leafForOffset(at)
	default:
		return rope.right.leafForOffset(at - rope.left.length)
	}
}

func (rope Rope) rebalanceIfNeeded() Rope {
	if rope.isBalanced() || abs(rope.left.depth-rope.right.depth) < maxDepth {
		return rope
	}

	return rope.Rebalance()
}

func (rope Rope) walk(callback func(Rope)) {
	if rope.isLeaf() {
		callback(rope)
	} else {
		rope.left.walk(callback)
		rope.right.walk(callback)
	}
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func merge(leaves []Rope, start, end int) Rope {
	length := end - start
	switch length {
	case 1:
		return leaves[start]
	case 2:
		return leaves[start].concat(leaves[start+1])
	default:
		mid := start + length/2
		return merge(leaves, start, mid).concat(merge(leaves, mid, end))
	}
}

// A Reader provides an implementation of io.Reader for ropes.
type Reader struct {
	rope     Rope
	position int64
}

// Return a new reader attached to the given rope.
func NewReader(rope Rope) *Reader {
	return rope.Reader()
}

// Read implements the standard Read interface:
// it reads data from the rope, populating p, and returns
// the number of bytes actually read.
func (reader *Reader) Read(p []byte) (n int, err error) {
	n, err = reader.rope.ReadAt(p, reader.position)
	if err == nil {
		reader.position += int64(n)
	}
	return
}

// Return a new Reader attached to this rope.
func (rope Rope) Reader() *Reader {
	return rope.OffsetReader(0)
}

func (rope Rope) OffsetReader(offset int) *Reader {
	return &Reader{rope: rope, position: int64(offset)}
}

// ReadAt implements the standard ReadAt interface:
// it reads len(p) bytes from offset off into p, and returns
// the number of bytes actually read. If n < len(p), err will
// explain the shortfall.
func (rope Rope) ReadAt(p []byte, off int64) (n int, err error) {
	o := int(off)
	for n < len(p) && o+n < rope.Length() {
		leaf, at := rope.leafForOffset(o + n)
		n += copy(p[n:], leaf.content[at:])
	}

	if n < len(p) {
		err = io.EOF
	}

	return
}
