from bisect import bisect_left, bisect_right
import math
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# LIS模板
def LIS(nums: List[int], isStrict=True) -> int:
    """求LIS长度"""
    n = len(nums)
    if n <= 1:
        return n

    res = [nums[0]]
    for i in range(1, n):
        pos = bisect_left(res, nums[i]) if isStrict else bisect_right(res, nums[i])
        if pos >= len(res):
            res.append(nums[i])
        else:
            res[pos] = nums[i]

    return len(res)


def caldp(nums: List[int], isStrict=True) -> List[int]:
    """求以每个位置为结尾的LIS长度(包括自身)"""
    if not nums:
        return []
    res = [1] * len(nums)
    LIS = [nums[0]]
    for i in range(1, len(nums)):
        if nums[i] > LIS[-1]:
            LIS.append(nums[i])
            res[i] = len(LIS)
        else:
            pos = bisect_left(LIS, nums[i]) if isStrict else bisect_right(LIS, nums[i])
            LIS[pos] = nums[i]
            res[i] = pos + 1
    return res


# dp[i] = max(nums[i], nums[i] + dp[x])  where   i-k <= x <= i-1


# The subsequence is strictly increasing and
# The difference between adjacent elements in the subsequence is at most k.
# Return the length of the longest subsequence that meets the requirements.

# 对每个数，求出以它为结尾的最长子序列长度
class Solution:
    def lengthOfLIS(self, nums: List[int], k: int) -> int:
        ...


print(Solution().lengthOfLIS(nums=[4, 2, 1, 4, 3, 4, 5, 8, 15], k=3))
