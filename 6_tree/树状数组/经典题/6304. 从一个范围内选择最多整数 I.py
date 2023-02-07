# 6304. 从一个范围内选择最多整数 II
# 给你一个整数数组 banned 和两个整数 n 和 maxSum 。你需要按照以下规则选择一些整数：

# 被选择整数的范围是 [1, n] 。
# 每个整数 至多 选择 一次 。
# 被选择整数不能在数组 banned 中。
# 被选择整数的和不超过 maxSum 。
# 请你返回按照上述规则 最多 可以选择的整数数目。
# n<=1e9

from bisect import bisect_right
from itertools import accumulate
from typing import List


def maxCount2(banned: List[int], n: int, maxSum: int) -> int:
    def check(mid: int) -> bool:
        # 选取[1,mid]内的合法整数，使得和不超过maxSum
        pos = bisect_right(bad, mid)
        return (1 + mid) * mid // 2 - preSum[pos] <= maxSum

    bad = sorted(set(banned))
    preSum = [0] + list(accumulate(bad))
    left, right = 1, n
    while left <= right:
        mid = (left + right) // 2
        if check(mid):
            left = mid + 1
        else:
            right = mid - 1

    return right - bisect_right(bad, right)


class Solution:
    def maxCount(self, banned: List[int], n: int, maxSum: int) -> int:
        def check(mid: int) -> bool:
            # 选取[1,mid]内的合法整数，使得和不超过maxSum
            return (1 + mid) * mid // 2 - bit1.query(mid) <= maxSum

        bit1 = BIT1(n)  # 统计区间内被禁止的数的和
        bit2 = BIT1(n)  # 统计区间内被禁止的数的个数
        badSet = set(banned)
        for num in badSet:
            bit1.add(num, num)
            bit2.add(num, 1)
        left, right = 1, n
        while left <= right:
            mid = (left + right) // 2
            if check(mid):
                left = mid + 1
            else:
                right = mid - 1

        return right - bit2.query(right)


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
