from typing import List


def kanade(nums: List[int], getMax=True) -> int:
    """求最大/最小子数组和"""
    n = len(nums)
    if n == 0:
        raise Exception("nums is empty")

    res = -int(1e20) if getMax else int(1e20)
    dp = 0
    for num in nums:
        if getMax:
            dp = max(dp, 0) + num
            res = max(res, dp)
        else:
            dp = min(dp, 0) + num
            res = min(res, dp)
    return res
