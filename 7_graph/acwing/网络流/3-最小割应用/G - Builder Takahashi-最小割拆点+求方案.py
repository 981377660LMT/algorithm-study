# 最小割拆点+求最小割的方案
# n<=100 稠密图
# 割每个`点`的成本为ci
# 求`阻断0 - (n-1)路径所需的最小代价`、 `需要割的点数` 以及 `需要割哪些点`
# 注意不能割0或n-1(即代价无穷大)

# !把边拆成点 in out  即把每个点拆成 `i 和 i + OFFSET`
import sys
import os

from collections import defaultdict
from typing import Set

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)
INF = int(4e18)


def main() -> None:
    n, m = map(int, input().split())
    OFFSET = 500
    maxFlow = MaxFlow(start=0 + OFFSET, end=n - 1)

    for _ in range(m):
        u, v = map(int, input().split())
        u, v = u - 1, v - 1
        # !点间：两点相连 出点连入点 容量为INF
        maxFlow.addEdge(v + OFFSET, u, INF)
        maxFlow.addEdge(u + OFFSET, v, INF)

    costs = list(map(int, input().split()))

    # !点内：割(1-n-2)每个点的成本为ci 入点连出点 容量为ci
    for i in range(1, n - 1):
        maxFlow.addEdge(i, i + OFFSET, costs[i])

    # !点内：注意0和n-1不能割 所以容量为INF
    maxFlow.addEdge(0, 0 + OFFSET, INF)
    maxFlow.addEdge(n - 1, n - 1 + OFFSET, INF)

    print(maxFlow.calMaxFlow())  # 最小割

    points = []
    visited = maxFlow.getPath()
    for i in range(1, n - 1):
        if (i in visited) ^ (i + OFFSET in visited):
            points.append(i + 1)

    print(len(points))
    print(*points)


if __name__ == "__main__":

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

    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
