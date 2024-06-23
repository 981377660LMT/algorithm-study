package main

const INF32 int32 = 1e9

// 100332. 包含所有 1 的最小矩形面积 II
// https://leetcode.cn/problems/find-the-minimum-area-to-cover-all-ones-ii/submissions/
func minimumSum(grid [][]int) int {
	cacl := func(b BoundingRect) int32 {
		top, bottom, left, right := b[0], b[1], b[2], b[3]
		minTop, maxBottom, minLeft, maxRight := INF32, -INF32, INF32, -INF32
		for r := top; r <= bottom; r++ {
			for c := left; c <= right; c++ {
				if grid[r][c] == 1 {
					minTop = min32(minTop, r)
					maxBottom = max32(maxBottom, r)
					minLeft = min32(minLeft, c)
					maxRight = max32(maxRight, c)
				}
			}
		}
		if minTop == INF32 || maxBottom == -INF32 || minLeft == INF32 || maxRight == -INF32 {
			return 0
		}
		return (maxBottom - minTop + 1) * (maxRight - minLeft + 1)
	}

	row, col := int32(len(grid)), int32(len(grid[0]))
	res := row * col
	EnumerateBoundingRect3(row, col, func(b1, b2, b3 BoundingRect) {
		res = min32(res, cacl(b1)+cacl(b2)+cacl(b3))
	})
	return int(res)
}

type BoundingRect = [4]int32 // (top,bottom,left,right)

// 给定一个row*col的矩阵,分割成3个不重合的矩形,返回所有可能的分割方法.
func EnumerateBoundingRect3(row, col int32, consumer func(b1, b2, b3 BoundingRect)) {
	// 三横
	for r1 := int32(0); r1 < row-2; r1++ {
		for r2 := r1 + 1; r2 < row-1; r2++ {
			consumer([4]int32{0, r1, 0, col - 1}, [4]int32{r1 + 1, r2, 0, col - 1}, [4]int32{r2 + 1, row - 1, 0, col - 1})
		}
	}

	// 三竖
	for c1 := int32(0); c1 < col-2; c1++ {
		for c2 := c1 + 1; c2 < col-1; c2++ {
			consumer([4]int32{0, row - 1, 0, c1}, [4]int32{0, row - 1, c1 + 1, c2}, [4]int32{0, row - 1, c2 + 1, col - 1})
		}
	}

	// 先一横 然后再切一竖
	for r := int32(0); r < row-1; r++ {
		for c := int32(0); c < col-1; c++ {
			consumer([4]int32{0, r, 0, c}, [4]int32{0, r, c + 1, col - 1}, [4]int32{r + 1, row - 1, 0, col - 1})
		}
		for c := int32(0); c < col-1; c++ {
			consumer([4]int32{0, r, 0, col - 1}, [4]int32{r + 1, row - 1, c + 1, col - 1}, [4]int32{r + 1, row - 1, 0, c})
		}
	}

	// 先一竖 再切一横
	for c := int32(0); c < col-1; c++ {
		for r := int32(0); r < row-1; r++ {
			consumer([4]int32{0, r, 0, c}, [4]int32{r + 1, row - 1, 0, c}, [4]int32{0, row - 1, c + 1, col - 1})
		}
		for r := int32(0); r < row-1; r++ {
			consumer([4]int32{0, row - 1, 0, c}, [4]int32{0, r, c + 1, col - 1}, [4]int32{r + 1, row - 1, c + 1, col - 1})
		}
	}
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

func min32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func max32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}
