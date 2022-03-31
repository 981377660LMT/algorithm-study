# 1 ≤ n, m ≤ 20
from itertools import product

# 注意此题不可记忆化
# 1219. 黄金矿工


class Solution:
    def solve(self, matrix):
        def dfs(x: int, y: int) -> int:
            raw = matrix[x][y]
            matrix[x][y] = 0

            res = raw
            for dx, dy in zip((0, 1, 0, -1), (1, 0, -1, 0)):
                nx, ny = x + dx, y + dy
                if 0 <= nx < row and 0 <= ny < col and matrix[nx][ny] != 0:
                    cand = dfs(nx, ny) + raw
                    res = max(res, cand)

            matrix[x][y] = raw
            return res

        row, col = len(matrix), len(matrix[0])

        return max(
            (dfs(r, c) for r, c in product(range(row), range(col)) if matrix[r][c] != 0), default=0
        )


print(Solution().solve(matrix=[[1, 3, 2], [2, 5, 0], [1, 0, 10]]))
