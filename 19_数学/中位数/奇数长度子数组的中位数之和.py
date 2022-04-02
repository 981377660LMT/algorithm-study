# Sum of Odd Length Medians
from heapq import heappush, heappushpop


class MedianFinder:
    def __init__(self):
        self.small = []  # the smaller half of the list, max heap (invert min-heap)
        self.large = []  # the larger half of the list, min heap
        self.size = 0

    def addNum(self, num: int) -> None:
        if len(self.small) == len(self.large):
            heappush(self.large, -heappushpop(self.small, -num))
        elif len(self.small) < len(self.large):
            heappush(self.small, -heappushpop(self.large, num))
        self.size += 1

    def findMedian(self) -> float:
        if len(self.small) == len(self.large):
            return float(self.large[0] - self.small[0]) / 2
        return float(self.large[0])


# n^2logn
class Solution:
    def solve(self, nums):
        res = 0
        for i in range(len(nums)):
            medianFinder = MedianFinder()
            for j in range(i, len(nums)):
                medianFinder.addNum(nums[j])
                if medianFinder.size & 1:
                    res += medianFinder.findMedian()

        return res
