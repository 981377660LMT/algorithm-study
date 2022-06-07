# 约翰共有 N 头奶牛，其中第 i 头奶牛有 Fi 种喜欢的食物以及 Di 种喜欢的饮料。
# 约翰需要给每头奶牛分配一种食物和一种饮料，并使得有吃有喝的奶牛数量尽可能大。
# 每种食物或饮料都只有一份，所以只能分配给一头奶牛食用
# 输出一个整数，表示能够有吃有喝的奶牛的最大数量。
'''
把牛拆成两个点，两个点连接容量是1的边，限制每一头牛只能用一次
'''
# !边有限制 容量设在边上 点有限制 点拆成in 和 out


from collections import defaultdict
from typing import DefaultDict, List, Set
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

                child = next(curArc[cur], None)
                if child is not None:
                    if (depth[child] == depth[cur] + 1) and (self._reGraph[cur][child] > 0):
                        nextFlow = dfsWithCurArc(
                            child, min(minFlow - flow, self._reGraph[cur][child])
                        )
                        if nextFlow == 0:
                            depth[child] = -1
                        self._reGraph[cur][child] -= nextFlow
                        self._reGraph[child][cur] += nextFlow
                        flow += nextFlow
                else:
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

    def getRemainOfEdge(self, v1: int, v2: int) -> int:
        assert v1 in self._graph and v2 in self._graph[v1]
        return self._reGraph[v1][v2]

    def _updateRedisualGraph(self) -> None:
        self._reGraph = defaultdict(lambda: defaultdict(int))
        for cur in self._graph:
            for next in self._graph[cur]:
                self._reGraph[cur][next] = self._graph[cur][next]
                self._reGraph[next].setdefault(cur, 0)


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
adjMap = defaultdict(lambda: defaultdict(int))

# 虚拟源点到food
for fid in range(1, food + 1):
    adjMap[START][fid + OFFSET1] = 1
# food到牛的in
for cowId in foodLike:
    for fid in foodLike[cowId]:
        adjMap[fid + OFFSET1][cowId + OFFSET2] = 1
# 牛的in到牛的out
for cowId in range(n):
    adjMap[cowId + OFFSET2][cowId + OFFSET3] = 1
# 牛的out到drink
for cowId in drinkLike:
    for did in drinkLike[cowId]:
        adjMap[cowId + OFFSET3][did + OFFSET4] = 1
# drink到虚拟汇点
for did in range(1, drink + 1):
    adjMap[did + OFFSET4][END] = 1

maxFlow = Dinic(adjMap)
print(maxFlow.calMaxFlow(START, END))
