# 有向带权图
# 时间复杂度：V*E^2
# bfs寻找增广路径

from typing import DefaultDict, List
from collections import defaultdict, deque


WeightedDirectedGraph = DefaultDict[int, DefaultDict[int, int]]


class MaxFlow:
    """edmond-karp算法求解有向带权图的最大流，时间复杂度"""

    def __init__(self, graph: WeightedDirectedGraph, /, *, start: int, target: int):
        assert len(graph) >= 2 and start != target
        self._graph = graph
        self._start = start
        self._target = target
        # 1.构建残量图
        self._rGraph = self._buildRedisualGraph(graph)

    def getResult(self) -> int:
        # 2.bfs在残量图中寻找增广路径
        res = 0

        while True:
            augPath = self._findAugmentingPath(self._rGraph, self._start, self._target)
            if not augPath:
                break

            # 3.计算增广路径上的最小值
            minWeight = int(1e20)
            for i in range(0, len(augPath) - 1):
                cur, next = augPath[i], augPath[i + 1]
                weight = self._rGraph[cur][next]
                minWeight = min(minWeight, weight)
            res += minWeight

            # 4.根据增广路径更新残量图
            for i in range(0, len(augPath) - 1):
                cur, next = augPath[i], augPath[i + 1]
                # 正向边，消耗了流量，残量减少；反向边，抵消了流量，残量增多
                if next in self._graph[cur]:
                    self._rGraph[cur][next] -= minWeight
                    self._rGraph[next][cur] += minWeight
                else:
                    self._rGraph[cur][next] += minWeight
                    self._rGraph[next][cur] -= minWeight

        return res

    def getFlowOfEdge(self, cur: int, next: int) -> int:
        """获取某条边上的`流量`"""
        return self._rGraph[cur][next]

    def getCapacityOfEdge(self, cur: int, next: int) -> int:
        """获取某条边上的`容量`"""
        return self._graph[cur][next]

    def _buildRedisualGraph(self, graph: WeightedDirectedGraph) -> WeightedDirectedGraph:
        """构建残量图"""
        rGraph = defaultdict(lambda: defaultdict(int))
        for cur, mapping in graph.items():
            for next, weight in mapping.items():
                rGraph[cur][next] = weight
                rGraph[next][cur] = 0
        return rGraph

    def _findAugmentingPath(
        self, graph: WeightedDirectedGraph, start: int, target: int
    ) -> List[int]:
        """bfs在残量图中寻找增广路径"""
        queue = deque()
        pre = defaultdict(lambda: -1)
        queue.append(start)
        pre[start] = start

        while queue:
            cur = queue.popleft()
            if cur == target:
                break
            for next in graph[cur].keys():
                if pre[next] == -1 and graph[cur][next] > 0:
                    pre[next] = cur
                    queue.append(next)

        res = []
        if pre[target] == -1:
            return res

        cur = target
        while cur != start:
            res.append(cur)
            cur = pre[cur]
        res.append(start)
        return res[::-1]


if __name__ == '__main__':
    adjMap1 = defaultdict(lambda: defaultdict(int))
    adjMap1[0][1] = 3
    adjMap1[0][2] = 2
    adjMap1[1][2] = 5
    adjMap1[1][3] = 2
    adjMap1[2][3] = 3

    maxFlow1 = MaxFlow(adjMap1, start=0, target=3)
    assert maxFlow1.getResult() == 5

    adjMap2 = defaultdict(lambda: defaultdict(int))
    adjMap2[0][1] = 9
    adjMap2[0][3] = 9
    adjMap2[1][2] = 8
    adjMap2[1][3] = 10
    adjMap2[2][5] = 10
    adjMap2[3][2] = 1
    adjMap2[3][4] = 3
    adjMap2[4][2] = 8
    adjMap2[4][5] = 7

    maxFlow2 = MaxFlow(adjMap2, start=0, target=5)
    assert maxFlow2.getResult() == 12
