from math import comb
from typing import List, Tuple, Optional
from collections import defaultdict, Counter


MOD = int(1e9 + 7)
INF = int(1e20)


class Solution:
    def countBadPairs(self, nums: List[int]) -> int:
        n = len(nums)
        res = comb(n, 2)
        arr = [num - i for i, num in enumerate(nums)]
        counter = Counter(arr)
        for v in counter.values():
            res -= comb(v, 2)
        return res
