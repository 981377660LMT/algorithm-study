# 参加考试的最大学生数最大独立集解法
# https://leetcode.cn/problems/maximum-students-taking-exam/solution/wang-luo-liu-er-fen-tu-zui-da-du-li-ji-jian-dan-fa/

# !O((ROW*COL)^3)

from typing import List, Set
from collections import defaultdict, deque

INF = int(1e18)


class ATCMaxFlow:
    """Dinic算法 数组+边存图"""

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

    def __init__(self, n: int, start: int, end: int) -> None:
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


# 学生可以看到左侧、右侧、左上、右上这四个方向上紧邻他的学生的答卷 (即这几个方向不能有两个学生)
# 只关心可以坐人的座位，我们把作为按照`列下标`的奇偶建二分图，
# !S向偶数下标的座位连边，奇数下标的座位向T连边，有冲突的座位偶数座位向奇数座位连边
DIR8 = set([(0, 1), (1, 0), (0, -1), (-1, 0), (1, 1), (1, -1), (-1, 1), (-1, -1)])
DIR2 = set([(-1, 0), (1, 0)])
DIR6 = DIR8 - DIR2


class Solution:
    def maxStudents(self, seats: List[List[str]]) -> int:
        """注意这里要按照`列下标`的奇偶性来建立二分图 而不是按照(行+列)奇偶性建立二分图"""
        ROW, COL = len(seats), len(seats[0])
        START, END = ROW * COL + 2, ROW * COL + 3
        maxFlow = ATCMaxFlow(ROW * COL + 5, START, END)
        res = 0
        for r in range(ROW):
            for c in range(COL):
                if seats[r][c] == "#":
                    continue
                res += 1
                cur = r * COL + c
                if c & 1:
                    maxFlow.addEdge(cur, END, 1)
                else:
                    maxFlow.addEdge(START, cur, 1)
                    for dr, dc in DIR6:  # !不能与这六个方向有连接
                        nr, nc = r + dr, c + dc
                        if 0 <= nr < ROW and 0 <= nc < COL and seats[nr][nc] == ".":
                            maxFlow.addEdgeIfAbsent(cur, nr * COL + nc, 1)

        return res - maxFlow.calMaxFlow()
