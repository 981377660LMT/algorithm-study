from heapq import heappop, heappush
from typing import List


# 数据流的第K大元素 维护一个容量为k的小顶堆
class KthLargest:
    def __init__(self, k: int, nums: List[int]):
        self.cap = k
        self.pq = []
        for num in nums:
            self.add(num)

    def add(self, val: int) -> int:
        heappush(self.pq, val)
        if len(self.pq) > self.cap:
            heappop(self.pq)
        return self.pq[0]


from sortedcontainers import SortedList


class KthLargest2:
    def __init__(self, k: int, nums: List[int]):
        self.k = k
        self.sl = SortedList(nums)

    def add(self, val: int) -> int:
        self.sl.add(val)
        return self.sl[-self.k]
