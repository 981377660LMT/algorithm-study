from typing import List

# 对于 d 从 1 到 n-1 ，请你找到城市间 最大距离 恰好为 d 的所有`子树`数目。
# 请你返回一个大小为 n-1 的数组，其中第 d 个元素（下标从 1 开始）是城市间 最大距离 恰好等于 d 的子树数目。
# 2 <= n <= 15 可以枚举子集


class Solution:
    def countSubgraphsForEachDiameter(self, n: int, edges: List[List[int]]) -> List[int]:
        """
        1.求每个点到所有点的最短距离--多源最短路径算法 floyd O(n^3)
        2.枚举子集看哪些是子树
        """

        dist = [[int(1e20)] * n for _ in range(n)]
        for u, v in edges:
            dist[u - 1][v - 1] = 1
            dist[v - 1][u - 1] = 1

        for k in range(n):
            for i in range(n):
                for j in range(n):
                    dist[i][j] = min(dist[i][j], dist[i][k] + dist[k][j])

        res = [0] * n
        for state in range(1, 1 << n):
            points = []
            for i in range(n):
                if state & (1 << i):
                    points.append(i)

            edgeCount = 0
            vertexCount = len(points)
            maxDist = 0
            for i in range(vertexCount):
                for j in range(i + 1, vertexCount):
                    curDist = dist[points[i]][points[j]]
                    if curDist == 1:
                        edgeCount += 1
                    maxDist = max(maxDist, curDist)

            if vertexCount == edgeCount + 1:
                res[maxDist] += 1

        return res[1:]


print(Solution().countSubgraphsForEachDiameter(n=4, edges=[[1, 2], [2, 3], [2, 4]]))

