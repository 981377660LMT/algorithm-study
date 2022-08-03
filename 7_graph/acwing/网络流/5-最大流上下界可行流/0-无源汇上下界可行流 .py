# 给定一个包含 n 个点 m 条边的有向图，每条边都有一个流量下界和流量上界。
# !可行流：求一种`可行方案`使得在所有点满足流量平衡条件的前提下，所有边满足流量限制
# https://www.acwing.com/problem/content/submission/2190/

# 1≤n≤200,
# 1≤m≤10200,

# TODO 有问题


from collections import defaultdict, deque
from typing import Set

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

    def addEdge(self, v1: int, v2: int, w: int, *, cover=False) -> None:
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


# 无源汇上下界可行流
n, m = map(int, input().split())
edge = defaultdict(lambda: [0, 0])
diff = defaultdict(int)
START, END, OFFSET = -1, -2, int(1e4)
maxFlow = MaxFlow(START, END)
for _ in range(m):
    u, v, lower, upper = map(int, input().split())
    edge[(u, v)][0] = lower
    edge[(u, v)][1] = upper
    diff[u] -= lower
    diff[v] += lower
    maxFlow.addEdge(u, v, upper - lower)


fullFlow = 0
for key, delta in diff.items():
    if delta > 0:
        fullFlow += delta
        maxFlow.addEdge(START, key, delta)
    elif delta < 0:
        maxFlow.addEdge(key, END, -delta)

# ///////////////////////////////////////////

res = maxFlow.calMaxFlow()
if res != fullFlow:
    print("NO")
    exit(0)
print("YES")
for (u, v) in edge:
    print(maxFlow.getFlowOfEdge(u, v) + edge[(u, v)][0])
