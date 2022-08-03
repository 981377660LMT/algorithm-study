# 1≤N≤500,
# 1≤M≤5000,

# n个点m条边的网络，给定源点S和汇点T，求如果有这样边：
# !只给其扩大容量之后整个流网络的最大流能够变大，
# 对于这样的边我们称之为`关键边`。求这样的边的个数
# https://www.acwing.com/solution/content/73355/


# 1.先对原图跑一次最大流，求出原图的残量网络
# 2.对于边(u,v)，判断从s->u和从v->t是否都存在流量大于0的路
# 3.如果某条边上的流量是满的 并且左右端点可以分别到达源点、汇点 就是关键边


from collections import defaultdict, deque
from typing import DefaultDict, Set


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


# 城市编号从 0 到 N−1。
# 生产日常商品的城市为 0 号城市，首都为 N−1 号城市。
def dfs(
    cur: int, visited: Set[int], graph: DefaultDict[int, DefaultDict[int, int]], isRevered: bool
) -> None:
    visited.add(cur)
    for next in graph[cur]:
        remain = (
            maxFlow.getRemainOfEdge(cur, next)
            if not isRevered
            else maxFlow.getRemainOfEdge(next, cur)
        )
        if next not in visited and remain != 0:  # 可以继续走
            dfs(next, visited, graph, isRevered)


n, m = map(int, input().split())
adjMap = defaultdict(lambda: defaultdict(int))
rAdjMap = defaultdict(lambda: defaultdict(int))
maxFlow = MaxFlow(0, n - 1)
for _ in range(m):
    u, v, w = map(int, input().split())
    adjMap[u][v] += w
    rAdjMap[v][u] += w
    maxFlow.addEdge(u, v, w)

maxFlow.calMaxFlow()

# !如果某条边上的流量是满的 就可以作为候选
visited1, visited2 = set(), set()  # 可以到源点的点 和 可以到汇点的点
dfs(0, visited1, adjMap, False)
dfs(n - 1, visited2, rAdjMap, True)

res = 0
for cur in adjMap:
    for next in adjMap[cur]:
        if (
            (cur in visited1)
            and (next in visited2)
            and maxFlow.getRemainOfEdge(cur, next) == 0  # 满流
        ):
            res += 1

print(res)
