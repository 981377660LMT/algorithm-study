from heapq import heappop, heappush
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个二维数组 edges 表示一个 n 个点的无向图，其中 edges[i] = [ui, vi, lengthi] 表示节点 ui 和节点 vi 之间有一条需要 lengthi 单位时间通过的无向边。

# 同时给你一个数组 disappear ，其中 disappear[i] 表示节点 i 从图中消失的时间点，在那一刻及以后，你无法再访问这个节点。

# 注意，图有可能一开始是不连通的，两个节点之间也可能有多条边。


# 请你返回数组 answer ，answer[i] 表示从节点 0 到节点 i 需要的 最少 单位时间。如果从节点 0 出发 无法 到达节点 i ，那么 answer[i] 为 -1 。


class Solution:
    def minimumTime(self, n: int, edges: List[List[int]], disappear: List[int]) -> List[int]:
        adjList = [[] for _ in range(n)]
        for u, v, w in edges:
            adjList[u].append((v, w))
            adjList[v].append((u, w))

        dist = [INF] * n
        dist[0] = 0
        pq = [(0, 0)]

        while pq:
            curDist, cur = heappop(pq)
            if dist[cur] < curDist:
                continue
            for next, weight in adjList[cur]:
                cand = dist[cur] + weight
                if cand < dist[next] and disappear[next] > cand:
                    dist[next] = cand
                    heappush(pq, (dist[next], next))
        dist = [-1 if x == INF else x for x in dist]
        return dist
