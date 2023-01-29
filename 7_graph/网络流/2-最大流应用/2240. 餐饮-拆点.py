# 约翰共有 N 头奶牛，其中第 i 头奶牛有 Fi 种喜欢的食物以及 Di 种喜欢的饮料。
# 约翰需要给每头奶牛分配一种食物和一种饮料，并使得有吃有喝的奶牛数量尽可能大。
# 每种食物或饮料都只有一份，所以只能分配给一头奶牛食用
# 输出一个整数，表示能够有吃有喝的奶牛的最大数量。
"""
把牛拆成两个点，两个点连接容量是1的边，限制每一头牛只能用一次
"""
# !边有限制 容量设在边上 点有限制 点拆成in 和 out


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


n, food, drink = map(int, input().split())
foodLike = defaultdict(set)
drinkLike = defaultdict(set)
for cowId in range(n):
    f, d, *rest = map(int, input().split())
    for fid in rest[:f]:
        foodLike[cowId].add(fid)
    for did in rest[f:]:
        drinkLike[cowId].add(did)

START, END = -1, -2
OFFSET1 = int(1e5)
OFFSET2 = int(2e5)
OFFSET3 = int(3e5)
OFFSET4 = int(4e5)
maxFlow = MaxFlow(START, END)

# 虚拟源点到food
for fid in range(1, food + 1):
    maxFlow.addEdge(START, fid + OFFSET1, 1, cover=True)
# food到牛的in
for cowId in foodLike:
    for fid in foodLike[cowId]:
        maxFlow.addEdge(fid + OFFSET1, cowId + OFFSET2, 1, cover=True)
# 牛的in到牛的out
for cowId in range(n):
    maxFlow.addEdge(cowId + OFFSET2, cowId + OFFSET3, 1, cover=True)
# 牛的out到drink
for cowId in drinkLike:
    for did in drinkLike[cowId]:
        maxFlow.addEdge(cowId + OFFSET3, did + OFFSET4, 1, cover=True)
# drink到虚拟汇点
for did in range(1, drink + 1):
    maxFlow.addEdge(did + OFFSET4, END, 1, cover=True)

print(maxFlow.calMaxFlow())
