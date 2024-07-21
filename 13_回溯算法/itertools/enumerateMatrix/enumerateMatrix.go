// 二维矩阵遍历工具函数
// https://github.com/EndlessCheng/codeforces-go/blob/ff168d8e767e1d09ce47b4107d5c2e511b8bf41d/copypasta/search.go#L1
// Api:
//  EnumerateMatrixSpirally(row, col int, f func(r, c int))
//  EnumerateMatrixAroundByDepth(mat [][]int, depth int, f func(r, c int))
//  EnumerateMatrixZigZag(row, col int, f func(r, c int))
//  EnumerateMatrixAroundByManhattan(n, m, ox, oy, dist int, f func(x, y int))
//  EnumerateMatrixInnerByManhattan(n, m, ox, oy, dist int, f func(x, y int))
//  EnumerateMatrixByManhattan(n, m, ox, oy int, f func(x, y int))
//  EnumerateMatrixAroundByChebyshev(n, m, ox, oy, dist int, f func(x, y int))
//  EnumerateMatrixBorder(x0, x1, y0, y1 int, f func(x, y int))

package main

func main() {

	EnumerateMatrixBorder(1, 3, 1, 3, func(x, y int) {
		println(x, y)
	})

}

func spiralOrder(matrix [][]int) []int {
	if len(matrix) == 0 {
		return nil
	}
	row, col := len(matrix), len(matrix[0])
	res := make([]int, 0, row*col)
	EnumerateMatrixSpirally(row, col, func(r, c int) {
		res = append(res, matrix[r][c])
	})
	return res
}

// 1914. 循环轮转矩阵
// LC1914 https://leetcode.cn/problems/cyclically-rotating-a-grid/
func rotateGrid(grid [][]int, k int) [][]int {
	row, col := len(grid), len(grid[0])
	newGrid := make([][]int, row)
	for i := range newGrid {
		newGrid[i] = make([]int, col)
	}
	for d := 0; d < min(row, col)/2; d++ {
		var points [][2]int
		EnumerateMatrixAroundByDepth(grid, d, func(r, c int) {
			points = append(points, [2]int{r, c})
		})
		m := len(points)
		for i, p := range points {
			to := (i - k) % m
			if to < 0 {
				to += m
			}
			newP := points[to]
			newGrid[newP[0]][newP[1]] = grid[p[0]][p[1]]
		}
	}
	return newGrid
}

func allCellsDistOrder(rows int, cols int, rCenter int, cCenter int) [][]int {
	points := make([][]int, 0, rows*cols)
	EnumerateMatrixByManhattan(rows, cols, rCenter, cCenter, func(x, y int) {
		points = append(points, []int{x, y})
	})
	return points
}

//
//

// 顺时针遍历row行col列的螺旋矩阵(SpiralMatrix).
func EnumerateMatrixSpirally(row, col int, f func(r, c int)) {
	dir4 := []struct{ r, c int }{{0, 1}, {1, 0}, {0, -1}, {-1, 0}} // 右下左上
	mat := make([][]int, row)
	for i := range mat {
		mat[i] = make([]int, col)
		for j := range mat[i] {
			mat[i][j] = -1
		}
	}
	r, c, di := 0, 0, 0
	for id := 0; id < row*col; id++ {
		f(r, c)
		mat[r][c] = id
		d := dir4[di]
		if x, y := r+d.r, c+d.c; x < 0 || x >= row || y < 0 || y >= col || mat[x][y] != -1 {
			di = (di + 1) & 3
			d = dir4[di]
		}
		r += d.r
		c += d.c
	}
}

// 顺时针遍历矩阵从外向内的第 d 圈（保证不自交）.
// 一共 (row+col-d*4-2)*2 个元素.
func EnumerateMatrixAroundByDepth(mat [][]int, depth int, f func(r, c int)) {
	row, col := len(mat), len(mat[0])
	for j := depth; j < col-depth; j++ { // →
		f(depth, j)
	}
	for i := depth + 1; i < row-depth; i++ { // ↓
		f(i, col-1-depth)
	}
	for j := col - depth - 2; j >= depth; j-- { // ←
		f(row-1-depth, j)
	}
	for i := row - depth - 2; i > depth; i-- { // ↑
		f(i, depth)
	}
}

