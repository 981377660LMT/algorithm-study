# 所有数字都不同 求使得数组递增的最小交换次数
# 此题类似于情侣牵手
# !把哪些冲突的放在一组，解决这些冲突需要(size-1)次交换


from collections import defaultdict
from typing import DefaultDict, List, Set


class Solution:
    def solve(self, nums: List[int]) -> int:
        n = len(nums)
        uf = UnionFind(n)
        target = sorted(nums)
        indexes = {num: i for i, num in enumerate(nums)}
        for i in range(n):
            uf.union(i, indexes[target[i]])

        res = 0
        for group in uf.getGroups().values():
            if len(group) > 1:
                res += len(group) - 1
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

    def getGroups(self) -> DefaultDict[int, Set[int]]:
        groups = defaultdict(set)
        for key in range(self.n):
            root = self.find(key)
            groups[root].add(key)
        return groups


print(Solution().solve(nums=[3, 2, 1, 4]))
# We can swap 3 and 1.
