//
// Randomized KD Tree (for d = 2)
//
// Description
//   Randomized KD tree is a binary space partition tree.
//   Each tree node has a point, direction, and two childs.
//   The left descendants are located to the `left' of the point p
//   and the right descendants are located to the `right' of the point p,
//   where q is `left' of p means q[dir] < p[dir].
//
//   By randomizing the direction, we can easily implement
//   insertion and deletion.
//
// Complexity
//   O(log n) insertion and deletion.
//   O(log n) search for a random points.
//
// !非常快

// API:
//   NewRandomKdTree(calDist func(p1, p2 Point2D) V) *RandomKdTree
//   Insert(p Point2D)
//   Remove(p Point2D)
//   Nearest(p Point2D, allowOverlap bool) (Point2D, V)
//   Size() int

package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
)

func demo() {
	tree := NewRandomKdTree(func(p1, p2 Point2D) V {
		return (p1.x-p2.x)*(p1.x-p2.x) + (p1.y-p2.y)*(p1.y-p2.y)
	})

	tree.Insert(Point2D{1, 2})
	tree.Insert(Point2D{1, 2})
	fmt.Println(tree.Size())
	tree.Remove(Point2D{1, 2})
	fmt.Println(tree.Size())
	fmt.Println(tree.Nearest(Point2D{1, 2}, true))
	tree.Remove(Point2D{1, 2})
	fmt.Println(tree.Size())
	fmt.Println(tree.Nearest(Point2D{1, 2}, true))
}

// https://atcoder.jp/contests/abc283/tasks/abc283_f
// 569 ms
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	points := make([]Point2D, n)
	for i := 0; i < n; i++ {
		var x int
		fmt.Fscan(in, &x)
		points[i] = Point2D{x, i}
	}

	kdtree := NewRandomKdTree(func(p1, p2 Point2D) V {
		dist := abs(p1.x-p2.x) + abs(p1.y-p2.y)
		return dist * dist
	})
	for i := 0; i < n; i++ {
		kdtree.Insert(points[i])
	}

	for i := 0; i < n; i++ {
		kdtree.Remove(points[i])
		p, _ := kdtree.Nearest(points[i], true)
		dist := abs(p.x-points[i].x) + abs(p.y-points[i].y)
		kdtree.Insert(points[i])
		fmt.Fprintln(out, dist)
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

const INF V = 1e18

type V = int

type RandomKdTree struct {
	root     *TNode
	calDist2 func(p1, p2 Point2D) V // 计算两点距离的平方
	size     int
}

type TNode struct {
	p    Point2D
	d, s int
	l, r *TNode
}

func (t *TNode) IsLeftOf(x *TNode) bool {
	if x.d != 0 {
		return t.p.x < x.p.x
	}
	return t.p.y < x.p.y
}

type Point2D struct{ x, y V }

// 计算两点距离的平方(norm in cpp).
func NewRandomKdTree(calDist2 func(p1, p2 Point2D) V) *RandomKdTree {
	return &RandomKdTree{calDist2: calDist2}
}

// 查询距离p最近的点.需要保证树不为空.
//  allowOverlap: 是否允许p和最近点重合.
func (tree *RandomKdTree) Nearest(p Point2D, allowOverlap bool) (res Point2D, ok bool) {
	if tree.Size() == 0 {
		ok = false
		return
	}
	pair1 := INF
	var pair2 *TNode
	tree._closest(tree.root, p, &pair1, &pair2, allowOverlap)
	if pair2 == nil {
		ok = false
		return
	}
	ok = true
	res = pair2.p
	return
}

func (tree *RandomKdTree) Insert(p Point2D) {
	tree.root = tree._insert(tree.root, &TNode{p: p, d: rand.Int() & 1})
	tree.size++
}

// 删除点p.需要保证p在树中.
func (tree *RandomKdTree) Remove(p Point2D) {
	tree.root = tree._remove(tree.root, &TNode{p: p})
	tree.size--
}

func (tree *RandomKdTree) Size() int {
	return tree.size
}

func (tree *RandomKdTree) _size(t *TNode) int {
	if t == nil {
		return 0
	}
	return t.s
}

func (tree *RandomKdTree) _update(t *TNode) *TNode {
	t.s = 1 + tree._size(t.l) + tree._size(t.r)
	return t
}

func (tree *RandomKdTree) _split(t, x *TNode) (*TNode, *TNode) {
	if t == nil {
		return nil, nil
	}
	if t.d == x.d {
		if t.IsLeftOf(x) {
			l, r := tree._split(t.r, x)
			t.r = l
			return tree._update(t), r
		}
		l, r := tree._split(t.l, x)
		t.l = r
		return l, tree._update(t)
	}
	l, r := tree._split(t.l, x)
	l2, r2 := tree._split(t.r, x)
	if t.IsLeftOf(x) {
		t.l = l
		t.r = l2
		return tree._update(t), tree._join(r, r2, t.d)
	}
	t.l = r
	t.r = r2
	return tree._join(l, l2, t.d), tree._update(t)
}

func (tree *RandomKdTree) _join(l, r *TNode, d int) *TNode {
	if l == nil {
		return r
	}
	if r == nil {
		return l
	}
	if sl, sr := tree._size(l), tree._size(r); rand.Int()%(sl+sr) < sl {
		if l.d == d {
			l.r = tree._join(l.r, r, d)
			return tree._update(l)
		}
		l2, r2 := tree._split(r, l)
		l.l = tree._join(l.l, l2, d)
		l.r = tree._join(l.r, r2, d)
		return tree._update(l)
	}
	if r.d == d {
		r.l = tree._join(l, r.l, d)
		return tree._update(r)
	}
	l2, r2 := tree._split(l, r)
	r.l = tree._join(l2, r.l, d)
	r.r = tree._join(r2, r.r, d)
	return tree._update(r)
}

func (tree *RandomKdTree) _insert(t, x *TNode) *TNode {
	if rand.Int()%(tree._size(t)+1) == 0 {
		l, r := tree._split(t, x)
		x.l = l
		x.r = r
		return tree._update(x)
	}
	if x.IsLeftOf(t) {
		t.l = tree._insert(t.l, x)
	} else {
		t.r = tree._insert(t.r, x)
	}
	return tree._update(t)
}

func (tree *RandomKdTree) _remove(t, x *TNode) *TNode {
	if t == nil {
		return nil
	}
	if t.p == x.p {
		return tree._join(t.l, t.r, t.d)
	}
	if x.IsLeftOf(t) {
		t.l = tree._remove(t.l, x)
	} else {
		t.r = tree._remove(t.r, x)
	}
	return tree._update(t)
}

func (tree *RandomKdTree) _closest(t *TNode, p Point2D, pair1 *V, pair2 **TNode, allowOverlap bool) {
	if t == nil {
		return
	}
	r := tree.calDist2(t.p, p)
	if r == 0 && !allowOverlap {
		r = INF
	}
	if r < *pair1 {
		*pair1 = r
		*pair2 = t
	}
	fst, snd := t.r, t.l
	var w V
	if t.d != 0 {
		w = p.x - t.p.x
	} else {
		w = p.y - t.p.y
	}
	if w < 0 {
		fst, snd = snd, fst
	}
	tree._closest(fst, p, pair1, pair2, allowOverlap)
	if *pair1 > w*w {
		tree._closest(snd, p, pair1, pair2, allowOverlap)
	}
}
