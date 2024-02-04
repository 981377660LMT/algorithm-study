from itertools import accumulate
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个长度为 n 的数组 nums 和一个 正 整数 k 。

# 如果 nums 的一个子数组中，第一个元素和最后一个元素 差的绝对值恰好 为 k ，我们称这个子数组为 好 的。换句话说，如果子数组 nums[i..j] 满足 |nums[i] - nums[j]| == k ，那么它是一个好子数组。


# 请你返回 nums 中 好 子数组的 最大 和，如果没有好子数组，返回 0 。
def max2(a: int, b: int) -> int:
    return a if a > b else b


def min2(a: int, b: int) -> int:
    return a if a < b else b


class Solution:
    def maximumSubarraySum(self, nums: List[int], k: int) -> int:
        preMin = defaultdict(lambda: INF)
        res, curSum = -INF, 0
        for i, num in enumerate(nums):
            curSum += num
            if num - k in preMin:
                res = max2(res, curSum - preMin[num - k])
            if num + k in preMin:
                res = max2(res, curSum - preMin[num + k])
            preMin[num] = min2(preMin[num], curSum - num)

        return res if res != -INF else 0
