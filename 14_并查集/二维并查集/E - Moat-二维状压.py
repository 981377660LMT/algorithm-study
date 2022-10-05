# Moat(护城河/壕沟)

# !有一个 4*4 的网格和 若干房子(1)，要求在沿轮廓线上围一个圈，围住所有房子，求方案数
# !状压+并查集
# https://blog.csdn.net/weixin_45750972/article/details/120413237

# 先状压，然后所有包含给定输入的情况，
# !先将4×4的矩阵外围一圈变成6×6的大小，然后用dsu进行判断，
# 若区域外集合-大小和区域内集合大小刚好是36则满足条件。
# 这里用6×6的目的，方便以(0,0)作为区域外的一个点，
# 然后我们对区域内的一个点做标记，然后进行并查集即可。
# !最后要求恰好两个连通块,且两个连通块的大小之和为36
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

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


if __name__ == "__main__":
    grid = []
    for _ in range(4):
        row = list(map(int, input().split()))
        grid.append(row)

    need = 0
    for r in range(4):
        for c in range(4):
            if grid[r][c] == 1:
                id = r * 4 + c
                need |= 1 << id

    res = 0
    for state in range(1 << 16):
        if state & need != need:
            continue

        color = [[0] * 6 for _ in range(6)]
        for r in range(4):
            for c in range(4):
                color[r + 1][c + 1] = (state >> (r * 4 + c)) & 1

        uf = UnionFindArray(36)
        for r in range(6):
            for c in range(6):
                if c + 1 < 6 and color[r][c] == color[r][c + 1]:
                    uf.union(r * 6 + c, r * 6 + c + 1)
                if r + 1 < 6 and color[r][c] == color[r + 1][c]:
                    uf.union(r * 6 + c, (r + 1) * 6 + c)

        roots = uf.getRoots()
        if len(roots) == 2 and uf.rank[roots[0]] + uf.rank[roots[1]] == 36:
            res += 1

    print(res)
