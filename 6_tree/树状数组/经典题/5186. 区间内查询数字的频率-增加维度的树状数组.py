# 区间内查询数字的频率
# 求出给定子数组内一个给定值的 频率

from typing import List
from collections import defaultdict


class BIT1:
    """单点修改
    
    https://github.com/981377660LMT/algorithm-study/blob/master/6_tree/%E6%A0%91%E7%8A%B6%E6%95%B0%E7%BB%84/%E7%BB%8F%E5%85%B8%E9%A2%98/BIT.py
    """

    def __init__(self, n: int):
        self.size = n
        self.tree = defaultdict(lambda: defaultdict(int))  # 对每一个值，维护一个树状数组

    @staticmethod
    def _lowbit(index: int) -> int:
        return index & -index

    def add(self, index: int, delta: int, digit: int) -> None:
        if index <= 0:
            raise ValueError('index 必须是正整数')
        while index <= self.size:
            self.tree[index][digit] += delta
            index += self._lowbit(index)

    def query(self, index: int, digit: int) -> int:
        if index > self.size:
            index = self.size
        res = 0
        while index > 0:
            res += self.tree[index][digit]
            index -= self._lowbit(index)
        return res

    def sumRange(self, left: int, right: int, digit: int) -> int:
        return self.query(right, digit) - self.query(left - 1, digit)


class RangeFreqQuery:
    """在线查询子数组内一个值的出现次数:树状数组的前缀和"""

    def __init__(self, arr: List[int]):
        self.bit = BIT1(len(arr) + 10)
        for i, num in enumerate(arr, start=1):
            self.bit.add(i, 1, num)

    def query(self, left: int, right: int, value: int) -> int:
        return self.bit.sumRange(left + 1, right + 1, value)


# 应用场景:对每笔订单，有订单类型和订单时间
# 方便查询某段时间内订单的数量
# 用一个defaultdict(SortedList)，这个数据结构就可以在线查询某段时间内某类订单的数量
