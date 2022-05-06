from functools import lru_cache
import sys
from typing import List

sys.setrecursionlimit(int(1e9))
# 这只青蛙从点 0 处跑道 2 出发，并想到达点 n 处的 任一跑道 ，请你返回 最少侧跳次数 。
# 注意：点 0 处和点 n 处的任一跑道都不会有障碍。
# 1 <= n <= 5 * 105
# dp[i][j]表示第i点第j道最少的侧跳次数


class Solution:
    def minSideJumps(self, obstacles: List[int]) -> int:
        """TLE"""

        @lru_cache(None)
        def dfs(col: int, row: int) -> int:
            if col == n - 1:
                return 0

            res = int(1e20)
            for nextRow in range(1, 4):
                if obstacles[col + 1] == nextRow or obstacles[col] == nextRow:
                    continue
                res = min(res, dfs(col + 1, nextRow) + int(nextRow != row))
            return res

        n = len(obstacles) - 1
        res = dfs(0, 2)
        dfs.cache_clear()
        return res

    def minSideJumps2(self, obstacles: List[int]) -> int:
        """AC
        
        用 dp 和 ndp 数组 来枚举行间状态转移，逻辑会清晰一些
        """
        n = len(obstacles)
        dp = [0, 0, 0]
        dp[0] = 1 if obstacles[1] != 1 else int(1e20)
        dp[1] = 0 if obstacles[1] != 2 else int(1e20)
        dp[2] = 1 if obstacles[1] != 3 else int(1e20)

        for i in range(2, n):
            ndp = [int(1e20)] * 3
            for cur in range(3):
                if obstacles[i] == cur + 1 or obstacles[i - 1] == cur + 1:
                    continue
                for pre in range(3):
                    ndp[cur] = min(ndp[cur], dp[pre] + int(cur != pre))
            dp = ndp

        return min(dp)


print(Solution().minSideJumps(obstacles=[0, 1, 2, 3, 0]))
print(Solution().minSideJumps2(obstacles=[0, 1, 2, 3, 0]))
# 输出：2
# 解释：最优方案如上图箭头所示。总共有 2 次侧跳（红色箭头）。
# 注意，这只青蛙只有当侧跳时才可以跳过障碍（如上图点 2 处所示）。
