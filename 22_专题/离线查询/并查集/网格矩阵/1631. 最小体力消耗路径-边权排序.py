# 一条路径耗费的 体力值 是路径上相邻格子之间 高度差绝对值 的 最大值 决定的。
# 请你返回从左上角走到右下角的最小 体力消耗值 。
# !并查集+按边权排序连通

from collections import defaultdict
from typing import DefaultDict, List


class Solution:
    def minimumEffortPath(self, heights: List[List[int]]) -> int:
        ROW, COL = len(heights), len(heights[0])
        edges = []
        for r in range(ROW):
            for c in range(COL):
                if r + 1 < ROW:
                    edges.append(
                        (abs(heights[r][c] - heights[r + 1][c]), r * COL + c, (r + 1) * COL + c)
                    )
                if c + 1 < COL:
                    edges.append(
                        (abs(heights[r][c] - heights[r][c + 1]), r * COL + c, r * COL + c + 1)
                    )

        edges.sort(key=lambda x: x[0])

        uf = UnionFindArray(ROW * COL)
        for num, p1, p2 in edges:
            uf.union(p1, p2)
            if uf.isConnected(0, ROW * COL - 1):
                return num

        return 0  # 不存在边, 说明只有一个点


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


assert Solution().minimumEffortPath([[1, 2, 2], [3, 8, 2], [5, 3, 5]]) == 2
