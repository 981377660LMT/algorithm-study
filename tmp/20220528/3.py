from typing import List, Optional, Tuple

MOD = int(1e9 + 7)
INF = int(1e20)


class Solution:
    def maximumImportance(self, n: int, roads: List[List[int]]) -> int:
        deg = [0] * n
        for u, v in roads:
            deg[u] += 1
            deg[v] += 1
        deg.sort()
        res = 0
        for a, b in zip(range(1, n + 1), deg):
            res += a * b
        return res

