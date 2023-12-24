from itertools import accumulate
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个长度为 n 的 正 整数数组 nums 。

# 多边形 指的是一个至少有 3 条边的封闭二维图形。多边形的 最长边 一定 小于 所有其他边长度之和。

# 如果你有 k （k >= 3）个 正 数 a1，a2，a3, ...，ak 满足 a1 <= a2 <= a3 <= ... <= ak 且 a1 + a2 + a3 + ... + ak-1 > ak ，那么 一定 存在一个 k 条边的多边形，每条边的长度分别为 a1 ，a2 ，a3 ， ...，ak 。

# 一个多边形的 周长 指的是它所有边之和。


# 请你返回从 nums 中可以构造的 多边形 的 最大周长 。如果不能构造出任何多边形，请你返回 -1 。
class Solution:
    def largestPerimeter(self, nums: List[int]) -> int:
        nums.sort()
        # 枚举最长边
        preSum = [0] + list(accumulate(nums))
        for i in range(len(nums) - 1, 1, -1):
            if nums[i] < preSum[i + 1] - nums[i]:
                return preSum[i + 1]
        return -1
