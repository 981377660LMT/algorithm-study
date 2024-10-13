# 一定在最短路径上的边
# !无向图最短路割边
#
# 对于任意一条边(u, v, w), 如果满足
# dist1[u] + w + dist2[v] == dist1[target] 且 count1[u] * count2[v] % MOD == count1[target]
# 或者
# dist1[v] + w + dist2[u] == dist1[target] 且 count1[v] * count2[u] % MOD == count1[target]
# 那么这条边一定在从start到target的最短路径上.

from heapq import heappop, heappush
from typing import List, Tuple


INF = int(4e18)
MOD = int(1e9 + 7)


def edgesMustOnShortestPath(
    n: int, edges: List[Tuple[int, int, int]], start: int, target: int
) -> List[bool]:
    """
    给定一个n个点m条边的无向带权图.
    对于每条边(u, v, w), 判断是否一定在从start到target的最短路径上.
    """
    adjList = [[] for _ in range(n)]
    for u, v, w in edges:
        adjList[u].append((v, w))
        adjList[v].append((u, w))
    dist1, count1 = dijkstraWithCount(n, adjList, start, mod=MOD)
    dist2, count2 = dijkstraWithCount(n, adjList, target, mod=MOD)
    res = [False] * len(edges)
    for i, (u, v, w) in enumerate(edges):
        ok1 = (
            dist1[u] + w + dist2[v] == dist1[target]
            and count1[u] * count2[v] % MOD == count1[target]
        )
        if ok1:
            res[i] = True
            continue
        ok2 = (
            dist1[v] + w + dist2[u] == dist1[target]
            and count1[v] * count2[u] % MOD == count1[target]
        )
        if ok2:
            res[i] = True

    return res


def dijkstraWithCount(
    n: int, adjList: List[List[Tuple[int, int]]], start: int, *, mod=int(1e9 + 7)
) -> Tuple[List[int], List[int]]:
    """dijkstra求出起点到各点的最短距离和最短路径数量(模mod).

    时间复杂度O((V+E)logV).
    """
    dist = [INF] * n
    count = [0] * n
    dist[start] = 0
    count[start] = 1
    pq = [(0, start)]
    while pq:
        curDist, cur = heappop(pq)
        if dist[cur] < curDist:
            continue
        for next, weight in adjList[cur]:
            cand = dist[cur] + weight
            if cand < dist[next]:
                dist[next] = cand
                count[next] = count[cur]
                heappush(pq, (dist[next], next))
            elif cand == dist[next]:
                count[next] += count[cur]
                if count[next] >= mod:
                    count[next] -= mod
    return dist, count


if __name__ == "__main__":
    # G - Road Blocked 2
    # https://atcoder.jp/contests/abc375/tasks/abc375_g

    import sys

    input = lambda: sys.stdin.readline().rstrip("\r\n")

    n, m = map(int, input().split())
    edges = []
    for _ in range(m):
        u, v, w = map(int, input().split())
        u, v = u - 1, v - 1
        edges.append((u, v, w))

    res = edgesMustOnShortestPath(n, edges, 0, n - 1)
    for ok in res:
        print("Yes" if ok else "No")
