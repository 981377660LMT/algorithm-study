from typing import List
from collections import defaultdict
from heapq import heappop, heappush
from functools import lru_cache

# 用 distanceToLastNode(x) 表示节点 n 和 x 之间路径的最短距离
# 受限路径 为满足 distanceToLastNode(cur) > distanceToLastNode(next) 的一条路径 (cur,next对任何边)
# 返回从节点 1 出发到节点 n 的 受限路径数 。由于数字可能很大，请返回对 109 + 7 取余 的结果。

# 1.dijkstra求出各个点到n的最短距离
# 2.dfs求出u到n的受限路径数(每个点相邻的最短距离小于自身的节点)


MOD = int(1e9 + 7)
INF = int(1e20)


class Solution:
    def countRestrictedPaths(self, n: int, edges: List[List[int]]) -> int:
        """返回从节点 1 出发到节点 n 的 受限路径数"""
        adjMap = defaultdict(lambda: defaultdict(lambda: INF))
        for u, v, w in edges:
            adjMap[u][v] = w
            adjMap[v][u] = w

        dist = defaultdict(lambda: INF, {n: 0})
        pq = [(0, n)]
        while pq:
            curDist, cur = heappop(pq)
            if dist[cur] < curDist:
                continue
            for next, weight in adjMap[cur].items():
                if dist[cur] + weight < dist[next]:
                    dist[next] = dist[cur] + weight
                    heappush(pq, (dist[next], next))

        @lru_cache(None)
        def dfs(cur: int) -> int:
            if cur == n:
                return 1
            res = 0
            for next in adjMap[cur]:
                if dist[cur] > dist[next]:
                    res += dfs(next)
                    res %= MOD
            return res

        return dfs(1)


print(
    Solution().countRestrictedPaths(
        n=5, edges=[[1, 2, 3], [1, 3, 3], [2, 3, 1], [1, 4, 2], [5, 2, 2], [3, 5, 1], [5, 4, 10]]
    )
)
