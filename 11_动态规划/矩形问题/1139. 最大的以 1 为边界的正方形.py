# 给你一个由若干 0 和 1 组成的二维网格 grid，
# 请你找出`边界`全部由 1 组成的最大 正方形 子网格，
# 并返回该子网格中的元素数量。如果不存在，则返回 0。
# ROW,COL<=100

# 解:
# !预处理出每个点的左边和上边的连续1的个数，然后从最大的边长开始遍历
# 如果当前边长的正方形的四个角都满足条件，那么就是最大的正方形


from typing import List


class Solution:
    def largest1BorderedSquare(self, grid: List[List[int]]) -> int:
        ROW, COL = len(grid), len(grid[0])
        up, left = [[0] * COL for _ in range(ROW)], [[0] * COL for _ in range(ROW)]
        for row in range(ROW):
            for col in range(COL):
                if grid[row][col] == 1:
                    up[row][col] = up[row - 1][col] + 1 if row else 1
                    left[row][col] = left[row][col - 1] + 1 if col else 1

        res = 0
        for r in range(ROW):
            for c in range(COL):
                # 从最大的边长开始遍历
                for k in range(min(up[r][c], left[r][c]), res, -1):
                    if up[r][c - k + 1] >= k and left[r - k + 1][c] >= k:
                        res = k
                        break

        return res * res


assert (Solution().largest1BorderedSquare(grid=[[1, 1, 1], [1, 0, 1], [1, 1, 1]])) == 9
