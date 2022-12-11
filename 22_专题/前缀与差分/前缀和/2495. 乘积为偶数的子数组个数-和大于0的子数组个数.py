# 2495. 乘积为偶数的子数组个数-和大于0的子数组个数


from typing import List
from sortedcontainers import SortedList


def subarraysWithSumMoreThanZero(nums: List[int]) -> int:
    """和大于0的子数组个数"""
    res, curSum, sl = 0, 0, SortedList([0])
    for num in nums:
        curSum += num
        res += sl.bisect_left(curSum)
        sl.add(curSum)
    return res


class Solution:
    def evenProduct(self, nums: List[int]) -> int:
        fac2 = [0] * len(nums)
        for i, num in enumerate(nums):
            while num % 2 == 0:
                fac2[i] += 1
                num //= 2
        return subarraysWithSumMoreThanZero(fac2)


print(Solution().evenProduct(nums=[9, 6, 7, 13]))
