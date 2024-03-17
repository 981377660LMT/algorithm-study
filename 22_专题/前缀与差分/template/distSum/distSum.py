from typing import Callable, List, Optional
from itertools import accumulate
from bisect import bisect_left, bisect_right


def distSum(sortedNums: List[int]) -> Callable[[int], int]:
    """`有序数组`所有点到x=k的距离之和

    排序+二分+前缀和 O(logn)
    """
    preSum = [0] + list(accumulate(sortedNums))

    def query(k: int) -> int:
        pos = bisect_right(sortedNums, k)
        leftSum = k * pos - preSum[pos]
        rightSum = preSum[-1] - preSum[pos] - k * (len(sortedNums) - pos)
        return leftSum + rightSum

    return query


def distSumRange(sortedNums: List[int]) -> Callable[[int, int, int], int]:
    """`有序数组`切片[start:end)中所有点到`x=k`的距离之和.

    排序+二分+前缀和 O(logn)
    """
    preSum = [0] + list(accumulate(sortedNums))

    def query(k: int, start: int, end: int) -> int:
        if start < 0:
            start = 0
        if end > len(sortedNums):
            end = len(sortedNums)
        if start >= end:
            return 0
        pos = bisect_left(sortedNums, k)
        if pos <= start:
            return (preSum[end] - preSum[start]) - k * (end - start)
        if pos >= end:
            return k * (end - start) - (preSum[end] - preSum[start])
        leftSum = k * (pos - start) - (preSum[pos] - preSum[start])
        rightSum = preSum[end] - preSum[pos] - k * (end - pos)
        return leftSum + rightSum

    return query


def distSumOfAllPairs(sortedNums: List[int]) -> int:
    """`有序数组`中所有点对两两距离之和.一共有n*(n-1)//2对点对."""
    res, preSum = 0, 0
    for i, v in enumerate(sortedNums):
        res += v * i - preSum
        preSum += v
    return res


def distSumOfAllPairsRange(sortedNums: List[int], start: int, end: int) -> int:
    """`有序数组`切片[start:end)中所有点对两两距离之和."""
    res, preSum = 0, 0
    for i in range(start, end):
        v = sortedNums[i]
        res += v * i - preSum
        preSum += v
    return res


def getMedian(sortedNums: List[int], start: Optional[int] = None, end: Optional[int] = None) -> int:
    """有序数组的中位数(向下取整)."""
    if start is None:
        start = 0
    if end is None:
        end = len(sortedNums)
    if start < 0:
        start = 0
    if end > len(sortedNums):
        end = len(sortedNums)
    if start >= end:
        return 0
    if (end - start) & 1 == 0:
        return (sortedNums[(end + start) // 2 - 1] + sortedNums[(end + start) // 2]) // 2
    return sortedNums[(end + start) // 2]


if __name__ == "__main__":
    # https://leetcode.cn/problems/apply-operations-to-maximize-frequency-score/description/
    class Solution:
        def maxFrequencyScore(self, nums: List[int], k: int) -> int:
            nums.sort()
            D = distSumRange(nums)
            res, left = 0, 0
            for right in range(len(nums)):
                while left <= right:
                    median = getMedian(nums, left, right + 1)
                    if D(median, left, right + 1) <= k:
                        break
                    left += 1
                res = max(res, right - left + 1)
            return res

    nums = [1, 2, 3, 4, 5]
    D = distSumRange(nums)
    print(D(1, 0, 5))
