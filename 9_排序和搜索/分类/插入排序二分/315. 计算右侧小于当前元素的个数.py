from typing import List
from sortedcontainers import SortedList


def countSmaller(nums: List[int]) -> List[int]:
    """求各个位置处的逆序对数量"""
    n = len(nums)
    res = [0] * n
    visited = SortedList()
    for i in range(n - 1, -1, -1):
        smaller = visited.bisect_left(nums[i])
        res[i] = smaller
        visited.add(nums[i])

    return res


def countInv(nums: List[int]) -> int:
    """求数组逆序对数量之和"""
    n = len(nums)
    res = 0
    visited = SortedList()
    for i in range(n - 1, -1, -1):
        smaller = visited.bisect_left(nums[i])
        res += smaller
        visited.add(nums[i])

    return res
