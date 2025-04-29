# 827. 最大人工岛
# https://leetcode.cn/problems/making-a-large-island/description/
#
# 给你一个大小为 n x n 二进制矩阵 grid 。最多 只能将一格 0 变成 1 。
# 返回执行此操作后，grid 中最大的岛屿面积是多少？
# 岛屿 由一组上、下、左、右四个方向相连的 1 形成。
#
# !把原本矩阵格子改成2以上的数字，能作为标记区别岛屿、也能避免重复访问

from typing import List


class Solution:
    def largestIsland(self, grid: List[List[int]]) -> int:
        n = len(grid)

        def dfs(i: int, j: int) -> int:
            size = 1
            grid[i][j] = len(area) + 2
            for x, y in (i - 1, j), (i + 1, j), (i, j - 1), (i, j + 1):
                if 0 <= x < n and 0 <= y < n and grid[x][y] == 1:
                    size += dfs(x, y)
            return size

        area = []
        for i, row in enumerate(grid):
            for j, v in enumerate(row):
                if v == 1:
                    area.append(dfs(i, j))

        if not area:
            return 1

        res = 0
        for i, row in enumerate(grid):
            for j, v in enumerate(row):
                if v:
                    continue
                s = set()
                for x, y in (i - 1, j), (i + 1, j), (i, j - 1), (i, j + 1):
                    if 0 <= x < n and 0 <= y < n and grid[x][y] > 1:
                        s.add(grid[x][y])
                res = max(res, 1 + sum(area[i - 2] for i in s))

        return res if res else n * n
