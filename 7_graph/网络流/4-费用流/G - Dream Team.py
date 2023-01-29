# https://zhuanlan.zhihu.com/p/496282947
# 题意
# 有n个人来自不同的学校ai,擅长不同的学科bi,每个人有一个能力值ci
# 要求组建一支i个人的梦之队最大化队员的能力值,并且满足队伍中所有人来自的学校和擅长的学科都不同.
# n<=3e4
# !ai,bi<=150 (暗示作为顶点的数据量)
# ci<=1e9

# 分析
# 把学校和学科看作点,把一个人看成匹配边,能力值看作边权,其实就是求有i条匹配边的最优匹配.可以用费用流解决.
# 此外题目要求输出匹配数为1,2,…k个匹配时的最优匹配.
# !在spfa费用流算法中一次spfa只会找到一条费用最小的增广流,
# !所以每次增广之后得到的费用就对应匹配数为1,2,…k个匹配时的答案.

# !能力值对应 流的费用
# !队伍里每个人学校和学科都不同:学校学科为虚拟源汇点，容量为1，这样就不会取到重复的学生了

import sys


sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)
INF = int(1e18)

from collections import deque
from typing import List


class Edge:
    __slots__ = ("fromV", "toV", "cap", "cost", "flow")

    def __init__(self, fromV: int, toV: int, cap: int, cost: int, flow: int) -> None:
        self.fromV = fromV
        self.toV = toV
        self.cap = cap
        self.cost = cost
        self.flow = flow


class MinCostMaxFlow:
    """最小费用流的连续最短路算法复杂度为流量*最短路算法复杂度"""

    __slots__ = ("_n", "_start", "_end", "_edges", "_reGraph", "_dist", "_visited", "_curEdges")

    def __init__(self, n: int, start: int, end: int):
        """
        Args:
            n (int): 包含虚拟点在内的总点数
            start (int): (虚拟)源点
            end (int): (虚拟)汇点
        """
        assert 0 <= start < n and 0 <= end < n
        self._n = n
        self._start = start
        self._end = end
        self._edges: List["Edge"] = []
        self._reGraph: List[List[int]] = [[] for _ in range(n + 10)]  # 残量图存储的是边的下标

        self._dist = [INF] * (n + 10)
        self._visited = [False] * (n + 10)
        self._curEdges = [0] * (n + 10)

    def addEdge(self, fromV: int, toV: int, cap: int, cost: int) -> None:
        """原边索引为i 反向边索引为i^1"""
        self._edges.append(Edge(fromV, toV, cap, cost, 0))
        self._edges.append(Edge(toV, fromV, 0, -cost, 0))
        len_ = len(self._edges)
        self._reGraph[fromV].append(len_ - 2)
        self._reGraph[toV].append(len_ - 1)

    def work(self) -> List[int]:
        """
        Returns:
             List[int]: 每个k限流(匹配了k条边)的最大花费
        """
        res = []
        minCost = 0
        while self._spfa():
            # !限定流量为1 每次dfs只找到一条费用最小的增广流
            flow = self._dfs(self._start, self._end, 1)
            minCost += flow * self._dist[self._end]
            res.append(minCost)
        return res

    def _spfa(self) -> bool:
        """spfa沿着最短路寻找增广路径  有负cost的边不能用dijkstra"""
        n, start, end, edges, reGraph, visited = (
            self._n,
            self._start,
            self._end,
            self._edges,
            self._reGraph,
            self._visited,
        )

        self._curEdges = [0] * n
        self._dist = dist = [INF] * n
        dist[start] = 0
        queue = deque([start])

        while queue:
            cur = queue.popleft()
            visited[cur] = False
            for edgeIndex in reGraph[cur]:
                edge = edges[edgeIndex]
                cost, remain, next = edge.cost, edge.cap - edge.flow, edge.toV
                if remain > 0 and dist[cur] + cost < dist[next]:
                    dist[next] = dist[cur] + cost
                    if not visited[next]:
                        visited[next] = True
                        if queue and dist[queue[0]] > dist[next]:
                            queue.appendleft(next)
                        else:
                            queue.append(next)

        return dist[end] != INF

    def _dfs(self, cur: int, end: int, flow: int) -> int:
        if cur == end:
            return flow

        visited, reGraph, curEdges, edges, dist = (
            self._visited,
            self._reGraph,
            self._curEdges,
            self._edges,
            self._dist,
        )

        visited[cur] = True
        res = flow
        index = curEdges[cur]
        while res and index < len(reGraph[cur]):
            edgeIndex = reGraph[cur][index]
            next, remain = edges[edgeIndex].toV, edges[edgeIndex].cap - edges[edgeIndex].flow
            if remain > 0 and not visited[next] and dist[next] == dist[cur] + edges[edgeIndex].cost:
                delta = self._dfs(next, end, remain if remain < res else res)
                res -= delta
                edges[edgeIndex].flow += delta
                edges[edgeIndex ^ 1].flow -= delta
            curEdges[cur] += 1
            index = curEdges[cur]

        visited[cur] = False
        return flow - res


#####################################

n = int(input())
v = 150
START, END, OFFSET = 2 * v, 2 * v + 1, v
mcmf = MinCostMaxFlow(2 * v + 2, START, END)
for i in range(v):
    mcmf.addEdge(i + OFFSET, END, 1, 0)
    mcmf.addEdge(START, i, 1, 0)

for _ in range(n):
    u, v, w = map(int, input().split())
    u, v = u - 1, v - 1
    mcmf.addEdge(u, v + OFFSET, 1, -w)  # !要求最大费用流

res = mcmf.work()
print(len(res))
for i, cost in enumerate(res, 1):
    print(-cost)
