// 二维线段树

package main

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
	row, col := len(rectangle), len(rectangle[0])
	tree := NewLazySegmentTree2D(0, 0, row-1, col-1)
	for i := 0; i < row; i++ {
		for j := 0; j < col; j++ {
			tree.Set(i, j, rectangle[i][j])
		}
	}

	return SubrectangleQueries{tree}
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

func e(row1, col1, row2, col2 int) E { return INF }
func id() Id                         { return INF }
func op(e1, e2 E) E {
	if e1 == INF {
		return e2
	}
	return e1
}
func mapping(f Id, e E) E {
	if f == INF {
		return e
	}
	return f
}
func composition(f, g Id) Id {
	if f == INF {
		return g
	}
	return f
}

//
//
//
type Node struct {
	child1, child2, child3, child4 *Node
	data                           E
	lazy                           Id
}

type SegmentTree2D struct {
	row1, col1, row2, col2 int
	root                   *Node
}

// !指定区间上下界建立线段树
//  0<=row1<=row2, 0<=col1<=col2
func NewLazySegmentTree2D(row1, col1, row2, col2 int) *SegmentTree2D {
	root := newNode(row1, col1, row2, col2)
	return &SegmentTree2D{row1, col1, row2, col2, root}
}

func (t *SegmentTree2D) query(root *Node, ROW1, COL1, ROW2, COL2, row1, col1, row2, col2 int) E {
	if ROW1 <= row1 && row2 <= ROW2 && COL1 <= col1 && col2 <= COL2 {
		return root.data
	}

	rowMid, colMid := (row1+row2)>>1, (col1+col2)>>1
	t.pushDown(root, row1, col1, row2, col2, rowMid, colMid)
	res := e(row1, col1, row2, col2)
	if ROW1 <= rowMid {
		if COL1 <= colMid {
			res = op(res, t.query(root.child1, ROW1, COL1, ROW2, COL2, row1, col1, rowMid, colMid))
		}
		if colMid < COL2 {
			res = op(res, t.query(root.child2, ROW1, COL1, ROW2, COL2, row1, colMid+1, rowMid, col2))
		}
	}
	if rowMid < ROW2 {
		if COL1 <= colMid {
			res = op(res, t.query(root.child3, ROW1, COL1, ROW2, COL2, rowMid+1, col1, row2, colMid))
		}
		if colMid < COL2 {
			res = op(res, t.query(root.child4, ROW1, COL1, ROW2, COL2, rowMid+1, colMid+1, row2, col2))
		}
	}

	return res
}

func (t *SegmentTree2D) update(root *Node, ROW1, COL1, ROW2, COL2, row1, col1, row2, col2 int, lazy Id) {
	if ROW1 <= row1 && row2 <= ROW2 && COL1 <= col1 && col2 <= COL2 {
		root.data = mapping(lazy, root.data)
		root.lazy = composition(lazy, root.lazy)
		return
	}

	rowMid, colMid := (row1+row2)>>1, (col1+col2)>>1
	t.pushDown(root, row1, col1, row2, col2, rowMid, colMid)
	if ROW1 <= rowMid {
		if COL1 <= colMid {
			t.update(root.child1, ROW1, COL1, ROW2, COL2, row1, col1, rowMid, colMid, lazy)
		}
		if colMid < COL2 {
			t.update(root.child2, ROW1, COL1, ROW2, COL2, row1, colMid+1, rowMid, col2, lazy)
		}
	}
	if rowMid < ROW2 {
		if COL1 <= colMid {
			t.update(root.child3, ROW1, COL1, ROW2, COL2, rowMid+1, col1, row2, colMid, lazy)
		}
		if colMid < COL2 {
			t.update(root.child4, ROW1, COL1, ROW2, COL2, rowMid+1, colMid+1, row2, col2, lazy)
		}
	}

	t.pushUp(root)
}

func (t *SegmentTree2D) get(root *Node, ROW, COL, row1, col1, row2, col2 int) E {
	if row1 == row2 && col1 == col2 {
		return root.data
	}

	rowMid, colMid := (row1+row2)>>1, (col1+col2)>>1
	t.pushDown(root, row1, col1, row2, col2, rowMid, colMid)
	if ROW <= rowMid {
		if COL <= colMid {
			return t.get(root.child1, ROW, COL, row1, col1, rowMid, colMid)
		}
		return t.get(root.child2, ROW, COL, row1, colMid+1, rowMid, col2)
	} else {
		if COL <= colMid {
			return t.get(root.child3, ROW, COL, rowMid+1, col1, row2, colMid)
		}
		return t.get(root.child4, ROW, COL, rowMid+1, colMid+1, row2, col2)
	}
}

