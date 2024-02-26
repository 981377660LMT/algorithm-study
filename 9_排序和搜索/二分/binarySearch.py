# 34. 在排序数组中查找元素的第一个和最后一个位置
# https://leetcode.cn/problems/find-first-and-last-position-of-element-in-sorted-array/description/
# 查找非递减数组中目标值的第一个和最后一个位置


from typing import List, Sequence


def binarySearch(nums: Sequence[int], target: int, findFirst=True) -> int:
    """查询非递减数组中目标值的第一个或最后一个位置.

    Args:

    - nums: 非递减数组
    - target: 目标值
    - findFirst: 是否查询第一个位置. 默认为 True.

    Returns:
    - int: 目标值的第一个或最后一个位置. 如果目标值不存在, 返回 -1.

    """

    if not nums or (nums[0] > target or nums[-1] < target):
        return -1
    if findFirst:
        left, right = 0, len(nums) - 1
        while left <= right:
            mid = left + (right - left) // 2
            if nums[mid] < target:
                left = mid + 1
            else:
                right = mid - 1
        return left if left < len(nums) and nums[left] == target else -1
    else:
        left, right = 0, len(nums) - 1
        while left <= right:
            mid = left + (right - left) // 2
            if nums[mid] <= target:
                left = mid + 1
            else:
                right = mid - 1
        return left - 1 if left > 0 and nums[left - 1] == target else -1


class Solution:
    def searchRange(self, nums: List[int], target: int) -> List[int]:
        return [binarySearch(nums, target, True), binarySearch(nums, target, False)]
