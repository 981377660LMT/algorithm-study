# https://www.luogu.com.cn/blog/LawrenceSivan/cf545e-paths-and-trees-zui-duan-lu-jing-shu-post
# 求出最短路径树
# 给定一张带正权的无向图和一个源点，求边权和最小的最短路径树。
# 输出该树权值和和和每条边的编号


from heapq import heappop, heappush
from typing import List, Sequence, Tuple

INF = int(1e18)


# !dijkstra求出路径上每个点的前驱点和前驱边
# 其中邻接表的每个元素是一个三元组，分别是邻接点，边权，边的编号
def dijkstraShortestPathTree(
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

    def CF545E(n: int, edges: List[Tuple[int, int, int]], start: int) -> Tuple[int, List[int]]:
        """求出最短路径树的最小边权之和和每条边的编号"""
        adjList = [[] for _ in range(n)]
        for i, (u, v, w) in enumerate(edges):
            adjList[u].append((v, w, i))
            adjList[v].append((u, w, i))
        *_, preE = dijkstraShortestPathTree(n, adjList, start)

        res, eids = 0, []
        for i in range(n):
            if i == start:
                continue
            pre = preE[i]
            res += edges[pre][2]
            eids.append(pre)
        return res, eids

    n, m = map(int, input().split())
    edges = []
    for i in range(m):
        u, v, w = map(int, input().split())
        u, v = u - 1, v - 1
        edges.append((u, v, w))
    start = int(input())
    start -= 1
    res, eids = CF545E(n, edges, start)
    print(res)
    print(*[i + 1 for i in eids])
