from typing import List


MOD = int(1e9 + 7)
INF = int(1e20)


class UnionFind:
    def __init__(self, n: int):
        self.n = n
        self.part = n
        self.parent = list(range(n))
        self.size = [1] * n
        self.partSum = [0] * n

    def find(self, x: int) -> int:
        if x != self.parent[x]:
            self.parent[x] = self.find(self.parent[x])
        return self.parent[x]

    def union(self, x: int, y: int) -> bool:
        rootX = self.find(x)
        rootY = self.find(y)
        if rootX == rootY:
            return False
        if self.size[rootX] > self.size[rootY]:
            rootX, rootY = rootY, rootX
        self.parent[rootX] = rootY
        self.size[rootY] += self.size[rootX]
        self.partSum[rootY] += self.partSum[rootX]
        self.part -= 1
        return True

    def isConnected(self, x: int, y: int) -> bool:
        return self.find(x) == self.find(y)

    def getPartSum(self, x: int) -> int:
        root = self.find(x)
        return self.partSum[root]


class Solution:
    def maximumSegmentSum(self, nums: List[int], removeQueries: List[int]) -> List[int]:
        """倒序"""

        n = len(nums)
        uf = UnionFind(n)
        visited = [False] * n
        res = [0] * n

        for i in range(n - 1, 0, -1):
            curIndex = removeQueries[i]
            visited[curIndex] = True
            uf.partSum[curIndex] = nums[curIndex]

            if curIndex - 1 >= 0 and visited[curIndex - 1]:
                uf.union(curIndex - 1, curIndex)
            if curIndex + 1 < n and visited[curIndex + 1]:
                uf.union(curIndex + 1, curIndex)

            res[i - 1] = max(res[i], uf.getPartSum(curIndex))

        return res
