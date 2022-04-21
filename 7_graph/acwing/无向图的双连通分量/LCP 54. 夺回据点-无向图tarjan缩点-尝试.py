from typing import Counter, DefaultDict, List, Set, Tuple
from collections import defaultdict


INF = int(1e20)


def getCuttingPointAndCuttingEdge(
    n: int, adjMap: DefaultDict[int, Set[int]]
) -> Tuple[Set[int], Set[Tuple[int, int]]]:
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
                    cuttingEdge.add((cur, next))
                if parent != -1 and low[next] >= order[cur]:
                    cuttingPoint.add(cur)
                elif parent == -1 and dfsChild > 1:
                    cuttingPoint.add(cur)
            else:
                low[cur] = min(low[cur], low[next])

    dfsId = 0
    order, low = [INF] * n, [INF] * n
    visited = [False] * n

    cuttingPoint = set()
    cuttingEdge = set()

    for i in range(n):
        if not visited[i]:
            dfs(i, -1)

    return cuttingPoint, cuttingEdge


class UnionFindArray:
    def __init__(self, n: int):
        self.n = n
        self.count = n
        self.parent = list(range(n))
        self.rank = [1] * n

    def find(self, x: int) -> int:
        if x != self.parent[x]:
            self.parent[x] = self.find(self.parent[x])
        return self.parent[x]

    def union(self, x: int, y: int) -> bool:
        rootX = self.find(x)
        rootY = self.find(y)
        if rootX == rootY:
            return False
        if self.rank[rootX] > self.rank[rootY]:
            rootX, rootY = rootY, rootX
        self.parent[rootX] = rootY
        self.rank[rootY] += self.rank[rootX]
        self.count -= 1
        return True

    def isConnected(self, x: int, y: int) -> bool:
        return self.find(x) == self.find(y)

    def getGroups(self, bad: Set[int]) -> DefaultDict[int, List[int]]:
        groups = defaultdict(list)
        for key in range(self.n):
            if key in bad:
                continue
            root = self.find(key)
            groups[root].append(key)
        return groups


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

        # 割点
        cuttingPoint = getCuttingPointAndCuttingEdge(n, adjMap)[0]

        # 并查集分组(不可以用并查集，因为并查集是边双连通不是点双连通)
        uf = UnionFindArray(n)
        for u, v in roads:
            if u in cuttingPoint or v in cuttingPoint:
                continue
            uf.union(u, v)

        counter = defaultdict(set)  # 每个联通分量连了几个割点
        for cur in cuttingPoint:
            for next in adjMap[cur]:
                groupRoot = uf.find(next)
                counter[groupRoot].add(cur)

        goodGroups = defaultdict(list)
        for i in range(n):
            if i in cuttingPoint:
                continue
            root = uf.find(i)
            # print(i, root, 999, counter[root])
            if len(counter[root]) <= 1:
                goodGroups[root].append(i)

        print(cuttingPoint, goodGroups)
        # # 各个分量
        # groups = uf.getGroups(cuttingPoint).values()
        costs = [min(cost[i] for i in group) for group in goodGroups.values()]

        if len(costs) == 1:
            return costs[0]
        else:
            return sum(costs) - max(costs)


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

