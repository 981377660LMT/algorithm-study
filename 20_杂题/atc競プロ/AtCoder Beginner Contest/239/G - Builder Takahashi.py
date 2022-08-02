# 最小割拆点
# n<=100 稠密图
# 割每个`点`的成本为ci
# 求`阻断0 - (n-1)路径所需的最小代价`、 `需要割的点数` 以及 `需要割哪些点`
# 注意不能割0或n-1(即代价无穷大)

# !把边拆成点 in out  即把每个点拆成 `i 和 i + OFFSET`
import sys
import os

from collections import defaultdict
from typing import DefaultDict, Set

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)
INF = int(4e18)


def main() -> None:
    n, m = map(int, input().split())
    OFFSET = 500
    maxFlow = Dinic(int(2e3), start=0 + OFFSET, end=n - 1)

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

    class Dinic:
        INF = int(1e18)

        def __init__(self, n: int, *, start: int, end: int) -> None:
            """n为节点个数 start为源点 end为汇点"""
            self._n = n
            self._start = start
            self._end = end
            self._graph = [[] for _ in range(n)]

        def addEdge(self, v1: int, v2: int, w: int) -> None:
            """添加边 v1->v2, 容量为w"""
            forward = [v2, w, None]
            backward = [v1, 0, forward]
            forward[2] = backward  # type: ignore
            self._graph[v1].append(forward)
            self._graph[v2].append(backward)

        def calMaxFlow(self) -> int:
            flow = 0
            graph, INF, start, end = self._graph, self.INF, self._start, self._end
            while self._bfs():
                (*self._it,) = map(iter, graph)
                delta = self.INF
                while delta:
                    delta = self._dfs(start, end, INF)
                    flow += delta
            return flow

        def getEdgeRemain(self) -> DefaultDict[int, DefaultDict[int, int]]:
            """边的残量(剩余的容量)"""
            res = defaultdict(lambda: defaultdict(int))
            for edge in self._graph:
                pre, next, remain = edge[2][0], edge[0], edge[1]
                res[pre][next] = remain
            return res

        def getPath(self) -> Set[int]:
            """最大流经过了哪些点"""
            visited = set()
            stack = [self._start]
            while stack:
                cur = stack.pop()
                visited.add(cur)
                for next, remain, _ in self._graph[cur]:
                    if next not in visited and remain > 0:
                        visited.add(next)
                        stack.append(next)
            return visited

        def _bfs(self) -> bool:
            """建立分层图"""
            self._level = level = [None] * self._n
            start, end = self._start, self._end
            queue = deque([start])
            level[start] = 0  # type: ignore
            graph = self._graph
            while queue:
                cur = queue.popleft()
                dist = level[cur] + 1  # type: ignore
                for next, remain, _ in graph[cur]:
                    if level[next] is None and remain > 0:
                        level[next] = dist
                        queue.append(next)
            return level[end] is not None

        def _dfs(self, cur: int, target: int, flow: int) -> int:
            """寻找增广路"""
            if cur == target:
                return flow
            level = self._level
            for edge in self._it[cur]:  # type: ignore
                next, remain, rEdge = edge
                if remain and level[cur] < level[next]:  # type: ignore
                    delta = self._dfs(next, target, min(flow, remain))
                    if delta:
                        edge[1] -= delta
                        rEdge[1] += delta
                        return delta
            return 0

    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
