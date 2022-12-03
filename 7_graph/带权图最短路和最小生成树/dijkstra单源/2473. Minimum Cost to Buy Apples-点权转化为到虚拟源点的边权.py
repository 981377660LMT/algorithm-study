# 购买苹果的最小花费
# 2473. Minimum Cost to Buy Apples-虚拟源点

# Description
# 给定一个无向带权图
# !求出每个点 i 到其他点 j 的 dist[i][j] * (k+1) + appleCost[j] 的最小值

# O((V+E)log(V+E)) Solution:
# 点权转化为到虚拟源点的边权
# !添加一个虚拟原点 到每个点的距离为appleCost[i] (虚拟原点到城市的距离就是苹果的价格)
# !那么就可以将 `每个点到其他点的最短路` 转化为 `虚拟原点到其他点的最短路`


from typing import List, Sequence, Tuple
from heapq import heappop, heappush

INF = int(1e18)


class Solution:
    def minCost(self, n: int, roads: List[List[int]], appleCost: List[int], k: int) -> List[int]:
        START = n  # 虚拟源点
        adjList = [[] for _ in range(n + 1)]
        for u, v, w in roads:
            u, v = u - 1, v - 1
            adjList[u].append((v, w * (k + 1)))
            adjList[v].append((u, w * (k + 1)))
        for i in range(n):
            adjList[START].append((i, appleCost[i]))
            adjList[i].append((START, appleCost[i]))
        dist = dijkstra(n + 1, adjList, START)
        return dist[:-1]


def dijkstra(n: int, adjList: Sequence[Sequence[Tuple[int, int]]], start: int) -> List[int]:
    dist = [INF] * n
    dist[start] = 0
    pq = [(0, start)]

    while pq:
        curDist, cur = heappop(pq)
        if dist[cur] < curDist:
            continue
        for next, weight in adjList[cur]:
            cand = dist[cur] + weight
            if cand < dist[next]:
                dist[next] = cand
                heappush(pq, (dist[next], next))
    return dist


print(
    Solution().minCost(
        n=4,
        roads=[[1, 2, 4], [2, 3, 2], [2, 4, 5], [3, 4, 1], [1, 3, 4]],
        appleCost=[56, 42, 102, 301],
        k=2,
    )
)
