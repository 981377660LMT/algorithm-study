# 快速选择第k小元素(k>=1)
from random import randint
from typing import List


def quickSelect(nums: List[int], left: int, right: int, k: int) -> int:
    """[left,right]闭区间找到第k小的数(k从0开始)的下标"""

    def partition(nums: List[int], left: int, right: int) -> int:
        randIndex = randint(left, right)
        nums[left], nums[randIndex] = nums[randIndex], nums[left]

        pivotIndex, pivot = left, nums[left]
        # 小于基准数的放左边，大于基准数的放右边
        for i in range(left + 1, right + 1):
            if nums[i] < pivot:
                pivotIndex += 1
                nums[i], nums[pivotIndex] = nums[pivotIndex], nums[i]
        nums[left], nums[pivotIndex] = nums[pivotIndex], nums[left]
        return pivotIndex

    pos = partition(nums, left, right)
    if pos == k:
        return nums[pos]
    if pos > k:
        return quickSelect(nums, left, pos - 1, k)
    return quickSelect(nums, pos + 1, right, k)


class Solution:
    def minMoves2(self, nums: List[int]) -> int:
        mid = quickSelect(nums, 0, len(nums) - 1, len(nums) // 2)
        return sum(abs(num - mid) for num in nums)

