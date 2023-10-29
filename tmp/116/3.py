from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个下标从 0 开始的整数数组 nums 和一个整数 target 。

# 返回和为 target 的 nums 子序列中，子序列 长度的最大值 。如果不存在和为 target 的子序列，返回 -1 。


# 子序列 指的是从原数组中删除一些或者不删除任何元素后，剩余元素保持原来的顺序构成的数组。


def max(a, b):
    return a if a > b else b


class Solution:
    def lengthOfLongestSubsequence(self, nums: List[int], target: int) -> int:
        dp = defaultdict(int, {0: 0})
        for num in nums:
            ndp = dp.copy()
            for k, v in dp.items():
                tmp = k + num
                if tmp <= target:
                    ndp[tmp] = max(ndp[tmp], v + 1)
            dp = ndp
        if target not in dp:
            return -1
        return dp[target]


# nums = [1,2,3,4,5], target = 9
print(Solution().lengthOfLongestSubsequence([1, 2, 3, 4, 5], 9))
