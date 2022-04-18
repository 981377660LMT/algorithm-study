# 从第一行的任意0出发，在最后一行0结束 1表示障碍物
# 求可以通过的最长路径 你可以下、左、右走 不可以重复走


# 1 ≤ n * m ≤ 200,000
from functools import lru_cache
from typing import Literal

Direction = Literal[-1, 0, 1]


class Solution:
    def solve(self, matrix):
        row = len(matrix)
        col = len(matrix[0])

        @lru_cache(None)
        def dfs(r, c, direction: Direction):
            # dc表示前一次横向移动的方向 每行一旦选好了，就不能再回头
            if r == row:
                return 0
            if matrix[r][c] == 1:
                return -int(1e20)

            res = -int(1e20)
            if direction == 0:
                res = max(res, dfs(r, c, 1))
                res = max(res, dfs(r, c, -1))
            elif 0 <= c + direction < col:
                res = max(res, dfs(r, c + direction, direction) + 1)
            res = max(res, dfs(r + 1, c, 0) + 1)

            return res

        res = 0
        for start in range(col):
            res = max(res, dfs(0, start, 0))
        return res


print(Solution().solve(matrix=[[0, 0, 0, 0], [1, 0, 0, 0], [0, 0, 0, 0]]))
# We can move (0, 0), (0, 1), (0, 2), (0, 3), (1, 3), (1, 2), (1, 1), (2, 1), (2, 2), (2, 3).
