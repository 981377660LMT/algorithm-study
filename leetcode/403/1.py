from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)
# 给你一个二维 二进制 数组 grid。请你找出一个边在水平方向和竖直方向上、面积 最小 的矩形，并且满足 grid 中所有的 1 都在矩形的内部。


# 返回这个矩形可能的 最小 面积。


def max2(a: int, b: int) -> int:
    return a if a > b else b


def min2(a: int, b: int) -> int:
    return a if a < b else b


class Solution:
    def minimumArea(self, grid: List[List[int]]) -> int:
        minTop, maxBottom, minLeft, maxRight = INF, 0, INF, 0
        for r, row in enumerate(grid):
            for c, v in enumerate(row):
                if v == 1:
                    minTop = min2(minTop, r)
                    maxBottom = max2(maxBottom, r)
                    minLeft = min2(minLeft, c)
                    maxRight = max2(maxRight, c)
        return (maxRight - minLeft + 1) * (maxBottom - minTop + 1)
