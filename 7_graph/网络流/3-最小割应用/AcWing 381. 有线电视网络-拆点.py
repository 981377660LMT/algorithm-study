# 给定一张 n 个点 m 条边的无向图，求最少去掉多少个点，可以使图不连通。
# 如果不管去掉多少个点，都无法使原图不连通，则直接返回 n。
# !0≤n≤50


from collections import defaultdict, deque
from typing import DefaultDict, Set
import sys


sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


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

    for start in range(n):
        for end in range(n):
            if start == end:
                continue
            maxFlow = MaxFlow(start, end)
            for cur, nexts in rawAdjMap.items():
                for next in nexts:
                    maxFlow.addEdge(cur + OFFSET, next, INF)  # a出点 => b入点容量无穷 因为要割的是点，不是边
            for i in range(n):
                maxFlow.addEdge(i, OFFSET + n, 1)  # a入点=>a出点容量为1 因为要割的是点

            res = min(res, maxFlow.calMaxFlow())  # 源点是start的出点 ，汇点是end的入点
    return res


while True:
    try:
        n, m, *rest = list(input().split())
        adjMap = defaultdict(set)  # 无向图
        for item in rest:
            u, v = map(int, item[1:-1].split(","))
            adjMap[u].add(v)
            adjMap[v].add(u)
        print(main(int(n), adjMap))
    except EOFError:
        break
