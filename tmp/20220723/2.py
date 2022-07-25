from itertools import groupby
from typing import List


MOD = int(1e9 + 7)
INF = int(1e20)


class Solution:
    def zeroFilledSubarray(self, nums: List[int]) -> int:
        """全 0 子数组的数目"""
        groups = [(char, len(list(group))) for char, group in groupby(nums)]
        return sum(count * (count + 1) // 2 for char, count in groups if char == 0)
