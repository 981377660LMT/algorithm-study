# 最小割拆点
# n<=100 稠密图
# 割每个`点`的成本为ci
# 求`阻断0 - (n-1)路径所需的最小代价`、 `需要割的点数` 以及 `需要割哪些点`
# 注意不能割0或n-1(即代价无穷大)

# !把边拆成点 in out  即把每个点拆成 `i 和 i + OFFSET`
import sys
import os

from collections import defaultdict

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)
INF = int(4e18)


def main() -> None:
    n, m = map(int, input().split())
    OFFSET = int(1e5)
    maxFlow = Dinic(0 + OFFSET, n - 1)

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
    for i in range(n):
        remain = maxFlow.getRemainOfEdge(i, i + OFFSET)
        if remain == 0:
            points.append(i + 1)

    print(len(points))
    print(*points)


if __name__ == "__main__":

    from collections import defaultdict, deque

    class Dinic:
        INF = int(1e18)

        def __init__(self, start: int, end: int) -> None:
            self._graph = defaultdict(lambda: defaultdict(int))
            self._start = start
            self._end = end

        def calMaxFlow(self) -> int:
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
            start, end = self._start, self._end
            res = 0
            depth = defaultdict(lambda: -1, {start: 0})
            curArc = dict()

            while True:
                bfs()
                if depth[end] != -1:
                    while True:
                        delta = dfsWithCurArc(start, Dinic.INF)
                        if delta == 0:
                            break
                        res += delta
                else:
                    break
            return res

        def addEdge(self, v1: int, v2: int, w: int) -> None:
            """添加边 v1->v2, 容量为w"""
            self._graph[v1][v2] += w

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

    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
