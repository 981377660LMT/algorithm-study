from collections import defaultdict
from typing import List


class BIT:
    def __init__(self, n: int):
        self.size = n
        self._tree1 = defaultdict(int)
        self._tree2 = defaultdict(int)

    @staticmethod
    def _lowbit(index: int) -> int:
        return index & -index

    def add(self, left: int, right: int, delta: int) -> None:
        """闭区间[left, right]加delta"""
        self._add(left, delta)
        self._add(right + 1, -delta)

    def query(self, left: int, right: int) -> int:
        """闭区间[left, right]的和"""
        return self._query(right) - self._query(left - 1)

    def _add(self, index: int, delta: int) -> None:
        if index <= 0:
            raise ValueError('index 必须是正整数')

        rawIndex = index
        while index <= self.size:
            self._tree1[index] += delta
            self._tree2[index] += (rawIndex - 1) * delta
            index += self._lowbit(index)

    def _query(self, index: int) -> int:
        if index > self.size:
            index = self.size

        rawIndex = index
        res = 0
        while index > 0:
            res += rawIndex * self._tree1[index] - self._tree2[index]
            index -= self._lowbit(index)
        return res


class Solution:
    def amountPainted(self, paint: List[List[int]]) -> List[int]:
        res = []
        bit = BIT(int(1e5))
        for start, end in paint:
            start, end = start + 1, end + 1
            # 注意是end-1 因为一段区间对应一个数
            diff = bit.query(start, end - 1)
            print(diff, end - start)
            res.append(max(0, end - start - diff))
            bit.add(start, end - 1, 1)
        return res


# print(Solution().amountPainted(paint=[[1, 4], [4, 7], [5, 8]]))
# print(Solution().amountPainted(paint=[[1, 4], [5, 8], [4, 7]]))
# print(Solution().amountPainted(paint=[[1, 5], [2, 4]]))
print(Solution().amountPainted(paint=[[2, 4], [0, 4], [1, 4], [1, 5], [0, 2]]))
# [2, 2, 0, 1, 0]
print(Solution().amountPainted(paint=[[0, 5], [0, 2], [0, 3], [0, 4], [0, 5]]))
# [5, 0, 0, 0, 0]
