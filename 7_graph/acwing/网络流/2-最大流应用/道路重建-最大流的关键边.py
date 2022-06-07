# 1≤N≤500,
# 1≤M≤5000,

# n个点m条边的网络，给定源点S和汇点T，求如果有这样边：
# !只给其扩大容量之后整个流网络的最大流能够变大，
# 对于这样的边我们称之为`关键边`。求这样的边的个数
# https://www.acwing.com/solution/content/73355/


# 1.先对原图跑一次最大流，求出原图的残量网络
# 2.对于边(u,v)，判断从s->u和从v->t是否都存在流量大于0的路
# 3.如果某条边上的流量是满的 并且左右端点可以分别到达源点、汇点 就是关键边

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


# 城市编号从 0 到 N−1。
# 生产日常商品的城市为 0 号城市，首都为 N−1 号城市。
def dfs(cur: int, visited: Set[int], graph: Graph, isRevered: bool) -> None:
    visited.add(cur)
    for next in graph[cur]:
        remain = (
            maxFlow.getRemainOfEdge(cur, next)
            if not isRevered
            else maxFlow.getRemainOfEdge(next, cur)
        )
        if next not in visited and remain != 0:  # 可以继续走
            dfs(next, visited, graph, isRevered)


n, m = map(int, input().split())
adjMap = defaultdict(lambda: defaultdict(int))
rAdjMap = defaultdict(lambda: defaultdict(int))
for _ in range(m):
    u, v, w = map(int, input().split())
    adjMap[u][v] += w
    rAdjMap[v][u] += w

maxFlow = Dinic(adjMap)
maxFlow.calMaxFlow(0, n - 1)

# !如果某条边上的流量是满的 就可以作为候选
visited1, visited2 = set(), set()  # 可以到源点的点 和 可以到汇点的点
dfs(0, visited1, adjMap, False)
dfs(n - 1, visited2, rAdjMap, True)

res = 0
for cur in adjMap:
    for next in adjMap[cur]:
        if (
            (cur in visited1)
            and (next in visited2)
            and maxFlow.getRemainOfEdge(cur, next) == 0  # 满流
        ):
            res += 1
print(res)
