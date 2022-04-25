# 乘积不超过k的子数组个数
# 和不超过k的子数组个数
from typing import List


def noMoreThan1(nums: List[int], k: int) -> int:
    """子数组的和不大于k的子数组个数"""
    res = 0
    left = 0
    curSum = 0
    for right, num in enumerate(nums):
        curSum += num
        while left <= right and curSum > k:
            curSum -= nums[left]
            left += 1
        if left <= right:
            res += right - left + 1
    return res


def noMoreThan2(nums: List[int], k: int) -> int:
    """子数组的积不大于k的子数组个数"""
    res = 0
    left = 0
    curMul = 1
    for right, num in enumerate(nums):
        curMul *= num
        while left <= right and curMul > k:
            curMul //= nums[left]
            left += 1
        if left <= right:
            res += right - left + 1
    return res


print(noMoreThan1([1, 2, 3], 4))
print(noMoreThan2([1, 2, 3], 4))
