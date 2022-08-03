from collections import defaultdict, deque
from typing import Set

INF = int(1e18)


class MaxFlow:
    """Dinic算法 字典存残量图 比较慢"""

    __slots__ = ("graph", "_reGraph", "_start", "_end", "_levels", "_iters")

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
        """最大流经过了哪些点"""
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
        self._levels = level = defaultdict(lambda: -1, {self._start: 0})
        reGraph, start, end = self._reGraph, self._start, self._end
        queue = deque([start])

        while queue:
            cur = queue.popleft()
            nextDist = level[cur] + 1
            for next, remain in reGraph[cur].items():
                if remain > 0 and level[next] == -1:
                    level[next] = nextDist
                    if next == end:
                        return True
                    queue.append(next)

        return False

    def _dfs(self, cur: int, end: int, flow: int) -> int:
        if cur == end:
            return flow
        res = flow
        reGraph, level, iters = self._reGraph, self._levels, self._iters
        for next in iters[cur]:
            remain = reGraph[cur][next]
            if remain > 0 and level[cur] + 1 == level[next]:
                delta = self._dfs(next, end, min(res, remain))
                reGraph[cur][next] -= delta
                reGraph[next][cur] += delta
                res -= delta
                if res == 0:
                    return flow

        return flow - res


class ATCMaxFlow:
    """Dinic算法 数组+边存图 速度较快"""

    __slots__ = (
        "_n",
        "_start",
        "_end",
        "_reGraph",
        "_edges",
        "_visitedEdge",
        "_levels",
        "_curEdges",
    )

    def __init__(self, n: int, *, start: int, end: int) -> None:
        if not (0 <= start < n and 0 <= end < n):
            raise ValueError(f"start: {start}, end: {end} out of range [0,{n}]")

        self._n = n
        self._start = start
        self._end = end
        self._reGraph = [[] for _ in range(n)]  # 残量图存边的序号
        self._edges = []  # [next,capacity]

        self._visitedEdge = set()

        self._levels = [0] * n
        self._curEdges = [0] * n

    def addEdge(self, v1: int, v2: int, capacity: int) -> None:
        """添加边 v1->v2, 容量为w 注意会添加重边"""
        self._visitedEdge.add((v1, v2))
        self._reGraph[v1].append(len(self._edges))
        self._edges.append([v2, capacity])
        self._reGraph[v2].append(len(self._edges))
        self._edges.append([v1, 0])

    def addEdgeIfAbsent(self, v1: int, v2: int, capacity: int) -> None:
        """如果边不存在则添加边 v1->v2, 容量为w"""
        if (v1, v2) in self._visitedEdge:
            return
        self._visitedEdge.add((v1, v2))
        self._reGraph[v1].append(len(self._edges))
        self._edges.append([v2, capacity])
        self._reGraph[v2].append(len(self._edges))
        self._edges.append([v1, 0])

    def calMaxFlow(self) -> int:
        n, start, end = self._n, self._start, self._end
        res = 0

        while self._bfs():
            self._curEdges = [0] * n
            res += self._dfs(start, end, INF)
        return res

    def getPath(self) -> Set[int]:
        """最大流经过了哪些点"""
        visited = set()
        queue = [self._start]
        reGraph, edges = self._reGraph, self._edges
        while queue:
            cur = queue.pop()
            visited.add(cur)
            for ei in reGraph[cur]:
                edge = edges[ei]
                next, remain = edge
                if remain > 0 and next not in visited:
                    visited.add(next)
                    queue.append(next)
        return visited

    def useQueryRemainOfEdge(self):
        """求边的残量(剩余的容量)::

        ```python
        maxFlow = ATCMaxFlow(n, start, end)
        query = maxFlow.useQueryRemainOfEdge()
        edgeRemain = query(v1, v2)
        ```
        """

        def query(v1: int, v2: int) -> int:
            return adjList[v1][v2]

        n, reGraph, edges = self._n, self._reGraph, self._edges
        adjList = [defaultdict(int) for _ in range(n)]
        for cur in range(n):
            for ei in reGraph[cur]:
                edge = edges[ei]
                next, remain = edge
                adjList[cur][next] += remain

        return query

    def _bfs(self) -> bool:
        n, reGraph, start, end, edges = self._n, self._reGraph, self._start, self._end, self._edges
        self._levels = level = [-1] * n
        level[start] = 0
        queue = deque([start])

        while queue:
            cur = queue.popleft()
            nextDist = level[cur] + 1
            for ei in reGraph[cur]:
                next, remain = edges[ei]
                if remain > 0 and level[next] == -1:
                    level[next] = nextDist
                    if next == end:
                        return True
                    queue.append(next)

        return False

    def _dfs(self, cur: int, end: int, flow: int) -> int:
        if cur == end:
            return flow
        res = flow
        reGraph, level, curEdges, edges = self._reGraph, self._levels, self._curEdges, self._edges
        ei = curEdges[cur]
        while ei < len(reGraph[cur]):
            ej = reGraph[cur][ei]
            next, remain = edges[ej]
            if remain > 0 and level[cur] + 1 == level[next]:
                delta = self._dfs(next, end, min(res, remain))
                edges[ej][1] -= delta
                edges[ej ^ 1][1] += delta
                res -= delta
                if res == 0:
                    return flow
            curEdges[cur] += 1
            ei = curEdges[cur]

        return flow - res


if __name__ == "__main__":
    # 给定一个包含 n 个点 m 条边的有向图，并给定每条边的容量，边的容量非负。
    # 图中可能存在重边和自环。求从点 S 到点 T 的最大流。
    import sys

    sys.setrecursionlimit(int(1e9))

    input = sys.stdin.readline
    n, m, start, end = map(int, input().split())
    maxFlow = MaxFlow(start, end)

    for _ in range(m):
        u, v, c = map(int, input().split())
        maxFlow.addEdge(u, v, c)  # 可能存在重边

    print(maxFlow.calMaxFlow())
