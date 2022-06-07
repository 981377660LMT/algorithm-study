# W 公司有 m 个仓库和 n 个零售商店。
# 第 i 个仓库有 ai 个单位的货物；第 j 个零售商店需要 bj 个单位的货物。
# 货物供需平衡，即∑i=1mai=∑j=1nbj。
# !从第 i 个仓库运送每单位货物到第 j 个零售商店的费用为 cij。
# 试设计一个将仓库中所有货物运送到零售商店的运输方案。
# 对于给定的 m 个仓库和 n 个零售商店间运送货物的费用，计算最优运输方案和最差运输方案。
# https://www.acwing.com/problem/content/description/2194/


from typing import Generic, Hashable, Tuple, TypeVar
from collections import defaultdict, deque


V = TypeVar('V', bound=Hashable)


class Edge(Generic[V]):
    __slots__ = ('fromV', 'toV', 'cap', 'cost', 'flow')

    def __init__(self, fromV: V, toV: V, cap: int, cost: int, flow: int) -> None:
        self.fromV = fromV
        self.toV = toV
        self.cap = cap
        self.cost = cost
        self.flow = flow


class MinCostMaxFlow(Generic[V]):
    """最小费用流的连续最短路算法复杂度为流量*最短路算法复杂度"""

    def __init__(self, start: V, end: V):
        self._start = start
        self._end = end
        self._edges = []
        self._reGraph = defaultdict(list)  # 残量图存储的是边的下标

    def addEdge(self, fromV: V, toV: V, cap: int, cost: int) -> None:
        """原边索引为i 反向边索引为i^1
        
        Args:
            fromV (V): 起点
            toV (V): 终点
            cap (int): 容量，仓库提供货物的数量
            cost (int): 每单位货物的费用
        """
        self._edges.append(Edge(fromV, toV, cap, cost, 0))
        self._edges.append(Edge(toV, fromV, 0, -cost, 0))
        len_ = len(self._edges)
        self._reGraph[fromV].append(len_ - 2)
        self._reGraph[toV].append(len_ - 1)

    def work(self) -> Tuple[int, int]:
        """
        Returns:
            Tuple[int, int]: [最大流,最小费用]
        """

        def spfa() -> int:
            """spfa沿着最短路寻找增广路径  有负cost的边不能用dijkstra"""
            nonlocal dist
            dist = defaultdict(lambda: int(1e20), {self._start: 0})
            inQueue = defaultdict(lambda: False)
            queue = deque([self._start])
            inFlow = defaultdict(int, {self._start: int(1e20)})  # 到每条边上的流量
            pre = defaultdict(lambda: -1)

            while queue:
                cur = queue.popleft()
                inQueue[cur] = False
                for edgeIndex in self._reGraph[cur]:
                    edge = self._edges[edgeIndex]
                    cost, flow, cap, next = edge.cost, edge.flow, edge.cap, edge.toV
                    if dist[cur] + cost < dist[next] and (cap - flow) > 0:
                        dist[next] = dist[cur] + cost
                        pre[next] = edgeIndex
                        inFlow[next] = min(inFlow[cur], cap - flow)
                        if not inQueue[next]:
                            inQueue[next] = True
                            queue.append(next)

            resDelta = inFlow[self._end]
            if resDelta > 0:  # 找到可行流
                cur = self._end
                while cur != self._start:
                    preEdgeIndex = pre[cur]
                    self._edges[preEdgeIndex].flow += resDelta
                    self._edges[preEdgeIndex ^ 1].flow -= resDelta
                    cur = self._edges[preEdgeIndex].fromV
            return resDelta

        dist = defaultdict(lambda: int(1e20), {self._start: 0})
        flow = cost = 0
        while True:
            delta = spfa()
            if delta == 0:
                break
            flow += delta
            cost += delta * dist[self._end]
        return flow, cost


# 仓库数和零售商店数
# 1≤m≤100,
# 1≤n≤50,
m, n = map(int, input().split())

# 表示第 i 个仓库有 ai 个单位的货物。
stores = list(map(int, input().split()))

# 第 j 个零售商店需要 bj 个单位的货物。
needs = list(map(int, input().split()))

# 表示从第 i 个仓库运送每单位货物到第 j 个零售商店的费用
dist = defaultdict(lambda: defaultdict(lambda: int(1e20)))
for i in range(m):
    nums = list(map(int, input().split()))
    for j, num in enumerate(nums):
        dist[i][j] = num


START, END, OFFSET = -1, -2, int(1e4)

# 最小费用
mcmf1 = MinCostMaxFlow(START, END)
for i in range(m):
    mcmf1.addEdge(START, i, stores[i], 0)  # 虚拟源点提货物
for i in range(n):
    mcmf1.addEdge(i + OFFSET, END, needs[i], 0)  # 虚拟汇点接受货物
for i in dist:
    for j in dist[i]:
        mcmf1.addEdge(i, j + OFFSET, stores[i], dist[i][j])  # 仓库转移虚拟源点的货物
print(mcmf1.work()[1])

# 最大费用
mcmf2 = MinCostMaxFlow(START, END)
for i in range(m):
    mcmf2.addEdge(START, i, stores[i], 0)
for i in range(n):
    mcmf2.addEdge(i + OFFSET, END, needs[i], 0)
for i in dist:
    for j in dist[i]:
        mcmf2.addEdge(i, j + OFFSET, stores[i], -dist[i][j])
print(-mcmf2.work()[1])
