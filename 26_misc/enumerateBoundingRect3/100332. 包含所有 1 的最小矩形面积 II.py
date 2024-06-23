# 100332. 包含所有 1 的最小矩形面积 II
# https://leetcode.cn/problems/find-the-minimum-area-to-cover-all-ones-ii/description/
# 找到 3 个 不重叠、面积 非零 、边在水平方向和竖直方向上的矩形，并且满足 grid 中所有的 1 都在这些矩形的内部。
# 返回这些矩形面积之和的 最小 可能值。

from typing import List
from enumerateBoundingRect3 import BoundingRect, enumerateBoundingRect3


INF = int(1e18)


def max2(a: int, b: int) -> int:
    return a if a > b else b


def min2(a: int, b: int) -> int:
    return a if a < b else b


class Solution:
    def minimumSum(self, grid: List[List[int]]) -> int:
        def calc(boundingRect: BoundingRect) -> int:
            """TODO: 前缀优化成 O(nm)"""
            top, bottom, left, right = boundingRect
            minTop, maxBottom, minLeft, maxRight = INF, -INF, INF, -INF
            for r in range(top, bottom + 1):
                for c in range(left, right + 1):
                    if grid[r][c] == 1:
                        minTop = min2(minTop, r)
                        maxBottom = max2(maxBottom, r)
                        minLeft = min2(minLeft, c)
                        maxRight = max2(maxRight, c)
            if minTop == INF or maxBottom == -INF or minLeft == INF or maxRight == -INF:
                return 0
            return (maxRight - minLeft + 1) * (maxBottom - minTop + 1)

        res = INF
        row, col = len(grid), len(grid[0])
        for b1, b2, b3 in enumerateBoundingRect3(row, col):
            res = min2(res, calc(b1) + calc(b2) + calc(b3))
        return res
