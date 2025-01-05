from math import gcd
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)


class Solution:
    def maxLength(self, nums: List[int]) -> int:
        n = len(nums)
        res = 1
        for i in range(n):
            cur_gcd = nums[i]
            cur_lcm = nums[i]
            cur_mul = nums[i]
            if cur_mul == cur_lcm * cur_gcd:
                res = max(res, 1)
            for j in range(i + 1, n):
                x = nums[j]
                cur_mul *= x
                cur_gcd = gcd(cur_gcd, x)
                cur_lcm = (cur_lcm * x) // gcd(cur_lcm, x)
                if cur_mul == cur_lcm * cur_gcd:
                    res = max(res, j - i + 1)
        return res
