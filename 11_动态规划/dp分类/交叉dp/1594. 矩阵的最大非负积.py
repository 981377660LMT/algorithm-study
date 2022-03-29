from typing import List, Tuple
from functools import lru_cache

MOD = int(1e9 + 7)
INF = 0x7FFFFFFF
# 在从左上角 (0, 0) 开始到右下角 (rows - 1, cols - 1) 结束的所有路径中，找出具有 最大非负积 的路径。
# 每一步，你可以在矩阵中 向右 或 向下 移动。
# 如果最大积为负数，则返回 -1 。
# 1 <= rows, cols <= 15


class Solution:
    def maxProductPath(self, grid: List[List[int]]) -> int:
        m, n = len(grid), len(grid[0])

        @lru_cache(None)
        def dfs(i: int, j: int) -> Tuple[int, int]:
            if i == 0 and j == 0:
                return grid[0][0], grid[0][0]
            if i < 0 or j < 0:
                return -INF, INF
            if grid[i][j] == 0:
                return 0, 0

            # from top
            maxPos1, minNeg1 = dfs(i - 1, j)
            # from left
            maxPos2, minNeg2 = dfs(i, j - 1)

            if grid[i][j] > 0:
                return max(maxPos1, maxPos2) * grid[i][j], min(minNeg1, minNeg2) * grid[i][j]
            else:
                return min(minNeg1, minNeg2) * grid[i][j], max(maxPos1, maxPos2) * grid[i][j]

        # 注意要倒序，表示目标为(m-1,n-1)
        maxPos, _ = dfs(m - 1, n - 1)
        return -1 if maxPos < 0 else maxPos % MOD


print(Solution().maxProductPath(grid=[[-1, -2, -3], [-2, -3, -3], [-3, -3, -2]]))
# 输出：-1
# 解释：从 (0, 0) 到 (2, 2) 的路径中无法得到非负积，所以返回 -1
