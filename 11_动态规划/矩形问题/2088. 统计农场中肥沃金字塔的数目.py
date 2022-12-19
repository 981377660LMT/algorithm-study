# 2088. 统计农场中肥沃金字塔的数目

from typing import List
from functools import lru_cache

# How to find the biggest piramide with given top (x, y)?
# We need to look at points (x +1, y - 1) and (x+1, y+1) and
# find the biggest pyramides for these points:
# 1 <= m, n <= 1000
# 1 <= m * n <= 1e5  (暗示O(m*n为复杂度))


class Solution:
    def countPyramids(self, grid: List[List[int]]) -> int:
        @lru_cache(None)
        def dfs(r: int, c: int, dr: int) -> int:
            """以 (r, c) 为顶点,方向为 dr 的金字塔的最大`高度`"""
            if grid[r][c] == 0:
                return 0
            if r + dr < 0 or r + dr >= ROW or c - 1 < 0 or c + 1 >= COL:
                return 1
            if grid[r + dr][c] == 0:
                return 1
            left, right = dfs(r + dr, c - 1, dr), dfs(r + dr, c + 1, dr)
            return min(left, right) + 1

        ROW, COL = len(grid), len(grid[0])
        res = 0
        for r in range(ROW):
            for c in range(COL):
                # !高度为h时,有h-1个金字塔(高为1的金字塔不算)
                if grid[r][c] == 1:
                    res += (dfs(r, c, 1) - 1) + (dfs(r, c, -1) - 1)
        dfs.cache_clear()
        return res


print(Solution().countPyramids(grid=[[0, 1, 1, 0], [1, 1, 1, 1]]))
# [[0,1,1],[1,1,1],[1,0,1],[1,1,1]]
print(Solution().countPyramids(grid=[[0, 1, 1], [1, 1, 1], [1, 0, 1], [1, 1, 1]]))
# 输出：2
# 解释：
# 2 个可能的金字塔区域分别如上图蓝色和红色区域所示。
# 这个网格图中没有倒金字塔区域。
# 所以金字塔区域总数为 2 + 0 = 2 。
