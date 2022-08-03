"""Dinic算法 数组存边"""

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
        self._n = n + 5
        self._start = start
        self._end = end
        self._graph = [[] for _ in range(n + 5)]  # [next, remain, rEdge][]

    def addEdge(self, v1: int, v2: int, w: int) -> None:
        """添加边 v1->v2, 容量为w 注意重边影响"""
        forward = [v2, w, None]
        backward = [v1, 0, forward]
        forward[2] = backward  # type: ignore
        self._graph[v1].append(forward)
        self._graph[v2].append(backward)

    def calMaxFlow(self) -> int:
        flow = 0
        graph, INF, start, end = self._graph, self.INF, self._start, self._end
        while self._bfs():
            (*self._it,) = map(iter, graph)  # 当前弧优化
            delta = INF
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
        level, iters = self._level, self._it
        for edge in iters[cur]:  # type: ignore
            next, remain, rEdge = edge
            if remain and level[cur] < level[next]:  # type: ignore
                delta = self._dfs(next, target, min(flow, remain))
                if delta:
                    edge[1] -= delta
                    rEdge[1] += delta
                    return delta
        return 0


if __name__ == "__main__":

    # 给定一个包含 n 个点 m 条边的有向图，并给定每条边的容量，边的容量非负。
    # 图中可能存在重边和自环。求从点 S 到点 T 的最大流。

    import sys

    sys.setrecursionlimit(int(1e9))

    input = sys.stdin.readline
    n, m, start, end = map(int, input().split())
    adjMap = defaultdict(lambda: defaultdict(int))  # 用于去重

    # 从点 u 到点 v 存在一条有向边，容量为 c。
    for _ in range(m):
        u, v, w = map(int, input().split())
        adjMap[u][v] += w  # 可能存在重边

    maxFlow = Dinic(n + 10, start=start, end=end)
    for cur in adjMap:
        for next, cap in adjMap[cur].items():
            maxFlow.addEdge(cur, next, cap)

    print(maxFlow.calMaxFlow())
