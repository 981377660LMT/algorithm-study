# 3369. 设计数组统计跟踪器
# https://leetcode.cn/problems/design-an-array-statistics-tracker/description/
# 设计一个数据结构来跟踪它其中的值，并回答一些有关其平均值、中位数和众数的询问。
# 实现 StatisticsTracker 类。
# StatisticsTracker()：用空数组初始化 StatisticsTracker 对象。
# void addNumber(int number)：将 number 添加到数据结构中。
# void removeFirstAddedNumber()：从数据结构删除最早添加的数字。
# int getMean()：返回数据结构中数字向下取整的 平均值。
# int getMedian()：返回数据结构中数字的 中位数。
# int getMode()：返回数据结构中数字的 众数。如果有多个众数，返回最小的那个。


from collections import defaultdict, deque
from sortedcontainers import SortedList


class StatisticsTracker:
    __slots__ = ("_sum", "_queue", "_sl", "_freq", "_freqSl")

    def __init__(self):
        self._sum = 0
        self._queue = deque()
        self._sl = SortedList()
        self._freq = defaultdict(int)
        self._freqSl = SortedList(key=lambda x: (-x[1], x[0]))  # (number, freq)

    def addNumber(self, number: int) -> None:
        self._sum += number
        self._queue.append(number)
        self._sl.add(number)
        f = self._freq[number]
        self._freq[number] = f + 1
        self._freqSl.add((number, f + 1))
        if f:
            self._freqSl.remove((number, f))

    def removeFirstAddedNumber(self) -> None:
        removed = self._queue.popleft()
        self._sum -= removed
        self._sl.remove(removed)
        f = self._freq[removed]
        self._freq[removed] = f - 1
        self._freqSl.remove((removed, f))
        if f > 1:
            self._freqSl.add((removed, f - 1))

    def getMean(self) -> int:
        return self._sum // len(self._queue)

    def getMedian(self) -> int:
        return self._sl[len(self._sl) // 2]

    def getMode(self) -> int:
        return self._freqSl[0][0]
