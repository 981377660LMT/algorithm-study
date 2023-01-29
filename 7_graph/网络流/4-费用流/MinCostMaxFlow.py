# https://blog.csdn.net/Ratina/article/details/94748922
# Dinic算法（弧优化），最小费用最大流
# !golang实现
# !https://github.dev/EndlessCheng/codeforces-go/blob/6d127a66c2a11651e8d11783af687264780e82a8/copypasta/graph.go#L3584

from heapq import heappop, heappush
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


class MinCostMaxFlowDinic:
    """最小费用流的连续最短路算法复杂度为流量*最短路算法复杂度"""

    __slots__ = (
        "_n",
        "_start",
        "_end",
        "_edges",
        "_reGraph",
        "_dist",
        "_visited",
        "_curEdges",
        "_hasMinusCost",
    )

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
        self._reGraph: List[List[int]] = [[] for _ in range(n)]  # 残量图存储的是边的下标
        self._hasMinusCost = False

        self._dist = [INF] * n
        self._visited = [False] * n
        self._curEdges = [0] * n

    def addEdge(self, fromV: int, toV: int, cap: int, cost: int) -> None:
        """原边索引为i 反向边索引为i^1"""
        self._edges.append(Edge(fromV, toV, cap, cost, 0))
        self._edges.append(Edge(toV, fromV, 0, -cost, 0))
        len_ = len(self._edges)
        self._reGraph[fromV].append(len_ - 2)
        self._reGraph[toV].append(len_ - 1)
        if cost < 0:
            self._hasMinusCost = True

    def work(self) -> Tuple[int, int]:
        """
        Returns:
            Tuple[int, int]: [最大流,最小费用]
        """
        maxFlow, minCost = 0, 0
        fun = self._spfaStrategy if self._hasMinusCost else self._dijkstraStrategy
        while fun():
            self._curEdges = [0] * self._n
            # !如果流量限定为1，那么一次dfs只会找到一条费用最小的增广流
            # !如果流量限定为INF，那么一次dfs不只会找到一条费用最小的增广流
            flow = self._dfs(self._start, self._end, INF)
            maxFlow += flow
            minCost += flow * self._dist[self._end]
        return maxFlow, minCost

    def slope(self) -> List[Tuple[int, int]]:
        """
        Returns:
            List[Tuple[int, int]]: 流量为a时,最小费用是b
        """
        res = [(0, 0)]
        flow, cost = 0, 0
        fun = self._spfaStrategy if self._hasMinusCost else self._dijkstraStrategy
        while fun():
            self._curEdges = [0] * self._n
            deltaFlow = self._dfs(self._start, self._end, INF)
            flow += deltaFlow
            cost += deltaFlow * self._dist[self._end]
            res.append((flow, cost))  # type: ignore
        return res

    def _spfaStrategy(self) -> bool:
        """spfa沿着最短路寻找增广路径  有负cost的边不能用dijkstra"""
        n, start, end, edges, reGraph, visited = (
            self._n,
            self._start,
            self._end,
            self._edges,
            self._reGraph,
            self._visited,
        )

        self._dist = dist = [INF] * n
        dist[start] = 0
        visited = [False] * n
        visited[start] = True
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

    def _dijkstraStrategy(self) -> bool:
        """dijkstra沿着最短路寻找增广路径  有负cost的边不能用dijkstra"""
        n, start, end, edges, reGraph = (
            self._n,
            self._start,
            self._end,
            self._edges,
            self._reGraph,
        )

        self._dist = dist = [INF] * n
        dist[start] = 0
        pq = [(0, start)]

        while pq:
            curDist, cur = heappop(pq)
            if cur == end:
                return True
            if dist[cur] < curDist:
                continue
            for edgeIndex in reGraph[cur]:
                edge = edges[edgeIndex]
                cost, remain, next = edge.cost, edge.cap - edge.flow, edge.toV
                if remain > 0 and curDist + cost < dist[next]:
                    dist[next] = curDist + cost
                    heappush(pq, (dist[next], next))

        return False

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
        n = len(reGraph[cur])
        while res and index < n:
            edgeIndex = reGraph[cur][index]
            edge = edges[edgeIndex]
            next, remain = edge.toV, edge.cap - edge.flow
            if remain > 0 and not visited[next] and dist[next] == dist[cur] + edge.cost:
                delta = self._dfs(next, end, remain if remain < res else res)
                res -= delta
                edge.flow += delta
                edges[edgeIndex ^ 1].flow -= delta
            curEdges[cur] += 1
            index = curEdges[cur]

        visited[cur] = False
        return flow - res


class MinCostMaxFlowEK:
    """最小费用流的复杂度为流量*spfa的复杂度"""

    __slots__ = ("_n", "_start", "_end", "_edges", "_reGraph", "_dist", "_pre", "_flow")

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
        self._reGraph: List[List[int]] = [[] for _ in range(n)]  # 残量图存储的是边的下标

        self._dist = [INF] * n
        self._flow = [0] * n
        self._pre = [-1] * n

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
            delta = self._flow[self._end]
            minCost += delta * self._dist[self._end]
            maxFlow += delta
            cur = self._end
            while cur != self._start:
                edgeIndex = self._pre[cur]
                self._edges[edgeIndex].flow += delta
                self._edges[edgeIndex ^ 1].flow -= delta
                cur = self._edges[edgeIndex].fromV
        return maxFlow, minCost

    def _spfa(self) -> bool:
        """spfa沿着最短路寻找增广路径  有负cost的边不能用dijkstra"""
        n, start, end, edges, reGraph = (
            self._n,
            self._start,
            self._end,
            self._edges,
            self._reGraph,
        )

        self._flow = flow = [0] * n
        self._pre = pre = [-1] * n
        self._dist = dist = [INF] * n
        dist[start] = 0
        flow[start] = INF
        inQueue = [False] * n
        inQueue[start] = True
        queue = deque([start])

        while queue:
            cur = queue.popleft()
            inQueue[cur] = False
            for edgeIndex in reGraph[cur]:
                edge = edges[edgeIndex]
                cost, remain, next = edge.cost, edge.cap - edge.flow, edge.toV
                if remain > 0 and dist[cur] + cost < dist[next]:
                    dist[next] = dist[cur] + cost
                    pre[next] = edgeIndex
                    flow[next] = remain if remain < flow[cur] else flow[cur]
                    if not inQueue[next]:
                        inQueue[next] = True
                        if queue and dist[queue[0]] > dist[next]:
                            queue.appendleft(next)
                        else:
                            queue.append(next)

        return pre[end] != -1
