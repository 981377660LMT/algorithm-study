from typing import List, Tuple
from collections import defaultdict, deque, Counter
from heapq import heapify, heappop, heappush
from sortedcontainers import SortedList, SortedDict
from bisect import bisect_left, bisect_right
from functools import lru_cache, reduce
from itertools import accumulate, groupby, combinations, permutations, product
from math import gcd, sqrt, ceil, floor, comb
from string import ascii_lowercase, ascii_uppercase, ascii_letters, digits


MOD = int(1e9 + 7)
INF = 0x3F3F3F3F
EPS = int(1e-6)
dirs4 = [[-1, 0], [0, 1], [1, 0], [0, -1]]
dirs8 = [[-1, 0], [-1, 1], [0, 1], [1, 1], [1, 0], [1, -1], [0, -1], [-1, -1]]


class Solution:
    def minSwaps(self, nums: List[int]) -> int:
        winLen = nums.count(1)
        nums = nums * 2
        res = INF
        curOne = 0
        left = 0
        for right, cur in enumerate(nums):
            if cur == 1:
                curOne += 1
            if right >= winLen:
                if nums[left] == 1:
                    curOne -= 1
                left += 1

            if right >= winLen - 1:
                res = min(res, winLen - curOne)
        return res


print(Solution().minSwaps(nums=[0, 1, 0, 1, 1, 0, 0]))
print(Solution().minSwaps(nums=[0, 1, 1, 1, 0, 0, 1, 1, 0]))
print(Solution().minSwaps(nums=[1, 1, 0, 0, 1]))
