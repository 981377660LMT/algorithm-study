from typing import List
from heapq import heapify, heapreplace


def max2(a: int, b: int) -> int:
    return a if a > b else b


def smallestRange(nums: List[List[int]]) -> List[int]:
    """
    https://leetcode.cn/problems/smallest-range-covering-elements-from-k-lists

    有 k 个非递减排列的整数列表。
    找到一个最小区间，使得 k 个列表中的每个列表至少有一个数包含在其中。
    如果有多个长度相同的区间，返回起始点最小的那个。

    时间复杂度 O(n log k)，其中 n 是所有列表中元素的总数，k 是列表的数量。
    空间复杂度 O(k)，其中 k 是列表的数量。
    """
    if not nums:
        return []
    h = [(arr[0], i, 0) for i, arr in enumerate(nums)]
    heapify(h)
    res1 = h[0][0]
    res2 = r = max(arr[0] for arr in nums)
    while h[0][2] + 1 < len(nums[h[0][1]]):
        _, i, j = h[0]
        x = nums[i][j + 1]
        heapreplace(h, (x, i, j + 1))
        r = max2(r, x)
        l = h[0][0]
        if r - l < res2 - res1:
            res1, res2 = l, r
    return [res1, res2]
