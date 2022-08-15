from typing import List


# 给你一个 m x n 的矩阵 matrix ，
# !matrix元素互不相同
# 请你返回一个新的矩阵 answer ，
# 其中 answer[row][col] 是 matrix[row][col] 的`秩`。

# 1 <= m, n <= 1000
# m*n<=1e5
# !时间复杂度：O(m*n*log(m*n))
class Solution:
    def minScore(self, grid: List[List[int]]) -> List[List[int]]:
        """
        每次填入的数字为该行和该列的最大值再加1
        用新填入的数字更新当前行和列的最大值
        """
        ROW, COL = len(grid), len(grid[0])
        rowMax, colMax = [0] * ROW, [0] * COL
        nums = []

        for r in range(ROW):
            for c in range(COL):
                nums.append((grid[r][c], r, c))

        nums.sort(key=lambda x: x[0])
        res = [[0] * COL for _ in range(ROW)]
        for _, r, c in nums:
            res[r][c] = max(rowMax[r], colMax[c]) + 1
            rowMax[r] = res[r][c]
            colMax[c] = res[r][c]
        return res


print(Solution().minScore(grid=[[3, 1], [2, 5]]))
