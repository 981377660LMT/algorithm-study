from math import comb
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)


# 给你一个字符串 s 和一个字符 c 。返回在字符串 s 中并且以 c 字符开头和结尾的非空子字符串的总数。
class Solution:
    def countSubstrings(self, s: str, c: str) -> int:
        mp = defaultdict(list)
        for i, ch in enumerate(s):
            mp[ch].append(i)
        indexes = mp[c]
        return len(indexes) * (len(indexes) + 1) // 2
