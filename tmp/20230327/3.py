from bisect import bisect_right
from itertools import accumulate
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个正整数数组 nums 。

# 同时给你一个长度为 m 的整数数组 queries 。第 i 个查询中，你需要将 nums 中所有元素变成 queries[i] 。你可以执行以下操作 任意 次：

# 将数组里一个元素 增大 或者 减小 1 。
# 请你返回一个长度为 m 的数组 answer ，其中 answer[i]是将 nums 中所有元素变成 queries[i] 的 最少 操作次数。

# 注意，每次查询后，数组变回最开始的值。


def calDistSum(nums: List[int], k: int, preSum: List[int]) -> int:
    """有序数组所有点到x=k的距离之和

    排序+二分+前缀和 O(logn)
    """
    pos = bisect_right(nums, k)
    leftSum = k * pos - preSum[pos]
    rightSum = preSum[-1] - preSum[pos] - k * (len(nums) - pos)
    return leftSum + rightSum


class Solution:
    def minOperations(self, nums: List[int], queries: List[int]) -> List[int]:
        nums.sort()
        preSum = [0] + list(accumulate(nums))
        return [calDistSum(nums, x, preSum) for x in queries]
