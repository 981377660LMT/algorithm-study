from sortedcontainers import SortedList


class MedianFinder:
    def __init__(self):
        """
        initialize your data structure here.
        """
        self.nums = SortedList()

    def addNum(self, num: int) -> None:
        self.nums.add(num)

    def findMedian(self) -> float:
        return (
            self.nums[n // 2]
            if (n := len(self.nums)) % 2
            else float(self.nums[n // 2] + self.nums[(n - 1) // 2]) / 2
        )

