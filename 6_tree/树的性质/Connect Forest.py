# 1. 并查集获取各个树
# 2. 获取各个树的直径 讨论直径数 1 2 >=3
# 3. 连接直径
from collections import defaultdict, deque
from typing import DefaultDict, List


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

    def getGroups(self) -> DefaultDict[int, List[int]]:
        groups = defaultdict(list)
        for key in range(self.n):
            root = self.find(key)
            groups[root].append(key)
        return groups


def getDiameter(adjList: List[List[int]], start: int) -> int:
    queue = deque([start])
    visited = set([start])
    lastVisited = 0  # 全局变量，好记录第一次BFS最后一个点的ID
    while queue:
        curLen = len(queue)
        for _ in range(curLen):
            lastVisited = queue.popleft()
            for next in adjList[lastVisited]:
                if next not in visited:
                    visited.add(next)
                    queue.append(next)

    queue = deque([lastVisited])  # 第一次最后一个点，作为第二次BFS的起点
    visited = set([lastVisited])
    level = -1  # 记好距离
    while queue:
        curLen = len(queue)
        for _ in range(curLen):
            cur = queue.popleft()
            for next in adjList[cur]:
                if next not in visited:
                    visited.add(next)
                    queue.append(next)
        level += 1

    return level


class Solution:
    def solve(self, graph: List[List[int]]) -> int:
        n = len(graph)
        uf = UnionFindArray(n)
        for i, nexts in enumerate(graph):
            for j in nexts:
                uf.union(i, j)

        groups = uf.getGroups()
        diameters = [getDiameter(graph, root) for root in groups]

        if len(diameters) == 1:
            return diameters[0]

        if len(diameters) == 2:
            return max(max(diameters), (diameters[0] + 1) // 2 + (diameters[1] + 1) // 2 + 1)

        diameters = sorted(diameters, reverse=True)[:3]
        # 考虑第二第三大的直径
        return max(
            diameters[0],
            (diameters[0] + 1) // 2 + (diameters[1] + 1) // 2 + 1,
            (diameters[1] + 1) // 2 + (diameters[2] + 1) // 2 + 2,
        )


print(Solution().solve(graph=[[1, 2], [0], [0, 3], [2], [5], [4]]))
