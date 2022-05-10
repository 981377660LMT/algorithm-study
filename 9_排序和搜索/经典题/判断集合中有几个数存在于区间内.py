from bisect import bisect_left, bisect_right
from typing import List, Tuple

# 二分


def cal(nums: List[int], interval: Tuple[int, int]) -> int:
    """判断集合中是否有几个数存在于区间内"""
    nums = sorted(nums)
    left, right = interval
    pos1 = bisect_left(nums, left)
    pos2 = bisect_right(nums, right) - 1
    return pos2 - pos1 + 1


def cal2(nums: List[int], interval: Tuple[int, int]) -> bool:
    """判断集合中是否有数存在于区间内"""
    nums = sorted(nums)
    left, right = interval
    pos = bisect_right(nums, right) - 1
    if pos >= 0 and nums[pos] >= left:
        return True
    return False


print(cal([1, 3, 4, 8], (0, 3)))
print(cal([1, 3, 4, 8], (100, 101)))
print(cal([1, 1, 1, 2], (1, 1)))
print(cal2([1, 3, 6, 9], (3, 4)))


# 前缀和 (略)

