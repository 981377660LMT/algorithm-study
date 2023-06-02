// 二维莫队
package main

import (
	"fmt"
	"math"
	"sort"
)

func main() {
	// https://oi-wiki.org/misc/mo-algo-2dimen/#%E5%9D%97%E9%95%BF%E9%80%89%E5%AE%9A
	// https://hydro.ac/d/bzoj/p/2639
	mo2d := NewMo2D(3, 3, 3)
	mo2d.AddQuery(0, 0, 2, 2)
	mo2d.AddQuery(0, 0, 1, 1)
	mo2d.AddQuery(1, 1, 2, 2)
	addRow := func(row int, delta int, col1, col2 int) {
		fmt.Println("addRow", row, delta, col1, col2)
	}
	addCol := func(col int, delta int, row1, row2 int) {
		fmt.Println("addCol", col, delta, row1, row2)
	}
	removeRow := func(row int, delta int, col1, col2 int) {
		fmt.Println("removeRow", row, delta, col1, col2)
	}
	removeCol := func(col int, delta int, row1, row2 int) {
		fmt.Println("removeCol", col, delta, row1, row2)
	}
	query := func(qid int) {
		fmt.Println("query", qid)
	}
	mo2d.Run(addRow, addCol, removeRow, removeCol, query)

}

type query struct{ x1, y1, x2, y2, qid int }

type Mo2D struct {
	chunkSize int
	queries   []query
	blockId   []int
}

func NewMo2D(row, col int, q int) *Mo2D {
	chunkSize := int(math.Pow(float64(row*col), 0.5) / math.Pow(float64(q), 0.25))
	if chunkSize < 1 {
		chunkSize = 1
	}
	return &Mo2D{chunkSize: chunkSize}
}

func (mo *Mo2D) AddQuery(x1, y1, x2, y2 int) {
	mo.queries = append(mo.queries, query{x1: x1, y1: y1, x2: x2, y2: y2, qid: len(mo.queries)})
}

// 返回每个查询的结果.
//  addRow: 将新的行添加到窗口. delta: 1 表示row变大，-1 表示row变小.
//  addCol: 将新的列添加到窗口. delta: 1 表示col变大，-1 表示col变小.
// 	removeRow: 将行从窗口移除. delta: 1 表示row变大，-1 表示row变小.
// 	removeCol: 将列从窗口移除. delta: 1 表示col变大，-1 表示col变小.
//  query: 查询窗口内的数据.
func (mo *Mo2D) Run(
	addRow func(row int, delta int, col1, col2 int),
	addCol func(col int, delta int, row1, row2 int),
	removeRow func(row int, delta int, col1, col2 int),
	removeCol func(col int, delta int, row1, row2 int),
	query func(qid int),
) {
	chunkSize, queries := mo.chunkSize, mo.queries
	sort.Slice(queries, func(i, j int) bool {
		q1, q2 := queries[i], queries[j]
		bid1, bid2 := q1.x1/chunkSize, q2.x1/chunkSize
		if bid1 != bid2 {
			return bid1 < bid2
		}
		bid3, bid4 := q1.y1/chunkSize, q2.y1/chunkSize
		if bid3 != bid4 {
			if bid1&1 == 1 {
				return q1.y1 < q2.y1
			}
			return q1.y1 > q2.y1
		}
		bid5, bid6 := q1.y2/chunkSize, q2.y2/chunkSize
		if bid5 != bid6 {
			if bid3&1 == 1 {
				return q1.y2 < q2.y2
			}
			return q1.y2 > q2.y2
		}
		if bid5&1 == 1 {
			return q1.x2 < q2.x2
		}
		return q1.x2 > q2.x2
	})

	x1, y1, x2, y2 := 0, 0, 0, 0
	for _, q := range queries {
		qx1, qy1, qx2, qy2, qid := q.x1, q.y1, q.x2, q.y2, q.qid
		for x1 > qx1 {
			x1--
			addRow(x1, -1, y1, y2)
		}
		for x2 < qx2 {
			x2++
			addRow(x2, 1, y1, y2)
		}
		for y1 > qy1 {
			y1--
			addCol(y1, -1, x1, x2)
		}
		for y2 < qy2 {
			y2++
			addCol(y2, 1, x1, x2)
		}
		for x1 < qx1 {
			removeRow(x1, 1, y1, y2)
			x1++
		}
		for x2 > qx2 {
			removeRow(x2, -1, y1, y2)
			x2--
		}
		for y1 < qy1 {
			removeCol(y1, 1, x1, x2)
			y1++
		}
		for y2 > qy2 {
			removeCol(y2, -1, x1, x2)
			y2--
		}
		query(qid)
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
