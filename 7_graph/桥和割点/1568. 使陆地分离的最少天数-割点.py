# 一共3种情况，0，1，2， 并查集求岛屿数量如果大于2 返回0， 如果岛屿数量为1， tarjan算法求割点，
# 如果找到割点返回 1，没有割点则返回，2
from typing import DefaultDict, List, Set, Tuple
from collections import defaultdict
from copy import deepcopy


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

        for i in range(n):
            if not visited[i]:
                dfs(i, -1)

        return list(cuttingPoint), cuttingEdge


class Solution:
    def minDays(self, grid: List[List[int]]) -> int:
        # 特判
        oneCount = sum(row.count(1) for row in grid)
        if oneCount == 0:
            return 0
        elif oneCount == 1:
            return 1

        # 连通分量不为1的情况
        gridCopy = deepcopy(grid)
        part = self.floodFill(gridCopy)
        if part != 1:
            return 0

        # tarjan寻找割点
        adjMap = defaultdict(set)
        row, col = len(grid), len(grid[0])
        for i in range(row):
            for j in range(col):
                if grid[i][j] == 1:
                    cur = i * col + j
                    if i - 1 >= 0 and grid[i - 1][j] == 1:
                        next = (i - 1) * col + j
                        adjMap[cur].add(next)
                        adjMap[next].add(cur)
                    if j + 1 < col and grid[i][j + 1] == 1:
                        next = i * col + j + 1
                        adjMap[cur].add(next)
                        adjMap[next].add(cur)

        cuttingPoints, _ = Tarjan.getCuttingPointAndCuttingEdge(row * col, adjMap)
        if cuttingPoints:
            return 1
        else:
            return 2

    def floodFill(self, grid: List[List[int]]) -> int:
        def dfs(r: int, c: int) -> None:
            if grid[r][c] == 0:
                return
            grid[r][c] = 0
            for dr, dc in [(1, 0), (-1, 0), (0, 1), (0, -1)]:
                nr, nc = r + dr, c + dc
                if 0 <= nr < row and 0 <= nc < col and grid[nr][nc] == 1:
                    dfs(nr, nc)

        res = 0
        row, col = len(grid), len(grid[0])
        for r in range(row):
            for c in range(col):
                if grid[r][c] == 1:
                    res += 1
                    dfs(r, c)
        return res


print(Solution().minDays(grid=[[0, 1, 1, 0], [0, 1, 1, 0], [0, 0, 0, 0]]))
print(Solution().minDays(grid=[[0, 1, 0, 1, 1], [1, 1, 1, 1, 1], [1, 1, 1, 1, 1], [1, 1, 1, 1, 0]]))
