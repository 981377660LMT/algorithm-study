# 750. 角矩形的数量
# https://leetcode.cn/problems/number-of-corner-rectangles/description/
# 给定一个只包含 0 和 1 的 m x n 整数矩阵 grid ，返回 其中 「角矩形 」的数量 。
# 一个「角矩形」是由四个不同的在矩阵上的 1 形成的 轴对齐 的矩形。注意只有角的位置才需要为 1。
# 注意：4 个 1 的位置需要是不同的。
#
# m*n<=1e5

from typing import List
from collections import defaultdict
from itertools import combinations


class Solution:
    def countCornerRectangles1(self, grid: List[List[int]]) -> int:
        """枚举两个列, O(r*c^2)"""
        row, col = len(grid), len(grid[0])
        if row < col:
            grid = list(zip(*grid))  # type: ignore
            row, col = col, row

        dp = defaultdict(int)
        res = 0
        for r in grid:
            for c1, v1 in enumerate(r):
                if v1:
                    for c2 in range(c1 + 1, len(r)):
                        if r[c2]:
                            res += dp[(c1, c2)]
                            dp[(c1, c2)] += 1
        return res

    def countCornerRectangles2(self, grid: List[List[int]]) -> int:
        """枚举两个行，求列交集, O(r^2*c/64). 最快."""
        row, col = len(grid), len(grid[0])
        if row > col:
            grid = list(zip(*grid))  # type: ignore
            row, col = col, row

        colOnes = [0] * row  # row -> mask of col with 1
        for i, r in enumerate(grid):
            mask = sum((1 << j for j, v in enumerate(r) if v))
            colOnes[i] = mask

        res = 0
        for ones1, ones2 in combinations(colOnes, 2):
            intersection = (ones1 & ones2).bit_count()
            res += intersection * (intersection - 1) // 2
        return res

    def countCornerRectangles(self, grid: List[List[int]]) -> int:
        """
        根号分治, O(Nsqrt(N)+r*c), 其中 N 是网格中 1 的数量.
        这一行1很少时, 枚举两个列(解法1);
        这一行1很多时, 这样的行不多时, 枚举两个行计算列交集(解法2).
        """
        row = len(grid)
        allOnes = 0
        colOnes1 = [[] for _ in range(row)]  # row -> col with 1
        colOnes2 = [0] * row  # row -> mask of col with 1
        for i, cols in enumerate(grid):
            mask = 0
            for j, v in enumerate(cols):
                if v:
                    allOnes += 1
                    mask |= 1 << j
                    colOnes1[i].append(j)
            colOnes2[i] = mask

        sqrt_ = max(1, int(allOnes**0.5))

        dp = defaultdict(int)
        res = 0
        for i, cols in enumerate(colOnes1):
            if len(cols) < sqrt_:
                for pair in combinations(cols, 2):
                    res += dp[pair]
                    dp[pair] += 1
            else:
                for j, cols2 in enumerate(colOnes1):
                    if j <= i and len(cols2) >= sqrt_:
                        continue
                    intersection = (colOnes2[i] & colOnes2[j]).bit_count()
                    res += intersection * (intersection - 1) // 2
        return res
