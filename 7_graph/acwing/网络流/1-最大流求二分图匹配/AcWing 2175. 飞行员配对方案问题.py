from collections import defaultdict
from typing import DefaultDict
from collections import defaultdict, deque

Graph = DefaultDict[int, DefaultDict[int, int]]  # 有向带权图,权值为容量


class Dinic:
    def __init__(self, graph: Graph) -> None:
        self._graph = graph

    def calMaxFlow(self, start: int, end: int) -> int:
        def bfs() -> None:
            nonlocal depth, curArc
            depth = defaultdict(lambda: -1, {start: 0})
            visted = set([start])
            queue = deque([start])
            curArc = {cur: iter(self._reGraph[cur].keys()) for cur in self._reGraph.keys()}
            while queue:
                cur = queue.popleft()
                for child in self._reGraph[cur]:
                    if (child not in visted) and (self._reGraph[cur][child] > 0):
                        visted.add(child)
                        depth[child] = depth[cur] + 1
                        queue.append(child)

        def dfsWithCurArc(cur: int, minFlow: int) -> int:
            if cur == end:
                return minFlow
            flow = 0
            while True:
                if flow >= minFlow:
                    break
                try:
                    child = next(curArc[cur])
                    if (depth[child] == depth[cur] + 1) and (self._reGraph[cur][child] > 0):
                        nextFlow = dfsWithCurArc(
                            child, min(minFlow - flow, self._reGraph[cur][child])
                        )
                        if nextFlow == 0:
                            depth[child] = -1
                        self._reGraph[cur][child] -= nextFlow
                        self._reGraph[child][cur] += nextFlow
                        flow += nextFlow
                except StopIteration:
                    break
            return flow

        self._updateRedisualGraph()

        res = 0
        depth = defaultdict(lambda: -1, {start: 0})
        curArc = dict()

        while True:
            bfs()
            if depth[end] != -1:
                while True:
                    delta = dfsWithCurArc(start, int(1e20))
                    if delta == 0:
                        break
                    res += delta
            else:
                break
        return res

    def _updateRedisualGraph(self) -> None:
        self._reGraph = defaultdict(lambda: defaultdict(int))
        for cur in self._graph:
            for next in self._graph[cur]:
                self._reGraph[cur][next] = self._graph[cur][next]
                self._reGraph[next].setdefault(cur, 0)

    def getFlowOfEdge(self, v1: int, v2: int) -> int:
        return self._graph[v1][v2] - self._reGraph[v1][v2]


m, n = map(int, input().split())
adjMap = defaultdict(lambda: defaultdict(int))
A, B = set(), set()
while True:
    u, v = map(int, input().split())
    if u == -1 and v == -1:
        break
    adjMap[u][v] = 1
    adjMap[0][u] = 1
    adjMap[v][n + 10] = 1
    A.add(u)


maxFlow = Dinic(adjMap)
print(maxFlow.calMaxFlow(0, n + 10))
res = []
for v1 in A:
    for v2 in adjMap[v1]:
        if maxFlow.getFlowOfEdge(v1, v2) > 0:
            print(v1, v2)

