from collections import deque


class Solution:

    def numIslands(self, grid):
        if not grid or not grid[0]:
            return 0
        
        m, n = len(grid), len(grid[0])
        visited = [[False] * n for _ in range(m)]
        directions = [(0, 1), (1, 0), (0, -1), (-1, 0)]
        res = 0

        for i in range(m):
            for j in range(n):
                if grid[i][j] == '1' and \
                    not visited[i][j]:
                    res += 1
                    queue = deque([(i, j)])
                    visited[i][j] = True
                    while queue:
                        ci, cj = queue.popleft()
                        for di, dj in directions:
                            newi, newj = ci + di, cj + dj
                            if 0 <= newi < m \
                                and 0 <= newj < n \
                                and grid[newi][newj] == '1' \
                                and not visited[newi][newj]:
                                queue.append((newi, newj))
                                visited[newi][newj] = 1
        
        return res
         
                    