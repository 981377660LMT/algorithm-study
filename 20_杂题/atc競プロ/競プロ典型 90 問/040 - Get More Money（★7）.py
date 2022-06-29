# n个房间 进入每个房间必须有钥匙且交钱w 房间内有解锁其他房间的钥匙和钱
# 房间i里的钥匙编号大于等于i+1
# 问如何进入若干个房间 使得最后收益最大
# n<=100
# https://twitter.com/e869120/status/1393336369540341760/photo/1
# !燃やす埋める問題に帰着させる 最小割为答案(房子分成两类 进还是不进 再计算最小罚款 让每个房子都有`进/不进`的归属)
# !最大流最小割 时间复杂度O(V^2E)


import sys
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


# 最大流代表消耗的金钱
input = sys.stdin.readline
n, w = map(int, input().split())
moneys = list(map(int, input().split()))
keys = [[] for _ in range(n)]
for i in range(n):
    count, *rest = list(map(int, input().split()))
    keys[i].extend([num - 1 for num in rest])

adjMap = defaultdict(lambda: defaultdict(int))
START, END = -1, -2

# !在sum(moneys)的基础上计算罚款
# 源点:不访问
# 汇点:访问
for i in range(n):
    adjMap[START][i] = w  # 访问　需要罚款w
    adjMap[i][END] = moneys[i]  # 不访问　需要罚款moneys[i]

for i in range(n):
    for key in keys[i]:
        adjMap[i][key] = int(1e20)  # 去掉这条边表示 访问需要key的房子却没有访问key所在的房子 罚款1e20

maxFlow = Dinic(adjMap)
res = maxFlow.calMaxFlow(START, END)
print(sum(moneys) - res)  # 总收益-最小割(分成两半的最少罚款)
