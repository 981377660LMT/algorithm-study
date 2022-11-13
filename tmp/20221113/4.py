from functools import lru_cache
import time
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个字符串 s 和一个 正 整数 k 。

# 从字符串 s 中选出一组满足下述条件且 不重叠 的子字符串：

# 每个子字符串的长度 至少 为 k 。
# 每个子字符串是一个 回文串 。
# 返回最优方案中能选择的子字符串的 最大 数目。


class Solution:
    def maxPalindromes(self, s: str, k: int) -> int:
        def helper(left, right):
            """中心扩展法求所有回文子串"""
            while left >= 0 and right < len(s) and s[left] == s[right]:
                left -= 1
                right += 1
                if right - left - 1 >= k:
                    intervals.append((left + 1, right - 1))

        n = len(s)
        intervals = []
        for i in range(n):
            helper(i, i)
            helper(i, i + 1)

        intervals.sort(key=lambda x: x[1])
        res = 0
        preEnd = -1
        for start, end in intervals:
            if start > preEnd:
                res += 1
                preEnd = end
        return res


print(Solution().maxPalindromes(s="abaccdbbd", k=3))
print(Solution().maxPalindromes(s="iqqibcecvrbxxj", k=1))
print(Solution().maxPalindromes(s="i" * 2000, k=1))
