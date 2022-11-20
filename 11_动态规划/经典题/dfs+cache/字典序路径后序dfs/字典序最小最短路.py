# 字典序最小最短路/字典序最短路
# !这里是求路径顶点字典序最小的最短路

from heapq import heappop, heappush
from typing import List, Tuple

INF = int(1e18)


def lexicographicallySmallestShortestPath(
    n: int, graph: List[List[Tuple[int, int]]], start: int, end: int
) -> List[int]:
    """求路径顶点字典序最小的最短路
    :return: 最短路的顶点列表
    """
    dist = [INF] * n
    dist[end] = 0  # !倒序保证每个后缀都是字典序最小的
    pre = [-1] * n  # pre[i]表示i的前驱节点(字典序要最小)
    pq = [(0, end)]
    while pq:
        curDist, cur = heappop(pq)
        if curDist > dist[cur]:
            continue
        for next, weight in graph[cur]:
            cand = curDist + weight
            if cand < dist[next]:
                dist[next] = cand
                pre[next] = cur
                heappush(pq, (cand, next))
            elif cand == dist[next] and pre[next] > cur:
                pre[next] = cur

    path, cur = [], start
    while cur != -1:
        path.append(cur)
        cur = pre[cur]
    return path


if __name__ == "__main__":
    n = 5
    edges = [[1, 2, 1], [2, 3, 1], [2, 4, 1], [3, 4, 1], [3, 5, 1], [4, 5, 1]]
    adjList = [[] for _ in range(n)]
    for u, v, w in edges:
        adjList[u - 1].append((v - 1, w))
        adjList[v - 1].append((u - 1, w))
    start = 0
    end = 4
    print(lexicographicallySmallestShortestPath(n, adjList, start, end))
