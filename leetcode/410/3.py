from functools import lru_cache
from itertools import accumulate
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个长度为 n 的 正 整数数组 nums 。

# 如果两个 非负 整数数组 (arr1, arr2) 满足以下条件，我们称它们是 单调 数组对：

# 两个数组的长度都是 n 。
# arr1 是单调 非递减 的，换句话说 arr1[0] <= arr1[1] <= ... <= arr1[n - 1] 。
# arr2 是单调 非递增 的，换句话说 arr2[0] >= arr2[1] >= ... >= arr2[n - 1] 。
# 对于所有的 0 <= i <= n - 1 都有 arr1[i] + arr2[i] == nums[i] 。
# 请你返回所有 单调 数组对的数目。


# 由于答案可能很大，请你将它对 109 + 7 取余 后返回。


def max2(a: int, b: int) -> int:
    return a if a > b else b


def min2(a: int, b: int) -> int:
    return a if a < b else b


# dp[i][j] 表示前 i 个数中，arr1 的最后一个数是 j 的方案数
class Solution:
    def countOfPairs(self, nums: List[int]) -> int:
        max_ = max(nums)
        dp = [0] * (max_ + 1)
        dp[: nums[0] + 1] = [1] * (nums[0] + 1)
        ndp = [0] * (max_ + 1)
        for pre, cur in zip(nums, nums[1:]):
            ndp[:] = [0] * (max_ + 1)
            dpPreSum = list(accumulate(dp, initial=0))
            for v1 in range(cur + 1):
                v2 = cur - v1
                upper = max2(min2(v1, pre - v2) + 1, 0)
                ndp[v1] = dpPreSum[upper] % MOD
            dp, ndp = ndp, dp
        return sum(dp) % MOD


# nums = [2,3,2]

print(Solution().countOfPairs([2, 3, 2]))  # 4
# print(Solution().countOfPairs([5, 5, 5, 5]))  # 2
