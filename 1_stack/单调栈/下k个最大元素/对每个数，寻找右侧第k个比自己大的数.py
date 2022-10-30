"""对每个数,寻找右侧第k个比自己大的数"""

from typing import List
from bisect import bisect_right
from collections import defaultdict
from sortedcontainers import SortedList


def kthGreaterElement(nums: List[int], k: int) -> List[int]:
    """
    求每个数右侧下一个严格大于它的第k个数 (kth next greater)
    时间复杂度 O(n*k)

    !k次单调栈
    !第一个单调栈pop出去的元素放到第二个单调栈里面
    !第二个单调栈pop出去的元素放到第三个单调栈里面
    !...
    !第k个单调栈再被pop时统计
    """
    n = len(nums)
    stacks = {level: [] for level in range(k)}
    res = [-1] * n
    for i in range(n):
        while stacks[k - 1] and nums[stacks[k - 1][-1]] < nums[i]:
            res[stacks[k - 1].pop()] = nums[i]
        for level in range(k - 2, -1, -1):
            tmp = []
            while stacks[level] and nums[stacks[level][-1]] < nums[i]:
                tmp.append(stacks[level].pop())
            stacks[level + 1].extend(tmp[::-1])
            stacks[level].append(i)
    return res


def kthGreaterElement2(nums: List[int], k: int) -> List[int]:
    """
    求每个数右侧下一个严格大于它的第k个数 (kth next greater)
    时间复杂度 O(n*logn)

    !将相同的数字分为一组.
    !按照数字从大到小的顺序遍历分组,保证添加到有序集合中的元素都是`比当前数字更大的数字`的`下标`.
    !对每个分组的每个下标,从有序集合中找出`右侧第二个比它更大的下标`所对应的数字即可.
    """

    group = defaultdict(list)
    for i, num in enumerate(nums):
        group[num].append(i)

    sl = SortedList()  # 存index
    res = [-1] * len(nums)
    for num in sorted(group, reverse=True):  # 按照key从大到小遍历
        for index in group[num]:
            pos = bisect_right(sl, index) - 1
            if pos + k < len(sl):
                res[index] = nums[sl[pos + k]]

        sl.update(group[num])

    return res


for func in [kthGreaterElement, kthGreaterElement2]:
    assert func(nums=[11, 13, 15, 12, 0, 15, 12, 11, 9], k=2) == [
        15,
        15,
        -1,
        -1,
        12,
        -1,
        -1,
        -1,
        -1,
    ]


########################################################################################
# Given an array arr, and an integer k,
# find the kth next greater element for any element arr[i],
# or -1 if doesn’t exist.

# input:
# arr = [1,4,2,5,3]
# k = 2

# output:
# [3, -1, 5, -1, -1]

# constraint
# len(arr) < 10^5, k < 50, arr[i] < 10^9


def findNextKthLarge(nums: List[int], k: int) -> List[int]:
    """
    对每个数,寻找右侧`值域`中比自己`严格`大的数中的第k个
    倒序遍历+SortedList二分查找
    """
    n = len(nums)
    res = [-1] * n
    sl = SortedList()
    for i in range(n - 1, -1, -1):
        cur = nums[i]
        pos = sl.bisect_right(cur)
        if pos + k - 1 < len(sl):
            res[i] = sl[pos + k - 1]  # type: ignore
        sl.add(cur)
    return res


assert findNextKthLarge([1, 4, 2, 2, 2], 2) == [2, -1, -1, -1, -1]
assert findNextKthLarge([1, 4, 2, 5, 3], 2) == [3, -1, 5, -1, -1]
