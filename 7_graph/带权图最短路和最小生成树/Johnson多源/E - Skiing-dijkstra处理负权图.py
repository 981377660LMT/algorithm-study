# 给出n个点,m条无向边,每个点都有个高度h[i]
# 对于一条边,如果从h[u]> h[v],则经过这条边的幸福度为h[u]- h[v],否则为-2*(h[(v] - h[u])
# 问从1开始的`最长路(最大幸福度)`。


# !n<=2e5 m<=min(2e5,n*(n-1)/2)
# !稀疏图

"""
dijk求最长路
最长路，可以通过边权取反后跑最短路，不过取反后的边权是负数

添加点势能的dijkstra算法
为了让每条边边权>=0 负图中 u->v 边权需要加上 h[u]-h[v] (h[i]表示i点的势能)
那么负图的边权为max(0,h[v]-h[u])

最后答案为 - (dist[u] - (h[start] - h[u]))
"""

from collections import defaultdict
from heapq import heappop, heappush

import sys


sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


n, m = map(int, input().split())
heights = list(map(int, input().split()))
rAdjMap = defaultdict(lambda: defaultdict(lambda: INF))
for _ in range(m):
    u, v = map(int, input().split())
    u, v = u - 1, v - 1
    rAdjMap[u][v] = max(0, heights[v] - heights[u])
    rAdjMap[v][u] = max(0, heights[u] - heights[v])

dist = [INF] * n
dist[0] = 0
pq = [(0, 0)]
while pq:
    curDist, cur = heappop(pq)
    if dist[cur] < curDist:
        continue
    for next, weight in rAdjMap[cur].items():
        if dist[next] > dist[cur] + weight:
            dist[next] = dist[cur] + weight
            heappush(pq, (dist[next], next))

res = 0
for i in range(n):
    res = max(res, -(dist[i] - (heights[0] - heights[i])))  # 还原势能
print(res)
