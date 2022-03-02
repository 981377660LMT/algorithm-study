# 求出给定子数组内一个给定值的 频率

from typing import List
from collections import defaultdict
from bisect import bisect_left, bisect_right


class RangeFreqQuery:
    """查询子数组内一个值的出现次数"""

    def __init__(self, arr: List[int]):
        self.indexMap = defaultdict(list)
        for index, value in enumerate(arr):
            self.indexMap[value].append(index)

    def query(self, left: int, right: int, value: int) -> int:
        return bisect_right(self.indexMap[value], right) - bisect_left(self.indexMap[value], left)


# 应用场景:对每笔订单，有订单类型和订单时间
# 方便查询某段时间内订单的数量
# 用一个defaultdict(SortedList)，这个数据结构就可以在线查询某段时间内某类订单的数量
