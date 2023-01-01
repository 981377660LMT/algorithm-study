// https://atcoder.jp/contests/abc233/submissions/37628379
// n<=1e5
// 4558 ms  O(nlogn*logn)

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	type Point struct{ x, y int }
	points := make([]Point, n)
	for i := range points {
		var x, y int
		fmt.Fscan(in, &x, &y)
		points[i] = Point{x + y + 1, x - y + 1} // 变为切比雪夫距离然后平移1
	}

	sort.Slice(points, func(i, j int) bool {
		return points[i].y < points[j].y
	})
	bit := NewRectangleBIT(2e5 + 10)
	for _, p := range points {
		bit.Add(p.x, p.y)
	}

	// 到(x,y)的切比雪夫距离不超过mid的点的个数
	countNGT := func(mid, x, y int) int {
		qx, qy := x+y+1, x-y+1
		return bit.Query(qx-mid, qy-mid, qx+mid, qy+mid)
	}

	var q int
	fmt.Fscan(in, &q)
	for i := 0; i < q; i++ {
		var x, y, k int
		fmt.Fscan(in, &x, &y, &k)
		// 二分答案
		left, right := 0, int(2e5+10)
		for left <= right {
			mid := (left + right) >> 1
			if countNGT(mid, x, y) < k {
				left = mid + 1
			} else {
				right = mid - 1
			}
		}
		fmt.Fprintln(out, left)
	}
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
