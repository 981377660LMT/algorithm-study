from collections import defaultdict, deque
from functools import lru_cache
from typing import DefaultDict, Generic, Hashable, Protocol, TypeVar

Vertex = TypeVar('Vertex', bound=Hashable)
Graph = DefaultDict[Vertex, DefaultDict[Vertex, int]]  # 有向带权图,权值为容量


class MaxFlowStrategy(Protocol):
    """interface of MaxFlow Strategy"""

    def calMaxFlow(self, start: Vertex, end: Vertex) -> int:
        """求出从start到end的最大流"""
        ...

    def getFlowOfEdge(self, v1: Vertex, v2: Vertex) -> int:
        """获取某条边上的`流量`"""
        ...


###################################################################
class MaxFlow(Generic[Vertex]):
    def __init__(self, graph: Graph, *, strategy: MaxFlowStrategy | None = None) -> None:
        self._graph = graph
        self._strategy = strategy if strategy is not None else self._useDefaultStrategy()

    @lru_cache(None)
    def calMaxFlow(self, start: Vertex, end: Vertex) -> int:
        """求出从start到end的最大流

        如果无法到达end则返回0
        """
        return self._strategy.calMaxFlow(start, end)

    def getFlowOfEdge(self, v1: Vertex, v2: Vertex) -> int:
        """获取某条边上的`流量`"""
        return self._strategy.getFlowOfEdge(v1, v2)

    def getCapacityOfEdge(self, v1: int, v2: int) -> int:
        """获取某条边上的`容量`"""
        return self._graph[v1][v2]

    def switchTo(self, strategy: MaxFlowStrategy) -> None:
        self._strategy = strategy

    def _useDefaultStrategy(self) -> MaxFlowStrategy:
        """默认采用Dinic算法"""
        return Dinic(self._graph)


###################################################################
# region EK
class EK(MaxFlowStrategy):
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

    @lru_cache(None)
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


# endregion

# region Dinic
class Dinic(MaxFlowStrategy):
    """Dinic 求最大流
    
    如果一个流的残量网络里面没有可行流，那么这个流就是最大流
    
    时间复杂度:O(V^2*E)
    """

    def __init__(self, graph: Graph) -> None:
        self._graph = graph
        self._reGraph: Graph = None  # type: ignore

    def calMaxFlow(self, start: Vertex, end: Vertex) -> int:
        def bfs() -> None:
            """bfs建立分层处理出每个点的深度 以及更新当前弧"""
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

        def dfsWithCurArc(cur: Vertex, minFlow: int) -> int:
            """dfs增广 采用当前弧优化

            Args:
                cur (int): 当前点
                minFlow (int): 当前增广路上最小流量

            Returns:
                int: 增广路上的容量

            注意到字典存储的键是插入有序的 因此可以用作记录当前弧
            每次分配完的边就不再dfs了

            每个点的当前弧初始化为head
            每次我们找过某条边(弧)时,修改cur数组,改成该边(弧)的编号,
            下次到达该点时,会直接从cur对应的边开始(也)是说从head到cur中间的那一些边(弧)我们就不走了）。
            """

            if cur == end:
                return minFlow

            flow = 0  # 往后面流的流量可以是多少
            while True:
                if flow >= minFlow:  # !重要的优化
                    break
                try:
                    # 当前弧优化
                    # 先被遍历到的边肯定是已经增广过了（或者已经确定无法继续增广了），那么这条边就可以视为“废边”
                    child = next(curArc[cur])
                    if (depth[child] == depth[cur] + 1) and (self._reGraph[cur][child] > 0):
                        # minFlow - flow 是还可以分配的流量
                        min_ = minFlow - flow
                        if self._reGraph[cur][child] < minFlow - flow:
                            min_ = self._reGraph[cur][child]
                        nextFlow = dfsWithCurArc(child, min_)
                        # !优化：不存在路径 删掉这条边
                        if nextFlow == 0:
                            depth[child] = -1
                        # 找到了增广路 正向边减残量，反向边加残量
                        self._reGraph[cur][child] -= nextFlow
                        self._reGraph[child][cur] += nextFlow
                        flow += nextFlow
                except StopIteration:
                    break
            return flow

        self._updateRedisualGraph()

        res = 0
        depth = defaultdict(lambda: -1, {start: 0})  # 分层,防止环的影响
        curArc = dict()  # !在dfs搜索中第一次开始搜索的边，也称当前弧，用于优化dfs速度

        # 只要有增广路 就在分层图里把增广路全部找出来
        while True:
            bfs()
            if depth[end] != -1:
                while True:
                    delta = dfsWithCurArc(start, int(1e20))  # 多次dfs耗尽流量
                    if delta == 0:
                        break
                    res += delta
            else:
                break
        return res

    def getFlowOfEdge(self, v1: int, v2: int) -> int:
        """获取某条边上的`流量`"""
        return self._graph[v1][v2] - self._reGraph[v1][v2]

    # !depracated 不使用弧优化的版本
    def _calMaxFlow(self, start: Vertex, end: Vertex) -> int:
        def bfs() -> bool:
            """bfs建立分层处理出每个点的深度并判断是否有增广路"""
            nonlocal depth
            depth = defaultdict(lambda: int(1e20), {start: 0})
            visted = set([start])
            queue = deque([start])
            while queue:
                cur = queue.popleft()
                for next in self._reGraph[cur]:
                    if (next not in visted) and (self._reGraph[cur][next] > 0):
                        visted.add(next)
                        depth[next] = depth[cur] + 1
                        queue.append(next)
            return depth[end] != int(1e20)

        def dfs(cur: Vertex, minFlow: int) -> int:
            """dfs增广 不采用当前弧优化

            Args:
                cur (int): 当前点
                minFlow (int): 当前增广路上最小边权

            Returns:
                int: 增广路上的容量
            """

            if cur == end:
                return minFlow
            for next in self._reGraph[cur]:
                if (depth[next] == depth[cur] + 1) and (self._reGraph[cur][next] > 0):
                    nextRes = dfs(next, min(minFlow, self._reGraph[cur][next]))
                    if nextRes > 0:  # 找到了增广路 正向边减残量，反向边加残量
                        self._reGraph[cur][next] -= nextRes
                        self._reGraph[next][cur] += nextRes
                        return nextRes
            return 0  # 没有增广路

        self._updateRedisualGraph()

        res = 0
        depth = defaultdict(lambda: int(1e20), {start: 0})  # 分层,防止环的影响

        # 只要有增广路 就在分层图里把增广路全部找出来
        while True:
            hasAugmentingPath = bfs()
            if hasAugmentingPath:
                while True:
                    delta = dfs(start, int(1e20))  # 多次dfs耗尽流量
                    if delta == 0:
                        break
                    res += delta
            else:
                break
        return res

    def _updateRedisualGraph(self) -> None:
        self._reGraph = defaultdict(lambda: defaultdict(int))
        for cur in self._graph:
            for next in self._graph[cur]:
                self._reGraph[cur][next] = self._graph[cur][next]
                self._reGraph[next].setdefault(cur, 0)  # 注意这里 因为_graph可能存在自环


# endregion

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

    maxFlow = MaxFlow(adjMap)
    assert maxFlow.calMaxFlow(start, end) == 14
    maxFlow.switchTo(EK(adjMap))
    assert maxFlow.calMaxFlow(start, end) == 14

