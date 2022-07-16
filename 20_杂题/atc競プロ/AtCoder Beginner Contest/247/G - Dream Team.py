# https://www.zhihu.com/search?type=content&q=AtCoder%20Beginner%20Contest%20247
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


from typing import DefaultDict, Generic, Hashable, List, TypeVar
from collections import defaultdict, deque


V = TypeVar("V", bound=Hashable)


class Edge(Generic[V]):
    __slots__ = ("fromV", "toV", "cap", "cost", "flow")

    def __init__(self, fromV: V, toV: V, cap: int, cost: int, flow: int) -> None:
        self.fromV = fromV
        self.toV = toV
        self.cap = cap
        self.cost = cost
        self.flow = flow


# https://github.com/981377660LMT
class MinCostMaxFlow(Generic[V]):
    def __init__(self, start: V, end: V):
        self._start: V = start
        self._end: V = end
        self._edges: List["Edge"[V]] = []
        self._reGraph: DefaultDict[V, List[int]] = defaultdict(list)  # 残量图存储的是边的下标

    def addEdge(self, fromV: V, toV: V, cap: int, cost: int) -> None:
        """原边索引为i 反向边索引为i^1"""
        self._edges.append(Edge(fromV, toV, cap, cost, 0))
        self._edges.append(Edge(toV, fromV, 0, -cost, 0))
        len_ = len(self._edges)
        self._reGraph[fromV].append(len_ - 2)
        self._reGraph[toV].append(len_ - 1)

    def work(self) -> List[int]:
        """
        Returns:
            List[int]: 输出的最优匹配数(增广次数)为1,2,…i个匹配时的答案
        """

        def spfa() -> int:
            """spfa沿着最短路寻找增广路径  有负cost的边不能用dijkstra"""
            nonlocal dist
            dist = defaultdict(lambda: int(1e18), {self._start: 0})
            inQueue = defaultdict(lambda: False)
            queue = deque([self._start])
            inFlow = defaultdict(int, {self._start: int(1e18)})  # 到每条边上的流量
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

        dist = defaultdict(lambda: int(1e18), {self._start: 0})
        cost = 0
        res = []
        while True:
            delta = spfa()  # 一次spfa只会找到一条费用最小的增广流
            if delta == 0:
                break
            cost += delta * dist[self._end]
            res.append(cost)
        return res


import sys
import os

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)


def main() -> None:
    n = int(input())
    START, END, OFFSET = -1, -2, 1000
    mcmf = MinCostMaxFlow(START, END)
    for i in range(n):
        mcmf.addEdge(START, i, 1, 0)
        mcmf.addEdge(i + OFFSET, END, 1, 0)
    for _ in range(n):
        u, v, w = map(int, input().split())
        u, v = u - 1, v - 1
        mcmf.addEdge(u, v + OFFSET, 1, -w)  # 要求最大费用流 所以负数费用

    res = mcmf.work()
    print(len(res))
    for num in res:
        print(-num)


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()

# TODO
