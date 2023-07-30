from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# 给你两个正整数 low 和 high ，都用字符串表示，请你统计闭区间 [low, high] 内的 步进数字 数目。

# 如果一个整数相邻数位之间差的绝对值都 恰好 是 1 ，那么这个数字被称为 步进数字 。

# 请你返回一个整数，表示闭区间 [low, high] 之间步进数字的数目。

# 由于答案可能很大，请你将它对 109 + 7 取余 后返回。


# 注意：步进数字不能有前导 0 。

from functools import lru_cache


def cal(upper: int) -> int:
    @lru_cache(None)
    def dfs(pos: int, hasLeadingZero: bool, isLimit: bool, pre: int) -> int:
        """当前在第pos位，hasLeadingZero表示有前导0，isLimit表示是否贴合上界"""
        if pos == len(nums):
            return int(not hasLeadingZero)

        res = 0
        up = nums[pos] if isLimit else 9
        for cur in range(up + 1):
            if pre == -1 or abs(cur - pre) == 1:
                if hasLeadingZero and cur == 0:
                    res += dfs(pos + 1, True, (isLimit and cur == up), -1)
                else:
                    res += dfs(pos + 1, False, (isLimit and cur == up), cur)
        return res % MOD

    nums = list(map(int, str(upper)))
    return dfs(0, True, True, -1) % MOD


class Solution:
    def countSteppingNumbers(self, low: str, high: str) -> int:
        return (cal(int(high)) - cal(int(low) - 1)) % MOD


print(Solution().countSteppingNumbers(low="1", high="11"))
