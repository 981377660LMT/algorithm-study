from typing import List
from collections import deque

# 如果两个正方形具有相同的颜色并且在 4 个方向中的任何一个方向上彼此相邻，则它们属于同一个连通分量。


# 寻找border:
# if not (0 <= nx < m and 0 <= ny < n and grid[nx][ny] == originalColor):
#                     isBorder = True
class Solution:
    def colorBorder(self, grid: List[List[int]], row: int, col: int, color: int) -> List[List[int]]:
        originalColor = grid[row][col]
        m, n = len(grid), len(grid[0])
        visited = [[False] * n for _ in range(m)]
        borders = []
        direc = ((-1, 0), (1, 0), (0, -1), (0, 1))
        q = deque([(row, col)])
        visited[row][col] = True

        while q:
            x, y = q.popleft()
            isBorder = False
            for dx, dy in direc:
                nx, ny = x + dx, y + dy
                if not (0 <= nx < m and 0 <= ny < n and grid[nx][ny] == originalColor):
                    isBorder = True
                else:
                    if not visited[nx][ny]:
                        visited[nx][ny] = True
                        q.append((nx, ny))
            if isBorder:
                borders.append((x, y))

        for x, y in borders:
            grid[x][y] = color

        return grid


print(Solution().colorBorder(grid=[[1, 2, 2], [2, 3, 2]], row=0, col=1, color=3))
