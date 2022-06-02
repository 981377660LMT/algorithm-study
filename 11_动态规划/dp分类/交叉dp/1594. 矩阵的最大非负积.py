from typing import List, Tuple
from functools import lru_cache

MOD = int(1e9 + 7)
INF = int(1e20)
# 在从左上角 (0, 0) 开始到右下角 (rows - 1, cols - 1) 结束的所有路径中，找出具有 最大非负积 的路径。
# 每一步，你可以在矩阵中 向右 或 向下 移动。
# 如果最大积为负数，则返回 -1 。
# 1 <= rows, cols <= 15

DIR2 = ((1, 0), (0, 1))


class Solution:
    def maxProductPath(self, grid: List[List[int]]) -> int:
        """注意中途不能取模。会破坏“最”的特性"""

        @lru_cache(None)
        def dfs(row: int, col: int) -> Tuple[int, int]:
            if (row, col) == (ROW - 1, COL - 1):
                return grid[row][col], grid[row][col]
            # if grid[row][col] == 0:
            #     return 0, 0

            resPos, resNeg = -INF, INF
            for dRow, dCol in DIR2:
                nRow, nCol = row + dRow, col + dCol
                if 0 <= nRow < ROW and 0 <= nCol < COL:
                    nextPos, nextNeg = dfs(nRow, nCol)
                    resPos = max(resPos, nextPos * grid[row][col], nextNeg * grid[row][col])
                    resNeg = min(resNeg, nextPos * grid[row][col], nextNeg * grid[row][col])
            return resPos, resNeg

        ROW, COL = len(grid), len(grid[0])
        pos, _ = dfs(0, 0)
        return -1 if pos < 0 else pos % MOD


print(Solution().maxProductPath(grid=[[-1, -2, -3], [-2, -3, -3], [-3, -3, -2]]))
# 输出：-1
# 解释：从 (0, 0) 到 (2, 2) 的路径中无法得到非负积，所以返回 -1
