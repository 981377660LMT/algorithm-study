from collections import defaultdict
from heapq import heappop, heappush
from typing import DefaultDict, Hashable, Optional, TypeVar, overload
import sys
import os

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)


INF = int(1e18)
Vertex = TypeVar("Vertex", bound=Hashable)
Graph = DefaultDict[Vertex, DefaultDict[Vertex, int]]


def dijkstra(adjMap: Graph[Vertex], start: Vertex, end: Optional[Vertex] = None):
    """时间复杂度O((V+E)logV)"""
    dist = defaultdict(lambda: INF, {start: 0})
    pq = [(0, start)]
    pre = dict()

    while pq:
        curDist, cur = heappop(pq)
        if dist[cur] < curDist:  # 剪枝，有的题目不加就TLE
            continue
        if end is not None and cur == end:
            return curDist
        for next in adjMap[cur]:
            if dist[next] > dist[cur] + adjMap[cur][next]:
                dist[next] = dist[cur] + adjMap[cur][next]
                heappush(pq, (dist[next], next))
                pre[next] = cur

    return pre


def main() -> None:
    n, m = map(int, input().split())
    adjMap = defaultdict(lambda: defaultdict(lambda: int(1e18)))
    edges = {}
    for i in range(m):
        u, v, w = map(int, input().split())
        adjMap[u][v] = w
        adjMap[v][u] = w
        edges[tuple(sorted((u, v)))] = i + 1

    pre = dijkstra(adjMap, 1)
    for j in range(2, n + 1):
        i = pre[j]  # type: ignore
        print(edges[tuple(sorted((i, j)))])


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
