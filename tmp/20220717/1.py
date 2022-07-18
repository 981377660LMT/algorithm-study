from typing import List, Tuple, Optional
from collections import defaultdict, Counter, deque

MOD = int(1e9 + 7)
INF = int(1e20)


class Solution:
    def numberOfPairs(self, nums: List[int]) -> List[int]:
        counter = Counter(nums)
        n = len(nums)
        res = sum(v // 2 for v in counter.values())
        return [res, n - res * 2]
