from typing import List
from functools import lru_cache

# 0 表示这个格子是空的，所以你可以穿过它。
# 1 表示这个格子里装着一个樱桃，你可以摘到樱桃然后穿过它。
# -1 表示这个格子里有荆棘，挡着你的路。
# 从位置 (0, 0) 出发，最后到达 (N-1, N-1) ，只能`向下或向右走`，并且只能穿越有效的格子（即只可以穿过值为0或者1的格子）；
# 当到达 (N-1, N-1) 后，你要继续走，直到返回到 (0, 0) ，只能`向上或向左走`，并且只能穿越有效的格子；
# 当你经过一个格子且这个格子包含一个樱桃时，你将摘到樱桃并且这个格子会变成空的（值变为0）；

# 1 <= N <= 50。

# 等效为两个人一起摘草莓，他们都要到终点；并且都只能向右或向下
class Solution:
    def cherryPickup(self, grid: List[List[int]]) -> int:
        '''
        Convert problem to two workers picking cherries from (0, 0) to (m-1, n-1)
        '''
        m, n = len(grid), len(grid[0])
        direction = [(0, 1), (1, 0)]

        @lru_cache(None)
        def dfs(x1, y1, x2, y2) -> int:
            if (x1, y1, x2, y2) == (m - 1, n - 1, m - 1, n - 1):
                return grid[x1][y1]

            if x1 >= m or x2 >= m or y1 >= n or y2 >= n:
                return -0x7FFFFFFF
                
            if -1 in (grid[x1][y1], grid[x2][y2]):
                return -0x7FFFFFFF

            res = -0x7FFFFFFF
            for dx1, dy1 in direction:
                for dx2, dy2 in direction:
                    # 两个人一起
                    cur = grid[x1][y1] if (x1, y1) == (x2, y2) else grid[x1][y1] + grid[x2][y2]
                    next = dfs(x1 + dx1, y1 + dy1, x2 + dx2, y2 + dy2)
                    res = max(res, cur + next)
            return res

        return max(0, dfs(0, 0, 0, 0))


print(Solution().cherryPickup(grid=[[0, 1, -1], [1, 0, -1], [1, 1, 1]]))
