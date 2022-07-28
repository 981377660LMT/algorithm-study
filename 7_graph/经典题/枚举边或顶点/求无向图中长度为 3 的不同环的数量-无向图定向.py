# 2077. Paths in Maze That Lead to Same Room
# https://leetcode-cn.com/problems/paths-in-maze-that-lead-to-same-room/
# 无向图三元环的数量
# n<=1000
# !边数<=5e4
# 如果不限定边数 应该用枚举两个点 然后bitSet求交集 O(n^3/64)

from collections import defaultdict
from typing import List


class Solution:
    def numberOfPaths(self, n: int, edges: List[List[int]]) -> int:
        """三元环计数 枚举边+无向图定向 104ms O(E*n/64)"""
        adjMap = [0] * (n + 1)
        for u, v in edges:
            if u < v:
                adjMap[u] |= 1 << v
            else:
                adjMap[v] |= 1 << u

        res = 0
        for p1, p2 in edges:
            res += (adjMap[p1] & adjMap[p2]).bit_count()
        return res

    def numberOfPaths2(self, n: int, edges: List[List[int]]) -> int:
        """枚举点+无向边定向 O(E^3/2)"""
        adjMap = [[] for _ in range(n + 1)]
        deg = [0] * (n + 1)
        for u, v in edges:
            deg[u] += 1
            deg[v] += 1

        for u, v in edges:
            # u, v = sorted((u, v), key=lambda x: (deg[x], x))
            if deg[u] > deg[v] or deg[u] == deg[v] and u > v:
                u, v = v, u
            adjMap[u].append(v)

        res = 0
        for p1 in range(1, n + 1):
            for p2 in adjMap[p1]:
                for p3 in adjMap[p2]:
                    if p3 in adjMap[p1]:
                        res += 1
        return res


print(Solution().numberOfPaths(n=5, edges=[[1, 2], [5, 2], [4, 1], [2, 4], [3, 1], [3, 4]]))


print(Solution().numberOfPaths2(n=5, edges=[[1, 2], [5, 2], [4, 1], [2, 4], [3, 1], [3, 4]]))
