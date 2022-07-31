from typing import List


class UnionFind:
    def __init__(self, n):
        self.n = n
        self.part = n
        self.parent = list(range(n))
        self.size = [1] * n

    def find(self, x):
        if x != self.parent[x]:
            self.parent[x] = self.find(self.parent[x])
        return self.parent[x]

    def union(self, x, y):
        root_x = self.find(x)
        root_y = self.find(y)
        if root_x == root_y:
            return False
        if self.size[root_x] > self.size[root_y]:
            root_x, root_y = root_y, root_x
        self.parent[root_x] = root_y
        self.size[root_y] += self.size[root_x]
        self.part -= 1
        return True

    def isconnected(self, p, q):
        return self.find(p) == self.find(q)


# x 和 y 的两座城市直接连通的前提是： x 和 y 的公因数中，至少有一个 严格大于 某个阈值 threshold
# 给你两个整数 n 和 threshold ，以及一个待查询数组，
# 请你判断每个查询 queries[i] = [ai, bi] 指向的城市 ai 和 bi 是否连通
# 2 <= n <= 104
# 0 <= threshold <= n
class Solution:
    def areConnected(self, n: int, threshold: int, queries: List[List[int]]) -> List[bool]:
        """因子与倍数合并(这样具有传递性)"""
        ok = set(range(1, n + 1))
        uf = UnionFind(n + 1)
        res = []

        # 枚举因子与倍数合并nlogn
        for i in range(threshold + 1, n + 1):  # 按照threshold预处理并查集
            for j in range(2 * i, n + 1, i):
                if j in ok:
                    uf.union(i, j)

        for x, y in queries:  # 判断连通性，在不在一个连通域里
            res.append(uf.isconnected(x, y))

        return res


print(Solution().areConnected(n=100, threshold=1, queries=[[2, 9]]))
