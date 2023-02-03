// 二维线段树-区间修改区间查询

package main

import (
	"fmt"
	"math/bits"
)

// 单组测试数据时可以禁用gc加速
// func init() {
// 	debug.SetGCPercent(-1)
// }

const INF int = 1e18

// 子矩形查询 区间染色单点查询
// https://leetcode.cn/problems/subrectangle-queries/
type SubrectangleQueries struct {
	tree *SegmentTree2D
}

func Constructor(rectangle [][]int) SubrectangleQueries {
	return SubrectangleQueries{NewLazySegmentTree2D(rectangle)}
}

// 用 newValue 更新以 (row1,col1) 为左上角且以 (row2,col2) 为右下角的子矩形。
func (this *SubrectangleQueries) UpdateSubrectangle(row1 int, col1 int, row2 int, col2 int, newValue int) {
	this.tree.Update(row1, col1, row2, col2, newValue)
}

// 返回矩形中坐标 (row,col) 的当前值。
func (this *SubrectangleQueries) GetValue(row int, col int) int {
	return this.tree.Get(row, col)
}

/**
 * Your SubrectangleQueries object will be instantiated and called as such:
 * obj := Constructor(rectangle);
 * obj.UpdateSubrectangle(row1,col1,row2,col2,newValue);
 * param_2 := obj.GetValue(row,col);
 */

type E = int  // 单点查询
type Id = int // 区间染色

func (*SegmentTree2D) e() E   { return INF }
func (*SegmentTree2D) id() Id { return INF }
func (*SegmentTree2D) op(e1, e2 E) E {
	if e1 == INF {
		return e2
	}
	return e1
}
func (*SegmentTree2D) mapping(f Id, e E) E {
	if f == INF {
		return e
	}
	return f
}
func (*SegmentTree2D) composition(f, g Id) Id {
	if f == INF {
		return g
	}
	return f
}

//
//
//
//
type SegmentTree2D struct {
	row, col int
	data     []E
	lazy     []Id
}

func NewLazySegmentTree2D(leaves [][]E) *SegmentTree2D {
	tree := &SegmentTree2D{}
	n, m := len(leaves), len(leaves[0])
	tree.row, tree.col = n, m
	log1, log2 := int(bits.Len(uint(n-1))), int(bits.Len(uint(m-1)))
	size1, size2 := 1<<log1, 1<<log2
	tree.data = make([]E, 4*size1*size2)
	tree.lazy = make([]Id, 2*size1*size2) // !叶子结点不需要更新lazy (composition)
	for i := range tree.data {
		tree.data[i] = tree.e()
	}
	for i := range tree.lazy {
		tree.lazy[i] = tree.id()
	}
	tree.build(1, 1, 1, n, m, leaves)
	return tree
}

func (t *SegmentTree2D) build(root, row1, col1, row2, col2 int, leaves [][]E) {
	if row1 > row2 || col1 > col2 {
		return
	}
	if row1 == row2 && col1 == col2 {
		t.data[root] = leaves[row1-1][col1-1]
		return
	}
	rowMid, colMid := (row1+row2)>>1, (col1+col2)>>1
	t.build((root<<2)-2, row1, col1, rowMid, colMid, leaves)
	t.build((root<<2)-1, row1, colMid+1, rowMid, col2, leaves)
	t.build((root << 2), rowMid+1, col1, row2, colMid, leaves)
	t.build((root<<2)+1, rowMid+1, colMid+1, row2, col2, leaves)
	t.pushUp(root)
}

func (t *SegmentTree2D) query(root int, ROW1, COL1, ROW2, COL2, row1, col1, row2, col2 int) E {
	if ROW1 <= row1 && row2 <= ROW2 && COL1 <= col1 && col2 <= COL2 {
		return t.data[root]
	}

	t.pushDown(root)
	rowMid, colMid := (row1+row2)>>1, (col1+col2)>>1
	res := t.e()
	if ROW1 <= rowMid {
		if COL1 <= colMid {
			res = t.op(res, t.query((root<<2)-2, ROW1, COL1, ROW2, COL2, row1, col1, rowMid, colMid))
		}
		if colMid < COL2 {
			res = t.op(res, t.query((root<<2)-1, ROW1, COL1, ROW2, COL2, row1, colMid+1, rowMid, col2))
		}
	}
	if rowMid < ROW2 {
		if COL1 <= colMid {
			res = t.op(res, t.query((root<<2), ROW1, COL1, ROW2, COL2, rowMid+1, col1, row2, colMid))
		}
		if colMid < COL2 {
			res = t.op(res, t.query((root<<2)+1, ROW1, COL1, ROW2, COL2, rowMid+1, colMid+1, row2, col2))
		}
	}

	return res
}

