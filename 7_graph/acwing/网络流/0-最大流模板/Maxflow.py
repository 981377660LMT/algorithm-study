from collections import defaultdict
from typing import DefaultDict, Generic, Hashable, Protocol, TypeVar

Vertex = TypeVar('Vertex', bound=Hashable)
Graph = DefaultDict[Vertex, DefaultDict[Vertex, int]]  # 有向带权图


class MaxFlowStrategy(Protocol):
    """interface"""

    graph: Graph

    def getMaxFlow(self, start: Vertex, end: Vertex) -> int:
        """need to be lazy"""
        ...


class MaxFlow(Generic[Vertex]):
    def __init__(self, graph: Graph, *, strategy: MaxFlowStrategy | None = None) -> None:
        self._graph = graph
        self._strategy = strategy if strategy is not None else self._createDefaultStrategy()

    def getMaxFlow(self, start: Vertex, end: Vertex) -> int:
        return self._strategy.getMaxFlow(start, end)

    def switchTo(self, strategy: MaxFlowStrategy) -> None:
        self._strategy = strategy

    def _createDefaultStrategy(self) -> MaxFlowStrategy:
        """to do"""
        return EK(self._graph)


###################################################################


class EK(MaxFlowStrategy):
    """EK 求最大流"""

    def __init__(self, graph: Graph) -> None:
        self.graph = graph

    def getMaxFlow(self, start: Vertex, end: Vertex) -> int:
        return 1


#################################################################
if __name__ == '__main__':
    adjMap = defaultdict(lambda: defaultdict(int))
    adjMap['a']['b'] = 1
    maxFlow = MaxFlow[int](adjMap)
    print(maxFlow.getMaxFlow(1, 2))

