from collections import defaultdict
from itertools import combinations, product
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

        adjMap = defaultdict(lambda: defaultdict(lambda: int(1e20)))
        for i in range(n):
            adjMap[i][i] = 0
        for u, v in edges:
            adjMap[u - 1][v - 1] = 1
            adjMap[v - 1][u - 1] = 1

        for k, i, j in product(range(n), repeat=3):
            adjMap[i][j] = min(adjMap[i][j], adjMap[i][k] + adjMap[k][j])

        res = [0] * n
        for state in range(1, 1 << n):
            nodes = [i for i in range(n) if (state >> i) & 1]
            edgeCount = sum(adjMap[n1][n2] == 1 for n1, n2 in combinations(nodes, 2))
            if len(nodes) == edgeCount + 1:
                maxDist = max((adjMap[n1][n2] for n1, n2 in combinations(nodes, 2)), default=0)
                res[maxDist] += 1

        return res[1:]


print(Solution().countSubgraphsForEachDiameter(n=4, edges=[[1, 2], [2, 3], [2, 4]]))
