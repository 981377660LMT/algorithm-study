# 在由 1 x 1 方格组成的 n x n 网格 grid 中，每个 1 x 1 方块由 '/'、'\' 或空格构成。
# 这些字符会将方块划分为一些共边的区域。
# !给定网格 grid 表示为一个字符串数组，返回 `区域的数量`` 。
# 请注意，反斜杠字符是转义的，因此 '\' 用 '\\' 表示。


# 把每一个小方块看成如下图
#    0
#  3   1
#    2


from collections import defaultdict
from typing import DefaultDict, List


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


class Solution:
    def regionsBySlashes(self, grid: List[str]) -> int:
        ROW, COL = len(grid), len(grid[0])
        uf = UnionFindArray(ROW * COL * 4)

        for r in range(ROW):
            for c in range(COL):
                pos = r * COL + c
                top, right, bottom, left = pos * 4, pos * 4 + 1, pos * 4 + 2, pos * 4 + 3
                cur = grid[r][c]

                # 方格内部连通(top, right, bottom, left)
                if cur == " ":
                    uf.union(top, right)
                    uf.union(top, bottom)
                    uf.union(top, left)
                elif cur == "/":
                    uf.union(top, left)
                    uf.union(right, bottom)
                elif cur == "\\":
                    uf.union(top, right)
                    uf.union(bottom, left)

                # 方格之间连通(看右边和下面的方格)
                if c + 1 < COL:
                    rightPos = r * COL + c + 1
                    rightLeft = rightPos * 4 + 3
                    uf.union(right, rightLeft)
                if r + 1 < ROW:
                    bottomPos = (r + 1) * COL + c
                    bottomUp = bottomPos * 4
                    uf.union(bottom, bottomUp)

        return uf.part


print(Solution().regionsBySlashes([" /", "/ "]))  # 2
# print(Solution().regionsBySlashes(grid=["/\\", "\\/"]))  # 5
