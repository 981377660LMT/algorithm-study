# 3510. 移除最小数对使数组有序 II
# https://leetcode.cn/problems/minimum-pair-removal-to-sort-array-ii/description/
# 给你一个数组 nums，你可以执行以下操作任意次数：
#
# 选择 相邻 元素对中 和最小 的一对。如果存在多个这样的对，选择最左边的一个。
# 用它们的和替换这对元素。
# 返回将数组变为 非递减 所需的 最小操作次数 。
#
# 1.SortedList 维护 (pairSum, index)
# 2.Finder 维护每个下标 i 左侧最近剩余下标，以及右侧最近剩余下标.
# 3.dec 维护左边元素大于右边元素的pair数.

from typing import List
from itertools import pairwise

from sortedcontainers import SortedList


class Finder:
    __slots__ = ("_n", "_exist", "_prev", "_next")

    def __init__(self, n: int):
        self._n = n
        self._exist = [True] * n
        self._prev = [i - 1 for i in range(n)]
        self._next = [i + 1 for i in range(n)]

    def has(self, i: int) -> bool:
        return 0 <= i < self._n and self._exist[i]

    def erase(self, i: int) -> bool:
        if not self.has(i):
            return False
        l, r = self._prev[i], self._next[i]
        if l >= 0:
            self._next[l] = r
        if r < self._n:
            self._prev[r] = l
        self._exist[i] = False
        return True

    def prev(self, i: int) -> int:
        """
        返回`严格`小于i的最大元素,如果不存在,返回-1.
        !调用时需保证i存在.
        """
        return self._prev[i]

    def next(self, i: int) -> int:
        """
        返回`严格`大于i的最小元素.如果不存在,返回n.
        !调用时需保证i存在.
        """
        return self._next[i]


class Solution:
    def minimumPairRemoval(self, nums: List[int]) -> int:
        n = len(nums)
        sl = SortedList()  # (pairSum, index)
        remain = Finder(n)
        dec = 0

        for i, (x, y) in enumerate(pairwise(nums)):
            dec += x > y
            sl.add((x + y, i))

        res = 0
        while dec > 0:
            res += 1

            s, i = sl.pop(0)  # 删除相邻元素和最小的一对

            j = remain.next(i)
            dec -= nums[i] > nums[j]

            prev = remain.prev(i)
            if prev != -1:
                dec -= nums[prev] > nums[i]
                dec += nums[prev] > s
                sl.remove((nums[prev] + nums[i], prev))
                sl.add((nums[prev] + s, prev))

            next = remain.next(j)
            if next != n:
                dec -= nums[j] > nums[next]
                dec += s > nums[next]
                sl.remove((nums[j] + nums[next], j))
                sl.add((s + nums[next], i))

            nums[i] = s
            remain.erase(j)

        return res
