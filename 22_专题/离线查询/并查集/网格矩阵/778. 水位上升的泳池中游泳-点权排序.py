# 778. 水位上升的泳池中游泳

# 当开始下雨时，在时间为 t 时，水池中的水位为 t 。
# 你可以从一个平台游向四周相邻的任意一个平台，
# 但是前提是此时水位必须同时淹没这两个平台。
# 假定你可以瞬间移动无限距离，也就是默认在方格内部游动是不耗时的。
# 当然，在你游泳的时候你必须待在坐标方格里面。


# !你从坐标方格的左上平台 (0，0) 出发。
# !返回 你到达坐标方格的右下平台 (n-1, n-1) 所需的最少时间 。
# 并查集+排序模拟水位上升

from collections import defaultdict
from typing import DefaultDict, List

DIR4 = [(0, 1), (0, -1), (1, 0), (-1, 0)]


class Solution:
    def swimInWater(self, grid: List[List[int]]) -> int:
        ROW, COL = len(grid), len(grid[0])
        nums = []
        for r in range(ROW):
            for c in range(COL):
                nums.append((grid[r][c], r, c))
        nums.sort(key=lambda x: x[0])

        uf = UnionFindArray(ROW * COL)
        for num, r, c in nums:  # 从小到大模拟升高水位
            for dr, dc in DIR4:
                nr, nc = r + dr, c + dc
                if 0 <= nr < ROW and 0 <= nc < COL and grid[nr][nc] <= num:
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
