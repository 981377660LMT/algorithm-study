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

    def getFlowOfEdge(self, v1: int, v2: int) -> int:
        return self._graph[v1][v2] - self._reGraph[v1][v2]

    def _updateRedisualGraph(self) -> None:
        self._reGraph = defaultdict(lambda: defaultdict(int))
        for cur in self._graph:
            for next in self._graph[cur]:
                self._reGraph[cur][next] = self._graph[cur][next]
                self._reGraph[next].setdefault(cur, 0)


# 假设有来自 m 个不同单位的代表参加一次国际会议。
# 每个单位的代表数分别为 ri(i=1,2,…,m)。
# 会议餐厅共有 n 张餐桌，每张餐桌可容纳 ci(i=1,2,…,n) 个代表就餐。
# !为了使代表们充分交流，希望从同一个单位来的代表不在同一个餐桌就餐。
# 试设计一个算法，给出满足要求的代表就餐方案。
# 如果问题有解，在第 1 行输出 1，否则输出 0。
# 接下来的 m 行给出每个单位代表的就餐桌号。
# 如果有多个满足要求的方案，只要求输出 1 个方案。


m, n = map(int, input().split())  # m: 代表数量, n: 餐桌数量
adjMap = defaultdict(lambda: defaultdict(int))
nums = list(map(int, input().split()))
caps = list(map(int, input().split()))

START, END = -1, int(1e9)
for i, num in enumerate(nums):
    adjMap[-1][i] = num
for i, cap in enumerate(caps):
    adjMap[i + 1000][END] = cap
for i in range(len(nums)):
    for j in range(len(caps)):
        adjMap[i][j + 1000] = 1

maxFlow = Dinic(adjMap)
count = maxFlow.calMaxFlow(START, END)
if count == sum(nums):
    print(1)
else:
    print(0)
    exit(0)
for i in range(len(nums)):
    res = []
    for j in range(len(caps)):
        if maxFlow.getFlowOfEdge(i, j + 1000) > 0:
            res.append(j + 1)
    print(*res)
