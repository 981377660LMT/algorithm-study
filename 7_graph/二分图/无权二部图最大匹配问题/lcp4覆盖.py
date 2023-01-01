from typing import List
from 匈牙利算法 import Hungarian

DIR4 = [(0, 1), (1, 0), (0, -1), (-1, 0)]


class Solution:
    def domino(self, row: int, col: int, broken: List[List[int]]) -> int:
        H = Hungarian(row * col, row * col)
        grid = [[0] * col for _ in range(row)]
        for r, c in broken:
            grid[r][c] = 1

        for r in range(row):
            for c in range(col):
                if grid[r][c] == 1 or (r + c) & 1:
                    continue
                cur = r * col + c
                for dr, dc in DIR4:
                    nr, nc = r + dr, c + dc
                    if 0 <= nr < row and 0 <= nc < col and grid[nr][nc] == 0:
                        H.addEdge(cur, nr * col + nc)
        return len(H.work())
