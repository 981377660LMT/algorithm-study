# 魔法饰品
# bfs+状压dp
# n<=1e5 k<=17
# !包含所有元素的最短路径和(路径可以不连续)
# !建图，用bfs预处理出k种颜元素两两之间的最短距离，再用状压DP求解

from functools import lru_cache
from typing import List, Tuple
from collections import deque

INF = int(1e18)


def magicalOrnament(n: int, edges: List[Tuple[int, int]], need: List[int]) -> int:
    @lru_cache(None)
    def dfs(index: int, visited: int) -> int:
        if visited == (1 << k) - 1:
            return 0
        res = INF
        for next in range(k):
            if visited & (1 << next):
                continue
            res = min(res, dist[index][next] + dfs(next, visited | (1 << next)))
        return res

    adjList = [[] for _ in range(n)]
    for u, v in edges:
        adjList[u].append(v)
        adjList[v].append(u)

    dist = []  # dist[i][j]表示i,j两种颜元素之间的最短距离
    for start in need:
        curDist = bfs(start, adjList)
        dist.append([curDist[i] for i in need])

    res = min(dfs(i, 1 << i) for i in range(k))  # 枚举起点
    dfs.cache_clear()
    return res + 1 if res != INF else -1


def bfs(start: int, adjList: List[List[int]]) -> List[int]:
    n = len(adjList)
    dist = [INF] * n
    dist[start] = 0
    queue = deque([start])
    while queue:
        cur = queue.popleft()
        for next in adjList[cur]:
            cand = dist[cur] + 1
            if cand < dist[next]:
                dist[next] = cand
                queue.append(next)
    return dist


if __name__ == "__main__":
    import sys

    sys.setrecursionlimit(int(1e6))
    input = lambda: sys.stdin.readline().rstrip("\r\n")

    n, m = map(int, input().split())
    edges = []
    for _ in range(m):
        u, v = map(int, input().split())
        u, v = u - 1, v - 1
        edges.append((u, v))
    k = int(input())
    need = [int(x) - 1 for x in input().split()]
    print(magicalOrnament(n, edges, need))
