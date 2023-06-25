from functools import lru_cache
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个下标从 0 开始的数组 words ，它包含 n 个字符串。

# 定义 连接 操作 join(x, y) 表示将字符串 x 和 y 连在一起，得到 xy 。如果 x 的最后一个字符与 y 的第一个字符相等，连接后两个字符中的一个会被 删除 。

# 比方说 join("ab", "ba") = "aba" ， join("ab", "cde") = "abcde" 。

# 你需要执行 n - 1 次 连接 操作。令 str0 = words[0] ，从 i = 1 直到 i = n - 1 ，对于第 i 个操作，你可以执行以下操作之一：

# 令 stri = join(stri - 1, words[i])
# 令 stri = join(words[i], stri - 1)
# 你的任务是使 strn - 1 的长度 最小 。


# 请你返回一个整数，表示 strn - 1 的最小长度。


def max(a, b):
    if a > b:
        return a
    return b


class Solution:
    def minimizeConcatenatedLength(self, words: List[str]) -> int:
        if len(words) == 1:
            return len(words[0])

        @lru_cache(None)
        def dfs(index: int, first: str, last: str) -> int:
            """消除的最多字符数"""
            if index >= n - 1:
                return 0
            next_ = words[index + 1]
            res1 = dfs(index + 1, first, next_[-1]) + (last == next_[0])
            res2 = dfs(index + 1, next_[0], last) + (first == next_[-1])
            return max(res1, res2)

        n = len(words)
        res = dfs(0, words[0][0], words[0][-1])
        dfs.cache_clear()
        return sum(len(word) for word in words) - res
