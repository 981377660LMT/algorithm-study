# https://nyaannyaan.github.io/library/random_graph/gen.hpp
# https://nyaannyaan.github.io/library/random_graph/graph.hpp

# TODO


import random
import time
from typing import Optional


class Edge:
    __slots__ = ("u", "v", "w", "id")

    def __init__(self, u: int, v: int, w=1, id=-1):
        self.u = u
        self.v = v
        self.w = w
        self.id = id

    def __repr__(self) -> str:
        return f"Edge({self.u}, {self.v}, {self.w}, {self.id})"


class Graph:
    __slots__ = ("n", "m", "edges", "weighted")

    def __init__(self, n=0, weighted=False) -> None:
        self.n = n
        self.m = 0
        self.edges = []
        self.weighted = weighted

    def addDirectedEdge(self, u: int, v: int, w=1, id=-1) -> None:
        """添加从u到v的边"""
        self.edges.append(Edge(u=u, v=v, w=w, id=id))
        self.m += 1

    def addUndirectedEdge(self, u: int, v: int, w=1, id=-1) -> None:
        """添加从min(u,v)到max(u,v)的边"""
        min_, max_ = min(u, v), max(u, v)
        self.addDirectedEdge(min_, max_, w=w, id=id)

    def getAdjList(self, directed=False) -> list[list[Edge]]:
        adjList = [[] for _ in range(self.n)]
        for edge in self.edges:
            adjList[edge.u].append(edge)
            if not directed:
                adjList[edge.v].append(edge)
        return adjList

    def getAdjMatrix(self, directed=False) -> list[list[int]]:
        adjMatrix = [[0] * self.n for _ in range(self.n)]
        for edge in self.edges:
            adjMatrix[edge.u][edge.v] = edge.w
            if not directed:
                adjMatrix[edge.v][edge.u] = edge.w
        return adjMatrix

    @property
    def edgesSize(self) -> int:
        return self.m


# TODO
class UndirectedGraphGenerator:
    """随机无向图生成器"""

    __slots__ = ("_wMin", "_wMax")

    def __init__(self, seed: Optional[int] = None) -> None:
        if seed is None:
            seed = int(time.time())
        random.seed(seed)
        self._wMin = 1
        self._wMax = 1

    def test(self, n: int, isTree=True, weighted=False, wMin=1, wMax=1):
        ...

    def _setWeight(self, weighted: bool, wMin: int, wMax: int) -> None:
        self._wMin = wMin
        self._wMax = wMax
        if not weighted:
            self._wMin = 1
            self._wMax = 1

    def _addEdge(self, g: Graph, u: int, v: int) -> None:
        w = self._wMin if self._wMin == self._wMax else random.randint(self._wMin, self._wMax)
        g.addUndirectedEdge(u, v, w)
