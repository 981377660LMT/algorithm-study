from typing import List, Tuple
from collections import defaultdict, deque, Counter
from heapq import heapify, heappop, heappush
from sortedcontainers import SortedList
from bisect import bisect_left, bisect_right
from functools import lru_cache
from itertools import accumulate, groupby, combinations
from math import gcd

MOD = int(1e9 + 7)
INF = 0x7FFFFFFF


class Solution:
    def checkString(self, s: str) -> bool:
        return 'ba' not in s


print(Solution().checkString("aaabbb"))
print(Solution().checkString("abab"))
print(Solution().checkString("bbb"))
