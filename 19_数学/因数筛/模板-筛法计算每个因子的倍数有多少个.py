# 直接从1到num开始筛法 复杂度nlogn


from collections import Counter
from typing import List

# 这类题的特点是nums[i]<=10^5 筛法枚举因子是nlogn


class Solution:
    def getMulti(self, nums: List[int]) -> Counter:
        """统计对于每个因子，原数组中有多少个他的倍数"""
        MAX = max(nums)
        counter = Counter(nums)
        multiCouner = Counter()
        for factor in range(1, MAX + 1):
            for multi in range(factor, MAX + 1, factor):
                multiCouner[factor] += counter[multi]
        return multiCouner


# [1,2,3,4,5]
# Counter({1: 5, 2: 2, 3: 1, 4: 1, 5: 1})
