from typing import List
from functools import lru_cache

# 从格子 (i,j) 出发，机器人可以移动到格子 (i+1, j-1)，(i+1, j) 或者 (i+1, j+1) 。
# 两个机器人最后都要到达 grid 最底下一行。


class Solution:
    def cherryPickup(self, grid: List[List[int]]) -> int:
        m = len(grid)
        n = len(grid[0])
        direction = [(1, 0), (1, -1), (1, 1)]

        @lru_cache(None)
        def dfs(i1, j1, i2, j2):
            if i1 >= m:
                return 0

            if j1 < 0 or j1 >= n or j2 < 0 or j2 >= n:
                return -0x7FFFFFFF

            res = 0

            for di1, dj1 in direction:
                for di2, dj2 in direction:
                    res = (
                        max(res, grid[i1][j1] + dfs(i1 + di1, j1 + dj1, i2 + di2, j2 + dj2))
                        if (i1 == i2 and j1 == j2)
                        else max(
                            res,
                            grid[i1][j1]
                            + grid[i2][j2]
                            + dfs(i1 + di1, j1 + dj1, i2 + di2, j2 + dj2),
                        )
                    )

            return res

        return dfs(0, 0, 0, n - 1)


print(Solution().cherryPickup(grid=[[3, 1, 1], [2, 5, 1], [1, 5, 5], [2, 1, 1]]))
# 输出：24
# 解释：机器人 1 和机器人 2 的路径在上图中分别用绿色和蓝色表示。
# 机器人 1 摘的樱桃数目为 (3 + 2 + 5 + 2) = 12 。
# 机器人 2 摘的樱桃数目为 (1 + 5 + 5 + 1) = 12 。
# 樱桃总数为： 12 + 12 = 24 。

