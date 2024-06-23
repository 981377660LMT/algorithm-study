package main

// 2257. 统计网格图中没有被保卫的格子数
// https://leetcode.cn/problems/count-unguarded-cells-in-the-grid/description/
func countUnguarded(m int, n int, guards [][]int, walls [][]int) int {
	row, col := int32(m), int32(n)
	isGood := make(map[[2]int32]struct{})
	for _, g := range guards {
		isGood[[2]int32{int32(g[0]), int32(g[1])}] = struct{}{}
	}
	nearestGood := FindNearestSpecial(row, col, func(r, c int32) bool {
		_, ok := isGood[[2]int32{r, c}]
		return ok
	})
	isBad := make(map[[2]int32]struct{})
	for _, w := range walls {
		isBad[[2]int32{int32(w[0]), int32(w[1])}] = struct{}{}
	}
	nearestBad := FindNearestSpecial(row, col, func(r, c int32) bool {
		_, ok := isBad[[2]int32{r, c}]
		return ok
	})

	check := func(r, c int32, nearestGood, nearestBad [][]BoundingRect) bool {
		top1, bottom1, left1, right1 := nearestGood[r][c][0], nearestGood[r][c][1], nearestGood[r][c][2], nearestGood[r][c][3]
		top2, bottom2, left2, right2 := nearestBad[r][c][0], nearestBad[r][c][1], nearestBad[r][c][2], nearestBad[r][c][3]
		if top1 >= top2 && top1 != -1 {
			return true
		}
		if left1 >= left2 && left1 != -1 {
			return true
		}
		if right1 <= right2 && right1 != col {
			return true
		}
		if bottom1 <= bottom2 && bottom1 != row {
			return true
		}
		return false
	}

	res := 0
	for r := int32(0); r < row; r++ {
		for c := int32(0); c < col; c++ {
			if _, ok := isGood[[2]int32{r, c}]; ok {
				continue
			}
			if _, ok := isBad[[2]int32{r, c}]; ok {
				continue
			}
			if !check(r, c, nearestGood, nearestBad) {
				res++
			}
		}
	}
	return res
}

type BoundingRect = [4]int32 // (top,bottom,left,right)

// 给定一个row*col的矩阵, 对每个点, 找到 top,bottom,left,right 四个方向上最近的特殊点(包含自身).
// 若top不存在, 则top=-1;
// 若bottom不存在, 则bottom=row;
// 若left不存在, 则left=-1;
// 若right不存在, 则right=col.
func FindNearestSpecial(row, col int32, isSpecialFn func(r, c int32) bool) [][]BoundingRect {
	isSpecial := make([][]bool, row)
	for r := int32(0); r < row; r++ {
		isSpecial[r] = make([]bool, col)
		for c := int32(0); c < col; c++ {
			isSpecial[r][c] = isSpecialFn(r, c)
		}
	}

	leftRes := make([][]int32, row)
	rightRes := make([][]int32, row)
	topRes := make([][]int32, col)
	bottomRes := make([][]int32, col)
	for r := int32(0); r < row; r++ {
		leftRes[r] = make([]int32, col)
		rightRes[r] = make([]int32, col)
		for c := int32(0); c < col; c++ {
			leftRes[r][c] = -1
			rightRes[r][c] = col
		}
	}
	for c := int32(0); c < col; c++ {
		topRes[c] = make([]int32, row)
		bottomRes[c] = make([]int32, row)
		for r := int32(0); r < row; r++ {
			topRes[c][r] = -1
			bottomRes[c][r] = row
		}
	}

	for r := int32(0); r < row; r++ {
		for c := int32(0); c < col; c++ {
			if r-1 >= 0 {
				topRes[c][r] = topRes[c][r-1]
			}
			if c-1 >= 0 {
				leftRes[r][c] = leftRes[r][c-1]
			}
			if isSpecial[r][c] {
				topRes[c][r] = r
				leftRes[r][c] = c
			}
		}
	}

	for r := row - 1; r >= 0; r-- {
		for c := col - 1; c >= 0; c-- {
			if r+1 < row {
				bottomRes[c][r] = bottomRes[c][r+1]
			}
			if c+1 < col {
				rightRes[r][c] = rightRes[r][c+1]
			}
			if isSpecial[r][c] {
				bottomRes[c][r] = r
				rightRes[r][c] = c
			}
		}
	}

	res := make([][]BoundingRect, row)
	for r := int32(0); r < row; r++ {
		res[r] = make([]BoundingRect, col)
		for c := int32(0); c < col; c++ {
			res[r][c] = BoundingRect{topRes[c][r], bottomRes[c][r], leftRes[r][c], rightRes[r][c]}
		}
	}

	return res
}
