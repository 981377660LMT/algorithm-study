from typing import DefaultDict, List, Set, Tuple
from collections import defaultdict, deque
import sys

from numpy import flip

# Q个询问，问有向图中从 a  是否能够到达 b 输出'Yes'或'No'
# n,m<=1e5


# !如果是拓扑图:拓扑序+位运算记录祖先结点
# !如果不是拓扑图:缩点成DAG


sys.setrecursionlimit(int(1e9))
input = sys.stdin.readline
MOD = int(1e9 + 7)


class Tarjan:
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
        order, low = [int(1e9)] * n, [int(1e9)] * n

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


n, m, q = map(int, input().split())
adjMap = defaultdict(set)
for _ in range(m):
    a, b = map(int, input().split())
    a, b = a - 1, b - 1
    adjMap[a].add(b)

# 缩点成DAG
sccCount, sccGroupById, sccIdByNode = Tarjan.getSCC(n, adjMap)
newAdjMap = defaultdict(set)
deg = defaultdict(int)
visitedPair = set()
for cur in range(n):
    for next in adjMap[cur]:
        g1, g2 = sccIdByNode[cur], sccIdByNode[next]
        if g1 == g2 or (g1, g2) in visitedPair:
            continue
        visitedPair.add((g1, g2))
        newAdjMap[g1].add(g2)
        deg[g2] += 1


dp = [1 << i for i in range(sccCount)]
queue = [i for i in range(sccCount) if deg[i] == 0]
while queue:
    nextQueue = []
    len_ = len(queue)
    for _ in range(len_):
        cur = queue.pop()
        for next in newAdjMap[cur]:
            deg[next] -= 1
            dp[next] |= dp[cur]
            if deg[next] == 0:
                nextQueue.append(next)
    queue = nextQueue

for _ in range(q):
    a, b = map(int, input().split())
    a, b = a - 1, b - 1
    scc1, scc2 = sccIdByNode[a], sccIdByNode[b]
    print('Yes' if dp[scc2] & (1 << scc1) else 'No', flush=True)
