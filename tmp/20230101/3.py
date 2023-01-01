from functools import lru_cache
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)
# 给你一个字符串 s ，它每一位都是 1 到 9 之间的数字组成，同时给你一个整数 k 。

# 如果一个字符串 s 的分割满足以下条件，我们称它是一个 好 分割：

# s 中每个数位 恰好 属于一个子字符串。
# 每个子字符串的值都小于等于 k 。
# 请你返回 s 所有的 好 分割中，子字符串的 最少 数目。如果不存在 s 的 好 分割，返回 -1 。

# 注意：

# 一个字符串的 值 是这个字符串对应的整数。比方说，"123" 的值为 123 ，"1" 的值是 1 。
# 子字符串 是字符串中一段连续的字符序列。


class Solution:
    def minimumPartition(self, s: str, k: int) -> int:
        @lru_cache(None)
        def dfs(index: int) -> int:
            if index == len(s):
                return 0
            res = INF
            cur = 0
            for i in range(index, len(s)):
                cur = cur * 10 + int(s[i])
                if cur > k:
                    break
                res = min(res, dfs(i + 1) + 1)
            return res

        res = dfs(0)
        dfs.cache_clear()
        return res if res != INF else -1
