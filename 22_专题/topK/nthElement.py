# 快速选择算法

from random import randint
from typing import List


def quickSelect(nums: List[int], lower: int, upper: int, nth: int) -> int:
    """`[lower,upper]`闭区间找到第`nth`小的数(`nth`从0开始)的数"""
    pos = partition(nums, lower, upper)
    if pos == nth:
        return nums[pos]
    if pos > nth:
        return quickSelect(nums, lower, pos - 1, nth)
    return quickSelect(nums, pos + 1, upper, nth)


def partition(nums: List[int], left: int, right: int) -> int:
    randIndex = randint(left, right)
    nums[left], nums[randIndex] = nums[randIndex], nums[left]

    pivotIndex, pivot = left, nums[left]
    # !小于基准数的放左边，大于基准数的放右边
    for i in range(left + 1, right + 1):
        if nums[i] < pivot:
            pivotIndex += 1
            nums[i], nums[pivotIndex] = nums[pivotIndex], nums[i]
    nums[left], nums[pivotIndex] = nums[pivotIndex], nums[left]
    return pivotIndex


def nthElement(nums: List[int], lower: int, upper: int, nth: int) -> int:
    """`[lower,upper]`闭区间找到第`nth`小的数(`nth`从0开始)"""
    nums = nums[:]
    return quickSelect(nums, lower, upper, nth)


if __name__ == "__main__":
    nums = [4, 3, 2, 1, 5, 6, 7, 8]
    print(nthElement(nums, 0, len(nums) - 1, 0))
