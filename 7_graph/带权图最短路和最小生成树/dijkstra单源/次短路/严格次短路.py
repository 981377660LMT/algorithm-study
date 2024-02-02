# https://www.luogu.com.cn/problem/P2865
# P2865 [USACO06NOV] Roadblocks G
# !1.如果 dist1[next] > dist1[cur] + d(cur,next)，则更新 dist1[next]
# !2.如果 dist1[next] < dist1[cur] + d(cur,next) < dist2[next]，则更新 dist2[next]
#    注意不能取等，否则dist2[cur]和dist1[cur]可能相等

from heapq import heappop, heappush
from typing import Tuple, List


INF = int(1e18)


def p2865(n: int, edges: List[Tuple[int, int, int]], start: int) -> Tuple[List[int], List[int]]:
    """
    给定一个无向带权图,求start到其他点的最短路和严格次短路.
    不存在则为INF.
    """
    adjList = [[] for _ in range(n)]
    for u, v, w in edges:
        adjList[u].append((v, w))
        adjList[v].append((u, w))
    dist1, dist2 = [INF] * n, [INF] * n
    dist1[start] = 0
    pq = [(0, start)]
    while pq:
        curDist, cur = heappop(pq)
        if curDist > dist2[cur]:  # !注意是dist2
            continue
        for next, weight in adjList[cur]:
            cand = curDist + weight
            if cand < dist1[next]:
                dist1[next], cand = cand, dist1[next]  # !注意 cand 被更新了
                heappush(pq, (dist1[next], next))
            if dist1[next] < cand < dist2[next]:  # dist1[next] < cand ，严格次短路
                dist2[next] = cand
                heappush(pq, (dist2[next], next))
    return dist1, dist2


if __name__ == "__main__":
    import sys

    input = lambda: sys.stdin.readline().rstrip("\r\n")

    INF = int(4e18)
    n, m = map(int, input().split())
    edges = []
    for _ in range(m):
        u, v, w = map(int, input().split())
        edges.append((u - 1, v - 1, w))

    _, dist2 = p2865(n, edges, 0)
    print(dist2[n - 1])
