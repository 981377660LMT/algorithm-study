"""
n,q<=2e5
Ai<=1e9
Xi<=1e9

每次操作可以是+1或-1
对每个询问,求将数组所有数变为Xi的最小操作次数

有序数组所有点到x=k的距离之和
"""

from bisect import bisect_right
from itertools import accumulate

from typing import List


def calDistSum(nums: List[int], k: int, preSum: List[int]) -> int:
    """`有序数组`所有点到x=k的距离之和

    排序+二分+前缀和 O(logn)
    """
    pos = bisect_right(nums, k)
    leftSum = k * pos - preSum[pos]
    rightSum = preSum[-1] - preSum[pos] - k * (len(nums) - pos)
    return leftSum + rightSum


def calDistSumOfAllPairs(nums: List[int]) -> int:
    """`有序数组`中所有点对两两距离之和.一共有n*(n-1)//2对点对."""
    res, preSum = 0, 0
    for i, v in enumerate(nums):
        res += v * i - preSum
        preSum += v
    return res


def main() -> None:
    import sys
    import os

    sys.setrecursionlimit(int(1e6))
    input = sys.stdin.readline
    MOD = int(1e9 + 7)

    n, q = map(int, input().split())
    nums = list(map(int, input().split()))
    nums.sort()
    preSum = [0] + list(accumulate(nums))
    for _ in range(q):
        x = int(input())
        print(calDistSum(nums, x, preSum))
