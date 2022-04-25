from typing import List, Tuple
from collections import defaultdict

MOD = int(1e9 + 7)
INF = int(1e20)

# 6044. 花期内花的数目-树状数组区间更新


class BIT:
    """范围修改"""

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
    def fullBloomFlowers(self, flowers: List[List[int]], persons: List[int]) -> List[int]:
        res = [0] * len(persons)
        bit = BIT(int(1e9 + 10))
        for l, r in flowers:
            bit.add(l, r, 1)
        for i, p in enumerate(persons):
            res[i] = bit.query(p, p)
        return res

