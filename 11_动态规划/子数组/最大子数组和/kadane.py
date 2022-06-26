
from typing import List


def kanade(nums: List[int], getMax=True) -> int:
    """求最大/最小子数组和"""
    n = len(nums)
    if n == 0:
        raise ValueError("nums is empty")
    if n == 1:
        return nums[0]

    res = -int(1e20) if getMax else int(1e20)
    cur = 0
    for num in nums:
        if getMax:
            cur = max(cur, 0) + num
            res = max(res, cur)
        else:
            cur = min(cur, 0) + num
            res = min(res, cur)
    return res

