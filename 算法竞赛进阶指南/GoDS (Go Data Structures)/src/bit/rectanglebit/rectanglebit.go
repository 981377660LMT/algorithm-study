// 静态二维矩形区域计数 O(logn)

package main

import (
	"fmt"
	"sort"
)

func main() {
	points := [][]int{
		{1, 1}, {1, 2}, {1, 3}, {1, 4}, {1, 5},
		{2, 1}, {2, 2}, {2, 3}, {2, 4}, {2, 5},
		{3, 1}, {3, 2}, {3, 3}, {3, 4}, {3, 5},
		{4, 1}, {4, 2}, {4, 3}, {4, 4}, {4, 5},
		{5, 1}, {5, 2}, {5, 3}, {5, 4}, {5, 5},
	}

	sort.Slice(points, func(i, j int) bool {
		return points[i][1] < points[j][1]
	})

	bit := NewRectangleBIT(6)
	for _, p := range points {
		bit.Add(p[0], p[1])
	}

	fmt.Println(bit.Query(1, 1, 5, 5)) // 25

}

var bisectLeft = sort.SearchInts

var bisectRight = func(a []int, x int) int {
	return sort.Search(len(a), func(i int) bool { return a[i] > x })
}

// 二维矩形计数 更新和查询时间复杂度O(logk)
type RectangleBIT struct {
	row  int
	tree [][]int
}

func NewRectangleBIT(row int) *RectangleBIT {
	return &RectangleBIT{
		row:  row,
		tree: make([][]int, row),
	}
}

// 加入点(row,col) 注意加入过程中需要保证col递增 时间复杂度log(k)
func (b *RectangleBIT) Add(row, col int) {
	if row <= 0 {
		panic("row 必须是正整数")
	}
	for row < b.row {
		b.tree[row] = append(b.tree[row], col)
		row += row & -row
	}
}

// 计算矩形内的点数 时间复杂度log(k)
func (b *RectangleBIT) Query(r1, c1, r2, c2 int) int {
	if r1 >= b.row {
		r1 = b.row - 1
	}
	if r2 >= b.row {
		r2 = b.row - 1
	}
	return b.query(r2, c1, c2) - b.query(r1-1, c1, c2)
}

// row不超过rowUpper,col在[c1,c2]间的点数
func (b *RectangleBIT) query(rowUpper, col1, col2 int) int {
	res := 0
	index := rowUpper
	for index > 0 {
		res += bisectRight(b.tree[index], col2) - bisectLeft(b.tree[index], col1)
		index -= index & -index
	}
	return res
}