func (t *SegmentTree2D) update(root int, ROW1, COL1, ROW2, COL2, row1, col1, row2, col2 int, lazy Id) {
	if ROW1 <= row1 && row2 <= ROW2 && COL1 <= col1 && col2 <= COL2 {
		t.propagate(root, lazy)
		return
	}

	t.pushDown(root)
	rowMid, colMid := (row1+row2)>>1, (col1+col2)>>1
	if ROW1 <= rowMid {
		if COL1 <= colMid {
			t.update((root<<2)-2, ROW1, COL1, ROW2, COL2, row1, col1, rowMid, colMid, lazy)
		}
		if colMid < COL2 {
			t.update((root<<2)-1, ROW1, COL1, ROW2, COL2, row1, colMid+1, rowMid, col2, lazy)
		}
	}
	if rowMid < ROW2 {
		if COL1 <= colMid {
			t.update((root << 2), ROW1, COL1, ROW2, COL2, rowMid+1, col1, row2, colMid, lazy)
		}
		if colMid < COL2 {
			t.update((root<<2)+1, ROW1, COL1, ROW2, COL2, rowMid+1, colMid+1, row2, col2, lazy)
		}
	}

	t.pushUp(root)
}

func (t *SegmentTree2D) get(root int, ROW, COL, row1, col1, row2, col2 int) E {
	if row1 == row2 && col1 == col2 {
		return t.data[root]
	}

	t.pushDown(root)
	rowMid, colMid := (row1+row2)>>1, (col1+col2)>>1
	if ROW <= rowMid {
		if COL <= colMid {
			return t.get((root<<2)-2, ROW, COL, row1, col1, rowMid, colMid)
		}
		return t.get((root<<2)-1, ROW, COL, row1, colMid+1, rowMid, col2)
	} else {
		if COL <= colMid {
			return t.get((root << 2), ROW, COL, rowMid+1, col1, row2, colMid)
		}
		return t.get((root<<2)+1, ROW, COL, rowMid+1, colMid+1, row2, col2)
	}
}

func (t *SegmentTree2D) set(root int, ROW, COL, row1, col1, row2, col2 int, e E) {
	if row1 == row2 && col1 == col2 {
		t.data[root] = e
		return
	}

	t.pushDown(root)
	rowMid, colMid := (row1+row2)>>1, (col1+col2)>>1
	if ROW <= rowMid {
		if COL <= colMid {
			t.set((root<<2)-2, ROW, COL, row1, col1, rowMid, colMid, e)
		} else {
			t.set((root<<2)-1, ROW, COL, row1, colMid+1, rowMid, col2, e)
		}
	} else {
		if COL <= colMid {
			t.set((root << 2), ROW, COL, rowMid+1, col1, row2, colMid, e)
		} else {
			t.set((root<<2)+1, ROW, COL, rowMid+1, colMid+1, row2, col2, e)
		}
	}

	t.pushUp(root)
}

func (t *SegmentTree2D) pushUp(root int) {
	t.data[root] = t.op(t.op(t.data[(root<<2)-2], t.data[(root<<2)-1]), t.op(t.data[root<<2], t.data[(root<<2)+1]))
}

func (t *SegmentTree2D) pushDown(root int) {
	if t.lazy[root] != t.id() {
		t.propagate((root<<2)-2, t.lazy[root])
		t.propagate((root<<2)-1, t.lazy[root])
		t.propagate(root<<2, t.lazy[root])
		t.propagate((root<<2)+1, t.lazy[root])
		t.lazy[root] = t.id()
	}
}

func (t *SegmentTree2D) propagate(root int, f Id) {
	t.data[root] = t.mapping(f, t.data[root])
	// !叶子结点不需要更新lazy标记
	if root < len(t.lazy) {
		t.lazy[root] = t.composition(f, t.lazy[root])
	}
}

func (t *SegmentTree2D) QueryAll() E { return t.data[1] }

// !查询闭区间[row1, col1, row2, col2]的值
//  0<=row1<=row2<ROW, 0<=col1<=col2<COL
func (t *SegmentTree2D) Query(row1, col1, row2, col2 int) E {
	return t.query(1, row1+1, col1+1, row2+1, col2+1, 1, 1, t.row, t.col)
}

// !更新闭区间[row1, col1, row2, col2]的值
//  0<=row1<=row2<ROW, 0<=col1<=col2<COL
func (t *SegmentTree2D) Update(row1, col1, row2, col2 int, lazy Id) {
	t.update(1, row1+1, col1+1, row2+1, col2+1, 1, 1, t.row, t.col, lazy)
}

// !单点查询
//  0<=row<ROW, 0<=col<COL
func (t *SegmentTree2D) Get(row, col int) E {
	return t.get(1, row+1, col+1, 1, 1, t.row, t.col)
}

// !单点赋值
//  0<=row<ROW, 0<=col<COL
func (t *SegmentTree2D) Set(row, col int, e E) {
	t.set(1, row+1, col+1, 1, 1, t.row, t.col, e)
}

func main() {
	// ["SubrectangleQueries","getValue","updateSubrectangle","getValue","getValue","updateSubrectangle","getValue"]
	// [[[[1,1,1],[2,2,2],[3,3,3]]],[0,0],[0,0,2,2,100],[0,0],[2,2],[1,1,2,2,20],[2,2]]
	sub := Constructor([][]int{{1, 1, 1}, {2, 2, 2}, {3, 3, 3}})
	fmt.Println(sub.GetValue(0, 0))
	sub.UpdateSubrectangle(0, 0, 2, 2, 100)
	fmt.Println(sub.GetValue(0, 0))
	fmt.Println(sub.GetValue(2, 2))
	sub.UpdateSubrectangle(1, 1, 2, 2, 20)
	fmt.Println(sub.GetValue(2, 2))
}
