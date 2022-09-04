"""
find kth mex
"""

from typing import List


def findKthMex1(nums: List[int], k: int) -> int:
    """二分搜索有序数组缺失的第k个`正整数`

    Args:
        nums: List[int] 正整数数组
        k: int 第k个正整数 k>=1
    """
    left, right = 0, len(nums) - 1
    while left <= right:
        mid = (left + right) // 2
        diff = nums[mid] - (mid + 1)
        if diff >= k:
            right = mid - 1
        else:
            left = mid + 1
    return left + k


def findKthMex2(nums: List[int], k: int) -> int:
    """二分搜索有序数组缺失的第k个`非负整数`

    Args:
        nums: List[int] 非负整数数组
        k: int 第k个非负整数 k>=1
    """
    left, right = 0, len(nums) - 1
    while left <= right:
        mid = (left + right) // 2
        diff = nums[mid] - mid
        if diff >= k:
            right = mid - 1
        else:
            left = mid + 1
    return left + k - 1


if __name__ == "__main__":
    nums = [2, 3, 4, 5, 6, 7, 8, 9, 10, 10]
    print(findKthMex1(nums, 5))
    print(findKthMex2(nums, 2))
