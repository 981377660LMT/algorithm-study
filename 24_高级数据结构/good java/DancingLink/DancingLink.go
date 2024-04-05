// 用于解决精准覆盖问题.

// https://zhuanlan.zhihu.com/p/617477033
// api:
//	NewSparseInstance(points [][2]int32, row, col int32)
//	NewDenseInstance(adjMatrix [][]bool, row, col int32)
//	GetSolution() []int32

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	P4929()
}

func demo() {
	// 0 1 2
	// 3 4 5
	// 0 1 2 4 5
	// 0 1
	// 2
	grid := [][]bool{
		{true, true, true, false, false, false},
		{false, false, false, true, true, true},
		{true, true, false, true, true, false},
		{true, true, false, false, false, false},
		{false, false, true, false, false, false},
	}
	res := NewDancingLinkDense(grid, 5, 6)
	fmt.Println(res.GetSolution()) // [1,3,4]
}

// P4929 【模板】舞蹈链（DLX）精准覆盖问题
// 给定一个N 行M 列的矩阵，
// 矩阵中每个元素要么是1，要么是0。
// 你需要在矩阵中挑选出若干行，使得对于矩阵的每一列j，
// 在你挑选的这些行中，有且仅有一行的第j 个元素为1。
func P4929() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var row, col int32
	fmt.Fscan(in, &row, &col)
	cover := make([][]bool, row)
	for i := int32(0); i < row; i++ {
		cover[i] = make([]bool, col)
		for j := int32(0); j < col; j++ {
			var x int
			fmt.Fscan(in, &x)
			cover[i][j] = x == 1
		}
	}

	res := NewDancingLinkDense(cover, row, col).GetSolution()
	if res == nil {
		fmt.Fprintln(out, "No Solution!")
	} else {
		for _, v := range res {
			fmt.Fprint(out, v+1, " ")
		}
	}
}

// https://leetcode.cn/problems/sudoku-solver/
func solveSudoku(board [][]byte) {
	row, col := 9, 9
	grid := make([][]int32, row)
	for i := 0; i < row; i++ {
		grid[i] = make([]int32, col)
	}
	for i := 0; i < row; i++ {
		for j := 0; j < col; j++ {
			if board[i][j] != '.' {
				grid[i][j] = int32(board[i][j] - '0')
			} else {
				grid[i][j] = 0
			}
		}
	}

	res := SolveSudoku(grid, func(i, j int32) bool { return grid[i][j] == 0 })

	for i := 0; i < row; i++ {
		for j := 0; j < col; j++ {
			board[i][j] = byte(res[i][j] + '0')
		}
	}
}

type DancingLink struct {
	res      []int32
	rowDummy []*DNode
	colDummy []*DNode
	colHead  *DNode
	rowSize  []int32
	colSize  []int32
	stack    []*DNode
	resStack []*DNode
}

func NewDancingLinkSparse(points [][2]int32, row, col int32) *DancingLink {
	res := &DancingLink{}
	res.init(row, col)
	points = append(points[:0:0], points...)
	sort.Slice(points, func(i, j int) bool {
		if points[i][0] == points[j][0] {
			return points[i][1] < points[j][1]
		}
		return points[i][0] < points[j][0]
	})
	for _, xy := range points {
		res.add(xy[0], xy[1])
	}
	res.stack = make([]*DNode, 0, len(points)+int(row)+int(col))
	res.dance()
	return res
}

// adjMatrix[i][j] 表示第i行是否包含元素j.
func NewDancingLinkDense(adjMatrix [][]bool, row, col int32) *DancingLink {
	res := &DancingLink{}
	res.init(row, col)
	count := 0
	for i := int32(0); i < row; i++ {
		for j := int32(0); j < col; j++ {
			if adjMatrix[i][j] {
				count++
				res.add(i, j)
			}
		}
	}
	res.stack = make([]*DNode, 0, count+int(row)+int(col))
	res.dance()
	return res
}

// nil 表示无解.
func (dl *DancingLink) GetSolution() []int32 {
	return dl.res
}

func (dl *DancingLink) newNode(row, col int32) *DNode {
	return newNode(row, col)
}

func (dl *DancingLink) check(node, skip *DNode) bool {
	return node.l.r == node && node.r.l == node &&
		node.d.u == node && node.u.d == node || node == skip
}

func (dl *DancingLink) add(i, j int32) {
	node := dl.newNode(i, j)
	node.l = dl.rowDummy[i].l
	dl.rowDummy[i].l.r = node
	node.r = dl.rowDummy[i]
	dl.rowDummy[i].l = node

	node.u = dl.colDummy[j].u
	dl.colDummy[j].u.d = node
	node.d = dl.colDummy[j]
	dl.colDummy[j].u = node

	dl.rowSize[i]++
	dl.colSize[j]++
}

func (dl *DancingLink) remove(node *DNode) {
	node.l.r = node.r
	node.r.l = node.l
	node.u.d = node.d
	node.d.u = node.u

	if node.row >= 0 {
		dl.rowSize[node.row]--
	}
	if node.col >= 0 {
		dl.colSize[node.col]--
	}
	node.deleted = true
}

func (dl *DancingLink) recover(node *DNode) {
	node.l.r = node
	node.r.l = node
	node.u.d = node
	node.d.u = node

	if node.row >= 0 {
		dl.rowSize[node.row]++
	}
	if node.col >= 0 {
		dl.colSize[node.col]++
	}
	node.deleted = false
}

