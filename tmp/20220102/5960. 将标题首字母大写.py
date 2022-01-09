from typing import List, Tuple
from collections import defaultdict, deque, Counter
from heapq import heapify, heappop, heappush
from sortedcontainers import SortedList, SortedDict
from bisect import bisect_left, bisect_right
from functools import lru_cache, reduce
from itertools import accumulate, groupby, combinations, permutations, product, chain
from math import gcd, sqrt, ceil, floor, comb

MOD = int(1e9 + 7)
INF = 0x3F3F3F3F
EPS = int(1e-6)
dirs4 = [[-1, 0], [0, 1], [1, 0], [0, -1]]
dirs8 = [[-1, 0], [-1, 1], [0, 1], [1, 1], [1, 0], [1, -1], [0, -1], [-1, -1]]


class Solution:
    def capitalizeTitle(self, title: str) -> str:
        words = title.split()
        tmp = [w.capitalize() if len(w) >= 3 else w.lower() for w in words]
        return ' '.join(tmp)


print(Solution().capitalizeTitle("First leTTeR of EACH Word"))
print(Solution().capitalizeTitle("First leTTeR of EACH Word"))
print(Solution().capitalizeTitle(title="i lOve leetcode"))
print('assA'.capitalize())
print('assA'.lower())
print('assA'.upper())
