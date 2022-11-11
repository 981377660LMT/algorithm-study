# https://oi-wiki.org/graph/min-circle/#floyd
from collections import defaultdict
from heapq import heappop, heappush
from typing import DefaultDict, List

INF = int(1e20)


def minCycle1(n: int, adjMatrix: List[List[int]]) -> int:
    """floyd 求带权图最小环 O(n^3)

    Args:
        n (int): 顶点数
        adjMatrix (List[List[int]]): 表示距离的邻接矩阵,adjMatrix[u][v] == INF 表示没有 u->v 边

    Returns:
        int: 最小环长度, INF 表示不存在

    Notes:
        https://oi-wiki.org/graph/min-circle/#floyd

        1. 在最外层循环到点 k 时（尚未开始第 k 次循环），最短路数组 dist 中，
        表示的是从 i 到 j 且仅经过编号在 [0,k) 区间中的点的最短路。
        2. 由最小环的定义可知其至少有三个顶点,设其中编号最大的顶点为 v3 ,环上与 v3 相邻两侧的两个点为v1/v2,
        则在最外层循环枚举到 k=v3 时，该环的长度即 dist[v1][v2] + weights[v1][k] + weights[k][v2]。
    """
    distCopy = [row[:] for row in adjMatrix]
    res = INF
    for k in range(n):
        for i in range(k):
            for j in range(i):
                cand = distCopy[i][j] + adjMatrix[i][k] + adjMatrix[k][j]
                if cand < res:
                    res = cand
        for i in range(n):
            for j in range(n):
                cand = distCopy[i][k] + distCopy[k][j]
                if cand < distCopy[i][j]:
                    distCopy[i][j] = cand
    return res


def minCycle2(n: int, edges: List[List[int]]) -> int:
    """枚举所有边，删除这条边之后以这条边的端点为起点终点跑一次 Dijkstra

    O(E*E*Log(V+E))
    """

    def dijkstra(
        n: int, adjMap: DefaultDict[int, DefaultDict[int, int]], start: int, end: int
    ) -> int:
        """时间复杂度O((V+E)logV)"""
        dist = [INF] * n
        dist[start] = 0
        pq = [(0, start)]
        while pq:
            curDist, cur = heappop(pq)
            if dist[cur] < curDist:
                continue
            if cur == end:
                return curDist
            for next in adjMap[cur]:
                cand = curDist + adjMap[cur][next]
                if dist[next] > cand:
                    dist[next] = cand
                    heappush(pq, (cand, next))
        return INF

    adjMap = defaultdict(lambda: defaultdict(lambda: INF))
    for u, v, w in edges:
        adjMap[u][v] = min(adjMap[u][v], w)
        adjMap[v][u] = min(adjMap[v][u], w)

    res = INF
    for u, v, _ in edges:
        tmp = adjMap[u][v]
        adjMap[u][v] = INF
        adjMap[v][u] = INF
        res = min(res, dijkstra(n, adjMap, u, v) + tmp)
        adjMap[u][v] = tmp
        adjMap[v][u] = tmp

    return res


if __name__ == "__main__":
    n = 3
    adjMatrix = [
        [0, 1, 3],
        [1, 0, 1],
        [3, 1, 0],
    ]
    assert minCycle1(n, adjMatrix) == 5
    edges = [
        [0, 1, 1],
        [1, 2, 1],
        [2, 0, 3],
    ]
    assert minCycle2(n, edges) == 5
