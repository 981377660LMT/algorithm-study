# n个房间 进入每个房间必须有钥匙且交钱w 房间内有解锁其他房间的钥匙和钱
# 房间i里的钥匙编号大于等于i+1
# 问如何进入若干个房间 使得最后收益最大
# n<=100
# https://twitter.com/e869120/status/1393336369540341760/photo/1
# !燃やす埋める問題に帰着させる 最小割为答案(房子分成两类 进还是不进 再计算最小罚款 让每个房子都有`进/不进`的归属)
# !最大流最小割 时间复杂度O(V^2E)


import sys


sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(1e18)

from collections import defaultdict, deque
from typing import Set


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

    def addEdge(self, v1: int, v2: int, w: int, *, cover=True) -> None:
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


# 最大流代表消耗的金钱

n, w = map(int, input().split())
moneys = list(map(int, input().split()))
keys = [[] for _ in range(n)]
for i in range(n):
    count, *rest = list(map(int, input().split()))
    keys[i].extend([num - 1 for num in rest])

START, END = -1, -2
maxFlow = MaxFlow(START, END)


# !在sum(moneys)的基础上计算罚款
# !源点:不访问
# !汇点:访问
for i in range(n):
    maxFlow.addEdge(START, i, w, cover=True)  # 割不访问　需要罚款w
    maxFlow.addEdge(i, END, moneys[i], cover=True)  # 割访问　需要罚款moneys[i]

for i in range(n):
    for key in keys[i]:
        maxFlow.addEdge(i, key, INF, cover=True)  # 割去掉这条边表示 访问需要key的房子却没有访问key所在的房子 罚款1e20


minCut = maxFlow.calMaxFlow()
print(sum(moneys) - minCut)  # 总收益-最小割(分成两半的最少罚款)
