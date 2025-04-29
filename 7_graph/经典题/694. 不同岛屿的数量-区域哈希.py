# 694. 不同岛屿的数量-区域哈希
# 岛屿哈希值
# https://leetcode.cn/problems/number-of-distinct-islands/solutions/2418255/bu-tong-dao-yu-de-shu-liang-by-leetcode-kft01/
# !我们不考虑旋转、翻转操作。
#
# !1. 根据本地坐标哈希
#     每个岛屿的坐标相对于其左上角的坐标进行哈希.
# !2. 根据路径签名进行哈希
#     每次dfs都从左上角开始，最后生成一个dfs的哈希(用每次的转向表示)，如果岛屿一样则dfs必定一样

from typing import List, Tuple


class Solution:
    def numDistinctIslands(self, grid: List[List[int]]) -> int:
        """根据路径签名进行哈希."""
        row, col = len(grid), len(grid[0])
        dirs = [(0, 1), (1, 0), (0, -1), (-1, 0)]

        def dfs(r: int, c: int, path: List[int]) -> None:
            if not grid[r][c]:
                return

            grid[r][c] = 0
            for i, (dx, dy) in enumerate(dirs):
                nx, ny = r + dx, c + dy
                if 0 <= nx < row and 0 <= ny < col and grid[nx][ny] == 1:
                    path.append(i)
                    dfs(nx, ny, path)
                    path.append(-1)

        row, col = len(grid), len(grid[0])

        res = set()
        for r in range(row):
            for c in range(col):
                if grid[r][c] == 1:
                    path = []
                    dfs(r, c, path)
                    res.add(tuple(path))

        return len(res)

    def numDistinctIslands2(self, grid: List[List[int]]) -> int:
        """根据本地坐标哈希."""
        res = set()

        dirs = [(0, 1), (1, 0), (0, -1), (-1, 0)]
        m = len(grid)
        n = len(grid[0])

        def dfs(x: int, y: int, sx: int, sy: int, path: List[Tuple[int, int]]) -> None:
            grid[x][y] = 0
            path.append((x - sx, y - sy))
            for dx, dy in dirs:
                nx, ny = x + dx, y + dy
                if 0 <= nx < m and 0 <= ny < n and grid[nx][ny] == 1:
                    dfs(nx, ny, sx, sy, path)

        for x in range(m):
            for y in range(n):
                if grid[x][y] == 1:
                    path = []
                    dfs(x, y, x, y, path)
                    res.add(tuple(path))

        return len(res)
