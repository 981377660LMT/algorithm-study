// 树套树(外层为线段树),支持区间更新,单点查询.
// 保证内层树的col不大于row.

package main

import (
	"fmt"
	"time"
)

func main() {
	// 区间赋值,单点查询
	row, col := 500, 500
	seg := NewSegmentTree2DRangeUpdatePointGet(
		row, col,
		func(n int) IRangeUpdatePointGet1D { return NewNavieTree(n) },
		func(a E, b E) E {
			if a.time > b.time {
				return a
			}
			return b
		},
	)

	time1 := time.Now()
	for i := 0; i < 1e5; i++ {
		seg.Update(0, row, 0, col, E{i, i})
		seg.Get(0, 0)
		seg.Set(0, 0, E{0, 0})
	}
	time2 := time.Now()
	fmt.Println(time2.Sub(time1)) // 214.725792ms
}

const INF int = 1e18

// RangeUpdatePointGet1D

type E = struct{ time, value int }
type Id = E
type NaiveTree struct {
	data []E
}

func NewNavieTree(n int) IRangeUpdatePointGet1D {
	data := make([]E, n)
	for i := range data {
		data[i] = E{time: -1} // 初始更新时间为-1,表示没有更新
	}
	return &NaiveTree{data: data}
}

func (nat *NaiveTree) Update(start, end int, lazy Id) {
	for i := start; i < end; i++ {
		nat.data[i] = lazy
	}
}
func (nat *NaiveTree) Get(pos int) E        { return nat.data[pos] }
func (nat *NaiveTree) Set(pos int, value E) { nat.data[pos] = value }

// #region template SegmentTree2DRangeUpdatePointGet
type IRangeUpdatePointGet1D interface {
	Update(start int, end int, lazy Id)
	Get(index int) E
	Set(index int, value E)
}

type SegmentTree2DRangeUpdatePointGet struct {
	_seg        []IRangeUpdatePointGet1D
	_mergeRow   func(a E, b E) E
	_init1D     func() IRangeUpdatePointGet1D
	_needRotate bool
	_rawRow     int
	_size       int
}

// 二维区间更新，单点查询的线段树(树套树).
//
//	row 行数.对时间复杂度贡献为`O(log(row))`.
//	col 列数.内部树的大小为.列数越小,对内部树的时间复杂度要求越低.
//	createRangeUpdatePointGet1D 初始化内层"树"的函数.入参为内层"树"的大小.
//	mergeRow 合并两个内层"树"的结果.
func NewSegmentTree2DRangeUpdatePointGet(
	row, col int,
	createRangeUpdatePointGet1D func(n int) IRangeUpdatePointGet1D,
	mergeRow func(a E, b E) E,
) *SegmentTree2DRangeUpdatePointGet {
	res := &SegmentTree2DRangeUpdatePointGet{}
	res._rawRow = row
	res._needRotate = row < col
	if res._needRotate {
		row ^= col
		col ^= row
		row ^= col
	}

	size := 1
	for size < row {
		size <<= 1
	}
	res._seg = make([]IRangeUpdatePointGet1D, 2*size-1)
	res._mergeRow = mergeRow
	res._init1D = func() IRangeUpdatePointGet1D {
		return createRangeUpdatePointGet1D(col)
	}
	res._size = size
	return res
}

// 将`[row1,row2)`x`[col1,col2)`的区间值与`lazy`作用.
func (this *SegmentTree2DRangeUpdatePointGet) Update(row1, row2, col1, col2 int, lazy Id) {
	if this._needRotate {
		tmp1 := row1
		tmp2 := row2
		row1 = col1
		row2 = col2
		col1 = this._rawRow - tmp2
		col2 = this._rawRow - tmp1
	}

	this._update(row1, row2, col1, col2, lazy, 0, 0, this._size)
}

func (this *SegmentTree2DRangeUpdatePointGet) Get(row, col int) E {
	if this._needRotate {
		tmp := row
		row = col
		col = this._rawRow - tmp - 1
	}

	row += this._size - 1
	if this._seg[row] == nil {
		this._seg[row] = this._init1D()
	}
	res := this._seg[row].Get(col)
	for row > 0 {
		row = (row - 1) >> 1
		if this._seg[row] != nil {
			res = this._mergeRow(res, this._seg[row].Get(col))
		}
	}
	return res
}

func (this *SegmentTree2DRangeUpdatePointGet) Set(row, col int, value E) {
	if this._needRotate {
		tmp := row
		row = col
		col = this._rawRow - tmp - 1
	}

	row += this._size - 1
	if this._seg[row] == nil {
		this._seg[row] = this._init1D()
	}
	this._seg[row].Set(col, value)
	for row > 0 {
		row = (row - 1) >> 1
		if this._seg[row] == nil {
			this._seg[row] = this._init1D()
		}
		this._seg[row].Set(col, value)
	}
}

func (this *SegmentTree2DRangeUpdatePointGet) _update(R, C, start, end int, lazy Id, pos, r, c int) {
	if c <= R || C <= r {
		return
	}
	if R <= r && c <= C {
		if this._seg[pos] == nil {
			this._seg[pos] = this._init1D()
		}
		this._seg[pos].Update(start, end, lazy)
	} else {
		mid := (r + c) >> 1
		this._update(R, C, start, end, lazy, 2*pos+1, r, mid)
		this._update(R, C, start, end, lazy, 2*pos+2, mid, c)
	}
}

// #endregion
//
//
//1476. 子矩形查询

type SubrectangleQueries struct {
	seg  *SegmentTree2DRangeUpdatePointGet
	time int
}

func Constructor(rectangle [][]int) SubrectangleQueries {
	res := SubrectangleQueries{}
	row, col := len(rectangle), len(rectangle[0])
	seg := NewSegmentTree2DRangeUpdatePointGet(
		row, col,
		func(n int) IRangeUpdatePointGet1D { return NewNavieTree(n) },
		func(a E, b E) E {
			if a.time > b.time {
				return a
			}
			return b
		})
	res.seg = seg
	res.time = 1

	for i, row := range rectangle {
		for j, val := range row {
			res.UpdateSubrectangle(i, j, i, j, val)
		}
	}

	return res
}

func (this *SubrectangleQueries) UpdateSubrectangle(row1 int, col1 int, row2 int, col2 int, newValue int) {
	this.seg.Update(row1, row2+1, col1, col2+1, E{time: this.time, value: newValue})
	this.time++
}

func (this *SubrectangleQueries) GetValue(row int, col int) int {
	return this.seg.Get(row, col).value
}

/**
 * Your SubrectangleQueries object will be instantiated and called as such:
 * obj := Constructor(rectangle);
 * obj.UpdateSubrectangle(row1,col1,row2,col2,newValue);
 * param_2 := obj.GetValue(row,col);
 */
