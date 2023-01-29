from collections import defaultdict, deque
from typing import List, Set

INF = int(1e18)


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


# 相邻两个1组成一条边，每条边都要去掉一个端点，其实是找最小点覆盖，即求二分图的最大匹配，跑匈牙利算法
# !起点 => 奇数黑格 => 偶数黑格 => 终点
DIR2 = [(0, 1), (1, 0)]


class Solution:
    def minimumOperations(self, grid: List[List[int]]) -> int:
        ROW, COL = len(grid), len(grid[0])
        n = ROW * COL
        START, END = n + 4, n + 5
        maxFlow = ATCMaxFlow(n + 10, start=START, end=END)

        for r in range(ROW):
            for c in range(COL):
                if grid[r][c] == 1:
                    cur = r * COL + c
                    for dr, dc in DIR2:
                        nr, nc = r + dr, c + dc
                        if 0 <= nr < ROW and 0 <= nc < COL and grid[nr][nc] == 1:
                            next = nr * COL + nc
                            v1, v2 = (next, cur) if (r + c) & 1 else (cur, next)
                            maxFlow.addEdgeIfAbsent(v1, v2, 1)
                            maxFlow.addEdgeIfAbsent(START, v1, 1)
                            maxFlow.addEdgeIfAbsent(v2, END, 1)

        return maxFlow.calMaxFlow()


print(Solution().minimumOperations(grid=[[1, 1, 0], [0, 1, 1], [1, 1, 1]]))
