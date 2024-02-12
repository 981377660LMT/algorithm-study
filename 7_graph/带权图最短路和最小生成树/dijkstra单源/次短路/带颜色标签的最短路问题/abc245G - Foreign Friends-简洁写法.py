# https://atcoder.jp/contests/abc245/tasks/abc245_g
# 给 N 个点，M 条边，每个点的颜色（值域为 [1,K]）。并给定 L 个点作为特殊点。
# 询问 每个点到最近的与其颜色不同的特殊点的距离（无解输出 -1） 。
# !这个次短路指的是“前两个不同节点到达这里的最短路”

from heapq import heappop, heappush
from typing import List, Tuple


INF = int(4e18)


def min2(a: int, b: int) -> int:
    return a if a < b else b


def abc245g(
    n: int, edges: List[Tuple[int, int, int]], colors: List[int], criticals: List[int]
) -> List[int]:
    adjList = [[] for _ in range(n)]
    for u, v, w in edges:
        adjList[u].append((v, w))
        adjList[v].append((u, w))

    dist = [INF] * n
    source1, source2 = [-1] * n, [-1] * n
    pq = [(0, v, colors[v]) for v in criticals]

    while pq:
        curDist, cur, curSource = heappop(pq)
        if curSource == source1[cur] or curSource == source2[cur]:
            continue
        if source1[cur] == -1:
            source1[cur] = curSource
        elif source2[cur] == -1:
            source2[cur] = curSource
        else:
            continue

        if curSource != colors[cur]:  # 出发点颜色与自己颜色不同时，更新距离
            dist[cur] = min2(dist[cur], curDist)
        for next, weight in adjList[cur]:
            heappush(pq, (curDist + weight, next, curSource))  # type: ignore

    return [dist[v] if dist[v] != INF else -1 for v in range(n)]


if __name__ == "__main__":
    import sys

    input = lambda: sys.stdin.readline().rstrip("\r\n")

    N, M, K, L = map(int, input().split())
    colors = [v - 1 for v in map(int, input().split())]
    criticals = [v - 1 for v in map(int, input().split())]
    edges = []
    for _ in range(M):
        u, v, w = map(int, input().split())
        edges.append((u - 1, v - 1, w))

    print(*abc245g(N, edges, colors, criticals))
