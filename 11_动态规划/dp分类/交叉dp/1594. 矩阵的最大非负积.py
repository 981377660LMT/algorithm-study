from typing import List, Tuple
from functools import lru_cache

MOD = int(1e9 + 7)
INF = int(1e20)
# 在从左上角 (0, 0) 开始到右下角 (rows - 1, cols - 1) 结束的所有路径中，找出具有 最大非负积 的路径。
# 每一步，你可以在矩阵中 向右 或 向下 移动。
# 如果最大积为负数，则返回 -1 。
# 1 <= rows, cols <= 15


class Solution:
    def maxProductPath(self, grid: List[List[int]]) -> int:
        """注意中途不能取模。会破坏“最”的特性"""

        @lru_cache(None)
        def dfs(row: int, col: int) -> Tuple[int, int]:
            if row == 0 and col == 0:
                return grid[0][0], grid[0][0]
            if row < 0 or col < 0:
                return -INF, INF
            if grid[row][col] == 0:
                return 0, 0

            # from top
            maxPos1, minNeg1 = dfs(row - 1, col)
            # from left
            maxPos2, minNeg2 = dfs(row, col - 1)

            if grid[row][col] > 0:
                return (
                    max(maxPos1, maxPos2) * grid[row][col],
                    min(minNeg1, minNeg2) * grid[row][col],
                )
            else:
                return (
                    min(minNeg1, minNeg2) * grid[row][col],
                    max(maxPos1, maxPos2) * grid[row][col],
                )

        # 注意要倒序，表示目标为(m-1,n-1)
        ROW, COL = len(grid), len(grid[0])
        maxPos, _ = dfs(ROW - 1, COL - 1)
        return -1 if maxPos < 0 else maxPos % MOD


print(Solution().maxProductPath(grid=[[-1, -2, -3], [-2, -3, -3], [-3, -3, -2]]))
# 输出：-1
# 解释：从 (0, 0) 到 (2, 2) 的路径中无法得到非负积，所以返回 -1
