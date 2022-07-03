from typing import List, Tuple
from collections import defaultdict

MOD = int(1e9 + 7)
INF = int(1e20)

# 6044. 花期内花的数目-树状数组区间更新


class BIT:
    """单点修改"""

    def __init__(self, n: int):
        self.size = n
        self.tree = defaultdict(int)

    @staticmethod
    def _lowbit(index: int) -> int:
        return index & -index

    def add(self, index: int, delta: int) -> None:
        index += 1
        while index <= self.size:
            self.tree[index] += delta
            index += self._lowbit(index)

    def query(self, index: int) -> int:
        index += 1
        res = 0
        while index > 0:
            res += self.tree[index]
            index -= self._lowbit(index)
        return res

    def sumRange(self, left: int, right: int) -> int:
        return self.query(right) - self.query(left - 1)


class Solution:
    def fullBloomFlowers(
        self, flowers: List[List[int]], persons: List[int]
    ) -> List[int]:
        res = [0] * len(persons)
        bit = BIT(int(1e9 + 10))
        for l, r in flowers:  # 差分修改
            bit.add(l - 1, 1)
            bit.add(r, -1)
        for i, p in enumerate(persons):
            res[i] = bit.query(p - 1)
        return res
