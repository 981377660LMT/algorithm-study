from collections import defaultdict
from typing import List

# 一个 连通三元组 指的是 三个 节点组成的集合且这三个点之间 两两 有边。
# 请你返回所有连通三元组中度数的 最小值 ，如果图中没有连通三元组，那么返回 -1 。
# 2 <= n <= 400
# 看数据量不超过500 直接暴力枚举三个联通的点对

# Time complexity O(V^3)
# Space complexity O(V^2)

INF = int(1e20)


class Solution:
    def minTrioDegree1(self, n: int, edges: List[List[int]]) -> int:
        """时间复杂度O(n^3)"""
        adjList = [[False] * n for _ in range(n)]
        deg = [0] * n
        for u, v in edges:
            adjList[u - 1][v - 1] = adjList[v - 1][u - 1] = True
            deg[u - 1] += 1
            deg[v - 1] += 1

        res = INF
        for i in range(n):
            for j in range(i + 1, n):
                if adjList[i][j]:
                    for k in range(j + 1, n):
                        if adjList[j][k] and adjList[k][i]:
                            res = min(res, deg[i] + deg[j] + deg[k] - 6)
        return res if res < INF else -1

    # 给无向图定向减少重复枚举次数
    # 边的方向定为从度数小的点连向度数大的点
    # https://leetcode-cn.com/problems/minimum-degree-of-a-connected-trio-in-a-graph/solution/gei-wu-xiang-tu-ding-xiang-by-lucifer100-c72d/
    def minTrioDegree(self, n: int, edges: List[List[int]]) -> int:
        """三元环计数 n<=400 稠密图

        边定向 度数为第一关键字，点的标号为第二关键字
        此时任意点的出度不会超过 (sqrt(2E))
        时间复杂度O(E^(3/2))
        """
        deg = [0] * n
        for u, v in edges:
            deg[u - 1] += 1
            deg[v - 1] += 1

        adjMap = [set() for _ in range(n)]
        for u, v in edges:
            u, v = u - 1, v - 1
            # u, v = sorted((u, v), key=lambda x: (deg[x], x))
            if deg[u] > deg[v] or deg[u] == deg[v] and u > v:
                u, v = v, u
            adjMap[u].add(v)

        res = int(1e20)
        for i in range(n):
            for j in adjMap[i]:
                for k in adjMap[j]:
                    if k in adjMap[i]:
                        res = min(res, deg[i] + deg[j] + deg[k] - 6)

        return res if res < int(1e20) else -1


print(Solution().minTrioDegree(n=6, edges=[[1, 2], [1, 3], [3, 2], [4, 1], [5, 2], [3, 6]]))
# 输出：3
# 解释：只有一个三元组 [1,2,3] 。构成度数的边在上图中已被加粗。
