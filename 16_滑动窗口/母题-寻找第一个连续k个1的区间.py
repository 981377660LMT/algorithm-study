# 母题-寻找第一个连续k个1的区间
# 1. dp
# 2. 滑动窗口
from typing import List


def indexOfKOnes1(nums: List[int], k: int) -> int:
    """dp寻找第一个连续k个1的区间起点,不存在则返回-1"""
    dp = 0
    for i, num in enumerate(nums):
        if num == 1:
            dp += 1
            if dp == k:
                return i - k + 1
        else:
            dp = 0
    return -1


def indexOfKOnes2(nums: List[int], k: int) -> int:
    """定长滑窗寻找第一个连续k个1的区间起点,不存在则返回-1"""
    n = len(nums)
    curSum = 0
    for right in range(n):
        curSum += nums[right] == 1
        if right >= k:
            curSum -= nums[right - k] == 1
        if right >= k - 1:
            if curSum == k:
                return right - k + 1
    return -1


assert indexOfKOnes1(nums=[1, 2, 1, 1, 3], k=2) == 2
assert indexOfKOnes2(nums=[1, 2, 1, 1, 3], k=2) == 2
