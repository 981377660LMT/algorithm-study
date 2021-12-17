from typing import List


class Solution:
    def largestArea(self, grid: List[str]) -> int:
        m, n = len(grid), len(grid[0])

        res = 0
        visit = [[False] * n for _ in range(m)]
        valid = [[True] * n for _ in range(m)]

        count = [0]

        def dfs(i, j) -> bool:
            # 看完结束了
            if visit[i][j]:
                return True

            visit[i][j] = True
            count[0] += 1

            res = True
            if not valid[i][j]:
                res = False

            for ii, jj in [(i - 1, j), (i + 1, j), (i, j - 1), (i, j + 1)]:
                if (
                    ii >= 0
                    and ii < m
                    and jj >= 0
                    and jj < n
                    and grid[ii][jj] != '0'
                    and grid[ii][jj] == grid[i][j]
                ):
                    if not dfs(ii, jj):
                        res = False

            return res

        for i in range(m):
            for j in range(n):
                if grid[i][j] == '0':
                    valid[i][j] = False
                    for ii, jj in [(i - 1, j), (i + 1, j), (i, j - 1), (i, j + 1)]:
                        if ii >= 0 and ii < m and jj >= 0 and jj < n:
                            valid[ii][jj] = False

                if i == 0 or j == 0 or i == m - 1 or j == n - 1:
                    valid[i][j] = False

        for i in range(m):
            for j in range(n):
                if not visit[i][j]:
                    count[0] = 0
                    if dfs(i, j):
                        res = max(count[0], res)

        return res
