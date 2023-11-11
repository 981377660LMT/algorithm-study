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
import sys

from typing import List

sys.setrecursionlimit(int(1e6))
input = sys.stdin.readline
MOD = int(1e9 + 7)


def distanceSum(sortedArray: List[int], k: int, preSum: List[int]) -> int:
    """`有序数组`所有点到x=k的绝对值距离之和

    排序+二分+前缀和 O(logn)
    """
    pos = bisect_right(sortedArray, k)
    leftDist = k * pos - preSum[pos]
    rightDist = preSum[-1] - preSum[pos] - k * (len(sortedArray) - pos)
    return leftDist + rightDist


def main() -> None:
    n, q = map(int, input().split())
    nums = list(map(int, input().split()))
    nums.sort()
    preSum = [0] + list(accumulate(nums))
    for _ in range(q):
        x = int(input())
        print(distanceSum(nums, x, preSum))
