from typing import DefaultDict, List, Set, Tuple
from collections import defaultdict


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
        VBCCIdByNode = [set() for _ in range(n)]  # 每个点属于哪一(几)个点双，属于多个点双的点就是割点

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
          - 将所有的割边删掉剩下的都是边连通分量了(其实可以用并查集做)
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
                EBCCIdByEdge[cur] = EBCCId
                dfs(next, cur)

        _, cuttingEdges = Tarjan.getCuttingPointAndCuttingEdge(n, adjMap)

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


# 为了防止魔物暴动，勇者在每一次夺回据点后（包括花费资源夺回据点后），
# 需要保证剩余的所有魔物据点之间是相连通的（不经过「已夺回据点」）。


# Tarjan缩点
# 先用 Tarjan 算法找出割点，去掉这些点会剩下若干个连通块。
# 抛弃掉同时与多个割点相连的连通块。
# 求出剩余的连通块中的最小权值。
# 如果仅有一个连通块，答案就是这个最小权值；否则，答案为所有最小权值之和减去它们的最大值。
class Solution:
    def minimumCost(self, cost: List[int], roads: List[List[int]]) -> int:
        n = len(cost)
        adjMap = defaultdict(set)
        for u, v in roads:
            adjMap[u].add(v)
            adjMap[v].add(u)

        # 找VBCC和割点
        VBCCId, VBCCGroup, VBCCIdByNode = Tarjan.getVBCC(n, adjMap)
        cuttingPoints = set(i for i in range(n) if len(VBCCIdByNode[i]) > 1)

        # 统计连通分量里包含几个原图的割点，不能选连了两个以上个点的分量
        counter = [sum(node in cuttingPoints for node in VBCCGroup[i]) for i in range(VBCCId)]
        goodGroups = [i for i in range(VBCCId) if counter[i] <= 1]

        # 不能选割点
        costs = [min(cost[v] for v in VBCCGroup[k] if v not in cuttingPoints) for k in goodGroups]
        return costs[0] if len(costs) == 1 else sum(costs) - max(costs)


print(
    Solution().minimumCost(
        cost=[1, 2, 3, 4, 5, 6], roads=[[0, 1], [0, 2], [1, 3], [2, 3], [1, 2], [2, 4], [2, 5]]
    )
)

print(Solution().minimumCost(cost=[3, 2, 1, 4], roads=[[0, 2], [2, 3], [3, 1]]))


# print(Solution().minimumCost(cost=[0, 1, 2, 3], roads=[[0, 1], [1, 2], [2, 3], [0, 3]]))
print(
    Solution().minimumCost(
        cost=[9, 2, 3, 4, 5, 6, 7],
        roads=[[1, 2], [1, 3], [2, 3], [3, 6], [6, 0], [0, 3], [4, 2], [2, 5], [4, 5]],
    )
)

