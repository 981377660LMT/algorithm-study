from collections import defaultdict
from typing import DefaultDict, List, Set, Tuple
import sys

sys.setrecursionlimit(int(1e9))


class Tarjan:
    INF = int(1e20)

    @staticmethod
    def getSCC(
        n: int, adjMap: DefaultDict[int, Set[Tuple[int, int]]]
    ) -> Tuple[int, List[List[int]], List[int]]:
        """Tarjan求解有向图的强连通分量

        Args:
            n (int): 结点0-n-1
            adjMap (DefaultDict[int, Set[int]]): 图

        Returns:
            Tuple[int, List[List[int]], List[int]]: SCC的数量、分组、每个结点对应的SCC编号
        """

        def dfs(cur: int) -> None:
            if visited[cur]:
                return
            visited[cur] = True

            nonlocal dfsId, SCCId
            order[cur] = low[cur] = dfsId
            dfsId += 1
            stack.append(cur)
            inStack[cur] = True

            for next, _ in adjMap[cur]:
                if not visited[next]:
                    dfs(next)
                    low[cur] = min(low[cur], low[next])
                elif inStack[next]:
                    low[cur] = min(low[cur], low[next])
            if order[cur] == low[cur]:
                while stack:
                    top = stack.pop()
                    inStack[top] = False
                    SCCGroupById[SCCId].append(top)
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
        SCCGroupById = [[] for _ in range(n)]
        SCCIdByNode = [-1] * n

        for cur in range(n):
            if not visited[cur]:
                dfs(cur)

        return SCCId, SCCGroupById, SCCIdByNode

    @staticmethod
    def getCuttingPointAndCuttingEdge(
        n: int, adjMap: DefaultDict[int, Set[int]]
    ) -> Tuple[List[int], List[Tuple[int, int]]]:
        """Tarjan求解无向图的割点和割边(桥)

        Args:
            n (int): 结点0-n-1
            adjMap (DefaultDict[int, Set[int]]): 图

        Returns:
            Tuple[List[int], List[Tuple[int, int]]]: 割点、桥
        """

        def dfs(cur: int, parent: int) -> None:
            if visited[cur]:
                return
            visited[cur] = True

            nonlocal dfsId
            order[cur] = low[cur] = dfsId
            dfsId += 1

            dfsChild = 0
            for next in adjMap[cur]:
                if next == parent:
                    continue
                if not visited[next]:
                    dfsChild += 1
                    dfs(next, cur)
                    low[cur] = min(low[cur], low[next])
                    if low[next] > order[cur]:
                        cuttingEdge.append((cur, next))
                    if parent != -1 and low[next] >= order[cur]:
                        cuttingPoint.add(cur)
                    elif parent == -1 and dfsChild > 1:
                        cuttingPoint.add(cur)
                else:
                    low[cur] = min(low[cur], low[next])

        dfsId = 0
        order, low = [Tarjan.INF] * n, [Tarjan.INF] * n
        visited = [False] * n

        cuttingPoint = set()
        cuttingEdge = []

        dfs(0, -1)

        return list(cuttingPoint), cuttingEdge


class Solution:
    def criticalConnections(self, n: int, connections: List[List[int]]) -> List[List[int]]:
        adjMap = defaultdict(set)
        for u, v in connections:
            adjMap[u].add(v)
            adjMap[v].add(u)

        return Tarjan.getCuttingPointAndCuttingEdge(n, adjMap)[1]

