from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# 给你两个数字字符串 num1 和 num2 ，以及两个整数 max_sum 和 min_sum 。如果一个整数 x 满足以下条件，我们称它是一个好整数：

# num1 <= x <= num2
# min_sum <= digit_sum(x) <= max_sum.
# 请你返回好整数的数目。答案可能很大，请返回答案对 109 + 7 取余后的结果。


# 注意，digit_sum(x) 表示 x 各位数字之和。

from functools import lru_cache


def cal(upper: int, minSum: int, maxSum: int) -> int:
    @lru_cache(None)
    def dfs(pos: int, hasLeadingZero: int, isLimit: bool, curSum: int) -> int:
        """当前在第pos位,hasLeadingZero表示有前导0,isLimit表示是否贴合上界"""
        if curSum > maxSum:
            return 0
        if pos == len(nums):
            return 1 if curSum >= minSum else 0

        res = 0
        up = nums[pos] if isLimit else 9
        for cur in range(up + 1):
            if hasLeadingZero and cur == 0:
                res += dfs(pos + 1, True, (isLimit and cur == up), curSum)
            else:
                res += dfs(pos + 1, False, (isLimit and cur == up), curSum + cur)
        return res % MOD

    nums = list(map(int, str(upper)))
    return dfs(0, True, True, 0) % MOD


class Solution:
    def count(self, num1: str, num2: str, min_sum: int, max_sum: int) -> int:
        return (cal(int(num2), min_sum, max_sum) - cal(int(num1) - 1, min_sum, max_sum)) % MOD
