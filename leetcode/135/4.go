package main

func main() {

}

// def max2(a: int, b: int) -> int:
//     return a if a > b else b

// def min2(a: int, b: int) -> int:
//     return a if a < b else b

// class Solution:
//     def maximumScore(self, grid: List[List[int]]) -> int:
//         n = len(grid)
//         ROW, COL = n, n
//         colPreSum = [list(accumulate(col, initial=0)) for col in zip(*grid)]

//         def cal(col: int, rowCount1: int, rowCount2: int) -> int:
//             if col < 0 or col >= COL:
//                 return 0
//             if rowCount1 >= rowCount2:
//                 return 0
//             arr = colPreSum[col]
//             return arr[rowCount2] - arr[rowCount1]

//         @lru_cache(None)
//         def dfs(index: int, preBlackRowCount1: int, preBlackRowCount2: int) -> int:
//             if index == COL:
//                 return cal(index - 1, preBlackRowCount1, preBlackRowCount2)
//             res = 0
//             for i in range(ROW + 1):
//                 pre = cal(index - 1, preBlackRowCount1, max2(i, preBlackRowCount2))  # 前一列白色得分
//                 res = max2(res, pre + dfs(index + 1, i, preBlackRowCount1))
//             return res

// res = dfs(0, 0, 0)
// dfs.cache_clear()
// return res

var memo [101][101][101]int

func maximumScore(grid [][]int) int64 {
	n := int32(len(grid))
	ROW, COL := n, n
	colPreSum := make([][]int, COL)
	for c := range colPreSum {
		colPreSum[c] = make([]int, ROW+1)
		for r := int32(1); r <= ROW; r++ {
			colPreSum[c][r] = colPreSum[c][r-1] + grid[r-1][c]
		}
	}

	cal := func(col, rowCount1, rowCount2 int32) int {
		if col < 0 || col >= COL {
			return 0
		}
		if rowCount1 >= rowCount2 {
			return 0
		}
		arr := colPreSum[col]
		return arr[rowCount2] - arr[rowCount1]
	}

	for i := int32(0); i < COL; i++ {
		for j := int32(0); j <= ROW; j++ {
			for k := int32(0); k <= ROW; k++ {
				memo[i][j][k] = -1
			}
		}
	}

	var dfs func(int32, int32, int32) int
	dfs = func(index, preBlackRowCount1, preBlackRowCount2 int32) int {
		if index == COL {
			return cal(index-1, preBlackRowCount1, preBlackRowCount2)
		}
		ptr := &memo[index][preBlackRowCount1][preBlackRowCount2]
		if *ptr != -1 {
			return *ptr
		}
		res := 0
		for i := int32(0); i <= ROW; i++ {
			pre := cal(index-1, preBlackRowCount1, max32(i, preBlackRowCount2))
			res = max(res, pre+dfs(index+1, i, preBlackRowCount1))
		}
		*ptr = res
		return res
	}

	res := dfs(0, 0, 0)
	return int64(res)
}

func max(a, b int) int {
	if a > b {
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
