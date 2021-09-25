import bisect
class MyCalendar:

    def __init__(self):
        self.intervals = []

    def book(self, start: int, end: int) -> bool:
      """
        :type start: int
        :type end: int
        :rtype: bool
        """
        if end <= start:
            return False
        
        i = bisect.bisect_right(self.intervals, start)
        if i % 2:            # start is in some stored interval
            return False
        
        j = bisect.bisect_left(self.intervals, end)
        if i != j:
            return False
        
        self.intervals[i:i] = [start, end]
        return True



# Your MyCalendar object will be instantiated and called as such:
# obj = MyCalendar()
# param_1 = obj.book(start,end)