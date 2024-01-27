# 3009.折线图上的最大交点数量
# https://leetcode.cn/problems/maximum-number-of-intersections-on-the-chart/description/
# 有一条由 n 个点连接而成的折线图。给定一个 下标从 1 开始 的整数数组 y，第 k 个点的坐标是 (k, y[k])。
# 图中没有水平线，即没有两个相邻的点有相同的 y 坐标。
# 假设在图中任意画一条无限长的水平线。请返回这条水平线与折线相交的最多交点数。
# 2 <= y.length <= 1e5
# 1 <= y[i] <= 1e9
# 对于范围 [1, n - 1] 内的所有 i，都有 y[i] != y[i + 1]


# 对每根线段，在y上利用差分进行区间更新.

from collections import Counter
from itertools import pairwise
from typing import List, Tuple


def minmax(a: int, b: int) -> Tuple[int, int]:
    return (a, b) if a < b else (b, a)


class Solution:
    def maxIntersectionCount(self, y: List[int]) -> int:
        diff = Counter()
        for i, (a, b) in enumerate(pairwise(y)):
            start = 2 * a
            end = 2 * b + (0 if i == len(y) - 2 else -1 if b > a else 1)  # 重合需要-1(最后一条线段不会重合)
            min_, max_ = minmax(start, end)
            diff[min_] += 1
            diff[max_ + 1] -= 1
        res = 0
        curCount = 0
        for key in sorted(diff):
            curCount += diff[key]
            res = max(res, curCount)
        return res
