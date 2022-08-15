from typing import List, Tuple, Optional
from collections import defaultdict, Counter


MOD = int(1e9 + 7)
INF = int(1e20)

DIR8 = ((1, 0), (-1, 0), (0, 1), (0, -1), (1, 1), (-1, -1), (1, -1), (-1, 1))


class Solution:
    def largestLocal(self, grid: List[List[int]]) -> List[List[int]]:
        n = len(grid)
        res = [[0] * (n - 2) for _ in range((n - 2))]
        for i in range(1, n - 1):
            for j in range(1, n - 1):
                cur = [grid[i][j]]
                for dr, dc in DIR8:
                    cur.append(grid[i + dr][j + dc])
                res[i - 1][j - 1] = max(cur)
        return res
