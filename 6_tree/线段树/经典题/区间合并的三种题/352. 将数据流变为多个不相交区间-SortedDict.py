from typing import List
from sortedcontainers import SortedDict

# 0 <= val <= 104
# 最多调用 addNum 和 getIntervals 方法 3 * 104 次

# SortedDict 二分之后是得到键的插入位置，要查找对应的键需要divt.keys[index]
# SortedDict 二分之后是得到键的插入位置，要查找对应的键值对需要divt.peekitem(index)
# SortedDict 二分之后是得到键的插入位置，要删除对应的键值对需要divt.popitem(index)
# SortedDict 二分之后是得到键的插入位置，要删除对应的键需要divt.pop(key) 或者 del divt[key]


class SummaryRanges:
    def __init__(self):
        self.data = SortedDict()  # start=>end

    def addNum(self, val: int) -> None:
        """添加[val,val]"""
        right = self.data.bisect_right(val)
        left = right - 1

        # 先判断val是否已经在某个区间内
        if self.data and left >= 0:
            start, end = self.data.peekitem(left)
            if start <= val <= end:
                return

        leftConn = left >= 0 and self.data.peekitem(left)[1] + 1 == val
        rightConn = right < len(self.data) and self.data.peekitem(right)[0] - 1 == val

        if rightConn and leftConn:
            l, r = self.data.peekitem(left)[0], self.data.peekitem(right)[1]
            self.data[l] = r
            # 删除item，即右边的区间
            self.data.popitem(right)
        elif leftConn:
            l = self.data.peekitem(left)[0]
            self.data[l] = val
        elif rightConn:
            l, r = self.data.peekitem(right)
            self.data.pop(l)
            self.data[l - 1] = r
        else:
            self.data[val] = val

    def getIntervals(self) -> List[List[int]]:
        return list(self.data.items())  # type: ignore


if __name__ == '__main__':
    # intervalManager = SummaryRanges()
    # intervalManager.addNum(1)
    # intervalManager.addNum(3)
    # assert intervalManager.getIntervals() == [[1, 1], [3, 3]]
    # intervalManager.addNum(2)
    # intervalManager.addNum(7)
    # assert intervalManager.getIntervals() == [[1, 3], [7, 7]]
    m = SortedDict({10: 2, 30: 4, 52: 6})
    pos = m.bisect_right(22) - 1
    pos = m.bisect_right(112) - 1
    pos = m.bisect_right(2) - 1
    print(m.keys()[pos])
