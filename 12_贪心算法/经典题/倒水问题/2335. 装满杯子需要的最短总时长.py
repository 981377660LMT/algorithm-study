"""
现有一台饮水机，可以制备冷水、温水和热水。
每秒钟，可以装满 2 杯 不同 类型的水或者 1 杯任意类型的水
返回装满所有杯子所需的 最少 秒数

1753. 移除石子的最大得分
"""

from heapq import heappop, heappush
from typing import List


class Solution:
    def fillCups(self, amount: List[int]) -> int:
        """贪心
        一种是有一种水特别多，那么答案就是这种水的数量。
        否则，一定可以匹配到只剩一杯，或匹配完。
        """
        max_ = max(amount)
        sum_ = sum(amount)
        restSum = sum_ - max_

        # !最大的
        if max_ > restSum:
            return max_

        # !尽量多的做一次装2杯水的操作,两两抵消
        return (sum_ + 1) // 2

    def fillCups2(self, amount: List[int]) -> int:
        """优先队列模拟 每次取两个最大的消去"""
        pq = []
        for num in amount:  # !不要用heapq.heapify 这个是变小根堆
            if num > 0:  # !注意 [0,0,0]的情况
                heappush(pq, -num)

        res = 0
        while len(pq) >= 2:
            a, b = -heappop(pq), -heappop(pq)
            a, b = a - 1, b - 1
            if a > 0:
                heappush(pq, -a)
            if b > 0:
                heappush(pq, -b)
            res += 1
        if pq:
            res += -pq[0]
        return res
