# 1 <= points.length <= 1000
# 1584. 连接所有点的最小费用
# O(nlogn)曼哈顿距离最小生成树
# https://leetcode.cn/problems/min-cost-to-connect-all-points/solution/lian-jie-suo-you-dian-de-zui-xiao-fei-yo-kcx7/


from typing import List
from 模板 import kruskal


class Solution:
    def minCostConnectPoints(self, points: List[List[int]]) -> int:
        """!O(n^2logn^2) n为顶点数"""
        n = len(points)
        edges = []
        P = [(x, y, i) for i, (x, y) in enumerate(points)]
        for i in range(n):
            for j in range(i + 1, n):
                x1, y1, i1 = P[i]
                x2, y2, i2 = P[j]
                edges.append((i1, i2, abs(x1 - x2) + abs(y1 - y2)))
        return kruskal(n, edges)[0]


print(Solution().minCostConnectPoints([[0, 0], [2, 2], [3, 10], [5, 2], [7, 0]]))  # 20
