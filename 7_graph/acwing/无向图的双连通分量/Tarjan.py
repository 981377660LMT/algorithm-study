# golang 实现：
# !算法竞赛进阶指南/GoDS (Go Data Structures)/src/graph/tarjan/tarjan.go

# !deprecated
# !see:无向图的双连通分量/Tarjan


from collections import defaultdict
from typing import DefaultDict, Iterable, List, Mapping, Sequence, Set, Tuple, Union
import sys

Edge = Tuple[int, int]
AdjList = Sequence[Iterable[int]]
AdjMap = Mapping[int, Iterable[int]]
Graph = Union[AdjList, AdjMap]


class Tarjan:
    INF = int(1e20)
    sys.setrecursionlimit(int(1e9))

    @staticmethod
    def getSCC(n: int, graph: Graph) -> Tuple[int, DefaultDict[int, Set[int]], List[int]]:
        """Tarjan求解有向图的强连通分量

        Args:
            n (int): 结点0-n-1
            adjMap (DefaultDict[int, Set[int]]): 图

        Returns:
            Tuple[int, DefaultDict[int, Set[int]], List[int]]: SCC的数量、分组、每个结点对应的SCC编号
        """

        def dfs(cur: int) -> None:
            nonlocal dfsId, SCCId
            if visited[cur]:
                return
            visited[cur] = True

            order[cur] = low[cur] = dfsId
            dfsId += 1
            stack.append(cur)
            inStack[cur] = True

            for next in graph[cur]:
                if not visited[next]:
                    dfs(next)
                    low[cur] = min(low[cur], low[next])
                elif inStack[next]:
                    low[cur] = min(low[cur], order[next])  # 注意这里是order

            if order[cur] == low[cur]:
                while stack:
                    top = stack.pop()
                    inStack[top] = False
                    SCCGroupById[SCCId].add(top)
                    SCCIdByNode[top] = SCCId
                    if top == cur:
                        break
                SCCId += 1

        dfsId = 0
        order, low = [Tarjan.INF] * n, [Tarjan.INF] * n

        visited = [False] * n
        stack = []
        inStack = [False] * n

        SCCId = 0
        SCCGroupById = defaultdict(set)
        SCCIdByNode = [-1] * n

        for cur in range(n):
            if not visited[cur]:
                dfs(cur)

        return SCCId, SCCGroupById, SCCIdByNode

    @staticmethod
    def getCuttingPointAndCuttingEdge(n: int, graph: Graph) -> Tuple[Set[int], Set[Edge]]:
        """Tarjan求解无向图的割点和割边(桥)

        Args:
            n (int): 结点0-n-1
            adjMap (DefaultDict[int, Set[int]]): 图

        Returns:
            Tuple[Set[int], Set[Edge]]: 割点、桥

        - 边对 (u,v) 中 u < v
        """

        def dfs(cur: int, parent: int) -> None:
            nonlocal dfsId
            if visited[cur]:
                return
            visited[cur] = True

            order[cur] = low[cur] = dfsId
            dfsId += 1

            dfsChild = 0
            for next in graph[cur]:
                if next == parent:
                    continue
                if not visited[next]:
                    dfsChild += 1
                    dfs(next, cur)
                    low[cur] = min(low[cur], low[next])
                    if low[next] > order[cur]:
                        cuttingEdge.add(tuple(sorted([cur, next])))
                    if parent != -1 and low[next] >= order[cur]:
                        cuttingPoint.add(cur)
                    elif parent == -1 and dfsChild > 1:  # 出发点没有祖先order肯定小啊，所以特判一下
                        cuttingPoint.add(cur)
                else:
                    low[cur] = min(low[cur], order[next])  # 注意这里是order

        dfsId = 0
        order, low = [Tarjan.INF] * n, [Tarjan.INF] * n
        visited = [False] * n

        cuttingPoint = set()
        cuttingEdge = set()

        for i in range(n):
            if not visited[i]:
                dfs(i, -1)

        return cuttingPoint, cuttingEdge

    @staticmethod
    def getVBCC(n: int, graph: Graph) -> Tuple[int, DefaultDict[int, Set[int]], List[Set[int]]]:
        """Tarjan求解无向图的点双联通分量

        Args:
            n (int): 结点0-n-1
            adjMap (DefaultDict[int, Set[int]]): 图

        Returns:
            Tuple[int, DefaultDict[int, Set[int]], List[Set[int]]]: VBCC的数量、分组、每个结点对应的VBCC编号

        - 我们将深搜时遇到的所有边加入到栈里面，
        当找到一个割点的时候，
        就将这个割点往下走到的所有边弹出，
        而这些边所连接的点就是一个点双了

        - 两个点和一条边构成的图也称为(V)BCC,因为两个点均不为割点

        - VBCC编号多余1个的都是割点
        """

        def dfs(cur: int, parent: int) -> None:
            nonlocal dfsId, VBCCId
            if visited[cur]:
                return
            visited[cur] = True

            order[cur] = low[cur] = dfsId
            dfsId += 1

            dfsChild = 0
            for next in graph[cur]:
                if next == parent:
                    continue

                if not visited[next]:
                    dfsChild += 1
                    stack.append((cur, next))
                    dfs(next, cur)
                    low[cur] = min(low[cur], low[next])

                    # 遇到了割点(根和非根两种)
                    if (parent == -1 and dfsChild > 1) or (
                        parent != -1 and low[next] >= order[cur]
                    ):
                        while stack:
                            top = stack.pop()
                            VBCCGroupById[VBCCId].add(top[0])
                            VBCCGroupById[VBCCId].add(top[1])
                            VBCCIdByNode[top[0]].add(VBCCId)
                            VBCCIdByNode[top[1]].add(VBCCId)
                            if top == (cur, next):
                                break
                        VBCCId += 1

                elif low[cur] > order[next]:
                    low[cur] = min(low[cur], order[next])
                    stack.append((cur, next))

        dfsId = 0
        order, low = [Tarjan.INF] * n, [Tarjan.INF] * n

        visited = [False] * n
        stack = []

        VBCCId = 0  # 点双个数
        VBCCGroupById = defaultdict(set)  # 每个点双包含哪些点
        VBCCIdByNode = [set() for _ in range(n)]  # 每个点属于哪一(几)个点双

        for cur in range(n):
            if not visited[cur]:
                dfs(cur, -1)

            if stack:
                while stack:
                    top = stack.pop()
                    VBCCGroupById[VBCCId].add(top[0])
                    VBCCGroupById[VBCCId].add(top[1])
                    VBCCIdByNode[top[0]].add(VBCCId)
                    VBCCIdByNode[top[1]].add(VBCCId)
                VBCCId += 1

        return VBCCId, VBCCGroupById, VBCCIdByNode

    @staticmethod
    def getEBCC(
        n: int, graph: Graph
    ) -> Tuple[int, DefaultDict[int, Set[Edge]], DefaultDict[Edge, int]]:
        """Tarjan求解无向图的边双联通分量

        Args:
            n (int): 结点0-n-1
            adjMap (DefaultDict[int, Set[int]]): 图

        Returns:
            Tuple[int, DefaultDict[int, Set[Edge]], DefaultDict[Edge, int]]]: EBCC的数量、分组、每条边对应的EBCC编号

        - 边对 (u,v) 中 u < v

        - 实现思路：
          - 将所有的桥删掉剩下的都是边连通分量了(其实可以用并查集做)
          - 处理出割边,再对整个无向图进行一次DFS,对于节点cur的出边(cur,next),如果它是割边,则跳过这条边不沿着它往下走
        """

        def dfs(cur: int, parent: int) -> None:
            nonlocal EBCCId
            if visited[cur]:
                return
            visited[cur] = True

            for next in graph[cur]:
                if next == parent:
                    continue

                edge = tuple(sorted([cur, next]))
                if edge in cuttingEdges:
                    continue

                EBCCGroupById[EBCCId].add(edge)
                EBCCIdByEdge[edge] = EBCCId
                dfs(next, cur)

        _, cuttingEdges = Tarjan.getCuttingPointAndCuttingEdge(n, graph)

        visited = [False] * n

        EBCCId = 0  # 边双个数
        EBCCGroupById = defaultdict(set)  # 每个边双包含哪些边
        EBCCIdByEdge = defaultdict(int)  # 每条边属于哪一个边双

        for cur in range(n):
            if not visited[cur]:
                dfs(cur, -1)
                EBCCId += 1

        for edge in cuttingEdges:
            EBCCGroupById[EBCCId].add(edge)
            EBCCIdByEdge[edge] = EBCCId
            EBCCId += 1

        return EBCCId, EBCCGroupById, EBCCIdByEdge


