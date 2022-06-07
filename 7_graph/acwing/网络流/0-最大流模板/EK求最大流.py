# 给定一个包含 n 个点 m 条边的有向图，并给定每条边的容量，边的容量非负。
# 图中可能存在重边和自环。求从点 S 到点 T 的最大流。

# 点的编号从 1 到 n。
# 输出点 S 到点 T 的最大流。
# 如果从点 S 无法到达点 T 则输出 0。

# !EK:
# 2≤n≤1e3,
# 1≤m≤1e4,
# !DINIC:
# 2≤n≤1e4,
# 1≤m≤1e5,
# region EK

from collections import defaultdict, deque
from typing import DefaultDict, Hashable, TypeVar

Vertex = TypeVar('Vertex', bound=Hashable)
Graph = DefaultDict[Vertex, DefaultDict[Vertex, int]]  # 有向带权图,权值为容量


class EK:
    """EK 求最大流
    
    如果一个流的残量网络里面没有可行流，那么这个流就是最大流

    时间复杂度:O(V*E^2)
    1. 找增广路
    2. 更新残量网络增广路上的流量
    3. 重复执行1、2直到网络里没有增广路
    """

    def __init__(self, graph: Graph) -> None:
        self._graph = graph  # 容量原图
        self._reGraph: Graph = None  # type: ignore

    def calMaxFlow(self, start: Vertex, end: Vertex) -> int:
        def bfs() -> int:
            """bfs在残量网络上寻找增广路径"""
            visited = set([start])
            queue = deque([(start, int(1e20))])
            pre = {start: start}
            resDelta = 0
            while queue:
                cur, delta = queue.popleft()
                if cur == end:
                    resDelta = delta
                    break
                for next in self._reGraph[cur]:
                    if (next not in visited) and (self._reGraph[cur][next] > 0):
                        visited.add(next)
                        queue.append((next, min(delta, self._reGraph[cur][next])))
                        pre[next] = cur

            if resDelta > 0:  # 找到可行流
                cur = end
                while cur != start:
                    parent = pre[cur]
                    self._reGraph[parent][cur] -= resDelta
                    self._reGraph[cur][parent] += resDelta
                    cur = parent
            return resDelta

        self._updateRedisualGraph()
        res = 0
        while True:
            delta = bfs()
            if delta == 0:
                break
            res += delta
        return res

    def getFlowOfEdge(self, v1: int, v2: int) -> int:
        """获取某条边上的`流量`"""
        return self._graph[v1][v2] - self._reGraph[v1][v2]

    def getCapacityOfEdge(self, v1: int, v2: int) -> int:
        """获取某条边上的`容量`"""
        return self._graph[v1][v2]

    def _updateRedisualGraph(self) -> None:
        self._reGraph = defaultdict(lambda: defaultdict(int))
        for cur in self._graph:
            for next in self._graph[cur]:
                self._reGraph[cur][next] = self._graph[cur][next]
                self._reGraph[next].setdefault(cur, 0)  # 注意这里 因为可能存在自环


###################################################################
if __name__ == '__main__':
    # # 图中可能存在重边和自环
    n, m, start, end = [7, 14, 1, 7]
    adjMap = defaultdict(lambda: defaultdict(int))
    inputs = [
        [1, 2, 5],
        [1, 3, 6],
        [1, 4, 5],
        [2, 3, 2],
        [2, 5, 3],
        [3, 2, 2],
        [3, 4, 3],
        [3, 5, 3],
        [3, 6, 7],
        [4, 6, 5],
        [5, 6, 1],
        [6, 5, 1],
        [5, 7, 8],
        [6, 7, 7],
    ]
    # 从点 u 到点 v 存在一条有向边，容量为 c。
    for i in range(m):
        u, v, c = inputs[i]
        adjMap[u][v] += c  # 可能存在重边

    maxFlow = EK(adjMap)
    assert maxFlow.calMaxFlow(start, end) == 14
