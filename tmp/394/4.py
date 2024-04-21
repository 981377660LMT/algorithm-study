from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个 n 个节点的无向带权图，节点编号为 0 到 n - 1 。图中总共有 m 条边，用二维数组 edges 表示，其中 edges[i] = [ai, bi, wi] 表示节点 ai 和 bi 之间有一条边权为 wi 的边。

# 对于节点 0 为出发点，节点 n - 1 为结束点的所有最短路，你需要返回一个长度为 m 的 boolean 数组 answer ，如果 edges[i] 至少 在其中一条最短路上，那么 answer[i] 为 true ，否则 answer[i] 为 false 。

# 请你返回数组 answer 。


# 注意，图可能不连通。
from heapq import heappop, heappush
from typing import List, Sequence, Tuple

from typing import List, Sequence, Tuple
from heapq import heappop, heappush


def dijkstra(n: int, adjList: Sequence[Sequence[Tuple[int, int, int]]], start: int) -> List[int]:
    dist = [INF] * n
    dist[start] = 0
    pq = [(0, start)]

    while pq:
        curDist, cur = heappop(pq)
        if dist[cur] < curDist:
            continue
        for next, weight, _ in adjList[cur]:
            cand = dist[cur] + weight
            if cand < dist[next]:
                dist[next] = cand
                heappush(pq, (dist[next], next))
    return dist


# 可能在最短路径上的边
# https://zhuanlan.zhihu.com/p/166173790
# 一定在最短路径上的边
# https://www.luogu.com.cn/problem/CF567E
# 删边最短路
# P3238 [HNOI2014] 道路堵塞
# https://www.luogu.com.cn/problem/P3238
# 题意:
# 给你一个n个点m条边的有向带权图，问你删去每一条边后从0到n-1的最短路长度是多少。
# !本题可能不存在正确解法，题解均已被 hack。
# n<= 1e5, m <= 2e5


class Solution:
    def findAnswer(self, n: int, edges: List[List[int]]) -> List[bool]:
        adjList = [[] for _ in range(n)]
        for i, (u, v, w) in enumerate(edges):
            adjList[u].append((v, w, i))
            adjList[v].append((u, w, i))
        dist1 = dijkstra(n, adjList, 0)
        dist2 = dijkstra(n, adjList, n - 1)
        res = [False] * len(edges)
        for i, (u, v, w) in enumerate(edges):
            if dist1[u] + w + dist2[v] == dist1[n - 1] or dist1[v] + w + dist2[u] == dist1[n - 1]:
                res[i] = True
        return res


# n = 4, edges = [[2,0,1],[0,1,1],[0,3,4],[3,2,2]]

print(Solution().findAnswer(4, [[2, 0, 1], [0, 1, 1], [0, 3, 4], [3, 2, 2]]))
