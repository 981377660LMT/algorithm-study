from typing import List
from collections import defaultdict
from itertools import combinations


# 求无向图中长度为 3 的不同环的数量
# 由于环的长度为 3，因此每个相同的环会重复统计 3 次，答案需除 3。
class Solution:
    def numberOfPaths(self, n: int, corridors: List[List[int]]) -> int:
        adj = defaultdict(set)

        for u, v in corridors:
            adj[u].add(v)
            adj[v].add(u)

        res = 0

        # 对邻居两两判断
        for p1 in range(1, n + 1):
            for p2, p3 in combinations(adj[p1], 2):
                if p2 in adj[p3]:
                    res += 1

        return res // 3


# Output: 2
print(Solution().numberOfPaths(n=5, corridors=[[1, 2], [5, 2], [4, 1], [2, 4], [3, 1], [3, 4]]))
