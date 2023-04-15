// https://kopricky.github.io/code/Academic/fully_retroactive_deque.html
// 完全永続のように聞こえるが,
// !完全永続は過去のバージョンを更新した際にその更新を行わなかった世界と更新を行う世界を同時に保持するが,
// fully retroactive の方は更新を行った場合の 1つの世界のみを保持するという違いがある(更新は不可逆操作となる).
// 完全永続において更新は枝の分岐で表されるのに対して, fully retroactive において更新は過去の時刻での枝の挿入で表され,
// 常にバージョンは 1本のパスをなす.

// 每个操作都是均摊O(logn)的
// https://codeforces.com/gym/100451/problem/H

package main

import (
	"fmt"
)

func main() {
	qeue := NewRetroactiveDeque()

	qeue.PushBack(1, 0)
	fmt.Println(qeue.Size(0))
	qeue.PushBack(2, 1)
	fmt.Println(qeue.Size(1))
	qeue.PushBack(3, 2)
	fmt.Println(qeue.Query(0, 2))
	fmt.Println(qeue.Size(3))
	qeue.PopBack(1)
	fmt.Println(qeue.Size(3))
}

type T = int

// 完全可追溯化双端队列.
//  time: 0, 1, 2, 3, ...
type RetroactiveDeque struct {
	lroot, rroot *_Node
}

func NewRetroactiveDeque() *RetroactiveDeque {
	return &RetroactiveDeque{}
}

// 查询time时刻的双端队列的大小.
func (rd *RetroactiveDeque) Size(time int) int {
	return rd._leftPos(time) + rd._rightPos(time)
}

// 查询time时刻的双端队列的左端元素.
func (rd *RetroactiveDeque) Front(time int) T {
	return rd._queryImpl(0, time, false)
}

// 查询time时刻的双端队列的右端元素.
func (rd *RetroactiveDeque) Back(time int) T {
	return rd._queryImpl(0, time, true)
}

func (rd *RetroactiveDeque) PushFront(data T, time int) {
	rd._pushFrontImpl(data, time)
}

func (rd *RetroactiveDeque) PopFront(time int) {
	rd._popFrontImpl(time)
}

func (rd *RetroactiveDeque) PushBack(data T, time int) {
	rd._pushBackImpl(data, time)
}

func (rd *RetroactiveDeque) PopBack(time int) {
	rd._popBackImpl(time)
}

// 查询time时刻的双端队列的第index个元素.
func (rd *RetroactiveDeque) Query(index int, time int) T {
	return rd._queryImpl(index, time, false)
}

type _Node struct {
	key                  int
	val, sum, pmin, pmax int
	data                 T
	left, right, par     *_Node
}

func _NewNode(key int, data T, val int) *_Node {
	return &_Node{
		key:  key,
		val:  val,
		sum:  val,
		pmin: min(val, 0),
		pmax: max(val, 0),
		data: data,
	}
}

func (o *_Node) isRoot() bool {
	return o.par == nil
}

func (o *_Node) eval() {
	o.sum = 0
	o.pmin = 0
	o.pmax = 0
	if o.left != nil {
		o.sum += o.left.sum
		o.pmin = min(o.pmin, o.left.pmin)
		o.pmax = max(o.pmax, o.left.pmax)
	}
	o.sum += o.val
	o.pmin = min(o.pmin, o.sum)
	o.pmax = max(o.pmax, o.sum)
	if o.right != nil {
		o.pmin = min(o.pmin, o.sum+o.right.pmin)
		o.pmax = max(o.pmax, o.sum+o.right.pmax)
		o.sum += o.right.sum
	}
}

func (o *_Node) rotate(right bool) {
	p := o.par
	g := p.par
	if right {
		if p.left = o.right; p.left != nil {
			p.left.par = p
		}
		o.right = p
		p.par = o
	} else {
		if p.right = o.left; p.right != nil {
			p.right.par = p
		}
		o.left = p
		p.par = o
	}
	p.eval()
	o.eval()
	o.par = g
	if g == nil {
		return
	}
	if g.left == p {
		g.left = o
	}
	if g.right == p {
		g.right = o
	}
	g.eval()
}

