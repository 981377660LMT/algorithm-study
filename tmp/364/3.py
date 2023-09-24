from itertools import accumulate
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个长度为 n 下标从 0 开始的整数数组 maxHeights 。

# 你的任务是在坐标轴上建 n 座塔。第 i 座塔的下标为 i ，高度为 heights[i] 。

# 如果以下条件满足，我们称这些塔是 美丽 的：

# 1 <= heights[i] <= maxHeights[i]
# heights 是一个 山状 数组。
# 如果存在下标 i 满足以下条件，那么我们称数组 heights 是一个 山状 数组：


# 对于所有 0 < j <= i ，都有 heights[j - 1] <= heights[j]
# 对于所有 i <= k < n - 1 ，都有 heights[k + 1] <= heights[k]
# 请你返回满足 美丽塔 要求的方案中，高度和的最大值 。

from typing import List, Sequence, Union


class BIT2:
    """范围修改,0-indexed"""

    __slots__ = "size", "_tree1", "_tree2"

    def __init__(self, n: int):
        self.size = n + 5
        self._tree1 = dict()
        self._tree2 = dict()

    def add(self, left: int, right: int, delta: int) -> None:
        """区间[left, right)加delta."""
        right -= 1
        self._add(left, delta)
        self._add(right + 1, -delta)

    def query(self, left: int, right: int) -> int:
        """区间[left, right)的和."""
        right -= 1
        return self._query(right) - self._query(left - 1)

    def _add(self, index: int, delta: int) -> None:
        index += 1
        rawIndex = index
        while index <= self.size:
            self._tree1[index] = self._tree1.get(index, 0) + delta
            self._tree2[index] = self._tree2.get(index, 0) + (rawIndex - 1) * delta
            index += index & -index

    def _query(self, index: int) -> int:
        index += 1
        if index > self.size:
            index = self.size
        rawIndex = index
        res = 0
        while index > 0:
            res += rawIndex * self._tree1.get(index, 0) - self._tree2.get(index, 0)
            index &= index - 1
        return res

    def __repr__(self) -> str:
        arr = []
        for i in range(self.size):
            arr.append(self.query(i, i + 1))
        return str(arr)

    def __len__(self) -> int:
        return self.size


class Solution:
    def maximumSumOfHeights(self, maxHeights: List[int]) -> int:
        def makeDp(nums: List[int]) -> List[int]:
            """以nums[i]为peek的山脉数组，后缀和的最大值"""
            bit = BIT2(len(nums))
            threshold = INF
            for i in range(len(nums)):
                cur = min(nums[i], threshold)
                bit.add(0, i, cur)
                threshold = cur
            return [bit.query(i, i + 1) for i in range(len(nums))]

        dp1 = makeDp(maxHeights[:])
        dp2 = makeDp(maxHeights[::-1])[::-1]
        print(dp1, dp2)

        return


# maxHeights = [5,3,4,1,1]

# print(Solution().maximumSumOfHeights(maxHeights=[5, 3, 4, 1, 1]))
# maxHeights = [6,5,3,9,2,7]
# print(Solution().maximumSumOfHeights(maxHeights=[6, 5, 3, 9, 2, 7]))
# maxHeights = [3,2,5,5,2,3]

print(Solution().maximumSumOfHeights(maxHeights=[3, 2, 5, 5, 2, 3]))
