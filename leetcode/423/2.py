from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)


def max2(a: int, b: int) -> int:
    return a if a > b else b


def min2(a: int, b: int) -> int:
    return a if a < b else b


class Solution:
    def maxIncreasingSubarrays(self, nums: List[int]) -> int:
        def getGroups(nums: List[int]) -> List[int]:
            res = []
            ptr = 0
            for i in range(1, n + 1):
                if i == n or nums[i] <= nums[i - 1]:
                    len_ = i - ptr
                    res.append(len_)
                    ptr = i
            return res

        n = len(nums)
        groups = getGroups(nums)
        res = 0
        for i in range(len(groups) - 1):
            res = max2(res, min2(groups[i], groups[i + 1]))
        for v in groups:
            res = max2(res, v // 2)
        return res
