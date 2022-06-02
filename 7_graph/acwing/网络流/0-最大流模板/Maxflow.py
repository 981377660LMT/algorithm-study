from collections import defaultdict, deque
from functools import lru_cache
from typing import DefaultDict, Generic, Hashable, Protocol, TypeVar

Vertex = TypeVar('Vertex', bound=Hashable)
Graph = DefaultDict[Vertex, DefaultDict[Vertex, int]]  # 有向带权图,权值为容量


class MaxFlowStrategy(Protocol):
    """interface of MaxFlow Strategy"""

    def calMaxFlow(self, start: Vertex, end: Vertex) -> int:
        """求出从start到end的最大流"""
        ...

    def getFlowOfEdge(self, v1: Vertex, v2: Vertex) -> int:
        """获取某条边上的`流量`"""
        ...


class MaxFlow(Generic[Vertex]):
    def __init__(self, graph: Graph, *, strategy: MaxFlowStrategy | None = None) -> None:
        self._graph = graph
        self._strategy = strategy if strategy is not None else self._useDefaultStrategy()

    def calMaxFlow(self, start: Vertex, end: Vertex) -> int:
        return self._strategy.calMaxFlow(start, end)

    def switchTo(self, strategy: MaxFlowStrategy) -> None:
        self._strategy = strategy

    def _useDefaultStrategy(self) -> MaxFlowStrategy:
        """to do"""
        return EK(self._graph)


###################################################################
class EK(MaxFlowStrategy):
    """EK 求最大流
    
    如果一个流的残量网络里面没有可行流，那么这个流就是最大流

    时间复杂度:O(V*E^2)
    1. 找增广路
    2. 更新残量网络增广路上的流量
    3. 重复执行1、2直到网络里没有增广路
    """

    def __init__(self, graph: Graph) -> None:
        self._graph = graph  # 容量原图

    @lru_cache(None)
    def calMaxFlow(self, start: Vertex, end: Vertex) -> int:
        def bfs() -> int:
            """bfs在残量网络上寻找增广路径"""
            visited = set([start])
            queue = deque([(start, int(1e20))])
            pre = {start: start}
            resDelta = 0
            while queue:
                cur, delta = queue.popleft()
                if cur == end:
                    resDelta = delta
                    break
                for next in self._reGraph[cur]:
                    if (next not in visited) and (self._reGraph[cur][next] > 0):
                        visited.add(next)
                        queue.append((next, min(delta, self._reGraph[cur][next])))
                        pre[next] = cur

            if resDelta > 0:  # 找到可行流
                cur = end
                while cur != start:
                    parent = pre[cur]
                    # 正向边，消耗了流量，残量减少；反向边，抵消了流量，残量增多
                    # if cur in self._graph:
                    self._reGraph[parent][cur] -= resDelta
                    self._reGraph[cur][parent] += resDelta
                    # else:
                    # self._reGraph[parent][cur] += resDelta
                    # self._reGraph[cur][parent] -= resDelta
                    cur = parent
            return resDelta

        self._updateRedisualGraph()
        res = 0
        while True:
            delta = bfs()
            if delta == 0:
                break
            res += delta
        return res

    def getFlowOfEdge(self, v1: int, v2: int) -> int:
        """获取某条边上的`流量`"""
        return self._reGraph[v1][v2]

    def getCapacityOfEdge(self, v1: int, v2: int) -> int:
        """获取某条边上的`容量`"""
        return self._graph[v1][v2]

    def _updateRedisualGraph(self) -> None:
        self._reGraph = defaultdict(lambda: defaultdict(int))
        for cur in self._graph:
            for next in self._graph[cur]:
                self._reGraph[cur][next] = self._graph[cur][next]
                self._reGraph[next].setdefault(cur, 0)  # 注意这里 因为可能存在重边


class Dinic(MaxFlowStrategy):
    """Dinic 求最大流
    
    如果一个流的残量网络里面没有可行流，那么这个流就是最大流
    
    时间复杂度:O(V^2*E)
    1. 找增广路
    2. 更新残量网络增广路上的流量
    3. 重复执行1、2直到网络里没有增广路
    """

    def __init__(self, graph: Graph) -> None:
        self._graph = graph

    def calMaxFlow(self, start: Vertex, end: Vertex) -> int:
        raise NotImplementedError('not implemented')

    def getFlowOfEdge(self, v1: int, v2: int) -> int:
        raise NotImplementedError('not implemented')


if __name__ == '__main__':
    # 图中可能存在重边和自环
    n, m, start, end = map(int, input().split())
    adjMap = defaultdict(lambda: defaultdict(int))

    # 从点 u 到点 v 存在一条有向边，容量为 c。
    for _ in range(m):
        u, v, c = map(int, input().split())
        adjMap[u][v] += c

    maxFlow = MaxFlow(adjMap)
    print(maxFlow.calMaxFlow(start, end))

