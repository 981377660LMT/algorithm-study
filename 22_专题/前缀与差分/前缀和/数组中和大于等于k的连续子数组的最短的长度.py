# 数是有正有负的
# 有一组数和一个整数k， 返回数组中和大于等于k的连续子数组的个数
#  prefixSum[R] - prefixSum[L] >= K, which means prefixSum[L] <= prefixSum[R] - K.
# treeMap 每次找当前前缀和-k在哪个位置 再插入 不断更新最小值
# 维护前缀和，用树状数组/平衡树处理

from typing import List
from sortedcontainers import SortedDict

INF = 0x3F3F3F3F


def minSubArray(nums: List[int], k: int) -> int:
    res = INF
    curSum = 0
    treeMap = SortedDict({0: -1})

    for i, num in enumerate(nums):
        curSum += num
        lower = curSum - k
        cand = treeMap.bisect_left(lower)
        res += cand
        treeMap[curSum] = i

    return res


print(minSubArray([1, 2, -1, 4, 2], 3))

# 有一组数和一个整数k， 返回数组中和大于等于k的连续子数组的最短的长度
# 单调队列# 862. 和至少为 K 的最短子数组

