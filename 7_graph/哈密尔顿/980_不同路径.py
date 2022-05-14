from functools import lru_cache
from typing import List, Tuple

# 1 <= grid.length * grid[0].length <= 20
# 返回在四个方向（上、下、左、右）上行走时，从起始方格到结束方格的不同路径的数目。
# 每一个无障碍方格都要通过一次，但是一条路径中不能重复通过同一个方格。

DIR4 = ((0, 1), (0, -1), (1, 0), (-1, 0))


class Solution:
    def uniquePathsIII(self, grid: List[List[int]]) -> int:
        @lru_cache(None)
        def dfs(cur: Tuple[int, int], state: int) -> int:
            if state == target:
                return int(cur == end)
            if cur == end:
                return int(state == target)

            res = 0
            for dr, dc in DIR4:
                nr, nc = cur[0] + dr, cur[1] + dc
                if 0 <= nr < ROW and 0 <= nc < COL and grid[nr][nc] != -1:
                    next = nr * COL + nc
                    if (state >> next) & 1:
                        continue
                    res += dfs((nr, nc), state | (1 << next))
            return res

        ROW, COL = len(grid), len(grid[0])
        start, end = (0, 0), (0, 0)
        target = 0
        for r in range(ROW):
            for c in range(COL):
                if grid[r][c] != -1:
                    target |= 1 << (r * COL + c)
                if grid[r][c] == 1:
                    start = (r, c)
                elif grid[r][c] == 2:
                    end = (r, c)
        return dfs(start, 1 << (start[0] * COL + start[1]))


print(Solution().uniquePathsIII([[1, 0, 0, 0], [0, 0, 0, 0], [0, 0, 2, -1]]))