func (dl *DancingLink) init(n, m int32) {
	dl.colSize = make([]int32, m)
	dl.rowSize = make([]int32, n)
	dl.rowDummy = make([]*DNode, n)
	dl.colDummy = make([]*DNode, m)
	for i := int32(0); i < n; i++ {
		dl.rowDummy[i] = newNode(i, -1)
	}
	for i := int32(0); i < n; i++ {
		next := i + 1
		if next == n {
			next = 0
		}
		dl.rowDummy[i].d = dl.rowDummy[next]
		dl.rowDummy[next].u = dl.rowDummy[i]
	}
	for i := int32(0); i < m; i++ {
		dl.colDummy[i] = newNode(-1, i)
	}
	for i := int32(0); i+1 < m; i++ {
		next := i + 1
		dl.colDummy[i].r = dl.colDummy[next]
		dl.colDummy[next].l = dl.colDummy[i]
	}
	dl.colHead = newNode(-1, -1)
	dl.colHead.l = dl.colDummy[m-1]
	dl.colDummy[m-1].r = dl.colHead
	dl.colHead.r = dl.colDummy[0]
	dl.colDummy[0].l = dl.colHead
	dl.resStack = make([]*DNode, 0, n)
}

func (dl *DancingLink) dfs0(root *DNode) {
	if root.deleted {
		return
	}
	dl.remove(root)
	dl.stack = append(dl.stack, root)
	if root.col == -1 || root.row == -1 {
		return
	}
	if root.u != root {
		dl.dfs1(root.u)
	}
	if root.d != root {
		dl.dfs1(root.d)
	}
	if root.l != root {
		dl.dfs0(root.l)
	}
	if root.r != root {
		dl.dfs0(root.r)
	}
}

func (dl *DancingLink) dfs1(root *DNode) {
	if root.deleted {
		return
	}
	dl.remove(root)
	dl.stack = append(dl.stack, root)
	if root.col == -1 || root.row == -1 {
		return
	}
	if root.u != root {
		dl.dfs1(root.u)
	}
	if root.d != root {
		dl.dfs1(root.d)
	}
	if root.l != root {
		dl.dfs2(root.l)
	}
	if root.r != root {
		dl.dfs2(root.r)
	}
}

func (dl *DancingLink) dfs2(root *DNode) {
	if root.deleted {
		return
	}
	dl.remove(root)
	dl.stack = append(dl.stack, root)
	if root.col == -1 || root.row == -1 {
		return
	}
	if root.l != root {
		dl.dfs2(root.l)
	}
	if root.r != root {
		dl.dfs2(root.r)
	}
}

func (dl *DancingLink) removeRow(i int32) {
	dl.resStack = append(dl.resStack, dl.rowDummy[i])
	dl.dfs0(dl.rowDummy[i].r)
}

func (dl *DancingLink) undo(size int32) {
	dl.resStack = dl.resStack[:len(dl.resStack)-1]
	for len(dl.stack) > int(size) {
		dl.recover(dl.stack[len(dl.stack)-1])
		dl.stack = dl.stack[:len(dl.stack)-1]
	}
}

func (dl *DancingLink) dance() bool {
	if dl.colHead.r == dl.colHead {
		dl.res = make([]int32, 0, len(dl.resStack))
		for _, node := range dl.resStack {
			dl.res = append(dl.res, node.row)
		}
		sort.Slice(dl.res, func(i, j int) bool { return dl.res[i] < dl.res[j] })
		return true
	}
	bestCol := dl.colHead.r
	for node := dl.colHead.r; node != dl.colHead; node = node.r {
		if dl.colSize[node.col] < dl.colSize[bestCol.col] {
			bestCol = node
		}
	}
	if dl.colSize[bestCol.col] == 0 {
		return false
	}
	bestRow := bestCol.d
	for node := bestCol.d; node != bestCol; node = node.d {
		if dl.rowSize[node.row] > dl.rowSize[bestRow.row] {
			bestRow = node
		}
	}
	iter := bestRow
	for {
		iter = iter.d
		if iter.row == -1 {
			iter = iter.d
		}
		now := int32(len(dl.stack))
		dl.removeRow(iter.row)
		if dl.dance() {
			return true
		}
		dl.undo(now)
		if iter == bestRow {
			break
		}
	}
	return false
}

type DNode struct {
	u, d, l, r *DNode
	row, col   int32
	deleted    bool
}

func newNode(row, col int32) *DNode {
	res := &DNode{row: row, col: col}
	res.u, res.d, res.l, res.r = res, res, res, res
	return res
}

func (node *DNode) String() string {
	return fmt.Sprintf("(%d,%d)", node.row, node.col)
}

// 9*9解数独.
// 转化成精准覆盖问题：
func SolveSudoku(mat [][]int32, empty func(i, j int32) bool) [][]int32 {
	cover := make([][]bool, 9*9*9)
	for i := 0; i < 9*9*9; i++ {
		cover[i] = make([]bool, 9*9*4)
	}
	for i := int32(0); i < 9; i++ {
		for j := int32(0); j < 9; j++ {
			for k := int32(0); k < 9; k++ {
				row := (i*9+j)*9 + k
				cell := (i/3)*3 + (j / 3)
				cover[row][i*9+k] = true
				cover[row][9*9+j*9+k] = true
				cover[row][9*9*2+cell*9+k] = true
				if empty(i, j) || mat[i][j]-1 == k {
					cover[row][9*9*3+i*9+j] = true
				}
			}
		}
	}
	dl := NewDancingLinkDense(cover, int32(len(cover)), int32(len(cover[0])))
	selected := dl.GetSolution()
	if selected == nil {
		return nil
	}
	res := make([][]int32, 9)
	for i := 0; i < 9; i++ {
		res[i] = make([]int32, 9)
	}
	for _, choice := range selected {
		k := choice % 9
		choice /= 9
		j := choice % 9
		choice /= 9
		i := choice
		res[i][j] = k + 1
	}
	return res
}