func (t *SegmentTree2D) set(root *Node, ROW, COL, row1, col1, row2, col2 int, e E) {
	if row1 == row2 && col1 == col2 {
		root.data = e
		return
	}

	rowMid, colMid := (row1+row2)>>1, (col1+col2)>>1
	t.pushDown(root, row1, col1, row2, col2, rowMid, colMid)

	if ROW <= rowMid {
		if COL <= colMid {
			t.set(root.child1, ROW, COL, row1, col1, rowMid, colMid, e)
		} else {
			t.set(root.child2, ROW, COL, row1, colMid+1, rowMid, col2, e)
		}
	} else {
		if COL <= colMid {
			t.set(root.child3, ROW, COL, rowMid+1, col1, row2, colMid, e)
		} else {
			t.set(root.child4, ROW, COL, rowMid+1, colMid+1, row2, col2, e)
		}
	}

	t.pushUp(root)
}

func newNode(row1, col1, row2, col2 int) *Node {
	return &Node{data: e(row1, col1, row2, col2), lazy: id()}
}

// op
func (t *SegmentTree2D) pushUp(root *Node) {
	root.data = op(op(root.child1.data, root.child2.data), op(root.child3.data, root.child4.data))
}

func (t *SegmentTree2D) pushDown(root *Node, row1, col1, row2, col2, rowMid, colMid int) {
	if root.child1 == nil {
		root.child1 = newNode(row1, col1, rowMid, colMid)
	}
	if root.child2 == nil {
		root.child2 = newNode(row1, colMid+1, rowMid, col2)
	}
	if root.child3 == nil {
		root.child3 = newNode(rowMid+1, col1, row2, colMid)
	}
	if root.child4 == nil {
		root.child4 = newNode(rowMid+1, colMid+1, row2, col2)
	}

	if root.lazy != id() {
		// propagate  mapping+composition
		root.child1.data = mapping(root.lazy, root.child1.data)
		root.child1.lazy = composition(root.lazy, root.child1.lazy)
		root.child2.data = mapping(root.lazy, root.child2.data)
		root.child2.lazy = composition(root.lazy, root.child2.lazy)
		root.child3.data = mapping(root.lazy, root.child3.data)
		root.child3.lazy = composition(root.lazy, root.child3.lazy)
		root.child4.data = mapping(root.lazy, root.child4.data)
		root.child4.lazy = composition(root.lazy, root.child4.lazy)
		root.lazy = id()
	}
}

// !查询闭区间[row1, col1, row2, col2]的值
//  ROW1<=row1<=row2<=ROW2, COL1<=col1<=col2<=COL2
func (t *SegmentTree2D) Query(row1, col1, row2, col2 int) E {
	return t.query(t.root, row1, col1, row2, col2, t.row1, t.col1, t.row2, t.col2)
}

// !更新闭区间[row1, col1, row2, col2]的值
//  ROW1<=row1<=row2<=ROW2, COL1<=col1<=col2<=COL2
func (t *SegmentTree2D) Update(row1, col1, row2, col2 int, lazy Id) {
	t.update(t.root, row1, col1, row2, col2, t.row1, t.col1, t.row2, t.col2, lazy)
}

func (t *SegmentTree2D) QueryAll() E { return t.root.data }

// !单点查询
func (t *SegmentTree2D) Get(row, col int) E {
	return t.get(t.root, row, col, t.row1, t.col1, t.row2, t.col2)
}

// !单点赋值
func (t *SegmentTree2D) Set(row, col int, e E) {
	t.set(t.root, row, col, t.row1, t.col1, t.row2, t.col2, e)
}

// func main() {
// 	// ["SubrectangleQueries","getValue","updateSubrectangle","getValue","getValue","updateSubrectangle","getValue"]
// 	// [[[[1,1,1],[2,2,2],[3,3,3]]],[0,0],[0,0,2,2,100],[0,0],[2,2],[1,1,2,2,20],[2,2]]
// 	sub := Constructor([][]int{{1, 1, 1}, {2, 2, 2}, {3, 3, 3}})
// 	fmt.Println(sub.GetValue(0, 0))
// 	sub.UpdateSubrectangle(0, 0, 2, 2, 100)
// 	fmt.Println(sub.GetValue(0, 0))
// 	fmt.Println(sub.GetValue(2, 2))
// 	sub.UpdateSubrectangle(1, 1, 2, 2, 20)
// 	fmt.Println(sub.GetValue(2, 2))
// }
