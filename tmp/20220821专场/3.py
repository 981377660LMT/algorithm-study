from math import gcd
from typing import List, Tuple, Optional
from collections import defaultdict, Counter


MOD = int(1e9 + 7)
INF = int(1e20)

# 二和三的倍数


class Solution:
    def minOperations(self, numbers: List[int]) -> int:
        n = len(numbers)
        gcd_ = gcd(*numbers)
        fac2 = [0] * n
        fac3 = [0] * n
        for i, num in enumerate(numbers):
            cur = num // gcd_
            while cur % 2 == 0:
                cur //= 2
                fac2[i] += 1
            while cur % 3 == 0:
                cur //= 3
                fac3[i] += 1
            if cur != 1:
                return -1

        max2, max3 = max(fac2), max(fac3)
        res1 = sum(max2 - i for i in fac2)
        res2 = sum(max3 - i for i in fac3)
        return res1 + res2


print(Solution().minOperations(numbers=[50, 75, 100]))
print(Solution().minOperations(numbers=[10, 14]))
