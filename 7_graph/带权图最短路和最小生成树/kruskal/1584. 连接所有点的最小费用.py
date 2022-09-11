from itertools import combinations
from typing import List
from 模板 import kruskal

# 1 <= points.length <= 1000
# 1584. 连接所有点的最小费用


class Solution:
    def minCostConnectPoints(self, points: List[List[int]]) -> int:
        n = len(points)
        edges = []
        P = [tuple(p) for p in points]
        for p1, p2 in combinations(P, 2):
            edges.append((p1, p2, abs(p1[0] - p2[0]) + abs(p1[1] - p2[1])))
        return kruskal(n, edges)[0]
