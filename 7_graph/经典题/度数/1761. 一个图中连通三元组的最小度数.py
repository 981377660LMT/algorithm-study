from typing import List

# 一个 连通三元组 指的是 三个 节点组成的集合且这三个点之间 两两 有边。
# 请你返回所有连通三元组中度数的 最小值 ，如果图中没有连通三元组，那么返回 -1 。
# 2 <= n <= 400
# 看数据量不超过500 直接暴力枚举三个联通的点对

# Time complexity O(V^3)
# Space complexity O(V^2)

INF = 0x7FFFFFFF


class Solution:
    def minTrioDegree(self, n: int, edges: List[List[int]]) -> int:
        adjList = [[False] * n for _ in range(n)]
        degree = [0] * n
        for u, v in edges:
            adjList[u - 1][v - 1] = adjList[v - 1][u - 1] = True
            degree[u - 1] += 1
            degree[v - 1] += 1

        res = INF
        for i in range(n):
            for j in range(i + 1, n):
                if adjList[i][j]:
                    for k in range(j + 1, n):
                        if adjList[j][k] and adjList[k][i]:
                            res = min(res, degree[i] + degree[j] + degree[k] - 6)
        return res if res < INF else -1


print(Solution().minTrioDegree(n=6, edges=[[1, 2], [1, 3], [3, 2], [4, 1], [5, 2], [3, 6]]))
# 输出：3
# 解释：只有一个三元组 [1,2,3] 。构成度数的边在上图中已被加粗。
