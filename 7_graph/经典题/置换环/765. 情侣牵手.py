# n 对情侣坐在连续排列的 2n 个座位上，想要牵到对方的手。
# 人和座位由一个整数数组 row 表示，其中 row[i] 是坐在第 i 个座位上的人的 ID。
# 情侣们按顺序编号，第一对是 (0, 1)，第二对是 (2, 3)，以此类推，最后一对是 (2n-2, 2n-1)。
# 返回 最少交换座位的次数，以便每对情侣可以并肩坐在一起。
# 每次交换可选择任意两人，让他们站起来交换座位。
# !情侣牵手 - 置换环

from collections import defaultdict
from typing import DefaultDict, List, Set


class Solution:
    def minSwapsCouples(self, nums: List[int]) -> int:
        """并查集寻找置换环

        如果我们有 k 对情侣形成了错误环，需要交换 k - 1 次才能让情侣牵手。
        问题转化成 n / 2 对情侣中，有多少个这样的环
        """

        n = len(nums)
        uf = UnionFind(n // 2)
        for i in range(0, n, 2):
            uf.union(nums[i] // 2, nums[i + 1] // 2)  # 除以2表示对应的分组
        return sum(len(group) - 1 for group in uf.getGroups().values() if len(group) > 1)


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
