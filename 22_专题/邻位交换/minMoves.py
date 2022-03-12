from itertools import accumulate
from typing import List


def minMoves(indexes: List[int], target: int) -> int:
    """得到连续 target 个 相同字符 的最少相邻交换次数"""
    indexes = [num - i for i, num in enumerate(indexes)]
    preSum = [0] + list(accumulate(indexes))

    res = int(1e20)
    # left+k-1<len(ones)
    for left in range(len(indexes) + 1 - target):
        right = left + target - 1
        mid = (left + right) >> 1
        leftSum = indexes[mid] * (mid - left) - (preSum[mid] - preSum[left])
        rightSum = preSum[right + 1] - preSum[mid + 1] - indexes[mid] * (right - mid)
        res = min(res, leftSum + rightSum)

    return res
