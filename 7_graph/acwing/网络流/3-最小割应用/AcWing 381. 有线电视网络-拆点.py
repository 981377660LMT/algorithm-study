# 给定一张 n 个点 m 条边的无向图，求最少去掉多少个点，可以使图不连通。
# 如果不管去掉多少个点，都无法使原图不连通，则直接返回 n。
# !0≤n≤50


from typing import DefaultDict, Set
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

    def getFlowOfEdge(self, v1: int, v2: int) -> int:
        return self._graph[v1][v2] - self._reGraph[v1][v2]

    def _updateRedisualGraph(self) -> None:
        self._reGraph = defaultdict(lambda: defaultdict(int))
        for cur in self._graph:
            for next in self._graph[cur]:
                self._reGraph[cur][next] = self._graph[cur][next]
                self._reGraph[next].setdefault(cur, 0)


def main(n: int, rawAdjMap: DefaultDict[int, Set[int]]) -> int:
    """最少去掉多少个点，可以使图不连通

    不连通：最小割模型
    
    枚举源点与汇点

    0-n 入点
    OFFSET-OFFSET+n 出点
    rawAdjMap 无向图
    """
    res = n
    OFFSET = int(1e4)
    adjMap = defaultdict(lambda: defaultdict(int))
    for cur, nexts in rawAdjMap.items():
        for next in nexts:
            adjMap[cur + OFFSET][next] = int(1e20)  # a出点 => b入点容量无穷 因为要割的是点，不是边
    for i in range(n):
        adjMap[i][OFFSET + i] = 1  # a入点=>a出点容量为1 因为要割的是点

    for start in range(n):
        for end in range(n):
            if start == end:
                continue
            # 源点是start的出点 ，汇点是end的入点
            res = min(res, Dinic(adjMap).calMaxFlow(start + OFFSET, end))
    return res


while True:
    try:
        n, m, *rest = list(input().split())
        adjMap = defaultdict(set)  # 无向图
        for item in rest:
            u, v = map(int, item[1:-1].split(','))
            adjMap[u].add(v)
            adjMap[v].add(u)
        print(main(int(n), adjMap))
    except EOFError:
        break
