# 类似 2045. 到达目的地的第二短时间-dijkstra变形
from heapq import heappop, heappush
from math import ceil
from typing import List, Tuple


INF = int(4e18)


def train(n: int, edges: List[Tuple[int, int, int, int]], START: int, END: int) -> int:
    """
    时间为k的倍数时,从都市u和v双向发车,花费时间w
    坐火车,求从START到END的最短时间
    """
    adjList = [[] for _ in range(n)]
    for u, v, w, k in edges:  # !时间为k的倍数时,从都市u和v双向发车,花费时间w
        adjList[u].append((v, w, k))
        adjList[v].append((u, w, k))

    dist = [INF] * n
    dist[START] = 0
    pq = [(0, START)]
    while pq:
        curDist, cur = heappop(pq)
        if cur == END:
            return curDist
        if curDist > dist[cur]:
            continue
        for next, weight, k in adjList[cur]:
            cand = ceil(curDist / k) * k + weight
            if cand < dist[next]:
                dist[next] = cand
                heappush(pq, (cand, next))
    return INF


if __name__ == "__main__":
    import sys

    sys.setrecursionlimit(int(1e6))
    input = lambda: sys.stdin.readline().rstrip("\r\n")
    n, m, start, end = map(int, input().split())
    start, end = start - 1, end - 1
    edges = []
    for _ in range(m):
        u, v, w, k = map(int, input().split())
        u, v = u - 1, v - 1
        edges.append((u, v, w, k))

    minCost = train(n, edges, start, end)
    print(minCost if minCost != INF else -1)
