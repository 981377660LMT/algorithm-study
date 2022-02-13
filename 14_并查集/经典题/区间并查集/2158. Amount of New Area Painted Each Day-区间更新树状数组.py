from typing import List


class UnionFindArray:
    def __init__(self, n: int):
        self.n = n
        self.count = n
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
        if self.rank[rootX] > self.rank[rootY]:
            rootX, rootY = rootY, rootX
        self.parent[rootX] = rootY
        self.rank[rootY] += self.rank[rootX]
        self.count -= 1
        return True

    def isConnected(self, x: int, y: int) -> bool:
        return self.find(x) == self.find(y)


# 对于每天涂的区域，从右向左遍历，将右边合并至左边。
# 这样，当从右向左访问到已被涂过的位置时，
# 可以通过并查集直接跳到当前这块已被涂过的区域的左边界，
# 直到跳到（也可能跳过）当天涂色的左边界为止
# 每走一步，都会对当天的答案产生1点贡献；


class Solution:
    def amountPainted(self, paint: List[List[int]]) -> List[int]:
        res = []

        return res


# print(Solution().amountPainted(paint=[[1, 4], [4, 7], [5, 8]]))
# print(Solution().amountPainted(paint=[[1, 4], [5, 8], [4, 7]]))
# print(Solution().amountPainted(paint=[[1, 5], [2, 4]]))
print(Solution().amountPainted(paint=[[2, 4], [0, 4], [1, 4], [1, 5], [0, 2]]))
# [2, 2, 0, 1, 0]
print(Solution().amountPainted(paint=[[0, 5], [0, 2], [0, 3], [0, 4], [0, 5]]))
# [5, 0, 0, 0, 0]
