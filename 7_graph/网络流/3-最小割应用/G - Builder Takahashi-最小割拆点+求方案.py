# 最小割拆点+求最小割的方案
# n<=100 稠密图
# 割每个`点`的成本为ci(使与i相连的边全部消失)
# 求`阻断0 - (n-1)路径所需的最小代价`、 `需要割的点数` 以及 `需要割哪些点`
# 注意不能割0或n-1(即代价无穷大)
# https://yuuko.moe/index.php/archives/175/
# !把边拆成点 in out  即把每个点拆成 `i 和 i + OFFSET`
import sys
import os

from collections import defaultdict
from typing import Set

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
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
        # 残量网络上出点和入点是否连通，如果不连通才能说明选择了
        if (i in visited) ^ (i + OFFSET in visited):
            points.append(i + 1)

    print(len(points))
    print(*points)


if __name__ == "__main__":

    from collections import defaultdict, deque
    from typing import Set

    class MaxFlow:
        """Dinic算法 字典存残量图 速度一般"""

        def __init__(self, start: int, end: int) -> None:
            self.graph = defaultdict(lambda: defaultdict(int))  # 原图
            self._reGraph = defaultdict(lambda: defaultdict(int))  # 残量图
            self._start = start
            self._end = end

        def calMaxFlow(self) -> int:
            reGraph, start, end = self._reGraph, self._start, self._end
            res = 0

            while self._bfs():
                self._iters = {cur: iter(reGraph[cur].keys()) for cur in reGraph}
                res += self._dfs(start, end, INF)
            return res

        def addEdge(self, v1: int, v2: int, w: int, *, cover=True) -> None:
            """添加边 v1->v2, 容量为w

            Args:
                v1: 边的起点
                v2: 边的终点
                w: 边的容量
                cover: 是否覆盖原有边 默认为覆盖
            """
            if cover:
                self.graph[v1][v2] = w
                self._reGraph[v1][v2] = w
                self._reGraph[v2].setdefault(v1, 0)  # 注意自环边
            else:
                self.graph[v1][v2] += w
                self._reGraph[v1][v2] += w
                self._reGraph[v2].setdefault(v1, 0)

        def getFlowOfEdge(self, v1: int, v2: int) -> int:
            """边的流量=容量-残量"""
            assert v1 in self.graph and v2 in self.graph[v1]
            return self.graph[v1][v2] - self._reGraph[v1][v2]

        def getRemainOfEdge(self, v1: int, v2: int) -> int:
            """边的残量(剩余的容量)"""
            assert v1 in self.graph and v2 in self.graph[v1]
            return self._reGraph[v1][v2]

        def getPath(self) -> Set[int]:
            """最大流经过了(残量网络上的)哪些点"""
            visited = set()
            stack = [self._start]
            reGraph = self._reGraph
            while stack:
                cur = stack.pop()
                visited.add(cur)
                for next, remain in reGraph[cur].items():
                    if remain > 0 and next not in visited:
                        visited.add(next)
                        stack.append(next)
            return visited

        def _bfs(self) -> bool:
            self._depth = depth = defaultdict(lambda: -1, {self._start: 0})
            reGraph, start, end = self._reGraph, self._start, self._end
            queue = deque([start])

            while queue:
                cur = queue.popleft()
                nextDist = depth[cur] + 1
                for next, remain in reGraph[cur].items():
                    if remain > 0 and depth[next] == -1:
                        depth[next] = nextDist
                        if next == end:
                            return True
                        queue.append(next)

            return False

        def _dfs(self, cur: int, end: int, flow: int) -> int:
            if cur == end:
                return flow
            res = flow
            reGraph, depth, iters = self._reGraph, self._depth, self._iters
            for next in iters[cur]:
                remain = reGraph[cur][next]
                if remain > 0 and depth[cur] + 1 == depth[next]:
                    delta = self._dfs(next, end, min(res, remain))
                    reGraph[cur][next] -= delta
                    reGraph[next][cur] += delta
                    res -= delta
                    if res == 0:
                        return flow

            return flow - res

    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
