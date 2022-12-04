# 求无向图中1到n中所有路径的最小边权
# 一条路径可以 多次 包含同一条道路，你也可以沿着路径多次到达城市 1 和城市 n 。
# 测试数据保证城市 1 和城市n 之间 至少 有一条路径。

# !连通分量中的最小边权 => 并查集/bfs/dfs

from collections import deque
from typing import List

INF = int(1e18)


class Solution:
    def minScore(self, n: int, roads: List[List[int]]) -> int:
        """并查集"""
        uf = UnionFind(n + 1)
        for u, v, _ in roads:
            uf.union(u, v)
        res = INF
        for u, v, w in roads:
            if uf.isConnected(u, 1):  # 这条边与1连通
                res = min(res, w)
        return res

    def minScore2(self, n: int, roads: List[List[int]]) -> int:
        """bfs求连通块最大边权"""
        adjList = [[] for _ in range(n)]
        for u, v, w in roads:
            u, v = u - 1, v - 1
            adjList[u].append((v, w))
            adjList[v].append((u, w))

        res = INF
        queue = deque([0])
        visited = set([0])
        while queue:
            cur = queue.popleft()
            for next, weight in adjList[cur]:
                res = min(res, weight)
                if next in visited:
                    continue
                visited.add(next)
                queue.append(next)
        return res

    def minScore3(self, n: int, roads: List[List[int]]) -> int:
        """dfs求连通块最大边权"""

        def dfs(cur: int) -> None:
            nonlocal res
            for next, weight in adjList[cur]:
                res = min(res, weight)
                if next in visited:
                    continue
                visited.add(next)
                dfs(next)

        adjList = [[] for _ in range(n)]
        for u, v, w in roads:
            u, v = u - 1, v - 1
            adjList[u].append((v, w))
            adjList[v].append((u, w))

        res = INF
        visited = set([0])
        dfs(0)
        return res


class UnionFind:

    __slots__ = ("n", "part", "parent", "rank")

    def __init__(self, n: int):
        self.n = n
        self.part = n
        self.parent = list(range(n))
        self.rank = [1] * n

    def find(self, x: int) -> int:
        while x != self.parent[x]:
            self.parent[x] = self.parent[self.parent[x]]
            x = self.parent[x]
        return x

    def union(self, x: int, y: int) -> bool:
        rootX = self.find(x)
        rootY = self.find(y)
        if rootX == rootY:
            return False
        if self.rank[rootX] > self.rank[rootY]:
            rootX, rootY = rootY, rootX
        self.parent[rootX] = rootY
        self.rank[rootY] += self.rank[rootX]
        self.part -= 1
        return True

    def isConnected(self, x: int, y: int) -> bool:
        return self.find(x) == self.find(y)


# 4
# [[1,2,9],[2,3,6],[2,4,5],[1,4,7]]
print(Solution().minScore2(4, [[1, 2, 9], [2, 3, 6], [2, 4, 5], [1, 4, 7]]))
