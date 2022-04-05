# Intersection of Two Maps
# 两个矩阵相同岛屿的数量
from typing import List, Set


class Solution:
    def solve(self, a, b):
        def getHashOfIslands(grid: List[List[int]]) -> Set[str]:
            """获取矩阵中各个岛屿的哈希值"""

            def dfs(r: int, c: int, path: List[str]) -> None:
                if visited[r][c]:
                    return

                visited[r][c] = True
                dirs = [(r, c + 1), (r + 1, c), (r, c - 1), (r - 1, c)]
                for i in range(4):
                    nr, nc = dirs[i]
                    if 0 <= nr < row and 0 <= nc < col and grid[nr][nc] == 1:
                        path.append(str(i + 1))
                        dfs(nr, nc, path)

            row, col = len(grid), len(grid[0])
            visited = [[False for _ in range(col)] for _ in range(row)]

            res = set()
            for r in range(row):
                for c in range(col):
                    if grid[r][c] == 1 and not visited[r][c]:
                        path = []
                        dfs(r, c, path)
                        path.extend(['_', str(r), '_', str(c)])
                        res.add(''.join(path))

            return res

        s1, s2 = getHashOfIslands(a), getHashOfIslands(b)
        return len(s1 & s2)


print(
    Solution().solve(
        a=[[0, 1], [1, 0], [1, 0], [1, 1], [1, 0]], b=[[0, 1], [1, 0], [1, 0], [1, 1], [0, 1]]
    )
)

# 1 ≤ n * m ≤ 200,000
