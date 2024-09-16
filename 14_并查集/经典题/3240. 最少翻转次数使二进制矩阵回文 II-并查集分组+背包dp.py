# 3240. 最少翻转次数使二进制矩阵回文 II
# https://leetcode.cn/problems/minimum-number-of-flips-to-make-binary-grid-palindromic-ii/description/
# 给你一个 m x n 的二进制矩阵 grid 。
# 如果矩阵中一行或者一列从前往后与从后往前读是一样的，那么我们称这一行或者这一列是 回文 的。
# 你可以将 grid 中任意格子的值 翻转 ，也就是将格子里的值从 0 变成 1 ，或者从 1 变成 0 。
# 请你返回 最少 翻转次数，使得矩阵中 所有 行和列都是 回文的 ，且矩阵中 1 的数目可以被 4 整除 。
#
# !1.行和列回文 => a[i][j] == a[i][~j] == a[~i][j] == a[~i][~j]
# !2.并查集分组+背包dp

from collections import defaultdict
from functools import lru_cache
from typing import List

INF = int(1e18)


def min2(a: int, b: int) -> int:
    return a if a < b else b


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
    def minFlips(self, grid: List[List[int]]) -> int:
        ROW, COL = len(grid), len(grid[0])
        uf = UnionFindArraySimple(ROW * COL)
        for i in range(ROW):
            for j in range(COL):
                uf.union(i * COL + j, i * COL + (COL - 1 - j))
                uf.union(i * COL + j, (ROW - 1 - i) * COL + j)

        groups = defaultdict(list)
        for i in range(ROW):
            for j in range(COL):
                groups[uf.find(i * COL + j)].append((i, j))

        # !每个分组内对余数dp,dp[i]表示1的数量模4为i的最小操作次数
        @lru_cache(None)
        def dfs(gi: int, mod: int) -> int:
            if gi == len(groups):
                return 0 if mod == 0 else INF

            m = len(groups[gi])
            ones = sum(grid[i][j] for i, j in groups[gi])
            res1 = dfs(gi + 1, mod) + ones  # 全为0
            res2 = dfs(gi + 1, (mod + m) % 4) + (m - ones)  # 全为1
            return min2(res1, res2)

        res = dfs(0, 0)
        dfs.cache_clear()
        return res
