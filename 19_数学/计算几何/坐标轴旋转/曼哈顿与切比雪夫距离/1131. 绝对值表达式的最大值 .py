# Return the maximum value for all 0 ≤ i < j < n:
# |arr1[i] - arr1[j]| + |arr2[i] - arr2[j]| + |i - j|
# !即 |arr1[i] - arr1[j]| + |arr2[i] - arr2[j]| - i + j
# 三维曼哈顿距离

# 枚举正负号 既然不知道到底正还是负 那就枚举一下
# abs(a) + abs(b) = max(a+b,a-b,-a+b,-a-b)
# !因为最后最大的那个肯定是答案
# 去绝对值+公式变形

from typing import List


INF = int(1e18)


class Solution:
    def maxAbsValExpr(self, arr1: List[int], arr2: List[int]) -> int:
        res = 0
        n = len(arr1)
        for dr, dc in [(-1, -1), (-1, 1), (1, -1), (1, 1)]:
            min_ = INF
            max_ = -INF
            for i in range(n):
                cand = dr * arr1[i] + dc * arr2[i] - i
                min_ = min(min_, cand)
                max_ = max(max_, cand)
            res = max(res, max_ - min_)
        return res
