# 最大子数组和/最小子数组和

from typing import List, Tuple

INF = int(1e20)


def maxSubarraySum(nums: List[int]) -> Tuple[int, int, int]:
    n = len(nums)
    if n == 0:
        return 0, 0, 0  # !根据题意返回

    maxSum = nums[0]
    curSum = nums[0]
    curStart = 0
    start = 0
    end = 0
    for i in range(1, n):
        if curSum < 0:
            curSum = 0
            curStart = i
        curSum += nums[i]
        if curSum > maxSum:
            maxSum = curSum
            start = curStart
            end = i + 1

    if maxSum < 0:
        return 0, 0, 0  # !根据题意返回

    return maxSum, start, end


def maxSubarraySumTwoSum(nums: List[int], gap: int) -> int:
    """最大两段子段和（两段必须间隔至少 gap 个数）."""

    def max(a: int, b: int) -> int:
        return a if a > b else b

    n = len(nums)
    sufSumMax = [0] * n
    sufSumMax[n - 1] = nums[n - 1]
    curSumMax = nums[n - 1]
    for i in range(n - 2, -1, -1):
        v = nums[i]
        curSumMax = max(curSumMax + v, v)
        sufSumMax[i] = max(sufSumMax[i + 1], curSumMax)
    curSumMax = nums[0]
    preSumMax = nums[0]
    res = preSumMax + sufSumMax[1 + gap]
    for i in range(1, n - 1 - gap):
        v = nums[i]
        curSumMax = max(curSumMax + v, v)
        preSumMax = max(preSumMax, curSumMax)
        res = max(res, preSumMax + sufSumMax[i + 1 + gap])
    return res


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
