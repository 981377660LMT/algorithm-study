from collections import defaultdict
from typing import DefaultDict, List
from collections import defaultdict, deque

Graph = DefaultDict[int, DefaultDict[int, int]]  # 有向带权图,权值为容量

# 1 <= m, n <= 300


class Dinic:
    def __init__(self, graph: Graph) -> None:
        self._graph = graph

    def calMaxFlow(self, start: int, end: int) -> int:
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
                try:
                    child = next(curArc[cur])
                    if (depth[child] == depth[cur] + 1) and (self._reGraph[cur][child] > 0):
                        nextFlow = dfsWithCurArc(
                            child, min(minFlow - flow, self._reGraph[cur][child])
                        )
                        if nextFlow == 0:
                            depth[child] = -1
                        self._reGraph[cur][child] -= nextFlow
                        self._reGraph[child][cur] += nextFlow
                        flow += nextFlow
                except StopIteration:
                    break
            return flow

        self._updateRedisualGraph()

        res = 0
        depth = defaultdict(lambda: -1, {start: 0})
        curArc = dict()

        while True:
            bfs()
            if depth[end] != -1:
                while True:
                    delta = dfsWithCurArc(start, int(1e20))
                    if delta == 0:
                        break
                    res += delta
            else:
                break
        return res

    def getFlowOfEdge(self, v1: int, v2: int) -> int:
        assert v1 in self._graph and v2 in self._graph[v1]
        return self._graph[v1][v2] - self._reGraph[v1][v2]

    def _updateRedisualGraph(self) -> None:
        self._reGraph = defaultdict(lambda: defaultdict(int))
        for cur in self._graph:
            for next in self._graph[cur]:
                self._reGraph[cur][next] = self._graph[cur][next]
                self._reGraph[next].setdefault(cur, 0)


# 相邻两个1组成一条边，每条边都要去掉一个端点，其实是找最小点覆盖，即求二分图的最大匹配，跑匈牙利算法
class Solution:
    def minimumOperations(self, grid: List[List[int]]) -> int:
        ROW, COL = len(grid), len(grid[0])
        adjMap = defaultdict(lambda: defaultdict(int))
        for r in range(ROW):
            for c in range(COL):
                if grid[r][c] == 1:
                    cur = r * COL + c
                    for dr, dc in [(0, 1), (1, 0)]:
                        nr, nc = r + dr, c + dc
                        if 0 <= nr < ROW and 0 <= nc < COL and grid[nr][nc] == 1:
                            next = nr * COL + nc
                            v1, v2 = (next, cur) if (r + c) & 1 else (cur, next)
                            adjMap[-1][v1] = 1
                            adjMap[v1][v2] = 1
                            adjMap[v2][int(1e9)] = 1
        return Dinic(adjMap).calMaxFlow(-1, int(1e9))


print(Solution().minimumOperations(grid=[[1, 1, 0], [0, 1, 1], [1, 1, 1]]))
