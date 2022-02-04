from typing import List, Sequence, Tuple
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
LEFT = 0
RIGHT = 1
# 找出平面中所有矩形叠加覆盖后的总面积。 由于答案可能太大，请返回它对 10 ^ 9 + 7 取模的结果。
# 1.从左往右，线性扫描 注意排序 让x重合的左边界排在右边界前，先加后删
# 2.扫描，看左边界还是有右边界
# 3.计算底边宽*高，高度使用有序容器维护


class Solution:
    def rectangleArea(self, rectangles: List[List[int]]) -> int:
        def calHeight(intervals: Sequence[Tuple[int, int]]) -> int:
            """"求合并后的排序的区间长度之和"""
            if not intervals:
                return 0
            preStart, preEnd = intervals[0]
            res = preEnd - preStart
            for i in range(1, len(intervals)):
                curStart, curEnd = intervals[i]
                if curEnd <= preEnd:
                    continue
                elif curStart <= preEnd and curEnd > preEnd:
                    res += curEnd - preEnd
                    preEnd = curEnd
                else:
                    res += curEnd - curStart
                    preStart = curStart
                    preEnd = curEnd

            return res

        events = []
        for x1, y1, x2, y2 in rectangles:
            events.append((x1, LEFT, y1, y2))
            events.append((x2, RIGHT, y1, y2))
        events.sort()

        res = 0
        intervals = SortedList()
        preX = events[0][0]
        for x, type, y1, y2 in events:
            res += calHeight(intervals) * (x - preX)
            res %= MOD
            if type == LEFT:
                intervals.add((y1, y2))
            else:
                intervals.discard((y1, y2))
            preX = x
        return res


print(Solution().rectangleArea([[0, 0, 2, 2], [1, 0, 2, 3], [1, 0, 3, 1]]))
