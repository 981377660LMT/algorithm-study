from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个字符串 s ，请你判断字符串 s 是否存在一个长度为 2 的子字符串，在其反转后的字符串中也出现。


# 如果存在这样的子字符串，返回 true；如果不存在，返回 false 。
class Solution:
    def isSubstringPresent(self, s: str) -> bool:
        if len(s) < 2:
            return False
        rev = s[::-1]
        for i in range(len(s) - 1):
            if s[i : i + 2] in rev:
                return True
        return False
