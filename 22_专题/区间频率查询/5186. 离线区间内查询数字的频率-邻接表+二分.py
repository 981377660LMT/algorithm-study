# 区间内查询数字的频率
# 求出给定子数组内一个给定值的 频率

from typing import List
from collections import defaultdict
from bisect import bisect_left, bisect_right


class RangeFreqQuery:
    """离线查询子数组内一个值的出现次数，如果要在线，使用树状数组"""

    def __init__(self, arr: List[int]):
        self.indexMap = defaultdict(list)
        for index, value in enumerate(arr):
            self.indexMap[value].append(index)

    def query(self, left: int, right: int, value: int) -> int:
        return (
            (bisect_right(self.indexMap[value], right) - 1)
            - bisect_left(self.indexMap[value], left)
            + 1
        )


# 应用场景:对每笔订单，有订单类型和订单时间
# 方便查询某段时间内订单的数量
# 用一个defaultdict(SortedList)，这个数据结构就可以在线查询某段时间内某类订单的数量
