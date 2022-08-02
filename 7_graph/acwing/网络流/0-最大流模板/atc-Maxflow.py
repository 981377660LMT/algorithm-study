# atcoder模板

from collections import defaultdict, deque
from typing import DefaultDict, Set


class Dinic:
    """Dinic 求最大流

    如果一个流的残量网络里面没有可行流，那么这个流就是最大流

    时间复杂度:O(V^2*E)
    """

    INF = int(1e18)

    def __init__(self, n: int, *, start: int, end: int) -> None:
        """n为节点个数 start为源点 end为汇点"""
        self._n = n
        self._start = start
        self._end = end
        self._graph = [[] for _ in range(n)]  # [next, remain, rEdge][]

    def addEdge(self, v1: int, v2: int, w: int) -> None:
        """添加边 v1->v2, 容量为w"""
        forward = [v2, w, None]
        backward = [v1, 0, forward]
        forward[2] = backward  # type: ignore
        self._graph[v1].append(forward)
        self._graph[v2].append(backward)

    def addMultiEdge(self, v1: int, v2: int, w1: int, w2: int) -> None:
        """针对重边的情况 需要将重边处理后一起添加

        添加边 v1->v2, 容量为w1
        添加边 v2->v1, 容量为w2
        """
        edge1 = [v2, w1, None]
        edge2 = [v1, w2, edge1]
        edge1[2] = edge2  # type: ignore
        self._graph[v1].append(edge1)
        self._graph[v2].append(edge2)

    def calMaxFlow(self) -> int:
        flow = 0
        graph, INF, start, end = self._graph, self.INF, self._start, self._end
        while self._bfs():
            (*self._it,) = map(iter, graph)  # 当前弧优化
            delta = self.INF
            while delta:
                delta = self._dfs(start, end, INF)
                flow += delta
        return flow

    def getEdgeRemain(self) -> DefaultDict[int, DefaultDict[int, int]]:
        """边的残量(剩余的容量)"""
        res = defaultdict(lambda: defaultdict(int))
        for edges in self._graph:
            for edge in edges:
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
