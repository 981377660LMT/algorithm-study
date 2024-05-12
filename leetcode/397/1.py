from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你两个字符串 s 和 t，每个字符串中的字符都不重复，且 t 是 s 的一个排列。

# 排列差 定义为 s 和 t 中每个字符在两个字符串中位置的绝对差值之和。


# 返回 s 和 t 之间的 排列差 。
class Solution:
    def findPermutationDifference(self, s: str, t: str) -> int:
        mp1, mp2 = defaultdict(int), defaultdict(int)
        for i, c in enumerate(s):
            mp1[c] = i
        for i, c in enumerate(t):
            mp2[c] = i
        return sum(abs(mp1[c] - mp2[c]) for c in s)
