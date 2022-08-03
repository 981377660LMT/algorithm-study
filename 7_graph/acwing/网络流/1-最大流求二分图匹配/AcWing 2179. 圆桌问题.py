from collections import defaultdict, deque
from typing import Set


INF = int(1e18)


class MaxFlow:
    def __init__(self, start: int, end: int) -> None:
        self.graph = defaultdict(lambda: defaultdict(int))  # 原图
        self._start = start
        self._end = end

    def calMaxFlow(self) -> int:
        self._updateRedisualGraph()
        start, end = self._start, self._end
        flow = 0

        while self._bfs():
            delta = INF
            while delta:
                delta = self._dfs(start, end, INF)
                flow += delta
        return flow

    def addEdge(self, v1: int, v2: int, w: int, *, cover=False) -> None:
        """添加边 v1->v2, 容量为w

        Args:
            v1: 边的起点
            v2: 边的终点
            w: 边的容量
            cover: 是否覆盖原有边
        """
        if cover:
            self.graph[v1][v2] = w
        else:
            self.graph[v1][v2] += w

    def getFlowOfEdge(self, v1: int, v2: int) -> int:
        """边的流量=容量-残量"""
        assert v1 in self.graph and v2 in self.graph[v1]
        return self.graph[v1][v2] - self._reGraph[v1][v2]

    def getRemainOfEdge(self, v1: int, v2: int) -> int:
        """边的残量(剩余的容量)"""
        assert v1 in self.graph and v2 in self.graph[v1]
        return self._reGraph[v1][v2]

    def getPath(self) -> Set[int]:
        """最大流经过了哪些点"""
        visited = set()
        stack = [self._start]
        reGraph = self._reGraph
        while stack:
            cur = stack.pop()
            visited.add(cur)
            for next, remain in reGraph[cur].items():
                if next not in visited and remain > 0:
                    visited.add(next)
                    stack.append(next)
        return visited

    def _updateRedisualGraph(self) -> None:
        """残量图 存储每条边的剩余流量"""
        self._reGraph = defaultdict(lambda: defaultdict(int))
        for cur in self.graph:
            for next, cap in self.graph[cur].items():
                self._reGraph[cur][next] = cap
                self._reGraph[next].setdefault(cur, 0)  # 注意自环边

    def _bfs(self) -> bool:
        self._depth = depth = defaultdict(lambda: -1, {self._start: 0})
        reGraph, start, end = self._reGraph, self._start, self._end
        queue = deque([start])
        self._iters = {cur: iter(reGraph[cur].keys()) for cur in reGraph.keys()}
        while queue:
            cur = queue.popleft()
            nextDist = depth[cur] + 1
            for next, remain in reGraph[cur].items():
                if depth[next] == -1 and remain > 0:
                    depth[next] = nextDist
                    queue.append(next)

        return depth[end] != -1

    def _dfs(self, cur: int, end: int, flow: int) -> int:
        if cur == end:
            return flow
        reGraph, depth, iters = self._reGraph, self._depth, self._iters
        for next in iters[cur]:
            remain = reGraph[cur][next]
            if remain and depth[cur] < depth[next]:
                nextFlow = self._dfs(next, end, min(flow, remain))
                if nextFlow:
                    reGraph[cur][next] -= nextFlow
                    reGraph[next][cur] += nextFlow
                    return nextFlow
        return 0


# 假设有来自 m 个不同单位的代表参加一次国际会议。
# 每个单位的代表数分别为 ri(i=1,2,…,m)。
# 会议餐厅共有 n 张餐桌，每张餐桌可容纳 ci(i=1,2,…,n) 个代表就餐。
# !为了使代表们充分交流，希望从同一个单位来的代表不在同一个餐桌就餐。
# 试设计一个算法，给出满足要求的代表就餐方案。
# 如果问题有解，在第 1 行输出 1，否则输出 0。
# 接下来的 m 行给出每个单位代表的就餐桌号。
# 如果有多个满足要求的方案，只要求输出 1 个方案。

# m<=150
# n<=270
# 最大流建模： 源点 => 人分组 => 餐桌分组 => 汇点
m, n = map(int, input().split())  # m: 代表数量, n: 餐桌数量
adjMap = defaultdict(lambda: defaultdict(int))
nums = list(map(int, input().split()))
caps = list(map(int, input().split()))
START, END, OFFSET = -1, int(1e9), 1000
maxFlow = MaxFlow(START, END)
for i, num in enumerate(nums):
    maxFlow.addEdge(START, i, num)
for i, cap in enumerate(caps):
    maxFlow.addEdge(i + OFFSET, END, cap)
for i in range(len(nums)):
    for j in range(len(caps)):
        maxFlow.addEdge(i, j + OFFSET, 1)


count = maxFlow.calMaxFlow()
if count == sum(nums):
    print(1)
else:
    print(0)
    exit(0)

for i in range(len(nums)):
    res = []
    for j in range(len(caps)):
        if maxFlow.getFlowOfEdge(i, j + OFFSET) > 0:
            res.append(j + 1)
    print(*res)
