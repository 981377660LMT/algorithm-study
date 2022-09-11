from typing import List, Mapping, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)
Interval = Tuple[int, int]


# 请你将该字符串划分成一个或多个 子字符串 ，并满足每个子字符串中的字符都是 唯一 的。
# 也就是说，在单个子字符串中，字母的出现次数都不超过 一次
class Solution:
    def partitionString(self, s: str) -> int:
        counter = Counter()
        res = 1
        for c in s:
            counter[c] += 1
            if counter[c] > 1:
                res += 1
                counter.clear()
                counter[c] += 1
        return res
