from math import gcd
from typing import Counter, List

# 0 <= i < j <= n - 1 且
# nums[i] * nums[j] 能被 k 整除。


class Solution:
    def coutPairs(self, nums: List[int], k: int) -> int:
        MAX = max(nums)
        counter = Counter(nums)
        multiCouner = Counter()
        for factor in range(1, MAX + 1):
            for multi in range(factor, MAX + 1, factor):
                multiCouner[factor] += counter[multi]

        res = 0
        for num in nums:
            need = k // gcd(num, k)
            res += multiCouner[need]
            if num * num % k == 0:
                res -= 1
        return res // 2


print(Solution().coutPairs(nums=[1, 2, 3, 4, 5], k=2))
