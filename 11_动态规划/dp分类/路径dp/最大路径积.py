# Maximum Product Path in Matrix
# 最大路径积

# 1 ≤ n, m ≤ 20
# -2 ≤ matrix[r][c] ≤ 2

from functools import lru_cache


MOD = int(1e9 + 7)
INF = int(1e20)


class Solution:
    def solve(self, matrix):
        row, col = len(matrix), len(matrix[0])

        @lru_cache(None)
        def dfs(i=0, j=0):
            """返回 最大，最小 值"""
            if i >= row or j >= col:
                return -INF, INF

            if i == row - 1 and j == col - 1:
                return matrix[-1][-1], matrix[-1][-1]

            max_, min_ = -INF, INF
            cur = matrix[i][j]
            if i + 1 < row:
                nMax, nMin = dfs(i + 1, j)
                max_ = max(max_, cur * nMax, cur * nMin)
                min_ = min(min_, cur * nMax, cur * nMin)
            if j + 1 < col:
                nMax, nMin = dfs(i, j + 1)
                max_ = max(max_, cur * nMax, cur * nMin)
                min_ = min(min_, cur * nMax, cur * nMin)

            return max_, min_

        res = dfs()[0]
        return -1 if res < 0 else res % MOD


print(Solution().solve(matrix=[[2, 1, -2], [-1, -1, -2], [1, 1, 1]]))
