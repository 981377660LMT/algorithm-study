from collections import defaultdict
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


# 给你一个大小为 n x n 二进制矩阵 grid 。最多 只能将一格 0 变成 1 。
# 返回执行此操作后，grid 中最大的岛屿面积是多少？
DIR4 = ((0, 1), (1, 0), (0, -1), (-1, 0))


class Solution:
    def largestIsland(self, grid: List[List[int]]) -> int:
        ROW, COL = len(grid), len(grid[0])
        uf = UnionFindArray(ROW * COL + 10)
        for r in range(ROW):
            for c in range(COL):
                if grid[r][c] == 1:
                    cur = r * COL + c
                    for dr, dc in DIR4:
                        nr, nc = r + dr, c + dc
                        if 0 <= nr < ROW and 0 <= nc < COL:
                            next = nr * COL + nc
                            if grid[nr][nc] == 1:
                                uf.union(cur, next)

        res = max(uf.rank)
        for r in range(ROW):
            for c in range(COL):
                if grid[r][c] == 1:
                    continue

                # !不实际添加点 而是看周围点的分组
                roots = set()
                for dr, dc in DIR4:
                    nr, nc = r + dr, c + dc
                    if 0 <= nr < ROW and 0 <= nc < COL and grid[nr][nc] == 1:
                        roots.add(uf.find(nr * COL + nc))
                res = max(res, sum(uf.rank[root] for root in roots) + 1)
        return res


print(Solution().largestIsland(grid=[[1, 0], [0, 1]]))
