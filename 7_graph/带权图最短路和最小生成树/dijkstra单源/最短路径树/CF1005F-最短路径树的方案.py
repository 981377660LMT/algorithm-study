# Berland and the Shortest Paths
# 求最短路径树的方案数,并输出方案

from heapq import heappop, heappush
from typing import List, Sequence, Tuple

INF = int(1e18)


def solve(n: int, edges: List[Tuple[int, int, int]], k: int, root: int) -> None:
    """求最短路径树的方案数,并输出方案(k种)"""

    adjList = [[] for _ in range(n)]
    for i, (u, v, w) in enumerate(edges):
        adjList[u].append((v, w, i))
        adjList[v].append((u, w, i))
    *_, preE = dijkstraSPT(n, adjList, root)

    res = 1
    for i in range(n):
        if i == root:
            continue
        res *= len(preE[i])
        if res >= k:
            res = k
            break
    print(res)  # 输出方案数

    visited = [0] * len(edges)
    count = 0

    # dfs输出所有方案(为每个顶点选择前驱边)
    def dfs(step: int) -> None:
        if step == root:
            dfs(step + 1)
            return
        if step == n:
            print("".join(map(str, visited)))
            nonlocal count
            count += 1
            if count >= res:
                exit(0)
            return
        for e in preE[step]:
            visited[e] = 1
            dfs(step + 1)
            visited[e] = 0

    dfs(0)


# !dijkstra求出路径上每个点的前驱边
# 其中邻接表的每个元素是一个三元组，分别是邻接点，边权，边的编号
def dijkstraSPT(
    n: int, adjList: Sequence[Sequence[Tuple[int, int, int]]], start: int
) -> Tuple[List[int], List[List[int]]]:
    dist = [INF] * n
    preE = [[] for _ in range(n)]
    dist[start] = 0
    pq = [(0, start)]
    while pq:
        curDist, cur = heappop(pq)
        if dist[cur] < curDist:
            continue
        for next, weight, eid in adjList[cur]:
            cand = dist[cur] + weight
            if cand < dist[next]:
                dist[next] = cand
                heappush(pq, (dist[next], next))
                preE[next].append(eid)
            elif cand == dist[next]:
                preE[next].append(eid)
    return dist, preE


if __name__ == "__main__":
    import sys

    sys.setrecursionlimit(int(1e9))
    input = lambda: sys.stdin.readline().rstrip("\r\n")

    # 你需要找到 k 种不同的可行解，如果解的个数少于 k 种，你需要输出所有可行解。
    n, m, k = map(int, input().split())
    edges = []
    for _ in range(m):
        u, v = map(int, input().split())
        edges.append((u - 1, v - 1, 1))
    solve(n, edges, k, 0)
