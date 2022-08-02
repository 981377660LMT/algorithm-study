# 0表示空地 1表示墙
# !现在要阻断左上角到右下角的路 问最少需要加多少墙
# https://binarysearch.com/problems/Walled-Off

# 最小割问题：
# you have a graph with two vertices,
# and you want to remove the minimum number of vertices such that the two original vertices are disconnected


from typing import DefaultDict, List
from collections import defaultdict, deque

Graph = DefaultDict[int, DefaultDict[int, int]]  # 有向带权图,权值为容量


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
        return self._graph[v1][v2] - self._reGraph[v1][v2]

    def _updateRedisualGraph(self) -> None:
        self._reGraph = defaultdict(lambda: defaultdict(int))
        for cur in self._graph:
            for next in self._graph[cur]:
                self._reGraph[cur][next] = self._graph[cur][next]
                self._reGraph[next].setdefault(cur, 0)


DIR4 = [[1, 0], [-1, 0], [0, 1], [0, -1]]

# 2 ≤ n, m ≤ 250


# !把边拆成点 in out  即把每个点拆成 `i 和 i + OFFSET`
# https://binarysearch.com/problems/Walled-Off/solutions/2865567
class Solution:
    def solve(self, matrix: List[List[int]]):
        ROW, COL = len(matrix), len(matrix[0])
        OFFSET = ROW * COL
        adjMap = defaultdict(lambda: defaultdict(int))
        for r in range(ROW):
            for c in range(COL):
                if matrix[r][c] == 1:
                    continue
                cur = r * COL + c
                adjMap[cur][cur + OFFSET] = 1  # !in 容量为1(割边费用为1)
                for dr, dc in DIR4:
                    nr, nc = r + dr, c + dc
                    if (0 <= nr < ROW) and (0 <= nc < COL) and (matrix[nr][nc] == 0):
                        next = nr * COL + nc
                        adjMap[cur + OFFSET][next] = int(1e20)  # !out 容量无限大

        return Dinic(adjMap).calMaxFlow(0 + OFFSET, ROW * COL - 1)


print(Solution().solve(matrix=[[0, 1, 0, 0], [0, 1, 0, 0], [0, 0, 0, 0]]))
