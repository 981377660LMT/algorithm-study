from bisect import bisect_left, bisect_right
from typing import List, Tuple


def getLis(nums: List[int], strict=True) -> Tuple[List[int], List[int]]:
    """求LIS 返回(LIS,LIS的组成下标)"""
    n = len(nums)

    lis = []  # lis[i] 表示长度为 i 的上升子序列的最小末尾值
    dpIndex = [0] * n  # 每个元素对应的LIS长度
    f = bisect_left if strict else bisect_right
    for i in range(n):
        pos = f(lis, nums[i])
        if pos == len(lis):
            lis.append(nums[i])
        else:
            lis[pos] = nums[i]
        dpIndex[i] = pos

    res, resIndex = [], []
    j = len(lis) - 1
    for i in range(n - 1, -1, -1):
        if dpIndex[i] == j:
            res.append(nums[i])
            resIndex.append(i)
            j -= 1
    return res[::-1], resIndex[::-1]
