from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个字符串 word。如果 word 中同时出现某个字母 c 的小写形式和大写形式，并且 每个 小写形式的 c 都出现在第一个大写形式的 c 之前，则称字母 c 是一个 特殊字母 。


# 返回 word 中 特殊字母 的数量。
class Solution:
    def numberOfSpecialChars(self, word: str) -> int:
        s = set(word)
        first, last = dict(), dict()
        for i, c in enumerate(word):
            if c not in first:
                first[c] = i
            last[c] = i
        res = 0
        for c in s:
            if c.lower() in s and c.upper() in s and last[c.lower()] < first[c.upper()]:
                res += 1
        return res // 2
