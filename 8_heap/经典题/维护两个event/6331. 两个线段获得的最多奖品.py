# 在 X轴 上有一些奖品。给你一个整数数组 prizePositions ，
# 它按照 非递减 顺序排列，其中 prizePositions[i] 是第 i 件奖品的位置。
# 数轴上一个位置可能会有多件奖品。再给你一个整数 k 。

# 你可以选择两个端点为整数的线段。每个线段的长度都必须是 k 。
# 你可以获得位置在任一线段上的所有奖品（包括线段的两个端点）。注意，两个线段可能会有相交。


# !分类讨论:
# !1. 两个线段相交
# 此时可以看成一个长为2k的线段,枚举起点计算
# !2. 两个线段不相交,类似 `2054. 两个最好的不重叠活动`
# 用堆维护event的结束时间, 用一个变量维护之前的最大值


from heapq import heappop, heappush
from typing import List, Tuple


class Solution:
    def maximizeWin(self, prizePositions: List[int], k: int) -> int:
        res = 0
        bit = BIT1(maxs(prizePositions))
        for num in prizePositions:
            bit.add(num, 1)
        starts = sorted(set(prizePositions))

        # 两个线段相交
        for start in starts:
            end = start + 2 * k
            res = max(res, bit.queryRange(start, end))

        # 两个线段不相交
        events = [(start, start + k, bit.queryRange(start, start + k)) for start in starts]
        res = max(res, maxTwoEvents(events))
        return res


def maxTwoEvents(events: List[Tuple[int, int, int]]) -> int:
    events = sorted(events, key=lambda x: x[0])
    pq = []  # (end, val)
    res, pre_max = 0, 0
    for start, end, val in events:
        heappush(pq, (end, val))
        while pq and pq[0][0] < start:
            _, pre_val = heappop(pq)
            pre_max = max(pre_max, pre_val)
        res = max(res, pre_max + val)
    return res


def max(x, y):
    if x > y:
        return x
    return y


def maxs(seq):
    res = seq[0]
    for i in range(1, len(seq)):
        if seq[i] > res:
            res = seq[i]
    return res


class BIT1:
    """单点修改"""

    __slots__ = "size", "bit", "tree"

    def __init__(self, n: int):
        self.size = n
        self.bit = n.bit_length()
        self.tree = dict()

    def add(self, index: int, delta: int) -> None:
        # assert index >= 1, 'index must be greater than 0'
        while index <= self.size:
            self.tree[index] = self.tree.get(index, 0) + delta
            index += index & -index

    def query(self, index: int) -> int:
        if index > self.size:
            index = self.size
        res = 0
        while index > 0:
            res += self.tree.get(index, 0)
            index -= index & -index
        return res

    def queryRange(self, left: int, right: int) -> int:
        return self.query(right) - self.query(left - 1)

    def bisectLeft(self, k: int) -> int:
        """返回第一个前缀和大于等于k的位置pos

        1 <= pos <= self.size + 1
        """
        curSum, pos = 0, 0
        for i in range(self.bit, -1, -1):
            nextPos = pos + (1 << i)
            if nextPos <= self.size and curSum + self.tree.get(nextPos, 0) < k:
                pos = nextPos
                curSum += self.tree.get(pos, 0)
        return pos + 1

    def bisectRight(self, k: int) -> int:
        """返回第一个前缀和大于k的位置pos

        1 <= pos <= self.size + 1
        """
        curSum, pos = 0, 0
        for i in range(self.bit, -1, -1):
            nextPos = pos + (1 << i)
            if nextPos <= self.size and curSum + self.tree.get(nextPos, 0) <= k:
                pos = nextPos
                curSum += self.tree.get(pos, 0)
        return pos + 1

    def __repr__(self) -> str:
        preSum = []
        for i in range(self.size):
            preSum.append(self.query(i))
        return str(preSum)

    def __len__(self) -> int:
        return self.size
