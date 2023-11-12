from functools import lru_cache
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)


# 给你一个整数 n 。

# 如果一个字符串 s 只包含小写英文字母，且 将 s 的字符重新排列后，新字符串包含 子字符串 "leet" ，那么我们称字符串 s 是一个 好 字符串。

# 比方说：

# 字符串 "lteer" 是好字符串，因为重新排列后可以得到 "leetr" 。
# "letl" 不是好字符串，因为无法重新排列并得到子字符串 "leet" 。
# 请你返回长度为 n 的好字符串 总 数目。

# 由于答案可能很大，将答案对 109 + 7 取余 后返回。

# 子字符串 是一个字符串中一段连续的字符序列。


def min(a, b):
    return a if a < b else b


class Solution:
    def stringCount(self, n: int) -> int:
        if n <= 3:
            return 0

        @lru_cache(None)
        def dfs(index: int, a: int, b: int, c: int) -> int:
            if index == n:
                return a >= 1 and b >= 2 and c >= 1
            res = 0
            for i in range(26):
                if i == 0:
                    res += dfs(index + 1, min(a + 1, 1), b, c)
                elif i == 1:
                    res += dfs(index + 1, a, min(b + 1, 2), c)
                elif i == 2:
                    res += dfs(index + 1, a, b, min(c + 1, 1))
                else:
                    res += dfs(index + 1, a, b, c)
                res %= MOD
            return res

        res = dfs(0, 0, 0, 0)
        dfs.cache_clear()
        return res % MOD
