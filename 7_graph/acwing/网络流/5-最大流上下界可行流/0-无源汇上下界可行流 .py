# 给定一个包含 n 个点 m 条边的有向图，每条边都有一个流量下界和流量上界。
# !可行流：求一种`可行方案`使得在所有点满足流量平衡条件的前提下，所有边满足流量限制
# 1≤n≤200,
# 1≤m≤10200,

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
        """边的流量=容量-残量"""
        assert v1 in self._graph and v2 in self._graph[v1]
        return self._graph[v1][v2] - self._reGraph[v1][v2]

    def getRemainOfEdge(self, v1: int, v2: int) -> int:
        """边的残量(剩余的容量)"""
        assert v1 in self._graph and v2 in self._graph[v1]
        return self._reGraph[v1][v2]

    def _updateRedisualGraph(self) -> None:
        """残量图 存储每条边的剩余流量"""
        self._reGraph = defaultdict(lambda: defaultdict(int))
        for cur in self._graph:
            for next in self._graph[cur]:
                self._reGraph[cur][next] = self._graph[cur][next]
                self._reGraph[next].setdefault(cur, 0)


# 无源汇上下界可行流
n, m = map(int, input().split())
edge = defaultdict(lambda: [0, 0])
diff = defaultdict(int)
adjMap = defaultdict(lambda: defaultdict(int))
for _ in range(m):
    u, v, lower, upper = map(int, input().split())
    edge[(u, v)][0] = lower
    edge[(u, v)][1] = upper
    diff[u] -= lower
    diff[v] += lower
    adjMap[u][v] += upper - lower

fullFlow = 0
START, END, OFFSET = -1, -2, int(1e4)
for key, delta in diff.items():
    if delta > 0:
        fullFlow += delta
        adjMap[START][key] += delta
    elif delta < 0:
        adjMap[key][END] += -delta

# ///////////////////////////////////////////
maxFlow = Dinic(adjMap)
res = maxFlow.calMaxFlow(START, END)
if res != fullFlow:
    print('NO')
    exit(0)
print('YES')
for (u, v) in edge:
    print(maxFlow.getFlowOfEdge(u, v) + edge[(u, v)][0])


# 有问题
