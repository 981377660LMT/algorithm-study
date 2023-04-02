# 返回满足 0 <= i < j < n ，nums[i] == nums[j] 且 (i * j) 能被 k 整除的数对 (i, j) 的 数目 。
# 如果1 <= nums.length <= 100 暴力
# 如果1 <= nums.length <= 10^5 筛法

from collections import defaultdict
from math import gcd
from typing import Counter, List


class Solution:
    def countPairs(self, nums: List[int], k: int) -> int:
        n = len(nums)
        multiCounter = defaultdict(lambda: defaultdict(int))
        for factor in range(1, n):
            for multi in range(factor, n, factor):
                multiCounter[factor][nums[multi]] += 1

        counter = Counter(nums)
        res1, res2 = 0, 0
        for index, value in enumerate(nums):
            if index == 0:
                res1 += counter[value] - 1
            else:
                need = k // gcd(index, k)
                res2 += multiCounter[need][value]
                if index**2 % k == 0:
                    res2 -= 1

        return res1 + res2 // 2


print(Solution().countPairs(nums=[3, 1, 2, 2, 2, 1, 3], k=2))
print(Solution().countPairs(nums=[1, 2, 3, 4], k=1))
print(Solution().countPairs(nums=[5, 5, 9, 2, 5, 5, 9, 2, 2, 5, 5, 6, 2, 2, 5, 2, 5, 4, 3], k=7))
# 18
