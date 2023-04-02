from typing import List, Tuple, Optional
from collections import defaultdict, Counter, deque
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个整数 n 和一个在范围 [0, n - 1] 以内的整数 p ，它们表示一个长度为 n 且下标从 0 开始的数组 arr ，数组中除了下标为 p 处是 1 以外，其他所有数都是 0 。

# 同时给你一个整数数组 banned ，它包含数组中的一些位置。banned 中第 i 个位置表示 arr[banned[i]] = 0 ，题目保证 banned[i] != p 。

# 你可以对 arr 进行 若干次 操作。一次操作中，你选择大小为 k 的一个 子数组 ，并将它 翻转 。在任何一次翻转操作后，你都需要确保 arr 中唯一的 1 不会到达任何 banned 中的位置。换句话说，arr[banned[i]] 始终 保持 0 。

# 请你返回一个数组 ans ，对于 [0, n - 1] 之间的任意下标 i ，ans[i] 是将 1 放到位置 i 处的 最少 翻转操作次数，如果无法放到位置 i 处，此数为 -1 。

# 子数组 指的是一个数组里一段连续 非空 的元素序列。
# 对于所有的 i ，ans[i] 相互之间独立计算。
# 将一个数组中的元素 翻转 指的是将数组中的值变成 相反顺序 。

from typing import List, Sequence, Union


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


class Solution:
    def minReverseOperations(self, n: int, p: int, banned: List[int], k: int) -> List[int]:
        def getNextPos(cur: int):
            """反转长度为k的子数组后,cur可以到哪些位置"""
            for posInArray in range(k):  # 当前子数组中的位置
                leftBound = cur - posInArray
                rightBound = leftBound + k - 1
                cand = leftBound + rightBound - cur
                if 0 <= leftBound < n and 0 <= rightBound < n:
                    if not isBanned[cand]:
                        yield cand
                else:
                    break

        if k == 1:
            res = [-1] * n
            res[p] = 0
            return res

        bit = BITArray(n)
        for b in banned:
            bit.add(b, 1)

        # 1可以到哪些位置
        # bfs??


# n = 5, p = 0, banned = [2,4], k = 3
print(Solution().minReverseOperations(5, 0, [2, 4], 3))
