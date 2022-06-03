import sys
from typing import DefaultDict
from collections import defaultdict, deque


Graph = DefaultDict[int, DefaultDict[int, int]]  # 有向带权图,权值为容量


# region Dinic
class Dinic:
    """Dinic 求最大流
    
    时间复杂度:O(V^2*E)
    """

    def __init__(self, graph: Graph, maxNode: int, maxEdge: int) -> None:
        self._graph = graph
        self._maxNode = maxNode
        self._maxEdge = maxEdge

    def calMaxFlow(self, start: int, end: int) -> int:
        def bfs() -> None:
            """bfs建立分层处理出每个点的深度 以及更新当前弧"""
            nonlocal depth, curArc
            depth = [-1] * (self._maxNode + 5)
            depth[start] = 0
            queue = deque([start])
            curArc = [iter(self._reGraph[i].keys()) for i in range(self._maxNode + 5)]

            while queue:
                cur = queue.popleft()
                for child in self._reGraph[cur]:
                    if (depth[child] == -1) and (self._reGraph[cur][child] > 0):
                        depth[child] = depth[cur] + 1
                        queue.append(child)

        def dfsWithCurArc(cur: int, minFlow: int) -> int:
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
                if flow >= minFlow:  # !重要的优化 当前弧优化建立在这个优化的基础之上的
                    break

                # 当前弧优化
                # 先被遍历到的边肯定是已经增广过了（或者已经确定无法继续增广了），那么这条边就可以视为“废边”
                child = next(curArc[cur], None)
                if child is None:
                    break
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

            return flow

        self._updateRedisualGraph()

        res = 0
        depth = [-1] * (self._maxNode + 5)  # 分层,防止环的影响
        curArc = []  # !在dfs搜索中第一次开始搜索的边，也称当前弧，用于优化dfs速度

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

    def _updateRedisualGraph(self) -> None:
        self._reGraph = defaultdict(lambda: defaultdict(int))
        for cur in self._graph:
            for next in self._graph[cur]:
                self._reGraph[cur][next] = self._graph[cur][next]
                self._reGraph[next].setdefault(cur, 0)  # 注意这里 因为_graph可能存在自环


# endregion

# 图中可能存在重边和自环
input = sys.stdin.readline
n, m, start, end = map(int, input().split())
adjMap = defaultdict(lambda: defaultdict(int))

# 从点 u 到点 v 存在一条有向边，容量为 c。
for _ in range(m):
    u, v, c = map(int, input().split())
    adjMap[u][v] += c  # 可能存在重边

maxFlow = Dinic(adjMap, n, m)
print(maxFlow.calMaxFlow(start, end))
