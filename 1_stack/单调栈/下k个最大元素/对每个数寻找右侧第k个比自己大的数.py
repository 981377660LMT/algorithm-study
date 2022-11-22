"""对每个数,寻找右侧第k个比自己大的数"""

from typing import List
from bisect import bisect_right
from collections import defaultdict
from sortedcontainers import SortedList


def kthGreaterElement(nums: List[int], k: int) -> List[int]:
    """
    - 求每个数右侧下一个`严格大于`它的第k个数的`索引` (kth next greater)
    - `不存在为n`
    - 时间复杂度 O(n*k)

    !k次单调栈
    !第一个单调栈pop出去的元素放到第二个单调栈里面
    !第二个单调栈pop出去的元素放到第三个单调栈里面
    !...
    !第k个单调栈再被pop时统计
    """

    n = len(nums)
    res = [n] * n
    stacks = [[] for _ in range(k)]
    tmp = []
    for i in range(n):
        # 从最后一个单调栈开始处理
        for j in range(k - 1, -1, -1):
            while stacks[j] and nums[stacks[j][-1]] < nums[i]:  # 严格大于
                top = stacks[j].pop()
                if j == k - 1:
                    res[top] = i
                else:
                    tmp.append(top)
            if j + 1 < k:
                # 倒序进入下一个单调栈，保证所有单调栈的单调性
                while tmp:
                    stacks[j + 1].append(tmp.pop())
        stacks[0].append(i)
    return res


def kthGreaterElement2(nums: List[int], k: int) -> List[int]:
    """
    - 求每个数右侧下一个`严格大于`它的第k个数的`索引` (kth next greater)
    - `不存在为n`
    - 时间复杂度 O(n*logn)

    !将相同的数字分为一组.
    !按照数字从大到小的顺序遍历分组,保证添加到有序集合中的元素都是`比当前数字更大的数字`的`下标`.
    !对每个分组的每个下标,从有序集合中找出`右侧第二个比它更大的下标`所对应的数字即可.
    """

    group = defaultdict(list)
    for i, num in enumerate(nums):
        group[num].append(i)

    sl = SortedList()  # 存index
    res = [len(nums)] * len(nums)
    for num in sorted(group, reverse=True):  # 按照key从大到小遍历
        for index in group[num]:
            pos = bisect_right(sl, index) - 1
            if pos + k < len(sl):
                res[index] = sl[pos + k]

        sl.update(group[num])

    return res


for func in [kthGreaterElement, kthGreaterElement2]:
    assert func(nums=[11, 13, 15, 12, 0, 15, 12, 11, 9], k=2) == [2, 5, 9, 9, 6, 9, 9, 9, 9]
    assert func(nums=[1, 2, 3, 4, 5, 6, 7, 8, 9], k=3) == [3, 4, 5, 6, 7, 8, 9, 9, 9]
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
