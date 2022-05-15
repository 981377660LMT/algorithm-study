from itertools import accumulate
from typing import List, Tuple


MOD = int(1e9 + 7)
INF = int(1e20)


class Solution:
    def waysToSplitArray(self, nums: List[int]) -> int:
        preSum = [0] + list(accumulate(nums))
        sum_ = preSum[-1]
        res = 0
        n = len(nums)
        for i in range(n - 1):
            cur = preSum[i + 1]
            if cur >= sum_ - cur:
                res += 1
        return res

