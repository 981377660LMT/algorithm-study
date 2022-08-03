# 0表示空地 1表示墙
# !现在要阻断左上角到右下角的路 问最少需要加多少墙
# https://binarysearch.com/problems/Walled-Off

# 最小割问题：
# you have a graph with two vertices,
# and you want to remove the minimum number of vertices such that the two original vertices are disconnected


from collections import defaultdict, deque
from typing import List, Set

INF = int(1e18)


class MaxFlow:
    def __init__(self, start: int, end: int) -> None:
        self.graph = defaultdict(lambda: defaultdict(int))  # 原图
        self._start = start
        self._end = end

    def calMaxFlow(self) -> int:
        self._updateRedisualGraph()
        start, end = self._start, self._end
        flow = 0

        while self._bfs():
            delta = INF
            while delta:
                delta = self._dfs(start, end, INF)
                flow += delta
        return flow

    def addEdge(self, v1: int, v2: int, w: int, *, cover=True) -> None:
        """添加边 v1->v2, 容量为w

        Args:
            v1: 边的起点
            v2: 边的终点
            w: 边的容量
            cover: 是否覆盖原有边
        """
        if cover:
            self.graph[v1][v2] = w
        else:
            self.graph[v1][v2] += w

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
                if next not in visited and remain > 0:
                    visited.add(next)
                    stack.append(next)
        return visited

    def _updateRedisualGraph(self) -> None:
        """残量图 存储每条边的剩余流量"""
        self._reGraph = defaultdict(lambda: defaultdict(int))
        for cur in self.graph:
            for next, cap in self.graph[cur].items():
                self._reGraph[cur][next] = cap
                self._reGraph[next].setdefault(cur, 0)  # 注意自环边

    def _bfs(self) -> bool:
        self._depth = depth = defaultdict(lambda: -1, {self._start: 0})
        reGraph, start, end = self._reGraph, self._start, self._end
        queue = deque([start])
        self._iters = {cur: iter(reGraph[cur].keys()) for cur in reGraph.keys()}
        while queue:
            cur = queue.popleft()
            nextDist = depth[cur] + 1
            for next, remain in reGraph[cur].items():
                if depth[next] == -1 and remain > 0:
                    depth[next] = nextDist
                    queue.append(next)

        return depth[end] != -1

    def _dfs(self, cur: int, end: int, flow: int) -> int:
        if cur == end:
            return flow
        reGraph, depth, iters = self._reGraph, self._depth, self._iters
        for next in iters[cur]:
            remain = reGraph[cur][next]
            if remain and depth[cur] < depth[next]:
                nextFlow = self._dfs(next, end, min(flow, remain))
                if nextFlow:
                    reGraph[cur][next] -= nextFlow
                    reGraph[next][cur] += nextFlow
                    return nextFlow
        return 0


DIR4 = [[1, 0], [0, 1]]

# 2 ≤ n, m ≤ 250


# !把边拆成点 in out  即把每个点拆成 `i 和 i + OFFSET`
# https://binarysearch.com/problems/Walled-Off/solutions/2865567
class Solution:
    def solve(self, matrix: List[List[int]]):
        ROW, COL = len(matrix), len(matrix[0])
        n = OFFSET = ROW * COL

        maxFlow = MaxFlow(0 + OFFSET, n - 1)
        for r in range(ROW):
            for c in range(COL):
                if matrix[r][c] == 1:
                    continue
                cur = r * COL + c
                maxFlow.addEdge(cur, cur + OFFSET, 1, cover=True)  # !in 容量为1(割边费用为1)
                for dr, dc in DIR4:
                    nr, nc = r + dr, c + dc
                    if (0 <= nr < ROW) and (0 <= nc < COL) and (matrix[nr][nc] == 0):
                        next = nr * COL + nc
                        maxFlow.addEdge(cur + OFFSET, next, INF, cover=True)  # !out 容量无限大
                        maxFlow.addEdge(next + OFFSET, cur, INF, cover=True)  # !out 容量无限大

        return maxFlow.calMaxFlow()


print(Solution().solve(matrix=[[0, 1, 0, 0], [0, 1, 0, 0], [0, 0, 0, 0]]))
