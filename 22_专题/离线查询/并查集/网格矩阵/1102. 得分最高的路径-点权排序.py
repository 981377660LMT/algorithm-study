# !路径的得分是`该路径上的 最小 值`。例如，路径 8 →  4 →  5 →  9 的值为 4
# !找出所有路径中得分 `最高` 的那条路径，返回其 得分。
# 并查集+排序连通路径

from collections import defaultdict
from typing import DefaultDict, List

DIR4 = [(0, 1), (0, -1), (1, 0), (-1, 0)]


class Solution:
    def maximumMinimumPath(self, grid: List[List[int]]) -> int:
        ROW, COL = len(grid), len(grid[0])
        nums = []
        for r in range(ROW):
            for c in range(COL):
                nums.append((grid[r][c], r, c))
        nums.sort(key=lambda x: x[0], reverse=True)

        uf = UnionFindArray(ROW * COL)
        for num, r, c in nums:  # 从大到小连通路径
            for dr, dc in DIR4:
                nr, nc = r + dr, c + dc
                if 0 <= nr < ROW and 0 <= nc < COL and grid[nr][nc] >= num:
                    uf.union(r * COL + c, nr * COL + nc)
                if uf.isConnected(0, ROW * COL - 1):
                    return num

        raise Exception("No Solution")


class UnionFindArray:

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

    def getGroups(self) -> DefaultDict[int, List[int]]:
        groups = defaultdict(list)
        for key in range(self.n):
            root = self.find(key)
            groups[root].append(key)
        return groups

    def getRoots(self) -> List[int]:
        return list(set(self.find(key) for key in self.parent))

    def __repr__(self) -> str:
        return "\n".join(f"{root}: {member}" for root, member in self.getGroups().items())

    def __len__(self) -> int:
        return self.part


print(Solution().maximumMinimumPath([[5, 4, 5], [1, 2, 6], [7, 4, 6]]))
