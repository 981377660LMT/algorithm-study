# 352. 将数据流变为多个不相交区间
# 给你一个由非负整数 a1, a2, ..., an 组成的数据流输入，
# 请你将到目前为止看到的数字总结为不相交的区间列表。

# 实现 SummaryRanges 类：

# SummaryRanges() 使用一个空数据流初始化对象。
# void addNum(int val) 向数据流中加入整数 val 。
# int[][] getIntervals() 以不相交区间 [starti, endi] 的列表形式返回对数据流中整数的总结。


from typing import List
from SegmentSet import SegmentSet


class SummaryRanges:
    def __init__(self):
        self.seg = SegmentSet()

    def addNum(self, value: int) -> None:
        self.seg.insert(value, value + 1)

    def getIntervals(self) -> List[List[int]]:
        res = []
        for a, b in self.seg:
            res.append([a, b - 1])
        return res


# Your SummaryRanges object will be instantiated and called as such:
# obj = SummaryRanges()
# obj.addNum(value)
# param_2 = obj.getIntervals()
