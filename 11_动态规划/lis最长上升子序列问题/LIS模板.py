"""
贪心 + 二分查找
LIS[i]表示长度为 i+1 的子序列尾部元素的值
每次遍历到一个新元素,用二分查找法找到第一个大于等于它的元素,然后更新LIS
"""

from typing import List
from bisect import bisect_left, bisect_right


# LIS模板
def LIS(nums: List[int], isStrict=True) -> int:
    """求LIS长度"""
    n = len(nums)
    if n <= 1:
        return n

    res = [nums[0]]
    for i in range(1, n):
        pos = bisect_left(res, nums[i]) if isStrict else bisect_right(res, nums[i])
        if pos >= len(res):
            res.append(nums[i])
        else:
            res[pos] = nums[i]

    return len(res)


def caldp(nums: List[int], isStrict=True) -> List[int]:
    """求以每个位置为结尾的LIS长度(包括自身)"""
    if not nums:
        return []
    res = [1] * len(nums)
    LIS = [nums[0]]
    for i in range(1, len(nums)):
        if nums[i] > LIS[-1]:
            LIS.append(nums[i])
            res[i] = len(LIS)
        else:
            pos = bisect_left(LIS, nums[i]) if isStrict else bisect_right(LIS, nums[i])
            LIS[pos] = nums[i]
            res[i] = pos + 1
    return res


if __name__ == "__main__":
    assert LIS([10, 9, 2, 5, 3, 7, 101, 18]) == 4
    print(caldp([10, 9, 2, 5, 3, 7, 101, 18]))
