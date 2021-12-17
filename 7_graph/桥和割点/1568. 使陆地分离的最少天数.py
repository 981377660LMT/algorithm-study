# 一共3种情况，0，1，2， 并查集求岛屿数量如果大于2 返回0， 如果岛屿数量为1， tarjan算法求割点，
# 如果找到割点返回 1，没有割点则返回，2
from typing import List
import copy

# Day 0: Check islands at day 0, return 0 if you have less than or greater than one island.
# Day 1: If not, try to add water at any given location, and check if that gives you a valid island formation.
# Day 2: Else, just return 2!


class Solution:
    # 消除1
    def no_islands_recur(self, grid, i, j, m, n):
        if grid[i][j] == 0:
            return
        grid[i][j] = 0
        if i - 1 >= 0:
            self.no_islands_recur(grid, i - 1, j, m, n)
        if i + 1 < m:
            self.no_islands_recur(grid, i + 1, j, m, n)
        if j - 1 >= 0:
            self.no_islands_recur(grid, i, j - 1, m, n)
        if j + 1 < n:
            self.no_islands_recur(grid, i, j + 1, m, n)

    # dfs 寻找连通分量
    def no_islands(self, grid):
        res = 0
        m, n = len(grid), len(grid[0])
        for i in range(m):
            for j in range(n):
                if grid[i][j] == 1:
                    res += 1
                    self.no_islands_recur(grid, i, j, m, n)
        return res

    def minDays(self, grid: List[List[int]]) -> int:
        # if we have 0 or more than 1 islands at day 0, return day 0
        time = 0
        grid_copy = copy.deepcopy(grid)
        n = self.no_islands(grid_copy)
        if n != 1:
            return time

        # try to remove any land any see if it works
        # 这一步可用tarjan寻找割点
        time = 1
        for i in range(len(grid)):
            for j in range(len(grid[0])):
                grid_copy = copy.deepcopy(grid)
                grid_copy[i][j] = 0
                n = self.no_islands(grid_copy)
                if n != 1:
                    return time

        # well then just return 2
        time = 2
        return time


print(Solution().minDays(grid=[[0, 1, 1, 0], [0, 1, 1, 0], [0, 0, 0, 0]]))
