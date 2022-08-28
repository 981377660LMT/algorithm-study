"""
# 1. 反向并查集：正着删除元素变成倒着添加元素
# 2. 每次添加元素，检查左右邻居是否可以合并
# 3. 当前的最大值就是 `max(之前所有区域的和中的最大值，当前区域的和)`

如果有负数 需要有序容器维护partSum
"""

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

    def setPartSum(self, x: int, val: int) -> None:
        root = self.find(x)
        self.partSum[root] = val


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
            uf.setPartSum(curIndex, nums[curIndex])

            # 合并左右邻居
            if curIndex - 1 >= 0 and visited[curIndex - 1]:
                uf.union(curIndex - 1, curIndex)
            if curIndex + 1 < n and visited[curIndex + 1]:
                uf.union(curIndex + 1, curIndex)
            # 反向更新
            res[i - 1] = max(res[i], uf.getPartSum(curIndex))

        return res


n = int(input())
nums = [int(v) for v in input().split()]
queries = [int(v) - 1 for v in input().split()]


# 输出答案
print(*Solution().maximumSegmentSum(nums, queries), sep="\n")
