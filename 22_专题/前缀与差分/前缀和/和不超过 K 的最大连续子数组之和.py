from itertools import accumulate
from typing import List
from sortedcontainers import SortedList

# 在一维数组中，求解和不超过 K 的最大连续子数组之和
# 数组元素可以是负数

# 如果和为非负数，滑动窗口
# 如果和为负数，考虑单调队列/有序集合等数据结构
# 和不大于 K 的最大连续子数组之和


def maxSubArraySumNoMoreThanK(nums: List[int], k: int) -> int:
    res = -int(1e20)
    n = len(nums)
    preSum = [0] + list(accumulate(nums))

    sl = SortedList([0])
    for i in range(1, n + 1):
        pre = preSum[i]
        left = sl.bisect_left(pre - k)
        if left < len(sl):
            res = max(res, pre - sl[left])
        sl.add(pre)

    return res


print(maxSubArraySumNoMoreThanK([1, 2, 3, 4], 8))
print(maxSubArraySumNoMoreThanK([1, -2, 3, 4], 5))
