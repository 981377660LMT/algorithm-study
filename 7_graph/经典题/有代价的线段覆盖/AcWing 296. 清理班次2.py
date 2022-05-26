# https://www.acwing.com/solution/content/44273/
# 一张纸带，被划分成很多长度为 1 的网格，现在有 N 张贴纸，
# 第 i 张可以覆盖 ai 到 bi 范围内的格子，代价为 ci。
# 现在要用这些贴纸覆盖 [L,R] 的所有格子，需要最小化总代价。
# n<=1e4
# 0<=start,end<=1e5

# !最短路
# 把每一个网格看做有向图的顶点，从每个 ai 向 (bi+1) 连一条边权为 ci 的边。
# 然后从每个点向前一个点连一条边权为 0 的边。
# 求start到end的最短路即可

from collections import defaultdict
from heapq import heappop, heappush
from typing import DefaultDict


def dijkstra(adjMap: DefaultDict[int, DefaultDict[int, int]], start: int, end: int) -> int:
    dist = defaultdict(lambda: int(1e20))
    dist[start] = 0
    pq = [(0, start)]
    while pq:
        curDist, cur = heappop(pq)
        if dist[cur] < curDist:
            continue
        if cur == end:
            return curDist
        for next in adjMap[cur]:
            if dist[next] > dist[cur] + adjMap[cur][next]:
                dist[next] = dist[cur] + adjMap[cur][next]
                heappush(pq, (dist[next], next))

    return -1


n, start, end = map(int, input().split())
adjMap = defaultdict(lambda: defaultdict(lambda: int(1e20)))
for _ in range(n):
    a, b, cost = map(int, input().split())
    adjMap[a][b + 1] = min(adjMap[a][b + 1], cost)  # 注意b+1
for i in range(start, end):
    adjMap[i + 1][i] = min(adjMap[i + 1][i], 0)

print(dijkstra(adjMap, start, end + 1))  # 注意end+1

