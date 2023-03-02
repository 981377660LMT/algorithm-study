# https://www.acwing.com/problem/content/description/3238/
# 将现有的一部分铁路改造成高速铁路，使得任何两个城市间都可以通过高速铁路到达，
# !而且从所有城市乘坐高速铁路到首都的最短路程和原来一样长。
# 在这些条件下最少要改造多长的铁路。
# 所有的城市由 1 到 n 编号，首都为 1 号。
# 1≤n≤1e4 ，1≤m≤1e5


# 翻译：
# !给我们一个无向图，我们需要从无向图中选出来若干条边，使得选出来边的总权值和最小，
# !并且满足所有点到1号点的最短距离不变。
# 是最短路径树SPT
# !最短路径树SPT（Short Path Tree）是网络的源点到所有结点的最短路径构成的树。
# 求出最短路径树的边权和即可


from heapq import heappop, heappush
from typing import List, Sequence, Tuple

INF = int(1e18)


def shortestPathTree(
    n: int, edges: List[Tuple[int, int, int]], start: int
) -> Tuple[int, List[int]]:
    """求出最短路径树的最小边权之和和每条边的编号"""
    adjList = [[] for _ in range(n)]
    for i, (u, v, w) in enumerate(edges):
        adjList[u].append((v, w, i))
        adjList[v].append((u, w, i))
    *_, preE = dijkstraSPT(n, adjList, start)

    res, eids = 0, []
    for i in range(n):
        if i == start:
            continue
        pre = preE[i]
        res += edges[pre][2]
        eids.append(pre)
    return res, eids


# !dijkstra求出路径上每个点的前驱点和前驱边
# 其中邻接表的每个元素是一个三元组，分别是邻接点，边权，边的编号
def dijkstraSPT(
    n: int, adjList: Sequence[Sequence[Tuple[int, int, int]]], start: int
) -> Tuple[List[int], List[int], List[int]]:
    dist = [INF] * n
    preV, preE = [-1] * n, [-1] * n
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
                preE[next] = eid
                preV[next] = cur
                heappush(pq, (dist[next], next))
            elif cand == dist[next]:  # 在最短路相等的情况下，扩展到同一个节点，后出堆的点连的边权值一定更小
                preE[next] = eid
                preV[next] = cur
    return dist, preV, preE


if __name__ == "__main__":
    n, m = map(int, input().split())
    edges = []
    for i in range(m):
        u, v, w = map(int, input().split())
        edges.append((u - 1, v - 1, w))
    print(shortestPathTree(n, edges, 0)[0])
