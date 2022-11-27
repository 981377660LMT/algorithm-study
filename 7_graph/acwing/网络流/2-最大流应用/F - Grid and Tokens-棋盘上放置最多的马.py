# F - Grid and Tokens-棋盘上放置最多的马

# 有 n 颗棋子欲放到 H ∗ W 的网格中，每颗棋子可放的`范围`是一个矩阵，
# 且不能有两颗棋子待在同一行或同一列，求最多可以放多少颗棋子。
# row,col,n<=100

# !源点 => 行 => 棋子顶点u => 棋子顶点v => 列 => 汇点
from typing import List, Tuple


def gridAndTokens(ROW: int, COL: int, tokens: List[Tuple[int, int, int, int]]) -> int:
    n = len(tokens)
    START = ROW + n + n + COL
    END = START + 1
    mf = ATCMaxFlow(END + 1, START, END)

    for r in range(ROW):
        mf.addEdge(START, r, 1)  # !源点 => 行

    for i, (r1, _, r2, _) in enumerate(tokens):
        for r in range(r1, r2 + 1):
            mf.addEdge(r, ROW + i, 1)  # !行 => 棋子顶点u

    for i in range(n):
        mf.addEdge(ROW + i, ROW + n + i, 1)  # !棋子拆点 u => v

    for i, (_, c1, _, c2) in enumerate(tokens):
        for c in range(c1, c2 + 1):
            mf.addEdge(ROW + n + i, ROW + n + n + c, 1)  # !棋子顶点v => 列

    for c in range(COL):
        mf.addEdge(ROW + n + n + c, END, 1)  # !列 => 汇点

    return mf.calMaxFlow()


if __name__ == "__main__":
    from collections import defaultdict, deque
    from typing import Set

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
            n, reGraph, start, end, edges = (
                self._n,
                self._reGraph,
                self._start,
                self._end,
                self._edges,
            )
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
            reGraph, level, curEdges, edges = (
                self._reGraph,
                self._levels,
                self._curEdges,
                self._edges,
            )
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

    row, col, n = map(int, input().split())
    tokens = []
    for _ in range(n):
        r1, c1, r2, c2 = map(int, input().split())
        r1, c1, r2, c2 = r1 - 1, c1 - 1, r2 - 1, c2 - 1
        tokens.append((r1, c1, r2, c2))
    print(gridAndTokens(row, col, tokens))
