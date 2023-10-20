from typing import Callable, List
from itertools import accumulate
from bisect import bisect_right


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


def distSumOfAllPairs(sortedNums: List[int]) -> int:
    """`有序数组`中所有点对两两距离之和.一共有n*(n-1)//2对点对."""
    res, preSum = 0, 0
    for i, v in enumerate(sortedNums):
        res += v * i - preSum
        preSum += v
    return res
