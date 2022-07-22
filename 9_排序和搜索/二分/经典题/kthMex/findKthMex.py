from typing import List


def findKthMex1(nums: List[int], k: int) -> int:
    """二分搜索缺失的第k个正整数

    Args:
        nums: List[int] 正整数数组
        k: int 第k个正整数 k>=1
    """
    nums = sorted(set(nums))
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
    """二分搜索缺失的第k个非负整数

    Args:
        nums: List[int] 非负整数数组
        k: int 第k个非负整数 k>=1
    """
    nums = sorted(set(nums))
    left, right = 0, len(nums) - 1
    while left <= right:
        mid = (left + right) >> 1
        diff = nums[mid] - mid
        if diff >= k:
            right = mid - 1
        else:
            left = mid + 1
    return left + k - 1


def calMex(nums: List[int]) -> int:
    """求非负整数数组的Mex(最小的非负整数)"""
    nums = sorted(set(nums))
    left, right = 0, len(nums) - 1
    while left <= right:
        mid = (left + right) // 2
        diff = nums[mid] - mid
        if diff >= 1:
            right = mid - 1
        else:
            left = mid + 1
    return left


if __name__ == "__main__":
    nums = [2, 3, 4, 5, 6, 7, 8, 9, 10, 10]
    print(findKthMex1(nums, 5))
    print(findKthMex2(nums, 2))
