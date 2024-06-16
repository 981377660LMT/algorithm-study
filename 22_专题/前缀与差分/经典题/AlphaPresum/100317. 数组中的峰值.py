# 100317. 数组中的峰值(子数组峰值元素)
# https://leetcode.cn/problems/peaks-in-array/
# 数组 arr 中 大于 前面和后面相邻元素的元素被称为 峰值 元素。
# 给你一个整数数组 nums 和一个二维整数数组 queries 。
# 你需要处理以下两种类型的操作：
# queries[i] = [1, li, ri] ，求出子数组 nums[li..ri] 中 峰值 元素的数目。
# queries[i] = [2, indexi, vali] ，将 nums[indexi] 变为 vali 。
# 请你返回一个数组 answer ，它依次包含每一个第一种操作的答案。
# 注意：
# !子数组中 第一个 和 最后一个 元素都 不是 峰值元素。


from typing import List, Sequence, Union


class Solution:
    def countOfPeaks(self, nums: List[int], queries: List[List[int]]) -> List[int]:
        n = len(nums)
        bit = BITArray(n - 1)

        def add(i: int, v: int):
            if not (1 <= i <= n - 2):
                return
            if nums[i - 1] < nums[i] > nums[i + 1]:
                bit.add(i, v)

        for i in range(1, n - 1):
            add(i, 1)

        res = []
        for op, *args in queries:
            if op == 1:
                start, end = args
                end += 1
                if end - start <= 1:
                    res.append(0)
                else:
                    res.append(bit.queryRange(start + 1, end - 1))
            else:
                index, value = args
                if nums[index] == value:
                    continue
                for j in range(index - 1, index + 2):
                    add(j, -1)
                nums[index] = value
                for j in range(index - 1, index + 2):
                    add(j, 1)

        return res


class BITArray:
    """Point Add Range Sum, 0-indexed."""

    @staticmethod
    def _build(sequence: Sequence[int]) -> List[int]:
        tree = [0] * (len(sequence) + 1)
        for i in range(1, len(tree)):
            tree[i] += sequence[i - 1]
            parent = i + (i & -i)
            if parent < len(tree):
                tree[parent] += tree[i]
        return tree

    __slots__ = ("_n", "_tree")

    def __init__(self, lenOrSequence: Union[int, Sequence[int]]):
        if isinstance(lenOrSequence, int):
            self._n = lenOrSequence
            self._tree = [0] * (lenOrSequence + 1)
        else:
            self._n = len(lenOrSequence)
            self._tree = self._build(lenOrSequence)

    def add(self, index: int, delta: int) -> None:
        index += 1
        while index <= self._n:
            self._tree[index] += delta
            index += index & -index

    def query(self, right: int) -> int:
        """Query sum of [0, right)."""
        if right > self._n:
            right = self._n
        res = 0
        while right > 0:
            res += self._tree[right]
            right -= right & -right
        return res

    def queryRange(self, left: int, right: int) -> int:
        """Query sum of [left, right)."""
        return self.query(right) - self.query(left)

    def __len__(self) -> int:
        return self._n

    def __repr__(self) -> str:
        nums = []
        for i in range(1, self._n + 1):
            nums.append(self.queryRange(i, i + 1))
        return f"BITArray({nums})"
