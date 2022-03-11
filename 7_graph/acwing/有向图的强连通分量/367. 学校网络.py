# 最少需要将一个新软件直接提供给多少个学校，才能使软件能够通过网络被传送到所有学校？
# 最少需要添加几条新的支援关系，使得将一个新软件提供给任何一个学校，其他所有学校就都可以通过网络获得该软件？

# 缩点成DAG之后，
# 第一问相当于问多少个入度为0的点
# 第二问相当于问max(出度为零点数，入度为零点数)

from typing import DefaultDict, List, Set, Tuple
from collections import defaultdict
import sys

sys.setrecursionlimit(1000000)


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
        order, low = [int(1e20)] * n, [int(1e20)] * n

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
n = int(input())

for i in range(n):
    nums = list(map(int, input().split()))
    nums.pop()
    for num in nums:
        adjMap[i].add(num - 1)

# 1. tarjan缩点
SCCId, SCCGroupById, SCCIdByNode = Tarjan.getSCC(n, adjMap)
# 存每一个SCC的出度，后面要找出度为0的SCC
ind, outd = [0] * SCCId, [0] * SCCId
for cur in range(n):
    for next in adjMap[cur]:
        if SCCIdByNode[next] != SCCIdByNode[cur]:
            ind[SCCIdByNode[next]] += 1
            outd[SCCIdByNode[cur]] += 1

# 2. 找入度出度为0的SCC个数
start, end = ind.count(0), outd.count(0)
print(start)
# 特判:最后缩成了一个点
print(0 if SCCId == 1 else max(start, end))

