# 区间内查询数字的频率
# 求出给定子数组内一个给定值的 频率

from typing import List
from collections import defaultdict
from sortedcontainers import SortedList


class RangeFreqQuery:
    """在线区间频率查询,支持单点修改.基于SortedList实现."""

    def __init__(self, arr: List[int]):
        self._nums = arr[:]
        self._sl = defaultdict(SortedList)
        for i, num in enumerate(arr):
            self._sl[num].add(i)

    def query(self, left: int, right: int, value: int) -> int:
        """[left,right]区间内value的频率."""
        indexes = self._sl[value]
        return indexes.bisect_right(right) - indexes.bisect_left(left)

    def update(self, index: int, value: int) -> None:
        pre = self._nums[index]
        if pre == value:
            return
        self._sl[pre].remove(index)
        self._sl[value].add(index)
        self._nums[index] = value


# 应用场景:对每笔订单，有订单类型和订单时间
# 方便查询某段时间内订单的数量
# 用一个defaultdict(SortedList)，这个数据结构就可以在线查询某段时间内某类订单的数量
