# 多源汇添加虚拟点后就变成了单源汇
'''
添加一个虚拟源点SS和虚拟汇点TT，
SS 和 所有原图源点间连接容量无穷大的边
原图汇点和TT 连接容量无穷大的边
求SS到TT的最大流
'''
# 2≤n≤10000,
# 1≤m≤105,
from collections import defaultdict
from collections import defaultdict
from typing import DefaultDict, List, Set
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

                child = next(curArc[cur], None)
                if child is not None:
                    if (depth[child] == depth[cur] + 1) and (self._reGraph[cur][child] > 0):
                        nextFlow = dfsWithCurArc(
                            child, min(minFlow - flow, self._reGraph[cur][child])
                        )
                        if nextFlow == 0:
                            depth[child] = -1
                        self._reGraph[cur][child] -= nextFlow
                        self._reGraph[child][cur] += nextFlow
                        flow += nextFlow
                else:
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

    def getFlowOfEdge(self, v1: int, v2: int) -> int:
        assert v1 in self._graph and v2 in self._graph[v1]
        return self._graph[v1][v2] - self._reGraph[v1][v2]

    def getRemainOfEdge(self, v1: int, v2: int) -> int:
        assert v1 in self._graph and v2 in self._graph[v1]
        return self._reGraph[v1][v2]

    def _updateRedisualGraph(self) -> None:
        self._reGraph = defaultdict(lambda: defaultdict(int))
        for cur in self._graph:
            for next in self._graph[cur]:
                self._reGraph[cur][next] = self._graph[cur][next]
                self._reGraph[next].setdefault(cur, 0)


n, m, startCount, endCount = map(int, input().split())
starts = list(map(int, input().split()))
ends = list(map(int, input().split()))
adjMap = defaultdict(lambda: defaultdict(int))
for _ in range(m):
    u, v, w = map(int, input().split())
    adjMap[u][v] += w

START, END, OFFSET = -1, -2, int(1e6)
for num in starts:
    adjMap[START][num] += int(1e9)
for num in ends:
    adjMap[num][END] += int(1e9)

maxFlow = Dinic(adjMap)
print(maxFlow.calMaxFlow(START, END))