def toDAG(n: int, graph: Graph):
    """有向图缩点成DAG"""
    sccCount, sccGroupById, sccIdByNode = Tarjan.getSCC(n, graph)
    newAdjMap = defaultdict(set)
    deg = defaultdict(int)
    visitedPair = set()
    for cur in range(n):
        for next in graph[cur]:
            g1, g2 = sccIdByNode[cur], sccIdByNode[next]
            if g1 == g2 or (g1, g2) in visitedPair:
                continue
            visitedPair.add((g1, g2))
            newAdjMap[g1].add(g2)
            deg[g2] += 1

    return sccCount, newAdjMap, deg, sccIdByNode, sccGroupById


if __name__ == "__main__":
    # 无向图割点和桥
    adjMap1 = defaultdict(set)
    edges = [[0, 1], [0, 2], [1, 2], [2, 5], [2, 4], [3, 4], [3, 5], [4, 5]]
    for u, v in edges:
        adjMap1[u].add(v)
        adjMap1[v].add(u)
    assert Tarjan.getCuttingPointAndCuttingEdge(6, adjMap1) == ({2}, set())

    # 无向图VBCC
    adjMap2 = defaultdict(set)
    edges = [[0, 1], [0, 2], [1, 2], [2, 3], [2, 4], [3, 4]]
    for u, v in edges:
        adjMap2[u].add(v)
        adjMap2[v].add(u)
    assert Tarjan.getVBCC(5, adjMap2)[2] == [{1}, {1}, {0, 1}, {0}, {0}]

    # 无向图EBCC
    adjMap2 = defaultdict(set)
    edges = [[0, 1], [0, 2], [1, 2], [2, 3], [3, 4], [3, 5], [4, 5]]
    for u, v in edges:
        adjMap2[v].add(u)
        adjMap2[u].add(v)
    assert list(Tarjan.getEBCC(6, adjMap2)[1].values()) == [
        {(0, 1), (0, 2), (1, 2)},
        {(4, 5), (3, 4), (3, 5)},
        {(2, 3)},
    ]

    # 有向图SCC
    adjMap2 = defaultdict(set)
    edges = [[1, 0], [0, 2], [2, 1], [0, 3], [3, 4]]
    for u, v in edges:
        adjMap2[u].add(v)
    assert Tarjan.getSCC(5, adjMap2)[2] == [2, 2, 2, 1, 0]

    edges = [(0, 1), (1, 2), (0, 5), (1, 4), (2, 3), (3, 4), (4, 5)]
    adjMap = defaultdict(set)
    for u, v in edges:
        adjMap[u].add(v)
        adjMap[v].add(u)
    print(Tarjan.getEBCC(6, adjMap)[1])
