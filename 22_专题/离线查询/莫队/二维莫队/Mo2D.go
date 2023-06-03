// 二维莫队
// 时间复杂度O(row*col*(q**0.75))
// !一般 n,m<=200 q<=1e5

package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
)

// https://oi-wiki.org/misc/mo-algo-2dimen/#%E5%9D%97%E9%95%BF%E9%80%89%E5%AE%9A
// https://hydro.ac/d/bzoj/p/2639
// 查询子矩阵的`权值`.
// !这里的`权值`是这样定义的:对于一个整数 x，如果它在该矩阵中出现了 p 次，那么它给该矩阵的权值就贡献 p^2。
// 1<=n,m<=200,1<=q<=1e5,1<=a[i][j]<=2e9
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var ROW, COL int
	fmt.Fscan(in, &ROW, &COL)
	grid := make([][]int, ROW)
	for i := 0; i < ROW; i++ {
		grid[i] = make([]int, COL)
		for j := 0; j < COL; j++ {
			var v int
			fmt.Fscan(in, &v)
			grid[i][j] = v
		}
	}

	var q int
	fmt.Fscan(in, &q)
	mo2d := NewMo2D(ROW, COL, q)
	for i := 0; i < q; i++ {
		var x1, y1, x2, y2 int // >=1
		fmt.Fscan(in, &x1, &y1, &x2, &y2)
		if x1 > x2 {
			x1, x2 = x2, x1
		}
		if y1 > y2 {
			y1, y2 = y2, y1
		}
		mo2d.AddQuery(x1-1, x2, y1-1, y2)
	}

	// 离散化
	_pool := make(map[interface{}]int)
	id := func(o interface{}) int {
		if v, ok := _pool[o]; ok {
			return v
		}
		v := len(_pool)
		_pool[o] = v
		return v
	}

	newGrid := make([]int, ROW*COL)
	for i := 0; i < ROW; i++ {
		for j := 0; j < COL; j++ {
			newGrid[i*COL+j] = id(grid[i][j])
		}
	}

	cur := 0
	res := make([]int, q)
	counter := make([]int, len(_pool))
	addRow := func(row int, dir int, col1, col2 int) {
		for c := col1; c < col2; c++ {
			x := newGrid[row*COL+c]
			pre := counter[x]
			cur -= (pre * pre)
			counter[x] += 1
			cur += (pre + 1) * (pre + 1)
		}
	}
	addCol := func(col int, dir int, row1, row2 int) {
		for r := row1; r < row2; r++ {
			x := newGrid[r*COL+col]
			pre := counter[x]
			cur -= (pre * pre)
			counter[x] += 1
			cur += (pre + 1) * (pre + 1)
		}
	}
	removeRow := func(row int, dir int, col1, col2 int) {
		for c := col1; c < col2; c++ {
			x := newGrid[row*COL+c]
			pre := counter[x]
			cur -= (pre * pre)
			counter[x] -= 1
			cur += (pre - 1) * (pre - 1)
		}
	}
	removeCol := func(col int, dir int, row1, row2 int) {
		for r := row1; r < row2; r++ {
			x := newGrid[r*COL+col]
			pre := counter[x]
			cur -= (pre * pre)
			counter[x] -= 1
			cur += (pre - 1) * (pre - 1)
		}
	}
	query := func(qid int) {
		res[qid] = cur
	}
	mo2d.Run(addRow, addCol, removeRow, removeCol, query)

	for _, v := range res {
		fmt.Fprintln(out, v)
	}
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
	blockId := make([]int, max(row, col)+5)
	for i := range blockId {
		blockId[i] = i / chunkSize
	}
	return &Mo2D{chunkSize: chunkSize, blockId: blockId}
}

// 添加查询矩形区域`[x1, x2) * [y1, y2)`.
//  0 <= x1 < x2 <= row, 0 <= y1 < y2 <= col.
func (mo *Mo2D) AddQuery(x1, x2, y1, y2 int) {
	x2, y2 = x2-1, y2-1
	mo.queries = append(mo.queries, query{x1: x1, y1: y1, x2: x2, y2: y2, qid: len(mo.queries)})
}

// 返回每个查询的结果.
//  addRow: 将新的行添加到窗口. dir: 1 表示row变大，-1 表示row变小. 对应列的范围是[col1, col2).
//  addCol: 将新的列添加到窗口. dir: 1 表示col变大，-1 表示col变小. 对应行的范围是[row1, row2).
// 	removeRow: 将行从窗口移除. dir: 1 表示row变大，-1 表示row变小. 对应列的范围是[col1, col2).
// 	removeCol: 将列从窗口移除. dir: 1 表示col变大，-1 表示col变小. 对应行的范围是[row1, row2).
//  query: 查询窗口内的数据.
func (mo *Mo2D) Run(
	addRow func(row int, dir int, col1, col2 int),
	addCol func(col int, dir int, row1, row2 int),
	removeRow func(row int, dir int, col1, col2 int),
	removeCol func(col int, dir int, row1, row2 int),
	query func(qid int),
) {
	bid, queries := mo.blockId, mo.queries
	sort.Slice(queries, func(i, j int) bool {
		q1, q2 := queries[i], queries[j]
		bid1, bid2 := bid[q1.x1], bid[q2.x1]
		if bid1 != bid2 {
			return bid1 < bid2
		}
		bid3, bid4 := bid[q1.y1], bid[q2.y1]
		if bid3 != bid4 {
			if bid1&1 == 1 {
				return q1.y1 < q2.y1
			}
			return q1.y1 > q2.y1
		}
		bid5, bid6 := bid[q1.y2], bid[q2.y2]
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

	x1, y1, x2, y2 := 0, 0, -1, -1
	for _, q := range queries {
		qx1, qy1, qx2, qy2, qid := q.x1, q.y1, q.x2, q.y2, q.qid
		for x1 > qx1 {
			x1--
			addRow(x1, -1, y1, y2+1)
		}
		for x2 < qx2 {
			x2++
			addRow(x2, 1, y1, y2+1)
		}
		for y1 > qy1 {
			y1--
			addCol(y1, -1, x1, x2+1)
		}
		for y2 < qy2 {
			y2++
			addCol(y2, 1, x1, x2+1)
		}
		for x1 < qx1 {
			removeRow(x1, 1, y1, y2+1)
			x1++
		}
		for x2 > qx2 {
			removeRow(x2, -1, y1, y2+1)
			x2--
		}
		for y1 < qy1 {
			removeCol(y1, 1, x1, x2+1)
			y1++
		}
		for y2 > qy2 {
			removeCol(y2, -1, x1, x2+1)
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
