from math import comb
from typing import List
from UnionFind import UnionFindArray

# 无向图中不连通的点对数


class Solution:
    def countPairs(self, n: int, edges: List[List[int]]) -> int:
        """不连通的点对数=点对数-联通的点对数"""
        uf = UnionFindArray(n)

        for u, v in edges:
            uf.union(u, v)
        groups = uf.getGroups()
        return comb(n, 2) - sum(comb(len(g), 2) for g in groups.values())


print(Solution().countPairs(n=7, edges=[[0, 2], [0, 5], [2, 4], [1, 6], [5, 4]]))
print(Solution().countPairs(n=3, edges=[[0, 1], [0, 2], [1, 2]]))
