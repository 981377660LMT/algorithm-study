from typing import List, Tuple
from collections import defaultdict, deque, Counter
from heapq import heapify, heappop, heappush
from sortedcontainers import SortedList, SortedDict
from bisect import bisect_left, bisect_right
from functools import lru_cache, reduce
from itertools import accumulate, groupby, combinations, permutations, product
from math import gcd, sqrt, ceil, floor, comb

MOD = int(1e9 + 7)
INF = 0x3F3F3F3F
EPS = int(1e-6)
dirs4 = [[-1, 0], [0, 1], [1, 0], [0, -1]]
dirs8 = [[-1, 0], [-1, 1], [0, 1], [1, 1], [1, 0], [1, -1], [0, -1], [-1, -1]]


class Solution:
    def numberOfBeams(self, bank: List[str]) -> int:
        nums = [c for row in bank if (c := row.count('1')) != 0]
        return sum(a * b for a, b in zip(nums, nums[1:]))


print(Solution().numberOfBeams(bank=["011001", "000000", "010100", "001000"]))
print(Solution().numberOfBeams(bank=["000", "111", "000"]))
