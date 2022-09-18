# !1.一般定位元素直接调最右二分，然后看i和i-1
# !pos = bisect_right(nums, target) - 1
# !pos = bisect_left(nums, target) - 1
# !左边的元素下标为pos,右边的元素下标为pos+1

# !2. 如果不确定,可以看三个位置 i-1 i i+1
# for cand in (pos - 1, pos, pos + 1):
#     if 0 <= cand < len(nums):
#         dist = abs(nums[cand] - target)
#         if dist < res[0]:
#             res = (dist, cand)

from bisect import bisect_left
from typing import List

INF = int(1e18)


def findNearest(nums: List[int], target: int) -> int:
    """
    二分查找有序数组nums中与距离target最近的元素的下标
    如果有多个元素与target距离相同,则返回最`左边`的元素的下标
    注意数组中可能有重复元素

    返回最`左边`的元素的下标,需要最左二分定位
    """
    pos = bisect_left(nums, target) - 1
    res, dist = -1, INF
    if pos >= 0:
        cand = abs(nums[pos] - target)
        if cand < dist:
            dist = cand
            res = pos

    if pos + 1 < len(nums):
        cand = abs(nums[pos + 1] - target)
        if cand < dist:
            dist = cand
            res = pos + 1

    return res


print(findNearest([1, 2, 3, 4, 5], 3))
