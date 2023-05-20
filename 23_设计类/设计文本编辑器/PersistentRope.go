// Persistent Rope
//
// Description:
//   Rope is a binary tree data structure to maintains a sequence.
//

package main

import (
	"fmt"
	"math/rand"
	"strings"
)

func main() {
	root := NewPersistentRope([]string{"a", "b", "c", "d", "e"})
	root = root.Insert(0, []string{"A", "B", "C", "D", "E"})
	fmt.Println(root)
	fmt.Println(root.Substring(0, 3).Erase(-1, 30))
}

// https://leetcode.cn/problems/design-a-text-editor/
// 2296. 设计一个文本编辑器
type TextEditor struct {
	pos  int
	rope *PersistentRope
}

func Constructor() TextEditor {
	return TextEditor{rope: NewPersistentRope(nil)}
}

func (this *TextEditor) AddText(text string) {
	cur := make([]string, len(text))
	for i := 0; i < len(text); i++ {
		cur[i] = string(text[i])
	}
	this.rope = this.rope.Insert(this.pos, cur)
	this.pos += len(text)
}

func (this *TextEditor) DeleteText(k int) int {
	res := min(k, this.pos)
	this.rope = this.rope.Erase(this.pos-res, this.pos)
	this.pos -= res
	return res
}

func (this *TextEditor) CursorLeft(k int) string {
	this.pos = max(0, this.pos-k)
	res := this.rope.Substring(max(0, this.pos-10), this.pos)
	return res.String()
}

func (this *TextEditor) CursorRight(k int) string {
	this.pos = min(this.rope.Size(), this.pos+k)
	res := this.rope.Substring(max(0, this.pos-10), this.pos)
	return res.String()
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

type R = string
type PersistentRope struct {
	root *_RNode
}

func NewPersistentRope(bytes []R) *PersistentRope {
	res := &PersistentRope{}
	res.root = res._build(0, len(bytes), bytes)
	return res
}

func (rp *PersistentRope) Insert(i int, bytes []R) *PersistentRope {
	p1, p2 := rp._split(rp.root, i)
	newRoot := rp._join(p1, rp._join(NewPersistentRope(bytes).root, p2))
	return &PersistentRope{newRoot}
}

func (rp *PersistentRope) Substring(start, end int) *PersistentRope {
	first, _ := rp._cut(rp.root, start, end)
	return &PersistentRope{first}
}

func (rp *PersistentRope) Erase(start, end int) *PersistentRope {
	_, second := rp._cut(rp.root, start, end)
	return &PersistentRope{second}
}

func (rp *PersistentRope) Enumetate(f func(R)) {
	var dfs func(node *_RNode)
	dfs = func(node *_RNode) {
		if node == nil {
			return
		}
		dfs(node.l)
		f(node.v)
		dfs(node.r)
	}
	dfs(rp.root)
}

// 0<=i<rp.Size()
func (rp *PersistentRope) At(i int) R {
	return rp._at(i, rp.root)
}

func (rp *PersistentRope) Size() int {
	if rp.root == nil {
		return 0
	}
	return rp.root.s
}

func (rp *PersistentRope) String() string {
	sb := []string{}
	rp._string(rp.root, &sb)
	return strings.Join(sb, "")
}

func (rp *PersistentRope) _string(node *_RNode, sb *[]string) {
	if node == nil {
		return
	}
	rp._string(node.l, sb)
	*sb = append(*sb, node.v)
	rp._string(node.r, sb)
}

func (rp *PersistentRope) _build(l, r int, bytes []R) *_RNode {
	if l >= r {
		return nil
	}
	m := (l + r) >> 1
	return _NewNode(bytes[m], rp._build(l, m, bytes), rp._build(m+1, r, bytes))
}

func (rp *PersistentRope) _at(i int, node *_RNode) R {
	s := 0
	if node.l != nil {
		s = node.l.s
	}
	if i == s {
		return node.v
	}
	if i < s {
		return rp._at(i, node.l)
	}
	return rp._at(i-s-1, node.r)
}

func (rp *PersistentRope) _join(a, b *_RNode) *_RNode {
	if a == nil || b == nil {
		if a != nil {
			return a
		}
		return b
	}
	if r, s1, s2 := rand.Int(), a.s, b.s; r%(s1+s2) < s1 {
		return _NewNode(a.v, a.l, rp._join(a.r, b))
	}
	return _NewNode(b.v, rp._join(a, b.l), b.r)
}

func (rp *PersistentRope) _split(a *_RNode, s int) (*_RNode, *_RNode) {
	if a == nil || s <= 0 {
		return nil, a
	}
	if a.s <= s {
		return a, nil
	}
	if a.l != nil && s <= a.l.s {
		p1, p2 := rp._split(a.l, s)
		return p1, _NewNode(a.v, p2, a.r)
	}
	tmp := 0
	if a.l != nil {
		tmp = a.l.s
	}
	l, r := rp._split(a.r, s-tmp-1)
	return _NewNode(a.v, a.l, l), r
}

// (sub,rest)
func (rp *PersistentRope) _cut(a *_RNode, l, r int) (*_RNode, *_RNode) {
	if l >= r {
		return nil, a
	}
	p1, p2 := rp._split(a, l)
	q1, q2 := rp._split(p2, r-l)
	return q1, rp._join(p1, q2)
}

type _RNode struct {
	v    R // 叶子结点存储的字符
	l, r *_RNode
	s    int // 非叶子结点存储的子树字符长度之和
}

func _NewNode(v R, l, r *_RNode) *_RNode {
	res := &_RNode{v: v, l: l, r: r, s: 1}
	if l != nil {
		res.s += l.s
	}
	if r != nil {
		res.s += r.s
	}
	return res
}

// https://leetcode.cn/problems/extract-kth-character-from-the-rope-tree/
