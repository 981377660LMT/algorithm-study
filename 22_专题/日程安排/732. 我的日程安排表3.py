from sortedcontainers import SortedDict


class MyCalendarThree:
    def __init__(self):
        self.diff = SortedDict()

    def book(self, start: int, end: int) -> int:
        """日程安排 [start, end) ，请你在每个日程安排添加后，
        返回一个整数 k ，表示所有先前日程安排会产生的最大 k 次预订"""
        self.diff[start] = self.diff.get(start, 0) + 1
        self.diff[end] = self.diff.get(end, 0) - 1
        res, cur = 0, 0
        for key in self.diff:
            cur += self.diff[key]
            res = max(res, cur)
        return res


# Your MyCalendarThree object will be instantiated and called as such:
# obj = MyCalendarThree()
# param_1 = obj.book(start,end)
