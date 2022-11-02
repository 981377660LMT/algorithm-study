"""
给定一个 0-n-1 (n<=1e5) 的排列 其中元素0表示空

每次操作可以将任意一个元素移动到空的位置
使得元素变为[1,2,3,...,n-1,0] 或者 [0,1,2,3...n-1]
问最少需要多少次操作
"""


from collections import defaultdict
from typing import DefaultDict, List, Set


class Solution:
    def sortArray(self, nums: List[int]) -> int:
        """并查集找到所有的置换环

        如果 0 所在位置位于环中，那么这个环产生的交换次数为 `环的大小 - 1`.
        否则为 `环的大小 + 1`.
        """

        def cal(target: List[int]) -> int:
            uf = UnionFind(n)
            for i, num in enumerate(target):
                uf.union(i, indexes[num])

            zeroIndex = target.index(0)
            res = 0
            for group in uf.getGroups().values():
                if len(group) > 1:
                    res += len(group) - 1 if zeroIndex in group else len(group) + 1
            return res

        n = len(nums)
        indexes = {v: i for i, v in enumerate(nums)}
        target1 = [0] + list(range(1, n))
        target2 = list(range(1, n)) + [0]
        return min(cal(target1), cal(target2))


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


assert Solution().sortArray(nums=[4, 2, 0, 3, 1]) == 3
assert Solution().sortArray(nums=[1, 2, 3, 4, 0]) == 0
assert Solution().sortArray(nums=[3, 0, 1, 2]) == 3
assert Solution().sortArray(nums=[3, 2, 1, 6, 7, 4, 5, 0]) == 8

# Input: nums = [4,2,0,3,1]
# Output: 3
# Explanation:
# - Move item 2 to the empty space. Now, nums = [4,0,2,3,1].
# - Move item 1 to the empty space. Now, nums = [4,1,2,3,0].
# - Move item 4 to the empty space. Now, nums = [0,1,2,3,4].
# It can be proven that 3 is the minimum number of operations needed.
