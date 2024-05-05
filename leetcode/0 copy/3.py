from math import gcd
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个字符串 s ，它由某个字符串 t 和它的 同位字符串 连接而成。

# 请你返回字符串 t 的 最小 可能长度。


# 同位字符串 指的是重新排列一个单词得到的另外一个字符串，原来字符串中的每个字符在新字符串中都恰好只使用一次。
class Solution:
    def minAnagramLength(self, s: str) -> int:
        counter = Counter(s)
        freqSet = set(counter.values())
        gcd_ = gcd(*freqSet)
        return len(s) // gcd_


# "oionssonoi"

print(Solution().minAnagramLength("oionssonoi"))
