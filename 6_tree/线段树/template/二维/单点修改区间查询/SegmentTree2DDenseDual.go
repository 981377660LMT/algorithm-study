package main

import "fmt"

func main() {
	seg := NewSegmentTree2DDenseDual(3, 4)
	seg.Update(1, 2, 1, 3, 1)
	fmt.Println(seg.GetAll())
}

// !TODO: 答案不对
// 2536. 子矩阵元素加 1
// https://leetcode.cn/problems/increment-submatrices-by-one/description/
func rangeAddQueries(n int, queries [][]int) [][]int {
	row, col := int32(n), int32(n)
	seg := NewSegmentTree2DDenseDual(row, col)
	for _, q := range queries {
		r1, c1, r2, c2 := int32(q[0]), int32(q[1]), int32(q[2]), int32(q[3])
		seg.Update(r1, r2+1, c1, c2+1, 1)
	}
	return seg.GetAll()
}

const INF int = 1e18

type E = int

func e() E        { return 0 }
func op(a, b E) E { return a + b }

type SegmentTree2DDenseDual struct {
	row, col int32
	data     []E
}

func NewSegmentTree2DDenseDual(row, col int32) *SegmentTree2DDenseDual {
	data := make([]E, 4*row*col)
	for i := range data {
		data[i] = e()
	}
	return &SegmentTree2DDenseDual{row: row, col: col, data: data}
}

func (st *SegmentTree2DDenseDual) Get(row, col int32) E {
	res := e()
	curRow := row + st.row
	for curRow > 0 {
		curCol := col + st.col
		for curCol > 0 {
			res = op(res, st.data[st.idx(curRow, curCol)])
			curCol >>= 1
		}
		curRow >>= 1
	}
	return res
}

func (st *SegmentTree2DDenseDual) GetAll() [][]E {
	for r := int32(1); r < st.row; r++ {
		for c := int32(0); c < st.col; c++ {
			st.pushDown(r, c)
		}
	}
	for r := int32(1); r < st.row; r++ {
		for c := int32(0); c < 2*st.col; c++ {
			cur, left, right := st.idx(r, c), st.idx(r<<1, c), st.idx(r<<1|1, c)
			st.data[left] = op(st.data[left], st.data[cur])
			st.data[right] = op(st.data[right], st.data[cur])
			st.data[cur] = e()
		}
	}
	res := make([][]E, st.row)
	for r := int32(0); r < st.row; r++ {
		res[r] = make([]E, st.col)
		for c := int32(0); c < st.col; c++ {
			res[r][c] = st.data[st.idx(r+st.row, c+st.col)]
		}
	}
	return res
}

func (st *SegmentTree2DDenseDual) Update(rowStart, rowEnd, colStart, colEnd int32, x E) {
	rowStart += st.row
	rowEnd += st.row
	for rowStart < rowEnd {
		if rowStart&1 == 1 {
			st.update(rowStart, colStart, colEnd, x)
			rowStart++
		}
		if rowEnd&1 == 1 {
			rowEnd--
			st.update(rowEnd, colStart, colEnd, x)
		}
		rowStart >>= 1
		rowEnd >>= 1
	}
}

func (st *SegmentTree2DDenseDual) idx(r, c int32) int32 {
	return r*st.col<<1 + c
}

func (st *SegmentTree2DDenseDual) update(r, cl, cr int32, e E) {
	cl += st.col
	cr += st.col
	for cl < cr {
		if cl&1 == 1 {
			id := st.idx(r, cl)
			st.data[id] = op(st.data[id], e)
			cl++
		}
		if cr&1 == 1 {
			cr--
			id := st.idx(r, cr)
			st.data[id] = op(st.data[id], e)
		}
		cl >>= 1
		cr >>= 1
	}
}

func (st *SegmentTree2DDenseDual) pushDown(x, k int32) {
	cur, left, right := st.idx(x, k), st.idx(x, k<<1), st.idx(x, k<<1|1)
	st.data[left] = op(st.data[left], st.data[cur])
	st.data[right] = op(st.data[right], st.data[cur])
	st.data[cur] = e()
}
