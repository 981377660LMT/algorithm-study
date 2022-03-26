from typing import List
from functools import lru_cache


MOD = int(1e9 + 7)


@lru_cache(None)
def pow2(x: int) -> int:
    return pow(2, x, MOD)


class Solution:
    def solve(self, nums: List[int]) -> int:
        n = len(nums)
        nums = sorted(nums)
        res = 0

        for i in range(n):
            pos = (pow2(i) - 1) * nums[i]
            neg = (pow2(n - i - 1) - 1) * nums[i]
            res += pos - neg
            res %= MOD

        return res

