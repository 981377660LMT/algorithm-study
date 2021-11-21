# 求出给定子数组内一个给定值的 频率

from typing import List
from collections import defaultdict
from bisect import bisect_left, bisect_right


class RangeFreqQuery:
    def __init__(self, arr: List[int]):
        self.record = defaultdict(list)
        for index, value in enumerate(arr):
            self.record[value].append(index)

    def query(self, left: int, right: int, value: int) -> int:
        return bisect_right(self.record[value], right) - bisect_left(self.record[value], left)


# 应用场景:对每笔订单，有订单类型和订单时间
# 方便查询某段时间内订单的数量
