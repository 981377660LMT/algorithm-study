from collections import defaultdict
from math import gcd
from typing import List


# 0 <= i < j <= n - 1 且
# nums[i] * nums[j] 能被 k 整除。

# !nums[i]<=10^5
# !美服的字典比较慢 最好用数组做counter

N = int(1e5 + 10)


class Solution:
    def countPairs(self, nums: List[int], k: int) -> int:
        counter, multiCounter = [0] * N, [0] * N
        max_ = max(nums)
        for num in nums:
            counter[num] += 1
        for factor in range(1, max_ + 1):
            for multi in range(factor, max_ + 1, factor):
                multiCounter[factor] += counter[multi]

        res = 0
        for num in nums:
            need = k // gcd(num, k)
            res += multiCounter[need]
            if num * num % k == 0:
                res -= 1

        return res // 2


print(Solution().countPairs(nums=[1, 2, 3, 4, 5], k=2))
