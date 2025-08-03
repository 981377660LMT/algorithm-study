import math
from typing import List, Callable, Tuple


class MoRollback:
    """回滚莫队."""

    __slot__ = "_left", "_right"

    def __init__(self):
        self._left = []
        self._right = []

    def addQuery(self, l: int, r: int):
        """添加查询区间 [l, r)."""
        self._left.append(l)
        self._right.append(r)

    def run(
        self,
        addLeft: Callable[[int], None],
        addRight: Callable[[int], None],
        reset: Callable[[], None],
        snapshot: Callable[[], None],
        rollback: Callable[[], None],
        query: Callable[[int], None],
        blockSize=-1,
    ):
        """
        addLeft: 添加左端点的函数，参数为左端点的索引。
        addRight: 添加右端点的函数，参数为右端点的索引。
        reset: 重置区间的函数，无参数。
        snapshot: 快照当前状态的函数，无参数。
        rollback: 回滚到快照状态的函数，无参数。
        query: 执行查询的函数，参数为查询的索引。
        blockSize: 分块大小，默认为 -1，表示使用默认值。
        """
        q = len(self._left)
        if q == 0:
            return
        n = max(self._right)
        if blockSize == -1:
            # n / sqrt(q*2/3)
            blockSize = max(1, n // max(1, int(math.sqrt(q * 2 / 3))))

        nb = (n - 1) // blockSize + 1
        buckets = [[] for _ in range(nb)]

        def naive(qid: int):
            snapshot()
            for i in range(self._left[qid], self._right[qid]):
                addRight(i)
            query(qid)
            rollback()

        # 小区间暴力
        for qid in range(q):
            l, r = self._left[qid], self._right[qid]
            iL, iR = l // blockSize, r // blockSize
            if iL == iR:
                naive(qid)
            else:
                buckets[iL].append(qid)

        # 处理每个大块
        for bucket in buckets:
            if not bucket:
                continue
            bucket.sort(key=lambda qid: self._right[qid])
            lMax = max(self._left[qid] for qid in bucket)
            reset()
            l = lMax
            r = lMax
            for qid in bucket:
                L, R = self._left[qid], self._right[qid]
                while r < R:
                    addRight(r)
                    r += 1
                snapshot()
                while l > L:
                    l -= 1
                    addLeft(l)
                query(qid)
                rollback()
                l = lMax


def discretize(nums: List[int]) -> Tuple[List[int], List[int]]:
    """
    将 nums 中的元素进行离散化，返回新的数组和对应的原始值.

    Args:
        nums (List[int]): 输入的整数数组。

    Returns:
        Tuple[List[int], List[int]]: 返回一个元组，包含两个列表：
            - newNums: 离散化后的数组，其中每个元素是原数组中元素的离散值。
            - origin: 原始值列表，表示每个离散值对应的原始值。

    Example:
        >>> nums = [3, 1, 2, 3, 2]
        >>> newNums, origin = discretize(nums)
        >>> newNums
        [2, 0, 1, 2, 1]
        >>> [origin[x] for x in newNums]
        [3, 1, 2, 3, 2]
    """

    n = len(nums)
    order = sorted(range(n), key=lambda i: nums[i])
    origin = []
    newNums = [0] * n
    for i in order:
        if not origin or origin[-1] != nums[i]:
            origin.append(nums[i])
        newNums[i] = len(origin) - 1
    return newNums, origin


if __name__ == "__main__":

    class Solution:
        # https://leetcode.cn/problems/threshold-majority-queries/description/
        def subarrayMajority(self, nums: List[int], queries: List[List[int]]) -> List[int]:
            mo = MoRollback()
            for l, r, _ in queries:
                mo.addQuery(l, r + 1)

            newNums, origin = discretize(nums)
            size = len(origin)

            counter = [0] * size  # 维护元素出现次数
            counterHistory = []  # 用于回滚counter

            maxCount, maxValue = 0, 0  # 当前最大出现次数和对应的值
            snapVersion = 0  # 快照版本
            snapMaxCount, snapMaxValue = 0, 0  # 快照的最大出现次数和对应的值

            res = [0] * len(queries)

            def add(idx: int):
                nonlocal maxCount, maxValue
                v = newNums[idx]
                counter[v] += 1
                counterHistory.append(v)
                if counter[v] > maxCount or (counter[v] == maxCount and v < maxValue):
                    maxCount, maxValue = counter[v], v

            def reset():
                nonlocal maxCount, maxValue
                for x in counterHistory:
                    counter[x] = 0
                counterHistory.clear()
                maxCount, maxValue = 0, 0

            def snapshot():
                nonlocal snapVersion, snapMaxCount, snapMaxValue
                snapVersion = len(counterHistory)
                snapMaxCount, snapMaxValue = maxCount, maxValue

            def rollback():
                nonlocal maxCount, maxValue
                while len(counterHistory) > snapVersion:
                    x = counterHistory.pop()
                    counter[x] -= 1
                maxCount, maxValue = snapMaxCount, snapMaxValue

            def query(qid: int):
                t = queries[qid][2]
                if maxCount >= t:
                    res[qid] = origin[maxValue]
                else:
                    res[qid] = -1

            mo.run(add, add, reset, snapshot, rollback, query, -1)
            return res
