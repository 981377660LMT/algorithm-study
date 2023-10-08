from functools import lru_cache
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你两个下标从 0 开始的二进制字符串 s1 和 s2 ，两个字符串的长度都是 n ，再给你一个正整数 x 。

# 你可以对字符串 s1 执行以下操作 任意次 ：

# 选择两个下标 i 和 j ，将 s1[i] 和 s1[j] 都反转，操作的代价为 x 。
# 选择满足 i < n - 1 的下标 i ，反转 s1[i] 和 s1[i + 1] ，操作的代价为 1 。
# 请你返回使字符串 s1 和 s2 相等的 最小 操作代价之和，如果无法让二者相等，返回 -1 。


# 注意 ，反转字符的意思是将 0 变成 1 ，或者 1 变成 0 。


def min(a, b):
    if a < b:
        return a
    return b


# !好题,难题，要意识到无后效性
# TODO: O(n)/费用流


class Solution:
    def minOperations(self, s1: str, s2: str, x: int) -> int:
        @lru_cache(None)
        def dfs(index: int, flip1Count: int, prevFlip2: bool) -> int:
            if index == n:
                return 0 if flip1Count == 0 else INF

            cur = nums1[index]
            if prevFlip2:
                cur = 1 - cur
            target = nums2[index]

            if cur == target:
                return dfs(index + 1, flip1Count, False)

            # 反转1
            res = dfs(index + 1, flip1Count + 1, False) + x
            if flip1Count > 0:
                res = min(res, dfs(index + 1, flip1Count - 1, False))
            # 反转2
            if index < n - 1:
                res = min(res, dfs(index + 1, flip1Count, True) + 1)
            return res

        nums1 = [int(i) for i in s1]
        nums2 = [int(i) for i in s2]
        n = len(s1)
        res = dfs(0, 0, False)
        dfs.cache_clear()
        return res if res < INF else -1


# "101101"
# "000000"
# 6
print(Solution().minOperations("101101", "000000", 6))
