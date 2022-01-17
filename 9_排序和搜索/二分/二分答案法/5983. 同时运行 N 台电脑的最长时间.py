from typing import List, Tuple
from collections import defaultdict, deque, Counter
from heapq import heapify, heappop, heappush
from sortedcontainers import SortedList, SortedDict, SortedSet
from bisect import bisect_left, bisect_right
from functools import lru_cache, reduce
from itertools import accumulate, groupby, combinations, permutations, product, chain
from math import gcd, sqrt, ceil, floor, comb
from string import ascii_lowercase, ascii_uppercase, ascii_letters, digits
from operator import xor, or_, and_, not_

MOD = int(1e9 + 7)
INF = 0x3F3F3F3F
EPS = int(1e-8)
dirs4 = [[-1, 0], [0, 1], [1, 0], [0, -1]]
dirs8 = [[-1, 0], [-1, 1], [0, 1], [1, 1], [1, 0], [1, -1], [0, -1], [-1, -1]]

# 1 <= n <= batteries.length <= 105
# 1 <= batteries[i] <= 10^9
# 优先队列模拟不太好做，因为 1 <= batteries[i] <= 10^9


class Solution:
    def maxRunTime(self, n: int, batteries: List[int]) -> int:
        batteries = sorted(batteries, reverse=True)
        # 备用电池
        spareSum = sum(batteries[n:])

        def check(needTime: int) -> bool:
            """"备用电池能否提供足够的储备"""
            needSpare = 0
            for i in range(n):
                supply = batteries[i]
                if supply < needTime:
                    needSpare += needTime - supply
            return needSpare <= spareSum

        left, right = 0, sum(batteries)
        while left <= right:
            mid = (left + right) >> 1
            if check(mid):
                left = mid + 1
            else:
                right = mid - 1
        return right


print(Solution().maxRunTime(n=2, batteries=[3, 3, 3]))
print(Solution().maxRunTime(n=2, batteries=[1, 1, 1, 1]))
print(Solution().maxRunTime(n=2, batteries=[1, 1, 1, 1]))

