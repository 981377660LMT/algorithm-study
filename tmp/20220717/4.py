from math import gcd
from typing import List, Tuple, Optional
from collections import defaultdict, Counter, deque

MOD = int(1e9 + 7)
INF = int(1e20)


class Solution:
    def minOperations(self, nums: List[int], numsDivide: List[int]) -> int:
        gcd_ = gcd(*numsDivide)
        counter = Counter(nums)
        res = 0
        keys = sorted(counter)
        for key in keys:
            if gcd_ % key == 0:
                return res
            res += counter[key]
        return -1


print(Solution().minOperations(nums=[2, 3, 2, 4, 3], numsDivide=[9, 6, 9, 3, 15]))
