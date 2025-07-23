from typing import Generic, List, TypeVar
from collections import defaultdict

from sortedcontainers import SortedList

T = TypeVar("T")


class RangeFreqDynamic(Generic[T]):
    __slots__ = "_data", "_valueToIndexes"

    def __init__(self, data: List[T]) -> None:
        self._data = data[:]
        self._valueToIndexes = defaultdict(SortedList)
        mp = defaultdict(list)
        for i, v in enumerate(data):
            mp[v].append(i)
        for v, indexes in mp.items():
            self._valueToIndexes[v] = SortedList(indexes)

    def query(self, start: int, end: int, value: T) -> int:
        """[start,end)区间内value的频率."""
        if start >= end:
            return 0
        pos = self._valueToIndexes[value]
        return pos.bisect_left(end) - pos.bisect_left(start)

    def set(self, index: int, value: T) -> None:
        pre = self._data[index]
        if pre == value:
            return
        self._valueToIndexes[pre].remove(index)
        self._data[index] = value
        self._valueToIndexes[value].add(index)


if __name__ == "__main__":
    # 此外，给定一个二维整数数组 queries，其中每个 queries[i] 可以是以下两种类型之一：

    # [1, l, r, k] - 计算在区间 [l, r] 中，满足 nums[j] 的 popcount-depth 等于 k 的索引 j 的数量。
    # [2, idx, val] - 将 nums[idx] 更新为 val。
    # 返回一个整数数组 answer，其中 answer[i] 表示第 i 个类型为 [1, l, r, k] 的查询的结果。
    class Solution:
        # 3624. 位计数深度为 K 的整数数目 II
        # https://leetcode.cn/problems/number-of-integers-with-popcount-depth-equal-to-k-ii/
        def popcountDepth(self, nums: List[int], queries: List[List[int]]) -> List[int]:
            def depth(x: int) -> int:
                res = 0
                while x > 1:
                    res += 1
                    x = x.bit_count()
                return res

            R = RangeFreqDynamic([depth(v) for v in nums])
            res = []
            for q in queries:
                if q[0] == 1:
                    res.append(R.query(q[1], q[2] + 1, q[3]))
                else:
                    R.set(q[1], depth(q[2]))
            return res
