from collections import defaultdict
from typing import List
from typing import DefaultDict, Generic, Hashable, List, Tuple, TypeVar
from dataclasses import dataclass
from collections import defaultdict, deque


V = TypeVar('V', bound=Hashable)


@dataclass(slots=True)
class Edge(Generic[V]):
    fromV: V
    toV: V
    cap: int
    cost: int
    flow: int


# !图中可能存在自环边、平行边
class MinCostMaxFlow(Generic[V]):
    """最小费用流的连续最短路算法复杂度为流量*最短路算法复杂度"""

    def __init__(self, start: V, end: V):
        self._start: V = start
        self._end: V = end
        self._edges: List['Edge'[V]] = []
        self._reGraph: DefaultDict[V, List[int]] = defaultdict(list)  # 残量图存储的是边的下标

    def addEdge(self, fromV: V, toV: V, cap: int, cost: int) -> None:
        """原边索引为i 反向边索引为i^1"""
        self._edges.append(Edge(fromV, toV, cap, cost, 0))
        self._edges.append(Edge(toV, fromV, 0, -cost, 0))
        len_ = len(self._edges)
        self._reGraph[fromV].append(len_ - 2)
        self._reGraph[toV].append(len_ - 1)

    def work(self) -> Tuple[int, int]:
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


# 注意网络流的流量守恒
# !如果某个点流入流出的流量不一致，那么就需要从源点加一条虚拟边来平衡流量
# 虚拟边流量为abs(入度-出度) 费用为0
# !出度大于入度的点看做源点，每个入度大于出度的点看做汇点，然后再添加虚拟源点和汇点，是多源多汇的
# !核心是费用流，多源多汇最后也是转换成单源单汇来解的


class Solution:
    def maximumRequests(self, n: int, requests: List[List[int]]) -> int:
        """请你从原请求列表中选出若干个请求，使得它们是一个可行的请求列表，并返回所有可行列表中最大请求数目。"""
        START, END = n + 1, n + 2
        flowDiff = defaultdict(int)
        mcmf = MinCostMaxFlow(START, END)
        for u, v in requests:
            mcmf.addEdge(u, v, 1, 1)
            flowDiff[u] -= 1
            flowDiff[v] += 1

        for key, count in flowDiff.items():
            if count > 0:
                mcmf.addEdge(key, END, count, 0)
            elif count < 0:
                mcmf.addEdge(START, key, -count, 0)

        # !因为虚拟源点和汇点是需要去掉的(这两个点流量不守恒) 去掉之后剩下的部分就成了题中需要求的循环自洽的整体
        # !剩下要最大 那么就要在原来图中删去经过边数最少的流 即原图的最小费用流
        return len(requests) - mcmf.work()[1]


print(Solution().maximumRequests(n=5, requests=[[0, 1], [1, 0], [0, 1], [1, 2], [2, 0], [3, 4]]))
print(
    Solution().maximumRequests(
        n=3, requests=[[1, 2], [1, 2], [2, 2], [0, 2], [2, 1], [1, 1], [1, 2]]
    )
)
# 4
