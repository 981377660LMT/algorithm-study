from typing import MutableSequence
from collections import defaultdict, deque
from sortedcontainers import SortedList


def minAdjacentSwap1(nums1: MutableSequence[int], nums2: MutableSequence[int]) -> int:
    """求使两个数组相等的最少邻位交换次数 映射+求逆序对 时间复杂度`O(nlogn)`"""

    def countInversionPair(nums: MutableSequence[int]) -> int:
        """计算逆序对的个数 时间复杂度`O(nlogn)`"""
        res = 0
        sl = SortedList()
        for num in reversed(nums):
            pos = sl.bisect_left(num)
            res += pos
            sl.add(num)
        return res

    # 含有重复元素的映射 例如nums [1,3,2,1,4] 表示已经排序的数组  [0,1,2,3,4]
    # 那么nums1 [1,1,3,4,2] 就 映射到 [0,3,1,4,2]
    mapping = defaultdict(deque)
    for index, num in enumerate(nums2):
        mapping[num].append(index)

    for index, num in enumerate(nums1):
        mapped = mapping[num].popleft()
        nums1[index] = mapped

    res = countInversionPair(nums1)

    return res


def minAdjacentSwap2(nums1: MutableSequence[int], nums2: MutableSequence[int]) -> int:
    """求使两个数组相等的最少邻位交换次数

    对每个数，贪心找到对应的最近位置交换

    时间复杂度`O(n^2)`
    """
    res = 0

    for num in nums1:
        index = nums2.index(num)  # 最左边的第一个位置
        res += index
        nums2.pop(index)  # 已经被换到最左边了，所以减1

    return res
