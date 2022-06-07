from typing import DefaultDict, Generic, Hashable, List, Tuple, TypeVar
from dataclasses import dataclass
from collections import defaultdict, deque


V = TypeVar('V', bound=Hashable)


@dataclass(slots=True)
class Edge(Generic[V]):
    fromV: V
    toV: V
    cap: int
    cost: int
    flow: int


# class Edge(Generic[V]):
#     __slots__ = ('fromV', 'toV', 'cap', 'cost', 'flow')

#     def __init__(self, fromV: V, toV: V, cap: int, cost: int, flow: int) -> None:
#         self.fromV = fromV
#         self.toV = toV
#         self.cap = cap
#         self.cost = cost
#         self.flow = flow


class MinCostMaxFlow(Generic[V]):
    """最小费用流的连续最短路算法复杂度为流量*最短路算法复杂度"""

    def __init__(self, start: V, end: V):
        self._start: V = start
        self._end: V = end
        self._edges: List['Edge'[V]] = []
        self._reGraph: DefaultDict[V, List[int]] = defaultdict(list)  # 残量图存储的是边的下标

    def addEdge(self, fromV: V, toV: V, cap: int, cost: int) -> None:
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

        def spfa() -> int:
            """spfa沿着最短路寻找增广路径  有负cost的边不能用dijkstra"""
            nonlocal dist
            dist = defaultdict(lambda: int(1e20), {self._start: 0})
            inQueue = defaultdict(lambda: False)
            queue = deque([self._start])
            inFlow = defaultdict(int, {self._start: int(1e20)})  # 到每条边上的流量
            pre = defaultdict(lambda: -1)

            while queue:
                cur = queue.popleft()
                inQueue[cur] = False
                for edgeIndex in self._reGraph[cur]:
                    edge = self._edges[edgeIndex]
                    cost, flow, cap, next = edge.cost, edge.flow, edge.cap, edge.toV
                    if dist[cur] + cost < dist[next] and (cap - flow) > 0:
                        dist[next] = dist[cur] + cost
                        pre[next] = edgeIndex
                        inFlow[next] = min(inFlow[cur], cap - flow)
                        if not inQueue[next]:
                            inQueue[next] = True
                            queue.append(next)

            resDelta = inFlow[self._end]
            if resDelta > 0:  # 找到可行流
                cur = self._end
                while cur != self._start:
                    preEdgeIndex = pre[cur]
                    self._edges[preEdgeIndex].flow += resDelta
                    self._edges[preEdgeIndex ^ 1].flow -= resDelta
                    cur = self._edges[preEdgeIndex].fromV
            return resDelta

        dist = defaultdict(lambda: int(1e20), {self._start: 0})
        flow = cost = 0
        while True:
            delta = spfa()
            if delta == 0:
                break
            flow += delta
            cost += delta * dist[self._end]
        return flow, cost
