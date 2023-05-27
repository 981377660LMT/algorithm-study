# 静态区间内查询数字的频率(区间频率查询)
# 求出给定子数组内一个给定值的 频率

from typing import List
from collections import defaultdict
from bisect import bisect_left, bisect_right


class RangeFreqQuery:
    """静态查询子数组内一个值的出现次数."""

    def __init__(self, arr: List[int]):
        self._mp = defaultdict(list)
        for i, v in enumerate(arr):
            self._mp[v].append(i)

    def query(self, left: int, right: int, value: int) -> int:
        """
        查询闭区间[left,right]内value的频率
        0 <= left <= right < len(arr)
        """
        return bisect_right(self._mp[value], right) - bisect_left(self._mp[value], left)


# 应用场景:对每笔订单，有订单类型和订单时间
# 方便查询某段时间内订单的数量
# 用一个defaultdict(SortedList)，这个数据结构就可以在线查询某段时间内某类订单的数量
