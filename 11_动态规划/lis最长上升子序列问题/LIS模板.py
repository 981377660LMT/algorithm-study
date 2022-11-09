"""
贪心 + 二分查找
LIS[i]表示长度为 i+1 的子序列尾部元素的值
每次遍历到一个新元素,用二分查找法找到第一个大于等于它的元素,然后更新LIS
"""
# LIS模板

from typing import List
from bisect import bisect_left, bisect_right


def LIS(nums: List[int], isStrict=True) -> int:
    """求LIS长度"""
    n = len(nums)
    if n <= 1:
        return n

    lis = []  # lis[i] 表示长度为 i 的上升子序列的最小末尾值
    for i in range(n):
        pos = bisect_left(lis, nums[i]) if isStrict else bisect_right(lis, nums[i])
        if pos == len(lis):
            lis.append(nums[i])
        else:
            lis[pos] = nums[i]

    return len(lis)


def caldp(nums: List[int], isStrict=True) -> List[int]:
    """求以每个位置为结尾的LIS长度(包括自身)"""
    if not nums:
        return []

    n = len(nums)
    res = [0] * n
    lis = []
    for i in range(n):
        pos = bisect_left(lis, nums[i]) if isStrict else bisect_right(lis, nums[i])
        if pos == len(lis):
            lis.append(nums[i])
            res[i] = len(lis)
        else:
            lis[pos] = nums[i]
            res[i] = pos + 1
    return res


if __name__ == "__main__":
    assert LIS([10, 9, 2, 5, 3, 7, 101, 18]) == 4
    print(caldp([10, 9, 2, 5, 3, 7, 101, 18]))
