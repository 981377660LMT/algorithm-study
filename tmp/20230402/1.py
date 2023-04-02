from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)
# 给你一个仅由 0 和 1 组成的二进制字符串 s 。

# 如果子字符串中 所有的 0 都在 1 之前 且其中 0 的数量等于 1 的数量，则认为 s 的这个子字符串是平衡子字符串。请注意，空子字符串也视作平衡子字符串。

# 返回  s 中最长的平衡子字符串长度。

# 子字符串是字符串中的一个连续字符序列。


class Solution:
    def findTheLongestBalancedSubstring(self, s: str) -> int:
        res = 0
        for start in range(len(s)):
            for end in range(start, len(s)):
                cur = s[start : end + 1]
                if cur.count("0") == cur.count("1") and "10" not in cur:
                    res = max(res, end - start + 1)
        return res
