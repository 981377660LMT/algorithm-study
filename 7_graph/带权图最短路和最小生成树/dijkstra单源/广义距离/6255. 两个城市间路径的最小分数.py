# 两个城市之间一条路径的 分数 定义为这条路径中道路的 最小 距离。
# !城市 1 和城市 n 之间的`所有路径`的 最小 分数。
# !因为是所有路径，所以是到连通分量中`任意一个点`的最小分数
# 边权:取最小值
# 当然更好的做法是并查集

from heapq import heappop
from typing import List

INF = int(1e18)


class Solution:
    def minScore(self, n: int, roads: List[List[int]]) -> int:
        adjList = [[] for _ in range(n)]
        for u, v, w in roads:
            u, v = u - 1, v - 1
            adjList[u].append((v, w))
            adjList[v].append((u, w))

        dist = [INF] * n
        dist[0] = INF
        pq = [(INF, 0)]  # (dist, cur)
        while pq:
            curDist, cur = heappop(pq)
            if curDist > dist[cur]:
                continue
            for next, weight in adjList[cur]:
                cand = min(curDist, weight)  # !边权:取最小值
                if cand < dist[next]:
                    dist[next] = cand
                    pq.append((cand, next))

        return min(dist)  # 到连通分量中任意一个点的最小分数


print(Solution().minScore(n=4, roads=[[1, 2, 2], [1, 3, 4], [3, 4, 7]]))
