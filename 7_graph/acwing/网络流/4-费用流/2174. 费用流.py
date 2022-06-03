# 给定一个包含 n 个点 m 条边的有向图，并给定每条边的容量和费用，边的容量非负。
# 图中可能存在重边和自环，保证费用不会存在负环。
# !求从 S 到 T 的最大流，以及在流量最大时的最小费用。
# 2≤n≤5000,
# 1≤m≤50000,

import sys
from typing import DefaultDict, Tuple
from collections import defaultdict, deque

Graph = DefaultDict[int, DefaultDict[int, int]]  # 有向带权图,权值为容量


class EK:
    """EK 求最小费用最大流"""

    def __init__(self, graph: Graph, cost: Graph) -> None:
        self._graph = graph  # 容量原图
        self._cost = cost

    def calMinCostMaxFlow(self, start: int, end: int) -> Tuple[int, int]:
        def spfa() -> int:
            """spfa沿着最短路寻找增广路径"""
            nonlocal dist
            dist = defaultdict(lambda: int(1e20), {start: 0})
            inFlow = defaultdict(int, {start: int(1e20)})  # 到每条边上流量
            inQueue = defaultdict(lambda: False)
            queue = deque([start])
            pre = {start: start}

            while queue:
                cur = queue.popleft()
                inQueue[cur] = False
                for next in self._reGraph[cur]:
                    weight = self._reGraph[cur][next][1]
                    if dist[cur] + weight < dist[next] and (self._reGraph[cur][next][0] > 0):
                        dist[next] = dist[cur] + weight
                        pre[next] = cur
                        inFlow[next] = min(inFlow[cur], self._reGraph[cur][next][0])
                        if not inQueue[next]:
                            inQueue[next] = True
                            queue.append(next)

            resDelta = inFlow[end]
            if resDelta > 0:  # 找到可行流
                cur = end
                while cur != start:
                    parent = pre[cur]
                    self._reGraph[parent][cur][0] -= resDelta
                    self._reGraph[cur][parent][0] += resDelta
                    cur = parent
            return resDelta

        self._updateRedisualGraph()
        dist = defaultdict(lambda: int(1e20), {start: 0})  # 最短路
        flow = cost = 0
        while True:
            delta = spfa()
            if delta == 0:
                break
            flow += delta
            cost += delta * dist[end]
        return flow, cost

    def getFlowOfEdge(self, v1: int, v2: int) -> int:
        """获取某条边上的`流量`"""
        return self._graph[v1][v2] - self._reGraph[v1][v2][0]

    def _updateRedisualGraph(self) -> None:
        self._reGraph = defaultdict(lambda: defaultdict(list))
        for cur in self._graph:
            for next in self._graph[cur]:
                self._reGraph[cur][next] = [self._graph[cur][next], self._cost[cur][next]]
                self._reGraph[next][cur] = [0, -self._cost[cur][next]]


# endregion

# !图中不存在重边和自环
input = sys.stdin.readline
n, m, start, end = map(int, input().split())
adjMap = defaultdict(lambda: defaultdict(int))
costMap = defaultdict(lambda: defaultdict(int))

# 从点 u 到点 v 存在一条有向边，容量为 c。
for _ in range(m):
    u, v, c, cost = map(int, input().split())
    adjMap[u][v] = c
    costMap[u][v] = cost

maxFlow = EK(adjMap, costMap)
flow, cost = maxFlow.calMinCostMaxFlow(start, end)
print(flow, cost, sep=' ')
