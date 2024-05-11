from functools import lru_cache
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个字符串 s ，你需要将它分割成一个或者更多的 平衡 子字符串。比方说，s == "ababcc" 那么 ("abab", "c", "c") ，("ab", "abc", "c") 和 ("ababcc") 都是合法分割，但是 ("a", "bab", "cc") ，("aba", "bc", "c") 和 ("ab", "abcc") 不是，不平衡的子字符串用粗体表示。

# 请你返回 s 最少 能分割成多少个平衡子字符串。


# 注意：一个 平衡 字符串指的是字符串中所有字符出现的次数都相同。
def min2(a: int, b: int) -> int:
    return a if a < b else b


class Solution:
    def minimumSubstringsInPartition(self, s: str) -> int:
        @lru_cache(None)
        def dfs(index: int) -> int:
            if index >= n:
                return 0
            res = INF
            counter = dict()
            for i in range(index, n):
                counter[s[i]] = counter.get(s[i], 0) + 1
                if all(v == counter[s[index]] for v in counter.values()):
                    res = min2(res, 1 + dfs(i + 1))
            return res

        n = len(s)
        res = dfs(0)
        dfs.cache_clear()
        return res
