from functools import lru_cache
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个整数 k 和一个整数 x 。

# 令 s 为整数 num 的下标从 1 开始的二进制表示。我们说一个整数 num 的 价值 是满足 i % x == 0 且 s[i] 是 设置位 的 i 的数目。

# 请你返回 最大 整数 num ，满足从 1 到 num 的所有整数的 价值 和小于等于 k 。

# 注意：


# 一个整数二进制表示下 设置位 是值为 1 的数位。
# 一个整数的二进制表示下标从右到左编号，比方说如果 s == 11100 ，那么 s[4] == 1 且 s[2] == 0 。
class Solution:
    def findMaximumNumber(self, k: int, x: int) -> int:
        def check(mid: int) -> bool:
            """满足从 1 到 mid 的所有整数的 价值 和小于等于 k 。"""

            def cal(upper: int, k: int) -> int:
                """[0,upper]中二进制第k(k>=0)位为1的数的个数.
                即满足 `num & (1 << k) > 0` 的数的个数
                """
                bit = upper.bit_length()
                if k >= bit:
                    return 0

                @lru_cache(None)
                def dfs(pos: int, hasLeadingZero: bool, isLimit: bool) -> int:
                    """当前在第pos位,hasLeadingZero表示有前导0，isLimit表示是否贴合上界"""
                    if pos == -1:
                        return not hasLeadingZero
                    if pos == k:
                        if isLimit and nums[pos] == 0:
                            return 0
                        return dfs(pos - 1, False, isLimit)
                    res = 0
                    up = nums[pos] if isLimit else 1
                    for cur in range(up + 1):
                        if hasLeadingZero and cur == 0:
                            res += dfs(pos - 1, True, (isLimit and cur == up))
                        else:
                            res += dfs(pos - 1, False, (isLimit and cur == up))
                    return res

                nums = list(map(int, bin(upper)[2:]))[::-1]
                res = dfs(len(nums) - 1, True, True)
                dfs.cache_clear()
                return res

            res = 0
            bitLen = mid.bit_length()
            for i in range(1, bitLen + 1):
                if i % x == 0:
                    res += cal(mid, i - 1)
            return res <= k

        left, right = 1, min(2**x * k, int(1e19))
        while left <= right:
            mid = (left + right) // 2
            if check(mid):
                left = mid + 1
            else:
                right = mid - 1
        return right
