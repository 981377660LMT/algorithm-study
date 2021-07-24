class Solution:
    def color_border(self, grid, r0, c0, color):
        if not grid or not grid[0]:
            return grid

        m, n = len(grid), len(grid[0])
        old_color = grid[r0][c0]
        stack = []
        visited = set()
        stack.append((r0, c0))
        visited.add((r0, c0))
        
        res = []
        while stack:
            ci, cj = stack.pop()
            for di, dj in [(0, 1), (1, 0), (0, -1), (-1, 0)]:
                newi, newj = ci + di, cj + dj
                if not 0 <= newi < m or not 0 <= newj < n:
                    res.append((ci, cj))
                    continue
                if grid[newi][newj] != old_color:
                    res.append((ci, cj))
                    continue
                if (newi, newj) in visited:
                    continue
                stack.append((newi, newj))
                visited.add((newi, newj))
        
        for i, j in res:
            grid[i][j] = color
        
        return grid