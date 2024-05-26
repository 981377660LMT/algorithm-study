from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)


from typing import List, Sequence, Union


class BIT1:
    """单点修改"""

    __slots__ = "size", "bit", "tree"

    def __init__(self, n: int):
        self.size = n + 5
        self.bit = n.bit_length()
        self.tree = dict()

    def add(self, index: int, delta: int) -> None:
        index += 1
        while index <= self.size:
            self.tree[index] = self.tree.get(index, 0) + delta
            index += index & -index

    def query(self, right: int) -> int:
        """Query sum of [0, right)."""
        if right > self.size:
            right = self.size
        res = 0
        while right > 0:
            res += self.tree.get(right, 0)
            right -= right & -right
        return res

    def queryRange(self, left: int, right: int) -> int:
        """Query sum of [left, right)."""
        return self.query(right) - self.query(left)


class Solution:
    def getResults(self, queries: List[List[int]]) -> List[bool]:
        n = len(queries)
        bit = BIT1(int(1e6))
        res = []
        for i in range(n):
            t, *x = queries[i]
            if t == 1:
                bit.add(x[0], 1)
            else:
                pos, size = x
                res.append(bit.queryRange(0, x) > 0)
        return res
