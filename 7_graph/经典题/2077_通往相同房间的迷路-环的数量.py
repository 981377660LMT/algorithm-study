from typing import List
from collections import defaultdict
from itertools import combinations


# 求无向图中长度为 3 的不同环的数量
class Solution:
    def numberOfPaths(self, n: int, corridors: List[List[int]]) -> int:
        adjMap = defaultdict(set)
        for u, v in corridors:
            adjMap[u].add(v)
            adjMap[v].add(u)

        res = 0
        for p1, p2 in corridors:
            res += len(adjMap[p1] & adjMap[p2])
        return res

        # res = 0
        # # 对邻居两两判断
        # for p1 in range(1, n + 1):
        #     for p2, p3 in combinations(adj[p1], 2):
        #         if p2 in adj[p3]:
        #             res += 1
        # return res // 3


# Output: 2
print(Solution().numberOfPaths(n=5, corridors=[[1, 2], [5, 2], [4, 1], [2, 4], [3, 1], [3, 4]]))
