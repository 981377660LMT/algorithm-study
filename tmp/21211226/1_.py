from typing import List, Tuple
from collections import defaultdict, deque, Counter
from heapq import heapify, heappop, heappush
from sortedcontainers import SortedList
from bisect import bisect_left, bisect_right
from functools import lru_cache
from itertools import accumulate, groupby, combinations
from math import gcd


class Solution:
    def mostWordsFound(self, sentences: List[str]) -> int:
        return max([len(w.split(' ')) for w in sentences])


print(
    Solution().mostWordsFound(
        ["alice and bob love leetcode", "i think so too", "this is great thanks very much"]
    )
)
