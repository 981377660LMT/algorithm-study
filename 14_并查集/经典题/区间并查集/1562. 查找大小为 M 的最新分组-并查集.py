from collections import defaultdict
from typing import List
from sortedcontainers import SortedList


class UnionFind:
    """区间并查集维护各个区域的sum值"""

    def __init__(self, n: int):
        self.n = n
        self.part = n
        self.parent = list(range(n))
        self.size = [1] * n
        self.partSum = [0] * n
        self.partSumCounter = defaultdict(int)

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
        preSum1, preSum2 = self.partSum[rootX], self.partSum[rootY]
        self.partSum[rootY] += self.partSum[rootX]
        self.partSumCounter[preSum1] -= 1
        self.partSumCounter[preSum2] -= 1
        self.partSumCounter[self.partSum[rootY]] += 1
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
        self.partSumCounter[val] += 1


class Solution:
    def findLatestStep(self, arr: List[int], m: int) -> int:
        """
        返回存在长度 恰好 为 m 的 一组 1  的最后步骤
        如果不存在这样的步骤，请返回 -1 。

        正向并查集记录每组1的长度
        """
        n = len(arr)
        if m == n:
            return n
        visited = [False] * n
        uf = UnionFind(n)
        res = -1
        for i, pos in enumerate(arr, 1):
            pos -= 1
            visited[pos] = True
            uf.setPartSum(pos, 1)
            if pos - 1 >= 0 and visited[pos - 1]:
                uf.union(pos - 1, pos)
            if pos + 1 < n and visited[pos + 1]:
                uf.union(pos + 1, pos)
            if uf.partSumCounter[m] > 0:
                res = i
        return res

    def findLatestStep2(self, arr: List[int], m: int) -> int:
        """反向查找第一个 出现m个1的情况"""
        n = len(arr)

        if m > n:
            return -1
        if m == n:
            return m

        zeros = SortedList([-1, n])

        for i in range(n - 1, -1, -1):
            pos = arr[i] - 1
            zeros.add(pos)
            index = zeros.bisect_left(pos)

            if index + 1 < len(zeros):
                rightPos = zeros[index + 1]
                if rightPos - pos == m + 1:
                    return i

            if index > 0:
                leftPos = zeros[index - 1]
                if pos - leftPos == m + 1:
                    return i

        return -1


print(Solution().findLatestStep(arr=[3, 5, 1, 2, 4], m=1))
# 输出：4
# 解释：
# 步骤 1："00100"，由 1 构成的组：["1"]
# 步骤 2："00101"，由 1 构成的组：["1", "1"]
# 步骤 3："10101"，由 1 构成的组：["1", "1", "1"]
# 步骤 4："11101"，由 1 构成的组：["111", "1"]
# 步骤 5："11111"，由 1 构成的组：["11111"]
# 存在长度为 1 的一组 1 的最后步骤是步骤 4 。
