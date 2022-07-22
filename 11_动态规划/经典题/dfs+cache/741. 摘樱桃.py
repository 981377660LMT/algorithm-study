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

DIR2 = [[0, 1], [1, 0]]


class Solution:
    def cherryPickup(self, grid: List[List[int]]) -> int:
        @lru_cache(None)
        def dfs(r1: int, c1: int, r2: int, c2: int) -> int:
            if (r1, c1, r2, c2) == (ROW - 1, COL - 1, ROW - 1, COL - 1):
                return grid[r1][c1]

            if r1 >= ROW or r2 >= ROW or c1 >= COL or c2 >= COL:
                return -int(1e20)

            if -1 in (grid[r1][c1], grid[r2][c2]):
                return -int(1e20)

            res = -int(1e20)
            for dr1, dc1 in DIR2:
                for dr2, dc2 in DIR2:
                    # 两个人一起
                    cur = grid[r1][c1] if (r1, c1) == (r2, c2) else grid[r1][c1] + grid[r2][c2]
                    next = dfs(r1 + dr1, c1 + dc1, r2 + dr2, c2 + dc2)
                    if cur + next > res:
                        res = cur + next
            return res

        ROW, COL = len(grid), len(grid[0])
        res = dfs(0, 0, 0, 0)
        dfs.cache_clear()
        return res if res > 0 else 0


print(Solution().cherryPickup(grid=[[0, 1, -1], [1, 0, -1], [1, 1, 1]]))
