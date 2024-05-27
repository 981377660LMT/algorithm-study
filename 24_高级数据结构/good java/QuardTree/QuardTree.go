package main

import "fmt"

func main() {
	node := NewQuadTree(1, 2, 3, 4)
	node.Update(1, 2, 3, 4, 1, 2, 3, 4)
}

func (tree *QuadTree) Modify() {
	fmt.Println("modify")
}

func (tree *QuadTree) PushUp() {
	fmt.Println("pushUp")
}

func (tree *QuadTree) PushDown() {
	if tree == NIL {
		return
	}
	fmt.Println("pushDown")
}

var NIL = &QuadTree{}

// 四叉树
// 00 01
// 10 11
type QuadTree struct {
	s00, s01, s10, s11 *QuadTree
}

func NewQuadTree(x1, x2, y1, y2 int32) *QuadTree {
	if x1 > x2 || y1 > y2 {
		return NIL
	}
	res := &QuadTree{}
	xm := (x1 + x2) >> 1
	ym := (y1 + y2) >> 1
	if x1 < x2 || y1 < y2 {
		res.s00 = NewQuadTree(x1, xm, y1, ym)
		res.s01 = NewQuadTree(x1, xm, ym+1, y2)
		res.s10 = NewQuadTree(xm+1, x2, y1, ym)
		res.s11 = NewQuadTree(xm+1, x2, ym+1, y2)
		res.PushUp()
	}
	return res
}

func (tree *QuadTree) Update(tx1, tx2, ty1, ty2, x1, x2, y1, y2 int32) {
	if tx1 > x2 || tx2 < x1 || ty1 > y2 || ty2 < y1 {
		return
	}
	if tx1 <= x1 && x2 <= tx2 && ty1 <= y1 && y2 <= ty2 {
		return
	}
	mx := (x1 + x2) >> 1
	my := (y1 + y2) >> 1
	tree.PushDown()
	tree.s00.Update(tx1, tx2, ty1, ty2, x1, mx, y1, my)
	tree.s01.Update(tx1, tx2, ty1, ty2, x1, mx, my+1, y2)
	tree.s10.Update(tx1, tx2, ty1, ty2, mx+1, x2, y1, my)
	tree.s11.Update(tx1, tx2, ty1, ty2, mx+1, x2, my+1, y2)
	tree.PushUp()
}

func (tree *QuadTree) Query(tx1, tx2, ty1, ty2, x1, x2, y1, y2 int32) {
	if tx1 > x2 || tx2 < x1 || ty1 > y2 || ty2 < y1 {
		return
	}
	if tx1 <= x1 && x2 <= tx2 && ty1 <= y1 && y2 <= ty2 {
		tree.Modify()
		return
	}
	mx := (x1 + x2) >> 1
	my := (y1 + y2) >> 1
	tree.PushDown()
	tree.s00.Query(tx1, tx2, ty1, ty2, x1, mx, y1, my)
	tree.s01.Query(tx1, tx2, ty1, ty2, x1, mx, my+1, y2)
	tree.s10.Query(tx1, tx2, ty1, ty2, mx+1, x2, y1, my)
	tree.s11.Query(tx1, tx2, ty1, ty2, mx+1, x2, my+1, y2)
}
