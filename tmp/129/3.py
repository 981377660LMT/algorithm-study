from functools import lru_cache
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你 3 个正整数 zero ，one 和 limit 。

# 一个 二进制数组 arr 如果满足以下条件，那么我们称它是 稳定的 ：

# 0 在 arr 中出现次数 恰好 为 zero 。
# 1 在 arr 中出现次数 恰好 为 one 。
# arr 中每个长度超过 limit 的 子数组 都 同时 包含 0 和 1 。
# 请你返回 稳定 二进制数组的 总 数目。


# 由于答案可能很大，将它对 109 + 7 取余 后返回。


def min2(a: int, b: int) -> int:
    return a if a < b else b


class Solution:
    def numberOfStableArrays(self, zero: int, one: int, limit: int) -> int:
        n = zero + one

        @lru_cache(None)
        def dfs(index: int, remainZero: int, pre: int, sameCount: int) -> int:
            if remainZero < 0:
                return 0
            if index == n:
                return 1 if remainZero == 0 else 0
            res = 0
            nextSameCount = sameCount + 1 if pre == 1 else 1
            if nextSameCount <= limit:
                res += dfs(index + 1, remainZero, 1, nextSameCount)
            nextSameCount = sameCount + 1 if pre == 0 else 1
            if nextSameCount <= limit:
                res += dfs(index + 1, remainZero - 1, 0, nextSameCount)
            return res % MOD

        res = dfs(0, zero, -1, 0)
        dfs.cache_clear()
        return res % MOD


# zero = 1, one = 1, limit = 2
# print(Solution().numberOfStableArrays(1, 1, 2))  # 1
print(Solution().numberOfStableArrays(1, 2, 1))  # 1
print(Solution().numberOfStableArrays(3, 3, 2))  # 1
