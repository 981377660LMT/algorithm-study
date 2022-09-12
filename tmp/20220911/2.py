from typing import List, Mapping, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)
Interval = Tuple[int, int]


# 请你将该字符串划分成一个或多个 子字符串 ，并满足每个子字符串中的字符都是 唯一 的。
# !也就是说，在单个子字符串中，字母的出现次数都不超过 一次
# 贪心
class Solution:
    def partitionString(self, s: str) -> int:
        visited = set()
        res = 1
        for char in s:
            if char in visited:
                res += 1
                visited = set([char])
            else:
                visited.add(char)
        return res
