"""
Y和.组成的字符串 最多邻位交换 k 次,求最大连续 Y 的个数

fix模型/fixK模型
1703. 得到连续 K 个 1 的最少相邻交换次数-母题
"""

from itertools import accumulate
from typing import Sequence, TypeVar

T = TypeVar("T")
INF = int(1e20)


def fix(seq: Sequence[T], target: T, k: int) -> int:
    """
    最多邻位交换 k 次,求 seq 中最大连续 target 的个数
    k<=1e12 len(seq)<=2e5
    """

    left, right = 1, int(1e18)
    while left <= right:
        mid = (left + right) // 2
        if minMoves(seq, k, target) <= k:
            left = mid + 1
        else:
            right = mid - 1
    return right


def minMoves(nums: Sequence[T], k: int, target: T) -> int:
    """得到连续 k 个 target 的最少相邻交换次数"""
    if k == 0:
        return 0

    dist = []
    for i in range(len(nums)):
        if nums[i] == target:
            dist.append(i - len(dist))  # 移动到对应target位置的距离
    preSum = [0] + list(accumulate(dist))

    res = INF
    # 枚举哪k个数移动到一起
    for left in range(len(dist) + 1 - k):
        right = left + k - 1
        mid = (left + right) // 2
        leftSum = dist[mid] * (mid - left) - (preSum[mid] - preSum[left])
        rightSum = preSum[right + 1] - preSum[mid] - dist[mid] * (right - mid + 1)
        res = min(res, leftSum + rightSum)

    return res


if __name__ == "__main__":
    s = input()
    k = int(input())
    print(fix(s, "Y", k))  # Y和.组成的字符串 最多邻位交换 k 次,求最大连续 Y 的个数
