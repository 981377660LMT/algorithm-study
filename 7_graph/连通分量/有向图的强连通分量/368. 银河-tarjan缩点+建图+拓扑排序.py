# 银河中的恒星浩如烟海，但是我们只关注那些最亮的恒星。
# 我们用一个正整数来表示恒星的亮度，数值越大则恒星就越亮，恒星的亮度最暗是 1。
# 现在对于 N 颗我们关注的恒星，有 M 对亮度之间的相对关系已经判明。
# 你的任务就是求出这 N 颗恒星的亮度值总和至少有多大。


from typing import DefaultDict, List, Set, Tuple
from collections import defaultdict
import sys

sys.setrecursionlimit(int(1e6))


class Tarjan:
    INF = int(1e20)

    @staticmethod
    def getSCC(
        n: int, adjMap: DefaultDict[int, Set[Tuple[int, int]]]
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
            for next, _ in adjMap[cur]:
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


n, m = map(int, input().split())
adjMap = defaultdict(set)
for _ in range(m):
    kind, a, b = map(int, input().split())
    if kind == 1:
        adjMap[a].add((b, 0))
        adjMap[b].add((a, 0))
    elif kind == 2:
        # A 的亮度小于 B 的亮度
        adjMap[a].add((b, 1))
    elif kind == 3:
        # A 的亮度不小于 B 的亮度
        adjMap[b].add((a, 0))
    elif kind == 4:
        adjMap[b].add((a, 1))
    else:
        adjMap[a].add((b, 0))

# 0号虚拟源点
for i in range(1, n + 1):
    adjMap[0].add((i, 1))

# 1. tarjan缩点+建新图
SCCId, SCCGroupById, SCCIdByNode = Tarjan.getSCC(n + 1, adjMap)

newAdjMap = defaultdict(set)
for cur in range(n):
    for next, weight in adjMap[cur]:
        id1, id2 = SCCIdByNode[cur], SCCIdByNode[next]
        if id1 == id2:
            # 有正环，无解
            if weight != 0:
                print(-1)
                exit(0)
        else:
            newAdjMap[id1].add((id2, weight))

# 2. 在新图上求拓扑排序最长路
# 这里不必拓扑排序，因为强连通分量反序结果就是拓扑序
n = SCCId
dist = [-int(1e20)] * n
dist[n - 1] = 0
for cur in range(n - 1, -1, -1):
    for next, weight in newAdjMap[cur]:
        dist[next] = max(dist[next], dist[cur] + weight)

res = 0
for d, g in zip(dist, SCCGroupById):
    res += d * len(g)

print(res)