// 获取之字遍历的所有坐标.
// 例如 3 行 3 列的矩阵：
// 0 -> 1 -> 2
//           ↓
// 3 <- 4 <- 5
// ↓
// 6 -> 7 -> 8

func EnumerateMatrixZigZag(row, col int, f func(r, c int)) {
	for i := 0; i < row; i++ {
		for j := 0; j < col; j++ {
			f(i, j)
		}
		i++
		if i == row {
			break
		}
		for j := col - 1; j >= 0; j-- {
			f(i, j)
		}
	}
}

// 遍历以 (ox, oy) 为中心的曼哈顿距离为 dist 的边界上的格点.
// 从最右顶点出发，逆时针移动.
func EnumerateMatrixAroundByManhattan(n, m, ox, oy, dist int, f func(x, y int)) {
	if dist == 0 {
		f(ox, oy)
		return
	}
	dir4r := []struct{ x, y int }{{-1, 1}, {-1, -1}, {1, -1}, {1, 1}} // 逆时针
	x, y := ox+dist, oy                                               // 从最右顶点出发，逆时针移动
	for _, d := range dir4r {
		for k := 0; k < dist; k++ {
			if 0 <= x && x < n && 0 <= y && y < m {
				f(x, y)
			}
			x += d.x
			y += d.y
		}
	}
}

// 从上到下，遍历以 (ox, oy) 为中心的曼哈顿距离 <= dist 的所有格点.
func EnumerateMatrixInnerByManhattan(n, m, ox, oy, dist int, f func(x, y int)) {
	for i := max(ox-dist, 0); i <= ox+dist && i < n; i++ {
		d := dist - abs(ox-i)
		for j := max(oy-d, 0); j <= oy+d && j < m; j++ {
			if i == ox && j == oy {
				continue
			}
			f(i, j)
		}
	}
}

// 曼哈顿圈序遍历.
// 从最右顶点出发，逆时针移动.
func EnumerateMatrixByManhattan(n, m, ox, oy int, f func(x, y int)) {
	f(ox, oy)
	dir4r := []struct{ x, y int }{{-1, 1}, {-1, -1}, {1, -1}, {1, 1}} // 逆时针
	maxDist := max(ox, n-1-ox) + max(oy, m-1-oy)
	for dis := 1; dis <= maxDist; dis++ {
		x, y := ox+dis, oy // 从最右顶点出发，逆时针移动
		for _, d := range dir4r {
			for k := 0; k < dis; k++ {
				if 0 <= x && x < n && 0 <= y && y < m {
					f(x, y)
				}
				x += d.x
				y += d.y
			}
		}
	}
}

// 遍历以 (ox, oy) 为中心的切比雪夫距离为 dist 的边界上的格点.
// 从最右顶点出发，逆时针移动.
func EnumerateMatrixAroundByChebyshev(n, m, ox, oy, dist int, f func(x, y int)) {
	// 上下
	for _, x := range []int{ox - dist, ox + dist} {
		if 0 <= x && x < n {
			for y := max(oy-dist, 0); y <= min(oy+dist, m-1); y++ {
				f(x, y)
			}
		}
	}
	// 左右（注意四角已经被上面的循环枚举到了）
	for _, y := range []int{oy - dist, oy + dist} {
		if 0 <= y && y < m {
			for x := max(ox-dist, 0) + 1; x <= min(ox+dist, n-1)-1; x++ {
				f(x, y)
			}
		}
	}
}

// 遍历矩阵的边界.
// x0<=x<=x1, y0<=y<=y1.
func EnumerateMatrixBorder(x0, x1, y0, y1 int, f func(x, y int)) {
	if y0 == y1 {
		for i := x0; i <= x1; i++ {
			f(i, y0)
		}
		return
	}
	for i := x0; i <= x1; i++ {
		for j := y0; j <= y1; {
			f(i, j)
			if i == x0 || i == x1 {
				j++
			} else {
				j += y1 - y0
			}
		}
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
