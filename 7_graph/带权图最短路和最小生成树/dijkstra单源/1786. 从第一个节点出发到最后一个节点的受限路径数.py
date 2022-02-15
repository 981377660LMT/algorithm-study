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


class Solution:
    def countRestrictedPaths(self, n: int, edges: List[List[int]]) -> int:
        adjMap = defaultdict(list)
        for u, v, w in edges:
            adjMap[u - 1].append((v - 1, w))
            adjMap[v - 1].append((u - 1, w))

        pq = [(0, n - 1)]
        dist = [0x7FFFFFFF] * (n - 1) + [0]
        while pq:
            _, cur = heappop(pq)
            for next, weight in adjMap[cur]:
                if dist[cur] + weight < dist[next]:
                    dist[next] = dist[cur] + weight
                    heappush(pq, (dist[next], next))

        @lru_cache(None)
        def dfs(cur: int) -> int:
            if cur == n - 1:
                return 1
            res = 0
            for next, _ in adjMap[cur]:
                if dist[cur] > dist[next]:
                    res += dfs(next)
            return res

        return dfs(0) % MOD


print(
    Solution().countRestrictedPaths(
        n=5, edges=[[1, 2, 3], [1, 3, 3], [2, 3, 1], [1, 4, 2], [5, 2, 2], [3, 5, 1], [5, 4, 10]]
    )
)
