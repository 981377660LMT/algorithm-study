from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

from typing import List, Sequence, Tuple
from heapq import heappop, heappush


def dijkstra(n: int, adjList: Sequence[Sequence[Tuple[int, int]]], start: int, end: int) -> int:
    dist = [INF] * n
    dist[start] = 0
    pq = [(0, start)]

    while pq:
        curDist, cur = heappop(pq)
        if dist[cur] < curDist:
            continue
        if cur == end:
            return dist[end]
        for next, weight in adjList[cur]:
            cand = dist[cur] + weight
            if cand < dist[next]:
                dist[next] = cand
                heappush(pq, (dist[next], next))
    return -1


# 给你一个有 n 个节点的 有向带权 图，节点编号为 0 到 n - 1 。图中的初始边用数组 edges 表示，其中 edges[i] = [fromi, toi, edgeCosti] 表示从 fromi 到 toi 有一条代价为 edgeCosti 的边。

# 请你实现一个 Graph 类：

# Graph(int n, int[][] edges) 初始化图有 n 个节点，并输入初始边。
# addEdge(int[] edge) 向边集中添加一条边，其中 edge = [from, to, edgeCost] 。数据保证添加这条边之前对应的两个节点之间没有有向边。
# int shortestPath(int node1, int node2) 返回从节点 node1 到 node2 的路径 最小 代价。如果路径不存在，返回 -1 。一条路径的代价是路径中所有边代价之和。


class Graph:
    def __init__(self, n: int, edges: List[List[int]]):
        self.g = [[] for _ in range(n)]
        for u, v, w in edges:
            self.g[u].append((v, w))

    def addEdge(self, edge: List[int]) -> None:
        from_, to_, cost = edge
        self.g[from_].append((to_, cost))

    def shortestPath(self, node1: int, node2: int) -> int:
        return dijkstra(len(self.g), self.g, node1, node2)


# Your Graph object will be instantiated and called as such:
# obj = Graph(n, edges)
# obj.addEdge(edge)
# param_2 = obj.shortestPath(node1,node2)
