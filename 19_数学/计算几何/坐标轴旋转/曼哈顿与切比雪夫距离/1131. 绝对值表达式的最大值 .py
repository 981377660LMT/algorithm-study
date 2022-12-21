# Return the maximum value for all 0 ≤ i < j < n:
# |arr1[i] - arr1[j]| + |arr2[i] - arr2[j]| + |i - j|
# !即 |arr1[i] - arr1[j]| + |arr2[i] - arr2[j]| - i + j
# 三维曼哈顿距离

# 枚举正负号 既然不知道到底正还是负 那就枚举一下
# abs(a) + abs(b) = max(a+b,a-b,-a+b,-a-b)
# !因为最后最大的那个肯定是答案
# 去绝对值+公式变形

from itertools import product
from typing import List


INF = int(1e18)


class Solution:
    def maxAbsValExpr(self, arr1: List[int], arr2: List[int]) -> int:
        res = 0
        n = len(arr1)
        for (d1, d2, d3) in product([-1, 1], repeat=3):
            min_, max_ = INF, -INF
            for i in range(n):
                cand = d1 * arr1[i] + d2 * arr2[i] - i * d3
                min_ = min(min_, cand)
                max_ = max(max_, cand)
            res = max(res, max_ - min_)
        return res


assert Solution().maxAbsValExpr(arr1=[1, 2, 3, 4], arr2=[-1, 4, 5, 6]) == 13