func _splay(u *_Node) *_Node {
	if u == nil {
		return nil
	}
	for !u.isRoot() {
		p := u.par
		gp := p.par
		if p.isRoot() {
			u.rotate(u == p.left)
		} else {
			flag := u == p.left
			if flag == (p == gp.left) {
				p.rotate(flag)
				u.rotate(flag)
			} else {
				u.rotate(flag)
				u.rotate(!flag)
			}
		}
	}
	return u
}

func _get(key int, root *_Node) (*_Node, bool) {
	var cur, res *_Node
	nx := root
	for nx != nil {
		cur = nx
		if cur.key <= key {
			nx = cur.right
			res = cur
		} else {
			nx = cur.left
		}
	}
	tmp := _splay(cur)
	if res != nil {
		return _splay(res), true
	}
	return tmp, false
}

func _insert(ver *_Node, root *_Node) *_Node {
	if root == nil {
		return ver
	}
	var cur *_Node
	nx := root
	for nx != nil {
		cur = nx
		if cur.key > ver.key {
			nx = cur.left
		} else {
			nx = cur.right
		}
	}
	if cur.key > ver.key {
		cur.left = ver
		ver.par = cur
	} else {
		cur.right = ver
		ver.par = cur
	}
	cur.eval()
	return _splay(ver)
}

func (rd *RetroactiveDeque) _leftPos(_time int) int {
	first, second := _get(_time, rd.lroot)
	rd.lroot = first
	l := 0
	if second {
		if rd.lroot.left != nil {
			l = rd.lroot.left.sum
		}
		l += rd.lroot.val
	}
	return l
}

func (rd *RetroactiveDeque) _rightPos(_time int) int {
	first, second := _get(_time, rd.rroot)
	rd.rroot = first
	r := 0
	if second {
		if rd.rroot.left != nil {
			r = rd.rroot.left.sum
		}
		r += rd.rroot.val
	}
	return r
}

func (rd *RetroactiveDeque) _pushFrontImpl(data T, _time int) {
	newNode := _NewNode(_time, data, 1)
	rd.lroot = _insert(newNode, rd.lroot)
}

func (rd *RetroactiveDeque) _popFrontImpl(_time int) {
	var e T
	newNode := _NewNode(_time, e, -1)
	rd.lroot = _insert(newNode, rd.lroot)
}

func (rd *RetroactiveDeque) _pushBackImpl(data T, _time int) {
	newNode := _NewNode(_time, data, 1)
	rd.rroot = _insert(newNode, rd.rroot)
}

func (rd *RetroactiveDeque) _popBackImpl(_time int) {
	var e T
	newNode := _NewNode(_time, e, -1)
	rd.rroot = _insert(newNode, rd.rroot)
}

func (rd *RetroactiveDeque) _find(index int, root *_Node) (*_Node, bool) {
	if root == nil {
		return nil, false
	}
	cur, nx := root, root.left
	var res *_Node
	if nx == nil && index == 0 {
		return root, true
	}
	if nx == nil || index < nx.pmin || nx.pmax < index {
		return root, false
	}
	for nx != nil {
		cur = nx
		curSum := 0
		if cur.left != nil {
			curSum += cur.left.sum
		}
		curSum += cur.val
		if cur.right != nil {
			if curSum+cur.right.pmin <= index && index <= curSum+cur.right.pmax {
				nx = cur.right
				index -= curSum
				continue
			}
		}
		if curSum == index {
			res = cur
			break
		} else {
			nx = cur.left
		}
	}

	tmp := _splay(cur)
	if res == nil {
		return tmp, true
	}
	cur = _splay(res)
	nx = cur.right
	for nx != nil {
		cur = nx
		nx = cur.left
	}
	return _splay(cur), true
}

func (rd *RetroactiveDeque) _queryImpl(index int, _time int, back bool) T {
	lpos := rd._leftPos(_time)
	rpos := rd._rightPos(_time)
	if back {
		index = lpos + rpos - 1
	}
	lid := lpos - index
	rid := index + 1 - lpos
	lFirst, lSecond := rd._find(lid-1, rd.lroot)
	rFirst, rSecond := rd._find(rid-1, rd.rroot)
	rd.lroot = lFirst
	rd.rroot = rFirst
	if lSecond {
		if rSecond {
			if rd.lroot.key < rd.rroot.key {
				return rd.rroot.data
			}
			return rd.lroot.data
		}
		return rd.lroot.data
	}
	return rd.rroot.data
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
