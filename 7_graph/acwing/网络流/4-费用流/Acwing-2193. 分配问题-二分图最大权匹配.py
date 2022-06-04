# 有 n 件工作要分配给 n 个人做。
# 第 i 个人做第 j 件工作产生的效益为 cij。
# 试设计一个将 n 件工作分配给 n 个人做的分配方案。
# 对于给定的 n 件工作和 n 个人，计算最优分配方案和最差分配方案。
# 1≤n≤50,

# !最优解肯定是最大匹配 即取最大流

from typing import DefaultDict, Tuple, List
from collections import defaultdict, deque

Graph = DefaultDict[int, DefaultDict[int, int]]  # 有向带权图,权值为容量


class EK:
    """EK 求最小费用最大流"""

    def __init__(self, flowGraph: Graph, costGraph: Graph) -> None:
        self._graph = flowGraph  # 容量原图
        self._cost = costGraph

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
        dist = defaultdict(lambda: int(1e20), {start: 0})
        flow = cost = 0
        while True:
            delta = spfa()
            if delta == 0:
                break
            flow += delta
            cost += delta * dist[end]
        return flow, cost

    def getFlowOfEdge(self, v1: int, v2: int) -> int:
        return self._graph[v1][v2] - self._reGraph[v1][v2][0]

    def _updateRedisualGraph(self) -> None:
        self._reGraph = defaultdict(lambda: defaultdict(list))
        for cur in self._graph:
            for next in self._graph[cur]:
                self._reGraph[cur][next] = [self._graph[cur][next], self._cost[cur][next]]
                self._reGraph[next][cur] = [0, -self._cost[cur][next]]


# 最小费用、最大费用
adjMap1 = defaultdict(lambda: defaultdict(int))  # 容量
costMap1 = defaultdict(lambda: defaultdict(int))  # 费用
adjMap2 = defaultdict(lambda: defaultdict(int))  # 容量
costMap2 = defaultdict(lambda: defaultdict(int))  # 费用
n = int(input())
for i in range(n):
    nums = list(map(int, input().split()))
    for j, cost in enumerate(nums):
        adjMap1[i][j + 1000] = 1
        costMap1[i][j + 1000] = cost
        adjMap2[i][j + 1000] = 1
        costMap2[i][j + 1000] = -cost
for i in range(n):
    adjMap1[-1][i] = 1
    costMap1[-1][i] = 0
    adjMap1[i + 1000][int(1e9)] = 1
    costMap1[i + 1000][int(1e9)] = 0
    adjMap2[-1][i] = 1
    costMap2[-1][i] = 0
    adjMap2[i + 1000][int(1e9)] = 1
    costMap2[i + 1000][int(1e9)] = 0
MCMF1 = EK(adjMap1, costMap1)
MCMF2 = EK(adjMap2, costMap2)
_, cost1 = MCMF1.calMinCostMaxFlow(-1, int(1e9))
_, cost2 = MCMF2.calMinCostMaxFlow(-1, int(1e9))
print(cost1, -cost2, sep='\n')
