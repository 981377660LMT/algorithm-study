# 最大子数组和/最小子数组和

from typing import List

INF = int(1e20)


def kanade(nums: List[int], getMax=True) -> int:
    """求最大/最小子数组和"""
    n = len(nums)
    if n == 0:
        raise Exception("nums is empty")

    res = -INF if getMax else INF
    dp = 0
    for num in nums:
        if getMax:
            dp = max(dp, 0) + num
            res = max(res, dp)
        else:
            dp = min(dp, 0) + num
            res = min(res, dp)
    return res
