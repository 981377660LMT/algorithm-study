from typing import List


class UnionFindArray:
    def __init__(self, n: int):
        self.n = n
        self.part = n
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

        rootX, rootY = sorted([rootX, rootY], reverse=True)
        # 小的总是指向大的
        self.parent[rootY] = rootX
        self.rank[rootX] += self.rank[rootY]
        self.part -= 1
        return True

    def isConnected(self, x: int, y: int) -> bool:
        return self.find(x) == self.find(y)


class Solution:
    def validSubarraySize(self, nums: List[int], threshold: int) -> int:
        """并查集维护区间标记 从大到小遍历 把看过的区间串起来"""
        n = len(nums)
        uf = UnionFindArray(n + 10)
        Q = sorted(((num, i) for i, num in enumerate(nums)), reverse=True)  # 数组中的元素越大越好
        for num, i in Q:
            uf.union(i, i + 1)  # 向右连接
            length = uf.rank[uf.find(i + 1)] - 1  # 串联区间的长度
            if num * length > threshold:
                return length
        return -1
