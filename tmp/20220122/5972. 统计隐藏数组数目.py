from typing import List, Tuple
from collections import defaultdict, deque, Counter
from heapq import heapify, heappop, heappush
from sortedcontainers import SortedList, SortedDict, SortedSet
from bisect import bisect_left, bisect_right
from functools import lru_cache, reduce
from itertools import accumulate, groupby, combinations, permutations, product, chain, islice
from math import gcd, sqrt, ceil, floor, comb
from string import ascii_lowercase, ascii_uppercase, ascii_letters, digits
from operator import xor, or_, and_, not_

MOD = int(1e9 + 7)
INF = 0x3F3F3F3F
EPS = int(1e-8)
dirs4 = [[-1, 0], [0, 1], [1, 0], [0, -1]]
dirs8 = [[-1, 0], [-1, 1], [0, 1], [1, 1], [1, 0], [1, -1], [0, -1], [-1, -1]]

# 讨论数组第一个元素的可行的取值范围
class Solution:
    def numberOfArrays(self, differences: List[int], lower: int, upper: int) -> int:
        preSum = list(accumulate(differences))
        minVal, maxVal = min(preSum), max(preSum)
        minStart = max(lower, lower - minVal)
        maxStart = min(upper, upper - maxVal)
        return max(0, maxStart - minStart + 1)


print(Solution().numberOfArrays(differences=[1, -3, 4], lower=1, upper=6))
print(Solution().numberOfArrays(differences=[3, -4, 5, 1, -2], lower=-4, upper=5))
print(Solution().numberOfArrays(differences=[4, -7, 2], lower=3, upper=6))
print(Solution().numberOfArrays(differences=[-40], lower=-46, upper=53))
print(
    Solution().numberOfArrays(differences=[57587, 47629, -7859, 32782], lower=-63763, upper=35288)
)
