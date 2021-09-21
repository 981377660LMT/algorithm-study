from typing import List
import bisect


class SummaryRanges:
    def __init__(self):
        """
        Initialize your data structure here.
        """
        self.intervals = [[float('-inf'), float('-inf')], [float('inf'), float('inf')]]

    # 合并左右/左/右/不合并
    def addNum(self, val: int) -> None:
        i = bisect.bisect_right(self.intervals, [val, val])
        ps, pe = self.intervals[i - 1]
        ns, ne = self.intervals[i]
        if pe == val - 1 and ns == val + 1:
            self.intervals = self.intervals[: i - 1] + [[ps, ne]] + self.intervals[i + 1 :]
        elif pe == val - 1:
            self.intervals[i - 1][1] = val
        elif ns == val + 1:
            self.intervals[i][0] = val
        elif pe < val - 1 and ns > val + 1:
            self.intervals = self.intervals[:i] + [[val, val]] + self.intervals[i:]

    def getIntervals(self) -> List[List[int]]:
        return self.intervals[1:-1]


sr = SummaryRanges()

sr.addNum(1)
print(sr.getIntervals())  # [1, 1]
sr.addNum(3)
print(sr.getIntervals())  # [1, 1], [3, 3]
sr.addNum(7)
print(sr.getIntervals())  # [1, 1], [3, 3], [7, 7]
sr.addNum(2)
print(sr.getIntervals())  # [1, 3], [7, 7]
sr.addNum(6)
print(sr.getIntervals())  # [1, 3], [6, 7]

