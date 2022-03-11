# 1. tarjan缩点成DAG
# 2. 当一个强连通的出度为0,则该强连通分量中的所有点都被其他强连通分量的牛欢迎
# 但假如存在两及以上个出度=0的牛(强连通分量) 则必然有一头牛(强连通分量)不被所有牛欢迎


# 现在有 N 头牛，编号从 1 到 N，给你 M 对整数 (A,B)，表示牛 A 认为牛 B 受欢迎。
# 这种关系是具有传递性的，如果 A 认为 B 受欢迎，B 认为 C 受欢迎，那么牛 A 也认为牛 C 受欢迎。
# 你的任务是求出有多少头牛被除自己之外的所有牛认为是受欢迎的。
# 1≤N≤104 ,
# 1≤M≤5×104

from typing import DefaultDict, List, Set, Tuple
from collections import defaultdict
import sys

sys.setrecursionlimit(1000000)


class Tarjan:
    INF = int(1e20)

    @staticmethod
    def getSCC(
        n: int, adjMap: DefaultDict[int, Set[int]]
    ) -> Tuple[int, List[List[int]], List[int]]:
        """Tarjan求解有向图的强连通分量

        Args:
            n (int): 结点数
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
            for next in adjMap[cur]:
                if not visited[next]:
                    dfs(next)
                    # 回退阶段
                    low[cur] = min(low[cur], low[next])
                elif inStack[next]:
                    # next被访问过而且也在stack里面，找到了一个环
                    low[cur] = min(low[cur], low[next])
                # 访问过但是不在stack里的点，说明是在别的SCC里面被统计过了

            # 这个条件说明我们在当前这一轮找到了一个SCC并且cur是SCC的最高点
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
        stack = []  # 用来存当前DFS访问的点
        inStack = [False] * n

        SCCId = 0
        SCCGroupById = [[] for _ in range(n)]
        SCCIdByNode = [-1] * n

        for cur in range(n):
            if not visited[cur]:
                dfs(cur)

        return SCCId, SCCGroupById, SCCIdByNode


adjMap = defaultdict(set)
n, m = map(int, input().split())
for _ in range(m):
    a, b = map(int, input().split())
    a, b = a - 1, b - 1
    adjMap[a].add(b)

# 1. tarjan缩点
SCCId, SCCGroupById, SCCIdByNode = Tarjan.getSCC(n, adjMap)
# 存每一个SCC的出度，后面要找出度为0的SCC
outd = [0] * SCCId
for cur in range(n):
    for next in adjMap[cur]:
        # 不是同一个SCC 就连边
        if SCCIdByNode[next] != SCCIdByNode[cur]:
            outd[SCCIdByNode[cur]] += 1

# 2. 找出度为0的SCC个数
res = 0
zero = 0
for id in range(SCCId):
    if outd[id] == 0:
        res += len(SCCGroupById[id])
        zero += 1
        if zero > 1:
            res = 0
            break

print(res)

