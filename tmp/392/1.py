from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个整数数组 nums 。


# 返回数组 nums 中 严格递增 或 严格递减 的最长非空子数组的长度。
class Solution:
    def longestMonotonicSubarray(self, nums: List[int]) -> int:
        res = 1
        n = len(nums)
        for i in range(n):
            for j in range(i + 1, n):
                cur = nums[i : j + 1]
                if all(cur[k] < cur[k + 1] for k in range(len(cur) - 1)) or all(
                    cur[k] > cur[k + 1] for k in range(len(cur) - 1)
                ):
                    res = max(res, len(cur))
        return res
