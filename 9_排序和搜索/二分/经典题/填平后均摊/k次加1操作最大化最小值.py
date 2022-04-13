# k次加1操作最大化一个数组的最小值  二分的check函数O(1)


from typing import List
from itertools import accumulate


def maximizeMinValue(nums: List[int], delta: int) -> int:
    """delta次加1操作，让最小值最大化，返回最小值"""
    n = len(nums)
    nums = sorted(nums)
    preSum = [0] + list(accumulate(nums))
    nums = [0] + nums

    # 最右二分求最后能和哪个数齐平
    left, right = 0, n
    while left <= right:
        mid = (left + right) >> 1
        diff = mid * nums[mid] - preSum[mid]
        if diff <= delta:
            left = mid + 1
        else:
            right = mid - 1

    min_ = nums[right]
    overflow = delta - (right * nums[right] - preSum[right])
    min_ += overflow // right if right else 0
    return min_


def maximizeMinValue2(nums: List[int], delta: int) -> List[int]:
    """delta次加1操作，让最小值最大化，返回操作后的数组"""
    n = len(nums)
    copy = nums[:]
    nums = sorted(nums)
    preSum = [0] + list(accumulate(nums))
    nums = [0] + nums

    # 最右二分求最后能和哪个数齐平
    left, right = 0, n
    while left <= right:
        mid = (left + right) >> 1
        diff = mid * nums[mid] - preSum[mid]
        if diff <= delta:
            left = mid + 1
        else:
            right = mid - 1

    min_ = nums[right]
    overflow = delta - (right * nums[right] - preSum[right])
    div, mod = 0, 0
    if right:
        div, mod = divmod(overflow, right)
    min_ += div

    for i in range(n):
        if copy[i] < min_ + int(mod > 0):
            copy[i] = min_ + int(mod > 0)
            mod -= 1

    return copy


print(maximizeMinValue2(list(reversed([1, 2, 3, 4, 5, 6, 7, 8, 9, 10])), delta=5))
