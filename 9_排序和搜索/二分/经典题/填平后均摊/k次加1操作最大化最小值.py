# k次加1操作最大化一个数组的最小值  二分的check函数O(1)
# 花园模型.

from typing import List
from itertools import accumulate


def maximizeMinValue(nums: List[int], k: int) -> int:
    """k次加1操作,让最小值最大化,返回最小值"""
    n = len(nums)
    nums = sorted(nums)
    preSum = [0] + list(accumulate(nums))
    nums = [0] + nums

    # !最右二分求最后能和哪个数齐平
    left, right = 0, n
    while left <= right:
        mid = (left + right) // 2
        diff = mid * nums[mid] - preSum[mid]
        if diff <= k:
            left = mid + 1
        else:
            right = mid - 1

    max_ = nums[right]
    overflow = k - (right * nums[right] - preSum[right])
    max_ += overflow // right if right else 0
    return max_


def maximizeMinValue2(nums: List[int], k: int) -> List[int]:
    """k次加1操作,让最小值最大化,返回操作后的数组"""
    n = len(nums)
    copy = nums[:]
    nums = sorted(nums)
    preSum = [0] + list(accumulate(nums))
    nums = [0] + nums

    # !最右二分求最后能和哪个数齐平
    left, right = 0, n
    while left <= right:
        mid = (left + right) // 2
        diff = mid * nums[mid] - preSum[mid]
        if diff <= k:
            left = mid + 1
        else:
            right = mid - 1

    max_ = nums[right]
    overflow = k - (right * nums[right] - preSum[right])
    div, mod = 0, 0
    if right:
        div, mod = divmod(overflow, right)  # mod个数需要再加1
    max_ += div

    for i in range(n):
        if copy[i] < max_ + int(mod > 0):
            copy[i] = max_ + int(mod > 0)
            mod -= 1

    return copy


def minimizeMaxValue(nums: List[int], k: int) -> List[int]:
    """k次-1操作,让最大值最小化,返回操作后的数组"""
    n = len(nums)
    copy = nums[:]
    nums = sorted(nums)
    preSum = [0] + list(accumulate(nums))
    nums = [0] + nums  # [0]表示哨兵

    # !最左二分求最后能和哪个数齐平
    left, right = 0, n
    while left <= right:
        mid = (left + right) // 2
        diff = preSum[n] - preSum[mid] - (n - mid) * nums[mid]
        if k < diff:
            left = mid + 1
        else:
            right = mid - 1

    # 如果最小值可以小到0 那么就直接返回[0]*n
    min_ = nums[left]
    if min_ == 0:
        return [0] * n

    overflow = k - (preSum[n] - preSum[left] - (n - left) * nums[left])
    div, mod = 0, 0
    count = n - left + 1

    if count:
        div, mod = divmod(overflow, count)  # mod个数需要再减1
    min_ -= div

    for i in range(n):
        if copy[i] > min_ - int(mod > 0):
            copy[i] = min_ - int(mod > 0)
            mod -= 1

    return copy


print(maximizeMinValue2(list(reversed([1, 2, 3, 4, 5, 6, 7, 8, 9, 10])), k=5))
