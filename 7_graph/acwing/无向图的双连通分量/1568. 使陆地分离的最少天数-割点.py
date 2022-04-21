# 一共3种情况，0，1，2， 并查集求岛屿数量如果大于2 返回0， 如果岛屿数量为1， tarjan算法求割点，
# 如果找到割点返回 1，没有割点则返回，2
from typing import DefaultDict, List, Set, Tuple
from collections import defaultdict
from copy import deepcopy


class Tarjan:
    INF = int(1e20)

    @staticmethod
    def getSCC(
        n: int, adjMap: DefaultDict[int, Set[int]]
    ) -> Tuple[int, DefaultDict[int, Set[int]], List[int]]:
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

            for next in adjMap[cur]:
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
    def getCuttingPointAndCuttingEdge(
        n: int, adjMap: DefaultDict[int, Set[int]]
    ) -> Tuple[Set[int], Set[Tuple[int, int]]]:
        """Tarjan求解无向图的割点和割边(桥)

        Args:
            n (int): 结点0-n-1
            adjMap (DefaultDict[int, Set[int]]): 图

        Returns:
            Tuple[Set[int], Set[Tuple[int, int]]]: 割点、桥

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
            for next in adjMap[cur]:
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
                    elif parent == -1 and dfsChild > 1:  # 出发点没有祖先啊，所以特判一下
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
    def getVBCC(
        n: int, adjMap: DefaultDict[int, Set[int]]
    ) -> Tuple[int, DefaultDict[int, Set[int]], List[Set[int]]]:
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
            for next in adjMap[cur]:
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
        n: int, adjMap: DefaultDict[int, Set[int]]
    ) -> Tuple[int, DefaultDict[int, Set[Tuple[int, int]]], DefaultDict[int, int]]:
        """Tarjan求解无向图的边双联通分量

        Args:
            n (int): 结点0-n-1
            adjMap (DefaultDict[int, Set[int]]): 图

        Returns:
            Tuple[int, DefaultDict[int, Set[Tuple[int, int]]], List[int]]: EBCC的数量、分组、每条边对应的EBCC编号

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

            for next in adjMap[cur]:
                if next == parent:
                    continue

                edge = tuple(sorted([cur, next]))
                if edge in cuttingEdges:
                    continue

                EBCCGroupById[EBCCId].add(edge)
                EBCCIdByNode[cur] = EBCCId
                dfs(next, cur)

        _, cuttingEdges = Tarjan.getCuttingPointAndCuttingEdge(n, adjMap)

        visited = [False] * n

        EBCCId = 0  # 边双个数
        EBCCGroupById = defaultdict(set)  # 每个边双包含哪些边
        EBCCIdByNode = defaultdict(int)  # 每条边属于哪一个边双

        for cur in range(n):
            if not visited[cur]:
                dfs(cur, -1)
                EBCCId += 1

        for edge in cuttingEdges:
            EBCCGroupById[EBCCId].add(edge)
            EBCCIdByNode[edge] = EBCCId
            EBCCId += 1

        return EBCCId, EBCCGroupById, EBCCIdByNode


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
