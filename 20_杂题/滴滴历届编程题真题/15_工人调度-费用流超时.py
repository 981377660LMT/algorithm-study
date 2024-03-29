# 某工厂有 N 名工人和 M 台机器，每名工人都有一个能力值，
# 且只懂得操作其中某两台机器。
# 另外，每名工人最多只允许操作一台机器，
# 且每台机器最多只允许被一名工人操作。
# 那么如何分配每名工人操作哪台机器(或者不操作机器)，
# 才能让所有操作机器的工人的能力值总和最大?
# !二分图最大权匹配 => 最小费用最大流 会超时
# n,m<=1e5

from typing import List, Tuple
from collections import deque

INF = int(1e18)


class Edge:
    __slots__ = ("fromV", "toV", "cap", "cost", "flow")

    def __init__(self, fromV: int, toV: int, cap: int, cost: int, flow: int) -> None:
        self.fromV = fromV
        self.toV = toV
        self.cap = cap
        self.cost = cost
        self.flow = flow


class MinCostMaxFlow:
    """最小费用流的连续最短路算法复杂度为流量*最短路算法复杂度"""

    __slots__ = ("_n", "_start", "_end", "_edges", "_reGraph", "_dist", "_visited", "_curEdges")

    def __init__(self, n: int, start: int, end: int):
        """
        Args:
            n (int): 包含虚拟点在内的总点数
            start (int): (虚拟)源点
            end (int): (虚拟)汇点
        """
        assert 0 <= start < n and 0 <= end < n
        self._n = n
        self._start = start
        self._end = end
        self._edges: List["Edge"] = []
        self._reGraph: List[List[int]] = [[] for _ in range(n + 10)]  # 残量图存储的是边的下标

        self._dist = [INF] * (n + 10)
        self._visited = [False] * (n + 10)
        self._curEdges = [0] * (n + 10)

    def addEdge(self, fromV: int, toV: int, cap: int, cost: int) -> None:
        """原边索引为i 反向边索引为i^1"""
        self._edges.append(Edge(fromV, toV, cap, cost, 0))
        self._edges.append(Edge(toV, fromV, 0, -cost, 0))
        len_ = len(self._edges)
        self._reGraph[fromV].append(len_ - 2)
        self._reGraph[toV].append(len_ - 1)

    def work(self) -> Tuple[int, int]:
        """
        Returns:
            Tuple[int, int]: [最大流,最小费用]
        """
        maxFlow, minCost = 0, 0
        while self._spfa():
            # !如果流量限定为1，那么一次dfs只会找到一条费用最小的增广流
            # !如果流量限定为INF，那么一次dfs不只会找到一条费用最小的增广流
            flow = self._dfs(self._start, self._end, INF)
            maxFlow += flow
            minCost += flow * self._dist[self._end]
        return maxFlow, minCost

    def _spfa(self) -> bool:
        """spfa沿着最短路寻找增广路径  有负cost的边不能用dijkstra"""
        n, start, end, edges, reGraph, visited = (
            self._n,
            self._start,
            self._end,
            self._edges,
            self._reGraph,
            self._visited,
        )

        self._curEdges = [0] * n
        self._dist = dist = [INF] * n
        dist[start] = 0
        queue = deque([start])

        while queue:
            cur = queue.popleft()
            visited[cur] = False
            for edgeIndex in reGraph[cur]:
                edge = edges[edgeIndex]
                cost, remain, next = edge.cost, edge.cap - edge.flow, edge.toV
                if remain > 0 and dist[cur] + cost < dist[next]:
                    dist[next] = dist[cur] + cost
                    if not visited[next]:
                        visited[next] = True
                        if queue and dist[queue[0]] > dist[next]:
                            queue.appendleft(next)
                        else:
                            queue.append(next)

        return dist[end] != INF

    def _dfs(self, cur: int, end: int, flow: int) -> int:
        if cur == end:
            return flow

        visited, reGraph, curEdges, edges, dist = (
            self._visited,
            self._reGraph,
            self._curEdges,
            self._edges,
            self._dist,
        )

        visited[cur] = True
        res = flow
        index = curEdges[cur]
        while res and index < len(reGraph[cur]):
            edgeIndex = reGraph[cur][index]
            next, remain = edges[edgeIndex].toV, edges[edgeIndex].cap - edges[edgeIndex].flow
            if remain > 0 and not visited[next] and dist[next] == dist[cur] + edges[edgeIndex].cost:
                delta = self._dfs(next, end, remain if remain < res else res)
                res -= delta
                edges[edgeIndex].flow += delta
                edges[edgeIndex ^ 1].flow -= delta
            curEdges[cur] += 1
            index = curEdges[cur]

        visited[cur] = False
        return flow - res


n, m = map(int, input().split())
START, END, OFFSET = n + m + 2, n + m + 3, n
mcmf = MinCostMaxFlow(n + m + 5, START, END)
for i in range(n):
    mcmf.addEdge(START, i, 1, 0)
for i in range(m):
    mcmf.addEdge(i + OFFSET, END, 1, 0)
for i in range(n):
    # !第 i 名工人只懂得操作第 a 台和第 b 台机器，且其能力值为 c 。
    a, b, c = map(int, input().split())
    a, b = a - 1, b - 1


print(-mcmf.work()[1])
