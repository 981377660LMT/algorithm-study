# 3608. 包含 K 个连通分量需要的最小时间
# https://leetcode.cn/problems/minimum-time-for-k-connected-components/description/
#
# 给你一个整数 n，表示一个包含 n 个节点（从 0 到 n - 1 编号）的无向图。该图由一个二维数组 edges 表示，
# 其中 edges[i] = [ui, vi, timei] 表示一条连接节点 ui 和节点 vi 的无向边，该边会在时间 timei 被移除。
#
# 同时，另给你一个整数 k。
#
# 最初，图可能是连通的，也可能是非连通的。
# 你的任务是找到一个 最小 的时间 t，使得在移除所有满足条件 time <= t 的边之后，该图包含 至少 k 个连通分量。
#
# 返回这个 最小 时间 t。
#
# 删边变加边
from typing import List


class UnionFindArraySimple:
    __slots__ = ("part", "n", "_data")

    def __init__(self, n: int):
        self.part = n
        self.n = n
        self._data = [-1] * n

    def union(self, key1: int, key2: int) -> bool:
        root1, root2 = self.find(key1), self.find(key2)
        if root1 == root2:
            return False
        if self._data[root1] > self._data[root2]:
            root1, root2 = root2, root1
        self._data[root1] += self._data[root2]
        self._data[root2] = root1
        self.part -= 1
        return True

    def find(self, key: int) -> int:
        if self._data[key] < 0:
            return key
        self._data[key] = self.find(self._data[key])
        return self._data[key]

    def getSize(self, key: int) -> int:
        return -self._data[self.find(key)]


class Solution:
    def minTime(self, n: int, edges: List[List[int]], k: int) -> int:
        edges.sort(key=lambda x: x[2], reverse=True)
        uf = UnionFindArraySimple(n)
        for u, v, t in edges:
            uf.union(u, v)
            if uf.part < k:
                return t
        return 0
