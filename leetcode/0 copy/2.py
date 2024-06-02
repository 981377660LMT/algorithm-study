from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)


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

    def bisectLeft(self, k: int) -> int:
        """返回第一个前缀和大于等于k的位置pos
        0 <= pos <= self.size"""
        curSum, pos = 0, 0
        for i in range(self.bit, -1, -1):
            nextPos = pos + (1 << i)
            if nextPos <= self.size and curSum + self.tree.get(nextPos, 0) < k:
                pos = nextPos
                curSum += self.tree.get(pos, 0)
        return pos

    def bisectRight(self, k: int) -> int:
        """返回第一个前缀和大于k的位置pos
        0 <= pos <= self.size"""
        curSum, pos = 0, 0
        for i in range(self.bit, -1, -1):
            nextPos = pos + (1 << i)
            if nextPos <= self.size and curSum + self.tree.get(nextPos, 0) <= k:
                pos = nextPos
                curSum += self.tree.get(pos, 0)
        return pos

    def __repr__(self) -> str:
        arr = []
        for i in range(self.size):
            arr.append(self.queryRange(i, i + 1))
        return str(arr)

    def __len__(self) -> int:
        return self.size


# TODO: 试一下01树状数组、SegmentSet
class Solution:
    def countDays(self, days: int, meetings: List[List[int]]) -> int:
        max_ = max(e for _, e in meetings) + 10
        bit = BIT1(max_)
        for s, e in meetings:
            bit.add(s, 1)
            bit.add(e + 1, -1)
        res = bit.query(max_)
        return days - res


# days = 10, meetings = [[5,7],[1,3],[9,10]]
