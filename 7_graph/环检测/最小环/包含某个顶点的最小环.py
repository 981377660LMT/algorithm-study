from heapq import heappop, heappush
from typing import List, Tuple


INF = int(4e18)


def min2(a: int, b: int) -> int:
    return a if a < b else b


def minCycleThroughVertex(adjList: List[List[Tuple[int, int]]], vertex: int) -> int:
    """经过指定顶点的最小环."""
    n = len(adjList)
    dist = [INF] * n
    dist[vertex] = 0
    pq = [(0, vertex)]
    while pq:
        curDist, cur = heappop(pq)
        if curDist > dist[cur]:
            continue
        for next, weight in adjList[cur]:
            cand = curDist + weight
            if cand < dist[next]:
                dist[next] = cand
                heappush(pq, (cand, next))
    res = INF
    for cur in range(n):
        for next, weight in adjList[cur]:
            if next == vertex:
                res = min2(res, dist[cur] + weight)
    return res


if __name__ == "__main__":
    # D - Cycle
    # https://atcoder.jp/contests/abc376/tasks/abc376_d
    N, M = map(int, input().split())
    adjList = [[] for _ in range(N)]
    for _ in range(M):
        u, v = map(int, input().split())
        u -= 1
        v -= 1
        adjList[u].append((v, 1))
    res = minCycleThroughVertex(adjList, 0)
    print(res if res != INF else -1)
